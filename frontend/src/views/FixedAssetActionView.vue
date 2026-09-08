<template>
  <LayoutDefault :title="page.title">
    <section class="asset-action-page bg-white rounded-xl border border-gray-200 p-6 space-y-6">
      <header class="space-y-2">
        <router-link to="/fixed-assets" class="text-sm font-semibold text-emerald-700">← Back to Asset Register</router-link>
        <h1 class="text-2xl font-bold text-gray-900">{{ page.title }}</h1>
        <p class="text-sm text-gray-600">{{ page.description }}</p>
      </header>
      <p v-if="loading" role="status">Loading asset details…</p>
      <div v-if="error" role="alert" class="rounded-lg border border-red-200 bg-red-50 p-4 text-red-800">{{ error }}</div>
      <form v-if="!loading && (!needsAsset || asset)" @submit.prevent="submit" class="space-y-6" :aria-busy="saving">
        <fieldset :disabled="saving" class="space-y-6">
          <div v-if="asset" class="rounded-lg border border-gray-200 bg-gray-50 p-4 text-sm text-gray-800">
            <p class="font-bold">{{ asset.asset_number }} · {{ asset.tag_number }}</p>
            <p>{{ asset.asset_description }}</p>
          </div>
          <div v-if="action === 'new'" class="grid grid-cols-1 md:grid-cols-2 gap-4">

            <label class="text-xs font-semibold text-gray-700">Asset Number<input required v-model="form.asset_number" class="form-input" placeholder="M1009999" /></label>
            <label class="text-xs font-semibold text-gray-700">Tag Number<input required v-model="form.tag_number" class="form-input" placeholder="NCS-TAG-001" /></label>
            <label class="md:col-span-2 text-xs font-semibold text-gray-700">Description<input required v-model="form.asset_description" class="form-input" placeholder="Asset description" /></label>
            <label class="text-xs font-semibold text-gray-700">Category Segment 1<input v-model="form.category_segment1" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Category Segment 3<select v-model="form.category_segment3" class="form-input"><option v-for="cat in categoriesList" :key="cat" :value="cat">{{ cat }}</option></select></label>
            <label class="text-xs font-semibold text-gray-700">Category Segment 4<input v-model="form.category_segment4" class="form-input" placeholder="Subclass" /></label>
            <label class="text-xs font-semibold text-gray-700">Units<input required v-model.number="form.asset_units" type="number" min="1" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">FB Cost<input required v-model.number="form.fb_cost" type="number" min="0" step="0.01" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Adjusted Cost<input required v-model.number="form.adjusted_cost" type="number" min="0" step="0.01" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Date In Service<input required v-model="form.date_placed_in_service" type="date" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Department<input v-model="form.custodian_department" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Location<input v-model="form.location_building" class="form-input" /></label>
            <label class="text-xs font-semibold text-gray-700">Useful Life Years<input required v-model.number="form.useful_life_years" type="number" min="0" class="form-input" /></label>

          </div>
          <div v-else-if="action === 'depreciation'" class="space-y-4">
            <label class="field-label">Period<input required v-model="form.period" type="month" class="form-input" /></label>
          </div>
          <div v-else-if="action === 'verify'" class="space-y-4">
            <label class="field-label">Verification Status<select v-model="form.status" class="form-input"><option value="VERIFIED">VERIFIED</option><option value="DISCREPANCY">DISCREPANCY</option><option value="MISSING">MISSING</option></select></label>
            <label class="field-label">Notes<textarea v-model="form.notes" rows="4" class="form-input" placeholder="Spot-check notes"></textarea></label>
          </div>
          <div v-else class="space-y-4">
            <label class="field-label">Current Adjusted Cost (UGX)<input :value="formatUGX(asset.adjusted_cost)" readonly class="form-input bg-gray-100" /></label>
            <label class="field-label">New Adjusted Valuation (UGX)<input required v-model.number="form.new_cost" type="number" min="0" step="0.01" class="form-input" /></label>
            <label class="field-label">Revaluation Notes &amp; Justification<textarea v-model="form.notes" rows="4" class="form-input" placeholder="Reason for revaluation..."></textarea></label>
          </div>
          <div class="flex flex-wrap justify-end gap-3 border-t border-gray-200 pt-5">
            <router-link to="/fixed-assets" class="px-4 py-2 rounded-lg bg-gray-100 text-gray-700 font-semibold text-sm">Cancel</router-link>
            <button type="submit" :disabled="saving" class="px-4 py-2 rounded-lg bg-emerald-700 text-white font-semibold text-sm disabled:opacity-50">{{ saving ? 'Saving…' : page.submit }}</button>
          </div>
        </fieldset>
      </form>
    </section>
  </LayoutDefault>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import client from '@/api/client'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'

const props = defineProps({ action: { type: String, required: true } })
const route = useRoute()
const router = useRouter()
const pages = {
  new: { title: 'Register New Fixed Asset', description: 'Enter the asset details and opening values.', submit: 'Create Asset' },
  revalue: { title: 'Post Asset Revaluation', description: 'Record an adjusted valuation and its justification.', submit: 'Save Revaluation' },
  verify: { title: 'Physical Tag Verification', description: 'Record the physical asset check and any discrepancies.', submit: 'Save Verification' },
  depreciation: { title: 'Run Monthly Depreciation', description: 'Select the period for the depreciation run.', submit: 'Run Depreciation' },
}
const page = computed(() => pages[props.action])
const needsAsset = computed(() => ['revalue', 'verify'].includes(props.action))
const asset = ref(null)
const form = ref({})
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const categoriesList = [
  'CYCLES',
  'ELECTRICAL MACHINERY',
  'FURNITURE AND FITTINGS',
  'LAND',
  'LIGHT ICT HARDWARE',
  'LIGHT VEHICLES',
  'NON RESIDENTIAL BUILDINGS',
  'OFFICE EQUIPMENT',
  'OTHER ICT EQUIPMENT',
  'RESIDENTIAL BUILDINGS',
]
const apiData = response => response?.data?.data ?? response?.data ?? {}
const errorMessage = (e, fallback) => e.response?.data?.error?.message || e.response?.data?.message || (typeof e.response?.data?.error === 'string' ? e.response.data.error : fallback)
const formatUGX = value => new Intl.NumberFormat('en-UG', { maximumFractionDigits: 2 }).format(value || 0)

function defaultNewAsset() {
  return {
    asset_book: 'NCS FA BOOK',
    asset_number: '',
    tag_number: '',
    asset_description: '',
    category_segment1: 'MACHINERY AND EQUIPMENT',
    category_segment3: 'LIGHT ICT HARDWARE',
    category_segment4: '',
    asset_units: 1,
    fb_cost: 0,
    adjusted_cost: 0,
    date_placed_in_service: '2023-07-01',
    custodian_department: 'General Administration',
    location_building: 'NCS Lugogo Head Office',
    location_room: 'Main Facility',
    depreciation_method: 'STRAIGHT_LINE',
    useful_life_years: 5,
    worksheet_source: 'LIGHT ICT HARDWARE',
  }
}

let loadVersion = 0
watch(() => [props.action, route.params.id], async () => {
  const version = ++loadVersion
  error.value = ''
  asset.value = null
  form.value = props.action === 'new' ? defaultNewAsset() : { period: new Date().toISOString().slice(0, 7) }
  loading.value = needsAsset.value
  if (!needsAsset.value) return
  try {
    const response = await client.get(`/api/v1/assets/${encodeURIComponent(route.params.id)}`)
    if (version !== loadVersion) return
    const data = apiData(response)
    if (!data.asset?.id) throw new Error('Asset not found')
    asset.value = data.asset
    form.value = props.action === 'revalue'
      ? { new_cost: data.asset.adjusted_cost, notes: '' }
      : { status: ['VERIFIED', 'DISCREPANCY', 'MISSING'].includes(data.asset.verification_status) ? data.asset.verification_status : 'VERIFIED', notes: '' }
  } catch (e) {
    if (version === loadVersion) error.value = errorMessage(e, 'Unable to load this asset. Return to the register and try again.')
  } finally {
    if (version === loadVersion) loading.value = false
  }
}, { immediate: true })

async function submit() {
  if (saving.value || loading.value || (needsAsset.value && !asset.value)) return
  saving.value = true
  error.value = ''
  const action = props.action
  try {
    let endpoint, payload
    if (action === 'new') {
      endpoint = '/api/v1/assets'
      payload = { ...form.value, worksheet_source: form.value.category_segment3 }
      if (!payload.adjusted_cost && payload.fb_cost) payload.adjusted_cost = payload.fb_cost
    } else if (action === 'depreciation') {
      endpoint = '/api/v1/assets/depreciate'
      payload = { period: form.value.period }
    } else {
      endpoint = `/api/v1/assets/${action}`
      payload = { asset_id: asset.value.id, ...form.value }
    }
    await client.post(endpoint, payload)
    await router.push({ name: 'FixedAssets', query: { completed: action } })
  } catch (e) {
    error.value = errorMessage(e, 'Unable to save. Please try again.')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.field-label { display: block; font-size: .875rem; font-weight: 600; color: #374151; }
.form-input { display: block; width: 100%; margin-top: .35rem; border: 1px solid #d1d5db; border-radius: .5rem; padding: .625rem; font-size: .875rem; }
.form-input:focus { outline: 2px solid #047857; outline-offset: 2px; }
</style>
