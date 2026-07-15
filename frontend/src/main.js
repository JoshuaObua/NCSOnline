import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router/index.js'
import clickOnce from './directives/clickOnce.js'
import 'icofont/dist/icofont.min.css'
import './style.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.directive('click-once', clickOnce)
app.mount('#app')
