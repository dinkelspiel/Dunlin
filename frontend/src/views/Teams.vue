<script setup lang="ts">
import Logo from '@/components/ui/Logo.vue'
import { useAuthUser } from '@/router/auth/AuthUserProvider'
import { Label } from '@/components/ui/label'
import Button from '@/components/ui/Button.vue'
import {
  Archive,
  ChevronDown,
  Home,
  Loader2,
  PanelLeft,
  Plus,
  Search,
  Upload,
} from 'lucide-vue-next'
import { Input } from '@/components/ui/input'
import {
  TableHeader,
  Table,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from '@/components/ui/table'
import DashboardLayout from '@/components/DashboardLayout.vue'
import { useQuery } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'
import { watchEffect } from 'vue'
import { StatusCodes } from 'http-status-codes'
import type { StatusError, Team } from '@/lib/types'

type TeamsResponse = {
  message: string
  teams: Team[]
}

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
  <DashboardLayout>
    <header class="h-[72px] py-4 px-6 flex justify-between items-center">
      <div class="flex gap-4 font-medium items-center">Teams</div>
      <div class="flex items-center gap-4">
        <!-- <div class="relative w-[350px] h-8">
          <Search class="size-4 stroke-neutral-400 absolute top-1/2 -translate-y-1/2 left-2" />
          <Input class="px-8" placeholder="Search" />
        </div> -->
        <Button size="sm"><Plus class="size-4" /> New Team </Button>
      </div>
    </header>
    <div class="p-4">
      <Loader2 class="size-4 animate-spin" v-if="teamsIsLoading" />
      <Table class="rounded-t-lg overflow-clip" v-if="!teamsIsLoading && teams">
        <TableHeader>
          <TableRow>
            <TableHead> Name </TableHead>
            <TableHead> Created </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="team in teams.teams"
            @click="() => router.push(`/-/${team.slug}`)"
            class="cursor-pointer hover:underline"
          >
            <TableCell>{{ team.name }}</TableCell>
            <TableCell>{{ new Date(team.createdAt).toDateString() }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </DashboardLayout>
</template>
