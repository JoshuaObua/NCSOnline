<template>
  <LayoutDefault title="Applications">
    <div class="space-y-5">
      <!-- Status Filter Tabs -->
      <div class="flex items-center gap-1 bg-white border border-gray-200 rounded-xl p-1 overflow-x-auto shadow-sm">
        <button
          v-for="tab in statusTabs"
          :key="tab.value"
          @click="setStatusFilter(tab.value)"
          :class="[
            activeStatus === tab.value
              ? 'bg-primary-700 text-white shadow-sm'
              : 'text-gray-600 hover:bg-gray-100',
            'px-3 py-1.5 rounded-lg text-sm font-medium transition-all whitespace-nowrap'
          ]"
        >
          {{ tab.label }}
          <span v-if="tab.value === '' && meta" class="ml-1 text-xs opacity-70">({{ meta.total }})</span>
        </button>
      </div>

      <!-- Alert -->
      <div v-if="alertMsg" :class="alertType === 'success' ? 'bg-green-50 border-green-200 text-green-700' : 'bg-red-50 border-red-200 text-red-700'" class="flex items-center gap-2 p-3 border rounded-lg text-sm">
        {{ alertMsg }}
      </div>

      <!-- Table -->
      <div class="admin-card overflow-hidden">
        <div v-if="loading" class="flex items-center justify-center py-20">
          <svg class="animate-spin w-7 h-7 text-primary-600" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>

        <table v-else class="w-full">
          <thead>
            <tr class="border-b border-gray-100">
              <th class="table-th">Reference</th>
              <th class="table-th">Applicant</th>
              <th class="table-th">Form Type</th>
              <th class="table-th">Status</th>
              <th class="table-th">Payment</th>
              <th class="table-th">Submitted</th>
              <th class="table-th text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-if="applications.length === 0">
              <td colspan="7" class="text-center py-16 text-gray-400">
                <div class="flex flex-col items-center gap-2">
                  <svg class="w-12 h-12 text-gray-300" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                  </svg>
                  <p class="text-sm">No applications found</p>
                </div>
              </td>
            </tr>
            <tr
              v-for="app in applications"
              :key="app.id"
              class="hover:bg-gray-50 cursor-pointer transition-colors"
              @click="openDetail(app)"
            >
              <td class="px-6 py-3 text-sm font-medium text-primary-700">
                {{ app.reference_number || app.id?.substring(0, 8)?.toUpperCase() || 'N/A' }}
              </td>
              <td class="px-6 py-3 text-sm text-gray-800">
                {{ getApplicantName(app) }}
              </td>
              <td class="px-6 py-3 text-sm text-gray-600">
                {{ formatFormType(app.form_type) }}
              </td>
              <td class="px-6 py-3">
                <StatusBadge :status="app.status" />
              </td>
              <td class="px-6 py-3">
                <span :class="app.payment_status === 'PAID' ? 'text-green-600' : 'text-gray-400'" class="text-sm">
                  {{ app.payment_status || 'N/A' }}
                </span>
              </td>
              <td class="px-6 py-3 text-sm text-gray-500">
                {{ formatDate(app.submitted_at || app.created_at) }}
              </td>
              <td class="px-6 py-3 text-right" @click.stop>
                <button @click="openDetail(app)" class="p-1.5 text-gray-400 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors" title="View Details">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Pagination -->
        <div v-if="meta && meta.total > 0" class="admin-card-footer flex items-center justify-between">
          <p class="text-xs text-gray-500">
            Showing <span class="font-medium text-gray-700">{{ (meta.page - 1) * meta.per_page + 1 }}–{{ Math.min(meta.page * meta.per_page, meta.total) }}</span> of <span class="font-medium text-gray-700">{{ meta.total }}</span>
          </p>
          <div class="flex items-center gap-1.5">
            <button @click="prevPage" :disabled="meta.page <= 1" class="px-3 py-1 text-xs font-medium border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors">← Prev</button>
            <span class="text-xs text-gray-500 px-2 font-medium">{{ meta.page }}</span>
            <button @click="nextPage" :disabled="meta.page * meta.per_page >= meta.total" class="px-3 py-1 text-xs font-medium border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors">Next →</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Application Detail Modal -->
    <Modal :show="!!selectedApp" :title="`Application — ${selectedApp?.reference_number || selectedApp?.id?.substring(0,8)?.toUpperCase() || 'Detail'}`" size="xl" @close="selectedApp = null">
      <div v-if="selectedApp" class="space-y-4">
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span class="text-gray-500">Form Type</span>
            <div class="font-medium text-gray-900">{{ formatFormType(selectedApp.form_type) }}</div>
          </div>
          <div>
            <span class="text-gray-500">Status</span>
            <div class="mt-0.5"><StatusBadge :status="selectedApp.status" /></div>
          </div>
          <div>
            <span class="text-gray-500">Applicant</span>
            <div class="font-medium text-gray-900">{{ getApplicantName(selectedApp) }}</div>
          </div>
          <div>
            <span class="text-gray-500">Submitted</span>
            <div class="font-medium text-gray-900">{{ formatDate(selectedApp.submitted_at || selectedApp.created_at) }}</div>
          </div>
          <div v-if="selectedApp.payment_status">
            <span class="text-gray-500">Payment</span>
            <div class="font-medium text-gray-900">{{ selectedApp.payment_status }}</div>
          </div>
          <div v-if="selectedApp.notes">
            <span class="text-gray-500">Notes</span>
            <div class="font-medium text-gray-900">{{ selectedApp.notes }}</div>
          </div>
        </div>

        <!-- Admin Actions for reviewable statuses -->
        <div v-if="isAdmin && canReview(selectedApp)" class="border-t border-gray-100 pt-4">
          <h4 class="text-sm font-medium text-gray-700 mb-3">Admin Actions</h4>
          <div class="space-y-3">
            <textarea
              v-model="actionNotes"
              rows="2"
              placeholder="Add notes (optional)..."
              class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm resize-none focus:outline-none focus:ring-2 focus:ring-primary-500"
            ></textarea>
            <div class="flex gap-3">
              <button
                @click="approveApplication(selectedApp)"
                :disabled="actionLoading"
                class="flex-1 py-2 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-lg transition-colors disabled:opacity-60"
              >
                {{ actionLoading === 'approve' ? 'Approving...' : 'Approve' }}
              </button>
              <button
                @click="rejectApplication(selectedApp)"
                :disabled="actionLoading"
                class="flex-1 py-2 bg-red-600 hover:bg-red-700 text-white text-sm font-medium rounded-lg transition-colors disabled:opacity-60"
              >
                {{ actionLoading === 'reject' ? 'Rejecting...' : 'Reject' }}
              </button>
            </div>
          </div>
          <div v-if="actionError" class="mt-2 p-2 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700">{{ actionError }}</div>
        </div>
      </div>
      <template #footer>
        <button @click="selectedApp = null" class="px-4 py-2 text-sm border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors">Close</button>
      </template>
    </Modal>
  </LayoutDefault>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'
const breadcrumbStore = useBreadcrumbStore()
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import Modal from '@/components/ui/Modal.vue'
import apiClient from '@/api/client.js'
import { useAuthStore } from '@/stores/auth.js'

const authStore = useAuthStore()
const isAdmin = computed(() => authStore.isAdmin)

const applications = ref([])
const loading = ref(true)
const meta = ref(null)
const activeStatus = ref('')
const currentPage = ref(1)
const alertMsg = ref('')
const alertType = ref('success')
const selectedApp = ref(null)
const actionNotes = ref('')
const actionLoading = ref('')
const actionError = ref('')

const statusTabs = [
  { label: 'All', value: '' },
  { label: 'Draft', value: 'DRAFT' },
  { label: 'Submitted', value: 'SUBMITTED' },
  { label: 'Under Review', value: 'UNDER_REVIEW' },
  { label: 'Approved', value: 'APPROVED' },
  { label: 'Rejected', value: 'REJECTED' },
  { label: 'Queried', value: 'NEEDS_INFORMATION' }
]

function canReview(app) {
  return ['SUBMITTED', 'UNDER_REVIEW', 'NEEDS_INFORMATION'].includes(app.status)
}

function setStatusFilter(status) {
  activeStatus.value = status
  currentPage.value = 1
  loadApplications()
}

async function loadApplications() {
  loading.value = true
  try {
    const params = new URLSearchParams({ page: currentPage.value, per_page: 20 })
    if (activeStatus.value) params.set('status', activeStatus.value)

    let res
    try {
      res = await apiClient.get(`/api/v1/admin/applications?${params}`)
    } catch {
      res = await apiClient.get(`/api/v1/applications?${params}`)
    }

    applications.value = Array.isArray(res.data.data) ? res.data.data : []
    meta.value = res.data.meta || null
  } catch (err) {
    showAlert('Failed to load applications: ' + (err.response?.data?.error?.message || err.message), 'error')
  } finally {
    loading.value = false
  }
}

function getApplicantName(app) {
  if (app.applicant_name) return app.applicant_name
  if (app.user) {
    const u = app.user
    if (u.first_name || u.last_name) return `${u.first_name || ''} ${u.last_name || ''}`.trim()
    return u.email || 'N/A'
  }
  return app.user_email || 'N/A'
}

function formatFormType(type) {
  if (!type) return 'N/A'
  return type.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

function formatDate(dateStr) {
  if (!dateStr) return 'N/A'
  try {
    return new Date(dateStr).toLocaleDateString('en-UG', { day: '2-digit', month: 'short', year: 'numeric' })
  } catch { return dateStr }
}

function openDetail(app) {
  selectedApp.value = app
  actionNotes.value = ''
  actionError.value = ''
  actionLoading.value = ''
}

function showAlert(msg, type = 'success') {
  alertMsg.value = msg
  alertType.value = type
  setTimeout(() => { alertMsg.value = '' }, 5000)
}

async function approveApplication(app) {
  actionLoading.value = 'approve'
  actionError.value = ''
  try {
    await apiClient.post(`/api/v1/admin/applications/${app.id}/approve`, { notes: actionNotes.value })
    selectedApp.value = null
    showAlert('Application approved successfully.')
    loadApplications()
  } catch (err) {
    actionError.value = err.response?.data?.error?.message || 'Failed to approve application.'
  } finally {
    actionLoading.value = ''
  }
}

async function rejectApplication(app) {
  actionLoading.value = 'reject'
  actionError.value = ''
  try {
    await apiClient.post(`/api/v1/admin/applications/${app.id}/reject`, { notes: actionNotes.value })
    selectedApp.value = null
    showAlert('Application rejected.')
    loadApplications()
  } catch (err) {
    actionError.value = err.response?.data?.error?.message || 'Failed to reject application.'
  } finally {
    actionLoading.value = ''
  }
}

function prevPage() {
  if (currentPage.value > 1) { currentPage.value--; loadApplications() }
}
function nextPage() {
  if (meta.value && currentPage.value * meta.value.per_page < meta.value.total) { currentPage.value++; loadApplications() }
}

onMounted(() => {
  breadcrumbStore.set('Applications', [{ label: 'Management' }, { label: 'Applications' }])
  loadApplications()
})
</script>
