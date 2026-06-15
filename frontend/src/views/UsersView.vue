<template>
  <LayoutDefault title="User Management">
    <div class="space-y-4">

      <!-- Alert -->
      <div
        v-if="alertMsg"
        :class="alertType === 'success' ? 'bg-green-50 border-green-200 text-green-700' : 'bg-red-50 border-red-200 text-red-700'"
        class="flex items-center gap-2 p-3 border rounded-lg text-sm"
      >
        <i :class="alertType === 'success' ? 'icofont-check-circled' : 'icofont-warning-alt'" class="text-lg flex-shrink-0"></i>
        {{ alertMsg }}
      </div>

      <AdminTable
        title="Users"
        :loading="loading"
        :is-empty="!loading && users.length === 0"
        :meta="meta"
        :total="meta?.total ?? null"
        create-label="Create User"
        empty-icon="icofont-people"
        empty-message="No users found"
        :searchable="true"
        v-model:search-value="searchQuery"
        search-placeholder="Search users…"
        :colspan="5"
        @create="openCreateModal"
        @page-change="changePage"
      >
        <template #head>
          <th class="table-th">User</th>
          <th class="table-th">Status</th>
          <th class="table-th">Roles</th>
          <th class="table-th">Last Login</th>
          <th class="table-th text-right">Actions</th>
        </template>

        <template #rows>
          <tr v-for="user in users" :key="user.id" class="hover:bg-gray-50 transition-colors">
            <!-- User -->
            <td class="table-td">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 bg-primary-100 text-primary-700 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0">
                  {{ getInitials(user) }}
                </div>
                <div class="min-w-0">
                  <div class="text-sm font-semibold text-gray-900 truncate">{{ getFullName(user) }}</div>
                  <div class="text-xs text-gray-400 truncate">{{ user.email }}</div>
                </div>
              </div>
            </td>
            <!-- Status -->
            <td class="table-td">
              <div class="flex flex-col items-start gap-1">
                <span :class="statusBadgeClass(user)">
                  <span class="w-1.5 h-1.5 rounded-full mr-1" :class="statusDotClass(user)"></span>
                  {{ formatAccountStatus(user) }}
                </span>
                <span v-if="user.fraud_flag" class="badge-red">Fraud flagged</span>
                <span v-if="user.suspended_until" class="text-[11px] text-gray-500">Until {{ formatDateTime(user.suspended_until) }}</span>
              </div>
            </td>
            <!-- Roles -->
            <td class="table-td">
              <div class="flex flex-wrap gap-1">
                <span v-for="role in (user.roles || [])" :key="getRoleName(role)" class="badge-blue">
                  {{ formatRoleName(getRoleName(role)) }}
                </span>
                <span v-if="!user.roles?.length" class="badge-gray">No roles</span>
              </div>
            </td>
            <!-- Last Login -->
            <td class="table-td text-gray-500">{{ formatDate(user.last_login_at) }}</td>
            <!-- Actions -->
            <td class="table-td">
              <select
                aria-label="User actions"
                class="form-input min-w-[9rem] py-1.5 text-xs"
                @change="handleActionSelection($event, user)"
              >
                <option value="">Actions...</option>
                <option value="edit">Edit user</option>
                <option v-if="canManageAccount(user)" value="reset">Reset password</option>
                <option v-if="canManageAccount(user) && accountStatus(user) === 'ACTIVE'" value="SUSPEND">Suspend</option>
                <option v-if="canManageAccount(user) && accountStatus(user) !== 'BANNED'" value="BAN">Ban</option>
                <option v-if="canManageAccount(user) && !user.fraud_flag" value="MARK_FRAUD">Mark as fraud</option>
                <option v-if="canManageAccount(user) && user.fraud_flag" value="CLEAR_FRAUD">Clear fraud flag</option>
                <option v-if="canManageAccount(user) && accountStatus(user) !== 'ACTIVE' && !user.fraud_flag" value="REACTIVATE">Reactivate</option>
                <option v-if="canManageAccount(user)" value="delete">Delete user</option>
              </select>
            </td>
          </tr>
        </template>
      </AdminTable>
    </div>

    <!-- Create User Modal -->
    <Modal :show="showCreateModal" title="Create New User" size="md" @close="showCreateModal = false">
      <form @submit.prevent="submitCreate" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="form-label">First Name</label>
            <input v-model="createForm.first_name" type="text" required class="form-input" />
          </div>
          <div>
            <label class="form-label">Last Name</label>
            <input v-model="createForm.last_name" type="text" required class="form-input" />
          </div>
        </div>
        <div>
          <label class="form-label">Email Address</label>
          <input v-model="createForm.email" type="email" required class="form-input" />
        </div>
        <div>
          <label class="form-label">Phone Number</label>
          <input v-model="createForm.phone" type="tel" class="form-input" placeholder="+256..." />
        </div>
        <div>
          <label class="form-label">Password</label>
          <input v-model="createForm.password" type="password" required class="form-input" />
        </div>
        <div>
          <label class="form-label">Role</label>
          <select v-model="createForm.role_id" class="form-input">
            <option value="">— No role —</option>
            <option v-for="r in assignableRoles" :key="r.id" :value="r.id">{{ formatRoleName(r.name) }}</option>
          </select>
        </div>
        <div v-if="modalError" class="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{{ modalError }}</div>
      </form>
      <template #footer>
        <button @click="showCreateModal = false" class="btn-ghost">Cancel</button>
        <button @click="submitCreate" :disabled="submitting" class="btn-primary disabled:opacity-60">
          {{ submitting ? 'Creating…' : 'Create User' }}
        </button>
      </template>
    </Modal>

    <!-- Edit User Modal -->
    <Modal :show="showEditModal" title="Edit User" size="md" @close="showEditModal = false">
      <form @submit.prevent="submitEdit" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="form-label">First Name</label>
            <input v-model="editForm.first_name" type="text" class="form-input" />
          </div>
          <div>
            <label class="form-label">Last Name</label>
            <input v-model="editForm.last_name" type="text" class="form-input" />
          </div>
        </div>
        <div>
          <label class="form-label">Phone Number</label>
          <input v-model="editForm.phone" type="tel" class="form-input" />
        </div>
        <div class="border-t border-gray-100 pt-4">
          <label class="form-label mb-2">Roles</label>
          <div class="flex flex-wrap gap-1.5 mb-3 min-h-[1.5rem]">
            <span v-for="role in (editTarget?.roles || [])" :key="getRoleName(role)" class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-blue-50 text-blue-700">
              {{ formatRoleName(getRoleName(role)) }}
              <button v-if="getRoleName(role) !== 'Super Admin'" type="button" @click="revokeRole(role)" class="hover:text-red-600 leading-none">×</button>
            </span>
            <span v-if="!editTarget?.roles?.length" class="badge-gray">No roles assigned</span>
          </div>
          <div class="flex gap-2">
            <select v-model="newRoleId" class="form-input flex-1">
              <option value="">— Select role to assign —</option>
              <option v-for="r in assignableRolesForEdit" :key="r.id" :value="r.id">{{ formatRoleName(r.name) }}</option>
            </select>
            <button type="button" @click="assignRole" :disabled="!newRoleId || roleActionLoading" class="btn-primary disabled:opacity-50 flex-shrink-0">
              {{ roleActionLoading ? '…' : 'Assign' }}
            </button>
          </div>
        </div>
        <div v-if="modalError" class="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{{ modalError }}</div>
      </form>
      <template #footer>
        <button @click="showEditModal = false" class="btn-ghost">Cancel</button>
        <button @click="submitEdit" :disabled="submitting" class="btn-primary disabled:opacity-60">
          {{ submitting ? 'Saving…' : 'Save Changes' }}
        </button>
      </template>
    </Modal>

	<!-- Reset Password Modal -->
	<Modal :show="showResetModal" title="Reset User Password" size="sm" @close="closeResetModal">
	  <form @submit.prevent="submitPasswordReset" class="space-y-4">
		<div class="p-3 bg-yellow-50 border border-yellow-200 rounded-lg text-sm text-yellow-800">
		  Resetting <strong>{{ resetTarget?.email }}</strong> will immediately sign the user out of every active session.
		</div>
		<div>
		  <label class="form-label">New Password</label>
		  <input v-model="resetForm.new_password" type="password" minlength="12" autocomplete="new-password" required class="form-input" />
		  <p class="text-xs text-gray-500 mt-1">Use at least 12 characters.</p>
		</div>
		<div>
		  <label class="form-label">Confirm New Password</label>
		  <input v-model="resetForm.confirm_password" type="password" minlength="12" autocomplete="new-password" required class="form-input" />
		</div>
		<div v-if="resetForm.confirm_password && resetForm.new_password !== resetForm.confirm_password" class="text-xs text-red-600">
		  Passwords do not match.
		</div>
		<div v-if="modalError" class="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{{ modalError }}</div>
	  </form>
	  <template #footer>
		<button @click="closeResetModal" class="btn-ghost">Cancel</button>
		<button
		  @click="submitPasswordReset"
		  :disabled="submitting || resetForm.new_password.length < 12 || resetForm.new_password !== resetForm.confirm_password"
		  class="btn-primary disabled:opacity-60"
		>
		  {{ submitting ? 'Resetting…' : 'Reset Password' }}
		</button>
	  </template>
	</Modal>

    <!-- Account Action Modal -->
    <Modal :show="showAccountActionModal" :title="accountActionTitle" size="sm" @close="closeAccountActionModal">
      <form @submit.prevent="submitAccountAction" class="space-y-4">
        <div class="p-3 bg-yellow-50 border border-yellow-200 rounded-lg text-sm text-yellow-800">
          {{ accountActionDescription }}
        </div>
        <div v-if="accountActionRequiresReason">
          <label class="form-label">Reason</label>
          <textarea v-model="accountActionForm.reason" rows="3" maxlength="500" required class="form-input" placeholder="Enter a clear reason for the audit record"></textarea>
        </div>
        <div v-if="accountActionForm.action === 'SUSPEND'">
          <label class="form-label">Suspended Until <span class="text-gray-400">(optional)</span></label>
          <input v-model="accountActionForm.suspended_until" type="datetime-local" class="form-input" />
        </div>
        <div v-if="modalError" class="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{{ modalError }}</div>
      </form>
      <template #footer>
        <button @click="closeAccountActionModal" class="btn-ghost">Cancel</button>
        <button @click="submitAccountAction" :disabled="submitting || (accountActionRequiresReason && !accountActionForm.reason.trim())" class="btn-primary disabled:opacity-60">
          {{ submitting ? 'Applying...' : accountActionButtonLabel }}
        </button>
      </template>
    </Modal>

    <!-- Delete Confirmation Modal -->
    <Modal :show="showDeleteModal" title="Delete User" size="sm" @close="showDeleteModal = false">
      <p class="text-sm text-gray-600">
        Are you sure you want to delete <strong>{{ deleteTarget?.email }}</strong>? This action cannot be undone.
      </p>
      <template #footer>
        <button @click="showDeleteModal = false" class="btn-ghost">Cancel</button>
        <button @click="submitDelete" :disabled="submitting" class="px-4 py-2 text-sm font-semibold bg-red-600 hover:bg-red-700 text-white rounded-lg disabled:opacity-60">
          {{ submitting ? 'Deleting…' : 'Delete' }}
        </button>
      </template>
    </Modal>
  </LayoutDefault>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import AdminTable from '@/components/ui/AdminTable.vue'
import Modal from '@/components/ui/Modal.vue'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'
import { useAuthStore } from '@/stores/auth.js'
import apiClient from '@/api/client.js'

const breadcrumbStore = useBreadcrumbStore()
const authStore = useAuthStore()

const users = ref([])
const loading = ref(true)
const meta = ref(null)
const searchQuery = ref('')
const currentPage = ref(1)
const alertMsg = ref('')
const alertType = ref('success')

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const showResetModal = ref(false)
const showAccountActionModal = ref(false)
const submitting = ref(false)
const modalError = ref('')
const editTarget = ref(null)
const deleteTarget = ref(null)
const resetTarget = ref(null)
const accountActionTarget = ref(null)

const createForm = ref({ first_name: '', last_name: '', email: '', phone: '', password: '', role_id: '' })
const editForm = ref({ first_name: '', last_name: '', phone: '' })
const resetForm = ref({ new_password: '', confirm_password: '' })
const accountActionForm = ref({ action: '', reason: '', suspended_until: '' })
const availableRoles = ref([])
const newRoleId = ref('')
const roleActionLoading = ref(false)

const assignableRoles = computed(() => availableRoles.value.filter(r => r.name !== 'super_admin'))
const assignableRolesForEdit = computed(() => {
  const has = new Set((editTarget.value?.roles || []).map(r => typeof r === 'string' ? r : r.name))
  return assignableRoles.value.filter(r => !has.has(r.name))
})
const accountActionRequiresReason = computed(() => ['SUSPEND', 'BAN', 'MARK_FRAUD'].includes(accountActionForm.value.action))
const accountActionLabels = {
  SUSPEND: 'Suspend Account',
  BAN: 'Ban Account',
  MARK_FRAUD: 'Mark as Fraud',
  CLEAR_FRAUD: 'Clear Fraud Flag',
  REACTIVATE: 'Reactivate Account',
}
const accountActionTitle = computed(() => accountActionLabels[accountActionForm.value.action] || 'Account Action')
const accountActionButtonLabel = computed(() => accountActionLabels[accountActionForm.value.action] || 'Apply Action')
const accountActionDescription = computed(() => {
  const email = accountActionTarget.value?.email || 'this user'
  const descriptions = {
    SUSPEND: `Suspend ${email} and immediately revoke all active sessions.`,
    BAN: `Ban ${email} and immediately revoke all active sessions.`,
    MARK_FRAUD: `Flag ${email} for fraud, ban access, and revoke all active sessions.`,
    CLEAR_FRAUD: `Remove the fraud flag from ${email}. This does not automatically reactivate the account.`,
    REACTIVATE: `Restore access for ${email}.`,
  }
  return descriptions[accountActionForm.value.action] || ''
})

let searchTimer = null
watch(searchQuery, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => { currentPage.value = 1; loadUsers() }, 400)
})

function changePage(p) { currentPage.value = p; loadUsers() }

function formatRoleName(name) {
  if (!name) return ''
  return name.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

async function loadRoles() {
  try {
    const res = await apiClient.get('/api/v1/admin/roles-list')
    availableRoles.value = Array.isArray(res.data.data) ? res.data.data : []
  } catch { availableRoles.value = [] }
}

async function loadUsers() {
  loading.value = true
  try {
    const params = new URLSearchParams({ page: currentPage.value, per_page: 20 })
    if (searchQuery.value) params.set('search', searchQuery.value)
    const res = await apiClient.get(`/api/v1/admin/users?${params}`)
    users.value = Array.isArray(res.data.data) ? res.data.data : []
    meta.value = res.data.meta || null
  } catch (err) {
    showAlert('Failed to load users: ' + (err.response?.data?.error?.message || err.message), 'error')
  } finally { loading.value = false }
}

function getFullName(user) {
  if (user.first_name || user.last_name) return `${user.first_name || ''} ${user.last_name || ''}`.trim()
  return user.email
}
function getInitials(user) {
  if (user.first_name && user.last_name) return `${user.first_name[0]}${user.last_name[0]}`.toUpperCase()
  return user.email?.[0]?.toUpperCase() || 'U'
}
function getRoleName(role) {
  const name = typeof role === 'string' ? role : role.name
  return name.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}
function isProtectedUser(user) {
  return (user.roles || []).some(role => (typeof role === 'string' ? role : role.name) === 'super_admin')
}
function canResetPassword(user) {
  return !isProtectedUser(user) && user.id !== authStore.user?.id
}
function canManageAccount(user) { return canResetPassword(user) }
function accountStatus(user) { return user.account_status || (user.is_active ? 'ACTIVE' : 'SUSPENDED') }
function formatAccountStatus(user) { return formatRoleName(accountStatus(user)) }
function statusBadgeClass(user) {
  return accountStatus(user) === 'ACTIVE'
    ? 'badge-green'
    : accountStatus(user) === 'SUSPENDED'
      ? 'inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-yellow-50 text-yellow-700'
      : 'badge-red'
}
function statusDotClass(user) {
  return accountStatus(user) === 'ACTIVE' ? 'bg-green-500' : accountStatus(user) === 'SUSPENDED' ? 'bg-yellow-500' : 'bg-red-500'
}
function formatDate(dateStr) {
  if (!dateStr) return 'Never'
  try { return new Date(dateStr).toLocaleDateString('en-UG', { day: '2-digit', month: 'short', year: 'numeric' }) }
  catch { return dateStr }
}
function formatDateTime(dateStr) {
  if (!dateStr) return ''
  try { return new Date(dateStr).toLocaleString('en-UG', { dateStyle: 'medium', timeStyle: 'short' }) }
  catch { return dateStr }
}
function apiMessage(err, fallback) {
  const responseData = err.response?.data
  const apiError = err.response?.data?.error
  const details = apiError?.details
  const detailText = details && typeof details === 'object'
    ? Object.entries(details).map(([field, message]) => `${field}: ${message}`).join('; ')
    : ''
  if (apiError?.message) return detailText ? `${apiError.message} (${detailText})` : apiError.message
  if (typeof responseData === 'string' && responseData.trim()) return `${responseData.trim()} (HTTP ${err.response.status})`
  if (responseData?.message) return `${responseData.message} (HTTP ${err.response.status})`
  if (err.response) {
    const requestPath = err.config?.url ? ` for ${err.config.url}` : ''
    const statusText = err.response.statusText ? ` ${err.response.statusText}` : ''
    return `${fallback} (HTTP ${err.response.status}${statusText}${requestPath})`
  }
  if (err.request) return 'The backend could not be reached. Check that the API is running and try again.'
  return err.message || fallback
}
function showAlert(msg, type = 'success') {
  alertMsg.value = msg; alertType.value = type
  setTimeout(() => { alertMsg.value = '' }, 5000)
}

function openCreateModal() {
  createForm.value = { first_name: '', last_name: '', email: '', phone: '', password: '', role_id: '' }
  modalError.value = ''
  if (!availableRoles.value.length) loadRoles()
  showCreateModal.value = true
}
function openEditModal(user) {
  editTarget.value = user
  editForm.value = { first_name: user.first_name || '', last_name: user.last_name || '', phone: user.phone || '' }
  newRoleId.value = ''; modalError.value = ''
  if (!availableRoles.value.length) loadRoles()
  showEditModal.value = true
}
function confirmDelete(user) { deleteTarget.value = user; showDeleteModal.value = true }
function openResetModal(user) {
  resetTarget.value = user
  resetForm.value = { new_password: '', confirm_password: '' }
  modalError.value = ''
  showResetModal.value = true
}
function closeResetModal() {
  showResetModal.value = false
  resetTarget.value = null
  resetForm.value = { new_password: '', confirm_password: '' }
  modalError.value = ''
}
function handleActionSelection(event, user) {
  const action = event.target.value
  event.target.value = ''
  if (!action) return
  if (action === 'edit') return openEditModal(user)
  if (action === 'reset') return openResetModal(user)
  if (action === 'delete') return confirmDelete(user)
  accountActionTarget.value = user
  accountActionForm.value = { action, reason: '', suspended_until: '' }
  modalError.value = ''
  showAccountActionModal.value = true
}
function closeAccountActionModal() {
  showAccountActionModal.value = false
  accountActionTarget.value = null
  accountActionForm.value = { action: '', reason: '', suspended_until: '' }
  modalError.value = ''
}

async function submitCreate() {
  submitting.value = true; modalError.value = ''
  try {
    const { role_id, ...userPayload } = createForm.value
    const res = await apiClient.post('/api/v1/admin/users', userPayload)
    const newUser = res.data?.data || res.data
    if (role_id && newUser?.id) {
      try { await apiClient.post(`/api/v1/admin/users/${newUser.id}/roles`, { role_id }) }
      catch (e) { showAlert('User created but role assignment failed.', 'error') }
    }
    showCreateModal.value = false; showAlert('User created successfully.'); loadUsers()
  } catch (err) { modalError.value = err.response?.data?.error?.message || 'Failed to create user.' }
  finally { submitting.value = false }
}
async function assignRole() {
  if (!newRoleId.value || !editTarget.value?.id) return
  roleActionLoading.value = true
  try {
    await apiClient.post(`/api/v1/admin/users/${editTarget.value.id}/roles`, { role_id: newRoleId.value })
    const assigned = availableRoles.value.find(r => r.id === newRoleId.value)
    if (assigned) editTarget.value = { ...editTarget.value, roles: [...(editTarget.value.roles || []), { id: assigned.id, name: assigned.name }] }
    newRoleId.value = ''; showAlert('Role assigned.'); loadUsers()
  } catch (err) { modalError.value = err.response?.data?.error?.message || 'Failed to assign role.' }
  finally { roleActionLoading.value = false }
}
async function revokeRole(role) {
  if (!editTarget.value?.id) return
  const roleId = typeof role === 'string' ? null : role.id
  if (!roleId) { modalError.value = 'Cannot revoke role without ID.'; return }
  roleActionLoading.value = true
  try {
    await apiClient.delete(`/api/v1/admin/users/${editTarget.value.id}/roles/${roleId}`)
    editTarget.value = { ...editTarget.value, roles: (editTarget.value.roles || []).filter(r => (typeof r === 'string' ? r : r.id) !== roleId) }
    showAlert('Role revoked.'); loadUsers()
  } catch (err) { modalError.value = err.response?.data?.error?.message || 'Failed to revoke role.' }
  finally { roleActionLoading.value = false }
}
async function submitEdit() {
  submitting.value = true; modalError.value = ''
  try {
    await apiClient.put(`/api/v1/admin/users/${editTarget.value.id}`, editForm.value)
    showEditModal.value = false; showAlert('User updated successfully.'); loadUsers()
  } catch (err) { modalError.value = err.response?.data?.error?.message || 'Failed to update user.' }
  finally { submitting.value = false }
}
async function submitDelete() {
  submitting.value = true
  try {
    await apiClient.delete(`/api/v1/admin/users/${deleteTarget.value.id}`)
    showDeleteModal.value = false; showAlert('User deleted successfully.'); loadUsers()
  } catch (err) { showAlert(err.response?.data?.error?.message || 'Failed to delete user.', 'error') }
  finally { submitting.value = false }
}
async function submitPasswordReset() {
  if (!resetTarget.value?.id || resetForm.value.new_password.length < 12 || resetForm.value.new_password !== resetForm.value.confirm_password) return
  submitting.value = true; modalError.value = ''
  try {
    const res = await apiClient.post(`/api/v1/admin/users/${resetTarget.value.id}/reset-password`, {
      new_password: resetForm.value.new_password,
    })
    closeResetModal()
    showAlert(res.data?.message || 'Password reset successfully. All existing sessions were revoked.')
  } catch (err) {
    modalError.value = apiMessage(err, 'Failed to reset password.')
  } finally { submitting.value = false }
}
async function submitAccountAction() {
  if (!accountActionTarget.value?.id) return
  submitting.value = true; modalError.value = ''
  try {
    const payload = {
      action: accountActionForm.value.action,
      reason: accountActionForm.value.reason.trim(),
      suspended_until: accountActionForm.value.suspended_until
        ? new Date(accountActionForm.value.suspended_until).toISOString()
        : null,
    }
    const res = await apiClient.post(`/api/v1/admin/users/${accountActionTarget.value.id}/account-action`, payload)
    closeAccountActionModal()
    showAlert(res.data?.message || 'Account action completed successfully.')
    loadUsers()
  } catch (err) {
    modalError.value = apiMessage(err, 'Account action failed.')
  } finally { submitting.value = false }
}

onMounted(() => {
  breadcrumbStore.set('User Management', [{ label: 'User Management' }])
  loadUsers(); loadRoles()
})
</script>
