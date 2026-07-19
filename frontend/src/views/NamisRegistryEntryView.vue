<template>
  <main class="registry-entry-page">
    <section v-if="definition" class="registry-card">
      <header class="page-header">
        <div>
          <p>Sports Registry</p>
          <h2>Add {{ definition.label }}</h2>
          <span>Complete the fields below to create a new {{ definition.label.toLowerCase() }} record.</span>
        </div>
        <router-link class="back-button" :to="listRoute"><i class="icofont-rounded-left"></i> Back to {{ listLabel }}</router-link>
      </header>

      <div v-if="definition.sensitive" class="sensitive-notice">
        <i class="icofont-shield"></i>
        <span>This module contains restricted personal information. Only enter authorised records.</span>
      </div>
      <p v-if="loadError" class="alert alert-error">{{ loadError }}</p>
      <p v-if="successMessage" class="alert alert-success">{{ successMessage }}</p>

      <form class="registry-form" @submit.prevent="saveRecord">
        <div v-for="item in definition.fields" :key="item.key" class="field" :class="{ wide: isWide(item) }">
          <label :for="`registry-${item.key}`">{{ item.label }} <strong v-if="item.required">*</strong></label>

          <select v-if="item.type === 'select'" :id="`registry-${item.key}`" v-model="form[item.key]" :required="item.required">
            <option value="">Select {{ item.label.toLowerCase() }}</option>
            <option v-for="option in item.options" :key="option" :value="option">{{ readable(option) }}</option>
          </select>
          <select v-else-if="item.type === 'federation'" :id="`registry-${item.key}`" v-model="form[item.key]" :required="item.required" :disabled="loadingReferences">
            <option value="">Select federation</option>
            <option v-for="option in federations" :key="option.id" :value="option.id">{{ option.name }}{{ option.acronym ? ` (${option.acronym})` : '' }}</option>
          </select>
          <select v-else-if="item.type === 'athlete'" :id="`registry-${item.key}`" v-model="form[item.key]" :required="item.required" :disabled="loadingReferences">
            <option value="">Select athlete</option>
            <option v-for="option in athletes" :key="option.id" :value="option.id">{{ option.full_name }}{{ option.athlete_number ? ` (${option.athlete_number})` : '' }}</option>
          </select>
          <select v-else-if="item.type === 'competition'" :id="`registry-${item.key}`" v-model="form[item.key]" :required="item.required" :disabled="loadingReferences">
            <option value="">Select competition</option>
            <option v-for="option in competitions" :key="option.id" :value="option.id">{{ option.name }}{{ option.venue ? ` - ${option.venue}` : '' }}</option>
          </select>
          <label v-else-if="item.type === 'boolean'" class="checkbox-field">
            <input :id="`registry-${item.key}`" v-model="form[item.key]" type="checkbox" />
            <span>Yes</span>
          </label>
          <textarea v-else-if="item.type === 'textarea'" :id="`registry-${item.key}`" v-model.trim="form[item.key]" rows="4" :required="item.required"></textarea>
          <textarea v-else-if="item.type === 'json'" :id="`registry-${item.key}`" v-model="form[item.key]" rows="5" placeholder='{"name":"value"}' :required="item.required"></textarea>
          <input v-else :id="`registry-${item.key}`" v-model="form[item.key]" :type="item.type || 'text'" :min="item.type === 'number' ? 0 : undefined" :step="item.type === 'number' ? 'any' : undefined" :required="item.required" />
        </div>

        <p v-if="formError" class="alert alert-error wide">{{ formError }}</p>
        <footer class="form-actions wide">
          <router-link class="cancel-button" :to="listRoute">Cancel</router-link>
          <button class="save-button" type="submit" :disabled="saving || loadingReferences"><i class="icofont-save"></i> {{ saving ? 'Saving...' : `Save ${definition.label}` }}</button>
        </footer>
      </form>
    </section>
    <section v-else class="registry-card"><p class="alert alert-error">Unknown Sports Registry form.</p></section>
  </main>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { createNsmisDomain, listNsmisDomain } from '@/api/nsmis.js'
import { registryResources } from '@/utils/namisRegistryConfig.js'

const route = useRoute()
const router = useRouter()
const resource = computed(() => String(route.params.resource || ''))
const definition = computed(() => registryResources[resource.value])
const listRoute = computed(() => ({ path: '/portal', query: { section: `namis-${resource.value}` } }))
const listLabel = computed(() => definition.value?.label || 'registry')
const form = reactive({})
const federations = ref([])
const athletes = ref([])
const competitions = ref([])
const loadingReferences = ref(false)
const saving = ref(false)
const loadError = ref('')
const formError = ref('')
const successMessage = ref('')

watch(definition, initialiseForm, { immediate: true })
onMounted(loadReferences)

function initialiseForm() {
  for (const key of Object.keys(form)) delete form[key]
  for (const item of definition.value?.fields || []) form[item.key] = item.type === 'boolean' ? false : ''
  if (resource.value === 'equipment') form.unit = 'ITEM'
}
function isWide(item) { return ['textarea', 'json'].includes(item.type) || ['full_name','item_name','name','reference'].includes(item.key) }
function readable(value) { return String(value).toLowerCase().replaceAll('_', ' ').replace(/\b\w/g, char => char.toUpperCase()) }
async function loadReferences() {
  const types = new Set((definition.value?.fields || []).map(item => item.type))
  const requests = []
  if (types.has('federation')) requests.push(listNsmisDomain('federations', { per_page: 100 }).then(res => { federations.value = res.data?.data?.items || res.data?.items || [] }))
  if (types.has('athlete')) requests.push(listNsmisDomain('athletes', { per_page: 100 }).then(res => { athletes.value = res.data?.data?.items || res.data?.items || [] }))
  if (types.has('competition')) requests.push(listNsmisDomain('competitions', { per_page: 100 }).then(res => { competitions.value = res.data?.data?.items || res.data?.items || [] }))
  if (!requests.length) return
  loadingReferences.value = true
  try { await Promise.all(requests) } catch (err) { loadError.value = err.response?.data?.error?.message || 'Could not load form options.' } finally { loadingReferences.value = false }
}
function buildPayload() {
  const payload = {}
  for (const item of definition.value.fields) {
    let value = form[item.key]
    if (item.type === 'json') {
      if (!String(value || '').trim()) value = {}
      else { try { value = JSON.parse(value) } catch { throw new Error(`${item.label} must contain valid JSON.`) } }
    }
    if (item.type === 'number' && value !== '') value = Number(value)
    if (value !== '' || item.type === 'boolean') payload[item.key] = value
  }
  if (resource.value === 'equipment' && Number(payload.quantity_distributed) > Number(payload.quantity_received)) throw new Error('Quantity distributed cannot exceed quantity received.')
  if (resource.value === 'competitions' && payload.starts_on && payload.ends_on && payload.ends_on < payload.starts_on) throw new Error('Competition end date cannot be before its start date.')
  return payload
}
async function saveRecord() {
  formError.value = ''; successMessage.value = ''
  saving.value = true
  try {
    await createNsmisDomain(resource.value, buildPayload())
    successMessage.value = `${definition.value.label} saved successfully.`
    window.setTimeout(() => router.push(listRoute.value), 500)
  } catch (err) { formError.value = err.response?.data?.error?.message || err.message || 'Could not save this record.' } finally { saving.value = false }
}
</script>

<style scoped>
.registry-entry-page{padding:0;color:#34395e}.registry-card{width:100%;padding:30px;box-sizing:border-box;background:#fff;border:1px solid #e8eaf0;border-radius:14px;box-shadow:0 8px 30px rgba(44,62,80,.08)}.page-header{display:flex;justify-content:space-between;gap:24px;align-items:flex-start;padding-bottom:24px;border-bottom:1px solid #eceef3}.page-header p{margin:0 0 5px;color:#6777ef;font-size:12px;font-weight:800;letter-spacing:.08em;text-transform:uppercase}.page-header h2{margin:0 0 8px;font-size:28px}.page-header span{color:#6c757d}.back-button,.cancel-button{display:inline-flex;align-items:center;gap:7px;padding:10px 15px;border:1px solid #dfe2ea;border-radius:7px;color:#4f5d73;text-decoration:none;font-weight:700;white-space:nowrap}.registry-form{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:22px;padding-top:26px}.field{display:flex;flex-direction:column;gap:8px}.wide{grid-column:1/-1}.field label{font-size:13px;font-weight:700}.field label strong{color:#e74c3c}.field input:not([type=checkbox]),.field select,.field textarea{width:100%;box-sizing:border-box;padding:10px 12px;border:1px solid #d9dce5;border-radius:7px;background:#fff;color:#34395e;font:inherit}.field input:not([type=checkbox]),.field select{min-height:44px}.field textarea{resize:vertical}.field input:focus,.field select:focus,.field textarea:focus{border-color:#6777ef;outline:3px solid rgba(103,119,239,.12)}.checkbox-field{min-height:44px;display:flex;align-items:center;gap:10px;padding:0 12px;border:1px solid #d9dce5;border-radius:7px}.checkbox-field input{width:18px;height:18px}.sensitive-notice{display:flex;gap:10px;align-items:center;margin-top:20px;padding:13px 15px;border-radius:8px;background:#fff8e7;color:#856404;border:1px solid #ffe6a7}.alert{padding:12px 14px;margin:20px 0 0;border-radius:7px;font-weight:600}.alert-error{background:#fff1f0;color:#c0392b;border:1px solid #ffd4d0}.alert-success{background:#edf9f0;color:#218838;border:1px solid #ccebd4}.form-actions{display:flex;justify-content:flex-end;gap:12px;padding-top:22px;border-top:1px solid #eceef3}.save-button{display:inline-flex;align-items:center;gap:8px;padding:11px 18px;border:0;border-radius:7px;background:#6777ef;color:#fff;font-weight:700;cursor:pointer}.save-button:disabled{opacity:.65;cursor:wait}@media(max-width:720px){.registry-card{padding:20px}.page-header{flex-direction:column}.registry-form{grid-template-columns:1fr}.wide{grid-column:auto}.form-actions{flex-direction:column-reverse}.back-button,.cancel-button,.save-button{justify-content:center}}
</style>
