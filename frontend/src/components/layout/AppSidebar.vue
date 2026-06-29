<template>
  <aside
    :class="[
      'fixed inset-y-0 left-0 z-40 w-64 flex flex-col bg-primary-700 shadow-xl transition-transform duration-300',
      open ? 'translate-x-0' : '-translate-x-full md:translate-x-0'
    ]"
  >
    <!-- Brand -->
    <div class="flex items-center gap-3 px-4 py-4 border-b border-primary-600/60 flex-shrink-0">
      <div class="flex-shrink-0 w-11 h-11 bg-white rounded-xl flex items-center justify-center overflow-hidden p-1">
        <img src="/main-logo.png" alt="NCS Logo" class="w-full h-full object-contain" />
      </div>
      <div class="min-w-0">
        <div class="text-white font-bold text-sm leading-tight">National Council</div>
        <div class="text-primary-300 text-xs mt-0.5 leading-tight">of Sports · NCSMS</div>
      </div>
      <!-- Mobile close -->
      <button @click="$emit('close')" class="ml-auto md:hidden text-primary-300 hover:text-white p-1 rounded-lg">
        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Profile block -->
    <div class="px-4 py-3 border-b border-primary-600/40 flex-shrink-0">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-primary-600 border-2 border-primary-500 flex items-center justify-center flex-shrink-0">
          <span class="text-white text-sm font-bold">{{ userInitials }}</span>
        </div>
        <div class="min-w-0 flex-1">
          <div class="text-white text-sm font-semibold truncate leading-tight">{{ userName }}</div>
          <div class="text-primary-300 text-xs truncate mt-0.5">{{ userRole }}</div>
        </div>
        <router-link to="/profile" class="flex-shrink-0 text-primary-300 hover:text-white transition-colors" title="My Profile">
          <i class="icofont-settings text-base leading-none"></i>
        </router-link>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto pb-4">

      <!-- MAIN -->
      <div class="section-label">Main</div>
      <div class="px-2 space-y-0.5">
        <NavItem :to="'/dashboard'" label="Dashboard" icon="icofont-dashboard-web" />
      </div>

      <template v-if="authStore.canUseNSMIS">
        <div class="section-label">NSMIS Reporting</div>
        <div class="px-2 space-y-0.5">
          <NavItem v-if="authStore.hasAnyRole('super_admin','admin','ncs_general_secretary','general_secretary','technical_department','federation_president','federation_general_secretary','auditor')" :to="'/nsmis/governance'" label="Governance" icon="icofont-chart-histogram" />
          <NavItem :to="'/nsmis/reports'" label="Federation Reports" icon="icofont-file-document" />
          <NavItem v-if="authStore.hasAnyRole('super_admin','admin','ncs_general_secretary','general_secretary','technical_department')" :to="'/nsmis/insights/athletes'" label="Athletes" icon="icofont-runner-alt-1" />
          <NavItem v-if="authStore.hasAnyRole('super_admin','admin','ncs_general_secretary','general_secretary','technical_department')" :to="'/nsmis/insights/performance'" label="Performance" icon="icofont-medal" />
          <NavItem v-if="authStore.hasAnyRole('super_admin','admin','ncs_general_secretary','general_secretary','finance_department')" :to="'/nsmis/insights/finance'" label="Finance" icon="icofont-money" />
          <NavItem v-if="authStore.hasAnyRole('super_admin','admin','ncs_general_secretary','general_secretary','technical_department')" :to="'/nsmis/insights/talent'" label="Talent" icon="icofont-search-user" />
          <NavItem :to="'/nsmis/data/federation-officers'" label="Federation Profiles" icon="icofont-building-alt" />
          <NavItem :to="'/nsmis/data/athletes'" label="Athlete Registry" icon="icofont-users-alt-4" />
          <NavItem :to="'/nsmis/data/competitions'" label="Competitions" icon="icofont-trophy" />
          <NavItem :to="'/nsmis/data/medals'" label="Medal Records" icon="icofont-medal-sport" />
          <NavItem :to="'/nsmis/data/coaches'" label="Coaches" icon="icofont-teacher" />
          <NavItem :to="'/nsmis/data/technical-officials'" label="Technical Officials" icon="icofont-judge" />
          <NavItem :to="'/nsmis/data/equipment'" label="Equipment" icon="icofont-box" />
        </div>
      </template>

      <!-- MANAGEMENT -->
      <div v-if="canManage" class="section-label">Management</div>
      <div v-if="canManage" class="px-2 space-y-0.5">
        <NavItem :to="'/applications'" label="Applications" icon="icofont-files-stack" />
        <NavItem :to="'/admin/forms'" label="Form Builder" icon="icofont-edit" />
        <NavItem :to="'/admin/forms/submissions'" label="Form Submissions" icon="icofont-inbox" />
        <NavItem v-if="isAdminPlus" :to="'/users'" label="Users" icon="icofont-people" />
        <NavItem v-if="authStore.isSuperAdmin" :to="'/roles'" label="Roles & Permissions" icon="icofont-safety" />
      </div>

      <!-- WEBSITE CONTENT -->
      <template v-if="isContentManager">
        <div class="section-label">Website Content</div>
        <div class="px-2 space-y-0.5">
          <NavItem :to="'/cms'" label="Content Manager" icon="icofont-newspaper" />
          <NavItem :to="'/cms/page-builder'" label="Page Builder" icon="icofont-layout" />
        </div>
      </template>

      <!-- SECURITY -->
      <template v-if="isAdminPlus">
        <div class="section-label">Security</div>
        <div class="px-2 space-y-0.5">
          <NavItem :to="'/audit-logs'" label="Audit Logs" icon="icofont-history" />
        </div>
      </template>

      <template v-if="isAdminPlus">
        <div class="section-label">Operations</div>
        <div class="px-2 space-y-0.5"><NavItem :to="'/maintenance/command-center'" label="Command Center" icon="icofont-dashboard" /><NavItem :to="'/maintenance'" label="Maintenance Core" icon="icofont-tools-alt-2" /><NavItem v-if="authStore.isSuperAdmin" :to="'/settings/storage'" label="File Storage" icon="icofont-cloud-upload" /><NavItem v-if="authStore.isSuperAdmin" :to="'/maintenance/backups'" label="Backup & Restore" icon="icofont-database" /><NavItem v-if="authStore.isSuperAdmin" :to="'/maintenance/updates'" label="Smart Updates" icon="icofont-download-alt" /></div>
      </template>

      <!-- ACCOUNT -->
      <div class="section-label">Account</div>
      <div class="px-2 space-y-0.5">
        <NavItem :to="'/profile'" label="My Profile" icon="icofont-user-alt-5" />
        <NavItem :to="'/me/activities'" label="My Activities" icon="icofont-history" />
        <NavItem :to="'/me/security'" label="Security Settings" icon="icofont-shield-alt" />
      </div>
    </nav>

    <!-- Sign out -->
    <div class="flex-shrink-0 border-t border-primary-600/60 px-3 py-3">
      <button
        @click="handleLogout"
        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-primary-200 hover:bg-red-500/80 hover:text-white transition-all duration-150"
      >
        <i class="icofont-logout text-lg leading-none flex-shrink-0"></i>
        Sign Out
      </button>
    </div>
  </aside>
</template>

<script setup>
import { computed, defineComponent, h } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'

defineProps({ open: { type: Boolean, default: true } })
defineEmits(['close'])

const router = useRouter()
const authStore = useAuthStore()

const isAdminPlus     = computed(() => authStore.isSuperAdmin || authStore.isAdmin)
const isContentManager = computed(() => authStore.isSuperAdmin || authStore.isAdmin || authStore.isContentManager)
const canManage       = computed(() => isAdminPlus.value)

const userName = computed(() => {
  const u = authStore.user
  if (!u) return 'User'
  if (u.first_name || u.last_name) return `${u.first_name || ''} ${u.last_name || ''}`.trim()
  return u.email || 'User'
})

const userInitials = computed(() => {
  const u = authStore.user
  if (!u) return 'U'
  if (u.first_name && u.last_name) return `${u.first_name[0]}${u.last_name[0]}`.toUpperCase()
  if (u.first_name) return u.first_name[0].toUpperCase()
  if (u.email) return u.email[0].toUpperCase()
  return 'U'
})

const userRole = computed(() => {
  const u = authStore.user
  if (!u || !u.roles || u.roles.length === 0) return 'User'
  const role = u.roles[0]
  const name = typeof role === 'string' ? role : role.name
  return name.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
})

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}

// Inline nav item component to keep template clean
const NavItem = defineComponent({
  name: 'NavItem',
  props: {
    to: String,
    label: String,
    icon: String,
    highlight: { type: Boolean, default: false }
  },
  render() {
    return h(RouterLink, { to: this.to, custom: true }, {
      default: ({ href, navigate, isActive }) => h('a', {
        href,
        onClick: navigate,
        class: [
          this.highlight
            ? (isActive ? 'bg-yellow-400 text-primary-900' : 'bg-yellow-500/90 hover:bg-yellow-400 text-primary-900')
            : (isActive
                ? 'bg-primary-600 text-white border-l-2 border-white/70 pl-[10px]'
                : 'text-primary-200 hover:bg-primary-600/50 hover:text-white'),
          'group flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150'
        ].flat().filter(Boolean).join(' ')
      }, [
        h('i', { class: `text-base leading-none flex-shrink-0 ${this.icon}` }),
        this.label
      ])
    })
  }
})
</script>
