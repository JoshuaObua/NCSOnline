<template>
  <div class="otika-app">
    <div class="main-wrapper main-wrapper-1" :class="{ 'sidebar-mini': sidebarCollapsed }">
      
      <!-- Navbar Background -->
      <div class="navbar-bg"></div>

      <!-- ── Top Navbar ──────────────────────────────────────────────── -->
      <nav class="main-navbar sticky">
        <!-- Left Side: Toggles & Search -->
        <div class="flex items-center gap-3">
          <!-- Sidebar Toggle Button -->
          <button
            type="button"
            class="cms-top-icon collapse-btn"
            @click="toggleSidebar"
            :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'"
          >
            <i class="icofont-navigation-menu text-lg"></i>
          </button>

          <!-- Refresh Button -->
          <button
            type="button"
            class="cms-top-icon"
            @click="refreshPage"
            title="Refresh View"
          >
            <i class="icofont-refresh text-lg"></i>
          </button>

          <!-- Otika Search Element -->
          <div class="search-element hidden md:flex">
            <input
              type="text"
              v-model="searchQuery"
              class="form-control"
              placeholder="Search registry, athletes, files..."
              @keyup.enter="handleSearch"
            />
            <button type="button" class="btn" @click="handleSearch">
              <i class="icofont-search-1"></i>
            </button>
          </div>
        </div>

        <!-- Right Side: Theme, Notifications, Profile -->
        <div class="flex items-center gap-2">
          <!-- Dark Mode / Theme Toggle -->
          <ThemeToggle class="cms-top-icon" />

          <!-- Notifications Bell -->
          <div class="relative">
            <button
              type="button"
              class="cms-top-icon"
              @click="notificationsOpen = !notificationsOpen; profileOpen = false"
              title="Notifications"
            >
              <i class="icofont-notification text-lg"></i>
              <span class="headerBadge1">3</span>
            </button>

            <!-- Notifications Dropdown -->
            <div
              v-if="notificationsOpen"
              class="absolute right-0 mt-2 w-80 bg-white dark:bg-slate-900 rounded-xl shadow-xl border border-slate-200 dark:border-slate-800 py-2 z-50 animate-fadeIn"
            >
              <div class="px-4 py-2 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between">
                <span class="text-xs font-bold uppercase tracking-wider text-slate-800 dark:text-slate-200">Notifications</span>
                <span class="text-[10px] bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full font-bold">3 New</span>
              </div>
              <div class="max-h-60 overflow-y-auto divide-y divide-slate-100 dark:divide-slate-800">
                <div class="px-4 py-2.5 hover:bg-slate-50 dark:hover:bg-slate-800/50 cursor-pointer transition-colors">
                  <div class="text-xs font-semibold text-slate-800 dark:text-slate-200">Form 5 Vetting Requisition</div>
                  <div class="text-[11px] text-slate-500 mt-0.5">Procurement submitted PPDA Form 5 for approval</div>
                </div>
                <div class="px-4 py-2.5 hover:bg-slate-50 dark:hover:bg-slate-800/50 cursor-pointer transition-colors">
                  <div class="text-xs font-semibold text-slate-800 dark:text-slate-200">Lugogo Arena Booking</div>
                  <div class="text-[11px] text-slate-500 mt-0.5">UAF requested match readiness inspection</div>
                </div>
                <div class="px-4 py-2.5 hover:bg-slate-50 dark:hover:bg-slate-800/50 cursor-pointer transition-colors">
                  <div class="text-xs font-semibold text-slate-800 dark:text-slate-200">Asset Valuation Sign-Off</div>
                  <div class="text-[11px] text-slate-500 mt-0.5">Finance Department updated fixed asset ledger</div>
                </div>
              </div>
            </div>
          </div>

          <!-- User Profile Dropdown -->
          <div class="relative">
            <button
              type="button"
              class="flex items-center gap-2 pl-2 pr-3 py-1.5 rounded-xl hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
              @click="profileOpen = !profileOpen; notificationsOpen = false"
            >
              <div class="w-8 h-8 rounded-full bg-blue-600 text-white font-bold text-xs flex items-center justify-center border border-blue-400">
                {{ userInitials }}
              </div>
              <div class="hidden sm:block text-left">
                <span class="block text-xs font-bold text-slate-800 dark:text-slate-200 leading-tight">{{ userName }}</span>
                <span class="block text-[10px] text-slate-400 font-medium leading-tight">{{ userRole }}</span>
              </div>
              <i class="icofont-thin-down text-xs text-slate-400"></i>
            </button>

            <!-- Profile Menu Dropdown -->
            <div
              v-if="profileOpen"
              class="absolute right-0 mt-2 w-56 bg-white dark:bg-slate-900 rounded-xl shadow-xl border border-slate-200 dark:border-slate-800 py-1.5 z-50 animate-fadeIn"
            >
              <div class="px-4 py-2 border-b border-slate-100 dark:border-slate-800">
                <div class="text-xs font-bold text-slate-800 dark:text-slate-200">{{ userName }}</div>
                <div class="text-[11px] text-slate-500 truncate">{{ authStore.user?.email }}</div>
              </div>
              <router-link
                to="/profile"
                class="flex items-center gap-2 px-4 py-2 text-xs text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors"
                @click="profileOpen = false"
              >
                <i class="icofont-user text-sm text-blue-600"></i> My Profile
              </router-link>
              <router-link
                to="/me/security"
                class="flex items-center gap-2 px-4 py-2 text-xs text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors"
                @click="profileOpen = false"
              >
                <i class="icofont-lock text-sm text-amber-600"></i> Security Settings
              </router-link>
              <div class="border-t border-slate-100 dark:border-slate-800 my-1"></div>
              <button
                type="button"
                class="w-full flex items-center gap-2 px-4 py-2 text-xs text-red-600 hover:bg-red-50 dark:hover:bg-red-950/30 transition-colors text-left font-semibold"
                @click="handleLogout"
              >
                <i class="icofont-logout text-sm text-red-600"></i> Sign Out
              </button>
            </div>
          </div>
        </div>
      </nav>

      <!-- ── Mobile Drawer Overlay ───────────────────────────────────── -->
      <div
        v-if="isMobile && !sidebarCollapsed"
        class="fixed inset-0 bg-black/40 z-40 md:hidden"
        @click="sidebarCollapsed = true"
      ></div>

      <!-- ── Otika Sidebar ───────────────────────────────────────────── -->
      <div class="main-sidebar sidebar-style-2">
        <aside id="sidebar-wrapper">
          <!-- Sidebar Brand -->
          <div class="sidebar-brand">
            <router-link to="/dashboard">
              <img src="/main-logo.png" alt="NCS Logo" class="header-logo" />
              <span v-if="!sidebarCollapsed" class="sidebar-brand-text">NCS ONLINE</span>
            </router-link>
          </div>

          <!-- Sidebar User Card -->
          <div v-if="!sidebarCollapsed" class="sidebar-user">
            <div class="w-9 h-9 rounded-full bg-blue-600 text-white font-bold text-xs flex items-center justify-center flex-shrink-0">
              {{ userInitials }}
            </div>
            <div class="sidebar-user-info">
              <span class="sidebar-user-name" :title="userName">{{ userName }}</span>
              <span class="sidebar-user-role" :title="userRole">{{ userRole }}</span>
            </div>
          </div>

          <!-- Sidebar Menu Links -->
          <ul class="sidebar-menu">
            <!-- MAIN COMMAND -->
            <li class="menu-header">Main Command</li>
            <li :class="{ active: currentPath === '/dashboard' }">
              <router-link to="/dashboard" title="Dashboard">
                <i class="icofont-dashboard-web"></i>
                <span v-if="!sidebarCollapsed">Dashboard Overview</span>
              </router-link>
            </li>

            <!-- UNIVERSAL APPLICATION MODULES: visible across all dashboards -->
            <li class="menu-header">Universal Applications</li>
            <li :class="{ active: currentPath === '/reception/leave/apply' }">
              <router-link to="/reception/leave/apply" title="Apply for Leave">
                <i class="icofont-calendar"></i>
                <span v-if="!sidebarCollapsed">Leave Application</span>
              </router-link>
            </li>
            <li :class="{ active: currentPath === '/it/ppda/new' }">
              <router-link to="/it/ppda/new" title="New PPDA Form 5 Application">
                <i class="icofont-document-folder"></i>
                <span v-if="!sidebarCollapsed">PPDA Application</span>
              </router-link>
            </li>

            <!-- RECEPTIONIST SPECIFIC DESK MENU -->
            <template v-if="isReceptionistOnly">
              <li class="menu-header">Front Desk & Visitor Control</li>
              <li :class="{ active: currentPath === '/reception/visitors' }">
                <router-link to="/reception/visitors" title="Visitor Clearance Registry">
                  <i class="icofont-id-card"></i>
                  <span v-if="!sidebarCollapsed">Visitor Clearance</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/reception/leave/apply' }">
                <router-link to="/reception/leave/apply" title="Apply for Leave">
                  <i class="icofont-calendar"></i>
                  <span v-if="!sidebarCollapsed">Apply for Leave</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/reception/leave/status' }">
                <router-link to="/reception/leave/status" title="Leave Status & History">
                  <i class="icofont-history"></i>
                  <span v-if="!sidebarCollapsed">Leave Status & History</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/reception/reports' }">
                <router-link to="/reception/reports" title="Application & Desk Reports">
                  <i class="icofont-chart-bar-graph"></i>
                  <span v-if="!sidebarCollapsed">Application Reports</span>
                </router-link>
              </li>
            </template>

            <!-- IT & ENGINEERING DEPARTMENT SPECIFIC WORKSTATION MENU -->
            <template v-else-if="isITOfficerOnly || isEngineeringOfficerOnly">
              <li class="menu-header">Procurement & PPDA</li>
              <li :class="{ active: currentPath === '/it/ppda/status' || currentPath === '/it/ppda' }">
                <router-link to="/it/ppda/status" title="PPDA Form 5 Requisitions">
                  <i class="icofont-document-folder"></i>
                  <span v-if="!sidebarCollapsed">PPDA Form 5 Status</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/it/ppda/new' }">
                <router-link to="/it/ppda/new" title="New PPDA Form 5 Requisition">
                  <i class="icofont-plus-circle"></i>
                  <span v-if="!sidebarCollapsed">New PPDA Form 5</span>
                </router-link>
              </li>

              <li class="menu-header">Staff Self-Service</li>
              <li :class="{ active: currentPath === '/reception/leave/apply' }">
                <router-link to="/reception/leave/apply" title="Apply for Leave">
                  <i class="icofont-calendar"></i>
                  <span v-if="!sidebarCollapsed">Apply for Leave</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/reception/leave/status' }">
                <router-link to="/reception/leave/status" title="Leave Status & History">
                  <i class="icofont-history"></i>
                  <span v-if="!sidebarCollapsed">Leave Status & History</span>
                </router-link>
              </li>
            </template>

            <!-- ACCOUNTING & FINANCE DEPARTMENT SPECIFIC WORKSTATION MENU -->
            <template v-else-if="isAccountantOnly">
              <li class="menu-header">Financial & Asset Ledgers</li>
              <li :class="{ active: currentPath.startsWith('/fixed-assets') }">
                <router-link to="/fixed-assets" title="Fixed Assets & IPSAS 17 Register">
                  <i class="icofont-building-alt"></i>
                  <span v-if="!sidebarCollapsed">Fixed Assets Register</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/stores/inventory' }">
                <router-link to="/stores/inventory" title="Stores & Inventory">
                  <i class="icofont-box"></i>
                  <span v-if="!sidebarCollapsed">Stores & Inventory</span>
                </router-link>
              </li>

              <li class="menu-header">Procurement & Vetting</li>
              <li :class="{ active: currentPath === '/it/ppda/status' || currentPath === '/it/ppda' }">
                <router-link to="/it/ppda/status" title="PPDA Form 5 Requisitions">
                  <i class="icofont-document-folder"></i>
                  <span v-if="!sidebarCollapsed">PPDA Form 5 Status</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/it/ppda/new' }">
                <router-link to="/it/ppda/new" title="New PPDA Form 5 Requisition">
                  <i class="icofont-plus-circle"></i>
                  <span v-if="!sidebarCollapsed">New PPDA Form 5</span>
                </router-link>
              </li>

              <li class="menu-header">Staff Self-Service</li>
              <li :class="{ active: currentPath === '/reception/leave/apply' }">
                <router-link to="/reception/leave/apply" title="Apply for Leave">
                  <i class="icofont-calendar"></i>
                  <span v-if="!sidebarCollapsed">Apply for Leave</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/reception/leave/status' }">
                <router-link to="/reception/leave/status" title="Leave Status & History">
                  <i class="icofont-history"></i>
                  <span v-if="!sidebarCollapsed">Leave Status & History</span>
                </router-link>
              </li>
            </template>

            <!-- GENERAL / EXECUTIVE / DEPARTMENTAL MENUS (HIDDEN FOR RECEPTIONIST ONLY) -->
            <template v-else>
              <!-- EXECUTIVE LEADERSHIP & GOVERNANCE -->
              <template v-if="authStore.hasAnyRole('super_admin', 'admin', 'general_secretary', 'ags_technical', 'ags_admin')">
                <li class="menu-header">Executive Leadership & Governance</li>
                <li v-if="authStore.hasAnyRole('super_admin', 'admin', 'general_secretary')" :class="{ active: currentPath === '/executive/appraisals' }">
                  <router-link to="/executive/appraisals" title="Online Appraisal">
                    <i class="icofont-law-order"></i>
                    <span v-if="!sidebarCollapsed">Online Appraisal</span>
                  </router-link>
                </li>

                <!-- DEPARTMENTAL COMPILED REPORTS DROPDOWN -->
                <li v-if="authStore.hasAnyRole('super_admin', 'admin', 'general_secretary')" class="dropdown" :class="{ active: isReportsGroupActive }">
                  <a href="javascript:void(0)" @click="reportsDropdownOpen = !reportsDropdownOpen" class="nav-link has-dropdown flex items-center justify-between" title="Departmental Reports">
                    <div class="flex items-center gap-3">
                      <i class="icofont-files-stack"></i>
                      <span v-if="!sidebarCollapsed">Departmental Reports</span>
                    </div>
                    <i v-if="!sidebarCollapsed" class="icofont-rounded-down text-xs transition-transform duration-200" :class="{ 'rotate-180': reportsDropdownOpen || isReportsGroupActive }"></i>
                  </a>
                  <ul v-show="(reportsDropdownOpen || isReportsGroupActive) && !sidebarCollapsed" class="dropdown-menu pl-4 space-y-1 mt-1">
                    <li :class="{ active: currentPath === '/executive/reports/engineering' }">
                      <router-link to="/executive/reports/engineering" class="flex items-center gap-2 py-1.5 text-xs">
                        <i class="icofont-wrench text-[11px]"></i> Engineering & Works
                      </router-link>
                    </li>
                    <li :class="{ active: currentPath === '/executive/reports/human-resources' }">
                      <router-link to="/executive/reports/human-resources" class="flex items-center gap-2 py-1.5 text-xs">
                        <i class="icofont-users text-[11px]"></i> Human Resources
                      </router-link>
                    </li>
                    <li :class="{ active: currentPath === '/executive/reports/finance-accounts' }">
                      <router-link to="/executive/reports/finance-accounts" class="flex items-center gap-2 py-1.5 text-xs">
                        <i class="icofont-money text-[11px]"></i> Finance & Accounts
                      </router-link>
                    </li>
                    <li :class="{ active: currentPath === '/executive/reports/technical-sports' }">
                      <router-link to="/executive/reports/technical-sports" class="flex items-center gap-2 py-1.5 text-xs">
                        <i class="icofont-badge text-[11px]"></i> Technical & Sports
                      </router-link>
                    </li>
                    <li :class="{ active: currentPath === '/executive/reports/medical-science' }">
                      <router-link to="/executive/reports/medical-science" class="flex items-center gap-2 py-1.5 text-xs">
                        <i class="icofont-heart-beat text-[11px]"></i> Sports Medicine / WADA
                      </router-link>
                    </li>
                    <li :class="{ active: currentPath === '/executive/reports/legal-logistics' }">
                      <router-link to="/executive/reports/legal-logistics" class="flex items-center gap-2 py-1.5 text-xs">
                        <i class="icofont-law-document text-[11px]"></i> Legal & Logistics
                      </router-link>
                    </li>
                  </ul>
                </li>

                <li v-if="authStore.hasAnyRole('super_admin', 'admin', 'ags_technical')" :class="{ active: currentPath === '/executive/ags-technical' }">
                  <router-link to="/executive/ags-technical" title="AGS Technical">
                    <i class="icofont-badge"></i>
                    <span v-if="!sidebarCollapsed">AGS - Technical</span>
                  </router-link>
                </li>
                <li v-if="authStore.hasAnyRole('super_admin', 'admin', 'ags_admin')" :class="{ active: currentPath === '/executive/ags-admin' }">
                  <router-link to="/executive/ags-admin" title="AGS Administration">
                    <i class="icofont-architecture-alt"></i>
                    <span v-if="!sidebarCollapsed">AGS - Administration</span>
                  </router-link>
                </li>
              </template>

              <!-- OPERATIONS & LOGISTICS -->
              <li class="menu-header">Operations & Front Desk</li>
              <li :class="{ active: currentPath.startsWith('/fixed-assets') }">
                <router-link to="/fixed-assets" title="Fixed Assets & IPSAS 17 Register">
                  <i class="icofont-building-alt"></i>
                  <span v-if="!sidebarCollapsed">Fixed Assets Register</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/reception/visitors' }">
                <router-link to="/reception/visitors" title="Reception & Visitor Clearance">
                  <i class="icofont-id-card"></i>
                  <span v-if="!sidebarCollapsed">Reception & Visitors</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/stores/inventory' }">
                <router-link to="/stores/inventory" title="Stores & Inventory">
                  <i class="icofont-box"></i>
                  <span v-if="!sidebarCollapsed">Stores & Inventory</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/facilities/venues' }">
                <router-link to="/facilities/venues" title="Facilities & Venues">
                  <i class="icofont-building"></i>
                  <span v-if="!sidebarCollapsed">Venues & Facilities</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/fleet/transport' }">
                <router-link to="/fleet/transport" title="Transport & Fleet">
                  <i class="icofont-truck"></i>
                  <span v-if="!sidebarCollapsed">Fleet & Transport</span>
                </router-link>
              </li>

              <!-- SPORTS & COMPLIANCE -->
              <li class="menu-header">Sports & Compliance</li>
              <li :class="{ active: currentPath === '/medical/sports-science' }">
                <router-link to="/medical/sports-science" title="Sports Medicine & WADA">
                  <i class="icofont-heart-beat"></i>
                  <span v-if="!sidebarCollapsed">Sports Medicine / WADA</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/legal/compliance' }">
                <router-link to="/legal/compliance" title="Legal & Compliance">
                  <i class="icofont-law-document"></i>
                  <span v-if="!sidebarCollapsed">Legal & Compliance</span>
                </router-link>
              </li>

              <template v-if="authStore.canUseNSMIS">
                <li :class="{ active: currentPath === '/nsmis/governance' }">
                  <router-link to="/nsmis/governance" title="NSMIS Governance">
                    <i class="icofont-chart-histogram"></i>
                    <span v-if="!sidebarCollapsed">NSMIS Governance</span>
                  </router-link>
                </li>
                <li :class="{ active: currentPath === '/nsmis/reports' }">
                  <router-link to="/nsmis/reports" title="Federation Reports">
                    <i class="icofont-file-document"></i>
                    <span v-if="!sidebarCollapsed">Federation Reports</span>
                  </router-link>
                </li>
                <li :class="{ active: currentPath.startsWith('/nsmis/data') }">
                  <router-link to="/nsmis/data/athletes" title="Athletes Registry">
                    <i class="icofont-runner-alt-1"></i>
                    <span v-if="!sidebarCollapsed">Athlete Registry</span>
                  </router-link>
                </li>
              </template>

              <!-- ADMINISTRATION & SYSTEM -->
              <li class="menu-header">Administration</li>
              <li v-if="authStore.isAdmin" :class="{ active: currentPath === '/users' }">
                <router-link to="/users" title="Staff & Users">
                  <i class="icofont-users-alt-5"></i>
                  <span v-if="!sidebarCollapsed">Staff & Users</span>
                </router-link>
              </li>
              <li :class="{ active: currentPath === '/applications' }">
                <router-link to="/applications" title="Applications & Licences">
                  <i class="icofont-files-stack"></i>
                  <span v-if="!sidebarCollapsed">Licence Applications</span>
                </router-link>
              </li>
              <li v-if="authStore.isAdmin" :class="{ active: currentPath === '/maintenance/command-center' }">
                <router-link to="/maintenance/command-center" title="Infrastructure Command">
                  <i class="icofont-tools-alt-2"></i>
                  <span v-if="!sidebarCollapsed">Infrastructure Hub</span>
                </router-link>
              </li>
              <li v-if="authStore.isAdmin" :class="{ active: currentPath === '/cms' }">
                <router-link to="/cms" title="CMS Website Editor">
                  <i class="icofont-ui-browser"></i>
                  <span v-if="!sidebarCollapsed">Public CMS Editor</span>
                </router-link>
              </li>
              <li v-if="authStore.isAdmin" :class="{ active: currentPath === '/audit-logs' }">
                <router-link to="/audit-logs" title="Audit Trail">
                  <i class="icofont-shield-alt"></i>
                  <span v-if="!sidebarCollapsed">Audit Logs</span>
                </router-link>
              </li>
              <li v-if="authStore.isSuperAdmin" :class="{ active: currentPath === '/maintenance/backups' }">
                <router-link to="/maintenance/backups" title="Cloud Backups">
                  <i class="icofont-database"></i>
                  <span v-if="!sidebarCollapsed">DB Cloud Backups</span>
                </router-link>
              </li>
              <li v-if="authStore.isSuperAdmin" :class="{ active: currentPath === '/maintenance/updates' }">
                <router-link to="/maintenance/updates" title="Smart Updates">
                  <i class="icofont-download"></i>
                  <span v-if="!sidebarCollapsed">Smart Updates</span>
                </router-link>
              </li>
            </template>
          </ul>
        </aside>
      </div>

      <!-- ── Main Content Area ───────────────────────────────────────── -->
      <div class="main-content">
        <section class="section">
          <!-- Otika Section Header & Breadcrumb -->
          <div class="section-header">
            <div>
              <h1>{{ displayTitle }}</h1>
            </div>
            <div class="flex items-center gap-4">
              <!-- Custom Action Slot -->
              <slot name="actions"></slot>

              <!-- Dynamic Breadcrumbs -->
              <ul class="section-header-breadcrumb hidden sm:flex">
                <li class="breadcrumb-item">
                  <router-link to="/dashboard"><i class="icofont-home mr-1"></i>Home</router-link>
                </li>
                <template v-if="crumbs.length">
                  <li v-for="(crumb, i) in crumbs" :key="i" class="breadcrumb-item" :class="{ active: i === crumbs.length - 1 }">
                    <router-link v-if="crumb.to && i !== crumbs.length - 1" :to="crumb.to">{{ crumb.label }}</router-link>
                    <span v-else>{{ crumb.label }}</span>
                  </li>
                </template>
                <li v-else class="breadcrumb-item active">{{ displayTitle }}</li>
              </ul>
            </div>
          </div>

          <!-- Section Body Default Slot -->
          <div class="section-body">
            <UniversalDashboardApplications class="mb-6" />
            <slot></slot>
          </div>
        </section>
      </div>

      <!-- ── Otika Master Footer ─────────────────────────────────────── -->
      <footer class="main-footer">
        <div class="footer-left">
          Copyright &copy; {{ year }} <a href="https://ncs.go.ug" target="_blank" class="text-blue-600 font-bold hover:underline">National Council of Sports</a>. All rights reserved.
        </div>
        <div class="footer-right">
          NCSMS Intranet v1.3.0 · <span class="text-blue-600 font-semibold">National Council of Sports</span>
        </div>
      </footer>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'
import ThemeToggle from '@/components/ui/ThemeToggle.vue'
import UniversalDashboardApplications from '@/components/layout/UniversalDashboardApplications.vue'

const props = defineProps({
  title: { type: String, default: '' }
})

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const breadcrumbStore = useBreadcrumbStore()

const year = new Date().getFullYear()
const currentPath = computed(() => route.path)
const isReportsGroupActive = computed(() => route.path.startsWith('/executive/reports'))
const displayTitle = computed(() => props.title || breadcrumbStore.title || 'National Council of Sports')
const crumbs = computed(() => breadcrumbStore.crumbs)

// UI Dropdowns & Sidebar state
const sidebarCollapsed = ref(false)
const profileOpen = ref(false)
const notificationsOpen = ref(false)
const reportsDropdownOpen = ref(false)
const searchQuery = ref('')
const isMobile = ref(false)

const userName = computed(() => {
  if (authStore.user?.first_name || authStore.user?.last_name) {
    return `${authStore.user.first_name || ''} ${authStore.user.last_name || ''}`.trim()
  }
  return authStore.user?.email || 'NCS Staff'
})

const userInitials = computed(() => {
  const f = authStore.user?.first_name?.[0] || ''
  const l = authStore.user?.last_name?.[0] || ''
  return (f + l).toUpperCase() || 'NC'
})

const userRole = computed(() => {
  const rawRoles = authStore.user?.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
  if (roles.includes('super_admin')) return 'Super Admin'
  if (roles.includes('general_secretary')) return 'General Secretary (CEO)'
  if (roles.includes('ags_technical')) return 'AGS Technical'
  if (roles.includes('ags_admin')) return 'AGS Administration'
  if (roles.includes('stores_officer')) return 'Stores Officer'
  if (roles.includes('facilities_manager')) return 'Facilities Manager'
  if (roles.includes('transport_officer')) return 'Transport Officer'
  if (roles.includes('medical_officer')) return 'Medical Officer'
  if (roles.includes('legal_counsel')) return 'Legal Counsel'
  if (roles.includes('senior_engineer')) return 'Senior Engineer'
  if (roles.includes('it_officer')) return 'ICT Systems & Database Administrator'
  if (roles.includes('accountant') || roles.includes('senior_accountant') || roles.includes('finance_department')) return 'Senior Accountant / Head of Finance'
  if (roles.includes('receptionist') || roles.includes('helpdesk')) return 'Front Desk Receptionist'
  const first = roles[0]
  return typeof first === 'string' && first ? first.replace(/_/g, ' ') : 'Staff Member'
})

const isReceptionistOnly = computed(() => {
  const rawRoles = authStore.user?.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
  const isReception = roles.includes('receptionist') || roles.includes('helpdesk')
  const isPrivileged = roles.includes('admin') || roles.includes('super_admin') || roles.includes('general_secretary') || roles.includes('ags_admin') || roles.includes('ags_technical')
  return isReception && !isPrivileged
})

const isITOfficerOnly = computed(() => {
  const rawRoles = authStore.user?.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
  const isIT = roles.includes('it_officer')
  const isPrivileged = roles.includes('admin') || roles.includes('super_admin') || roles.includes('general_secretary') || roles.includes('ags_admin') || roles.includes('ags_technical')
  return isIT && !isPrivileged
})

const isEngineeringOfficerOnly = computed(() => {
  const rawRoles = authStore.user?.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
  const isEng = roles.includes('senior_engineer') ||
    roles.includes('assistant_engineer_civil') ||
    roles.includes('assistant_engineer_electrical') ||
    roles.includes('engineering_officer_civil') ||
    roles.includes('engineering_officer_electrical') ||
    roles.includes('plumber')
  const isPrivileged = roles.includes('admin') || roles.includes('super_admin') || roles.includes('general_secretary') || roles.includes('ags_admin') || roles.includes('ags_technical')
  return isEng && !isPrivileged
})

const isAccountantOnly = computed(() => {
  const rawRoles = authStore.user?.roles || []
  const roles = rawRoles.map(r => (typeof r === 'string' ? r : (r?.name || r?.role || ''))).filter(Boolean)
  const isAcc = roles.includes('accountant') || roles.includes('senior_accountant') || roles.includes('finance_department')
  const isPrivileged = roles.includes('admin') || roles.includes('super_admin') || roles.includes('general_secretary') || roles.includes('ags_admin') || roles.includes('ags_technical')
  return isAcc && !isPrivileged
})

function checkMobile() {
  isMobile.value = window.innerWidth < 1024
  if (isMobile.value) {
    sidebarCollapsed.value = true
  } else {
    sidebarCollapsed.value = localStorage.getItem('ncs_sidebar_collapsed') === '1'
  }
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  if (!isMobile.value) {
    localStorage.setItem('ncs_sidebar_collapsed', sidebarCollapsed.value ? '1' : '0')
  }
}

function refreshPage() {
  window.location.reload()
}

function handleSearch() {
  if (!searchQuery.value.trim()) return
  router.push({ path: '/nsmis/data/athletes', query: { q: searchQuery.value } })
}

function handleLogout() {
  authStore.clearAuth()
  try {
    authStore.logout()
  } catch {}
  window.location.href = '/login'
}

// Close dropdowns on outside click
function handleDocumentClick(e) {
  if (!e.target.closest('.main-navbar')) {
    profileOpen.value = false
    notificationsOpen.value = false
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  document.addEventListener('click', handleDocumentClick)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
  document.removeEventListener('click', handleDocumentClick)
})
</script>

<style scoped>
.animate-fadeIn {
  animation: fadeIn 0.15s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-6px); }
  to   { opacity: 1; transform: translateY(0); }
}
</style>
