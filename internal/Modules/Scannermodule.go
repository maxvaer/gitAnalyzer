// Package Modules contains all business logic modules/components of the application.
package Modules

import (
	"GitAnalyzer/api/Analyzer"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ricochet2200/go-disk-usage/du"
	"golang.org/x/net/websocket"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

// webServer starts a local webserver for the webui frontend
// The provided channel monitorStats can be used, to send stats to the frontend via a websocket.
func webServer(monitorStats <-chan Analyzer.MonitorStat) {
	e := echo.New()

	// Attach middlewares for logging and recovering of failed requests.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define listening GET endpoint for the websocket of the frontend
	e.GET("/ws", func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				select {
				// Read stats from the channel
				case stat := <-monitorStats:
					fmt.Println("Sending stats:", stat)
					// Serialise the stats as json
					statJson, err := json.Marshal(stat)
					if err != nil {
						fmt.Println("Error marshalling json:", err)
						ws.Close()
					}
					// Send the json data via the websocket
					err = websocket.Message.Send(ws, string(statJson))
					if err != nil {
						c.Logger().Error(err)
						fmt.Println("Error closing...")
						ws.Close()
					}
				}
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))
}

// cloneRepoAndRunTemplates reads tasks from the given tasks channel until the channel is empty.
// the provided repoHandler is used to clone the repository from the read task.
// the templateHandler is used to run all templates on the cloned repository
// Updates of the state of the running task are shared via the cTasks channel.
func cloneRepoAndRunTemplates(repoHandler *RepoHandler, templateHandler *TemplateHandler,
	tasks <-chan Analyzer.Task, cTasks chan<- Analyzer.Task) {
	for {
		select {
		// Read task from tasks channel
		case task := <-tasks:
			// Wait a bit to debounce
			time.Sleep(time.Millisecond * 200)
			// Update state to cloning
			task.State = "cloning"
			cTasks <- task
			// Clone the repository from the task
			repo := repoHandler.CloneRepositories(task)
			if repo == nil {
				// If the repository is nil the clone process failed
				// Update state to failed
				task.State = "failed"
				cTasks <- task
				// Continue with next task
				continue
			}
			// Run templates
			templateHandler.RunAllTemplates(task, repo, cTasks)
			// Reset task and repo to prevent memory leak
			task = Analyzer.Task{}
			repo = nil
		}
		// Wait a bit for debounce
		time.Sleep(time.Millisecond * 200)
	}
}

// Run is the main function to start the process of cloning and scanning a repository.
// Therefore, the provided config object is used to set different settings.
func Run(config Analyzer.Config) {

	// Initialize handlers
	repoHandler := NewRepoHandler(&GitHelper{}, config)
	templateHandler := NewTemplateHandler(repoHandler, config)

	// Load templates
	templateHandler.LoadTemplates(config.TemplatesPath)

	// Check if all requirements are met
	templateHandler.CheckRequirements()

	// Run preScripts
	templateHandler.preScript()

	// Load tasks (urls.csv)
	fmt.Println("Loading tasks...")
	tasksSlc := templateHandler.FileHelper.GetTasks(config.UrlFilePath, "./results/checked.csv")
	numberOfTasks := int32(len(tasksSlc))
	fmt.Println("Tasks loaded:", strconv.Itoa(int(numberOfTasks)))

	// Initialize variables for monitoring of stats
	var failedScans, finishedScans, resultsFound, queuedScans, runningScans, cloning int32
	failedScans = 0
	queuedScans = numberOfTasks
	runningScans = 0
	finishedScans = 0
	resultsFound = 0
	cloning = 0

	// Create result directory
	templateHandler.FileHelper.PrepareResultsFolder(config.ResultsDir)

	// Start disk usage monitoring
	go checkDiskUsage()

	// Initialize necessary channels
	cloneQueue := make(chan Analyzer.Task, numberOfTasks)
	cTasks := make(chan Analyzer.Task, numberOfTasks)
	monitor := make(chan Analyzer.MonitorStat, numberOfTasks)

	// Start the webserver
	go webServer(monitor)

	// Start the goroutine to send the stats frequently (every second) to monitoring channel.
	go func(numberOfTasks int32, failedScans, queuedScans, runningScans, finishedScans,
		resultsFound *int32, monitor chan<- Analyzer.MonitorStat) {
		for {
			monitor <- Analyzer.MonitorStat{NumberOfTasks: numberOfTasks, QueuedScans: atomic.LoadInt32(queuedScans),
				RunningScans: atomic.LoadInt32(runningScans), FailedScans: atomic.LoadInt32(failedScans),
				FinishedScans: atomic.LoadInt32(finishedScans), ResultsFound: atomic.LoadInt32(resultsFound)}
			time.Sleep(1 * time.Second)
		}

	}(numberOfTasks, &failedScans, &queuedScans, &runningScans, &finishedScans, &resultsFound, monitor)

	// Add the loaded task to the cloneQueue channel
	for _, task := range tasksSlc {
		// Set state to queued
		task.State = "queued"
		// Add to channel
		cloneQueue <- task
	}
	// Reset the taskSlc to prevent memory leaking
	tasksSlc = nil

	// Start the worker pool
	for i := 0; i < config.WorkerCount; i++ {
		go cloneRepoAndRunTemplates(repoHandler, templateHandler, cloneQueue, cTasks)
	}

	// Read all updates from the cTasks channel
	for {
		select {
		// Read task
		case task := <-cTasks:
			// Wait a bit for debounce
			time.Sleep(time.Millisecond * 200)
			// Read the current state of the task update
			switch task.State {
			case "finished":
				//write result
				if task.Results != nil {
					templateHandler.FileHelper.MarshalMultipleResults(task.Results)
					atomic.AddInt32(&resultsFound, 1)
				}
				//write to checked file
				templateHandler.FileHelper.MarshalStat(Analyzer.Stat{URL: task.URL, ElapsedTime: task.ElapsedTime})

				//Increment finishedScans
				atomic.AddInt32(&finishedScans, 1)
				//Decrement runningScans
				atomic.AddInt32(&runningScans, -int32(1))
			case "failed":
				//Increment failedScans
				atomic.AddInt32(&failedScans, 1)
				//Decrement cloning
				atomic.AddInt32(&cloning, -int32(1))
			case "running":
				//Increment runningScans
				atomic.AddInt32(&runningScans, 1)
				//Decrement cloning
				atomic.AddInt32(&cloning, -int32(1))
			case "cloning":
				//Increment cloning
				atomic.AddInt32(&cloning, 1)
				//Decrement queuedScans
				atomic.AddInt32(&queuedScans, -int32(1))
			}

		}
		// Check if all scans are finished
		if allScansFinished(numberOfTasks, &finishedScans, &failedScans) {
			// Shut down worker pool
			// Break for loop
			break
		}
		// Wait a bit for debounce
		time.Sleep(time.Millisecond * 200)
	}

	// Run post process script at the end.
	templateHandler.postProcess()

	fmt.Println("Finished")
}

// allScansFinished checks if any scans are still running or if all are finished.
// numberOfTasks is the total number of all repos to scan.
// finishedScans is the amount of scans which finished successfully
// failedScans is the amount of repos which could not be cloned.
func allScansFinished(numberOfTasks int32, finishedScans, failedScans *int32) bool {
	return numberOfTasks == (atomic.LoadInt32(finishedScans) + atomic.LoadInt32(failedScans))
}

// checkDiskUsage checks every 5 seconds if the current
// disk usage is above 90%.
// If this is the case, the application will exit.
func checkDiskUsage() {
	// Check for exiting conditions
	for {
		// Critical disk usage
		diskUsage := du.NewDiskUsage("/")
		if diskUsage.Usage() > 90 {
			fmt.Print("Aborting Disk usage critical!")
			os.Exit(1)
		}
		time.Sleep(5 * time.Second)
	}
}
