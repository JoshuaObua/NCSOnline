<template>
  <div class="federation-editor-card">
    <div class="card">
      <div class="card-header d-flex justify-content-between align-items-center">
        <h4>
          <i class="icofont-trophy text-primary"></i>
          {{ isEditing ? 'Edit Sports Federation' : 'Add New Sports Federation' }}
        </h4>
        <button type="button" class="btn btn-outline-secondary btn-sm" @click="$emit('cancel')">
          <i class="icofont-close"></i> Cancel
        </button>
      </div>

      <form @submit.prevent="submitForm">
        <div class="card-body">
          <!-- Notification / Collision Alert -->
          <div v-if="collisionError" class="alert alert-danger d-flex align-items-center mb-4">
            <i class="icofont-warning font-20 mr-2"></i>
            <div>{{ collisionError }}</div>
          </div>

          <!-- Section 1: Core Identification & Hybrid ID Override -->
          <div class="section-title mt-0">
            <i class="icofont-id-card"></i> Federation Identification & ID Override
          </div>
          <div class="row">
            <!-- Federation ID with Hybrid Override -->
            <div class="form-group col-lg-6">
              <label class="font-weight-bold">
                Federation ID (Unique Database Identifier)
                <span class="text-danger">*</span>
              </label>
              <div class="input-group">
                <input
                  v-model="form.id"
                  type="text"
                  class="form-control font-monospace"
                  placeholder="e.g. assoc_fufa, assoc_ukf, NCS-FED-001"
                  required
                  @input="onIdChanged"
                />
                <div class="input-group-append">
                  <button
                    type="button"
                    class="btn btn-outline-secondary"
                    title="Auto-generate ID from acronym/name"
                    @click="autoGenerateId"
                  >
                    <i class="icofont-magic"></i> Auto
                  </button>
                </div>
              </div>
              <small class="form-text text-muted">
                {{ isEditing ? 'You can override/edit this existing federation ID as long as it does not collide with another federation. All relational data will cascade automatically.' : 'Custom ID for this federation or auto-generated based on acronym.' }}
              </small>
            </div>

            <!-- Full Federation Name -->
            <div class="form-group col-lg-6">
              <label class="font-weight-bold">Full Federation / Association Name <span class="text-danger">*</span></label>
              <input
                v-model="form.name"
                type="text"
                class="form-control"
                placeholder="e.g. Federation of Uganda Football Associations"
                required
                @input="onNameChanged"
              />
            </div>

            <!-- Abbreviation / Acronym -->
            <div class="form-group col-lg-4">
              <label class="font-weight-bold">Acronym / Abbreviation <span class="text-danger">*</span></label>
              <input
                v-model="form.abbreviation"
                type="text"
                class="form-control"
                placeholder="e.g. FUFA, UKF, AUUS"
                required
                @input="onAcronymChanged"
              />
            </div>

            <!-- Category -->
            <div class="form-group col-lg-4">
              <label class="font-weight-bold">Federation Category <span class="text-danger">*</span></label>
              <select v-model="form.category" class="form-control" required>
                <option value="">Select Category</option>
                <option v-for="cat in categories" :key="cat.slug || cat.name || cat" :value="cat.name || cat">
                  {{ cat.name || cat }}
                </option>
              </select>
            </div>

            <!-- Slug -->
            <div class="form-group col-lg-4">
              <label class="font-weight-bold">Web Slug URL</label>
              <input
                v-model="form.slug"
                type="text"
                class="form-control font-monospace"
                placeholder="e.g. federation-of-uganda-football-associations"
              />
            </div>
          </div>

          <!-- Section 2: Leadership & Governance with Dynamic User/NIN Search -->
          <div class="section-title">
            <i class="icofont-business-man"></i> Leadership & Governance (Dynamic User Search by Name / NIN)
          </div>
          <div class="row">
            <!-- President / Chairperson Searchable Dropdown -->
            <div class="form-group col-lg-6">
              <label class="font-weight-bold">
                President / Chairperson Name
              </label>
              <SearchableUserSelect
                v-model="form.president"
                placeholder="Search registered user by Name, NIN, or Email..."
                :users="userList"
                @select="onPresidentSelected"
              />
              <small class="form-text text-muted">
                Dynamically searchable across registered user profiles by Full Name, NIN, or Email. You can also assign custom external names.
              </small>
            </div>

            <!-- General Secretary / CEO Searchable Dropdown -->
            <div class="form-group col-lg-6">
              <label class="font-weight-bold">
                General Secretary / CEO Name
              </label>
              <SearchableUserSelect
                v-model="form.secretary"
                placeholder="Search registered user by Name, NIN, or Email..."
                :users="userList"
                @select="onSecretarySelected"
              />
              <small class="form-text text-muted">
                Dynamically searchable across registered user profiles by Full Name, NIN, or Email. You can also assign custom external names.
              </small>
            </div>
          </div>

          <!-- Section 3: Contact & Headquarters Details -->
          <div class="section-title">
            <i class="icofont-location-pin"></i> Contact & Headquarters Details
          </div>
          <div class="row">
            <!-- Official Email -->
            <div class="form-group col-lg-6">
              <label class="font-weight-bold">Official Email Address</label>
              <input
                v-model="form.email"
                type="email"
                class="form-control"
                placeholder="e.g. admin@fufa.co.ug, info@federation.ug"
              />
              <small class="form-text text-muted">Primary administrative email for statutory notifications and communications.</small>
            </div>

            <!-- Official Phone -->
            <div class="form-group col-lg-6">
              <label class="font-weight-bold">Official Phone Number</label>
              <input
                v-model="form.phone"
                type="text"
                class="form-control"
                placeholder="e.g. +256 414 345 678"
              />
            </div>

            <!-- Physical Address -->
            <div class="form-group col-lg-6">
              <label>Physical Address / Headquarters</label>
              <input
                v-model="form.address"
                type="text"
                class="form-control"
                placeholder="e.g. FUFA House, Plot 879, Albert Cook Rd, Mengo, Kampala"
              />
            </div>

            <!-- Official Website URL -->
            <div class="form-group col-lg-6">
              <label>Official Website URL</label>
              <input
                v-model="form.website_url"
                type="url"
                class="form-control"
                placeholder="e.g. https://fufa.co.ug"
              />
            </div>
          </div>

          <!-- Section 4: Federation Logo & Media Upload -->
          <div class="section-title">
            <i class="icofont-image"></i> Federation Logo & Emblem
          </div>
          <div class="row">
            <div class="form-group col-12">
              <label class="font-weight-bold">Upload Federation Logo / Crest</label>
              <DropzoneUpload
                v-model="form.logo_url"
                label="Federation Logo"
                hint="Drag and drop or click to upload SVG, PNG, WebP or JPEG (transparent background recommended)."
                accept="image/png,image/jpeg,image/webp,image/svg+xml"
                @error="onUploadError"
              />
              <small class="form-text text-muted mt-1">
                The logo will be displayed on the national sports directory, federation profile, and public website.
              </small>
            </div>
          </div>

          <!-- Section 5: Mandate & About Details -->
          <div class="section-title">
            <i class="icofont-document-folder"></i> Statutory Mandate & Description
          </div>
          <div class="row">
            <div class="form-group col-12">
              <label>About the Federation & Statutory Mandate</label>
              <textarea
                v-model="form.description"
                class="form-control"
                rows="4"
                placeholder="Describe the background, core sporting objectives, mandate under the National Sports Act, and affiliated disciplines..."
              ></textarea>
            </div>

            <div class="form-group col-lg-3">
              <label>Display Sort Order</label>
              <input
                v-model.number="form.sort_order"
                type="number"
                class="form-control"
              />
            </div>

            <div class="form-group col-lg-3">
              <label class="d-block">Publication Status</label>
              <div class="custom-control custom-checkbox mt-2">
                <input
                  id="fed-active-toggle"
                  v-model="form.is_active"
                  type="checkbox"
                  class="custom-control-input"
                />
                <label class="custom-control-label" for="fed-active-toggle">Active & Published</label>
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
            {{ isEditing ? 'Update Federation Details' : 'Save New Federation' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import DropzoneUpload from '@/components/cms/DropzoneUpload.vue'
import SearchableUserSelect from '@/components/ui/SearchableUserSelect.vue'

const props = defineProps({
  initialModel: {
    type: Object,
    default: () => ({}),
  },
  categories: {
    type: Array,
    default: () => [],
  },
  existingFederations: {
    type: Array,
    default: () => [],
  },
  userList: {
    type: Array,
    default: () => [],
  },
  saving: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['save', 'cancel'])

const collisionError = ref('')
const originalId = ref('')

const form = reactive({
  id: '',
  name: '',
  slug: '',
  abbreviation: '',
  category: '',
  president: '',
  secretary: '',
  email: '',
  phone: '',
  address: '',
  website_url: '',
  logo_url: '',
  description: '',
  sort_order: 0,
  is_active: true,
})

const isEditing = computed(() => !!originalId.value)

function syncFromProps() {
  collisionError.value = ''
  if (props.initialModel && (props.initialModel.id || props.initialModel.name)) {
    originalId.value = props.initialModel.id || ''
    const catVal = (typeof props.initialModel.category === 'object' ? props.initialModel.category?.name : props.initialModel.category) || 'Other'
    Object.assign(form, {
      id: props.initialModel.id || '',
      name: props.initialModel.name || '',
      slug: props.initialModel.slug || '',
      abbreviation: props.initialModel.abbreviation || props.initialModel.acronym || '',
      category: catVal,
      president: props.initialModel.president || '',
      secretary: props.initialModel.secretary || '',
      email: props.initialModel.email || '',
      phone: props.initialModel.phone || '',
      address: props.initialModel.address || props.initialModel.physical_address || '',
      website_url: props.initialModel.website_url || props.initialModel.website || '',
      logo_url: props.initialModel.logo_url || '',
      description: props.initialModel.description || '',
      sort_order: props.initialModel.sort_order || 0,
      is_active: props.initialModel.is_active !== false,
    })
  } else {
    originalId.value = ''
    Object.assign(form, {
      id: '',
      name: '',
      slug: '',
      abbreviation: '',
      category: props.categories?.[0]?.name || 'Ball Sports',
      president: '',
      secretary: '',
      email: '',
      phone: '',
      address: '',
      website_url: '',
      logo_url: '',
      description: '',
      sort_order: (props.existingFederations?.length || 0) + 1,
      is_active: true,
    })
  }
}

function onPresidentSelected(user) {
  if (user && user.email && !form.email) {
    form.email = user.email
  }
  if (user && user.phone && !form.phone) {
    form.phone = user.phone
  }
}

function onSecretarySelected(user) {
  if (user && user.email && !form.email) {
    form.email = user.email
  }
  if (user && user.phone && !form.phone) {
    form.phone = user.phone
  }
}

function onUploadError(err) {
  collisionError.value = typeof err === 'string' ? err : (err?.message || 'File upload failed')
}

function slugify(text) {
  return String(text || '')
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

function onNameChanged() {
  if (!form.slug) {
    form.slug = slugify(form.name)
  }
  if (!originalId.value && !form.id) {
    autoGenerateId()
  }
}

function onAcronymChanged() {
  if (!originalId.value && (!form.id || form.id.startsWith('assoc_'))) {
    autoGenerateId()
  }
}

function onIdChanged() {
  collisionError.value = ''
}

function autoGenerateId() {
  if (form.abbreviation) {
    form.id = 'assoc_' + slugify(form.abbreviation).replace(/-/g, '_')
  } else if (form.name) {
    form.id = 'assoc_' + slugify(form.name).replace(/-/g, '_')
  }
}

function validateCollision() {
  collisionError.value = ''
  const currentId = String(form.id || '').trim().toLowerCase()
  if (!currentId) {
    collisionError.value = 'Federation ID is required.'
    return false
  }

  // Check collision with existing federations
  const duplicate = (props.existingFederations || []).find(
    f => f.id && f.id.toLowerCase() === currentId && f.id !== originalId.value
  )

  if (duplicate) {
    collisionError.value = `Federation ID "${form.id}" is already used by "${duplicate.name}". Please choose a distinct identifier.`
    return false
  }

  return true
}

function submitForm() {
  if (!validateCollision()) {
    return
  }

  if (!form.name.trim()) {
    collisionError.value = 'Federation name is required.'
    return
  }

  if (!form.slug.trim()) {
    form.slug = slugify(form.name)
  }

  const payload = {
    name: form.name.trim(),
    slug: form.slug.trim(),
    abbreviation: form.abbreviation.trim(),
    category: form.category,
    president: form.president.trim(),
    secretary: form.secretary.trim(),
    email: form.email.trim(),
    phone: form.phone.trim(),
    address: form.address.trim(),
    website_url: form.website_url.trim(),
    logo_url: form.logo_url.trim(),
    description: form.description.trim(),
    sort_order: Number(form.sort_order) || 0,
    is_active: form.is_active,
  }

  if (isEditing.value) {
    payload.id = originalId.value
    payload.new_id = form.id.trim()
  } else {
    payload.id = form.id.trim()
  }

  emit('save', payload)
}

watch(() => props.initialModel, syncFromProps, { immediate: true })
</script>

<style scoped>
.federation-editor-card {
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
