<template>
  <main class="otika-cms">
    <div class="otika-app">
      <div class="main-wrapper main-wrapper-1" :class="{ 'sidebar-mini': sidebarCollapsed }">
        <div class="navbar-bg"></div>

        <!-- Header / Top Navbar -->
        <nav class="navbar navbar-expand-lg main-navbar sticky">
          <div class="form-inline mr-auto">
            <ul class="navbar-nav mr-3">
              <li>
                <button type="button" class="nav-link nav-link-lg cms-top-icon collapse-btn" @click="sidebarCollapsed = !sidebarCollapsed" :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'">
                  <i class="icofont-navigation-menu"></i>
                </button>
              </li>
              <li>
                <button type="button" class="nav-link nav-link-lg cms-top-icon fullscreen-btn" @click="resetDashboard" title="Refresh">
                  <i class="icofont-refresh"></i>
                </button>
              </li>
            </ul>
          </div>
          <ul class="navbar-nav navbar-right align-items-center">
            <!-- Notifications Dropdown -->
            <li class="dropdown dropdown-list-toggle" :class="{ show: notificationsOpen }">
              <button type="button" class="nav-link nav-link-lg message-toggle cms-top-icon" @click="notificationsOpen = !notificationsOpen" title="Notifications" style="position: relative;">
                <i class="icofont-notification"></i>
                <span v-if="unreadNotificationsCount > 0" class="badge badge-primary" style="position: absolute; top: 12px; right: 2px; padding: 2px 4px; font-size: 8px;">
                  {{ unreadNotificationsCount }}
                </span>
              </button>
              <div class="dropdown-menu dropdown-list dropdown-menu-right pullDown" :class="{ show: notificationsOpen }" style="width: 320px;">
                <div class="dropdown-header">
                  Notifications
                  <div class="float-right">
                    <button type="button" class="btn btn-link btn-sm p-0 text-primary" style="font-size: 10px; text-decoration: none;" @click="markAllRead">Mark All As Read</button>
                  </div>
                </div>
                <div class="dropdown-list-content dropdown-list-icons" style="max-height: 280px; overflow-y: auto; padding: 10px;">
                  <div v-if="notifications.length === 0" class="text-center py-4 text-muted small">
                    No new notifications.
                  </div>
                  <div v-else v-for="notif in notifications" :key="notif.id" class="dropdown-item d-flex align-items-start p-2 border-bottom" style="gap: 10px; cursor: pointer;" @click="goNotificationsPage">
                    <div class="p-2 rounded-circle" :style="notif.status === 'unread' ? 'background: #eff6ff;' : 'background: #f8fafc;'" style="font-size: 14px;">
                      <i :class="notif.icon_key || 'icofont-notification'" :style="notif.status === 'unread' ? 'color: #2563eb;' : 'color: #94a3b8;'"></i>
                    </div>
                    <div>
                      <strong class="d-block text-xs" :style="notif.status === 'unread' ? 'color: #1e293b; font-weight: 700;' : 'color: #64748b; font-weight: 400;'">{{ notif.title }}</strong>
                      <span class="d-block text-[10px] text-slate-400 mt-1 leading-tight">{{ notif.message }}</span>
                    </div>
                  </div>
                </div>
                <div class="dropdown-footer text-center border-top">
                  <button type="button" class="btn btn-link btn-sm text-xs" @click="goNotificationsPage">View All Notifications</button>
                </div>
              </div>
            </li>
            <li>
              <ThemeToggle class="nav-link nav-link-lg cms-top-icon portal-theme-toggle" />
            </li>
            <li class="dropdown" :class="{ show: profileOpen }">
              <button type="button" class="nav-link dropdown-toggle nav-link-lg nav-link-user cms-top-icon" @click="profileOpen = !profileOpen">
                <img alt="" src="/otika-assets/img/user.png" class="user-img-radious-style" />
                <span class="d-sm-none d-lg-inline-block"></span>
              </button>
              <div class="dropdown-menu dropdown-menu-right pullDown" :class="{ show: profileOpen }">
                <div class="dropdown-title">Hello {{ user.first_name }}</div>
                <button type="button" class="dropdown-item has-icon text-danger" @click="logout">
                  <i class="fas fa-sign-out-alt"></i> Logout
                </button>
              </div>
            </li>
          </ul>
        </nav>

        <!-- Sidebar Navigation -->
        <div class="main-sidebar sidebar-style-2">
          <aside id="sidebar-wrapper">
            <div class="sidebar-brand">
              <router-link to="/portal/helpdesk">
                <img alt="NCS" src="/main-logo.png" class="header-logo" />
              </router-link>
            </div>
            <ul class="sidebar-menu">
              <li class="menu-header">Helpdesk Workspace</li>
              <li :class="{ active: activeTab === 'overview' }">
                <button type="button" class="nav-link" @click="activeTab = 'overview'">
                  <i class="icofont-dashboard"></i><span>Dashboard Overview</span>
                </button>
              </li>
              <li class="menu-header">Universal Applications</li>
              <li>
                <router-link to="/reception/leave/apply" class="nav-link">
                  <i class="icofont-calendar"></i><span>Leave Application</span>
                </router-link>
              </li>
              <li>
                <router-link to="/it/ppda/new" class="nav-link">
                  <i class="icofont-document-folder"></i><span>PPDA Application</span>
                </router-link>
              </li>
              <li class="menu-header">Helpdesk Workspace</li>
              <li :class="{ active: activeTab === 'queue' }">
                <button type="button" class="nav-link" @click="activeTab = 'queue'">
                  <i class="icofont-ticket"></i><span>Issues Queue</span>
                </button>
              </li>
              <li :class="{ active: activeTab === 'leave-approvals' }">
                <button type="button" class="nav-link" @click="activeTab = 'leave-approvals'">
                  <i class="icofont-calendar"></i><span>Leave Approvals</span>
                </button>
              </li>
              <li :class="{ active: activeTab === 'visitor-clearance' }">
                <button type="button" class="nav-link" @click="activeTab = 'visitor-clearance'">
                  <i class="icofont-badge"></i><span>Visitors Clearance Form</span>
                </button>
              </li>
              <li :class="{ active: activeTab === 'reports' }">
                <button type="button" class="nav-link" @click="activeTab = 'reports'">
                  <i class="icofont-chart-histogram"></i><span>Reports & Analytics</span>
                </button>
              </li>
              <li :class="{ active: activeTab === 'notifications' }">
                <button type="button" class="nav-link" @click="activeTab = 'notifications'">
                  <i class="icofont-notification"></i><span>Notifications</span>
                  <span v-if="unreadNotificationsCount > 0" class="badge badge-primary">{{ unreadNotificationsCount }}</span>
                </button>
              </li>
              <li :class="{ active: activeTab === 'activities' }">
                <button type="button" class="nav-link" @click="activeTab = 'activities'">
                  <i class="icofont-history"></i><span>My Activities</span>
                </button>
              </li>
              <li class="menu-header">Session</li>
              <li>
                <button type="button" class="nav-link text-danger" @click="logout">
                  <i class="icofont-logout text-danger"></i><span>Log Out</span>
                </button>
              </li>
            </ul>
          </aside>
        </div>

        <!-- Main Content Area -->
        <div class="main-content">
          <section class="section">
            <div class="section-body">
              <!-- Tab 1: Dashboard Overview -->
              <div v-if="activeTab === 'overview'" class="space-y-6">
                <!-- Welcome/SLA banner -->
                <div class="card bg-gradient-to-r from-emerald-900 via-teal-900 to-slate-900 border-0 text-white p-6 shadow-sm rounded-lg mb-6">
                  <h2 class="text-xl font-bold">Helpdesk Support Hub</h2>
                  <p class="text-sm text-slate-300 mt-1">Manage and resolve user feedback, system complaints, and application payment support issues.</p>
                </div>

                <!-- KPI Grid -->
                <div class="row">
                  <div class="col-xl-3 col-lg-6 col-md-6 col-sm-6 col-xs-12">
                    <div class="card card-statistic-4 hover:shadow-lg transition-shadow">
                      <div class="card-body">
                        <h5 class="text-xs text-slate-400 font-semibold uppercase">Open Tickets</h5>
                        <h2 class="text-2xl font-bold mt-2 text-slate-800 dark:text-slate-100">{{ tickets.length }}</h2>
                        <p class="text-xs text-rose-500 mt-1">2 urgent issues</p>
                      </div>
                    </div>
                  </div>
                  <div class="col-xl-3 col-lg-6 col-md-6 col-sm-6 col-xs-12">
                    <div class="card card-statistic-4 hover:shadow-lg transition-shadow">
                      <div class="card-body">
                        <h5 class="text-xs text-slate-400 font-semibold uppercase">Avg Response Time</h5>
                        <h2 class="text-2xl font-bold mt-2 text-slate-800 dark:text-slate-100">14m</h2>
                        <p class="text-xs text-emerald-500 mt-1">Within SLA targets</p>
                      </div>
                    </div>
                  </div>
                  <div class="col-xl-3 col-lg-6 col-md-6 col-sm-6 col-xs-12">
                    <div class="card card-statistic-4 hover:shadow-lg transition-shadow">
                      <div class="card-body">
                        <h5 class="text-xs text-slate-400 font-semibold uppercase">Resolved Today</h5>
                        <h2 class="text-2xl font-bold mt-2 text-slate-800 dark:text-slate-100">19</h2>
                        <p class="text-xs text-emerald-500 mt-1">95% resolution rate</p>
                      </div>
                    </div>
                  </div>
                  <div class="col-xl-3 col-lg-6 col-md-6 col-sm-6 col-xs-12">
                    <div class="card card-statistic-4 hover:shadow-lg transition-shadow">
                      <div class="card-body">
                        <h5 class="text-xs text-slate-400 font-semibold uppercase">Active Agents</h5>
                        <h2 class="text-2xl font-bold mt-2 text-slate-800 dark:text-slate-100">3 / 4</h2>
                        <p class="text-xs text-emerald-500 mt-1">Online now</p>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Recent Tickets Snapshot -->
                <div class="row mt-4">
                  <div class="col-12">
                    <div class="card p-4">
                      <div class="card-header p-0 pb-3 mb-3 border-b border-slate-100 dark:border-slate-800 flex justify-between items-center">
                        <h4 class="text-sm font-bold text-slate-800 dark:text-slate-100">High Priority Issues</h4>
                        <button class="btn btn-sm btn-link" @click="activeTab = 'queue'">Go to queue</button>
                      </div>
                      <div class="space-y-3">
                        <div v-for="ticket in tickets.filter(t => t.priority === 'Urgent')" :key="ticket.id" class="p-3 rounded border border-slate-100 dark:border-slate-800 bg-rose-50/10 dark:bg-rose-950/10 flex justify-between items-center">
                          <div>
                            <p class="text-sm font-semibold text-rose-600 dark:text-rose-400">{{ ticket.subject }}</p>
                            <p class="text-xs text-slate-400">By {{ ticket.user }} • {{ ticket.time }}</p>
                          </div>
                          <button class="btn btn-xs btn-outline-danger" @click="activeTab = 'queue'; selectedTicket = ticket">Respond</button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Tab 2: Issues Queue -->
              <div v-else-if="activeTab === 'queue'" class="row">
                <!-- Sidebar: Queue List -->
                <div class="col-lg-4 col-md-12 mb-4">
                  <div class="card p-4 h-[550px] flex flex-col">
                    <h3 class="text-base font-bold text-slate-800 dark:text-slate-100 mb-3">Ticket Queue</h3>
                    <div class="flex gap-2 mb-3">
                      <button @click="filter = 'all'" :class="filter === 'all' ? 'btn-primary' : 'btn-outline-light'" class="btn btn-xs">All</button>
                      <button @click="filter = 'Urgent'" :class="filter === 'Urgent' ? 'btn-danger' : 'btn-outline-light'" class="btn btn-xs">Urgent</button>
                    </div>
                    <div class="overflow-y-auto flex-1 space-y-3 pr-1">
                      <div
                        v-for="ticket in filteredTickets"
                        :key="ticket.id"
                        @click="selectTicket(ticket)"
                        :class="[
                          selectedTicket?.id === ticket.id ? 'border-emerald-500 bg-emerald-50/30 dark:bg-emerald-950/20' : 'border-slate-100 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-900/10',
                          'p-3.5 rounded border cursor-pointer hover:border-slate-300 dark:hover:border-slate-700 transition-all space-y-2'
                        ]"
                      >
                        <div class="flex justify-between items-center">
                          <span :class="ticket.priority === 'Urgent' ? 'badge badge-danger' : 'badge badge-light'" class="text-[9px] uppercase font-bold">
                            {{ ticket.priority }}
                          </span>
                          <span class="text-[10px] text-slate-400">{{ ticket.time }}</span>
                        </div>
                        <p class="text-xs font-semibold text-slate-800 dark:text-slate-200 truncate">{{ ticket.subject }}</p>
                        <div class="text-[11px] text-slate-400 flex justify-between">
                          <span>{{ ticket.user }}</span>
                          <span class="text-emerald-500 font-medium">Open</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Main Panel: Ticket Detail -->
                <div class="col-lg-8 col-md-12">
                  <div class="card p-5 h-[550px] flex flex-col justify-between">
                    <div v-if="selectedTicket" class="flex-1 flex flex-col space-y-4 overflow-hidden">
                      <div class="border-b border-slate-100 dark:border-slate-800 pb-3 flex justify-between items-start">
                        <div>
                          <h4 class="text-base font-bold text-slate-800 dark:text-slate-100">{{ selectedTicket.subject }}</h4>
                          <p class="text-xs text-slate-400 mt-1">Submitted by: <strong>{{ selectedTicket.user }}</strong> ({{ selectedTicket.email }})</p>
                        </div>
                        <button class="btn btn-success btn-xs" @click="closeTicket">Resolve Ticket</button>
                      </div>

                      <div class="flex-1 overflow-y-auto space-y-3 pr-1">
                        <!-- Initial Message -->
                        <div class="bg-slate-50 dark:bg-slate-950 p-4 rounded border border-slate-100 dark:border-slate-800 space-y-1">
                          <div class="flex justify-between text-xs text-slate-400">
                            <span>{{ selectedTicket.user }}</span>
                            <span>{{ selectedTicket.time }}</span>
                          </div>
                          <p class="text-sm text-slate-700 dark:text-slate-300 mt-2">{{ selectedTicket.message }}</p>
                        </div>

                        <!-- Replies timeline -->
                        <div v-for="reply in replies" :key="reply.id" class="bg-slate-100/50 dark:bg-slate-800/40 p-4 rounded border border-slate-100 dark:border-slate-800 space-y-1 ml-6">
                          <div class="flex justify-between text-xs text-indigo-500">
                            <span>{{ user.first_name }} (Agent)</span>
                            <span>Just now</span>
                          </div>
                          <p class="text-sm text-slate-700 dark:text-slate-300 mt-2">{{ reply.text }}</p>
                        </div>
                      </div>

                      <!-- Reply Input -->
                      <div class="pt-3 border-t border-slate-100 dark:border-slate-800">
                        <form @submit.prevent="sendReply" class="flex gap-3">
                          <input v-model="replyText" type="text" placeholder="Type a response to resolve the client query..." class="form-control text-sm" required />
                          <button type="submit" class="btn btn-primary btn-sm">Send</button>
                        </form>
                      </div>
                    </div>

                    <!-- Empty State -->
                    <div v-else class="flex-1 flex flex-col items-center justify-center text-center space-y-3">
                      <i class="icofont-ticket text-5xl text-slate-300"></i>
                      <h4 class="text-sm font-semibold text-slate-600 dark:text-slate-400">No Ticket Selected</h4>
                      <p class="text-xs text-slate-400 max-w-xs">Select a support case from the queue to start responding.</p>
                    </div>
                  </div>
                </div>
              </div>

              <div v-else-if="activeTab === 'leave-approvals'">
                <LeaveApprovalsPanel />
              </div>

              <div v-else-if="activeTab === 'visitor-clearance'">
                <VisitorClearancePanel />
              </div>

              <div v-else-if="activeTab === 'reports'">
                <ReportsPanel />
              </div>

              <!-- Tab 5: Notifications Page -->
              <div v-else-if="activeTab === 'notifications'" class="card p-5">
                <div class="d-flex justify-content-between align-items-center mb-4">
                  <div>
                    <h3 class="text-base font-bold text-slate-800 dark:text-slate-100 mb-1">My Notifications</h3>
                    <p class="text-xs text-slate-400 mb-0">Stay updated on leave approvals, support ticket assignments, and portal logs.</p>
                  </div>
                  <button type="button" class="btn btn-outline-primary btn-sm" @click="markAllRead">
                    Mark all as read
                  </button>
                </div>

                <div v-if="notifications.length === 0" class="text-center py-5">
                  <i class="icofont-notification text-slate-300" style="font-size: 48px;"></i>
                  <p class="text-slate-400 text-sm mt-3">You don't have any notifications at the moment.</p>
                </div>
                <div v-else class="space-y-3">
                  <div v-for="notif in notifications" :key="notif.id" class="p-3 rounded border d-flex justify-content-between align-items-center" :style="notif.status === 'unread' ? 'border-color: #2563eb; background: rgba(37, 99, 235, 0.05);' : 'border-color: #f1f5f9;'">
                    <div class="d-flex align-items-start gap-3">
                      <div class="p-2 rounded-circle bg-light dark:bg-slate-800" style="font-size: 16px;">
                        <i :class="[notif.icon_key || 'icofont-notification', notif.status === 'unread' ? 'text-primary' : 'text-slate-400']"></i>
                      </div>
                      <div>
                        <h4 class="text-sm font-semibold text-slate-800 dark:text-slate-200 mb-1">
                          {{ notif.title }}
                          <span v-if="notif.status === 'unread'" class="badge badge-primary ml-2" style="font-size: 9px; padding: 2px 5px;">New</span>
                        </h4>
                        <p class="text-xs text-slate-500 dark:text-slate-400 mb-0">{{ notif.message }}</p>
                      </div>
                    </div>
                    <button v-if="notif.status === 'unread'" type="button" class="btn btn-link btn-xs text-primary" @click="markSingleRead(notif.id)">
                      Mark read
                    </button>
                  </div>
                </div>
              </div>

              <!-- Tab 6: My Activities Page -->
              <div v-else-if="activeTab === 'activities'" class="card p-5">
                <div class="d-flex justify-content-between align-items-center mb-4">
                  <div>
                    <h3 class="text-base font-bold text-slate-800 dark:text-slate-100 mb-1">My Activities & Audit Trail</h3>
                    <p class="text-xs text-slate-400 mb-0">Verifiable write-once-read-many log pipeline complying with platform security policies.</p>
                  </div>
                  <button type="button" class="btn btn-outline-secondary btn-sm" @click="loadAuditLogs">
                    <i class="icofont-refresh"></i> Refresh Logs
                  </button>
                </div>
                
                <div class="table-responsive">
                  <table class="table table-bordered table-striped official-table text-xs">
                    <thead>
                      <tr>
                        <th>Timestamp (UTC)</th>
                        <th>Event & Action</th>
                        <th>Actor Principal</th>
                        <th>Client Medium / IP</th>
                        <th>Target Resource</th>
                        <th>Security Signature</th>
                        <th class="text-center">Action</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="log in auditLogsList" :key="log.event_id">
                        <td>
                          <span class="font-mono text-slate-600 dark:text-slate-300">{{ formatDateLong(log.timestamp) }}</span>
                        </td>
                        <td>
                          <div class="font-bold text-indigo-600">{{ log.action }}</div>
                          <span class="badge" :class="log.status === 'SUCCESS' ? 'badge-success' : 'badge-danger'" style="font-size: 8px; padding: 2px 4px;">{{ log.status }}</span>
                          <span class="badge badge-info ml-1" style="font-size: 8px; padding: 2px 4px;">{{ log.severity }}</span>
                        </td>
                        <td>
                          <div class="font-semibold">{{ log.actor_username }}</div>
                          <span class="text-[9px] text-slate-400 d-block">{{ log.actor_type }} • {{ log.actor_roles?.join(', ') }}</span>
                        </td>
                        <td>
                          <div><code>{{ log.client_ip }}</code></div>
                          <span class="text-[9px] text-slate-400 d-block">{{ log.access_medium }} • {{ log.parsed_os }}</span>
                        </td>
                        <td>
                          <div class="font-semibold">{{ log.resource_type }}</div>
                          <code class="text-[10px] text-slate-500">{{ log.resource_id }}</code>
                        </td>
                        <td>
                          <code class="text-[9px] text-slate-400" :title="log.signature">{{ log.signature ? log.signature.slice(0, 18) + '...' : 'N/A' }}</code>
                        </td>
                        <td class="text-center">
                          <button type="button" class="btn btn-xs btn-outline-primary" @click="selectedAudit = log">
                            Inspect
                          </button>
                        </td>
                      </tr>
                      <tr v-if="!auditLogsList.length">
                        <td colspan="7" class="text-center py-6 text-slate-400">No audit logs found in storage queue.</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <!-- Audit Log Inspection Modal Drawer -->
                <div v-if="selectedAudit" class="fixed inset-0 bg-slate-900/60 dark:bg-black/80 flex justify-end z-[999] no-print" @click.self="selectedAudit = null">
                  <div class="w-full max-w-2xl bg-white dark:bg-slate-900 h-full p-6 shadow-2xl overflow-y-auto flex flex-col space-y-6 animate-slide-left">
                    <div class="flex justify-between items-center pb-3 border-b border-slate-100 dark:border-slate-800">
                      <div>
                        <h4 class="text-sm font-bold text-slate-800 dark:text-slate-100 uppercase tracking-wider">Audit Log Details</h4>
                        <span class="text-[10px] font-mono text-slate-400">{{ selectedAudit.event_id }}</span>
                      </div>
                      <button type="button" class="btn btn-xs btn-outline-danger" @click="selectedAudit = null">&times; Close</button>
                    </div>

                    <div class="space-y-4 text-xs">
                      <div class="grid grid-cols-2 gap-4">
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Timestamp (UTC)</strong>
                          <span class="font-mono">{{ selectedAudit.timestamp }}</span>
                        </div>
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Origin Timezone</strong>
                          <span>{{ selectedAudit.timezone }}</span>
                        </div>
                      </div>

                      <div class="grid grid-cols-2 gap-4 border-t pt-3 border-slate-100 dark:border-slate-800">
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Correlation ID</strong>
                          <span class="font-mono text-slate-500">{{ selectedAudit.correlation_id }}</span>
                        </div>
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Operation Severity</strong>
                          <span class="badge badge-info">{{ selectedAudit.severity }}</span>
                        </div>
                      </div>

                      <div class="grid grid-cols-2 gap-4 border-t pt-3 border-slate-100 dark:border-slate-800">
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Action Route</strong>
                          <strong class="text-indigo-600 font-mono">{{ selectedAudit.action }}</strong>
                        </div>
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Client IP Address</strong>
                          <code>{{ selectedAudit.client_ip }}</code>
                        </div>
                      </div>

                      <div class="border-t pt-3 border-slate-100 dark:border-slate-800">
                        <strong class="text-slate-400 uppercase text-[9px] block">Full User-Agent String</strong>
                        <code class="text-[10px] leading-tight block bg-slate-50 dark:bg-slate-950 p-2 rounded border border-slate-100 dark:border-slate-900">{{ selectedAudit.user_agent_raw }}</code>
                      </div>

                      <div class="grid grid-cols-3 gap-3 border-t pt-3 border-slate-100 dark:border-slate-800">
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Client Browser</strong>
                          <span>{{ selectedAudit.parsed_browser }}</span>
                        </div>
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Operating System</strong>
                          <span>{{ selectedAudit.parsed_os }}</span>
                        </div>
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Access Medium</strong>
                          <span>{{ selectedAudit.access_medium }}</span>
                        </div>
                      </div>

                      <div class="grid grid-cols-2 gap-4 border-t pt-3 border-slate-100 dark:border-slate-800">
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Target Resource Type</strong>
                          <span>{{ selectedAudit.resource_type }}</span>
                        </div>
                        <div>
                          <strong class="text-slate-400 uppercase text-[9px] block">Target Resource ID</strong>
                          <code>{{ selectedAudit.resource_id }}</code>
                        </div>
                      </div>

                      <div class="border-t pt-3 border-slate-100 dark:border-slate-800">
                        <strong class="text-slate-400 uppercase text-[9px] block mb-1">State Payload Changes</strong>
                        <div class="grid grid-cols-2 gap-3">
                          <div>
                            <span class="text-[9px] font-bold text-slate-400 block mb-1">BEFORE CHANGES</span>
                            <pre class="bg-slate-50 dark:bg-slate-950 p-2 rounded text-[10px] border max-h-36 overflow-y-auto">{{ JSON.stringify(selectedAudit.payload_before, null, 2) }}</pre>
                          </div>
                          <div>
                            <span class="text-[9px] font-bold text-slate-400 block mb-1">AFTER CHANGES</span>
                            <pre class="bg-slate-50 dark:bg-slate-950 p-2 rounded text-[10px] border max-h-36 overflow-y-auto">{{ JSON.stringify(selectedAudit.payload_after, null, 2) }}</pre>
                          </div>
                        </div>
                      </div>

                      <div class="border-t pt-3 border-slate-100 dark:border-slate-800">
                        <strong class="text-slate-400 uppercase text-[9px] block">Cryptographic Log Signature (HMAC-SHA256)</strong>
                        <code class="text-[10px] break-all block bg-light-panel p-2 rounded font-mono text-emerald-600 dark:text-emerald-400 border border-dashed border-emerald-500/20">{{ selectedAudit.signature }}</code>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

            </div>
          </section>
        </div>

        <!-- App Footer -->
        <footer class="main-footer cms-main-footer">
          <div class="footer-left">
            Copyright &copy; 2026 <div class="bullet"></div> National Council of Sports, Uganda
          </div>
          <div class="footer-right">
            Helpdesk Module v1.0.0
          </div>
        </footer>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { storedPortalUser } from '@/utils/portalAuth.js'
import { ensureOtikaStyles } from '@/utils/otikaAssets.js'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'
import LeaveApprovalsPanel from '@/components/portal/LeaveApprovalsPanel.vue'
import VisitorClearancePanel from '@/components/portal/VisitorClearancePanel.vue'
import ReportsPanel from '@/components/portal/ReportsPanel.vue'
import * as cms from '@/api/cms.js'

import { createAuditLog } from '@/utils/auditLogger.js'

ensureOtikaStyles()

const router = useRouter()
const route = useRoute()
const user = ref(storedPortalUser())

const sidebarCollapsed = ref(false)
const profileOpen = ref(false)
const activeTab = ref(route.query.tab || 'overview')
const notificationsOpen = ref(false)
const notifications = ref([])

const auditLogsList = ref([])
const selectedAudit = ref(null)

onMounted(() => {
  fetchNotifications()
  loadAuditLogs()
})

function loadAuditLogs() {
  let list = []
  const stored = localStorage.getItem('ncsms_audit_logs')
  if (stored) {
    try {
      list = JSON.parse(stored)
    } catch (e) {}
  }
  
  if (list.length === 0) {
    // Generate initial seeds using standard logger
    const seed1 = createAuditLog({
      action: 'user.security:mfa_reset',
      status: 'SUCCESS',
      severity: 'WARN',
      actor: user.value,
      resourceId: 'usr_sarah_chemutai',
      resourceType: 'UserSecurity',
      payloadBefore: { mfa_enabled: true, mfa_keys_configured: true },
      payloadAfter: { mfa_enabled: false, mfa_keys_configured: false }
    })
    const seed2 = createAuditLog({
      action: 'user.security:session_unlock',
      status: 'SUCCESS',
      severity: 'INFO',
      actor: user.value,
      resourceId: 'usr_baker_nsubuga',
      resourceType: 'UserSession',
      payloadBefore: { session_locked: true, attempts_count: 5 },
      payloadAfter: { session_locked: false, attempts_count: 0 }
    })
    list = [seed1, seed2]
  }
  auditLogsList.value = list
}

function reset2FA() {
  createAuditLog({
    action: 'user.security:mfa_reset',
    status: 'SUCCESS',
    severity: 'WARN',
    actor: user.value,
    resourceId: selectedTicket.value ? selectedTicket.value.user : 'Sarah Chemutai',
    resourceType: 'UserSecurity',
    payloadBefore: { mfa_enabled: true, mfa_keys_configured: true },
    payloadAfter: { mfa_enabled: false, mfa_keys_configured: false }
  })
  loadAuditLogs()
  alert('MFA token keys have been successfully reset. The user will be requested to scan a new QR code on their next log in.')
}

function unlockSession() {
  createAuditLog({
    action: 'user.security:session_unlock',
    status: 'SUCCESS',
    severity: 'INFO',
    actor: user.value,
    resourceId: selectedTicket.value ? selectedTicket.value.user : 'Baker Nsubuga',
    resourceType: 'UserSession',
    payloadBefore: { session_locked: true, attempts_count: 5 },
    payloadAfter: { session_locked: false, attempts_count: 0 }
  })
  loadAuditLogs()
  alert('System security login attempts have been cleared. User profile is now unlocked.')
}

function formatDateLong(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

const unreadNotificationsCount = computed(() => {
  return notifications.value.filter(n => n.status === 'unread').length
})

async function fetchNotifications() {
  let list = []
  try {
    const res = await cms.listNotifications({ page: 1, per_page: 50 })
    const apiData = res?.data || res?.entries || res || []
    const parsedApi = (Array.isArray(apiData) ? apiData : []).map(n => ({
      id: n.id,
      title: n.title,
      message: n.message,
      status: n.status || 'unread',
      icon_key: n.icon_key || 'icofont-notification'
    }))
    list = [...parsedApi]
  } catch (e) {}

  const local = localStorage.getItem('ncsms_local_notifications')
  if (local) {
    try {
      const parsedLocal = JSON.parse(local)
      list = [...parsedLocal, ...list]
    } catch (e) {}
  }
  
  notifications.value = list
}

async function markAllRead() {
  try {
    await cms.markAllNotificationsRead()
  } catch (e) {}
  const local = localStorage.getItem('ncsms_local_notifications')
  if (local) {
    try {
      const parsed = JSON.parse(local)
      parsed.forEach(n => { n.status = 'read' })
      localStorage.setItem('ncsms_local_notifications', JSON.stringify(parsed))
    } catch (e) {}
  }
  await fetchNotifications()
}

async function markSingleRead(id) {
  if (typeof id === 'string' && id.startsWith('notif_')) {
    // Local storage notification
    const local = localStorage.getItem('ncsms_local_notifications')
    if (local) {
      try {
        const parsed = JSON.parse(local)
        const found = parsed.find(n => n.id === id)
        if (found) {
          found.status = 'read'
          localStorage.setItem('ncsms_local_notifications', JSON.stringify(parsed))
        }
      } catch (e) {}
    }
  } else {
    // API notification
    try {
      await cms.updateNotification(id, { status: 'read' })
    } catch (e) {}
  }
  await fetchNotifications()
}

function goNotificationsPage() {
  notificationsOpen.value = false
  activeTab.value = 'notifications'
}

watch(() => route.query.tab, (newTab) => {
  if (newTab) {
    activeTab.value = newTab
  }
})
const filter = ref('all')
const replyText = ref('')
const replies = ref([])
const selectedTicket = ref(null)

const activeTabLabel = computed(() => {
  switch (activeTab.value) {
    case 'queue': return 'Issues Queue'
    case 'leave-request':
    case 'leave-approvals': return 'Leave Approvals'
    case 'visitor-clearance': return 'Visitor Clearance'
    case 'reports': return 'Reports'
    case 'activities': return 'My Activities'
    default: return 'Overview'
  }
})

const tickets = ref([
  { id: 1, subject: 'Federation account locked out', user: 'Baker Nsubuga', email: 'b.nsubuga@fufa.co.ug', priority: 'Urgent', status: 'open', time: '10 mins ago', message: 'Hi support team, I entered my login details wrong three times and now our official FUFA profile is locked out. Please unlock it so we can submit the disbursement audit report by tomorrow.' },
  { id: 2, subject: 'Reset Google Authenticator 2FA seed', user: 'Sarah Chemutai', email: 's.chemutai@athleticsug.org', priority: 'Urgent', status: 'open', time: '22 mins ago', message: 'Hello, I got a new phone today and lost my authentication keys. Please reset the 2FA configurations on Sarah Chemutai federation profile.' },
  { id: 3, subject: 'Payment proof upload failed', user: 'Okello Dennis', email: 'dennis.okello@rowingug.com', priority: 'Normal', status: 'open', time: '1 hour ago', message: 'I tried uploading our bank slip proof of payment for the annual license renewal, but the system keeps throwing a file format error. The file is a PDF of size 12MB. Please help.' },
  { id: 4, subject: 'Incorrect club logo displayed', user: 'Musa Kirumira', email: 'mkirumira@expressfc.co.ug', priority: 'Low', status: 'open', time: '3 hours ago', message: 'Our Express FC page displays the logo from 2024. We recently updated the branding. How do we submit the new SVG logo for our club page?' }
])

const filteredTickets = computed(() => {
  if (filter.value === 'all') return tickets.value
  return tickets.value.filter(t => t.priority === filter.value)
})

const selectTicket = (ticket) => {
  selectedTicket.value = ticket
  replies.value = []
}

const sendReply = () => {
  if (!replyText.value.trim()) return
  replies.value.push({ id: replies.value.length + 1, text: replyText.value })
  replyText.value = ''
}

const closeTicket = () => {
  tickets.value = tickets.value.filter(t => t.id !== selectedTicket.value.id)
  selectedTicket.value = null
  replies.value = []
}

const showActionAlert = (message) => {
  alert(message)
}

const resetDashboard = () => {
  activeTab.value = 'overview'
}

const logout = () => {
  localStorage.removeItem('ncsms_access_token')
  localStorage.removeItem('ncsms_user')
  router.push('/login')
}
</script>

<style scoped>
.navbar-bg {
  z-index: 1 !important;
}
.main-sidebar {
  z-index: 5 !important;
}
</style>
