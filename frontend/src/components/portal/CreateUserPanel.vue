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
        <span class="active-trail">Create New User</span>
      </div>
    </nav>

    <!-- Page Header -->
    <header class="page-header">
      <div class="header-content">
        <p class="section-tag">User Provisioning</p>
        <h2 class="page-title">Add New User Account</h2>
        <p class="page-desc">Create a new user account with credentials, contact details, and role assignments in the portal database.</p>
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

    <!-- Main Creation Form Card -->
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
          </div>
        </div>

        <hr class="form-divider" />

        <!-- Security & Credentials Section -->
        <div class="form-section">
          <div class="section-header">
            <div class="section-icon icon-security"><i class="icofont-key"></i></div>
            <div>
              <h3 class="section-title">Security & Password</h3>
              <p class="section-sub">Configure the initial login password (minimum 12 characters)</p>
            </div>
          </div>

          <div class="form-group">
            <label>Initial Password <span class="req">*</span></label>
            <div class="password-group">
              <div class="input-with-icon flex-1">
                <i class="icofont-lock input-icon"></i>
                <input
                  v-model="form.password"
                  :type="showPassword ? 'text' : 'password'"
                  required
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

        <!-- Role Assignment Section -->
        <div class="form-section">
          <div class="section-header">
            <div class="section-icon icon-roles"><i class="icofont-shield-alt"></i></div>
            <div>
              <h3 class="section-title">Role Assignment</h3>
              <p class="section-sub">Select the initial role and permission level for this user</p>
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
                <p class="role-card-desc">{{ role.description || 'Custom administrative permissions role' }}</p>
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
          <button type="submit" class="btn btn-primary btn-submit" :disabled="saving || form.password.length < 12">
            <i v-if="saving" class="icofont-spinner icofont-spin"></i>
            <i v-else class="icofont-check-circled"></i>
            {{ saving ? 'Creating User Account...' : 'Create User Account' }}
          </button>
        </footer>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import * as cms from '@/api/cms.js'

const emit = defineEmits(['saved', 'cancel', 'message', 'error'])

const roles = ref([])
const saving = ref(false)
const showPassword = ref(false)
const localError = ref('')
const localMessage = ref('')

const form = reactive({
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  password: '',
  role_id: '',
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
  if (form.password.length < 12) {
    localError.value = 'Password must be at least 12 characters long.'
    return
  }

  saving.value = true
  localError.value = ''
  localMessage.value = ''

  try {
    const res = await cms.adminCreateUser({
      first_name: form.first_name,
      last_name: form.last_name,
      email: form.email,
      phone: form.phone,
      password: form.password,
    })

    const newUser = res.data?.data || res.data

    if (form.role_id && newUser?.id) {
      await cms.adminAssignUserRole(newUser.id, form.role_id)
    }

    const successMsg = `User ${form.email} was successfully created.`
    emit('message', successMsg)
    emit('saved', newUser)
  } catch (err) {
    localError.value = err.response?.data?.error?.message || err.message || 'Failed to create user.'
    emit('error', localError.value)
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
  gap: 22px;
  width: 100%;
  max-width: 960px;
  margin: 0 auto;
}

/* Navigation & Breadcrumbs */
.page-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
  padding-bottom: 10px;
}

.back-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 6px;
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #cbd5e1;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.back-link-btn:hover {
  background: #6777ef;
  color: #ffffff;
  border-color: #6777ef;
}

.breadcrumb-trail {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  font-weight: 600;
  color: #94a3b8;
}

.breadcrumb-trail .active-trail {
  color: #6777ef;
  font-weight: 700;
}

/* Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.section-tag {
  margin: 0 0 4px;
  font-size: 11px;
  font-weight: 800;
  color: #6777ef;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.page-title {
  margin: 0 0 6px;
  font-size: 26px;
  font-weight: 800;
  color: #1e293b;
}

.page-desc {
  margin: 0;
  color: #64748b;
  font-size: 13px;
}

/* Card */
.form-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
  padding: 30px 34px;
}

.user-create-form {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

/* Form Sections */
.form-section {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 4px;
}

.section-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.icon-security { background: #fef3c7; color: #d97706; }
.icon-roles { background: #f5f3ff; color: #7c3aed; }

.section-title {
  margin: 0;
  font-size: 16px;
  font-weight: 800;
  color: #1e293b;
}

.section-sub {
  margin: 2px 0 0;
  font-size: 12px;
  color: #64748b;
}

.form-divider {
  border: none;
  border-top: 1px solid #f1f5f9;
  margin: 0;
}

/* Grid & Inputs */
.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 13px;
  font-weight: 700;
  color: #334155;
}

.req {
  color: #ef4444;
}

.field-help {
  font-size: 11px;
  color: #94a3b8;
}

.form-input {
  width: 100%;
  height: 44px;
  padding: 0 14px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-family: inherit;
  font-size: 13px;
  background: #ffffff;
  color: #1e293b;
  box-sizing: border-box;
  transition: all 0.2s ease;
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
  left: 14px;
  color: #94a3b8;
  font-size: 16px;
  pointer-events: none;
}

.form-input.has-icon {
  padding-left: 40px;
}

/* Password Group */
.password-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.flex-1 { flex: 1; }

.password-input {
  padding-right: 40px;
}

.password-eye-btn {
  position: absolute;
  right: 12px;
  border: none;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;
  font-size: 18px;
}

.btn-generate {
  height: 44px;
  white-space: nowrap;
  font-size: 12px;
}

.password-strength-hint {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 11px;
  font-weight: 700;
  margin-top: 4px;
}

.pwd-weak { color: #ef4444; }
.pwd-medium { color: #f59e0b; }
.pwd-strong { color: #10b981; }

/* Roles Grid Selection */
.roles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 14px;
}

.role-card-label {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #f8fafc;
  cursor: pointer;
  transition: all 0.15s ease;
}

.role-card-label:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.role-card-label.is-selected {
  background: #f5f3ff;
  border-color: #8b5cf6;
  box-shadow: 0 0 0 2px rgba(139, 92, 246, 0.2);
}

.role-radio {
  margin-top: 3px;
  accent-color: #7c3aed;
  cursor: pointer;
}

.role-card-body {
  display: flex;
  flex-direction: column;
  gap: 3px;
  min-width: 0;
}

.role-card-header {
  display: flex;
  align-items: center;
  gap: 6px;
}

.role-card-title {
  font-size: 13px;
  color: #1e293b;
}

.role-system-badge {
  font-size: 9px;
  font-weight: 800;
  padding: 1px 5px;
  background: #e2e8f0;
  color: #475569;
  border-radius: 4px;
  text-transform: uppercase;
}

.role-card-desc {
  margin: 0;
  font-size: 11px;
  color: #64748b;
  line-height: 1.4;
}

/* Footer & Buttons */
.form-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 14px;
  padding-top: 10px;
  border-top: 1px solid #f1f5f9;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 8px;
  font-family: inherit;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  border: 1px solid transparent;
}

.btn-primary { background: #6777ef; color: #ffffff; }
.btn-primary:hover:not(:disabled) { background: #5563db; }

.btn-secondary { background: #eff6ff; color: #2563eb; border-color: #bfdbfe; }
.btn-secondary:hover:not(:disabled) { background: #2563eb; color: #ffffff; }

.btn-outline { background: #ffffff; color: #475569; border-color: #cbd5e1; }
.btn-outline:hover:not(:disabled) { background: #f8fafc; border-color: #94a3b8; }

.btn-submit { padding: 11px 24px; font-size: 14px; }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }

/* Alerts */
.alert {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 18px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
}

.alert-error { background: #fff1f0; color: #dc2626; border: 1px solid #fecaca; }
.alert-success { background: #ecfdf5; color: #059669; border: 1px solid #a7f3d0; }
.alert-close { border: none; background: transparent; font-size: 18px; cursor: pointer; color: inherit; }

@media (max-width: 768px) {
  .form-grid-2 { grid-template-columns: 1fr; }
  .password-group { flex-direction: column; align-items: stretch; }
  .form-card { padding: 20px; }
}
</style>
