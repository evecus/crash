import { defineStore } from 'pinia'
import { ref } from 'vue'
import { request } from '@/api/request'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')

  async function login(password: string) {
    const res = await request.post('/auth/login', { password })
    token.value = res.data.token
    localStorage.setItem('token', token.value)
    router.push('/dashboard')
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('token')
    router.push('/login')
  }

  return { token, login, logout }
})
