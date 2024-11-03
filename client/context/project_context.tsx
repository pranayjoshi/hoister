"use client";

import { createContext, useContext, useState, ReactNode } from "react";

interface ProjectContextProps {
  domain: string;
  setDomain: (domain: string) => void;
  projectName: string;
  setProjectName: (name: string) => void;
  username: string;
  setUsername: (username: string) => void;
  projectSlug: string;
  setProjectSlug: (username: string) => void;
  outputURL: string;
  setOutputURL: (username: string) => void;
}

const ProjectContext = createContext<ProjectContextProps | undefined>(
  undefined
);

export const ProjectProvider = ({ children }: { children: ReactNode }) => {
  const [domain, setDomain] = useState<string>("");
  const [projectName, setProjectName] = useState<string>("");
  const [username, setUsername] = useState<string>("");
  const [outputURL, setOutputURL] = useState<string>("");
  const [projectSlug, setProjectSlug] = useState<string>("");

  return (
    <ProjectContext.Provider
      value={{
        domain,
        setDomain,
        projectName,
        setProjectName,
        username,
        setUsername,
        projectSlug,
        setProjectSlug,
        outputURL,
        setOutputURL,
      }}
    >
      {children}
    </ProjectContext.Provider>
  );
};

export const useProject = () => {
  const context = useContext(ProjectContext);
  if (!context) {
    throw new Error("useProject must be used within a ProjectProvider");
  }
  return context;
};
