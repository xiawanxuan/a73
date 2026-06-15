import { ref, onBeforeUnmount } from 'vue'
import type { WSMessage, OnlineUser, CursorData } from '@/types'
import { getUserColor } from '@/utils/color'

export function useWebSocket() {
  const connected = ref(false)
  const onlineUsers = ref<OnlineUser[]>([])
  const cursorMap = ref<Map<string, CursorData>>(new Map())
  let ws: WebSocket | null = null
  let heartbeatTimer: number | null = null
  let reconnectTimer: number | null = null
  let reconnectCount = 0
  let reconnectDelay = 1000
  let currentPointCloudId = ''
  let currentUserId = ''
  let currentUsername = ''
  const listeners: Map<string, ((data: any) => void)[]> = new Map()

  function connect(pointCloudId: string, userId: string, username: string) {
    currentPointCloudId = pointCloudId
    currentUserId = userId
    currentUsername = username

    if (ws) {
      ws.close()
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    ws = new WebSocket(`${protocol}//${host}/ws?point_cloud_id=${pointCloudId}&user_id=${userId}&username=${encodeURIComponent(username)}`)

    ws.onopen = () => {
      connected.value = true
      reconnectCount = 0
      reconnectDelay = 1000
      startHeartbeat()
    }

    ws.onclose = () => {
      connected.value = false
      if (heartbeatTimer) {
        clearInterval(heartbeatTimer)
        heartbeatTimer = null
      }
      scheduleReconnect()
    }

    ws.onerror = () => {
      if (ws) {
        ws.close()
      }
    }

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        handleMessage(msg)
      } catch (e) {
        console.error('WebSocket message parse error:', e)
      }
    }
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
    if (ws) {
      ws.close()
      ws = null
    }
    connected.value = false
    onlineUsers.value = []
    cursorMap.value.clear()
  }

  function send(type: string, data: any) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type, data, timestamp: Date.now() }))
    }
  }

  function on(type: string, callback: (data: any) => void) {
    if (!listeners.has(type)) {
      listeners.set(type, [])
    }
    listeners.get(type)!.push(callback)
  }

  function off(type: string, callback: (data: any) => void) {
    const callbacks = listeners.get(type)
    if (callbacks) {
      const index = callbacks.indexOf(callback)
      if (index !== -1) {
        callbacks.splice(index, 1)
      }
    }
  }

  function emit(type: string, data: any) {
    const callbacks = listeners.get(type)
    if (callbacks) {
      callbacks.forEach(cb => cb(data))
    }
  }

  function broadcastAnnotation(type: 'annotation_create' | 'annotation_update' | 'annotation_delete' | 'annotation_rollback', annotation: any) {
    send(type, annotation)
  }

  function broadcastCursor(x: number, y: number, z: number) {
    send('cursor_move', { userId: currentUserId, pointCloudId: currentPointCloudId, x, y, z })
  }

  function broadcastDraft(data: string) {
    send('draft_sync', { userId: currentUserId, pointCloudId: currentPointCloudId, data })
  }

  function handleMessage(msg: WSMessage) {
    switch (msg.type) {
      case 'user_join':
        const joinUser: OnlineUser = {
          userId: msg.data.userId,
          username: msg.data.username,
          color: msg.data.color
        }
        if (!onlineUsers.value.find(u => u.userId === joinUser.userId)) {
          onlineUsers.value.push(joinUser)
        }
        break
      case 'user_leave':
        onlineUsers.value = onlineUsers.value.filter(u => u.userId !== msg.data.userId)
        cursorMap.value.delete(msg.data.userId)
        break
      case 'online_users':
        if (msg.data && msg.data.users) {
          onlineUsers.value = msg.data.users.map((u: any) => ({
            userId: u.userId,
            username: u.username,
            color: u.color
          }))
        }
        break
      case 'cursor_move':
        const cursor: CursorData = {
          userId: msg.data.userId,
          pointCloudId: msg.data.pointCloudId,
          x: msg.data.x,
          y: msg.data.y,
          z: msg.data.z,
          color: getUserColor(msg.data.userId)
        }
        cursorMap.value.set(cursor.userId, cursor)
        break
      case 'annotation_create':
      case 'annotation_update':
      case 'annotation_delete':
      case 'annotation_rollback':
        emit(msg.type, msg.data)
        break
      case 'draft_sync':
        emit('draft_sync', msg.data)
        break
      case 'ping':
        send('pong', {})
        break
      case 'pong':
        break
    }
  }

  function startHeartbeat() {
    heartbeatTimer = window.setInterval(() => {
      send('ping', {})
    }, 25000)
  }

  function scheduleReconnect() {
    if (reconnectTimer) {
      return
    }
    reconnectCount++
    reconnectDelay = Math.min(reconnectDelay * 2, 30000)
    reconnectTimer = window.setTimeout(() => {
      reconnectTimer = null
      connect(currentPointCloudId, currentUserId, currentUsername)
    }, reconnectDelay)
  }

  onBeforeUnmount(() => {
    disconnect()
  })

  return {
    connected,
    onlineUsers,
    cursorMap,
    connect,
    disconnect,
    send,
    on,
    off,
    broadcastAnnotation,
    broadcastCursor,
    broadcastDraft
  }
}
