/// <reference types="vite/client" />

import { ComponentCustomProperties } from 'vue'

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    window: Window
  }
}
