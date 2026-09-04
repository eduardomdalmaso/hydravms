import { createApp } from 'vue'
import App from './App.vue'
import './assets/css/main.css'

// Bloqueia o menu de contexto padrao do navegador para comportamento nativo de VMS
if (typeof window !== 'undefined') {
  window.addEventListener('contextmenu', (e) => {
    e.preventDefault()
  })
}

createApp(App).mount('#app')
