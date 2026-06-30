<template>
  <div class="dz" :class="{ active: dragging, filled: !!modelValue, uploading }">
    <input
      ref="fileInput"
      class="sr-only"
      type="file"
      :accept="acceptTypes"
      @change="onPick"
    />

    <div
      class="dz-zone"
      role="button"
      tabindex="0"
      :aria-label="modelValue ? 'Replace media' : 'Choose or drop media'"
      @click="fileInput?.click()"
      @keydown.enter.prevent="fileInput?.click()"
      @keydown.space.prevent="fileInput?.click()"
      @dragover.prevent="dragging = true"
      @dragleave.prevent="dragging = false"
      @drop.prevent="onDrop"
    >
      <div v-if="modelValue" class="dz-preview">
        <video v-if="mediaType === 'video'" :src="resolvedUrl" muted playsinline />
        <img v-else :src="resolvedUrl" :alt="label" />
      </div>
      <div v-else class="dz-empty">
        <i class="icofont-upload-alt" aria-hidden="true"></i>
      </div>

      <div class="dz-meta">
        <strong>{{ statusTitle }}</strong>
        <span>{{ uploadStatus || hint }}</span>
      </div>
    </div>

    <div class="dz-actions">
      <button type="button" class="dz-choose" @click="fileInput?.click()">{{ modelValue ? 'Replace' : 'Choose file' }}</button>
      <button v-if="modelValue" type="button" class="dz-remove" @click="clear" aria-label="Remove media">
        <i class="icofont-trash" aria-hidden="true"></i>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { uploadMedia } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const props = defineProps({
  modelValue: { type: String, default: '' },
  mediaType: { type: String, default: 'image' },
  label: { type: String, default: 'media' },
  hint: { type: String, default: 'Drag a file here, or click to browse.' },
  accept: { type: String, default: 'image/png,image/jpeg,image/webp' },
  allowVideo: { type: Boolean, default: false },
  videoAccept: { type: String, default: 'video/mp4,video/webm' },
  maxImageMb: { type: Number, default: 8 },
  maxVideoMb: { type: Number, default: 24 },
  minWidth: { type: Number, default: 0 },
  minHeight: { type: Number, default: 0 },
  aspectRatio: { type: Number, default: 0 },
  aspectTolerance: { type: Number, default: 0.18 },
})

const emit = defineEmits(['update:modelValue', 'update:mediaType', 'error'])

const fileInput = ref(null)
const dragging = ref(false)
const uploading = ref(false)
const uploadStatus = ref('')

const resolvedUrl = computed(() => mediaUrl(props.modelValue))
const acceptTypes = computed(() => props.allowVideo ? `${props.accept},${props.videoAccept}` : props.accept)
const statusTitle = computed(() => {
  if (uploading.value) return 'Uploading…'
  return props.modelValue ? `${props.label} selected` : `Drop ${props.label} here`
})

function onPick(event) {
  handleFile(event.target.files?.[0])
  event.target.value = ''
}

function onDrop(event) {
  dragging.value = false
  handleFile(event.dataTransfer.files?.[0])
}

function clear() {
  emit('update:modelValue', '')
  uploadStatus.value = ''
}

function readImageDimensions(file) {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => {
      URL.revokeObjectURL(img.src)
      resolve({ width: img.naturalWidth, height: img.naturalHeight })
    }
    img.onerror = reject
    img.src = URL.createObjectURL(file)
  })
}

async function handleFile(file) {
  if (!file) return
  const imageTypes = props.accept.split(',').map(t => t.trim()).filter(Boolean)
  const videoTypes = props.allowVideo ? props.videoAccept.split(',').map(t => t.trim()).filter(Boolean) : []
  const isImage = imageTypes.includes(file.type)
  const isVideo = videoTypes.includes(file.type)

  if (!isImage && !isVideo) {
    fail(`Use ${[...imageTypes, ...videoTypes].map(t => t.split('/')[1]?.toUpperCase()).join(', ')}.`)
    return
  }

  if (isImage) {
    if (file.size > props.maxImageMb * 1024 * 1024) {
      fail(`Image must be ${props.maxImageMb}MB or smaller.`)
      return
    }
    if (props.minWidth || props.minHeight || props.aspectRatio) {
      try {
        const dims = await readImageDimensions(file)
        if (props.minWidth && dims.width < props.minWidth) { fail(`Image must be at least ${props.minWidth}px wide.`); return }
        if (props.minHeight && dims.height < props.minHeight) { fail(`Image must be at least ${props.minHeight}px tall.`); return }
        if (props.aspectRatio && Math.abs(dims.width / dims.height - props.aspectRatio) > props.aspectTolerance) {
          fail('Image aspect ratio does not match the required shape.')
          return
        }
      } catch {
        fail('Could not read image dimensions.')
        return
      }
    }
  } else {
    if (file.size > props.maxVideoMb * 1024 * 1024) {
      fail(`Video must be ${props.maxVideoMb}MB or smaller.`)
      return
    }
  }

  uploading.value = true
  uploadStatus.value = 'Uploading…'
  try {
    const res = await uploadMedia(file)
    const url = res.data?.data?.url || res.data?.url || ''
    emit('update:modelValue', url)
    if (props.allowVideo) emit('update:mediaType', isVideo ? 'video' : 'image')
    uploadStatus.value = 'Upload complete'
  } catch (error) {
    fail(error.response?.data?.error?.message || error.message || 'Upload failed')
  } finally {
    uploading.value = false
  }
}

function fail(message) {
  uploadStatus.value = message
  emit('error', new Error(message))
}
</script>

<style scoped>
.dz {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: .85rem;
  width: 100%;
  min-width: 0;
  box-sizing: border-box;
  border: 2px dashed #cbd5e1;
  border-radius: .5rem;
  background: #f8fafc;
  padding: .85rem;
  transition: border-color 150ms ease, background 150ms ease;
}
.dz.active,
.dz.filled { border-color: #f5a623; background: #fff7ed; }
.dz.uploading { opacity: .75; }

.dz-zone {
  display: flex;
  align-items: center;
  gap: .85rem;
  flex: 1 1 16rem;
  min-width: 0;
  cursor: pointer;
}
.dz-zone:focus-visible { outline: 2px solid #1a365d; outline-offset: 2px; border-radius: .35rem; }

.dz-preview {
  flex: 0 0 auto;
  width: 4.5rem;
  height: 4.5rem;
  border-radius: .4rem;
  overflow: hidden;
  background: #0f172a;
  display: grid;
  place-items: center;
}
.dz-preview img,
.dz-preview video { width: 100%; height: 100%; object-fit: cover; }

.dz-empty {
  flex: 0 0 auto;
  width: 4.5rem;
  height: 4.5rem;
  border-radius: .4rem;
  background: white;
  border: 1px solid #e2e8f0;
  display: grid;
  place-items: center;
  font-size: 1.4rem;
  color: #94a3b8;
}

.dz-meta { display: block; min-width: 0; word-break: break-word; }
.dz-meta strong { display: block; color: #1a365d; }
.dz-meta span { display: block; color: #64748b; font-size: .82rem; }

.dz-actions {
  display: flex;
  align-items: center;
  gap: .5rem;
  flex: 0 0 auto;
  margin-left: auto;
}
.dz-choose,
.dz-remove {
  border: 1px solid #cbd5e1;
  border-radius: .4rem;
  padding: .55rem .8rem;
  font-weight: 800;
  background: white;
  white-space: nowrap;
}
.dz-remove { color: #b91c1c; border-color: #fecaca; display: grid; place-items: center; padding: .55rem; }

@media (max-width: 30rem) {
  .dz { flex-direction: column; align-items: stretch; }
  .dz-zone { flex-direction: column; text-align: center; }
  .dz-actions { margin-left: 0; justify-content: stretch; }
  .dz-choose { flex: 1; }
}

:global(.dark) .dz { background: #0f172a; border-color: #334155; }
:global(.dark) .dz.active,
:global(.dark) .dz.filled { background: #422006; border-color: #f5a623; }
:global(.dark) .dz-empty { background: #111827; border-color: #334155; }
:global(.dark) .dz-meta strong { color: #f8fafc; }
:global(.dark) .dz-choose,
:global(.dark) .dz-remove { background: #111827; border-color: #475569; color: #e5e7eb; }
</style>
