<template>
  <div ref="containerRef" class="point-cloud-canvas"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, shallowRef } from 'vue'
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import type { Point3D, PointBounds, Annotation } from '@/types'

interface Props {
  pointCloudId: string
  pointBounds?: PointBounds
  points?: Point3D[]
}

const props = withDefaults(defineProps<Props>(), {
  pointBounds: undefined,
  points: () => []
})

const emit = defineEmits<{
  'point-click': [point: Point3D]
  'polygon-created': [polygon: Point3D[]]
  'canvas-ready': []
}>()

const containerRef = ref<HTMLDivElement | null>(null)

const scene = shallowRef<THREE.Scene | null>(null)
const camera = shallowRef<THREE.PerspectiveCamera | null>(null)
const renderer = shallowRef<THREE.WebGLRenderer | null>(null)
const controls = shallowRef<OrbitControls | null>(null)
const pointsMesh = shallowRef<THREE.Points | null>(null)
const clock = shallowRef<THREE.Clock | null>(null)
const raycaster = shallowRef<THREE.Raycaster | null>(null)
const mouse = shallowRef<THREE.Vector2 | null>(null)
const gridHelper = shallowRef<THREE.GridHelper | null>(null)
const axesHelper = shallowRef<THREE.AxesHelper | null>(null)
const annotationGroup = shallowRef<THREE.Group | null>(null)
const drawingGroup = shallowRef<THREE.Group | null>(null)

const animationFrameId = ref<number | null>(null)
const isPolygonDrawing = ref(false)
const currentPolygonPoints = ref<Point3D[]>([])
const currentPolygonLine = shallowRef<THREE.LineLoop | null>(null)
const currentPolygonVertices = ref<THREE.Mesh[]>([])
const currentPolygonFill = shallowRef<THREE.Mesh | null>(null)
const loadedChunks = ref<Set<number>>(new Set())

const CHUNK_SIZE = 100000

const initScene = () => {
  if (!containerRef.value) return

  scene.value = new THREE.Scene()
  scene.value.background = new THREE.Color(0x0a1628)

  const ambientLight = new THREE.AmbientLight(0xffffff, 0.6)
  scene.value.add(ambientLight)

  const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8)
  directionalLight.position.set(50, 100, 50)
  scene.value.add(directionalLight)

  const width = containerRef.value.clientWidth
  const height = containerRef.value.clientHeight

  camera.value = new THREE.PerspectiveCamera(60, width / height, 0.1, 100000)
  camera.value.position.set(0, 50, 100)

  renderer.value = new THREE.WebGLRenderer({
    antialias: true,
    alpha: true
  })
  renderer.value.setPixelRatio(window.devicePixelRatio)
  renderer.value.setSize(width, height)
  containerRef.value.appendChild(renderer.value.domElement)

  controls.value = new OrbitControls(camera.value, renderer.value.domElement)
  controls.value.enableDamping = true
  controls.value.dampingFactor = 0.05
  controls.value.screenSpacePanning = true
  controls.value.touches = {
    ONE: THREE.TOUCH.ROTATE,
    TWO: THREE.TOUCH.DOLLY_PAN
  }

  clock.value = new THREE.Clock()
  raycaster.value = new THREE.Raycaster()
  mouse.value = new THREE.Vector2()

  annotationGroup.value = new THREE.Group()
  scene.value.add(annotationGroup.value)

  drawingGroup.value = new THREE.Group()
  scene.value.add(drawingGroup.value)

  setupHelpers()
  setupEventListeners()
  startAnimationLoop()

  emit('canvas-ready')
}

const setupHelpers = () => {
  if (!scene.value) return

  const bounds = props.pointBounds || {
    minX: -100, maxX: 100,
    minY: -100, maxY: 100,
    minZ: -50, maxZ: 50
  }

  const width = Math.abs(bounds.maxX - bounds.minX)
  const depth = Math.abs(bounds.maxY - bounds.minY)
  const size = Math.max(width, depth, 100)
  const divisions = Math.max(10, Math.floor(size / 10))

  gridHelper.value = new THREE.GridHelper(size, divisions, 0x1a3a5c, 0x0f2840)
  gridHelper.value.position.y = bounds.minZ
  gridHelper.value.rotation.x = Math.PI / 2
  scene.value.add(gridHelper.value)

  axesHelper.value = new THREE.AxesHelper(size / 5)
  axesHelper.value.position.set(bounds.minX, bounds.minY, bounds.minZ)
  scene.value.add(axesHelper.value)

  const centerX = (bounds.minX + bounds.maxX) / 2
  const centerY = (bounds.minY + bounds.maxY) / 2
  const centerZ = (bounds.minZ + bounds.maxZ) / 2
  const maxDim = Math.max(width, depth, Math.abs(bounds.maxZ - bounds.minZ))

  if (camera.value) {
    camera.value.position.set(centerX + maxDim * 0.8, centerY + maxDim * 0.8, centerZ + maxDim * 1.2)
    camera.value.lookAt(centerX, centerY, centerZ)
  }
  if (controls.value) {
    controls.value.target.set(centerX, centerY, centerZ)
  }
}

const setupEventListeners = () => {
  if (!renderer.value || !containerRef.value) return

  const canvas = renderer.value.domElement

  const onResize = () => {
    if (!containerRef.value || !camera.value || !renderer.value) return
    const width = containerRef.value.clientWidth
    const height = containerRef.value.clientHeight
    camera.value.aspect = width / height
    camera.value.updateProjectionMatrix()
    renderer.value.setSize(width, height)
  }

  const onCanvasClick = (event: MouseEvent) => {
    if (!containerRef.value) return
    const rect = containerRef.value.getBoundingClientRect()
    const screenX = event.clientX - rect.left
    const screenY = event.clientY - rect.top
    const point = screenToWorld(screenX, screenY)
    if (point) {
      emit('point-click', point)
      if (isPolygonDrawing.value) {
        addPolygonPoint(point)
      }
    }
  }

  const onTouchEnd = (event: TouchEvent) => {
    if (!containerRef.value || event.touches.length !== 1) return
    const rect = containerRef.value.getBoundingClientRect()
    const touch = event.touches[0]
    const screenX = touch.clientX - rect.left
    const screenY = touch.clientY - rect.top
    const point = screenToWorld(screenX, screenY)
    if (point) {
      emit('point-click', point)
      if (isPolygonDrawing.value) {
        addPolygonPoint(point)
      }
    }
  }

  window.addEventListener('resize', onResize)
  canvas.addEventListener('click', onCanvasClick)
  canvas.addEventListener('touchend', onTouchEnd)

  onUnmounted(() => {
    window.removeEventListener('resize', onResize)
    canvas.removeEventListener('click', onCanvasClick)
    canvas.removeEventListener('touchend', onTouchEnd)
  })
}

const startAnimationLoop = () => {
  const animate = () => {
    animationFrameId.value = requestAnimationFrame(animate)
    if (clock.value && controls.value) {
      clock.value.getDelta()
      controls.value.update()
    }
    if (renderer.value && scene.value && camera.value) {
      renderer.value.render(scene.value, camera.value)
    }
  }
  animate()
}

const computeDepthColor = (z: number, minZ: number, maxZ: number): THREE.Color => {
  const range = maxZ - minZ || 1
  const t = Math.max(0, Math.min(1, (z - minZ) / range))
  const shallow = new THREE.Color(0x87ceeb)
  const deep = new THREE.Color(0x00008b)
  return shallow.clone().lerp(deep, 1 - t)
}

const loadPointsInChunks = (pointsData: Point3D[]) => {
  if (!scene.value) return
  clearPointCloud()

  const totalChunks = Math.ceil(pointsData.length / CHUNK_SIZE)
  const bounds = props.pointBounds || {
    minX: -100, maxX: 100,
    minY: -100, maxY: 100,
    minZ: -50, maxZ: 50
  }

  let minZ = bounds.minZ
  let maxZ = bounds.maxZ
  if (pointsData.length > 0) {
    minZ = Infinity
    maxZ = -Infinity
    for (const p of pointsData) {
      if (p.z < minZ) minZ = p.z
      if (p.z > maxZ) maxZ = p.z
    }
  }

  const totalCount = pointsData.length
  const positions = new Float32Array(totalCount * 3)
  const colors = new Float32Array(totalCount * 3)

  for (let i = 0; i < totalCount; i++) {
    const p = pointsData[i]
    positions[i * 3] = p.x
    positions[i * 3 + 1] = p.y
    positions[i * 3 + 2] = p.z

    const color = computeDepthColor(p.z, minZ, maxZ)
    colors[i * 3] = color.r
    colors[i * 3 + 1] = color.g
    colors[i * 3 + 2] = color.b
  }

  const geometry = new THREE.BufferGeometry()
  geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3))
  geometry.setAttribute('color', new THREE.BufferAttribute(colors, 3))
  geometry.setDrawRange(0, totalCount)
  geometry.computeBoundingSphere()

  const material = new THREE.PointsMaterial({
    size: 1.5,
    sizeAttenuation: true,
    vertexColors: true,
    transparent: true,
    opacity: 0.9
  })

  pointsMesh.value = new THREE.Points(geometry, material)
  pointsMesh.value.frustumCulled = false
  scene.value.add(pointsMesh.value)

  loadedChunks.value.clear()
  for (let i = 0; i < totalChunks; i++) {
    loadedChunks.value.add(i)
  }
}

const clearPointCloud = () => {
  if (pointsMesh.value && scene.value) {
    scene.value.remove(pointsMesh.value)
    if (pointsMesh.value.geometry) {
      pointsMesh.value.geometry.dispose()
    }
    if (pointsMesh.value.material) {
      const mat = pointsMesh.value.material as THREE.Material
      mat.dispose()
    }
    pointsMesh.value = null
  }
  loadedChunks.value.clear()
}

const screenToWorld = (screenX: number, screenY: number): Point3D | null => {
  if (!containerRef.value || !camera.value || !raycaster.value || !mouse.value) {
    return null
  }

  const width = containerRef.value.clientWidth
  const height = containerRef.value.clientHeight

  mouse.value.x = (screenX / width) * 2 - 1
  mouse.value.y = -(screenY / height) * 2 + 1

  raycaster.value.setFromCamera(mouse.value, camera.value)

  if (pointsMesh.value) {
    const intersects = raycaster.value.intersectObject(pointsMesh.value)
    if (intersects.length > 0) {
      const hit = intersects[0]
      return {
        x: hit.point.x,
        y: hit.point.y,
        z: hit.point.z
      }
    }
  }

  const bounds = props.pointBounds || {
    minX: -100, maxX: 100,
    minY: -100, maxY: 100,
    minZ: -50, maxZ: 50
  }
  const planeZ = bounds.minZ
  const plane = new THREE.Plane(new THREE.Vector3(0, 0, 1), -planeZ)
  const target = new THREE.Vector3()
  raycaster.value.ray.intersectPlane(plane, target)
  if (target) {
    return {
      x: target.x,
      y: target.y,
      z: target.z
    }
  }

  return null
}

const startPolygonDrawing = () => {
  isPolygonDrawing.value = true
  currentPolygonPoints.value = []
  clearCurrentDrawing()
}

const addPolygonPoint = (p: Point3D) => {
  if (!isPolygonDrawing.value || !drawingGroup.value) return

  currentPolygonPoints.value.push(p)
  updateDrawingPolygon()
}

const cancelPolygonDrawing = () => {
  isPolygonDrawing.value = false
  currentPolygonPoints.value = []
  clearCurrentDrawing()
}

const finishPolygonDrawing = () => {
  if (currentPolygonPoints.value.length >= 3) {
    emit('polygon-created', [...currentPolygonPoints.value])
  }
  isPolygonDrawing.value = false
  currentPolygonPoints.value = []
  clearCurrentDrawing()
}

const clearCurrentDrawing = () => {
  if (!drawingGroup.value) return

  if (currentPolygonLine.value) {
    drawingGroup.value.remove(currentPolygonLine.value)
    if (currentPolygonLine.value.geometry) currentPolygonLine.value.geometry.dispose()
    if (currentPolygonLine.value.material) {
      const mat = currentPolygonLine.value.material as THREE.Material
      mat.dispose()
    }
    currentPolygonLine.value = null
  }

  for (const v of currentPolygonVertices.value) {
    drawingGroup.value.remove(v)
    if (v.geometry) v.geometry.dispose()
    if (v.material) {
      const mat = v.material as THREE.Material
      mat.dispose()
    }
  }
  currentPolygonVertices.value = []

  if (currentPolygonFill.value) {
    drawingGroup.value.remove(currentPolygonFill.value)
    if (currentPolygonFill.value.geometry) currentPolygonFill.value.geometry.dispose()
    if (currentPolygonFill.value.material) {
      const mat = currentPolygonFill.value.material as THREE.Material
      mat.dispose()
    }
    currentPolygonFill.value = null
  }
}

const updateDrawingPolygon = () => {
  if (!drawingGroup.value) return
  clearCurrentDrawing()

  const points = currentPolygonPoints.value
  if (points.length < 1) return

  const sphereGeo = new THREE.SphereGeometry(0.8, 16, 16)
  for (const p of points) {
    const mat = new THREE.MeshBasicMaterial({ color: 0x00ff00 })
    const sphere = new THREE.Mesh(sphereGeo, mat)
    sphere.position.set(p.x, p.y, p.z)
    drawingGroup.value.add(sphere)
    currentPolygonVertices.value.push(sphere)
  }

  if (points.length >= 2) {
    const positions = new Float32Array(points.length * 3)
    for (let i = 0; i < points.length; i++) {
      positions[i * 3] = points[i].x
      positions[i * 3 + 1] = points[i].y
      positions[i * 3 + 2] = points[i].z
    }
    const geo = new THREE.BufferGeometry()
    geo.setAttribute('position', new THREE.BufferAttribute(positions, 3))
    const mat = new THREE.LineBasicMaterial({ color: 0x00ff00, linewidth: 2 })
    const line = new THREE.LineLoop(geo, mat)
    drawingGroup.value.add(line)
    currentPolygonLine.value = line

    if (points.length >= 3) {
      const n = points.length
      const triCount = n - 2
      const fillPositions = new Float32Array(triCount * 9)
      const fillIndices = new Uint32Array(triCount * 3)

      let vi = 0
      let ii = 0
      for (let t = 0; t < triCount; t++) {
        const i0 = 0
        const i1 = t + 1
        const i2 = t + 2
        const p0 = points[i0]
        const p1 = points[i1]
        const p2 = points[i2]
        fillPositions[vi++] = p0.x; fillPositions[vi++] = p0.y; fillPositions[vi++] = p0.z
        fillPositions[vi++] = p1.x; fillPositions[vi++] = p1.y; fillPositions[vi++] = p1.z
        fillPositions[vi++] = p2.x; fillPositions[vi++] = p2.y; fillPositions[vi++] = p2.z
        fillIndices[ii++] = t * 3
        fillIndices[ii++] = t * 3 + 1
        fillIndices[ii++] = t * 3 + 2
      }

      const fillGeo = new THREE.BufferGeometry()
      fillGeo.setAttribute('position', new THREE.BufferAttribute(fillPositions, 3))
      fillGeo.setIndex(new THREE.BufferAttribute(fillIndices, 1))
      fillGeo.computeVertexNormals()
      const fillMat = new THREE.MeshBasicMaterial({
        color: 0x00ff00,
        transparent: true,
        opacity: 0.2,
        side: THREE.DoubleSide
      })
      const fillMesh = new THREE.Mesh(fillGeo, fillMat)
      drawingGroup.value.add(fillMesh)
      currentPolygonFill.value = fillMesh
    }
  }
}

const addAnnotationPolygon = (annotation: Annotation) => {
  if (!annotationGroup.value) return

  const points = annotation.polygon
  if (!points || points.length < 3) return

  const color = annotation.label?.color || '#ff6600'
  const threeColor = new THREE.Color(color)

  const positions = new Float32Array(points.length * 3)
  for (let i = 0; i < points.length; i++) {
    positions[i * 3] = points[i].x
    positions[i * 3 + 1] = points[i].y
    positions[i * 3 + 2] = points[i].z
  }

  const lineGeo = new THREE.BufferGeometry()
  lineGeo.setAttribute('position', new THREE.BufferAttribute(positions, 3))
  const lineMat = new THREE.LineBasicMaterial({ color: threeColor, linewidth: 2 })
  const lineLoop = new THREE.LineLoop(lineGeo, lineMat)
  lineLoop.userData = { annotationId: annotation.id }
  annotationGroup.value.add(lineLoop)

  const sphereGeo = new THREE.SphereGeometry(0.5, 8, 8)
  for (const p of points) {
    const sphereMat = new THREE.MeshBasicMaterial({ color: threeColor })
    const sphere = new THREE.Mesh(sphereGeo, sphereMat)
    sphere.position.set(p.x, p.y, p.z)
    sphere.userData = { annotationId: annotation.id }
    annotationGroup.value.add(sphere)
  }

  const n = points.length
  const triCount = n - 2
  const fillPositions = new Float32Array(triCount * 9)
  const fillIndices = new Uint32Array(triCount * 3)

  let vi = 0
  let ii = 0
  for (let t = 0; t < triCount; t++) {
    const i0 = 0
    const i1 = t + 1
    const i2 = t + 2
    const p0 = points[i0]
    const p1 = points[i1]
    const p2 = points[i2]
    fillPositions[vi++] = p0.x; fillPositions[vi++] = p0.y; fillPositions[vi++] = p0.z
    fillPositions[vi++] = p1.x; fillPositions[vi++] = p1.y; fillPositions[vi++] = p1.z
    fillPositions[vi++] = p2.x; fillPositions[vi++] = p2.y; fillPositions[vi++] = p2.z
    fillIndices[ii++] = t * 3
    fillIndices[ii++] = t * 3 + 1
    fillIndices[ii++] = t * 3 + 2
  }

  const fillGeo = new THREE.BufferGeometry()
  fillGeo.setAttribute('position', new THREE.BufferAttribute(fillPositions, 3))
  fillGeo.setIndex(new THREE.BufferAttribute(fillIndices, 1))
  fillGeo.computeVertexNormals()
  const fillMat = new THREE.MeshBasicMaterial({
    color: threeColor,
    transparent: true,
    opacity: 0.25,
    side: THREE.DoubleSide
  })
  const fillMesh = new THREE.Mesh(fillGeo, fillMat)
  fillMesh.userData = { annotationId: annotation.id }
  annotationGroup.value.add(fillMesh)
}

const clearAnnotations = () => {
  if (!annotationGroup.value) return

  while (annotationGroup.value.children.length > 0) {
    const child = annotationGroup.value.children[0]
    annotationGroup.value.remove(child)
    if ((child as THREE.Mesh).geometry) {
      (child as THREE.Mesh).geometry.dispose()
    }
    if ((child as THREE.Mesh).material) {
      const mat = (child as THREE.Mesh).material as THREE.Material | THREE.Material[]
      if (Array.isArray(mat)) {
        mat.forEach(m => m.dispose())
      } else {
        mat.dispose()
      }
    }
  }
}

const cleanup = () => {
  if (animationFrameId.value !== null) {
    cancelAnimationFrame(animationFrameId.value)
    animationFrameId.value = null
  }

  clearCurrentDrawing()
  clearAnnotations()
  clearPointCloud()

  if (gridHelper.value && scene.value) {
    scene.value.remove(gridHelper.value)
    gridHelper.value.geometry.dispose()
    gridHelper.value = null
  }

  if (axesHelper.value && scene.value) {
    scene.value.remove(axesHelper.value)
    axesHelper.value.geometry.dispose()
    axesHelper.value = null
  }

  if (annotationGroup.value && scene.value) {
    scene.value.remove(annotationGroup.value)
    annotationGroup.value = null
  }

  if (drawingGroup.value && scene.value) {
    scene.value.remove(drawingGroup.value)
    drawingGroup.value = null
  }

  if (controls.value) {
    controls.value.dispose()
    controls.value = null
  }

  if (renderer.value && containerRef.value) {
    containerRef.value.removeChild(renderer.value.domElement)
    renderer.value.dispose()
    renderer.value = null
  }

  camera.value = null
  scene.value = null
  clock.value = null
  raycaster.value = null
  mouse.value = null
}

watch(() => props.points, (newPoints) => {
  if (newPoints && newPoints.length > 0) {
    loadPointsInChunks(newPoints)
  }
}, { immediate: true })

watch(() => props.pointBounds, () => {
  if (scene.value && gridHelper.value && axesHelper.value) {
    scene.value.remove(gridHelper.value)
    scene.value.remove(axesHelper.value)
    gridHelper.value.geometry.dispose()
    axesHelper.value.geometry.dispose()
    setupHelpers()
  }
})

onMounted(() => {
  initScene()
  if (props.points && props.points.length > 0) {
    loadPointsInChunks(props.points)
  }
})

onUnmounted(() => {
  cleanup()
})

defineExpose({
  screenToWorld,
  startPolygonDrawing,
  addPolygonPoint,
  cancelPolygonDrawing,
  finishPolygonDrawing,
  clearAnnotations,
  addAnnotationPolygon
})
</script>

<style lang="scss" scoped>
.point-cloud-canvas {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
}
</style>
