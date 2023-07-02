'use client'

import React from 'react'
import useSWR from 'swr';
import { Select, Spinner } from 'grommet';
import Table from './table';

export default function Results() {
  interface Option {
    file: string;
    size: number;
  }
  const [option, setOption] = React.useState<Option>({
    file: '',
    size: 0
  });

  const fetcher = (url:string) => fetch(url).then((res) => {
    return res.json()
  });

  const { data, error } = useSWR('/api/results/', fetcher);
  
  if (error) return <div>Failed to load</div>;
  if (!data) return <div><Spinner /></div>;
  if (data.message) return <div>{data.message}</div>
  return (
    <div>
      <div className='flex justify-between items-center'>
        <Select
          options={data}
          value={option.file}
          placeholder="Please select a file."
          onChange={({ option }) => setOption(option)}
        />
        {option.file != '' ? <div className='font-semibold'>Size: {option.size}B</div> : null}
      </div>
      {option.file != '' ? <Table fileName={option.file} /> : null}
    </div>
  );
}
