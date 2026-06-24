<template>
  <LayoutDefault title="Form Submissions">
    <div class="p-6 max-w-7xl mx-auto">
      <!-- Toast -->
      <Transition name="toast">
        <div v-if="toast" class="fixed top-5 right-5 z-[100] px-4 py-3 bg-gray-900 text-white text-sm rounded-xl shadow-xl">{{ toast }}</div>
      </Transition>

      <div class="flex items-center justify-between mb-5">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Form Submissions</h2>
          <p class="text-xs text-gray-500 mt-1">Submissions from applicants. You can only see entries for forms in your department.</p>
        </div>
      </div>

      <!-- Filters -->
      <div class="flex flex-wrap gap-3 mb-4">
        <select v-model="filters.template_id" @change="load" class="text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white">
          <option value="">All forms in my department</option>
          <option v-for="f in forms" :key="f.id" :value="f.id">{{ f.title }}</option>
        </select>
        <select v-model="filters.status" @change="load" class="text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white">
          <option value="">All statuses</option>
          <option value="SUBMITTED">Submitted</option>
          <option value="UNDER_REVIEW">Under review</option>
          <option value="NEEDS_INFORMATION">Needs information</option>
          <option value="APPROVED">Approved</option>
          <option value="REJECTED">Rejected</option>
        </select>
      </div>

      <!-- Table -->
      <div class="bg-white border border-gray-200 rounded-xl overflow-hidden">
        <div v-if="loading" class="p-6 space-y-2">
          <div v-for="i in 4" :key="i" class="h-12 bg-gray-100 rounded animate-pulse"/>
        </div>
        <div v-else-if="!submissions.length" class="p-10 text-center text-sm text-gray-500">
          No submissions match these filters.
        </div>
        <table v-else class="w-full text-sm">
          <thead class="bg-gray-50 text-xs uppercase tracking-wider text-gray-500">
            <tr>
              <th class="px-4 py-3 text-left">Reference</th>
              <th class="px-4 py-3 text-left">Form</th>
              <th class="px-4 py-3 text-left">Applicant</th>
              <th class="px-4 py-3 text-left">Status</th>
              <th class="px-4 py-3 text-left">Payment</th>
              <th class="px-4 py-3 text-left">Submitted</th>
              <th class="px-4 py-3"></th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="s in submissions" :key="s.id" class="hover:bg-gray-50 cursor-pointer" @click="openDetail(s)">
              <td class="px-4 py-3 font-mono text-xs text-gray-700">{{ s.submission_reference || '—' }}</td>
              <td class="px-4 py-3">{{ s.template_title }}</td>
              <td class="px-4 py-3">
                <div class="text-gray-900">{{ s.applicant_name || '—' }}</div>
                <div class="text-xs text-gray-400">{{ s.applicant_email }}</div>
              </td>
              <td class="px-4 py-3"><span :class="statusClass(s.status)" class="text-[10px] font-bold uppercase tracking-wider px-2 py-1 rounded-full">{{ s.status }}</span></td>
              <td class="px-4 py-3"><span :class="paymentClass(s.payment_status)" class="text-[10px] font-bold uppercase tracking-wider px-2 py-1 rounded-full">{{ s.payment_status }}</span></td>
              <td class="px-4 py-3 text-xs text-gray-500">{{ s.submitted_at ? formatDate(s.submitted_at) : '—' }}</td>
              <td class="px-4 py-3 text-right"><span class="text-primary-600 text-xs">View →</span></td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Detail drawer -->
      <Transition name="drawer">
        <div v-if="detail" class="fixed inset-0 z-50 flex">
          <div class="flex-1 bg-black/30" @click="detail = null"/>
          <aside class="w-full max-w-2xl bg-white shadow-2xl overflow-y-auto">
            <header class="px-6 py-4 border-b border-gray-100 flex items-start justify-between">
              <div>
                <h3 class="text-lg font-semibold text-gray-900">{{ detail.template_title }}</h3>
                <p class="text-xs text-gray-500 mt-1 font-mono">{{ detail.submission_reference || detail.id }}</p>
              </div>
              <button @click="detail = null" class="text-gray-400 hover:text-gray-700 text-2xl leading-none">×</button>
            </header>

            <div class="px-6 py-5 space-y-5">
              <div class="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <div class="text-[10px] uppercase tracking-wider text-gray-400">Applicant</div>
                  <div class="text-gray-900">{{ detail.applicant_name }}</div>
                  <div class="text-xs text-gray-500">{{ detail.applicant_email }}</div>
                </div>
                <div>
                  <div class="text-[10px] uppercase tracking-wider text-gray-400">Status</div>
                  <span :class="statusClass(detail.status)" class="text-[10px] font-bold uppercase tracking-wider px-2 py-1 rounded-full">{{ detail.status }}</span>
                </div>
                <div>
                  <div class="text-[10px] uppercase tracking-wider text-gray-400">Payment</div>
                  <span :class="paymentClass(detail.payment_status)" class="text-[10px] font-bold uppercase tracking-wider px-2 py-1 rounded-full">{{ detail.payment_status }}</span>
                  <span v-if="detail.payment_reference" class="ml-2 text-xs text-gray-500">{{ detail.payment_reference }}</span>
                </div>
                <div>
                  <div class="text-[10px] uppercase tracking-wider text-gray-400">Submitted</div>
                  <div class="text-xs text-gray-700">{{ detail.submitted_at ? formatDate(detail.submitted_at) : '—' }}</div>
                </div>
              </div>

              <section>
                <h4 class="text-xs uppercase tracking-wider text-gray-500 font-semibold mb-2">Answers</h4>
                <div class="bg-gray-50 border border-gray-200 rounded-lg p-4 space-y-2 text-sm">
                  <div v-for="(val, key) in detail.answers || {}" :key="key">
                    <div class="text-[10px] uppercase tracking-wider text-gray-400">{{ key }}</div>
                    <div class="text-gray-900 break-words whitespace-pre-wrap">{{ formatVal(val) }}</div>
                  </div>
                  <div v-if="!Object.keys(detail.answers || {}).length" class="text-xs text-gray-400">No answers recorded.</div>
                </div>
              </section>

              <section v-if="detail.review_notes">
                <h4 class="text-xs uppercase tracking-wider text-gray-500 font-semibold mb-2">Reviewer Notes</h4>
                <p class="text-sm text-gray-800 bg-amber-50 border border-amber-200 rounded-lg p-3">{{ detail.review_notes }}</p>
              </section>

              <section>
                <h4 class="text-xs uppercase tracking-wider text-gray-500 font-semibold mb-2">Take action</h4>
                <textarea v-model="reviewNotes" rows="3" placeholder="Notes (required for reject / request info)"
                  class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 mb-3"></textarea>
                <div class="flex flex-wrap gap-2">
                  <button @click="act('APPROVED')" :disabled="actBusy" class="text-xs font-medium bg-green-600 hover:bg-green-700 text-white px-3 py-1.5 rounded-lg disabled:opacity-50">Approve</button>
                  <button @click="act('REJECTED')" :disabled="actBusy" class="text-xs font-medium bg-red-600 hover:bg-red-700 text-white px-3 py-1.5 rounded-lg disabled:opacity-50">Reject</button>
                  <button @click="act('NEEDS_INFORMATION')" :disabled="actBusy" class="text-xs font-medium bg-amber-500 hover:bg-amber-600 text-white px-3 py-1.5 rounded-lg disabled:opacity-50">Request Info</button>
                  <button @click="act('UNDER_REVIEW')" :disabled="actBusy" class="text-xs font-medium bg-gray-700 hover:bg-gray-800 text-white px-3 py-1.5 rounded-lg disabled:opacity-50">Mark Under Review</button>
                  <button v-if="detail.payment_status === 'PROOF_UPLOADED'" @click="verifyPayment" :disabled="actBusy"
                    class="text-xs font-medium bg-blue-600 hover:bg-blue-700 text-white px-3 py-1.5 rounded-lg disabled:opacity-50">Verify Payment</button>
                </div>
              </section>
            </div>
          </aside>
        </div>
      </Transition>
    </div>
  </LayoutDefault>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import {
  adminListSubmissions, adminGetSubmission,
  adminReviewSubmission, adminVerifySubmissionPayment,
  adminListForms,
} from '@/api/forms'

const submissions = ref([])
const forms = ref([])
const loading = ref(false)
const detail = ref(null)
const reviewNotes = ref('')
const actBusy = ref(false)
const toast = ref('')
const filters = reactive({ template_id: '', status: '' })

function showToast(msg) { toast.value = msg; setTimeout(() => (toast.value = ''), 2500) }
function formatDate(d) { return new Date(d).toLocaleString() }
function formatVal(v) {
  if (v === null || v === undefined || v === '') return '—'
  if (Array.isArray(v)) return v.join(', ')
  if (typeof v === 'object') return JSON.stringify(v, null, 2)
  return String(v)
}
function statusClass(s) {
  return {
    DRAFT: 'bg-gray-200 text-gray-700',
    SUBMITTED: 'bg-blue-100 text-blue-800',
    UNDER_REVIEW: 'bg-amber-100 text-amber-800',
    NEEDS_INFORMATION: 'bg-orange-100 text-orange-800',
    APPROVED: 'bg-green-100 text-green-800',
    REJECTED: 'bg-red-100 text-red-800',
  }[s] || 'bg-gray-200 text-gray-700'
}
function paymentClass(s) {
  return {
    UNPAID: 'bg-gray-200 text-gray-700',
    PROOF_UPLOADED: 'bg-amber-100 text-amber-800',
    PAID: 'bg-green-100 text-green-800',
  }[s] || 'bg-gray-200 text-gray-700'
}

async function load() {
  loading.value = true
  try {
    const r = await adminListSubmissions(filters)
    submissions.value = (r.data || []).map(s => ({
      ...s,
      answers: tryParse(s.answers),
    }))
  } finally {
    loading.value = false
  }
}

function tryParse(v) {
  if (!v) return {}
  if (typeof v === 'object') return v
  try { return JSON.parse(v) } catch { return {} }
}

async function loadForms() {
  try { forms.value = await adminListForms('') } catch {}
}

async function openDetail(s) {
  reviewNotes.value = ''
  try {
    const full = await adminGetSubmission(s.id)
    detail.value = { ...full, answers: tryParse(full.answers) }
  } catch (e) {
    showToast('Failed to load submission')
  }
}

async function act(status) {
  if (!detail.value) return
  if ((status === 'REJECTED' || status === 'NEEDS_INFORMATION') && !reviewNotes.value.trim()) {
    showToast('Notes are required for this action'); return
  }
  actBusy.value = true
  try {
    await adminReviewSubmission(detail.value.id, status, reviewNotes.value)
    showToast(`Marked ${status}`)
    detail.value = null
    load()
  } catch (e) {
    showToast(e?.response?.data?.message || 'Action failed')
  } finally {
    actBusy.value = false
  }
}

async function verifyPayment() {
  if (!detail.value) return
  actBusy.value = true
  try {
    await adminVerifySubmissionPayment(detail.value.id)
    showToast('Payment verified')
    detail.value = null
    load()
  } catch (e) {
    showToast('Failed to verify')
  } finally {
    actBusy.value = false
  }
}

onMounted(async () => { await loadForms(); load() })
</script>

<style scoped>
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(-8px); }
.toast-enter-active, .toast-leave-active { transition: all .25s ease; }
.drawer-enter-from { opacity: 0; }
.drawer-enter-from aside { transform: translateX(100%); }
.drawer-leave-to { opacity: 0; }
.drawer-leave-to aside { transform: translateX(100%); }
.drawer-enter-active, .drawer-leave-active { transition: all .25s ease; }
.drawer-enter-active aside, .drawer-leave-active aside { transition: transform .3s ease; }
</style>
