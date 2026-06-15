<template>
  <el-container class="workspace">
    <el-aside width="280px" class="left-panel">
      <div class="panel-header">
        <span class="title-icon">🌊</span>
        <span>海底勘测点云</span>
      </div>
      <div class="panel-content">
        <PointCloudTree
          :point-clouds="pcStore.filteredList"
          :selected-id="pcStore.currentPointCloud?.id || ''"
          @select="handleSelectPointCloud"
          @upload="handleUpload"
        />
      </div>
    </el-aside>
    <el-container>
      <el-header class="center-header" height="48px">
        <div class="header-left">
          <el-icon v-if="pcStore.currentPointCloud"><Document /></el-icon>
          <span class="current-name">{{ pcStore.currentPointCloud?.name || '未选择点云' }}</span>
          <el-tag v-if="syncStore.connected" type="success" size="small" effect="dark" class="sync-tag">
            <el-icon><Connection /></el-icon> 协同在线 {{ syncStore.onlineUsers.length }} 人
          </el-tag>
          <el-tag v-else type="info" size="small" class="sync-tag">
            <el-icon><Close /></el-icon> 离线
          </el-tag>
        </div>
        <div class="header-right">
          <el-tag v-for="u in syncStore.onlineUsers" :key="u.userId" size="small" :color="u.color" effect="light" style="margin-left: 4px;">
            {{ u.username }}
          </el-tag>
        </div>
      </el-header>
      <el-main class="center-panel">
        <PointCloudCanvas
          v-if="pcStore.currentPointCloud"
          ref="canvasRef"
          :point-cloud-id="pcStore.currentPointCloud.id"
          :point-bounds="{
            minX: pcStore.currentPointCloud.boundsMinX,
            maxX: pcStore.currentPointCloud.boundsMaxX,
            minY: pcStore.currentPointCloud.boundsMinY,
            maxY: pcStore.currentPointCloud.boundsMaxY,
            minZ: pcStore.currentPointCloud.boundsMinZ,
            maxZ: pcStore.currentPointCloud.boundsMaxZ
          }"
          :points="pcStore.mockPoints"
          @point-click="handlePointClick"
          @polygon-created="handlePolygonCreated"
          @canvas-ready="handleCanvasReady"
        />
        <div v-else class="empty-canvas">
          <el-empty description="请从左侧选择海底点云文件" :image-size="120" />
        </div>
      </el-main>
      <el-aside width="320px" class="right-panel">
        <div class="panel-header">
          <span class="title-icon">🏷️</span>
          <span>标注工具</span>
        </div>
        <div class="panel-content">
          <AnnotationToolbar
            :disabled="!pcStore.currentPointCloud"
            :drawing-point-count="annStore.drawingPoints.length"
            @start-draw="handleStartDraw"
            @cancel-draw="handleCancelDraw"
            @finish-draw="handleFinishDraw"
            @undo-point="handleUndoPoint"
            @clear-all="handleClearAll"
            @delete-selected="handleDeleteSelected"
          />

          <el-divider content-position="left">地形分类</el-divider>
          <LabelPanel
            :labels="annStore.labels"
            :selected-label-id="annStore.currentLabelId || ''"
            @select-label="handleSelectLabel"
          />

          <el-divider content-position="left">已保存标注</el-divider>
          <AnnotationList
            :annotations="annStore.filteredAnnotations"
            :selected-id="annStore.selectedAnnotationId || ''"
            @select="handleSelectAnnotation"
            @edit="handleEditAnnotation"
            @delete="handleDeleteAnnotation"
            @view-history="handleViewHistory"
          />
        </div>
      </el-aside>
    </el-container>
    <VersionHistory
      :visible="historyVisible"
      :snapshots="annStore.snapshots"
      @close="historyVisible = false"
      @rollback="handleRollback"
    />
  </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Connection, Close } from '@element-plus/icons-vue'
import PointCloudCanvas from '@/components/PointCloudCanvas.vue'
import AnnotationToolbar from '@/components/AnnotationToolbar.vue'
import LabelPanel from '@/components/LabelPanel.vue'
import AnnotationList from '@/components/AnnotationList.vue'
import PointCloudTree from '@/components/PointCloudTree.vue'
import VersionHistory from '@/components/VersionHistory.vue'
import { usePointCloudStore } from '@/stores/pointCloud'
import { useAnnotationStore } from '@/stores/annotation'
import { useSyncStore } from '@/stores/sync'
import type { PointCloud, Annotation, Point3D } from '@/types'
import { useDraftCache } from '@/composables/useDraftCache'

const pcStore = usePointCloudStore()
const annStore = useAnnotationStore()
const syncStore = useSyncStore()
const draftCache = useDraftCache()

const canvasRef = ref<InstanceType<typeof PointCloudCanvas> | null>(null)
const historyVisible = ref(false)
const annotationNameCounter = ref(1)

onMounted(async () => {
  syncStore.ensureUser()
  await Promise.all([
    pcStore.fetchList(),
    annStore.fetchLabels()
  ])
})

watch(() => pcStore.currentPointCloud, async (pc, oldPc) => {
  if (oldPc?.id) {
    syncStore.disconnect()
  }
  if (!pc) return
  await annStore.fetchByPointCloud(pc.id)
  await syncStore.connect(pc.id)
  syncStore.onAnnotationEvent(handleRemoteAnnotationEvent)
  nextTick(() => {
    if (canvasRef.value) {
      annStore.annotations.forEach(a => canvasRef.value?.addAnnotationPolygon(a))
    }
  })
  const draft = draftCache.getDraft(pc.id, syncStore.currentUserId)
  if (draft && draft.points?.length > 0) {
    ElMessage.info(`检测到本地草稿，已恢复 ${draft.points.length} 个标注点`)
  }
}, { immediate: false })

function handleRemoteAnnotationEvent(action: string, data: any) {
  switch (action) {
    case 'create':
      if (!annStore.annotations.find(a => a.id === data.id)) {
        const created = { ...data, label: annStore.labels.find(l => l.id === data.labelId) }
        annStore.annotations.push(created)
        canvasRef.value?.addAnnotationPolygon(created)
        ElMessage.info(`协同更新: 新增标注 "${data.name}"`)
      }
      break
    case 'update':
      const idx = annStore.annotations.findIndex(a => a.id === data.id)
      if (idx !== -1) {
        annStore.annotations[idx] = { ...data, label: annStore.labels.find(l => l.id === data.labelId) }
        canvasRef.value?.clearAnnotations()
        annStore.annotations.forEach(a => canvasRef.value?.addAnnotationPolygon(a))
        ElMessage.info(`协同更新: 标注 "${data.name}" 已修改`)
      }
      break
    case 'delete':
      annStore.annotations = annStore.annotations.filter(a => a.id !== data.id)
      canvasRef.value?.clearAnnotations()
      annStore.annotations.forEach(a => canvasRef.value?.addAnnotationPolygon(a))
      ElMessage.info(`协同更新: 标注已删除`)
      break
    case 'rollback':
      const ridx = annStore.annotations.findIndex(a => a.id === data.id)
      if (ridx !== -1) {
        annStore.annotations[ridx] = { ...data, label: annStore.labels.find(l => l.id === data.labelId) }
      } else {
        annStore.annotations.push({ ...data, label: annStore.labels.find(l => l.id === data.labelId) })
      }
      canvasRef.value?.clearAnnotations()
      annStore.annotations.forEach(a => canvasRef.value?.addAnnotationPolygon(a))
      ElMessage.info(`协同更新: 标注已回滚`)
      break
  }
}

async function handleSelectPointCloud(pc: PointCloud) {
  if (annStore.drawing) {
    try {
      await ElMessageBox.confirm('当前正在绘制标注，切换点云将丢失进度，是否继续？', '提示', {
        confirmButtonText: '继续',
        cancelButtonText: '取消',
        type: 'warning'
      })
    } catch {
      return
    }
    annStore.cancelDrawing()
  }
  pcStore.selectPointCloud(pc)
}

function handleUpload() {
  ElMessage.info('点云上传功能开发中')
}

function handleCanvasReady() {
  annStore.annotations.forEach(a => canvasRef.value?.addAnnotationPolygon(a))
}

function handlePointClick(p: Point3D) {
  if (annStore.drawing) {
    annStore.addDrawingPoint(p)
    canvasRef.value?.addPolygonPoint(p)
    draftCache.setDraft(pcStore.currentPointCloud!.id, syncStore.currentUserId, {
      points: annStore.drawingPoints
    })
  }
  syncStore.broadcastCursor(p)
}

function handlePolygonCreated(polygon: Point3D[]) {
  saveAnnotation(polygon)
}

async function handleStartDraw() {
  if (!annStore.currentLabelId) {
    ElMessage.warning('请先选择地形分类标签')
    return
  }
  annStore.startDrawing()
  canvasRef.value?.startPolygonDrawing()
  ElMessage.info('开始绘制：在点云上点击添加多边形顶点，至少 3 个点')
}

function handleCancelDraw() {
  annStore.cancelDrawing()
  canvasRef.value?.cancelPolygonDrawing()
  draftCache.clearDraft(pcStore.currentPointCloud!.id, syncStore.currentUserId)
}

function handleUndoPoint() {
  annStore.undoDrawingPoint()
  canvasRef.value?.cancelPolygonDrawing()
  annStore.drawingPoints.forEach(p => canvasRef.value?.addPolygonPoint(p))
}

async function handleFinishDraw() {
  if (annStore.drawingPoints.length < 3) {
    ElMessage.warning('至少需要 3 个点才能创建多边形标注')
    return
  }
  const polygon = annStore.finishDrawing()
  canvasRef.value?.finishPolygonDrawing()
  await saveAnnotation(polygon)
}

async function saveAnnotation(polygon: Point3D[]) {
  const label = annStore.labels.find(l => l.id === annStore.currentLabelId)
  const name = `${label?.name || '标注'}_${annotationNameCounter.value++}`
  try {
    const created = await annStore.createAnnotation(pcStore.currentPointCloud!.id, {
      labelId: annStore.currentLabelId!,
      name,
      polygon
    })
    canvasRef.value?.addAnnotationPolygon(created)
    syncStore.broadcastAnnotation('create', created)
    draftCache.clearDraft(pcStore.currentPointCloud!.id, syncStore.currentUserId)
    ElMessage.success(`标注 "${name}" 已保存`)
  } catch (e: any) {
    ElMessage.error('保存标注失败: ' + (e?.message || e))
  }
}

function handleClearAll() {
  ElMessageBox.confirm('确定清除画布上的所有标注？此操作不会删除已保存的标注。', '提示', {
    type: 'warning'
  }).then(() => {
    canvasRef.value?.clearAnnotations()
    annStore.annotations.forEach(a => canvasRef.value?.addAnnotationPolygon(a))
  }).catch(() => {})
}

async function handleDeleteSelected() {
  if (!annStore.selectedAnnotationId) {
    ElMessage.warning('请先选择要删除的标注')
    return
  }
  const ann = annStore.selectedAnnotation
  try {
    await ElMessageBox.confirm(`确定删除标注 "${ann?.name}"？`, '删除确认', { type: 'warning' })
    await annStore.deleteAnnotation(annStore.selectedAnnotationId)
    syncStore.broadcastAnnotation('delete', { id: annStore.selectedAnnotationId } as any)
    canvasRef.value?.clearAnnotations()
    annStore.annotations.forEach(a => canvasRef.value?.addAnnotationPolygon(a))
    ElMessage.success('已删除')
  } catch {}
}

function handleSelectLabel(labelId: string) {
  annStore.selectLabel(labelId)
}

function handleSelectAnnotation(id: string) {
  annStore.selectAnnotation(id)
}

function handleEditAnnotation(_ann: Annotation) {
  ElMessage.info('编辑标注功能开发中')
}

async function handleDeleteAnnotation(id: string) {
  const ann = annStore.annotations.find(a => a.id === id)
  try {
    await ElMessageBox.confirm(`确定删除标注 "${ann?.name}"？`, '删除确认', { type: 'warning' })
    await annStore.deleteAnnotation(id)
    syncStore.broadcastAnnotation('delete', { id } as any)
    canvasRef.value?.clearAnnotations()
    annStore.annotations.forEach(a => canvasRef.value?.addAnnotationPolygon(a))
    ElMessage.success('已删除')
  } catch {}
}

async function handleViewHistory(ann: Annotation) {
  await annStore.fetchSnapshots(ann.id)
  historyVisible.value = true
}

async function handleRollback(version: number) {
  if (!annStore.selectedAnnotationId) return
  try {
    const rolled = await annStore.rollbackAnnotation(annStore.selectedAnnotationId, version)
    syncStore.broadcastAnnotation('rollback', rolled)
    canvasRef.value?.clearAnnotations()
    annStore.annotations.forEach(a => canvasRef.value?.addAnnotationPolygon(a))
    historyVisible.value = false
    ElMessage.success(`已回滚到版本 V${version}`)
  } catch (e: any) {
    ElMessage.error('回滚失败: ' + (e?.message || e))
  }
}
</script>

<style lang="scss" scoped>
.workspace {
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.left-panel,
.right-panel {
  background-color: #f5f7fa;
  border-right: 1px solid #e4e7ed;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.right-panel {
  border-right: none;
  border-left: 1px solid #e4e7ed;
}

.panel-header {
  padding: 12px 16px;
  font-weight: 600;
  font-size: 14px;
  border-bottom: 1px solid #e4e7ed;
  background-color: #fff;
  display: flex;
  align-items: center;
  gap: 8px;

  .title-icon {
    font-size: 16px;
  }
}

.panel-content {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
}

.center-header {
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;

  .header-left {
    display: flex;
    align-items: center;
    gap: 10px;

    .current-name {
      font-weight: 600;
      font-size: 15px;
    }

    .sync-tag {
      margin-left: 8px;
    }
  }

  .header-right {
    display: flex;
    align-items: center;
  }
}

.center-panel {
  padding: 0;
  background-color: #0a1628;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.empty-canvas {
  color: #888;
  font-size: 16px;
}
</style>
