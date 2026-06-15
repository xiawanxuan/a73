<template>
  <el-card class="label-panel" shadow="never">
    <template #header>
      <div class="panel-header">
        <el-icon><CollectionTag /></el-icon>
        <span>地形分类</span>
      </div>
    </template>
    <div v-if="labels.length === 0" class="loading-tip">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
    <el-radio-group
      v-else
      v-model="selectedId"
      @change="handleSelect"
      class="label-list"
    >
      <div
        v-for="label in labels"
        :key="label.id"
        class="label-item"
        :class="{ 'is-selected': label.id === selectedId }"
        @click="handleSelect(label.id)"
      >
        <div class="label-color" :style="{ background: label.color }"></div>
        <div class="label-info">
          <div class="label-name">{{ label.name }}</div>
          <div v-if="label.description" class="label-desc">{{ label.description }}</div>
        </div>
        <el-radio
          :value="label.id"
          style="visibility: hidden; position: absolute"
        />
      </div>
    </el-radio-group>
  </el-card>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { CollectionTag, Loading } from '@element-plus/icons-vue'
import type { TerrainLabel } from '@/types'

const props = defineProps<{
  labels: TerrainLabel[]
  selectedLabelId: string
}>()

const emit = defineEmits<{
  'select-label': [labelId: string]
}>()

const selectedId = ref(props.selectedLabelId)

watch(() => props.selectedLabelId, (val) => {
  selectedId.value = val
})

const handleSelect = (labelId: string) => {
  selectedId.value = labelId
  emit('select-label', labelId)
}
</script>

<style lang="scss" scoped>
.label-panel {
  border: none;

  :deep(.el-card__header) {
    padding: 12px 16px;
    border-bottom: 1px solid #ebeef5;
  }

  :deep(.el-card__body) {
    padding: 8px;
  }
}

.panel-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  font-size: 14px;
}

.loading-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 32px 16px;
  color: #909399;
  font-size: 13px;
}

.label-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.label-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border: 1px solid transparent;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background-color: #f5f7fa;
  }

  &.is-selected {
    border-color: #409eff;
    background-color: #ecf5ff;
  }
}

.label-color {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  flex-shrink: 0;
  border: 1px solid rgba(0, 0, 0, 0.1);
}

.label-info {
  flex: 1;
  min-width: 0;
}

.label-name {
  font-size: 13px;
  font-weight: 500;
  color: #303133;
  line-height: 1.4;
}

.label-desc {
  font-size: 12px;
  color: #909399;
  margin-top: 2px;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
