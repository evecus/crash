import { request } from './request'

// ── 核心 ──
export const coreApi = {
  status:  ()        => request.get('/core/status'),
  start:   ()        => request.post('/core/start'),
  stop:    ()        => request.post('/core/stop'),
  restart: ()        => request.post('/core/restart'),
  log:     (n = 200) => request.get(`/core/log?lines=${n}`)
}

// ── 设置 ──
export const settingsApi = {
  get:    ()    => request.get('/settings'),
  update: (d: any) => request.put('/settings', d)
}

// ── 订阅 ──
export const subApi = {
  list:    ()       => request.get('/subscriptions'),
  create:  (d: any) => request.post('/subscriptions', d),
  update:  (id: number, d: any) => request.put(`/subscriptions/${id}`, d),
  remove:  (id: number)         => request.delete(`/subscriptions/${id}`),
  refresh: (id: number)         => request.post(`/subscriptions/${id}/refresh`)
}

// ── 防火墙 ──
export const firewallApi = {
  rules:      ()       => request.get('/firewall/rules'),
  addRule:    (d: any) => request.post('/firewall/rules', d),
  deleteRule: (id: number) => request.delete(`/firewall/rules/${id}`),
  apply:  ()           => request.post('/firewall/apply'),
  flush:  ()           => request.post('/firewall/flush')
}

// ── DNS ──
export const dnsApi = {
  get:    ()       => request.get('/dns'),
  update: (d: any) => request.put('/dns', d)
}

// ── 计划任务 ──
export const taskApi = {
  list:   ()            => request.get('/tasks'),
  create: (d: any)      => request.post('/tasks', d),
  update: (id: number, d: any) => request.put(`/tasks/${id}`, d),
  remove: (id: number)  => request.delete(`/tasks/${id}`),
  run:    (id: number)  => request.post(`/tasks/${id}/run`)
}

// ── 系统 ──
export const systemApi = {
  info:    () => request.get('/system/info'),
  network: () => request.get('/system/network')
}
