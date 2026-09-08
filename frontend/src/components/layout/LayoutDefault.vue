<template>
  <div class="otika-app">
    <div class="main-wrapper main-wrapper-1" :class="{ 'sidebar-mini': sidebarCollapsed }">
      
      <!-- Navbar Background -->
      <div class="navbar-bg"></div>

      <!-- ── Top Navbar ──────────────────────────────────────────────── -->
      <nav class="main-navbar sticky">
        <!-- Left Side: Toggles & Search -->
        <div class="navbar-left">
          <!-- Sidebar Toggle Button -->
          <button
            type="button"
            class="cms-top-icon collapse-btn"
            @click="toggleSidebar"
            :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'"
          >
            <svg class="feather" viewBox="0 0 24 24" aria-hidden="true"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
          </button>

          <!-- Fullscreen Toggle Button -->
          <button
            type="button"
            class="cms-top-icon fullscreen-btn"
            @click="toggleFullscreen"
            title="Toggle Fullscreen"
          >
            <svg class="feather" viewBox="0 0 24 24" aria-hidden="true"><path d="M8 3H5a2 2 0 0 0-2 2v3M16 3h3a2 2 0 0 1 2 2v3M8 21H5a2 2 0 0 1-2-2v-3M16 21h3a2 2 0 0 0 2-2v-3"/></svg>
          </button>

          <!-- Otika Search Element -->
          <form class="search-element" @submit.prevent="handleSearch">
            <input
              type="text"
              v-model="searchQuery"
              class="form-control"
              placeholder="Search"
              @keyup.enter="handleSearch"
            />
            <button type="submit" class="btn" title="Search">
              <svg class="feather" viewBox="0 0 24 24" aria-hidden="true"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
            </button>
          </form>
        </div>

        <!-- Right Side: Messages, Notifications, Profile -->
        <div class="navbar-right">
          <div class="relative">
            <button
              type="button"
              class="cms-top-icon message-toggle"
              @click="messagesOpen = !messagesOpen; notificationsOpen = false; profileOpen = false"
              title="Messages"
              aria-label="Messages"
            >
              <svg class="feather" viewBox="0 0 24 24" aria-hidden="true"><path d="M4 4h16a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
              <span class="headerBadge1">6</span>
            </button>
            <div v-if="messagesOpen" class="otika-dropdown dropdown-list absolute right-0 mt-2 w-80 z-50 animate-fadeIn">
              <div class="dropdown-header">
                <span>Messages</span>
                <button type="button" @click="messagesOpen = false">Mark All As Read</button>
              </div>
              <div class="dropdown-list-content">
                <router-link to="/reception/leave/status" class="dropdown-message" @click="messagesOpen = false">
                  <span class="dropdown-avatar bg-blue-500">HR</span>
                  <span><strong>Human Resources</strong><small>Leave application updates</small><em>2 Min Ago</em></span>
                </router-link>
                <router-link to="/it/ppda/status" class="dropdown-message" @click="messagesOpen = false">
                  <span class="dropdown-avatar bg-emerald-500">PP</span>
                  <span><strong>Procurement Desk</strong><small>PPDA Form 5 status updates</small><em>5 Min Ago</em></span>
                </router-link>
                <router-link to="/me/activities" class="dropdown-message" @click="messagesOpen = false">
                  <span class="dropdown-avatar bg-amber-500">NC</span>
                  <span><strong>NCS Intranet</strong><small>Review your recent activity</small><em>12 Min Ago</em></span>
                </router-link>
              </div>
              <router-link to="/me/activities" class="dropdown-footer" @click="messagesOpen = false">
                View All <i class="icofont-rounded-right"></i>
              </router-link>
            </div>
          </div>

          <!-- Notifications Bell -->
          <div class="relative">
            <button
              type="button"
              class="cms-top-icon"
              @click="notificationsOpen = !notificationsOpen; messagesOpen = false; profileOpen = false"
              title="Notifications"
            >
              <svg class="feather bell" viewBox="0 0 24 24" aria-hidden="true"><path d="M18 8a6 6 0 0 0-12 0c0 7-3 7-3 9h18c0-2-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>
            </button>

            <!-- Notifications Dropdown -->
            <div
              v-if="notificationsOpen"
              class="otika-dropdown absolute right-0 mt-2 w-80 z-50 animate-fadeIn"
            >
              <div class="px-4 py-2 border-b border-slate-100 dark:border-slate-800 flex items-center justify-between">
                <span class="text-xs font-bold uppercase tracking-wider text-slate-800 dark:text-slate-200">System Notifications</span>
                <span class="text-[10px] bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full font-bold">3 New</span>
              </div>
              <div class="max-h-60 overflow-y-auto divide-y divide-slate-100 dark:divide-slate-800">
                <div class="px-4 py-2.5 hover:bg-slate-50 dark:hover:bg-slate-800/50 cursor-pointer transition-colors">
                  <div class="text-xs font-semibold text-slate-800 dark:text-slate-200 flex items-center justify-between">
                    <span>Form 5 Vetting Requisition</span>
                    <span class="text-[10px] text-slate-400">10m ago</span>
                  </div>
                  <div class="text-[11px] text-slate-500 mt-0.5">Procurement submitted PPDA Form 5 for approval</div>
                </div>
                <div class="px-4 py-2.5 hover:bg-slate-50 dark:hover:bg-slate-800/50 cursor-pointer transition-colors">
                  <div class="text-xs font-semibold text-slate-800 dark:text-slate-200 flex items-center justify-between">
                    <span>Lugogo Arena Inspection</span>
                    <span class="text-[10px] text-slate-400">1h ago</span>
                  </div>
                  <div class="text-[11px] text-slate-500 mt-0.5">UAF requested match readiness inspection</div>
                </div>
                <div class="px-4 py-2.5 hover:bg-slate-50 dark:hover:bg-slate-800/50 cursor-pointer transition-colors">
                  <div class="text-xs font-semibold text-slate-800 dark:text-slate-200 flex items-center justify-between">
                    <span>Fixed Asset Valuation Logged</span>
                    <span class="text-[10px] text-slate-400">2h ago</span>
                  </div>
                  <div class="text-[11px] text-slate-500 mt-0.5">Finance Department updated fixed asset ledger</div>
                </div>
              </div>
            </div>
          </div>

          <!-- User Profile Dropdown -->
          <div class="relative">
            <button
              type="button"
              class="nav-link-user"
              @click="profileOpen = !profileOpen; notificationsOpen = false; messagesOpen = false"
            >
              <div class="w-8 h-8 rounded-full bg-blue-600 text-white font-bold text-xs flex items-center justify-center border border-blue-400 shadow-sm">
                {{ userInitials }}
              </div>
              <div class="hidden sm:block text-left">
                <span class="block text-xs font-bold text-slate-800 dark:text-slate-200 leading-tight">{{ userName }}</span>
                <span class="block text-[10px] text-blue-600 dark:text-blue-400 font-semibold leading-tight">{{ userRole }}</span>
              </div>
              <i class="icofont-thin-down text-xs text-slate-400"></i>
            </button>

            <!-- Profile Menu Dropdown -->
            <div
              v-if="profileOpen"
              class="absolute right-0 mt-2 w-60 bg-white dark:bg-slate-900 rounded-xl shadow-xl border border-slate-200 dark:border-slate-800 py-2 z-50 animate-fadeIn"
            >
              <div class="px-4 py-2.5 border-b border-slate-100 dark:border-slate-800">
                <div class="text-xs font-extrabold text-slate-900 dark:text-white">{{ userName }}</div>
                <div class="text-[11px] text-slate-500 truncate mt-0.5">{{ authStore.user?.email }}</div>
                <span class="inline-block mt-1.5 px-2 py-0.5 rounded bg-blue-50 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300 text-[10px] font-bold">
                  {{ userRole }}
                </span>
              </div>
              <div class="py-1">
                <router-link
                  to="/profile"
                  class="flex items-center gap-2.5 px-4 py-2 text-xs text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors font-medium"
                  @click="profileOpen = false"
                >
                  <i class="icofont-user text-sm text-blue-600"></i> My Account Profile
                </router-link>
                <router-link
                  to="/me/security"
                  class="flex items-center gap-2.5 px-4 py-2 text-xs text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors font-medium"
                  @click="profileOpen = false"
                >
                  <i class="icofont-lock text-sm text-amber-600"></i> Security & Password
                </router-link>
                <router-link
                  to="/audit-logs"
                  class="flex items-center gap-2.5 px-4 py-2 text-xs text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors font-medium"
                  @click="profileOpen = false"
                >
                  <i class="icofont-history text-sm text-emerald-600"></i> Audit Activity Logs
                </router-link>
              </div>
              <div class="border-t border-slate-100 dark:border-slate-800 my-1"></div>
              <button
                type="button"
                class="w-full flex items-center gap-2.5 px-4 py-2 text-xs text-red-600 hover:bg-red-50 dark:hover:bg-red-950/30 transition-colors text-left font-bold"
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

            <template v-if="canUseExpenses">
              <li class="menu-header">Expense Management</li>
              <li :class="{ active: currentPath === '/expenses' || (currentPath.startsWith('/expenses/') && !['/expenses/new','/expenses/categories'].includes(currentPath)) }"><router-link to="/expenses" title="Expense Register"><i class="icofont-money"></i><span v-if="!sidebarCollapsed">Expenses</span></router-link></li>
              <li :class="{ active: currentPath === '/expenses/new' }"><router-link to="/expenses/new" title="Record Expense"><i class="icofont-plus-circle"></i><span v-if="!sidebarCollapsed">Record Expense</span></router-link></li>
              <li v-if="authStore.isAdmin" :class="{ active: currentPath === '/expenses/categories' }"><router-link to="/expenses/categories" title="Expense Categories"><i class="icofont-list"></i><span v-if="!sidebarCollapsed">Expense Categories</span></router-link></li>
            </template>

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
              <!-- Fixed Assets Register Dropdown -->
              <li :class="{ active: currentPath.startsWith('/fixed-assets'), dropdown: true }">
                <a href="javascript:void(0)"
                   class="nav-link has-dropdown flex items-center justify-between"
                   title="Fixed Assets Register"
                   @click.prevent="fixedAssetsOpen = !fixedAssetsOpen">
                  <div class="flex items-center gap-3">
                    <i class="icofont-building-alt"></i>
                    <span v-if="!sidebarCollapsed">Fixed Assets Register</span>
                  </div>
                  <i v-if="!sidebarCollapsed"
                     :class="fixedAssetsOpen ? 'icofont-rounded-up' : 'icofont-rounded-down'"
                     class="text-xs transition-transform duration-200"></i>
                </a>
                <ul class="dropdown-menu pl-4 space-y-1 mt-1" :style="{ display: fixedAssetsOpen ? 'block' : 'none' }">
                  <li :class="{ active: currentPath === '/fixed-assets' }">
                    <router-link to="/fixed-assets" class="nav-link flex items-center gap-2 py-1.5 px-3 text-xs rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800">
                      <i class="icofont-listine-dots text-xs"></i>
                      <span v-if="!sidebarCollapsed">Manage Assets</span>
                    </router-link>
                  </li>
                  <li :class="{ active: currentPath === '/fixed-assets/value-adjustments' }">
                    <router-link to="/fixed-assets/value-adjustments" class="nav-link flex items-center gap-2 py-1.5 px-3 text-xs rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800">
                      <i class="icofont-history text-xs"></i>
                      <span v-if="!sidebarCollapsed">Value Adjustments</span>
                    </router-link>
                  </li>
                  <li :class="{ active: currentPath === '/fixed-assets/pivot-engine' }">
                    <router-link to="/fixed-assets/pivot-engine" class="nav-link flex items-center gap-2 py-1.5 px-3 text-xs rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800">
                      <i class="icofont-chart-histogram text-xs"></i>
                      <span v-if="!sidebarCollapsed">Dynamic Pivot Engine</span>
                    </router-link>
                  </li>
                </ul>
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
              <!-- Fixed Assets Register Dropdown -->
              <li :class="{ active: currentPath.startsWith('/fixed-assets'), dropdown: true }">
                <a href="javascript:void(0)"
                   class="nav-link has-dropdown flex items-center justify-between"
                   title="Fixed Assets Register"
                   @click.prevent="fixedAssetsOpen = !fixedAssetsOpen">
                  <div class="flex items-center gap-3">
                    <i class="icofont-building-alt"></i>
                    <span v-if="!sidebarCollapsed">Fixed Assets Register</span>
                  </div>
                  <i v-if="!sidebarCollapsed"
                     :class="fixedAssetsOpen ? 'icofont-rounded-up' : 'icofont-rounded-down'"
                     class="text-xs transition-transform duration-200"></i>
                </a>
                <ul class="dropdown-menu pl-4 space-y-1 mt-1" :style="{ display: fixedAssetsOpen ? 'block' : 'none' }">
                  <li :class="{ active: currentPath === '/fixed-assets' }">
                    <router-link to="/fixed-assets" class="nav-link flex items-center gap-2 py-1.5 px-3 text-xs rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800">
                      <i class="icofont-listine-dots text-xs"></i>
                      <span v-if="!sidebarCollapsed">Manage Assets</span>
                    </router-link>
                  </li>
                  <li :class="{ active: currentPath === '/fixed-assets/value-adjustments' }">
                    <router-link to="/fixed-assets/value-adjustments" class="nav-link flex items-center gap-2 py-1.5 px-3 text-xs rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800">
                      <i class="icofont-history text-xs"></i>
                      <span v-if="!sidebarCollapsed">Value Adjustments</span>
                    </router-link>
                  </li>
                  <li :class="{ active: currentPath === '/fixed-assets/pivot-engine' }">
                    <router-link to="/fixed-assets/pivot-engine" class="nav-link flex items-center gap-2 py-1.5 px-3 text-xs rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800">
                      <i class="icofont-chart-histogram text-xs"></i>
                      <span v-if="!sidebarCollapsed">Dynamic Pivot Engine</span>
                    </router-link>
                  </li>
                </ul>
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
          <div class="sidebar-theme-control" :class="{ 'is-collapsed': sidebarCollapsed }">
            <ThemeToggle />
          </div>
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
import { expenseRoles } from '@/api/expenses'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import { useBreadcrumbStore } from '@/stores/breadcrumb.js'
import ThemeToggle from '@/components/ui/ThemeToggle.vue'

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
const messagesOpen = ref(false)
const reportsDropdownOpen = ref(false)
const searchQuery = ref('')
const isMobile = ref(false)
const fixedAssetsOpen = ref(true)

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

const canUseExpenses = computed(() => (authStore.user?.roles || []).some(role => expenseRoles.includes(typeof role === 'string' ? role : role?.name || role?.role)))

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

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen().catch(() => {})
  } else {
    document.exitFullscreen().catch(() => {})
  }
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
    messagesOpen.value = false
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
