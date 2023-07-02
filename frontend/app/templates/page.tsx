'use client'

import React from 'react'
import dynamic from 'next/dynamic'
import useSWR from 'swr';
import { Select, Spinner, Button } from 'grommet';

export default function Templates() {
  interface TemplateOption {
    file: string;
    size: number;
  }
  const [templateOption, setTemplateOption] = React.useState<TemplateOption>({
    file: '',
    size: 0
  });

  const fetcher = (url:string) => fetch(url).then((res) => {
    return res.json()
  });

  const { data, error } = useSWR('/api/templates/', fetcher);
  
  if (error) return <div>Failed to load</div>;
  if (!data) return <div><Spinner /></div>;
  if (data.message) return <div>{data.message}</div>

  const Ace = dynamic(
    () => import("./editor"),
    { ssr: false }
  )

  const onClickNewTemplate = () =>  {
    setTemplateOption({file: "template.yaml", size: 0})
  };


  return (
    <div>
      <div className='flex flex-col'>
        <div>
          <Select
            options={data}
            value={templateOption.file}
            placeholder="Edit existing template."
            onChange={({ option }) => setTemplateOption(option)}
          />
        </div>
        {templateOption.file == '' && <div className='my-5 text-center'>or</div>}
        {templateOption.file == '' && <Button primary label="Create new template" onClick={onClickNewTemplate}/>} 
      </div>
      {templateOption.file != '' ? <Ace file={templateOption.file} /> : null}
    </div>
  )
}
