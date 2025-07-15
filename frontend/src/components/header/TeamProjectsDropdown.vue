<script setup lang="ts">
import type { TeamProjectsResponse, TeamsResponse } from '@/lib/types'
import { useQuery } from '@tanstack/vue-query'
import { StatusCodes } from 'http-status-codes'
import { useRoute, useRouter } from 'vue-router'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

const router = useRouter()
const route = useRoute()

const {
  data: teamProjects,
  isLoading: teamProjectsIsLoading,
  error: teamProjectsError,
  status: teamProjectsStatus,
} = useQuery<TeamProjectsResponse>({
  queryKey: ['teamProjects', route.params.team],
  queryFn: async () => {
    const response = await fetch(
      `http://localhost:8080/api/v1/teams/${route.params.team}/projects`,
      {
        credentials: 'include',
      },
    )
    if (!response.ok) {
      if (response.status === StatusCodes.UNAUTHORIZED) {
        router.push('/auth/login')
      }
      throw new Error((await response.json()).message)
    }
    return response.json() as Promise<TeamProjectsResponse>
  },
})
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger :as-child="true">
      <slot></slot>
    </DropdownMenuTrigger>
    <DropdownMenuContent v-if="teamProjects">
      <DropdownMenuItem
        v-for="teamProject in teamProjects.teamProjects"
        @click="() => router.push(`/-/${teamProject.team.slug}/${teamProject.slug}`)"
      >
        {{ teamProject.name }}
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
