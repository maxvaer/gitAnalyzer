'use client'

import React from 'react'
import useSWR from 'swr';
import { DataTable, Spinner, Text } from 'grommet';

interface Filename {
    fileName: string;
}

export default function Table({fileName}:Filename) {
  const fetcher = (url:string) => fetch(url).then((res) => {
    return res.json()
  });

  const getColumns = (entry:any) => {
    let columns = []
    for(let key in entry){
      columns.push({
        property: key,
        header: <Text className='capitalize'>{key.replaceAll("_", " ")}</Text>,
      })
    }
    return columns
  }
  const { data, error } = useSWR('/api/results/'+fileName, fetcher);
  
  if (error) return <div>Failed to load</div>;
  if (!data) return <div><Spinner /></div>;
  if (data.message) return <div>{data.message}</div>
  
  return (
    <div>
      <DataTable
        columns={getColumns(data[0])}
        data={data}
        paginate={true}
        sortable={true}
      />
    </div>
  );
}
