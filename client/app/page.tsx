'use client'

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { GithubIcon } from "lucide-react"
import Link from "next/link"
import Image from "next/image"
import { useState } from "react"
import { useRouter } from "next/navigation"
import { useGitURL } from "../context/git_url_context"

export default function HomePage() {
  const [gitURLInput, setGitURLInput] = useState("")
  const router = useRouter()
  const { setGitURL } = useGitURL()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    console.log(gitURLInput);

    try {
      const parsedURL = new URL(gitURLInput);
      const pathSegments = parsedURL.pathname.split('/').filter(segment => segment);
      console.log(pathSegments);
      if (pathSegments.length >= 2) {
        setGitURL(gitURLInput);
        router.push('/deploy');
      }
    } catch (error) {
      console.error('Invalid URL:', error);
    }
  };

  return (
    <div className="flex flex-col min-h-screen">
      <header className="px-4 lg:px-6 h-14 flex items-center">
        <Link className="flex items-center justify-center" href="/">
          <Image src="/logo.png" alt="Logo" width={20} height={20} className="mr-2" />
          <span className="font-bold text-md">Hoister</span>
        </Link>
        <nav className="ml-auto flex gap-4 sm:gap-6">
          <Link className="text-sm font-medium hover:underline underline-offset-4" href="#">
            Features
          </Link>
          <Link className="text-sm font-medium hover:underline underline-offset-4" href="#">
            Pricing
          </Link>
          <Link className="text-sm font-medium hover:underline underline-offset-4" href="#">
            Documentation
          </Link>
        </nav>
      </header>
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32 xl:py-48">
          <div className="container px-4 md:px-6">
            <div className="flex flex-col items-center space-y-4 text-center">
              <div className="space-y-2">
                <h1 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl lg:text-6xl/none">
                  Deploy Your GitHub Projects with Ease
                </h1>
                <p className="mx-auto max-w-[700px] text-gray-500 md:text-xl dark:text-gray-400">
                  Hoister makes it simple to deploy your GitHub projects. Just paste your repository URL and we'll handle the rest.
                </p>
              </div>
              <div className="w-full max-w-sm space-y-2">
                <form className="flex space-x-2" onSubmit={handleSubmit}>
                  <Input
                    className="flex-1"
                    placeholder="Enter GitHub URL"
                    type="url"
                    value={gitURLInput}
                    onChange={(e) => setGitURLInput(e.target.value)}
                  />
                  <Button type="submit">Deploy</Button>
                </form>
                <p className="text-xs text-gray-500 dark:text-gray-400">
                  By deploying, you agree to our Terms of Service and Privacy Policy.
                </p>
              </div>
            </div>
          </div>
        </section>
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