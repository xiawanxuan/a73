export interface ApiResponse<T> {
  code: number
  msg?: string
  data: T
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

export interface Point3D {
  x: number
  y: number
  z: number
}

export interface PointBounds {
  minX: number
  maxX: number
  minY: number
  maxY: number
  minZ: number
  maxZ: number
}

export interface TerrainLabel {
  id: string
  name: string
  color: string
  description: string
  icon: string
  isSystem: boolean
  createdAt: string
  updatedAt: string
}

export interface Annotation {
  id: string
  pointCloudId: string
  labelId: string
  name: string
  polygon: Point3D[]
  boundsCenterX: number
  boundsCenterY: number
  boundsCenterZ: number
  creatorId: string
  createdAt: string
  updatedAt: string
  label?: TerrainLabel
}

export type OperationType = 'create' | 'update' | 'delete' | 'rollback'

export interface AnnotationSnapshot {
  id: string
  annotationId: string
  pointCloudId: string
  version: number
  snapshot: Record<string, any>
  operatorId: string
  operation: OperationType
  createdAt: string
  operatorName?: string
}

export type PointCloudStatus = 'ready' | 'processing' | 'error'

export interface PointCloud {
  id: string
  name: string
  filePath: string
  fileSize: number
  pointCount: number
  boundsMinX: number
  boundsMinY: number
  boundsMinZ: number
  boundsMaxX: number
  boundsMaxY: number
  boundsMaxZ: number
  uploadedBy: string
  status: PointCloudStatus
  createdAt: string
  updatedAt: string
}

export interface OnlineUser {
  userId: string
  username: string
  color: string
}

export interface CursorData {
  userId: string
  pointCloudId?: string
  x: number
  y: number
  z: number
  color?: string
}

export type WSMessageType =
  | 'annotation_create'
  | 'annotation_update'
  | 'annotation_delete'
  | 'annotation_rollback'
  | 'user_join'
  | 'user_leave'
  | 'online_users'
  | 'cursor_move'
  | 'draft_sync'
  | 'ping'
  | 'pong'

export interface WSMessage {
  type: WSMessageType
  data: any
  timestamp?: number
}
