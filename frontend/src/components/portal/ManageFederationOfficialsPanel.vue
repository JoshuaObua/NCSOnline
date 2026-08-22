<template>
  <div class="manage-officials-panel">
    <!-- Header bar -->
    <div class="panel-top-bar">
      <div>
        <h2 class="panel-title"><i class="icofont-users-social text-primary"></i> Federation Officials Registry</h2>
        <p class="panel-subtitle">Manage elected leaders, executive committee members, and officers across all registered sports federations.</p>
      </div>
      <div class="header-actions">
        <button type="button" class="btn-new-official" @click="$emit('create-official')">
          <i class="icofont-plus"></i> Add Federation Official
        </button>
      </div>
    </div>

    <!-- Filter & Search Toolbar -->
    <div class="table-toolbar">
      <div class="search-box-wrapper">
        <i class="icofont-search search-icon"></i>
        <input
          v-model="searchQuery"
          type="text"
          class="form-control official-search-input"
          placeholder="Search officials by name, role, email, phone, or federation..."
        />
        <button v-if="searchQuery" type="button" class="btn-clear-search" @click="searchQuery = ''">
          <i class="icofont-close-line"></i>
        </button>
      </div>

      <div class="filter-group">
        <select v-model="selectedFederationFilter" class="form-control filter-select">
          <option value="">All Federations ({{ federations.length }})</option>
          <option v-for="fed in federations" :key="fed.id" :value="fed.id">
            {{ fed.name }} ({{ (fed.abbreviation || fed.acronym || fed.id).toUpperCase() }})
          </option>
        </select>
      </div>

      <div class="filter-group">
        <select v-model="selectedPositionFilter" class="form-control filter-select">
          <option value="">All Positions</option>
          <option value="PRESIDENT">President / Chairperson</option>
          <option value="VICE_PRESIDENT">Vice President</option>
          <option value="GENERAL_SECRETARY">General Secretary</option>
          <option value="TREASURER">Treasurer</option>
          <option value="TECHNICAL_DIRECTOR">Technical Director</option>
          <option value="COMMITTEE_MEMBER">Committee Member</option>
          <option value="OTHER">Other Roles</option>
        </select>
      </div>
    </div>

    <!-- Officials Table Card -->
    <div class="card table-card">
      <div class="table-responsive">
        <table class="table officials-data-table">
          <thead>
            <tr>
              <th>Federation</th>
              <th>Official Name</th>
              <th>Designated Position</th>
              <th>Contact Details</th>
              <th>Tenure / Term</th>
              <th>Status</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="7" class="text-center py-5">
                <i class="icofont-spinner icofont-spin font-24 text-primary"></i>
                <div class="mt-2 text-muted">Loading federation officials from database...</div>
              </td>
            </tr>
            <tr v-else-if="filteredItems.length === 0">
              <td colspan="7" class="text-center py-5 text-muted">
                <i class="icofont-user-alt-7 font-36 mb-2 d-block text-muted"></i>
                <p>No federation officials found matching your filter criteria.</p>
                <button type="button" class="btn btn-outline-primary btn-sm mt-2" @click="$emit('create-official')">
                  Add First Official
                </button>
              </td>
            </tr>
            <tr v-for="off in filteredItems" :key="off.id" class="official-row">
              <td>
                <div class="fed-cell">
                  <span class="badge-acronym" :title="getFederationName(off.federation_id)">
                    {{ getFederationAcronym(off.federation_id) }}
                  </span>
                  <a
                    href="javascript:void(0)"
                    class="fed-link-name"
                    @click="$emit('view-federation-profile', off.federation_id)"
                  >
                    {{ getFederationName(off.federation_id) }}
                  </a>
                </div>
              </td>
              <td>
                <div class="official-name-cell">
                  <div class="official-avatar">
                    <i class="icofont-user-suited"></i>
                  </div>
                  <div class="official-name-meta">
                    <span class="official-full-name">{{ off.full_name }}</span>
                    <span v-if="off.nin" class="official-nin-badge">
                      <i class="icofont-id-card"></i> {{ off.nin }}
                    </span>
                  </div>
                </div>
              </td>
              <td>
                <span class="badge-position" :class="positionClass(off.position)">
                  {{ off.position_label || off.position }}
                </span>
              </td>
              <td>
                <div class="contact-stack">
                  <div v-if="off.email" class="contact-item">
                    <i class="icofont-email text-muted"></i> <span>{{ off.email }}</span>
                  </div>
                  <div v-if="off.phone" class="contact-item">
                    <i class="icofont-phone text-muted"></i> <span>{{ off.phone }}</span>
                  </div>
                  <span v-if="!off.email && !off.phone" class="text-muted small">None recorded</span>
                </div>
              </td>
              <td>
                <div class="tenure-stack small">
                  <span v-if="off.appointed_on">From: {{ formatDate(off.appointed_on) }}</span>
                  <span v-if="off.term_ends_on">To: {{ formatDate(off.term_ends_on) }}</span>
                  <span v-if="!off.appointed_on && !off.term_ends_on" class="text-muted">Permanent / Indefinite</span>
                </div>
              </td>
              <td>
                <span class="badge-status" :class="off.is_active !== false ? 'success' : 'danger'">
                  {{ off.is_active !== false ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="text-right">
                <div class="action-buttons-group">
                  <button
                    type="button"
                    class="btn-action btn-action-profile"
                    title="View Federation Profile"
                    @click="$emit('view-federation-profile', off.federation_id)"
                  >
                    <i class="icofont-trophy"></i> Federation
                  </button>
                  <button
                    type="button"
                    class="btn-action btn-action-edit"
                    title="Edit Official Details"
                    @click="$emit('edit-official', off)"
                  >
                    <i class="icofont-edit"></i> Edit
                  </button>
                  <button
                    type="button"
                    class="btn-action btn-action-delete"
                    title="Delete Official"
                    @click="$emit('delete-official', off)"
                  >
                    <i class="icofont-trash"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="card-footer d-flex justify-content-between align-items-center">
        <span class="text-muted small">Showing {{ filteredItems.length }} of {{ items.length }} officials</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => [],
  },
  federations: {
    type: Array,
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['create-official', 'edit-official', 'delete-official', 'view-federation-profile'])

const searchQuery = ref('')
const selectedFederationFilter = ref('')
const selectedPositionFilter = ref('')

const federationMap = computed(() => {
  const map = {}
  for (const fed of props.federations || []) {
    map[fed.id] = fed
  }
  return map
})

function getFederationName(fedId) {
  return federationMap.value[fedId]?.name || fedId || 'Sports Federation'
}

function getFederationAcronym(fedId) {
  const fed = federationMap.value[fedId]
  if (fed) return (fed.abbreviation || fed.acronym || fed.id).toUpperCase().replace(/^ASSOC_/, '')
  return String(fedId || 'FED').toUpperCase().replace(/^ASSOC_/, '')
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return String(dateStr).slice(0, 10)
}

function positionClass(pos) {
  const p = String(pos || '').toUpperCase()
  if (p === 'PRESIDENT') return 'pos-president'
  if (p === 'VICE_PRESIDENT') return 'pos-vp'
  if (p === 'GENERAL_SECRETARY') return 'pos-sg'
  if (p === 'TREASURER') return 'pos-treasurer'
  return 'pos-default'
}

const filteredItems = computed(() => {
  let list = props.items || []
  if (selectedFederationFilter.value) {
    list = list.filter(item => item.federation_id === selectedFederationFilter.value)
  }
  if (selectedPositionFilter.value) {
    list = list.filter(item => item.position === selectedPositionFilter.value)
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(item => {
      const fedName = getFederationName(item.federation_id).toLowerCase()
      const fedAcronym = getFederationAcronym(item.federation_id).toLowerCase()
      return (
        String(item.full_name || '').toLowerCase().includes(q) ||
        String(item.nin || '').toLowerCase().includes(q) ||
        String(item.position || '').toLowerCase().includes(q) ||
        String(item.position_label || '').toLowerCase().includes(q) ||
        String(item.email || '').toLowerCase().includes(q) ||
        String(item.phone || '').toLowerCase().includes(q) ||
        fedName.includes(q) ||
        fedAcronym.includes(q)
      )
    })
  }
  return list
})
</script>

<style scoped>
.manage-officials-panel {
  padding: 0 4px;
}

.panel-top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}

.panel-title {
  font-size: 20px;
  font-weight: 800;
  color: #1e293b;
  margin: 0 0 4px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.panel-subtitle {
  font-size: 13px;
  color: #64748b;
  margin: 0;
}

.btn-new-official {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 9px 18px;
  border-radius: 8px;
  background: #6777ef;
  border: 1px solid #6777ef;
  color: #ffffff;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(103, 119, 239, 0.3);
  transition: all 0.15s ease;
}

.btn-new-official:hover {
  background: #5566de;
}

/* Toolbar */
.table-toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.search-box-wrapper {
  position: relative;
  flex: 1;
  min-width: 260px;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #94a3b8;
  font-size: 15px;
}

.official-search-input {
  padding-left: 36px;
  padding-right: 32px;
  height: 40px;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  font-size: 13px;
}

.btn-clear-search {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
}

.filter-select {
  height: 40px;
  border-radius: 8px;
  border: 1px solid #cbd5e1;
  font-size: 13px;
  min-width: 170px;
}

/* Table */
.table-card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
}

.officials-data-table {
  width: 100%;
  margin: 0;
  border-collapse: collapse;
}

.officials-data-table th {
  background: #f8fafc;
  padding: 12px 16px;
  font-size: 12px;
  font-weight: 700;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
  white-space: nowrap;
}

.officials-data-table td {
  padding: 12px 16px;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}

.fed-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  max-width: 220px;
}

.badge-acronym {
  padding: 2px 6px;
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 800;
  flex-shrink: 0;
}

.fed-link-name {
  font-size: 12px;
  font-weight: 600;
  color: #334155;
  text-decoration: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fed-link-name:hover {
  color: #6777ef;
  text-decoration: underline;
}

.official-name-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.official-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: #f1f5f9;
  color: #475569;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}

.official-full-name {
  font-size: 13px;
  font-weight: 700;
  color: #1e293b;
}

.badge-position {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
}

.badge-position.pos-president {
  background: #faf5ff;
  color: #7e22ce;
  border: 1px solid #e9d5ff;
}

.badge-position.pos-vp {
  background: #eff6ff;
  color: #1d4ed8;
  border: 1px solid #bfdbfe;
}

.badge-position.pos-sg {
  background: #ecfdf5;
  color: #047857;
  border: 1px solid #a7f3d0;
}

.badge-position.pos-treasurer {
  background: #fffbeb;
  color: #b45309;
  border: 1px solid #fde68a;
}

.badge-position.pos-default {
  background: #f8fafc;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.contact-stack,
.tenure-stack {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 12px;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #475569;
}

.badge-status {
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 800;
  text-transform: uppercase;
}

.badge-status.success {
  background: #ecfdf5;
  color: #059669;
}

.badge-status.danger {
  background: #fef2f2;
  color: #dc2626;
}

/* Action buttons */
.action-buttons-group {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.btn-action {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 5px 10px;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-action-profile {
  background: #eff6ff !important;
  border: 1px solid #93c5fd !important;
  color: #1d4ed8 !important;
  font-weight: 700 !important;
}

.btn-action-profile:hover {
  background: #2563eb !important;
  color: #ffffff !important;
}

.btn-action-edit {
  background: #f8fafc !important;
  border: 1px solid #cbd5e1 !important;
  color: #334155 !important;
  font-weight: 700 !important;
}

.btn-action-edit:hover {
  background: #e2e8f0 !important;
  color: #0f172a !important;
}

.btn-action-delete {
  background: #fff1f2 !important;
  border: 1px solid #fecdd3 !important;
  color: #e11d48 !important;
  font-weight: 700 !important;
}

.btn-action-delete:hover {
  background: #e11d48 !important;
  color: #ffffff !important;
}
</style>

<style scoped>
.official-name-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.official-nin-badge {
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
  width: fit-content;
}
:global(body.dark-theme) .official-nin-badge {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
  border-color: rgba(16, 185, 129, 0.35);
}
</style>
