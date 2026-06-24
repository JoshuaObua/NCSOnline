<template>
  <div class="bg-gray-50 min-h-screen py-10">
    <div class="max-w-3xl mx-auto px-4">

      <!-- Toast -->
      <Transition name="toast">
        <div v-if="toast" class="fixed top-5 right-5 z-[100] px-4 py-3 bg-gray-900 text-white text-sm rounded-xl shadow-xl">{{ toast }}</div>
      </Transition>

      <div v-if="loading" class="bg-white rounded-2xl shadow-sm overflow-hidden">
        <div class="h-48 bg-gray-100 animate-pulse"/>
        <div class="p-6 space-y-3">
          <div class="h-6 w-2/3 bg-gray-100 rounded animate-pulse"/>
          <div class="h-3 w-full bg-gray-100 rounded animate-pulse"/>
          <div class="h-3 w-5/6 bg-gray-100 rounded animate-pulse"/>
        </div>
      </div>

      <div v-else-if="!form" class="bg-white rounded-2xl p-10 text-center shadow-sm">
        <p class="text-gray-600">This form is not available.</p>
        <router-link to="/my-portal" class="text-primary-600 hover:underline text-sm mt-3 inline-block">← Back to My Portal</router-link>
      </div>

      <div v-else class="bg-white rounded-2xl shadow-sm overflow-hidden">
        <!-- Banner -->
        <div v-if="form.banner_image_url" class="h-48 bg-cover bg-center" :style="{ backgroundImage: `url(${mediaUrl(form.banner_image_url)})` }"/>
        <div v-else class="h-32 bg-gradient-to-br from-primary-500 to-primary-700"/>

        <!-- Header -->
        <header class="px-6 md:px-10 pt-8 pb-4 border-b border-gray-100">
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0">
              <h1 class="text-2xl md:text-3xl font-bold text-gray-900">{{ form.title }}</h1>
              <p v-if="form.description" class="mt-2 text-gray-600 text-sm md:text-base">{{ form.description }}</p>
            </div>
            <div class="flex-shrink-0 text-right">
              <div class="text-xs text-gray-400 uppercase tracking-wider">Fee</div>
              <div class="text-lg font-bold" :class="form.price_ugx > 0 ? 'text-gray-900' : 'text-green-600'">
                {{ form.price_ugx > 0 ? `UGX ${formatNum(form.price_ugx)}` : 'Free' }}
              </div>
            </div>
          </div>
          <div v-if="submitted" class="mt-4 p-3 bg-green-50 border border-green-200 rounded-lg text-sm text-green-800">
            ✓ Submitted successfully. Reference: <strong>{{ submitted.submission_reference }}</strong>
          </div>
        </header>

        <!-- Fields -->
        <form v-if="!submitted" @submit.prevent="onSubmit" class="px-6 md:px-10 py-6 space-y-5">
          <div v-for="field in form.fields" :key="field.id">
            <label class="block text-sm font-medium text-gray-800 mb-1">
              {{ field.label }}
              <span v-if="field.is_required" class="text-red-500">*</span>
            </label>
            <p v-if="field.help_text" class="text-xs text-gray-500 mb-1.5">{{ field.help_text }}</p>

            <!-- short_text / email / number / phone / date -->
            <input v-if="['short_text','email','number','phone','date'].includes(field.field_type)"
              :type="inputType(field.field_type)"
              v-model="answers[field.field_key]"
              :placeholder="field.placeholder"
              :required="field.is_required"
              class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/>

            <!-- long_text -->
            <textarea v-else-if="field.field_type === 'long_text'"
              v-model="answers[field.field_key]"
              :placeholder="field.placeholder"
              :required="field.is_required"
              rows="4"
              class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"></textarea>

            <!-- file / image (dropzone) -->
            <DropzoneUpload v-else-if="['file','image'].includes(field.field_type)"
              v-model="answers[field.field_key]"
              :accept="field.field_type === 'image' ? 'image/*' : '*/*'"
              :label="field.field_type === 'image' ? 'Drop an image here or click to upload' : 'Drop a file here or click to upload'"
              :hint="field.placeholder"
              preview-class="h-40"
            />

            <!-- dropdown -->
            <select v-else-if="field.field_type === 'dropdown'"
              v-model="answers[field.field_key]"
              :required="field.is_required"
              class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white">
              <option value="">— Select —</option>
              <option v-for="opt in optionsOf(field)" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
            </select>

            <!-- radio -->
            <div v-else-if="field.field_type === 'radio'" class="space-y-1.5">
              <label v-for="opt in optionsOf(field)" :key="opt.value"
                class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer">
                <input type="radio" :name="field.field_key" :value="opt.value"
                  v-model="answers[field.field_key]" :required="field.is_required"
                  class="text-primary-600"/>
                {{ opt.label }}
              </label>
            </div>

            <!-- checkbox (multi-select) -->
            <div v-else-if="field.field_type === 'checkbox'" class="space-y-1.5">
              <label v-for="opt in optionsOf(field)" :key="opt.value"
                class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer">
                <input type="checkbox" :value="opt.value"
                  v-model="checkboxModel(field.field_key).value"
                  class="rounded border-gray-300 text-primary-600"/>
                {{ opt.label }}
              </label>
            </div>
          </div>

          <!-- Payment proof (only when fee > 0 and not yet paid) -->
          <section v-if="form.price_ugx > 0" class="mt-6 pt-6 border-t border-gray-100">
            <h3 class="text-sm font-semibold text-gray-800 mb-3">Payment Proof</h3>
            <p class="text-xs text-gray-500 mb-3">
              Deposit <strong>UGX {{ formatNum(form.price_ugx) }}</strong> and enter your reference.
              An officer verifies the payment after submission.
            </p>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <input v-model="payment.reference" type="text" placeholder="Bank reference / transaction ID"
                class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2"/>
              <input v-model.number="payment.amount" type="number" min="0" :placeholder="form.price_ugx"
                class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2"/>
            </div>
          </section>

          <!-- Actions -->
          <div class="flex flex-wrap items-center gap-3 pt-4">
            <button type="button" @click="saveDraft" :disabled="busy"
              class="px-4 py-2 text-sm font-medium border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-60">
              Save Draft
            </button>
            <button type="submit" :disabled="busy"
              class="px-5 py-2 text-sm font-medium bg-primary-600 hover:bg-primary-700 disabled:opacity-60 text-white rounded-lg">
              {{ form.price_ugx > 0 ? 'Submit & Record Payment' : 'Submit' }}
            </button>
            <span v-if="submission?.id" class="text-xs text-gray-400">
              Draft saved · ref {{ submission.id.slice(0, 8) }}
            </span>
          </div>
        </form>

        <!-- Post-submit summary -->
        <div v-else class="px-6 md:px-10 py-8">
          <router-link to="/my-portal" class="inline-flex items-center gap-1 text-primary-600 hover:underline text-sm">
            ← Back to My Portal
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  portalGetForm, portalSaveDraft, portalSubmit, portalUploadPaymentProof,
} from '@/api/forms'
import DropzoneUpload from '@/components/ui/DropzoneUpload.vue'
import { mediaUrl } from '@/api/client.js'

const route = useRoute()
const form = ref(null)
const loading = ref(true)
const busy = ref(false)
const toast = ref('')
const answers = reactive({})
const submission = ref(null)
const submitted = ref(null)
const payment = reactive({ reference: '', amount: 0 })

function showToast(msg) {
  toast.value = msg
  setTimeout(() => (toast.value = ''), 2800)
}

function formatNum(n) { return new Intl.NumberFormat().format(n || 0) }

function inputType(t) {
  return { email: 'email', number: 'number', phone: 'tel', date: 'date', short_text: 'text' }[t] || 'text'
}

function optionsOf(field) {
  if (field.config && Array.isArray(field.config.options)) return field.config.options
  return []
}

function checkboxModel(key) {
  if (!Array.isArray(answers[key])) answers[key] = []
  return computed({
    get: () => answers[key],
    set: (v) => { answers[key] = v },
  })
}

async function load() {
  loading.value = true
  try {
    form.value = await portalGetForm(route.params.slug)
    // Initialize answers map for known field keys
    for (const f of form.value.fields || []) {
      if (!(f.field_key in answers)) {
        answers[f.field_key] = f.field_type === 'checkbox' ? [] : ''
      }
    }
    payment.amount = Number(form.value.price_ugx) || 0
  } catch (e) {
    form.value = null
  } finally {
    loading.value = false
  }
}

async function saveDraft() {
  if (!form.value) return
  busy.value = true
  try {
    submission.value = await portalSaveDraft(form.value.id, JSON.parse(JSON.stringify(answers)))
    showToast('Draft saved')
  } catch (e) {
    showToast(e?.response?.data?.message || 'Save failed')
  } finally {
    busy.value = false
  }
}

async function onSubmit() {
  if (!form.value) return
  busy.value = true
  try {
    // Always save latest answers first
    submission.value = await portalSaveDraft(form.value.id, JSON.parse(JSON.stringify(answers)))
    if (form.value.price_ugx > 0) {
      if (!payment.reference || !payment.amount) {
        showToast('Enter your payment reference and amount before submitting')
        return
      }
      await portalUploadPaymentProof(submission.value.id, payment.reference, payment.amount)
    }
    submitted.value = await portalSubmit(submission.value.id)
    showToast('Submitted')
  } catch (e) {
    showToast(e?.response?.data?.message || 'Submission failed')
  } finally {
    busy.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(-8px); }
.toast-enter-active, .toast-leave-active { transition: all .25s ease; }
</style>
