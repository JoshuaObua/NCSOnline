<template>
  <LayoutDefault title="Roles & Permissions">
    <div class="flex gap-5 min-h-0">

      <!-- ── Left: Role list ─────────────────────────────────────── -->
      <div class="w-64 flex-shrink-0 space-y-2">
        <button
          @click="openCreate"
          class="w-full flex items-center justify-center gap-2 bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-4 py-2.5 rounded-xl transition-colors"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
          New Role
        </button>

        <div v-if="loading" class="flex items-center justify-center py-10">
          <svg class="animate-spin w-6 h-6 text-primary-600" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
          </svg>
        </div>

        <button
          v-for="role in roles"
          :key="role.id"
          @click="selectRole(role)"
          :class="selectedRole?.id === role.id
            ? 'border-primary-500 bg-primary-50 ring-1 ring-primary-300'
            : 'border-gray-200 hover:border-gray-300 bg-white'"
          class="w-full text-left rounded-xl border p-3.5 transition-all cursor-pointer"
        >
          <div class="flex items-start justify-between gap-2">
            <div class="min-w-0">
              <div class="flex items-center gap-1.5">
                <span class="font-medium text-sm text-gray-900 truncate">{{ fmt(role.name) }}</span>
                <span v-if="role.is_system" class="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-semibold bg-orange-100 text-orange-700 flex-shrink-0">SYS</span>
              </div>
              <p v-if="role.description" class="text-xs text-gray-400 mt-0.5 truncate">{{ role.description }}</p>
            </div>
          </div>
          <div class="mt-2 text-xs text-gray-400">
            {{ permCounts[role.id] !== undefined ? permCounts[role.id] + ' permission' + (permCounts[role.id] !== 1 ? 's' : '') : '—' }}
          </div>
        </button>

        <div v-if="!loading && roles.length === 0" class="text-center py-8 text-sm text-gray-400">
          No roles found
        </div>
      </div>

      <!-- ── Right: Permission matrix ────────────────────────────── -->
      <div class="flex-1 min-w-0">
        <!-- Empty state -->
        <div v-if="!selectedRole" class="flex flex-col items-center justify-center h-64 text-gray-400">
          <svg class="w-12 h-12 mb-3 text-gray-300" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" />
          </svg>
          <p class="text-sm">Select a role to manage its permissions</p>
        </div>

        <div v-else class="space-y-4">

          <!-- Role header card -->
          <div class="bg-white rounded-xl border border-gray-200 shadow-sm p-5 flex items-start justify-between gap-4">
            <div>
              <div class="flex items-center gap-2">
                <h2 class="text-lg font-bold text-gray-900">{{ fmt(selectedRole.name) }}</h2>
                <span v-if="selectedRole.is_system" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-orange-100 text-orange-700">System Role</span>
              </div>
              <p class="text-sm text-gray-500 mt-0.5">{{ selectedRole.description || 'No description' }}</p>
              <p class="text-xs text-gray-400 mt-1">
                {{ rolePermSet.size }} of {{ allPerms.length }} permissions granted
              </p>
            </div>
            <div v-if="!selectedRole.is_system" class="flex items-center gap-2 flex-shrink-0">
              <button
                @click="openEdit"
                class="flex items-center gap-1.5 text-sm text-gray-600 hover:text-primary-700 border border-gray-200 hover:border-primary-300 px-3 py-1.5 rounded-lg transition-colors"
              >
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125" />
                </svg>
                Edit
              </button>
              <button
                @click="confirmDelete"
                class="flex items-center gap-1.5 text-sm text-red-600 hover:text-red-700 border border-red-200 hover:border-red-300 px-3 py-1.5 rounded-lg transition-colors"
              >
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                </svg>
                Delete
              </button>
            </div>
          </div>

          <!-- Super admin notice -->
          <div v-if="selectedRole.name === 'super_admin'" class="flex items-center gap-2 p-3 bg-amber-50 border border-amber-200 rounded-xl text-sm text-amber-700">
            <svg class="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
            </svg>
            Super Admin has all permissions and cannot be restricted.
          </div>

          <!-- Loading permissions -->
          <div v-if="loadingPerms" class="flex items-center justify-center py-10">
            <svg class="animate-spin w-6 h-6 text-primary-600" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
            </svg>
          </div>

          <!-- Permission grid by resource -->
          <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div
              v-for="(perms, resource) in allPermsByResource"
              :key="resource"
              class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden"
            >
              <!-- Resource header -->
              <div class="px-5 py-3 border-b border-gray-100 flex items-center justify-between">
                <h3 class="font-semibold text-sm text-gray-800 uppercase tracking-wide">
                  {{ fmt(resource) }}
                </h3>
                <div class="flex items-center gap-3 text-xs text-gray-500">
                  <span>{{ countGranted(perms) }}/{{ perms.length }}</span>
                  <button
                    v-if="selectedRole.name !== 'super_admin'"
                    @click="toggleAll(perms, !allGranted(perms))"
                    class="text-primary-600 hover:text-primary-800 font-medium transition-colors"
                  >
                    {{ allGranted(perms) ? 'None' : 'All' }}
                  </button>
                </div>
              </div>

              <!-- Checkboxes -->
              <div class="px-5 py-4 space-y-3">
                <label
                  v-for="perm in perms"
                  :key="perm.id"
                  class="flex items-start gap-3 cursor-pointer group"
                  :class="selectedRole.name === 'super_admin' ? 'cursor-default opacity-70' : ''"
                >
                  <div class="flex-shrink-0 mt-0.5">
                    <input
                      type="checkbox"
                      :checked="rolePermSet.has(perm.id)"
                      :disabled="selectedRole.name === 'super_admin' || saving.has(perm.id)"
                      @change="togglePerm(perm, $event.target.checked)"
                      class="w-4 h-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500 cursor-pointer"
                    />
                  </div>
                  <div class="min-w-0">
                    <div class="flex items-center gap-1.5">
                      <span class="text-sm font-medium text-gray-800 group-hover:text-primary-700 transition-colors">
                        {{ fmtAction(perm.action) }}
                      </span>
                      <svg v-if="saving.has(perm.id)" class="animate-spin w-3 h-3 text-primary-500" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                      </svg>
                    </div>
                    <p v-if="perm.description" class="text-xs text-gray-400 mt-0.5">{{ perm.description }}</p>
                  </div>
                </label>
              </div>
            </div>
          </div>

          <!-- Error -->
          <div v-if="permError" class="flex items-center gap-2 p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">
            {{ permError }}
          </div>
        </div>
      </div>
    </div>

    <!-- ── Role create/edit modal ──────────────────────────────────── -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-2xl shadow-2xl w-full max-w-md">
          <div class="px-6 py-5 border-b border-gray-100">
            <h2 class="text-lg font-bold text-gray-900">{{ editing ? 'Edit Role' : 'New Role' }}</h2>
          </div>
          <form @submit.prevent="saveRole" class="px-6 py-5 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Role Name</label>
              <input
                v-model="form.name"
                :disabled="editing && selectedRole?.is_system"
                type="text"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent disabled:bg-gray-50 disabled:text-gray-500"
                placeholder="e.g. review_officer"
                required
              />
              <p class="text-xs text-gray-400 mt-1">Use lowercase with underscores. Cannot be changed after creation.</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <textarea
                v-model="form.description"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent resize-none"
                rows="3"
                placeholder="Brief description of this role's purpose"
              />
            </div>
            <div v-if="formError" class="text-sm text-red-600 bg-red-50 border border-red-200 rounded-lg px-3 py-2">
              {{ formError }}
            </div>
          </form>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-2">
            <button
              @click="showModal = false"
              class="px-4 py-2 text-sm text-gray-600 hover:text-gray-800 border border-gray-200 rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              @click="saveRole"
              :disabled="formSaving"
              class="px-4 py-2 text-sm font-medium bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors disabled:opacity-60"
            >
              {{ formSaving ? 'Saving…' : editing ? 'Save Changes' : 'Create Role' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ── Delete confirm modal ────────────────────────────────────── -->
    <Teleport to="body">
      <div v-if="showDeleteConfirm" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-2xl shadow-2xl w-full max-w-sm p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-red-600" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-gray-900">Delete Role</h3>
              <p class="text-sm text-gray-500">This action cannot be undone.</p>
            </div>
          </div>
          <p class="text-sm text-gray-700 mb-5">
            Are you sure you want to delete <strong>{{ fmt(selectedRole?.name) }}</strong>? Users with this role will lose their access.
          </p>
          <div class="flex justify-end gap-2">
            <button @click="showDeleteConfirm = false" class="px-4 py-2 text-sm border border-gray-200 text-gray-600 rounded-lg hover:bg-gray-50 transition-colors">
              Cancel
            </button>
            <button @click="deleteRole" :disabled="formSaving" class="px-4 py-2 text-sm font-medium bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors disabled:opacity-60">
              {{ formSaving ? 'Deleting…' : 'Delete' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </LayoutDefault>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import apiClient from '@/api/client.js'

const roles = ref([])
const allPerms = ref([])
const loading = ref(true)
const loadingPerms = ref(false)
const selectedRole = ref(null)
const rolePermSet = ref(new Set())
const permCounts = reactive({})
const saving = ref(new Set())
const permError = ref('')

const showModal = ref(false)
const editing = ref(false)
const form = reactive({ name: '', description: '' })
const formSaving = ref(false)
const formError = ref('')
const showDeleteConfirm = ref(false)

const allPermsByResource = computed(() => {
  const grouped = {}
  for (const p of allPerms.value) {
    if (!grouped[p.resource]) grouped[p.resource] = []
    grouped[p.resource].push(p)
  }
  return grouped
})

function fmt(name) {
  if (!name) return ''
  return name.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

function fmtAction(action) {
  if (!action) return ''
  return action.replace(/:/g, ' › ').replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

function countGranted(perms) {
  return perms.filter(p => rolePermSet.value.has(p.id)).length
}

function allGranted(perms) {
  return perms.every(p => rolePermSet.value.has(p.id))
}

async function selectRole(role) {
  selectedRole.value = role
  permError.value = ''
  loadingPerms.value = true
  try {
    const res = await apiClient.get(`/api/v1/admin/roles/${role.id}`)
    const data = res.data?.data ?? res.data
    const perms = Array.isArray(data?.permissions) ? data.permissions : []
    rolePermSet.value = new Set(perms.map(p => p.id))
    permCounts[role.id] = perms.length
  } catch {
    permError.value = 'Failed to load permissions for this role.'
  } finally {
    loadingPerms.value = false
  }
}

async function togglePerm(perm, granted) {
  if (saving.value.has(perm.id)) return
  const newSet = new Set(rolePermSet.value)
  const roleId = selectedRole.value.id

  // optimistic update
  if (granted) newSet.add(perm.id)
  else newSet.delete(perm.id)
  rolePermSet.value = newSet

  const prev = new Set(rolePermSet.value)
  saving.value = new Set([...saving.value, perm.id])
  try {
    if (granted) {
      await apiClient.post(`/api/v1/admin/roles/${roleId}/permissions`, { permission_id: perm.id })
    } else {
      await apiClient.delete(`/api/v1/admin/roles/${roleId}/permissions/${perm.id}`)
    }
    permCounts[roleId] = rolePermSet.value.size
  } catch {
    // revert on failure
    rolePermSet.value = prev
    permError.value = `Failed to ${granted ? 'grant' : 'revoke'} permission.`
    setTimeout(() => { permError.value = '' }, 3000)
  } finally {
    const s = new Set(saving.value)
    s.delete(perm.id)
    saving.value = s
  }
}

async function toggleAll(perms, grant) {
  const tasks = perms
    .filter(p => grant ? !rolePermSet.value.has(p.id) : rolePermSet.value.has(p.id))
    .map(p => togglePerm(p, grant))
  await Promise.allSettled(tasks)
}

function openCreate() {
  editing.value = false
  form.name = ''
  form.description = ''
  formError.value = ''
  showModal.value = true
}

function openEdit() {
  editing.value = true
  form.name = selectedRole.value.name
  form.description = selectedRole.value.description || ''
  formError.value = ''
  showModal.value = true
}

function confirmDelete() {
  showDeleteConfirm.value = true
}

async function saveRole() {
  if (!form.name.trim()) { formError.value = 'Role name is required.'; return }
  formSaving.value = true
  formError.value = ''
  try {
    if (editing.value) {
      await apiClient.put(`/api/v1/admin/roles/${selectedRole.value.id}`, { name: form.name, description: form.description })
      selectedRole.value.name = form.name
      selectedRole.value.description = form.description
      const idx = roles.value.findIndex(r => r.id === selectedRole.value.id)
      if (idx !== -1) { roles.value[idx].name = form.name; roles.value[idx].description = form.description }
    } else {
      const res = await apiClient.post('/api/v1/admin/roles', { name: form.name, description: form.description })
      const newRole = res.data?.data ?? res.data
      roles.value.push(newRole)
      permCounts[newRole.id] = 0
    }
    showModal.value = false
  } catch (err) {
    formError.value = err.response?.data?.error?.message || 'Failed to save role.'
  } finally {
    formSaving.value = false
  }
}

async function deleteRole() {
  formSaving.value = true
  try {
    await apiClient.delete(`/api/v1/admin/roles/${selectedRole.value.id}`)
    roles.value = roles.value.filter(r => r.id !== selectedRole.value.id)
    delete permCounts[selectedRole.value.id]
    selectedRole.value = null
    showDeleteConfirm.value = false
  } catch (err) {
    permError.value = err.response?.data?.error?.message || 'Failed to delete role.'
    showDeleteConfirm.value = false
  } finally {
    formSaving.value = false
  }
}

async function loadData() {
  loading.value = true
  try {
    const [rolesRes, permsRes] = await Promise.all([
      apiClient.get('/api/v1/admin/roles'),
      apiClient.get('/api/v1/admin/permissions')
    ])
    roles.value = Array.isArray(rolesRes.data?.data) ? rolesRes.data.data : []
    allPerms.value = Array.isArray(permsRes.data?.data) ? permsRes.data.data : []
    // pre-select first role
    if (roles.value.length > 0) await selectRole(roles.value[0])
  } catch {
    // silently fail — user sees empty state
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>
