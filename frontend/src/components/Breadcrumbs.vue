<template>
  <div class="flex gap-4 items-center" v-for="(path, index) in pathSegments" :key="index">
    <router-link :to="getPathLink(index)">
      <div class="text-neutral-400">/</div>
    </router-link>
    <router-link :to="getPathLink(index)" class="flex gap-2 items-center cursor-pointer">
      {{ path }}
    </router-link>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  filepath: string
  teamSlug: string
  projectSlug: string
}>()

const pathSegments = computed(() => props.filepath.split('/').filter((e) => !!e))

const getPathLink = (index: number) => {
  const partialPath = pathSegments.value.slice(0, index + 1).join('/')
  return `/-/${props.teamSlug}/${props.projectSlug}/${partialPath}`
}
</script>
