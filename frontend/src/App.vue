<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'

// import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <!-- <header>
    <div>
      <nav>
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/about">About</RouterLink>
      </nav>
    </div>
  </header> -->

  <div class="text-neutral-900">
    <RouterView />
  </div>
</template>
