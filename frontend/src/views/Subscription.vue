<template>
  <div>
    <!-- 操作栏 -->
    <el-card class="mb-4">
      <div class="action-row">
        <el-button type="primary" :icon="Plus" @click="openDialog()">添加提供者</el-button>
        <el-button :icon="Refresh" :loading="refreshingAll" @click="doRefreshAll">更新全部</el-button>
        <el-button type="success" :icon="Document" :loading="generating" @click="doGenerate">
          生成配置文件
        </el-button>
        <el-tooltip content="更新所有提供者后，用选中的规则模板生成完整 config.yaml" placement="top">
          <el-icon class="help-icon"><QuestionFilled /></el-icon>
        </el-tooltip>
      </div>
    </el-card>

    <!-- 提供者列表 -->
    <el-card class="mb-4">
      <template #header>
        <div class="card-header-row">
          <span class="card-title">提供者列表</span>
          <el-tag type="info" size="small">{{ list.length }} 个</el-tag>
        </div>
      </template>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column label="类型" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="typeColor(row.link_type)" size="small">{{ typeLabel(row.link_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column label="链接/路径" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="link-text">{{ row.url || row.file_path }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="node_count" label="节点" width="70" align="center" />
        <el-table-column label="状态" width="85" align="center">
          <template #default="{ row }">
            <el-tag :type="statusColor(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="更新时间" width="160">
          <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.link_type !== 'file'"
              text type="primary" size="small"
              :loading="refreshing[row.id]"
              @click="doRefresh(row)"
            >更新</el-button>
            <el-button text size="small" @click="openDialog(row)">编辑</el-button>
            <el-button text type="danger" size="small" @click="doDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 规则模板 -->
    <el-card>
      <template #header>
        <div class="card-header-row">
          <span class="card-title">规则模板</span>
          <el-button size="small" :icon="Plus" @click="tmplDialogVisible = true">自定义模板</el-button>
        </div>
      </template>

      <el-alert type="info" :closable="false" class="mb-4"
        title="生成配置文件时，会使用当前选中（默认）的规则模板来构建 proxy-groups 和 rules" />

      <el-table :data="templates" stripe>
        <el-table-column label="" width="50" align="center">
          <template #default="{ row }">
            <el-icon v-if="row.is_default" color="#67c23a"><Select /></el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="模板名称" min-width="220" />
        <el-table-column prop="url" label="URL" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="link-text">{{ row.url || row.local_path || '本地文件' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center">
          <template #default="{ row }">
            <el-button
              v-if="!row.is_default"
              text type="primary" size="small"
              @click="setDefaultTemplate(row)"
            >设为默认</el-button>
            <el-tag v-else type="success" size="small">当前使用</el-tag>
            <el-button
              text type="danger" size="small"
              @click="deleteTemplate(row)"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加提供者对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑提供者' : '添加提供者'"
      width="580px"
      destroy-on-close
    >
      <el-form :model="form" label-width="110px">
        <!-- 类型选择 -->
        <el-form-item label="提供者类型">
          <el-radio-group v-model="form.link_type" @change="onTypeChange">
            <el-radio-button value="url">订阅链接</el-radio-button>
            <el-radio-button value="uri">分享链接</el-radio-button>
            <el-radio-button value="file">本地文件</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="名称" required>
          <el-input v-model="form.name" placeholder="如：我的机场" />
        </el-form-item>

        <!-- 订阅链接 -->
        <template v-if="form.link_type === 'url'">
          <el-form-item label="订阅链接" required>
            <el-input v-model="form.url" placeholder="https://..." type="textarea" :rows="2" />
          </el-form-item>
          <el-form-item label="SubConverter">
            <el-input v-model="form.sub_converter_url" placeholder="留空则直接使用原链接" />
          </el-form-item>
          <el-form-item label="目标格式">
            <el-select v-model="form.target" style="width:100%">
              <el-option label="Clash (Meta/mihomo)" value="clash" />
              <el-option label="sing-box" value="singbox" />
            </el-select>
          </el-form-item>
          <el-form-item label="User-Agent">
            <el-input v-model="form.user_agent" placeholder="留空使用默认值 clash.meta" />
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
          <el-form-item label="健康检查">
            <el-input-number v-model="form.health_interval" :min="1" :max="60" style="width:120px" />
            <span class="unit-label">分钟</span>
          </el-form-item>
          <el-form-item label="自动更新">
            <el-switch v-model="form.auto_update" />
            <el-input-number
              v-if="form.auto_update"
              v-model="form.interval"
              :min="1" :max="168"
              style="margin-left:12px;width:110px"
            />
            <span v-if="form.auto_update" class="unit-label">小时</span>
          </el-form-item>
        </template>

        <!-- 分享链接 URI -->
        <template v-if="form.link_type === 'uri'">
          <el-form-item label="分享链接" required>
            <el-input
              v-model="form.url"
              type="textarea" :rows="3"
              placeholder="vmess://... 或 ss://... 或 trojan://... 或 vless://..."
            />
          </el-form-item>
          <el-alert type="info" :closable="false" class="mb-2"
            title="支持 vmess / ss / trojan / vless / hysteria2 等格式的单节点分享链接" />
        </template>

        <!-- 本地文件 -->
        <template v-if="form.link_type === 'file'">
          <el-form-item label="选择文件">
            <el-upload
              :auto-upload="false"
              :on-change="onFileChange"
              :show-file-list="true"
              accept=".yaml,.yml,.json"
              :limit="1"
            >
              <el-button :icon="Upload">选择 yaml/json 文件</el-button>
            </el-upload>
          </el-form-item>
          <el-form-item label="文件路径" v-if="form.file_path">
            <el-input v-model="form.file_path" readonly />
          </el-form-item>
          <el-alert type="info" :closable="false"
            title="文件将上传到服务器的 providers 目录，之后配置生成时会直接引用" />
        </template>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="doSave">保存</el-button>
      </template>
    </el-dialog>

    <!-- 自定义模板对话框 -->
    <el-dialog v-model="tmplDialogVisible" title="添加自定义规则模板" width="500px" destroy-on-close>
      <el-form :model="tmplForm" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="tmplForm.name" placeholder="我的规则模板" />
        </el-form-item>
        <el-form-item label="URL" required>
          <el-input v-model="tmplForm.url" placeholder="https://... (.ini 格式)" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="tmplForm.description" placeholder="可选" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="tmplDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="doAddTemplate">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Refresh, Document, Upload, QuestionFilled, Select } from '@element-plus/icons-vue'
import { subApi, templateApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const list = ref<any[]>([])
const templates = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)
const generating = ref(false)
const refreshingAll = ref(false)
const dialogVisible = ref(false)
const tmplDialogVisible = ref(false)
const refreshing = reactive<Record<number, boolean>>({})
let pendingFile: File | null = null

const defaultForm = () => ({
  id: 0, name: '', link_type: 'url', url: '', file_path: '',
  target: 'clash', sub_converter_url: '', user_agent: '',
  include: '', exclude: '', config_url: '',
  auto_update: false, interval: 12, health_interval: 3,
})
const form = ref(defaultForm())
const tmplForm = ref({ name: '', url: '', description: '' })

async function load() {
  loading.value = true
  try {
    const [s, t] = await Promise.all([subApi.list(), templateApi.list()])
    list.value = s.data
    templates.value = t.data
  } finally {
    loading.value = false
  }
}

function openDialog(row?: any) {
  pendingFile = null
  form.value = row ? { ...defaultForm(), ...row } : defaultForm()
  dialogVisible.value = true
}

function onTypeChange() {
  form.value.url = ''
  form.value.file_path = ''
}

function onFileChange(file: any) {
  pendingFile = file.raw
}

async function doSave() {
  if (!form.value.name) return ElMessage.warning('请填写名称')

  saving.value = true
  try {
    // 本地文件：先上传文件
    if (form.value.link_type === 'file' && pendingFile) {
      const res = await subApi.upload(pendingFile)
      form.value.file_path = res.data.file_path
    }

    if (form.value.link_type === 'file' && !form.value.file_path) {
      ElMessage.warning('请选择文件')
      return
    }
    if (form.value.link_type !== 'file' && !form.value.url) {
      ElMessage.warning('请填写链接')
      return
    }

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
    ElMessage.success(res.data.message)
    load()
  } finally {
    refreshing[row.id] = false
  }
}

async function doRefreshAll() {
  refreshingAll.value = true
  try {
    const res = await subApi.refreshAll()
    const ok = res.data.results.filter((r: any) => !r.error).length
    const fail = res.data.results.length - ok
    ElMessage.success(`更新完成：${ok} 成功${fail > 0 ? `，${fail} 失败` : ''}`)
    load()
  } finally {
    refreshingAll.value = false
  }
}

async function doGenerate() {
  generating.value = true
  try {
    const res = await subApi.generateConfig()
    ElMessage.success(res.data.message)
  } finally {
    generating.value = false
  }
}

async function doDelete(row: any) {
  await ElMessageBox.confirm(`确定删除「${row.name}」？`, '确认', { type: 'warning' })
  await subApi.remove(row.id)
  ElMessage.success('已删除')
  load()
}

async function setDefaultTemplate(row: any) {
  await templateApi.setDefault(row.id)
  ElMessage.success(`已将「${row.name}」设为默认模板`)
  load()
}

async function deleteTemplate(row: any) {
  await ElMessageBox.confirm(`确定删除模板「${row.name}」？`, '确认', { type: 'warning' })
  await templateApi.remove(row.id)
  load()
}

async function doAddTemplate() {
  if (!tmplForm.value.name || !tmplForm.value.url) return ElMessage.warning('请填写名称和 URL')
  await templateApi.create(tmplForm.value)
  ElMessage.success('添加成功')
  tmplDialogVisible.value = false
  tmplForm.value = { name: '', url: '', description: '' }
  load()
}

function typeLabel(t: string) { return { url: '订阅', uri: '分享', file: '本地' }[t] || t }
function typeColor(t: string) { return { url: 'primary', uri: 'warning', file: 'success' }[t] || 'info' }
function statusLabel(s: string) { return { ok: '正常', error: '错误', pending: '待更新' }[s] || s }
function statusColor(s: string) { return { ok: 'success', error: 'danger', pending: 'info' }[s] || 'info' }
function formatTime(t: string) {
  if (!t || t.startsWith('0001')) return '—'
  return new Date(t).toLocaleString('zh-CN', { hour12: false })
}

onMounted(load)
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
.mb-2 { margin-bottom: 8px; }
.card-header-row { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 15px; font-weight: 600; }
.action-row { display: flex; align-items: center; gap: 10px; }
.help-icon { color: #8a8f98; cursor: help; font-size: 16px; }
.link-text { font-size: 12px; color: #8a8f98; font-family: monospace; }
.unit-label { margin-left: 8px; font-size: 13px; color: #8a8f98; }
</style>
