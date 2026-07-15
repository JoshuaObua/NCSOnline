<template>
  <main class="otika-cms application-wizard-page">
    <div class="wizard-navbar">
      <router-link to="/dashboard" class="wizard-brand">
        <img src="/main-logo.png" alt="NCS" />
        <span>Applicant Portal</span>
      </router-link>
      <div class="wizard-navbar-actions">
        <ThemeToggle class="wizard-icon-button" />
        <button type="button" class="wizard-icon-button" title="Back to dashboard" @click="backToDashboard">
          <i class="icofont-dashboard-web"></i>
        </button>
        <button type="button" class="wizard-profile" @click="backToDashboard">
          <span>{{ initials }}</span>
          <b>{{ firstName }}</b>
        </button>
      </div>
    </div>

    <section class="wizard-workspace">
      <header class="wizard-header">
        <div>
          <button type="button" class="back-link" @click="backToDashboard"><i class="icofont-rounded-left"></i> Dashboard</button>
          <p>{{ form.department_name || 'NCS Application' }}</p>
          <h1>{{ form.title || 'Application Form' }}</h1>
          <span>{{ form.description || 'Complete and submit this application online.' }}</span>
        </div>
        <div class="fee-panel">
          <small>Application fee</small>
          <strong>{{ Number(form.price_ugx) ? `UGX ${formatMoney(form.price_ugx)}` : 'Free' }}</strong>
        </div>
      </header>

      <p v-if="success" class="wizard-success">{{ success }}</p>
      <p v-if="formError" class="wizard-error">{{ formError }}</p>
      <div v-if="loading" class="wizard-loading">Loading application form...</div>

      <div v-else class="wizard-layout">
        <aside class="wizard-steps" aria-label="Application steps">
          <button
            v-for="(step, index) in steps"
            :key="step.id"
            type="button"
            :class="{ active: index === currentStepIndex, complete: index < currentStepIndex }"
            @click="goToStep(index)"
          >
            <span>{{ index + 1 }}</span>
            <div>
              <strong>{{ step.title }}</strong>
              <small>{{ step.subtitle || `${step.fields.length} field${step.fields.length === 1 ? '' : 's'}` }}</small>
            </div>
          </button>
        </aside>

        <form class="wizard-form" @submit.prevent="submitApplication">
          <section class="wizard-step-panel">
            <div class="step-heading">
              <small>Step {{ currentStepIndex + 1 }} of {{ steps.length }}</small>
              <h2>{{ currentStep.title }}</h2>
              <p v-if="currentStep.subtitle">{{ currentStep.subtitle }}</p>
              <p v-if="currentStep.description">{{ currentStep.description }}</p>
            </div>

            <div v-if="!currentStep.fields.length" class="empty-step">This step has no fields yet.</div>
            <div v-for="field in currentStep.fields" :key="field.id || field.field_key" class="wizard-field">
              <label :for="fieldId(field)">{{ field.label }} <b v-if="field.is_required">*</b></label>
              <p v-if="field.help_text">{{ field.help_text }}</p>
              <textarea
                v-if="field.field_type === 'long_text'"
                :id="fieldId(field)"
                v-model="formAnswers[field.field_key]"
                :placeholder="field.placeholder"
                rows="5"
              ></textarea>
              <select v-else-if="field.field_type === 'dropdown'" :id="fieldId(field)" v-model="formAnswers[field.field_key]">
                <option value="">Select an option</option>
                <option v-for="option in fieldOptions(field)" :key="option" :value="option">{{ option }}</option>
              </select>
              <div v-else-if="field.field_type === 'radio'" class="choice-list">
                <label v-for="option in fieldOptions(field)" :key="option">
                  <input v-model="formAnswers[field.field_key]" type="radio" :name="field.field_key" :value="option" />
                  {{ option }}
                </label>
              </div>
              <div v-else-if="field.field_type === 'checkbox'" class="choice-list">
                <label v-for="option in fieldOptions(field)" :key="option">
                  <input v-model="formAnswers[field.field_key]" type="checkbox" :value="option" />
                  {{ option }}
                </label>
              </div>
              <div
                v-else-if="['file', 'image'].includes(field.field_type)"
                class="dropify-zone"
                :class="{ filled: !!formAnswers[field.field_key], uploading: uploadingField === field.field_key }"
                @dragover.prevent
                @drop.prevent="dropFieldFile(field, $event)"
              >
                <input :id="fieldId(field)" class="dropify-input" type="file" :accept="fieldAccept(field)" @change="uploadFieldFile(field, $event)" />
                <span class="dropify-icon"><i :class="field.field_type === 'image' ? 'icofont-image' : 'icofont-upload-alt'"></i></span>
                <strong>{{ uploadingField === field.field_key ? 'Uploading...' : (formAnswers[field.field_key] ? 'Attachment uploaded' : 'Drop file here or choose one') }}</strong>
                <small>{{ field.field_type === 'image' ? 'Image files are accepted' : 'PDF, JPG, PNG, JPEG or configured file type' }}</small>
                <a v-if="formAnswers[field.field_key]" :href="mediaUrl(formAnswers[field.field_key])" target="_blank" rel="noopener">Preview attachment</a>
                <button v-if="formAnswers[field.field_key]" type="button" class="proof-clear" @click.stop="clearFieldFile(field)">Remove</button>
              </div>
              <input
                v-else
                :id="fieldId(field)"
                v-model="formAnswers[field.field_key]"
                :type="inputType(field.field_type)"
                :placeholder="field.placeholder"
              />
            </div>

            <section v-if="isLastStep && form.price_ugx > 0 && paymentRequired" class="payment-section">
              <h3>Payment proof</h3>
              <p>Use one proof method: enter the URA PRN for a payment already made, or upload a scanned receipt/proof file.</p>
              <label>URA PRN (Payment Reference Number)
                <input
                  v-model.trim="paymentReference"
                  :disabled="!!paymentProofURL || uploadingPaymentProof"
                  placeholder="Enter PRN generated from the URA portal"
                  @input="clearPaymentProofFile"
                />
              </label>
              <div
                class="dropify-zone payment-proof-drop"
                :class="{ filled: !!paymentProofURL, uploading: uploadingPaymentProof, disabled: !!paymentReference }"
                @dragover.prevent
                @drop.prevent="dropPaymentProofFile"
              >
                <input
                  ref="paymentProofInput"
                  class="dropify-input"
                  type="file"
                  accept=".pdf,.png,.jpg,.jpeg,application/pdf,image/png,image/jpeg"
                  :disabled="!!paymentReference || uploadingPaymentProof"
                  @change="uploadPaymentProofFile"
                />
                <span class="dropify-icon"><i class="icofont-paperclip"></i></span>
                <strong>{{ uploadingPaymentProof ? 'Uploading proof...' : (paymentProofURL ? 'Payment proof uploaded' : 'Drop receipt or proof here') }}</strong>
                <small>PDF, PNG, JPG or JPEG</small>
                <a v-if="paymentProofURL" :href="mediaUrl(paymentProofURL)" target="_blank" rel="noopener">{{ paymentProofFileName || 'Preview proof' }}</a>
                <button v-if="paymentProofURL" type="button" class="proof-clear" :disabled="uploadingPaymentProof" @click.stop="clearPaymentProofFile">Remove</button>
              </div>
              <label>Amount paid (UGX)
                <input v-model.number="paymentAmount" type="number" min="1" />
              </label>
            </section>
          </section>

          <footer class="wizard-actions">
            <button type="button" class="secondary-command" :disabled="formSaving" @click="saveDraft"><i class="icofont-save"></i> Save draft</button>
            <div>
              <button type="button" class="secondary-command" :disabled="currentStepIndex === 0 || formSaving" @click="previousStep"><i class="icofont-rounded-left"></i> Back</button>
              <button v-if="!isLastStep" type="button" class="primary-command" :disabled="formSaving || !!uploadingField" @click="nextStep">Next <i class="icofont-rounded-right"></i></button>
              <button v-else type="submit" class="primary-command" :disabled="formSaving || !!uploadingField || uploadingPaymentProof"><i class="icofont-paper-plane"></i> {{ formSaving ? 'Saving...' : 'Submit application' }}</button>
            </div>
          </footer>
        </form>
      </div>
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCurrentUser } from '@/api/auth.js'
import { mediaUrl } from '@/api/client.js'
import * as cms from '@/api/cms.js'
import {
  portalGetForm,
  portalListSubmissions,
  portalSaveDraft,
  portalSubmit,
  portalUploadPaymentProof,
} from '@/api/forms.js'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'
import { buildSectionSteps, parseSubmissionAnswers } from '@/utils/formBuilder.js'
import { ensureOtikaStyles } from '@/utils/otikaAssets.js'

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const formSaving = ref(false)
const uploadingField = ref('')
const formError = ref('')
const success = ref('')
const form = reactive({ fields: [], sections: [] })
const profile = reactive({})
const formAnswers = reactive({})
const activeSubmission = ref(null)
const currentStepIndex = ref(0)
const paymentReference = ref('')
const paymentAmount = ref(0)
const paymentProofURL = ref('')
const paymentProofFileName = ref('')
const paymentProofInput = ref(null)
const uploadingPaymentProof = ref(false)

ensureOtikaStyles()

const steps = computed(() => buildSectionSteps(form))
const currentStep = computed(() => steps.value[currentStepIndex.value] || { title: 'Application Details', fields: [] })
const isLastStep = computed(() => currentStepIndex.value >= steps.value.length - 1)
const firstName = computed(() => profile.first_name || String(profile.email || 'User').split('@')[0])
const initials = computed(() => `${profile.first_name?.[0] || ''}${profile.last_name?.[0] || ''}`.toUpperCase() || 'U')
const paymentRequired = computed(() => Number(form.price_ugx) > 0 && !['PAID', 'PROOF_UPLOADED'].includes(activeSubmission.value?.payment_status))

onMounted(loadWizard)

async function loadWizard() {
  loading.value = true
  formError.value = ''
  try {
    const [user, template] = await Promise.all([
      getCurrentUser(),
      portalGetForm(route.params.slug),
    ])
    Object.assign(profile, unwrap(user))
    Object.assign(form, unwrap(template))
    initializeAnswers()
    const submissions = asList(await portalListSubmissions({ page: 1, per_page: 200 }))
    const pending = submissions.find(item => item.template_id === form.id && isPendingSubmission(item))
    if (pending) {
      formError.value = `You already have a pending application for this form (${pending.submission_reference || titleize(pending.status)}).`
      return
    }
    const draft = submissions.find(item => item.template_id === form.id && ['DRAFT', 'NEEDS_INFORMATION'].includes(item.status))
    activeSubmission.value = draft || null
    Object.assign(formAnswers, parseSubmissionAnswers(draft?.answers))
    initializeAnswers()
    paymentReference.value = draft?.payment_reference || ''
    paymentProofURL.value = draft?.payment_proof_url || ''
    paymentAmount.value = draft?.payment_amount_ugx || form.price_ugx || 0
    document.title = `${form.title || 'Application'} - NCS Uganda`
  } catch (error) {
    formError.value = apiError(error, 'Could not load this application form.')
  } finally {
    loading.value = false
  }
}

function initializeAnswers() {
  for (const step of steps.value) {
    for (const field of step.fields) {
      if (field.field_type === 'checkbox') {
        if (!Array.isArray(formAnswers[field.field_key])) formAnswers[field.field_key] = []
      } else if (formAnswers[field.field_key] == null) {
        formAnswers[field.field_key] = ''
      }
    }
  }
}

async function nextStep() {
  if (!validateFields(currentStep.value.fields)) return
  await saveDraft(false)
  if (!formError.value) currentStepIndex.value = Math.min(currentStepIndex.value + 1, steps.value.length - 1)
}

function previousStep() {
  currentStepIndex.value = Math.max(0, currentStepIndex.value - 1)
}

function goToStep(index) {
  if (index <= currentStepIndex.value) {
    currentStepIndex.value = index
    return
  }
  if (validateFields(currentStep.value.fields)) currentStepIndex.value = index
}

async function saveDraft(showMessage = true) {
  formSaving.value = true
  formError.value = ''
  try {
    activeSubmission.value = await portalSaveDraft(form.id, { ...formAnswers })
    if (showMessage) success.value = 'Your draft was saved.'
  } catch (error) {
    formError.value = apiError(error, 'Could not save your draft.')
  } finally {
    formSaving.value = false
  }
}

async function submitApplication() {
  if (!validateAllSteps()) return
  formSaving.value = true
  formError.value = ''
  try {
    let submission = await portalSaveDraft(form.id, { ...formAnswers })
    activeSubmission.value = submission
    if (Number(form.price_ugx) > 0 && !['PAID', 'PROOF_UPLOADED'].includes(submission.payment_status)) {
      if (!validPaymentProof()) throw new Error('Enter a URA PRN or upload one payment proof file, then enter the amount paid.')
      await portalUploadPaymentProof(submission.id, {
        payment_reference: paymentReference.value.trim(),
        payment_proof_url: paymentProofURL.value.trim(),
        payment_amount_ugx: Number(paymentAmount.value),
      })
    }
    await portalSubmit(submission.id)
    success.value = 'Your application was submitted successfully.'
    router.replace('/dashboard?section=applications')
  } catch (error) {
    formError.value = apiError(error, error.message || 'Could not submit your application.')
  } finally {
    formSaving.value = false
  }
}

async function uploadFieldFile(field, event) {
  const file = event.target.files?.[0]
  if (event?.target) event.target.value = ''
  await uploadFieldFileObject(field, file)
}

async function dropFieldFile(field, event) {
  const file = event.dataTransfer?.files?.[0]
  await uploadFieldFileObject(field, file)
}

async function uploadFieldFileObject(field, file) {
  if (!file) return
  uploadingField.value = field.field_key
  formError.value = ''
  try {
    const res = await cms.uploadMedia(file, 'application')
    const uploaded = unwrap(res)
    formAnswers[field.field_key] = uploaded.url || uploaded.file_url || uploaded.path
  } catch (error) {
    formError.value = apiError(error, 'Could not upload this file.')
  } finally {
    uploadingField.value = ''
  }
}

async function uploadPaymentProofFile(event) {
  const file = event.target.files?.[0]
  await uploadPaymentProofObject(file)
}

async function dropPaymentProofFile(event) {
  if (paymentReference.value || uploadingPaymentProof.value) return
  const file = event.dataTransfer?.files?.[0]
  await uploadPaymentProofObject(file)
}

async function uploadPaymentProofObject(file) {
  if (!file) return
  formError.value = ''
  if (!validProofFile(file)) {
    formError.value = 'Upload payment proof as PDF, PNG, JPG, or JPEG.'
    clearPaymentProofFile()
    return
  }
  uploadingPaymentProof.value = true
  paymentReference.value = ''
  try {
    const res = await cms.uploadMedia(file, 'application')
    const uploaded = unwrap(res)
    paymentProofURL.value = uploaded.url || uploaded.file_url || uploaded.path || ''
    paymentProofFileName.value = file.name
    if (!paymentProofURL.value) throw new Error('Upload did not return a file URL.')
  } catch (error) {
    clearPaymentProofFile()
    formError.value = apiError(error, 'Could not upload the payment proof file.')
  } finally {
    uploadingPaymentProof.value = false
  }
}

function clearFieldFile(field) {
  formAnswers[field.field_key] = field.field_type === 'checkbox' ? [] : ''
}

function clearPaymentProofFile() {
  paymentProofURL.value = ''
  paymentProofFileName.value = ''
  if (paymentProofInput.value) paymentProofInput.value.value = ''
}

function validateAllSteps() {
  for (let index = 0; index < steps.value.length; index += 1) {
    if (!validateFields(steps.value[index].fields)) {
      currentStepIndex.value = index
      return false
    }
  }
  if (Number(form.price_ugx) > 0 && paymentRequired.value && !validPaymentProof()) {
    currentStepIndex.value = steps.value.length - 1
    formError.value = 'Enter a URA PRN or upload one payment proof file, then enter the amount paid.'
    return false
  }
  return true
}

function validateFields(fields) {
  const missing = fields.find(field => {
    if (!field.is_required) return false
    const value = formAnswers[field.field_key]
    return Array.isArray(value) ? value.length === 0 : !String(value ?? '').trim()
  })
  if (missing) {
    formError.value = `${missing.label} is required.`
    return false
  }
  formError.value = ''
  return true
}

function backToDashboard() {
  router.push('/dashboard')
}
function fieldId(field) { return `field-${field.field_key}` }
function fieldOptions(field) { return Array.isArray(field.config?.options) ? field.config.options : [] }
function fieldAccept(field) {
  const accept = field.config?.accept
  return Array.isArray(accept) ? accept.join(',') : (field.field_type === 'image' ? 'image/*' : '.pdf,.jpg,.jpeg,.png,.webp')
}
function inputType(type) { return ({ phone: 'tel', email: 'email', number: 'number', date: 'date' })[type] || 'text' }
function validPaymentProof() {
  const hasReference = !!paymentReference.value.trim()
  const hasFile = !!paymentProofURL.value.trim()
  return hasReference !== hasFile && Number(paymentAmount.value) > 0
}
function validProofFile(file) {
  const name = String(file?.name || '').toLowerCase()
  return ['application/pdf', 'image/png', 'image/jpeg'].includes(file?.type) || /\.(pdf|png|jpe?g)$/.test(name)
}
function isPendingSubmission(item) {
  return item?.status && !['DRAFT', 'NEEDS_INFORMATION', 'APPROVED', 'REJECTED'].includes(item.status)
}
function formatMoney(value) { return new Intl.NumberFormat('en-UG', { maximumFractionDigits: 0 }).format(Number(value || 0)) }
function unwrap(value) { return value?.data?.data ?? value?.data ?? value ?? {} }
function asList(value) { const data = unwrap(value); return Array.isArray(data) ? data : [] }
function apiError(error, fallback) { return error.response?.data?.error?.message || fallback }
function titleize(value) { return String(value || '').toLowerCase().replaceAll('_', ' ').replaceAll('-', ' ').replace(/\b\w/g, char => char.toUpperCase()) }
</script>

<style scoped>
.wizard-navbar{position:sticky!important;top:0!important;z-index:1000!important;display:flex!important;align-items:center!important;justify-content:space-between!important;gap:18px!important;min-height:70px!important;padding:12px 28px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.wizard-brand{display:flex!important;align-items:center!important;gap:10px!important;color:#34395e!important;text-decoration:none!important;font-weight:700!important}.wizard-brand img{width:52px!important;height:42px!important;object-fit:contain!important}.wizard-navbar-actions{display:flex!important;align-items:center!important;gap:8px!important}.wizard-icon-button,.wizard-profile{display:inline-flex!important;align-items:center!important;justify-content:center!important;min-height:38px!important;border:0!important;border-radius:30px!important;background:#f4f6f9!important;color:#34395e!important}.wizard-icon-button{width:38px!important}.wizard-profile{gap:8px!important;padding:4px 12px!important}.wizard-profile span{display:inline-flex!important;align-items:center!important;justify-content:center!important;width:30px!important;height:30px!important;border-radius:50%!important;background:#6777ef!important;color:#fff!important;font-size:11px!important;font-weight:700!important}.wizard-profile b{font-size:12px!important}.wizard-workspace{width:min(1180px,calc(100% - 32px))!important;margin:0 auto!important;padding:28px 0 40px!important}.wizard-header{display:flex!important;align-items:flex-end!important;justify-content:space-between!important;gap:20px!important;margin-bottom:20px!important}.back-link{display:inline-flex!important;align-items:center!important;gap:5px!important;margin-bottom:12px!important;border:0!important;background:transparent!important;color:#6777ef!important;font-size:12px!important;font-weight:700!important}.wizard-header p{margin:0!important;color:#6777ef!important;font-size:11px!important;font-weight:800!important;text-transform:uppercase!important}.wizard-header h1{margin:4px 0!important;color:#34395e!important;font-size:28px!important;font-weight:700!important}.wizard-header span{color:#6c757d!important;font-size:13px!important}.fee-panel{display:grid!important;gap:3px!important;min-width:180px!important;padding:16px 18px!important;border-left:3px solid #ffa426!important;border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.fee-panel small{color:#98a6ad!important;font-size:11px!important;font-weight:800!important;text-transform:uppercase!important}.fee-panel strong{color:#34395e!important;font-size:20px!important}.wizard-success,.wizard-error{margin:0 0 15px!important;padding:12px 14px!important;border-radius:3px!important;background:#e8f7f0!important;color:#47c363!important;box-shadow:0 4px 25px rgba(0,0,0,.05)!important;font-size:12px!important}.wizard-error{background:#fdeaea!important;color:#fc544b!important}.wizard-loading,.empty-step{padding:56px!important;border:1px dashed #e4e6fc!important;border-radius:3px!important;background:#fdfdff!important;text-align:center!important;color:#98a6ad!important}.wizard-layout{display:grid!important;grid-template-columns:300px minmax(0,1fr)!important;gap:20px!important;align-items:start!important}.wizard-steps{position:sticky!important;top:92px!important;display:grid!important;gap:8px!important}.wizard-steps button{display:grid!important;grid-template-columns:34px minmax(0,1fr)!important;gap:11px!important;align-items:center!important;width:100%!important;min-height:64px!important;padding:12px!important;border:0!important;border-radius:3px!important;background:#fff!important;text-align:left!important;box-shadow:0 4px 25px rgba(0,0,0,.07)!important}.wizard-steps button>span{display:inline-flex!important;align-items:center!important;justify-content:center!important;width:34px!important;height:34px!important;border-radius:50%!important;background:#eef2ff!important;color:#6777ef!important;font-size:12px!important;font-weight:800!important}.wizard-steps button.active{box-shadow:0 4px 25px rgba(103,119,239,.22)!important}.wizard-steps button.active>span,.wizard-steps button.complete>span{background:#6777ef!important;color:#fff!important}.wizard-steps strong,.wizard-steps small{display:block!important;min-width:0!important;overflow:hidden!important;text-overflow:ellipsis!important;white-space:nowrap!important}.wizard-steps strong{color:#34395e!important;font-size:13px!important}.wizard-steps small{color:#98a6ad!important;font-size:11px!important}.wizard-form{display:grid!important;gap:14px!important}.wizard-step-panel{display:grid!important;gap:16px!important;padding:24px!important;border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.1)!important}.step-heading{padding-bottom:14px!important;border-bottom:1px solid #f4f6f9!important}.step-heading small{color:#6777ef!important;font-size:10px!important;font-weight:800!important;text-transform:uppercase!important}.step-heading h2{margin:4px 0!important;color:#34395e!important;font-size:22px!important}.step-heading p{margin:4px 0 0!important;color:#6c757d!important;font-size:13px!important;line-height:1.5!important}.wizard-field{display:grid!important;gap:7px!important}.wizard-field>label,.payment-section label{color:#34395e!important;font-size:12px!important;font-weight:600!important}.wizard-field>label b{color:#fc544b!important}.wizard-field>p{margin:0!important;color:#98a6ad!important;font-size:11px!important}.wizard-field input,.wizard-field textarea,.wizard-field select,.payment-section input{width:100%!important;box-sizing:border-box!important;padding:11px 14px!important;border:1px solid #e4e6fc!important;border-radius:3px!important;background:#fdfdff!important;color:#495057!important;font:inherit!important;outline:none!important}.wizard-field input:focus,.wizard-field textarea:focus,.wizard-field select:focus,.payment-section input:focus{border-color:#6777ef!important;box-shadow:0 2px 6px #acb5f6!important}.choice-list{display:grid!important;gap:8px!important}.choice-list label{display:flex!important;align-items:center!important;gap:8px!important;color:#6c757d!important;font-size:12px!important}.choice-list input{width:auto!important}.file-input{display:flex!important;align-items:center!important;gap:10px!important;flex-wrap:wrap!important}.file-input a{color:#6777ef!important;font-size:12px!important;font-weight:700!important}.payment-section{display:grid!important;gap:11px!important;margin-top:4px!important;padding:16px!important;border:1px solid #ffe2ad!important;border-radius:3px!important;background:#fffaf0!important}.payment-section h3{margin:0!important;color:#34395e!important;font-size:16px!important}.payment-section p{margin:0!important;color:#6c757d!important;font-size:12px!important}.payment-section label{display:grid!important;gap:6px!important}.payment-proof-choice{display:grid!important;grid-template-columns:minmax(0,1fr) auto auto!important;align-items:end!important;gap:8px!important;min-width:0!important}.payment-proof-choice a,.payment-proof-choice span{align-self:center!important;min-width:0!important;overflow-wrap:anywhere!important;color:#6777ef!important;font-size:12px!important;font-weight:700!important}.proof-clear{align-self:center!important;border:0!important;background:transparent!important;color:#fc544b!important;font-size:12px!important;font-weight:700!important}.payment-section input:disabled{background:#eef1f7!important;color:#98a6ad!important}.wizard-actions{display:flex!important;justify-content:space-between!important;gap:12px!important;padding:18px!important;border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.wizard-actions>div{display:flex!important;gap:8px!important;flex-wrap:wrap!important}.primary-command,.secondary-command{display:inline-flex!important;align-items:center!important;justify-content:center!important;gap:7px!important;min-height:40px!important;padding:0 16px!important;border:0!important;border-radius:30px!important;background:#6777ef!important;color:#fff!important;font-size:12px!important;font-weight:600!important;box-shadow:0 2px 6px #acb5f6!important}.secondary-command{background:#f4f6f9!important;color:#34395e!important;box-shadow:none!important}.application-wizard-page button:disabled{cursor:not-allowed!important;opacity:.55!important}:global(.dark) .application-wizard-page{background:#0f172a!important;color:#e5e7eb!important}:global(.dark) .wizard-navbar,:global(.dark) .fee-panel,:global(.dark) .wizard-steps button,:global(.dark) .wizard-step-panel,:global(.dark) .wizard-actions{background:#1f2937!important;color:#e5e7eb!important;border-color:#334155!important}:global(.dark) .wizard-brand,:global(.dark) .wizard-header h1,:global(.dark) .fee-panel strong,:global(.dark) .wizard-steps strong,:global(.dark) .step-heading h2,:global(.dark) .wizard-field>label,:global(.dark) .payment-section h3,:global(.dark) .payment-section label{color:#f8fafc!important}:global(.dark) .wizard-header span,:global(.dark) .step-heading p,:global(.dark) .choice-list label{color:#cbd5e1!important}:global(.dark) .wizard-field input,:global(.dark) .wizard-field textarea,:global(.dark) .wizard-field select,:global(.dark) .payment-section input,:global(.dark) .wizard-icon-button,:global(.dark) .wizard-profile,:global(.dark) .secondary-command{background:#111827!important;color:#f8fafc!important;border-color:#475569!important}:global(.dark) .step-heading{border-color:#334155!important}:global(.dark) .payment-section{background:#111827!important;border-color:#92400e!important}@media(max-width:880px){.wizard-navbar{padding:12px 16px!important}.wizard-workspace{width:calc(100% - 24px)!important;padding-top:20px!important}.wizard-header{align-items:flex-start!important;flex-direction:column!important}.fee-panel{width:100%!important}.wizard-layout{grid-template-columns:1fr!important}.wizard-steps{position:static!important;grid-template-columns:repeat(2,minmax(0,1fr))!important}.wizard-actions{align-items:stretch!important;flex-direction:column!important}.wizard-actions>div{justify-content:space-between!important}.wizard-actions button{flex:1 1 150px!important}.payment-proof-choice{grid-template-columns:1fr!important;align-items:start!important}.proof-clear{justify-self:start!important}}@media(max-width:600px){.wizard-brand span,.wizard-profile b{display:none!important}.wizard-navbar{min-height:62px!important}.wizard-header h1{font-size:22px!important}.wizard-steps{grid-template-columns:1fr!important}.wizard-step-panel{padding:18px!important}.wizard-actions>div{flex-direction:column!important}.wizard-actions button{width:100%!important}}
.dropify-zone{position:relative!important;display:grid!important;place-items:center!important;gap:6px!important;min-height:142px!important;padding:18px!important;border:1px dashed #c9d0ff!important;border-radius:6px!important;background:#fdfdff!important;text-align:center!important;color:#6c757d!important;overflow:hidden!important}.dropify-zone.filled{border-style:solid!important;background:#f8fbff!important}.dropify-zone.uploading{opacity:.72!important}.dropify-zone.disabled{background:#eef1f7!important;color:#98a6ad!important}.dropify-input{position:absolute!important;inset:0!important;width:100%!important;height:100%!important;padding:0!important;border:0!important;opacity:0!important;cursor:pointer!important}.dropify-input:disabled{cursor:not-allowed!important}.dropify-icon{display:grid!important;place-items:center!important;width:42px!important;height:42px!important;border-radius:50%!important;background:#e8edff!important;color:#6777ef!important;font-size:20px!important}.dropify-zone strong,.dropify-zone small,.dropify-zone a{position:relative!important;z-index:1!important;min-width:0!important;max-width:100%!important;overflow-wrap:anywhere!important}.dropify-zone strong{color:#34395e!important;font-size:13px!important}.dropify-zone small{color:#98a6ad!important;font-size:11px!important}.dropify-zone a{color:#6777ef!important;font-size:12px!important;font-weight:800!important}.dropify-zone .proof-clear{position:relative!important;z-index:2!important;justify-self:center!important}.payment-proof-drop{min-height:150px!important}:global(.dark) .dropify-zone{background:#111827!important;border-color:#475569!important;color:#cbd5e1!important}:global(.dark) .dropify-zone strong{color:#f8fafc!important}:global(.dark) .dropify-icon{background:#1f2937!important}
</style>
