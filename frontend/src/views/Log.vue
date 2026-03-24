<template>
  <el-card>
    <template #header>
      <div class="card-header-row">
        <span class="card-title">运行日志</span>
        <div style="display:flex;gap:8px;align-items:center">
          <el-select v-model="lines" size="small" style="width:100px" @change="load">
            <el-option :value="100" label="最近 100 行" />
            <el-option :value="300" label="最近 300 行" />
            <el-option :value="500" label="最近 500 行" />
          </el-select>
          <el-switch v-model="autoRefresh" active-text="自动刷新" size="small" />
          <el-button size="small" :icon="Refresh" @click="load">刷新</el-button>
          <el-button size="small" :icon="Delete" @click="logData = []">清屏</el-button>
        </div>
      </div>
    </template>

    <div class="log-box" ref="logBox">
      <div v-if="logData.length === 0" class="log-empty">暂无日志，请先启动核心</div>
      <div
        v-for="(line, i) in logData"
        :key="i"
        class="log-line"
        :class="lineClass(line)"
      >{{ line }}</div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { Refresh, Delete } from '@element-plus/icons-vue'
import { coreApi } from '@/api'

const logData = ref<string[]>([])
const lines = ref(200)
const autoRefresh = ref(true)
const logBox = ref<HTMLElement>()

async function load() {
  try {
    const res = await coreApi.log(lines.value)
    logData.value = res.data.lines || []
    await nextTick()
    scrollToBottom()
  } catch {}
}

function scrollToBottom() {
  if (logBox.value) {
    logBox.value.scrollTop = logBox.value.scrollHeight
  }
}

function lineClass(line: string) {
  const l = line.toLowerCase()
  if (l.includes('error') || l.includes('fatal')) return 'log-error'
  if (l.includes('warn')) return 'log-warn'
  if (l.includes('info')) return 'log-info'
  return ''
}

let timer: ReturnType<typeof setInterval>
watch(autoRefresh, (v) => {
  clearInterval(timer)
  if (v) timer = setInterval(load, 3000)
})

onMounted(() => {
  load()
  timer = setInterval(load, 3000)
})
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.card-header-row { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 15px; font-weight: 600; }

.log-box {
  background: #1a1d23;
  border-radius: 8px;
  padding: 16px;
  height: calc(100vh - 220px);
  min-height: 400px;
  overflow-y: auto;
  font-family: 'Menlo', 'Consolas', 'SF Mono', monospace;
  font-size: 12.5px;
  line-height: 1.7;
}

.log-line       { color: #c8d3e0; white-space: pre-wrap; word-break: break-all; }
.log-error      { color: #ff6b6b; }
.log-warn       { color: #ffa94d; }
.log-info       { color: #69db7c; }
.log-empty      { color: #555; text-align: center; padding: 40px 0; }
</style>
