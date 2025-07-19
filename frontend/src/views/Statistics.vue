<script setup lang="ts">
import Logo from '@/components/ui/Logo.vue'
import { useAuthUser } from '@/router/auth/AuthUserProvider'
import { Label } from '@/components/ui/label'
import Button from '@/components/ui/Button.vue'
import {
  AlertCircle,
  Archive,
  ChartArea,
  ChartPie,
  ChevronDown,
  CornerDownRight,
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
import type { TeamProjectResponse, TeamResponse, StatisticsResponse } from '@/lib/types'
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
import { Card } from '@/components/ui/card'
import { DonutChart } from '@/components/ui/chart-donut'
import StatisticsView from './StatisticsView.vue'

const { authUser } = useAuthUser()
const route = useRoute()
const router = useRouter()

const error = ref('')
const queryClient = useQueryClient()

const { data: statistics } = useQuery<StatisticsResponse>({
  queryKey: ['statistics'],
  queryFn: async () => {
    const url = `${import.meta.env.VITE_API_URL}/api/v1/statistics`
    const response = await fetch(url, {
      credentials: 'include',
    })
    if (!response.ok) {
      error.value = (await response.json()).error
      throw new Error((await response.json()).error)
    }
    return response.json() as Promise<StatisticsResponse>
  },
})

document.title = `Statistics`
</script>

<template>
  <DashboardLayout>
    <header class="h-[72px] py-4 px-6 flex justify-between items-center">
      <div class="flex gap-4 font-medium items-center">
        <Logo />

        <div class="text-neutral-400">/</div>
        <div class="flex gap-2 items-center cursor-pointer">Statistics</div>
      </div>
      <router-link to="/auth/login" v-if="!authUser.value">
        <Button> Log in </Button>
      </router-link>
      <div class="flex items-center gap-4" v-if="authUser.value">
        <router-link :to="`/-/`">
          <Button size="sm"><CornerDownRight class="size-4" /> To Browser </Button>
        </router-link>
      </div>
    </header>
    <StatisticsView v-if="statistics" :statistics="statistics" />
  </DashboardLayout>
</template>
