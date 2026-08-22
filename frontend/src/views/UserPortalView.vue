<template>
  <main class="otika-cms user-portal">
    <div class="otika-app">
      <div class="main-wrapper main-wrapper-1" :class="{ 'sidebar-mini': sidebarCollapsed }">
        <div class="navbar-bg"></div>
        <nav class="navbar navbar-expand-lg main-navbar sticky">
          <div class="form-inline mr-auto">
            <ul class="navbar-nav mr-3">
              <li>
                <button type="button" class="nav-link nav-link-lg cms-top-icon collapse-btn" :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'" @click="onMenuToggle">
                  <i class="icofont-navigation-menu"></i>
                </button>
              </li>
              <li>
                <button type="button" class="nav-link nav-link-lg cms-top-icon fullscreen-btn" :disabled="loading" title="Refresh dashboard" @click="loadPortal">
                  <i class="icofont-refresh"></i>
                </button>
              </li>
              <li class="portal-navbar-title">
                <small>{{ userProfileLabel }}</small>
                <strong>{{ sectionTitle }}</strong>
              </li>
            </ul>
          </div>
          <ul class="navbar-nav navbar-right">
            <li>
              <ThemeToggle class="nav-link nav-link-lg cms-top-icon portal-theme-toggle" />
            </li>
            <li>
              <button type="button" class="nav-link nav-link-lg cms-top-icon portal-top-action" title="Messages" @click="select('messages')">
                <i class="icofont-envelope"></i>
                <span v-if="messageCount" class="badge headerBadge1">{{ messageCount }}</span>
              </button>
            </li>
            <li>
              <button type="button" class="nav-link nav-link-lg cms-top-icon portal-top-action" title="Notifications" @click="select('notifications')">
                <i class="icofont-notification"></i>
                <span v-if="unreadNotifications" class="badge headerBadge2">{{ unreadNotifications }}</span>
              </button>
            </li>
            <li class="dropdown" :class="{ show: profileOpen }">
              <button type="button" class="nav-link dropdown-toggle nav-link-lg nav-link-user cms-top-icon" @click="profileOpen = !profileOpen">
                <img v-if="profileAvatar" :src="profileAvatar" alt="" class="user-img-radious-style" />
                <span v-else class="nav-avatar-fallback">{{ initials }}</span>
              </button>
              <div class="dropdown-menu dropdown-menu-right pullDown" :class="{ show: profileOpen }">
                <div class="dropdown-title">Hello {{ firstName }}</div>
                <button type="button" class="dropdown-item has-icon" @click="select('profile')"><i class="fas fa-user"></i> My Profile</button>
                <button type="button" class="dropdown-item has-icon" @click="select('settings')"><i class="fas fa-cog"></i> Settings</button>
                <button type="button" class="dropdown-item has-icon" @click="select('activities')"><i class="fas fa-history"></i> My Activities</button>
                <button type="button" class="dropdown-item has-icon" @click="select('applications')"><i class="fas fa-file-alt"></i> My Applications</button>
                <button type="button" class="dropdown-item has-icon" @click="select('transactions')"><i class="fas fa-receipt"></i> Transactions</button>
                <div class="dropdown-divider"></div>
                <button type="button" class="dropdown-item has-icon text-danger" @click="logout"><i class="fas fa-sign-out-alt"></i> Logout</button>
              </div>
            </li>
          </ul>
        </nav>

        <div class="main-sidebar sidebar-style-2" :class="{ 'mobile-sidebar-open': mobileSidebarOpen }">
          <aside id="sidebar-wrapper">
            <div class="sidebar-brand">
              <router-link to="/dashboard" @click="select('dashboard')">
                <img alt="NCS" src="/main-logo.png" class="header-logo" />
              </router-link>
            </div>
            <div class="sidebar-user">
              <img v-if="profileAvatar" :src="profileAvatar" alt="" class="user-img-radious-style" />
              <span v-else class="sidebar-avatar-fallback">{{ initials }}</span>
              <div>
                <strong>{{ sidebarUserName }}</strong>
                <span>{{ userWorkspaceLabel }}</span>
              </div>
            </div>
            <ul class="sidebar-menu" aria-label="Applicant navigation">
              <li class="menu-header">{{ userProfileLabel }}</li>
              <li v-for="item in navigation" :key="item.id" :class="{ active: section === item.id }">
                <button type="button" class="nav-link" :title="item.label" @click="select(item.id)">
                  <i :class="item.icon"></i>
                  <span>{{ item.label }}</span>
                  <b v-if="item.badge" class="portal-nav-badge">{{ item.badge }}</b>
                </button>
              </li>
              <li class="menu-header">Account</li>
              <li>
                <button type="button" class="nav-link" title="Sign out" @click="logout">
                  <i class="icofont-logout"></i>
                  <span>Sign out</span>
                </button>
              </li>
            </ul>
          </aside>
        </div>
        <div v-if="mobileSidebarOpen" class="sidebar-scrim" @click="mobileSidebarOpen = false"></div>

        <div class="main-content">
          <section class="section">
            <div class="section-body">
              <div class="section-header">
                <h1>{{ sectionTitle }}</h1>
                <div class="section-header-breadcrumb">
                  <div class="breadcrumb-item active"><router-link to="/dashboard">Dashboard</router-link></div>
                  <div class="breadcrumb-item">{{ sectionTitle }}</div>
                </div>
              </div>
              <div class="cms-actions otika-page-actions">
                <button type="button" class="btn btn-icon icon-left btn-primary" @click="select('apply')"><i class="fas fa-plus"></i> Apply Now</button>
                <button type="button" class="btn btn-icon icon-left btn-info" :disabled="loading" @click="loadPortal"><i class="fas fa-sync"></i> Refresh</button>
              </div>

              <p v-if="success" class="portal-success cms-message">{{ success }}</p>
              <p v-if="error" class="portal-error cms-error">{{ error }}</p>

              <template v-if="section === 'dashboard'">
                <header class="page-heading">
                  <div><p>{{ greeting }}, {{ firstName }}</p><h1>Your applications at a glance</h1><span>Track progress and start an application from the open forms below.</span></div>
                  <button type="button" class="primary-command" @click="select('apply')"><i class="icofont-plus-circle"></i> Apply now</button>
                </header>
                <div class="user-kpis">
                  <article v-for="card in userKpis" :key="card.label">
                    <span :class="card.tone"><i :class="card.icon"></i></span>
                    <div><small>{{ card.label }}</small><strong>{{ card.value }}</strong><p>{{ card.note }}</p></div>
                  </article>
                </div>
                <OpenFormsPanel
                  title="Open Applications"
                  :forms="openForms"
                  :submissions="dynamicSubmissions"
                  :loading="loading"
                  @start="startApplication"
                  @view-all="select('apply')"
                />

                <!-- Athlete Dashboard Extension Panel -->
                <div v-if="isAthlete && athleteData" class="athlete-registry-dashboard mt-4">
                  <div class="card border-0 shadow-sm bg-gradient-primary-to-secondary text-white mb-4 overflow-hidden position-relative" style="border-radius: 12px;">
                    <div class="card-body p-4 position-relative z-index-1">
                      <div class="d-flex align-items-center gap-3">
                        <span class="athlete-badge"><i class="icofont-runner-alt-1"></i></span>
                        <div>
                          <span class="text-uppercase tracking-wider small opacity-75 d-block" style="font-size: 10px; font-weight: 700; letter-spacing: 1px;">NCS Verified Athlete Registry</span>
                          <h3 class="mb-0 fw-bold text-white fs-4">{{ athleteData.full_name }}</h3>
                          <span class="small fw-medium opacity-85">ID: {{ athleteData.athlete_number }} &middot; Discipline: {{ athleteData.discipline }}</span>
                        </div>
                      </div>
                    </div>
                    <div class="watermark-icon"><i class="icofont-trophy"></i></div>
                  </div>

                  <div class="athlete-stats-grid">
                    <!-- General Details -->
                    <article class="athlete-detail-card">
                      <header><i class="icofont-id-card text-primary"></i> <h3>Classification</h3></header>
                      <ul>
                        <li><span>Age Category</span> <strong>{{ athleteData.age_category || 'Senior' }}</strong></li>
                        <li><span>District / Region</span> <strong>{{ athleteData.district }} ({{ athleteData.region }})</strong></li>
                        <li><span>Affiliated Club</span> <strong>{{ athleteData.club || '-' }}</strong></li>
                        <li><span>License Status</span> <span class="badge bg-success-light text-success text-uppercase" style="font-size: 10px; font-weight: 700;">{{ athleteData.status }}</span></li>
                      </ul>
                    </article>

                    <!-- National Team Caps -->
                    <article class="athlete-detail-card">
                      <header><i class="icofont-flag text-danger"></i> <h3>National Duty</h3></header>
                      <ul v-if="athleteNationalTeam">
                        <li><span>Squad Tier</span> <strong>{{ athleteNationalTeam.category }}</strong></li>
                        <li><span>Team Name</span> <strong>{{ athleteNationalTeam.team_name }}</strong></li>
                        <li><span>Total Caps</span> <strong class="badge bg-danger text-white fs-6 py-1 px-2 rounded-circle" style="min-width: 24px;">{{ athleteNationalTeam.appearances_count }}</strong></li>
                        <li><span>First Call Up</span> <strong>{{ formatDate(athleteNationalTeam.first_call_up_on) }}</strong></li>
                      </ul>
                      <div v-else class="text-center py-4 text-muted small">
                        No official national team caps recorded yet.
                      </div>
                    </article>

                    <!-- Medical & Safeguarding -->
                    <article class="athlete-detail-card">
                      <header><i class="icofont-shield-alt text-warning"></i> <h3>Safeguarding & Medical</h3></header>
                      <ul>
                        <li><span>Blood Group</span> <strong>{{ athleteMedical?.blood_group || '-' }}</strong></li>
                        <li><span>Injury Clearance</span> <strong>{{ athleteMedical?.current_injury_status || 'Fit' }}</strong></li>
                        <li><span>Anti-Doping Ed.</span> <strong>{{ athleteSafeguarding?.anti_doping_education_completed ? 'Completed' : 'Pending' }}</strong></li>
                        <li><span>Safeguarding Status</span> <span class="badge" :class="athleteSafeguarding?.consent_forms_url ? 'bg-success-light text-success' : 'bg-warning-light text-warning'" style="font-size: 10px; font-weight: 700;">{{ athleteSafeguarding?.consent_forms_url ? 'Cleared' : 'Pending Consent' }}</span></li>
                      </ul>
                    </article>

                    <!-- Anti-Doping Compliance -->
                    <article class="athlete-detail-card">
                      <header><i class="icofont-test-bulb text-success"></i> <h3>WADA Compliance</h3></header>
                      <ul v-if="athleteAntiDoping">
                        <li><span>Testing Pool</span> <strong>{{ athleteAntiDoping.testing_status }}</strong></li>
                        <li><span>Last Tested On</span> <strong>{{ formatDate(athleteAntiDoping.last_tested_on) }}</strong></li>
                        <li><span>Test Result</span> <span class="badge" :class="athleteAntiDoping.last_test_result === 'NEGATIVE' ? 'bg-success-light text-success' : 'bg-danger-light text-danger'" style="font-size: 10px; font-weight: 700;">{{ athleteAntiDoping.last_test_result }}</span></li>
                        <li><span>WADA Education</span> <strong>{{ athleteAntiDoping.wada_education_completed ? 'Completed' : 'Pending' }}</strong></li>
                      </ul>
                      <div v-else class="text-center py-4 text-muted small">
                        No anti-doping tests or logs found.
                      </div>
                    </article>

                    <!-- Medal Standings -->
                    <article class="athlete-detail-card">
                      <header><i class="icofont-badge text-warning"></i> <h3>Medal Standings</h3></header>
                      <ul v-if="athleteMedals.length">
                        <li v-for="medal in athleteMedals" :key="medal.id">
                          <span>{{ medal.event }}</span>
                          <strong :class="getMedalBadgeClass(medal.medal_type)">{{ medal.medal_type }}</strong>
                        </li>
                      </ul>
                      <div v-else class="text-center py-4 text-muted small">
                        No medal records linked to this profile.
                      </div>
                    </article>

                    <!-- Competition Results -->
                    <article class="athlete-detail-card">
                      <header><i class="icofont-spreadsheet text-info"></i> <h3>Competition Results</h3></header>
                      <ul v-if="athleteResults.length" style="max-height: 180px; overflow-y: auto; padding-right: 4px;">
                        <li v-for="res in athleteResults" :key="res.id" style="border-bottom: 1px dashed rgba(0,0,0,0.06); padding-bottom: 6px; margin-bottom: 6px;">
                          <span>{{ res.event }}</span>
                          <strong>Rank: {{ res.position || 'N/A' }} ({{ res.time_result || res.score_result || res.distance_result || 'Result' }})</strong>
                        </li>
                      </ul>
                      <div v-else class="text-center py-4 text-muted small">
                        No competition results recorded yet.
                      </div>
                    </article>
                  </div>
                </div>
              </template>

              <template v-else-if="section === 'apply'">
                <header class="page-heading"><div><p>Apply Now</p><h1>Open application forms</h1><span>Choose a service to begin or continue your application.</span></div></header>
                <OpenFormsPanel :forms="openForms" :submissions="dynamicSubmissions" :loading="loading" expanded @start="startApplication" />
              </template>

              <template v-else-if="section === 'applications'">
                <header class="page-heading"><div><p>My Applications</p><h1>Application history</h1><span>All standard and custom-form applications linked to your account.</span></div></header>
                <div class="list-toolbar">
                  <label><i class="icofont-search-1"></i><input v-model="applicationSearch" type="search" placeholder="Search your applications" /></label>
                  <select v-model="applicationStatus"><option value="">All statuses</option><option v-for="item in applicationStatuses" :key="item" :value="item">{{ titleize(item) }}</option></select>
                </div>
                <div class="data-table">
                  <table><thead><tr><th>Application</th><th>Reference</th><th>Status</th><th>Payment</th><th>Updated</th><th></th></tr></thead>
                    <tbody>
                      <tr v-if="!filteredApplications.length"><td colspan="6" class="empty-cell">You have not started an application matching these filters.</td></tr>
                      <tr v-for="item in filteredApplications" :key="`${item.source}-${item.id}`">
                        <td data-label="Application"><strong>{{ item.title }}</strong><small>{{ item.source === 'custom' ? 'Custom form' : 'Standard form' }}</small></td>
                        <td data-label="Reference">{{ item.reference || 'Pending' }}</td>
                        <td data-label="Status"><span class="status" :class="statusClass(item.status)">{{ titleize(item.status) }}</span></td>
                        <td data-label="Payment">{{ titleize(item.payment_status || 'not required') }}</td>
                        <td data-label="Updated">{{ formatDate(item.updated_at) }}</td>
                        <td data-label="Action">
                          <span class="table-actions">
                            <button type="button" title="View application" @click="viewApplication(item)"><i class="icofont-eye-alt"></i></button>
                            <button type="button" title="Download application form" @click="downloadApplication(item)"><i class="icofont-download"></i></button>
                            <button v-if="item.source === 'custom' && isEditableSubmission(item)" type="button" title="Edit application" @click="continueApplication(item)"><i class="icofont-rounded-right"></i></button>
                          </span>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </template>

              <template v-else-if="section === 'activities'">
                <header class="page-heading">
                  <div>
                    <p>Security & Audit Trail</p>
                    <h1>My Activities & Audit Logs</h1>
                    <span>Comprehensive record of your authentication events, application updates, and security logs.</span>
                  </div>
                  <button type="button" class="secondary-command" :disabled="loadingActivities" @click="loadActivities(activityPage)">
                    <i class="icofont-refresh" :class="{ 'animate-spin': loadingActivities }"></i> Refresh logs
                  </button>
                </header>

                <!-- Activity Summary KPIs -->
                <div class="user-kpis activity-kpis">
                  <article>
                    <span class="blue"><i class="icofont-history"></i></span>
                    <div>
                      <small>Total Events</small>
                      <strong>{{ activityTotal }}</strong>
                      <p>All recorded actions</p>
                    </div>
                  </article>
                  <article>
                    <span class="green"><i class="icofont-shield-check"></i></span>
                    <div>
                      <small>Auth & Security</small>
                      <strong>{{ activityStats.authCount }}</strong>
                      <p>Logins & credentials</p>
                    </div>
                  </article>
                  <article>
                    <span class="amber"><i class="icofont-file-document"></i></span>
                    <div>
                      <small>Applications</small>
                      <strong>{{ activityStats.appCount }}</strong>
                      <p>Drafts & submissions</p>
                    </div>
                  </article>
                  <article>
                    <span class="cyan"><i class="icofont-clock-time"></i></span>
                    <div>
                      <small>Last Active</small>
                      <strong style="font-size: 15px; font-weight: 700; margin-top: 4px;">{{ activityStats.lastActive }}</strong>
                      <p>Latest event timestamp</p>
                    </div>
                  </article>
                </div>

                <!-- Activities Filter & Search Toolbar -->
                <div class="list-toolbar activity-toolbar">
                  <label class="search-label">
                    <i class="icofont-search-1"></i>
                    <input v-model="activitySearchQuery" type="search" placeholder="Search activities by action, IP, or location..." />
                  </label>
                  <div class="activity-category-pills">
                    <button
                      type="button"
                      class="category-pill"
                      :class="{ active: activityCategoryFilter === 'ALL' }"
                      @click="activityCategoryFilter = 'ALL'"
                    >
                      All ({{ activities.length }})
                    </button>
                    <button
                      type="button"
                      class="category-pill"
                      :class="{ active: activityCategoryFilter === 'AUTH' }"
                      @click="activityCategoryFilter = 'AUTH'"
                    >
                      <i class="icofont-key"></i> Auth & Security
                    </button>
                    <button
                      type="button"
                      class="category-pill"
                      :class="{ active: activityCategoryFilter === 'APPLICATIONS' }"
                      @click="activityCategoryFilter = 'APPLICATIONS'"
                    >
                      <i class="icofont-file-alt"></i> Applications
                    </button>
                    <button
                      type="button"
                      class="category-pill"
                      :class="{ active: activityCategoryFilter === 'PROFILE' }"
                      @click="activityCategoryFilter = 'PROFILE'"
                    >
                      <i class="icofont-user-alt-7"></i> Profile & Account
                    </button>
                  </div>
                </div>

                <!-- Activities Timeline List -->
                <div v-if="loadingActivities" class="text-center py-5">
                  <div class="spinner-border text-primary" role="status" style="width: 3rem; height: 3rem;"></div>
                  <p class="text-muted mt-3 font-weight-bold">Loading activity history...</p>
                </div>
                <div v-else-if="!filteredActivities.length" class="empty-state card text-center p-5 border-0 shadow-sm" style="border-radius: 10px; background: #fff;">
                  <i class="icofont-history" style="font-size: 48px; color: #cbd5e1; margin-bottom: 12px;"></i>
                  <h3 class="font-weight-bold" style="font-size: 18px; color: #1e293b;">No activities match your filters</h3>
                  <p class="text-muted small mb-0">Try clearing your search query or selecting a different category filter.</p>
                </div>
                <div v-else class="activity-timeline-feed">
                  <div v-for="item in filteredActivities" :key="item.id || item.created_at" class="activity-feed-card">
                    <div class="activity-icon-wrapper" :class="getActivityTone(item)">
                      <i :class="getActivityIcon(item)"></i>
                    </div>
                    <div class="activity-content">
                      <div class="activity-header-row">
                        <div class="activity-title-group">
                          <h4 class="activity-title">{{ activityTitle(item) }}</h4>
                          <span class="activity-category-badge" :class="getActivityTone(item)">
                            {{ getActivityCategory(item) }}
                          </span>
                        </div>
                        <time class="activity-time" :title="formatDate(item.created_at)">
                          <i class="icofont-clock-time"></i> {{ formatRelativeTime(item.created_at) }}
                        </time>
                      </div>
                      <p class="activity-description">{{ activityDescription(item) }}</p>
                      <div class="activity-meta-row">
                        <span v-if="item.ip_address" class="meta-pill">
                          <i class="icofont-globe"></i> IP: <strong>{{ item.ip_address }}</strong>
                        </span>
                        <span v-if="item.geo_city || item.geo_country" class="meta-pill">
                          <i class="icofont-location-pin"></i> {{ [item.geo_city, item.geo_country].filter(Boolean).join(', ') }}
                        </span>
                        <span v-if="item.browser || item.os_name" class="meta-pill">
                          <i class="icofont-laptop"></i> {{ [item.os_name, item.browser].filter(Boolean).join(' ') }}
                        </span>
                        <span class="meta-pill date-pill">
                          <i class="icofont-calendar"></i> {{ formatDate(item.created_at) }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Pagination -->
                <div v-if="activityTotalPages > 1" class="activity-pagination-bar mt-4">
                  <div class="pagination-info">
                    Showing Page <strong>{{ activityPage }}</strong> of <strong>{{ activityTotalPages }}</strong> ({{ activityTotal }} total events recorded)
                  </div>
                  <div class="pagination-buttons">
                    <button type="button" class="btn btn-sm btn-secondary" :disabled="activityPage <= 1 || loadingActivities" @click="loadActivities(activityPage - 1)">
                      <i class="icofont-arrow-left"></i> Previous
                    </button>
                    <button type="button" class="btn btn-sm btn-primary" :disabled="activityPage >= activityTotalPages || loadingActivities" @click="loadActivities(activityPage + 1)">
                      Next <i class="icofont-arrow-right"></i>
                    </button>
                  </div>
                </div>
              </template>

              <template v-else-if="section === 'notifications'">
                <header class="page-heading"><div><p>Notifications</p><h1>Your notifications</h1><span>Application and account updates from NCS.</span></div><button v-if="unreadNotifications" type="button" class="secondary-command" @click="markAllNotificationsRead"><i class="icofont-check"></i> Mark all read</button></header>
                <div class="feed-list">
                  <article v-if="!notifications.length" class="empty-state"><i class="icofont-notification"></i><h3>No notifications</h3><p>New updates will appear here.</p></article>
                  <article v-for="item in notifications" :key="item.id" :class="{ unread: item.status === 'unread' }">
                    <span><i :class="notificationIcon(item)"></i></span><div><strong>{{ item.title || 'Portal notification' }}</strong><p>{{ item.message }}</p><small>{{ formatDate(item.created_at) }}</small></div>
                    <button v-if="item.status === 'unread'" type="button" title="Mark as read" @click="markNotificationRead(item)"><i class="icofont-check"></i></button>
                  </article>
                </div>
              </template>

              <template v-else-if="section === 'messages'">
                <header class="page-heading"><div><p>Messages</p><h1>Application messages</h1><span>Review notes and decisions sent through your applications.</span></div></header>
                <div class="feed-list">
                  <article v-if="!userMessages.length" class="empty-state"><i class="icofont-envelope-open"></i><h3>No messages</h3><p>Messages from application reviewers will appear here.</p></article>
                  <article v-for="item in userMessages" :key="item.id"><span><i class="icofont-envelope"></i></span><div><strong>{{ item.title }}</strong><p>{{ item.message }}</p><small>{{ formatDate(item.created_at) }}</small></div></article>
                </div>
              </template>

              <template v-else-if="section === 'transactions'">
                <header class="page-heading"><div><p>My Transactions</p><h1>Payment records</h1><span>Payment references and verification status for your applications.</span></div></header>
                <div class="transaction-summary"><span>Total recorded</span><strong>UGX {{ formatMoney(transactionTotal) }}</strong></div>
                <div class="data-table">
                  <table><thead><tr><th>Application</th><th>Reference</th><th>Amount</th><th>Payment status</th><th>Method</th></tr></thead>
                    <tbody>
                      <tr v-if="!transactions.length"><td colspan="5" class="empty-cell">No transactions have been recorded.</td></tr>
                      <tr v-for="item in transactions" :key="`${item.source}-${item.id}`"><td data-label="Application"><strong>{{ item.title }}</strong></td><td data-label="Reference">{{ item.payment_reference || '-' }}</td><td data-label="Amount">UGX {{ formatMoney(item.amount) }}</td><td data-label="Payment status"><span class="status" :class="paymentClass(item.payment_status)">{{ titleize(item.payment_status || 'unpaid') }}</span></td><td data-label="Method">{{ titleize(item.payment_method || 'proof upload') }}</td></tr>
                    </tbody>
                  </table>
                </div>
              </template>

              <template v-else-if="section === 'my-files'">
                <header class="page-heading">
                  <div>
                    <p>Credential Wallet</p>
                    <h1>My Files & Certificates</h1>
                    <span>Access and download all official licenses, registrations, and certificates offered to you by NCS.</span>
                  </div>
                </header>

                <div class="user-kpis">
                  <article>
                    <span class="bg-primary text-white" style="display: flex; align-items: center; justify-content: center; width: 48px; height: 48px; border-radius: 50%;"><i class="icofont-folder-open" style="font-size: 20px;"></i></span>
                    <div>
                      <small>Total Credentials</small>
                      <strong>{{ userFiles.length }}</strong>
                      <p>Active NCS certificates</p>
                    </div>
                  </article>
                  <article>
                    <span class="bg-success text-white" style="display: flex; align-items: center; justify-content: center; width: 48px; height: 48px; border-radius: 50%;"><i class="icofont-check-circled" style="font-size: 20px;"></i></span>
                    <div>
                      <small>Verification Status</small>
                      <strong>100%</strong>
                      <p>All records verified</p>
                    </div>
                  </article>
                </div>

                <div v-if="!userFiles.length" class="empty-state py-5 card shadow-sm text-center border-0" style="border-radius: 12px; background: #fff;">
                  <div class="card-body p-5">
                    <i class="icofont-folder-open text-muted" style="font-size: 64px;"></i>
                    <h3 class="mt-3 fw-bold text-dark">No Credentials Found</h3>
                    <p class="text-muted max-w-md mx-auto">
                      Currently, there are no active athlete licenses, coach credentials, or technical official clearances linked to your email in the system.
                    </p>
                  </div>
                </div>

                <div v-else class="row">
                  <div v-for="file in userFiles" :key="file.id" class="col-md-6 mb-4">
                    <div class="card shadow-sm border-0 h-100 credential-card" style="border-radius: 12px; background: #fff; transition: transform 0.2s; box-shadow: 0 4px 20px rgba(0,0,0,0.05) !important;">
                      <div class="card-body p-4 d-flex flex-column h-100">
                        <header class="d-flex align-items-start justify-content-between mb-3">
                          <div class="d-flex align-items-center gap-3">
                            <span class="credential-icon p-2 rounded-circle" :class="getFileIconClass(file.category)" style="display: inline-flex; align-items: center; justify-content: center; width: 40px; height: 40px;">
                              <i :class="getFileIcon(file.category)" style="font-size: 20px;"></i>
                            </span>
                            <div>
                              <span class="badge bg-light text-muted text-uppercase mb-1" style="font-size: 9px; font-weight: 700; border: 1px solid rgba(0,0,0,0.06);">{{ file.type }}</span>
                              <h3 class="h5 mb-0 fw-bold text-dark" style="font-size: 15px; font-weight: 700;">{{ file.title }}</h3>
                            </div>
                          </div>
                          <span class="badge" :class="getStatusBadgeClass(file.status)" style="font-size: 10px; font-weight: 700;">{{ file.status }}</span>
                        </header>
                        
                        <p class="text-muted small mb-4 flex-grow-1" style="font-size: 12px; line-height: 1.5;">{{ file.description }}</p>
                        
                        <div class="mt-auto pt-3 border-top d-flex flex-column gap-2" style="border-color: rgba(0,0,0,0.06) !important;">
                          <div class="d-flex justify-content-between text-muted small" style="font-size: 12px;">
                            <span>License No:</span>
                            <strong class="text-dark">{{ file.number }}</strong>
                          </div>
                          <div class="d-flex justify-content-between text-muted small" style="font-size: 12px;">
                            <span>Issued On:</span>
                            <strong>{{ formatDate(file.issueDate) }}</strong>
                          </div>
                          <button type="button" class="btn btn-primary btn-block mt-3 d-flex align-items-center justify-content-center gap-2" style="background-color: #6777ef; border-color: #6777ef; border-radius: 30px; font-weight: 700; font-size: 12px; padding: 10px 18px;" @click="downloadFile(file)">
                            <i class="icofont-download"></i> Download Document ({{ file.fileType }})
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </template>

              <!-- ── Settings Section (Password Change, Sessions, Preferences & 2FA) ── -->
              <template v-else-if="section === 'settings'">
                <header class="page-heading">
                  <div>
                    <p>Account Settings</p>
                    <h1>Security, Active Sessions & Preferences</h1>
                    <span>Manage your credentials, review signed-in devices, and configure two-factor authentication.</span>
                  </div>
                </header>

                <!-- Settings Sub-Navigation Tabs -->
                <div class="settings-nav-tabs">
                  <button
                    type="button"
                    class="settings-tab-btn"
                    :class="{ active: settingsTab === 'password' }"
                    @click="settingsTab = 'password'"
                  >
                    <i class="icofont-key"></i>
                    <span>Password & Security</span>
                  </button>
                  <button
                    type="button"
                    class="settings-tab-btn"
                    :class="{ active: settingsTab === 'sessions' }"
                    @click="onSelectSessionsTab"
                  >
                    <i class="icofont-laptop"></i>
                    <span>Active Sessions</span>
                    <span v-if="sessions.length" class="badge bg-primary text-white ms-2" style="font-size: 11px;">{{ sessions.length }}</span>
                  </button>
                  <button
                    type="button"
                    class="settings-tab-btn"
                    :class="{ active: settingsTab === 'preferences' }"
                    @click="settingsTab = 'preferences'"
                  >
                    <i class="icofont-shield-check"></i>
                    <span>Preferences & 2FA</span>
                  </button>
                </div>

                <!-- Tab 1: Password & Security -->
                <div v-if="settingsTab === 'password'" class="settings-tab-content">
                  <div class="row g-4">
                    <div class="col-lg-7">
                      <div class="settings-card">
                        <div class="settings-card-header">
                          <i class="icofont-key text-primary fs-4"></i>
                          <div>
                            <h3>Change Password</h3>
                            <p>Choose a strong, unique password to secure your NCS Portal account.</p>
                          </div>
                        </div>

                        <form @submit.prevent="updatePassword" class="settings-form">
                          <div class="form-group mb-3">
                            <label class="form-label font-weight-bold">Current Password <span class="text-danger">*</span></label>
                            <div class="password-input-wrap">
                              <input
                                v-model="passwordCurrent"
                                :type="showPasswordCurrent ? 'text' : 'password'"
                                class="form-control"
                                placeholder="Enter your current password"
                                required
                              />
                              <button type="button" class="btn-toggle-eye" @click="showPasswordCurrent = !showPasswordCurrent">
                                <i :class="showPasswordCurrent ? 'icofont-eye' : 'icofont-eye-blocked'"></i>
                              </button>
                            </div>
                          </div>

                          <div class="form-group mb-3">
                            <label class="form-label font-weight-bold">New Password <span class="text-danger">*</span></label>
                            <div class="password-input-wrap">
                              <input
                                v-model="passwordNew"
                                :type="showPasswordNew ? 'text' : 'password'"
                                class="form-control"
                                placeholder="Enter new password (min 8 characters)"
                                required
                              />
                              <button type="button" class="btn-toggle-eye" @click="showPasswordNew = !showPasswordNew">
                                <i :class="showPasswordNew ? 'icofont-eye' : 'icofont-eye-blocked'"></i>
                              </button>
                            </div>

                            <!-- Live Password Strength Meter -->
                            <div v-if="passwordNew" class="password-strength-box mt-2">
                              <div class="strength-bar-track">
                                <div class="strength-bar-fill" :class="passwordStrengthClass" :style="{ width: passwordStrengthPercent + '%' }"></div>
                              </div>
                              <div class="d-flex justify-content-between align-items-center mt-1">
                                <small class="strength-label">Strength: <strong :class="passwordStrengthTextClass">{{ passwordStrengthLabel }}</strong></small>
                              </div>
                              <ul class="password-rules-list mt-2">
                                <li :class="{ met: passwordNew.length >= 8 }">
                                  <i :class="passwordNew.length >= 8 ? 'icofont-check-circled text-success' : 'icofont-close-circled text-muted'"></i>
                                  At least 8 characters
                                </li>
                                <li :class="{ met: /[A-Z]/.test(passwordNew) && /[a-z]/.test(passwordNew) }">
                                  <i :class="(/[A-Z]/.test(passwordNew) && /[a-z]/.test(passwordNew)) ? 'icofont-check-circled text-success' : 'icofont-close-circled text-muted'"></i>
                                  Uppercase & lowercase letters
                                </li>
                                <li :class="{ met: /[0-9]/.test(passwordNew) }">
                                  <i :class="/[0-9]/.test(passwordNew) ? 'icofont-check-circled text-success' : 'icofont-close-circled text-muted'"></i>
                                  At least one number
                                </li>
                                <li :class="{ met: /[^A-Za-z0-9]/.test(passwordNew) }">
                                  <i :class="/[^A-Za-z0-9]/.test(passwordNew) ? 'icofont-check-circled text-success' : 'icofont-close-circled text-muted'"></i>
                                  Special symbol (!@#$%^&*)
                                </li>
                              </ul>
                            </div>
                          </div>

                          <div class="form-group mb-4">
                            <label class="form-label font-weight-bold">Confirm New Password <span class="text-danger">*</span></label>
                            <div class="password-input-wrap">
                              <input
                                v-model="passwordConfirm"
                                :type="showPasswordConfirm ? 'text' : 'password'"
                                class="form-control"
                                placeholder="Re-enter your new password"
                                required
                              />
                              <button type="button" class="btn-toggle-eye" @click="showPasswordConfirm = !showPasswordConfirm">
                                <i :class="showPasswordConfirm ? 'icofont-eye' : 'icofont-eye-blocked'"></i>
                              </button>
                            </div>
                            <small v-if="passwordConfirm && passwordNew !== passwordConfirm" class="text-danger mt-1 d-block font-weight-bold">
                              <i class="icofont-warning"></i> Passwords do not match.
                            </small>
                            <small v-else-if="passwordConfirm && passwordNew === passwordConfirm" class="text-success mt-1 d-block font-weight-bold">
                              <i class="icofont-check-circled"></i> Passwords match!
                            </small>
                          </div>

                          <button
                            type="submit"
                            class="primary-command"
                            :disabled="changingPassword || (passwordNew && passwordNew !== passwordConfirm)"
                          >
                            <i v-if="changingPassword" class="icofont-spinner animate-spin"></i>
                            <i v-else class="icofont-key"></i>
                            {{ changingPassword ? 'Updating Password...' : 'Update Password' }}
                          </button>
                        </form>
                      </div>
                    </div>

                    <div class="col-lg-5">
                      <div class="settings-card security-tips-card">
                        <div class="settings-card-header">
                          <i class="icofont-shield-alt text-info fs-4"></i>
                          <div>
                            <h3>Password Security Tips</h3>
                            <p>Best practices for safeguarding your NCS account.</p>
                          </div>
                        </div>
                        <ul class="security-tips-list">
                          <li>
                            <i class="icofont-check text-success"></i>
                            <div>
                              <strong>Use Unique Passwords</strong>
                              <p>Avoid reusing passwords from other personal or work websites.</p>
                            </div>
                          </li>
                          <li>
                            <i class="icofont-check text-success"></i>
                            <div>
                              <strong>Combine Words & Symbols</strong>
                              <p>Passphrases with numbers and punctuation are significantly more resilient to automated attacks.</p>
                            </div>
                          </li>
                          <li>
                            <i class="icofont-check text-success"></i>
                            <div>
                              <strong>Enable 2FA Verification</strong>
                              <p>Activate email-based two-factor authentication in the Preferences tab for extra defense.</p>
                            </div>
                          </li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Tab 2: Sessions Management -->
                <div v-else-if="settingsTab === 'sessions'" class="settings-tab-content">
                  <div class="settings-card mb-4">
                    <div class="d-flex align-items-center justify-content-between flex-wrap gap-3 pb-3 border-bottom mb-4">
                      <div class="d-flex align-items-center gap-3">
                        <span class="p-3 bg-light-primary text-primary rounded-circle" style="display: flex; align-items: center; justify-content: center; width: 50px; height: 50px; background: #e0e7ff;">
                          <i class="icofont-laptop" style="font-size: 24px;"></i>
                        </span>
                        <div>
                          <h3 class="mb-1 font-weight-bold" style="font-size: 18px; color: #1e293b;">Active Login Sessions</h3>
                          <p class="text-muted small mb-0">These are web browsers and devices that have recently authenticated to your NCS portal.</p>
                        </div>
                      </div>
                      <div class="d-flex gap-2">
                        <button
                          type="button"
                          class="secondary-command"
                          :disabled="loadingSessions"
                          @click="loadSessions"
                        >
                          <i class="icofont-refresh" :class="{ 'animate-spin': loadingSessions }"></i> Refresh
                        </button>
                        <button
                          v-if="sessions.length > 1"
                          type="button"
                          class="btn btn-danger btn-sm"
                          style="border-radius: 30px; font-weight: 600;"
                          :disabled="revokingAllOtherSessions"
                          @click="revokeAllOtherSessions"
                        >
                          <i v-if="revokingAllOtherSessions" class="icofont-spinner animate-spin"></i>
                          <i v-else class="icofont-logout"></i>
                          Sign Out All Other Devices
                        </button>
                      </div>
                    </div>

                    <!-- Sessions Feed -->
                    <div v-if="loadingSessions" class="text-center py-5">
                      <div class="spinner-border text-primary" role="status"></div>
                      <p class="text-muted mt-3 font-weight-bold">Loading active sessions...</p>
                    </div>

                    <div v-else-if="!sessions.length" class="empty-state text-center py-5">
                      <i class="icofont-laptop text-muted" style="font-size: 48px;"></i>
                      <h4 class="mt-3 font-weight-bold">No active sessions found</h4>
                      <p class="text-muted small">Your current session details will appear on your next refresh.</p>
                    </div>

                    <div v-else class="sessions-list-grid">
                      <div
                        v-for="(sess, index) in sessions"
                        :key="sess.id || index"
                        class="session-item-card"
                        :class="{ 'current-session-card': index === 0 }"
                      >
                        <div class="session-device-icon" :class="{ 'current': index === 0 }">
                          <i :class="getSessionDeviceIcon(sess.user_agent)"></i>
                        </div>
                        <div class="session-details">
                          <div class="d-flex align-items-center gap-2 flex-wrap mb-1">
                            <h4 class="session-device-name">{{ parseSessionBrowser(sess.user_agent) }}</h4>
                            <span v-if="index === 0" class="badge bg-success text-white current-badge">
                              <i class="icofont-check-circled"></i> Current Session (This Device)
                            </span>
                          </div>
                          <div class="session-meta-items">
                            <span class="meta-item">
                              <i class="icofont-globe"></i> IP Address: <strong>{{ sess.ip_address || 'Current Network' }}</strong>
                            </span>
                            <span class="meta-item">
                              <i class="icofont-calendar"></i> Signed in: {{ formatDate(sess.created_at) }}
                            </span>
                            <span v-if="sess.expires_at" class="meta-item">
                              <i class="icofont-clock-time"></i> Expires: {{ formatDate(sess.expires_at) }}
                            </span>
                          </div>
                        </div>
                        <div v-if="index !== 0" class="session-actions">
                          <button
                            type="button"
                            class="btn btn-outline-danger btn-sm"
                            style="border-radius: 20px; font-weight: 600;"
                            :disabled="revokingSessionId === sess.id"
                            @click="revokeSession(sess.id)"
                          >
                            <i v-if="revokingSessionId === sess.id" class="icofont-spinner animate-spin"></i>
                            <i v-else class="icofont-ui-delete"></i> Revoke
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Tab 3: Preferences & 2FA Setup -->
                <div v-if="settingsTab === 'preferences'" class="settings-tab-content">
                  <div class="row g-4">
                    <!-- Two-Factor Authentication Box -->
                    <div class="col-lg-6">
                      <div class="settings-card h-100">
                        <div class="settings-card-header">
                          <i class="icofont-shield-check text-primary fs-4"></i>
                          <div>
                            <h3>Two-Factor Authentication (2FA)</h3>
                            <p>Enhance account protection with email one-time passcodes.</p>
                          </div>
                        </div>

                        <div class="twofa-status-banner p-3 rounded-3 mb-4" :class="twofaEnabled ? 'twofa-active' : 'twofa-inactive'">
                          <div class="d-flex align-items-center justify-content-between">
                            <div class="d-flex align-items-center gap-2">
                              <i :class="twofaEnabled ? 'icofont-check-circled text-success fs-4' : 'icofont-exclamation-circle text-warning fs-4'"></i>
                              <div>
                                <strong class="d-block font-weight-bold" :class="twofaEnabled ? 'text-success' : 'text-warning'">
                                  {{ twofaEnabled ? '2FA is currently ENABLED' : '2FA is currently DISABLED' }}
                                </strong>
                                <small class="text-muted">
                                  {{ twofaEnabled ? 'Your account requires email verification code upon login.' : 'Activate 2FA for mandatory login challenge verification.' }}
                                </small>
                              </div>
                            </div>
                            <span class="badge" :class="twofaEnabled ? 'bg-success text-white' : 'bg-warning text-dark'" style="font-weight: 700; font-size: 11px;">
                              {{ twofaEnabled ? 'PROTECTED' : 'ACTION RECOMMENDED' }}
                            </span>
                          </div>
                        </div>

                        <p class="text-muted small mb-4" style="line-height: 1.6;">
                          When 2FA is active, every time you sign in with your email and password, a secure 6-digit one-time code will be dispatched to <strong>{{ profile.email }}</strong>. You must enter this code to complete authentication.
                        </p>

                        <!-- If 2FA is Disabled: Show Enable Button -->
                        <div v-if="!twofaEnabled && !showTwoFAVerify">
                          <button
                            type="button"
                            class="primary-command"
                            :disabled="toggling2FA"
                            @click="toggle2FA"
                          >
                            <i v-if="toggling2FA" class="icofont-spinner animate-spin"></i>
                            <i v-else class="icofont-shield-check"></i>
                            {{ toggling2FA ? 'Sending verification code...' : 'Set Up Two-Factor Authentication' }}
                          </button>
                        </div>

                        <!-- 2FA Enrollment Verification Panel -->
                        <div v-if="showTwoFAVerify" class="p-3 border rounded-3 bg-light mt-3" style="border-color: #fef08a !important; background: #fefce8 !important;">
                          <h4 class="font-weight-bold text-dark mb-1" style="font-size: 14px;">Verify Activation Code</h4>
                          <p class="text-muted small mb-3">A 6-digit activation code was sent to <strong>{{ profile.email }}</strong>. Enter it below to activate 2FA:</p>
                          <div class="d-flex gap-2 align-items-center">
                            <input
                              v-model.trim="twofaVerifyCode"
                              type="text"
                              maxlength="6"
                              placeholder="123456"
                              class="form-control text-center font-weight-bold"
                              style="width: 140px; font-size: 18px; letter-spacing: 4px; border: 2px solid #6777ef; background: #fff;"
                            />
                            <button
                              type="button"
                              class="primary-command"
                              :disabled="verifying2FA || !twofaVerifyCode || twofaVerifyCode.length < 6"
                              @click="confirm2FA"
                            >
                              <i v-if="verifying2FA" class="icofont-spinner animate-spin"></i>
                              <i v-else class="icofont-check"></i>
                              Confirm & Activate
                            </button>
                            <button
                              type="button"
                              class="secondary-command"
                              @click="cancel2FAEnrollment"
                            >
                              Cancel
                            </button>
                          </div>
                        </div>

                        <!-- If 2FA is Enabled: Show Disable Option -->
                        <div v-if="twofaEnabled">
                          <button
                            type="button"
                            class="btn btn-outline-danger"
                            style="border-radius: 30px; font-weight: 600; font-size: 13px; padding: 8px 20px;"
                            :disabled="toggling2FA"
                            @click="toggle2FA"
                          >
                            <i v-if="toggling2FA" class="icofont-spinner animate-spin"></i>
                            <i v-else class="icofont-ui-close"></i>
                            {{ toggling2FA ? 'Disabling...' : 'Disable Two-Factor Authentication' }}
                          </button>
                        </div>
                      </div>
                    </div>

                    <!-- Notification & Communication Preferences -->
                    <div class="col-lg-6">
                      <div class="settings-card h-100">
                        <div class="settings-card-header">
                          <i class="icofont-notification text-primary fs-4"></i>
                          <div>
                            <h3>Notification Preferences</h3>
                            <p>Choose which alerts and email dispatches you wish to receive.</p>
                          </div>
                        </div>

                        <div class="preferences-list">
                          <label class="preference-item">
                            <input type="checkbox" v-model="prefEmailAppStatus" class="pref-checkbox" />
                            <div class="preference-text">
                              <strong>Application Status Updates</strong>
                              <p>Receive email alerts when reviewers approve, request information, or issue decisions on your applications.</p>
                            </div>
                          </label>

                          <label class="preference-item">
                            <input type="checkbox" v-model="prefEmailSecurity" class="pref-checkbox" />
                            <div class="preference-text">
                              <strong>Security & Login Alerts</strong>
                              <p>Receive notifications whenever a new browser or IP signs into your account.</p>
                            </div>
                          </label>

                          <label class="preference-item">
                            <input type="checkbox" v-model="prefEmailAnnouncements" class="pref-checkbox" />
                            <div class="preference-text">
                              <strong>NCS Announcements & Bulletins</strong>
                              <p>Receive official sports notices, circulars, and registry updates from NCS administration.</p>
                            </div>
                          </label>
                        </div>

                        <div class="mt-4 pt-3 border-top d-flex justify-content-between align-items-center">
                          <button
                            type="button"
                            class="primary-command"
                            :disabled="savingPreferences"
                            @click="savePreferences"
                          >
                            <i v-if="savingPreferences" class="icofont-spinner animate-spin"></i>
                            <i v-else class="icofont-save"></i>
                            Save Preferences
                          </button>
                          <span v-if="preferencesSaved" class="text-success font-weight-bold small">
                            <i class="icofont-check-circled"></i> Saved!
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </template>

              <!-- ── My Profile Section ── -->
              <template v-else-if="section === 'profile'">
                <header class="page-heading">
                  <div>
                    <p>My Profile</p>
                    <h1>Personal Information & Identity</h1>
                    <span>Keep your identity details current and manage account profile settings.</span>
                  </div>
                  <button type="button" class="secondary-command" @click="select('settings')">
                    <i class="icofont-gear"></i> Manage Account Settings
                  </button>
                </header>
                <div class="profile-layout">
                  <div class="profile-grid-3">
                    <!-- Profile Photo & Summary -->
                    <article class="profile-summary">
                      <div class="profile-avatar">
                        <img v-if="profile.avatar_url" :src="mediaUrl(profile.avatar_url)" alt="" />
                        <span v-else>{{ initials }}</span>
                      </div>
                      <h2>{{ fullName }}</h2>
                      <p>{{ profile.email }}</p>
                      <div class="d-flex justify-content-center gap-1 mt-2">
                        <span class="badge bg-primary text-white" style="font-size: 11px; text-transform: uppercase;">{{ userProfileLabel }}</span>
                      </div>
                      
                      <!-- Avatar upload -->
                      <div style="margin-top: 16px; width: 100%;">
                        <input type="file" ref="avatarInput" accept="image/*" style="display:none;" @change="onAvatarFileSelected" />
                        <button type="button" class="secondary-command" style="width:100%; justify-content:center;" :disabled="uploadingAvatar" @click="$refs.avatarInput.click()">
                          <i v-if="uploadingAvatar" class="icofont-spinner animate-spin"></i>
                          <i v-else class="icofont-camera"></i>
                          {{ uploadingAvatar ? 'Uploading...' : 'Upload Photo' }}
                        </button>
                      </div>
                    </article>

                    <!-- Profile Details Form -->
                    <form class="profile-form" @submit.prevent="saveProfile">
                      <h2>Profile details</h2>
                      <div>
                        <label>First name<input v-model.trim="profile.first_name" required /></label>
                        <label>Last name<input v-model.trim="profile.last_name" required /></label>
                      </div>
                      <label>Email address<input :value="profile.email" type="email" disabled /></label>
                      <div class="d-flex justify-content-between align-items-center mt-2 flex-wrap gap-2">
                        <button type="submit" class="primary-command" :disabled="savingProfile"><i class="icofont-save"></i> {{ savingProfile ? 'Saving...' : 'Save profile' }}</button>
                        <button type="button" class="secondary-command" @click="select('settings')"><i class="icofont-key"></i> Change Password / Security</button>
                      </div>
                    </form>
                  </div>
                </div>
              </template>
            </div>
          </section>
        </div>
        <footer class="main-footer cms-main-footer">
          <div class="footer-left">
            Design By: Ateni Media Technologies LLC
          </div>
          <div class="footer-right">
            National Council of Sports Portal
          </div>
        </footer>
      </div>
    </div>

  </main>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getCurrentUser, updateMyProfile, uploadProfileAvatar } from '@/api/auth.js'
import { listMyAuditLogs } from '@/api/account.js'
import { listMyLegacyApplications, listMyTransactions } from '@/api/applications.js'
import { getMySecurity, enroll2FA, verify2FA, disable2FA, changePassword, listMySessions, revokeMySession, revokeOtherSessions } from '@/api/security.js'
import * as cms from '@/api/cms.js'
import { portalListOpenForms, portalListSubmissions } from '@/api/forms.js'
import { listNsmisDomain } from '@/api/nsmis.js'
import { mediaUrl } from '@/api/client.js'
import OpenFormsPanel from '@/components/portal/OpenFormsPanel.vue'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'
import { downloadApplicationForm } from '@/utils/applicationDownload.js'
import { ensureOtikaStyles } from '@/utils/otikaAssets.js'
import { recordMenuNavigation } from '@/services/activityAudit.js'

const router = useRouter()
const portalSectionIds = ['dashboard', 'apply', 'applications', 'my-files', 'activities', 'notifications', 'messages', 'transactions', 'settings', 'profile']
const section = ref('dashboard')
const sidebarCollapsed = ref(localStorage.getItem('ncsms_sidebar_collapsed') === 'true')
const mobileSidebarOpen = ref(false)
const profileOpen = ref(false)
const loading = ref(true)
const error = ref('')
const success = ref('')
const openForms = ref([])
const dynamicSubmissions = ref([])
const legacyApplications = ref([])
const legacyTransactions = ref([])
const activities = ref([])
const activityPage = ref(1)
const activityPerPage = ref(10)
const activityTotal = ref(0)
const activityTotalPages = computed(() => Math.ceil(activityTotal.value / activityPerPage.value))
const notifications = ref([])
const profile = reactive(JSON.parse(localStorage.getItem('ncsms_user') || '{}'))
const applicationSearch = ref('')
const applicationStatus = ref('')
const savingProfile = ref(false)
const uploadingAvatar = ref(false)
const passwordCurrent = ref('')
const passwordNew = ref('')
const passwordConfirm = ref('')
const showPasswordCurrent = ref(false)
const showPasswordNew = ref(false)
const showPasswordConfirm = ref(false)
const changingPassword = ref(false)
const twofaEnabled = ref(false)
const toggling2FA = ref(false)
const showTwoFAVerify = ref(false)
const twofaVerifyCode = ref('')
const verifying2FA = ref(false)

const settingsTab = ref('password')
const sessions = ref([])
const loadingSessions = ref(false)
const revokingSessionId = ref('')
const revokingAllOtherSessions = ref(false)

const prefEmailAppStatus = ref(localStorage.getItem('ncs_pref_app_status') !== 'false')
const prefEmailSecurity = ref(localStorage.getItem('ncs_pref_security') !== 'false')
const prefEmailAnnouncements = ref(localStorage.getItem('ncs_pref_announcements') === 'true')
const savingPreferences = ref(false)
const preferencesSaved = ref(false)

const activityCategoryFilter = ref('ALL')
const activitySearchQuery = ref('')
const loadingActivities = ref(false)

const isAthlete = ref(false)
const athleteData = ref(null)
const athleteMedical = ref(null)
const athleteSafeguarding = ref(null)
const athleteAntiDoping = ref(null)
const athleteNationalTeam = ref(null)
const athleteResults = ref([])
const athleteMedals = ref([])

const isCoach = ref(false)
const coachData = ref(null)
const isOfficial = ref(false)
const officialData = ref(null)

ensureOtikaStyles()

const userRoles = computed(() => {
  const roles = profile.roles || profile.role || []
  const roleList = Array.isArray(roles) ? roles : [roles]
  return roleList.map(r => typeof r === 'string' ? r : r?.name).filter(Boolean)
})
const isAdminOrGenSec = computed(() => {
  return userRoles.value.some(r => ['admin', 'super_admin', 'general_secretary'].includes(r))
})

const userProfileLabel = computed(() => {
  const roles = userRoles.value
  if (roles.includes('athlete') || roles.includes('role_athlete')) return 'Athlete Portal'
  if (roles.includes('coach') || roles.includes('role_coach')) return 'Coach Portal'
  if (roles.includes('technical_official') || roles.includes('role_technical_official')) return 'Official Portal'
  return 'Ordinary User'
})

const userWorkspaceLabel = computed(() => {
  const roles = userRoles.value
  if (roles.includes('athlete') || roles.includes('role_athlete')) return 'Athlete workspace'
  if (roles.includes('coach') || roles.includes('role_coach')) return 'Coach workspace'
  if (roles.includes('technical_official') || roles.includes('role_technical_official')) return 'Official workspace'
  return 'Applicant workspace'
})

const navigation = computed(() => {
  const items = [
    { id: 'dashboard', label: 'Dashboard', icon: 'icofont-dashboard-web' },
    { id: 'apply', label: 'Apply Now', icon: 'icofont-plus-circle', badge: openForms.value.length || '' },
    { id: 'applications', label: 'My Applications', icon: 'icofont-file-document' },
  ]
  if (!isAdminOrGenSec.value) {
    items.push({ id: 'my-files', label: 'My Files', icon: 'icofont-folder-open' })
  }
  items.push(
    { id: 'activities', label: 'My Activities', icon: 'icofont-history' },
    { id: 'notifications', label: 'Notifications', icon: 'icofont-notification', badge: unreadNotifications.value || '' },
    { id: 'messages', label: 'Messages', icon: 'icofont-envelope', badge: messageCount.value || '' },
    { id: 'transactions', label: 'My Transactions', icon: 'icofont-money' },
    { id: 'settings', label: 'Settings', icon: 'icofont-gear' },
    { id: 'profile', label: 'My Profile', icon: 'icofont-user-alt-3' }
  )
  return items
})
const sectionTitle = computed(() => navigation.value.find(item => item.id === section.value)?.label || 'Dashboard')
const firstName = computed(() => profile.first_name || String(profile.email || 'User').split('@')[0])
const fullName = computed(() => `${profile.first_name || ''} ${profile.last_name || ''}`.trim() || profile.email || 'Portal user')
const initials = computed(() => `${profile.first_name?.[0] || ''}${profile.last_name?.[0] || ''}`.toUpperCase() || 'U')
const profileAvatar = computed(() => profile.avatar_url ? mediaUrl(profile.avatar_url) : '')
const sidebarUserName = computed(() => {
  const value = fullName.value || firstName.value || 'Portal user'
  return value.length > 18 ? `${value.slice(0, 15)}...` : value
})
const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Good morning'
  if (hour < 18) return 'Good afternoon'
  return 'Good evening'
})
const allApplications = computed(() => [
  ...dynamicSubmissions.value.map(item => ({ ...item, source: 'custom', title: item.template_title || 'Custom application', reference: item.submission_reference })),
  ...legacyApplications.value.map(item => ({ ...item, source: 'standard', title: titleize(item.application_type || item.form_type), reference: item.application_reference })),
].sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at)))
const filteredApplications = computed(() => allApplications.value.filter(item => {
  if (applicationStatus.value && item.status !== applicationStatus.value) return false
  const needle = applicationSearch.value.trim().toLowerCase()
  return !needle || [item.title, item.reference, item.status].some(value => String(value || '').toLowerCase().includes(needle))
}))
const applicationStatuses = computed(() => [...new Set(allApplications.value.map(item => item.status).filter(Boolean))].sort())

const passwordStrengthScore = computed(() => {
  const p = passwordNew.value || ''
  if (!p) return 0
  let score = 0
  if (p.length >= 8) score++
  if (/[A-Z]/.test(p) && /[a-z]/.test(p)) score++
  if (/[0-9]/.test(p)) score++
  if (/[^A-Za-z0-9]/.test(p)) score++
  return score
})

const passwordStrengthPercent = computed(() => (passwordStrengthScore.value / 4) * 100)

const passwordStrengthLabel = computed(() => {
  if (!passwordNew.value) return ''
  switch (passwordStrengthScore.value) {
    case 1: return 'Weak'
    case 2: return 'Fair'
    case 3: return 'Good'
    case 4: return 'Strong'
    default: return 'Very Weak'
  }
})

const passwordStrengthClass = computed(() => {
  switch (passwordStrengthScore.value) {
    case 1: return 'strength-weak'
    case 2: return 'strength-fair'
    case 3: return 'strength-good'
    case 4: return 'strength-strong'
    default: return 'strength-weak'
  }
})

const passwordStrengthTextClass = computed(() => {
  switch (passwordStrengthScore.value) {
    case 1: return 'text-danger'
    case 2: return 'text-warning'
    case 3: return 'text-info'
    case 4: return 'text-success'
    default: return 'text-muted'
  }
})

const filteredActivities = computed(() => {
  return activities.value.filter(item => {
    if (activityCategoryFilter.value !== 'ALL') {
      const cat = getActivityCategory(item)
      if (activityCategoryFilter.value === 'AUTH' && cat !== 'AUTH') return false
      if (activityCategoryFilter.value === 'APPLICATIONS' && cat !== 'APPLICATION') return false
      if (activityCategoryFilter.value === 'PROFILE' && cat !== 'PROFILE') return false
    }
    const query = activitySearchQuery.value.trim().toLowerCase()
    if (query) {
      const matchText = [
        activityTitle(item),
        activityDescription(item),
        item.ip_address,
        item.geo_city,
        item.geo_country,
        item.browser,
        item.os_name
      ].filter(Boolean).join(' ').toLowerCase()
      return matchText.includes(query)
    }
    return true
  })
})

const activityStats = computed(() => {
  const items = activities.value
  let authCount = 0
  let appCount = 0
  items.forEach(item => {
    const cat = getActivityCategory(item)
    if (cat === 'AUTH') authCount++
    if (cat === 'APPLICATION') appCount++
  })
  const latest = items[0]?.created_at ? formatDate(items[0].created_at) : 'None'
  return {
    authCount,
    appCount,
    lastActive: latest
  }
})
const userKpis = computed(() => {
  const apps = allApplications.value
  const inReview = apps.filter(item => ['SUBMITTED', 'RESUBMITTED', 'UNDER_REVIEW'].includes(item.status)).length
  const approved = apps.filter(item => item.status === 'APPROVED').length
  return [
    { label: 'Applications', value: apps.length, note: 'Total started', icon: 'icofont-file-document', tone: 'blue' },
    { label: 'In Review', value: inReview, note: 'With NCS reviewers', icon: 'icofont-clock-time', tone: 'amber' },
    { label: 'Approved', value: approved, note: 'Successful applications', icon: 'icofont-check-circled', tone: 'green' },
    { label: 'Transactions', value: transactions.value.length, note: 'Payment records', icon: 'icofont-money', tone: 'cyan' },
  ]
})
const unreadNotifications = computed(() => notifications.value.filter(item => item.status === 'unread').length)
const userMessages = computed(() => allApplications.value.filter(item => item.review_notes).map(item => ({
  id: `review-${item.id}`, title: `${item.title}: ${titleize(item.status)}`, message: item.review_notes, created_at: item.updated_at,
})))
const messageCount = computed(() => userMessages.value.length)
const transactions = computed(() => [
  ...dynamicSubmissions.value.filter(item => item.payment_status && item.payment_status !== 'UNPAID').map(item => ({
    ...item, source: 'custom', title: item.template_title || 'Custom application',
    amount: item.payment_amount_ugx || 0,
  })),
  ...legacyTransactions.value.map(item => ({ ...item, source: 'standard', title: titleize(item.form_type), amount: item.payment_amount_ugx || 0 })),
])
const transactionTotal = computed(() => transactions.value.reduce((sum, item) => sum + Number(item.amount || 0), 0))

const userFiles = computed(() => {
  const filesList = []
  
  if (isAthlete.value && athleteData.value) {
    filesList.push({
      id: 'cert_athlete_license',
      title: 'National Athlete License Certificate',
      type: 'License / Certificate',
      number: athleteData.value.athlete_number,
      issueDate: athleteData.value.created_at || new Date().toISOString(),
      expiryDate: 'N/A (Active)',
      status: athleteData.value.status || 'Active',
      description: `Official NCS verification for athlete classification under ${athleteData.value.discipline || 'sports registry'}.`,
      fileType: 'TXT Document',
      downloadName: `NCS_Athlete_License_${athleteData.value.athlete_number || 'Cert'}.txt`,
      category: 'Athlete Registry'
    })
  }

  if (isCoach.value && coachData.value) {
    filesList.push({
      id: 'cert_coach_license',
      title: `NCS Coach License - Level ${coachData.value.certification_level || 'Certified'}`,
      type: 'License / Certificate',
      number: coachData.value.license_number || 'NCS-COACH-TEMP',
      issueDate: coachData.value.created_at || new Date().toISOString(),
      expiryDate: coachData.value.expiry_date || 'N/A',
      status: coachData.value.status || 'Active',
      description: `NCS recognized coaching qualifications and certification credentials.`,
      fileType: 'TXT Document',
      downloadName: `NCS_Coach_License_${coachData.value.license_number || 'Cert'}.txt`,
      category: 'Coaches Registry'
    })
  }

  if (isOfficial.value && officialData.value) {
    filesList.push({
      id: 'cert_official_license',
      title: `Technical Official Certification - ${officialData.value.official_type || 'Official'}`,
      type: 'Official Credentials',
      number: `NCS-TO-${officialData.value.id?.substring(0, 8).toUpperCase() || 'TEMP'}`,
      issueDate: officialData.value.created_at || new Date().toISOString(),
      expiryDate: officialData.value.valid_until || 'N/A',
      status: officialData.value.status || 'Active',
      description: `Official registration for NCS Technical Officials and Referees.`,
      fileType: 'TXT Document',
      downloadName: `NCS_Official_Credentials_${officialData.value.id || 'Cert'}.txt`,
      category: 'Technical Officials'
    })
  }

  if (athleteMedical.value) {
    filesList.push({
      id: 'cert_medical_clearance',
      title: 'Athlete Medical Clearance File',
      type: 'Medical Records',
      number: `NCS-MED-${athleteMedical.value.id?.substring(0, 8).toUpperCase() || 'TEMP'}`,
      issueDate: athleteMedical.value.created_at || new Date().toISOString(),
      expiryDate: 'N/A',
      status: athleteMedical.value.current_injury_status || 'Fit to Compete',
      description: `NCS Medical Department validation and clearance logs.`,
      fileType: 'TXT Document',
      downloadName: `NCS_Medical_Clearance_${athleteMedical.value.id || 'Record'}.txt`,
      category: 'Medical Files'
    })
  }

  if (athleteAntiDoping.value) {
    filesList.push({
      id: 'cert_antidoping_clearance',
      title: 'WADA Anti-Doping Compliance Certificate',
      type: 'Compliance Record',
      number: `NCS-WADA-${athleteAntiDoping.value.id?.substring(0, 8).toUpperCase() || 'TEMP'}`,
      issueDate: athleteAntiDoping.value.last_tested_on || new Date().toISOString(),
      expiryDate: 'N/A',
      status: athleteAntiDoping.value.last_test_result || 'Compliant',
      description: `Verification certificate for completion of WADA Anti-Doping education and compliance test logs.`,
      fileType: 'TXT Document',
      downloadName: `NCS_AntiDoping_Certificate_${athleteAntiDoping.value.id || 'Record'}.txt`,
      category: 'Anti-Doping Compliance'
    })
  }

  if (athleteNationalTeam.value) {
    filesList.push({
      id: 'cert_national_team_cap',
      title: `National Team Appearance Certificate (${athleteNationalTeam.value.team_name || 'Uganda National Team'})`,
      type: 'National Representation',
      number: `NCS-NT-${athleteNationalTeam.value.id?.substring(0, 8).toUpperCase() || 'TEMP'}`,
      issueDate: athleteNationalTeam.value.first_call_up_on || new Date().toISOString(),
      expiryDate: 'N/A',
      status: 'Verified',
      description: `Official NCS certification recognizing sports representation at national squad tier: ${athleteNationalTeam.value.category || 'National'}.`,
      fileType: 'TXT Document',
      downloadName: `NCS_National_Duty_Certificate_${athleteNationalTeam.value.id || 'Record'}.txt`,
      category: 'National Squads'
    })
  }

  return filesList
})

function getFileIcon(category) {
  switch (category) {
    case 'Athlete Registry': return 'icofont-runner-alt-1'
    case 'Coaches Registry': return 'icofont-whistle'
    case 'Technical Officials': return 'icofont-referee'
    case 'Medical Files': return 'icofont-first-aid'
    case 'Anti-Doping Compliance': return 'icofont-test-bulb'
    case 'National Squads': return 'icofont-flag'
    default: return 'icofont-document-folder'
  }
}

function getFileIconClass(category) {
  switch (category) {
    case 'Athlete Registry': return 'bg-primary text-white'
    case 'Coaches Registry': return 'bg-warning text-dark'
    case 'Technical Officials': return 'bg-info text-white'
    case 'Medical Files': return 'bg-danger text-white'
    case 'Anti-Doping Compliance': return 'bg-success text-white'
    case 'National Squads': return 'bg-dark text-white'
    default: return 'bg-secondary text-white'
  }
}

function getStatusBadgeClass(status) {
  const s = String(status).toUpperCase()
  if (['ACTIVE', 'COMPLIANT', 'VERIFIED', 'FIT TO COMPETE', 'NEGATIVE'].includes(s)) {
    return 'badge-success'
  }
  if (['PENDING', 'PENDING CONSENT'].includes(s)) {
    return 'badge-warning'
  }
  return 'badge-danger'
}

function getMedalBadgeClass(type) {
  const t = String(type).toUpperCase()
  if (t.includes('GOLD')) return 'badge badge-warning text-dark text-uppercase'
  if (t.includes('SILVER')) return 'badge bg-secondary text-white text-uppercase'
  return 'badge bg-bronze text-white text-uppercase'
}

function downloadFile(file) {
  const content = `========================================================================
                      NATIONAL COUNCIL OF SPORTS (NCS) UGANDA
                                OFFICIAL CERTIFICATE
========================================================================

CERTIFICATE TITLE : ${file.title}
DOCUMENT TYPE     : ${file.type}
LICENSE/REF NO.   : ${file.number}
CATEGORY          : ${file.category}
STATUS            : ${file.status}

ISSUED TO         : ${fullName.value}
EMAIL ADDRESS     : ${profile.email}
DATE OF ISSUE     : ${formatDate(file.issueDate)}
EXPIRY DATE       : ${formatDate(file.expiryDate)}

------------------------------------------------------------------------
DESCRIPTION:
${file.description}
------------------------------------------------------------------------

VERIFICATION STATUS: VERIFIED BY NATIONAL COUNCIL OF SPORTS (NCS)
This document serves as the official digital credential issued by the National
Council of Sports (NCS) Uganda portal. To verify, contact ncs@ncs.go.ug.

Generated on      : ${new Date().toLocaleString()}
========================================================================`

  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = file.downloadName
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  setMsg(`Successfully downloaded ${file.title}`)
}

watch(section, id => {
  const query = id === 'dashboard' ? {} : { section: id }
  router.replace({ path: '/dashboard', query }).catch(() => {})
})

onMounted(async () => {
  if (!localStorage.getItem('ncsms_access_token')) return router.replace('/login')
  const requested = router.currentRoute.value.query.section
  if (requested && portalSectionIds.includes(requested)) section.value = requested
  await loadPortal()
})

async function loadPortal() {
  loading.value = true
  error.value = ''
  try {
    const results = await Promise.allSettled([
      getCurrentUser(), portalListOpenForms(), portalListSubmissions({ page: 1, per_page: 200 }),
      listMyLegacyApplications({ page: 1, per_page: 200 }), listMyTransactions({ page: 1, per_page: 200 }),
      listMyAuditLogs({ page: 1, per_page: activityPerPage.value }), cms.listNotifications({ page: 1, per_page: 100 }),
    ])
    if (results[0].status === 'fulfilled') {
      Object.assign(profile, unwrap(results[0].value))
      localStorage.setItem('ncsms_user', JSON.stringify(profile))
    }
    openForms.value = results[1].status === 'fulfilled' ? asList(results[1].value) : []
    dynamicSubmissions.value = results[2].status === 'fulfilled' ? asList(results[2].value) : []
    legacyApplications.value = results[3].status === 'fulfilled' ? asList(results[3].value) : []
    legacyTransactions.value = results[4].status === 'fulfilled' ? asList(results[4].value) : []
    if (results[5].status === 'fulfilled') {
      const unwrapped = results[5].value?.data ?? results[5].value ?? {}
      activities.value = Array.isArray(unwrapped.data) ? unwrapped.data : []
      activityTotal.value = unwrapped.meta?.total ?? activities.value.length
    } else {
      activities.value = []
      activityTotal.value = 0
    }
    notifications.value = results[6].status === 'fulfilled' ? asList(results[6].value) : []
  } catch (err) {
    console.warn('Initial dashboard load warning:', err)
  } finally {
    loading.value = false
  }

  isAthlete.value = false
  athleteData.value = null
  athleteMedical.value = null
  athleteSafeguarding.value = null
  athleteAntiDoping.value = null
  athleteNationalTeam.value = null
  athleteResults.value = []
  athleteMedals.value = []
  isCoach.value = false
  coachData.value = null
  isOfficial.value = false
  officialData.value = null

  try {
    const sec = await getMySecurity()
    twofaEnabled.value = sec?.twofa?.enabled ?? false
  } catch (e) {}

  // Run sports registry lookups asynchronously in the background
  loadSportsRegistryContext()
}

async function loadSportsRegistryContext() {
  if (!profile || !profile.email) return

  const roles = userRoles.value
  const hasRegistryRole = roles.some(r =>
    ['athlete', 'role_athlete', 'coach', 'role_coach', 'technical_official', 'role_technical_official', 'admin', 'super_admin', 'general_secretary', 'federation_admin'].includes(r)
  )
  if (!hasRegistryRole) return

  try {
    const athletesRes = await listNsmisDomain('athletes', { search: profile.email })
    const athletesList = asList(athletesRes)
    const match = athletesList.find(ath => String(ath.email_address || '').toLowerCase() === String(profile.email || '').toLowerCase())
    if (match) {
      isAthlete.value = true
      athleteData.value = match
      
      const athleteId = match.id
      const [medicalRes, safeguardingRes, antidopingRes, nationalTeamRes, resultsRes, medalsRes] = await Promise.allSettled([
        listNsmisDomain('medical-records', { search: athleteId }),
        listNsmisDomain('safeguarding-records', { search: athleteId }),
        listNsmisDomain('anti-doping', { search: athleteId }),
        listNsmisDomain('national-team', { search: athleteId }),
        listNsmisDomain('competition-results', { search: athleteId }),
        listNsmisDomain('medals', { search: athleteId })
      ])
      
      if (medicalRes.status === 'fulfilled') {
        const medItems = asList(medicalRes.value)
        athleteMedical.value = medItems.find(r => r.athlete_id === athleteId) || null
      }
      if (safeguardingRes.status === 'fulfilled') {
        const sgItems = asList(safeguardingRes.value)
        athleteSafeguarding.value = sgItems.find(r => r.athlete_id === athleteId) || null
      }
      if (antidopingRes.status === 'fulfilled') {
        const adItems = asList(antidopingRes.value)
        athleteAntiDoping.value = adItems.find(r => r.athlete_id === athleteId) || null
      }
      if (nationalTeamRes.status === 'fulfilled') {
        const ntItems = asList(nationalTeamRes.value)
        athleteNationalTeam.value = ntItems.find(r => r.athlete_id === athleteId) || null
      }
      if (resultsRes.status === 'fulfilled') {
        const resItems = asList(resultsRes.value)
        athleteResults.value = resItems.filter(r => r.athlete_id === athleteId)
      }
      if (medalsRes.status === 'fulfilled') {
        const medItems = asList(medalsRes.value)
        athleteMedals.value = medItems.filter(r => r.athlete_id === athleteId)
      }
    }
  } catch (e) {}

  try {
    const coachesRes = await listNsmisDomain('coaches', { search: profile.email })
    const coachesList = asList(coachesRes)
    const matchCoach = coachesList.find(c => String(c.email || '').toLowerCase() === String(profile.email || '').toLowerCase())
    if (matchCoach) {
      isCoach.value = true
      coachData.value = matchCoach
    }
  } catch (e) {}

  try {
    const officialsRes = await listNsmisDomain('technical-officials', { search: profile.email })
    const officialsList = asList(officialsRes)
    const nameKey = fullName.value.toLowerCase().trim()
    const matchOfficial = officialsList.find(o => String(o.full_name || '').toLowerCase().trim() === nameKey)
    if (matchOfficial) {
      isOfficial.value = true
      officialData.value = matchOfficial
    }
  } catch (e) {}
}

async function loadActivities(page = 1) {
  activityPage.value = page
  loadingActivities.value = true
  try {
    const res = await listMyAuditLogs({ page: activityPage.value, per_page: activityPerPage.value })
    const unwrapped = res?.data ?? res ?? {}
    activities.value = Array.isArray(unwrapped.data) ? unwrapped.data : []
    activityTotal.value = unwrapped.meta?.total ?? activities.value.length
  } catch (err) {
    error.value = 'Could not load activity logs.'
  } finally {
    loadingActivities.value = false
  }
}

async function onSelectSessionsTab() {
  settingsTab.value = 'sessions'
  await loadSessions()
}

async function loadSessions() {
  loadingSessions.value = true
  try {
    const res = await listMySessions()
    sessions.value = Array.isArray(res) ? res : (res?.data || [])
  } catch (err) {
    console.warn('Could not load sessions:', err)
  } finally {
    loadingSessions.value = false
  }
}

async function revokeSession(sessionId) {
  if (!sessionId) return
  if (!confirm('Are you sure you want to revoke this session? The device will be signed out.')) return
  revokingSessionId.value = sessionId
  try {
    await revokeMySession(sessionId)
    success.value = 'Session revoked successfully.'
    await loadSessions()
  } catch (err) {
    error.value = apiError(err, 'Could not revoke session.')
  } finally {
    revokingSessionId.value = ''
  }
}

async function revokeAllOtherSessions() {
  if (!confirm('Are you sure you want to sign out of all other devices? All other active sessions will be terminated.')) return
  revokingAllOtherSessions.value = true
  try {
    await revokeOtherSessions()
    success.value = 'All other devices have been signed out.'
    await loadSessions()
  } catch (err) {
    error.value = apiError(err, 'Could not revoke other sessions.')
  } finally {
    revokingAllOtherSessions.value = false
  }
}

function getSessionDeviceIcon(ua) {
  const str = String(ua || '').toLowerCase()
  if (str.includes('iphone') || str.includes('android') || str.includes('mobile')) return 'icofont-smart-phone'
  if (str.includes('ipad') || str.includes('tablet')) return 'icofont-tablet'
  if (str.includes('macintosh') || str.includes('mac os')) return 'icofont-brand-apple'
  if (str.includes('windows')) return 'icofont-brand-windows'
  if (str.includes('linux')) return 'icofont-brand-linux'
  return 'icofont-laptop'
}

function parseSessionBrowser(ua) {
  const str = String(ua || '')
  if (!str) return 'Web Browser'
  let browser = 'Browser'
  if (str.includes('Edg/')) browser = 'Microsoft Edge'
  else if (str.includes('Chrome/')) browser = 'Google Chrome'
  else if (str.includes('Safari/') && !str.includes('Chrome/')) browser = 'Apple Safari'
  else if (str.includes('Firefox/')) browser = 'Mozilla Firefox'
  
  let os = ''
  if (str.includes('Windows NT 10.0')) os = 'Windows 10/11'
  else if (str.includes('Windows')) os = 'Windows'
  else if (str.includes('Macintosh')) os = 'macOS'
  else if (str.includes('iPhone')) os = 'iOS'
  else if (str.includes('Android')) os = 'Android'
  else if (str.includes('Linux')) os = 'Linux'
  
  return os ? `${browser} on ${os}` : browser
}

async function savePreferences() {
  savingPreferences.value = true
  preferencesSaved.value = false
  try {
    localStorage.setItem('ncs_pref_app_status', String(prefEmailAppStatus.value))
    localStorage.setItem('ncs_pref_security', String(prefEmailSecurity.value))
    localStorage.setItem('ncs_pref_announcements', String(prefEmailAnnouncements.value))
    preferencesSaved.value = true
    success.value = 'Preferences saved successfully.'
    setTimeout(() => { preferencesSaved.value = false }, 3000)
  } catch (err) {
    error.value = 'Could not save preferences.'
  } finally {
    savingPreferences.value = false
  }
}

function getActivityCategory(item) {
  const endpoint = String(item.endpoint || '').toLowerCase()
  const eventType = String(item.event_type || '').toUpperCase()
  if (eventType.includes('AUTH') || endpoint.includes('/auth') || endpoint.includes('/2fa') || endpoint.includes('/security')) return 'AUTH'
  if (endpoint.includes('/forms') || endpoint.includes('/submissions') || endpoint.includes('/applications')) return 'APPLICATION'
  if (endpoint.includes('/transactions') || endpoint.includes('/payment')) return 'FINANCE'
  if (endpoint.includes('/profile') || endpoint.includes('/users/me')) return 'PROFILE'
  return 'SYSTEM'
}

function getActivityIcon(item) {
  const cat = getActivityCategory(item)
  switch (cat) {
    case 'AUTH': return 'icofont-key'
    case 'APPLICATION': return 'icofont-file-document'
    case 'FINANCE': return 'icofont-money'
    case 'PROFILE': return 'icofont-ui-user'
    default: return 'icofont-history'
  }
}

function getActivityTone(item) {
  const cat = getActivityCategory(item)
  switch (cat) {
    case 'AUTH': return 'blue'
    case 'APPLICATION': return 'amber'
    case 'FINANCE': return 'green'
    case 'PROFILE': return 'cyan'
    default: return 'purple'
  }
}

function formatRelativeTime(dateString) {
  if (!dateString) return '-'
  const date = new Date(dateString)
  const now = new Date()
  const diffSecs = Math.floor((now - date) / 1000)
  if (diffSecs < 60) return 'Just now'
  if (diffSecs < 3600) return `${Math.floor(diffSecs / 60)} mins ago`
  if (diffSecs < 86400) return `${Math.floor(diffSecs / 3600)} hours ago`
  if (diffSecs < 604800) return `${Math.floor(diffSecs / 86400)} days ago`
  return date.toLocaleDateString('en-UG', { month: 'short', day: 'numeric', year: 'numeric' })
}

async function select(id) {
  const previous = section.value
  const next = navigation.value.find(item => item.id === id)
  recordMenuNavigation({
    portal: 'Applicant Portal',
    fromSection: previous,
    toSection: id,
    label: next?.label || id,
    basePath: '/dashboard',
    routeName: 'UserDashboard',
  })
  section.value = id
  profileOpen.value = false
  mobileSidebarOpen.value = false
  success.value = ''
  error.value = ''
  if (id === 'activities') {
    loadActivities(1)
  }
  if (id === 'settings' || id === 'profile') {
    try {
      const sec = await getMySecurity()
      twofaEnabled.value = sec?.twofa?.enabled ?? false
    } catch (e) {}
    if (id === 'settings' && settingsTab.value === 'sessions') {
      loadSessions()
    }
  }
}

function onMenuToggle() {
  if (window.innerWidth <= 991) {
    mobileSidebarOpen.value = !mobileSidebarOpen.value
  } else {
    sidebarCollapsed.value = !sidebarCollapsed.value
    localStorage.setItem('ncsms_sidebar_collapsed', String(sidebarCollapsed.value))
  }
}

async function startApplication(form) {
  if (!form?.slug) {
    error.value = 'This application form is not available.'
    return
  }
  const pending = dynamicSubmissions.value.find(item => item.template_id === form.id && isPendingSubmission(item))
  if (pending) {
    error.value = `You already have a pending application for ${form.title}.`
    select('applications')
    return
  }
  router.push({ name: 'ApplicationWizard', params: { slug: form.slug } })
}
async function continueApplication(item) {
  const form = openForms.value.find(entry => entry.id === item.template_id)
  if (form) startApplication(form)
  else error.value = 'This application form is no longer open.'
}
function viewApplication(item) {
  router.push({ name: 'UserApplicationDetail', params: { id: item.id }, query: { source: item.source } })
}
function downloadApplication(item) {
  downloadApplicationForm(item)
}
async function refreshSubmissions() {
  const result = await portalListSubmissions({ page: 1, per_page: 200 })
  dynamicSubmissions.value = asList(result)
}
async function markNotificationRead(item) {
  try {
    await cms.updateNotification(item.id, { status: 'read' })
    item.status = 'read'
  } catch (err) { error.value = apiError(err, 'Could not update notification.') }
}
async function markAllNotificationsRead() {
  try {
    await cms.markAllNotificationsRead()
    notifications.value.forEach(item => { item.status = 'read' })
  } catch (err) { error.value = apiError(err, 'Could not update notifications.') }
}
async function saveProfile() {
  savingProfile.value = true
  try {
    const res = await updateMyProfile({ first_name: profile.first_name, last_name: profile.last_name, avatar_url: profile.avatar_url || '' })
    Object.assign(profile, unwrap(res))
    localStorage.setItem('ncsms_user', JSON.stringify(profile))
    success.value = 'Your profile was updated.'
  } catch (err) { error.value = apiError(err, 'Could not update your profile.') }
  finally { savingProfile.value = false }
}
function logout() { localStorage.removeItem('ncsms_access_token'); localStorage.removeItem('ncsms_user'); router.push('/login') }
function unwrap(value) { return value?.data?.data ?? value?.data ?? value ?? {} }
function asList(value) {
  const data = unwrap(value)
  if (Array.isArray(data)) return data
  if (data && Array.isArray(data.items)) return data.items
  return []
}
function apiError(err, fallback) { return err.response?.data?.error?.message || fallback }
function titleize(value) { return String(value || '').toLowerCase().replaceAll('_', ' ').replaceAll('-', ' ').replace(/\b\w/g, char => char.toUpperCase()) }
function formatDate(value) { return value ? new Intl.DateTimeFormat('en-UG', { dateStyle: 'medium', timeStyle: 'medium' }).format(new Date(value)) : '-' }
function formatMoney(value) { return new Intl.NumberFormat('en-UG', { maximumFractionDigits: 0 }).format(Number(value || 0)) }
function statusClass(status) { if (status === 'APPROVED' || status === 'COMPLETE') return 'green'; if (status === 'REJECTED') return 'red'; if (['NEEDS_INFORMATION', 'PENDING_PAYMENT'].includes(status)) return 'amber'; return 'blue' }
function paymentClass(status) { return status === 'PAID' || status === 'VERIFIED' ? 'green' : ['REJECTED', 'VERIFICATION_FAILED'].includes(status) ? 'red' : 'amber' }
function notificationIcon(item) { return item.icon_key ? `icofont-${item.icon_key}` : 'icofont-notification' }
function activityTitle(item) {
  const method = String(item.method || '').toUpperCase()
  const endpoint = String(item.endpoint || '').toLowerCase()
  const eventType = String(item.event_type || '').toUpperCase()
  if (eventType === 'AUTH_LOGIN' || endpoint.includes('/auth/login')) return 'Logged In'
  if (eventType === 'AUTH_LOGOUT' || endpoint.includes('/auth/logout')) return 'Logged Out'
  if (endpoint.includes('/auth/register')) return 'Registered Account'
  if (endpoint.includes('/auth/password/reset')) return 'Requested Password Reset'
  if (endpoint.includes('/portal/forms/open')) return 'Viewed Services'
  if (endpoint.includes('/portal/forms/')) return 'Started Application Draft'
  if (endpoint.includes('/portal/submissions')) {
    if (method === 'POST') return 'Created Application Draft'
    return 'Listed Applications'
  }
  if (endpoint.includes('/applications')) {
    if (method === 'POST') return 'Submitted Application'
    return 'Viewed Applications List'
  }
  if (endpoint.includes('/transactions')) {
    if (method === 'POST') return 'Uploaded Payment Proof'
    return 'Viewed Payments List'
  }
  if (endpoint.includes('/notifications')) {
    if (method === 'PUT' || method === 'POST') return 'Updated Notifications'
    return 'Viewed Notifications'
  }
  if (endpoint.includes('/audit-logs')) return 'Viewed Security Audit Logs'
  if (endpoint.includes('/account/profile') || endpoint.includes('/users/me')) {
    if (method === 'PUT' || method === 'POST' || method === 'PATCH') return 'Updated Profile'
    return 'Viewed Profile Page'
  }
  if (item.event_type) return titleize(item.event_type)
  return 'Portal Interaction'
}
function activityDescription(item) {
  const endpoint = String(item.endpoint || '').toLowerCase()
  let pageName = 'Dashboard'
  if (endpoint.includes('/profile') || endpoint.includes('/users/me')) pageName = 'Profile page'
  else if (endpoint.includes('/forms/open')) pageName = 'Apply page'
  else if (endpoint.includes('/submissions') || endpoint.includes('/applications')) pageName = 'Applications page'
  else if (endpoint.includes('/transactions')) pageName = 'Transactions page'
  else if (endpoint.includes('/audit-logs')) pageName = 'Security page'
  else if (endpoint.includes('/notifications')) pageName = 'Notifications center'
  else if (endpoint.includes('/messages')) pageName = 'Messages center'

  let location = ''
  const city = String(item.geo_city || '').trim()
  const country = String(item.geo_country || '').trim()
  if (city && country) {
    const displayCity = city.toLowerCase() === 'internal' ? 'Local network' : city
    location = `${displayCity}, ${country}`
  } else if (country) {
    location = country
  } else {
    location = 'Unknown location'
  }

  let device = ''
  const browser = String(item.browser || '').trim()
  const os = String(item.os_name || '').trim()
  if (browser && os) {
    device = `on ${os} ${browser} browser`
  } else if (browser) {
    device = `on ${browser} browser`
  } else if (os) {
    device = `on ${os} system`
  } else {
    device = 'on web browser'
  }

  const ip = item.ip_address ? ` (IP: ${item.ip_address})` : ''
  return `${pageName}, ${location}, ${device}${ip}`
}
function isEditableSubmission(item) { return ['DRAFT', 'NEEDS_INFORMATION'].includes(item?.status) }
function isPendingSubmission(item) { return item?.status && !['DRAFT', 'NEEDS_INFORMATION', 'APPROVED', 'REJECTED'].includes(item.status) }

async function onAvatarFileSelected(e) {
  const file = e.target.files?.[0]
  if (!file) return
  uploadingAvatar.value = true
  error.value = ''
  success.value = ''
  try {
    const formData = new FormData()
    formData.append('file', file)
    const res = await uploadProfileAvatar(formData)
    profile.avatar_url = unwrap(res).avatar_url || ''
    success.value = 'Profile photo updated successfully!'
    
    const stored = localStorage.getItem('ncsms_user')
    if (stored) {
      const parsed = JSON.parse(stored)
      parsed.avatar_url = profile.avatar_url
      localStorage.setItem('ncsms_user', JSON.stringify(parsed))
    }
  } catch (err) {
    error.value = err.response?.data?.error?.message || 'Could not upload avatar image.'
  } finally {
    uploadingAvatar.value = false
  }
}

async function updatePassword() {
  if (passwordNew.value.length < 8) {
    error.value = 'New password must be at least 8 characters long.'
    return
  }
  if (passwordNew.value !== passwordConfirm.value) {
    error.value = 'Confirm password does not match.'
    return
  }
  changingPassword.value = true
  error.value = ''
  success.value = ''
  try {
    await changePassword(passwordCurrent.value, passwordNew.value)
    success.value = 'Password changed successfully.'
    passwordCurrent.value = ''
    passwordNew.value = ''
    passwordConfirm.value = ''
  } catch (err) {
    error.value = err.response?.data?.error?.message || 'Could not change password. Check your current password.'
  } finally {
    changingPassword.value = false
  }
}

async function toggle2FA(e) {
  const checked = e.target.checked
  if (checked) {
    toggling2FA.value = true
    error.value = ''
    success.value = ''
    try {
      await enroll2FA()
      showTwoFAVerify.value = true
      twofaVerifyCode.value = ''
      success.value = 'Verification code sent to your email.'
    } catch (err) {
      error.value = err.response?.data?.error?.message || 'Could not enroll 2FA.'
    } finally {
      toggling2FA.value = false
    }
  } else {
    toggling2FA.value = true
    error.value = ''
    success.value = ''
    try {
      await disable2FA()
      twofaEnabled.value = false
      showTwoFAVerify.value = false
      success.value = 'Two-factor authentication disabled.'
    } catch (err) {
      error.value = err.response?.data?.error?.message || 'Could not disable 2FA.'
    } finally {
      toggling2FA.value = false
    }
  }
}

async function confirm2FA() {
  if (twofaVerifyCode.value.length < 6) {
    error.value = 'Enter a valid 6-digit code.'
    return
  }
  verifying2FA.value = true
  error.value = ''
  success.value = ''
  try {
    await verify2FA(twofaVerifyCode.value)
    twofaEnabled.value = true
    showTwoFAVerify.value = false
    success.value = 'Two-factor authentication successfully enabled!'
  } catch (err) {
    error.value = err.response?.data?.error?.message || 'Verification failed. Incorrect code.'
  } finally {
    verifying2FA.value = false
  }
}

function cancel2FAEnrollment() {
  showTwoFAVerify.value = false
  twofaVerifyCode.value = ''
  disable2FA().catch(() => {})
}
</script>

<style scoped>
/* =====================================================================
   SELF-CONTAINED NCS USER PORTAL STYLES (OTIKA DESIGN SYSTEM)
   - Pure color indication on active navigation (NO black borders)
   - Crisp self-hosted icon support
   - Responsive KPI card grid, timelines, and data tables
   ===================================================================== */

/* ── Root & Layout ─────────────────────────────────────────────────── */
.user-portal {
  min-height: 100vh;
  background-color: #f4f6f9;
  font-family: 'Poppins', 'Inter', system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  font-size: 14px;
  color: #1e293b;
}
.otika-app { display: flex; flex-direction: column; min-height: 100vh; }
.main-wrapper {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* ── Navbar ────────────────────────────────────────────────────────── */
.navbar-bg {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 60px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,0,0,.08);
  z-index: 899;
}
.main-navbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,0,0,.08);
  z-index: 900;
}
.main-navbar .navbar-nav {
  display: flex;
  align-items: center;
  list-style: none;
  margin: 0; padding: 0;
  gap: 4px;
}
.main-navbar .navbar-right { margin-left: auto; }
.main-navbar .form-inline { display: flex; align-items: center; }
.main-navbar .mr-auto { margin-right: auto; }
.main-navbar .mr-3 { margin-right: 12px; }

/* ── Navigation Links & Buttons ────────────────────────────────────── */
.nav-link {
  background: none;
  border: none;
  outline: none !important;
  box-shadow: none !important;
  cursor: pointer;
  color: #1e293b;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 10px;
  border-radius: 8px;
  font-size: 14px;
  text-decoration: none;
  transition: background 0.15s, color 0.15s;
  white-space: nowrap;
}
.nav-link:hover { background: #f1f5f9; color: #4f46e5; }

/* ── Top Navbar Action Buttons (Vibrant, Sharp, High-Contrast) ───── */
.main-navbar .nav-link-lg.cms-top-icon,
.main-navbar button.cms-top-icon,
.main-navbar .portal-top-action,
.main-navbar .portal-theme-toggle {
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  width: 40px !important;
  height: 40px !important;
  padding: 0 !important;
  border-radius: 10px !important;
  border: 1px solid #cbd5e1 !important;
  background: #ffffff !important;
  color: #0f172a !important;
  font-size: 20px !important;
  position: relative !important;
  cursor: pointer !important;
  transition: all 0.18s ease !important;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06) !important;
}

.main-navbar .nav-link-lg.cms-top-icon i,
.main-navbar button.cms-top-icon i,
.main-navbar .portal-top-action i,
.main-navbar .portal-theme-toggle i {
  color: #0f172a !important;
  font-size: 20px !important;
  font-weight: 700 !important;
  line-height: 1 !important;
  transition: color 0.18s ease, transform 0.18s ease !important;
}

.main-navbar .nav-link-lg.cms-top-icon:hover,
.main-navbar button.cms-top-icon:hover,
.main-navbar .portal-top-action:hover,
.main-navbar .portal-theme-toggle:hover {
  background: #eef2ff !important;
  border-color: #818cf8 !important;
  color: #4338ca !important;
  transform: translateY(-1px) !important;
  box-shadow: 0 4px 12px rgba(79, 70, 229, 0.2) !important;
}

.main-navbar .nav-link-lg.cms-top-icon:hover i,
.main-navbar button.cms-top-icon:hover i,
.main-navbar .portal-top-action:hover i,
.main-navbar .portal-theme-toggle:hover i {
  color: #4338ca !important;
}

.headerBadge1, .headerBadge2 {
  position: absolute !important;
  top: -4px !important;
  right: -4px !important;
  min-width: 18px !important;
  height: 18px !important;
  border-radius: 999px !important;
  font-size: 10px !important;
  font-weight: 700 !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  padding: 0 4px !important;
  border: 2px solid #ffffff !important;
  box-shadow: 0 2px 5px rgba(0,0,0,0.18) !important;
}
.headerBadge1 { background: #ef4444 !important; color: #ffffff !important; }
.headerBadge2 { background: #6777ef !important; color: #ffffff !important; }

.portal-navbar-title {
  display: flex; flex-direction: column; line-height: 1.3;
  padding: 0 10px; color: #0f172a;
}
.portal-navbar-title small { font-size: 11px; color: #64748b; font-weight: 600; text-transform: uppercase; }
.portal-navbar-title strong { font-size: 15px; font-weight: 700; color: #0f172a; }

/* ── Avatar / Initials ─────────────────────────────────────────────── */
.user-img-radious-style {
  width: 34px; height: 34px;
  border-radius: 50%; object-fit: cover;
}
.nav-avatar-fallback, .sidebar-avatar-fallback {
  display: inline-flex;
  align-items: center; justify-content: center;
  width: 34px; height: 34px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6777ef, #4f46e5);
  color: #fff; font-weight: 700; font-size: 13px;
}

/* ── Dropdown Menu ─────────────────────────────────────────────────── */
.dropdown { position: relative; }
.dropdown-menu {
  display: none;
  position: absolute; top: calc(100% + 6px); right: 0;
  min-width: 200px;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0,0,0,.12);
  padding: 6px 0;
  z-index: 1000;
}
.dropdown-menu.show, .dropdown-menu.d-block { display: block; }
.dropdown-title {
  padding: 10px 16px 8px;
  font-size: 12px; font-weight: 600; color: #94a3b8;
  text-transform: uppercase; letter-spacing: .05em;
}
.dropdown-item {
  display: flex; align-items: center; gap: 8px;
  padding: 9px 16px;
  background: none; border: none; cursor: pointer;
  width: 100%; text-align: left;
  font-size: 14px; color: #374151;
  transition: background 0.15s;
  outline: none !important;
}
.dropdown-item:hover { background: #f8fafc; color: #1e293b; }
.dropdown-item.has-icon i { width: 18px; text-align: center; color: #94a3b8; }
.dropdown-item.text-danger { color: #ef4444; }
.dropdown-divider {
  height: 1px; background: #f1f5f9;
  margin: 4px 0;
}

/* ── Sidebar ───────────────────────────────────────────────────────── */
.main-sidebar {
  position: fixed;
  top: 0; left: 0; bottom: 0;
  width: 240px;
  background: #fff;
  border-right: 1px solid #e8edf2;
  z-index: 850;
  display: flex;
  flex-direction: column;
  transition: transform 0.28s ease;
  overflow-y: auto;
  padding-top: 60px;
}

.sidebar-mini .main-sidebar { width: 64px; overflow: visible; }
.sidebar-mini .main-sidebar .sidebar-brand span,
.sidebar-mini .main-sidebar .sidebar-user > div,
.sidebar-mini .main-sidebar .sidebar-menu .menu-header,
.sidebar-mini .main-sidebar .sidebar-menu li button span,
.sidebar-mini .main-sidebar .sidebar-menu li a span { display: none; }
.sidebar-mini .main-sidebar .sidebar-brand img { max-width: 36px; }

.sidebar-brand {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px 16px;
  border-bottom: 1px solid #f1f5f9;
  min-height: 62px;
}
.header-logo, .sidebar-brand img {
  max-height: 42px !important;
  max-width: 160px !important;
  width: auto !important;
  height: auto !important;
  object-fit: contain !important;
}

.sidebar-user {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
}
.sidebar-user > div { flex: 1; min-width: 0; line-height: 1.4; }
.sidebar-user strong { display: block; font-size: 13px; font-weight: 600; color: #1e293b; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.sidebar-user span { font-size: 11px; color: #94a3b8; }

/* ── Sidebar Menu & Active State (NO Black Borders) ────────────────── */
.sidebar-menu {
  list-style: none;
  margin: 0; padding: 8px 0;
  flex: 1;
}
.sidebar-menu .menu-header {
  padding: 14px 16px 6px;
  font-size: 10px; font-weight: 700;
  text-transform: uppercase; letter-spacing: .08em;
  color: #94a3b8;
}
.sidebar-menu li { position: relative; }

.sidebar-menu li a,
.sidebar-menu li button,
.sidebar-menu li button.nav-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 16px;
  border: none !important;
  border-radius: 6px;
  background: transparent !important;
  cursor: pointer;
  width: 100%;
  text-align: left;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  text-decoration: none;
  transition: background 0.15s, color 0.15s;
  outline: none !important;
  box-shadow: none !important;
  -webkit-tap-highlight-color: transparent;
}

/* Explicitly eliminate browser focus outline on buttons */
.sidebar-menu button:focus,
.sidebar-menu button:focus-visible,
.sidebar-menu button:active,
.sidebar-menu a:focus,
.sidebar-menu a:focus-visible,
.sidebar-menu a:active {
  outline: none !important;
  box-shadow: none !important;
  border: none !important;
}

/* Hover state */
.sidebar-menu li a:hover,
.sidebar-menu li button:hover {
  background-color: #f1f5f9 !important;
  color: #6777ef !important;
}
.sidebar-menu li a:hover i,
.sidebar-menu li button:hover i {
  color: #6777ef !important;
}

/* PURE COLOR INDICATION ON ACTIVE MENU ITEM */
.sidebar-menu li.active > a,
.sidebar-menu li.active > button,
.sidebar-menu li.active > button.nav-link {
  background-color: #f0f3ff !important;
  color: #6777ef !important;
  font-weight: 600 !important;
  border: none !important;
  border-left: none !important;
  outline: none !important;
  box-shadow: none !important;
}
.sidebar-menu li.active > a i,
.sidebar-menu li.active > button i {
  color: #6777ef !important;
}

/* Icons within sidebar */
.sidebar-menu li a i,
.sidebar-menu li button i,
.sidebar-menu .nav-link i {
  font-size: 16px;
  width: 20px;
  text-align: center;
  color: inherit;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.portal-nav-badge {
  margin-left: auto;
  background: #ef4444; color: #fff;
  border-radius: 20px; font-size: 10px; font-weight: 700;
  padding: 2px 7px;
}

/* Mobile sidebar */
.mobile-sidebar-open { transform: translateX(0); }
@media (max-width: 768px) {
  .main-sidebar { transform: translateX(-100%); }
  .mobile-sidebar-open { transform: translateX(0); }
}
.sidebar-scrim {
  position: fixed; inset: 0;
  background: rgba(0,0,0,.35);
  z-index: 849;
}

/* ── Main Content Area (100% Page Fit - No Horizontal Overflow) ───── */
.main-content {
  margin-top: 60px;
  margin-left: 240px;
  padding: 24px 30px 48px;
  width: calc(100% - 240px) !important;
  max-width: calc(100% - 240px) !important;
  min-height: calc(100vh - 60px);
  box-sizing: border-box !important;
  transition: margin-left 0.28s ease, width 0.28s ease, max-width 0.28s ease;
  overflow-x: hidden;
}
.sidebar-mini .main-content {
  margin-left: 64px !important;
  width: calc(100% - 64px) !important;
  max-width: calc(100% - 64px) !important;
}
@media (max-width: 768px) {
  .main-content {
    margin-left: 0 !important;
    padding: 16px 16px 32px !important;
    width: 100% !important;
    max-width: 100% !important;
  }
}

/* ── Section Structure ─────────────────────────────────────────────── */
.section {
  width: 100% !important;
  max-width: 100% !important;
  margin: 0 !important;
  padding: 0 !important;
  box-sizing: border-box !important;
}
.section-body {
  width: 100% !important;
  max-width: 100% !important;
  margin: 0 !important;
  padding: 0 !important;
  box-sizing: border-box !important;
}

/* ── Section Header (Reset Otika Negative Margins) ─────────────────── */
.section-header {
  display: flex !important;
  align-items: center !important;
  justify-content: space-between !important;
  flex-wrap: wrap !important;
  gap: 12px !important;
  width: 100% !important;
  max-width: 100% !important;
  margin: 0 0 20px 0 !important;
  padding: 0 !important;
  box-sizing: border-box !important;
  border-bottom: none !important;
}
.section-header h1 {
  font-size: 20px; font-weight: 700; color: #1e293b;
  margin: 0;
}
.section-header-breadcrumb {
  display: flex; align-items: center; gap: 6px;
  flex-wrap: wrap;
}
.breadcrumb-item {
  font-size: 13px; color: #94a3b8;
}
.breadcrumb-item.active a { color: #6777ef; text-decoration: none; font-weight: 600; }
.breadcrumb-item + .breadcrumb-item::before {
  content: '/';
  margin-right: 6px;
  color: #cbd5e1;
}

/* ── Top Action Buttons ────────────────────────────────────────────── */
.cms-actions.otika-page-actions {
  display: flex; align-items: center; gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 6px;
  padding: 8px 18px;
  border-radius: 30px; border: none; cursor: pointer;
  font-size: 13px; font-weight: 600;
  transition: all 0.15s ease;
  outline: none !important;
}
.btn-primary {
  background: #6777ef !important; color: #fff !important;
  box-shadow: 0 2px 6px #acb5f6;
}
.btn-primary:hover { background: #5a67d8 !important; }
.btn-info {
  background: #3abaf4 !important; color: #fff !important;
  box-shadow: 0 2px 6px #a8e3fe;
}
.btn-info:hover { background: #1da1f2 !important; }
.btn-success { background: #47c363 !important; color: #fff !important; }
.btn-danger { background: #fc544b !important; color: #fff !important; }
.btn-secondary { background: #e2e8f0; color: #374151; }
.btn-sm { padding: 5px 12px; font-size: 12px; border-radius: 20px; }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }

/* ── Page Heading & Hero ───────────────────────────────────────────── */
.page-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #edf2f7;
  flex-wrap: wrap;
  gap: 12px;
}
.page-heading p {
  font-size: 12px; font-weight: 700; color: #6777ef;
  text-transform: uppercase; margin: 0 0 4px; letter-spacing: 0.05em;
}
.page-heading h1 {
  font-size: 22px; font-weight: 700; color: #1e293b; margin: 0 0 4px;
}
.page-heading span { font-size: 13px; color: #64748b; }

.primary-command {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 9px 20px;
  background: #6777ef;
  color: #fff;
  border: none;
  border-radius: 30px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 2px 6px #acb5f6;
  transition: all 0.15s ease;
  outline: none !important;
}
.primary-command:hover {
  background: #5a67d8;
  transform: translateY(-1px);
}
.secondary-command {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 9px 18px;
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #cbd5e1;
  border-radius: 30px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  outline: none !important;
}
.secondary-command:hover {
  background: #e2e8f0;
  color: #1e293b;
}

/* ── Toolbar, Search Boxes & Dropdowns ─────────────────────────────── */
.list-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 20px;
  background: #fff;
  padding: 16px 20px;
  border-radius: 10px;
  border: 1px solid #edf2f7;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
  box-sizing: border-box;
  width: 100%;
}
.list-toolbar label {
  position: relative;
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 240px;
  max-width: 440px;
  margin-bottom: 0;
}
.list-toolbar label i {
  position: absolute;
  left: 14px;
  color: #94a3b8;
  font-size: 16px;
  pointer-events: none;
  display: flex;
  align-items: center;
  justify-content: center;
}
.list-toolbar input[type="search"],
.list-toolbar input[type="text"] {
  width: 100%;
  height: 42px;
  padding: 8px 16px 8px 40px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 13px;
  color: #1e293b;
  background: #f8fafc;
  transition: all 0.15s ease;
  outline: none !important;
  box-sizing: border-box;
}
.list-toolbar input[type="search"]:focus,
.list-toolbar input[type="text"]:focus {
  border-color: #6777ef;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.15) !important;
}
.list-toolbar select {
  height: 42px;
  padding: 8px 36px 8px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  color: #334155;
  background-color: #f8fafc;
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%2364748b' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
  background-position: right 12px center;
  background-repeat: no-repeat;
  background-size: 16px 16px;
  appearance: none;
  -webkit-appearance: none;
  cursor: pointer;
  outline: none !important;
  min-width: 180px;
  transition: all 0.15s ease;
  box-sizing: border-box;
}
.list-toolbar select:focus {
  border-color: #6777ef;
  background-color: #fff;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.15) !important;
}

/* ── KPI Metrics Grid ──────────────────────────────────────────────── */
.user-kpis {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 20px;
  margin-bottom: 28px;
}
.user-kpis article {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 4px 25px rgba(0, 0, 0, 0.06);
  border: 1px solid #edf2f7;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}
.user-kpis article:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
}
.user-kpis article span {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}
.user-kpis article span.blue { background: #e0f2fe; color: #0284c7; }
.user-kpis article span.amber { background: #fef3c7; color: #d97706; }
.user-kpis article span.green { background: #dcfce7; color: #16a34a; }
.user-kpis article span.cyan { background: #ecfeff; color: #0891b2; }

.user-kpis article div { display: flex; flex-direction: column; min-width: 0; }
.user-kpis article small { font-size: 12px; color: #64748b; font-weight: 600; text-transform: uppercase; letter-spacing: 0.04em; }
.user-kpis article strong { font-size: 26px; font-weight: 700; color: #1e293b; line-height: 1.2; margin: 2px 0; }
.user-kpis article p { font-size: 11px; color: #94a3b8; margin: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

@media (max-width: 1100px) {
  .user-kpis { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 575px) {
  .user-kpis { grid-template-columns: 1fr; }
}

/* ── Feed List (Activities & Notifications) ────────────────────────── */
.feed-list {
  display: grid;
  gap: 12px;
  margin-top: 16px;
}
.feed-list article {
  background: #fff;
  border-radius: 8px;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.04);
  border: 1px solid #edf2f7;
}
.feed-list article span {
  width: 40px; height: 40px;
  border-radius: 50%;
  background: #f0f3ff;
  color: #6777ef;
  display: flex; align-items: center; justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}
.feed-list article div { flex: 1; min-width: 0; }
.feed-list article strong { display: block; font-size: 14px; color: #1e293b; margin-bottom: 2px; }
.feed-list article p { font-size: 12px; color: #64748b; margin: 0 0 4px; }
.feed-list article small { font-size: 11px; color: #94a3b8; }
.feed-list article.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 20px;
  text-align: center;
}
.feed-list article.empty-state i {
  font-size: 42px;
  color: #cbd5e1;
  margin-bottom: 12px;
}

/* ── Data Tables (Clean horizontal scrolling within card) ─────────── */
.data-table {
  width: 100% !important;
  max-width: 100% !important;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 4px 25px rgba(0, 0, 0, 0.05);
  border: 1px solid #edf2f7;
  overflow-x: auto !important;
  -webkit-overflow-scrolling: touch;
  box-sizing: border-box;
  margin-top: 16px;
}
.data-table table {
  width: 100%;
  min-width: 680px;
  border-collapse: collapse;
  text-align: left;
}
.data-table th {
  padding: 14px 20px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #64748b;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  white-space: nowrap;
}
.data-table td {
  padding: 16px 20px;
  font-size: 13px;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}
.data-table tr:last-child td { border-bottom: none; }
.data-table tr:hover td { background: #fbfcfe; }

/* ── Table Action Buttons ──────────────────────────────────────────── */
.table-actions {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.table-actions button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  color: #475569;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.15s ease;
  outline: none !important;
}
.table-actions button:hover {
  background: #6777ef;
  border-color: #6777ef;
  color: #fff;
  box-shadow: 0 2px 6px rgba(103, 119, 239, 0.3);
}

/* ── Status Badges ─────────────────────────────────────────────────── */
.status {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}
.status.green { background: #dcfce7; color: #16a34a; }
.status.amber { background: #fef3c7; color: #d97706; }
.status.red { background: #fee2e2; color: #dc2626; }
.status.blue { background: #e0f2fe; color: #0284c7; }

/* ── Profile & Settings Layouts ────────────────────────────────────── */
.profile-layout {
  display: grid;
  gap: 24px;
  margin-top: 16px;
}
.profile-grid-3 {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 24px;
}
@media (max-width: 900px) {
  .profile-grid-3 { grid-template-columns: 1fr; }
}
.profile-grid-2 {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 24px;
}
@media (max-width: 900px) {
  .profile-grid-2 { grid-template-columns: 1fr; }
}

.profile-summary {
  background: #fff;
  border-radius: 10px;
  padding: 24px;
  text-align: center;
  box-shadow: 0 4px 25px rgba(0,0,0,0.05);
  border: 1px solid #edf2f7;
}
.profile-avatar {
  margin: 0 auto 16px;
  width: 96px; height: 96px;
  border-radius: 50%;
  overflow: hidden;
  border: 3px solid #6777ef;
  display: flex; align-items: center; justify-content: center;
  background: #f0f3ff;
  font-size: 32px; font-weight: 700; color: #6777ef;
}
.profile-avatar img { width: 100%; height: 100%; object-fit: cover; }

.profile-form {
  background: #fff;
  border-radius: 10px;
  padding: 24px;
  box-shadow: 0 4px 25px rgba(0,0,0,0.05);
  border: 1px solid #edf2f7;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.profile-form h2 {
  font-size: 16px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 4px;
}
.profile-form label {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #475569;
  margin-bottom: 0;
}
.profile-form input,
.profile-form select,
.profile-form textarea {
  width: 100%;
  height: 42px;
  padding: 8px 14px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  font-size: 13px;
  color: #1e293b;
  background: #f8fafc;
  transition: all 0.15s ease;
  outline: none !important;
  box-sizing: border-box;
}
.profile-form input:focus,
.profile-form select:focus,
.profile-form textarea:focus {
  border-color: #6777ef;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(103, 119, 239, 0.15) !important;
}
.profile-form input:disabled {
  background: #f1f5f9;
  color: #94a3b8;
  cursor: not-allowed;
}
.profile-form div:has(> label + label) {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
@media (max-width: 575px) {
  .profile-form div:has(> label + label) { grid-template-columns: 1fr; }
}

/* ── Card ──────────────────────────────────────────────────────────── */
.card {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 4px 25px rgba(0,0,0,0.05);
  border: 1px solid #edf2f7;
  overflow: hidden;
}
.card-body { padding: 20px; }

/* ── Messages & Alerts ─────────────────────────────────────────────── */
.portal-success, .portal-error {
  padding: 12px 18px;
  border-radius: 8px; margin-bottom: 20px;
  font-size: 13px; font-weight: 500;
}
.portal-success { background: #d1fae5; color: #065f46; border: 1px solid #6ee7b7; }
.portal-error { background: #fee2e2; color: #991b1b; border: 1px solid #fca5a5; }

/* ── Footer ────────────────────────────────────────────────────────── */
footer.main-footer, .main-footer {
  text-align: center;
  font-size: 12px; color: #94a3b8;
  padding: 16px 24px;
  border-top: 1px solid #e8edf2;
  background: #fff;
  width: calc(100% - 240px) !important;
  max-width: calc(100% - 240px) !important;
  margin-left: 240px !important;
  box-sizing: border-box !important;
  transition: margin-left 0.28s, width 0.28s, max-width 0.28s;
}
.sidebar-mini footer.main-footer {
  margin-left: 64px !important;
  width: calc(100% - 64px) !important;
  max-width: calc(100% - 64px) !important;
}
@media (max-width: 768px) {
  footer.main-footer, .main-footer {
    margin-left: 0 !important;
    width: 100% !important;
    max-width: 100% !important;
  }
}

/* ── Dark Mode Overrides ───────────────────────────────────────────── */
:global(.dark) .user-portal { background: #0f172a; color: #e2e8f0; }
:global(.dark) .navbar-bg,
:global(.dark) .main-navbar { background: #1e293b; box-shadow: 0 1px 4px rgba(0,0,0,.3); }
:global(.dark) .main-sidebar { background: #1e293b; border-right-color: #334155; }
:global(.dark) .sidebar-brand,
:global(.dark) .sidebar-user { border-bottom-color: #334155; }
:global(.dark) .sidebar-menu li a,
:global(.dark) .sidebar-menu li button.nav-link { color: #cbd5e1; }
:global(.dark) .sidebar-menu li.active > a,
:global(.dark) .sidebar-menu li.active > button,
:global(.dark) .sidebar-menu li.active > button.nav-link {
  background-color: #1e293b !important;
  color: #93c5fd !important;
  border: none !important;
  border-left: none !important;
  outline: none !important;
  box-shadow: none !important;
}
:global(.dark) .sidebar-menu li.active > a i,
:global(.dark) .sidebar-menu li.active > button i {
  color: #93c5fd !important;
}
:global(.dark) .sidebar-menu li a:hover,
:global(.dark) .sidebar-menu li button:hover {
  background-color: #334155 !important;
  color: #f8fafc !important;
}
:global(.dark) .card,
:global(.dark) .user-kpis article,
:global(.dark) .feed-list article,
:global(.dark) .data-table,
:global(.dark) .profile-summary,
:global(.dark) .dropdown-menu {
  background: #1e293b;
  border-color: #334155;
  box-shadow: 0 4px 25px rgba(0,0,0,.3);
}
:global(.dark) .list-toolbar {
  background: #1e293b;
  border-color: #334155;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.25);
}
:global(.dark) .list-toolbar input[type="search"],
:global(.dark) .list-toolbar input[type="text"],
:global(.dark) .list-toolbar select {
  background-color: #0f172a;
  border-color: #334155;
  color: #f8fafc;
}
:global(.dark) .list-toolbar select {
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%2394a3b8' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
}
:global(.dark) .table-actions button {
  background: #0f172a;
  border-color: #334155;
  color: #cbd5e1;
}
:global(.dark) .profile-form {
  background: #1e293b;
  border-color: #334155;
}
:global(.dark) .profile-form h2 { color: #f8fafc; }
:global(.dark) .profile-form label { color: #cbd5e1; }
:global(.dark) .profile-form input,
:global(.dark) .profile-form select,
:global(.dark) .profile-form textarea {
  background-color: #0f172a;
  border-color: #334155;
  color: #f8fafc;
}
:global(.dark) .user-kpis article strong,
:global(.dark) .section-header h1,
:global(.dark) .page-heading h1 { color: #f8fafc; }
:global(.dark) .data-table th { background: #1e293b; color: #94a3b8; }
:global(.dark) .data-table td { color: #cbd5e1; border-bottom-color: #334155; }
:global(.dark) footer.main-footer { background: #1e293b; border-top-color: #334155; }
:global(.dark) .main-navbar .nav-link-lg.cms-top-icon,
:global(.dark) .main-navbar button.cms-top-icon,
:global(.dark) .main-navbar .portal-top-action,
:global(.dark) .main-navbar .portal-theme-toggle {
  background: #0f172a !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.4) !important;
}
:global(.dark) .main-navbar .nav-link-lg.cms-top-icon i,
:global(.dark) .main-navbar button.cms-top-icon i,
:global(.dark) .main-navbar .portal-top-action i,
:global(.dark) .main-navbar .portal-theme-toggle i {
  color: #f8fafc !important;
}
:global(.dark) .main-navbar .nav-link-lg.cms-top-icon:hover,
:global(.dark) .main-navbar button.cms-top-icon:hover,
:global(.dark) .main-navbar .portal-top-action:hover,
:global(.dark) .main-navbar .portal-theme-toggle:hover {
  background: #1e293b !important;
  border-color: #6777ef !important;
  color: #93c5fd !important;
}
:global(.dark) .main-navbar .nav-link-lg.cms-top-icon:hover i,
:global(.dark) .main-navbar button.cms-top-icon:hover i,
:global(.dark) .main-navbar .portal-top-action:hover i,
:global(.dark) .main-navbar .portal-theme-toggle:hover i {
  color: #93c5fd !important;
}
:global(.dark) .headerBadge1,
:global(.dark) .headerBadge2 {
  border-color: #1e293b !important;
}
:global(.dark) .portal-navbar-title strong {
  color: #f8fafc !important;
}

/* ── My Activities & Audit Trail Styling ───────────────────────────── */
.activity-kpis {
  margin-bottom: 24px;
}
.activity-toolbar {
  margin-bottom: 24px;
  gap: 16px;
}
.activity-category-pills {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.category-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 30px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  color: #475569;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  outline: none !important;
}
.category-pill:hover {
  background: #f1f5f9;
  color: #1e293b;
}
.category-pill.active {
  background: #6777ef;
  color: #fff;
  border-color: #6777ef;
  box-shadow: 0 2px 6px rgba(103, 119, 239, 0.35);
}

.activity-timeline-feed {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.activity-feed-card {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #edf2f7;
  padding: 18px 22px;
  display: flex;
  align-items: flex-start;
  gap: 18px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}
.activity-feed-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.06);
}

.activity-icon-wrapper {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}
.activity-icon-wrapper.blue { background: #e0f2fe; color: #0284c7; }
.activity-icon-wrapper.green { background: #dcfce7; color: #16a34a; }
.activity-icon-wrapper.amber { background: #fef3c7; color: #d97706; }
.activity-icon-wrapper.purple { background: #f3e8ff; color: #9333ea; }
.activity-icon-wrapper.cyan { background: #ecfeff; color: #0891b2; }

.activity-content {
  flex: 1;
  min-width: 0;
}
.activity-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 6px;
}
.activity-title-group {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.activity-title {
  font-size: 15px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}
.activity-category-badge {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  padding: 3px 8px;
  border-radius: 6px;
  letter-spacing: 0.04em;
}
.activity-category-badge.blue { background: #e0f2fe; color: #0284c7; }
.activity-category-badge.green { background: #dcfce7; color: #16a34a; }
.activity-category-badge.amber { background: #fef3c7; color: #d97706; }
.activity-category-badge.purple { background: #f3e8ff; color: #9333ea; }
.activity-category-badge.cyan { background: #ecfeff; color: #0891b2; }

.activity-time {
  font-size: 12px;
  color: #94a3b8;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.activity-description {
  font-size: 13px;
  color: #475569;
  margin: 0 0 10px 0;
  line-height: 1.5;
}
.activity-meta-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.meta-pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  padding: 4px 10px;
  font-size: 11px;
  color: #64748b;
}
.meta-pill i { font-size: 12px; }
.meta-pill strong { color: #334155; }
.date-pill { background: #f1f5f9; color: #475569; }

.activity-pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
  background: #fff;
  padding: 14px 20px;
  border-radius: 10px;
  border: 1px solid #edf2f7;
}
.pagination-info { font-size: 13px; color: #64748b; }
.pagination-buttons { display: flex; gap: 8px; }

/* ── Settings Sub-Navigation & Layout ──────────────────────────────── */
.settings-nav-tabs {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 24px;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 12px;
}
.settings-tab-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 30px;
  border: 1px solid transparent;
  background: transparent;
  color: #64748b;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  outline: none !important;
}
.settings-tab-btn:hover {
  background: #f1f5f9;
  color: #1e293b;
}
.settings-tab-btn.active {
  background: #6777ef;
  color: #fff;
  box-shadow: 0 3px 10px rgba(103, 119, 239, 0.35);
}

.settings-card {
  background: #fff;
  border-radius: 12px;
  padding: 26px;
  border: 1px solid #edf2f7;
  box-shadow: 0 4px 25px rgba(0, 0, 0, 0.04);
}
.settings-card-header {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  margin-bottom: 22px;
  padding-bottom: 16px;
  border-bottom: 1px solid #edf2f7;
}
.settings-card-header h3 {
  font-size: 17px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 4px;
}
.settings-card-header p {
  font-size: 12px;
  color: #64748b;
  margin: 0;
}

.password-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}
.password-input-wrap input {
  padding-right: 42px !important;
}
.btn-toggle-eye {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  font-size: 16px;
  outline: none !important;
  display: flex;
  align-items: center;
  justify-content: center;
}
.btn-toggle-eye:hover { color: #6777ef; }

.strength-bar-track {
  width: 100%;
  height: 6px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
}
.strength-bar-fill {
  height: 100%;
  transition: width 0.3s ease, background-color 0.3s ease;
}
.strength-weak { background-color: #ef4444; }
.strength-fair { background-color: #f59e0b; }
.strength-good { background-color: #06b6d4; }
.strength-strong { background-color: #10b981; }

.password-rules-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 6px;
}
@media (max-width: 575px) {
  .password-rules-list { grid-template-columns: 1fr; }
}
.password-rules-list li {
  font-size: 11px;
  color: #94a3b8;
  display: flex;
  align-items: center;
  gap: 6px;
}
.password-rules-list li.met {
  color: #16a34a;
  font-weight: 600;
}

.security-tips-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.security-tips-list li {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}
.security-tips-list strong {
  display: block;
  font-size: 13px;
  color: #1e293b;
  margin-bottom: 2px;
}
.security-tips-list p {
  font-size: 12px;
  color: #64748b;
  margin: 0;
  line-height: 1.4;
}

/* ── Sessions List ─────────────────────────────────────────────────── */
.sessions-list-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.session-item-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  transition: all 0.15s ease;
}
.session-item-card.current-session-card {
  background: #fff;
  border-color: #86efac;
  box-shadow: 0 2px 10px rgba(34, 197, 94, 0.08);
}
.session-device-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: #e2e8f0;
  color: #475569;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}
.session-device-icon.current {
  background: #dcfce7;
  color: #16a34a;
}
.session-details { flex: 1; min-width: 0; }
.session-device-name {
  font-size: 14px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}
.current-badge {
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  padding: 3px 8px;
}
.session-meta-items {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
  margin-top: 4px;
}
.session-meta-items .meta-item {
  font-size: 12px;
  color: #64748b;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

/* ── Preferences & 2FA ─────────────────────────────────────────────── */
.twofa-status-banner.twofa-active {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
}
.twofa-status-banner.twofa-inactive {
  background: #fffbeb;
  border: 1px solid #fef08a;
}
.preferences-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.preference-item {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  cursor: pointer;
  padding: 12px 14px;
  border-radius: 8px;
  border: 1px solid #edf2f7;
  background: #f8fafc;
  transition: all 0.15s ease;
}
.preference-item:hover {
  background: #fff;
  border-color: #cbd5e1;
}
.pref-checkbox {
  width: 18px;
  height: 18px;
  margin-top: 2px;
  accent-color: #6777ef;
  flex-shrink: 0;
  cursor: pointer;
}
.preference-text strong {
  display: block;
  font-size: 13px;
  color: #1e293b;
  margin-bottom: 2px;
}
.preference-text p {
  font-size: 12px;
  color: #64748b;
  margin: 0;
  line-height: 1.4;
}

/* ── Dark Mode Extensions for Activities & Settings ───────────────── */
:global(.dark) .category-pill {
  background: #1e293b;
  border-color: #334155;
  color: #cbd5e1;
}
:global(.dark) .category-pill:hover {
  background: #334155;
  color: #f8fafc;
}
:global(.dark) .category-pill.active {
  background: #6777ef;
  color: #fff;
}
:global(.dark) .activity-feed-card,
:global(.dark) .settings-card,
:global(.dark) .activity-pagination-bar {
  background: #1e293b;
  border-color: #334155;
  box-shadow: 0 4px 25px rgba(0, 0, 0, 0.25);
}
:global(.dark) .activity-title,
:global(.dark) .settings-card-header h3,
:global(.dark) .session-device-name,
:global(.dark) .security-tips-list strong,
:global(.dark) .preference-text strong {
  color: #f8fafc;
}
:global(.dark) .activity-description,
:global(.dark) .settings-card-header p,
:global(.dark) .security-tips-list p,
:global(.dark) .preference-text p {
  color: #94a3b8;
}
:global(.dark) .meta-pill,
:global(.dark) .session-item-card {
  background: #0f172a;
  border-color: #334155;
  color: #cbd5e1;
}
:global(.dark) .session-item-card.current-session-card {
  background: #1e293b;
  border-color: #16a34a;
}
:global(.dark) .preference-item {
  background: #0f172a;
  border-color: #334155;
}
:global(.dark) .preference-item:hover {
  background: #1e293b;
}

</style>

