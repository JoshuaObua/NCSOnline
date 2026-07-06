<template>
  <section class="website-settings">
    <div class="cms-panel-head"><h2>Website Settings</h2></div>

    <div class="cms-panel">
      <div class="cms-panel-head">
        <h3>Site Identity</h3>
        <button type="button" :disabled="savingIdentity" @click="saveIdentity">{{ savingIdentity ? 'Saving…' : 'Save identity' }}</button>
      </div>
      <div class="cms-two">
        <label>Site name<input v-model="identity.name" class="form-control" placeholder="NCS Uganda" /></label>
      </div>
      <div class="cms-logo-grid">
        <div class="cms-logo-field">
          <label>Favicon</label>
          <DropzoneUpload v-model="identity.faviconUrl" label="favicon" hint="Square PNG or ICO, at least 32x32." accept="image/png,image/x-icon,image/vnd.microsoft.icon" @error="onError" />
        </div>
        <div class="cms-logo-field">
          <label>Main logo <small>(light backgrounds)</small></label>
          <DropzoneUpload v-model="identity.logoUrl" label="main logo" hint="PNG or SVG with a transparent background." accept="image/png,image/svg+xml,image/webp" @error="onError" />
        </div>
        <div class="cms-logo-field">
          <label>White logo <small>(dark backgrounds)</small></label>
          <DropzoneUpload v-model="identity.whiteLogoUrl" label="white logo" hint="Used in the header when dark mode is active." accept="image/png,image/svg+xml,image/webp" @error="onError" />
        </div>
        <div class="cms-logo-field">
          <label>Footer logo</label>
          <DropzoneUpload v-model="identity.footerLogoUrl" label="footer logo" hint="Shown in the site footer. Falls back to the white/main logo if left empty." accept="image/png,image/svg+xml,image/webp" @error="onError" />
        </div>
      </div>
    </div>

    <div class="cms-panel">
      <div class="cms-panel-head">
        <h3>SEO Defaults</h3>
        <button type="button" :disabled="savingSeo" @click="saveSeo">{{ savingSeo ? 'Saving…' : 'Save SEO defaults' }}</button>
      </div>
      <p class="cms-hint">Used as a fallback wherever an individual page doesn't set its own title, description, or share image.</p>
      <div class="cms-two">
        <label>Default meta title<input v-model="seo.meta_title" class="form-control" maxlength="180" placeholder="National Council of Sports - Uganda" /></label>
        <label>Meta keywords <small>(comma separated)</small><input v-model="seo.meta_keywords" class="form-control" placeholder="sports, uganda, council" /></label>
        <label class="wide">Default meta description<textarea v-model="seo.meta_description" class="form-control" maxlength="300"></textarea></label>
      </div>
      <div class="cms-logo-field">
        <label>Default social share image <small>(Open Graph / Twitter card)</small></label>
        <DropzoneUpload v-model="seo.meta_image_url" label="social share image" hint="Recommended size 1200x630." accept="image/png,image/jpeg,image/webp" @error="onError" />
      </div>
    </div>
  </section>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { getSettings, adminUpdateSettings } from '@/api/cms.js'
import DropzoneUpload from '@/components/cms/DropzoneUpload.vue'

const emit = defineEmits(['message', 'error'])

const identity = reactive({ name: '', faviconUrl: '', logoUrl: '', whiteLogoUrl: '', footerLogoUrl: '' })
const seo = reactive({ meta_title: '', meta_description: '', meta_keywords: '', meta_image_url: '' })
const savingIdentity = ref(false)
const savingSeo = ref(false)

function onError(err) { emit('error', err) }

async function loadIdentity() {
  try {
    const r = await getSettings('site')
    const v = r.data?.data?.value
    if (v && typeof v === 'object') Object.assign(identity, v)
  } catch { /* first-time, no record yet */ }
}

async function loadSeo() {
  try {
    const r = await getSettings('seo')
    const v = r.data?.data?.value
    if (v && typeof v === 'object') Object.assign(seo, v)
  } catch { /* first-time, no record yet */ }
}

async function saveIdentity() {
  savingIdentity.value = true
  try {
    await adminUpdateSettings('site', { ...identity })
    emit('message', 'Site identity saved')
  } catch (err) {
    emit('error', err.response?.data?.error?.message || err.message || 'Could not save site identity')
  } finally {
    savingIdentity.value = false
  }
}

async function saveSeo() {
  savingSeo.value = true
  try {
    await adminUpdateSettings('seo', { ...seo })
    emit('message', 'SEO defaults saved')
  } catch (err) {
    emit('error', err.response?.data?.error?.message || err.message || 'Could not save SEO defaults')
  } finally {
    savingSeo.value = false
  }
}

onMounted(() => { loadIdentity(); loadSeo() })
</script>

<style scoped>
.cms-logo-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1.25rem; margin-top: .75rem; }
.cms-logo-field label { display: block; margin-bottom: .4rem; font-weight: 600; }
.cms-logo-field { margin-top: 1rem; }
.cms-hint { color: var(--text-muted, #6c757d); margin-top: -.25rem; margin-bottom: .75rem; }
</style>
