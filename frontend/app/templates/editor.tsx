'use client'

import React, { useEffect } from 'react'
import AceEditor from "react-ace";
import useSWR from 'swr';
import { Button, Notification, TextInput, Form, FormField, Box } from 'grommet';

import "ace-builds/src-noconflict/mode-yaml";
import "ace-builds/src-noconflict/theme-github";
import "ace-builds/src-noconflict/ext-language_tools";
import "ace-builds/src-noconflict/snippets/python";

interface Props {
    file: string
}

export default function Templates({file}:Props) {
  const [templateData, setTemplateData] = React.useState('');
  const [newFile, setNewFile] = React.useState('');
  const [errors, setErrors] = React.useState([]);
  const [showSaveButton, setShowSaveButton] = React.useState(false);
  
  useEffect(() => {
    fetch('/api/templates/'+file, {
      method: 'GET',
    })
    .then(response => response.json())
    .then(data => {
      setTemplateData(data)
      setNewFile(file)
    });
  },[file])

  const onChange = (newValue:string) =>  {
    setTemplateData(newValue)
    setShowSaveButton(false)
  };

  const onChangeFilename = (event:any) =>  {
    let value = event.target.value
    setNewFile(value)
  };

  const onClick = () =>  {
    fetch('/api/templates/validate', {
      method: 'POST',
      body: JSON.stringify(templateData),
      headers: {
        'Content-type': 'application/json; charset=UTF-8',
      },
    })
       .then((response) => response.json())
       .then((data) => {
          console.log("Post response:", data);
          setErrors(data)
          setShowSaveButton(data.length === 0)
       })
       .catch((err) => {
          console.log("Post error:", err.message);
       });
  };

  const onClickSave = () =>  {
    fetch('/api/templates/save', {
      method: 'POST',
      body: JSON.stringify({newFile, templateData}),
      headers: {
        'Content-type': 'application/json; charset=UTF-8',
      },
    })
       .then((response) => response.json())
       .then((data) => {
          console.log("Save resp:", data);
       })
       .catch((err) => {
          console.log("Post error:", err.message);
       });
  };


  return (
    <div>
      <AceEditor className='mt-5 mb-5'
        mode="yaml"
        theme="github"
        value={templateData}
        onChange={onChange}
        name="AceEditor"
        style={{ width: '900px' }}
        setOptions={{ useWorker: false }}
        editorProps={{ $blockScrolling: true }}
      />
      {!showSaveButton &&<Button primary label="Validate" onClick={onClick}/>}
      {showSaveButton && 
      <div className='flex'>
        <Form 
          onSubmit={onClickSave} 
          validate="change"
        >
        <Box direction="row" gap="medium">
          <Button type="submit" primary label="Save" className='mr-5'/> 
          <FormField 
            name="fileName" 
            htmlFor="filename-id" 
            label="Filename:" 
            required
            validate={(fileName:string) => {
                if (fileName && fileName.length <= 6)
                  return 'must be >6 character!';
                if (fileName && fileName == 'template.yaml')
                  return 'Filename template.yaml is not allowed!';
                if (fileName && !fileName.endsWith(".yaml"))
                  return 'Filename must end with .yaml';
                return;
              }
            }
          >
            <TextInput
              name='fileName'
              id='filename-id'
              className='w-fit grow-0'
              placeholder="Enter filename"
              value={newFile}
              onChange={onChangeFilename}
            />
          </FormField>
        </Box>
        </Form>
  
      </div>}
      <div className='mt-5'>
        {Array.isArray(errors) && errors.map((valError:any, index:number) => (
          <Notification
            key={index}
            title={valError.property}
            message={valError.message}
            onClose={() => {}}
          />
        ))}
      </div>
    </div>
  )
}
