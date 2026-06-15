import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as annotationApi from '@/api/annotation'
import * as labelApi from '@/api/label'
import type { Annotation, AnnotationSnapshot, TerrainLabel, Point3D } from '@/types'

export const useAnnotationStore = defineStore('annotation', () => {
  const annotations = ref<Annotation[]>([])
  const selectedAnnotationId = ref<string | null>(null)
  const currentLabelId = ref<string | null>(null)
  const labels = ref<TerrainLabel[]>([])
  const snapshots = ref<AnnotationSnapshot[]>([])
  const drawing = ref(false)
  const drawingPoints = ref<Point3D[]>([])
  const loading = ref(false)
  const searchKeyword = ref('')

  const selectedAnnotation = computed(() =>
    annotations.value.find(a => a.id === selectedAnnotationId.value) || null
  )

  const filteredAnnotations = computed(() => {
    if (!searchKeyword.value) return annotations.value
    const kw = searchKeyword.value.toLowerCase()
    return annotations.value.filter(a =>
      a.name.toLowerCase().includes(kw) ||
      a.label?.name.toLowerCase().includes(kw)
    )
  })

  async function fetchLabels() {
    labels.value = await labelApi.list()
    if (labels.value.length > 0 && !currentLabelId.value) {
      currentLabelId.value = labels.value[0].id
    }
  }

  async function createLabel(data: { name: string; color: string; description: string; icon: string }) {
    const newLabel = await labelApi.create(data)
    labels.value.push(newLabel)
    if (!currentLabelId.value) {
      currentLabelId.value = newLabel.id
    }
    return newLabel
  }

  async function updateLabel(id: string, data: Partial<TerrainLabel>) {
    const updated = await labelApi.update(id, data)
    const idx = labels.value.findIndex(l => l.id === id)
    if (idx !== -1) {
      labels.value[idx] = updated
    }
    annotations.value.forEach((ann, i) => {
      if (ann.labelId === id) {
        annotations.value[i] = { ...ann, label: updated }
      }
    })
    return updated
  }

  async function deleteLabel(id: string) {
    await labelApi.remove(id)
    labels.value = labels.value.filter(l => l.id !== id)
    if (currentLabelId.value === id && labels.value.length > 0) {
      currentLabelId.value = labels.value[0].id
    }
  }

  function addLabelFromRemote(label: TerrainLabel) {
    if (labels.value.find(l => l.id === label.id)) return
    labels.value.push(label)
  }

  function updateLabelFromRemote(label: TerrainLabel) {
    const idx = labels.value.findIndex(l => l.id === label.id)
    if (idx !== -1) {
      labels.value[idx] = label
    }
    annotations.value.forEach((ann, i) => {
      if (ann.labelId === label.id) {
        annotations.value[i] = { ...ann, label }
      }
    })
  }

  function deleteLabelFromRemote(id: string) {
    labels.value = labels.value.filter(l => l.id !== id)
    if (currentLabelId.value === id && labels.value.length > 0) {
      currentLabelId.value = labels.value[0].id
    }
  }

  async function fetchByPointCloud(pointCloudId: string) {
    loading.value = true
    try {
      annotations.value = await annotationApi.listByPointCloud(pointCloudId)
      selectedAnnotationId.value = null
    } finally {
      loading.value = false
    }
  }

  async function createAnnotation(pointCloudId: string, data: { labelId: string; name: string; polygon: Point3D[] }) {
    const res = await annotationApi.create(pointCloudId, data)
    const created = { ...res, label: labels.value.find(l => l.id === res.labelId) }
    annotations.value.push(created)
    selectedAnnotationId.value = created.id
    return created
  }

  async function updateAnnotation(id: string, data: Partial<Annotation>) {
    const res = await annotationApi.update(id, data)
    const idx = annotations.value.findIndex(a => a.id === id)
    if (idx !== -1) {
      annotations.value[idx] = { ...res, label: labels.value.find(l => l.id === res.labelId) }
    }
    return res
  }

  async function deleteAnnotation(id: string) {
    await annotationApi.remove(id)
    annotations.value = annotations.value.filter(a => a.id !== id)
    if (selectedAnnotationId.value === id) {
      selectedAnnotationId.value = null
    }
  }

  async function rollbackAnnotation(id: string, version: number) {
    const res = await annotationApi.rollback(id, version)
    const idx = annotations.value.findIndex(a => a.id === id)
    if (idx !== -1) {
      annotations.value[idx] = { ...res, label: labels.value.find(l => l.id === res.labelId) }
    }
    return res
  }

  async function fetchSnapshots(annotationId: string) {
    snapshots.value = await annotationApi.listSnapshots(annotationId)
  }

  function selectAnnotation(id: string | null) {
    selectedAnnotationId.value = id
  }

  function selectLabel(labelId: string) {
    currentLabelId.value = labelId
  }

  function startDrawing() {
    drawing.value = true
    drawingPoints.value = []
  }

  function addDrawingPoint(p: Point3D) {
    drawingPoints.value.push(p)
  }

  function undoDrawingPoint() {
    if (drawingPoints.value.length > 0) {
      drawingPoints.value.pop()
    }
  }

  function cancelDrawing() {
    drawing.value = false
    drawingPoints.value = []
  }

  function finishDrawing(): Point3D[] {
    const pts = [...drawingPoints.value]
    drawing.value = false
    drawingPoints.value = []
    return pts
  }

  return {
    annotations,
    selectedAnnotationId,
    selectedAnnotation,
    currentLabelId,
    labels,
    snapshots,
    drawing,
    drawingPoints,
    loading,
    searchKeyword,
    filteredAnnotations,
    fetchLabels,
    createLabel,
    updateLabel,
    deleteLabel,
    addLabelFromRemote,
    updateLabelFromRemote,
    deleteLabelFromRemote,
    fetchByPointCloud,
    createAnnotation,
    updateAnnotation,
    deleteAnnotation,
    rollbackAnnotation,
    fetchSnapshots,
    selectAnnotation,
    selectLabel,
    startDrawing,
    addDrawingPoint,
    undoDrawingPoint,
    cancelDrawing,
    finishDrawing
  }
})
