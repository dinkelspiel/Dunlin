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
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'
import { ref, watch, watchEffect } from 'vue'
import { StatusCodes } from 'http-status-codes'
import type { StatusError, Team, TeamsResponse } from '@/lib/types'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'

const { authUser } = useAuthUser()

const router = useRouter()

const {
  data: teams,
  isLoading: teamsIsLoading,
  error: teamsError,
  status: teamsStatus,
} = useQuery<TeamsResponse>({
  queryKey: ['teams'],
  queryFn: async () => {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/v1/user/teams`, {
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

watch([authUser], () => {
  if (!authUser.value) {
    document.title = 'Teams for'
    return
  }
  document.title = `Teams for ${authUser.value?.username}`
})

const createTeamOpen = ref(false)
const teamName = ref('')
const queryClient = useQueryClient()

const createTeam = useMutation({
  mutationKey: ['createProject'],
  mutationFn: async (teamName: string) => {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/v1/teams`, {
      method: 'POST',
      credentials: 'include',
      body: JSON.stringify({
        teamName,
      }),
    })
    if (!response.ok) {
      throw new Error((await response.json()).message)
    }
    return response.json()
  },
  onSuccess() {
    queryClient.invalidateQueries({ queryKey: ['teams'] })
    createTeamOpen.value = false
  },
})
</script>

<template>
  <DashboardLayout>
    <header class="h-[72px] py-4 px-6 flex justify-between items-center">
      <div class="flex gap-4 font-medium items-center">
        <Logo />
        <router-link :to="`/-`">
          <div class="text-neutral-400">/</div>
        </router-link>
        Teams
      </div>

      <div class="flex items-center gap-4">
        <Dialog v-model:open="createTeamOpen">
          <DialogTrigger :as-child="true">
            <Button size="sm"><Plus class="size-4" /> New Team </Button>
          </DialogTrigger>
          <DialogContent :show-close="true">
            <DialogHeader>
              <DialogTitle> New Team </DialogTitle>
            </DialogHeader>
            <form class="grid gap-4" @submit.prevent="() => createTeam.mutate(teamName)">
              <Input v-model="teamName" placeholder="Name" />
              <DialogFooter>
                <Button type="submit">Create</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
        <!-- <div class="relative w-[350px] h-8">
          <Search class="size-4 stroke-neutral-400 absolute top-1/2 -translate-y-1/2 left-2" />
          <Input class="px-8" placeholder="Search" />
        </div> -->
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
            v-bind:key="team.id"
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
