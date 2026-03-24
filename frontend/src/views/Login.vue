<template>
  <div class="login-bg">
    <div class="login-card">
      <div class="brand">
        <el-icon class="brand-icon"><Connection /></el-icon>
        <h1>CrashPanel</h1>
        <p>透明代理管理面板</p>
      </div>

      <el-form @submit.prevent="handleLogin" class="form">
        <el-form-item>
          <el-input
            v-model="password"
            type="password"
            placeholder="请输入管理密码"
            size="large"
            show-password
            @keyup.enter="handleLogin"
            :prefix-icon="Lock"
          />
        </el-form-item>
        <el-button
          type="primary"
          size="large"
          :loading="loading"
          @click="handleLogin"
          class="login-btn"
        >
          登 录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const auth = useAuthStore()
const password = ref('')
const loading = ref(false)

async function handleLogin() {
  if (!password.value) return
  loading.value = true
  try {
    await auth.login(password.value)
  } catch {
    ElMessage.error('密码错误')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-bg {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1d23 0%, #2d3139 100%);
}

.login-card {
  width: 380px;
  background: #fff;
  border-radius: 16px;
  padding: 48px 40px;
  box-shadow: 0 20px 60px rgba(0,0,0,.3);
}

.brand {
  text-align: center;
  margin-bottom: 36px;
}
.brand-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 12px;
}
.brand h1 {
  font-size: 26px;
  font-weight: 700;
  color: #1d2129;
  margin-bottom: 6px;
}
.brand p { font-size: 13px; color: #8a8f98; }

.login-btn { width: 100%; margin-top: 8px; }
</style>
