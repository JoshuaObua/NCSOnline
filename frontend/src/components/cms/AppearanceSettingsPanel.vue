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
        <label class="wide">Font file<input ref="fileInput" type="file" accept=".ttf,.otf,.woff,.woff2" @change="onFilePicked" /></label>
        <button type="button" :disabled="uploading || !newFontFile || !newFontName.trim()" @click="uploadFont">{{ uploading ? 'Uploading…' : 'Upload font' }}</button>
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

function quoted(name) {
  return /\s/.test(name) ? `'${name}'` : name
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
  newFontFile.value = event.target.files?.[0] || null
}

async function uploadFont() {
  if (!newFontFile.value || !newFontName.value.trim()) return
  uploading.value = true
  try {
    await adminUploadFont(newFontFile.value, newFontName.value.trim())
    newFontName.value = ''
    newFontFile.value = null
    if (fileInput.value) fileInput.value.value = ''
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

.font-upload-row { display: grid; grid-template-columns: 1fr 1fr auto; gap: 14px; align-items: end; }
.font-upload-row label { display: grid; gap: 7px; font-size: 12px; font-weight: 600; color: #34395e; min-width: 0; }
.font-upload-row input { border: 1px solid #e4e6fc; border-radius: 3px; padding: 10px 15px; font-weight: 500; color: #495057; background: #fdfdff; outline: none; }
.font-upload-row button { border: 0; border-radius: 30px; background: #6777ef; color: #fff; padding: 10px 18px; font-size: 12px; font-weight: 700; white-space: nowrap; }
.font-upload-row button:disabled { opacity: .5; }

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
:global(.dark .font-table th) { color: #64748b; }
:global(.dark .font-table td) { color: #f8fafc; border-color: #1e293b; }
:global(.dark .typography-row) { border-color: #1e293b; }
:global(.dark .typography-row__label) { color: #f8fafc; }
</style>
