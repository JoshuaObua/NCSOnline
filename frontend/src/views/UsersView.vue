<template>
  <LayoutDefault title="User Management">
    <div class="space-y-5">
      <!-- Toolbar -->
      <div class="flex flex-col sm:flex-row sm:items-center gap-3">
        <div class="relative flex-1 max-w-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
            </svg>
          </div>
          <input
            v-model="searchQuery"
            @input="debouncedSearch"
            type="text"
            placeholder="Search users..."
            class="w-full pl-9 pr-4 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          />
        </div>
        <button
          @click="openCreateModal"
          class="flex items-center gap-2 px-4 py-2 bg-primary-700 hover:bg-primary-600 text-white rounded-lg text-sm font-medium transition-colors"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          Create User
        </button>
      </div>

      <!-- Alert -->
      <div v-if="alertMsg" :class="alertType === 'success' ? 'bg-green-50 border-green-200 text-green-700' : 'bg-red-50 border-red-200 text-red-700'" class="flex items-center gap-2 p-3 border rounded-lg text-sm">
        <svg class="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path v-if="alertType === 'success'" stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          <path v-else stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
        </svg>
        {{ alertMsg }}
      </div>

      <!-- Table -->
      <div class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
        <!-- Loading overlay -->
        <div v-if="loading" class="flex items-center justify-center py-20">
          <svg class="animate-spin w-7 h-7 text-primary-600" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>

        <table v-else class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">User</th>
              <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Status</th>
              <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Roles</th>
              <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Last Login</th>
              <th class="text-right text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-if="users.length === 0">
              <td colspan="5" class="text-center py-16 text-gray-400">
                <div class="flex flex-col items-center gap-2">
                  <svg class="w-12 h-12 text-gray-300" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
                  </svg>
                  <p class="text-sm">No users found</p>
                </div>
              </td>
            </tr>
            <tr
              v-for="user in users"
              :key="user.id"
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 bg-primary-100 text-primary-700 rounded-full flex items-center justify-center text-sm font-semibold flex-shrink-0">
                    {{ getInitials(user) }}
                  </div>
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ getFullName(user) }}</div>
                    <div class="text-xs text-gray-500">{{ user.email }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-3">
                <StatusBadge :status="user.is_active ? 'active' : 'inactive'" />
              </td>
              <td class="px-6 py-3">
                <div class="flex flex-wrap gap-1">
                  <span
                    v-for="role in (user.roles || [])"
                    :key="getRoleName(role)"
                    class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-indigo-50 text-indigo-700"
                  >
                    {{ getRoleName(role) }}
                  </span>
                  <span v-if="!user.roles || user.roles.length === 0" class="text-xs text-gray-400">No roles</span>
                </div>
              </td>
              <td class="px-6 py-3 text-sm text-gray-500">
                {{ formatDate(user.last_login_at) }}
              </td>
              <td class="px-6 py-3">
                <div class="flex items-center justify-end gap-2">
                  <button
                    @click="openEditModal(user)"
                    class="p-1.5 text-gray-400 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors"
                    title="Edit"
                  >
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                    </svg>
                  </button>
                  <button
                    @click="toggleUserStatus(user)"
                    :class="user.is_active ? 'text-gray-400 hover:text-yellow-600 hover:bg-yellow-50' : 'text-gray-400 hover:text-green-600 hover:bg-green-50'"
                    class="p-1.5 rounded-lg transition-colors"
                    :title="user.is_active ? 'Deactivate' : 'Activate'"
                  >
                    <svg v-if="user.is_active" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                    </svg>
                    <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </button>
                  <button
                    @click="confirmDelete(user)"
                    class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    title="Delete"
                  >
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Pagination -->
        <div v-if="meta && meta.total > 0" class="flex items-center justify-between px-6 py-3 border-t border-gray-100 bg-gray-50">
          <p class="text-sm text-gray-500">
            Showing {{ (meta.page - 1) * meta.per_page + 1 }}–{{ Math.min(meta.page * meta.per_page, meta.total) }} of {{ meta.total }} users
          </p>
          <div class="flex items-center gap-2">
            <button
              @click="prevPage"
              :disabled="meta.page <= 1"
              class="px-3 py-1.5 text-sm border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >Previous</button>
            <span class="text-sm text-gray-600 px-2">Page {{ meta.page }}</span>
            <button
              @click="nextPage"
              :disabled="meta.page * meta.per_page >= meta.total"
              class="px-3 py-1.5 text-sm border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >Next</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create User Modal -->
    <Modal :show="showCreateModal" title="Create New User" size="md" @close="showCreateModal = false">
      <form @submit.prevent="submitCreate" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">First Name</label>
            <input v-model="createForm.first_name" type="text" required class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Last Name</label>
            <input v-model="createForm.last_name" type="text" required class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
          <input v-model="createForm.email" type="email" required class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
          <input v-model="createForm.phone" type="tel" class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" placeholder="+256..." />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input v-model="createForm.password" type="password" required class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
        </div>
        <div v-if="modalError" class="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{{ modalError }}</div>
      </form>
      <template #footer>
        <button @click="showCreateModal = false" class="px-4 py-2 text-sm border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors">Cancel</button>
        <button @click="submitCreate" :disabled="submitting" class="px-4 py-2 text-sm bg-primary-700 hover:bg-primary-600 text-white rounded-lg font-medium transition-colors disabled:opacity-60">
          {{ submitting ? 'Creating...' : 'Create User' }}
        </button>
      </template>
    </Modal>

    <!-- Edit User Modal -->
    <Modal :show="showEditModal" title="Edit User" size="md" @close="showEditModal = false">
      <form @submit.prevent="submitEdit" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">First Name</label>
            <input v-model="editForm.first_name" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Last Name</label>
            <input v-model="editForm.last_name" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
          <input v-model="editForm.phone" type="tel" class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500" />
        </div>
        <div v-if="modalError" class="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{{ modalError }}</div>
      </form>
      <template #footer>
        <button @click="showEditModal = false" class="px-4 py-2 text-sm border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors">Cancel</button>
        <button @click="submitEdit" :disabled="submitting" class="px-4 py-2 text-sm bg-primary-700 hover:bg-primary-600 text-white rounded-lg font-medium transition-colors disabled:opacity-60">
          {{ submitting ? 'Saving...' : 'Save Changes' }}
        </button>
      </template>
    </Modal>

    <!-- Delete Confirmation Modal -->
    <Modal :show="showDeleteModal" title="Delete User" size="sm" @close="showDeleteModal = false">
      <p class="text-sm text-gray-600">
        Are you sure you want to delete <strong>{{ deleteTarget?.email }}</strong>? This action cannot be undone.
      </p>
      <template #footer>
        <button @click="showDeleteModal = false" class="px-4 py-2 text-sm border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors">Cancel</button>
        <button @click="submitDelete" :disabled="submitting" class="px-4 py-2 text-sm bg-red-600 hover:bg-red-700 text-white rounded-lg font-medium transition-colors disabled:opacity-60">
          {{ submitting ? 'Deleting...' : 'Delete' }}
        </button>
      </template>
    </Modal>
  </LayoutDefault>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import Modal from '@/components/ui/Modal.vue'
import apiClient from '@/api/client.js'

const users = ref([])
const loading = ref(true)
const meta = ref(null)
const searchQuery = ref('')
const currentPage = ref(1)
const alertMsg = ref('')
const alertType = ref('success')

// Modal state
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const submitting = ref(false)
const modalError = ref('')
const editTarget = ref(null)
const deleteTarget = ref(null)

const createForm = ref({ first_name: '', last_name: '', email: '', phone: '', password: '' })
const editForm = ref({ first_name: '', last_name: '', phone: '' })

let searchTimer = null

function debouncedSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    currentPage.value = 1
    loadUsers()
  }, 400)
}

async function loadUsers() {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: currentPage.value,
      per_page: 20
    })
    if (searchQuery.value) params.set('search', searchQuery.value)

    const res = await apiClient.get(`/api/v1/admin/users?${params}`)
    users.value = Array.isArray(res.data.data) ? res.data.data : []
    meta.value = res.data.meta || null
  } catch (err) {
    showAlert('Failed to load users: ' + (err.response?.data?.error?.message || err.message), 'error')
  } finally {
    loading.value = false
  }
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
  try {
    return new Date(dateStr).toLocaleDateString('en-UG', { day: '2-digit', month: 'short', year: 'numeric' })
  } catch { return dateStr }
}

function showAlert(msg, type = 'success') {
  alertMsg.value = msg
  alertType.value = type
  setTimeout(() => { alertMsg.value = '' }, 5000)
}

function openCreateModal() {
  createForm.value = { first_name: '', last_name: '', email: '', phone: '', password: '' }
  modalError.value = ''
  showCreateModal.value = true
}

function openEditModal(user) {
  editTarget.value = user
  editForm.value = {
    first_name: user.first_name || '',
    last_name: user.last_name || '',
    phone: user.phone || ''
  }
  modalError.value = ''
  showEditModal.value = true
}

function confirmDelete(user) {
  deleteTarget.value = user
  showDeleteModal.value = true
}

async function submitCreate() {
  submitting.value = true
  modalError.value = ''
  try {
    await apiClient.post('/api/v1/admin/users', createForm.value)
    showCreateModal.value = false
    showAlert('User created successfully.')
    loadUsers()
  } catch (err) {
    modalError.value = err.response?.data?.error?.message || 'Failed to create user.'
  } finally {
    submitting.value = false
  }
}

async function submitEdit() {
  submitting.value = true
  modalError.value = ''
  try {
    await apiClient.put(`/api/v1/admin/users/${editTarget.value.id}`, editForm.value)
    showEditModal.value = false
    showAlert('User updated successfully.')
    loadUsers()
  } catch (err) {
    modalError.value = err.response?.data?.error?.message || 'Failed to update user.'
  } finally {
    submitting.value = false
  }
}

async function submitDelete() {
  submitting.value = true
  try {
    await apiClient.delete(`/api/v1/admin/users/${deleteTarget.value.id}`)
    showDeleteModal.value = false
    showAlert('User deleted successfully.')
    loadUsers()
  } catch (err) {
    showAlert(err.response?.data?.error?.message || 'Failed to delete user.', 'error')
  } finally {
    submitting.value = false
  }
}

async function toggleUserStatus(user) {
  try {
    if (user.is_active) {
      await apiClient.post(`/api/v1/admin/users/${user.id}/deactivate`)
      showAlert(`${getFullName(user)} deactivated.`)
    } else {
      await apiClient.post(`/api/v1/admin/users/${user.id}/activate`)
      showAlert(`${getFullName(user)} activated.`)
    }
    loadUsers()
  } catch (err) {
    showAlert(err.response?.data?.error?.message || 'Action failed.', 'error')
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
    loadUsers()
  }
}

function nextPage() {
  if (meta.value && currentPage.value * meta.value.per_page < meta.value.total) {
    currentPage.value++
    loadUsers()
  }
}

onMounted(loadUsers)
</script>
