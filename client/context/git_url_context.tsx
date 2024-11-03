'use client'

import { createContext, useContext, useState, ReactNode } from 'react';

interface GitURLContextProps {
  gitURL: string;
  setGitURL: (url: string) => void;
}

const GitURLContext = createContext<GitURLContextProps | undefined>(undefined);

export const GitURLProvider = ({ children }: { children: ReactNode }) => {
  const [gitURL, setGitURL] = useState<string>('');

  return (
    <GitURLContext.Provider value={{ gitURL, setGitURL }}>
      {children}
    </GitURLContext.Provider>
  );
};

export const useGitURL = () => {
  const context = useContext(GitURLContext);
  if (!context) {
    throw new Error('useGitURL must be used within a GitURLProvider');
  }
  return context;
};