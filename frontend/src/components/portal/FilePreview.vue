<template>
  <article v-if="url" class="file-preview">
    <a class="file-thumb" :href="resolvedUrl" target="_blank" rel="noopener" :title="`Preview ${displayName}`">
      <img v-if="isImage" :src="resolvedUrl" :alt="displayName" loading="lazy" />
      <i v-else :class="fileIcon"></i>
    </a>
    <div>
      <strong>{{ label || displayName }}</strong>
      <small>{{ fileType }}</small>
      <span class="file-actions">
        <a :href="resolvedUrl" target="_blank" rel="noopener"><i class="icofont-eye-alt"></i> Preview</a>
        <a :href="resolvedUrl" :download="displayName"><i class="icofont-download"></i> Download</a>
      </span>
    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue'
import { mediaUrl } from '@/api/client.js'

const props = defineProps({
  url: { type: String, default: '' },
  label: { type: String, default: '' },
})

const resolvedUrl = computed(() => mediaUrl(props.url))
const cleanUrl = computed(() => String(props.url || '').split('?')[0].split('#')[0])
const extension = computed(() => {
  const match = cleanUrl.value.match(/\.([a-z0-9]+)$/i)
  return match ? match[1].toLowerCase() : ''
})
const displayName = computed(() => decodeURIComponent(cleanUrl.value.split('/').pop() || 'Uploaded file'))
const isImage = computed(() => ['png', 'jpg', 'jpeg', 'webp', 'gif'].includes(extension.value))
const fileType = computed(() => extension.value ? extension.value.toUpperCase() : 'Attachment')
const fileIcon = computed(() => extension.value === 'pdf' ? 'icofont-file-pdf' : 'icofont-file-document')
</script>

<style scoped>
.file-preview{display:grid;grid-template-columns:70px minmax(0,1fr);gap:12px;align-items:center;min-width:0;padding:10px;border:1px solid #e4e6fc;border-radius:6px;background:#fdfdff}
.file-thumb{display:grid;place-items:center;width:70px;height:58px;overflow:hidden;border-radius:5px;background:#eef2ff;color:#6777ef;text-decoration:none}
.file-thumb img{width:100%;height:100%;object-fit:cover}
.file-thumb i{font-size:28px}
.file-preview div{min-width:0}
.file-preview strong,.file-preview small{display:block;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.file-preview strong{color:#34395e;font-size:12px}
.file-preview small{color:#98a6ad;font-size:10px;font-weight:800}
.file-actions{display:flex;flex-wrap:wrap;gap:8px;margin-top:5px}
.file-actions a{display:inline-flex;align-items:center;gap:4px;color:#6777ef;font-size:11px;font-weight:800;text-decoration:none}
:global(.dark) .file-preview{background:#111827;border-color:#334155}
:global(.dark) .file-preview strong{color:#f8fafc}
:global(.dark) .file-preview small{color:#cbd5e1}
@media(max-width:520px){.file-preview{grid-template-columns:56px minmax(0,1fr)}.file-thumb{width:56px;height:52px}}
</style>
