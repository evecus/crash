import axios from 'axios'
import { ElMessage } from 'element-plus'

export const request = axios.create({
  baseURL: '/api',
  timeout: 15000
})

// 请求拦截：自动附带 token
request.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截：统一错误提示
request.interceptors.response.use(
  (res) => res,
  (err) => {
    const msg = err.response?.data?.error || err.message || '请求失败'
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
      return Promise.reject(err)
    }
    ElMessage.error(msg)
    return Promise.reject(err)
  }
)
