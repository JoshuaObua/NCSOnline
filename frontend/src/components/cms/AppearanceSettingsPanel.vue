<template>
  <section class="appearance-settings">
    <div class="cms-panel-head">
      <h2>Appearance</h2>
      <button v-if="activeTab === 'typography'" type="button" :disabled="saving" @click="saveTypography">{{ saving ? 'Saving…' : 'Save typography' }}</button>
    </div>

    <div class="appearance-tabs">
      <button type="button" :class="{ active: activeTab === 'fonts' }" @click="activeTab = 'fonts'">Font Manager</button>
      <button type="button" :class="{ active: activeTab === 'typography' }" @click="activeTab = 'typography'">Typography Mapping</button>
    </div>

    <!-- ── Tab A: Font Manager ─────────────────────────────────────── -->
    <div v-if="activeTab === 'fonts'" class="cms-panel">
      <div class="cms-panel-head"><h3>Upload a custom font</h3></div>
      <p class="storage-hint">Accepts .ttf, .otf, .woff, and .woff2. The display name becomes the CSS font-family used in Typography Mapping.</p>
      <div class="font-upload-row">
        <label class="wide">Display name<input v-model="newFontName" class="form-control" placeholder="e.g. Custom Header Bold" /></label>
        <div class="otika-upload">
          <input ref="fileInput" class="sr-only" type="file" accept=".ttf,.otf,.woff,.woff2" @change="onFilePicked" />
          <div
            class="dropzone dz-clickable"
            :class="{ 'dz-drag-hover': dragging, 'dz-started': !!newFontFile, 'dz-uploading': uploading }"
            role="button"
            tabindex="0"
            :aria-label="newFontFile ? 'Replace font file' : 'Choose or drop a font file'"
            @click="fileInput?.click()"
            @keydown.enter.prevent="fileInput?.click()"
            @keydown.space.prevent="fileInput?.click()"
            @dragover.prevent="dragging = true"
            @dragleave.prevent="dragging = false"
            @drop.prevent="onDrop"
          >
            <div class="dz-message">
              <div v-if="newFontFile" class="dz-preview">
                <i class="fas fa-font" aria-hidden="true"></i>
                <div class="dz-details">
                  <div class="dz-filename">{{ newFontFile.name }}</div>
                  <div class="dz-size">{{ formatFileSize(newFontFile.size) }}</div>
                </div>
              </div>
              <div v-else class="otika-drop-empty">
                <i class="fas fa-cloud-upload-alt" aria-hidden="true"></i>
                <h6>{{ uploading ? 'Uploading…' : 'Drop a font file here' }}</h6>
                <span>.ttf, .otf, .woff, or .woff2 — click to browse</span>
              </div>
            </div>
          </div>
          <div class="upload-actions">
            <button type="button" class="btn-secondary" @click="fileInput?.click()">{{ newFontFile ? 'Replace' : 'Choose file' }}</button>
            <button v-if="newFontFile" type="button" class="btn-danger" aria-label="Remove selected file" @click="clearPickedFile"><i class="fas fa-trash" aria-hidden="true"></i></button>
          </div>
        </div>
        <button type="button" class="upload-submit" :disabled="uploading || !newFontFile || !newFontName.trim()" @click="uploadFont">{{ uploading ? 'Uploading…' : 'Upload font' }}</button>
      </div>

      <div class="cms-panel-head"><h3>Installed fonts</h3></div>
      <table class="font-table">
        <thead><tr><th>Preview</th><th>Display name</th><th>Format</th><th></th></tr></thead>
        <tbody>
          <tr v-for="f in fonts" :key="f.id">
            <td><span class="font-preview" :style="{ fontFamily: quoted(f.font_name) }">The quick brown fox</span></td>
            <td>{{ f.display_name }}</td>
            <td><span class="font-format-badge">{{ f.font_format }}</span></td>
            <td><button type="button" class="font-delete" @click="deleteFont(f)">Delete</button></td>
          </tr>
          <tr v-if="!fonts.length"><td colspan="4" class="font-empty">No custom fonts uploaded yet.</td></tr>
        </tbody>
      </table>
    </div>

    <!-- ── Tab B: Typography Mapping ───────────────────────────────── -->
    <div v-else class="typography-mapping">
      <article v-for="group in TYPOGRAPHY_GROUPS" :key="group.label" class="cms-panel">
        <div class="cms-panel-head"><h3>{{ group.label }}</h3></div>
        <div class="typography-row" v-for="el in group.elements" :key="el.key">
          <div class="typography-row__label">{{ el.label }}</div>
          <label>Font family
            <select v-model="typography[el.key].family">
              <option value="">Site default</option>
              <optgroup label="Web-safe">
                <option v-for="name in WEB_SAFE_FONTS" :key="name" :value="name">{{ name }}</option>
              </optgroup>
              <optgroup v-if="fonts.length" label="Custom fonts">
                <option v-for="f in fonts" :key="f.id" :value="f.font_name">{{ f.display_name }}</option>
              </optgroup>
            </select>
          </label>
          <label>Size
            <span class="typography-size">
              <input v-model="typography[el.key].size" type="number" min="0" step="0.1" placeholder="16" />
              <select v-model="typography[el.key].unit">
                <option v-for="u in FONT_SIZE_UNITS" :key="u" :value="u">{{ u }}</option>
              </select>
            </span>
          </label>
          <label>Weight
            <select v-model="typography[el.key].weight">
              <option value="">Site default</option>
              <option v-for="w in FONT_WEIGHTS" :key="w.value" :value="w.value">{{ w.label }}</option>
            </select>
          </label>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { adminDeleteFont, adminListFonts, adminUpdateTypography, adminUploadFont, getTypography } from '@/api/cms.js'
import { FONT_SIZE_UNITS, FONT_WEIGHTS, TYPOGRAPHY_GROUPS, WEB_SAFE_FONTS, buildFontFaceCSS, normalizeTypography } from '@/utils/typography.js'

const emit = defineEmits(['message', 'error'])

const activeTab = ref('fonts')
const fonts = ref([])
const typography = reactive(normalizeTypography({}))
const newFontFile = ref(null)
const newFontName = ref('')
const fileInput = ref(null)
const uploading = ref(false)
const saving = ref(false)
const dragging = ref(false)

const ALLOWED_FONT_EXT = ['.ttf', '.otf', '.woff', '.woff2']

function quoted(name) {
  return /\s/.test(name) ? `'${name}'` : name
}

function formatFileSize(bytes) {
  const n = Number(bytes) || 0
  if (n <= 0) return '0 B'
  const units = ['B', 'KB', 'MB']
  const i = Math.min(units.length - 1, Math.floor(Math.log(n) / Math.log(1024)))
  return `${(n / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
}

function isFontFile(file) {
  const name = file?.name?.toLowerCase() || ''
  return ALLOWED_FONT_EXT.some(ext => name.endsWith(ext))
}

function injectFontPreviewFaces() {
  let tag = document.getElementById('cms-appearance-font-preview')
  if (!tag) {
    tag = document.createElement('style')
    tag.id = 'cms-appearance-font-preview'
    document.head.appendChild(tag)
  }
  tag.textContent = buildFontFaceCSS(fonts.value)
}

async function loadFonts() {
  try {
    fonts.value = (await adminListFonts()).data?.data || []
    injectFontPreviewFaces()
  } catch (err) { emit('error', err) }
}

async function loadTypography() {
  try {
    const value = (await getTypography()).data?.data?.value || {}
    Object.assign(typography, normalizeTypography(value))
  } catch (err) { emit('error', err) }
}

function onFilePicked(event) {
  const file = event.target.files?.[0] || null
  event.target.value = ''
  if (file && !isFontFile(file)) {
    emit('error', new Error('Use a .ttf, .otf, .woff, or .woff2 file.'))
    return
  }
  newFontFile.value = file
}

function onDrop(event) {
  dragging.value = false
  const file = event.dataTransfer.files?.[0] || null
  if (file && !isFontFile(file)) {
    emit('error', new Error('Use a .ttf, .otf, .woff, or .woff2 file.'))
    return
  }
  newFontFile.value = file
}

function clearPickedFile() {
  newFontFile.value = null
  if (fileInput.value) fileInput.value.value = ''
}

async function uploadFont() {
  if (!newFontFile.value || !newFontName.value.trim()) return
  uploading.value = true
  try {
    await adminUploadFont(newFontFile.value, newFontName.value.trim())
    newFontName.value = ''
    clearPickedFile()
    await loadFonts()
    emit('message', 'Font uploaded')
  } catch (err) {
    emit('error', err)
  } finally {
    uploading.value = false
  }
}

async function deleteFont(f) {
  try {
    await adminDeleteFont(f.id)
    await loadFonts()
    emit('message', 'Font deleted')
  } catch (err) { emit('error', err) }
}

async function saveTypography() {
  saving.value = true
  try {
    await adminUpdateTypography(JSON.parse(JSON.stringify(typography)))
    emit('message', 'Typography saved')
  } catch (err) {
    emit('error', err)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await loadFonts()
  await loadTypography()
})
</script>

<style scoped>
/*
 * This component is rendered as a child of WebsiteContentManagerView, whose
 * <style scoped> only reaches this component's root element (Vue scoped CSS
 * does not cascade into child templates) — so shared classes are redefined
 * here rather than assumed to be inherited.
 */
.appearance-settings { display: grid; gap: 20px; min-width: 0; }

.cms-panel-head { display: flex; justify-content: space-between; align-items: center; gap: 16px; flex-wrap: wrap; border-bottom: 1px solid #f9f9f9; padding-bottom: 15px; }
.cms-panel-head h2,
.cms-panel-head h3 { margin: 0; font-size: 16px; font-weight: 700; color: #34395e; }
.cms-panel-head button { border: 0; border-radius: 30px; background: #6777ef; color: white; padding: 8px 18px; font-size: 12px; font-weight: 600; white-space: nowrap; box-shadow: 0 2px 6px #acb5f6; }
.cms-panel-head button:disabled { opacity: .6; }

.cms-panel { background: white; border: 0; border-radius: 3px; box-shadow: 0 4px 25px rgba(0,0,0,.1); padding: 25px; display: grid; gap: 16px; min-width: 0; }

.appearance-tabs { display: flex; gap: 8px; }
.appearance-tabs button { border: 1px solid #e4e6fc; border-radius: 30px; background: #fff; color: #34395e; padding: 8px 18px; font-size: 12px; font-weight: 700; }
.appearance-tabs button.active { background: #6777ef; color: #fff; border-color: #6777ef; }

.storage-hint { color: #6c757d; font-size: 13px; margin: 0; line-height: 1.7; }

.font-upload-row { display: grid; grid-template-columns: 1fr 1.4fr auto; gap: 14px; align-items: start; }
.font-upload-row > label { display: grid; gap: 7px; font-size: 12px; font-weight: 600; color: #34395e; min-width: 0; }
.font-upload-row > label input { border: 1px solid #e4e6fc; border-radius: 3px; padding: 10px 15px; font-weight: 500; color: #495057; background: #fdfdff; outline: none; }
.upload-submit { border: 0; border-radius: 30px; background: #6777ef; color: #fff; padding: 10px 18px; font-size: 12px; font-weight: 700; white-space: nowrap; align-self: end; }
.upload-submit:disabled { opacity: .5; }

.otika-upload { display: grid; gap: 10px; min-width: 0; }
.sr-only { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0,0,0,0); white-space: nowrap; border: 0; }
.dropzone { min-height: 90px; border: 2px dashed #6777ef; background: #fff; border-radius: 3px; padding: 14px; cursor: pointer; transition: border-color 150ms ease, background 150ms ease; }
.dropzone.dz-drag-hover,
.dropzone.dz-started { border-color: #ffa426; background: #fffdf7; }
.dropzone.dz-uploading { opacity: .75; }
.dropzone:focus-visible { outline: 2px solid #6777ef; outline-offset: 2px; }
.dz-message { margin: 0; text-align: center; }
.otika-drop-empty { display: grid; justify-items: center; gap: 6px; padding: 6px 10px; color: #6c757d; }
.otika-drop-empty i { color: #6777ef; font-size: 30px; }
.otika-drop-empty h6 { margin: 0; color: #34395e; font-size: 13px; font-weight: 700; }
.otika-drop-empty span { font-size: 11px; }
.dz-preview { display: flex; align-items: center; justify-content: center; gap: 12px; }
.dz-preview i { color: #6777ef; font-size: 28px; }
.dz-details { text-align: left; color: #34395e; font-size: 12px; }
.dz-filename { font-weight: 700; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 16rem; }
.dz-size { color: #94a3b8; }
.upload-actions { display: flex; flex-wrap: wrap; align-items: center; gap: 8px; }
.upload-actions .btn-secondary { border: 1px solid #e4e6fc; border-radius: 30px; background: #fff; color: #34395e; padding: 6px 14px; font-size: 11px; font-weight: 700; }
.upload-actions .btn-danger { border: 0; border-radius: 30px; background: #fc544b; color: #fff; padding: 6px 10px; font-size: 11px; font-weight: 700; }

.font-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.font-table th { text-align: left; color: #94a3b8; font-size: 11px; text-transform: uppercase; letter-spacing: .04em; padding: 8px 10px; border-bottom: 1px solid #f1f2fb; }
.font-table td { padding: 10px; border-bottom: 1px solid #f1f2fb; color: #34395e; vertical-align: middle; }
.font-preview { font-size: 16px; }
.font-format-badge { background: #eef0fd; color: #6777ef; border-radius: 30px; padding: 3px 10px; font-size: 11px; font-weight: 700; text-transform: uppercase; }
.font-delete { border: 0; border-radius: 30px; background: #fc544b; color: #fff; padding: 6px 14px; font-size: 11px; font-weight: 700; }
.font-empty { color: #94a3b8; font-style: italic; text-align: center; padding: 20px; }

.typography-mapping { display: grid; gap: 20px; }
.typography-row { display: grid; grid-template-columns: minmax(10rem, 1fr) 1.3fr 1fr 1fr; gap: 14px; align-items: end; padding: 12px 0; border-top: 1px solid #f1f2fb; }
.typography-row:first-of-type { border-top: 0; }
.typography-row__label { font-size: 13px; font-weight: 700; color: #34395e; align-self: center; }
.typography-row label { display: grid; gap: 6px; font-size: 11px; font-weight: 600; color: #94a3b8; text-transform: uppercase; letter-spacing: .03em; min-width: 0; }
.typography-row select,
.typography-row input { border: 1px solid #e4e6fc; border-radius: 3px; padding: 8px 12px; font-weight: 500; color: #495057; background: #fdfdff; outline: none; width: 100%; box-sizing: border-box; }
.typography-size { display: flex; gap: 6px; }
.typography-size input { width: 5rem; }
.typography-size select { width: 4.5rem; }

@media (max-width: 900px) {
  .font-upload-row { grid-template-columns: 1fr; }
  .typography-row { grid-template-columns: 1fr 1fr; }
}

:global(.dark .cms-panel-head h2),
:global(.dark .cms-panel-head h3) { color: #f8fafc; }
:global(.dark .cms-panel) { background: #111827; border-color: #334155; }
:global(.dark .appearance-tabs button) { background: #0f172a; color: #f8fafc; border-color: #475569; }
:global(.dark .font-upload-row input),
:global(.dark .typography-row select),
:global(.dark .typography-row input) { background: #0f172a; color: #f8fafc; border-color: #475569; }
:global(.dark .dropzone) { background: #1f2937; border-color: #6777ef; }
:global(.dark .dropzone.dz-drag-hover),
:global(.dark .dropzone.dz-started) { background: #2b240f; border-color: #ffa426; }
:global(.dark .otika-drop-empty h6),
:global(.dark .dz-details) { color: #f8fafc; }
:global(.dark .otika-drop-empty) { color: #cbd5e1; }
:global(.dark .upload-actions .btn-secondary) { background: #0f172a; color: #f8fafc; border-color: #475569; }
:global(.dark .font-table th) { color: #64748b; }
:global(.dark .font-table td) { color: #f8fafc; border-color: #1e293b; }
:global(.dark .typography-row) { border-color: #1e293b; }
:global(.dark .typography-row__label) { color: #f8fafc; }
</style>
