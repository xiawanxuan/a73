import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as pointCloudApi from '@/api/pointCloud'
import type { PointCloud, Point3D } from '@/types'

export const usePointCloudStore = defineStore('pointCloud', () => {
  const currentPointCloud = ref<PointCloud | null>(null)
  const pointCloudList = ref<PointCloud[]>([])
  const loading = ref(false)
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)
  const searchKeyword = ref('')
  const mockPoints = ref<Point3D[]>([])

  const filteredList = computed(() => {
    if (!searchKeyword.value) return pointCloudList.value
    const kw = searchKeyword.value.toLowerCase()
    return pointCloudList.value.filter(pc => pc.name.toLowerCase().includes(kw))
  })

  async function fetchList() {
    loading.value = true
    try {
      const res = await pointCloudApi.list(page.value, pageSize.value)
      pointCloudList.value = res.list
      total.value = res.total
      if (pointCloudList.value.length > 0 && !currentPointCloud.value) {
        selectPointCloud(pointCloudList.value[0])
      }
    } finally {
      loading.value = false
    }
  }

  function selectPointCloud(pointCloud: PointCloud) {
    currentPointCloud.value = pointCloud
    generateMockPoints(pointCloud)
  }

  function generateMockPoints(pc: PointCloud) {
    const count = Math.min(pc.pointCount || 100000, 500000)
    const points: Point3D[] = []
    const minX = pc.boundsMinX ?? -500
    const maxX = pc.boundsMaxX ?? 500
    const minY = pc.boundsMinY ?? -500
    const maxY = pc.boundsMaxY ?? 500
    const minZ = pc.boundsMinZ ?? -200
    const maxZ = pc.boundsMaxZ ?? 0

    for (let i = 0; i < count; i++) {
      const x = minX + Math.random() * (maxX - minX)
      const y = minY + Math.random() * (maxY - minY)
      const baseZ = minZ + Math.random() * (maxZ - minZ)
      const terrainWave = Math.sin(x * 0.01) * Math.cos(y * 0.01) * 30
      const trenchFeature = (x > 100 && x < 200 && y > -50 && y < 50) ? -60 : 0
      const rockFeature = (Math.sqrt((x + 200) ** 2 + (y + 150) ** 2) < 60) ? 40 : 0
      const z = baseZ + terrainWave + trenchFeature + rockFeature
      points.push({ x, y, z })
    }
    mockPoints.value = points
  }

  async function createPointCloud(data: Partial<PointCloud>) {
    const res = await pointCloudApi.create(data)
    pointCloudList.value.unshift(res)
    return res
  }

  async function updatePointCloud(id: string, data: Partial<PointCloud>) {
    const res = await pointCloudApi.update(id, data)
    const idx = pointCloudList.value.findIndex(p => p.id === id)
    if (idx !== -1) pointCloudList.value[idx] = res
    if (currentPointCloud.value?.id === id) currentPointCloud.value = res
    return res
  }

  async function deletePointCloud(id: string) {
    await pointCloudApi.remove(id)
    pointCloudList.value = pointCloudList.value.filter(p => p.id !== id)
    if (currentPointCloud.value?.id === id) {
      currentPointCloud.value = pointCloudList.value[0] || null
    }
  }

  return {
    currentPointCloud,
    pointCloudList,
    loading,
    total,
    page,
    pageSize,
    searchKeyword,
    filteredList,
    mockPoints,
    fetchList,
    selectPointCloud,
    createPointCloud,
    updatePointCloud,
    deletePointCloud
  }
})
