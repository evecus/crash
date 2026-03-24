<template>
  <div>
    <!-- 核心控制卡片 -->
    <el-row :gutter="16" class="mb-4">
      <el-col :span="8">
        <el-card class="stat-card">
          <div class="stat-header">
            <span class="stat-label">核心状态</span>
            <el-tag :type="status.running ? 'success' : 'danger'" size="small">
              {{ status.running ? '运行中' : '已停止' }}
            </el-tag>
          </div>
          <div class="stat-value">{{ status.core_type?.toUpperCase() || '—' }}</div>
          <div class="stat-sub">
            <span v-if="status.running">PID {{ status.pid }} · 运行 {{ status.uptime }}</span>
            <span v-else>核心未启动</span>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card class="stat-card">
          <div class="stat-header"><span class="stat-label">内存使用</span></div>
          <div class="stat-value">{{ formatBytes(sysInfo.mem_used) }}</div>
          <div class="stat-sub">共 {{ formatBytes(sysInfo.mem_total) }}</div>
          <el-progress
            :percentage="memPercent"
            :show-text="false"
            :stroke-width="4"
            class="mt-2"
          />
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card class="stat-card">
          <div class="stat-header"><span class="stat-label">系统</span></div>
          <div class="stat-value">{{ sysInfo.os?.toUpperCase() || '—' }}</div>
          <div class="stat-sub">{{ sysInfo.arch }} · {{ sysInfo.go_version }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 控制按钮 -->
    <el-card class="mb-4">
      <template #header><span class="card-title">核心控制</span></template>
      <div class="control-row">
        <el-button
          type="success"
          :icon="VideoPlay"
          :loading="actioning"
          :disabled="status.running"
          @click="action('start')"
        >启动</el-button>
        <el-button
          type="danger"
          :icon="VideoPause"
          :loading="actioning"
          :disabled="!status.running"
          @click="action('stop')"
        >停止</el-button>
        <el-button
          type="warning"
          :icon="RefreshRight"
          :loading="actioning"
          @click="action('restart')"
        >重启</el-button>
        <el-button
          :icon="Refresh"
          :loading="actioning"
          @click="action('config')"
        >重新生成配置</el-button>

        <div class="spacer" />
        <el-button text :icon="Refresh" @click="loadAll">刷新</el-button>
      </div>
    </el-card>

    <!-- 最近日志 -->
    <el-card>
      <template #header>
        <div class="card-header-row">
          <span class="card-title">最近日志</span>
          <el-button text size="small" @click="$router.push('/log')">查看全部 →</el-button>
        </div>
      </template>
      <div class="log-box">
        <div v-if="logs.length === 0" class="log-empty">暂无日志</div>
        <div v-for="(line, i) in logs.slice(-30)" :key="i" class="log-line">{{ line }}</div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { VideoPlay, VideoPause, RefreshRight, Refresh } from '@element-plus/icons-vue'
import { coreApi, systemApi } from '@/api'
import { ElMessage } from 'element-plus'

const status = ref<any>({ running: false })
const sysInfo = ref<any>({})
const logs = ref<string[]>([])
const actioning = ref(false)

const memPercent = computed(() => {
  if (!sysInfo.value.mem_total) return 0
  return Math.round((sysInfo.value.mem_used / sysInfo.value.mem_total) * 100)
})

function formatBytes(b: number) {
  if (!b) return '—'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  while (b >= 1024 && i < units.length - 1) { b /= 1024; i++ }
  return `${b.toFixed(1)} ${units[i]}`
}

async function loadAll() {
  const [s, sys, lg] = await Promise.allSettled([
    coreApi.status(),
    systemApi.info(),
    coreApi.log(50)
  ])
  if (s.status === 'fulfilled') status.value = s.value.data
  if (sys.status === 'fulfilled') sysInfo.value = sys.value.data
  if (lg.status === 'fulfilled') logs.value = lg.value.data.lines || []
}

async function action(type: string) {
  actioning.value = true
  try {
    if (type === 'start')   await coreApi.start()
    if (type === 'stop')    await coreApi.stop()
    if (type === 'restart') await coreApi.restart()
    if (type === 'config') {
      // 重新生成配置后重启
      await coreApi.restart()
    }
    ElMessage.success('操作成功')
    await loadAll()
  } catch {
  } finally {
    actioning.value = false
  }
}

let timer: ReturnType<typeof setInterval>
onMounted(() => {
  loadAll()
  timer = setInterval(loadAll, 8000)
})
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
.mt-2 { margin-top: 8px; }

.stat-card { height: 130px; }
.stat-header {
  display: flex; align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}
.stat-label { font-size: 13px; color: #8a8f98; }
.stat-value { font-size: 26px; font-weight: 700; color: #1d2129; margin-bottom: 4px; }
.stat-sub   { font-size: 12px; color: #a3a9b5; }

.card-title { font-size: 15px; font-weight: 600; }
.card-header-row { display: flex; align-items: center; justify-content: space-between; }

.control-row {
  display: flex; align-items: center; gap: 10px; flex-wrap: wrap;
}
.spacer { flex: 1; }

.log-box {
  background: #1a1d23;
  border-radius: 8px;
  padding: 12px 16px;
  max-height: 280px;
  overflow-y: auto;
  font-family: 'Menlo', 'Consolas', monospace;
  font-size: 12px;
}
.log-line { color: #a8d8a8; line-height: 1.7; white-space: pre-wrap; word-break: break-all; }
.log-empty { color: #555; text-align: center; padding: 20px 0; }
</style>
