import path from 'path';
import { NextResponse } from 'next/server'
import { NextApiRequest, NextApiResponse } from 'next'
import fs from 'fs'

export async function GET(request: NextApiRequest, res:NextApiResponse) {
  let files:any = []
  try {
    const folder = path.join(process.cwd(), '../', 'results');
    fs.readdirSync(folder).forEach(file => {
      const filePath = path.join(folder, file);
      let size:number = fs.statSync(filePath).size
      if (size > 0) {
        files.push({
          file: file,
          size 
        })
      }
    });
  } catch (error) {
    return NextResponse.json({ message: 'Result dir not found' }, { status: 404 })
  }
  
  return NextResponse.json(files)
}