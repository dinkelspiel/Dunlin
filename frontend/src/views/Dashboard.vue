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
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import type { TeamProjectResponse, TeamResponse } from '@/lib/types'
import { useRoute, useRouter } from 'vue-router'
import { computed, Fragment, onMounted, ref, watch } from 'vue'
import { watchDeep } from '@vueuse/core'
import normalize from 'path-normalize'
import TeamsDropdown from '@/components/header/TeamsDropdown.vue'
import TeamProjectsDropdown from '@/components/header/TeamProjectsDropdown.vue'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { DropdownMenuItem } from '@/components/ui/dropdown-menu'
import FileUploader from '@/components/FileUploader.vue'
import { humanFileSize } from '@/lib/bytes'
import Breadcrumbs from '@/components/Breadcrumbs.vue'

const { authUser } = useAuthUser()
const route = useRoute()
const router = useRouter()
const teamSlug = computed(() => route.params.team as string)
const projectSlug = computed(() => route.params.project as string)
const apiUrl = import.meta.env.VITE_API_URL

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

const { data: team } = useQuery<TeamResponse>({
  queryKey: ['team', teamSlug],
  queryFn: async () => {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/v1/teams/${teamSlug.value}`, {
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
  queryKey: ['teamProject', teamSlug, projectSlug],
  queryFn: async () => {
    const response = await fetch(
      `${import.meta.env.VITE_API_URL}/api/v1/teams/${teamSlug.value}/projects/${projectSlug.value}`,
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

const createFolderOpen = ref(false)
const folderPath = ref('')

const createFolder = useMutation({
  mutationKey: ['createProject'],
  mutationFn: async (path: string) => {
    const response = await fetch(
      `${import.meta.env.VITE_API_URL}/api/v1/teams/${teamSlug.value}/projects/${projectSlug.value}/folders`,
      {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
          path: `${filepathWithSlashes.value}/${path.startsWith('/') ? path.substring(1) : path}`,
        }),
      },
    )
    if (!response.ok) {
      throw new Error((await response.json()).message)
    }
    return response.json()
  },
  onSuccess() {
    queryClient.invalidateQueries({ queryKey: ['files'] })
    createFolderOpen.value = false
    folderPath.value = ''
  },
})

const error = ref('')
const queryClient = useQueryClient()

const rawFilepath = computed(() =>
  normalize(Array.isArray(route.params.filepath) ? route.params.filepath.join('/') : ''),
)

const filepathWithSlashes = computed(() => (rawFilepath.value ? `/${rawFilepath.value}/` : '/'))

const { data: files } = useQuery<FilesResponse>({
  queryKey: ['files', teamSlug, projectSlug, filepathWithSlashes],
  queryFn: async () => {
    const url = `${import.meta.env.VITE_API_URL}/api/v1/teams/${teamSlug.value}/projects/${projectSlug.value}/files/${filepathWithSlashes.value}`
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

watch([team, teamProject], () => {
  if (!team.value && !teamProject.value) {
    document.title = 'Index of'
    return
  }
  document.title = `Index of /${team.value?.team.slug}/${teamProject.value?.teamProject.slug}${filepathWithSlashes.value}`
})
</script>

<template>
  <DashboardLayout>
    <header class="h-[72px] py-4 px-6 flex justify-between items-center">
      <div class="flex gap-4 font-medium items-center">
        <Logo />

        <router-link :to="`/-`" v-if="authUser.value">
          <div class="text-neutral-400">/</div>
        </router-link>
        <div class="text-neutral-400" v-if="!authUser.value">/</div>

        <TeamsDropdown v-if="authUser.value">
          <div class="flex gap-2 items-center cursor-pointer text-neutral-400">
            {{ team && team.team.name }}
            <ChevronDown class="size-4 stroke-neutral-400" />
          </div>
        </TeamsDropdown>
        <div v-if="!authUser.value" class="flex gap-2 items-center text-neutral-400">
          {{ team && team.team.name }}
        </div>

        <router-link :to="`/-/${team && team.team.slug}`" v-if="authUser.value">
          <div class="text-neutral-400">/</div>
        </router-link>
        <div class="text-neutral-400" v-if="!authUser.value">/</div>

        <TeamProjectsDropdown v-if="authUser.value">
          <div class="flex gap-2 items-center cursor-pointer">
            {{ teamProject && teamProject.teamProject.name }}
            <ChevronDown class="size-4 stroke-neutral-600" />
          </div>
        </TeamProjectsDropdown>
        <div v-if="!authUser.value" class="flex gap-2 items-center">
          {{ teamProject && teamProject.teamProject.name }}
        </div>
        <Breadcrumbs
          :team-slug="route.params.team as string"
          :project-slug="route.params.project as string"
          :filepath="filepathWithSlashes"
        />
      </div>
      <router-link to="/auth/login" v-if="!authUser.value">
        <Button> Log in </Button>
      </router-link>
      <div class="flex items-center gap-4" v-if="authUser.value">
        <Dialog>
          <DialogTrigger>
            <div class="relative w-[350px] h-8">
              <Search class="size-4 stroke-neutral-400 absolute top-1/2 -translate-y-1/2 left-2" />
              <Input class="px-8" placeholder="Search" />
            </div>
          </DialogTrigger>
          <DialogContent
            :show-close="false"
            class="p-0 gap-0 divide-y divide-y-neutral-200 border-0"
          >
            <div class="relative h-12">
              <Search class="size-4 stroke-neutral-400 absolute top-1/2 -translate-y-1/2 left-2" />
              <Input class="px-8 h-12 rounded-b-none" placeholder="Search" />
            </div>
            <div class="h-[350px] flex flex-col gap-1 overflow-y-scroll p-2 no-scrollbar">
              <Button class="w-full" size="sm" variant="secondary"> Test </Button>
            </div>
          </DialogContent>
        </Dialog>
        <FileUploader
          :team-slug="route.params.team as string"
          :project-slug="route.params.project as string"
          :target-path="filepathWithSlashes"
        />
        <Dialog v-model:open="createFolderOpen">
          <DialogTrigger :as-child="true">
            <Button size="sm"><Folder class="size-4" /> New Folder </Button>
          </DialogTrigger>
          <DialogContent :show-close="true">
            <DialogHeader>
              <DialogTitle> Create a folder </DialogTitle>
            </DialogHeader>
            <form
              @submit.prevent="() => createFolder.mutate(folderPath)"
              class="flex flex-col gap-4"
            >
              <Input v-model="folderPath" placeholder="Name" />
              <DialogFooter>
                <Button> Create </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
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
              rawFilepath !== '.'
                ? {
                    type: 'dir',
                    name: '..',
                    lastModified: '',
                    size: '',
                  }
                : null,
              ...files.files.slice().sort((a, b) => {
                if (a.type === 'dir' && b.type !== 'dir') return -1
                if (a.type !== 'dir' && b.type === 'dir') return 1
                return a.name.localeCompare(b.name, undefined, { sensitivity: 'base' })
              }),
            ].filter((a) => !!a)"
            v-bind:key="`${filepathWithSlashes}${file.name}`"
            @click="
              () => {
                if (file.type === 'dir') {
                  queryClient.invalidateQueries({ queryKey: ['files'] })
                  router.replace(
                    `/-/${route.params.team}/${route.params.project}${filepathWithSlashes}${file.name}`,
                  )
                } else {
                  // Window is added but the types don't seem to work see /frontend/env.d.ts
                  // @ts-ignore
                  window.location.href = `${apiUrl}/files/${route.params.team}/${route.params.project}${filepathWithSlashes}${file.name}`
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
              ><div v-if="file.type === 'file'">
                {{ humanFileSize(file.size as number) }}
              </div></TableCell
            >
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </DashboardLayout>
</template>
