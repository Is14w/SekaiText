import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import { useDebugLog } from './composables/useDebugLog'
import './style.css'
import App from './App.vue'

const app = createApp(App)
app.use(createPinia())
app.use(router)

// Initialize console log capture for debug panel
const debug = useDebugLog()
debug.initConsoleCapture()

app.config.errorHandler = (err, instance, info) => {
  console.error('[Vue Error]', err, info)
}

app.mount('#app')
