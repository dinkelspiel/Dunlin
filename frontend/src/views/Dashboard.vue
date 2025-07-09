<script setup lang="ts">
import { useAuthUser } from '@/router/auth/AuthUserProvider'

const { user } = useAuthUser()
</script>

<template>
  asd
  {{ JSON.stringify(user) }}
</template>
