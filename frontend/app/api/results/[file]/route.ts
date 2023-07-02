import path from 'path';
import {csv} from "csvtojson";
import { NextResponse } from 'next/server'
import { NextApiRequest, NextApiResponse } from 'next'

export async function GET(request: NextApiRequest, context:any) {
  let data: any = []
  const { file } = context.params
  try {
    const folder = path.join(process.cwd(), '../', 'results');
    const filePath = path.join(folder, file);
    await csv()
    .fromFile(filePath)
    .then((rows:JSON)=>{
        data = rows
    })
  } catch (error) {
    return NextResponse.json({ message: `File not found: ${file}.csv` }, { status: 404 })
  }
  

  data = data.map((item:any, index:number) => Object.assign({id: index+1 }, item))
  return NextResponse.json(data)
}