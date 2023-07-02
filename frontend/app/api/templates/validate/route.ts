import path from 'path';
import { NextResponse } from 'next/server'
import { NextApiRequest, NextApiResponse } from 'next'
import YAML, { YAMLException } from 'js-yaml'
import schema from './schema';

export async function POST(request: Request, res:NextApiResponse) {
    let result = []
    try {
        const templateData = await request.json()
        const parsedYaml = YAML.load(templateData);
        const validate = require('jsonschema').validate;
        result = validate(parsedYaml, schema);
    } catch (error:YAMLException) {
        if (error instanceof YAMLException) {
            return NextResponse.json([{ property: 'Indentation error', message: `Error near line: ${error.mark.line}` }])
        }
        return NextResponse.json({ message: 'Validation error' }, { status: 500 })
    }
    
    

  return NextResponse.json(result.errors)
}