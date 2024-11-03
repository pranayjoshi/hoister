import './globals.css'
import { Inter } from 'next/font/google'
import { GitURLProvider } from '../context/git_url_context'
import { ProjectProvider } from '../context/project_context'

const inter = Inter({ subsets: ['latin'] })

export const metadata = {
  title: 'Hoister',
  description: 'Deploy Your GitHub Projects with Ease',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <GitURLProvider>
          <ProjectProvider>
            {children}
          </ProjectProvider>
        </GitURLProvider>
      </body>
    </html>
  )
}