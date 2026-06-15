<template>
  <div class="annotation-toolbar">
    <el-button
      type="primary"
      :disabled="disabled"
      @click="$emit('start-draw')"
    >
      <el-icon><Edit /></el-icon>
      <span>开始绘制多边形</span>
    </el-button>
    <el-button
      :disabled="disabled"
      @click="$emit('undo-point')"
    >
      <el-icon><RefreshLeft /></el-icon>
      <span>撤销上一点</span>
    </el-button>
    <el-button
      :disabled="disabled"
      @click="$emit('cancel-draw')"
    >
      <el-icon><Close /></el-icon>
      <span>取消绘制</span>
    </el-button>
    <el-button
      type="success"
      :disabled="disabled"
      @click="$emit('finish-draw')"
    >
      <el-icon><Check /></el-icon>
      <span>完成绘制</span>
    </el-button>
    <el-button
      type="danger"
      :disabled="disabled"
      @click="$emit('clear-all')"
    >
      <el-icon><Delete /></el-icon>
      <span>清除全部</span>
    </el-button>
    <el-button
      type="warning"
      :disabled="disabled"
      @click="$emit('delete-selected')"
    >
      <el-icon><Remove /></el-icon>
      <span>删除选中</span>
    </el-button>
    <div v-if="drawingPointCount > 0" class="point-count">
      <el-tag type="info" size="small">
        当前点数: {{ drawingPointCount }}
      </el-tag>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  Edit,
  RefreshLeft,
  Close,
  Check,
  Delete,
  Remove
} from '@element-plus/icons-vue'

withDefaults(defineProps<{
  disabled?: boolean
  drawingPointCount?: number
}>(), {
  disabled: false,
  drawingPointCount: 0
})

defineEmits<{
  'start-draw': []
  'cancel-draw': []
  'finish-draw': []
  'undo-point': []
  'clear-all': []
  'delete-selected': []
}>()
</script>

<style lang="scss" scoped>
.annotation-toolbar {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;

  .el-button {
    width: 100%;
    justify-content: flex-start;
    padding: 8px 12px;

    .el-icon {
      margin-right: 6px;
    }
  }

  .point-count {
    margin-top: 4px;
    text-align: center;
  }
}
</style>
