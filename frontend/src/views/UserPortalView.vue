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
import { getMySecurity, enroll2FA, verify2FA, disable2FA, changePassword } from '@/api/security.js'
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
const portalSectionIds = ['dashboard', 'apply', 'applications', 'my-files', 'activities', 'notifications', 'messages', 'transactions', 'profile']
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
const changingPassword = ref(false)
const twofaEnabled = ref(false)
const toggling2FA = ref(false)
const showTwoFAVerify = ref(false)
const twofaVerifyCode = ref('')
const verifying2FA = ref(false)

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
  try {
    const res = await listMyAuditLogs({ page: activityPage.value, per_page: activityPerPage.value })
    const unwrapped = res?.data ?? res ?? {}
    activities.value = Array.isArray(unwrapped.data) ? unwrapped.data : []
    activityTotal.value = unwrapped.meta?.total ?? activities.value.length
  } catch (err) {
    error.value = 'Could not load activity logs.'
  }
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
  if (id === 'profile') {
    try {
      const sec = await getMySecurity()
      twofaEnabled.value = sec?.twofa?.enabled ?? false
    } catch (e) {}
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

.portal-navbar-title {
  display: flex; flex-direction: column; line-height: 1.3;
  padding: 0 10px; color: #1e293b;
}
.portal-navbar-title small { font-size: 11px; color: #94a3b8; }
.portal-navbar-title strong { font-size: 14px; font-weight: 600; }

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
</style>

