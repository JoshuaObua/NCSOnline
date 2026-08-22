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

          <!-- Section 2: Leadership & Governance -->
          <div class="section-title">
            <i class="icofont-business-man"></i> Leadership & Governance
          </div>
          <div class="row">
            <div class="form-group col-lg-6">
              <label>President / Chairperson Name</label>
              <input
                v-model="form.president"
                type="text"
                class="form-control"
                placeholder="e.g. Eng. Moses Magogo"
              />
            </div>

            <div class="form-group col-lg-6">
              <label>General Secretary / CEO Name</label>
              <input
                v-model="form.secretary"
                type="text"
                class="form-control"
                placeholder="e.g. Edgar Watson"
              />
            </div>
          </div>

          <!-- Section 3: Contact & Headquarters -->
          <div class="section-title">
            <i class="icofont-location-pin"></i> Contact & Headquarters Details
          </div>
          <div class="row">
            <div class="form-group col-lg-6">
              <label>Physical Address / Headquarters</label>
              <input
                v-model="form.address"
                type="text"
                class="form-control"
                placeholder="e.g. FUFA House, Plot 879, Albert Cook Rd, Mengo, Kampala"
              />
            </div>

            <div class="form-group col-lg-3">
              <label>Official Phone Number</label>
              <input
                v-model="form.phone"
                type="text"
                class="form-control"
                placeholder="e.g. +256 414 345 678"
              />
            </div>

            <div class="form-group col-lg-3">
              <label>Official Website URL</label>
              <input
                v-model="form.website_url"
                type="url"
                class="form-control"
                placeholder="e.g. https://fufa.co.ug"
              />
            </div>

            <div class="form-group col-lg-8">
              <label>Federation Logo URL</label>
              <input
                v-model="form.logo_url"
                type="text"
                class="form-control"
                placeholder="e.g. /uploads/logos/fufa.png or https://..."
              />
            </div>

            <div class="form-group col-lg-2">
              <label>Sort Order</label>
              <input
                v-model.number="form.sort_order"
                type="number"
                class="form-control"
              />
            </div>

            <div class="form-group col-lg-2">
              <label class="d-block">Status</label>
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

            <div class="form-group col-12">
              <label>Federation Mandate & Overview Description</label>
              <textarea
                v-model="form.description"
                class="form-control otika-textarea"
                rows="4"
                placeholder="Describe the mandate, sport disciplines governed, and background of this national federation..."
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Footer actions -->
        <div class="card-footer text-right d-flex justify-content-between align-items-center">
          <button type="button" class="btn btn-secondary" @click="$emit('cancel')">
            Cancel
          </button>
          <button type="submit" class="btn btn-primary btn-lg" :disabled="saving">
            <i v-if="saving" class="icofont-spinner icofont-spin mr-1"></i>
            <i v-else class="icofont-save mr-1"></i>
            {{ isEditing ? 'Update Federation' : 'Create Federation' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'

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
  saving: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['save', 'cancel'])

const originalId = ref('')
const collisionError = ref('')
const isCustomId = ref(false)

const form = reactive({
  id: '',
  name: '',
  slug: '',
  abbreviation: '',
  category: 'Other',
  president: '',
  secretary: '',
  phone: '',
  address: '',
  website_url: '',
  description: '',
  logo_url: '',
  sort_order: 0,
  is_active: true,
})

const isEditing = computed(() => !!originalId.value)

function normalizeSlug(val) {
  return String(val || '').toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '')
}

function syncFromProps() {
  collisionError.value = ''
  if (props.initialModel && props.initialModel.id) {
    originalId.value = props.initialModel.id
    Object.assign(form, {
      id: props.initialModel.id,
      name: props.initialModel.name || '',
      slug: props.initialModel.slug || '',
      abbreviation: props.initialModel.abbreviation || props.initialModel.acronym || '',
      category: (typeof props.initialModel.category === 'object' ? props.initialModel.category?.name : props.initialModel.category) || 'Other',
      president: props.initialModel.president || '',
      secretary: props.initialModel.secretary || '',
      phone: props.initialModel.phone || '',
      address: props.initialModel.address || props.initialModel.physical_address || '',
      website_url: props.initialModel.website_url || props.initialModel.website || '',
      description: props.initialModel.description || '',
      logo_url: props.initialModel.logo_url || '',
      sort_order: Number(props.initialModel.sort_order || 0),
      is_active: props.initialModel.is_active !== false,
    })
    isCustomId.value = true
  } else {
    originalId.value = ''
    isCustomId.value = false
    Object.assign(form, {
      id: '',
      name: '',
      slug: '',
      abbreviation: '',
      category: props.categories[0]?.name || 'Other',
      president: '',
      secretary: '',
      phone: '',
      address: '',
      website_url: '',
      description: '',
      logo_url: '',
      sort_order: 0,
      is_active: true,
    })
  }
}

function onIdChanged() {
  isCustomId.value = true
  collisionError.value = ''
}

function onNameChanged() {
  if (!isEditing.value && !form.slug) {
    form.slug = normalizeSlug(form.name)
  }
  if (!isCustomId.value && !isEditing.value && !form.id) {
    autoGenerateId()
  }
}

function onAcronymChanged() {
  if (!isCustomId.value && !isEditing.value) {
    autoGenerateId()
  }
}

function autoGenerateId() {
  const code = form.abbreviation || form.slug || normalizeSlug(form.name)
  if (code) {
    form.id = 'assoc_' + normalizeSlug(code)
  }
}

function submitForm() {
  collisionError.value = ''
  const newId = String(form.id || '').trim()

  if (!newId) {
    collisionError.value = 'Federation ID cannot be blank.'
    return
  }

  // Client-side collision check against other federations
  const collision = props.existingFederations.find(f => f.id === newId && f.id !== originalId.value)
  if (collision) {
    collisionError.value = `Federation ID "${newId}" is already used by "${collision.name}". Please choose another unique identifier.`
    return
  }

  if (!form.slug) {
    form.slug = normalizeSlug(form.name)
  }

  emit('save', {
    ...form,
    id: originalId.value || newId,
    new_id: newId,
  })
}

watch(() => props.initialModel, syncFromProps, { immediate: true })
</script>

<style scoped>
.federation-editor-card {
  max-width: 100%;
}

.font-monospace {
  font-family: 'Courier New', Courier, monospace;
}
</style>
