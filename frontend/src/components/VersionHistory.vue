<template>
  <el-dialog
    v-model="dialogVisible"
    title="版本历史"
    width="560px"
    :close-on-click-modal="false"
    @closed="$emit('close')"
  >
    <div class="snapshot-list">
      <div
        v-for="snapshot in sortedSnapshots"
        :key="snapshot.id"
        class="snapshot-item"
      >
        <div class="snapshot-main">
          <div class="version-badge">V{{ snapshot.version }}</div>
          <div class="snapshot-info">
            <div class="snapshot-header">
              <el-tag
                :type="operationTagType(snapshot.operation)"
                size="small"
                effect="light"
              >
                {{ operationText(snapshot.operation) }}
              </el-tag>
              <span class="operator">{{ snapshot.operatorName || '未知用户' }}</span>
            </div>
            <div class="snapshot-time">{{ formatTime(snapshot.createdAt) }}</div>
          </div>
        </div>
        <el-button
          type="primary"
          size="small"
          plain
          @click="$emit('rollback', snapshot.version)"
        >
          <el-icon><RefreshLeft /></el-icon>
          <span>回滚到此版本</span>
        </el-button>
      </div>
      <div v-if="sortedSnapshots.length === 0" class="empty-tip">
        <el-empty description="暂无版本记录" :image-size="60" />
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { RefreshLeft } from '@element-plus/icons-vue'
import type { AnnotationSnapshot, OperationType } from '@/types'

const props = defineProps<{
  snapshots: AnnotationSnapshot[]
  visible: boolean
}>()

defineEmits<{
  'close': []
  'rollback': [version: number]
}>()

const dialogVisible = ref(props.visible)

watch(() => props.visible, (val) => {
  dialogVisible.value = val
})

const sortedSnapshots = computed(() => {
  return [...props.snapshots].sort((a, b) => b.version - a.version)
})

const operationTagType = (operation: OperationType): 'success' | 'warning' | 'danger' | 'info' => {
  const map: Record<OperationType, 'success' | 'warning' | 'danger' | 'info'> = {
    create: 'success',
    update: 'warning',
    delete: 'danger',
    rollback: 'info'
  }
  return map[operation] || 'info'
}

const operationText = (operation: OperationType) => {
  const map: Record<OperationType, string> = {
    create: '创建',
    update: '更新',
    delete: '删除',
    rollback: '回滚'
  }
  return map[operation] || operation
}

const formatTime = (timeStr: string) => {
  const date = new Date(timeStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}
</script>

<style lang="scss" scoped>
.snapshot-list {
  max-height: 400px;
  overflow-y: auto;
}

.snapshot-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #ebeef5;

  &:last-child {
    border-bottom: none;
  }
}

.snapshot-main {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.version-badge {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  flex-shrink: 0;
}

.snapshot-info {
  flex: 1;
  min-width: 0;
}

.snapshot-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.operator {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

.snapshot-time {
  font-size: 12px;
  color: #909399;
  font-family: monospace;
}

.empty-tip {
  padding: 40px 0;
}

:deep(.el-button .el-icon) {
  margin-right: 4px;
}
</style>
