<template>
  <div class="dropzone-document-upload" :class="{ 'is-disabled': disabled }">
    <input
      :id="id"
      ref="fileInputRef"
      type="file"
      class="hidden-file-input"
      :multiple="multiple"
      :accept="accept"
      :disabled="disabled"
      @change="handleFileSelect"
    />

    <!-- Dropzone Area -->
    <div
      class="dropzone-area"
      :class="{
        'is-dragging': isDragging,
        'is-uploading': isUploading,
        'has-files': fileList.length > 0
      }"
      role="button"
      tabindex="0"
      @click="triggerBrowse"
      @keydown.enter.prevent="triggerBrowse"
      @keydown.space.prevent="triggerBrowse"
      @dragover.prevent="onDragOver"
      @dragleave.prevent="onDragLeave"
      @drop.prevent="onDrop"
    >
      <div class="dropzone-content">
        <div class="cloud-icon-wrapper">
          <i v-if="isUploading" class="icofont-spinner icofont-spin"></i>
          <i v-else-if="isDragging" class="icofont-download"></i>
          <i v-else class="icofont-cloud-upload"></i>
        </div>

        <div class="dropzone-text">
          <h5 v-if="isUploading">Uploading documents...</h5>
          <h5 v-else-if="isDragging">Drop files here to upload</h5>
          <h5 v-else-if="fileList.length > 0">Drag & drop more documents or click to browse</h5>
          <h5 v-else>Drag & drop test result documents here, or click to browse</h5>

          <p class="dropzone-hint">
            Supports PDF, Medical Reports, Images, DOCX (Max 25MB per file)
          </p>
        </div>

        <button
          type="button"
          class="browse-btn"
          :disabled="disabled || isUploading"
          @click.stop="triggerBrowse"
        >
          <i class="icofont-folder-open"></i> Browse Files
        </button>
      </div>

      <!-- Uploading Progress Overlay -->
      <div v-if="isUploading" class="upload-progress-overlay">
        <div class="progress-bar-track">
          <div class="progress-bar-fill" :style="{ width: `${uploadProgress}%` }"></div>
        </div>
        <span class="progress-text">{{ uploadProgress }}% Uploading {{ currentUploadingName }}</span>
      </div>
    </div>

    <!-- Error Alert -->
    <p v-if="errorMessage" class="upload-error-alert">
      <i class="icofont-warning"></i> {{ errorMessage }}
    </p>

    <!-- Attached Files List -->
    <div v-if="fileList.length > 0" class="attached-files-container">
      <div class="files-header">
        <h6>Attached Test Documents ({{ fileList.length }})</h6>
        <button
          v-if="fileList.length > 1 && !disabled"
          type="button"
          class="clear-all-btn"
          @click="removeAllFiles"
        >
          Clear All
        </button>
      </div>

      <ul class="files-list">
        <li v-for="(file, idx) in fileList" :key="file.url || idx" class="file-item">
          <div class="file-icon" :class="getFileIconClass(file.name || file.url)">
            <i :class="getFileIcon(file.name || file.url)"></i>
          </div>

          <div class="file-info">
            <a :href="getResolvedUrl(file.url)" target="_blank" rel="noopener noreferrer" class="file-name" title="Click to view document">
              {{ file.name || getBasename(file.url) }}
            </a>
            <div class="file-meta">
              <span v-if="file.size" class="meta-tag">{{ formatSize(file.size) }}</span>
              <span v-if="getFileExt(file.name || file.url)" class="meta-tag ext-tag">{{ getFileExt(file.name || file.url).toUpperCase() }}</span>
              <span class="meta-status"><i class="icofont-check-circled"></i> Attached</span>
            </div>
          </div>

          <div class="file-actions">
            <a
              :href="getResolvedUrl(file.url)"
              target="_blank"
              rel="noopener noreferrer"
              class="action-btn view-btn"
              title="View/Download Document"
            >
              <i class="icofont-eye"></i> View
            </a>
            <button
              v-if="!disabled"
              type="button"
              class="action-btn remove-btn"
              title="Remove file"
              @click="removeFile(idx)"
            >
              <i class="icofont-trash"></i>
            </button>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { uploadMedia } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const props = defineProps({
  modelValue: {
    type: [String, Array, Object],
    default: ''
  },
  id: {
    type: String,
    default: ''
  },
  multiple: {
    type: Boolean,
    default: true
  },
  accept: {
    type: String,
    default: '.pdf,.doc,.docx,.jpg,.jpeg,.png,.webp,.txt'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  required: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const fileInputRef = ref(null)
const isDragging = ref(false)
const isUploading = ref(false)
const uploadProgress = ref(0)
const currentUploadingName = ref('')
const errorMessage = ref('')
const fileList = ref([])

// Initialize file list from modelValue
watch(
  () => props.modelValue,
  (val) => {
    fileList.value = parseModelValue(val)
  },
  { immediate: true }
)

function parseModelValue(val) {
  if (!val) return []
  if (Array.isArray(val)) {
    return val.map(normalizeFileItem).filter(Boolean)
  }
  if (typeof val === 'object') {
    return [normalizeFileItem(val)].filter(Boolean)
  }
  if (typeof val === 'string') {
    const trimmed = val.trim()
    if (!trimmed) return []
    if (trimmed.startsWith('[') || trimmed.startsWith('{')) {
      try {
        const parsed = JSON.parse(trimmed)
        return parseModelValue(parsed)
      } catch {
        // Fallback: comma separated URLs
        return trimmed.split(',').map(s => normalizeFileItem(s.trim())).filter(Boolean)
      }
    }
    // Plain URL or path string
    return trimmed.split(',').map(s => normalizeFileItem(s.trim())).filter(Boolean)
  }
  return []
}

function normalizeFileItem(item) {
  if (!item) return null
  if (typeof item === 'string') {
    return { url: item, name: getBasename(item), size: null }
  }
  if (typeof item === 'object' && item.url) {
    return {
      url: item.url,
      name: item.name || getBasename(item.url),
      size: item.size || null,
      type: item.type || ''
    }
  }
  return null
}

function emitChanges() {
  let val = ''
  if (fileList.value.length > 0) {
    val = JSON.stringify(fileList.value)
  }
  emit('update:modelValue', val)
  emit('change', val)
}

function triggerBrowse() {
  if (props.disabled || isUploading.value) return
  fileInputRef.value?.click()
}

function onDragOver() {
  if (props.disabled || isUploading.value) return
  isDragging.value = true
}

function onDragLeave() {
  isDragging.value = false
}

function onDrop(event) {
  isDragging.value = false
  if (props.disabled || isUploading.value) return
  const files = event.dataTransfer?.files
  if (files && files.length > 0) {
    uploadSelectedFiles(Array.from(files))
  }
}

function handleFileSelect(event) {
  const files = event.target?.files
  if (files && files.length > 0) {
    uploadSelectedFiles(Array.from(files))
  }
  if (fileInputRef.value) {
    fileInputRef.value.value = ''
  }
}

async function uploadSelectedFiles(files) {
  errorMessage.value = ''
  if (!props.multiple && files.length > 1) {
    files = [files[0]]
  }

  isUploading.value = true
  uploadProgress.value = 0

  const total = files.length
  let completed = 0

  for (const file of files) {
    currentUploadingName.value = file.name
    try {
      const res = await uploadMedia(file, 'anti-doping')
      const data = res?.data?.data || res?.data || {}
      const url = data.url || data.file_url || data.path || ''

      if (url) {
        const fileObj = {
          url,
          name: file.name,
          size: file.size,
          type: file.type
        }
        if (!props.multiple) {
          fileList.value = [fileObj]
        } else {
          fileList.value.push(fileObj)
        }
      }
    } catch (err) {
      console.error('[DropzoneDocumentUpload] Failed to upload file:', err)
      errorMessage.value = err?.response?.data?.error?.message || `Could not upload "${file.name}". Please try again.`
    } finally {
      completed++
      uploadProgress.value = Math.round((completed / total) * 100)
    }
  }

  isUploading.value = false
  currentUploadingName.value = ''
  emitChanges()
}

function removeFile(index) {
  fileList.value.splice(index, 1)
  emitChanges()
}

function removeAllFiles() {
  fileList.value = []
  emitChanges()
}

function getBasename(pathStr) {
  if (!pathStr) return 'Document'
  const clean = String(pathStr).split('?')[0]
  const parts = clean.split('/')
  return parts[parts.length - 1] || 'Document'
}

function getFileExt(filename) {
  if (!filename) return ''
  const parts = String(filename).split('.')
  return parts.length > 1 ? parts[parts.length - 1].toLowerCase() : ''
}

function getFileIcon(filename) {
  const ext = getFileExt(filename)
  if (['pdf'].includes(ext)) return 'icofont-file-pdf'
  if (['jpg', 'jpeg', 'png', 'webp', 'gif', 'svg'].includes(ext)) return 'icofont-file-jpg'
  if (['doc', 'docx'].includes(ext)) return 'icofont-file-word'
  if (['xls', 'xlsx', 'csv'].includes(ext)) return 'icofont-file-spreadsheet'
  return 'icofont-file-alt'
}

function getFileIconClass(filename) {
  const ext = getFileExt(filename)
  if (['pdf'].includes(ext)) return 'icon-pdf'
  if (['jpg', 'jpeg', 'png', 'webp', 'gif', 'svg'].includes(ext)) return 'icon-img'
  if (['doc', 'docx'].includes(ext)) return 'icon-doc'
  return 'icon-default'
}

function getResolvedUrl(url) {
  if (!url) return '#'
  if (url.startsWith('http://') || url.startsWith('https://') || url.startsWith('blob:')) return url
  return mediaUrl(url)
}

function formatSize(bytes) {
  if (!bytes) return ''
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}
</script>

<style scoped>
.dropzone-document-upload {
  width: 100%;
  font-family: inherit;
}

.hidden-file-input {
  display: none !important;
}

.dropzone-area {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 140px;
  padding: 24px 20px;
  border: 2px dashed #cbd5e1;
  border-radius: 12px;
  background-color: #fafbfc;
  cursor: pointer;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.dropzone-area:hover:not(.is-disabled) {
  border-color: #6777ef;
  background-color: #f4f6fd;
}

.dropzone-area.is-dragging {
  border-color: #6777ef;
  background-color: #eef2ff;
  box-shadow: 0 0 0 4px rgba(103, 119, 239, 0.15);
}

.dropzone-area.has-files {
  min-height: 110px;
  padding: 16px;
  border-style: solid;
  border-color: #e2e8f0;
}

.dropzone-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 10px;
}

.cloud-icon-wrapper {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6777ef, #3b82f6);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  box-shadow: 0 4px 12px rgba(103, 119, 239, 0.25);
}

.dropzone-text h5 {
  margin: 0 0 4px;
  font-size: 15px;
  font-weight: 700;
  color: #34395e;
}

.dropzone-hint {
  margin: 0;
  font-size: 12px;
  color: #6c757d;
}

.browse-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 16px;
  border: 1px solid #6777ef;
  border-radius: 6px;
  background: #ffffff;
  color: #6777ef;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.browse-btn:hover:not(:disabled) {
  background: #6777ef;
  color: #ffffff;
}

.upload-progress-overlay {
  position: absolute;
  inset: 0;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.92);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  gap: 10px;
  z-index: 10;
}

.progress-bar-track {
  width: 80%;
  height: 8px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, #6777ef, #3b82f6);
  transition: width 0.2s ease;
}

.progress-text {
  font-size: 12.5px;
  font-weight: 700;
  color: #34395e;
}

.upload-error-alert {
  margin: 10px 0 0;
  padding: 8px 12px;
  border-radius: 6px;
  background: #fff1f0;
  color: #e74c3c;
  font-size: 12.5px;
  font-weight: 600;
  border: 1px solid #ffd4d0;
}

.attached-files-container {
  margin-top: 16px;
}

.files-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.files-header h6 {
  margin: 0;
  font-size: 13px;
  font-weight: 700;
  color: #34395e;
}

.clear-all-btn {
  background: transparent;
  border: none;
  color: #e74c3c;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

.files-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 8px;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 10px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: #ffffff;
  gap: 12px;
  transition: border-color 0.15s;
}

.file-item:hover {
  border-color: #cbd5e1;
}

.file-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}

.icon-pdf { background: #fef2f2; color: #dc2626; }
.icon-img { background: #eff6ff; color: #2563eb; }
.icon-doc { background: #f0fdf4; color: #16a34a; }
.icon-default { background: #f8fafc; color: #64748b; }

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  display: block;
  font-size: 13.5px;
  font-weight: 600;
  color: #34395e;
  text-decoration: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-name:hover {
  color: #6777ef;

}

.file-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 2px;
  font-size: 11px;
}

.meta-tag {
  color: #6c757d;
}

.ext-tag {
  font-weight: 700;
  padding: 1px 4px;
  border-radius: 3px;
  background: #f1f5f9;
  color: #475569;
}

.meta-status {
  color: #059669;
  font-weight: 600;
}

.file-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 5px 10px;
  border-radius: 5px;
  font-size: 12px;
  font-weight: 700;
  text-decoration: none;
  cursor: pointer;
  border: 1px solid transparent;
}

.view-btn {
  background: #f4f6fd;
  color: #6777ef;
  border-color: #dbe0fc;
}

.view-btn:hover {
  background: #6777ef;
  color: #ffffff;
}

.remove-btn {
  background: #fff1f0;
  color: #e74c3c;
  border-color: #ffd4d0;
}

.remove-btn:hover {
  background: #e74c3c;
  color: #ffffff;
}

/* Dark Mode Support */
:global(html.dark) .dropzone-area,
:global(body.dark) .dropzone-area,
:global([data-theme="dark"]) .dropzone-area {
  background-color: #0f172a !important;
  border-color: #334155 !important;
}

:global(html.dark) .dropzone-text h5,
:global(body.dark) .dropzone-text h5,
:global([data-theme="dark"]) .dropzone-text h5,
:global(html.dark) .file-name,
:global(body.dark) .file-name,
:global([data-theme="dark"]) .file-name,
:global(html.dark) .files-header h6,
:global(body.dark) .files-header h6,
:global([data-theme="dark"]) .files-header h6 {
  color: #f8fafc !important;
}

:global(html.dark) .file-item,
:global(body.dark) .file-item,
:global([data-theme="dark"]) .file-item {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(html.dark) .upload-progress-overlay,
:global(body.dark) .upload-progress-overlay,
:global([data-theme="dark"]) .upload-progress-overlay {
  background-color: rgba(15, 23, 42, 0.92) !important;
}

:global(html.dark) .progress-text,
:global(body.dark) .progress-text,
:global([data-theme="dark"]) .progress-text {
  color: #f8fafc !important;
}
</style>
