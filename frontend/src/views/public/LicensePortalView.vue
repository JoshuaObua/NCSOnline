<template>
  <div>
    <!-- Hero Banner -->
    <div class="bg-cream pt-20 pb-12 px-4">
      <div class="max-w-7xl mx-auto sm:px-6 lg:px-8 text-center">
        <p class="section-tag mb-3">Online Services</p>
        <h1 class="text-4xl md:text-5xl font-bold text-darken mb-4">Apply for a <span class="text-accent">Sports Licence</span></h1>
        <p class="text-gray-500 max-w-xl mx-auto text-lg">Apply for a new sports licence or renew an existing one through our secure online portal.</p>

        <!-- Step indicator -->
        <div class="flex items-center justify-center gap-0 mt-8 max-w-md mx-auto">
          <div v-for="(step, i) in steps" :key="i" class="flex items-center flex-1">
            <div class="flex flex-col items-center gap-1 flex-1">
              <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold transition-colors"
                :class="i < currentStep ? 'bg-[#F48C06] text-white' : i === currentStep ? 'bg-[#112b4e] text-white' : 'bg-gray-200 text-gray-500'">
                <svg v-if="i < currentStep" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/></svg>
                <span v-else>{{ i + 1 }}</span>
              </div>
              <span class="text-xs text-gray-500 text-center leading-tight hidden sm:block">{{ step }}</span>
            </div>
            <div v-if="i < steps.length - 1" class="h-px flex-1 max-w-8 mb-4 transition-colors" :class="i < currentStep ? 'bg-[#F48C06]' : 'bg-gray-200'"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Wave -->
    <div class="text-white -mb-1">
      <svg viewBox="0 0 1200 80" preserveAspectRatio="none" class="w-full h-10 fill-white">
        <path d="M0,40 C300,80 900,0 1200,40 L1200,80 L0,80 Z"/>
      </svg>
    </div>

    <!-- Content -->
    <div class="bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">

        <!-- If not logged in -->
        <div v-if="!isAuthenticated" class="max-w-xl mx-auto text-center">
          <div class="bg-[#FEF9F2] border border-[#F48C06]/20 rounded-2xl p-8 mb-8 shadow-sm">
            <div class="w-16 h-16 bg-[#112b4e] rounded-2xl flex items-center justify-center mx-auto mb-5">
              <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
              </svg>
            </div>
            <h2 class="text-xl font-bold text-darken mb-3">Sign In Required</h2>
            <p class="text-gray-500 mb-6">You need to create an account or sign in to apply for or manage your sports licence.</p>
            <div class="flex flex-col sm:flex-row gap-3 justify-center">
              <router-link to="/login" class="bg-[#F48C06] hover:bg-[#d47b05] text-white font-bold px-6 py-3 rounded-full transition-colors">
                Sign In to Continue
              </router-link>
              <router-link to="/register" class="bg-[#112b4e] hover:bg-[#0e2240] text-white font-bold px-6 py-3 rounded-full transition-colors">
                Create Account
              </router-link>
            </div>
          </div>
        </div>

        <!-- STEP 1: Select Type -->
        <div v-else-if="currentStep === 0">
          <h2 class="text-2xl font-bold text-darken mb-8 text-center">Select Licence Type</h2>
          <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
            <div
              v-for="type in licenseTypes"
              :key="type.id"
              class="bg-white border border-gray-100 rounded-2xl p-6 hover:border-[#F48C06]/40 hover:shadow-md transition-all cursor-pointer group"
              @click="selectType(type)"
            >
              <div class="w-12 h-12 rounded-2xl mb-4 flex items-center justify-center shadow-sm transition-transform group-hover:scale-110" :class="type.color">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="type.icon"/>
                </svg>
              </div>
              <h3 class="font-bold text-darken mb-2 group-hover:text-accent transition-colors">{{ type.title }}</h3>
              <p class="text-sm text-gray-500 mb-4">{{ type.description }}</p>
              <span class="text-sm font-semibold text-accent flex items-center gap-1">Apply now →</span>
            </div>
          </div>

          <!-- Track applications -->
          <div class="bg-[#FEF9F2] border border-[#F48C06]/20 rounded-2xl p-6">
            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
              <div>
                <h3 class="font-bold text-darken mb-1">Track Your Application</h3>
                <p class="text-sm text-gray-500">View the status of all your submitted applications.</p>
              </div>
              <router-link to="/my-portal" class="bg-[#112b4e] hover:bg-[#0e2240] text-white font-semibold px-5 py-2.5 rounded-full transition-colors text-sm whitespace-nowrap">
                My Portal
              </router-link>
            </div>
          </div>
        </div>

        <!-- STEP 2: Fill Form -->
        <div v-else-if="currentStep === 1" class="max-w-2xl mx-auto">
          <div class="flex items-center gap-3 mb-8">
            <button aria-label="Return to application type selection" @click="currentStep = 0; selectedType = null" class="w-8 h-8 rounded-full bg-gray-100 hover:bg-gray-200 flex items-center justify-center transition-colors flex-shrink-0">
              <svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
            </button>
            <div>
              <h2 class="text-xl font-bold text-darken">{{ selectedType?.title }}</h2>
              <p class="text-sm text-gray-500">Complete the form below to begin your application</p>
            </div>
          </div>

          <div class="bg-white border border-gray-100 rounded-2xl p-8 shadow-sm space-y-5">

            <!-- Organisation details -->
            <div class="grid sm:grid-cols-2 gap-5">
              <div class="sm:col-span-2">
                <label class="block text-sm font-semibold text-darken mb-1.5">
                  {{ selectedType?.id === 'renewal' ? 'Organisation Name' : 'Organisation / Club Name' }}
                  <span class="text-red-500">*</span>
                </label>
                <input v-model="form.org_name" type="text" placeholder="Enter full name" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#F48C06]/30 focus:border-[#F48C06]"/>
                <p v-if="errors.org_name" class="text-xs text-red-500 mt-1">{{ errors.org_name }}</p>
              </div>

              <div v-if="selectedType?.id !== 'renewal'">
                <label class="block text-sm font-semibold text-darken mb-1.5">Sport Type <span class="text-red-500">*</span></label>
                <input v-model="form.sport_type" type="text" placeholder="e.g. Football, Basketball" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#F48C06]/30 focus:border-[#F48C06]"/>
                <p v-if="errors.sport_type" class="text-xs text-red-500 mt-1">{{ errors.sport_type }}</p>
              </div>

              <div v-if="selectedType?.id !== 'renewal'">
                <label class="block text-sm font-semibold text-darken mb-1.5">Year Founded</label>
                <input v-model="form.year_founded" type="number" placeholder="e.g. 2015" min="1900" :max="new Date().getFullYear()" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#F48C06]/30 focus:border-[#F48C06]"/>
              </div>

              <div v-if="selectedType?.id === 'renewal'">
                <label class="block text-sm font-semibold text-darken mb-1.5">Existing Licence Number <span class="text-red-500">*</span></label>
                <input v-model="form.existing_licence_no" type="text" placeholder="e.g. NCS/2023/0042" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#F48C06]/30 focus:border-[#F48C06]"/>
                <p v-if="errors.existing_licence_no" class="text-xs text-red-500 mt-1">{{ errors.existing_licence_no }}</p>
              </div>

              <div v-if="selectedType?.id === 'community_club'">
                <label class="block text-sm font-semibold text-darken mb-1.5">District / Region <span class="text-red-500">*</span></label>
                <input v-model="form.district" type="text" placeholder="e.g. Kampala, Gulu" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#F48C06]/30 focus:border-[#F48C06]"/>
                <p v-if="errors.district" class="text-xs text-red-500 mt-1">{{ errors.district }}</p>
              </div>

              <div>
                <label class="block text-sm font-semibold text-darken mb-1.5">Contact Person <span class="text-red-500">*</span></label>
                <input v-model="form.contact_person" type="text" placeholder="Full name" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#F48C06]/30 focus:border-[#F48C06]"/>
                <p v-if="errors.contact_person" class="text-xs text-red-500 mt-1">{{ errors.contact_person }}</p>
              </div>

              <div>
                <label class="block text-sm font-semibold text-darken mb-1.5">Contact Phone <span class="text-red-500">*</span></label>
                <input v-model="form.contact_phone" type="tel" placeholder="+256 700 000000" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#F48C06]/30 focus:border-[#F48C06]"/>
                <p v-if="errors.contact_phone" class="text-xs text-red-500 mt-1">{{ errors.contact_phone }}</p>
              </div>

              <div>
                <label class="block text-sm font-semibold text-darken mb-1.5">Contact Email <span class="text-red-500">*</span></label>
                <input v-model="form.contact_email" type="email" placeholder="email@example.com" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#F48C06]/30 focus:border-[#F48C06]"/>
                <p v-if="errors.contact_email" class="text-xs text-red-500 mt-1">{{ errors.contact_email }}</p>
              </div>

              <div :class="selectedType?.id === 'national_federation' ? 'sm:col-span-2' : ''">
                <label class="block text-sm font-semibold text-darken mb-1.5">Physical Address <span class="text-red-500">*</span></label>
                <input v-model="form.physical_address" type="text" placeholder="Street / Area, City" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#F48C06]/30 focus:border-[#F48C06]"/>
                <p v-if="errors.physical_address" class="text-xs text-red-500 mt-1">{{ errors.physical_address }}</p>
              </div>

              <div class="sm:col-span-2">
                <label class="block text-sm font-semibold text-darken mb-1.5">Additional Notes</label>
                <textarea v-model="form.notes" rows="3" placeholder="Any additional information relevant to your application..." class="w-full px-4 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#F48C06]/30 focus:border-[#F48C06] resize-none"></textarea>
              </div>
            </div>

            <div v-if="submitError" class="bg-red-50 border border-red-200 rounded-xl p-3 text-sm text-red-700">
              {{ submitError }}
            </div>

            <div class="flex gap-3 pt-2">
              <button
                @click="saveDraftAndContinue"
                :disabled="submitting"
                class="flex-1 bg-[#F48C06] hover:bg-[#d47b05] disabled:opacity-60 text-white font-bold px-6 py-3 rounded-full transition-colors text-sm"
              >
                {{ submitting ? 'Saving…' : 'Save & Continue →' }}
              </button>
              <button
                @click="currentStep = 0; selectedType = null"
                class="px-5 py-3 border border-gray-200 rounded-full text-sm text-gray-600 hover:bg-gray-50 transition-colors"
              >
                Back
              </button>
            </div>
          </div>
        </div>

        <!-- STEP 3: Submitted / Next Steps -->
        <div v-else-if="currentStep === 2" class="max-w-xl mx-auto text-center">
          <div class="bg-[#FEF9F2] border border-[#F48C06]/20 rounded-2xl p-10 shadow-sm">
            <div class="w-16 h-16 bg-green-500 rounded-2xl flex items-center justify-center mx-auto mb-5">
              <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/></svg>
            </div>
            <h2 class="text-xl font-bold text-darken mb-2">Application Draft Saved!</h2>
            <p class="text-gray-500 mb-2 text-sm">Your application for <strong>{{ selectedType?.title }}</strong> has been saved.</p>
            <p class="text-xs text-gray-400 mb-7">Reference: <strong class="text-darken">{{ draftRef }}</strong></p>

            <div class="bg-white rounded-xl p-5 text-left mb-7 space-y-3">
              <h3 class="font-bold text-darken text-sm mb-3">Next Steps to Complete Your Application</h3>
              <div class="flex items-start gap-3">
                <div class="w-6 h-6 rounded-full bg-[#F48C06] flex items-center justify-center flex-shrink-0 mt-0.5"><span class="text-white text-xs font-bold">1</span></div>
                <div>
                  <p class="text-sm font-semibold text-darken">Download & Sign Form</p>
                  <p class="text-xs text-gray-400 mt-0.5">Download the pre-filled application form, sign it, and upload the scanned copy.</p>
                </div>
              </div>
              <div class="flex items-start gap-3">
                <div class="w-6 h-6 rounded-full bg-[#F48C06] flex items-center justify-center flex-shrink-0 mt-0.5"><span class="text-white text-xs font-bold">2</span></div>
                <div>
                  <p class="text-sm font-semibold text-darken">Make Payment</p>
                  <p class="text-xs text-gray-400 mt-0.5">Pay the licence fee via bank transfer or mobile money and upload proof of payment.</p>
                </div>
              </div>
              <div class="flex items-start gap-3">
                <div class="w-6 h-6 rounded-full bg-[#112b4e] flex items-center justify-center flex-shrink-0 mt-0.5"><span class="text-white text-xs font-bold">3</span></div>
                <div>
                  <p class="text-sm font-semibold text-darken">Await Review</p>
                  <p class="text-xs text-gray-400 mt-0.5">NCS staff will review your application and notify you of the outcome.</p>
                </div>
              </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-3 justify-center">
              <router-link to="/my-portal" class="bg-[#112b4e] hover:bg-[#0e2240] text-white font-bold px-6 py-3 rounded-full transition-colors text-sm">
                Go to My Portal
              </router-link>
              <button @click="startNew" class="border border-gray-200 text-gray-700 hover:bg-gray-50 font-semibold px-6 py-3 rounded-full transition-colors text-sm">
                New Application
              </button>
            </div>
          </div>
        </div>

        <!-- Licence types info section (visible only on step 1) -->
        <div v-if="currentStep === 0 && isAuthenticated" class="mt-20">
          <div class="text-center mb-10">
            <p class="section-tag mb-2">Overview</p>
            <h2 class="text-2xl font-bold text-darken">Types of Licences</h2>
          </div>
          <div class="grid md:grid-cols-2 gap-5">
            <div
              v-for="info in licenseInfo"
              :key="info.title"
              class="bg-white border border-gray-100 rounded-2xl p-6 flex items-start gap-4 hover:shadow-sm transition-shadow"
            >
              <div class="w-10 h-10 bg-[#F48C06] text-white rounded-2xl flex items-center justify-center flex-shrink-0 font-bold text-sm shadow-sm">
                {{ info.no }}
              </div>
              <div>
                <h3 class="font-bold text-darken mb-1">{{ info.title }}</h3>
                <p class="text-sm text-gray-500">{{ info.desc }}</p>
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import apiClient from '@/api/client.js'

const route = useRoute()
const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)

const currentStep = ref(0) // 0 = select type, 1 = fill form, 2 = done
const selectedType = ref(null)
const submitting = ref(false)
const submitError = ref('')
const draftRef = ref('')

const steps = ['Select Type', 'Fill Form', 'Submit']

const licenseTypes = [
  {
    id: 'national_federation',
    title: 'National Federation',
    description: 'Register and license a national sports federation to organise competitions at a national level.',
    color: 'bg-[#112b4e]',
    icon: 'M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064',
  },
  {
    id: 'community_club',
    title: 'Community / Club',
    description: 'Register a community-level sports club or association.',
    color: 'bg-gray-600',
    icon: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0',
  },
  {
    id: 'renewal',
    title: 'Renew Existing Licence',
    description: 'Renew your sports organisation licence for the current season.',
    color: 'bg-[#F48C06]',
    icon: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15',
  },
]

const licenseInfo = [
  { no: '1', title: 'National Sports Federation Licence', desc: 'Required for bodies that oversee a specific sport at the national level and represent Uganda in international competitions.' },
  { no: '2', title: 'Regional/District Sports Association Licence', desc: 'For organisations that oversee a sport within a specific region or district.' },
  { no: '3', title: 'Community Club Registration', desc: 'For grassroots-level clubs, school sports clubs, and community associations.' },
  { no: '4', title: 'Sports Facility Licence', desc: 'Required for commercial sports facilities, stadiums and training grounds.' },
]

const form = ref({
  org_name: '',
  sport_type: '',
  year_founded: '',
  district: '',
  contact_person: '',
  contact_phone: '',
  contact_email: '',
  physical_address: '',
  existing_licence_no: '',
  notes: '',
})

const errors = ref({})

onMounted(() => {
  const typeParam = route.query.type
  if (typeParam && isAuthenticated.value) {
    const found = licenseTypes.find(t => t.id === typeParam)
    if (found) selectType(found)
  }
})

function selectType(type) {
  selectedType.value = type
  currentStep.value = 1
  errors.value = {}
  submitError.value = ''
}

function validate() {
  const e = {}
  if (!form.value.org_name.trim()) e.org_name = 'Organisation name is required'
  if (selectedType.value?.id !== 'renewal' && !form.value.sport_type.trim()) e.sport_type = 'Sport type is required'
  if (selectedType.value?.id === 'renewal' && !form.value.existing_licence_no.trim()) e.existing_licence_no = 'Existing licence number is required'
  if (selectedType.value?.id === 'community_club' && !form.value.district.trim()) e.district = 'District is required'
  if (!form.value.contact_person.trim()) e.contact_person = 'Contact person is required'
  if (!form.value.contact_phone.trim()) e.contact_phone = 'Contact phone is required'
  if (!form.value.contact_email.trim()) e.contact_email = 'Contact email is required'
  if (!form.value.physical_address.trim()) e.physical_address = 'Physical address is required'
  errors.value = e
  return !Object.keys(e).length
}

async function saveDraftAndContinue() {
  if (!validate()) return
  submitting.value = true
  submitError.value = ''
  try {
    const res = await apiClient.post(`/api/v1/applications/${selectedType.value.id}/draft`, {
      form_data: { ...form.value },
      last_saved_step: 1,
      application_type: selectedType.value.id,
      organisation_type: selectedType.value.id,
    })
    const app = res.data?.data || res.data
    draftRef.value = app?.reference_number || app?.id?.substring(0, 8)?.toUpperCase() || 'N/A'
    currentStep.value = 2
  } catch (e) {
    submitError.value = e?.response?.data?.message || 'Could not save your application. Please try again.'
  } finally {
    submitting.value = false
  }
}

function startNew() {
  selectedType.value = null
  currentStep.value = 0
  form.value = { org_name: '', sport_type: '', year_founded: '', district: '', contact_person: '', contact_phone: '', contact_email: '', physical_address: '', existing_licence_no: '', notes: '' }
  errors.value = {}
  draftRef.value = ''
}
</script>
