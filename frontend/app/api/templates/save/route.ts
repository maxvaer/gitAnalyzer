import path from 'path';
import { NextResponse } from 'next/server'
import { NextApiRequest, NextApiResponse } from 'next'
import fs from 'fs'

export async function POST(request: Request, res:NextApiResponse) {
    let fileName = ''
    try {
        const data = await request.json()
        let fileName = data.newFile
        let templateData = data.templateData
        const folder = path.join(process.cwd(), '../', 'templates');
        const filePath = path.join(folder, fileName)
        fs.writeFileSync(filePath, templateData, { flag: 'w+' });
    } catch (error) {
        return NextResponse.json({ message: 'Save template error', error}, { status: 500 })
    }

  return NextResponse.json({fileName}, { status: 200 })
}