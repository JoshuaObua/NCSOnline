<template>
  <main class="equipment-entry-page">
    <section class="equipment-card">
      <header class="page-header">
        <div>
          <p>NAMIS Registry Manager</p>
          <h1>Add Distributed Equipment</h1>
          <span>Record equipment received by a sports federation and the quantity already distributed.</span>
        </div>
        <router-link class="back-button" :to="listRoute">
          <i class="icofont-rounded-left"></i> Back to equipment
        </router-link>
      </header>

      <p v-if="loadError" class="alert alert-error">{{ loadError }}</p>
      <p v-if="successMessage" class="alert alert-success">{{ successMessage }}</p>

      <form class="equipment-form" @submit.prevent="saveEquipment">
        <div class="field field-wide">
          <label for="equipment-federation">Sports federation <strong>*</strong></label>
          <SearchableFederationSelect
            id="equipment-federation"
            v-model="form.federation_id"
            :federations="federations"
            :required="true"
            :disabled="loadingFederations"
            :loading="loadingFederations"
            placeholder="Search and select sports federation (e.g. FUFA, UAF, Boxing)..."
          />
        </div>

        <div class="field field-wide">
          <label for="equipment-name">Equipment name <strong>*</strong></label>
          <input id="equipment-name" v-model.trim="form.item_name" required maxlength="160" placeholder="e.g. Match footballs" />
        </div>

        <div class="field">
          <label for="equipment-unit">Unit <strong>*</strong></label>
          <select id="equipment-unit" v-model="form.unit" required>
            <option v-for="unit in units" :key="unit" :value="unit">{{ unit }}</option>
          </select>
        </div>

        <div class="field">
          <label for="equipment-received">Quantity received <strong>*</strong></label>
          <input id="equipment-received" v-model.number="form.quantity_received" required type="number" min="0" step="1" />
        </div>

        <div class="field">
          <label for="equipment-distributed">Quantity distributed <strong>*</strong></label>
          <input id="equipment-distributed" v-model.number="form.quantity_distributed" required type="number" min="0" step="1" />
          <small>Cannot be greater than the quantity received.</small>
        </div>

        <aside class="stock-summary">
          <span>Remaining stock</span>
          <strong>{{ remainingStock }}</strong>
        </aside>

        <p v-if="formError" class="alert alert-error field-wide">{{ formError }}</p>

        <footer class="form-actions field-wide">
          <router-link class="cancel-button" :to="listRoute">Cancel</router-link>
          <button type="submit" class="save-button" :disabled="saving || loadingFederations">
            <i class="icofont-save"></i> {{ saving ? 'Saving...' : 'Save equipment entry' }}
          </button>
        </footer>
      </form>
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { createNsmisDomain, listNsmisDomain } from '@/api/nsmis.js'
import apiClient from '@/api/client.js'
import SearchableFederationSelect from '@/components/ui/SearchableFederationSelect.vue'

const router = useRouter()
const listRoute = { path: '/portal', query: { section: 'namis-equipment' } }
const units = ['ITEM', 'PIECE', 'PAIR', 'SET', 'BOX', 'CARTON', 'KIT', 'BAG']

const federations = ref([])
const loadingFederations = ref(false)
const saving = ref(false)
const loadError = ref('')
const formError = ref('')
const successMessage = ref('')
const form = reactive({
  federation_id: '',
  item_name: '',
  unit: 'ITEM',
  quantity_received: 0,
  quantity_distributed: 0,
})

const remainingStock = computed(() => Math.max(0, Number(form.quantity_received || 0) - Number(form.quantity_distributed || 0)))

onMounted(loadFederations)

function extractArray(res) {
  if (!res || !res.data) return []
  if (Array.isArray(res.data.data)) return res.data.data
  if (Array.isArray(res.data.items)) return res.data.items
  if (Array.isArray(res.data.data?.items)) return res.data.data.items
  if (Array.isArray(res.data)) return res.data
  return []
}

async function loadFederations() {
  loadingFederations.value = true
  loadError.value = ''
  try {
    const res = await listNsmisDomain('federations', { per_page: 200 })
    federations.value = extractArray(res)
  } catch (err) {
    try {
      const fallback = await apiClient.get('/api/v1/nsmis/federations')
      federations.value = extractArray(fallback)
    } catch (e) {
      loadError.value = err.response?.data?.error?.message || 'Could not load sports federations.'
    }
  } finally {
    loadingFederations.value = false
  }
}

async function saveEquipment() {
  formError.value = ''
  successMessage.value = ''
  const received = Number(form.quantity_received)
  const distributed = Number(form.quantity_distributed)
  if (!Number.isInteger(received) || received < 0 || !Number.isInteger(distributed) || distributed < 0) {
    formError.value = 'Quantities must be whole numbers of zero or more.'
    return
  }
  if (distributed > received) {
    formError.value = 'Quantity distributed cannot exceed quantity received.'
    return
  }

  saving.value = true
  try {
    await createNsmisDomain('equipment', {
      federation_id: form.federation_id,
      item_name: form.item_name,
      unit: form.unit,
      quantity_received: received,
      quantity_distributed: distributed,
    })
    successMessage.value = 'Equipment entry saved successfully.'
    window.setTimeout(() => router.push(listRoute), 500)
  } catch (err) {
    formError.value = err.response?.data?.error?.message || 'Could not save the equipment entry.'
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.equipment-entry-page { padding: 0; background: transparent; color: #34395e; }
.equipment-card { width: 100%; padding: 30px; box-sizing: border-box; background: #fff; border: 1px solid #e8eaf0; border-radius: 14px; box-shadow: 0 8px 30px rgba(44, 62, 80, .08); }
.page-header { display: flex; justify-content: space-between; gap: 24px; align-items: flex-start; padding-bottom: 24px; border-bottom: 1px solid #eceef3; }
.page-header p { margin: 0 0 5px; color: #6777ef; font-size: 12px; font-weight: 800; letter-spacing: .08em; text-transform: uppercase; }
.page-header h1 { margin: 0 0 8px; font-size: 28px; }
.page-header span { color: #6c757d; }
.back-button, .cancel-button { display: inline-flex; align-items: center; gap: 7px; padding: 10px 15px; border: 1px solid #dfe2ea; border-radius: 7px; color: #4f5d73; text-decoration: none; font-weight: 700; white-space: nowrap; }
.equipment-form { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 22px; padding-top: 26px; }
.field { display: flex; flex-direction: column; gap: 8px; }
.field-wide { grid-column: 1 / -1; }
.field label { font-size: 13px; font-weight: 700; }
.field label strong { color: #e74c3c; }
.field input, .field select { width: 100%; min-height: 44px; box-sizing: border-box; padding: 10px 12px; border: 1px solid #d9dce5; border-radius: 7px; background: #fff; color: #34395e; font: inherit; }
.field input:focus, .field select:focus { border-color: #6777ef; outline: 3px solid rgba(103, 119, 239, .12); }
.field small { color: #7b8190; }
.stock-summary { align-self: end; min-height: 74px; padding: 14px 18px; border-radius: 9px; background: #f0f2ff; display: flex; justify-content: space-between; align-items: center; }
.stock-summary span { color: #596275; font-weight: 700; }
.stock-summary strong { color: #6777ef; font-size: 28px; }
.alert { padding: 12px 14px; margin: 20px 0 0; border-radius: 7px; font-weight: 600; }
.alert-error { background: #fff1f0; color: #c0392b; border: 1px solid #ffd4d0; }
.alert-success { background: #edf9f0; color: #218838; border: 1px solid #ccebd4; }
.form-actions { display: flex; justify-content: flex-end; gap: 12px; padding-top: 22px; border-top: 1px solid #eceef3; }
.save-button { display: inline-flex; align-items: center; gap: 8px; padding: 11px 18px; border: 0; border-radius: 7px; background: #6777ef; color: #fff; font-weight: 700; cursor: pointer; }
.save-button:disabled { opacity: .65; cursor: wait; }

/* Dark Mode */
:global(html.dark) .equipment-card,
:global(body.dark) .equipment-card,
:global([data-theme="dark"]) .equipment-card {
  background-color: #1e293b !important;
  border-color: #334155 !important;
  color: #e2e8f0 !important;
}

:global(html.dark) .page-header,
:global(body.dark) .page-header,
:global([data-theme="dark"]) .page-header {
  border-bottom-color: #334155 !important;
}

:global(html.dark) .page-header h1,
:global(body.dark) .page-header h1,
:global([data-theme="dark"]) .page-header h1 {
  color: #f8fafc !important;
}

:global(html.dark) .field label,
:global(body.dark) .field label,
:global([data-theme="dark"]) .field label {
  color: #cbd5e1 !important;
}

:global(html.dark) .field input,
:global(html.dark) .field select,
:global(body.dark) .field input,
:global(body.dark) .field select,
:global([data-theme="dark"]) .field input,
:global([data-theme="dark"]) .field select {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(html.dark) .stock-summary,
:global(body.dark) .stock-summary,
:global([data-theme="dark"]) .stock-summary {
  background-color: #0f172a !important;
  border: 1px solid #334155 !important;
}

:global(html.dark) .stock-summary span,
:global(body.dark) .stock-summary span,
:global([data-theme="dark"]) .stock-summary span {
  color: #94a3b8 !important;
}

:global(html.dark) .back-button,
:global(html.dark) .cancel-button,
:global(body.dark) .back-button,
:global(body.dark) .cancel-button,
:global([data-theme="dark"]) .back-button,
:global([data-theme="dark"]) .cancel-button {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #cbd5e1 !important;
}

:global(html.dark) .form-actions,
:global(body.dark) .form-actions,
:global([data-theme="dark"]) .form-actions {
  border-top-color: #334155 !important;
}

@media (max-width: 720px) { .equipment-entry-page { padding: 14px; } .equipment-card { padding: 20px; } .page-header { flex-direction: column; } .equipment-form { grid-template-columns: 1fr; } .field-wide { grid-column: auto; } .form-actions { flex-direction: column-reverse; } .back-button, .cancel-button, .save-button { justify-content: center; } }
</style>
