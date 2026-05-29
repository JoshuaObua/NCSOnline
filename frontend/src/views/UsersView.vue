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
              <span :class="user.is_active ? 'badge-green' : 'badge-red'">
                <span class="w-1.5 h-1.5 rounded-full mr-1" :class="user.is_active ? 'bg-green-500' : 'bg-red-500'"></span>
                {{ user.is_active ? 'Active' : 'Inactive' }}
              </span>
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
              <div class="flex items-center justify-end gap-1">
                <button @click="openEditModal(user)" title="Edit" class="p-1.5 text-gray-400 hover:text-primary-700 hover:bg-primary-50 rounded-lg transition-colors">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931z" />
                  </svg>
                </button>
                <button
                  @click="toggleUserStatus(user)"
                  :title="user.is_active ? 'Deactivate' : 'Activate'"
                  :class="user.is_active ? 'hover:text-yellow-600 hover:bg-yellow-50' : 'hover:text-green-600 hover:bg-green-50'"
                  class="p-1.5 text-gray-400 rounded-lg transition-colors"
                >
                  <svg v-if="user.is_active" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                  </svg>
                  <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </button>
                <button @click="confirmDelete(user)" title="Delete" class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                  </svg>
                </button>
              </div>
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
import apiClient from '@/api/client.js'

const breadcrumbStore = useBreadcrumbStore()

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
const submitting = ref(false)
const modalError = ref('')
const editTarget = ref(null)
const deleteTarget = ref(null)

const createForm = ref({ first_name: '', last_name: '', email: '', phone: '', password: '', role_id: '' })
const editForm = ref({ first_name: '', last_name: '', phone: '' })
const availableRoles = ref([])
const newRoleId = ref('')
const roleActionLoading = ref(false)

const assignableRoles = computed(() => availableRoles.value.filter(r => r.name !== 'super_admin'))
const assignableRolesForEdit = computed(() => {
  const has = new Set((editTarget.value?.roles || []).map(r => typeof r === 'string' ? r : r.name))
  return assignableRoles.value.filter(r => !has.has(r.name))
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
function formatDate(dateStr) {
  if (!dateStr) return 'Never'
  try { return new Date(dateStr).toLocaleDateString('en-UG', { day: '2-digit', month: 'short', year: 'numeric' }) }
  catch { return dateStr }
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
async function toggleUserStatus(user) {
  try {
    if (user.is_active) await apiClient.post(`/api/v1/admin/users/${user.id}/deactivate`)
    else await apiClient.post(`/api/v1/admin/users/${user.id}/activate`)
    showAlert(`${getFullName(user)} ${user.is_active ? 'deactivated' : 'activated'}.`); loadUsers()
  } catch (err) { showAlert(err.response?.data?.error?.message || 'Action failed.', 'error') }
}

onMounted(() => {
  breadcrumbStore.set('User Management', [{ label: 'User Management' }])
  loadUsers(); loadRoles()
})
</script>
