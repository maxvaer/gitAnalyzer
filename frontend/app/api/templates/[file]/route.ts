import path from 'path';
import { NextResponse } from 'next/server'
import { NextApiRequest, NextApiResponse } from 'next'
import fs from 'fs'

export async function GET(request: NextApiRequest, context:any) {
  let data: string = ''
  const { file } = context.params
  try {
    const folder = path.join(process.cwd(), '../', 'templates');
    const filePath = path.join(folder, file);
    data = fs.readFileSync(filePath, 'utf8');
  } catch (error) {
    return NextResponse.json({ message: `File not found: ${file}` }, { status: 404 })
  }
  

  return NextResponse.json(data)
}