'use client'

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import Link from "next/link"
import { useGitURL } from "../../context/git_url_context"
import { useProject } from "../../context/project_context"
import Image from "next/image"

export default function DeployPage() {
  const { gitURL } = useGitURL()
  const { domain, setDomain, projectName, setProjectName, username, setUsername } = useProject()

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    // Here you would typically send this data to your backend

    console.log({ gitURL, projectName, username, domain })
    
    // For now, we'll just log the data
    alert("Deployment information submitted!")
  }

  return (
    <div className="flex flex-col min-h-screen">
      <header className="px-4 lg:px-6 h-14 flex items-center border-b">
        <Link className="flex items-center justify-center" href="/">
          <Image src="/logo.png" alt="Logo" width={20} height={20} className="mr-2" />
          <span className="font-bold text-md">Hoister</span>
        </Link>
        <nav className="ml-auto flex gap-4 sm:gap-6">
          <Link className="text-sm font-medium hover:underline underline-offset-4" href="#">
            Dashboard
          </Link>
          <Link className="text-sm font-medium hover:underline underline-offset-4" href="#">
            Documentation
          </Link>
        </nav>
      </header>
      <main className="flex-1 p-4 md:p-6">
        <Card className="max-w-md mx-auto">
          <CardHeader>
            <CardTitle>Deploy Your Project</CardTitle>
            <CardDescription>Provide additional information about your GitHub project</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit}>
              <div className="grid w-full items-center gap-4">
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="githubUrl">GitHub URL</Label>
                  <Input id="githubUrl" value={gitURL} disabled />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="projectName">Project Name</Label>
                  <Input
                    id="projectName"
                    placeholder="Enter project name"
                    value={projectName}
                    onChange={(e) => setProjectName(e.target.value)}
                    required
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="username">GitHub Username</Label>
                  <Input
                    id="username"
                    placeholder="Enter GitHub username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                  />
                </div>
                <div className="flex flex-col space-y-1.5">
                  <Label htmlFor="domain">Custom Domain (optional)</Label>
                  <Input
                    id="domain"
                    placeholder="Enter custom domain"
                    value={domain}
                    onChange={(e) => setDomain(e.target.value)}
                  />
                </div>
              </div>
              <CardFooter className="flex justify-between mt-6">
                <Button variant="outline" onClick={() => window.history.back()}>Back</Button>
                <Button type="submit">Deploy Project</Button>
              </CardFooter>
            </form>
          </CardContent>
        </Card>
      </main>
      <footer className="flex flex-col gap-2 sm:flex-row py-6 w-full shrink-0 items-center px-4 md:px-6 border-t">
        <p className="text-xs text-gray-500 dark:text-gray-400">© 2024 Hoister. All rights reserved.</p>
        <nav className="sm:ml-auto flex gap-4 sm:gap-6">
          <Link className="text-xs hover:underline underline-offset-4" href="#">
            Terms of Service
          </Link>
          <Link className="text-xs hover:underline underline-offset-4" href="#">
            Privacy
          </Link>
        </nav>
      </footer>
    </div>
  )
}