<template>
  <LayoutDefault title="File Storage">
    <div class="max-w-6xl mx-auto p-6 space-y-6">
      <section class="bg-white border border-gray-100 rounded-xl shadow-sm">
        <div class="px-5 py-4 border-b border-gray-100 flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
          <div>
            <p class="text-xs font-bold uppercase tracking-wider text-primary-700">System configuration</p>
            <h1 class="text-xl font-bold text-gray-900 mt-1">File Storage</h1>
          </div>
          <button @click="save" :disabled="saving || loading" class="inline-flex items-center justify-center gap-2 rounded-lg bg-primary-700 px-4 py-2 text-sm font-semibold text-white hover:bg-primary-800 disabled:opacity-60">
            <i :class="saving ? 'icofont-spinner-alt-1 animate-spin' : 'icofont-save'"></i>
            {{ saving ? 'Saving' : 'Save settings' }}
          </button>
        </div>

        <div class="p-5 grid gap-5 lg:grid-cols-2">
          <section class="border border-gray-100 rounded-lg p-4">
            <div class="flex items-start justify-between gap-3">
              <div>
                <h2 class="font-semibold text-gray-900">Public Website Uploads</h2>
                <p class="text-xs text-gray-500 mt-1">Images, logos, PDFs, and dynamic public assets.</p>
              </div>
              <span class="rounded-full bg-blue-50 px-3 py-1 text-xs font-semibold text-blue-700">{{ form.public_provider }}</span>
            </div>
            <div class="mt-4 grid grid-cols-2 gap-2">
              <button type="button" class="provider-option" :class="form.public_provider === 'google_drive' && 'provider-option-active'" @click="form.public_provider = 'google_drive'">
                <i class="icofont-brand-google"></i>
                Google Drive
              </button>
              <button type="button" class="provider-option" :class="form.public_provider === 'local' && 'provider-option-active'" @click="form.public_provider = 'local'">
                <i class="icofont-hdd"></i>
                Local disk
              </button>
            </div>
          </section>

          <section class="border border-gray-100 rounded-lg p-4">
            <div class="flex items-start justify-between gap-3">
              <div>
                <h2 class="font-semibold text-gray-900">Application Uploads</h2>
                <p class="text-xs text-gray-500 mt-1">Applicant documents, form attachments, and evidence files.</p>
              </div>
              <span class="rounded-full bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-700">{{ form.application_provider }}</span>
            </div>
            <div class="mt-4 grid grid-cols-2 gap-2">
              <button type="button" class="provider-option" :class="form.application_provider === 's3' && 'provider-option-active'" @click="form.application_provider = 's3'">
                <i class="icofont-cloud-upload"></i>
                Amazon S3
              </button>
              <button type="button" class="provider-option" :class="form.application_provider === 'local' && 'provider-option-active'" @click="form.application_provider = 'local'">
                <i class="icofont-hdd"></i>
                Local disk
              </button>
            </div>
          </section>
        </div>
      </section>

      <section class="grid gap-6 lg:grid-cols-2">
        <article class="bg-white border border-gray-100 rounded-xl shadow-sm">
          <div class="px-5 py-4 border-b border-gray-100">
            <h2 class="font-semibold text-gray-900">Local Disk</h2>
          </div>
          <div class="p-5 grid gap-4">
            <Field label="Public files path">
              <input v-model="form.local_public_path" class="input" />
            </Field>
            <Field label="Public URL prefix">
              <input v-model="form.local_public_url_prefix" class="input" />
            </Field>
            <Field label="Application files path">
              <input v-model="form.local_app_path" class="input" />
            </Field>
            <Field label="Application URL prefix">
              <input v-model="form.local_app_url_prefix" class="input" />
            </Field>
          </div>
        </article>

        <article class="bg-white border border-gray-100 rounded-xl shadow-sm">
          <div class="px-5 py-4 border-b border-gray-100">
            <h2 class="font-semibold text-gray-900">Google Drive API</h2>
          </div>
          <div class="p-5 grid gap-4">
            <Field label="Folder ID">
              <input v-model="form.google_drive_folder_id" class="input" placeholder="Drive folder ID" />
            </Field>
            <div class="grid gap-4 md:grid-cols-2">
              <Field label="OAuth client ID">
                <input v-model="form.google_drive_client_id" class="input" placeholder="Google API client ID" />
              </Field>
              <Field label="OAuth client secret">
                <input v-model="form.google_drive_client_secret" type="password" class="input" placeholder="Leave blank to keep current value" />
              </Field>
            </div>
            <Field label="OAuth refresh token">
              <input v-model="form.google_drive_refresh_token" type="password" class="input" placeholder="Leave blank to keep current value" />
            </Field>
            <Field label="Service account JSON">
              <textarea v-model="form.google_drive_credentials_json" rows="5" class="input font-mono text-xs" placeholder="Paste JSON only when setting or replacing credentials"></textarea>
            </Field>
            <label class="flex items-center gap-2 text-sm font-medium text-gray-700">
              <input v-model="form.google_drive_make_public" type="checkbox" class="rounded text-primary-700" />
              Make uploaded public files readable by link
            </label>
          </div>
        </article>
      </section>

      <section class="bg-white border border-gray-100 rounded-xl shadow-sm">
        <div class="px-5 py-4 border-b border-gray-100">
          <h2 class="font-semibold text-gray-900">Amazon S3 Bucket</h2>
        </div>
        <div class="p-5 grid gap-4 md:grid-cols-2">
          <Field label="Bucket">
            <input v-model="form.s3_bucket" class="input" />
          </Field>
          <Field label="Region">
            <input v-model="form.s3_region" class="input" />
          </Field>
          <Field label="Key prefix">
            <input v-model="form.s3_prefix" class="input" placeholder="applications" />
          </Field>
          <Field label="Public base URL">
            <input v-model="form.s3_public_base_url" class="input" placeholder="https://cdn.example.org" />
          </Field>
          <Field label="Custom endpoint">
            <input v-model="form.s3_endpoint" class="input" placeholder="Optional S3-compatible endpoint" />
          </Field>
          <Field label="Access key ID">
            <input v-model="form.s3_access_key_id" class="input" placeholder="Leave blank to keep current value" />
          </Field>
          <Field label="Secret access key">
            <input v-model="form.s3_secret_access_key" type="password" class="input" placeholder="Leave blank to keep current value" />
          </Field>
          <label class="flex items-center gap-2 text-sm font-medium text-gray-700 md:self-end md:pb-2">
            <input v-model="form.s3_force_path_style" type="checkbox" class="rounded text-primary-700" />
            Use path-style addressing
          </label>
        </div>
      </section>

      <p v-if="message" class="rounded-lg border px-4 py-3 text-sm" :class="error ? 'border-red-200 bg-red-50 text-red-700' : 'border-emerald-200 bg-emerald-50 text-emerald-700'">{{ message }}</p>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { defineComponent, h, onMounted, reactive, ref } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'
import { getStorageSettings, updateStorageSettings } from '@/api/storageSettings.js'

const breadcrumbStore = useBreadcrumbStore()
const loading = ref(true)
const saving = ref(false)
const message = ref('')
const error = ref(false)

const form = reactive({
  public_provider: 'google_drive',
  application_provider: 's3',
  local_public_path: '/app/uploads',
  local_public_url_prefix: '/uploads',
  local_app_path: '/app/uploads/applications',
  local_app_url_prefix: '/uploads/applications',
  google_drive_folder_id: '',
  google_drive_credentials_json: '',
  google_drive_client_id: '',
  google_drive_client_secret: '',
  google_drive_refresh_token: '',
  google_drive_make_public: true,
  s3_bucket: '',
  s3_region: 'us-east-1',
  s3_prefix: 'applications',
  s3_endpoint: '',
  s3_public_base_url: '',
  s3_force_path_style: false,
  s3_access_key_id: '',
  s3_secret_access_key: '',
})

async function load() {
  loading.value = true
  try {
    Object.assign(form, await getStorageSettings())
  } catch (e) {
    show(e.response?.data?.error?.message || 'Could not load settings', true)
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    Object.assign(form, await updateStorageSettings({ ...form }))
    form.google_drive_credentials_json = ''
    form.google_drive_client_secret = ''
    form.google_drive_refresh_token = ''
    form.s3_secret_access_key = ''
    show('Storage settings saved')
  } catch (e) {
    show(e.response?.data?.error?.message || 'Could not save settings', true)
  } finally {
    saving.value = false
  }
}

function show(text, isError = false) {
  message.value = text
  error.value = isError
}

const Field = defineComponent({
  props: { label: String },
  setup(props, { slots }) {
    return () => h('label', { class: 'block text-sm font-medium text-gray-700' }, [
      h('span', props.label),
      h('div', { class: 'mt-1' }, slots.default?.()),
    ])
  },
})

onMounted(() => {
  breadcrumbStore.set('File Storage', [{ label: 'Operations' }, { label: 'File Storage' }])
  load()
})
</script>

<style scoped>
.input { width: 100%; border-radius: 0.5rem; border: 1px solid #d1d5db; padding: 0.5rem 0.75rem; font-size: 0.875rem; }
.provider-option { display: inline-flex; align-items: center; justify-content: center; gap: 0.5rem; min-height: 2.75rem; border-radius: 0.5rem; border: 1px solid #e5e7eb; color: #374151; font-size: 0.875rem; font-weight: 600; background: white; }
.provider-option-active { border-color: #1d4ed8; background: #eff6ff; color: #1d4ed8; }
</style>
