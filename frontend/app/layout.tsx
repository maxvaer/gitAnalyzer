import './globals.css'
import { AppWrapper } from './AppWrapper'

export const metadata = {
  title: 'GitAnalyzer',
  description: 'WebUI of GitAnalyzer',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>
        <AppWrapper>{children}</AppWrapper>
      </body>
    </html>
  )
}
