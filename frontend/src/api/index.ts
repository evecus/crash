import { request } from './request'

export const coreApi = {
  status:  ()        => request.get('/core/status'),
  start:   ()        => request.post('/core/start'),
  stop:    ()        => request.post('/core/stop'),
  restart: ()        => request.post('/core/restart'),
  log:     (n = 200) => request.get(`/core/log?lines=${n}`)
}

export const settingsApi = {
  get:    ()       => request.get('/settings'),
  update: (d: any) => request.put('/settings', d)
}

export const subApi = {
  list:           ()                   => request.get('/subscriptions'),
  create:         (d: any)             => request.post('/subscriptions', d),
  update:         (id: number, d: any) => request.put(`/subscriptions/${id}`, d),
  remove:         (id: number)         => request.delete(`/subscriptions/${id}`),
  refresh:        (id: number)         => request.post(`/subscriptions/${id}/refresh`),
  refreshAll:     ()                   => request.post('/subscriptions/refresh-all'),
  generateConfig: ()                   => request.post('/subscriptions/generate-config'),
  upload:         (file: File)         => {
    const fd = new FormData()
    fd.append('file', file)
    return request.post('/subscriptions/upload', fd, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}

export const templateApi = {
  list:       ()            => request.get('/templates'),
  create:     (d: any)      => request.post('/templates', d),
  setDefault: (id: number)  => request.post(`/templates/${id}/default`),
  remove:     (id: number)  => request.delete(`/templates/${id}`)
}

export const firewallApi = {
  rules:      ()            => request.get('/firewall/rules'),
  addRule:    (d: any)      => request.post('/firewall/rules', d),
  deleteRule: (id: number)  => request.delete(`/firewall/rules/${id}`),
  apply:      ()            => request.post('/firewall/apply'),
  flush:      ()            => request.post('/firewall/flush')
}

export const dnsApi = {
  get:    ()       => request.get('/dns'),
  update: (d: any) => request.put('/dns', d)
}

export const taskApi = {
  list:   ()                   => request.get('/tasks'),
  create: (d: any)             => request.post('/tasks', d),
  update: (id: number, d: any) => request.put(`/tasks/${id}`, d),
  remove: (id: number)         => request.delete(`/tasks/${id}`),
  run:    (id: number)         => request.post(`/tasks/${id}/run`)
}

export const systemApi = {
  info:    () => request.get('/system/info'),
  network: () => request.get('/system/network')
}
