"use client";

//Wrapper Component to fix the following bug:
//https://github.com/vercel/next.js/discussions/47547

import React from 'react'
import type { ReactNode } from "react";
import Header from './Header'
import styles from './layout.module.css'
import { Grommet } from 'grommet'

export const AppWrapper = ({ children }: { children: ReactNode }) => {
  return (
    <div>
        <Grommet full>
          <Header />
          <main className={styles.main}>
            {children}
          </main>
        </Grommet>
    </div>
  )
}
