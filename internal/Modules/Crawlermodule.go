// Package Modules contains all business logic modules/components of the application.
package Modules

import (
	"GitAnalyzer/api/Analyzer"
	"context"
	"database/sql"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)
import "github.com/google/go-github/v50/github"

// loadConfig loads the configuration from the local "config.env" file
func loadConfig() Analyzer.CrawlConfig {
	// Read the config.env file
	viper.SetConfigFile("config.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("error reading config.env:", err)
	}
	// Get the values
	gitHubToken := viper.GetString("GITHUB_API_KEY")
	if gitHubToken == "" {
		log.Fatalln("GITHUB_API_KEY not set in config.env!")
	}
	dbUser := viper.GetString("DB_USER")
	if dbUser == "" {
		log.Fatalln("DB_USER not set in config.env!")
	}
	dbPassword := viper.GetString("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatalln("DB_PASSWORD not set in config.env!")
	}
	dbIP := viper.GetString("DB_IP")
	if dbIP == "" {
		log.Fatalln("DB_IP not set in config.env!")
	}
	dbPort := viper.GetString("DB_PORT")
	if dbPort == "" {
		log.Fatalln("DB_PORT not set in config.env!")
	}

	// Return new CrawlConfig with loaded settings
	return Analyzer.CrawlConfig{GitHubAPIToken: gitHubToken, DBUser: dbUser, DBPassword: dbPassword, DBIP: dbIP, DBPort: dbPort}
}

// Crawl starts the process of crawling all public repositories of GitHub
func Crawl() {
	// Load config from config.env
	crawlConfig := loadConfig()

	// Open DB connection
	db := openDB(crawlConfig.DBUser, crawlConfig.DBPassword, crawlConfig.DBIP, crawlConfig.DBPort)
	defer db.Close()
	if db == nil {
		return
	}

	// Create new GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: crawlConfig.GitHubAPIToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opt := &github.RepositoryListAllOptions{
		Since: int64(loadMaxId(db)),
	}

	// Check the current rate limit of the used API token
	checkRateLimit(client, ctx)
	// Initialize the repository list
	repos, _, err := client.Repositories.ListAll(ctx, opt)
	if err != nil {
		log.Fatalln("Error getting list of public repositories:", err.Error())
		return
	}

	// Loop through pages of results
	for {
		// Move to unread repositories
		opt.Since = *repos[len(repos)-1].ID

		// Check the current rate limit
		checkRateLimit(client, ctx)

		// Fetch new repositories
		reposPage, resp, err := client.Repositories.ListAll(ctx, opt)
		if err != nil {
			log.Fatalln("Error getting list of public repositories:", err.Error())
			return
		}

		// Iterate over fetched repositories
		for _, repository := range reposPage {
			// Crawl the metadata of every repository
			getRepositoryDetails(db, client, ctx, repository)
		}
		// Append the crawled repositories to the list
		repos = append(repos, reposPage...)

		// Check if no new page is available
		if resp.NextPage == 0 {
			// Crawled all repositories, therefore break
			break
		}
	}

	fmt.Println("No more public repositories to fetch, exiting!")
	os.Exit(0)
}

// checkRateLimit uses the provided client and context to check if the rate limit is reached.
// If the rate limit is reached, the crawler will sleep for 10 minutes.
func checkRateLimit(client *github.Client, ctx context.Context) {

	for {
		// Check if some API request are available
		if fetchRemainingAPIRequests(client, ctx) > 10 {
			// Enough APi request are available, therefore break
			break
		}

		// If no API request are available, sleep for 10 minutes
		fmt.Println("API Limit reached, sleeping for 10 minutes...")
		time.Sleep(10 * time.Minute)
	}

}

// fetchRemainingAPIRequests uses the provided client and context to check
// how many API request are available
func fetchRemainingAPIRequests(client *github.Client, ctx context.Context) int {
	rate, _, err := client.RateLimits(ctx)
	if err != nil {
		fmt.Println("Error checking remaining rate limit:", err.Error(), rate)
		return 0
	}
	return rate.Core.Remaining
}

// getRepositoryDetails uses the provided client and context to crawl the details the repository.
// The results are saved inside the DB.
func getRepositoryDetails(db *sql.DB, client *github.Client, ctx context.Context, repository *github.Repository) {
	// Check the rate limit
	checkRateLimit(client, ctx)

	// Get the whole repository object
	repo, _, err := client.Repositories.Get(ctx, *repository.Owner.Login, *repository.Name)
	if err != nil {
		//fmt.Println("Error fetching details for repository:", err.Error())
		return
	}

	// Check the rate limit
	checkRateLimit(client, ctx)
	// Get the list of commits for the repository
	commits, _, err := client.Repositories.ListCommits(ctx, *repository.Owner.Login, *repository.Name, &github.CommitsListOptions{})
	if err != nil {
		//fmt.Println("Error fetching commits for repository:", err.Error())
		return
	}

	// Check if the repository ID is not already inside the DB
	if !idIsInData(db, strconv.Itoa(int(repo.GetID()))) {
		// Check if the repository is not a fork, has a size of more than 92 Bytes and at least 2 commits
		if !repo.GetFork() && repo.GetSize() > 92 && len(commits) >= 2 {
			// Add the data to the database
			insertDataFromRepo(db, repo, len(commits))
		}
	}
}

// insertDataFromRepo inserts data from the provided repository and commitCount into the database
func insertDataFromRepo(db *sql.DB, repo *github.Repository, commitCount int) {
	// Extract data from repo object
	url := repo.GetGitURL()
	id := strconv.FormatInt(repo.GetID(), 10)
	creationDate := repo.GetCreatedAt().Format("2006-01-02")
	forkCount := strconv.Itoa(repo.GetForksCount())
	size := strconv.Itoa(repo.GetSize())
	stars := strconv.Itoa(repo.GetStargazersCount())
	updateDate := repo.GetUpdatedAt().Format("2006-01-02")
	language := repo.GetLanguage()
	cCount := strconv.Itoa(commitCount)

	// Insert data into the table
	writeRowIntoDB(db, url, id, creationDate, forkCount, size, stars, updateDate, language, cCount)
	fmt.Println("Crawled Repo:", url)
}

// unmarshallRecords load all CrawlRecords from a provided file
// and returns them as a slice of CrawlRecords
func unmarshallRecords(csvPath string) []Analyzer.CrawlRecord {
	// Initialize empty slice of Analyzer.CrawlRecords
	var records []Analyzer.CrawlRecord
	// Open file
	recordsCSVFile, err := os.Open(csvPath)
	// Load all Analyzer.CrawlRecords from the file
	if err = gocsv.UnmarshalFile(recordsCSVFile, &records); err != nil {
		log.Fatalln("Error Unmarshalling records for crawler:", err.Error())
		return nil
	}

	defer recordsCSVFile.Close()

	// Return the loaded records
	return records
}

// InsertDataFromCSV loads all Analyzer.CrawlRecords from the provided file path and
// inserts them into the database
func InsertDataFromCSV(path string) {
	// Load config
	crawlConfig := loadConfig()
	// Read records from file
	records := unmarshallRecords(path)
	// Open database connection
	db := openDB(crawlConfig.DBUser, crawlConfig.DBPassword, crawlConfig.DBIP, crawlConfig.DBPort)
	defer db.Close()
	if db == nil {
		return
	}

	// Iterate over loaded records
	for _, record := range records {
		// Check if the ID of the loaded records is not already in the database
		if !idIsInData(db, record.ID) {
			// Write the records into the database
			writeRowIntoDB(db, record.URL, record.ID, record.CreationDate,
				record.ForkCount, record.Size, record.Stars, record.UpdateDate, record.Language, record.CommitCount)
		}
	}

	fmt.Println("Import finished!")
	os.Exit(0)
}

// writeRowIntoDB inserts a new row with the provided data into the database
func writeRowIntoDB(db *sql.DB, url, id, creationDate, forkCount, size, stars, updateDate, language, commitCount string) {
	// Prepare the insert statement
	stmt, err := db.Prepare("INSERT INTO repos (url, id, creation_date, fork_count, size, stars, update_date, language, commit_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalln("Error preparing statement to insert crawled data:", err.Error())
		return
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(url, id, creationDate, forkCount, size, stars, updateDate, language, commitCount)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// idIsInData checks if the provided ID is already inside the provided database
func idIsInData(db *sql.DB, id string) bool {
	// Prepare the database statement
	stmt, err := db.Prepare("Select url from repos Where id = ?")
	if err != nil {
		log.Fatalln("Error preparing statement to insert crawled data:", err.Error())
		return false
	}
	defer stmt.Close()

	// Execute the statement
	var url string
	err = stmt.QueryRow(id).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Fatalln("Error fetching id from crawled data:", err.Error())
		return false
	}

	// Check if the result is not an empty string
	return url != ""
}

// Loads the maximum(last) ID from the given database.
func loadMaxId(db *sql.DB) int {
	// Prepare the select statement
	stmt, err := db.Prepare("SELECT MAX(id) from repos")
	if err != nil {
		log.Fatalln("Error preparing statement to fetch max id:", err.Error())
		return 0
	}
	defer stmt.Close()

	//Execute the statement
	var id int
	err = stmt.QueryRow().Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		log.Fatalln("Error fetching id from crawled data:", err.Error())
		return 0
	}

	return id
}

// openDB uses the provided database credentials to open a new connection.
// If the repos table is not found, it will be created
func openDB(username, password, ip, port string) *sql.DB {
	// Connect to mariadb
	db, err := sql.Open("mysql", ""+username+":"+password+"@tcp("+ip+":"+port+")/gitanalyzer")
	if err != nil {
		log.Fatalln("Error opening database connection:", err.Error())
		return nil
	}

	// Create the table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS repos (
    url VARCHAR(255) NOT NULL,
    id INT NOT NULL PRIMARY KEY,
    creation_date DATE NOT NULL,
    fork_count INT NOT NULL,
    size INT NOT NULL,
    stars INT NOT NULL,
    update_date DATE NOT NULL,
    language VARCHAR(255) NOT NULL,
    commit_count INT NOT NULL);`)
	if err != nil {
		log.Fatalln("Error creating table in database:", err.Error())
		return nil
	}

	// Return the database connection
	return db
}
