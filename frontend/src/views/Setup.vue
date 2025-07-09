<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { H3, SubText } from '@/components/typography'
import Logo from '@/components/ui/Logo.vue'
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const adminUsername = ref('')
const adminEmail = ref('')

const router = useRouter()

onMounted(async () => {
  const response = await fetch('http://localhost:8080/api/v1/setup', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (response.status === 404) {
    router.replace('/auth/login')
    return
  }
})

const onSubmit = async () => {
  const response = await fetch('http://localhost:8080/api/v1/setup', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      adminUsername: adminUsername.value,
      adminEmail: adminEmail.value,
    }),
  })

  if (response.status === 200) {
    router.replace('/auth/login')
  }
}
</script>

<template>
  <div class="bg-neutral-100 grid h-[100dvh] w-full items-center justify-center">
    <form class="grid w-[350px] gap-8" @submit.prevent="onSubmit">
      <div class="grid gap-2">
        <Logo class="size-11 mb-2 rounded-lg" />
        <H3>Setup OpenCDN</H3>
        <SubText> OpenCDN is a lightweight, self-hosted CDN for personal projects. </SubText>
      </div>
      <div class="grid gap-4">
        <div class="grid gap-2">
          <Label>Administrator Account</Label>
          <Input v-model="adminUsername" placeholder="Username" name="username" />
          <Input v-model="adminEmail" placeholder="Email" name="email" />
        </div>
        <Button size="sm"> Setup </Button>
      </div>
    </form>
  </div>
</template>
