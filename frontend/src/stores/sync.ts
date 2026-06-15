import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useWebSocket } from '@/composables/useWebSocket'
import type { Annotation, Point3D } from '@/types'

export const useSyncStore = defineStore('sync', () => {
  const ws = useWebSocket()
  const currentUserId = ref<string>('')
  const currentUsername = ref<string>('')
  const activePointCloudId = ref<string>('')

  function ensureUser() {
    let uid: string = localStorage.getItem('user_id') || ''
    let uname: string = localStorage.getItem('username') || ''
    if (!uid) {
      uid = 'user_' + Math.random().toString(36).slice(2, 10)
      uname = '用户' + Math.floor(Math.random() * 1000)
      localStorage.setItem('user_id', uid)
      localStorage.setItem('username', uname)
    }
    currentUserId.value = uid
    currentUsername.value = uname
    return { uid, uname }
  }

  async function connect(pointCloudId: string) {
    const { uid, uname } = ensureUser()
    activePointCloudId.value = pointCloudId
    await ws.connect(pointCloudId, uid, uname)
  }

  function disconnect() {
    ws.disconnect()
    activePointCloudId.value = ''
  }

  function onAnnotationEvent(callback: (action: string, annotation: Annotation) => void) {
    const actions = ['annotation_create', 'annotation_update', 'annotation_delete', 'annotation_rollback']
    actions.forEach(type => {
      ws.on(type, (data: any) => {
        const action = type.replace('annotation_', '')
        callback(action, data)
      })
    })
  }

  function offAnnotationEvent() {
    const actions = ['annotation_create', 'annotation_update', 'annotation_delete', 'annotation_rollback']
    actions.forEach(type => ws.off(type, () => {}))
  }

  function broadcastAnnotation(action: 'create' | 'update' | 'delete' | 'rollback', annotation: Annotation) {
    ws.broadcastAnnotation(`annotation_${action}` as any, annotation)
  }

  function broadcastCursor(pos: Point3D) {
    ws.broadcastCursor(pos.x, pos.y, pos.z)
  }

  function broadcastDraft(data: string) {
    ws.broadcastDraft(data)
  }

  return {
    connected: ws.connected,
    onlineUsers: ws.onlineUsers,
    cursorMap: ws.cursorMap,
    currentUserId,
    currentUsername,
    activePointCloudId,
    ensureUser,
    connect,
    disconnect,
    onAnnotationEvent,
    offAnnotationEvent,
    broadcastAnnotation,
    broadcastCursor,
    broadcastDraft
  }
})
