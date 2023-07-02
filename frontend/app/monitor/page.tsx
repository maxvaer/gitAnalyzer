'use client';

import { useState } from 'react'
import { WebSocket } from 'nextjs-websocket'

function Monitor() {
  const [failedScans, setFailedScans] = useState(0)
  const [queuedScans, setQueuedScans] = useState(0)
  const [runningScans, setRunningScans] = useState(0)
  const [finishedScans, setFinishedScans] = useState(0)
  const [resultsFound, setResultsFound] = useState(0)
  const [numberOfTasks, setNumberOfTasks] = useState(0)

  const handleMessage = (jsonString: string) => {
    let result = JSON.parse(jsonString)
    console.log("Got", result)
    if(result.hasOwnProperty('failedScans')){
      setFailedScans(result['failedScans'])
    }
    if(result.hasOwnProperty('queuedScans')){
      setQueuedScans(result['queuedScans'])
    }
    if(result.hasOwnProperty('runningScans')){
      setRunningScans(result['runningScans'])
    }
    if(result.hasOwnProperty('finishedScans')){
      setFinishedScans(result['finishedScans'])
    }
    if(result.hasOwnProperty('resultsFound')){
      setResultsFound(result['resultsFound'])
    }
    if(result.hasOwnProperty('numberOfTasks')){
      setNumberOfTasks(result['numberOfTasks'])
    }
  }

  return (
    <div>
      <WebSocket
        reconnect={true}
        url={'ws://localhost:8080/ws'}
        onOpen={() => console.log('Connection esatblished')}
        onClose={(e: any) => console.log('Closed:', e)}
        onMessage={(data: string) => handleMessage(data)}
      />
      <table className="border-separate border-spacing-x-5 border-spacing-y-2 ">
        <thead className='uppercase '>
          <tr>
            <th className='text-start '>Statistic:</th>
            <th>Count:</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Repositories</td>
            <td>{numberOfTasks}</td>
          </tr>
          <tr>
            <td>Queued Scans</td>
            <td>{queuedScans}</td>
          </tr>
          <tr>
            <td>Failed Scans</td>
            <td>{failedScans}</td>
          </tr>
          <tr>
            <td>Running Scans</td>
            <td>{runningScans}</td>
          </tr>
          <tr>
            <td>Finished Scans</td>
            <td>{finishedScans}</td>
          </tr>
          <tr>
            <td>Results found</td>
            <td>{resultsFound}</td>
          </tr>
        </tbody>
      </table>
    </div>
  )
}


export default Monitor