<script setup lang="ts">
const props = defineProps(['class'])

import { Bird } from 'lucide-vue-next'
import { cn } from '@/lib/utils'
</script>

<template>
  <div
    :class="
      cn(
        'flex size-[40px] items-center justify-center rounded-lg border border-neutral-200 p-1 shadow-sm',
        props.class,
      )
    "
  >
    <div class="flex h-full w-full items-center justify-center rounded-[6px] bg-neutral-900">
      <Bird class="flex size-4 items-center justify-center stroke-white" />
    </div>
  </div>
</template>
