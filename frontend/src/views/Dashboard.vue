<script setup lang="ts">
import Logo from '@/components/ui/Logo.vue'
import { useAuthUser } from '@/router/auth/AuthUserProvider'
import { Label } from '@/components/ui/label'
import Button from '@/components/ui/Button.vue'
import { Archive, ChevronDown, Folder, Home, PanelLeft, Search, Upload } from 'lucide-vue-next'
import { Input } from '@/components/ui/input'
import { TableHeader, Table, TableRow, TableHead } from '@/components/ui/table'
import DashboardLayout from '@/components/DashboardLayout.vue'
import { useQuery } from '@tanstack/vue-query'
import type { TeamResponse } from '@/lib/types'
import { useRoute, useRouter } from 'vue-router'

const { user } = useAuthUser()
const route = useRoute()
const router = useRouter()

const { data: team } = useQuery<TeamResponse>({
  queryKey: ['team', route.params.team],
  queryFn: async () => {
    const response = await fetch(`http://localhost:8080/api/v1/teams?slug=${route.params.team}`, {
      credentials: 'include',
    })
    if (!response.ok) {
      router.push('/auth/login')
      throw new Error((await response.json()).message)
    }
    return response.json() as Promise<TeamResponse>
  },
})
</script>

<template>
  <DashboardLayout>
    <header class="h-[72px] py-4 px-6 flex justify-between items-center">
      <div class="flex gap-4 font-medium items-center">
        <div class="flex gap-2 items-center cursor-pointer text-neutral-400">
          {{ team && team.team.name }}
          <ChevronDown class="size-4 stroke-neutral-400" />
        </div>
        <div class="text-neutral-400">/</div>
        <div class="flex gap-2 items-center cursor-pointer">
          Project
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
      <Table class="rounded-t-lg overflow-clip">
        <TableHeader>
          <TableRow>
            <TableHead> Name </TableHead>
            <TableHead> Last Changed </TableHead>
            <TableHead> Size </TableHead>
          </TableRow>
        </TableHeader>
      </Table>
    </div>
  </DashboardLayout>
</template>
