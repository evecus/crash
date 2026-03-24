import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      component: () => import('@/views/Login.vue'),
      meta: { public: true }
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      redirect: '/dashboard',
      children: [
        { path: 'dashboard',    component: () => import('@/views/Dashboard.vue'),    meta: { title: '仪表盘' } },
        { path: 'subscription', component: () => import('@/views/Subscription.vue'), meta: { title: '订阅管理' } },
        { path: 'dns',          component: () => import('@/views/DNS.vue'),           meta: { title: 'DNS 配置' } },
        { path: 'firewall',     component: () => import('@/views/Firewall.vue'),      meta: { title: '防火墙' } },
        { path: 'tasks',        component: () => import('@/views/Tasks.vue'),          meta: { title: '计划任务' } },
        { path: 'settings',     component: () => import('@/views/Settings.vue'),       meta: { title: '全局设置' } },
        { path: 'log',          component: () => import('@/views/Log.vue'),            meta: { title: '运行日志' } },
      ]
    },
    { path: '/:pathMatch(.*)*', redirect: '/' }
  ]
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.token) {
    return '/login'
  }
})

export default router
