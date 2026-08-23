<template>
  <div class="assign-role-page">
    <!-- Breadcrumb & Top Bar -->
    <nav class="page-nav">
      <button type="button" class="back-link-btn" @click="handleBack">
        <i class="icofont-arrow-left"></i> Back to Manage Users
      </button>
      <div class="breadcrumb-trail">
        <span>Portal</span>
        <i class="icofont-simple-right"></i>
        <button type="button" class="trail-link" @click="handleBack">Users</button>
        <i class="icofont-simple-right"></i>
        <span class="active-trail">Assign Roles</span>
      </div>
    </nav>

    <!-- Page Header -->
    <header class="page-header">
      <div class="header-content">
        <p class="section-tag">Access Control & Role-Based Permissions</p>
        <h2 class="page-title">
          <i class="icofont-shield-alt"></i> Manage User Roles
        </h2>
        <p class="page-desc">
          Assign, reassign, or revoke administrative and sports registry roles for this user. Roles control portal modules, permissions, and operational capabilities across the system.
        </p>
      </div>
      <div class="header-actions">
        <button type="button" class="btn btn-outline" :disabled="loading" @click="refreshData">
          <i class="icofont-refresh" :class="{ 'icofont-spin': loading }"></i> Refresh
        </button>
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

    <!-- User Profile Header Card -->
    <div v-if="user" class="user-profile-banner">
      <div class="user-main-info">
        <div class="user-avatar-large" :class="avatarColorClass">
          {{ userInitials }}
        </div>
        <div class="user-details-col">
          <div class="user-title-row">
            <h3 class="user-name">{{ user.first_name }} {{ user.last_name }}</h3>
            <span class="status-badge" :class="user.is_active ? 'status-active' : 'status-suspended'">
              <span class="status-dot"></span>
              {{ user.account_status || (user.is_active ? 'ACTIVE' : 'SUSPENDED') }}
            </span>
          </div>
          <div class="user-meta-row">
            <span class="meta-item"><i class="icofont-ui-email"></i> {{ user.email }}</span>
            <span v-if="user.phone" class="meta-item"><i class="icofont-ui-touch-phone"></i> {{ user.phone }}</span>
            <span v-if="user.nin" class="meta-item nin-pill"><i class="icofont-id-card"></i> NIN: {{ user.nin }}</span>
          </div>
        </div>
      </div>

      <!-- Current Assigned Roles Summary -->
      <div class="current-roles-box">
        <div class="roles-box-header">
          <span class="box-label">Currently Assigned Roles ({{ assignedRoleNames.length }})</span>
        </div>
        <div v-if="assignedRoleNames.length" class="assigned-roles-wrap">
          <span
            v-for="rName in assignedRoleNames"
            :key="rName"
            class="assigned-badge"
            :class="getRoleBadgeClass(rName)"
          >
            <i class="icofont-shield"></i> {{ rName }}
            <button
              type="button"
              class="revoke-badge-btn"
              :disabled="savingRole === rName"
              title="Revoke this role"
              @click="revokeRoleByName(rName)"
            >
              <i v-if="savingRole === rName" class="icofont-spinner icofont-spin"></i>
              <i v-else class="icofont-close-line"></i>
            </button>
          </span>
        </div>
        <div v-else class="no-assigned-state">
          <i class="icofont-info-circle"></i> No roles currently assigned to this account.
        </div>
      </div>
    </div>

    <!-- Toolbar: Search & Category Filter -->
    <div class="role-toolbar">
      <div class="search-box">
        <i class="icofont-search"></i>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search roles by title, keyword, or capability..."
        />
        <button v-if="searchQuery" type="button" class="clear-search-btn" @click="searchQuery = ''">
          <i class="icofont-close-line"></i>
        </button>
      </div>

      <div class="filter-pills">
        <button
          type="button"
          class="filter-pill"
          :class="{ active: roleFilter === 'ALL' }"
          @click="roleFilter = 'ALL'"
        >
          All ({{ roles.length }})
        </button>
        <button
          type="button"
          class="filter-pill"
          :class="{ active: roleFilter === 'ASSIGNED' }"
          @click="roleFilter = 'ASSIGNED'"
        >
          Assigned ({{ assignedRoleNames.length }})
        </button>
        <button
          type="button"
          class="filter-pill"
          :class="{ active: roleFilter === 'UNASSIGNED' }"
          @click="roleFilter = 'UNASSIGNED'"
        >
          Unassigned ({{ unassignedRolesCount }})
        </button>
        <button
          type="button"
          class="filter-pill"
          :class="{ active: roleFilter === 'SYSTEM' }"
          @click="roleFilter = 'SYSTEM'"
        >
          System Roles
        </button>
        <button
          type="button"
          class="filter-pill"
          :class="{ active: roleFilter === 'REGISTRY' }"
          @click="roleFilter = 'REGISTRY'"
        >
          Sports Registry
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading && !roles.length" class="loading-container">
      <i class="icofont-spinner icofont-spin"></i> Loading roles from database...
    </div>

    <!-- Empty State -->
    <div v-else-if="!filteredRoles.length" class="empty-roles-card">
      <i class="icofont-search-folder empty-icon"></i>
      <h3>No matching roles found</h3>
      <p>Try clearing your search query or adjusting your filters.</p>
      <button type="button" class="btn btn-outline" @click="searchQuery = ''; roleFilter = 'ALL'">
        Reset Filters
      </button>
    </div>

    <!-- Role Cards Grid -->
    <div v-else class="role-cards-grid">
      <div
        v-for="role in filteredRoles"
        :key="role.id"
        class="role-card"
        :class="{ 'card-assigned': isRoleAssigned(role) }"
      >
        <div class="role-card-top">
          <div class="role-icon-box" :class="isRoleAssigned(role) ? 'icon-box-assigned' : 'icon-box-unassigned'">
            <i class="icofont-shield-alt"></i>
          </div>
          <div class="role-header-info">
            <div class="role-name-row">
              <h4 class="role-name">{{ role.name }}</h4>
              <span v-if="role.is_system" class="pill-badge pill-system">System</span>
              <span v-else-if="isSportsRole(role.name)" class="pill-badge pill-registry">Registry</span>
              <span v-else class="pill-badge pill-custom">Custom</span>
            </div>
            <span
              class="assignment-indicator"
              :class="isRoleAssigned(role) ? 'indicator-assigned' : 'indicator-unassigned'"
            >
              <i :class="isRoleAssigned(role) ? 'icofont-check-circled' : 'icofont-circle'"></i>
              {{ isRoleAssigned(role) ? 'Assigned' : 'Not Assigned' }}
            </span>
          </div>
        </div>

        <p class="role-description">
          {{ role.description || 'Custom administrative role for managing portal features and resources.' }}
        </p>

        <!-- Role Permissions Preview -->
        <div v-if="(role.permissions || []).length" class="permissions-preview">
          <span class="perms-label">Granted Permissions ({{ role.permissions.length }}):</span>
          <div class="perms-chips">
            <span
              v-for="p in role.permissions.slice(0, 5)"
              :key="p.id || p.name"
              class="perm-chip"
            >
              {{ p.name || p }}
            </span>
            <span v-if="role.permissions.length > 5" class="perm-chip perm-more">
              +{{ role.permissions.length - 5 }} more
            </span>
          </div>
        </div>

        <!-- Role Action Buttons -->
        <div class="role-card-actions">
          <button
            v-if="isRoleAssigned(role)"
            type="button"
            class="btn-card-action btn-revoke"
            :disabled="savingRole === (role.id || role.name)"
            @click="handleRevoke(role)"
          >
            <i v-if="savingRole === (role.id || role.name)" class="icofont-spinner icofont-spin"></i>
            <span v-else><i class="icofont-minus-circle"></i> Revoke Role</span>
          </button>

          <template v-else>
            <button
              type="button"
              class="btn-card-action btn-assign"
              :disabled="savingRole === (role.id || role.name)"
              @click="handleAssign(role)"
            >
              <i v-if="savingRole === (role.id || role.name)" class="icofont-spinner icofont-spin"></i>
              <span v-else><i class="icofont-plus-circle"></i> Assign Role</span>
            </button>
            <button
              v-if="assignedRoleNames.length > 0"
              type="button"
              class="btn-card-action btn-reassign"
              :disabled="savingRole === (role.id || role.name)"
              title="Replace all existing roles with this role"
              @click="handleReassign(role)"
            >
              <i class="icofont-refresh"></i> Reassign Primary
            </button>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Swal from 'sweetalert2'
import * as cms from '@/api/cms.js'

const props = defineProps({
  userId: {
    type: String,
    default: ''
  },
  initialUser: {
    type: Object,
    default: () => null
  }
})

const emit = defineEmits(['back', 'message', 'error', 'updated'])
const route = useRoute()
const router = useRouter()

const user = ref(props.initialUser ? JSON.parse(JSON.stringify(props.initialUser)) : null)
const roles = ref([])
const loading = ref(false)
const savingRole = ref('')
const localError = ref('')
const localMessage = ref('')
const searchQuery = ref('')
const roleFilter = ref('ALL')

const targetUserId = computed(() => {
  return props.userId || route.query.id || (user.value ? user.value.id : '')
})

const sportsRoleNames = new Set([
  'club_manager', 'coach', 'athlete', 'technical_official',
  'medical_officer', 'anti_doping_officer', 'federation_officer'
])

function isSportsRole(name) {
  return sportsRoleNames.has(name)
}

const assignedRoleNames = computed(() => {
  if (!user.value || !user.value.roles) return []
  return user.value.roles.map(r => typeof r === 'object' ? (r.name || '') : String(r)).filter(Boolean)
})

const unassignedRolesCount = computed(() => {
  return roles.value.filter(r => !isRoleAssigned(r)).length
})

const userInitials = computed(() => {
  if (!user.value) return 'U'
  const f = (user.value.first_name || '').charAt(0).toUpperCase()
  const l = (user.value.last_name || '').charAt(0).toUpperCase()
  return `${f}${l}` || 'U'
})

const avatarColorClass = computed(() => {
  const colors = ['avatar-purple', 'avatar-blue', 'avatar-emerald', 'avatar-amber', 'avatar-rose']
  const str = user.value?.email || user.value?.id || ''
  let hash = 0
  for (let i = 0; i < str.length; i++) hash += str.charCodeAt(i)
  return colors[hash % colors.length]
})

function isRoleAssigned(role) {
  if (!user.value || !user.value.roles) return false
  return user.value.roles.some(r => {
    if (typeof r === 'object') {
      return (r.id && r.id === role.id) || (r.name && r.name === role.name)
    }
    return r === role.name || r === role.id
  })
}

const filteredRoles = computed(() => {
  const q = (searchQuery.value || '').trim().toLowerCase()
  const filter = roleFilter.value

  return roles.value.filter(role => {
    // 1. Category Filter
    if (filter === 'ASSIGNED' && !isRoleAssigned(role)) return false
    if (filter === 'UNASSIGNED' && isRoleAssigned(role)) return false
    if (filter === 'SYSTEM' && !role.is_system) return false
    if (filter === 'REGISTRY' && !isSportsRole(role.name)) return false

    // 2. Text Search
    if (q) {
      const name = (role.name || '').toLowerCase()
      const desc = (role.description || '').toLowerCase()
      const perms = (role.permissions || []).map(p => (p.name || p).toLowerCase()).join(' ')
      if (!name.includes(q) && !desc.includes(q) && !perms.includes(q)) {
        return false
      }
    }
    return true
  })
})

onMounted(async () => {
  await refreshData()
})

watch(() => route.query.id, async (newId) => {
  if (newId && (!user.value || user.value.id !== newId)) {
    await fetchUserDetails(newId)
  }
})

async function refreshData() {
  loading.value = true
  localError.value = ''
  try {
    await Promise.all([
      loadRoles(),
      targetUserId.value ? fetchUserDetails(targetUserId.value) : Promise.resolve()
    ])
  } catch (err) {
    localError.value = err.response?.data?.error?.message || err.message || 'Failed to load role assignment data.'
  } finally {
    loading.value = false
  }
}

async function loadRoles() {
  try {
    const res = await cms.adminListRoles()
    const list = res.data?.data || res.data || []
    roles.value = Array.isArray(list) ? list : []
  } catch (e) {
    console.error('Could not load roles:', e)
  }
}

async function fetchUserDetails(id) {
  if (!id) return
  try {
    const res = await cms.adminGetUser(id)
    const data = res.data?.data || res.data
    if (data) {
      user.value = data
    }
  } catch (e) {
    console.error('Could not fetch user details:', e)
    localError.value = e.response?.data?.error?.message || 'Could not load user details.'
  }
}

function handleBack() {
  emit('back')
}

function notifySuccess(msg) {
  localMessage.value = msg
  localError.value = ''
  emit('message', msg)
  setTimeout(() => { localMessage.value = '' }, 4000)
}

async function handleAssign(role) {
  if (!user.value) return
  const roleId = role.id || role.name
  savingRole.value = role.id || role.name
  localError.value = ''
  try {
    await cms.adminAssignUserRole(user.value.id, roleId)
    notifySuccess(`Role '${role.name}' successfully assigned to ${user.value.email}.`)
    await fetchUserDetails(user.value.id)
    emit('updated', user.value)
  } catch (err) {
    const errorMsg = err.response?.data?.error?.message || err.message || `Failed to assign role '${role.name}'.`
    localError.value = errorMsg
    emit('error', errorMsg)
  } finally {
    savingRole.value = ''
  }
}

async function handleRevoke(role) {
  if (!user.value) return
  const confirmResult = await Swal.fire({
    title: `Revoke role '${role.name}'?`,
    text: `Are you sure you want to remove the '${role.name}' role from ${user.value.email}?`,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Yes, Revoke Role',
    cancelButtonText: 'Cancel',
    confirmButtonColor: '#fc544b',
    cancelButtonColor: '#6c757d',
    reverseButtons: true,
  })
  if (!confirmResult.isConfirmed) return

  const roleId = role.id || role.name
  savingRole.value = role.id || role.name
  localError.value = ''
  try {
    await cms.adminRemoveUserRole(user.value.id, roleId)
    notifySuccess(`Role '${role.name}' successfully revoked from ${user.value.email}.`)
    await fetchUserDetails(user.value.id)
    emit('updated', user.value)
  } catch (err) {
    const errorMsg = err.response?.data?.error?.message || err.message || `Failed to revoke role '${role.name}'.`
    localError.value = errorMsg
    emit('error', errorMsg)
  } finally {
    savingRole.value = ''
  }
}

async function revokeRoleByName(roleName) {
  const fullRole = roles.value.find(r => r.name === roleName) || { id: roleName, name: roleName }
  await handleRevoke(fullRole)
}

async function handleReassign(role) {
  if (!user.value) return
  const confirmResult = await Swal.fire({
    title: `Reassign to '${role.name}'?`,
    text: `This will revoke all other current roles (${assignedRoleNames.value.join(', ')}) and set '${role.name}' as the sole role for ${user.value.email}.`,
    icon: 'question',
    showCancelButton: true,
    confirmButtonText: 'Yes, Reassign Role',
    cancelButtonText: 'Cancel',
    confirmButtonColor: '#f59e0b',
    cancelButtonColor: '#6c757d',
  })
  if (!confirmResult.isConfirmed) return

  const roleId = role.id || role.name
  savingRole.value = role.id || role.name
  localError.value = ''
  try {
    // Revoke previous roles
    const prevRoles = [...(user.value.roles || [])]
    for (const r of prevRoles) {
      const rId = typeof r === 'object' ? (r.id || r.name) : r
      if (rId && rId !== roleId) {
        try { await cms.adminRemoveUserRole(user.value.id, rId) } catch (e) { console.warn('Could not revoke previous role:', e) }
      }
    }
    // Assign new role
    await cms.adminAssignUserRole(user.value.id, roleId)
    notifySuccess(`Roles for ${user.value.email} successfully reassigned to '${role.name}'.`)
    await fetchUserDetails(user.value.id)
    emit('updated', user.value)
  } catch (err) {
    const errorMsg = err.response?.data?.error?.message || err.message || `Failed to reassign role to '${role.name}'.`
    localError.value = errorMsg
    emit('error', errorMsg)
  } finally {
    savingRole.value = ''
  }
}

function getRoleBadgeClass(roleName) {
  const r = String(roleName || '').toLowerCase()
  if (r.includes('super_admin')) return 'badge-superadmin'
  if (r.includes('admin')) return 'badge-admin'
  if (r.includes('operator')) return 'badge-operator'
  if (r.includes('federation')) return 'badge-federation'
  return 'badge-user'
}
</script>

<style scoped>
.assign-role-page {
  display: flex;
  flex-direction: column;
  gap: 22px;
  width: 100%;
}

/* Breadcrumb Navigation */
.page-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.back-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: transparent;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  color: #475569;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.back-link-btn:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.breadcrumb-trail {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #64748b;
}

.trail-link {
  background: transparent;
  border: none;
  color: #6777ef;
  font-weight: 600;
  cursor: pointer;
  padding: 0;
  font-size: inherit;
}

.trail-link:hover {
  text-decoration: underline;
}

.breadcrumb-trail .active-trail {
  color: #1e293b;
  font-weight: 700;
}

/* Page Header */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  padding-bottom: 18px;
  border-bottom: 1px solid #e2e8f0;
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
  display: flex;
  align-items: center;
  gap: 10px;
}

.page-desc {
  color: #64748b;
  font-size: 13px;
  margin: 0;
  max-width: 820px;
  line-height: 1.5;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* User Profile Banner */
.user-profile-banner {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 22px 26px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.user-main-info {
  display: flex;
  align-items: center;
  gap: 18px;
}

.user-avatar-large {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 800;
  color: #ffffff;
  flex-shrink: 0;
  letter-spacing: 0.05em;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.avatar-purple { background: linear-gradient(135deg, #7c3aed, #a855f7); }
.avatar-blue { background: linear-gradient(135deg, #2563eb, #38bdf8); }
.avatar-emerald { background: linear-gradient(135deg, #059669, #34d399); }
.avatar-amber { background: linear-gradient(135deg, #d97706, #fbbf24); }
.avatar-rose { background: linear-gradient(135deg, #e11d48, #fb7185); }

.user-details-col {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-title-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-name {
  margin: 0;
  font-size: 20px;
  font-weight: 800;
  color: #1e293b;
}

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

.user-meta-row {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  font-size: 13px;
  color: #64748b;
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.nin-pill {
  background: rgba(16, 185, 129, 0.1);
  color: #059669;
  font-family: monospace;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 6px;
  border: 1px solid rgba(16, 185, 129, 0.2);
}

/* Current Roles Box */
.current-roles-box {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 14px 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.roles-box-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.box-label {
  font-size: 12px;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #475569;
}

.assigned-roles-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.assigned-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 700;
}

.badge-superadmin { background: #ede9fe; color: #5b21b6; border: 1px solid #c4b5fd; }
.badge-admin { background: #eff6ff; color: #1d4ed8; border: 1px solid #bfdbfe; }
.badge-operator { background: #ecfdf5; color: #047857; border: 1px solid #a7f3d0; }
.badge-federation { background: #fff7ed; color: #c2410c; border: 1px solid #fed7aa; }
.badge-user { background: #f1f5f9; color: #334155; border: 1px solid #cbd5e1; }

.revoke-badge-btn {
  border: none;
  background: transparent;
  color: inherit;
  opacity: 0.7;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  font-size: 13px;
  transition: opacity 0.15s;
}

.revoke-badge-btn:hover:not(:disabled) {
  opacity: 1;
  color: #dc2626;
}

.no-assigned-state {
  font-size: 13px;
  color: #94a3b8;
  font-style: italic;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* Toolbar */
.role-toolbar {
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

.filter-pills {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.filter-pill {
  padding: 8px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 20px;
  background: #ffffff;
  color: #475569;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}

.filter-pill:hover {
  background: #f8fafc;
  border-color: #94a3b8;
}

.filter-pill.active {
  background: #6777ef;
  color: #ffffff;
  border-color: #6777ef;
}

/* Role Cards Grid */
.role-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 18px;
}

.role-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.03);
  display: flex;
  flex-direction: column;
  gap: 14px;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.role-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
}

.role-card.card-assigned {
  border-color: #c4b5fd;
  background: #fdfcff;
  box-shadow: 0 4px 20px rgba(124, 58, 237, 0.08);
}

.role-card-top {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.role-icon-box {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.icon-box-assigned {
  background: #ede9fe;
  color: #7c3aed;
}

.icon-box-unassigned {
  background: #f1f5f9;
  color: #64748b;
}

.role-header-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
  flex: 1;
}

.role-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.role-name {
  margin: 0;
  font-size: 15px;
  font-weight: 800;
  color: #1e293b;
}

.pill-badge {
  font-size: 10px;
  font-weight: 800;
  text-transform: uppercase;
  padding: 1px 6px;
  border-radius: 4px;
  letter-spacing: 0.04em;
}

.pill-system { background: #e2e8f0; color: #475569; }
.pill-registry { background: #dbeafe; color: #1d4ed8; }
.pill-custom { background: #fef3c7; color: #92400e; }

.assignment-indicator {
  font-size: 11px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.indicator-assigned { color: #059669; }
.indicator-unassigned { color: #94a3b8; }

.role-description {
  margin: 0;
  font-size: 13px;
  color: #64748b;
  line-height: 1.4;
  flex: 1;
}

/* Permissions preview */
.permissions-preview {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 10px 12px;
  background: #f8fafc;
  border-radius: 8px;
  border: 1px solid #f1f5f9;
}

.perms-label {
  font-size: 11px;
  font-weight: 700;
  color: #64748b;
}

.perms-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.perm-chip {
  font-size: 10px;
  font-family: monospace;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  padding: 1px 5px;
  border-radius: 4px;
  color: #475569;
}

.perm-more {
  background: #e2e8f0;
  font-weight: 700;
  color: #475569;
}

/* Role Card Actions */
.role-card-actions {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}

.btn-card-action {
  flex: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  font-family: inherit;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  border: 1px solid transparent;
}

.btn-assign {
  background: #6777ef;
  color: #ffffff;
}

.btn-assign:hover:not(:disabled) {
  background: #5563db;
}

.btn-revoke {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
}

.btn-revoke:hover:not(:disabled) {
  background: #dc2626;
  color: #ffffff;
}

.btn-reassign {
  background: #fffbeb;
  color: #d97706;
  border-color: #fde68a;
}

.btn-reassign:hover:not(:disabled) {
  background: #d97706;
  color: #ffffff;
}

.btn-card-action:disabled {
  opacity: 0.6;
  cursor: not-allowed;
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
  padding: 8px 16px;
  border-radius: 8px;
  font-family: inherit;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  border: 1px solid transparent;
}

.btn-outline { background: #ffffff; color: #475569; border-color: #cbd5e1; }
.btn-outline:hover:not(:disabled) { background: #f8fafc; border-color: #94a3b8; }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }

.loading-container {
  padding: 50px 20px;
  text-align: center;
  color: #64748b;
  font-size: 15px;
  font-weight: 600;
}

.empty-roles-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 50px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.empty-icon {
  font-size: 44px;
  color: #94a3b8;
  margin-bottom: 8px;
}

/* =====================================================================
   DARK MODE STYLES
   ===================================================================== */
:global(html.dark) .assign-role-page,
:global(body.dark) .assign-role-page,
:global([data-theme="dark"]) .assign-role-page {
  color: #e2e8f0 !important;
}

:global(html.dark) .page-title,
:global(body.dark) .page-title,
:global([data-theme="dark"]) .page-title {
  color: #f8fafc !important;
}

:global(html.dark) .page-header,
:global(body.dark) .page-header,
:global([data-theme="dark"]) .page-header {
  border-bottom-color: #334155 !important;
}

:global(html.dark) .user-profile-banner,
:global(body.dark) .user-profile-banner,
:global([data-theme="dark"]) .user-profile-banner {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(html.dark) .user-name,
:global(body.dark) .user-name,
:global([data-theme="dark"]) .user-name {
  color: #f8fafc !important;
}

:global(html.dark) .current-roles-box,
:global(body.dark) .current-roles-box,
:global([data-theme="dark"]) .current-roles-box {
  background-color: #0f172a !important;
  border-color: #334155 !important;
}

:global(html.dark) .role-card,
:global(body.dark) .role-card,
:global([data-theme="dark"]) .role-card {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(html.dark) .role-card.card-assigned,
:global(body.dark) .role-card.card-assigned,
:global([data-theme="dark"]) .role-card.card-assigned {
  background-color: #1e1b4b !important;
  border-color: #6366f1 !important;
}

:global(html.dark) .role-name,
:global(body.dark) .role-name,
:global([data-theme="dark"]) .role-name {
  color: #f8fafc !important;
}

:global(html.dark) .permissions-preview,
:global(body.dark) .permissions-preview,
:global([data-theme="dark"]) .permissions-preview {
  background-color: #0f172a !important;
  border-color: #334155 !important;
}

:global(html.dark) .perm-chip,
:global(body.dark) .perm-chip,
:global([data-theme="dark"]) .perm-chip {
  background-color: #1e293b !important;
  border-color: #334155 !important;
  color: #94a3b8 !important;
}

:global(html.dark) .search-box input,
:global(body.dark) .search-box input,
:global([data-theme="dark"]) .search-box input {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(html.dark) .filter-pill:not(.active),
:global(body.dark) .filter-pill:not(.active),
:global([data-theme="dark"]) .filter-pill:not(.active) {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #94a3b8 !important;
}

:global(html.dark) .btn-outline,
:global(body.dark) .btn-outline,
:global([data-theme="dark"]) .btn-outline {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #cbd5e1 !important;
}

:global(html.dark) .back-link-btn,
:global(body.dark) .back-link-btn,
:global([data-theme="dark"]) .back-link-btn {
  border-color: #334155 !important;
  color: #cbd5e1 !important;
}

@media (max-width: 768px) {
  .user-main-info { flex-direction: column; align-items: flex-start; }
  .role-toolbar { flex-direction: column; align-items: stretch; }
  .role-cards-grid { grid-template-columns: 1fr; }
}
</style>
