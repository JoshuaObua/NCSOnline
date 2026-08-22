<template>
  <div class="official-editor-card">
    <div class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <h4>
          <i class="icofont-user-plus text-primary"></i>
          {{ isEditing ? 'Edit Federation Official' : 'Add New Federation Official' }}
        </h4>
        <button type="button" class="btn btn-outline-secondary btn-sm" @click="$emit('cancel')">
          <i class="icofont-close"></i> Cancel
        </button>
      </div>

      <form @submit.prevent="submitForm">
        <div class="card-body">
          <div v-if="errorMessage" class="alert alert-danger d-flex align-items-center mb-4">
            <i class="icofont-warning font-20 mr-2"></i>
            <div>{{ errorMessage }}</div>
          </div>

          <!-- Section 1: Federation Selection (Searchable Dropdown) -->
          <div class="section-title mt-0">
            <i class="icofont-trophy"></i> Sports Federation Assignment
          </div>
          <div class="row">
            <div class="form-group col-12">
              <label class="font-weight-bold">
                Select Sports Federation / Association <span class="text-danger">*</span>
              </label>
              <SearchableFederationSelect
                v-model="form.federation_id"
                :federations="federationsList"
                placeholder="Type to search federation name, acronym or ID..."
                required
              />
              <small class="form-text text-muted">
                Search and select the national sports governing body this official is assigned to.
              </small>
            </div>
          </div>

          <!-- Section 2: Official Identification & Position (User / NIN Searchable Select) -->
          <div class="section-title">
            <i class="icofont-tie"></i> Official Information & Identification (Name / NIN)
          </div>
          <div class="row">
            <!-- Full Name with User/NIN search -->
            <div class="form-group col-lg-6">
              <label class="font-weight-bold">
                Official User Profile / Full Name <span class="text-danger">*</span>
              </label>
              <SearchableUserSelect
                v-model="form.full_name"
                placeholder="Search registered user by Name, NIN, or Email..."
                @select="onUserSelected"
              />
              <small class="form-text text-muted">
                Dynamically searchable across registered user profiles by Full Name, NIN, or Email. You can also assign custom external names.
              </small>
            </div>

            <!-- National Identification Number (NIN) -->
            <div class="form-group col-lg-6">
              <label class="font-weight-bold">
                National Identification Number (NIN)
              </label>
              <div class="input-group">
                <div class="input-group-prepend">
                  <span class="input-group-text"><i class="icofont-id-card"></i></span>
                </div>
                <input
                  v-model="form.nin"
                  type="text"
                  class="form-control font-monospace"
                  placeholder="e.g. CM92018104NCS2"
                />
              </div>
              <small class="form-text text-muted">
                Assigned independent NIN linking official across all sporting and administrative profiles.
              </small>
            </div>

            <!-- Position / Role -->
            <div class="form-group col-lg-6">
              <label class="font-weight-bold">Designated Position / Role <span class="text-danger">*</span></label>
              <select v-model="form.position" class="form-control" required @change="onPositionSelectChange">
                <option value="">Select official position...</option>
                <option value="PRESIDENT">President / Chairperson</option>
                <option value="VICE_PRESIDENT">Vice President</option>
                <option value="GENERAL_SECRETARY">General Secretary / CEO</option>
                <option value="TREASURER">Treasurer / Finance Director</option>
                <option value="TECHNICAL_DIRECTOR">Technical Director</option>
                <option value="COMMITTEE_MEMBER">Executive Committee Member</option>
                <option value="PUBLIC_RELATIONS_OFFICER">Public Relations / Media Officer</option>
                <option value="MEDICAL_OFFICER">Medical Officer / Doctor</option>
                <option value="LEGAL_OFFICER">Legal Counsel</option>
                <option value="OTHER">Other / Custom Title</option>
              </select>
            </div>

            <!-- Custom Position Label -->
            <div v-if="form.position === 'OTHER'" class="form-group col-lg-6">
              <label class="font-weight-bold">Custom Position Title <span class="text-danger">*</span></label>
              <input
                v-model="form.position_label"
                type="text"
                class="form-control"
                placeholder="e.g. Head of Youth Development, Safeguarding Officer"
                required
              />
            </div>
          </div>

          <!-- Section 3: Contact & Term Information -->
          <div class="section-title">
            <i class="icofont-calendar"></i> Contact & Term Tenure
          </div>
          <div class="row">
            <!-- Email -->
            <div class="form-group col-lg-6">
              <label>Official Email Address</label>
              <input
                v-model="form.email"
                type="email"
                class="form-control"
                placeholder="e.g. official@federation.ug"
              />
            </div>

            <!-- Phone -->
            <div class="form-group col-lg-6">
              <label>Official Phone Number</label>
              <input
                v-model="form.phone"
                type="text"
                class="form-control"
                placeholder="e.g. +256 700 000 000"
              />
            </div>

            <!-- Appointed On -->
            <div class="form-group col-lg-4">
              <label>Appointment / Election Date</label>
              <input
                v-model="form.appointed_on"
                type="date"
                class="form-control"
              />
            </div>

            <!-- Term Ends On -->
            <div class="form-group col-lg-4">
              <label>Term Expiry Date</label>
              <input
                v-model="form.term_ends_on"
                type="date"
                class="form-control"
              />
            </div>

            <!-- Active Status -->
            <div class="form-group col-lg-4">
              <label class="d-block">Status</label>
              <div class="custom-control custom-checkbox mt-2">
                <input
                  id="official-active-toggle"
                  v-model="form.is_active"
                  type="checkbox"
                  class="custom-control-input"
                />
                <label class="custom-control-label" for="official-active-toggle">Currently Serving / Active</label>
              </div>
            </div>
          </div>
        </div>

        <!-- Card Footer -->
        <div class="card-footer text-right d-flex justify-content-between align-items-center">
          <button type="button" class="btn btn-secondary" @click="$emit('cancel')">
            Cancel
          </button>
          <button type="submit" class="btn btn-primary btn-lg" :disabled="saving">
            <i v-if="saving" class="icofont-spinner icofont-spin mr-1"></i>
            <i v-else class="icofont-save mr-1"></i>
            {{ isEditing ? 'Update Official' : 'Save Federation Official' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import SearchableFederationSelect from '@/components/ui/SearchableFederationSelect.vue'
import SearchableUserSelect from '@/components/ui/SearchableUserSelect.vue'

const props = defineProps({
  initialModel: {
    type: Object,
    default: () => ({}),
  },
  federations: {
    type: Array,
    default: () => [],
  },
  saving: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['save', 'cancel'])

const errorMessage = ref('')
const federationsList = computed(() => props.federations || [])

const form = reactive({
  id: '',
  federation_id: '',
  full_name: '',
  nin: '',
  user_id: '',
  position: 'PRESIDENT',
  position_label: 'President / Chairperson',
  email: '',
  phone: '',
  appointed_on: '',
  term_ends_on: '',
  is_active: true,
})

const isEditing = computed(() => !!props.initialModel?.id)

function syncFromProps() {
  errorMessage.value = ''
  if (props.initialModel && props.initialModel.id) {
    Object.assign(form, {
      id: props.initialModel.id,
      federation_id: props.initialModel.federation_id || '',
      full_name: props.initialModel.full_name || '',
      nin: props.initialModel.nin || '',
      user_id: props.initialModel.user_id || '',
      position: props.initialModel.position || 'OTHER',
      position_label: props.initialModel.position_label || '',
      email: props.initialModel.email || '',
      phone: props.initialModel.phone || '',
      appointed_on: props.initialModel.appointed_on ? String(props.initialModel.appointed_on).slice(0, 10) : '',
      term_ends_on: props.initialModel.term_ends_on ? String(props.initialModel.term_ends_on).slice(0, 10) : '',
      is_active: props.initialModel.is_active !== false,
    })
  } else {
    Object.assign(form, {
      id: '',
      federation_id: '',
      full_name: '',
      nin: '',
      user_id: '',
      position: 'PRESIDENT',
      position_label: 'President / Chairperson',
      email: '',
      phone: '',
      appointed_on: '',
      term_ends_on: '',
      is_active: true,
    })
  }
}

function onUserSelected(user) {
  if (user) {
    if (user.nin) form.nin = user.nin
    if (user.email && !form.email) form.email = user.email
    if (user.phone && !form.phone) form.phone = user.phone
    if (user.id) form.user_id = user.id
  }
}

function onPositionSelectChange() {
  const map = {
    PRESIDENT: 'President / Chairperson',
    VICE_PRESIDENT: 'Vice President',
    GENERAL_SECRETARY: 'General Secretary / CEO',
    TREASURER: 'Treasurer',
    TECHNICAL_DIRECTOR: 'Technical Director',
    COMMITTEE_MEMBER: 'Executive Committee Member',
    PUBLIC_RELATIONS_OFFICER: 'Public Relations Officer',
    MEDICAL_OFFICER: 'Medical Officer',
    LEGAL_OFFICER: 'Legal Counsel',
  }
  if (form.position !== 'OTHER') {
    form.position_label = map[form.position] || form.position
  } else {
    form.position_label = ''
  }
}

function submitForm() {
  errorMessage.value = ''
  if (!form.federation_id) {
    errorMessage.value = 'Please select a sports federation.'
    return
  }
  if (!form.full_name.trim()) {
    errorMessage.value = 'Official full name is required.'
    return
  }
  if (!form.position) {
    errorMessage.value = 'Position is required.'
    return
  }
  if (form.position === 'OTHER' && !form.position_label.trim()) {
    errorMessage.value = 'Please specify the custom position title.'
    return
  }

  const payload = {
    federation_id: form.federation_id,
    full_name: form.full_name.trim(),
    nin: form.nin.trim(),
    user_id: form.user_id || undefined,
    position: form.position,
    position_label: form.position_label.trim() || form.position,
    email: form.email.trim(),
    phone: form.phone.trim(),
    appointed_on: form.appointed_on || null,
    term_ends_on: form.term_ends_on || null,
    is_active: form.is_active,
  }

  if (form.id) {
    payload.id = form.id
  }

  emit('save', payload)
}

watch(() => props.initialModel, syncFromProps, { immediate: true })
</script>

<style scoped>
.official-editor-card {
  max-width: 100%;
}
.section-title {
  font-size: 14px;
  font-weight: 700;
  color: #34395e;
  margin-top: 24px;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f4f6f9;
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
