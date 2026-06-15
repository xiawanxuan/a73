<template>
  <el-card class="label-panel" shadow="never">
    <template #header>
      <div class="panel-header">
        <el-icon><CollectionTag /></el-icon>
        <span>地形分类</span>
        <el-button
          type="primary"
          size="small"
          class="add-btn"
          icon="Plus"
          @click="openCreateDialog"
        >
          新增
        </el-button>
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
        :class="{ 'is-selected': label.id === selectedId, 'is-system': label.isSystem }"
        @click="handleSelect(label.id)"
      >
        <div class="label-color" :style="{ background: label.color }"></div>
        <div class="label-info">
          <div class="label-name">
            {{ label.name }}
            <el-tag v-if="label.isSystem" type="info" size="small" effect="plain" class="system-tag">系统</el-tag>
          </div>
          <div v-if="label.description" class="label-desc">{{ label.description }}</div>
        </div>
        <div class="label-actions" @click.stop>
          <el-button
            v-if="!label.isSystem"
            text
            size="small"
            type="primary"
            @click="openEditDialog(label)"
          >
            编辑
          </el-button>
          <el-button
            v-if="!label.isSystem"
            text
            size="small"
            type="danger"
            @click="handleDelete(label)"
          >
            删除
          </el-button>
        </div>
        <el-radio
          :value="label.id"
          style="visibility: hidden; position: absolute"
        />
      </div>
    </el-radio-group>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑地形标签' : '新增地形标签'"
      width="420px"
      @closed="resetForm"
    >
      <el-form :model="form" label-width="80px" ref="formRef">
        <el-form-item label="名称" prop="name" :rules="[{ required: true, message: '请输入标签名称', trigger: 'blur' }]">
          <el-input v-model="form.name" placeholder="请输入标签名称" maxlength="32" show-word-limit />
        </el-form-item>
        <el-form-item label="配色" prop="color">
          <el-color-picker v-model="form.color" show-alpha />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-select v-model="form.icon" placeholder="选择图标">
            <el-option label="地形" value="terrain" />
            <el-option label="海沟" value="trending_down" />
            <el-option label="礁石" value="warning" />
            <el-option label="管线" value="cable" />
            <el-option label="标记" value="flag" />
            <el-option label="定位" value="location" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述信息"
            maxlength="128"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { CollectionTag, Loading, Plus } from '@element-plus/icons-vue'
import type { TerrainLabel, FormInstance, FormRules } from '@/types'

const props = defineProps<{
  labels: TerrainLabel[]
  selectedLabelId: string
}>()

const emit = defineEmits<{
  'select-label': [labelId: string]
  'create-label': [data: { name: string; color: string; description: string; icon: string }]
  'update-label': [id: string, data: { name: string; color: string; description: string; icon: string }]
  'delete-label': [id: string]
}>()

const selectedId = ref(props.selectedLabelId)
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const editingId = ref('')
const formRef = ref<FormInstance>()

const form = ref({
  name: '',
  color: '#1E88E5',
  icon: 'flag',
  description: ''
})

watch(() => props.selectedLabelId, (val) => {
  selectedId.value = val
})

const handleSelect = (labelId: string) => {
  selectedId.value = labelId
  emit('select-label', labelId)
}

const openCreateDialog = () => {
  isEdit.value = false
  editingId.value = ''
  resetForm()
  dialogVisible.value = true
}

const openEditDialog = (label: TerrainLabel) => {
  isEdit.value = true
  editingId.value = label.id
  form.value = {
    name: label.name,
    color: label.color,
    icon: label.icon,
    description: label.description
  }
  dialogVisible.value = true
}

const resetForm = () => {
  form.value = {
    name: '',
    color: '#1E88E5',
    icon: 'flag',
    description: ''
  }
  formRef.value?.resetFields()
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const data = {
        name: form.value.name,
        color: form.value.color,
        icon: form.value.icon,
        description: form.value.description
      }
      if (isEdit.value) {
        emit('update-label', editingId.value, data)
      } else {
        emit('create-label', data)
      }
      dialogVisible.value = false
    } finally {
      submitting.value = false
    }
  })
}

const handleDelete = async (label: TerrainLabel) => {
  try {
    await ElMessageBox.confirm(
      `确定删除标签「${label.name}」吗？删除后不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    emit('delete-label', label.id)
  } catch {
    // 取消删除
  }
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

  .add-btn {
    margin-left: auto;
  }
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
  position: relative;

  &:hover {
    background-color: #f5f7fa;

    .label-actions {
      opacity: 1;
    }
  }

  &.is-selected {
    border-color: #409eff;
    background-color: #ecf5ff;
  }

  &.is-system {
    .label-name {
      font-weight: 600;
    }
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
  display: flex;
  align-items: center;
  gap: 6px;

  .system-tag {
    font-weight: 400;
    transform: scale(0.85);
    transform-origin: left center;
  }
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

.label-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s ease;
  flex-shrink: 0;
}
</style>
