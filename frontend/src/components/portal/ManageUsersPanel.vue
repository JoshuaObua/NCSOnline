<template>
  <div class="manage-users-panel">
    <!-- Header & Action Bar -->
    <header class="panel-header">
      <div>
        <p class="panel-subtitle">Access & Identity Management</p>
        <h2 class="panel-title">Manage Users</h2>
        <span class="panel-desc">View, create, edit, suspend, reset passwords, and manage roles for all portal users in the database.</span>
      </div>
      <div class="header-actions">
        <button type="button" class="btn btn-outline" :disabled="loading" @click="loadUsers">
          <i class="icofont-refresh" :class="{ 'icofont-spin': loading }"></i> Refresh
        </button>
        <button id="btn-add-new-user" type="button" class="btn btn-primary btn-add-new-user" @click="handleCreateUser">
          <i class="icofont-plus-circle"></i> Add New User
        </button>
      </div>
    </header>

    <!-- KPI Summary Cards -->
    <section class="kpi-grid">
      <div class="kpi-card">
        <div class="kpi-icon icon-users"><i class="icofont-users-alt-2"></i></div>
        <div class="kpi-info">
          <span class="kpi-label">Total Users</span>
          <strong class="kpi-value">{{ users.length }}</strong>
        </div>
      </div>
      <div class="kpi-card">
        <div class="kpi-icon icon-active"><i class="icofont-check-circled"></i></div>
        <div class="kpi-info">
          <span class="kpi-label">Active Accounts</span>
          <strong class="kpi-value">{{ activeUsersCount }}</strong>
        </div>
      </div>
      <div class="kpi-card">
        <div class="kpi-icon icon-suspended"><i class="icofont-ban"></i></div>
        <div class="kpi-info">
          <span class="kpi-label">Suspended / Inactive</span>
          <strong class="kpi-value">{{ suspendedUsersCount }}</strong>
        </div>
      </div>
      <div class="kpi-card">
        <div class="kpi-icon icon-admin"><i class="icofont-shield-alt"></i></div>
        <div class="kpi-info">
          <span class="kpi-label">Administrators</span>
          <strong class="kpi-value">{{ adminUsersCount }}</strong>
        </div>
      </div>
    </section>

    <!-- Search & Filter Controls -->
    <section class="filter-toolbar">
      <div class="search-box">
        <i class="icofont-search"></i>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by name, email, or phone..."
          @input="handleSearchInput"
        />
        <button v-if="searchQuery" type="button" class="clear-search-btn" @click="searchQuery = ''; filterUsers()">
          <i class="icofont-close-line"></i>
        </button>
      </div>
      <div class="filter-dropdowns">
        <select v-model="statusFilter" class="filter-select" @change="filterUsers">
          <option value="ALL">All Statuses</option>
          <option value="ACTIVE">Active</option>
          <option value="SUSPENDED">Suspended / Inactive</option>
        </select>
        <select v-model="roleFilter" class="filter-select" @change="filterUsers">
          <option value="ALL">All Roles</option>
          <option v-for="r in roles" :key="r.id || r.name" :value="r.name">{{ r.name }}</option>
        </select>
      </div>
    </section>

    <!-- Error / Success Alerts -->
    <div v-if="localError" class="alert alert-error">
      <i class="icofont-warning"></i> {{ localError }}
      <button type="button" class="alert-close" @click="localError = ''">&times;</button>
    </div>
    <div v-if="localMessage" class="alert alert-success">
      <i class="icofont-check-circled"></i> {{ localMessage }}
      <button type="button" class="alert-close" @click="localMessage = ''">&times;</button>
    </div>

    <!-- Users Table Container -->
    <section class="table-card">
      <div class="table-responsive">
        <table class="users-table">
          <thead>
            <tr>
              <th>User</th>
              <th>Status</th>
              <th>Roles</th>
              <th>Created Date</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in filteredUsers" :key="user.id" class="user-row">
              <!-- User Details -->
              <td class="user-cell">
                <div class="user-identity">
                  <div class="user-avatar" :class="getAvatarColorClass(user)">
                    {{ getInitials(user) }}
                  </div>
                  <div class="user-text">
                    <strong class="user-full-name">
                      {{ user.first_name }} {{ user.last_name }}
                      <span v-if="user.id === currentUserId" class="badge-self">You</span>
                    </strong>
                    <span v-if="user.nin" class="user-nin-badge">
                      <i class="icofont-id-card"></i> {{ user.nin }}
                    </span>
                    <span class="user-email"><i class="icofont-ui-email"></i> {{ user.email }}</span>
                    <span v-if="user.phone" class="user-phone"><i class="icofont-ui-touch-phone"></i> {{ user.phone }}</span>
                  </div>
                </div>
              </td>

              <!-- Status -->
              <td>
                <span class="status-badge" :class="user.is_active ? 'status-active' : 'status-suspended'">
                  <span class="status-dot"></span>
                  {{ user.account_status || (user.is_active ? 'ACTIVE' : 'SUSPENDED') }}
                </span>
              </td>

              <!-- Roles -->
              <td>
                <div class="role-badges">
                  <span
                    v-for="role in user.roles || []"
                    :key="role.id || role.name"
                    class="role-pill role-pill-revokable"
                    :class="getRolePillClass(role.name || role)"
                  >
                    {{ role.name || role }}
                    <button
                      type="button"
                      class="revoke-inline-btn"
                      title="Revoke this role"
                      :disabled="savingRole === role.id"
                      @click.stop="quickRevokeRole(user, role)"
                    >
                      <i v-if="savingRole === role.id" class="icofont-spinner icofont-spin"></i>
                      <i v-else class="icofont-close-line"></i>
                    </button>
                  </span>
                  <span v-if="!(user.roles || []).length" class="no-roles-badge">No roles</span>
                  <button type="button" class="btn-icon-link" title="Manage roles" @click="handleAssignRoles(user)">
                    <i class="icofont-plus-circle"></i>
                  </button>
                </div>
              </td>

              <!-- Date -->
              <td class="date-cell">
                <span>{{ formatDate(user.created_at) }}</span>
              </td>

              <!-- Action Toolbar -->
              <td class="actions-cell text-right">
                <div class="action-buttons">
                  <button
                    type="button"
                    class="action-btn btn-edit"
                    title="Edit user details"
                    @click="handleEditUser(user)"
                  >
                    <i class="icofont-ui-edit"></i> Edit
                  </button>

                  <button
                    type="button"
                    class="action-btn"
                    :class="user.is_active ? 'btn-suspend' : 'btn-activate'"
                    :title="user.is_active ? 'Suspend account' : 'Reactivate account'"
                    @click="toggleUserSuspension(user)"
                  >
                    <i :class="user.is_active ? 'icofont-ban' : 'icofont-check'"></i>
                    {{ user.is_active ? 'Suspend' : 'Activate' }}
                  </button>

                  <button
                    type="button"
                    class="action-btn btn-reset"
                    title="Reset user password"
                    @click="openResetPasswordModal(user)"
                  >
                    <i class="icofont-key"></i> Reset PW
                  </button>

                  <button
                    type="button"
                    class="action-btn btn-roles"
                    title="Manage roles on dedicated page"
                    @click="handleAssignRoles(user)"
                  >
                    <i class="icofont-shield"></i> Roles
                  </button>

                  <button
                    type="button"
                    class="action-btn btn-delete"
                    title="Delete user account"
                    @click="deleteUser(user)"
                  >
                    <i class="icofont-ui-delete"></i>
                  </button>
                </div>
              </td>
            </tr>

            <!-- Empty State -->
            <tr v-if="!loading && !filteredUsers.length">
              <td colspan="5" class="empty-state-cell">
                <div class="empty-state">
                  <i class="icofont-search-folder empty-state-icon"></i>
                  <h3>No users found</h3>
                  <p>{{ searchQuery || statusFilter !== 'ALL' || roleFilter !== 'ALL' ? 'No user matches the current filters.' : 'There are no user accounts in the database.' }}</p>
                  <button v-if="searchQuery || statusFilter !== 'ALL' || roleFilter !== 'ALL'" type="button" class="btn btn-outline" @click="resetFilters">
                    Clear Filters
                  </button>
                </div>
              </td>
            </tr>

            <!-- Loading State -->
            <tr v-if="loading">
              <td colspan="5" class="loading-cell">
                <div class="loading-spinner">
                  <i class="icofont-spinner icofont-spin"></i> Loading users from database...
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- =====================================================================
         MODAL: RESET PASSWORD MODAL (Quick Action)
         ===================================================================== -->
    <div v-if="showResetModal" class="modal-backdrop" @click.self="showResetModal = false">
      <div class="modal-dialog">
        <header class="modal-header">
          <h3><i class="icofont-key"></i> Reset Password</h3>
          <button type="button" class="modal-close" @click="showResetModal = false">&times;</button>
        </header>
        <div class="modal-content-body">
          <p class="reset-target-info">
            Setting a new password for: <strong>{{ activeUser?.first_name }} {{ activeUser?.last_name }}</strong> ({{ activeUser?.email }}).
          </p>
          <p class="reset-warning">
            <i class="icofont-info-circle"></i> Resetting the password will automatically revoke all existing sessions for this user.
          </p>

          <form class="modal-form" @submit.prevent="submitResetPassword">
            <div class="form-group">
              <label>New Password <strong>*</strong> <small>(minimum 12 characters)</small></label>
              <div class="password-input-wrap">
                <input
                  v-model="resetPasswordValue"
                  :type="showResetPasswordText ? 'text' : 'password'"
                  required
                  minlength="12"
                  placeholder="Enter new password (min 12 chars)"
                />
                <button type="button" class="toggle-pwd-btn" @click="showResetPasswordText = !showResetPasswordText">
                  <i :class="showResetPasswordText ? 'icofont-eye-blocked' : 'icofont-eye'"></i>
                </button>
                <button type="button" class="btn-generate" @click="generateResetPassword">
                  Generate
                </button>
              </div>
            </div>

            <footer class="modal-footer">
              <button type="button" class="btn btn-outline" @click="showResetModal = false">Cancel</button>
              <button type="submit" class="btn btn-warning" :disabled="saving || resetPasswordValue.length < 12">
                <i class="icofont-key"></i> {{ saving ? 'Resetting...' : 'Confirm Reset Password' }}
              </button>
            </footer>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import * as cms from '@/api/cms.js'
import apiClient from '@/api/client.js'

const emit = defineEmits(['create-user', 'edit-user', 'assign-roles', 'message', 'error', 'navigate'])
const router = useRouter()

function handleAssignRoles(user) {
  emit('assign-roles', user)
}

function handleCreateUser() {
  emit('create-user')
}

function handleEditUser(user) {
  emit('edit-user', user)
}

const users = ref([])
const roles = ref([])
const filteredUsers = ref([])
const loading = ref(false)
const saving = ref(false)
const savingRole = ref('')
const localError = ref('')
const localMessage = ref('')

const searchQuery = ref('')
const statusFilter = ref('ALL')
const roleFilter = ref('ALL')

const activeUser = ref(null)
const showResetModal = ref(false)
const showRolesModal = ref(false)

const showResetPasswordText = ref(false)
const resetPasswordValue = ref('')

const currentUserId = computed(() => {
  try {
    const u = JSON.parse(localStorage.getItem('ncsms_user') || '{}')
    return u.id || ''
  } catch {
    return ''
  }
})

const activeUsersCount = computed(() => users.value.filter(u => u.is_active).length)
const suspendedUsersCount = computed(() => users.value.filter(u => !u.is_active).length)
const adminUsersCount = computed(() => users.value.filter(u => {
  const roleNames = (u.roles || []).map(r => r.name || r)
  return roleNames.includes('super_admin') || roleNames.includes('admin') || roleNames.includes('portal_operator')
}).length)

onMounted(async () => {
  await Promise.all([loadUsers(), loadRoles()])
})

async function loadRoles() {
  try {
    const res = await cms.adminListRoles()
    roles.value = res.data?.data || res.data || []
  } catch (e) {
    console.error('Could not load roles:', e)
  }
}

async function loadUsers() {
  loading.value = true
  localError.value = ''
  try {
    const res = await cms.adminListUsers({ page: 1, per_page: 200 })
    const list = res.data?.data || res.data?.items || res.data || []
    users.value = Array.isArray(list) ? list : []
    filterUsers()
  } catch (err) {
    localError.value = err.response?.data?.error?.message || err.message || 'Could not load users from database.'
    emit('error', localError.value)
  } finally {
    loading.value = false
  }
}

function handleSearchInput() {
  filterUsers()
}

function filterUsers() {
  const query = (searchQuery.value || '').trim().toLowerCase()
  const status = statusFilter.value
  const role = roleFilter.value

  filteredUsers.value = users.value.filter(user => {
    // 1. Text search
    if (query) {
      const name = `${user.first_name || ''} ${user.last_name || ''}`.toLowerCase()
      const email = (user.email || '').toLowerCase()
      const phone = (user.phone || '').toLowerCase()
      const nin = (user.nin || '').toLowerCase()
      if (!name.includes(query) && !email.includes(query) && !phone.includes(query) && !nin.includes(query)) {
        return false
      }
    }

    // 2. Status filter
    if (status === 'ACTIVE' && !user.is_active) return false
    if (status === 'SUSPENDED' && user.is_active) return false

    // 3. Role filter
    if (role !== 'ALL') {
      const userRoles = (user.roles || []).map(r => r.name || r)
      if (!userRoles.includes(role)) return false
    }

    return true
  })
}

function resetFilters() {
  searchQuery.value = ''
  statusFilter.value = 'ALL'
  roleFilter.value = 'ALL'
  filterUsers()
}

// ── SUSPEND / REACTIVATE ───────────────────────────────────────────────────

async function toggleUserSuspension(user) {
  const isActivating = !user.is_active
  const actionText = isActivating ? 'Activate' : 'Suspend'
  const confirmResult = await Swal.fire({
    title: `${actionText} ${user.first_name || user.email}?`,
    text: isActivating
      ? 'This user will regain access to the portal.'
      : 'This user will be immediately suspended and all active sessions will be revoked.',
    icon: isActivating ? 'question' : 'warning',
    showCancelButton: true,
    confirmButtonText: `Yes, ${actionText}`,
    cancelButtonText: 'Cancel',
    confirmButtonColor: isActivating ? '#28a745' : '#fc544b',
    cancelButtonColor: '#6c757d',
  })
  if (!confirmResult.isConfirmed) return

  loading.value = true
  try {
    if (isActivating) {
      await cms.adminActivateUser(user.id)
      notifySuccess(`User ${user.email} has been activated.`)
    } else {
      await cms.adminDeactivateUser(user.id)
      notifySuccess(`User ${user.email} has been suspended.`)
    }
    await loadUsers()
  } catch (err) {
    localError.value = err.response?.data?.error?.message || err.message || `Failed to ${actionText.toLowerCase()} user.`
    emit('error', localError.value)
  } finally {
    loading.value = false
  }
}

// ── RESET PASSWORD ─────────────────────────────────────────────────────────

function openResetPasswordModal(user) {
  activeUser.value = user
  resetPasswordValue.value = ''
  showResetPasswordText.value = false
  showResetModal.value = true
}

function generateResetPassword() {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789!@#$%^&*'
  let pwd = ''
  for (let i = 0; i < 16; i++) {
    pwd += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  resetPasswordValue.value = pwd
  showResetPasswordText.value = true
}

async function submitResetPassword() {
  if (resetPasswordValue.value.length < 12) {
    localError.value = 'Password must be at least 12 characters.'
    return
  }
  saving.value = true
  localError.value = ''
  try {
    await cms.adminResetUserPassword(activeUser.value.id, resetPasswordValue.value)
    showResetModal.value = false
    notifySuccess(`Password for ${activeUser.value.email} was reset successfully.`)
  } catch (err) {
    localError.value = err.response?.data?.error?.message || err.message || 'Failed to reset password.'
    emit('error', localError.value)
  } finally {
    saving.value = false
  }
}

// ── MANAGE ROLES ───────────────────────────────────────────────────────────

function openRolesModal(user) {
  activeUser.value = user
  showRolesModal.value = true
  if (!roles.value.length) loadRoles()
}

function userHasRole(user, role) {
  if (!user) return false
  return (user.roles || []).some(item => (item.id && item.id === role.id) || (item.name || item) === role.name)
}

// Resolve a role object (which may be a partial {id, name} or just name string) to the full role from the roles list
function resolveRoleObject(roleOrName) {
  if (!roleOrName) return null
  const name = typeof roleOrName === 'object' ? (roleOrName.name || '') : roleOrName
  const id = typeof roleOrName === 'object' ? (roleOrName.id || '') : ''
  return roles.value.find(r => r.id === id || r.name === name) || roleOrName
}

async function quickRevokeRole(user, role) {
  const fullRole = resolveRoleObject(role)
  if (!fullRole) return
  await toggleRoleAssignment(user, fullRole)
}

async function toggleRoleAssignment(user, role) {
  const isAssigned = userHasRole(user, role)
  // Always use UUID as the key and for the API call
  const roleUUID = role.id
  if (!roleUUID) {
    localError.value = `Cannot resolve UUID for role '${role.name || role}'. Please refresh and try again.`
    return
  }
  savingRole.value = roleUUID
  localError.value = ''
  try {
    if (isAssigned) {
      await cms.adminRemoveUserRole(user.id, roleUUID)
      notifySuccess(`Role '${role.name || role}' revoked from ${user.email}.`)
    } else {
      await cms.adminAssignUserRole(user.id, roleUUID)
      notifySuccess(`Role '${role.name || role}' assigned to ${user.email}.`)
    }
    // Refresh user in list
    const res = await cms.adminGetUser(user.id)
    const updated = res.data?.data || res.data
    const idx = users.value.findIndex(u => u.id === user.id)
    if (idx !== -1 && updated) {
      users.value[idx] = updated
      activeUser.value = updated
      filterUsers()
    } else {
      await loadUsers()
    }
  } catch (err) {
    localError.value = err.response?.data?.error?.message || err.message || 'Failed to update role.'
    emit('error', localError.value)
  } finally {
    savingRole.value = ''
  }
}

async function reassignUserRole(user, role) {
  const roleUUID = role.id
  if (!roleUUID) {
    localError.value = `Cannot resolve UUID for role '${role.name || role}'.`
    return
  }
  savingRole.value = roleUUID
  localError.value = ''
  try {
    const currentRoles = [...(user.roles || [])]
    for (const r of currentRoles) {
      const curId = typeof r === 'object' ? (r.id || '') : ''
      if (curId && curId !== roleUUID) {
        try { await cms.adminRemoveUserRole(user.id, curId) } catch (e) { console.warn('Could not revoke previous role:', e) }
      }
    }
    await cms.adminAssignUserRole(user.id, roleUUID)
    notifySuccess(`Role for ${user.email} reassigned to '${role.name || role}'.`)

    const res = await cms.adminGetUser(user.id)
    const updated = res.data?.data || res.data
    const idx = users.value.findIndex(u => u.id === user.id)
    if (idx !== -1 && updated) {
      users.value[idx] = updated
      activeUser.value = updated
      filterUsers()
    } else {
      await loadUsers()
    }
  } catch (err) {
    localError.value = err.response?.data?.error?.message || err.message || 'Failed to reassign role.'
    emit('error', localError.value)
  } finally {
    savingRole.value = ''
  }
}

// ── DELETE USER ────────────────────────────────────────────────────────────

async function deleteUser(user) {
  const confirmResult = await Swal.fire({
    title: `Delete user ${user.email}?`,
    text: 'This account will be permanently disabled and marked deleted in the database.',
    icon: 'error',
    showCancelButton: true,
    confirmButtonText: 'Yes, Delete User',
    cancelButtonText: 'Cancel',
    confirmButtonColor: '#fc544b',
    cancelButtonColor: '#6c757d',
    reverseButtons: true,
  })
  if (!confirmResult.isConfirmed) return

  loading.value = true
  try {
    await cms.adminDeleteUser(user.id)
    notifySuccess(`User ${user.email} has been deleted.`)
    await loadUsers()
  } catch (err) {
    localError.value = err.response?.data?.error?.message || err.message || 'Failed to delete user.'
    emit('error', localError.value)
  } finally {
    loading.value = false
  }
}

// ── UTILITIES ──────────────────────────────────────────────────────────────

function notifySuccess(msg) {
  localMessage.value = msg
  localError.value = ''
  emit('message', msg)
  setTimeout(() => { localMessage.value = '' }, 4000)
}

function getInitials(user) {
  const first = (user.first_name || '').charAt(0).toUpperCase()
  const last = (user.last_name || '').charAt(0).toUpperCase()
  return `${first}${last}` || 'U'
}

function getAvatarColorClass(user) {
  const colors = ['avatar-purple', 'avatar-blue', 'avatar-emerald', 'avatar-amber', 'avatar-rose']
  const str = user.email || user.id || ''
  let hash = 0
  for (let i = 0; i < str.length; i++) hash += str.charCodeAt(i)
  return colors[hash % colors.length]
}

function getRolePillClass(roleName) {
  const r = String(roleName || '').toLowerCase()
  if (r.includes('super_admin')) return 'role-superadmin'
  if (r.includes('admin')) return 'role-admin'
  if (r.includes('operator')) return 'role-operator'
  if (r.includes('federation')) return 'role-federation'
  return 'role-user'
}

function formatDate(isoStr) {
  if (!isoStr) return '-'
  try {
    const d = new Date(isoStr)
    return d.toLocaleDateString('en-UG', { year: 'numeric', month: 'short', day: 'numeric' })
  } catch {
    return isoStr
  }
}
</script>

<style scoped>
.manage-users-panel {
  display: flex;
  flex-direction: column;
  gap: 22px;
  width: 100%;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  padding-bottom: 18px;
  border-bottom: 1px solid #e2e8f0;
}

.panel-subtitle {
  margin: 0 0 4px;
  font-size: 11px;
  font-weight: 800;
  color: #6777ef;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.panel-title {
  margin: 0 0 6px;
  font-size: 26px;
  font-weight: 800;
  color: #1e293b;
}

.panel-desc {
  color: #64748b;
  font-size: 13px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

/* KPI Cards */
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.kpi-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 20px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.kpi-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.06);
}

.kpi-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  flex-shrink: 0;
}

.icon-users { background: #eff6ff; color: #2563eb; }
.icon-active { background: #ecfdf5; color: #059669; }
.icon-suspended { background: #fef2f2; color: #dc2626; }
.icon-admin { background: #f5f3ff; color: #7c3aed; }

.kpi-info {
  display: flex;
  flex-direction: column;
}

.kpi-label {
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
}

.kpi-value {
  font-size: 22px;
  font-weight: 800;
  color: #1e293b;
}

/* Filter Toolbar */
.filter-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.search-box {
  position: relative;
  flex: 1;
  min-width: 280px;
  display: flex;
  align-items: center;
}

.search-box i {
  position: absolute;
  left: 14px;
  color: #94a3b8;
  font-size: 16px;
  pointer-events: none;
}

.search-box input {
  width: 100%;
  height: 42px;
  padding: 8px 36px 8px 38px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-family: inherit;
  font-size: 13px;
  background: #ffffff;
  color: #1e293b;
  transition: all 0.2s ease;
}

.search-box input:focus {
  outline: none;
  border-color: #6777ef;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.15);
}

.clear-search-btn {
  position: absolute;
  right: 10px;
  border: none;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;
  font-size: 16px;
}

.filter-dropdowns {
  display: flex;
  align-items: center;
  gap: 10px;
}

.filter-select {
  height: 42px;
  padding: 0 14px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-family: inherit;
  font-size: 13px;
  background: #ffffff;
  color: #1e293b;
  cursor: pointer;
}

.filter-select:focus {
  outline: none;
  border-color: #6777ef;
}

/* Table Card */
.table-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.table-responsive {
  width: 100%;
  overflow-x: auto;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.users-table th {
  padding: 14px 18px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  font-size: 12px;
  font-weight: 700;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.users-table td {
  padding: 14px 18px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 13px;
  color: #334155;
  vertical-align: middle;
}

.user-row:hover {
  background: #f8fafc;
}

/* User identity cell */
.user-identity {
  display: flex;
  align-items: center;
  gap: 14px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 13px;
  color: #ffffff;
  flex-shrink: 0;
  letter-spacing: 0.05em;
}

.avatar-purple { background: linear-gradient(135deg, #7c3aed, #a855f7); }
.avatar-blue { background: linear-gradient(135deg, #2563eb, #38bdf8); }
.avatar-emerald { background: linear-gradient(135deg, #059669, #34d399); }
.avatar-amber { background: linear-gradient(135deg, #d97706, #fbbf24); }
.avatar-rose { background: linear-gradient(135deg, #e11d48, #fb7185); }

.user-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.user-full-name {
  font-size: 14px;
  font-weight: 700;
  color: #1e293b;
  display: flex;
  align-items: center;
  gap: 6px;
}

.badge-self {
  font-size: 10px;
  padding: 1px 6px;
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
  border-radius: 4px;
}

.user-email, .user-phone {
  font-size: 12px;
  color: #64748b;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

/* Status Badge */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.04em;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.status-active {
  background: #ecfdf5;
  color: #059669;
  border: 1px solid #a7f3d0;
}
.status-active .status-dot { background: #10b981; }

.status-suspended {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
}
.status-suspended .status-dot { background: #ef4444; }

/* Role Badges */
.role-badges {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
}

.role-pill {
  padding: 3px 9px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.02em;
}

.role-superadmin { background: #f5f3ff; color: #7c3aed; border: 1px solid #ddd6fe; }
.role-admin { background: #eff6ff; color: #2563eb; border: 1px solid #bfdbfe; }
.role-operator { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }
.role-federation { background: #fff7ed; color: #ea580c; border: 1px solid #fed7aa; }
.role-user { background: #f1f5f9; color: #475569; border: 1px solid #e2e8f0; }

.no-roles-badge {
  font-size: 11px;
  color: #94a3b8;
  font-style: italic;
}

.btn-icon-link {
  border: none;
  background: transparent;
  color: #6777ef;
  cursor: pointer;
  font-size: 14px;
  padding: 0;
  display: flex;
  align-items: center;
}

/* Date cell */
.date-cell {
  font-size: 12px;
  color: #64748b;
  white-space: nowrap;
}

/* Action Toolbar */
.actions-cell {
  white-space: nowrap;
}

.action-buttons {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 5px 9px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  border: 1px solid transparent;
}

.btn-edit {
  background: #eff6ff;
  color: #2563eb;
  border-color: #bfdbfe;
}
.btn-edit:hover {
  background: #2563eb;
  color: #ffffff;
}

.btn-suspend {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
}
.btn-suspend:hover {
  background: #dc2626;
  color: #ffffff;
}

.btn-activate {
  background: #ecfdf5;
  color: #059669;
  border-color: #a7f3d0;
}
.btn-activate:hover {
  background: #059669;
  color: #ffffff;
}

.btn-reset {
  background: #fffbeb;
  color: #d97706;
  border-color: #fde68a;
}
.btn-reset:hover {
  background: #d97706;
  color: #ffffff;
}

.btn-roles {
  background: #f5f3ff;
  color: #7c3aed;
  border-color: #ddd6fe;
}
.btn-roles:hover {
  background: #7c3aed;
  color: #ffffff;
}

.btn-delete {
  background: #f8fafc;
  color: #ef4444;
  border-color: #e2e8f0;
  padding: 5px 7px;
}
.btn-delete:hover {
  background: #ef4444;
  color: #ffffff;
}

/* Empty & Loading States */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.empty-state-icon {
  font-size: 38px;
  color: #94a3b8;
  margin-bottom: 10px;
}

.empty-state h3 {
  margin: 0 0 6px;
  font-size: 16px;
  font-weight: 700;
  color: #334155;
}

.empty-state p {
  margin: 0 0 14px;
  font-size: 13px;
  color: #64748b;
}

.loading-cell {
  text-align: center;
  padding: 40px 20px;
  color: #64748b;
  font-size: 14px;
  font-weight: 600;
}

/* Alerts */
.alert {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
}

.alert-error { background: #fff1f0; color: #dc2626; border: 1px solid #fecaca; }
.alert-success { background: #ecfdf5; color: #059669; border: 1px solid #a7f3d0; }
.alert-close { border: none; background: transparent; font-size: 18px; cursor: pointer; color: inherit; }

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 9px 16px;
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

.btn-outline { background: #ffffff; color: #475569; border-color: #cbd5e1; }
.btn-outline:hover:not(:disabled) { background: #f8fafc; border-color: #94a3b8; }

.btn-warning { background: #f59e0b; color: #ffffff; }
.btn-warning:hover:not(:disabled) { background: #d97706; }

.btn:disabled { opacity: 0.6; cursor: not-allowed; }

/* Modals */
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(4px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  box-sizing: border-box;
}

.modal-dialog {
  position: relative;
  z-index: 10000;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  width: 100%;
  max-width: 520px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.2);
  animation: modalFadeIn 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

/* Roles modal uses flex-column layout so header/footer are sticky */
.roles-modal-dialog {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.roles-modal-dialog .modal-header {
  flex-shrink: 0;
}

.roles-modal-dialog .roles-selection-list {
  flex: 1 1 auto;
  overflow-y: auto;
  min-height: 0;
  padding: 12px 20px;
}

.roles-modal-dialog .roles-modal-footer {
  flex-shrink: 0;
}

.modal-dialog-lg {
  max-width: 660px;
}

@keyframes modalFadeIn {
  from { opacity: 0; transform: scale(0.95) translateY(-10px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 800;
  color: #1e293b;
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-subhead {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
  display: block;
}

.modal-close {
  border: none;
  background: transparent;
  font-size: 24px;
  color: #94a3b8;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.modal-close:hover {
  color: #1e293b;
}

.modal-form, .modal-content-body {
  padding: 20px 24px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 16px;
}

.form-group label {
  font-size: 12px;
  font-weight: 700;
  color: #334155;
}

.form-group label strong {
  color: #ef4444;
}

.form-group input, .form-group select {
  width: 100%;
  height: 42px;
  padding: 0 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-family: inherit;
  font-size: 13px;
  background: #ffffff;
  color: #1e293b;
  box-sizing: border-box;
}

.form-group input:focus, .form-group select:focus {
  outline: none;
  border-color: #6777ef;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.15);
}

.input-disabled {
  background: #f1f5f9 !important;
  color: #94a3b8 !important;
  cursor: not-allowed;
}

.field-hint {
  font-size: 11px;
  color: #94a3b8;
}

.password-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
}

.password-input-wrap input {
  padding-right: 36px;
}

.toggle-pwd-btn {
  position: absolute;
  right: 90px;
  border: none;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;
  font-size: 16px;
}

.btn-generate {
  height: 42px;
  padding: 0 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  background: #f8fafc;
  color: #475569;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  white-space: nowrap;
}

.btn-generate:hover {
  background: #e2e8f0;
}

.reset-target-info {
  margin: 0 0 8px;
  font-size: 14px;
  color: #1e293b;
}

.reset-warning {
  margin: 0 0 18px;
  padding: 10px 14px;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 8px;
  color: #1e40af;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
  margin-top: 8px;
}

/* Roles Modal List */
.roles-selection-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-height: 55vh;
  overflow-y: auto !important;
  padding: 14px 20px;
  scrollbar-width: thin;
  scrollbar-color: #6777ef #e2e8f0;
}

.roles-selection-list::-webkit-scrollbar {
  width: 8px;
}

.roles-selection-list::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 4px;
}

.roles-selection-list::-webkit-scrollbar-thumb {
  background: #6777ef;
  border-radius: 4px;
}

.roles-selection-list::-webkit-scrollbar-thumb:hover {
  background: #4f46e5;
}

:global(html.dark) .roles-selection-list,
:global(body.dark) .roles-selection-list,
:global([data-theme="dark"]) .roles-selection-list {
  scrollbar-color: #6366f1 #1e293b;
}

:global(html.dark) .roles-selection-list::-webkit-scrollbar-track,
:global(body.dark) .roles-selection-list::-webkit-scrollbar-track,
:global([data-theme="dark"]) .roles-selection-list::-webkit-scrollbar-track {
  background: #0f172a !important;
}

:global(html.dark) .roles-selection-list::-webkit-scrollbar-thumb,
:global(body.dark) .roles-selection-list::-webkit-scrollbar-thumb,
:global([data-theme="dark"]) .roles-selection-list::-webkit-scrollbar-thumb {
  background: #6366f1 !important;
}

.role-option-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 14px 18px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  transition: all 0.15s ease;
}

.role-option-card.is-assigned {
  background: #f5f3ff;
  border-color: #c4b5fd;
}

.role-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 3px;
}

.role-name {
  font-size: 14px;
  color: #1e293b;
}

.system-pill {
  font-size: 10px;
  padding: 1px 6px;
  background: #e2e8f0;
  color: #475569;
  border-radius: 4px;
  font-weight: 700;
}

.role-desc {
  margin: 0;
  font-size: 12px;
  color: #64748b;
}

.role-toggle-btn {
  padding: 6px 14px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  border: 1px solid transparent;
}

.btn-add-role {
  background: #eff6ff;
  color: #2563eb;
  border-color: #bfdbfe;
}
.btn-add-role:hover {
  background: #2563eb;
  color: #ffffff;
}

.btn-remove-role {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
}
.btn-remove-role:hover {
  background: #dc2626;
  color: #ffffff;
}

/* Assigned roles summary bar inside modal */
.assigned-roles-summary {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
  flex-shrink: 0;
  font-size: 13px;
}

.assigned-empty {
  color: #94a3b8;
  font-style: italic;
}

.assigned-label {
  font-weight: 700;
  color: #059669;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.assigned-role-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px 3px 10px;
  background: #ede9fe;
  border: 1px solid #c4b5fd;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 700;
  color: #5b21b6;
}

.chip-revoke-btn {
  border: none;
  background: transparent;
  color: #7c3aed;
  cursor: pointer;
  padding: 0 2px;
  font-size: 11px;
  display: flex;
  align-items: center;
  transition: color 0.15s;
}

.chip-revoke-btn:hover:not(:disabled) {
  color: #dc2626;
}

/* Inline revoke button inside role pills in the table */
.role-pill-revokable {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding-right: 5px;
}

.revoke-inline-btn {
  border: none;
  background: transparent;
  cursor: pointer;
  color: inherit;
  opacity: 0.6;
  padding: 0 1px;
  font-size: 10px;
  display: flex;
  align-items: center;
  transition: opacity 0.15s;
  line-height: 1;
}

.revoke-inline-btn:hover:not(:disabled) {
  opacity: 1;
  color: #dc2626;
}

.roles-modal-footer {
  padding: 14px 20px;
  border-top: 1px solid #e2e8f0;
  background: #ffffff;
  display: flex;
  justify-content: flex-end;
}

/* Dark mode for new elements */
:global(html.dark) .assigned-roles-summary,
:global(body.dark) .assigned-roles-summary,
:global([data-theme="dark"]) .assigned-roles-summary {
  background: #0f172a !important;
  border-bottom-color: #334155 !important;
}

:global(html.dark) .assigned-role-chip,
:global(body.dark) .assigned-role-chip,
:global([data-theme="dark"]) .assigned-role-chip {
  background: #2e1065 !important;
  border-color: #4c1d95 !important;
  color: #c4b5fd !important;
}

:global(html.dark) .roles-modal-footer,
:global(body.dark) .roles-modal-footer,
:global([data-theme="dark"]) .roles-modal-footer {
  background: #1e293b !important;
  border-top-color: #334155 !important;
}

/* Dark mode for modal-dialog */
:global(html.dark) .roles-modal-dialog .modal-header,
:global(body.dark) .roles-modal-dialog .modal-header,
:global([data-theme="dark"]) .roles-modal-dialog .modal-header {
  background: #1e293b !important;
}

/* =====================================================================
   DARK MODE STYLES
   ===================================================================== */
:global(html.dark) .manage-users-panel,
:global(body.dark) .manage-users-panel,
:global([data-theme="dark"]) .manage-users-panel {
  color: #e2e8f0 !important;
}

:global(html.dark) .panel-title,
:global(body.dark) .panel-title,
:global([data-theme="dark"]) .panel-title {
  color: #f8fafc !important;
}

:global(html.dark) .panel-header,
:global(body.dark) .panel-header,
:global([data-theme="dark"]) .panel-header {
  border-bottom-color: #334155 !important;
}

:global(html.dark) .kpi-card,
:global(body.dark) .kpi-card,
:global([data-theme="dark"]) .kpi-card {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(html.dark) .kpi-value,
:global(body.dark) .kpi-value,
:global([data-theme="dark"]) .kpi-value {
  color: #f8fafc !important;
}

:global(html.dark) .search-box input,
:global(html.dark) .filter-select,
:global(body.dark) .search-box input,
:global(body.dark) .filter-select,
:global([data-theme="dark"]) .search-box input,
:global([data-theme="dark"]) .filter-select {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(html.dark) .table-card,
:global(body.dark) .table-card,
:global([data-theme="dark"]) .table-card {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(html.dark) .users-table th,
:global(body.dark) .users-table th,
:global([data-theme="dark"]) .users-table th {
  background-color: #0f172a !important;
  border-bottom-color: #334155 !important;
  color: #94a3b8 !important;
}

:global(html.dark) .users-table td,
:global(body.dark) .users-table td,
:global([data-theme="dark"]) .users-table td {
  border-bottom-color: #334155 !important;
  color: #cbd5e1 !important;
}

:global(html.dark) .user-row:hover,
:global(body.dark) .user-row:hover,
:global([data-theme="dark"]) .user-row:hover {
  background-color: #0f172a !important;
}

:global(html.dark) .user-full-name,
:global(body.dark) .user-full-name,
:global([data-theme="dark"]) .user-full-name {
  color: #f8fafc !important;
}

:global(html.dark) .modal-dialog,
:global(body.dark) .modal-dialog,
:global([data-theme="dark"]) .modal-dialog {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(html.dark) .modal-header,
:global(html.dark) .modal-footer,
:global(body.dark) .modal-header,
:global(body.dark) .modal-footer,
:global([data-theme="dark"]) .modal-header,
:global([data-theme="dark"]) .modal-footer {
  border-color: #334155 !important;
}

:global(html.dark) .modal-header h3,
:global(body.dark) .modal-header h3,
:global([data-theme="dark"]) .modal-header h3 {
  color: #f8fafc !important;
}

:global(html.dark) .form-group label,
:global(body.dark) .form-group label,
:global([data-theme="dark"]) .form-group label {
  color: #cbd5e1 !important;
}

:global(html.dark) .form-group input:not(.input-disabled),
:global(html.dark) .form-group select,
:global(body.dark) .form-group input:not(.input-disabled),
:global(body.dark) .form-group select,
:global([data-theme="dark"]) .form-group input:not(.input-disabled),
:global([data-theme="dark"]) .form-group select {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(html.dark) .role-option-card,
:global(body.dark) .role-option-card,
:global([data-theme="dark"]) .role-option-card {
  background-color: #0f172a !important;
  border-color: #334155 !important;
}

:global(html.dark) .role-name,
:global(body.dark) .role-name,
:global([data-theme="dark"]) .role-name {
  color: #f8fafc !important;
}

.role-btn-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-reassign-role {
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fde68a;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: all 0.15s ease;
}

.btn-reassign-role:hover:not(:disabled) {
  background: #f59e0b;
  color: #ffffff;
}

:global(html.dark) .btn-outline,
:global(body.dark) .btn-outline,
:global([data-theme="dark"]) .btn-outline {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #cbd5e1 !important;
}

:global(html.dark) .btn-generate,
:global(body.dark) .btn-generate,
:global([data-theme="dark"]) .btn-generate {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #cbd5e1 !important;
}

:global(html.dark) .reset-target-info,
:global(body.dark) .reset-target-info,
:global([data-theme="dark"]) .reset-target-info {
  color: #f8fafc !important;
}

@media (max-width: 900px) {
  .kpi-grid { grid-template-columns: repeat(2, 1fr); }
  .panel-header { flex-direction: column; }
  .header-actions { width: 100%; justify-content: flex-start; }
}

@media (max-width: 600px) {
  .kpi-grid { grid-template-columns: 1fr; }
  .form-row { grid-template-columns: 1fr; }
  .filter-toolbar { flex-direction: column; align-items: stretch; }
  .filter-dropdowns { flex-direction: column; }
}
</style>

<style scoped>
.user-nin-badge {
  font-size: 11px;
  font-weight: 600;
  font-family: monospace;
  padding: 1px 6px;
  border-radius: 4px;
  background: rgba(16, 185, 129, 0.12);
  color: #059669;
  border: 1px solid rgba(16, 185, 129, 0.25);
  display: inline-flex;
  align-items: center;
  gap: 3px;
  margin: 2px 0;
  width: fit-content;
}
:global(body.dark-theme) .user-nin-badge {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
  border-color: rgba(16, 185, 129, 0.35);
}
</style>
