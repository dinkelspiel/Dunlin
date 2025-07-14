<script setup lang="ts">
import Logo from '@/components/ui/Logo.vue'
import { useAuthUser } from '@/router/auth/AuthUserProvider'
import { Label } from '@/components/ui/label'
import Button from '@/components/ui/Button.vue'
import {
  AlertCircle,
  Archive,
  ChevronDown,
  File,
  Folder,
  Home,
  PanelLeft,
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
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'
import DashboardLayout from '@/components/DashboardLayout.vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import type { TeamProjectResponse, TeamResponse } from '@/lib/types'
import { useRoute, useRouter } from 'vue-router'
import { computed, ref } from 'vue'
import { watchDeep } from '@vueuse/core'
import normalize from 'path-normalize'

const { user } = useAuthUser()
const route = useRoute()
const router = useRouter()

type File = {
  type: 'dir' | 'file'
  name: string
  lastModified: string
  size: number
}

type FilesResponse = {
  message: string
  files: File[]
}

function sortFiles(files: File[]): File[] {
  return files.sort((a, b) => {
    if (a.type !== b.type) {
      return a.type === 'dir' ? 1 : -1
    }
    return a.name.localeCompare(b.name, undefined, { sensitivity: 'base' })
  })
}

const { data: team } = useQuery<TeamResponse>({
  queryKey: ['team', route.params.team],
  queryFn: async () => {
    const response = await fetch(`http://localhost:8080/api/v1/teams/${route.params.team}`, {
      credentials: 'include',
    })
    if (!response.ok) {
      router.push('/auth/login')
      throw new Error((await response.json()).message)
    }
    return response.json() as Promise<TeamResponse>
  },
})

const { data: teamProject } = useQuery<TeamProjectResponse>({
  queryKey: ['teamProject', route.params.team, route.params.project],
  queryFn: async () => {
    const response = await fetch(
      `http://localhost:8080/api/v1/teams/${route.params.team}/projects/${route.params.project}`,
      {
        credentials: 'include',
      },
    )
    if (!response.ok) {
      router.push('/auth/login')
      throw new Error((await response.json()).message)
    }
    return response.json() as Promise<TeamProjectResponse>
  },
})

const error = ref('')
const queryClient = useQueryClient()

const rawFilepath = computed(() =>
  normalize(Array.isArray(route.params.filepath) ? route.params.filepath.join('/') : ''),
)

const filepathWithSlashes = computed(() => (rawFilepath.value ? `/${rawFilepath.value}/` : '/'))

const { data: files } = useQuery<FilesResponse>({
  queryKey: () => ['files', route.params.team, route.params.project, filepathWithSlashes],
  queryFn: async () => {
    const url = `http://localhost:8080/api/v1/teams/${route.params.team}/projects/${route.params.project}/files/${filepathWithSlashes.value}`
    const response = await fetch(url, {
      credentials: 'include',
    })
    if (!response.ok) {
      error.value = (await response.json()).error
      throw new Error((await response.json()).error)
    }
    return response.json() as Promise<FilesResponse>
  },
})
</script>

<template>
  <DashboardLayout>
    <header class="h-[72px] py-4 px-6 flex justify-between items-center">
      <div class="flex gap-4 font-medium items-center">
        <Logo />

        <div class="flex gap-2 items-center cursor-pointer text-neutral-400">
          {{ team && team.team.name }}
          <ChevronDown class="size-4 stroke-neutral-400" />
        </div>
        <div class="text-neutral-400">/</div>
        <div class="flex gap-2 items-center cursor-pointer">
          {{ teamProject && teamProject.teamProject.name }}
          <ChevronDown class="size-4 stroke-neutral-600" />
        </div>
      </div>
      <div class="flex items-center gap-4">
        <div class="relative w-[350px] h-8">
          <Search class="size-4 stroke-neutral-400 absolute top-1/2 -translate-y-1/2 left-2" />
          <Input class="px-8" placeholder="Search" />
        </div>
        <Button size="sm"><Upload class="size-4" /> Upload </Button>
        <Button size="sm"><Folder class="size-4" /> New Folder </Button>
      </div>
    </header>
    <div class="p-4">
      <Alert v-if="error" variant="destructive">
        <AlertCircle class="w-4 h-4" />
        <AlertTitle>Error</AlertTitle>
        <AlertDescription> {{ error }} </AlertDescription>
      </Alert>
      <Table class="rounded-t-lg overflow-clip" v-if="!error">
        <TableHeader>
          <TableRow>
            <TableHead> Name </TableHead>
            <TableHead> Last Changed </TableHead>
            <TableHead> Size </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody v-if="files">
          <TableRow
            class="hover:underline cursor-pointer"
            v-for="file in [
              ...sortFiles(files.files),
              rawFilepath !== '.'
                ? {
                    type: 'dir',
                    name: '..',
                    lastModified: '',
                    size: '',
                  }
                : undefined,
            ].filter((a) => !!a)"
            v-bind:key="file.name"
            @click="
              () => {
                if (file.type === 'dir') {
                  queryClient.invalidateQueries({ queryKey: ['files'] })
                  router.replace(
                    `/-/${route.params.team}/${route.params.project}${filepathWithSlashes}${file.name}`,
                  )
                } else {
                  window.location.href = `http://localhost:8080/files/${route.params.team}/${route.params.project}${filepathWithSlashes}${file.name}`
                }
              }
            "
          >
            <TableCell>
              <div class="flex gap-2 items-center">
                <File class="size-4 stroke-neutral-600" v-if="file.type === 'file'" />
                <Folder class="size-4 stroke-neutral-600" v-if="file.type === 'dir'" />
                {{ file.name }}
              </div>
            </TableCell>
            <TableCell>{{
              file.lastModified !== '' ? new Date(file.lastModified).toDateString() : ''
            }}</TableCell>
            <TableCell
              ><div v-if="file.type === 'file'">{{ file.size }} B</div></TableCell
            >
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </DashboardLayout>
</template>
