<template>
  <aside
    :class="[
      'otika-sidebar fixed inset-y-0 left-0 z-40 flex flex-col bg-white dark:bg-slate-900 shadow-xl transition-all duration-300',
      open ? 'w-64 translate-x-0' : 'sidebar-collapsed w-64 -translate-x-full md:w-20 md:translate-x-0'
    ]"
  >
    <!-- Brand -->
    <div class="sidebar-brand flex items-center gap-3 px-4 border-b border-slate-100 dark:border-slate-800 flex-shrink-0">
      <div class="flex-shrink-0 w-11 h-11 bg-white rounded-lg flex items-center justify-center overflow-hidden p-1">
        <img src="/main-logo.png" alt="NCS Logo" class="w-full h-full object-contain" />
      </div>
      <div class="min-w-0" :class="!open ? 'md:hidden' : ''">
        <div class="text-slate-800 dark:text-white font-bold text-sm leading-tight">National Council</div>
        <div class="text-[#6777ef] text-xs mt-0.5 leading-tight">of Sports · NCSMS</div>
      </div>
      <!-- Mobile close -->
      <button @click="$emit('close')" class="ml-auto md:hidden text-slate-400 hover:text-[#6777ef] p-1 rounded-lg">
        <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto pb-4">

      <!-- MAIN -->
      <div class="section-label">Main</div>
      <div class="px-2 space-y-0.5">
        <NavItem :collapsed="!open" :to="'/dashboard'" label="Dashboard" icon="icofont-dashboard-web" />
      </div>

      <template v-if="authStore.canUseNSMIS">
        <div class="section-label">NSMIS Reporting</div>
        <div class="px-2 space-y-0.5">
          <NavItem :collapsed="!open" v-if="authStore.hasAnyRole('super_admin','admin','ncs_general_secretary','general_secretary','technical_department','federation_president','federation_general_secretary','auditor')" :to="'/nsmis/governance'" label="Governance" icon="icofont-chart-histogram" />
          <NavItem :collapsed="!open" :to="'/nsmis/reports'" label="Federation Reports" icon="icofont-file-document" />
          <NavItem :collapsed="!open" v-if="authStore.hasAnyRole('super_admin','admin','ncs_general_secretary','general_secretary','technical_department')" :to="'/nsmis/insights/athletes'" label="Athletes" icon="icofont-runner-alt-1" />
          <NavItem :collapsed="!open" v-if="authStore.hasAnyRole('super_admin','admin','ncs_general_secretary','general_secretary','technical_department')" :to="'/nsmis/insights/performance'" label="Performance" icon="icofont-medal" />
          <NavItem :collapsed="!open" v-if="authStore.hasAnyRole('super_admin','admin','ncs_general_secretary','general_secretary','finance_department')" :to="'/nsmis/insights/finance'" label="Finance" icon="icofont-money" />
          <NavItem :collapsed="!open" v-if="authStore.hasAnyRole('super_admin','admin','ncs_general_secretary','general_secretary','technical_department')" :to="'/nsmis/insights/talent'" label="Talent" icon="icofont-search-user" />
          <NavItem :collapsed="!open" :to="'/nsmis/data/federation-officers'" label="Federation Profiles" icon="icofont-building-alt" />
          <NavItem :collapsed="!open" :to="'/nsmis/data/athletes'" label="Athlete Registry" icon="icofont-users-alt-4" />
          <NavItem :collapsed="!open" :to="'/nsmis/data/competitions'" label="Competitions" icon="icofont-trophy" />
          <NavItem :collapsed="!open" :to="'/nsmis/data/medals'" label="Medal Records" icon="icofont-medal-sport" />
          <NavItem :collapsed="!open" :to="'/nsmis/data/coaches'" label="Coaches" icon="icofont-teacher" />
          <NavItem :collapsed="!open" :to="'/nsmis/data/technical-officials'" label="Technical Officials" icon="icofont-judge" />
          <NavItem :collapsed="!open" :to="'/nsmis/data/equipment'" label="Equipment" icon="icofont-box" />
        </div>
      </template>

      <!-- EXECUTIVE COMMAND -->
      <template v-if="authStore.hasAnyRole('super_admin', 'admin', 'general_secretary', 'ags_technical', 'ags_admin')">
        <div class="section-label">Executive Command</div>
        <div class="px-2 space-y-0.5">
          <NavItem :collapsed="!open" v-if="authStore.hasAnyRole('super_admin', 'admin', 'general_secretary', 'ags_technical')" :to="'/executive/ags-technical'" label="AGS - Technical" icon="icofont-badge" />
          <NavItem :collapsed="!open" v-if="authStore.hasAnyRole('super_admin', 'admin', 'general_secretary', 'ags_admin')" :to="'/executive/ags-admin'" label="AGS - Administration" icon="icofont-architecture-alt" />
          <NavItem :collapsed="!open" v-if="authStore.hasAnyRole('super_admin', 'admin', 'general_secretary')" :to="'/executive/appraisals'" label="Master Appraisals" icon="icofont-law-order" />
        </div>
      </template>

      <!-- DEPARTMENT OPERATIONS -->
      <div class="section-label">Department Portals</div>
      <div class="px-2 space-y-0.5">
        <NavItem :collapsed="!open" :to="'/stores/inventory'" label="Stores & Inventory" icon="icofont-box" />
        <NavItem :collapsed="!open" :to="'/facilities/venues'" label="Venues & Facilities" icon="icofont-building" />
        <NavItem :collapsed="!open" :to="'/legal/compliance'" label="Legal & Compliance" icon="icofont-law-document" />
        <NavItem :collapsed="!open" :to="'/medical/sports-science'" label="Sports Medicine & WADA" icon="icofont-heart-beat" />
        <NavItem :collapsed="!open" :to="'/fleet/transport'" label="Fleet & Logistics" icon="icofont-truck" />
      </div>

      <!-- FINANCIAL & ASSETS LEDGER -->
      <template v-if="canAccessFinance">
        <div class="section-label">Financial & Assets Ledger</div>
        <div class="px-2 space-y-0.5">
          <!-- Fixed Assets Register Dropdown Group -->
          <div class="space-y-0.5">
            <button
              @click="isFixedAssetsExpanded = !isFixedAssetsExpanded"
              type="button"
              :class="[
                isFixedAssetsRouteActive ? 'otika-menu-active' : 'otika-menu-link',
                open ? 'justify-between' : 'md:justify-center',
                'w-full flex items-center px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150 group'
              ]"
              :title="!open ? 'Fixed Assets Register' : undefined"
            >
              <div class="flex items-center gap-3 min-w-0">
                <i class="icofont-building text-base leading-none flex-shrink-0"></i>
                <span :class="!open ? 'md:hidden' : ''" class="truncate font-semibold">Fixed Assets Register</span>
              </div>
              <i
                :class="[
                  isFixedAssetsExpanded ? 'icofont-simple-up' : 'icofont-simple-down',
                  !open ? 'md:hidden' : ''
                ]"
                class="text-xs ml-2 flex-shrink-0 transition-transform duration-200 opacity-80 group-hover:opacity-100"
              ></i>
            </button>

            <div v-show="isFixedAssetsExpanded" :class="!open ? 'md:pl-0' : 'pl-3'" class="space-y-0.5 mt-0.5 border-l-2 border-slate-100 dark:border-slate-700 ml-2">
              <NavItem :collapsed="!open" :to="'/fixed-assets'" label="Manage Assets" icon="icofont-listine-dots" />
              <NavItem :collapsed="!open" :to="'/fixed-assets/value-adjustments'" label="Value Adjustments" icon="icofont-history" />
              <NavItem :collapsed="!open" :to="'/fixed-assets/pivot-engine'" label="Dynamic Pivot Engine" icon="icofont-chart-histogram" />
            </div>
          </div>

          <NavItem :collapsed="!open" :to="'/applications'" label="Applications" icon="icofont-files-stack" />
          <NavItem :collapsed="!open" :to="'/admin/forms'" label="Form Builder" icon="icofont-edit" />
          <NavItem :collapsed="!open" :to="'/admin/forms/submissions'" label="Form Submissions" icon="icofont-inbox" />
          <NavItem :collapsed="!open" v-if="isAdminPlus" :to="'/users'" label="Users" icon="icofont-people" />
          <NavItem :collapsed="!open" v-if="authStore.isSuperAdmin" :to="'/roles'" label="Roles & Permissions" icon="icofont-safety" />
        </div>
      </template>

      <!-- WEBSITE CONTENT -->
      <template v-if="isContentManager">
        <div class="section-label">Website Content</div>
        <div class="px-2 space-y-0.5">
          <NavItem :collapsed="!open" :to="'/cms'" label="Content Manager" icon="icofont-newspaper" />
          <NavItem :collapsed="!open" :to="'/cms/page-builder'" label="Page Builder" icon="icofont-layout" />
        </div>
      </template>

      <!-- SECURITY -->
      <template v-if="isAdminPlus">
        <div class="section-label">Security</div>
        <div class="px-2 space-y-0.5">
          <NavItem :collapsed="!open" :to="'/audit-logs'" label="Audit Logs" icon="icofont-history" />
        </div>
      </template>

      <template v-if="isAdminPlus">
        <div class="section-label">Operations</div>
        <div class="px-2 space-y-0.5"><NavItem :collapsed="!open" :to="'/maintenance/command-center'" label="Command Center" icon="icofont-dashboard" /><NavItem :collapsed="!open" :to="'/maintenance'" label="Maintenance Core" icon="icofont-tools-alt-2" /><NavItem :collapsed="!open" v-if="authStore.isSuperAdmin" :to="'/settings/storage'" label="File Storage" icon="icofont-cloud-upload" /><NavItem :collapsed="!open" v-if="authStore.isSuperAdmin" :to="'/maintenance/backups'" label="Backup & Restore" icon="icofont-database" /><NavItem :collapsed="!open" v-if="authStore.isSuperAdmin" :to="'/maintenance/updates'" label="Smart Updates" icon="icofont-download-alt" /></div>
      </template>

      <!-- ACCOUNT -->
      <div class="section-label">Account</div>
      <div class="px-2 space-y-0.5">
        <NavItem :collapsed="!open" :to="'/profile'" label="My Profile" icon="icofont-user-alt-5" />
        <NavItem :collapsed="!open" :to="'/me/activities'" label="My Activities" icon="icofont-history" />
        <NavItem :collapsed="!open" :to="'/me/security'" label="Security Settings" icon="icofont-shield-alt" />
      </div>
    </nav>

    <!-- Sign out -->
    <div class="sidebar-footer flex-shrink-0 border-t border-slate-100 dark:border-slate-800 px-3 py-3 space-y-2">
      <ThemeToggle class="sidebar-theme-toggle" />
      <button
        @click="handleLogout"
        class="signout-link w-full flex items-center gap-3 px-3 py-2.5 text-sm font-medium text-slate-600 dark:text-slate-300 hover:text-red-600 transition-all duration-150"
      >
        <i class="icofont-logout text-lg leading-none flex-shrink-0"></i>
        Sign Out
      </button>
    </div>
  </aside>
</template>

<script setup>
import { computed, defineComponent, h, ref, watch } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import ThemeToggle from '@/components/ui/ThemeToggle.vue'

defineProps({ open: { type: Boolean, default: true } })
defineEmits(['close'])

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const isAdminPlus     = computed(() => authStore.isSuperAdmin || authStore.isAdmin)
const isContentManager = computed(() => authStore.isSuperAdmin || authStore.isAdmin || authStore.isContentManager)
const canManage       = computed(() => isAdminPlus.value)
const canAccessFinance = computed(() => {
  return isAdminPlus.value || authStore.hasAnyRole('accountant', 'finance_officer', 'general_secretary', 'ags_admin', 'ags_technical', 'finance_department', 'senior_accountant')
})

const isFixedAssetsRouteActive = computed(() => route.path.startsWith('/fixed-assets'))
const isFixedAssetsExpanded = ref(true)

watch(() => route.path, (newPath) => {
  if (newPath.startsWith('/fixed-assets')) {
    isFixedAssetsExpanded.value = true
  }
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
    highlight: { type: Boolean, default: false },
    collapsed: { type: Boolean, default: false }
  },
  render() {
    return h(RouterLink, { to: this.to, custom: true }, {
      default: ({ href, navigate, isActive }) => h('a', {
        href,
        onClick: navigate,
        title: this.collapsed ? this.label : undefined,
        class: [
          this.highlight ? 'otika-nav-highlight' : (isActive ? 'otika-nav-active' : 'otika-nav-link'),
          this.collapsed ? 'md:justify-center md:gap-0' : '',
          'group flex items-center gap-3 px-3 py-2.5 text-sm font-medium transition-all duration-150'
        ].flat().filter(Boolean).join(' ')
      }, [
        h('i', { class: `text-base leading-none flex-shrink-0 ${this.icon}` }),
        h('span', { class: this.collapsed ? 'md:hidden' : '' }, this.label)
      ])
    })
  }
})
</script>

<style scoped>
.otika-sidebar { width: 250px; color: #6c757d; box-shadow: 0 4px 25px rgba(0, 0, 0, .1); }
.sidebar-brand { height: 70px; }
.otika-sidebar nav { scrollbar-width: thin; scrollbar-color: #d9d9d9 transparent; }
.otika-sidebar :deep(.section-label) { color: #a1a8ae; font-size: 10px; font-weight: 700; letter-spacing: 1.3px; padding: 18px 20px 6px; }
.otika-sidebar :deep(.otika-nav-link), .otika-sidebar .otika-menu-link { color: #60686f; border-radius: 0; border-left: 3px solid transparent; }
.otika-sidebar :deep(.otika-nav-link:hover), .otika-sidebar .otika-menu-link:hover { color: #6777ef; background: #f8f9fa; }
.otika-sidebar :deep(.otika-nav-active), .otika-sidebar .otika-menu-active { color: #6777ef; background: #f3f4fd; border-left: 3px solid #6777ef; font-weight: 600; }
.otika-sidebar :deep(.otika-nav-highlight) { color: #6777ef; background: #f3f4fd; border-left: 3px solid #6777ef; }
.otika-sidebar :deep(.otika-nav-link i), .otika-sidebar :deep(.otika-nav-active i) { width: 20px; text-align: center; font-size: 17px; }
.sidebar-footer :deep(.theme-toggle-btn), .sidebar-footer :deep(.theme-toggle-btn > div) { width: 100%; }
.sidebar-footer :deep(.theme-toggle-btn > div) { justify-content: flex-start; border-radius: 3px; }
.signout-link { border-radius: 3px; }
:global(html.dark) .otika-sidebar :deep(.section-label) { color: #64748b; }
:global(html.dark) .otika-sidebar :deep(.otika-nav-link), :global(html.dark) .otika-sidebar .otika-menu-link { color: #94a3b8; }
:global(html.dark) .otika-sidebar :deep(.otika-nav-link:hover), :global(html.dark) .otika-sidebar .otika-menu-link:hover { color: #a5b4fc; background: #1e293b; }
:global(html.dark) .otika-sidebar :deep(.otika-nav-active), :global(html.dark) .otika-sidebar .otika-menu-active { color: #a5b4fc; background: rgba(103,119,239,.16); border-left-color: #818cf8; }
@media (min-width: 768px) {
  .sidebar-collapsed.otika-sidebar { width: 80px; }
  .sidebar-collapsed :deep(.section-label) {
    font-size: 0;
    padding-left: 0;
    padding-right: 0;
    text-align: center;
  }
  .sidebar-collapsed :deep(.section-label)::after {
    content: "";
    display: block;
    width: 1.75rem;
    height: 1px;
    margin: 0.75rem auto 0.35rem;
    background: #d9d9d9;
  }
  .sidebar-collapsed .sidebar-footer :deep(.theme-toggle-btn span),
  .sidebar-collapsed .signout-link:not(:hover) { font-size: 0; }
}
</style>

