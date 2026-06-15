import http from './http'
import type { TerrainLabel } from '../types'

export function list(): Promise<TerrainLabel[]> {
  return http.get('/v1/labels')
}

export function getById(id: string): Promise<TerrainLabel> {
  return http.get(`/v1/labels/${id}`)
}

export function create(data: Partial<TerrainLabel>): Promise<TerrainLabel> {
  return http.post('/v1/labels', data)
}

export function update(id: string, data: Partial<TerrainLabel>): Promise<TerrainLabel> {
  return http.put(`/v1/labels/${id}`, data)
}

export function remove(id: string): Promise<void> {
  return http.delete(`/v1/labels/${id}`)
}
