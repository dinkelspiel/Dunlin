<script setup lang="ts">
import Logo from '@/components/ui/Logo.vue'
import { useAuthUser } from '@/router/auth/AuthUserProvider'
import { Label } from '@/components/ui/label'
import Button from '@/components/ui/Button.vue'
import { Archive, ChevronDown, Home, PanelLeft, Search, Upload } from 'lucide-vue-next'
import { Input } from '@/components/ui/input'
import { TableHeader, Table, TableRow, TableHead } from '@/components/ui/table'

const { user } = useAuthUser()
</script>

<template>
  <!-- grid-cols-[min-content_1fr] -->
  <div class="grid min-h-screen">
    <!-- <div class="w-[270px] border-e border-e-neutral-200 p-4 flex flex-col gap-8">
      <div class="flex justify-between items-center">
        <Logo />
        <Button variant="ghost" size="icon">
          <PanelLeft class="size-4" />
        </Button>
      </div>

      <div class="flex flex-col">
        <Button variant="secondary" size="sm" class="justify-start">
          <Archive class="size-4 stroke-neutral-600" /> Files
        </Button>
      </div>
    </div> -->
    <div>
      <slot></slot>
    </div>
  </div>
</template>
