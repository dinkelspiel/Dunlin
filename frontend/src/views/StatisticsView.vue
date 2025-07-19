<script setup lang="ts">
import { Card, CardHeader, CardTitle } from '@/components/ui/card'
import { DonutChart } from '@/components/ui/chart-donut'
import { humanFileSize } from '@/lib/bytes'
import type { StatisticsResponse } from '@/lib/types'

const props = defineProps<{
  statistics: StatisticsResponse
}>()

const diskUsageColors = ['#60a5fa', '#4ade80', '#facc15', '#e5e5e5']
const rainbowColors = [
  '#3b82f6',
  '#10b981',
  '#f43f5e',
  '#84cc16',
  '#d946ef',
  '#06b6d4',
  '#8b5cf6',
  '#f59e0b',
  '#ef4444',
]
const diskUsageData = [
  {
    name: 'Used by Dunlin',
    size: props.statistics.data.disk.dunlinFilesUsed,
  },
  {
    name: 'Used by Dunlin Cache',
    size: props.statistics.data.disk.dunlinCacheUsed,
  },
  {
    name: 'Used by Other',
    size: props.statistics.data.disk.hostUsed,
  },
  {
    name: 'Free',
    size: props.statistics.data.disk.free,
  },
]

const teamProjectsSize = props.statistics.data.teamProjectSizes
  ? props.statistics.data.teamProjectSizes.map((e) => ({
      name: e.teamProject.name,
      size: e.size,
    }))
  : []

function valueFormatter(tick: number | Date) {
  return typeof tick === 'number' ? `${humanFileSize(tick).toString()}` : ''
}
</script>
<template>
  <div class="px-6 pb-4">
    <div class="grid xl:grid-cols-2 gap-4">
      <Card>
        <CardHeader>
          <CardTitle>Size on disk</CardTitle>
        </CardHeader>
        <div class="p-4 flex justify-center gap-4">
          <DonutChart
            index="name"
            :colors="diskUsageColors"
            :category="'size'"
            :data="diskUsageData"
            :value-formatter="valueFormatter"
            class="w-[230px]"
          />
          <div class="flex flex-col gap-2">
            <div
              class="flex items-center gap-2"
              v-bind:key="entry.name"
              v-for="(entry, idx) in diskUsageData"
            >
              <div
                class="size-3 rounded-full"
                :style="`background-color: ${diskUsageColors[idx]}`"
              />
              {{ entry.name }}
            </div>
          </div>
        </div>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Team Projects by Size</CardTitle>
        </CardHeader>
        <div class="p-4 flex justify-center gap-4">
          <DonutChart
            index="name"
            :colors="rainbowColors"
            :category="'size'"
            :data="teamProjectsSize"
            :value-formatter="valueFormatter"
            class="w-[230px]"
          />
          <div class="flex flex-col gap-2">
            <div
              class="flex items-center gap-2"
              v-bind:key="entry.name"
              v-for="(entry, idx) in teamProjectsSize"
            >
              <div class="size-3 rounded-full" :style="`background-color: ${rainbowColors[idx]}`" />
              {{ entry.name }}
            </div>
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>
