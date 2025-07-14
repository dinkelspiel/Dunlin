import './assets/styles.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'

const app = createApp(App)

app.config.globalProperties.window = window

app.use(router)
app.use(VueQueryPlugin)

app.mount('#app')
