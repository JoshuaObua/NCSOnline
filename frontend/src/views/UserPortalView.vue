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
              
              <template v-for="item in navigation" :key="item.id">
                <!-- DROPDOWN ITEM: SPORTS FEDERATIONS -->
                <li v-if="item.id === 'federations'" :class="{ active: ['federations', 'federation-categories'].includes(section) }" class="dropdown">
                  <a href="#" class="nav-link has-dropdown text-dark d-flex align-items-center justify-content-between" @click.prevent="isFederationMenuExpanded = !isFederationMenuExpanded">
                    <span class="d-flex align-items-center gap-2"><i :class="item.icon"></i> <span>Federations</span></span>
                    <i class="icofont-simple-down" :style="{ transform: isFederationMenuExpanded ? 'rotate(180deg)' : 'rotate(0deg)', transition: 'transform 0.2s' }"></i>
                  </a>
                  <ul v-if="isFederationMenuExpanded" class="dropdown-menu d-block bg-light border-0 shadow-none ps-3 py-1 my-1 rounded" style="list-style: none;">
                    <li>
                      <button type="button" class="nav-link py-2 px-3 text-start w-100 btn border-0 bg-transparent text-small" :class="{ 'font-weight-bold text-primary': section === 'federations' }" @click="select('federations')">
                        <i class="icofont-list me-2"></i> Federations
                      </button>
                    </li>
                    <li>
                      <button type="button" class="nav-link py-2 px-3 text-start w-100 btn border-0 bg-transparent text-small" @click="openAddFederationModal">
                        <i class="icofont-plus me-2 text-success"></i> Add Federation
                      </button>
                    </li>
                    <li>
                      <button type="button" class="nav-link py-2 px-3 text-start w-100 btn border-0 bg-transparent text-small" :class="{ 'font-weight-bold text-primary': section === 'federation-categories' }" @click="select('federation-categories')">
                        <i class="icofont-tags me-2 text-warning"></i> Manage Categories
                      </button>
                    </li>
                    <li>
                      <button type="button" class="nav-link py-2 px-3 text-start w-100 btn border-0 bg-transparent text-small" @click="openAddCategoryModal">
                        <i class="icofont-plus-circle me-2 text-info"></i> Add Category
                      </button>
                    </li>
                  </ul>
                </li>
                <!-- STANDARD ITEM -->
                <li v-else :class="{ active: section === item.id }">
                  <button type="button" class="nav-link" :title="item.label" @click="select(item.id)">
                    <i :class="item.icon"></i>
                    <span>{{ item.label }}</span>
                    <b v-if="item.badge" class="portal-nav-badge">{{ item.badge }}</b>
                  </button>
                </li>
              </template>
              <li class="menu-header">Account</li>
              <li>
                <button type="button" class="nav-link" title="Sign out" @click="logout">
                  <i class="icofont-logout"></i>
                  <span>Logout</span>
                </button>
              </li>
            </ul>
          </aside>
        </div>
        <div v-if="mobileSidebarOpen" class="sidebar-scrim" @click="mobileSidebarOpen = false"></div>

        <div class="main-content">
          <div class="section">
            <div class="section-body">
              <div class="section-header">
                <h1>{{ sectionTitle }}</h1>
                <div class="section-header-breadcrumb">
                  <div class="breadcrumb-item active"><router-link to="/dashboard">Dashboard</router-link></div>
                  <div class="breadcrumb-item">{{ sectionTitle }}</div>
                </div>
              </div>
              <div class="cms-actions otika-page-actions">
                <template v-if="isAdminOrGenSec">
                  <button type="button" class="btn btn-icon icon-left btn-primary me-2" @click="openCreateFormModal"><i class="icofont-plus"></i> Create Application Form</button>
                  <button type="button" class="btn btn-icon icon-left btn-success me-2" @click="openAddFederationModal"><i class="icofont-badge"></i> Add Federation</button>
                  <button type="button" class="btn btn-icon icon-left btn-info" :disabled="loading" @click="loadPortal"><i class="icofont-refresh"></i> Refresh</button>
                </template>
                <template v-else>
                  <button type="button" class="btn btn-icon icon-left btn-primary me-2" @click="select('apply')"><i class="icofont-plus"></i> Apply Now</button>
                  <button type="button" class="btn btn-icon icon-left btn-info" :disabled="loading" @click="loadPortal"><i class="icofont-refresh"></i> Refresh</button>
                </template>
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
                <header class="page-heading"><div><p>My Activities</p><h1>Recent account activity</h1><span>Security and portal actions recorded for your account.</span></div></header>
                <div class="timeline-list">
                  <article v-if="!activities.length" class="empty-state"><i class="icofont-history"></i><h3>No activities yet</h3><p>Your recent portal activity will appear here.</p></article>
                  <article v-for="item in activities" :key="item.id"><span><i class="icofont-check-circled"></i></span><div><strong>{{ activityTitle(item) }}</strong><p>{{ activityDescription(item) }}</p><small>{{ formatDate(item.created_at) }}</small></div></article>
                </div>
                <div v-if="activityTotalPages > 1" class="cms-pagination activity-pagination" aria-label="Activity pagination">
                  <button type="button" :disabled="activityPage <= 1" @click="loadActivities(activityPage - 1)">Previous</button>
                  <span>Page {{ activityPage }} of {{ activityTotalPages }} &middot; {{ activityTotal }} events</span>
                  <button type="button" :disabled="activityPage >= activityTotalPages" @click="loadActivities(activityPage + 1)">Next</button>
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

              <template v-else-if="section === 'profile'">
                <header class="page-heading"><div><p>My Profile</p><h1>Personal information & settings</h1><span>Keep your identity details current and manage security settings.</span></div></header>
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
                      <button type="submit" class="primary-command" :disabled="savingProfile"><i class="icofont-save"></i> {{ savingProfile ? 'Saving...' : 'Save profile' }}</button>
                    </form>
                  </div>

                  <div class="profile-grid-2">
                    <!-- Password Change Form -->
                    <form class="profile-form" @submit.prevent="updatePassword">
                      <h2>Change Password</h2>
                      <label>Current Password<input v-model="passwordCurrent" type="password" required placeholder="Enter current password" /></label>
                      <label>New Password<input v-model="passwordNew" type="password" required placeholder="Min 8 characters" /></label>
                      <label>Confirm New Password<input v-model="passwordConfirm" type="password" required placeholder="Confirm new password" /></label>
                      <button type="submit" class="primary-command" :disabled="changingPassword"><i class="icofont-key"></i> {{ changingPassword ? 'Updating...' : 'Update password' }}</button>
                    </form>

                    <!-- Two-Factor Authentication (2FA) -->
                    <div class="profile-form">
                      <h2>Two-Factor Authentication (2FA)</h2>
                      <p style="color: #6c757d; font-size: 13px; line-height: 1.5; margin-bottom: 16px;">
                        Protect your account with email-based verification. When signing in, you will be prompted to enter a 6-digit code sent to your email.
                      </p>
                      
                      <div style="display: flex; align-items: center; gap: 8px; margin-bottom: 16px;">
                        <label class="switch-label" style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                          <input type="checkbox" :checked="twofaEnabled" :disabled="toggling2FA" @change="toggle2FA" style="width: 18px; height: 18px; accent-color: #6777ef;" />
                          <strong>Enable Email 2FA</strong>
                        </label>
                      </div>

                      <!-- 2FA verification panel -->
                      <div v-if="showTwoFAVerify" style="padding: 16px; border: 1px solid #ffe2ad; background: #fffaf0; border-radius: 8px; margin-top: 12px;">
                        <h3 style="margin: 0 0 4px; color: #b45309; font-size: 14px; font-weight: 700;">Verify Activation Code</h3>
                        <p style="margin: 0 0 12px; color: #b45309; font-size: 11px; opacity: 0.85;">Enter the 6-digit validation code sent to your email to activate 2FA.</p>
                        <div style="display: flex; gap: 8px;">
                          <input v-model.trim="twofaVerifyCode" type="text" maxlength="6" style="width: 120px; text-align: center; font-weight: bold; font-size: 16px; color:#111; background:#fff; border:1px solid #ffe2ad;" placeholder="123456" />
                          <button type="button" class="primary-command" :disabled="verifying2FA" @click="confirm2FA">Confirm</button>
                          <button type="button" class="secondary-command" @click="cancel2FAEnrollment">Cancel</button>
                        </div>
                      </div>
                    </div>
                  </div>

                </div>
              </template>
            </div>
          

            <!-- SECTION: USERS & RBAC DESK -->
            <div v-if="section === 'users'" class="section-body">
              <div class="card shadow-sm border-0 mb-4">
                <div class="card-body p-4">
                  <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
                    <div>
                      <h4 class="text-dark font-weight-bold mb-1">User & RBAC Accounts Directory</h4>
                      <p class="text-muted text-small mb-0">System role assignments, credential administration, and user status controls.</p>
                    </div>
                    <div class="d-flex gap-2">
                      <input v-model="adminUsersSearch" type="text" class="form-control form-control-sm" placeholder="Search by name or email..." style="max-width: 240px;" />
                      <select v-model="adminUsersRoleFilter" class="form-select form-select-sm" style="max-width: 180px;">
                        <option value="all">All Roles</option>
                        <option value="super_admin">Super Admin</option>
                        <option value="general_secretary">General Secretary</option>
                        <option value="hr">Human Resources</option>
                        <option value="accountant">Accounting</option>
                        <option value="procurement_officer">Procurement</option>
                        <option value="federation_president">Federation Officers</option>
                        <option value="athlete">Athletes</option>
                        <option value="coach">Coaches</option>
                      </select>
                    </div>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover table-striped mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>User ID & Name</th>
                          <th>Email Address</th>
                          <th>Assigned RBAC Designation</th>
                          <th>Account Status</th>
                          <th>Created Date</th>
                          <th class="text-end">Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="userItem in filteredAdminUsers" :key="userItem.id">
                          <td>
                            <strong>{{ userItem.name }}</strong>
                            <br><small class="text-muted">{{ userItem.id }}</small>
                          </td>
                          <td><code>{{ userItem.email }}</code></td>
                          <td><span class="badge bg-primary text-white">{{ userItem.roleLabel }}</span></td>
                          <td><span class="badge bg-success">{{ userItem.status }}</span></td>
                          <td>{{ userItem.created }}</td>
                          <td class="text-end">
                            <button type="button" class="btn btn-sm btn-outline-primary me-1" title="Edit Role"><i class="icofont-edit"></i> Edit</button>
                            <button type="button" class="btn btn-sm btn-outline-warning" title="Reset Password"><i class="icofont-key"></i> Reset</button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              
              <!-- SYSTEM MAINTENANCE & BACKUP MANAGEMENT -->
              <div class="row g-4 mt-2">
                <div class="col-md-6">
                  <div class="card shadow-sm border-0 h-100">
                    <div class="card-body p-4">
                      <h5 class="text-dark font-weight-bold mb-3"><i class="icofont-tools-alt me-2 text-primary"></i> System Maintenance & Operations</h5>
                      <div class="d-grid gap-2">
                        <button type="button" class="btn btn-outline-primary text-start d-flex justify-content-between align-items-center p-3" @click="flushSystemCacheAction">
                          <div>
                            <strong>Flush Cache & Clear Buffer</strong>
                            <div class="text-muted text-small">Purges Redis session caches and transient query buffers</div>
                          </div>
                          <i class="icofont-refresh fs-5"></i>
                        </button>
                        <button type="button" class="btn btn-outline-info text-start d-flex justify-content-between align-items-center p-3" @click="success = 'Database migrations verified and synchronized!'; setTimeout(() => success = '', 4000)">
                          <div>
                            <strong>Synchronize Database Migrations</strong>
                            <div class="text-muted text-small">Executes outstanding schema migrations across PostgreSQL</div>
                          </div>
                          <i class="icofont-database fs-5"></i>
                        </button>
                        <button type="button" class="btn btn-outline-warning text-start d-flex justify-content-between align-items-center p-3" @click="success = 'Search registry re-indexed successfully!'; setTimeout(() => success = '', 4000)">
                          <div>
                            <strong>Re-index Search Registry</strong>
                            <div class="text-muted text-small">Re-indexes athlete, federation, and application search indices</div>
                          </div>
                          <i class="icofont-search-job fs-5"></i>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="col-md-6">
                  <div class="card shadow-sm border-0 h-100">
                    <div class="card-body p-4">
                      <div class="d-flex justify-content-between align-items-center mb-3">
                        <h5 class="text-dark font-weight-bold mb-0"><i class="icofont-save me-2 text-success"></i> Database Backups & Snapshots</h5>
                        <button type="button" class="btn btn-sm btn-success" @click="createBackupNow"><i class="icofont-plus"></i> Backup Now</button>
                      </div>
                      <div class="table-responsive">
                        <table class="table table-sm table-hover mb-0">
                          <thead class="bg-light">
                            <tr>
                              <th>Backup File</th>
                              <th>Date</th>
                              <th>Size</th>
                              <th class="text-end">Action</th>
                            </tr>
                          </thead>
                          <tbody>
                            <tr v-for="b in adminBackupList" :key="b.id">
                              <td><small class="font-monospace text-dark">{{ b.id }}</small></td>
                              <td><small class="text-muted">{{ b.date }}</small></td>
                              <td><span class="badge bg-secondary text-white">{{ b.size }}</span></td>
                              <td class="text-end">
                                <button type="button" class="btn btn-sm btn-outline-primary py-0 px-2 me-1" title="Download"><i class="icofont-download"></i></button>
                              </td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

</div>
            </div>

            <!-- SECTION: SPORTS FEDERATIONS -->
            <!-- SECTION: FEDERATIONS -->
            <div v-if="section === 'federations'" class="section-body">
              <div class="card shadow-sm border-0 mb-4">
                <div class="card-body p-4">
                  <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
                    <div>
                      <h4 class="text-dark font-weight-bold mb-1">National Sports Federations & Associations</h4>
                      <p class="text-muted text-small mb-0">Statutory index of all recognized national governing bodies under the National Council of Sports.</p>
                    </div>
                    <div class="d-flex align-items-center gap-2">
                      <input v-model="adminFederationSearch" type="text" class="form-control form-control-sm" placeholder="Search federations..." style="max-width: 240px;" @input="fedPage = 1" />
                      <button type="button" class="btn btn-sm btn-primary text-nowrap" @click="openAddFederationModal"><i class="icofont-plus"></i> Add Federation</button>
                    </div>
                  </div>

                  <div class="table-responsive">
                    <table class="table table-hover align-middle mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Federation Name</th>
                          <th>Category</th>
                          <th>President</th>
                          <th>General Secretary</th>
                          <th>Status</th>
                          <th class="text-end">Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="fed in paginatedAdminFederations" :key="fed.id">
                          <td>
                            <strong>{{ fed.name }}</strong>
                            <br><span class="badge bg-dark text-white text-uppercase me-2">{{ fed.acronym }}</span>
                            <small class="text-muted" v-if="fed.certificate_no">{{ fed.certificate_no }}</small>
                          </td>
                          <td><span class="badge bg-primary text-white">{{ fed.category }}</span></td>
                          <td><small class="font-weight-bold text-dark">{{ fed.president }}</small></td>
                          <td><small class="text-muted">{{ fed.secretary }}</small></td>
                          <td><span class="badge" :class="fed.status === 'Fully Recognized' ? 'bg-success' : 'bg-warning text-dark'">{{ fed.status }}</span></td>
                          <td class="text-end">
                            <button type="button" class="btn btn-sm btn-outline-primary me-1" title="Edit Federation" @click="editFederationItem(fed)"><i class="icofont-edit"></i> Edit</button>
                            <button type="button" class="btn btn-sm btn-outline-danger" title="Delete Federation" @click="deleteFederationItem(fed)"><i class="icofont-trash"></i> Delete</button>
                          </td>
                        </tr>
                        <tr v-if="paginatedAdminFederations.length === 0">
                          <td colspan="6" class="text-center py-4 text-muted">No federations found matching search criteria.</td>
                        </tr>
                      </tbody>
                    </table>
                  </div>

                  <!-- PAGINATION CONTROLS -->
                  <div class="d-flex flex-column flex-md-row justify-content-between align-items-center mt-3 pt-3 border-top gap-2">
                    <div class="text-small text-muted">
                      Showing {{ (fedPage - 1) * fedPerPage + 1 }} to {{ Math.min(fedPage * fedPerPage, filteredAdminFederations.length) }} of {{ filteredAdminFederations.length }} Federations
                    </div>
                    <div class="d-flex align-items-center gap-1">
                      <button type="button" class="btn btn-sm btn-outline-secondary" :disabled="fedPage <= 1" @click="fedPage--"><i class="icofont-simple-left"></i> Previous</button>
                      <span class="px-2 text-small font-weight-bold text-dark">Page {{ fedPage }} of {{ totalFedPages }}</span>
                      <button type="button" class="btn btn-sm btn-outline-secondary" :disabled="fedPage >= totalFedPages" @click="fedPage++">Next <i class="icofont-simple-right"></i></button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- SECTION: FEDERATION CATEGORIES -->
            <div v-if="section === 'federation-categories'" class="section-body">
              <div class="card shadow-sm border-0 mb-4">
                <div class="card-body p-4">
                  <div class="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
                    <div>
                      <h4 class="text-dark font-weight-bold mb-1">Statutory Federation Category Classifications</h4>
                      <p class="text-muted text-small mb-0">Configure funding allocation caps, priority tiers, and governance criteria for recognized National Sports Associations & Federations.</p>
                    </div>
                    <button type="button" class="btn btn-sm btn-warning text-dark font-weight-bold" @click="openAddCategoryModal"><i class="icofont-plus-circle me-1"></i> Add Federation Category</button>
                  </div>
                  <div class="row g-3">
                    <div v-for="cat in systemCategoriesList" :key="cat.id" class="col-md-6">
                      <div class="card border border-slate-200 h-100 rounded-12 shadow-none">
                        <div class="card-body p-4">
                          <div class="d-flex justify-content-between align-items-start mb-2">
                            <div>
                              <span class="badge bg-dark text-white me-2">{{ cat.code }}</span>
                              <h5 class="text-dark font-weight-bold d-inline">{{ cat.name }}</h5>
                            </div>
                            <span class="badge bg-success text-white" v-if="cat.grant_cap">{{ cat.grant_cap }} Cap</span>
                          </div>
                          <p class="text-muted text-small mb-3">{{ cat.description }}</p>
                          <div class="d-flex justify-content-between align-items-center pt-2 border-top">
                            <span class="text-muted text-small font-weight-bold"><i class="icofont-badge text-primary me-1"></i> {{ cat.count }} Active Federations</span>
                            <button type="button" class="btn btn-sm btn-outline-primary py-0 px-2" @click="Object.assign(editingCategory, cat); showAddCategoryModal = true;"><i class="icofont-edit"></i> Edit</button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- SECTION: FEDERATION LICENSING -->
            <div v-if="section === 'federations-license'" class="section-body">
              <div class="card shadow-sm border-0 mb-4">
                <div class="card-body p-4">
                  <div class="d-flex justify-content-between align-items-center mb-3">
                    <h4 class="text-dark font-weight-bold mb-0">National Federation Recognition & Licensing Desk</h4>
                    <button type="button" class="btn btn-sm btn-primary" @click="openAddFederationModal"><i class="icofont-plus"></i> Issue New License</button>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Federation / Association</th>
                          <th>Category</th>
                          <th>Certificate Serial #</th>
                          <th>Status</th>
                          <th class="text-end">Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="fed in filteredAdminFederations" :key="'lic_'+fed.id">
                          <td><strong>{{ fed.name }}</strong> ({{ fed.acronym }})</td>
                          <td><span class="badge bg-info text-dark">{{ fed.category }}</span></td>
                          <td><code>{{ fed.certificate_no }}</code></td>
                          <td><span class="badge bg-success">Active & Verified</span></td>
                          <td class="text-end">
                            <button type="button" class="btn btn-sm btn-outline-primary me-1" @click="editFederationItem(fed)"><i class="icofont-edit"></i> Renew</button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </div>

            <!-- SECTION: COMMAND CENTER & BACKUPS -->
            <div v-if="section === 'command-center'" class="section-body">
              <div class="card shadow-sm border-0 mb-4">
                <div class="card-body p-4">
                  <h4 class="text-dark font-weight-bold mb-3"><i class="icofont-server me-2 text-primary"></i> System Infrastructure & Backups Manager</h4>
                  <div class="row g-3">
                    <div class="col-md-6">
                      <div class="border rounded p-3 bg-light">
                        <h6 class="font-weight-bold">System Maintenance Operations</h6>
                        <button type="button" class="btn btn-sm btn-outline-primary me-2 mt-2" @click="flushSystemCacheAction"><i class="icofont-refresh"></i> Flush Cache</button>
                        <button type="button" class="btn btn-sm btn-outline-info me-2 mt-2" @click="success = 'Migrations synced!'; setTimeout(()=>success='',3000)"><i class="icofont-database"></i> DB Migration Sync</button>
                      </div>
                    </div>
                    <div class="col-md-6">
                      <div class="border rounded p-3 bg-light">
                        <div class="d-flex justify-content-between align-items-center mb-2">
                          <h6 class="font-weight-bold mb-0">Database Backups</h6>
                          <button type="button" class="btn btn-sm btn-success" @click="createBackupNow"><i class="icofont-plus"></i> Backup Now</button>
                        </div>
                        <ul class="list-unstyled text-small mb-0">
                          <li v-for="b in adminBackupList" :key="b.id" class="d-flex justify-content-between border-bottom py-1">
                            <span><code>{{ b.id }}</code></span>
                            <span class="badge bg-secondary">{{ b.size }}</span>
                          </li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

          </div><!-- .section-body -->
        </div><!-- .main-content -->

        <footer class="main-footer cms-main-footer">
          <div class="footer-left">
            DESIGN BY: ATENI MEDIA TECHNOLOGIES LLC
          </div>
          <div class="footer-right">
            NATIONAL COUNCIL OF SPORTS PORTAL
          </div>
        </footer>
      </div><!-- .main-wrapper -->
    </div><!-- .otika-app -->

    <!-- MODAL: ADD / EDIT FEDERATION -->
    <div v-if="showAddFederationModal" class="form-backdrop" style="position: fixed; inset: 0; background: rgba(15,23,42,0.6); z-index: 1050; display: grid; place-items: center; padding: 20px;">
      <div class="card shadow-lg border-0" style="max-width: 680px; width: 100%; max-height: 90vh; overflow-y: auto; border-radius: 16px;">
        <div class="card-header bg-primary text-white d-flex justify-content-between align-items-center p-3">
          <h5 class="mb-0 font-weight-bold"><i class="icofont-badge me-2"></i> {{ editingFederation.name ? 'Edit Sports Federation' : 'Add New Sports Federation' }}</h5>
          <button type="button" class="btn-close btn-close-white" @click="showAddFederationModal = false"></button>
        </div>
        <div class="card-body p-4">
          <form @submit.prevent="saveFederation">
            <div class="row g-3">
              <div class="col-md-8">
                <label class="form-label font-weight-bold">Federation / Association Name</label>
                <input v-model="editingFederation.name" type="text" class="form-control" placeholder="e.g. Uganda Athletics Federation" required />
              </div>
              <div class="col-md-4">
                <label class="form-label font-weight-bold">Acronym</label>
                <input v-model="editingFederation.acronym" type="text" class="form-control" placeholder="e.g. UAF" required />
              </div>
              <div class="col-md-6">
                <label class="form-label font-weight-bold">Category Classification</label>
                <select v-model="editingFederation.category" class="form-select">
                  <option value="Category A (Priority)">Category A (Priority - High Impact)</option>
                  <option value="Category B (Established)">Category B (Established)</option>
                  <option value="Category C (Developing)">Category C (Developing)</option>
                  <option value="Category D (Recognized Bodies)">Category D (Recognized Bodies)</option>
                </select>
              </div>
              <div class="col-md-6">
                <label class="form-label font-weight-bold">Recognition Status</label>
                <select v-model="editingFederation.status" class="form-select">
                  <option value="Fully Recognized">Fully Recognized</option>
                  <option value="Provisional Recognition">Provisional Recognition</option>
                  <option value="Under Statutory Audit">Under Statutory Audit</option>
                  <option value="Suspended">Suspended</option>
                </select>
              </div>
              <div class="col-md-6">
                <label class="form-label font-weight-bold">President Name</label>
                <input v-model="editingFederation.president" type="text" class="form-control" placeholder="President full name" />
              </div>
              <div class="col-md-6">
                <label class="form-label font-weight-bold">General Secretary Name</label>
                <input v-model="editingFederation.secretary" type="text" class="form-control" placeholder="General Secretary full name" />
              </div>
              <div class="col-md-6">
                <label class="form-label font-weight-bold">Official Email</label>
                <input v-model="editingFederation.email" type="email" class="form-control" placeholder="info@federation.go.ug" />
              </div>
              <div class="col-12">
                <label class="form-label font-weight-bold">Recognition Certificate Serial #</label>
                <input v-model="editingFederation.certificate_no" type="text" class="form-control" readonly />
              </div>
            </div>
            <div class="d-flex justify-content-end gap-2 mt-4">
              <button type="button" class="btn btn-light" @click="showAddFederationModal = false">Cancel</button>
              <button type="submit" class="btn btn-primary px-4"><i class="icofont-save me-1"></i> Save Federation</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- MODAL: ADD / EDIT FEDERATION CATEGORY -->
    <div v-if="showAddCategoryModal" class="form-backdrop" style="position: fixed; inset: 0; background: rgba(15,23,42,0.6); z-index: 1050; display: grid; place-items: center; padding: 20px;">
      <div class="card shadow-lg border-0" style="max-width: 580px; width: 100%; border-radius: 16px;">
        <div class="card-header bg-dark text-white d-flex justify-content-between align-items-center p-3">
          <h5 class="mb-0 font-weight-bold"><i class="icofont-tags me-2"></i> {{ editingCategory.name ? 'Edit Category Classification' : 'Add New Federation Category' }}</h5>
          <button type="button" class="btn-close btn-close-white" @click="showAddCategoryModal = false"></button>
        </div>
        <div class="card-body p-4">
          <form @submit.prevent="saveFederationCategory">
            <div class="row g-3">
              <div class="col-md-4">
                <label class="form-label font-weight-bold">Category Code</label>
                <input v-model="editingCategory.code" type="text" class="form-control" placeholder="CAT_E" required />
              </div>
              <div class="col-md-8">
                <label class="form-label font-weight-bold">Category Name</label>
                <input v-model="editingCategory.name" type="text" class="form-control" placeholder="Category Name" required />
              </div>
              <div class="col-12">
                <label class="form-label font-weight-bold">Description & Criteria</label>
                <textarea v-model="editingCategory.description" class="form-control" rows="3" placeholder="Category definition, medal targets, grassroots criteria..."></textarea>
              </div>
            </div>
            <div class="d-flex justify-content-end gap-2 mt-4">
              <button type="button" class="btn btn-light" @click="showAddCategoryModal = false">Cancel</button>
              <button type="submit" class="btn btn-primary px-4"><i class="icofont-save me-1"></i> Save Category</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
/* =====================================================================
   SELF-CONTAINED NCS PORTAL LAYOUT
   All Otika-framework layout classes fully inlined — no external CSS
   dependency. Works in any deployment environment.
   ===================================================================== */

/* ── Root ─────────────────────────────────────────────────────────── */
.user-portal {
  min-height: 100vh;
  background-color: #f4f6f9;
  font-family: 'Poppins', 'Inter', system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  font-size: 14px;
  color: #1e293b;
}
.otika-app { display: flex; flex-direction: column; min-height: 100vh; }

/* ── Main wrapper (sidebar + content) ──────────────────────────────── */
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

/* ── Navbar icon buttons ───────────────────────────────────────────── */
.nav-link {
  background: none;
  border: none;
  cursor: pointer;
  color: #64748b;
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
.nav-link:hover { background: #f1f5f9; color: #1e293b; }
.nav-link-lg i { font-size: 18px; }
.cms-top-icon { font-size: 18px; }

/* ── Portal navbar title ───────────────────────────────────────────── */
.portal-navbar-title {
  display: flex; flex-direction: column; line-height: 1.3;
  padding: 0 10px; color: #1e293b;
}
.portal-navbar-title small { font-size: 11px; color: #94a3b8; }
.portal-navbar-title strong { font-size: 14px; font-weight: 600; }

/* ── Avatar / initials ─────────────────────────────────────────────── */
.user-img-radious-style {
  width: 34px; height: 34px;
  border-radius: 50%; object-fit: cover;
}
.nav-avatar-fallback, .sidebar-avatar-fallback {
  display: inline-flex;
  align-items: center; justify-content: center;
  width: 34px; height: 34px;
  border-radius: 50%;
  background: linear-gradient(135deg, #3b82f6, #6366f1);
  color: #fff; font-weight: 700; font-size: 13px;
}

/* ── Dropdown ──────────────────────────────────────────────────────── */
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
}
.dropdown-item:hover { background: #f8fafc; color: #1e293b; }
.dropdown-item.has-icon i { width: 18px; text-align: center; color: #94a3b8; }
.dropdown-item.text-danger { color: #ef4444; }
.dropdown-divider {
  height: 1px; background: #f1f5f9;
  margin: 4px 0;
}
.pullDown { margin-top: 4px; }

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
  padding-top: 60px; /* below navbar */
}

/* Mini sidebar */
.sidebar-mini .main-sidebar { width: 64px; overflow: visible; }
.sidebar-mini .main-sidebar .sidebar-brand span,
.sidebar-mini .main-sidebar .sidebar-user > div,
.sidebar-mini .main-sidebar .sidebar-menu .menu-header,
.sidebar-mini .main-sidebar .sidebar-menu li button span,
.sidebar-mini .main-sidebar .sidebar-menu li a span { display: none; }
.sidebar-mini .main-sidebar .sidebar-brand img { max-width: 36px; }

/* Sidebar brand (logo) */
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

/* Sidebar user info */
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

/* Sidebar menu */
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
.sidebar-menu li button.nav-link {
  display: flex; align-items: center; gap: 10px;
  padding: 9px 16px;
  border-radius: 0; border: none;
  background: none; cursor: pointer;
  width: 100%; text-align: left;
  font-size: 13px; font-weight: 500; color: #374151;
  text-decoration: none;
  transition: background 0.15s, color 0.15s;
}
.sidebar-menu li a:hover,
.sidebar-menu li button:hover {
  background-color: #f1f5f9 !important;
  color: #1e293b !important;
}
.sidebar-menu li.active > a,
.sidebar-menu li.active > button {
  background-color: #ede9fe !important;
  color: #4f46e5 !important;
  font-weight: 600 !important;
  border-left: 3px solid #4f46e5;
}
.sidebar-menu li a i,
.sidebar-menu li button i { font-size: 16px; width: 20px; text-align: center; color: inherit; }
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

/* ── Main content ──────────────────────────────────────────────────── */
.main-content {
  margin-top: 60px;
  margin-left: 240px;
  padding: 24px;
  min-height: calc(100vh - 60px);
  transition: margin-left 0.28s ease;
}
.sidebar-mini .main-content { margin-left: 64px; }
@media (max-width: 768px) {
  .main-content { margin-left: 0; padding: 16px; }
}

/* ── Section header ────────────────────────────────────────────────── */
.section-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 20px;
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
.breadcrumb-item.active a { color: #4f46e5; text-decoration: none; }
.breadcrumb-item + .breadcrumb-item::before {
  content: '/';
  margin-right: 6px;
  color: #cbd5e1;
}

/* ── Actions bar ───────────────────────────────────────────────────── */
.cms-actions.otika-page-actions {
  display: flex; align-items: center; gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

/* ── Buttons ───────────────────────────────────────────────────────── */
.btn {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 8px 18px;
  border-radius: 8px; border: none; cursor: pointer;
  font-size: 13px; font-weight: 600;
  transition: all 0.15s;
}
.btn-primary { background: #4f46e5; color: #fff; }
.btn-primary:hover { background: #4338ca; }
.btn-success { background: #10b981; color: #fff; }
.btn-success:hover { background: #059669; }
.btn-info { background: #0ea5e9; color: #fff; }
.btn-info:hover { background: #0284c7; }
.btn-danger { background: #ef4444; color: #fff; }
.btn-danger:hover { background: #dc2626; }
.btn-secondary { background: #e2e8f0; color: #374151; }
.btn-secondary:hover { background: #cbd5e1; }
.btn-warning { background: #f59e0b; color: #fff; }
.btn-icon { padding: 8px 14px; }
.btn-sm { padding: 5px 12px; font-size: 12px; border-radius: 6px; }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }
.me-2 { margin-right: 8px; }

/* ── Cards / tables ────────────────────────────────────────────────── */
.card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0,0,0,.06);
  border: 1px solid #e8edf2;
  overflow: hidden;
}
.card-body { padding: 20px; }

.table { width: 100%; border-collapse: collapse; }
.table th {
  padding: 11px 16px;
  font-size: 11px; font-weight: 700; text-transform: uppercase;
  letter-spacing: .05em; color: #94a3b8;
  background: #f8fafc;
  border-bottom: 1px solid #e8edf2;
  text-align: left;
}
.table td {
  padding: 12px 16px;
  font-size: 13px; color: #374151;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}
.table tr:last-child td { border-bottom: none; }
.table tr:hover td { background: #fafafa; }

/* ── Badges ────────────────────────────────────────────────────────── */
.badge {
  display: inline-flex; align-items: center;
  padding: 3px 9px;
  border-radius: 20px; font-size: 11px; font-weight: 700;
  text-transform: uppercase; letter-spacing: .03em;
}
.headerBadge1, .headerBadge2 {
  position: absolute; top: 2px; right: 2px;
  min-width: 16px; height: 16px;
  background: #ef4444; color: #fff;
  border-radius: 10px; font-size: 10px; font-weight: 700;
  padding: 0 4px;
  display: flex; align-items: center; justify-content: center;
}
.bg-success, .bg-success-light { background-color: #d1fae5 !important; color: #065f46 !important; }
.bg-primary { background-color: #ede9fe !important; color: #4f46e5 !important; }
.bg-warning { background-color: #fef3c7 !important; color: #92400e !important; }
.bg-danger { background-color: #fee2e2 !important; color: #991b1b !important; }
.bg-dark { background-color: #1e293b !important; color: #fff !important; }
.bg-light { background-color: #f8fafc !important; }

/* ── Messages ──────────────────────────────────────────────────────── */
.portal-success, .portal-error {
  padding: 12px 16px;
  border-radius: 8px; margin-bottom: 16px;
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
  margin-left: 240px;
  transition: margin-left 0.28s;
}
.sidebar-mini footer.main-footer { margin-left: 64px; }
@media (max-width: 768px) { footer.main-footer { margin-left: 0; } }

/* ── Dark mode ─────────────────────────────────────────────────────── */
:global(.dark) .user-portal { background: #0f172a; color: #e2e8f0; }
:global(.dark) .navbar-bg,
:global(.dark) .main-navbar { background: #1e293b; box-shadow: 0 1px 4px rgba(0,0,0,.3); }
:global(.dark) .main-sidebar { background: #1e293b; border-right-color: #334155; }
:global(.dark) .sidebar-brand,
:global(.dark) .sidebar-user { border-bottom-color: #334155; }
:global(.dark) .sidebar-menu li a,
:global(.dark) .sidebar-menu li button.nav-link { color: #cbd5e1; }
:global(.dark) .sidebar-menu li.active > a,
:global(.dark) .sidebar-menu li.active > button {
  background-color: #312e81 !important;
  color: #c7d2fe !important;
  border-left-color: #6366f1;
}
:global(.dark) .sidebar-menu li a:hover,
:global(.dark) .sidebar-menu li button:hover {
  background-color: #334155 !important;
  color: #f8fafc !important;
}
:global(.dark) .card,
:global(.dark) .dropdown-menu { background: #1e293b; border-color: #334155; }
:global(.dark) .table th { background: #1e293b; color: #64748b; }
:global(.dark) .table td { color: #cbd5e1; border-bottom-color: #334155; }
:global(.dark) footer.main-footer { background: #1e293b; border-top-color: #334155; }
:global(.dark) .nav-link { color: #94a3b8; }
:global(.dark) .nav-link:hover { background: #334155; color: #f1f5f9; }
:global(.dark) .section-header h1 { color: #f1f5f9; }

/* ── Misc overrides ────────────────────────────────────────────────── */
.dropdown-menu button.nav-link:hover { background-color: #f1f5f9 !important; }
.mt-4 { margin-top: 16px; }
.mb-4 { margin-bottom: 16px; }
.p-4 { padding: 16px; }
.fw-bold { font-weight: 700; }
</style>