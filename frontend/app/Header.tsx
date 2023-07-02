import Link from 'next/link'
import React from 'react'

export default function Header() {
  return (
    <div className='p-5 flex space-x-4 bg-blue-500 font-bold text-white'>
        <Link href='/monitor'>Monitor</Link>
        <Link href='/templates'>Templates</Link>
        <Link href='/results'>Results</Link>
        <Link href='/repos'>Repos</Link>
        <Link target="_blank" href='https://github.com/maxvaer/gitAnalyzer/wiki'>Wiki</Link>
    </div>
  )
}
