import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import i18n from './i18n'
import router from './router'
import { setupGuard } from './router/guard'
import { setupPermissionDirective } from './directives/permission'
import { configureHttp } from './utils/http'
import { useUserStore } from './stores/user'
import './styles/tokens.css'
import './styles/base.css'
import './styles/shell.css'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)

// 注入 http 与 user store 的协作（避免循环依赖）
const userStore = useUserStore(pinia)
configureHttp({
  getToken: () => userStore.accessToken,
  refresh: () => userStore.doRefresh(),
  onFail: () => { userStore.clear(); router.push('/login') },
})

app.use(i18n)
app.use(router)
app.use(ElementPlus)
setupPermissionDirective(app)
setupGuard(router)
app.mount('#app')
