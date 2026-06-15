<template>
  <div class="annotation-list">
    <div class="list-header">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索标注名称"
        clearable
        size="small"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>
    <div class="list-content">
      <div
        v-for="item in filteredAnnotations"
        :key="item.id"
        class="annotation-item"
        :class="{ 'is-selected': item.id === selectedId }"
        @click="$emit('select', item.id)"
      >
        <div class="item-main">
          <div
            class="color-dot"
            :style="{ background: item.label?.color || '#909399' }"
          ></div>
          <div class="item-info">
            <div class="item-name">{{ item.name }}</div>
            <div class="item-meta">
              <span class="terrain-type">{{ item.label?.name || '未分类' }}</span>
              <span class="create-time">{{ formatTime(item.createdAt) }}</span>
            </div>
          </div>
        </div>
        <div class="item-actions">
          <el-button
            link
            type="primary"
            size="small"
            @click.stop="$emit('edit', item)"
          >
            <el-icon><Edit /></el-icon>
          </el-button>
          <el-button
            link
            type="primary"
            size="small"
            @click.stop="$emit('view-history', item)"
          >
            <el-icon><Clock /></el-icon>
          </el-button>
          <el-button
            link
            type="danger"
            size="small"
            @click.stop="$emit('delete', item.id)"
          >
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
      </div>
      <div v-if="filteredAnnotations.length === 0" class="empty-tip">
        <el-empty :description="searchKeyword ? '无匹配结果' : '暂无标注'" :image-size="60" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Search, Edit, Clock, Delete } from '@element-plus/icons-vue'
import type { Annotation } from '@/types'

const props = defineProps<{
  annotations: Annotation[]
  selectedId: string
}>()

defineEmits<{
  'select': [id: string]
  'edit': [annotation: Annotation]
  'delete': [id: string]
  'view-history': [annotation: Annotation]
}>()

const searchKeyword = ref('')

const filteredAnnotations = computed(() => {
  if (!searchKeyword.value) {
    return props.annotations
  }
  const keyword = searchKeyword.value.toLowerCase()
  return props.annotations.filter(a =>
    a.name.toLowerCase().includes(keyword)
  )
})

const formatTime = (timeStr: string) => {
  const date = new Date(timeStr)
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${month}-${day} ${hours}:${minutes}`
}
</script>

<style lang="scss" scoped>
.annotation-list {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #fff;
}

.list-header {
  padding: 12px;
  border-bottom: 1px solid #ebeef5;
}

.list-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.annotation-item {
  padding: 10px 12px;
  border: 1px solid transparent;
  border-radius: 6px;
  margin-bottom: 4px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background-color: #f5f7fa;

    .item-actions {
      opacity: 1;
    }
  }

  &.is-selected {
    border-color: #409eff;
    background-color: #ecf5ff;
  }
}

.item-main {
  display: flex;
  align-items: center;
  gap: 10px;
}

.color-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 13px;
  font-weight: 500;
  color: #303133;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
  font-size: 11px;
  color: #909399;
}

.terrain-type {
  background-color: #f0f2f5;
  padding: 1px 6px;
  border-radius: 3px;
}

.create-time {
  font-family: monospace;
}

.item-actions {
  display: flex;
  gap: 4px;
  margin-top: 6px;
  padding-top: 6px;
  border-top: 1px dashed #ebeef5;
  opacity: 0;
  transition: opacity 0.2s ease;
  justify-content: flex-end;
}

.empty-tip {
  padding: 40px 0;
}
</style>
