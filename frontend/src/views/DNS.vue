<template>
  <el-card>
    <template #header>
      <div class="card-header-row">
        <span class="card-title">DNS 配置</span>
        <el-button type="primary" :loading="saving" @click="doSave">保存</el-button>
      </div>
    </template>

    <el-form :model="form" label-width="130px" v-loading="loading">
      <el-divider content-position="left">基础 DNS</el-divider>

      <el-form-item label="国内 DNS">
        <el-input
          v-model="form.nameserver"
          placeholder="223.5.5.5, 1.2.4.8"
          clearable
        />
        <div class="hint">多个 DNS 用逗号分隔，用于解析国内域名</div>
      </el-form-item>

      <el-form-item label="国外 DNS (Fallback)">
        <el-input
          v-model="form.fallback"
          placeholder="1.1.1.1, 8.8.8.8"
          clearable
        />
        <div class="hint">用于解析国外域名，支持 DoH/DoT 格式</div>
      </el-form-item>

      <el-form-item label="Fallback 过滤">
        <el-input
          v-model="form.fallback_filter"
          placeholder="geoip:cn"
          clearable
        />
        <div class="hint">命中条件的域名走 Fallback DNS，如 geoip:cn</div>
      </el-form-item>

      <el-divider content-position="left">Fake-IP 配置</el-divider>

      <el-form-item label="Fake-IP 过滤">
        <el-input
          v-model="form.fake_ip_filter"
          type="textarea"
          :rows="6"
          placeholder="*.lan&#10;*.local&#10;*.home.arpa"
        />
        <div class="hint">匹配这些规则的域名不使用 Fake-IP，每行一条或逗号分隔</div>
      </el-form-item>
    </el-form>

    <!-- DNS 模式说明 -->
    <el-alert type="info" :closable="false" style="margin-top:16px">
      <template #title>
        DNS 工作模式在<el-link @click="$router.push('/settings')" style="margin:0 4px">全局设置</el-link>中配置（redir-host / fake-ip / mix）
      </template>
    </el-alert>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { dnsApi } from '@/api'
import { ElMessage } from 'element-plus'

const form = ref({ nameserver: '', fallback: '', fallback_filter: '', fake_ip_filter: '' })
const loading = ref(false)
const saving = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await dnsApi.get()
    form.value = res.data
  } finally {
    loading.value = false
  }
}

async function doSave() {
  saving.value = true
  try {
    await dnsApi.update(form.value)
    ElMessage.success('DNS 配置已保存，重启核心后生效')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.card-header-row { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 15px; font-weight: 600; }
.hint { font-size: 12px; color: #8a8f98; margin-top: 4px; line-height: 1.4; }
</style>
