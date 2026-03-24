<template>
  <div>
    <!-- 操作按钮 -->
    <el-card class="mb-4">
      <template #header><span class="card-title">防火墙控制</span></template>
      <div class="control-row">
        <el-button type="primary" :loading="applying" @click="doApply">
          下发规则
        </el-button>
        <el-button type="danger" plain :loading="flushing" @click="doFlush">
          清空规则
        </el-button>
        <el-alert
          type="info"
          :closable="false"
          style="flex:1; margin:0"
          title="防火墙模式（iptables/nftables）和代理模式（Redir/Tproxy/Tun）在全局设置中配置"
          show-icon
        />
      </div>
    </el-card>

    <!-- MAC/IP 过滤规则 -->
    <el-card>
      <template #header>
        <div class="card-header-row">
          <div>
            <span class="card-title">MAC / IP 过滤规则</span>
            <el-tag size="small" style="margin-left:10px">{{ filterMode === 'whitelist' ? '白名单' : '黑名单' }}模式</el-tag>
          </div>
          <div style="display:flex;gap:8px">
            <el-select v-model="filterMode" size="small" style="width:110px">
              <el-option label="黑名单" value="blacklist" />
              <el-option label="白名单" value="whitelist" />
            </el-select>
            <el-button type="primary" size="small" :icon="Plus" @click="openDialog()">添加</el-button>
          </div>
        </div>
      </template>

      <el-alert
        :type="filterMode === 'whitelist' ? 'warning' : 'info'"
        :closable="false"
        class="mb-4"
        :title="filterMode === 'whitelist'
          ? '白名单模式：仅列表中的设备流量经过代理，其他设备直连'
          : '黑名单模式：列表中的设备流量直连，其他设备经过代理'"
      />

      <el-table :data="rules" v-loading="loading" stripe>
        <el-table-column prop="type" label="类型" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.type === 'mac' ? 'primary' : 'warning'" size="small">
              {{ row.type.toUpperCase() }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="value" label="地址" min-width="160" />
        <el-table-column prop="remark" label="备注" min-width="120" />
        <el-table-column label="操作" width="80" align="center">
          <template #default="{ row }">
            <el-button text type="danger" size="small" @click="doDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加规则对话框 -->
    <el-dialog v-model="dialogVisible" title="添加过滤规则" width="420px" destroy-on-close>
      <el-form :model="form" label-width="80px">
        <el-form-item label="类型">
          <el-radio-group v-model="form.type">
            <el-radio value="mac">MAC 地址</el-radio>
            <el-radio value="ip">IP 地址</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="地址" required>
          <el-input
            v-model="form.value"
            :placeholder="form.type === 'mac' ? 'AA:BB:CC:DD:EE:FF' : '192.168.1.100 或 192.168.1.0/24'"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" placeholder="可选" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="doSave">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { firewallApi } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const rules = ref<any[]>([])
const loading = ref(false)
const applying = ref(false)
const flushing = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const filterMode = ref('blacklist')
const form = ref({ type: 'mac', value: '', remark: '', filter_mode: 'blacklist' })

async function load() {
  loading.value = true
  try {
    const res = await firewallApi.rules()
    rules.value = res.data
    if (res.data.length > 0) filterMode.value = res.data[0].filter_mode
  } finally {
    loading.value = false
  }
}

function openDialog() {
  form.value = { type: 'mac', value: '', remark: '', filter_mode: filterMode.value }
  dialogVisible.value = true
}

async function doSave() {
  if (!form.value.value) return ElMessage.warning('请填写地址')
  saving.value = true
  try {
    await firewallApi.addRule({ ...form.value, filter_mode: filterMode.value })
    ElMessage.success('已添加')
    dialogVisible.value = false
    load()
  } finally {
    saving.value = false
  }
}

async function doDelete(row: any) {
  await ElMessageBox.confirm(`确定删除 ${row.value}？`, '确认', { type: 'warning' })
  await firewallApi.deleteRule(row.id)
  ElMessage.success('已删除')
  load()
}

async function doApply() {
  applying.value = true
  try {
    await firewallApi.apply()
    ElMessage.success('防火墙规则已下发')
  } finally {
    applying.value = false
  }
}

async function doFlush() {
  await ElMessageBox.confirm('确定清空所有防火墙规则？', '确认', { type: 'warning' })
  flushing.value = true
  try {
    await firewallApi.flush()
    ElMessage.success('防火墙规则已清空')
  } finally {
    flushing.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
.card-header-row { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 15px; font-weight: 600; }
.control-row { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
</style>
