<template>
  <el-container class="layout">
    <!-- 侧边栏 -->
    <el-aside :width="collapsed ? '64px' : '220px'" class="aside">
      <div class="logo">
        <el-icon class="logo-icon"><Connection /></el-icon>
        <span v-if="!collapsed" class="logo-text">CrashPanel</span>
      </div>

      <el-menu
        :default-active="currentPath"
        :collapse="collapsed"
        :collapse-transition="false"
        router
        class="menu"
      >
        <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
          <el-icon><component :is="item.icon" /></el-icon>
          <template #title>{{ item.title }}</template>
        </el-menu-item>
      </el-menu>

      <div class="collapse-btn" @click="collapsed = !collapsed">
        <el-icon><Fold v-if="!collapsed" /><Expand v-else /></el-icon>
      </div>
    </el-aside>

    <el-container class="main-container">
      <!-- 顶栏 -->
      <el-header class="header">
        <div class="header-left">
          <span class="page-title">{{ currentTitle }}</span>
        </div>
        <div class="header-right">
          <!-- 核心状态指示 -->
          <div class="core-badge" :class="coreRunning ? 'running' : 'stopped'">
            <span class="dot" />
            {{ coreRunning ? '运行中' : '已停止' }}
          </div>
          <el-dropdown @command="handleCommand">
            <div class="user-btn">
              <el-icon><UserFilled /></el-icon>
              <span>管理员</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 内容区 -->
      <el-main class="main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { coreApi } from '@/api'

const route = useRoute()
const auth = useAuthStore()
const collapsed = ref(false)
const coreRunning = ref(false)

const menuItems = [
  { path: '/dashboard',    title: '仪表盘',  icon: 'Odometer' },
  { path: '/subscription', title: '订阅管理', icon: 'Link' },
  { path: '/dns',          title: 'DNS 配置', icon: 'Discover' },
  { path: '/firewall',     title: '防火墙',   icon: 'Shield' },
  { path: '/tasks',        title: '计划任务', icon: 'Clock' },
  { path: '/log',          title: '运行日志', icon: 'Document' },
  { path: '/settings',     title: '全局设置', icon: 'Setting' },
]

const currentPath = computed(() => route.path)
const currentTitle = computed(() =>
  menuItems.find(m => m.path === route.path)?.title || 'CrashPanel'
)

// 轮询核心状态
let timer: ReturnType<typeof setInterval>
async function pollStatus() {
  try {
    const res = await coreApi.status()
    coreRunning.value = res.data.running
  } catch {}
}

onMounted(() => {
  pollStatus()
  timer = setInterval(pollStatus, 5000)
})
onUnmounted(() => clearInterval(timer))

function handleCommand(cmd: string) {
  if (cmd === 'logout') auth.logout()
}
</script>

<style scoped>
.layout { height: 100vh; overflow: hidden; }

.aside {
  background: #1a1d23;
  display: flex;
  flex-direction: column;
  transition: width 0.2s;
  overflow: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 20px;
  border-bottom: 1px solid #2d3139;
  flex-shrink: 0;
}
.logo-icon { font-size: 22px; color: #409eff; }
.logo-text { font-size: 16px; font-weight: 700; color: #fff; white-space: nowrap; }

.menu {
  flex: 1;
  border-right: none;
  background: transparent;
  --el-menu-text-color: #a3a9b5;
  --el-menu-active-color: #fff;
  --el-menu-hover-bg-color: #2d3139;
  --el-menu-bg-color: transparent;
  --el-menu-item-height: 48px;
  overflow-y: auto;
  overflow-x: hidden;
}
:deep(.el-menu-item.is-active) {
  background: #409eff22 !important;
  border-right: 3px solid #409eff;
}

.collapse-btn {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #a3a9b5;
  cursor: pointer;
  border-top: 1px solid #2d3139;
  transition: color 0.2s;
  flex-shrink: 0;
}
.collapse-btn:hover { color: #fff; }

.header {
  height: 60px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0,0,0,.06);
}
.page-title { font-size: 16px; font-weight: 600; color: #1d2129; }

.header-right { display: flex; align-items: center; gap: 16px; }

.core-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
}
.core-badge.running  { background: #f0faf0; color: #27ae60; }
.core-badge.stopped  { background: #fef0f0; color: #e74c3c; }
.dot {
  width: 7px; height: 7px;
  border-radius: 50%;
  background: currentColor;
  animation: pulse 2s infinite;
}
.core-badge.stopped .dot { animation: none; }
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: .4; }
}

.user-btn {
  display: flex; align-items: center; gap: 6px;
  cursor: pointer; color: #606266; font-size: 14px;
}
.user-btn:hover { color: #409eff; }

.main-container { overflow: hidden; }
.main {
  background: #f0f2f5;
  padding: 24px;
  overflow-y: auto;
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.15s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
