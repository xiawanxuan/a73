import http from './http'
import type { PointCloud, PageResult } from '../types'

export function list(page = 1, pageSize = 20): Promise<PageResult<PointCloud>> {
  return http.get('/v1/point-clouds', { params: { page, pageSize } })
}

export function getById(id: string): Promise<PointCloud> {
  return http.get(`/v1/point-clouds/${id}`)
}

export function create(data: Partial<PointCloud>): Promise<PointCloud> {
  return http.post('/v1/point-clouds', data)
}

export function update(id: string, data: Partial<PointCloud>): Promise<PointCloud> {
  return http.put(`/v1/point-clouds/${id}`, data)
}

export function remove(id: string): Promise<void> {
  return http.delete(`/v1/point-clouds/${id}`)
}
