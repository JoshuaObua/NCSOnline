<template>
  <div class="otika-upload">
    <input
      ref="fileInput"
      class="sr-only"
      type="file"
      :accept="acceptTypes"
      @change="onPick"
    />

    <div
      class="dropzone dz-clickable"
      :class="{ 'dz-drag-hover': dragging, 'dz-started': !!modelValue, 'dz-uploading': uploading }"
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
      <div class="dz-message" data-dz-message>
        <div v-if="modelValue" class="dz-preview dz-file-preview dz-processing dz-success dz-complete">
          <div class="dz-image">
            <video v-if="mediaType === 'video'" :src="resolvedUrl" muted playsinline />
            <div v-else-if="mediaType === 'file'" class="dz-document-preview">
              <i class="fas fa-file-alt" aria-hidden="true"></i>
            </div>
            <img v-else :src="resolvedUrl" :alt="label" />
          </div>
          <div class="dz-details">
            <div class="dz-filename"><span>{{ statusTitle }}</span></div>
            <div class="dz-size"><strong>{{ uploadStatus || 'Ready' }}</strong></div>
          </div>
        </div>
        <div v-else class="otika-drop-empty">
          <i class="fas fa-cloud-upload-alt" aria-hidden="true"></i>
          <h6>{{ statusTitle }}</h6>
          <span>{{ uploadStatus || hint }}</span>
        </div>
      </div>
    </div>

    <div class="upload-actions">
      <button type="button" class="btn btn-primary btn-sm mr-2" @click="fileInput?.click()">
        {{ modelValue ? 'Replace' : 'Choose file' }}
      </button>
      <button v-if="modelValue" type="button" class="btn btn-danger btn-sm" @click="clear" aria-label="Remove media">
        <i class="fas fa-trash" aria-hidden="true"></i>
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
  if (uploading.value) return 'Uploading...'
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
  } else if (file.size > props.maxVideoMb * 1024 * 1024) {
    fail(`Video must be ${props.maxVideoMb}MB or smaller.`)
    return
  }

  uploading.value = true
  uploadStatus.value = 'Uploading...'
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
.otika-upload {
  display: grid;
  gap: 12px;
  width: 100%;
}

.dropzone {
  min-height: 150px;
  border: 2px dashed #6777ef;
  background: #fff;
  border-radius: 3px;
  padding: 20px;
  cursor: pointer;
  transition: border-color 150ms ease, background 150ms ease;
}

.dropzone.dz-drag-hover,
.dropzone.dz-started {
  border-color: #ffa426;
  background: #fffdf7;
}

.dropzone.dz-uploading {
  opacity: .75;
}

.dropzone:focus-visible {
  outline: 2px solid #6777ef;
  outline-offset: 2px;
}

.dropzone .dz-message {
  margin: 0;
  text-align: center;
}

.otika-drop-empty {
  display: grid;
  justify-items: center;
  gap: 8px;
  padding: 14px 10px;
  color: #6c757d;
}

.otika-drop-empty i {
  color: #6777ef;
  font-size: 42px;
}

.otika-drop-empty h6 {
  margin: 0;
  color: #34395e;
  font-size: 14px;
  font-weight: 700;
}

.otika-drop-empty span {
  max-width: 32rem;
  font-size: 12px;
  line-height: 1.45;
}

.dz-preview {
  position: relative;
  display: inline-flex;
  flex-direction: column;
  width: 180px;
  min-height: 160px;
  margin: 0;
  vertical-align: top;
}

.dz-image {
  width: 180px;
  height: 120px;
  overflow: hidden;
  border-radius: 3px;
  background: #f2f2f2;
}

.dz-image img,
.dz-image video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.dz-document-preview {
  display: grid;
  width: 100%;
  height: 100%;
  place-items: center;
  color: #6777ef;
  font-size: 36px;
}

.dz-details {
  color: #34395e;
  font-size: 12px;
  padding: 10px 6px 0;
  text-align: center;
}

.dz-filename,
.dz-size {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.upload-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

@media (max-width: 420px) {
  .dropzone {
    padding: 14px;
  }

  .dz-preview,
  .dz-image {
    width: 100%;
  }
}

:global(.dark) .dropzone {
  background: #1f2937;
  border-color: #6777ef;
}

:global(.dark) .dropzone.dz-drag-hover,
:global(.dark) .dropzone.dz-started {
  background: #2b240f;
  border-color: #ffa426;
}

:global(.dark) .otika-drop-empty h6,
:global(.dark) .dz-details {
  color: #f8fafc;
}

:global(.dark) .otika-drop-empty {
  color: #cbd5e1;
}

:global(.dark) .dz-image {
  background: #111827;
}
</style>
