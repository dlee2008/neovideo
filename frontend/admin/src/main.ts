import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/es/components/message/style/css'
import App from './App.vue'
import './tailwind.css'

const app = createApp(App)

app.use(ElementPlus)

app.mount('#app')
