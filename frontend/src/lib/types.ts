export type User = {
  id: number
  email: string
  username: string
  updatedAt: string | null
  createdAt: string
}

export type Team = {
  id: number
  name: string
  slug: string
  owner: User
  ownerId: number
  updatedAt: string | null
  createdAt: string
}

export type TeamProject = {
  id: number
  name: string
  slug: string
  team: Team
  teamId: number
  updatedAt: string | null
  createdAt: string
}

export type StatusError = Error & { statusCode: number }

export type TeamProjectsResponse = {
  message: string
  teamProjects: TeamProject[]
}

export type TeamResponse = {
  message: string
  team: Team
}

export type TeamProjectResponse = {
  message: string
  teamProject: TeamProject
}

export type TeamsResponse = {
  message: string
  teams: Team[]
}

export type StatisticsResponse = {
  message: string
  data: {
    disk: {
      total: number
      free: number
      hostUsed: number
      dunlinFilesUsed: number
      dunlinCacheUsed: number
    }
    teamProjectSizes: {
      size: number
      teamProject: TeamProject
    }[]
  }
}
