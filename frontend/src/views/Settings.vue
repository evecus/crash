<template>
  <el-card>
    <template #header>
      <div class="card-header-row">
        <span class="card-title">全局设置</span>
        <el-button type="primary" :loading="saving" @click="doSave">保存设置</el-button>
      </div>
    </template>

    <el-form :model="form" label-width="140px" v-loading="loading">

      <el-divider content-position="left">核心设置</el-divider>

      <el-form-item label="核心类型">
        <el-radio-group v-model="form.core_type">
          <el-radio-button value="meta">Clash Meta (mihomo)</el-radio-button>
          <el-radio-button value="singbox">sing-box</el-radio-button>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="核心二进制路径">
        <el-input v-model="form.core_path" placeholder="/usr/local/bin/mihomo" />
        <div class="hint">核心可执行文件的完整路径，需用户自行放置</div>
      </el-form-item>

      <el-form-item label="核心工作目录">
        <el-input v-model="form.core_work_dir" placeholder="/etc/crashpanel/core" />
        <div class="hint">配置文件、规则文件等均放置于此目录</div>
      </el-form-item>

      <el-divider content-position="left">端口设置</el-divider>

      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="混合端口">
            <el-input-number v-model="form.mix_port" :min="1024" :max="65535" style="width:100%" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Redir 端口">
            <el-input-number v-model="form.redir_port" :min="1024" :max="65535" style="width:100%" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="TProxy 端口">
            <el-input-number v-model="form.tproxy_port" :min="1024" :max="65535" style="width:100%" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Dashboard 端口">
            <el-input-number v-model="form.dashboard_port" :min="1024" :max="65535" style="width:100%" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="DNS 端口">
            <el-input-number v-model="form.dns_port" :min="1024" :max="65535" style="width:100%" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-divider content-position="left">代理模式</el-divider>

      <el-form-item label="透明代理模式">
        <el-select v-model="form.redir_mod" style="width:200px">
          <el-option label="Redir（TCP 重定向）" value="Redir" />
          <el-option label="Tproxy（TCP+UDP）" value="Tproxy" />
          <el-option label="Tun（虚拟网卡）" value="Tun" />
          <el-option label="Mix（混合）" value="Mix" />
        </el-select>
      </el-form-item>

      <el-form-item label="DNS 模式">
        <el-select v-model="form.dns_mod" style="width:200px">
          <el-option label="redir-host（真实 IP）" value="redir-host" />
          <el-option label="fake-ip（虚假 IP）" value="fake-ip" />
          <el-option label="mix（混合模式）" value="mix" />
        </el-select>
      </el-form-item>

      <el-divider content-position="left">防火墙设置</el-divider>

      <el-form-item label="防火墙实现">
        <el-radio-group v-model="form.firewall_mod">
          <el-radio-button value="nftables">nftables（推荐）</el-radio-button>
          <el-radio-button value="iptables">iptables</el-radio-button>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="代理范围">
        <el-select v-model="form.firewall_area" style="width:260px">
          <el-option :value="1" label="全局代理（所有流量）" />
          <el-option :value="2" label="绕过 CN IP（推荐）" />
          <el-option :value="3" label="仅局域网设备" />
          <el-option :value="4" label="纯净模式（不拦截）" />
        </el-select>
      </el-form-item>

      <el-divider content-position="left">功能开关</el-divider>

      <el-form-item label="绕过 CN IP">
        <el-switch v-model="form.cn_ip_route" />
        <span class="switch-hint">开启后国内 IP 直连，需 GeoIP 数据库</span>
      </el-form-item>

      <el-form-item label="仅代理常用端口">
        <el-switch v-model="form.common_ports" />
        <span class="switch-hint">仅拦截 80/443 等常用端口流量</span>
      </el-form-item>

    </el-form>

    <el-alert
      type="warning"
      :closable="false"
      title="修改设置后需重新生成配置文件并重启核心才能生效"
      style="margin-top:16px"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { settingsApi } from '@/api'
import { ElMessage } from 'element-plus'

const form = ref<any>({
  core_type: 'meta', core_path: '', core_work_dir: '',
  mix_port: 7890, redir_port: 7892, tproxy_port: 7893,
  dashboard_port: 9999, dns_port: 1053,
  redir_mod: 'Redir', dns_mod: 'redir-host',
  firewall_mod: 'nftables', firewall_area: 2,
  cn_ip_route: true, common_ports: false
})
const loading = ref(false)
const saving = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await settingsApi.get()
    form.value = res.data
  } finally {
    loading.value = false
  }
}

async function doSave() {
  saving.value = true
  try {
    await settingsApi.update(form.value)
    ElMessage.success('设置已保存')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.card-header-row { display: flex; align-items: center; justify-content: space-between; }
.card-title { font-size: 15px; font-weight: 600; }
.hint { font-size: 12px; color: #8a8f98; margin-top: 4px; }
.switch-hint { font-size: 13px; color: #8a8f98; margin-left: 10px; }
</style>
