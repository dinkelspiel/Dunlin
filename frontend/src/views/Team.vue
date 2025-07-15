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
import { useRoute, useRouter } from 'vue-router'
import { QueryClient, useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import type {
  StatusError,
  Team,
  TeamProject,
  TeamProjectsResponse,
  TeamResponse,
  TeamsResponse,
} from '@/lib/types'
import { StatusCodes } from 'http-status-codes'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { ref } from 'vue'
import TeamsDropdown from '@/components/header/TeamsDropdown.vue'

const router = useRouter()
const route = useRoute()
const { authUser } = useAuthUser()

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

const { data: team } = useQuery<TeamResponse>({
  queryKey: ['team', route.params.team],
  queryFn: async () => {
    const response = await fetch(`http://localhost:8080/api/v1/teams/${route.params.team}`, {
      credentials: 'include',
    })
    if (!response.ok) {
      // router.push('/auth/l')
      throw new Error((await response.json()).message)
    }
    return response.json() as Promise<TeamResponse>
  },
})

// const { data: project } = useQuery<TeamProjectsResponse>({
//   queryKey: ['team', route.params.team],
//   queryFn: async () => {
//     const response = await fetch(`http://localhost:8080/api/v1/teams/${route.params.team}`, {
//       credentials: 'include',
//     })
//     if (!response.ok) {
//       router.push('/')
//       throw new Error((await response.json()).message)
//     }
//     return response.json() as Promise<TeamResponse>
//   },
// })

const projectName = ref('')
const queryClient = useQueryClient()

const createProject = useMutation({
  mutationKey: ['createProject'],
  mutationFn: async (projectName: string) => {
    const response = await fetch(
      `http://localhost:8080/api/v1/teams/${route.params.team}/projects`,
      {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
          teamSlug: route.params.team,
          projectName,
        }),
      },
    )
    if (!response.ok) {
      throw new Error((await response.json()).message)
    }
    return response.json()
  },
  onSuccess() {
    queryClient.invalidateQueries({ queryKey: ['teamProjects'] })
    createProjectOpen.value = false
  },
})

const createProjectOpen = ref(false)
</script>

<template>
  <DashboardLayout>
    <header class="h-[72px] py-4 px-6 flex justify-between items-center">
      <div class="flex gap-4 font-medium items-center">
        <Logo />
        <router-link :to="`/-`">
          <div class="text-neutral-400">/</div>
        </router-link>
        <TeamsDropdown>
          <div class="flex gap-2 items-center cursor-pointer">
            {{ team && team.team.name }}
            <ChevronDown class="size-4 stroke-neutral-600" />
          </div>
        </TeamsDropdown>
      </div>
      <div class="flex items-center gap-4">
        <Dialog v-model:open="createProjectOpen">
          <DialogTrigger :as-child="true">
            <Button size="sm"><Plus class="size-4" /> New Project </Button>
          </DialogTrigger>
          <DialogContent :show-close="true">
            <DialogHeader>
              <DialogTitle> New Project </DialogTitle>
            </DialogHeader>
            <form class="grid gap-4" @submit.prevent="() => createProject.mutate(projectName)">
              <Input v-model="projectName" placeholder="Name" />
              <DialogFooter>
                <Button type="submit">Create</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </div>
    </header>
    <div class="p-4">
      <Loader2 class="animate-spin size-4" v-if="teamProjectsIsLoading" />
      <Table class="rounded-t-lg overflow-clip" v-if="teamProjects">
        <TableHeader>
          <TableRow>
            <TableHead> Name </TableHead>
            <TableHead> Created </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-bind:key="teamProject.id"
            v-for="teamProject in teamProjects.teamProjects"
            @click="() => router.push(`/-/${route.params.team}/${teamProject.slug}`)"
            class="cursor-pointer hover:underline"
          >
            <TableCell>{{ teamProject.name }}</TableCell>
            <TableCell>{{ new Date(teamProject.createdAt).toDateString() }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </DashboardLayout>
</template>
