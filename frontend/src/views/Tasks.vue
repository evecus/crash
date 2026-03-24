<template>
  <el-card>
    <template #header>
      <div class="card-header-row">
        <span class="card-title">计划任务</span>
        <el-button type="primary" :icon="Plus" @click="openDialog()">添加任务</el-button>
      </div>
    </template>

    <el-table :data="tasks" v-loading="loading" stripe>
      <el-table-column prop="name" label="任务名称" min-width="130" />
      <el-table-column prop="cron" label="Cron 表达式" width="140">
        <template #default="{ row }">
          <el-tag size="small" type="info">{{ row.cron }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="command" label="命令" min-width="200" show-overflow-tooltip />
      <el-table-column prop="enabled" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-switch
            v-model="row.enabled"
            size="small"
            @change="toggleEnabled(row)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="last_run" label="上次运行" width="165">
        <template #default="{ row }">{{ formatTime(row.last_run) }}</template>
      </el-table-column>
      <el-table-column prop="last_code" label="结果" width="70" align="center">
        <template #default="{ row }">
          <el-tag
            v-if="row.last_run && !row.last_run.startsWith('0001')"
            :type="row.last_code === 0 ? 'success' : 'danger'"
            size="small"
          >{{ row.last_code === 0 ? '成功' : '失败' }}</el-tag>
          <span v-else style="color:#c0c4cc">—</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="190" fixed="right">
        <template #default="{ row }">
          <el-button text type="primary" size="small" :loading="running[row.id]" @click="doRun(row)">立即执行</el-button>
          <el-button text size="small" @click="openDialog(row)">编辑</el-button>
          <el-button text type="danger" size="small" @click="doDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 最近日志展示 -->
    <el-dialog v-model="logVisible" title="执行日志" width="600px">
      <pre class="task-log">{{ currentLog || '无输出' }}</pre>
    </el-dialog>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑任务' : '添加任务'"
      width="520px"
      destroy-on-close
    >
      <el-form :model="form" label-width="110px">
        <el-form-item label="任务名称" required>
          <el-input v-model="form.name" placeholder="订阅自动更新" />
        </el-form-item>
        <el-form-item label="Cron 表达式" required>
          <el-input v-model="form.cron" placeholder="0 3 * * *" />
          <div class="hint">标准 5 字段 cron，如 0 3 * * * 表示每天凌晨 3 点</div>
        </el-form-item>
        <el-form-item label="执行命令" required>
          <el-input
            v-model="form.command"
            type="textarea"
            :rows="3"
            placeholder="curl -s https://example.com/update.sh | sh"
          />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="doSave">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { taskApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const tasks = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const logVisible = ref(false)
const currentLog = ref('')
const running = reactive<Record<number, boolean>>({})

const defaultForm = () => ({ id: 0, name: '', cron: '', command: '', enabled: true })
const form = ref(defaultForm())

async function load() {
  loading.value = true
  try {
    const res = await taskApi.list()
    tasks.value = res.data
  } finally {
    loading.value = false
  }
}

function openDialog(row?: any) {
  form.value = row ? { ...defaultForm(), ...row } : defaultForm()
  dialogVisible.value = true
}

async function doSave() {
  if (!form.value.name || !form.value.cron || !form.value.command) {
    return ElMessage.warning('请填写完整信息')
  }
  saving.value = true
  try {
    if (form.value.id) {
      await taskApi.update(form.value.id, form.value)
    } else {
      await taskApi.create(form.value)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

async function toggleEnabled(row: any) {
  await taskApi.update(row.id, row)
}

async function doRun(row: any) {
  running[row.id] = true
  try {
    const res = await taskApi.run(row.id)
    currentLog.value = res.data.log
    logVisible.value = true
    ElMessage.success(`执行完成，退出码 ${res.data.exit_code}`)
    load()
  } finally {
    running[row.id] = false
  }
}

async function doDelete(row: any) {
  await ElMessageBox.confirm(`确定删除任务「${row.name}」？`, '确认', { type: 'warning' })
  await taskApi.remove(row.id)
  ElMessage.success('已删除')
  load()
}

function formatTime(t: string) {
  if (!t || t.startsWith('0001')) return '—'
  return new Date(t).toLocaleString('zh-CN', { hour12: false })
}

onMounted(load)
</script>

<style scoped>
.card-header-row { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 15px; font-weight: 600; }
.hint { font-size: 12px; color: #8a8f98; margin-top: 4px; }
.task-log {
  background: #1a1d23;
  color: #a8d8a8;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 12px;
  font-family: 'Menlo', 'Consolas', monospace;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 400px;
  overflow-y: auto;
}
</style>
