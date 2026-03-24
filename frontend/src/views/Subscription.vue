<template>
  <div>
    <el-card>
      <template #header>
        <div class="card-header-row">
          <span class="card-title">订阅管理</span>
          <el-button type="primary" :icon="Plus" @click="openDialog()">添加订阅</el-button>
        </div>
      </template>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column prop="target" label="类型" width="90">
          <template #default="{ row }">
            <el-tag size="small">{{ row.target || 'clash' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="node_count" label="节点数" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'ok' ? 'success' : row.status === 'error' ? 'danger' : 'info'"
              size="small"
            >{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="165">
          <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column prop="auto_update" label="自动更新" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="row.auto_update ? 'success' : 'info'" size="small">
              {{ row.auto_update ? `${row.interval}h` : '关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="190" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" size="small" :loading="refreshing[row.id]" @click="doRefresh(row)">更新</el-button>
            <el-button text size="small" @click="openDialog(row)">编辑</el-button>
            <el-button text type="danger" size="small" @click="doDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑订阅' : '添加订阅'"
      width="560px"
      destroy-on-close
    >
      <el-form :model="form" label-width="110px" class="sub-form">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="如：我的机场" />
        </el-form-item>
        <el-form-item label="订阅链接" required>
          <el-input v-model="form.url" placeholder="https://..." type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="目标格式">
          <el-select v-model="form.target" style="width:100%">
            <el-option label="Clash (Meta/mihomo)" value="clash" />
            <el-option label="sing-box" value="singbox" />
          </el-select>
        </el-form-item>
        <el-form-item label="SubConverter">
          <el-input v-model="form.sub_converter_url" placeholder="https://sub.example.com（留空则直接使用原链接）" />
        </el-form-item>
        <el-form-item label="User-Agent">
          <el-input v-model="form.user_agent" placeholder="留空使用默认值" />
        </el-form-item>
        <el-form-item label="节点过滤">
          <el-row :gutter="8">
            <el-col :span="12">
              <el-input v-model="form.include" placeholder="包含关键词（正则）" />
            </el-col>
            <el-col :span="12">
              <el-input v-model="form.exclude" placeholder="排除关键词（正则）" />
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item label="规则配置 URL">
          <el-input v-model="form.config_url" placeholder="远程规则配置链接（可选）" />
        </el-form-item>
        <el-form-item label="自动更新">
          <el-switch v-model="form.auto_update" />
          <el-input-number
            v-if="form.auto_update"
            v-model="form.interval"
            :min="1" :max="168"
            style="margin-left:12px; width:120px"
          />
          <span v-if="form.auto_update" style="margin-left:8px;color:#8a8f98;font-size:13px">小时</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="doSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { subApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const list = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const refreshing = reactive<Record<number, boolean>>({})

const defaultForm = () => ({
  id: 0, name: '', url: '', target: 'clash',
  sub_converter_url: '', user_agent: '',
  include: '', exclude: '', config_url: '',
  auto_update: false, interval: 12
})
const form = ref(defaultForm())

async function load() {
  loading.value = true
  try {
    const res = await subApi.list()
    list.value = res.data
  } finally {
    loading.value = false
  }
}

function openDialog(row?: any) {
  form.value = row ? { ...defaultForm(), ...row } : defaultForm()
  dialogVisible.value = true
}

async function doSave() {
  saving.value = true
  try {
    if (form.value.id) {
      await subApi.update(form.value.id, form.value)
    } else {
      await subApi.create(form.value)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

async function doRefresh(row: any) {
  refreshing[row.id] = true
  try {
    const res = await subApi.refresh(row.id)
    ElMessage.success(`更新成功，共 ${res.data.node_count} 个节点`)
    load()
  } finally {
    refreshing[row.id] = false
  }
}

async function doDelete(row: any) {
  await ElMessageBox.confirm(`确定删除订阅「${row.name}」？`, '确认', { type: 'warning' })
  await subApi.remove(row.id)
  ElMessage.success('已删除')
  load()
}

function statusLabel(s: string) {
  return { ok: '正常', error: '错误', pending: '待更新' }[s] || s
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
.sub-form { padding: 0 8px; }
</style>
