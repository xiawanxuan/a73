<template>
  <div class="point-cloud-tree">
    <div class="tree-header">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索点云"
        clearable
        size="small"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-button
        type="primary"
        size="small"
        @click="$emit('upload')"
      >
        <el-icon><Upload /></el-icon>
        <span>上传点云</span>
      </el-button>
    </div>
    <div class="tree-content">
      <el-tree
        :data="treeData"
        node-key="id"
        :highlight-current="true"
        :expand-on-click-node="false"
        :current-node-key="selectedId"
        @node-click="handleNodeClick"
      >
        <template #default="{ node, data }">
          <div class="tree-node" :class="{ 'is-selected': data.id === selectedId }">
            <el-icon class="file-icon"><Document /></el-icon>
            <span class="node-name">{{ node.label }}</span>
            <div class="node-tags">
              <el-tag
                :type="statusTagType(data.status)"
                size="small"
                effect="light"
              >
                {{ statusText(data.status) }}
              </el-tag>
              <span class="point-count">{{ formatPointCount(data.pointCount) }}</span>
            </div>
          </div>
        </template>
      </el-tree>
      <div v-if="treeData.length === 0" class="empty-tip">
        <el-empty description="暂无点云数据" :image-size="60" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Search, Upload, Document } from '@element-plus/icons-vue'
import type { PointCloud, PointCloudStatus } from '@/types'

const props = defineProps<{
  pointClouds: PointCloud[]
  selectedId: string
}>()

const emit = defineEmits<{
  'select': [pointCloud: PointCloud]
  'upload': []
}>()

const searchKeyword = ref('')

const treeData = computed(() => {
  if (!searchKeyword.value) {
    return props.pointClouds
  }
  const keyword = searchKeyword.value.toLowerCase()
  return props.pointClouds.filter(p =>
    p.name.toLowerCase().includes(keyword)
  )
})

const handleNodeClick = (data: PointCloud) => {
  emit('select', data)
}

const statusTagType = (status: PointCloudStatus) => {
  const map: Record<PointCloudStatus, 'success' | 'warning' | 'danger'> = {
    ready: 'success',
    processing: 'warning',
    error: 'danger'
  }
  return map[status] || 'info'
}

const statusText = (status: PointCloudStatus) => {
  const map: Record<PointCloudStatus, string> = {
    ready: '就绪',
    processing: '处理中',
    error: '错误'
  }
  return map[status] || status
}

const formatPointCount = (count: number) => {
  if (count >= 1000000) {
    return (count / 1000000).toFixed(1) + 'M'
  }
  if (count >= 1000) {
    return (count / 1000).toFixed(1) + 'K'
  }
  return String(count)
}
</script>

<style lang="scss" scoped>
.point-cloud-tree {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #fff;
}

.tree-header {
  padding: 12px;
  border-bottom: 1px solid #ebeef5;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.tree-content {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 2px 0;
}

.file-icon {
  color: #606266;
  flex-shrink: 0;
}

.node-name {
  flex: 1;
  font-size: 13px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.node-tags {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.point-count {
  font-size: 11px;
  color: #909399;
  font-family: monospace;
}

.empty-tip {
  padding: 40px 0;
}

:deep(.el-tree-node__content) {
  height: 36px;
}

:deep(.el-tree-node.is-current > .el-tree-node__content) {
  background-color: #ecf5ff;
}
</style>
