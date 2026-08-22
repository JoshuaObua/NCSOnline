<template>
  <div class="create-user-page">
    <!-- Breadcrumb & Back Navigation -->
    <nav class="page-nav">
      <button type="button" class="back-link-btn" @click="handleCancel">
        <i class="icofont-arrow-left"></i> Back to Manage Users
      </button>
      <div class="breadcrumb-trail">
        <span>Portal</span>
        <i class="icofont-simple-right"></i>
        <span>Users</span>
        <i class="icofont-simple-right"></i>
        <span class="active-trail">{{ isEditing ? 'Edit User' : 'Create New User' }}</span>
      </div>
    </nav>

    <!-- Page Header -->
    <header class="page-header">
      <div class="header-content">
        <p class="section-tag">{{ isEditing ? 'User Profile Management' : 'User Provisioning' }}</p>
        <h2 class="page-title">
          <i :class="isEditing ? 'icofont-ui-edit' : 'icofont-user-plus'"></i>
          {{ isEditing ? `Edit User Account: ${form.first_name} ${form.last_name}` : 'Add New User Account' }}
        </h2>
        <p class="page-desc">
          {{ isEditing
            ? `Update personal details, contact information, NIN, account status, and role assignments for ${form.email || 'this user'}.`
            : 'Create a new user account with credentials, contact details, and role assignments in the portal database.'
          }}
        </p>
      </div>
    </header>

    <!-- Alerts -->
    <div v-if="localError" class="alert alert-error">
      <i class="icofont-warning"></i> {{ localError }}
      <button type="button" class="alert-close" @click="localError = ''">&times;</button>
    </div>
    <div v-if="localMessage" class="alert alert-success">
      <i class="icofont-check-circled"></i> {{ localMessage }}
      <button type="button" class="alert-close" @click="localMessage = ''">&times;</button>
    </div>

    <!-- Main Creation / Editing Form Card -->
    <div class="card form-card">
      <form class="user-create-form" @submit.prevent="submitForm">
        <!-- Personal Information Section -->
        <div class="form-section">
          <div class="section-header">
            <div class="section-icon"><i class="icofont-user-alt-7"></i></div>
            <div>
              <h3 class="section-title">Personal Information</h3>
              <p class="section-sub">Enter the user's primary identity and contact details</p>
            </div>
          </div>

          <div class="form-grid-2">
            <div class="form-group">
              <label>First Name <span class="req">*</span></label>
              <input
                v-model.trim="form.first_name"
                type="text"
                required
                placeholder="e.g. Sarah"
                class="form-input"
              />
            </div>
            <div class="form-group">
              <label>Last Name <span class="req">*</span></label>
              <input
                v-model.trim="form.last_name"
                type="text"
                required
                placeholder="e.g. Namukasa"
                class="form-input"
              />
            </div>
          </div>

          <div class="form-grid-2">
            <div class="form-group">
              <label>Email Address <span class="req">*</span></label>
              <div class="input-with-icon">
                <i class="icofont-ui-email input-icon"></i>
                <input
                  v-model.trim="form.email"
                  type="email"
                  required
                  placeholder="user@ncs.go.ug"
                  class="form-input has-icon"
                />
              </div>
              <small class="field-help">Must be a valid, unique email address for authentication.</small>
            </div>

            <div class="form-group">
              <label>Phone Number</label>
              <div class="input-with-icon">
                <i class="icofont-ui-touch-phone input-icon"></i>
                <input
                  v-model.trim="form.phone"
                  type="tel"
                  placeholder="+256 700 000 000"
                  class="form-input has-icon"
                />
              </div>
            </div>

            <div class="form-group form-grid-full">
              <label>National Identification Number (NIN)</label>
              <div class="input-with-icon">
                <i class="icofont-id-card input-icon"></i>
                <input
                  v-model.trim="form.nin"
                  type="text"
                  placeholder="e.g. CM92018104NCS2"
                  class="form-input has-icon font-monospace"
                />
              </div>
              <small class="field-help">Assigned independent NIN for cross-profile identity verification (Athletes, Coaches, Officials).</small>
            </div>
          </div>
        </div>

        <hr class="form-divider" />

        <!-- Security & Credentials Section -->
        <div class="form-section">
          <div class="section-header">
            <div class="section-icon icon-security"><i class="icofont-key"></i></div>
            <div>
              <h3 class="section-title">Security & Password</h3>
              <p class="section-sub">
                {{ isEditing
                  ? 'Update or reset the user password if requested (minimum 12 characters)'
                  : 'Configure the initial login password (minimum 12 characters)'
                }}
              </p>
            </div>
          </div>

          <!-- Edit Mode Toggle for Changing Password -->
          <div v-if="isEditing" class="change-password-toggle">
            <label class="toggle-checkbox-label">
              <input v-model="changePassword" type="checkbox" />
              <span>Change User Password</span>
            </label>
            <small class="text-muted d-block mt-1">Leave unchecked to keep the user's current password unchanged.</small>
          </div>

          <div v-if="!isEditing || changePassword" class="form-group mt-3">
            <label>{{ isEditing ? 'New Password' : 'Initial Password' }} <span class="req">*</span></label>
            <div class="password-group">
              <div class="input-with-icon flex-1">
                <i class="icofont-lock input-icon"></i>
                <input
                  v-model="form.password"
                  :type="showPassword ? 'text' : 'password'"
                  :required="!isEditing || changePassword"
                  minlength="12"
                  placeholder="Enter or generate a secure password (min 12 chars)"
                  class="form-input has-icon password-input"
                />
                <button
                  type="button"
                  class="password-eye-btn"
                  :title="showPassword ? 'Hide password' : 'Show password'"
                  @click="showPassword = !showPassword"
                >
                  <i :class="showPassword ? 'icofont-eye-blocked' : 'icofont-eye'"></i>
                </button>
              </div>
              <button type="button" class="btn btn-secondary btn-generate" @click="generatePassword">
                <i class="icofont-magic"></i> Generate Strong Password
              </button>
            </div>
            <div v-if="form.password" class="password-strength-hint">
              <span :class="passwordStrengthClass">{{ passwordStrengthLabel }}</span>
              <small>{{ form.password.length }} characters</small>
            </div>
          </div>
        </div>

        <hr class="form-divider" />

        <!-- Account Status Section (Edit Mode) -->
        <div v-if="isEditing" class="form-section">
          <div class="section-header">
            <div class="section-icon icon-status"><i class="icofont-ui-settings"></i></div>
            <div>
              <h3 class="section-title">Account Status</h3>
              <p class="section-sub">Control account availability and login authorization</p>
            </div>
          </div>

          <div class="status-selection-grid">
            <label class="status-option-label" :class="{ 'is-selected': form.is_active }">
              <input v-model="form.is_active" type="radio" :value="true" name="user_status" />
              <div class="status-option-content">
                <div class="status-option-header">
                  <span class="status-indicator-dot active"></span>
                  <strong>Active Account</strong>
                </div>
                <p>User can log into the portal and perform authorized administrative actions.</p>
              </div>
            </label>

            <label class="status-option-label" :class="{ 'is-selected': !form.is_active }">
              <input v-model="form.is_active" type="radio" :value="false" name="user_status" />
              <div class="status-option-content">
                <div class="status-option-header">
                  <span class="status-indicator-dot suspended"></span>
                  <strong>Suspended / Inactive</strong>
                </div>
                <p>User login access is revoked. Current sessions will be terminated immediately.</p>
              </div>
            </label>
          </div>
        </div>

        <hr v-if="isEditing" class="form-divider" />

        <!-- Role Assignment Section -->
        <div class="form-section">
          <div class="section-header">
            <div class="section-icon icon-roles"><i class="icofont-shield-alt"></i></div>
            <div>
              <h3 class="section-title">Role Assignment</h3>
              <p class="section-sub">Select the assigned role and access tier for this user</p>
            </div>
          </div>

          <div class="roles-grid">
            <label
              v-for="role in roles"
              :key="role.id"
              class="role-card-label"
              :class="{ 'is-selected': form.role_id === role.id }"
            >
              <input
                v-model="form.role_id"
                type="radio"
                name="user_role"
                :value="role.id"
                class="role-radio"
              />
              <div class="role-card-body">
                <div class="role-card-header">
                  <strong class="role-card-title">{{ role.name }}</strong>
                  <span v-if="role.is_system" class="role-system-badge">System</span>
                </div>
                <p class="role-card-desc">{{ role.description || 'Administrative privileges role' }}</p>
              </div>
            </label>

            <!-- Default / None Option -->
            <label
              class="role-card-label"
              :class="{ 'is-selected': form.role_id === '' }"
            >
              <input
                v-model="form.role_id"
                type="radio"
                name="user_role"
                value=""
                class="role-radio"
              />
              <div class="role-card-body">
                <div class="role-card-header">
                  <strong class="role-card-title">Default User</strong>
                  <span class="role-system-badge">Standard</span>
                </div>
                <p class="role-card-desc">Standard user access with basic dashboard capabilities.</p>
              </div>
            </label>
          </div>
        </div>

        <!-- Form Actions Bar -->
        <footer class="form-footer">
          <button type="button" class="btn btn-outline" :disabled="saving" @click="handleCancel">
            Cancel
          </button>
          <button
            type="submit"
            class="btn btn-primary btn-submit"
            :disabled="saving || (!isEditing && form.password.length < 12) || (isEditing && changePassword && form.password.length < 12)"
          >
            <i v-if="saving" class="icofont-spinner icofont-spin"></i>
            <i v-else :class="isEditing ? 'icofont-check-circled' : 'icofont-plus-circle'"></i>
            {{ saving
              ? (isEditing ? 'Updating User Account...' : 'Creating User Account...')
              : (isEditing ? 'Update User Account' : 'Create User Account')
            }}
          </button>
        </footer>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as cms from '@/api/cms.js'

const props = defineProps({
  initialModel: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['saved', 'cancel', 'message', 'error'])
const route = useRoute()
const router = useRouter()

const roles = ref([])
const saving = ref(false)
const showPassword = ref(false)
const changePassword = ref(false)
const localError = ref('')
const localMessage = ref('')

const form = reactive({
  id: '',
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  nin: '',
  password: '',
  role_id: '',
  is_active: true,
})

const isEditing = computed(() => !!form.id || !!props.initialModel?.id)

function populateFromModel(model) {
  if (model && model.id) {
    form.id = model.id
    form.first_name = model.first_name || ''
    form.last_name = model.last_name || ''
    form.email = model.email || ''
    form.phone = model.phone || ''
    form.nin = model.nin || ''
    form.password = ''
    form.is_active = model.is_active !== undefined ? !!model.is_active : true
    changePassword.value = false

    // Determine initial role_id
    if (model.roles && model.roles.length) {
      const r = model.roles[0]
      form.role_id = typeof r === 'object' && r ? (r.id || '') : (typeof r === 'string' ? r : '')
    } else {
      form.role_id = ''
    }
  } else {
    form.id = ''
    form.first_name = ''
    form.last_name = ''
    form.email = ''
    form.phone = ''
    form.nin = ''
    form.password = ''
    form.role_id = ''
    form.is_active = true
    changePassword.value = false
  }
}

watch(() => props.initialModel, (newModel) => {
  populateFromModel(newModel)
}, { immediate: true, deep: true })

onMounted(async () => {
  try {
    const res = await cms.adminListRoles()
    roles.value = res.data?.data || res.data || []

    const queryUserId = route?.query?.id
    if (!props.initialModel?.id && queryUserId) {
      try {
        const uRes = await cms.adminGetUser(queryUserId)
        const u = uRes?.data?.data || uRes?.data
        if (u && u.id) {
          populateFromModel(u)
        }
      } catch (err) {
        console.warn('Could not fetch user by ID:', err)
      }
    }

    // If editing and role was a name string, resolve to role.id
    if (form.role_id && roles.value.length) {
      const matched = roles.value.find(r => r.name === form.role_id || r.id === form.role_id)
      if (matched) {
        form.role_id = matched.id
      }
    }
  } catch (e) {
    console.error('Could not load roles:', e)
  }
})

const passwordStrengthLabel = computed(() => {
  const len = form.password.length
  if (len < 12) return 'Too short (min 12 chars)'
  if (len < 16) return 'Good password strength'
  return 'Strong password'
})

const passwordStrengthClass = computed(() => {
  const len = form.password.length
  if (len < 12) return 'pwd-weak'
  if (len < 16) return 'pwd-medium'
  return 'pwd-strong'
})

onMounted(async () => {
  try {
    const res = await cms.adminListRoles()
    roles.value = res.data?.data || res.data || []

    // If editing and role was a name string, resolve to role.id
    if (props.initialModel?.roles?.length && roles.value.length) {
      const currentRole = props.initialModel.roles[0]
      const roleName = currentRole.name || (typeof currentRole === 'string' ? currentRole : '')
      if (roleName) {
        const matched = roles.value.find(r => r.name === roleName || r.id === currentRole.id)
        if (matched) {
          form.role_id = matched.id
        }
      }
    }
  } catch (e) {
    console.error('Could not load roles:', e)
  }
})

function generatePassword() {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789!@#$%^&*'
  let pwd = ''
  for (let i = 0; i < 16; i++) {
    pwd += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  form.password = pwd
  showPassword.value = true
}

async function submitForm() {
  saving.value = true
  localError.value = ''
  localMessage.value = ''

  try {
    if (isEditing.value) {
      // 1. Update user details
      await cms.adminUpdateUser(form.id, {
        first_name: form.first_name,
        last_name: form.last_name,
        email: form.email,
        phone: form.phone,
        nin: form.nin,
      })

      // 2. Update status if changed
      if (props.initialModel?.is_active !== form.is_active) {
        if (form.is_active) {
          await cms.adminActivateUser(form.id)
        } else {
          await cms.adminDeactivateUser(form.id)
        }
      }

      // 3. Update password if changed
      if (changePassword.value && form.password.length >= 12) {
        await cms.adminResetUserPassword(form.id, form.password)
      }

      // 4. Update role assignment if changed
      const oldRoleName = props.initialModel?.roles?.[0]?.name || props.initialModel?.roles?.[0] || ''
      const newRoleObj = roles.value.find(r => r.id === form.role_id)
      if (form.role_id && (!oldRoleName || oldRoleName !== newRoleObj?.name)) {
        await cms.adminAssignUserRole(form.id, form.role_id)
      }

      const successMsg = `User ${form.first_name} ${form.last_name} (${form.email}) was successfully updated.`
      emit('message', successMsg)
      emit('saved', successMsg)
    } else {
      // Create user
      if (form.password.length < 12) {
        localError.value = 'Password must be at least 12 characters long.'
        saving.value = false
        return
      }

      const res = await cms.adminCreateUser({
        first_name: form.first_name,
        last_name: form.last_name,
        email: form.email,
        phone: form.phone,
        nin: form.nin,
        password: form.password,
      })

      const newUser = res.data?.data || res.data

      if (form.role_id && newUser?.id) {
        await cms.adminAssignUserRole(newUser.id, form.role_id)
      }

      const successMsg = `User ${form.email} was successfully created.`
      emit('message', successMsg)
      emit('saved', successMsg)
    }
  } catch (err) {
    const errorMsg = err.response?.data?.error?.message || err.message || (isEditing.value ? 'Failed to update user account.' : 'Failed to create user account.')
    localError.value = errorMsg
    emit('error', errorMsg)
  } finally {
    saving.value = false
  }
}

function handleCancel() {
  emit('cancel')
}
</script>

<style scoped>
.create-user-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 960px;
  margin: 0 auto;
  padding: 10px 0 40px 0;
}

/* Nav & Breadcrumbs */
.page-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.back-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: transparent;
  border: 1px solid #cbd5e1;
  color: #334155;
  font-size: 13px;
  font-weight: 600;
  padding: 8px 14px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.back-link-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.breadcrumb-trail {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #64748b;
}

.breadcrumb-trail i {
  font-size: 10px;
}

.breadcrumb-trail .active-trail {
  color: #0f172a;
  font-weight: 700;
}

/* Header */
.page-header {
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 16px;
}

.section-tag {
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #6777ef;
  margin: 0 0 4px 0;
}

.page-title {
  font-size: 24px;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 6px 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.page-desc {
  font-size: 14px;
  color: #64748b;
  margin: 0;
}

/* Alerts */
.alert {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
}

.alert-error {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

.alert-success {
  background: #f0fdf4;
  color: #16a34a;
  border: 1px solid #bbf7d0;
}

.alert-close {
  margin-left: auto;
  background: transparent;
  border: none;
  font-size: 18px;
  line-height: 1;
  color: inherit;
  cursor: pointer;
  opacity: 0.7;
}

.alert-close:hover {
  opacity: 1;
}

/* Form Card */
.form-card {
  background: #ffffff;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
  padding: 28px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 14px;
}

.section-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.section-icon.icon-security {
  background: #fdf2f8;
  color: #db2777;
}

.section-icon.icon-roles {
  background: #f5f3ff;
  color: #7c3aed;
}

.section-icon.icon-status {
  background: #ecfdf5;
  color: #059669;
}

.section-title {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
  margin: 0 0 2px 0;
}

.section-sub {
  font-size: 13px;
  color: #64748b;
  margin: 0;
}

.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
}

.form-grid-full {
  grid-column: 1 / -1;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 13px;
  font-weight: 600;
  color: #334155;
}

.req {
  color: #ef4444;
}

.form-input {
  width: 100%;
  height: 42px;
  padding: 8px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 14px;
  color: #0f172a;
  background: #ffffff;
  transition: all 0.15s ease;
}

.font-monospace {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.form-input:focus {
  outline: none;
  border-color: #6777ef;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.15);
}

.input-with-icon {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 12px;
  color: #94a3b8;
  font-size: 16px;
  pointer-events: none;
}

.form-input.has-icon {
  padding-left: 38px;
}

.field-help {
  font-size: 12px;
  color: #64748b;
}

.form-divider {
  border: 0;
  border-top: 1px solid #f1f5f9;
  margin: 10px 0;
}

/* Password Group */
.password-group {
  display: flex;
  gap: 10px;
}

.flex-1 {
  flex: 1;
}

.password-input {
  padding-right: 40px;
}

.password-eye-btn {
  position: absolute;
  right: 12px;
  background: transparent;
  border: none;
  color: #64748b;
  font-size: 16px;
  cursor: pointer;
}

.password-eye-btn:hover {
  color: #0f172a;
}

.btn-generate {
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  color: #334155;
  font-size: 13px;
  font-weight: 600;
  padding: 0 16px;
  border-radius: 6px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
  transition: all 0.15s ease;
}

.btn-generate:hover {
  background: #f1f5f9;
  border-color: #94a3b8;
  color: #0f172a;
}

.password-strength-hint {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  margin-top: 4px;
}

.pwd-weak {
  color: #ef4444;
  font-weight: 700;
}

.pwd-medium {
  color: #f59e0b;
  font-weight: 700;
}

.pwd-strong {
  color: #10b981;
  font-weight: 700;
}

/* Change password toggle */
.change-password-toggle {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 12px 16px;
}

.toggle-checkbox-label {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  font-weight: 700;
  color: #1e293b;
  cursor: pointer;
}

.toggle-checkbox-label input[type='checkbox'] {
  width: 18px;
  height: 18px;
  accent-color: #6777ef;
  cursor: pointer;
}

/* Status selection grid */
.status-selection-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.status-option-label {
  position: relative;
  display: flex;
  border: 1.5px solid #e2e8f0;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  background: #ffffff;
  transition: all 0.15s ease;
}

.status-option-label input[type='radio'] {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 18px;
  height: 18px;
  accent-color: #6777ef;
}

.status-option-label:hover {
  border-color: #cbd5e1;
  background: #f8fafc;
}

.status-option-label.is-selected {
  border-color: #6777ef;
  background: #f0f4ff;
}

.status-option-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding-right: 28px;
}

.status-option-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #1e293b;
}

.status-indicator-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.status-indicator-dot.active {
  background: #10b981;
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
}

.status-indicator-dot.suspended {
  background: #ef4444;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.2);
}

.status-option-content p {
  margin: 0;
  font-size: 12px;
  color: #64748b;
  line-height: 1.4;
}

/* Roles Grid */
.roles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 14px;
}

.role-card-label {
  display: block;
  position: relative;
  border: 1.5px solid #e2e8f0;
  border-radius: 8px;
  padding: 14px;
  cursor: pointer;
  background: #ffffff;
  transition: all 0.15s ease;
}

.role-radio {
  position: absolute;
  top: 14px;
  right: 14px;
  width: 16px;
  height: 16px;
  accent-color: #6777ef;
}

.role-card-label:hover {
  border-color: #cbd5e1;
  background: #f8fafc;
}

.role-card-label.is-selected {
  border-color: #6777ef;
  background: #f0f4ff;
}

.role-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  padding-right: 24px;
}

.role-card-title {
  font-size: 13px;
  color: #0f172a;
}

.role-system-badge {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  padding: 2px 6px;
  border-radius: 4px;
  background: #f1f5f9;
  color: #475569;
}

.role-card-desc {
  font-size: 12px;
  color: #64748b;
  margin: 0;
  line-height: 1.35;
}

/* Footer Actions */
.form-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 14px;
  border-top: 1px solid #f1f5f9;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 700;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-outline {
  background: transparent;
  border: 1px solid #cbd5e1;
  color: #475569;
}

.btn-outline:hover {
  background: #f8fafc;
  color: #0f172a;
}

.btn-primary {
  background: #6777ef;
  border: 1px solid #6777ef;
  color: #ffffff;
  box-shadow: 0 2px 8px rgba(103, 119, 239, 0.35);
}

.btn-primary:hover:not(:disabled) {
  background: #5566de;
  border-color: #5566de;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Dark Mode Support */
:global(.dark) .back-link-btn {
  border-color: #334155;
  color: #cbd5e1;
}

:global(.dark) .back-link-btn:hover {
  background: #1e293b;
  color: #ffffff;
}

:global(.dark) .page-header {
  border-color: #334155;
}

:global(.dark) .page-title {
  color: #f8fafc;
}

:global(.dark) .page-desc {
  color: #94a3b8;
}

:global(.dark) .form-card {
  background: #1e293b;
  border-color: #334155;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

:global(.dark) .section-title {
  color: #f8fafc;
}

:global(.dark) .section-sub {
  color: #94a3b8;
}

:global(.dark) .form-group label {
  color: #e2e8f0;
}

:global(.dark) .form-input {
  background: #0f172a;
  border-color: #334155;
  color: #f8fafc;
}

:global(.dark) .form-input:focus {
  border-color: #818cf8;
  box-shadow: 0 0 0 3px rgba(129, 140, 248, 0.2);
}

:global(.dark) .btn-generate {
  background: #0f172a;
  border-color: #334155;
  color: #cbd5e1;
}

:global(.dark) .btn-generate:hover {
  background: #1e293b;
  color: #f8fafc;
}

:global(.dark) .change-password-toggle {
  background: #0f172a;
  border-color: #334155;
}

:global(.dark) .toggle-checkbox-label {
  color: #f8fafc;
}

:global(.dark) .status-option-label {
  background: #0f172a;
  border-color: #334155;
}

:global(.dark) .status-option-label:hover {
  background: #1e293b;
  border-color: #475569;
}

:global(.dark) .status-option-label.is-selected {
  background: rgba(103, 119, 239, 0.15);
  border-color: #818cf8;
}

:global(.dark) .status-option-header {
  color: #f8fafc;
}

:global(.dark) .role-card-label {
  background: #0f172a;
  border-color: #334155;
}

:global(.dark) .role-card-label:hover {
  background: #1e293b;
  border-color: #475569;
}

:global(.dark) .role-card-label.is-selected {
  background: rgba(103, 119, 239, 0.15);
  border-color: #818cf8;
}

:global(.dark) .role-card-title {
  color: #f8fafc;
}

:global(.dark) .role-system-badge {
  background: #1e293b;
  color: #94a3b8;
}

:global(.dark) .form-divider {
  border-color: #334155;
}

:global(.dark) .form-footer {
  border-color: #334155;
}

:global(.dark) .btn-outline {
  border-color: #334155;
  color: #cbd5e1;
}

:global(.dark) .btn-outline:hover {
  background: #0f172a;
  color: #f8fafc;
}

@media (max-width: 768px) {
  .form-grid-2 {
    grid-template-columns: 1fr;
  }
  .status-selection-grid {
    grid-template-columns: 1fr;
  }
  .password-group {
    flex-direction: column;
  }
  .form-card {
    padding: 18px;
  }
}
</style>
