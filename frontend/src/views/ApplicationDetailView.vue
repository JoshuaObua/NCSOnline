<template>
  <main class="application-detail-page">
    <nav class="detail-navbar">
      <button type="button" class="icon-button" title="Back" @click="goBack"><i class="icofont-rounded-left"></i></button>
      <div>
        <small>{{ isAdmin ? 'Application Review' : 'My Application' }}</small>
        <strong>{{ pageTitle }}</strong>
      </div>
      <button v-if="record" type="button" class="download-button" @click="downloadCurrent"><i class="icofont-download"></i> Download form</button>
    </nav>

    <section class="detail-shell">
      <p v-if="message" class="detail-success">{{ message }}</p>
      <p v-if="error" class="detail-error">{{ error }}</p>
      <div v-if="loading" class="detail-empty">Loading application...</div>
      <template v-else-if="record">
        <header class="detail-hero">
          <div>
            <p>{{ source === 'custom' ? 'Dynamic form' : 'Standard form' }}</p>
            <h1>{{ pageTitle }}</h1>
            <span>{{ reference || 'Reference pending' }}</span>
          </div>
          <span class="status-pill" :class="statusClass(record.status)">{{ statusLabel(record.status) }}</span>
        </header>

        <section class="meta-grid">
          <article><small>Applicant</small><strong>{{ record.applicant_name || record.applicant_email || 'Portal user' }}</strong></article>
          <article><small>Payment</small><strong>{{ paymentLabel(record.payment_status) }}</strong></article>
          <article><small>Handled by</small><strong>{{ handlerName }}</strong></article>
          <article><small>Updated</small><strong>{{ formatDate(record.updated_at) }}</strong></article>
        </section>

        <section class="detail-card">
          <header><h2>Payment Proof</h2></header>
          <div class="payment-grid">
            <div v-if="record.payment_reference"><small>URA PRN</small><strong>{{ record.payment_reference }}</strong></div>
            <div v-if="record.payment_amount_ugx"><small>Amount</small><strong>UGX {{ formatMoney(record.payment_amount_ugx) }}</strong></div>
            <FilePreview v-if="record.payment_proof_url" :url="record.payment_proof_url" label="Payment proof" />
            <p v-if="!record.payment_reference && !record.payment_proof_url" class="muted">No payment proof has been recorded for this application.</p>
          </div>
        </section>

        <section class="detail-card">
          <header><h2>Application Responses</h2></header>
          <div v-if="!answerRows.length" class="detail-empty compact">No response data is available.</div>
          <dl v-else class="answer-list">
            <div v-for="row in answerRows" :key="row.key">
              <dt>{{ row.label }}</dt>
              <dd>
                <div v-if="filesFor(row.value).length" class="file-list">
                  <FilePreview v-for="file in filesFor(row.value)" :key="file" :url="file" :label="row.label" />
                </div>
                <template v-else>{{ displayValue(row.value) }}</template>
              </dd>
            </div>
          </dl>
        </section>

        <section v-if="standardFiles.length" class="detail-card">
          <header><h2>Uploaded Files</h2></header>
          <div class="file-list">
            <FilePreview v-for="file in standardFiles" :key="file.url" :url="file.url" :label="file.label" />
          </div>
        </section>

        <section v-if="record.review_notes" class="detail-card note-card">
          <header><h2>Review Notes</h2></header>
          <p>{{ record.review_notes }}</p>
        </section>

        <section v-if="isAdmin" class="detail-card admin-review">
          <header><h2>Review Decision</h2></header>
          <div class="review-grid">
            <label>Application status
              <select v-model="reviewStatus">
                <option v-for="option in reviewOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
              </select>
            </label>
            <label>Payment verification
              <select v-model="paymentStatus">
                <option v-for="option in paymentOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
              </select>
            </label>
            <label class="notes-field">Reason or notes
              <textarea v-model="reviewNotes" rows="4" placeholder="Add a clear reason for query, rejection, or audit record"></textarea>
            </label>
          </div>
          <footer>
            <button type="button" class="secondary-command" :disabled="saving" @click="savePaymentStatus"><i class="icofont-money"></i> Save payment</button>
            <button type="button" class="primary-command" :disabled="saving" @click="saveReview"><i class="icofont-check-circled"></i> Save review</button>
          </footer>
        </section>
      </template>
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  adminGetApplication,
  adminRejectApplicationPayment,
  adminReviewApplication,
  adminVerifyApplicationPayment,
  getMyLegacyApplication,
} from '@/api/applications.js'
import {
  adminGetSubmission,
  adminReviewSubmission,
  adminUpdateSubmissionPaymentStatus,
  portalGetSubmission,
} from '@/api/forms.js'
import FilePreview from '@/components/portal/FilePreview.vue'
import { downloadApplicationForm } from '@/utils/applicationDownload.js'
import { ensureOtikaStyles } from '@/utils/otikaAssets.js'

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const saving = ref(false)
const error = ref('')
const message = ref('')
const record = ref(null)
const reviewStatus = ref('COMPLETE')
const paymentStatus = ref('PROOF_UPLOADED')
const reviewNotes = ref('')

ensureOtikaStyles()

const isAdmin = computed(() => route.path.startsWith('/portal'))
const source = computed(() => route.query.source === 'standard' ? 'standard' : 'custom')
const pageTitle = computed(() => record.value?.template_title || record.value?.title || titleize(record.value?.application_type || record.value?.form_type || 'Application'))
const reference = computed(() => record.value?.submission_reference || record.value?.application_reference || record.value?.reference || '')
const handlerName = computed(() => record.value?.reviewer_name || record.value?.payment_verifier_name || (record.value?.reviewer_id ? 'Assigned reviewer' : 'Not assigned'))
const answerRows = computed(() => Object.entries(parseObject(record.value?.answers ?? record.value?.form_data)).map(([key, value]) => ({ key, value, label: titleize(key) })))
const standardFiles = computed(() => {
  if (source.value !== 'standard' || !record.value) return []
  const files = []
  if (record.value.signed_form_url) files.push({ url: record.value.signed_form_url, label: 'Signed form' })
  for (const attachment of record.value.attachments || []) {
    if (attachment.file_url) files.push({ url: attachment.file_url, label: attachment.field_name || attachment.file_name || 'Attachment' })
  }
  return files
})
const reviewOptions = computed(() => source.value === 'custom'
  ? [
      { value: 'COMPLETE', label: 'Complete' },
      { value: 'NEEDS_INFORMATION', label: 'Queried' },
      { value: 'APPROVED', label: 'Approved' },
      { value: 'REJECTED', label: 'Rejected' },
    ]
  : [
      { value: 'NEEDS_INFORMATION', label: 'Queried' },
      { value: 'APPROVED', label: 'Approved' },
      { value: 'REJECTED', label: 'Rejected' },
    ])
const paymentOptions = [
  { value: 'PROOF_UPLOADED', label: 'Pending' },
  { value: 'VERIFICATION_FAILED', label: 'Verification failed' },
  { value: 'PAID', label: 'Payment Verified' },
]

onMounted(load)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const id = route.params.id
    const res = source.value === 'custom'
      ? (isAdmin.value ? await adminGetSubmission(id) : await portalGetSubmission(id))
      : (isAdmin.value ? await adminGetApplication(id) : await getMyLegacyApplication(id))
    record.value = normalizeRecord(unwrap(res))
    reviewStatus.value = normalizeReviewStatus(record.value.status)
    paymentStatus.value = normalizePaymentStatus(record.value.payment_status)
    reviewNotes.value = record.value.review_notes || ''
    document.title = `${pageTitle.value} - NCS Uganda`
  } catch (err) {
    error.value = apiError(err, 'Could not load this application.')
  } finally {
    loading.value = false
  }
}

async function saveReview() {
  if (!isAdmin.value || !record.value) return
  if (['NEEDS_INFORMATION', 'REJECTED'].includes(reviewStatus.value) && !reviewNotes.value.trim()) {
    error.value = 'Add a reason before querying or rejecting an application.'
    return
  }
  saving.value = true
  error.value = ''
  try {
    if (source.value === 'custom') {
      await adminReviewSubmission(record.value.id, reviewStatus.value, reviewNotes.value.trim())
    } else {
      const action = { NEEDS_INFORMATION: 'request-info', APPROVED: 'approve', REJECTED: 'reject' }[reviewStatus.value]
      if (!action) throw new Error('This status is not available for standard applications.')
      await adminReviewApplication(record.value.id, action, reviewNotes.value.trim())
    }
    message.value = 'Application review saved.'
    await load()
  } catch (err) {
    error.value = apiError(err, err.message || 'Could not save the review.')
  } finally {
    saving.value = false
  }
}

async function savePaymentStatus() {
  if (!isAdmin.value || !record.value) return
  saving.value = true
  error.value = ''
  try {
    if (source.value === 'custom') {
      await adminUpdateSubmissionPaymentStatus(record.value.id, paymentStatus.value)
    } else if (paymentStatus.value === 'PAID') {
      await adminVerifyApplicationPayment(record.value.id)
    } else if (paymentStatus.value === 'VERIFICATION_FAILED') {
      await adminRejectApplicationPayment(record.value.id, reviewNotes.value.trim())
    }
    message.value = 'Payment status saved.'
    await load()
  } catch (err) {
    error.value = apiError(err, 'Could not save the payment status.')
  } finally {
    saving.value = false
  }
}

function downloadCurrent() {
  if (record.value) downloadApplicationForm({ ...record.value, source: source.value, title: pageTitle.value })
}
function goBack() {
  router.push(isAdmin.value ? '/portal?section=applications' : '/dashboard?section=applications')
}
function normalizeRecord(value) {
  const data = unwrap(value)
  return { ...data, title: data.template_title || titleize(data.application_type || data.form_type || 'Application') }
}
function parseObject(value) {
  if (!value) return {}
  if (typeof value === 'object') return value
  try { return JSON.parse(value) } catch { return {} }
}
function filesFor(value) {
  if (Array.isArray(value)) return value.filter(isFileValue)
  return isFileValue(value) ? [value] : []
}
function isFileValue(value) {
  return typeof value === 'string' && /(\.pdf|\.png|\.jpe?g|\.webp|\/uploads\/|\/media\/)/i.test(value)
}
function displayValue(value) {
  if (Array.isArray(value)) return value.join(', ') || '-'
  if (value && typeof value === 'object') return JSON.stringify(value)
  return String(value ?? '-')
}
function normalizeReviewStatus(status) {
  if (status === 'NEEDS_INFORMATION') return 'NEEDS_INFORMATION'
  if (status === 'APPROVED') return 'APPROVED'
  if (status === 'REJECTED') return 'REJECTED'
  return source.value === 'custom' ? 'COMPLETE' : 'NEEDS_INFORMATION'
}
function normalizePaymentStatus(status) {
  if (status === 'PAID') return 'PAID'
  if (status === 'VERIFICATION_FAILED' || status === 'REJECTED') return 'VERIFICATION_FAILED'
  return 'PROOF_UPLOADED'
}
function paymentLabel(status) {
  if (status === 'PROOF_UPLOADED') return 'Pending'
  if (status === 'PAID') return 'Payment Verified'
  if (status === 'VERIFICATION_FAILED') return 'Verification failed'
  return titleize(status || 'not required')
}
function statusLabel(status) {
  if (status === 'NEEDS_INFORMATION') return 'Queried'
  return titleize(status)
}
function statusClass(status) {
  if (status === 'APPROVED' || status === 'COMPLETE') return 'green'
  if (status === 'REJECTED') return 'red'
  if (status === 'NEEDS_INFORMATION' || status === 'PENDING_PAYMENT') return 'amber'
  return 'blue'
}
function titleize(value) { return String(value || '').toLowerCase().replaceAll('_', ' ').replaceAll('-', ' ').replace(/\b\w/g, char => char.toUpperCase()) }
function formatDate(value) { return value ? new Intl.DateTimeFormat('en-UG', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) : '-' }
function formatMoney(value) { return new Intl.NumberFormat('en-UG', { maximumFractionDigits: 0 }).format(Number(value || 0)) }
function unwrap(value) { return value?.data?.data ?? value?.data ?? value ?? {} }
function apiError(err, fallback) { return err.response?.data?.error?.message || fallback }
</script>

<style scoped>
.application-detail-page{min-height:100vh;background:#f4f6f9;color:#34395e}.detail-navbar{position:sticky;top:0;z-index:20;display:flex;align-items:center;gap:12px;min-height:68px;padding:10px 24px;background:#fff;box-shadow:0 4px 25px rgba(0,0,0,.08)}.detail-navbar>div{min-width:0}.detail-navbar small,.detail-navbar strong{display:block}.detail-navbar small{color:#98a6ad;font-size:10px;font-weight:800;text-transform:uppercase}.detail-navbar strong{overflow:hidden;color:#34395e;text-overflow:ellipsis;white-space:nowrap}.icon-button{display:grid;place-items:center;width:38px;height:38px;border:0;border-radius:30px;background:#f4f6f9;color:#6777ef}.download-button,.primary-command,.secondary-command{display:inline-flex;align-items:center;justify-content:center;gap:7px;min-height:40px;margin-left:auto;padding:0 15px;border:0;border-radius:30px;background:#6777ef;color:#fff;font-size:12px;font-weight:700}.detail-shell{width:min(1120px,calc(100% - 28px));margin:0 auto;padding:26px 0 42px}.detail-success,.detail-error{margin:0 0 14px;padding:12px 14px;border-radius:3px;background:#e8f7f0;color:#47c363;font-size:12px}.detail-error{background:#fdeaea;color:#fc544b}.detail-empty{padding:48px;border:1px dashed #e4e6fc;border-radius:3px;background:#fdfdff;text-align:center;color:#98a6ad}.detail-empty.compact{padding:26px}.detail-hero{display:flex;align-items:flex-start;justify-content:space-between;gap:18px;margin-bottom:16px;padding:22px;border-radius:3px;background:#fff;box-shadow:0 4px 25px rgba(0,0,0,.08)}.detail-hero p{margin:0;color:#6777ef;font-size:11px;font-weight:800;text-transform:uppercase}.detail-hero h1{margin:4px 0;color:#34395e;font-size:26px}.detail-hero span:not(.status-pill){color:#6c757d}.status-pill{display:inline-flex;flex:0 0 auto;padding:6px 10px;border-radius:30px;font-size:10px;font-weight:800;text-transform:uppercase}.status-pill.blue{background:#e8edff;color:#6777ef}.status-pill.green{background:#e8f7f0;color:#47c363}.status-pill.amber{background:#fff4e6;color:#ffa426}.status-pill.red{background:#fdeaea;color:#fc544b}.meta-grid{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:12px;margin-bottom:16px}.meta-grid article,.detail-card{min-width:0;border-radius:3px;background:#fff;box-shadow:0 4px 25px rgba(0,0,0,.08)}.meta-grid article{padding:16px}.meta-grid small{display:block;color:#98a6ad;font-size:10px;font-weight:800;text-transform:uppercase}.meta-grid strong{display:block;margin-top:4px;overflow-wrap:anywhere;color:#34395e}.detail-card{display:grid;gap:14px;margin-bottom:16px;padding:20px}.detail-card>header h2{margin:0;color:#34395e;font-size:17px}.payment-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:12px}.payment-grid>div:not(.file-preview){padding:12px;border-radius:3px;background:#fdfdff}.payment-grid small{display:block;color:#98a6ad;font-size:10px;font-weight:800;text-transform:uppercase}.payment-grid strong{color:#34395e}.muted{margin:0;color:#98a6ad}.answer-list{display:grid;gap:0;margin:0}.answer-list>div{display:grid;grid-template-columns:minmax(180px,.34fr) minmax(0,1fr);gap:14px;padding:13px 0;border-bottom:1px solid #f4f6f9}.answer-list dt{color:#98a6ad;font-size:11px;font-weight:800;text-transform:uppercase}.answer-list dd{min-width:0;margin:0;color:#34395e;white-space:pre-wrap;overflow-wrap:anywhere}.file-list{display:grid;gap:10px}.note-card p{margin:0;color:#6c757d;line-height:1.55}.review-grid{display:grid;grid-template-columns:1fr 1fr;gap:12px}.review-grid label{display:grid;gap:6px;color:#34395e;font-size:12px;font-weight:700}.review-grid select,.review-grid textarea{width:100%;padding:10px 12px;border:1px solid #e4e6fc;border-radius:3px;background:#fdfdff;color:#495057}.notes-field{grid-column:1/-1}.admin-review footer{display:flex;justify-content:flex-end;gap:8px}.secondary-command{margin-left:0;background:#f4f6f9;color:#34395e}.primary-command{margin-left:0}:global(.dark) .application-detail-page{background:#0f172a;color:#e5e7eb}:global(.dark) .detail-navbar,:global(.dark) .detail-hero,:global(.dark) .meta-grid article,:global(.dark) .detail-card{background:#1f2937;color:#cbd5e1}:global(.dark) .detail-navbar strong,:global(.dark) .detail-hero h1,:global(.dark) .meta-grid strong,:global(.dark) .detail-card>header h2,:global(.dark) .answer-list dd,:global(.dark) .review-grid label{color:#f8fafc}:global(.dark) .payment-grid>div:not(.file-preview),:global(.dark) .review-grid select,:global(.dark) .review-grid textarea,:global(.dark) .detail-empty{background:#111827;border-color:#334155;color:#e5e7eb}@media(max-width:860px){.detail-navbar{padding:10px 14px}.download-button{margin-left:0}.detail-navbar{flex-wrap:wrap}.detail-navbar>div{flex:1 1 220px}.meta-grid,.payment-grid,.review-grid{grid-template-columns:1fr}.answer-list>div{grid-template-columns:1fr;gap:5px}.detail-hero{flex-direction:column}.admin-review footer{flex-direction:column}.admin-review footer button{width:100%}}@media(max-width:520px){.detail-shell{width:calc(100% - 20px);padding-top:18px}.detail-hero,.detail-card{padding:16px}.detail-hero h1{font-size:21px}.download-button{width:100%}}
</style>
