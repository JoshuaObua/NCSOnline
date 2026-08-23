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
                class="dropify-wrapper"
                :class="{
                  'has-preview': !!formAnswers[field.field_key],
                  'is-uploading': uploadingField === field.field_key
                }"
              >
                <div
                  class="dropify-drop-area"
                  @click="triggerFieldFileInput(field)"
                  @dragover.prevent
                  @drop.prevent="dropFieldFile(field, $event)"
                >
                  <input
                    :id="fieldId(field)"
                    :ref="el => setFieldInputRef(field.field_key, el)"
                    class="dropify-file-input"
                    type="file"
                    :accept="fieldAccept(field)"
                    @change="uploadFieldFile(field, $event)"
                  />

                  <!-- Uploading State Overlay -->
                  <div v-if="uploadingField === field.field_key" class="dropify-loader">
                    <div class="spinner-border text-primary spinner-border-sm" role="status"></div>
                    <span class="ms-2 fw-semibold">Uploading document...</span>
                  </div>

                  <!-- Filled / Resumed State Preview -->
                  <div v-else-if="formAnswers[field.field_key]" class="dropify-preview-container">
                    <div class="dropify-preview-card">
                      <div class="preview-icon-wrapper">
                        <img
                          v-if="isImageFile(formAnswers[field.field_key]) || field.field_type === 'image'"
                          :src="mediaUrl(formAnswers[field.field_key])"
                          alt="Uploaded Attachment"
                          class="preview-img-thumbnail"
                        />
                        <span v-else class="preview-doc-icon">
                          <i class="icofont-file-pdf"></i>
                        </span>
                      </div>

                      <div class="preview-meta-info">
                        <span class="preview-status-tag"><i class="icofont-check-circled"></i> Attached Document</span>
                        <h4 class="preview-file-name" :title="getFileName(formAnswers[field.field_key])">
                          {{ getFileName(formAnswers[field.field_key]) }}
                        </h4>
                        <div class="preview-actions-bar">
                          <a
                            :href="mediaUrl(formAnswers[field.field_key])"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="btn-preview-link"
                            @click.stop
                          >
                            <i class="icofont-eye-alt"></i> View File
                          </a>
                          <button
                            type="button"
                            class="btn-dropify-replace"
                            @click.stop="triggerFieldFileInput(field)"
                          >
                            <i class="icofont-refresh"></i> Replace
                          </button>
                          <button
                            type="button"
                            class="btn-dropify-remove"
                            @click.stop="clearFieldFile(field)"
                          >
                            <i class="icofont-trash"></i> Remove
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Default Empty State -->
                  <div v-else class="dropify-message">
                    <span class="dropify-cloud-icon">
                      <i :class="field.field_type === 'image' ? 'icofont-image' : 'icofont-cloud-upload'"></i>
                    </span>
                    <p class="dropify-primary-text">
                      Drag and drop your file here or <span>browse</span>
                    </p>
                    <p class="dropify-hint-text">
                      {{ field.field_type === 'image' ? 'Supports JPG, PNG, WEBP, GIF (Max 10MB)' : 'Supports PDF, JPG, PNG, DOCX (Max 25MB)' }}
                    </p>
                  </div>
                </div>
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
              <div class="payment-header-badge">
                <div class="badge-icon"><i class="icofont-credit-card"></i></div>
                <div>
                  <h3>Application Fee Payment</h3>
                  <p>A fee of <strong>UGX {{ formatMoney(form.price_ugx) }}</strong> is required for this application.</p>
                </div>
              </div>

              <!-- Payment Method Selector (if both or multiple allowed) -->
              <div v-if="allowsMoMo && allowsCounter" class="method-selector">
                <label class="method-card" :class="{ active: selectedPaymentMethod === 'MOBILE_MONEY' }">
                  <input v-model="selectedPaymentMethod" type="radio" value="MOBILE_MONEY" />
                  <div class="method-details">
                    <div class="method-title"><i class="icofont-smart-phone"></i> Instant Mobile Money</div>
                    <small>MTN & Airtel Money with automated instant approval.</small>
                  </div>
                  <span class="instant-tag">Instant</span>
                </label>
                <label class="method-card" :class="{ active: selectedPaymentMethod === 'OVER_THE_COUNTER' }">
                  <input v-model="selectedPaymentMethod" type="radio" value="OVER_THE_COUNTER" />
                  <div class="method-details">
                    <div class="method-title"><i class="icofont-bank-alt"></i> Bank Deposit / Over the Counter</div>
                    <small>Upload receipt or enter URA PRN for manual staff check.</small>
                  </div>
                </label>
              </div>

              <!-- Mobile Money Flow -->
              <div v-if="selectedPaymentMethod === 'MOBILE_MONEY'" class="momo-flow-card">
                <div class="momo-amount-display">
                  <span>Amount to Pay</span>
                  <strong>UGX {{ formatMoney(form.price_ugx) }}</strong>
                </div>

                <div class="momo-input-group">
                  <label for="momoPhone">Mobile Money Phone Number (Uganda)
                    <div class="phone-input-wrapper">
                      <span class="phone-country-code">+256</span>
                      <input
                        id="momoPhone"
                        v-model="momoPhoneNumber"
                        type="tel"
                        maxlength="15"
                        placeholder="770 000000"
                        :disabled="momoPolling || momoSuccess"
                        @input="onMomoPhoneInput"
                      />
                      <span v-if="detectedCarrier" class="carrier-badge" :class="detectedCarrier.toLowerCase()">
                        {{ detectedCarrier }}
                      </span>
                    </div>
                  </label>

                  <!-- Action / Status Controls -->
                  <div v-if="!momoPolling && !momoSuccess" class="momo-action-row">
                    <button
                      type="button"
                      class="btn-momo-pay"
                      :disabled="!isValidMomoPhone || formSaving"
                      @click="initiateMoMoPayment"
                    >
                      <i class="icofont-paper-plane"></i> Pay UGX {{ formatMoney(form.price_ugx) }} via Mobile Money
                    </button>
                  </div>

                  <!-- Polling / USSD Prompt Tracker Banner -->
                  <div v-if="momoPolling" class="momo-status-tracker">
                    <div class="pulse-spinner"></div>
                    <div class="tracker-info">
                      <h4>USSD Prompt Sent!</h4>
                      <p>A payment request has been sent to <strong>{{ momoPhoneNumber }}</strong>. Please check your phone and enter your Mobile Money PIN to complete payment.</p>
                      <span class="tracker-timer">Waiting for confirmation ({{ momoCountdown }}s remaining)...</span>
                    </div>
                  </div>

                  <!-- Success Banner -->
                  <div v-if="momoSuccess" class="momo-success-banner">
                    <i class="icofont-check-circled"></i>
                    <div>
                      <h4>Payment Confirmed!</h4>
                      <p>Your Mobile Money payment of <strong>UGX {{ formatMoney(form.price_ugx) }}</strong> was received successfully. (Ref: <code>{{ momoTxRef }}</code>)</p>
                    </div>
                  </div>

                  <!-- Error Banner -->
                  <div v-if="momoError" class="momo-error-banner">
                    <i class="icofont-warning-alt"></i>
                    <div>
                      <h4>Payment Failed</h4>
                      <p>{{ momoError }}</p>
                      <button type="button" class="btn-retry-momo" @click="retryMoMo">Try Again</button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Over the Counter Flow -->
              <div v-else-if="selectedPaymentMethod === 'OVER_THE_COUNTER'" class="counter-flow-card">
                <div class="bank-deposit-guide mb-3">
                  <div class="guide-icon"><i class="icofont-bank-alt"></i></div>
                  <div>
                    <h4 class="mb-1 fw-bold">Bank Deposit / Over the Counter Payment</h4>
                    <p class="mb-0 text-muted text-sm">
                      Please deposit the application fee of <strong>UGX {{ formatMoney(form.price_ugx) }}</strong> at any authorized bank or NCS cashier counter, then upload the stamped deposit receipt below.
                    </p>
                  </div>
                </div>

                <div class="form-group mb-2">
                  <label class="form-label fw-bold">Upload Stamped Bank Deposit Receipt / Proof <b class="text-danger">*</b></label>
                  
                  <div
                    class="dropify-wrapper payment-proof-dropify"
                    :class="{
                      'has-preview': !!paymentProofURL,
                      'is-uploading': uploadingPaymentProof
                    }"
                  >
                    <div
                      class="dropify-drop-area"
                      @click="triggerPaymentProofInput"
                      @dragover.prevent
                      @drop.prevent="dropPaymentProofFile"
                    >
                      <input
                        ref="paymentProofInput"
                        class="dropify-file-input"
                        type="file"
                        accept=".pdf,.png,.jpg,.jpeg,application/pdf,image/png,image/jpeg"
                        @change="uploadPaymentProofFile"
                      />

                      <!-- Uploading State Overlay -->
                      <div v-if="uploadingPaymentProof" class="dropify-loader">
                        <div class="spinner-border text-primary spinner-border-sm" role="status"></div>
                        <span class="ms-2 fw-semibold">Uploading payment proof...</span>
                      </div>

                      <!-- Filled / Resumed State Preview -->
                      <div v-else-if="paymentProofURL" class="dropify-preview-container">
                        <div class="dropify-preview-card">
                          <div class="preview-icon-wrapper">
                            <img
                              v-if="isImageFile(paymentProofURL)"
                              :src="mediaUrl(paymentProofURL)"
                              alt="Payment Proof"
                              class="preview-img-thumbnail"
                            />
                            <span v-else class="preview-doc-icon">
                              <i class="icofont-file-pdf"></i>
                            </span>
                          </div>

                          <div class="preview-meta-info">
                            <span class="preview-status-tag"><i class="icofont-check-circled"></i> Payment Receipt Attached</span>
                            <h4 class="preview-file-name" :title="paymentProofFileName || getFileName(paymentProofURL)">
                              {{ paymentProofFileName || getFileName(paymentProofURL) }}
                            </h4>
                            <div class="preview-actions-bar">
                              <a
                                :href="mediaUrl(paymentProofURL)"
                                target="_blank"
                                rel="noopener noreferrer"
                                class="btn-preview-link"
                                @click.stop
                              >
                                <i class="icofont-eye-alt"></i> View Receipt
                              </a>
                              <button
                                type="button"
                                class="btn-dropify-replace"
                                @click.stop="triggerPaymentProofInput"
                              >
                                <i class="icofont-refresh"></i> Replace Receipt
                              </button>
                              <button
                                type="button"
                                class="btn-dropify-remove"
                                @click.stop="clearPaymentProofFile"
                              >
                                <i class="icofont-trash"></i> Remove
                              </button>
                            </div>
                          </div>
                        </div>
                      </div>

                      <!-- Default Empty State -->
                      <div v-else class="dropify-message">
                        <span class="dropify-cloud-icon">
                          <i class="icofont-paperclip"></i>
                        </span>
                        <p class="dropify-primary-text">
                          Drag and drop bank deposit slip here or <span>browse</span>
                        </p>
                        <p class="dropify-hint-text">
                          Supports PDF, PNG, JPG or JPEG (Max 15MB)
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
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
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCurrentUser } from '@/api/auth.js'
import { mediaUrl } from '@/api/client.js'
import * as cms from '@/api/cms.js'
import {
  portalGetForm,
  portalGetPaymentStatus,
  portalInitiateMoMoPayment,
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

const selectedPaymentMethod = ref('MOBILE_MONEY')
const momoPhoneNumber = ref('')
const momoPolling = ref(false)
const momoSuccess = ref(false)
const momoError = ref('')
const momoTxRef = ref('')
const momoCountdown = ref(60)
let pollTimer = null
let countdownTimer = null

ensureOtikaStyles()

const steps = computed(() => buildSectionSteps(form))
const currentStep = computed(() => steps.value[currentStepIndex.value] || { title: 'Application Details', fields: [] })
const isLastStep = computed(() => currentStepIndex.value >= steps.value.length - 1)
const firstName = computed(() => profile.first_name || String(profile.email || 'User').split('@')[0])
const initials = computed(() => `${profile.first_name?.[0] || ''}${profile.last_name?.[0] || ''}`.toUpperCase() || 'U')
const paymentRequired = computed(() => Number(form.price_ugx) > 0 && !['PAID', 'PROOF_UPLOADED'].includes(activeSubmission.value?.payment_status))

const allowedMethods = computed(() => {
  let raw = form.allowed_payment_methods
  if (typeof raw === 'string') {
    try { raw = JSON.parse(raw) } catch { raw = [] }
  }
  if (!Array.isArray(raw) || !raw.length) {
    return ['OVER_THE_COUNTER', 'MOBILE_MONEY']
  }
  return raw
})
const allowsMoMo = computed(() => allowedMethods.value.includes('MOBILE_MONEY'))
const allowsCounter = computed(() => allowedMethods.value.includes('OVER_THE_COUNTER'))

const detectedCarrier = computed(() => {
  const num = momoPhoneNumber.value.replace(/\D/g, '')
  if (num.startsWith('0111')) return 'Sandbox Test'
  if (num.startsWith('77') || num.startsWith('78') || num.startsWith('76') || num.startsWith('077') || num.startsWith('078') || num.startsWith('076')) return 'MTN MoMo'
  if (num.startsWith('70') || num.startsWith('75') || num.startsWith('74') || num.startsWith('070') || num.startsWith('075') || num.startsWith('074')) return 'Airtel Money'
  return ''
})

const isValidMomoPhone = computed(() => {
  const num = momoPhoneNumber.value.replace(/\D/g, '')
  return num.length >= 9
})

onMounted(loadWizard)
onUnmounted(stopTimers)

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

    if (profile.phone) {
      momoPhoneNumber.value = profile.phone
    }

    if (allowsMoMo.value) {
      selectedPaymentMethod.value = 'MOBILE_MONEY'
    } else if (allowsCounter.value) {
      selectedPaymentMethod.value = 'OVER_THE_COUNTER'
    }

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
    if (draft?.payment_status === 'PAID') {
      momoSuccess.value = true
    }
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

const fieldInputs = reactive({})
function setFieldInputRef(key, el) { if (el) fieldInputs[key] = el }
function triggerFieldFileInput(field) { if (fieldInputs[field.field_key]) fieldInputs[field.field_key].click() }
function triggerPaymentProofInput() { if (paymentProofInput.value) paymentProofInput.value.click() }

function getFileName(url) {
  if (!url) return ''
  const clean = String(url).split('?')[0]
  const parts = clean.split('/')
  const raw = parts[parts.length - 1] || 'Document'
  try { return decodeURIComponent(raw) } catch { return raw }
}

function isImageFile(url) {
  if (!url) return false
  return /\.(png|jpe?g|webp|gif|svg)$/i.test(String(url).split('?')[0])
}

function onMomoPhoneInput() {
  momoError.value = ''
}

function retryMoMo() {
  momoError.value = ''
  momoPolling.value = false
  momoSuccess.value = false
  stopTimers()
}

function stopTimers() {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
  if (countdownTimer) { clearInterval(countdownTimer); countdownTimer = null }
}

async function initiateMoMoPayment() {
  if (!isValidMomoPhone.value) return
  formError.value = ''
  momoError.value = ''
  formSaving.value = true
  try {
    if (!activeSubmission.value) {
      activeSubmission.value = await portalSaveDraft(form.id, { ...formAnswers })
    }
    const tx = await portalInitiateMoMoPayment(activeSubmission.value.id, momoPhoneNumber.value)
    momoTxRef.value = tx.transaction_reference || ''

    if (tx.status === 'SUCCESS') {
      momoSuccess.value = true
      momoPolling.value = false
      activeSubmission.value.payment_status = 'PAID'
      success.value = 'Mobile money payment confirmed!'
      return
    }
    if (tx.status === 'FAILED') {
      momoError.value = tx.status_message || 'Mobile money collection failed.'
      momoPolling.value = false
      return
    }

    // Start polling & countdown
    momoPolling.value = true
    momoCountdown.value = 60
    stopTimers()

    countdownTimer = setInterval(() => {
      if (momoCountdown.value > 0) {
        momoCountdown.value -= 1
      } else {
        stopTimers()
        momoPolling.value = false
        momoError.value = 'Payment timed out waiting for phone PIN approval. You can try again.'
      }
    }, 1000)

    pollTimer = setInterval(async () => {
      try {
        const check = await portalGetPaymentStatus(activeSubmission.value.id)
        if (check.status === 'SUCCESS') {
          stopTimers()
          momoPolling.value = false
          momoSuccess.value = true
          activeSubmission.value.payment_status = 'PAID'
          success.value = 'Payment confirmed! You may now submit your application.'
        } else if (check.status === 'FAILED') {
          stopTimers()
          momoPolling.value = false
          momoError.value = check.status_message || 'Payment was declined or failed on your mobile phone.'
        }
      } catch (err) {
        console.warn('Poll status error', err)
      }
    }, 3000)

  } catch (error) {
    momoError.value = apiError(error, 'Could not initiate mobile money payment.')
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
      if (selectedPaymentMethod.value === 'MOBILE_MONEY') {
        throw new Error('Please complete your Mobile Money payment before submitting.')
      }
      if (!validPaymentProof()) throw new Error('Please upload your bank deposit receipt / payment proof before submitting.')
      await portalUploadPaymentProof(submission.id, {
        payment_reference: paymentProofFileName.value || getFileName(paymentProofURL.value) || 'BANK_DEPOSIT',
        payment_proof_url: paymentProofURL.value.trim(),
        payment_amount_ugx: Number(form.price_ugx),
        payment_method: 'OVER_THE_COUNTER',
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
  if (uploadingPaymentProof.value) return
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
  if (Number(form.price_ugx) > 0 && paymentRequired.value) {
    if (selectedPaymentMethod.value === 'MOBILE_MONEY' && !momoSuccess.value) {
      currentStepIndex.value = steps.value.length - 1
      formError.value = 'Please complete the Mobile Money payment on your phone before submitting.'
      return false
    }
    if (selectedPaymentMethod.value === 'OVER_THE_COUNTER' && !validPaymentProof()) {
      currentStepIndex.value = steps.value.length - 1
      formError.value = 'Please upload your bank deposit receipt / payment proof before submitting.'
      return false
    }
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
.payment-header-badge{display:flex;align-items:center;gap:14px;padding:16px;background:#f0f4ff;border:1px solid #d9e2fc;border-radius:6px;margin-bottom:14px}
.payment-header-badge .badge-icon{display:grid;place-items:center;width:44px;height:44px;border-radius:50%;background:#6366f1;color:#fff;font-size:22px;flex-shrink:0}
.payment-header-badge h3{margin:0;color:#1e293b;font-size:16px;font-weight:700}
.payment-header-badge p{margin:3px 0 0;color:#64748b;font-size:13px}
.method-selector{display:grid;grid-template-columns:repeat(auto-fit, minmax(240px, 1fr));gap:12px;margin-bottom:16px}
.method-card{position:relative;display:flex;align-items:flex-start;gap:12px;padding:14px;background:#fff;border:2px solid #e2e8f0;border-radius:6px;cursor:pointer;transition:all 0.2s ease}
.method-card:hover{border-color:#a5b4fc;background:#faf5ff}
.method-card.active{border-color:#6366f1;background:#f5f3ff;box-shadow:0 2px 8px rgba(99,102,241,0.12)}
.method-card input[type="radio"]{margin-top:3px;accent-color:#6366f1;cursor:pointer}
.method-details{display:flex;flex-direction:column;gap:3px}
.method-title{font-size:13.5px;font-weight:700;color:#1e293b;display:flex;align-items:center;gap:6px}
.method-details small{color:#64748b;font-size:11.5px;line-height:1.4}
.instant-tag{position:absolute;top:10px;right:10px;padding:2px 7px;border-radius:12px;background:#e0e7ff;color:#4338ca;font-size:10px;font-weight:800;text-transform:uppercase}

.momo-flow-card{background:#fff;border:1px solid #e2e8f0;border-radius:6px;padding:20px;display:grid;gap:16px}
.momo-amount-display{display:flex;align-items:center;justify-content:space-between;padding:12px 16px;background:#f8fafc;border-radius:6px;border-left:4px solid #6366f1}
.momo-amount-display span{font-size:12px;text-transform:uppercase;color:#64748b;font-weight:700}
.momo-amount-display strong{font-size:22px;color:#1e293b;font-weight:800}

.momo-input-group{display:grid;gap:12px}
.phone-input-wrapper{display:flex;align-items:center;position:relative;border:1px solid #cbd5e1;border-radius:4px;background:#fff;overflow:hidden;margin-top:6px}
.phone-input-wrapper:focus-within{border-color:#6366f1;box-shadow:0 0 0 3px rgba(99,102,241,0.15)}
.phone-country-code{padding:10px 12px;background:#f1f5f9;color:#475569;font-weight:700;font-size:13px;border-right:1px solid #cbd5e1;user-select:none}
.phone-input-wrapper input{border:0!important;padding:10px 14px!important;box-shadow:none!important;font-size:14px;font-weight:600;color:#1e293b;flex:1}
.carrier-badge{margin-right:10px;padding:4px 8px;border-radius:4px;font-size:11px;font-weight:700;text-transform:uppercase}
.carrier-badge.mtn{background:#fef08a;color:#854d0e}
.carrier-badge.airtel{background:#fee2e2;color:#991b1b}
.carrier-badge.sandbox{background:#e0e7ff;color:#3730a3}



.btn-momo-pay{display:inline-flex;align-items:center;justify-content:center;gap:8px;padding:12px 24px;border:0;border-radius:30px;background:linear-gradient(135deg, #4f46e5, #6366f1);color:#fff;font-size:13px;font-weight:700;cursor:pointer;box-shadow:0 4px 12px rgba(79,70,229,0.3);transition:all 0.2s ease}
.btn-momo-pay:hover:not(:disabled){transform:translateY(-1px);box-shadow:0 6px 16px rgba(79,70,229,0.4)}
.btn-momo-pay:disabled{opacity:0.5;cursor:not-allowed}

.momo-status-tracker{display:flex;align-items:flex-start;gap:14px;padding:16px;background:#eff6ff;border:1px solid #bfdbfe;border-radius:6px}
.pulse-spinner{width:24px;height:24px;border:3px solid #93c5fd;border-top-color:#2563eb;border-radius:50%;animation:spin 1s linear infinite;flex-shrink:0;margin-top:2px}
@keyframes spin{to{transform:rotate(360deg)}}
.tracker-info h4{margin:0;color:#1e40af;font-size:14px;font-weight:700}
.tracker-info p{margin:4px 0 6px;color:#1e3a8a;font-size:12px;line-height:1.4}
.tracker-timer{font-size:11px;font-weight:700;color:#3b82f6}

.momo-success-banner{display:flex;align-items:center;gap:12px;padding:14px 16px;background:#ecfdf5;border:1px solid #a7f3d0;border-radius:6px;color:#065f46}
.momo-success-banner i{font-size:26px;color:#059669;flex-shrink:0}
.momo-success-banner h4{margin:0;font-size:14px;font-weight:700}
.momo-success-banner p{margin:2px 0 0;font-size:12px}
.momo-success-banner code{background:#d1fae5;padding:2px 6px;border-radius:3px;color:#064e3b;font-weight:700}

.momo-error-banner{display:flex;align-items:flex-start;gap:12px;padding:14px 16px;background:#fef2f2;border:1px solid #fecaca;border-radius:6px;color:#991b1b}
.momo-error-banner i{font-size:24px;color:#dc2626;flex-shrink:0;margin-top:2px}
.momo-error-banner h4{margin:0;font-size:14px;font-weight:700}
.momo-error-banner p{margin:2px 0 8px;font-size:12px}
.btn-retry-momo{border:1px solid #ef4444;background:#fff;color:#dc2626;padding:4px 12px;border-radius:4px;font-size:11.5px;font-weight:700;cursor:pointer}
.btn-retry-momo:hover{background:#fee2e2}

:global(.dark) .payment-header-badge{background:#1e293b;border-color:#334155}
:global(.dark) .payment-header-badge h3{color:#f8fafc}
:global(.dark) .payment-header-badge p{color:#cbd5e1}
:global(.dark) .method-card{background:#0f172a;border-color:#334155}
:global(.dark) .method-card.active{background:#1e1b4b;border-color:#6366f1}
:global(.dark) .method-title{color:#f8fafc}
:global(.dark) .momo-flow-card{background:#0f172a;border-color:#334155}
:global(.dark) .momo-amount-display{background:#1e293b}
:global(.dark) .momo-amount-display strong{color:#f8fafc}
:global(.dark) .phone-input-wrapper{background:#1e293b;border-color:#475569}
:global(.dark) .phone-country-code{background:#334155;color:#f8fafc;border-color:#475569}
:global(.dark) .phone-input-wrapper input{background:#1e293b!important;color:#f8fafc!important}
:global(.dark) .momo-status-tracker{background:#172554;border-color:#1e40af}
:global(.dark) .tracker-info h4{color:#93c5fd}
:global(.dark) .tracker-info p{color:#bfdbfe}

/* ========================================================
   DROPIFY COMPONENT STYLES
   ======================================================== */
.dropify-wrapper{position:relative;width:100%}
.dropify-drop-area{position:relative;display:flex;align-items:center;justify-content:center;min-height:140px;padding:20px;background:#f8fafc;border:2px dashed #cbd5e1;border-radius:8px;cursor:pointer;transition:all .2s ease;text-align:center}
.dropify-drop-area:hover{border-color:#6366f1;background:#f5f3ff}
.dropify-wrapper.has-preview .dropify-drop-area{padding:10px;border-style:solid;border-color:#e2e8f0;background:#fafafa;cursor:default}
.dropify-file-input{display:none}

.dropify-loader{display:flex;align-items:center;justify-content:center;gap:10px;color:#4f46e5;font-size:13px}

.dropify-message{display:flex;flex-direction:column;align-items:center;gap:6px}
.dropify-cloud-icon{font-size:36px;color:#6366f1;line-height:1}
.dropify-primary-text{margin:0;font-size:13px;font-weight:600;color:#334155}
.dropify-primary-text span{color:#4f46e5;text-decoration:underline}
.dropify-hint-text{margin:0;font-size:11px;color:#94a3b8}

.dropify-preview-container{width:100%}
.dropify-preview-card{display:flex;align-items:center;gap:14px;background:#fff;border:1px solid #e2e8f0;border-radius:6px;padding:12px 16px;box-shadow:0 1px 4px rgba(0,0,0,0.04);width:100%;box-sizing:border-box}
.preview-icon-wrapper{width:56px;height:56px;border-radius:6px;overflow:hidden;background:#f1f5f9;display:grid;place-items:center;flex-shrink:0;border:1px solid #e2e8f0}
.preview-img-thumbnail{width:100%;height:100%;object-fit:cover}
.preview-doc-icon{font-size:30px;color:#e11d48;display:grid;place-items:center}
.preview-meta-info{flex:1;min-width:0;display:flex;flex-direction:column;gap:3px;text-align:left}
.preview-status-tag{font-size:10px;font-weight:800;text-transform:uppercase;color:#16a34a;display:inline-flex;align-items:center;gap:3px}
.preview-file-name{font-size:13px;font-weight:700;color:#1e293b;margin:0;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.preview-actions-bar{display:flex;align-items:center;gap:8px;flex-wrap:wrap;margin-top:4px}
.btn-preview-link{display:inline-flex;align-items:center;gap:4px;padding:3px 9px;background:#eff6ff;color:#2563eb;font-size:11.5px;font-weight:700;border-radius:4px;text-decoration:none}
.btn-preview-link:hover{background:#dbeafe;color:#1d4ed8}
.btn-dropify-replace{display:inline-flex;align-items:center;gap:4px;padding:3px 9px;background:#f8fafc;border:1px solid #cbd5e1;color:#475569;font-size:11.5px;font-weight:700;border-radius:4px;cursor:pointer}
.btn-dropify-replace:hover{background:#f1f5f9;color:#0f172a}
.btn-dropify-remove{display:inline-flex;align-items:center;gap:4px;padding:3px 9px;background:#fef2f2;border:1px solid #fecaca;color:#dc2626;font-size:11.5px;font-weight:700;border-radius:4px;cursor:pointer}
.btn-dropify-remove:hover{background:#fee2e2;color:#b91c1c}

/* Bank deposit guide */
.bank-deposit-guide{display:flex;align-items:flex-start;gap:12px;padding:12px 14px;background:#f8fafc;border:1px solid #e2e8f0;border-radius:6px}
.bank-deposit-guide .guide-icon{width:36px;height:36px;border-radius:50%;background:#e0e7ff;color:#4338ca;display:grid;place-items:center;font-size:18px;flex-shrink:0}
.text-sm{font-size:12px!important;line-height:1.4!important}

/* Dark mode for Dropify */
:global(.dark) .dropify-drop-area{background:#111827;border-color:#374151}
:global(.dark) .dropify-drop-area:hover{border-color:#818cf8;background:#1e1b4b}
:global(.dark) .dropify-wrapper.has-preview .dropify-drop-area{background:#111827;border-color:#374151}
:global(.dark) .dropify-primary-text{color:#e2e8f0}
:global(.dark) .dropify-primary-text span{color:#a5b4fc}
:global(.dark) .dropify-hint-text{color:#94a3b8}
:global(.dark) .dropify-preview-card{background:#1f2937;border-color:#374151}
:global(.dark) .preview-icon-wrapper{background:#111827;border-color:#374151}
:global(.dark) .preview-file-name{color:#f8fafc}
:global(.dark) .btn-dropify-replace{background:#374151;border-color:#4b5563;color:#f8fafc}
:global(.dark) .bank-deposit-guide{background:#1e293b;border-color:#374151}
:global(.dark) .bank-deposit-guide h4{color:#f8fafc}
:global(.dark) .bank-deposit-guide p{color:#cbd5e1!important}
</style>

