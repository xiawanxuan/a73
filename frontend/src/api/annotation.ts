import http from './http'
import type { Annotation, AnnotationSnapshot, Point3D } from '../types'

export function listByPointCloud(pointCloudId: string): Promise<Annotation[]> {
  return http.get(`/v1/point-clouds/${pointCloudId}/annotations`)
}

export function create(
  pointCloudId: string,
  data: { labelId: string; name: string; polygon: Point3D[] }
): Promise<Annotation> {
  return http.post(`/v1/point-clouds/${pointCloudId}/annotations`, data)
}

export function getById(id: string): Promise<Annotation> {
  return http.get(`/v1/annotations/${id}`)
}

export function update(
  id: string,
  data: Partial<{ labelId: string; name: string; polygon: Point3D[] }>
): Promise<Annotation> {
  return http.put(`/v1/annotations/${id}`, data)
}

export function remove(id: string): Promise<void> {
  return http.delete(`/v1/annotations/${id}`)
}

export function rollback(id: string, version: number): Promise<Annotation> {
  return http.post(`/v1/annotations/${id}/rollback`, null, { params: { version } })
}

export function listSnapshots(id: string): Promise<AnnotationSnapshot[]> {
  return http.get(`/v1/annotations/${id}/snapshots`)
}
