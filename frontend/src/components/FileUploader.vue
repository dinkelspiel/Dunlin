<script setup lang="ts">
import { Upload } from 'lucide-vue-next'
import Button from './ui/Button.vue'
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { useQueryClient } from '@tanstack/vue-query'

const fileInput = ref<HTMLInputElement | null>(null)
const queryClient = useQueryClient()

const props = defineProps<{
  teamSlug: string
  projectSlug: string
  targetPath: string
}>()

const triggerFilePicker = (): void => {
  fileInput.value?.click()
}

const handleFileChange = async (event: Event): Promise<void> => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  try {
    const response = await fetch(
      `http://localhost:8080/api/v1/teams/${props.teamSlug}/projects/${props.projectSlug}/files${props.targetPath}${file.name}`,
      {
        credentials: 'include',
        method: 'PUT',
        body: formData,
      },
    )

    const result = await response.json()

    if (response.ok) {
      queryClient.invalidateQueries({ queryKey: ['files'] })
    } else {
      alert(`Error: ${result.error}`)
    }
  } catch (error) {
    console.error('Upload failed:', error)
    alert('Upload failed.')
  } finally {
    // Reset the file input so same file can be picked again
    input.value = ''
  }
}
</script>

<template>
  <input type="file" ref="fileInput" style="display: none" @change="handleFileChange" />
  <Button size="sm" @click="triggerFilePicker"><Upload class="size-4" /> Upload </Button>
</template>
