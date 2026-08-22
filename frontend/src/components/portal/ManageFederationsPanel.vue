<template>
  <div class="manage-feds-panel">
    <!-- Header bar -->
    <div class="panel-top-bar">
      <div>
        <h2 class="panel-title"><i class="icofont-trophy text-primary"></i> Sports Federations & Associations</h2>
        <p class="panel-subtitle">Manage recognized sports federations, override IDs, update leadership, and view complete federation profiles.</p>
      </div>
      <div class="header-actions">
        <button type="button" class="btn-new-fed" @click="$emit('create-federation')">
          <i class="icofont-plus"></i> Add New Federation
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
          class="form-control fed-search-input"
          placeholder="Search federations by name, acronym, ID, president, or category..."
        />
        <button v-if="searchQuery" type="button" class="btn-clear-search" @click="searchQuery = ''">
          <i class="icofont-close-line"></i>
        </button>
      </div>

      <div class="filter-group">
        <select v-model="selectedCategory" class="form-control filter-select">
          <option value="">All Categories ({{ items.length }})</option>
          <option v-for="cat in categories" :key="cat.slug || cat.name || cat" :value="cat.name || cat">
            {{ cat.name || cat }}
          </option>
        </select>
      </div>
    </div>

    <!-- Federations Table Card -->
    <div class="card table-card">
      <div class="table-responsive">
        <table class="table fed-data-table">
          <thead>
            <tr>
              <th class="th-logo">Logo</th>
              <th>Federation Name & ID</th>
              <th>Category</th>
              <th>Leadership</th>
              <th>NCS Reg Number</th>
              <th>Status</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="7" class="text-center py-5">
                <i class="icofont-spinner icofont-spin font-24 text-primary"></i>
                <div class="mt-2 text-muted">Loading federations...</div>
              </td>
            </tr>
            <tr v-else-if="filteredItems.length === 0">
              <td colspan="7" class="text-center py-5 text-muted">
                <i class="icofont-folder-open font-36 mb-2 d-block text-muted"></i>
                <p>No federations match your filter criteria.</p>
              </td>
            </tr>
            <tr v-for="fed in filteredItems" :key="fed.id" class="fed-row">
              <td class="td-logo">
                <div class="logo-circle">
                  <img v-if="fed.logo_url" :src="fed.logo_url" :alt="fed.name" class="table-logo-img" />
                  <span v-else class="logo-initials">{{ (fed.abbreviation || fed.acronym || fed.name.slice(0, 3)).toUpperCase() }}</span>
                </div>
              </td>
              <td>
                <div class="fed-name-cell">
                  <a href="javascript:void(0)" class="fed-primary-name" @click="$emit('view-profile', fed.id, fed)">
                    {{ fed.name }}
                  </a>
                  <div class="fed-id-line">
                    <span v-if="fed.abbreviation || fed.acronym" class="badge-acronym-mini">
                      {{ fed.abbreviation || fed.acronym }}
                    </span>
                    <code class="fed-id-code" :title="'ID: ' + fed.id">{{ fed.id }}</code>
                  </div>
                </div>
              </td>
              <td>
                <span class="badge-category">{{ fed.category || 'Other' }}</span>
              </td>
              <td>
                <div class="leader-info">
                  <div v-if="fed.president" class="leader-line">
                    <i class="icofont-user-suited text-muted"></i> <small>Pres:</small> <strong>{{ fed.president }}</strong>
                  </div>
                  <div v-if="fed.secretary" class="leader-line">
                    <i class="icofont-business-man text-muted"></i> <small>Sec:</small> <span>{{ fed.secretary }}</span>
                  </div>
                  <span v-if="!fed.president && !fed.secretary" class="text-muted small">Not recorded</span>
                </div>
              </td>
              <td>
                <span class="badge-reg">{{ fed.ncs_registration_number || ('NCS-REG-' + (fed.abbreviation || fed.acronym || fed.id).toUpperCase().replace(/^ASSOC_/, '')) }}</span>
              </td>
              <td>
                <span class="badge-status" :class="fed.is_active !== false ? 'success' : 'danger'">
                  {{ fed.is_active !== false ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="text-right">
                <div class="action-buttons-group">
                  <button
                    type="button"
                    class="btn-action btn-action-profile"
                    title="View Federation Profile (General, Licenses, Team Members, Records)"
                    @click="$emit('view-profile', fed.id, fed)"
                  >
                    <i class="icofont-id-card"></i> Profile
                  </button>
                  <button
                    type="button"
                    class="btn-action btn-action-edit"
                    title="Edit Federation Details & Override ID"
                    @click="$emit('edit-federation', fed)"
                  >
                    <i class="icofont-edit"></i> Edit
                  </button>
                  <button
                    type="button"
                    class="btn-action btn-action-delete"
                    title="Delete Federation"
                    @click="$emit('delete-federation', fed)"
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
        <span class="text-muted small">Showing {{ filteredItems.length }} of {{ items.length }} federations</span>
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
  categories: {
    type: Array,
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['create-federation', 'edit-federation', 'delete-federation', 'view-profile'])

const searchQuery = ref('')
const selectedCategory = ref('')

const filteredItems = computed(() => {
  let list = props.items || []
  if (selectedCategory.value) {
    list = list.filter(item => item.category === selectedCategory.value)
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(item => {
      return (
        String(item.name || '').toLowerCase().includes(q) ||
        String(item.abbreviation || item.acronym || '').toLowerCase().includes(q) ||
        String(item.id || '').toLowerCase().includes(q) ||
        String(item.president || '').toLowerCase().includes(q) ||
        String(item.secretary || '').toLowerCase().includes(q) ||
        String(item.category || '').toLowerCase().includes(q) ||
        String(item.ncs_registration_number || '').toLowerCase().includes(q)
      )
    })
  }
  return list
})
</script>

<style scoped>
.manage-feds-panel {
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

.btn-new-fed {
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

.btn-new-fed:hover {
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
  min-width: 280px;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #94a3b8;
  font-size: 15px;
}

.fed-search-input {
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
  min-width: 180px;
}

/* Table */
.table-card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
}

.fed-data-table {
  width: 100%;
  margin: 0;
  border-collapse: collapse;
}

.fed-data-table th {
  background: #f8fafc;
  padding: 12px 16px;
  font-size: 12px;
  font-weight: 700;
  color: #475569;
  border-bottom: 1px solid #e2e8f0;
  white-space: nowrap;
}

.th-logo {
  width: 50px;
  text-align: center;
}

.fed-data-table td {
  padding: 12px 16px;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}

.td-logo {
  text-align: center;
}

.logo-circle {
  width: 38px;
  height: 38px;
  border-radius: 8px;
  background: #eff6ff;
  border: 1px solid #dbeafe;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.table-logo-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.logo-initials {
  font-size: 12px;
  font-weight: 800;
  color: #2563eb;
}

.fed-name-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.fed-primary-name {
  font-size: 13px;
  font-weight: 700;
  color: #1e293b;
  text-decoration: none;
}

.fed-primary-name:hover {
  color: #6777ef;
  text-decoration: underline;
}

.fed-id-line {
  display: flex;
  align-items: center;
  gap: 6px;
}

.badge-acronym-mini {
  padding: 1px 5px;
  background: #e0e7ff;
  color: #4338ca;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 800;
}

.fed-id-code {
  font-size: 11px;
  color: #64748b;
  background: #f1f5f9;
  padding: 1px 4px;
  border-radius: 4px;
}

.badge-category {
  padding: 4px 8px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  color: #334155;
}

.leader-info {
  font-size: 12px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.leader-line {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #475569;
}

.badge-reg {
  font-size: 11px;
  font-weight: 700;
  color: #0284c7;
  background: #e0f2fe;
  padding: 3px 6px;
  border-radius: 4px;
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
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  border: 1px solid transparent;
}

.btn-action-profile {
  background: #eff6ff;
  border-color: #bfdbfe;
  color: #2563eb;
}

.btn-action-profile:hover {
  background: #2563eb;
  color: #ffffff;
}

.btn-action-edit {
  background: #f1f5f9;
  border-color: #cbd5e1;
  color: #334155;
}

.btn-action-edit:hover {
  background: #e2e8f0;
}

.btn-action-delete {
  background: #fff1f2;
  border-color: #fecdd3;
  color: #e11d48;
}

.btn-action-delete:hover {
  background: #e11d48;
  color: #ffffff;
}
</style>
