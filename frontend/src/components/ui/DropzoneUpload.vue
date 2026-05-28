<template>
  <div>
    <div
      class="relative border-2 border-dashed rounded-xl transition-all duration-200 overflow-hidden select-none"
      :class="[
        dragover ? 'border-primary-400 bg-primary-50' : 'border-gray-300 hover:border-primary-400 hover:bg-gray-50/60',
        uploading ? 'opacity-75 pointer-events-none' : 'cursor-pointer'
      ]"
      @click="!uploading && $refs.fileInput.click()"
      @dragover.prevent="dragover = true"
      @dragleave.prevent="dragover = false"
      @drop.prevent="onDrop"
    >
      <input ref="fileInput" type="file" :accept="accept" class="hidden" @change="onFileChange" />

      <!-- Preview: image -->
      <template v-if="modelValue && !uploading && !isPdf">
        <img :src="modelValue" class="w-full object-cover rounded-xl" :class="previewClass" @click.stop="$refs.fileInput.click()"/>
        <button type="button" @click.stop="$emit('update:modelValue', '')"
          class="absolute top-2 right-2 bg-white/90 hover:bg-white shadow rounded-full w-7 h-7 flex items-center justify-center text-gray-500 hover:text-red-500 transition-colors z-10">
          <i class="icofont-close text-sm leading-none"></i>
        </button>
        <div class="absolute bottom-0 inset-x-0 bg-black/40 py-1.5 flex items-center justify-center gap-1 opacity-0 hover:opacity-100 transition-opacity"
          @click.stop="$refs.fileInput.click()">
          <i class="icofont-edit text-white text-sm"></i>
          <span class="text-white text-xs font-medium">Replace</span>
        </div>
      </template>

      <!-- Preview: PDF -->
      <template v-else-if="modelValue && !uploading && isPdf">
        <div class="flex items-center gap-4 px-5 py-4" @click.stop="$refs.fileInput.click()">
          <i class="icofont-file-pdf text-4xl text-red-500 flex-shrink-0"></i>
          <div class="min-w-0 flex-1">
            <p class="text-sm font-medium text-gray-800 truncate">{{ fileName }}</p>
            <p class="text-xs text-gray-400 mt-0.5">Click to replace file</p>
          </div>
          <button type="button" @click.stop="$emit('update:modelValue', '')"
            class="flex-shrink-0 bg-gray-100 hover:bg-red-100 rounded-full w-7 h-7 flex items-center justify-center text-gray-400 hover:text-red-500 transition-colors">
            <i class="icofont-close text-sm leading-none"></i>
          </button>
        </div>
      </template>

      <!-- Uploading -->
      <div v-else-if="uploading" class="flex flex-col items-center justify-center py-10 gap-3">
        <i class="icofont-spinner-alt-1 text-3xl text-primary-500 animate-spin"></i>
        <p class="text-sm font-medium text-primary-600">Uploading…</p>
      </div>

      <!-- Empty -->
      <div v-else class="flex flex-col items-center justify-center py-9 gap-2 text-center px-6">
        <i :class="['text-5xl text-gray-300', emptyIcon]"></i>
        <p class="text-sm font-medium text-gray-500 mt-1">{{ label }}</p>
        <p class="text-xs text-gray-400">
          Drag & drop or <span class="text-primary-600 font-medium underline underline-offset-2">browse files</span>
        </p>
        <p v-if="hint" class="text-xs text-gray-400 mt-0.5">{{ hint }}</p>
      </div>
    </div>

    <p v-if="uploadError" class="flex items-center gap-1 text-xs text-red-500 mt-1.5">
      <i class="icofont-warning-alt"></i> {{ uploadError }}
    </p>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import * as cmsApi from '@/api/cms.js'

const props = defineProps({
  modelValue: { type: String, default: '' },
  accept: { type: String, default: 'image/*' },
  label: { type: String, default: 'Drop image here or click to upload' },
  hint: { type: String, default: '' },
  previewClass: { type: String, default: 'h-40' }
})
const emit = defineEmits(['update:modelValue'])

const fileInput = ref(null)
const dragover = ref(false)
const uploading = ref(false)
const uploadError = ref('')

const isPdf = computed(() => {
  if (props.accept?.includes('pdf')) return true
  return props.modelValue?.toLowerCase().endsWith('.pdf')
})
const emptyIcon = computed(() => isPdf.value ? 'icofont-file-pdf' : 'icofont-image')
const fileName = computed(() => props.modelValue?.split('/').pop() || 'document.pdf')

async function upload(file) {
  if (!file) return
  uploading.value = true
  uploadError.value = ''
  dragover.value = false
  try {
    const r = await cmsApi.uploadMedia(file)
    const url = r.data.data?.url || r.data.url
    emit('update:modelValue', url)
  } catch (err) {
    uploadError.value = err.response?.data?.error?.message || err.message || 'Upload failed'
  } finally {
    uploading.value = false
  }
}

function onFileChange(e) {
  const file = e.target.files[0]
  if (file) upload(file)
  e.target.value = ''
}

function onDrop(e) {
  dragover.value = false
  const file = e.dataTransfer.files[0]
  if (file) upload(file)
}
</script>
