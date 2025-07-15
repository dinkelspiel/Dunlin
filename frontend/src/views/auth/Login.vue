<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
import { Input } from '@/components/ui/input'
import { H3, SubText, P } from '@/components/typography'
import Logo from '@/components/ui/Logo.vue'
import { onMounted, ref } from 'vue'
import { PinInput, PinInputGroup, PinInputSlot } from '@/components/ui/pin-input'
import { useRouter } from 'vue-router'

onMounted(async () => {
  document.title = 'Log in'
})

const verificationCode = ref<string[]>([])
const handleVerificationCodeComplete = async (e: string[]) => {
  const response = await fetch(`${import.meta.env.VITE_API_URL}/api/v1/auth/verify-code`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      code: parseInt(e.join('')),
    }),
  })

  if (response.status === 201) {
    router.push('/')
  }
}

const email = ref('')
const message = ref<string | null>(null)
const router = useRouter()

const onSubmitSendCode = async () => {
  const response = await fetch(`${import.meta.env.VITE_API_URL}/api/v1/auth/send-code`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: email.value,
    }),
  })

  if (response.status === 201) {
    const json = await response.json()
    message.value = json.message
  }
}
</script>

<template>
  <div class="bg-neutral-100 grid h-[100dvh] w-full items-center justify-center">
    <form class="grid w-[350px] gap-8" @submit.prevent="onSubmitSendCode">
      <div class="grid gap-2">
        <Logo class="size-11 mb-2 rounded-lg" />
        <H3>Login to OpenCDN</H3>
        <SubText> OpenCDN is a lightweight, self-hosted CDN for personal projects. </SubText>
      </div>
      <div class="grid gap-4">
        <div class="grid gap-2">
          <Input v-if="!message" v-model="email" placeholder="Email" name="email" />
          <div class="flex gap-2 flex-col items-center text-center" v-if="message">
            <P v-if="message">{{ message }}</P>
            <PinInput
              id="pin-input"
              v-model="verificationCode"
              placeholder="○"
              @complete="handleVerificationCodeComplete"
            >
              <PinInputGroup>
                <PinInputSlot v-for="(id, index) in 5" :key="id" :index="index" />
              </PinInputGroup>
            </PinInput>
          </div>
        </div>
        <Button size="sm" v-if="!message" type="submit"> Log in </Button>
        <!-- <SubText>
          Don't have an account? <a class="text-primary" href="/auth/sign-up">Sign up</a>
        </SubText> -->
      </div>
    </form>
  </div>
</template>
