<script setup lang="ts">
import type { TeamsResponse } from '@/lib/types'
import { useQuery } from '@tanstack/vue-query'
import { StatusCodes } from 'http-status-codes'
import { useRouter } from 'vue-router'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

const router = useRouter()

const {
  data: teams,
  isLoading: teamsIsLoading,
  error: teamsError,
  status: teamsStatus,
} = useQuery<TeamsResponse>({
  queryKey: ['teams'],
  queryFn: async () => {
    const response = await fetch('http://localhost:8080/api/v1/user/teams', {
      credentials: 'include',
    })
    if (!response.ok) {
      if (response.status === StatusCodes.UNAUTHORIZED) {
        router.push('/auth/login')
      }
      throw new Error((await response.json()).message)
    }
    return response.json() as Promise<TeamsResponse>
  },
})
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger :as-child="true">
      <slot></slot>
    </DropdownMenuTrigger>
    <DropdownMenuContent v-if="teams">
      <DropdownMenuItem v-for="team in teams.teams" @click="() => router.push(`/-/${team.slug}`)">
        {{ team.name }}
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
