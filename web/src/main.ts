import { createApp } from 'vue'
import App from './App.vue'
import './assets/css/main.css'
import 'leaflet/dist/leaflet.css'
import './assets/css/map-tooltips.css'
import { initEventSocket } from './services/eventSocket'

// Bloqueia o menu de contexto padrao do navegador para comportamento nativo de VMS
if (typeof window !== 'undefined') {
  window.addEventListener('contextmenu', (e) => {
    e.preventDefault()
  })
}

// Inicializa conexao multiplexada CloudEvents WebSocket
initEventSocket()

createApp(App).mount('#app')
