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
              <button
                type="button"
                class="btn-close-sidebar"
                title="Close menu"
                aria-label="Close menu"
                @click="mobileSidebarOpen = false"
              >
                <i class="icofont-close"></i>
              </button>
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

                <!-- Federation President / Secretary Recognition License Dashboard Panel -->
                <div v-if="(isFederationLeader || federationLicenseData) && federationData" class="federation-leader-dashboard mt-4">
                  <div class="card border-0 shadow-sm text-white mb-4 overflow-hidden position-relative" style="background: linear-gradient(135deg, #1e3a8a 0%, #0f172a 100%); border-radius: 12px;">
                    <div class="card-body p-4 position-relative z-index-1">
                      <div class="d-flex align-items-center justify-content-between flex-wrap gap-3">
                        <div class="d-flex align-items-center gap-3">
                          <div class="fed-icon-seal" style="width: 50px; height: 50px; border-radius: 12px; background: rgba(255,255,255,0.12); display: flex; align-items: center; justify-content: center; font-size: 26px; color: #f59e0b;">
                            <i class="icofont-certificate-alt-1"></i>
                          </div>
                          <div>
                            <span class="badge bg-warning text-dark text-uppercase mb-1" style="font-size: 10px; font-weight: 700; letter-spacing: 0.5px;">
                              {{ federationLeaderRole || 'Federation Leadership Portal' }}
                            </span>
                            <h3 class="mb-1 fw-bold text-white fs-4">{{ federationData.name }}</h3>
                            <span class="small opacity-85 text-light">
                              NCS Reg No: <strong>{{ federationData.ncs_registration_number || federationData.registration_number || 'NCS-STATUTORY' }}</strong>
                              <span v-if="federationData.acronym" class="ms-2">({{ federationData.acronym }})</span>
                            </span>
                          </div>
                        </div>

                        <div class="d-flex align-items-center gap-2">
                          <button
                            v-if="federationLicenseData"
                            type="button"
                            class="btn btn-warning text-dark fw-bold btn-sm px-3 shadow-sm"
                            @click="printFederationCertificate(federationLicenseData)"
                          >
                            <i class="icofont-print me-1"></i> Print Recognition Certificate
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- License Details Grid -->
                  <div v-if="federationLicenseData" class="row g-3 mb-4">
                    <div class="col-md-3">
                      <div class="p-3 bg-white border rounded shadow-sm h-100">
                        <span class="text-muted small d-block text-uppercase fw-semibold" style="font-size: 11px;">Official License No</span>
                        <strong class="text-dark fs-6 font-monospace">{{ federationLicenseData.license_number }}</strong>
                        <small class="text-muted d-block mt-1">Statutory Identifier</small>
                      </div>
                    </div>

                    <div class="col-md-3">
                      <div class="p-3 bg-white border rounded shadow-sm h-100">
                        <span class="text-muted small d-block text-uppercase fw-semibold" style="font-size: 11px;">License Status</span>
                        <span class="badge" :class="getLicStatusClass(federationLicenseData.status)" style="font-size: 11px; font-weight: 700;">
                          {{ federationLicenseData.status }}
                        </span>
                        <small v-if="federationLicenseData.status === 'ACTIVE'" class="text-success d-block mt-1">
                          <i class="icofont-check-circled"></i> Fully Recognized
                        </small>
                        <small v-else-if="federationLicenseData.status === 'REVOKED'" class="text-danger d-block mt-1">
                          <i class="icofont-ban"></i> Recognition Suspended
                        </small>
                      </div>
                    </div>

                    <div class="col-md-3">
                      <div class="p-3 bg-white border rounded shadow-sm h-100">
                        <span class="text-muted small d-block text-uppercase fw-semibold" style="font-size: 11px;">Issue Date</span>
                        <strong class="text-dark">{{ formatDate(federationLicenseData.issue_date) }}</strong>
                        <small class="text-muted d-block mt-1">National Sports Act 2023</small>
                      </div>
                    </div>

                    <div class="col-md-3">
                      <div class="p-3 bg-white border rounded shadow-sm h-100">
                        <span class="text-muted small d-block text-uppercase fw-semibold" style="font-size: 11px;">Valid Until</span>
                        <strong class="text-dark">{{ formatDate(federationLicenseData.expiry_date) }}</strong>
                        <small class="text-muted d-block mt-1">
                          <span v-if="federationLicenseData.extended_at" class="text-info">Extended</span>
                          <span v-else>Standard Term</span>
                        </small>
                      </div>
                    </div>

                    <!-- Revocation Warning Alert if revoked -->
                    <div v-if="federationLicenseData.status === 'REVOKED'" class="col-12">
                      <div class="alert alert-danger mb-0 d-flex align-items-center gap-3">
                        <i class="icofont-warning fs-3"></i>
                        <div>
                          <strong class="d-block">Statutory Recognition License Revoked</strong>
                          <span>Reason: {{ federationLicenseData.revocation_reason || 'Compliance violation or statutory directive.' }} Please contact NCS Secretariat for rectification.</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

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
                        <td data-label="Payment">
                          <div class="d-flex align-items-center gap-2 flex-wrap">
                            <span class="status" :class="paymentClass(item.payment_status)">{{ titleize(item.payment_status || 'not required') }}</span>
                            <button
                              v-if="canPayApplication(item)"
                              type="button"
                              class="btn-pay-action-pill"
                              title="Make direct payment"
                              @click="openDirectPaymentModal(item)"
                            >
                              <i class="icofont-credit-card"></i> Pay Now
                            </button>
                          </div>
                        </td>
                        <td data-label="Updated">{{ formatDate(item.updated_at) }}</td>
                        <td data-label="Action">
                          <span class="table-actions">
                            <button v-if="canPayApplication(item)" type="button" class="btn-pay-icon text-success" title="Pay application fee" @click="openDirectPaymentModal(item)"><i class="icofont-credit-card"></i></button>
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
                      :class="{ active: activityCategoryFilter === 'PROFILE' }"
                      @click="activityCategoryFilter = 'PROFILE'"
                    >
                      <i class="icofont-user-alt-7"></i> Profile & Settings
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
                      :class="{ active: activityCategoryFilter === 'FINANCE' }"
                      @click="activityCategoryFilter = 'FINANCE'"
                    >
                      <i class="icofont-money"></i> Payments
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
                      :class="{ active: activityCategoryFilter === 'NAVIGATION' }"
                      @click="activityCategoryFilter = 'NAVIGATION'"
                    >
                      <i class="icofont-compass"></i> Page Visits
                    </button>
                  </div>
                </div>

                <!-- Activities Timeline List -->
                <div v-if="loadingActivities" class="text-center py-5">
                  <div class="spinner-border text-primary" role="status" style="width: 3rem; height: 3rem;"></div>
                  <p class="text-muted mt-3 font-weight-bold">Loading activity history...</p>
                </div>
                <div v-else-if="!filteredActivities.length" class="empty-state card text-center p-5 border-0 shadow-sm">
                  <i class="icofont-history" style="font-size: 48px; color: #cbd5e1; margin-bottom: 12px;"></i>
                  <h3 class="font-weight-bold section-subheading">No activities match your filters</h3>
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
                        <span v-if="formatLocationPill(item)" class="meta-pill">
                          <i class="icofont-location-pin"></i> {{ formatLocationPill(item) }}
                        </span>
                        <span v-if="formatDevicePill(item)" class="meta-pill">
                          <i class="icofont-laptop"></i> {{ formatDevicePill(item) }}
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

                <div v-if="!userFiles.length" class="empty-state py-5 card shadow-sm text-center border-0">
                  <div class="card-body p-5">
                    <i class="icofont-folder-open text-muted" style="font-size: 64px;"></i>
                    <h3 class="mt-3 fw-bold empty-title">No Credentials Found</h3>
                    <p class="text-muted max-w-md mx-auto">
                      Currently, there are no active athlete licenses, coach credentials, or technical official clearances linked to your email in the system.
                    </p>
                  </div>
                </div>

                <div v-else class="row">
                  <div v-for="file in userFiles" :key="file.id" class="col-md-6 mb-4">
                    <div class="card shadow-sm border-0 h-100 credential-card">
                      <div class="card-body p-4 d-flex flex-column h-100">
                        <header class="d-flex align-items-start justify-content-between mb-3">
                          <div class="d-flex align-items-center gap-3">
                            <span class="credential-icon p-2 rounded-circle" :class="getFileIconClass(file.category)" style="display: inline-flex; align-items: center; justify-content: center; width: 40px; height: 40px;">
                              <i :class="getFileIcon(file.category)" style="font-size: 20px;"></i>
                            </span>
                            <div>
                              <span class="badge bg-light text-muted text-uppercase mb-1" style="font-size: 9px; font-weight: 700; border: 1px solid rgba(0,0,0,0.06);">{{ file.type }}</span>
                              <h3 class="h5 mb-0 fw-bold credential-title">{{ file.title }}</h3>
                            </div>
                          </div>
                          <span class="badge" :class="getStatusBadgeClass(file.status)" style="font-size: 10px; font-weight: 700;">{{ file.status }}</span>
                        </header>
                        
                        <p class="text-muted small mb-4 flex-grow-1" style="font-size: 12px; line-height: 1.5;">{{ file.description }}</p>
                        
                        <div class="mt-auto pt-3 border-top d-flex flex-column gap-2" style="border-color: rgba(0,0,0,0.06) !important;">
                          <div class="d-flex justify-content-between text-muted small" style="font-size: 12px;">
                            <span>License No:</span>
                            <strong class="credential-number">{{ file.number }}</strong>
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

              <!-- ── Athlete Medical Records & Clearance Files (Table View) ── -->
              <template v-else-if="section === 'athlete-medical'">
                <header class="page-heading">
                  <div>
                    <p>Sports Medicine & Health</p>
                    <h1>Medical Records & Clearance Files</h1>
                    <span>Official NCS medical profile, injury clearance status, and medical clearance certificates.</span>
                  </div>
                </header>

                <div class="user-kpis mb-4">
                  <article>
                    <span class="bg-danger text-white kpi-icon-circle"><i class="icofont-first-aid"></i></span>
                    <div>
                      <small>Blood Group</small>
                      <strong>{{ athleteMedical?.blood_group || 'Not Recorded' }}</strong>
                      <p>Verified profile</p>
                    </div>
                  </article>
                  <article>
                    <span class="bg-success text-white kpi-icon-circle"><i class="icofont-heart-beat"></i></span>
                    <div>
                      <small>Injury Clearance</small>
                      <strong>{{ athleteMedical?.current_injury_status || 'Fit to Compete' }}</strong>
                      <p>NCS medical status</p>
                    </div>
                  </article>
                  <article>
                    <span class="bg-info text-white kpi-icon-circle"><i class="icofont-shield"></i></span>
                    <div>
                      <small>Insurance Status</small>
                      <strong>{{ athleteMedical?.medical_insurance ? 'Covered' : 'Standard Cover' }}</strong>
                      <p>Medical protection</p>
                    </div>
                  </article>
                </div>

                <!-- Medical Profile Table -->
                <div class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom d-flex align-items-center justify-content-between">
                    <h5 class="mb-0 fw-bold"><i class="icofont-id-card text-primary me-2"></i> Clinical & Physical Profile Records</h5>
                    <span class="badge" :class="getStatusBadgeClass(athleteMedical?.current_injury_status || 'FIT TO COMPETE')">
                      {{ athleteMedical?.current_injury_status || 'FIT TO COMPETE' }}
                    </span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Medical Parameter / Check</th>
                          <th>Record Value / Detail</th>
                          <th>Status / Classification</th>
                          <th>Clinical Verification</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-blood-drop text-danger me-1"></i> Blood Group & Genotype</td>
                          <td>
                            <strong class="badge bg-danger text-white fs-6 py-1 px-3">{{ athleteMedical?.blood_group || 'Not Recorded' }}</strong>
                          </td>
                          <td><span class="badge bg-light text-dark">Primary Blood Profile</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Verified</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-warning-alt text-warning me-1"></i> Allergies & Medical Alerts</td>
                          <td><strong>{{ athleteMedical?.allergies || 'None recorded' }}</strong></td>
                          <td><span class="badge" :class="athleteMedical?.allergies ? 'bg-warning text-dark' : 'bg-success-light text-success'">{{ athleteMedical?.allergies ? 'Alert Active' : 'No Allergies' }}</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Verified</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-runner-alt-1 text-primary me-1"></i> Competition Clearance & Fitness</td>
                          <td>
                            <span class="badge" :class="getStatusBadgeClass(athleteMedical?.current_injury_status || 'FIT TO COMPETE')">
                              {{ athleteMedical?.current_injury_status || 'FIT TO COMPETE' }}
                            </span>
                          </td>
                          <td><span class="badge bg-light text-dark">Sports Medicine Clearance</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Active Clearance</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-shield text-info me-1"></i> Medical Insurance Policy</td>
                          <td><strong>{{ athleteMedical?.medical_insurance || 'NCS National Sports Group Insurance' }}</strong></td>
                          <td><span class="badge bg-info text-white">Covered</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Active Policy</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-phone text-secondary me-1"></i> Emergency Medical Contact</td>
                          <td><strong>{{ athleteData?.emergency_contact || athleteData?.phone_contact || 'Registered with Federation' }}</strong></td>
                          <td><span class="badge bg-light text-dark">Primary Contact</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> On File</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-users text-secondary me-1"></i> Next of Kin Contact</td>
                          <td><strong>{{ athleteData?.next_of_kin || 'Registered with Federation' }}</strong></td>
                          <td><span class="badge bg-light text-dark">Guardian / Kin</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> On File</span></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>

                <!-- Injury History Table -->
                <div v-if="athleteMedical?.injury_history" class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom">
                    <h5 class="mb-0 fw-bold"><i class="icofont-history text-danger me-2"></i> Injury History & Treatment Log</h5>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Record #</th>
                          <th>Injury Details & Treatment History</th>
                          <th>Current Status</th>
                          <th>Clearance Status</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td class="fw-bold">#1</td>
                          <td style="white-space: pre-wrap; line-height: 1.6;">{{ athleteMedical.injury_history }}</td>
                          <td>
                            <span class="badge" :class="getStatusBadgeClass(athleteMedical.current_injury_status || 'FIT')">
                              {{ athleteMedical.current_injury_status || 'FIT TO COMPETE' }}
                            </span>
                          </td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Medically Cleared</span></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>

                <!-- Medical Clearance Documents Table -->
                <div class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom d-flex align-items-center justify-content-between">
                    <h5 class="mb-0 fw-bold"><i class="icofont-certificate-alt text-primary me-2"></i> Official Medical Clearance Documents</h5>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Document Title</th>
                          <th>Type</th>
                          <th>Reference Number</th>
                          <th>Issued Date</th>
                          <th>Fitness Status</th>
                          <th>Action</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td class="fw-bold text-dark">
                            <i class="icofont-file-document text-primary me-2"></i> Official Medical Clearance Certificate
                          </td>
                          <td><span class="badge bg-light text-dark">Medical File</span></td>
                          <td><strong class="font-monospace text-primary">NCS-MED-{{ athleteMedical?.id?.substring(0, 8).toUpperCase() || 'VERIFIED' }}</strong></td>
                          <td>{{ athleteMedical?.created_at ? formatDate(athleteMedical.created_at) : formatDate(new Date()) }}</td>
                          <td>
                            <span class="badge" :class="getStatusBadgeClass(athleteMedical?.current_injury_status || 'FIT TO COMPETE')">
                              {{ athleteMedical?.current_injury_status || 'FIT TO COMPETE' }}
                            </span>
                          </td>
                          <td>
                            <button
                              type="button"
                              class="btn btn-primary btn-sm d-flex align-items-center gap-1"
                              style="border-radius: 20px; font-weight: 600; font-size: 11px;"
                              @click="downloadFile({
                                title: 'Official Medical Clearance Certificate',
                                type: 'Medical Records',
                                number: `NCS-MED-${athleteMedical?.id?.substring(0, 8).toUpperCase() || 'VERIFIED'}`,
                                category: 'Medical Files',
                                status: athleteMedical?.current_injury_status || 'Fit to Compete',
                                issueDate: athleteMedical?.created_at || new Date().toISOString(),
                                expiryDate: 'N/A',
                                description: 'Official NCS medical clearance verifying athlete physical condition, blood group, allergies, and injury fitness.',
                                fileType: 'TXT Document',
                                downloadName: `NCS_Medical_Clearance_${athleteData?.athlete_number || 'Athlete'}.txt`
                              })"
                            >
                              <i class="icofont-download"></i> Download Certificate
                            </button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </template>

              <!-- ── Athlete Anti-Doping Records (Table View) ── -->
              <template v-else-if="section === 'athlete-antidoping'">
                <header class="page-heading">
                  <div>
                    <p>Clean Sport & Integrity</p>
                    <h1>Anti-Doping Records & Compliance</h1>
                    <span>WADA testing pool status, sample testing log, education certification, and clean athlete records.</span>
                  </div>
                </header>

                <div class="user-kpis mb-4">
                  <article>
                    <span class="bg-success text-white kpi-icon-circle"><i class="icofont-test-bulb"></i></span>
                    <div>
                      <small>Testing Status</small>
                      <strong>{{ athleteAntiDoping?.testing_status || 'Standard Pool' }}</strong>
                      <p>WADA category</p>
                    </div>
                  </article>
                  <article>
                    <span class="bg-primary text-white kpi-icon-circle"><i class="icofont-certificate"></i></span>
                    <div>
                      <small>WADA Education</small>
                      <strong>{{ athleteAntiDoping?.wada_education_completed || athleteSafeguarding?.anti_doping_education_completed ? 'Certified' : 'Required' }}</strong>
                      <p>Education compliance</p>
                    </div>
                  </article>
                  <article>
                    <span class="bg-dark text-white kpi-icon-circle"><i class="icofont-shield-alt"></i></span>
                    <div>
                      <small>Sanctions / Suspensions</small>
                      <strong>{{ athleteAntiDoping?.suspension_history || 'Clean Record' }}</strong>
                      <p>Disciplinary status</p>
                    </div>
                  </article>
                </div>

                <!-- Anti-Doping Compliance Table -->
                <div class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom d-flex align-items-center justify-content-between">
                    <h5 class="mb-0 fw-bold"><i class="icofont-shield-check text-success me-2"></i> Anti-Doping & WADA Compliance Records</h5>
                    <span class="badge bg-success-light text-success fw-bold px-3 py-2">
                      <i class="icofont-check-circled me-1"></i> Clean Athlete
                    </span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Compliance Parameter</th>
                          <th>Status / Detail</th>
                          <th>Classification Tier</th>
                          <th>Official Verification Status</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-badge text-primary me-1"></i> Testing Pool Tier</td>
                          <td><strong>{{ athleteAntiDoping?.testing_status || 'National Testing Pool (NTP)' }}</strong></td>
                          <td><span class="badge bg-light text-dark">WADA Testing Pool</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Registered</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-calendar text-secondary me-1"></i> Last Tested Date</td>
                          <td><strong>{{ athleteAntiDoping?.last_tested_on ? formatDate(athleteAntiDoping.last_tested_on) : 'Not Sampled (Clean Roster)' }}</strong></td>
                          <td><span class="badge bg-light text-dark">Sample History</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Verified</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-test-bulb text-success me-1"></i> Last Laboratory Test Result</td>
                          <td>
                            <span class="badge" :class="getStatusBadgeClass(athleteAntiDoping?.last_test_result || 'COMPLIANT')">
                              {{ athleteAntiDoping?.last_test_result || 'COMPLIANT / NEGATIVE' }}
                            </span>
                          </td>
                          <td><span class="badge bg-light text-dark">Laboratory Clearance</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Passed</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-certificate text-primary me-1"></i> WADA Anti-Doping Education</td>
                          <td>
                            <span class="badge" :class="athleteAntiDoping?.wada_education_completed || athleteSafeguarding?.anti_doping_education_completed ? 'badge-success' : 'badge-warning'">
                              {{ athleteAntiDoping?.wada_education_completed || athleteSafeguarding?.anti_doping_education_completed ? 'COMPLETED & VERIFIED' : 'PENDING EDUCATION' }}
                            </span>
                          </td>
                          <td><span class="badge bg-light text-dark">Integrity Education</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Certified</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-file-document text-info me-1"></i> Code of Conduct Submission</td>
                          <td>
                            <span class="badge" :class="athleteSafeguarding?.code_of_conduct_signed !== false ? 'badge-success' : 'badge-warning'">
                              {{ athleteSafeguarding?.code_of_conduct_signed !== false ? 'SIGNED & ON FILE' : 'PENDING' }}
                            </span>
                          </td>
                          <td><span class="badge bg-light text-dark">Ethics Agreement</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> On File</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold text-dark"><i class="icofont-shield-alt text-dark me-1"></i> Suspension / Ineligibility History</td>
                          <td><strong>{{ athleteAntiDoping?.suspension_history || 'None (In good standing)' }}</strong></td>
                          <td><span class="badge bg-success text-white">Clean Record</span></td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Zero Sanctions</span></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>

                <!-- Anti-Doping Documents Table -->
                <div class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom">
                    <h5 class="mb-0 fw-bold"><i class="icofont-certificate-alt text-primary me-2"></i> Official Anti-Doping Certificates & Documents</h5>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Certificate Title</th>
                          <th>Type</th>
                          <th>Reference Number</th>
                          <th>Issued Date</th>
                          <th>Status</th>
                          <th>Action</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td class="fw-bold text-dark">
                            <i class="icofont-certificate text-success me-2"></i> WADA Anti-Doping Compliance Certificate
                          </td>
                          <td><span class="badge bg-light text-dark">Compliance Record</span></td>
                          <td><strong class="font-monospace text-primary">NCS-WADA-{{ athleteAntiDoping?.id?.substring(0, 8).toUpperCase() || 'VERIFIED' }}</strong></td>
                          <td>{{ athleteAntiDoping?.last_tested_on ? formatDate(athleteAntiDoping.last_tested_on) : formatDate(new Date()) }}</td>
                          <td>
                            <span class="badge bg-success-light text-success">
                              <i class="icofont-check"></i> Compliant
                            </span>
                          </td>
                          <td>
                            <button
                              type="button"
                              class="btn btn-success btn-sm d-flex align-items-center gap-1"
                              style="border-radius: 20px; font-weight: 600; font-size: 11px;"
                              @click="downloadFile({
                                title: 'WADA Anti-Doping Compliance Certificate',
                                type: 'Compliance Record',
                                number: `NCS-WADA-${athleteAntiDoping?.id?.substring(0, 8).toUpperCase() || 'VERIFIED'}`,
                                category: 'Anti-Doping Compliance',
                                status: athleteAntiDoping?.last_test_result || 'Compliant',
                                issueDate: athleteAntiDoping?.last_tested_on || new Date().toISOString(),
                                expiryDate: 'N/A',
                                description: 'Official NCS verification certificate for WADA anti-doping compliance, testing pool registration, and education completion.',
                                fileType: 'TXT Document',
                                downloadName: `NCS_AntiDoping_Certificate_${athleteData?.athlete_number || 'Athlete'}.txt`
                              })"
                            >
                              <i class="icofont-download"></i> Download Certificate
                            </button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </template>

              <!-- ── Athlete Medals & Accolades (Table View) ── -->
              <template v-else-if="section === 'athlete-medals'">
                <header class="page-heading">
                  <div>
                    <p>Honours & Trophies</p>
                    <h1>Medals & Accolades Registry</h1>
                    <span>Official registry of podium finishes, championship titles, and medals won representing Uganda and affiliated sports federations.</span>
                  </div>
                </header>

                <div class="user-kpis mb-4">
                  <article>
                    <span class="bg-warning text-dark kpi-icon-circle"><i class="icofont-medal"></i></span>
                    <div>
                      <small>Gold Medals</small>
                      <strong>{{ athleteMedals.filter(m => String(m.medal_type || '').toUpperCase().includes('GOLD')).length }}</strong>
                      <p>1st place titles</p>
                    </div>
                  </article>
                  <article>
                    <span class="bg-secondary text-white kpi-icon-circle"><i class="icofont-medal"></i></span>
                    <div>
                      <small>Silver Medals</small>
                      <strong>{{ athleteMedals.filter(m => String(m.medal_type || '').toUpperCase().includes('SILVER')).length }}</strong>
                      <p>Runners up</p>
                    </div>
                  </article>
                  <article>
                    <span class="bg-danger text-white kpi-icon-circle"><i class="icofont-medal"></i></span>
                    <div>
                      <small>Bronze Medals</small>
                      <strong>{{ athleteMedals.filter(m => String(m.medal_type || '').toUpperCase().includes('BRONZE')).length }}</strong>
                      <p>3rd place finishes</p>
                    </div>
                  </article>
                  <article>
                    <span class="bg-primary text-white kpi-icon-circle"><i class="icofont-trophy-alt"></i></span>
                    <div>
                      <small>Total Medals</small>
                      <strong>{{ athleteMedals.length }}</strong>
                      <p>All-time podiums</p>
                    </div>
                  </article>
                </div>

                <div class="d-flex align-items-center justify-content-between mb-4 flex-wrap gap-2">
                  <div class="filter-pills">
                    <button type="button" class="filter-pill" :class="{ active: selectedMedalFilter === 'ALL' }" @click="selectedMedalFilter = 'ALL'">
                      All Medals ({{ athleteMedals.length }})
                    </button>
                    <button type="button" class="filter-pill" :class="{ active: selectedMedalFilter === 'GOLD' }" @click="selectedMedalFilter = 'GOLD'">
                      Gold
                    </button>
                    <button type="button" class="filter-pill" :class="{ active: selectedMedalFilter === 'SILVER' }" @click="selectedMedalFilter = 'SILVER'">
                      Silver
                    </button>
                    <button type="button" class="filter-pill" :class="{ active: selectedMedalFilter === 'BRONZE' }" @click="selectedMedalFilter = 'BRONZE'">
                      Bronze
                    </button>
                  </div>
                </div>

                <div v-if="!filteredAthleteMedals.length" class="empty-state py-5 card shadow-sm text-center border-0">
                  <div class="card-body p-5">
                    <i class="icofont-medal text-muted" style="font-size: 64px;"></i>
                    <h3 class="mt-3 fw-bold empty-title">No Medals Found</h3>
                    <p class="text-muted max-w-md mx-auto">
                      No medals are currently recorded for this filter. As championship and competition results are officially ratified by NCS, your accolades will appear here.
                    </p>
                  </div>
                </div>

                <div v-else class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom d-flex align-items-center justify-content-between">
                    <h5 class="mb-0 fw-bold"><i class="icofont-medal text-warning me-2"></i> Ratified Medals & Championship Accolades ({{ filteredAthleteMedals.length }})</h5>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Medal Tier</th>
                          <th>Event / Discipline</th>
                          <th>Competition Title</th>
                          <th>Host Venue / Country</th>
                          <th>Date / Year Won</th>
                          <th>Coach Responsible</th>
                          <th>Prize / Recognition</th>
                          <th>NCS Ratification</th>
                          <th>Action</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="medal in filteredAthleteMedals" :key="medal.id">
                          <td>
                            <span :class="getMedalBadgeClass(medal.medal_type)">
                              <i class="icofont-medal me-1"></i> {{ medal.medal_type }}
                            </span>
                          </td>
                          <td class="fw-bold text-dark">{{ medal.event || 'Championship Event' }}</td>
                          <td>{{ medal.competition_name || 'National Championship' }}</td>
                          <td><i class="icofont-location-pin text-muted me-1"></i> {{ medal.country || 'Uganda' }}</td>
                          <td><strong>{{ medal.won_on ? formatDate(medal.won_on) : (medal.year || 'Accredited') }}</strong></td>
                          <td>{{ medal.coach_responsible || '-' }}</td>
                          <td>
                            <strong v-if="medal.prize_money" class="text-success">{{ formatUGX(medal.prize_money) }}</strong>
                            <span v-else class="text-muted small">-</span>
                          </td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Ratified</span></td>
                          <td>
                            <button
                              type="button"
                              class="btn btn-outline-primary btn-sm d-flex align-items-center gap-1"
                              style="border-radius: 20px; font-weight: 600; font-size: 11px;"
                              @click="downloadFile({
                                title: `Certificate of Achievement - ${medal.medal_type} Medal`,
                                type: 'Honours & Awards',
                                number: `NCS-MEDAL-${medal.id?.substring(0, 8).toUpperCase() || 'HONOUR'}`,
                                category: 'Athlete Registry',
                                status: 'Verified',
                                issueDate: medal.won_on || new Date().toISOString(),
                                expiryDate: 'Lifetime Recognition',
                                description: `Official NCS Certificate recognizing ${medal.medal_type} medal achievement in ${medal.event || 'Sports Championship'}.`,
                                fileType: 'TXT Document',
                                downloadName: `NCS_Medal_Certificate_${medal.medal_type}_${athleteData?.athlete_number || 'Award'}.txt`
                              })"
                            >
                              <i class="icofont-download"></i> Certificate
                            </button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </template>

              <!-- ── Athlete Competition Logs & Results (Table View) ── -->
              <template v-else-if="section === 'athlete-competitions'">
                <header class="page-heading">
                  <div>
                    <p>Performance & Tournaments</p>
                    <h1>Competitions & Results Logs</h1>
                    <span>Complete history of event entries, tournament stages, race timings, match scores, and national team appearances.</span>
                  </div>
                </header>

                <!-- National Team Duty Table -->
                <div v-if="athleteNationalTeam" class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom d-flex align-items-center justify-content-between">
                    <h5 class="mb-0 fw-bold"><i class="icofont-flag text-danger me-2"></i> National Team Appearances & International Caps</h5>
                    <span class="badge bg-warning text-dark fw-bold px-3 py-1">
                      {{ athleteNationalTeam.category || 'Senior National Squad' }}
                    </span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Squad Tier / Category</th>
                          <th>National Team Name</th>
                          <th>First Call-Up Date</th>
                          <th>Last Appearance Date</th>
                          <th>Total International Caps</th>
                          <th>Official Notes</th>
                          <th>Status</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td><span class="badge bg-primary text-white">{{ athleteNationalTeam.category || 'Senior National Team' }}</span></td>
                          <td class="fw-bold text-dark">{{ athleteNationalTeam.team_name || 'Uganda National Sports Team' }}</td>
                          <td>{{ athleteNationalTeam.first_call_up_on ? formatDate(athleteNationalTeam.first_call_up_on) : 'Verified Member' }}</td>
                          <td>{{ athleteNationalTeam.last_appearance_on ? formatDate(athleteNationalTeam.last_appearance_on) : 'Active Squad' }}</td>
                          <td><strong class="fs-6 text-primary">{{ athleteNationalTeam.appearances_count || 1 }} Caps</strong></td>
                          <td>{{ athleteNationalTeam.notes || 'Official Squad Selection' }}</td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Active Roster</span></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>

                <div v-if="!athleteResults.length" class="empty-state py-5 card shadow-sm text-center border-0 mb-4">
                  <div class="card-body p-5">
                    <i class="icofont-trophy text-muted" style="font-size: 64px;"></i>
                    <h3 class="mt-3 fw-bold empty-title">No Competition Logs Recorded</h3>
                    <p class="text-muted max-w-md mx-auto">
                      Your competition entries and timing results will appear here as they are entered and verified by tournament officials and federation administrators.
                    </p>
                  </div>
                </div>

                <!-- Competition Results Table -->
                <div v-else class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom d-flex align-items-center justify-content-between">
                    <h5 class="mb-0 fw-bold"><i class="icofont-listine-dots text-primary me-2"></i> Ratified Competition Results & Performance Logs ({{ athleteResults.length }})</h5>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Event / Discipline</th>
                          <th>Competition Name</th>
                          <th>Stage</th>
                          <th>Position / Rank</th>
                          <th>Result (Time / Score / Distance)</th>
                          <th>Milestone Breaks</th>
                          <th>Ratification Status</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="res in athleteResults" :key="res.id">
                          <td class="fw-bold text-dark">{{ res.event || 'Athletics Discipline' }}</td>
                          <td>{{ res.competition_name || 'National Championship' }}</td>
                          <td><span class="badge bg-light text-dark">{{ res.stage || 'Final' }}</span></td>
                          <td>
                            <strong class="text-primary fs-6">{{ res.position ? `#${res.position}` : '-' }}</strong>
                          </td>
                          <td>
                            <span class="font-monospace fw-bold">{{ res.time_result || res.score_result || res.distance_result || '-' }}</span>
                          </td>
                          <td>
                            <div class="d-flex gap-1 flex-wrap">
                              <span v-if="res.is_national_record" class="badge bg-danger text-white">NR</span>
                              <span v-if="res.is_personal_best" class="badge bg-success text-white">PB</span>
                              <span v-if="res.is_seasonal_best" class="badge bg-info text-white">SB</span>
                              <span v-if="!res.is_national_record && !res.is_personal_best && !res.is_seasonal_best" class="text-muted small">-</span>
                            </div>
                          </td>
                          <td>
                            <span class="badge bg-success-light text-success"><i class="icofont-check"></i> Ratified</span>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </template>

              <!-- ── Athlete Talent Pathways & Transfers (Table View) ── -->
              <template v-else-if="section === 'athlete-pathways'">
                <header class="page-heading">
                  <div>
                    <p>Development & Scouting</p>
                    <h1>Talent Pathways & Career Progression</h1>
                    <span>Scouting evaluation history, talent identification center records, development milestones, and transfer logs.</span>
                  </div>
                </header>

                <!-- Talent Scouting Records Table -->
                <div class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom d-flex align-items-center justify-content-between">
                    <h5 class="mb-0 fw-bold"><i class="icofont-chart-growth text-success me-2"></i> Talent Scouting & Identification Records</h5>
                    <span class="badge bg-success-light text-success fw-bold px-3 py-1">
                      <i class="icofont-check-circled me-1"></i> NCS Pathway Roster
                    </span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Record Reference</th>
                          <th>Identified By (Scout)</th>
                          <th>Identification Date</th>
                          <th>Age at Identification</th>
                          <th>Talent Centre / School</th>
                          <th>Recommended Pathway</th>
                          <th>Sports Scholarship</th>
                          <th>Status</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-if="athleteTalent.length" v-for="t in athleteTalent" :key="t.id">
                          <td><strong class="font-monospace text-primary">TR-{{ t.id?.substring(0, 8).toUpperCase() || 'NCS' }}</strong></td>
                          <td class="fw-bold text-dark">{{ t.identified_by || 'National Talent Scout' }}</td>
                          <td>{{ t.identified_on ? formatDate(t.identified_on) : 'N/A' }}</td>
                          <td>{{ t.age_at_identification || 'Junior' }} yrs</td>
                          <td>{{ t.talent_centre || t.school || athleteData?.education_institution || 'National Development Center' }}</td>
                          <td><span class="badge bg-primary text-white">{{ t.recommended_pathway || 'Elite National Squad' }}</span></td>
                          <td>
                            <span class="badge" :class="t.scholarship_status || athleteData?.sports_scholarship_status ? 'badge-success' : 'badge-secondary'">
                              {{ t.scholarship_status || (athleteData?.sports_scholarship_status ? 'Active Scholarship' : 'None') }}
                            </span>
                          </td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Enrolled</span></td>
                        </tr>
                        <tr v-else>
                          <td><strong class="font-monospace text-primary">TR-{{ athleteData?.id?.substring(0, 8).toUpperCase() || 'PROFILE' }}</strong></td>
                          <td class="fw-bold text-dark">National Federation Scout</td>
                          <td>{{ athleteData?.created_at ? formatDate(athleteData.created_at) : formatDate(new Date()) }}</td>
                          <td>{{ athleteData?.age_category || 'Senior' }}</td>
                          <td>{{ athleteData?.education_institution || 'Federation Training Center' }}</td>
                          <td><span class="badge bg-primary text-white">Elite National Pathway</span></td>
                          <td>
                            <span class="badge" :class="athleteData?.sports_scholarship_status ? 'badge-success' : 'badge-secondary'">
                              {{ athleteData?.sports_scholarship_status ? 'Active Scholarship' : 'Standard Development' }}
                            </span>
                          </td>
                          <td><span class="badge bg-success-light text-success"><i class="icofont-check"></i> Active</span></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>

                <!-- Career Progression Pathway Milestones Table -->
                <div class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom">
                    <h5 class="mb-0 fw-bold"><i class="icofont-chart-growth text-primary me-2"></i> Career Progression Pathway Milestones</h5>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Stage #</th>
                          <th>Development Pathway Tier</th>
                          <th>Focus & Milestone Target</th>
                          <th>Organization / Facility</th>
                          <th>Status</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td class="fw-bold">Stage 1</td>
                          <td class="fw-bold text-dark">Grassroots & District School Identification</td>
                          <td>Identification at district tournaments and school talent showcases.</td>
                          <td>{{ athleteData?.education_institution || 'District School / Academy' }}</td>
                          <td><span class="badge bg-success text-white"><i class="icofont-check"></i> Completed</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold">Stage 2</td>
                          <td class="fw-bold text-dark">Affiliated Club Academy & Regional Training</td>
                          <td>Enrolled with {{ athleteData?.club || 'Club Academy' }} under NCS licensing.</td>
                          <td>{{ athleteData?.club || 'Club Academy Center' }}</td>
                          <td><span class="badge bg-success text-white"><i class="icofont-check"></i> Completed</span></td>
                        </tr>
                        <tr>
                          <td class="fw-bold">Stage 3</td>
                          <td class="fw-bold text-dark">Junior / Senior National Squad Development</td>
                          <td>{{ athleteNationalTeam ? 'Active national squad member.' : 'Target pathway for upcoming trials and selection.' }}</td>
                          <td>National High-Performance Center</td>
                          <td>
                            <span class="badge" :class="athleteNationalTeam ? 'bg-success text-white' : 'bg-warning text-dark'">
                              <i :class="athleteNationalTeam ? 'icofont-check' : 'icofont-clock-time'"></i>
                              {{ athleteNationalTeam ? 'Active Member' : 'In Progress' }}
                            </span>
                          </td>
                        </tr>
                        <tr>
                          <td class="fw-bold">Stage 4</td>
                          <td class="fw-bold text-dark">International Podium & Elite Tier</td>
                          <td>Continental and Olympic representation with NCS High-Performance support.</td>
                          <td>Uganda Olympic Committee / NCS</td>
                          <td><span class="badge bg-light text-dark"><i class="icofont-star"></i> Target Elite Tier</span></td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </template>

              <!-- ── Athlete Clubs & Academies Record (Table View Only) ── -->
              <template v-else-if="section === 'athlete-clubs'">
                <header class="page-heading">
                  <div>
                    <p>Affiliations & Teams</p>
                    <h1>Clubs & Academies Record</h1>
                    <span>Registered primary club, academy training center, and federation affiliation credentials.</span>
                  </div>
                </header>

                <!-- Clubs & Academies Table -->
                <div class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom d-flex align-items-center justify-content-between">
                    <h5 class="mb-0 fw-bold"><i class="icofont-building-alt text-primary me-2"></i> Clubs & Academies Affiliation Records</h5>
                    <span class="badge bg-success-light text-success fw-bold px-3 py-2">
                      <i class="icofont-check-circled me-1"></i> Active Registration
                    </span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Affiliation Type</th>
                          <th>Club / Academy Name</th>
                          <th>Acronym</th>
                          <th>District & Region</th>
                          <th>Contact Person</th>
                          <th>Official Email</th>
                          <th>Contact Phone</th>
                          <th>NCS Recognition Status</th>
                          <th>Action</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td><span class="badge bg-primary text-white">Primary Club</span></td>
                          <td class="fw-bold text-dark">{{ athleteClubDetails?.name || athleteData?.club || 'National Sports Club' }}</td>
                          <td><span class="badge bg-light text-dark">{{ athleteClubDetails?.acronym || 'CLUB' }}</span></td>
                          <td>{{ athleteClubDetails?.district || athleteData?.district || 'Kampala' }} ({{ athleteClubDetails?.region || athleteData?.region || 'Central' }})</td>
                          <td>{{ athleteClubDetails?.contact_person || 'Federation Secretariat' }}</td>
                          <td>{{ athleteClubDetails?.email || 'club@ncs.ug' }}</td>
                          <td>{{ athleteClubDetails?.phone || athleteData?.phone_contact || '-' }}</td>
                          <td><span class="badge bg-success text-white">RECOGNIZED & LICENSED</span></td>
                          <td>
                            <button
                              type="button"
                              class="btn btn-outline-primary btn-sm d-flex align-items-center gap-1"
                              style="border-radius: 20px; font-weight: 600; font-size: 11px;"
                              @click="downloadFile({
                                title: 'Club Affiliation & Membership Certificate',
                                type: 'Club Affiliation',
                                number: `NCS-CLUB-${athleteClubDetails?.id?.substring(0, 8).toUpperCase() || 'AFFIL'}`,
                                category: 'Athlete Registry',
                                status: 'Verified Active',
                                issueDate: athleteData?.created_at || new Date().toISOString(),
                                expiryDate: 'Active',
                                description: `Official NCS verification for athlete club registration with ${athleteClubDetails?.name || athleteData?.club || 'Sports Club'}.`,
                                fileType: 'TXT Document',
                                downloadName: `NCS_Club_Affiliation_${athleteData?.athlete_number || 'Cert'}.txt`
                              })"
                            >
                              <i class="icofont-download"></i> Certificate
                            </button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>

                <!-- Athlete License & Federation Link Table -->
                <div class="card shadow-sm border-0 mb-4 athlete-panel-card">
                  <div class="card-header bg-white py-3 border-bottom d-flex align-items-center justify-content-between">
                    <h5 class="mb-0 fw-bold"><i class="icofont-license text-info me-2"></i> Athlete License & Federation Registry</h5>
                    <span class="badge" :class="athleteData?.status === 'ACTIVE' || !athleteData ? 'badge-success' : 'badge-warning'">
                      {{ athleteData?.status || 'ACTIVE' }}
                    </span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover mb-0">
                      <thead class="bg-light">
                        <tr>
                          <th>Athlete License Number</th>
                          <th>Discipline / Sport</th>
                          <th>Age Division</th>
                          <th>National Federation Linkage</th>
                          <th>License Expiry</th>
                          <th>Status</th>
                          <th>Action</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td><strong class="font-monospace text-primary fs-6">{{ athleteData?.athlete_number || 'NCS-ATH-PENDING' }}</strong></td>
                          <td class="fw-bold text-dark">{{ athleteData?.discipline || 'Athletics' }}</td>
                          <td><span class="badge bg-light text-dark">{{ athleteData?.age_category || 'Senior' }}</span></td>
                          <td>Uganda National Sports Federation</td>
                          <td>Active (Renewable Annually)</td>
                          <td>
                            <span class="badge" :class="athleteData?.status === 'ACTIVE' || !athleteData ? 'badge-success' : 'badge-warning'">
                              {{ athleteData?.status || 'ACTIVE' }}
                            </span>
                          </td>
                          <td>
                            <button
                              type="button"
                              class="btn btn-primary btn-sm d-flex align-items-center gap-1"
                              style="border-radius: 20px; font-weight: 600; font-size: 11px;"
                              @click="downloadFile({
                                title: 'National Athlete License Certificate',
                                type: 'License / Certificate',
                                number: athleteData?.athlete_number || 'NCS-ATH-CERT',
                                category: 'Athlete Registry',
                                status: athleteData?.status || 'Active',
                                issueDate: athleteData?.created_at || new Date().toISOString(),
                                expiryDate: 'N/A (Active)',
                                description: `Official NCS verification for athlete classification under ${athleteData?.discipline || 'sports registry'}.`,
                                fileType: 'TXT Document',
                                downloadName: `NCS_Athlete_License_${athleteData?.athlete_number || 'Cert'}.txt`
                              })"
                            >
                              <i class="icofont-download"></i> Download License
                            </button>
                          </td>
                        </tr>
                      </tbody>
                    </table>
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
                          <h3 class="mb-1 font-weight-bold sessions-heading">Active Login Sessions</h3>
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

    <!-- Direct Application Payment Modal -->
    <div v-if="paymentModalOpen" class="portal-modal-backdrop" @click.self="closeDirectPaymentModal">
      <div class="portal-modal-card">
        <div class="portal-modal-header">
          <div class="d-flex align-items-center gap-3">
            <span class="modal-header-icon">
              <i class="icofont-credit-card"></i>
            </span>
            <div>
              <h3 class="modal-title">Pay Application Fee</h3>
              <p class="modal-subtitle">{{ paymentModalApp?.title || 'Online Application' }}</p>
            </div>
          </div>
          <button type="button" class="btn-close-modal" @click="closeDirectPaymentModal">
            <i class="icofont-close"></i>
          </button>
        </div>

        <div class="portal-modal-body">
          <!-- Fee summary banner -->
          <div class="modal-fee-banner">
            <div class="fee-meta">
              <small>Required Application Fee</small>
              <div class="fee-val">UGX {{ formatMoney(getAppPrice(paymentModalApp)) }}</div>
            </div>
            <div class="fee-ref">
              <small>Ref / Submission</small>
              <strong>{{ paymentModalApp?.reference || (paymentModalApp?.id ? paymentModalApp.id.slice(0, 8) : 'Pending') }}</strong>
            </div>
          </div>

          <!-- Method Selector (Only rendered if multiple channels are allowed) -->
          <div v-if="getAppAllowedMethods(paymentModalApp).length > 1" class="payment-method-selector-tabs mb-4">
            <button
              type="button"
              class="selector-tab-btn"
              :class="{ active: paymentModalMethod === 'MOBILE_MONEY' }"
              @click="paymentModalMethod = 'MOBILE_MONEY'"
            >
              <i class="icofont-smart-phone"></i>
              <div>
                <strong>Mobile Money</strong>
                <small>Instant USSD Prompt (MTN / Airtel)</small>
              </div>
            </button>
            <button
              type="button"
              class="selector-tab-btn"
              :class="{ active: paymentModalMethod === 'OVER_THE_COUNTER' }"
              @click="paymentModalMethod = 'OVER_THE_COUNTER'"
            >
              <i class="icofont-bank-alt"></i>
              <div>
                <strong>Bank / Over Counter</strong>
                <small>URA PRN / Bank Deposit Slip</small>
              </div>
            </button>
          </div>

          <!-- Alert / Error / Success Banners -->
          <div v-if="paymentModalError" class="modal-alert-box alert-error mb-3">
            <i class="icofont-warning-alt me-2"></i> {{ paymentModalError }}
          </div>
          <div v-if="paymentModalSuccess" class="modal-alert-box alert-success mb-3">
            <i class="icofont-check-circled me-2"></i> {{ paymentModalSuccess }}
          </div>

          <!-- Mobile Money Flow -->
          <div v-if="paymentModalMethod === 'MOBILE_MONEY' && !paymentModalSuccess" class="momo-flow-content">
            <div v-if="!paymentModalTracking">
              <div class="form-group mb-3">
                <label class="form-label fw-bold">Mobile Money Phone Number</label>
                <div class="input-group">
                  <span class="input-group-text bg-light fw-bold"><i class="icofont-phone me-1"></i> +256</span>
                  <input
                    v-model="paymentModalPhone"
                    type="tel"
                    class="form-control"
                    placeholder="770000000 or 0111777771"
                  />
                </div>
                <div class="d-flex align-items-center justify-content-between mt-2">
                  <small class="text-muted">Enter your registered MTN or Airtel Uganda phone number.</small>
                  <span v-if="getCarrierName(paymentModalPhone)" class="carrier-tag" :class="getCarrierName(paymentModalPhone).toLowerCase()">
                    {{ getCarrierName(paymentModalPhone) }}
                  </span>
                </div>
              </div>

              <!-- Quick test numbers -->
              <div class="sandbox-hints-box mb-4">
                <span class="hints-label"><i class="icofont-info-circle"></i> Quick Test Numbers:</span>
                <div class="d-flex gap-2 flex-wrap">
                  <button type="button" class="btn btn-xs btn-outline-success" @click="paymentModalPhone = '0111777771'">
                    0111777771 (Success)
                  </button>
                  <button type="button" class="btn btn-xs btn-outline-danger" @click="paymentModalPhone = '0111777991'">
                    0111777991 (Failed)
                  </button>
                </div>
              </div>

              <button
                type="button"
                class="btn btn-primary w-100 py-3 fw-bold btn-momo-submit"
                :disabled="paymentModalSubmitting"
                @click="submitModalMoMoPayment"
              >
                <i v-if="paymentModalSubmitting" class="icofont-spinner-alt-3 animate-spin me-2"></i>
                <i v-else class="icofont-smart-phone me-2"></i>
                <span>Pay UGX {{ formatMoney(getAppPrice(paymentModalApp)) }} via Mobile Money</span>
              </button>
            </div>

            <!-- Active USSD Tracking & Polling State -->
            <div v-else class="ussd-tracking-card text-center py-4">
              <div class="ussd-spinner-wrap mb-3">
                <div class="pulsing-circle"></div>
                <i class="icofont-smart-phone ussd-center-icon"></i>
              </div>
              <h4 class="fw-bold mb-2">USSD Prompt Sent!</h4>
              <p class="text-muted max-w-sm mx-auto mb-3">
                Please check your phone (<strong>{{ paymentModalPhone }}</strong>) and enter your Mobile Money PIN to approve the payment of <strong>UGX {{ formatMoney(getAppPrice(paymentModalApp)) }}</strong>.
              </p>
              <div class="tracking-progress-box mx-auto mb-3">
                <span class="spinner-border spinner-border-sm text-primary me-2"></span>
                <span>Awaiting payment approval... ({{ paymentModalCountdown }}s)</span>
              </div>
              <p class="text-xs text-muted">
                Transaction Reference: <code>{{ paymentModalTxRef || 'NCS-TXN-PENDING' }}</code>
              </p>
            </div>
          </div>

          <!-- Over the Counter Flow -->
          <div v-else-if="paymentModalMethod === 'OVER_THE_COUNTER' && !paymentModalSuccess" class="otc-flow-content">
            <div class="form-group mb-3">
              <label class="form-label fw-bold">URA Payment Registration Number (PRN)</label>
              <input
                v-model="paymentModalPRN"
                type="text"
                class="form-control"
                placeholder="e.g. 224000123456"
              />
            </div>
            <div class="form-group mb-4">
              <label class="form-label fw-bold">Upload Bank Deposit Slip / Proof (PDF, PNG, JPG)</label>
              <input
                type="file"
                class="form-control"
                accept=".pdf,.png,.jpg,.jpeg"
                @change="onModalFileSelected"
              />
            </div>
            <button
              type="button"
              class="btn btn-primary w-100 py-3 fw-bold"
              :disabled="paymentModalSubmitting"
              @click="submitModalOTCProof"
            >
              <i v-if="paymentModalSubmitting" class="icofont-spinner-alt-3 animate-spin me-2"></i>
              <i v-else class="icofont-upload me-2"></i>
              <span>Submit Bank Deposit Proof</span>
            </button>
          </div>

          <!-- Success State Action -->
          <div v-if="paymentModalSuccess" class="text-center py-3">
            <button type="button" class="btn btn-success px-4 py-2 fw-bold" @click="closeDirectPaymentModal">
              <i class="icofont-check me-2"></i> Done / Continue
            </button>
          </div>
        </div>
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
import {
  portalListOpenForms,
  portalListSubmissions,
  portalInitiateMoMoPayment,
  portalGetPaymentStatus,
  portalUploadPaymentProof,
} from '@/api/forms.js'
import { listNsmisDomain } from '@/api/nsmis.js'
import { getFederationActiveLicense } from '@/api/federationLicenses.js'
import { mediaUrl } from '@/api/client.js'
import OpenFormsPanel from '@/components/portal/OpenFormsPanel.vue'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'
import { downloadApplicationForm } from '@/utils/applicationDownload.js'
import { ensureOtikaStyles } from '@/utils/otikaAssets.js'
import { recordMenuNavigation, recordActivityEvent } from '@/services/activityAudit.js'

const router = useRouter()
const portalSectionIds = [
  'dashboard', 'apply', 'applications', 'my-files',
  'athlete-medical', 'athlete-antidoping', 'athlete-medals', 'athlete-competitions', 'athlete-pathways', 'athlete-clubs',
  'activities', 'notifications', 'messages', 'transactions', 'settings', 'profile'
]
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
const athleteTalent = ref([])
const athleteClubDetails = ref(null)
const selectedMedalFilter = ref('ALL')

const isCoach = ref(false)
const coachData = ref(null)
const isOfficial = ref(false)
const officialData = ref(null)
const isFederationLeader = ref(false)
const federationLeaderRole = ref('')
const federationData = ref(null)
const federationLicenseData = ref(null)

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
  if (isFederationLeader.value) return `${federationLeaderRole.value || 'Federation Executive'} Portal`
  const roles = userRoles.value
  if (roles.includes('athlete') || roles.includes('role_athlete')) return 'Athlete Portal'
  if (roles.includes('coach') || roles.includes('role_coach')) return 'Coach Portal'
  if (roles.includes('technical_official') || roles.includes('role_technical_official')) return 'Official Portal'
  return 'Ordinary User'
})

const userWorkspaceLabel = computed(() => {
  if (isFederationLeader.value && federationData.value?.name) return `${federationData.value.name} workspace`
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
    items.push({ id: 'my-files', label: 'My Files & Credentials', icon: 'icofont-folder-open' })
  }

  // Dedicated Athlete Sections on Sidebar
  if (isAthlete.value || userRoles.value.includes('athlete') || userRoles.value.includes('role_athlete')) {
    items.push(
      { id: 'athlete-medical', label: 'Medical Records', icon: 'icofont-first-aid' },
      { id: 'athlete-antidoping', label: 'Anti-Doping Records', icon: 'icofont-test-bulb' },
      { id: 'athlete-medals', label: 'Medals & Accolades', icon: 'icofont-medal', badge: athleteMedals.value.length || '' },
      { id: 'athlete-competitions', label: 'Competitions & Results', icon: 'icofont-trophy', badge: athleteResults.value.length || '' },
      { id: 'athlete-pathways', label: 'Talent Pathways', icon: 'icofont-chart-growth' },
      { id: 'athlete-clubs', label: 'Clubs & Academies', icon: 'icofont-building-alt' }
    )
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
      if (activityCategoryFilter.value === 'FINANCE' && cat !== 'FINANCE') return false
      if (activityCategoryFilter.value === 'NAVIGATION' && cat !== 'NAVIGATION') return false
    }
    const query = activitySearchQuery.value.trim().toLowerCase()
    if (query) {
      const matchText = [
        activityTitle(item),
        activityDescription(item),
        getActivityCategory(item),
        item.action,
        item.resource,
        item.endpoint,
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
  let profileCount = 0
  let financeCount = 0
  items.forEach(item => {
    const cat = getActivityCategory(item)
    if (cat === 'AUTH') authCount++
    if (cat === 'APPLICATION') appCount++
    if (cat === 'PROFILE') profileCount++
    if (cat === 'FINANCE') financeCount++
  })
  const latest = items[0]?.created_at ? formatDate(items[0].created_at) : 'None'
  return {
    authCount,
    appCount,
    profileCount,
    financeCount,
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

  if (federationLicenseData.value) {
    filesList.push({
      id: federationLicenseData.value.id,
      title: `Statutory Recognition License - ${federationLicenseData.value.federation_name || federationData.value?.name || 'Federation'}`,
      type: 'Statutory License',
      number: federationLicenseData.value.license_number,
      issueDate: federationLicenseData.value.issue_date,
      expiryDate: federationLicenseData.value.expiry_date,
      status: federationLicenseData.value.status,
      description: `Official statutory license of recognition under the National Sports Act 2023. Category: ${federationLicenseData.value.category || 'National Sports Federation'}.`,
      fileType: 'Official License',
      downloadName: `NCS_Federation_License_${(federationLicenseData.value.license_number || '').replaceAll('/', '_')}.txt`,
      category: 'Federation Recognition',
      licenseObj: federationLicenseData.value
    })
  }

  return filesList
})

function getLicStatusClass(status) {
  const s = String(status || '').toUpperCase()
  if (s === 'ACTIVE') return 'bg-success text-white'
  if (s === 'EXTENDED') return 'bg-info text-white'
  if (s === 'REVOKED') return 'bg-danger text-white'
  if (s === 'EXPIRED') return 'bg-warning text-dark'
  return 'bg-secondary text-white'
}

function printFederationCertificate(lic) {
  if (!lic) return
  const printWin = window.open('', '_blank', 'width=900,height=800')
  if (!printWin) return

  const issueDateStr = formatDate(lic.issue_date)
  const expiryDateStr = formatDate(lic.expiry_date)

  const html = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>NCS Certificate of Recognition - ${lic.license_number}</title>
  <style>
    @page { size: A4 landscape; margin: 10mm; }
    body { font-family: "Georgia", "Times New Roman", serif; background: #fafafa; margin: 0; padding: 20px; color: #1e293b; }
    .cert-frame { border: 8px double #b45309; padding: 30px 40px; background: #fff; text-align: center; border-radius: 4px; box-shadow: 0 0 20px rgba(0,0,0,0.05); }
    .logo-row { margin-bottom: 12px; }
    .logo-row img { max-height: 70px; }
    .republic-title { font-size: 16px; font-weight: bold; letter-spacing: 2px; text-transform: uppercase; color: #b45309; margin: 0; }
    .ncs-title { font-size: 26px; font-weight: bold; color: #0f172a; margin: 6px 0 16px; text-transform: uppercase; letter-spacing: 1px; }
    .cert-heading { font-size: 20px; font-style: italic; color: #475569; margin: 0 0 10px; }
    .cert-body { font-size: 15px; color: #334155; margin: 0 auto 16px; max-width: 700px; line-height: 1.6; }
    .fed-name { font-size: 28px; font-weight: bold; color: #1e3a8a; margin: 10px 0; text-decoration: underline; text-underline-offset: 6px; }
    .reg-tag { font-size: 13px; color: #64748b; margin-bottom: 16px; }
    .meta-box { display: flex; justify-content: space-around; margin: 24px auto; max-width: 650px; background: #fefce8; border: 1px solid #fef08a; padding: 12px; border-radius: 6px; }
    .meta-item strong { display: block; font-size: 14px; color: #713f12; }
    .meta-item span { font-size: 11px; color: #854d0e; text-transform: uppercase; }
    .conditions { font-size: 11px; font-style: italic; color: #64748b; margin: 14px auto; max-width: 600px; }
    .signatures { display: flex; justify-content: space-between; margin-top: 40px; padding: 0 40px; }
    .sig-line { width: 220px; border-top: 1px solid #334155; padding-top: 6px; font-size: 12px; font-weight: bold; text-align: center; }
    .sig-title { font-size: 11px; color: #64748b; font-weight: normal; }
  </style>
</head>
<body>
  <div class="cert-frame">
    <div class="logo-row">
      <img src="/main-logo.png" alt="National Council of Sports" />
    </div>
    <div class="republic-title">Republic of Uganda</div>
    <div class="ncs-title">National Council of Sports</div>
    <div class="cert-heading">Certificate of Statutory Recognition & Licensing</div>
    
    <div class="cert-body">
      This is to certify that under the provisions of the <strong>National Sports Act, 2023</strong>, the national sports organisation:
    </div>

    <div class="fed-name">${lic.federation_name || federationData.value?.name || 'National Sports Federation'}</div>
    <div class="reg-tag">Registration Number: <strong>${lic.federation_reg_no || federationData.value?.ncs_registration_number || 'NCS-STATUTORY'}</strong> · Category: <strong>${lic.category || 'National Sports Federation'}</strong></div>

    <div class="meta-box">
      <div class="meta-item">
        <span>License Number</span>
        <strong>${lic.license_number}</strong>
      </div>
      <div class="meta-item">
        <span>Issue Date</span>
        <strong>${issueDateStr}</strong>
      </div>
      <div class="meta-item">
        <span>Valid Until</span>
        <strong>${expiryDateStr}</strong>
      </div>
      <div class="meta-item">
        <span>Status</span>
        <strong style="color: ${lic.status === 'REVOKED' ? '#dc2626' : '#16a34a'};">${lic.status}</strong>
      </div>
    </div>

    <div class="conditions">
      ${lic.conditions || 'Granted subject to compliance with the National Sports Act 2023, anti-doping protocols, and financial transparency regulations.'}
    </div>

    <div class="signatures">
      <div class="sig-line">
        General Secretary<br>
        <span class="sig-title">National Council of Sports</span>
      </div>
      <div class="sig-line">
        Chairman / Board President<br>
        <span class="sig-title">National Council of Sports</span>
      </div>
    </div>
  </div>

  <script>
    window.onload = function() {
      setTimeout(function() { window.print(); }, 400);
    };
  <\/script>
</body>
</html>`

  printWin.document.write(html)
  printWin.document.close()
}

function getFileIcon(category) {
  switch (category) {
    case 'Federation Recognition': return 'icofont-certificate'
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
    case 'Federation Recognition': return 'bg-warning text-dark'
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
  
  recordActivityEvent({
    type: 'download_document',
    action: `Downloaded ${file.title || 'Document'}`,
    page_name: 'Credential Wallet',
    section: section.value || 'my-files',
    resource: file.category || 'Documents',
    label: `Download ${file.fileType || 'Certificate'}`,
    metadata: {
      document_title: file.title,
      document_number: file.number,
      category: file.category,
    },
  })

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
  const hasAthleteRole = roles.includes('athlete') || roles.includes('role_athlete')

  try {
    const athletesRes = await listNsmisDomain('athletes', { search: profile.email || profile.id })
    const athletesList = asList(athletesRes)
    let match = athletesList.find(ath => ath.user_id === profile.id || (ath.email_address && String(ath.email_address).toLowerCase() === String(profile.email || '').toLowerCase()))
    
    // If not matched by user_id or email, try matching by NIN
    if (!match && profile.nin) {
      const ninRes = await listNsmisDomain('athletes', { search: profile.nin })
      const ninList = asList(ninRes)
      match = ninList.find(ath => String(ath.national_id_passport || '').trim().toLowerCase() === String(profile.nin || '').trim().toLowerCase())
    }

    if (match || hasAthleteRole || profile.athlete_profile) {
      isAthlete.value = true
      if (match || profile.athlete_profile) {
        athleteData.value = match || profile.athlete_profile
        const athleteId = (match || profile.athlete_profile).id
        const [medicalRes, safeguardingRes, antidopingRes, nationalTeamRes, resultsRes, medalsRes, talentRes, clubsRes] = await Promise.allSettled([
          listNsmisDomain('medical-records', { search: athleteId }),
          listNsmisDomain('safeguarding-records', { search: athleteId }),
          listNsmisDomain('anti-doping', { search: athleteId }),
          listNsmisDomain('national-team', { search: athleteId }),
          listNsmisDomain('competition-results', { search: athleteId }),
          listNsmisDomain('medals', { search: athleteId }),
          listNsmisDomain('talent', { search: athleteId }),
          listNsmisDomain('clubs', { per_page: 200 })
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
        if (talentRes.status === 'fulfilled') {
          const talItems = asList(talentRes.value)
          athleteTalent.value = talItems.filter(r => r.athlete_id === athleteId || (r.athlete_name && (match || profile.athlete_profile).full_name && r.athlete_name.toLowerCase() === (match || profile.athlete_profile).full_name.toLowerCase()))
        }
        if (clubsRes.status === 'fulfilled') {
          const clubsList = asList(clubsRes.value)
          const clubName = (match || profile.athlete_profile).club
          if (clubName) {
            athleteClubDetails.value = clubsList.find(c => c.name === clubName || c.acronym === clubName || c.id === clubName) || null
          }
        }
      }
    }
  } catch (e) {
    console.warn('Could not load athlete records:', e)
  }

  // Coaches lookup
  try {
    const coachesRes = await listNsmisDomain('coaches', { search: profile.id || profile.email })
    const coachesList = asList(coachesRes)
    const matchCoach = coachesList.find(c => c.user_id === profile.id || (c.email && String(c.email).toLowerCase() === String(profile.email || '').toLowerCase())) || profile.coach_profile
    if (matchCoach) {
      isCoach.value = true
      coachData.value = matchCoach
    }
  } catch (e) {}

    // Officials lookup
    try {
      const officialsRes = await listNsmisDomain('technical-officials', { search: profile.id || profile.email })
      const officialsList = asList(officialsRes)
      const nameKey = fullName.value.toLowerCase().trim()
      const matchOfficial = officialsList.find(o => o.user_id === profile.id || (o.full_name && String(o.full_name).toLowerCase().trim() === nameKey)) || profile.official_profile
      if (matchOfficial) {
        isOfficial.value = true
        officialData.value = matchOfficial
      }
    } catch (e) {}

    // Federation Leadership & Executive lookup (President, General Secretary, Officials)
    try {
      const fedsRes = await listNsmisDomain('federations', { per_page: 200 })
      const fedsList = asList(fedsRes)
      const userEmail = String(profile.email || '').toLowerCase().trim()
      const userName = fullName.value.toLowerCase().trim()

      let matchFed = null
      let matchRole = ''

      for (const fed of fedsList) {
        if (fed.email && String(fed.email).toLowerCase().trim() === userEmail) {
          matchFed = fed
          matchRole = 'Federation Executive'
          break
        }
        if (fed.president && String(fed.president).toLowerCase().trim() === userName) {
          matchFed = fed
          matchRole = 'Federation President'
          break
        }
        if (fed.secretary && String(fed.secretary).toLowerCase().trim() === userName) {
          matchFed = fed
          matchRole = 'General Secretary'
          break
        }
      }

      if (!matchFed) {
        try {
          const officersRes = await listNsmisDomain('federation-officials', { per_page: 200 })
          const officersList = asList(officersRes)
          const officer = officersList.find(o => o.user_id === profile.id || (o.email && String(o.email).toLowerCase().trim() === userEmail) || (o.full_name && String(o.full_name).toLowerCase().trim() === userName))
          if (officer) {
            matchFed = fedsList.find(f => f.id === officer.federation_id) || { id: officer.federation_id, name: officer.federation_name }
            matchRole = officer.position || 'Federation Official'
          }
        } catch (e) {}
      }

      if (matchFed) {
        isFederationLeader.value = true
        federationLeaderRole.value = matchRole
        federationData.value = matchFed

        try {
          const lic = await getFederationActiveLicense(matchFed.id)
          if (lic) {
            federationLicenseData.value = lic
          }
        } catch (e) {}
      }
    } catch (e) {
      console.warn('Could not load federation leadership records:', e)
    }
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
    recordActivityEvent({
      type: 'revoke_session',
      action: 'Revoked Active Login Session',
      page_name: 'Account Settings',
      section: 'settings',
      resource: 'Active Sessions',
      label: 'Revoke Session',
      metadata: { session_id: sessionId }
    })
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
    recordActivityEvent({
      type: 'revoke_session',
      action: 'Signed Out All Other Devices',
      page_name: 'Account Settings',
      section: 'settings',
      resource: 'Active Sessions',
      label: 'Sign Out All Other Devices'
    })
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
    recordActivityEvent({
      type: 'update_preferences',
      action: 'Saved Notification & Privacy Preferences',
      page_name: 'Account Settings',
      section: 'settings',
      resource: 'Preferences',
      label: 'Save Preferences'
    })
    setTimeout(() => { preferencesSaved.value = false }, 3000)
  } catch (err) {
    error.value = 'Could not save preferences.'
  } finally {
    savingPreferences.value = false
  }
}

function parseAuditPayload(item) {
  if (!item) return {}
  if (item.payload_excerpt && typeof item.payload_excerpt === 'object') return item.payload_excerpt
  if (typeof item.payload_excerpt === 'string') {
    try { return JSON.parse(item.payload_excerpt) } catch {}
  }
  if (item.metadata && typeof item.metadata === 'object') return item.metadata
  return {}
}

function getActivityCategory(item) {
  if (!item) return 'SYSTEM'
  const action = String(item.action || '').toLowerCase()
  const endpoint = String(item.endpoint || '').toLowerCase()
  const eventType = String(item.event_type || '').toUpperCase()
  const meta = parseAuditPayload(item)
  const metaSec = String(meta.section || '').toLowerCase()

  if (eventType.includes('AUTH') || endpoint.includes('/auth') || action.includes('login') || action.includes('logout') || action.includes('sign in') || action.includes('sign out') || action.includes('password') || action.includes('2fa') || action.includes('security')) return 'AUTH'
  if (action.includes('profile') || action.includes('avatar') || action.includes('photo') || action.includes('preference') || metaSec.includes('profile') || metaSec.includes('settings') || endpoint.includes('/profile') || endpoint.includes('/users/me') || endpoint.includes('/security')) return 'PROFILE'
  if (action.includes('pay') || action.includes('momo') || action.includes('money') || action.includes('transaction') || endpoint.includes('/pay') || endpoint.includes('/transactions') || metaSec.includes('transactions')) return 'FINANCE'
  if (action.includes('application') || action.includes('draft') || action.includes('submission') || action.includes('wizard') || action.includes('service') || action.includes('download') || endpoint.includes('/forms') || endpoint.includes('/submissions') || endpoint.includes('/applications') || metaSec.includes('apply') || metaSec.includes('applications') || metaSec.includes('my-files')) return 'APPLICATION'
  if (eventType.includes('NAV') || action.includes('visited') || action.includes('navigate') || action.includes('viewed') || meta.type === 'navigation') return 'NAVIGATION'
  return 'SYSTEM'
}

function getActivityIcon(item) {
  const cat = getActivityCategory(item)
  switch (cat) {
    case 'AUTH': return 'icofont-key'
    case 'APPLICATION': return 'icofont-file-document'
    case 'FINANCE': return 'icofont-money'
    case 'PROFILE': return 'icofont-user-alt-7'
    case 'NAVIGATION': return 'icofont-compass'
    default: return 'icofont-history'
  }
}

function getActivityTone(item) {
  const cat = getActivityCategory(item)
  switch (cat) {
    case 'AUTH': return 'amber'
    case 'APPLICATION': return 'blue'
    case 'FINANCE': return 'green'
    case 'PROFILE': return 'cyan'
    case 'NAVIGATION': return 'purple'
    default: return 'secondary'
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

// --- Direct Application Payment Modal Functions ---
const paymentModalOpen = ref(false)
const paymentModalApp = ref(null)
const paymentModalMethod = ref('MOBILE_MONEY')
const paymentModalPhone = ref('')
const paymentModalPRN = ref('')
const paymentModalFile = ref(null)
const paymentModalSubmitting = ref(false)
const paymentModalError = ref('')
const paymentModalSuccess = ref('')
const paymentModalTracking = ref(false)
const paymentModalCountdown = ref(90)
const paymentModalPollingTimer = ref(null)
const paymentModalCountdownTimer = ref(null)
const paymentModalTxRef = ref('')

function canPayApplication(item) {
  if (item.source !== 'custom') return false
  if (item.payment_status === 'PAID') return false
  const fee = getAppPrice(item)
  return fee > 0
}

function getAppPrice(item) {
  if (!item) return 0
  if (item.price_ugx && Number(item.price_ugx) > 0) return Number(item.price_ugx)
  if (item.payment_amount_ugx && Number(item.payment_amount_ugx) > 0) return Number(item.payment_amount_ugx)
  const template = openForms.value.find(f => f.id === item.template_id)
  return template && Number(template.price_ugx) ? Number(template.price_ugx) : 0
}

function getAppAllowedMethods(item) {
  if (!item) return ['OVER_THE_COUNTER', 'MOBILE_MONEY']
  let methods = item.allowed_payment_methods
  if (!methods || (Array.isArray(methods) && !methods.length)) {
    const template = openForms.value.find(f => f.id === item.template_id)
    methods = template?.allowed_payment_methods
  }
  if (typeof methods === 'string') {
    try { methods = JSON.parse(methods) } catch { methods = [] }
  }
  if (!Array.isArray(methods) || !methods.length) {
    return ['OVER_THE_COUNTER', 'MOBILE_MONEY']
  }
  return methods
}

function getCarrierName(phone) {
  const clean = String(phone || '').replace(/[\s\-\+]/g, '')
  if (clean.startsWith('0111') || clean.startsWith('256111')) return 'Sandbox Test'
  if (clean.startsWith('077') || clean.startsWith('078') || clean.startsWith('076') || clean.startsWith('25677') || clean.startsWith('25678') || clean.startsWith('25676')) return 'MTN MoMo'
  if (clean.startsWith('070') || clean.startsWith('075') || clean.startsWith('074') || clean.startsWith('25670') || clean.startsWith('25675') || clean.startsWith('25674')) return 'Airtel Money'
  return ''
}

function openDirectPaymentModal(item) {
  paymentModalApp.value = item
  paymentModalError.value = ''
  paymentModalSuccess.value = ''
  paymentModalTracking.value = false
  paymentModalTxRef.value = ''
  paymentModalPRN.value = ''
  paymentModalFile.value = null
  paymentModalPhone.value = profile.phone || ''

  const allowed = getAppAllowedMethods(item)
  if (allowed.includes('MOBILE_MONEY')) {
    paymentModalMethod.value = 'MOBILE_MONEY'
  } else {
    paymentModalMethod.value = 'OVER_THE_COUNTER'
  }
  paymentModalOpen.value = true
}

function closeDirectPaymentModal() {
  if (paymentModalPollingTimer.value) {
    clearInterval(paymentModalPollingTimer.value)
    paymentModalPollingTimer.value = null
  }
  if (paymentModalCountdownTimer.value) {
    clearInterval(paymentModalCountdownTimer.value)
    paymentModalCountdownTimer.value = null
  }
  paymentModalOpen.value = false
  paymentModalApp.value = null
}

function onModalFileSelected(event) {
  const file = event.target?.files?.[0]
  if (file) {
    paymentModalFile.value = file
  }
}

async function submitModalMoMoPayment() {
  if (!paymentModalPhone.value.trim()) {
    paymentModalError.value = 'Please enter your Mobile Money phone number.'
    return
  }
  paymentModalError.value = ''
  paymentModalSubmitting.value = true
  try {
    const res = await portalInitiateMoMoPayment(paymentModalApp.value.id, {
      phone_number: paymentModalPhone.value.trim()
    })
    const data = res?.data || res || {}
    paymentModalTxRef.value = data.transaction_reference || ''
    paymentModalTracking.value = true
    paymentModalCountdown.value = 90

    if (paymentModalCountdownTimer.value) clearInterval(paymentModalCountdownTimer.value)
    paymentModalCountdownTimer.value = setInterval(() => {
      if (paymentModalCountdown.value > 0) {
        paymentModalCountdown.value--
      } else {
        clearInterval(paymentModalCountdownTimer.value)
      }
    }, 1000)

    if (paymentModalPollingTimer.value) clearInterval(paymentModalPollingTimer.value)
    paymentModalPollingTimer.value = setInterval(async () => {
      try {
        const pollRes = await portalGetPaymentStatus(paymentModalApp.value.id)
        const st = pollRes?.data || pollRes || {}
        if (st.status === 'SUCCESS' || st.payment_status === 'PAID') {
          clearInterval(paymentModalPollingTimer.value)
          clearInterval(paymentModalCountdownTimer.value)
          paymentModalPollingTimer.value = null
          paymentModalCountdownTimer.value = null
          paymentModalTracking.value = false
          paymentModalSuccess.value = `Payment of UGX ${formatMoney(getAppPrice(paymentModalApp.value))} confirmed! Reference: ${st.transaction_reference || paymentModalTxRef.value}`
          await loadPortal()
        } else if (st.status === 'FAILED') {
          clearInterval(paymentModalPollingTimer.value)
          clearInterval(paymentModalCountdownTimer.value)
          paymentModalPollingTimer.value = null
          paymentModalCountdownTimer.value = null
          paymentModalTracking.value = false
          paymentModalError.value = st.status_message || 'Payment was declined or failed. Please retry.'
        }
      } catch (e) {
        console.warn('Payment polling:', e)
      }
    }, 3000)
  } catch (err) {
    paymentModalError.value = err.response?.data?.message || err.response?.data?.error?.message || err.message || 'Failed to initiate Mobile Money payment'
  } finally {
    paymentModalSubmitting.value = false
  }
}

async function submitModalOTCProof() {
  if (!paymentModalPRN.value.trim() && !paymentModalFile.value) {
    paymentModalError.value = 'Please enter a URA PRN or choose a bank slip file to upload.'
    return
  }
  paymentModalError.value = ''
  paymentModalSubmitting.value = true
  try {
    const formData = new FormData()
    if (paymentModalPRN.value.trim()) formData.append('reference', paymentModalPRN.value.trim())
    if (paymentModalFile.value) formData.append('file', paymentModalFile.value)
    formData.append('amount', String(getAppPrice(paymentModalApp.value)))

    await portalUploadPaymentProof(paymentModalApp.value.id, formData)
    paymentModalSuccess.value = 'Bank deposit proof uploaded successfully for verification!'
    await loadPortal()
  } catch (err) {
    paymentModalError.value = err.response?.data?.message || err.response?.data?.error?.message || err.message || 'Failed to upload payment proof'
  } finally {
    paymentModalSubmitting.value = false
  }
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
    recordActivityEvent({
      type: 'update_profile',
      action: 'Updated Profile Information',
      page_name: 'Profile',
      section: 'profile',
      resource: 'Profile Details',
      label: 'Save profile',
      metadata: {
        first_name: profile.first_name,
        last_name: profile.last_name,
      }
    })
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
function formatUGX(value) { return 'UGX ' + formatMoney(value) }

const filteredAthleteMedals = computed(() => {
  if (selectedMedalFilter.value === 'ALL') return athleteMedals.value
  return athleteMedals.value.filter(m => String(m.medal_type || '').toUpperCase().includes(selectedMedalFilter.value))
})
function statusClass(status) { if (status === 'APPROVED' || status === 'COMPLETE') return 'green'; if (status === 'REJECTED') return 'red'; if (['NEEDS_INFORMATION', 'PENDING_PAYMENT'].includes(status)) return 'amber'; return 'blue' }
function paymentClass(status) { return status === 'PAID' || status === 'VERIFIED' ? 'green' : ['REJECTED', 'VERIFICATION_FAILED'].includes(status) ? 'red' : 'amber' }
function notificationIcon(item) { return item.icon_key ? `icofont-${item.icon_key}` : 'icofont-notification' }

function activityTitle(item) {
  if (!item) return 'Portal Activity'
  const action = String(item.action || '').trim()
  const resource = String(item.resource || '').trim()
  const resourceId = String(item.resource_id || '').trim()
  const method = String(item.method || '').toUpperCase()
  const endpoint = String(item.endpoint || '').toLowerCase()
  const eventType = String(item.event_type || '').toUpperCase()
  const meta = parseAuditPayload(item)
  const metaPage = meta.page || meta.page_name || ''
  const metaLabel = meta.label || ''
  const metaSection = meta.section || ''

  // 1. If action is already descriptive and human-readable
  if (action && !['ui:click', 'ui:navigate', 'ui:interaction', 'click', 'navigation', 'navigate', 'interaction', 'UI'].includes(action)) {
    if (action.startsWith('ui:')) {
      const rest = action.slice(3)
      return titleize(rest)
    }
    return action
  }

  // 2. High-priority auth & account endpoints
  if (eventType === 'AUTH_LOGIN' || endpoint.includes('/auth/login')) return 'Logged In to Account'
  if (eventType === 'AUTH_LOGOUT' || endpoint.includes('/auth/logout') || metaLabel.toLowerCase().includes('logout') || metaLabel.toLowerCase().includes('sign out')) return 'Logged Out of Portal'
  if (endpoint.includes('/auth/register')) return 'Registered Portal Account'
  if (endpoint.includes('/auth/password/reset')) return 'Requested Password Reset'
  if (endpoint.includes('/account/password') || metaLabel.toLowerCase().includes('password')) return 'Changed Account Password'
  if (endpoint.includes('/account/avatar') || metaLabel.toLowerCase().includes('photo') || metaLabel.toLowerCase().includes('avatar') || metaLabel.toLowerCase().includes('camera')) return 'Uploaded Profile Photo'
  if (endpoint.includes('/account/profile') || endpoint.includes('/users/me')) {
    if (['PUT', 'POST', 'PATCH'].includes(method)) return 'Updated Profile Information'
    return 'Viewed Profile Page'
  }

  // 3. Applications & Forms
  if (endpoint.includes('/portal/submissions') && endpoint.includes('/pay/momo')) return 'Initiated Mobile Money Payment'
  if (endpoint.includes('/portal/submissions') && endpoint.includes('/pay/proof')) return 'Uploaded Bank Payment Proof'
  if (endpoint.includes('/portal/submissions') && endpoint.includes('/submit')) return 'Submitted Form Application'
  if (endpoint.includes('/portal/submissions') && (method === 'POST' || method === 'PUT')) return 'Saved Application Draft'
  if (endpoint.includes('/portal/forms/open')) return 'Browsed Services Catalog'
  if (endpoint.includes('/portal/forms/') || endpoint.includes('/dashboard/apply/')) return 'Opened Application Form Wizard'
  if (endpoint.includes('/submissions') || endpoint.includes('/applications')) {
    if (method === 'POST') return 'Submitted Application'
    return 'Viewed My Applications List'
  }
  if (endpoint.includes('/media/upload')) return 'Uploaded Document Attachment'

  // 4. Navigation & Page / Section Visits
  if (action === 'ui:navigate' || eventType === 'UI_NAVIGATION' || meta.type === 'navigation') {
    if (metaPage) return `Visited ${metaPage}`
    if (metaSection) return `Visited ${titleize(metaSection)}`
    if (endpoint.includes('section=profile')) return 'Visited Profile'
    if (endpoint.includes('section=settings')) return 'Visited Account Settings'
    if (endpoint.includes('section=applications')) return 'Visited My Applications'
    if (endpoint.includes('section=apply')) return 'Visited Open Services'
    if (endpoint.includes('section=activities')) return 'Visited My Activities & Audit Logs'
    if (endpoint.includes('section=notifications')) return 'Visited Notification Center'
    if (endpoint.includes('section=messages')) return 'Visited Messages Center'
    if (endpoint.includes('section=transactions')) return 'Visited My Transactions'
    if (endpoint.includes('section=my-files')) return 'Visited Credential Wallet & Files'
    if (endpoint.includes('section=athlete-')) return `Visited Athlete ${titleize(endpoint.split('section=athlete-')[1] || 'Records')}`
    if (resource && resource !== 'Frontend Activity' && resource !== 'Portal Navigation') return `Visited ${resource}`
    return 'Visited Dashboard'
  }

  // 5. Clicks / Interactions
  if (action === 'ui:click' || eventType === 'UI_INTERACTION' || eventType === 'UI_CLICK' || meta.type === 'click') {
    const label = metaLabel || resourceId
    if (label) {
      const low = label.toLowerCase()
      if (low.includes('save') && (metaSection === 'profile' || endpoint.includes('profile'))) return 'Updated Profile Information'
      if (low.includes('save') && (metaSection === 'settings' || endpoint.includes('settings'))) return 'Updated Account Settings'
      if (low.includes('password')) return 'Updated Password'
      if (low.includes('2fa') || low.includes('two-factor')) return 'Configured Two-Factor Authentication'
      if (low.includes('submit')) return 'Submitted Form Application'
      if (low.includes('draft')) return 'Saved Application Draft'
      if (low.includes('pay') || low.includes('momo') || low.includes('mobile money')) return 'Initiated Mobile Money Payment'
      if (low.includes('upload')) return 'Uploaded File'
      if (low.includes('download')) return `Downloaded ${label.replace(/^Download\s+/i, '') || 'Document'}`
      if (low.includes('refresh')) return `Refreshed ${metaPage || titleize(metaSection) || 'Data'}`
      if (low.includes('replace')) return 'Replaced Uploaded File'
      if (low.includes('remove')) return 'Removed Uploaded File'
      return `Clicked "${label}" on ${metaPage || titleize(metaSection) || 'Dashboard'}`
    }
    if (metaPage) return `Interacted with ${metaPage}`
  }

  if (resource && resource !== 'Frontend Activity') return `${resource} Activity`
  if (item.event_type) return titleize(item.event_type)
  return 'Portal Activity'
}

function activityDescription(item) {
  if (!item) return ''
  const meta = parseAuditPayload(item)
  const title = activityTitle(item)
  
  let detail = ''
  if (meta.from_page && meta.page_name && meta.from_page !== meta.page_name) {
    detail = `Navigated from ${meta.from_page} to ${meta.page_name}`
  } else if (meta.label && !title.includes(meta.label)) {
    detail = `Action: "${meta.label}"`
  } else if (meta.page_name) {
    detail = `Section: ${meta.page_name}`
  } else if (meta.section) {
    detail = `Section: ${titleize(meta.section)}`
  }

  const loc = formatLocationPill(item)
  const dev = formatDevicePill(item)

  const parts = []
  if (detail) parts.push(detail)
  if (loc) parts.push(`From ${loc}`)
  if (dev && dev !== 'Web Client') parts.push(`via ${dev}`)

  return parts.length ? parts.join(' • ') : 'Authenticated user portal action recorded securely.'
}

function formatLocationPill(item) {
  const city = String(item?.geo_city || '').trim()
  const country = String(item?.geo_country || '').trim()
  if (city && country && country.toLowerCase() !== 'unknown') {
    const displayCity = city.toLowerCase() === 'internal' ? 'Kampala' : city
    return `${displayCity}, ${country}`
  }
  if (country && country.toLowerCase() !== 'unknown') return country
  const ip = String(item?.ip_address || '').trim()
  if (ip.startsWith('41.') || ip.startsWith('197.') || ip.startsWith('154.') || ip.startsWith('102.') || ip.startsWith('105.')) {
    return 'Kampala, Uganda'
  }
  return ''
}

function formatDevicePill(item) {
  const os = String(item?.os_name || '').trim()
  const browser = String(item?.browser || '').trim()
  const validOs = os && !['unknown', 'other'].includes(os.toLowerCase()) ? os : ''
  const validBrowser = browser && !['unknown', 'other'].includes(browser.toLowerCase()) ? browser : ''
  if (validBrowser && validOs) return `${validBrowser} on ${validOs}`
  if (validBrowser) return validBrowser
  if (validOs) return validOs
  return 'Web Client'
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
    
    recordActivityEvent({
      type: 'upload_avatar',
      action: 'Uploaded Profile Photo',
      page_name: 'Profile',
      section: 'profile',
      resource: 'Profile Details',
      label: 'Upload Photo',
      metadata: { file_name: file.name }
    })

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
    recordActivityEvent({
      type: 'change_password',
      action: 'Changed Account Password',
      page_name: 'Account Settings',
      section: 'settings',
      resource: 'Security Settings',
      label: 'Update Password',
    })
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
      recordActivityEvent({
        type: 'security_2fa',
        action: 'Disabled Two-Factor Authentication (2FA)',
        page_name: 'Account Settings',
        section: 'settings',
        resource: 'Security Settings',
        label: 'Disable 2FA',
      })
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
    recordActivityEvent({
      type: 'security_2fa',
      action: 'Enabled Two-Factor Authentication (2FA)',
      page_name: 'Account Settings',
      section: 'settings',
      resource: 'Security Settings',
      label: 'Activate 2FA',
    })
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
  width: 250px;
  background: #ffffff;
  border-right: 1px solid #e8edf2;
  z-index: 850;
  display: flex;
  flex-direction: column;
  transition: width 0.25s cubic-bezier(0.4, 0, 0.2, 1), transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  overflow-y: auto;
  overflow-x: hidden;
  padding-top: 60px;
}

/* ── Collapsed Sidebar (sidebar-mini: Icons only, titles on hover) ─── */
.sidebar-mini .main-sidebar {
  width: 70px !important;
  overflow: visible !important;
}
.sidebar-mini .main-sidebar #sidebar-wrapper {
  overflow: visible !important;
}
.sidebar-mini .main-sidebar .sidebar-brand {
  padding: 10px 6px;
  justify-content: center;
}
.sidebar-mini .main-sidebar .sidebar-brand img {
  max-width: 38px !important;
  max-height: 38px !important;
}
.sidebar-mini .main-sidebar .btn-close-sidebar {
  display: none !important;
}
.sidebar-mini .main-sidebar .sidebar-user {
  padding: 12px 6px;
  justify-content: center;
  gap: 0;
}
.sidebar-mini .main-sidebar .sidebar-user > div,
.sidebar-mini .main-sidebar .sidebar-menu .menu-header,
.sidebar-mini .main-sidebar .portal-nav-badge {
  display: none !important;
}
.sidebar-mini .main-sidebar .sidebar-menu {
  padding: 8px 6px;
  overflow: visible !important;
}
.sidebar-mini .main-sidebar .sidebar-menu li {
  position: relative;
  margin-bottom: 4px;
}
.sidebar-mini .main-sidebar .sidebar-menu li button.nav-link,
.sidebar-mini .main-sidebar .sidebar-menu li a.nav-link {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px 0 !important;
  width: 100% !important;
  border-radius: 10px;
  position: relative;
}
.sidebar-mini .main-sidebar .sidebar-menu li button.nav-link i,
.sidebar-mini .main-sidebar .sidebar-menu li a.nav-link i {
  margin: 0 !important;
  font-size: 20px !important;
  width: 32px;
  text-align: center;
}

/* Tooltip on Hover for Collapsed Menu Items */
.sidebar-mini .main-sidebar .sidebar-menu li button.nav-link span,
.sidebar-mini .main-sidebar .sidebar-menu li a.nav-link span {
  display: none;
  position: absolute;
  left: calc(100% + 12px);
  top: 50%;
  transform: translateY(-50%);
  background: #0f172a;
  color: #f8fafc;
  font-size: 13px;
  font-weight: 600;
  padding: 8px 14px;
  border-radius: 8px;
  white-space: nowrap;
  box-shadow: 0 10px 25px -3px rgba(0, 0, 0, 0.3), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  pointer-events: none;
  z-index: 99999;
  letter-spacing: 0.2px;
  line-height: 1.2;
}
.sidebar-mini .main-sidebar .sidebar-menu li button.nav-link span::before,
.sidebar-mini .main-sidebar .sidebar-menu li a.nav-link span::before {
  content: '';
  position: absolute;
  right: 100%;
  top: 50%;
  transform: translateY(-50%);
  border-width: 6px;
  border-style: solid;
  border-color: transparent #0f172a transparent transparent;
}
.sidebar-mini .main-sidebar .sidebar-menu li:hover > button.nav-link span,
.sidebar-mini .main-sidebar .sidebar-menu li:hover > a.nav-link span {
  display: block !important;
  animation: flyoutSlide 0.16s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
@keyframes flyoutSlide {
  from {
    opacity: 0;
    transform: translateY(-50%) translateX(-6px);
  }
  to {
    opacity: 1;
    transform: translateY(-50%) translateX(0);
  }
}

.sidebar-brand {
  display: flex;
  align-items: center;
  justify-content: space-between;
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

.btn-close-sidebar {
  display: none;
  background: #f1f5f9;
  border: none;
  color: #475569;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  cursor: pointer;
  margin-left: auto;
  transition: all 0.15s;
}
.btn-close-sidebar:hover {
  background: #e2e8f0;
  color: #0f172a;
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
  padding: 10px 16px;
  border: none !important;
  border-radius: 8px;
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
  font-size: 17px;
  width: 22px;
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
/* ── Mobile Sidebar & Scrim Base ───────────────────────────────────── */
.mobile-sidebar-open { transform: translateX(0); }
.sidebar-scrim {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.6);
  backdrop-filter: blur(3px);
  z-index: 1040;
}

/* ── Main Content Area Base (Desktop) ──────────────────────────────── */
.main-content {
  margin-top: 60px;
  margin-left: 250px;
  padding: 24px 30px 48px;
  width: calc(100% - 250px);
  max-width: calc(100% - 250px);
  min-height: calc(100vh - 60px);
  box-sizing: border-box !important;
  transition: margin-left 0.25s cubic-bezier(0.4, 0, 0.2, 1), width 0.25s cubic-bezier(0.4, 0, 0.2, 1), max-width 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  overflow-x: hidden;
}
.sidebar-mini .main-content {
  margin-left: 70px;
  width: calc(100% - 70px);
  max-width: calc(100% - 70px);
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


/* =====================================================================
   COMPREHENSIVE PORTAL DARK MODE ENGINE
   Works dynamically across all sections, cards, forms, tables & components
   ===================================================================== */
:global(.dark) .user-portal,
:global([data-theme="dark"]) .user-portal {
  background-color: #0f172a !important;
  color: #e2e8f0 !important;
}

:global(.dark) .navbar-bg,
:global(.dark) .main-navbar,
:global([data-theme="dark"]) .navbar-bg,
:global([data-theme="dark"]) .main-navbar {
  background-color: #1e293b !important;
  border-bottom: 1px solid #334155 !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4) !important;
}

:global(.dark) .main-sidebar,
:global([data-theme="dark"]) .main-sidebar {
  background-color: #1e293b !important;
  border-right: 1px solid #334155 !important;
}

:global(.dark) .sidebar-brand,
:global(.dark) .sidebar-user,
:global([data-theme="dark"]) .sidebar-brand,
:global([data-theme="dark"]) .sidebar-user {
  border-bottom-color: #334155 !important;
}

:global(.dark) .sidebar-user strong,
:global([data-theme="dark"]) .sidebar-user strong {
  color: #f8fafc !important;
}

:global(.dark) .sidebar-user span,
:global([data-theme="dark"]) .sidebar-user span {
  color: #94a3b8 !important;
}

:global(.dark) .sidebar-menu .menu-header,
:global([data-theme="dark"]) .sidebar-menu .menu-header {
  color: #64748b !important;
}

:global(.dark) .sidebar-menu li a,
:global(.dark) .sidebar-menu li button.nav-link,
:global([data-theme="dark"]) .sidebar-menu li a,
:global([data-theme="dark"]) .sidebar-menu li button.nav-link {
  color: #cbd5e1 !important;
}

:global(.dark) .sidebar-menu li.active > a,
:global(.dark) .sidebar-menu li.active > button,
:global(.dark) .sidebar-menu li.active > button.nav-link,
:global([data-theme="dark"]) .sidebar-menu li.active > a,
:global([data-theme="dark"]) .sidebar-menu li.active > button,
:global([data-theme="dark"]) .sidebar-menu li.active > button.nav-link {
  background-color: #0f172a !important;
  color: #93c5fd !important;
}

:global(.dark) .sidebar-menu li.active > a i,
:global(.dark) .sidebar-menu li.active > button i,
:global([data-theme="dark"]) .sidebar-menu li.active > a i,
:global([data-theme="dark"]) .sidebar-menu li.active > button i {
  color: #93c5fd !important;
}

:global(.dark) .sidebar-menu li a:hover,
:global(.dark) .sidebar-menu li button:hover,
:global([data-theme="dark"]) .sidebar-menu li a:hover,
:global([data-theme="dark"]) .sidebar-menu li button:hover {
  background-color: #334155 !important;
  color: #f8fafc !important;
}

/* Page Headers & Titles */
:global(.dark) .section-header h1,
:global(.dark) .page-heading h1,
:global([data-theme="dark"]) .section-header h1,
:global([data-theme="dark"]) .page-heading h1 {
  color: #f8fafc !important;
}

:global(.dark) .page-heading p,
:global([data-theme="dark"]) .page-heading p {
  color: #93c5fd !important;
}

:global(.dark) .page-heading span,
:global([data-theme="dark"]) .page-heading span {
  color: #94a3b8 !important;
}

/* Cards, Panels & KPI Containers */
:global(.dark) .card,
:global(.dark) .user-kpis article,
:global(.dark) .activity-kpis article,
:global(.dark) .feed-list article,
:global(.dark) .data-table,
:global(.dark) .profile-summary,
:global(.dark) .profile-form,
:global(.dark) .settings-card,
:global(.dark) .empty-state,
:global(.dark) .credential-card,
:global(.dark) .athlete-detail-card,
:global([data-theme="dark"]) .card,
:global([data-theme="dark"]) .user-kpis article,
:global([data-theme="dark"]) .activity-kpis article,
:global([data-theme="dark"]) .feed-list article,
:global([data-theme="dark"]) .data-table,
:global([data-theme="dark"]) .profile-summary,
:global([data-theme="dark"]) .profile-form,
:global([data-theme="dark"]) .settings-card,
:global([data-theme="dark"]) .empty-state,
:global([data-theme="dark"]) .credential-card,
:global([data-theme="dark"]) .athlete-detail-card {
  background-color: #1e293b !important;
  border-color: #334155 !important;
  color: #e2e8f0 !important;
  box-shadow: 0 4px 25px rgba(0, 0, 0, 0.3) !important;
}

:global(.dark) .user-kpis article strong,
:global(.dark) .activity-kpis article strong,
:global(.dark) .credential-title,
:global(.dark) .credential-number,
:global(.dark) .empty-title,
:global(.dark) .sessions-heading,
:global(.dark) .section-subheading,
:global([data-theme="dark"]) .user-kpis article strong,
:global([data-theme="dark"]) .activity-kpis article strong,
:global([data-theme="dark"]) .credential-title,
:global([data-theme="dark"]) .credential-number,
:global([data-theme="dark"]) .empty-title,
:global([data-theme="dark"]) .sessions-heading,
:global([data-theme="dark"]) .section-subheading {
  color: #f8fafc !important;
}

:global(.dark) .user-kpis article small,
:global(.dark) .activity-kpis article small,
:global([data-theme="dark"]) .user-kpis article small,
:global([data-theme="dark"]) .activity-kpis article small {
  color: #94a3b8 !important;
}

:global(.dark) .user-kpis article p,
:global(.dark) .activity-kpis article p,
:global([data-theme="dark"]) .user-kpis article p,
:global([data-theme="dark"]) .activity-kpis article p {
  color: #cbd5e1 !important;
}

/* Toolbars, Inputs, Selects & Searches */
:global(.dark) .list-toolbar,
:global([data-theme="dark"]) .list-toolbar {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(.dark) .list-toolbar input[type="search"],
:global(.dark) .list-toolbar input[type="text"],
:global(.dark) .list-toolbar select,
:global(.dark) .profile-form input,
:global(.dark) .profile-form select,
:global(.dark) .profile-form textarea,
:global(.dark) .settings-form input,
:global(.dark) .settings-form select,
:global(.dark) .settings-form textarea,
:global([data-theme="dark"]) .list-toolbar input[type="search"],
:global([data-theme="dark"]) .list-toolbar input[type="text"],
:global([data-theme="dark"]) .list-toolbar select,
:global([data-theme="dark"]) .profile-form input,
:global([data-theme="dark"]) .profile-form select,
:global([data-theme="dark"]) .profile-form textarea,
:global([data-theme="dark"]) .settings-form input,
:global([data-theme="dark"]) .settings-form select,
:global([data-theme="dark"]) .settings-form textarea {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #f8fafc !important;
}

:global(.dark) .list-toolbar select,
:global([data-theme="dark"]) .list-toolbar select {
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%2394a3b8' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e") !important;
}

:global(.dark) .profile-form label,
:global(.dark) .settings-form label,
:global([data-theme="dark"]) .profile-form label,
:global([data-theme="dark"]) .settings-form label {
  color: #cbd5e1 !important;
}

:global(.dark) .profile-form input:disabled,
:global([data-theme="dark"]) .profile-form input:disabled {
  background-color: #1e293b !important;
  color: #64748b !important;
}

/* Data Tables */
:global(.dark) .data-table,
:global([data-theme="dark"]) .data-table {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(.dark) .data-table th,
:global([data-theme="dark"]) .data-table th {
  background-color: #0f172a !important;
  color: #94a3b8 !important;
  border-bottom-color: #334155 !important;
}

:global(.dark) .data-table td,
:global([data-theme="dark"]) .data-table td {
  color: #cbd5e1 !important;
  border-bottom-color: #334155 !important;
}

:global(.dark) .data-table tr:hover td,
:global([data-theme="dark"]) .data-table tr:hover td {
  background-color: #243044 !important;
}

:global(.dark) .data-table td strong,
:global([data-theme="dark"]) .data-table td strong {
  color: #f8fafc !important;
}

:global(.dark) .table-actions button,
:global([data-theme="dark"]) .table-actions button {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #cbd5e1 !important;
}

:global(.dark) .table-actions button:hover,
:global([data-theme="dark"]) .table-actions button:hover {
  background-color: #6777ef !important;
  border-color: #6777ef !important;
  color: #fff !important;
}

/* Feed Lists & Notifications & Messages */
:global(.dark) .feed-list article,
:global([data-theme="dark"]) .feed-list article {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(.dark) .feed-list article strong,
:global([data-theme="dark"]) .feed-list article strong {
  color: #f8fafc !important;
}

:global(.dark) .feed-list article p,
:global([data-theme="dark"]) .feed-list article p {
  color: #cbd5e1 !important;
}

:global(.dark) .feed-list article small,
:global([data-theme="dark"]) .feed-list article small {
  color: #94a3b8 !important;
}

:global(.dark) .feed-list article.unread,
:global([data-theme="dark"]) .feed-list article.unread {
  background-color: #1e3a8a33 !important;
  border-color: #3b82f6 !important;
}

/* Transactions */
:global(.dark) .transaction-summary,
:global([data-theme="dark"]) .transaction-summary {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(.dark) .transaction-summary span,
:global([data-theme="dark"]) .transaction-summary span {
  color: #94a3b8 !important;
}

:global(.dark) .transaction-summary strong,
:global([data-theme="dark"]) .transaction-summary strong {
  color: #f8fafc !important;
}

/* Activities & Timeline Feed */
:global(.dark) .activity-feed-card,
:global([data-theme="dark"]) .activity-feed-card {
  background-color: #1e293b !important;
  border-color: #334155 !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.25) !important;
}

:global(.dark) .activity-title,
:global([data-theme="dark"]) .activity-title {
  color: #f8fafc !important;
}

:global(.dark) .activity-description,
:global([data-theme="dark"]) .activity-description {
  color: #cbd5e1 !important;
}

:global(.dark) .meta-pill,
:global([data-theme="dark"]) .meta-pill {
  background-color: #0f172a !important;
  border-color: #334155 !important;
  color: #cbd5e1 !important;
}

:global(.dark) .meta-pill strong,
:global([data-theme="dark"]) .meta-pill strong {
  color: #f8fafc !important;
}

:global(.dark) .category-pill,
:global([data-theme="dark"]) .category-pill {
  background-color: #1e293b !important;
  border-color: #334155 !important;
  color: #cbd5e1 !important;
}

:global(.dark) .category-pill:hover,
:global([data-theme="dark"]) .category-pill:hover {
  background-color: #334155 !important;
  color: #f8fafc !important;
}

:global(.dark) .category-pill.active,
:global([data-theme="dark"]) .category-pill.active {
  background-color: #6777ef !important;
  color: #fff !important;
  border-color: #6777ef !important;
}

:global(.dark) .activity-pagination-bar,
:global([data-theme="dark"]) .activity-pagination-bar {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(.dark) .pagination-info,
:global([data-theme="dark"]) .pagination-info {
  color: #94a3b8 !important;
}

/* Settings, Sessions & Preferences */
:global(.dark) .settings-nav-tabs,
:global([data-theme="dark"]) .settings-nav-tabs {
  border-bottom-color: #334155 !important;
}

:global(.dark) .settings-tab-btn,
:global([data-theme="dark"]) .settings-tab-btn {
  color: #94a3b8 !important;
}

:global(.dark) .settings-tab-btn:hover,
:global([data-theme="dark"]) .settings-tab-btn:hover {
  background-color: #334155 !important;
  color: #f8fafc !important;
}

:global(.dark) .settings-tab-btn.active,
:global([data-theme="dark"]) .settings-tab-btn.active {
  background-color: #6777ef !important;
  color: #fff !important;
}

:global(.dark) .settings-card-header,
:global([data-theme="dark"]) .settings-card-header {
  border-bottom-color: #334155 !important;
}

:global(.dark) .settings-card-header h3,
:global([data-theme="dark"]) .settings-card-header h3 {
  color: #f8fafc !important;
}

:global(.dark) .settings-card-header p,
:global([data-theme="dark"]) .settings-card-header p {
  color: #94a3b8 !important;
}

:global(.dark) .security-tips-card,
:global([data-theme="dark"]) .security-tips-card {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(.dark) .security-tips-list strong,
:global([data-theme="dark"]) .security-tips-list strong {
  color: #f8fafc !important;
}

:global(.dark) .security-tips-list p,
:global([data-theme="dark"]) .security-tips-list p {
  color: #94a3b8 !important;
}

:global(.dark) .session-item-card,
:global([data-theme="dark"]) .session-item-card {
  background-color: #0f172a !important;
  border-color: #334155 !important;
}

:global(.dark) .session-item-card.current-session-card,
:global([data-theme="dark"]) .session-item-card.current-session-card {
  background-color: #1e293b !important;
  border-color: #16a34a !important;
}

:global(.dark) .session-device-name,
:global([data-theme="dark"]) .session-device-name {
  color: #f8fafc !important;
}

:global(.dark) .meta-item,
:global([data-theme="dark"]) .meta-item {
  color: #94a3b8 !important;
}

:global(.dark) .session-device-icon,
:global([data-theme="dark"]) .session-device-icon {
  background-color: #1e293b !important;
  color: #cbd5e1 !important;
}

:global(.dark) .preference-item,
:global([data-theme="dark"]) .preference-item {
  background-color: #0f172a !important;
  border-color: #334155 !important;
}

:global(.dark) .preference-item:hover,
:global([data-theme="dark"]) .preference-item:hover {
  background-color: #1e293b !important;
}

:global(.dark) .preference-text strong,
:global([data-theme="dark"]) .preference-text strong {
  color: #f8fafc !important;
}

:global(.dark) .preference-text p,
:global([data-theme="dark"]) .preference-text p {
  color: #94a3b8 !important;
}

:global(.dark) .twofa-status-banner.twofa-active,
:global([data-theme="dark"]) .twofa-status-banner.twofa-active {
  background-color: #052e16 !important;
  border-color: #166534 !important;
}

:global(.dark) .twofa-status-banner.twofa-inactive,
:global([data-theme="dark"]) .twofa-status-banner.twofa-inactive {
  background-color: #451a03 !important;
  border-color: #9a3412 !important;
}

/* Athlete Registry Extension */
:global(.dark) .athlete-detail-card,
:global([data-theme="dark"]) .athlete-detail-card {
  background-color: #1e293b !important;
  border-color: #334155 !important;
}

:global(.dark) .athlete-detail-card header h3,
:global([data-theme="dark"]) .athlete-detail-card header h3 {
  color: #f8fafc !important;
}

:global(.dark) .athlete-detail-card ul li,
:global([data-theme="dark"]) .athlete-detail-card ul li {
  border-bottom-color: #334155 !important;
}

:global(.dark) .athlete-detail-card ul li span,
:global([data-theme="dark"]) .athlete-detail-card ul li span {
  color: #94a3b8 !important;
}

:global(.dark) .athlete-detail-card ul li strong,
:global([data-theme="dark"]) .athlete-detail-card ul li strong {
  color: #f8fafc !important;
}

/* Dropdown Menu */
:global(.dark) .dropdown-menu,
:global([data-theme="dark"]) .dropdown-menu {
  background-color: #1e293b !important;
  border-color: #334155 !important;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.4) !important;
}

:global(.dark) .dropdown-title,
:global([data-theme="dark"]) .dropdown-title {
  color: #64748b !important;
}

:global(.dark) .dropdown-item,
:global([data-theme="dark"]) .dropdown-item {
  color: #e2e8f0 !important;
}

:global(.dark) .dropdown-item:hover,
:global([data-theme="dark"]) .dropdown-item:hover {
  background-color: #334155 !important;
  color: #f8fafc !important;
}

:global(.dark) .dropdown-divider,
:global([data-theme="dark"]) .dropdown-divider {
  background-color: #334155 !important;
}

/* Footer */
:global(.dark) footer.main-footer,
:global([data-theme="dark"]) footer.main-footer {
  background-color: #1e293b !important;
  border-top-color: #334155 !important;
  color: #94a3b8 !important;
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

/* ── Athlete Portal Dedicated Styling ───────────────────────────── */
.athlete-panel-card {
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s ease;
  border: 1px solid rgba(0, 0, 0, 0.06) !important;
}
.athlete-panel-card:hover {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08) !important;
}
.kpi-icon-circle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  font-size: 20px;
}
.athlete-detail-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.athlete-detail-list li {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 12px;
  border-bottom: 1px dashed rgba(0, 0, 0, 0.07);
  font-size: 13px;
}
.athlete-detail-list li:last-child {
  border-bottom: none;
  padding-bottom: 0;
}
.athlete-detail-list li span {
  color: #64748b;
  font-weight: 500;
}
.athlete-detail-list li strong {
  color: #1e293b;
  font-weight: 600;
  text-align: right;
}
.athlete-national-banner {
  background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
  border-radius: 14px;
  color: #fff;
  border-left: 5px solid #eab308 !important;
}
.national-flag-icon {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  background: rgba(234, 179, 8, 0.2);
  color: #eab308;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}
.club-logo-box {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: #e0e7ff;
  color: #4f46e5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
}
.bg-bronze {
  background-color: #cd7f32 !important;
}
.pathway-timeline {
  display: flex;
  flex-direction: column;
  gap: 20px;
  position: relative;
  padding-left: 20px;
}
.pathway-timeline::before {
  content: '';
  position: absolute;
  left: 31px;
  top: 10px;
  bottom: 20px;
  width: 2px;
  background: #e2e8f0;
}
.pathway-step {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  position: relative;
}
.step-dot {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #cbd5e1;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  z-index: 1;
  flex-shrink: 0;
}
.pathway-step.completed .step-dot {
  background: #16a34a;
}
.step-content h6 {
  margin: 0 0 4px 0;
  font-size: 14px;
  font-weight: 700;
  color: #1e293b;
}
.filter-pills {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.filter-pill {
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  color: #475569;
  cursor: pointer;
  transition: all 0.15s ease;
}
.filter-pill:hover {
  background: #e2e8f0;
}
.filter-pill.active {
  background: #6777ef;
  border-color: #6777ef;
  color: #fff;
}

/* ── Dark Mode for Athlete Registry ── */
:global(.dark) .athlete-panel-card {
  background: #1e293b;
  border-color: #334155 !important;
}
:global(.dark) .athlete-panel-card .card-header {
  background: #1e293b !important;
  border-color: #334155 !important;
}
:global(.dark) .athlete-detail-list li span {
  color: #94a3b8;
}
:global(.dark) .athlete-detail-list li strong {
  color: #f8fafc;
}
:global(.dark) .athlete-detail-list li {
  border-bottom-color: #334155;
}
:global(.dark) .club-logo-box {
  background: #312e81;
  color: #a5b4fc;
}
:global(.dark) .step-content h6 {
  color: #f8fafc;
}
:global(.dark) .filter-pill {
  background: #1e293b;
  border-color: #334155;
  color: #cbd5e1;
}
:global(.dark) .filter-pill:hover {
  background: #334155;
}
:global(.dark) .filter-pill.active {
  background: #6777ef;
  border-color: #6777ef;
  color: #fff;
}
:global(.dark) .pathway-timeline::before {
  background: #334155;
}
:global(.dark) .table thead.bg-light th {
  background: #0f172a !important;
  color: #cbd5e1 !important;
  border-color: #334155 !important;
}
:global(.dark) .table td {
  border-color: #334155 !important;
  color: #cbd5e1 !important;
}

/* Pay Now Action Buttons */
.btn-pay-action-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: #10b981;
  color: #ffffff;
  border: none;
  padding: 3px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 4px rgba(16, 185, 129, 0.2);
}
.btn-pay-action-pill:hover {
  background: #059669;
  transform: translateY(-1px);
  box-shadow: 0 4px 6px rgba(16, 185, 129, 0.3);
}
.btn-pay-icon {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s ease;
}
.btn-pay-icon:hover {
  background: #ecfdf5;
  color: #059669 !important;
}

/* Portal Payment Modal */
.portal-modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(4px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}
.portal-modal-card {
  background: #ffffff;
  border-radius: 16px;
  width: 100%;
  max-width: 520px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.2), 0 10px 10px -5px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
  overflow: hidden;
  animation: modalPop 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}
@keyframes modalPop {
  0% { opacity: 0; transform: scale(0.95) translateY(10px); }
  100% { opacity: 1; transform: scale(1) translateY(0); }
}
.portal-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid #f1f5f9;
  background: #fafbfc;
}
.modal-header-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: #eff6ff;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
}
.modal-title {
  margin: 0;
  font-size: 17px;
  font-weight: 700;
  color: #0f172a;
}
.modal-subtitle {
  margin: 0;
  font-size: 12px;
  color: #64748b;
}
.btn-close-modal {
  background: none;
  border: none;
  font-size: 20px;
  color: #94a3b8;
  cursor: pointer;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}
.btn-close-modal:hover {
  background: #f1f5f9;
  color: #0f172a;
}
.portal-modal-body {
  padding: 24px;
}

/* Modal Fee Banner */
.modal-fee-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 18px;
  background: linear-gradient(135deg, #f8fafc 0%, #edf2f7 100%);
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  margin-bottom: 20px;
}
.fee-meta small {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  color: #64748b;
  letter-spacing: 0.5px;
}
.fee-val {
  font-size: 20px;
  font-weight: 800;
  color: #047857;
}
.fee-ref small {
  font-size: 10px;
  text-transform: uppercase;
  color: #94a3b8;
  display: block;
}
.fee-ref strong {
  font-size: 12px;
  color: #334155;
}

/* Payment Method Tabs */
.payment-method-selector-tabs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}
.selector-tab-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: #f8fafc;
  border: 2px solid #e2e8f0;
  border-radius: 10px;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s;
}
.selector-tab-btn i {
  font-size: 24px;
  color: #64748b;
}
.selector-tab-btn strong {
  display: block;
  font-size: 13px;
  color: #1e293b;
}
.selector-tab-btn small {
  font-size: 10px;
  color: #64748b;
}
.selector-tab-btn.active {
  background: #eff6ff;
  border-color: #3b82f6;
}
.selector-tab-btn.active i {
  color: #2563eb;
}
.selector-tab-btn.active strong {
  color: #1d4ed8;
}

/* Carrier Tags */
.carrier-tag {
  font-size: 10px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 10px;
  text-transform: uppercase;
}
.carrier-tag.mtn {
  background: #fef08a;
  color: #854d0e;
}
.carrier-tag.airtel {
  background: #fee2e2;
  color: #991b1b;
}
.carrier-tag.sandbox {
  background: #dbeafe;
  color: #1e40af;
}

/* Alerts */
.modal-alert-box {
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
}
.modal-alert-box.alert-error {
  background: #fef2f2;
  color: #991b1b;
  border: 1px solid #fecaca;
}
.modal-alert-box.alert-success {
  background: #ecfdf5;
  color: #065f46;
  border: 1px solid #a7f3d0;
}

/* Sandbox Hints */
.sandbox-hints-box {
  background: #f8fafc;
  padding: 10px 12px;
  border-radius: 8px;
  border: 1px dashed #cbd5e1;
}
.hints-label {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  display: block;
  margin-bottom: 6px;
}
.btn-xs {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 6px;
}

/* USSD Tracking Spinner */
.ussd-tracking-card {
  background: #fafbfc;
  border: 1px solid #f1f5f9;
  border-radius: 12px;
}
.ussd-spinner-wrap {
  position: relative;
  width: 72px;
  height: 72px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
}
.pulsing-circle {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: rgba(37, 99, 235, 0.15);
  animation: pulseAnim 1.6s infinite ease-out;
}
@keyframes pulseAnim {
  0% { transform: scale(0.8); opacity: 0.8; }
  100% { transform: scale(1.4); opacity: 0; }
}
.ussd-center-icon {
  font-size: 32px;
  color: #2563eb;
  position: relative;
  z-index: 1;
}
.tracking-progress-box {
  display: inline-flex;
  align-items: center;
  background: #ffffff;
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  color: #334155;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

/* Dark Mode Overrides for Payment Modal */
:global(.dark) .portal-modal-card {
  background: #1e293b;
  border-color: #334155;
}
:global(.dark) .portal-modal-header {
  background: #0f172a;
  border-color: #334155;
}
:global(.dark) .modal-title {
  color: #f8fafc;
}
:global(.dark) .modal-subtitle {
  color: #94a3b8;
}
:global(.dark) .modal-fee-banner {
  background: #0f172a;
  border-color: #334155;
}
:global(.dark) .fee-val {
  color: #34d399;
}
:global(.dark) .selector-tab-btn {
  background: #0f172a;
  border-color: #334155;
}
:global(.dark) .selector-tab-btn strong {
  color: #f8fafc;
}
:global(.dark) .selector-tab-btn.active {
  background: #1e3a8a;
  border-color: #60a5fa;
}
:global(.dark) .selector-tab-btn.active strong {
  color: #bfdbfe;
}
:global(.dark) .sandbox-hints-box {
  background: #0f172a;
  border-color: #334155;
}
:global(.dark) .ussd-tracking-card {
  background: #0f172a;
  border-color: #334155;
}
/* ── Comprehensive Mobile & Tablet Responsive Overrides ───────────── */
@media (max-width: 991px) {
  .btn-close-sidebar {
    display: flex;
  }
  .main-sidebar {
    position: fixed !important;
    top: 0 !important;
    left: 0 !important;
    bottom: 0 !important;
    width: 280px !important;
    padding-top: 0 !important;
    transform: translateX(-100%);
    box-shadow: none;
    z-index: 1050;
  }
  .main-sidebar.mobile-sidebar-open {
    transform: translateX(0) !important;
    box-shadow: 0 0 40px rgba(0,0,0,0.3) !important;
  }

  .main-content,
  .sidebar-mini .main-content {
    margin-left: 0 !important;
    width: 100% !important;
    max-width: 100% !important;
    min-width: 100% !important;
    padding: 16px 14px 40px !important;
    box-sizing: border-box !important;
  }

  .main-navbar {
    padding: 0 12px;
  }
  .main-navbar .mr-3 {
    margin-right: 6px;
  }
  .portal-navbar-title {
    padding: 0 6px;
  }
  .portal-navbar-title small {
    display: none;
  }
  .portal-navbar-title strong {
    font-size: 14px;
    max-width: 140px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .section-header {
    margin-bottom: 16px !important;
  }
  .section-header h1 {
    font-size: 18px !important;
  }
  .cms-actions.otika-page-actions {
    margin-bottom: 16px !important;
  }
  .profile-layout {
    display: block !important;
  }
  .profile-layout .profile-grid-3 {
    display: flex !important;
    flex-direction: column !important;
    gap: 16px !important;
  }
  .profile-summary, .profile-form {
    width: 100% !important;
  }
  .settings-tab-content .row {
    margin: 0 !important;
  }
  .settings-tab-content .col-lg-7,
  .settings-tab-content .col-lg-5,
  .settings-tab-content .col-lg-6 {
    padding: 0 !important;
    margin-bottom: 16px;
    width: 100% !important;
    max-width: 100% !important;
  }
}

@media (max-width: 768px) {
  .user-kpis {
    grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)) !important;
    gap: 12px !important;
    margin-bottom: 16px !important;
  }
  .user-kpis article {
    padding: 14px 12px !important;
  }
  .page-heading {
    flex-direction: column !important;
    align-items: flex-start !important;
    gap: 12px !important;
    padding-bottom: 14px !important;
    margin-bottom: 16px !important;
  }
  .page-heading > div {
    width: 100%;
  }
  .page-heading h1 {
    font-size: 18px !important;
  }
  .page-heading button.primary-command,
  .page-heading button.secondary-command {
    width: 100% !important;
    justify-content: center !important;
  }
  .athlete-panel-card {
    margin-bottom: 16px !important;
  }
  .athlete-panel-card .card-header {
    padding: 12px 14px !important;
    flex-wrap: wrap !important;
    gap: 8px !important;
  }
  .athlete-panel-card .card-header h5 {
    font-size: 14px !important;
  }
  .table-responsive {
    -webkit-overflow-scrolling: touch !important;
    overflow-x: auto !important;
    width: 100% !important;
    border-radius: 0 0 10px 10px;
  }
  .table-responsive table {
    min-width: 580px !important;
    margin-bottom: 0 !important;
  }
  .table-responsive table th,
  .table-responsive table td {
    padding: 10px 12px !important;
    font-size: 12px !important;
    white-space: nowrap !important;
  }
  .credential-card .card-body {
    padding: 16px !important;
  }
  .credential-card .credential-title {
    font-size: 14px !important;
  }
}

@media (max-width: 576px) {
  .main-navbar {
    height: 54px !important;
    padding: 0 10px !important;
  }
  .main-content {
    margin-top: 54px !important;
    padding: 12px 10px 36px !important;
  }
  .main-navbar .nav-link-lg.cms-top-icon,
  .main-navbar button.cms-top-icon,
  .main-navbar .portal-top-action,
  .main-navbar .portal-theme-toggle {
    width: 36px !important;
    height: 36px !important;
    font-size: 18px !important;
    border-radius: 8px !important;
  }
  .portal-navbar-title {
    display: none !important;
  }
  .user-kpis {
    grid-template-columns: 1fr !important;
    gap: 10px !important;
  }
  .user-kpis article span {
    width: 40px !important;
    height: 40px !important;
    font-size: 18px !important;
  }
  .activity-feed-card {
    padding: 12px 10px !important;
  }
  .activity-header-row {
    flex-direction: column !important;
    align-items: flex-start !important;
    gap: 6px !important;
  }
  .activity-title-group {
    flex-wrap: wrap !important;
    gap: 6px !important;
  }
  .activity-title {
    font-size: 13px !important;
  }
  .activity-description {
    font-size: 12px !important;
  }
  .activity-meta-row {
    flex-direction: column !important;
    align-items: flex-start !important;
    gap: 4px !important;
  }
  .activity-meta-row .meta-pill {
    font-size: 10px !important;
    padding: 2px 8px !important;
  }
  .activity-category-pills {
    display: flex !important;
    overflow-x: auto !important;
    padding-bottom: 6px !important;
    -webkit-overflow-scrolling: touch !important;
    gap: 6px !important;
    margin-bottom: 12px !important;
    flex-wrap: nowrap !important;
  }
  .category-pill {
    flex-shrink: 0 !important;
    font-size: 11px !important;
    padding: 6px 10px !important;
    white-space: nowrap !important;
  }
  .activity-pagination-bar {
    flex-direction: column !important;
    gap: 10px !important;
    align-items: center !important;
    text-align: center !important;
  }
  .portal-modal-card {
    max-width: 96vw !important;
    margin: 8px auto !important;
    border-radius: 12px !important;
  }
  .portal-modal-header {
    padding: 14px 16px !important;
  }
  .modal-header-icon {
    width: 36px !important;
    height: 36px !important;
    font-size: 18px !important;
  }
  .modal-title {
    font-size: 15px !important;
  }
  .portal-modal-body {
    padding: 16px 14px !important;
  }
  .modal-fee-banner {
    padding: 10px 12px !important;
    flex-direction: column !important;
    align-items: flex-start !important;
    gap: 6px !important;
  }
  .fee-val {
    font-size: 18px !important;
  }
  .payment-method-selector-tabs {
    grid-template-columns: 1fr !important;
    gap: 8px !important;
  }
}

/* ── Dark Mode Tooltip Flyout Overrides ────────────────────────────── */
:global(.dark) .sidebar-mini .main-sidebar .sidebar-menu li button.nav-link span,
:global(.dark) .sidebar-mini .main-sidebar .sidebar-menu li a.nav-link span {
  background: #1e293b !important;
  color: #f1f5f9 !important;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.6) !important;
  border: 1px solid #334155 !important;
}
:global(.dark) .sidebar-mini .main-sidebar .sidebar-menu li button.nav-link span::before,
:global(.dark) .sidebar-mini .main-sidebar .sidebar-menu li a.nav-link span::before {
  border-color: transparent #1e293b transparent transparent !important;
}
</style>


