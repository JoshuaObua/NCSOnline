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
                <small>Ordinary User</small>
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
                <span>Applicant workspace</span>
              </div>
            </div>
            <ul class="sidebar-menu" aria-label="Applicant navigation">
              <li class="menu-header">Applicant Portal</li>
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

const router = useRouter()
const portalSectionIds = ['dashboard', 'apply', 'applications', 'activities', 'notifications', 'messages', 'transactions', 'profile']
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

ensureOtikaStyles()

const navigation = computed(() => [
  { id: 'dashboard', label: 'Dashboard', icon: 'icofont-dashboard-web' },
  { id: 'apply', label: 'Apply Now', icon: 'icofont-plus-circle', badge: openForms.value.length || '' },
  { id: 'applications', label: 'My Applications', icon: 'icofont-file-document' },
  { id: 'activities', label: 'My Activities', icon: 'icofont-history' },
  { id: 'notifications', label: 'Notifications', icon: 'icofont-notification', badge: unreadNotifications.value || '' },
  { id: 'messages', label: 'Messages', icon: 'icofont-envelope', badge: messageCount.value || '' },
  { id: 'transactions', label: 'My Transactions', icon: 'icofont-money' },
  { id: 'profile', label: 'My Profile', icon: 'icofont-user-alt-3' },
])
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
  if (results.some(item => item.status === 'rejected')) error.value = 'Some dashboard information could not be loaded. Refresh to try again.'
  
  try {
    const sec = await getMySecurity()
    twofaEnabled.value = sec?.twofa?.enabled ?? false
  } catch (e) {}

  if (profile && profile.email) {
    try {
      const athletesRes = await listNsmisDomain('athletes', { search: profile.email })
      const athletesList = asList(athletesRes)
      const match = athletesList.find(ath => String(ath.email_address || '').toLowerCase() === String(profile.email || '').toLowerCase())
      if (match) {
        isAthlete.value = true
        athleteData.value = match
        
        const athleteId = match.id
        const [medicalRes, safeguardingRes, antidopingRes, nationalTeamRes] = await Promise.allSettled([
          listNsmisDomain('medical-records', { search: athleteId }),
          listNsmisDomain('safeguarding-records', { search: athleteId }),
          listNsmisDomain('anti-doping', { search: athleteId }),
          listNsmisDomain('national-team', { search: athleteId })
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
      }
    } catch (e) {
      console.warn('Failed to load athlete context details:', e)
    }
  }

  loading.value = false
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
function asList(value) { const data = unwrap(value); return Array.isArray(data) ? data : [] }
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
/* Scoped base styles */
.user-portal {
  min-height: 100vh!important;
  background: #f4f6f9!important;
  color: #34395e!important;
  overflow-x: hidden!important;
}
.user-portal button {
  font-family: inherit;
}
.user-portal .main-wrapper {
  min-height: 100vh!important;
  display: flex!important;
  flex-direction: column!important;
  --portal-sidebar-width: 250px;
  --portal-navbar-height: 70px;
  --portal-page-pad: 24px;
}
.user-portal .main-wrapper.sidebar-mini {
  --portal-sidebar-width: 65px;
}

/* Navbar */
.user-portal .navbar-bg {
  position: fixed!important;
  top: 0!important;
  right: 0!important;
  left: var(--portal-sidebar-width)!important;
  z-index: 1030!important;
  height: var(--portal-navbar-height)!important;
  background: #fff!important;
  box-shadow: 0 4px 25px rgba(0,0,0,.08)!important;
  transition: left 0.3s ease!important;
}
.user-portal .main-navbar {
  position: fixed!important;
  top: 0!important;
  right: 0!important;
  left: var(--portal-sidebar-width)!important;
  z-index: 1040!important;
  min-height: var(--portal-navbar-height)!important;
  background: #fff!important;
  color: #111827!important;
  box-shadow: 0 4px 25px rgba(0,0,0,.08)!important;
  padding: 8px var(--portal-page-pad)!important;
  display: flex!important;
  align-items: center!important;
  justify-content: space-between!important;
  transition: left 0.3s ease!important;
}
.user-portal .main-navbar .nav-link,
.user-portal .main-navbar i {
  color: #111827!important;
}

/* Sidebar */
.user-portal .main-sidebar {
  position: fixed!important;
  left: 0!important;
  top: 0!important;
  width: var(--portal-sidebar-width)!important;
  height: 100vh!important;
  z-index: 1045!important;
  overflow: hidden!important;
  background: #fff!important;
  border-right: 1px solid #f4f6f9!important;
  box-shadow: 0 4px 25px rgba(0,0,0,.04)!important;
  transition: width 0.3s ease!important;
}
.user-portal #sidebar-wrapper {
  display: flex!important;
  flex-direction: column!important;
  height: 100vh!important;
  min-height: 0!important;
}
.user-portal .sidebar-brand {
  display: flex!important;
  align-items: center!important;
  justify-content: center!important;
  flex: 0 0 var(--portal-navbar-height)!important;
  height: var(--portal-navbar-height)!important;
  border-bottom: 1px solid #f4f6f9!important;
}
.user-portal .header-logo {
  max-width: 84px!important;
  max-height: 48px!important;
  object-fit: contain!important;
}
.user-portal .sidebar-user {
  display: flex!important;
  gap: 10px!important;
  align-items: center!important;
  flex: 0 0 auto!important;
  min-width: 0!important;
  margin: 12px 14px!important;
  padding: 12px!important;
  border-radius: 8px!important;
  background: #f8f9fa!important;
  transition: padding 0.3s ease!important;
}
.user-portal .sidebar-user img,
.sidebar-avatar-fallback {
  width: 38px!important;
  height: 38px!important;
  flex: 0 0 38px!important;
  border-radius: 50%!important;
  object-fit: cover!important;
}
.sidebar-avatar-fallback,
.nav-avatar-fallback {
  display: inline-flex!important;
  align-items: center!important;
  justify-content: center!important;
  background: #6777ef!important;
  color: #fff!important;
  font-weight: 700!important;
}
.user-portal .sidebar-user strong {
  display: block!important;
  max-width: 130px!important;
  overflow: hidden!important;
  text-overflow: ellipsis!important;
  white-space: nowrap!important;
  color: #34395e!important;
  font-size: 13px!important;
}
.user-portal .sidebar-user span:last-child {
  display: block!important;
  color: #98a6ad!important;
  font-size: 11px!important;
}
.user-portal .sidebar-menu {
  flex: 1 1 auto!important;
  min-height: 0!important;
  overflow-y: auto!important;
  overflow-x: hidden!important;
  padding: 0 0 18px!important;
  scrollbar-width: thin!important;
}
.user-portal .sidebar-menu::-webkit-scrollbar {
  width: 5px;
}
.user-portal .sidebar-menu::-webkit-scrollbar-thumb {
  background: #d7dbea;
  border-radius: 10px;
}
.user-portal .sidebar-menu .menu-header {
  padding: 10px 18px 5px!important;
  color: #abb6ce!important;
  font-size: 10px!important;
  font-weight: 700!important;
  text-transform: uppercase!important;
  letter-spacing: 0.8px!important;
}
.user-portal .sidebar-menu button {
  border: 0!important;
  background: transparent!important;
  width: 100%!important;
  text-align: left!important;
}
.user-portal .sidebar-menu .nav-link {
  display: flex!important;
  align-items: center!important;
  gap: 12px!important;
  padding: 12px 18px!important;
  color: #555c70!important;
  font-size: 13px!important;
  font-weight: 500!important;
  transition: all 0.2s ease!important;
}
.user-portal .sidebar-menu .nav-link i {
  font-size: 16px!important;
  color: #78829d!important;
  transition: color 0.2s ease!important;
}
.user-portal .sidebar-menu li.active > button .nav-link {
  background: #f0f3ff!important;
  color: #6777ef!important;
  font-weight: 600!important;
  border-left: 3px solid #6777ef!important;
}
.user-portal .sidebar-menu li.active > button .nav-link i {
  color: #6777ef!important;
}
.portal-nav-badge {
  display: inline-flex!important;
  align-items: center!important;
  justify-content: center!important;
  min-width: 18px!important;
  height: 18px!important;
  margin-left: auto!important;
  padding: 0 4px!important;
  border-radius: 10px!important;
  background: #ffa426!important;
  color: #fff!important;
  font-size: 9px!important;
  font-weight: 700!important;
}

/* Layout Content */
.user-portal .main-content {
  flex: 1 0 auto!important;
  margin-left: var(--portal-sidebar-width)!important;
  width: calc(100% - var(--portal-sidebar-width))!important;
  max-width: calc(100% - var(--portal-sidebar-width))!important;
  padding: calc(var(--portal-navbar-height) + 24px) var(--portal-page-pad) 30px!important;
  background: #f4f6f9!important;
  transition: margin-left 0.3s ease, padding 0.3s ease, width 0.3s ease, max-width 0.3s ease!important;
}
.user-portal .cms-main-footer {
  flex: 0 0 auto!important;
  margin-top: auto!important;
  margin-left: var(--portal-sidebar-width)!important;
  width: calc(100% - var(--portal-sidebar-width))!important;
  max-width: calc(100% - var(--portal-sidebar-width))!important;
  border-top: 1px solid #e4e6fc!important;
  background: #fff!important;
  color: #6c757d!important;
  padding: 16px var(--portal-page-pad)!important;
  display: flex!important;
  align-items: center!important;
  justify-content: space-between!important;
  font-size: 12px!important;
  transition: margin-left 0.3s ease, width 0.3s ease, max-width 0.3s ease!important;
}
.user-portal .cms-main-footer .footer-left {
  font-weight: 700!important;
  color: #34395e!important;
}

/* Top Actions */
.cms-top-icon {
  display: inline-flex!important;
  align-items: center!important;
  justify-content: center!important;
  width: 38px!important;
  height: 38px!important;
  border-radius: 8px!important;
  border: 1px solid #e4e6fc!important;
  background: #fdfdff!important;
  cursor: pointer!important;
  transition: all 0.15s ease!important;
}
.cms-top-icon:hover {
  background: #f4f6f9!important;
  border-color: #cbd2f6!important;
}
.portal-navbar-title {
  display: flex!important;
  flex-direction: column!important;
  justify-content: center!important;
  margin-left: 10px!important;
}
.portal-navbar-title small {
  color: #98a6ad!important;
  font-size: 9px!important;
  font-weight: 800!important;
  text-transform: uppercase!important;
  letter-spacing: 0.5px!important;
}
.portal-navbar-title strong {
  color: #34395e!important;
  font-size: 14px!important;
  font-weight: 700!important;
}
.portal-top-action {
  position: relative!important;
}
.headerBadge1, .headerBadge2 {
  position: absolute!important;
  top: -4px!important;
  right: -4px!important;
  background: #fc544b!important;
  color: #fff!important;
  font-size: 8px!important;
  padding: 3px 5px!important;
  border-radius: 10px!important;
}
.nav-link-user,
.nav-link-user.cms-top-icon {
  position: relative!important;
  display: inline-flex!important;
  align-items: center!important;
  justify-content: center!important;
  width: 38px!important;
  height: 38px!important;
  min-width: 38px!important;
  max-width: 38px!important;
  min-height: 38px!important;
  max-height: 38px!important;
  flex: 0 0 38px!important;
  padding: 0!important;
  margin: 0!important;
  border-radius: 50%!important;
  background: #f8f9fa!important;
  cursor: pointer!important;
  overflow: hidden!important;
  border: 1px solid #e4e6fc!important;
}
.nav-link-user::after {
  display: none!important;
}
.nav-link-user img,
.nav-link-user.cms-top-icon img {
  position: absolute!important;
  top: 0!important;
  left: 0!important;
  width: 100%!important;
  height: 100%!important;
  min-width: 100%!important;
  max-width: 100%!important;
  min-height: 100%!important;
  max-height: 100%!important;
  border-radius: 50%!important;
  object-fit: cover!important;
}
.nav-avatar-fallback {
  position: absolute!important;
  top: 0!important;
  left: 0!important;
  display: inline-flex!important;
  align-items: center!important;
  justify-content: center!important;
  width: 100%!important;
  height: 100%!important;
  border-radius: 50%!important;
  font-size: 13px!important;
  font-weight: 700!important;
  background: #6777ef!important;
  color: #fff!important;
}

/* Alerts */
.cms-message, .cms-error {
  margin: 0 0 16px!important;
  padding: 12px 16px!important;
  border-radius: 8px!important;
  font-size: 13px!important;
  box-shadow: 0 4px 12px rgba(0,0,0,.03)!important;
}
.cms-message {
  background: #e8f7f0!important;
  color: #218b55!important;
  border: 1px solid #a7f3d0!important;
}
.cms-error {
  background: #fdeaea!important;
  color: #fc544b!important;
  border: 1px solid #fecaca!important;
}

/* Cards & KPIs */
.page-heading {
  display: flex!important;
  align-items: center!important;
  justify-content: space-between!important;
  gap: 20px!important;
  margin-bottom: 24px!important;
}
.page-heading p {
  margin: 0 0 4px!important;
  color: #6777ef!important;
  font-size: 11px!important;
  font-weight: 800!important;
  text-transform: uppercase!important;
  letter-spacing: 0.8px!important;
}
.page-heading h1 {
  margin: 0 0 6px!important;
  color: #34395e!important;
  font-size: 24px!important;
  font-weight: 700!important;
}
.page-heading span {
  color: #78829d!important;
  font-size: 13px!important;
}

.primary-command, .secondary-command {
  display: inline-flex!important;
  align-items: center!important;
  justify-content: center!important;
  gap: 8px!important;
  min-height: 42px!important;
  padding: 0 20px!important;
  border-radius: 30px!important;
  font-size: 13px!important;
  font-weight: 600!important;
  cursor: pointer!important;
  transition: all 0.2s ease!important;
}
.primary-command {
  background: #6777ef!important;
  color: #fff!important;
  border: 0!important;
  box-shadow: 0 4px 12px rgba(103,119,239,0.35)!important;
}
.primary-command:hover {
  background: #4e61e8!important;
  transform: translateY(-1px)!important;
}
.secondary-command {
  border: 1px solid #e4e6fc!important;
  background: #fff!important;
  color: #34395e!important;
}
.secondary-command:hover {
  background: #f8f9fa!important;
}

.otika-page-actions {
  display: flex!important;
  gap: 10px!important;
  margin-bottom: 20px!important;
}
.otika-page-actions .btn {
  border-radius: 30px!important;
  font-size: 12px!important;
  font-weight: 600!important;
  padding: 8px 16px!important;
}

.user-kpis {
  display: grid!important;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr))!important;
  gap: 16px!important;
  margin-bottom: 24px!important;
}
.user-kpis article {
  display: flex!important;
  align-items: center!important;
  gap: 14px!important;
  padding: 18px!important;
  border-radius: 12px!important;
  background: #fff!important;
  box-shadow: 0 4px 15px rgba(0,0,0,.03)!important;
  border: 1px solid #f0f2f8!important;
}
.user-kpis article > span {
  display: grid!important;
  place-items: center!important;
  flex: 0 0 48px!important;
  width: 48px!important;
  height: 48px!important;
  border-radius: 10px!important;
  font-size: 22px!important;
}
.user-kpis .blue { background: #e8edff!important; color: #6777ef!important; }
.user-kpis .amber { background: #fff4e6!important; color: #ffa426!important; }
.user-kpis .green { background: #e8f7f0!important; color: #47c363!important; }
.user-kpis .cyan { background: #eaf7ff!important; color: #3abaf4!important; }

.user-kpis small {
  color: #abb6ce!important;
  font-size: 10px!important;
  font-weight: 700!important;
  text-transform: uppercase!important;
}
.user-kpis strong {
  color: #34395e!important;
  font-size: 24px!important;
  font-weight: 700!important;
  margin: 2px 0!important;
}
.user-kpis p {
  color: #8a94ad!important;
  font-size: 11px!important;
  margin: 0!important;
}

/* Tables & Lists */
.list-toolbar {
  display: grid!important;
  grid-template-columns: 1fr 220px!important;
  gap: 12px!important;
  margin-bottom: 16px!important;
}
.list-toolbar label {
  display: flex!important;
  align-items: center!important;
  gap: 8px!important;
  background: #fff!important;
  padding: 0 14px!important;
  border: 1px solid #e4e6fc!important;
  border-radius: 8px!important;
}
.list-toolbar label input {
  border: 0!important;
  outline: 0!important;
  width: 100%!important;
  font-size: 13px!important;
}
.list-toolbar select {
  height: 44px!important;
  padding: 0 12px!important;
  border-radius: 8px!important;
  border: 1px solid #e4e6fc!important;
  background: #fff!important;
  color: #495057!important;
  font-size: 13px!important;
}

.data-table {
  overflow-x: auto!important;
  border-radius: 12px!important;
  background: #fff!important;
  box-shadow: 0 4px 15px rgba(0,0,0,.02)!important;
  border: 1px solid #edf0f5!important;
}
.data-table table {
  width: 100%!important;
  border-collapse: collapse!important;
}
.data-table th, .data-table td {
  padding: 16px 20px!important;
  font-size: 13px!important;
  border-bottom: 1px solid #edf0f5!important;
  text-align: left!important;
}
.data-table th {
  background: #fafbfe!important;
  color: #78829d!important;
  font-weight: 700!important;
  text-transform: uppercase!important;
  font-size: 11px!important;
  letter-spacing: 0.5px!important;
}
.data-table td {
  color: #555c70!important;
}
.data-table td strong {
  color: #34395e!important;
  font-weight: 600!important;
}
.data-table td small {
  color: #98a6ad!important;
  font-size: 11px!important;
}
.table-actions {
  display: flex!important;
  gap: 6px!important;
}
.table-actions button {
  display: inline-flex!important;
  align-items: center!important;
  justify-content: center!important;
  width: 32px!important;
  height: 32px!important;
  border-radius: 8px!important;
  border: 1px solid #e4e6fc!important;
  background: #fdfdff!important;
  color: #6777ef!important;
  cursor: pointer!important;
  transition: all 0.15s ease!important;
}
.table-actions button:hover {
  background: #6777ef!important;
  color: #fff!important;
  border-color: #6777ef!important;
}

.status {
  display: inline-flex!important;
  padding: 4px 10px!important;
  border-radius: 30px!important;
  font-size: 10px!important;
  font-weight: 700!important;
  text-transform: uppercase!important;
}
.status.blue { background: #e8edff!important; color: #6777ef!important; }
.status.green { background: #e8f7f0!important; color: #47c363!important; }
.status.amber { background: #fff4e6!important; color: #ffa426!important; }
.status.red { background: #fdeaea!important; color: #fc544b!important; }

/* Timeline & Feed */
.timeline-list, .feed-list {
  display: grid!important;
  gap: 12px!important;
}
.timeline-list article, .feed-list article {
  display: flex!important;
  align-items: flex-start!important;
  gap: 14px!important;
  padding: 18px!important;
  border-radius: 12px!important;
  background: #fff!important;
  box-shadow: 0 4px 15px rgba(0,0,0,.02)!important;
  border: 1px solid #f0f2f8!important;
}
.timeline-list article > span, .feed-list article > span {
  display: grid!important;
  place-items: center!important;
  flex: 0 0 38px!important;
  width: 38px!important;
  height: 38px!important;
  border-radius: 8px!important;
  background: #e8edff!important;
  color: #6777ef!important;
}
.timeline-list strong, .feed-list strong {
  color: #34395e!important;
  font-weight: 650!important;
}
.timeline-list p, .feed-list p {
  margin: 4px 0!important;
  color: #6c757d!important;
  font-size: 13px!important;
  line-height: 1.45!important;
}
.timeline-list small, .feed-list small {
  color: #a3abc0!important;
  font-size: 11px!important;
}
.feed-list article.unread {
  border-left: 4px solid #6777ef!important;
}
.feed-list article > button {
  margin-left: auto!important;
  border: 0!important;
  background: transparent!important;
  color: #6777ef!important;
  cursor: pointer!important;
}

/* Profile view */
.profile-layout {
  display: flex!important;
  flex-direction: column!important;
  gap: 24px!important;
}
.profile-grid-2 {
  display: grid!important;
  grid-template-columns: repeat(2, minmax(0, 1fr))!important;
  gap: 24px!important;
}
.profile-grid-3 {
  display: grid!important;
  grid-template-columns: 270px 1fr!important;
  gap: 24px!important;
}
@media (max-width: 768px) {
  .profile-grid-2, .profile-grid-3 {
    grid-template-columns: 1fr!important;
  }
}
.profile-summary, .profile-form {
  padding: 24px!important;
  border-radius: 12px!important;
  background: #fff!important;
  border: 1px solid #f0f2f8!important;
  box-shadow: 0 4px 15px rgba(0,0,0,.02)!important;
}
.profile-summary {
  text-align: center!important;
}
.profile-avatar {
  display: grid!important;
  place-items: center!important;
  width: 88px!important;
  height: 88px!important;
  margin: 0 auto 16px!important;
  overflow: hidden!important;
  border-radius: 50%!important;
  background: #6777ef!important;
  color: #fff!important;
  font-size: 26px!important;
  font-weight: 700!important;
}
.profile-avatar img {
  width: 100%!important;
  height: 100%!important;
  object-fit: cover!important;
}
.profile-summary h2, .profile-form h2 {
  margin: 0 0 6px!important;
  color: #34395e!important;
  font-size: 18px!important;
  font-weight: 700!important;
}
.profile-summary p {
  color: #78829d!important;
  font-size: 13px!important;
  margin-bottom: 12px!important;
}
.ordinary-badge {
  display: inline-flex!important;
  padding: 5px 12px!important;
  border-radius: 30px!important;
  background: #e8edff!important;
  color: #6777ef!important;
  font-size: 10px!important;
  font-weight: 700!important;
  text-transform: uppercase!important;
  letter-spacing: 0.5px!important;
}
.profile-form {
  display: grid!important;
  gap: 16px!important;
}
.profile-form > div {
  display: grid!important;
  grid-template-columns: 1fr 1fr!important;
  gap: 16px!important;
}
.profile-form label {
  display: grid!important;
  gap: 6px!important;
  color: #34395e!important;
  font-size: 12px!important;
  font-weight: 600!important;
}
.profile-form input, 
.form-field input, 
.form-field textarea, 
.form-field select, 
.payment-section input {
  width: 100%!important;
  padding: 10px 14px!important;
  border: 1px solid #e4e6fc!important;
  border-radius: 8px!important;
  background: #fdfdff!important;
  color: #495057!important;
  outline: none!important;
  font-size: 13px!important;
  transition: all 0.2s ease!important;
}
.profile-form input:focus,
.form-field input:focus,
.form-field textarea:focus,
.form-field select:focus {
  border-color: #6777ef!important;
  box-shadow: 0 2px 8px rgba(103,119,239,0.15)!important;
}
.profile-form input:disabled {
  background: #f4f6f9!important;
  cursor: not-allowed!important;
}
.profile-form button {
  justify-self: start!important;
}

/* Modals */
.form-backdrop {
  position: fixed!important;
  inset: 0!important;
  z-index: 1080!important;
  display: grid!important;
  place-items: center!important;
  padding: var(--portal-page-pad)!important;
  background: rgba(17, 24, 39, 0.6)!important;
  backdrop-filter: blur(4px)!important;
}
.application-modal {
  width: min(760px, 100%)!important;
  max-height: 90dvh!important;
  overflow-y: auto!important;
  border-radius: 12px!important;
  background: #fff!important;
  box-shadow: 0 20px 50px rgba(15,23,42,0.15)!important;
  border: 1px solid #edf0f5!important;
}
.application-modal > header {
  display: flex!important;
  justify-content: space-between!important;
  gap: 18px!important;
  padding: 20px 24px!important;
  border-bottom: 1px solid #f4f6f9!important;
}
.application-modal > header h2 {
  margin: 0 0 4px!important;
  color: #34395e!important;
  font-size: 20px!important;
  font-weight: 700!important;
}
.application-modal > header p {
  margin: 0!important;
  color: #78829d!important;
  font-size: 12px!important;
}
.application-modal > header button {
  display: grid!important;
  place-items: center!important;
  width: 32px!important;
  height: 32px!important;
  border-radius: 50%!important;
  background: #f4f6f9!important;
  border: 0!important;
  cursor: pointer!important;
}
.application-modal > form {
  display: grid!important;
  gap: 16px!important;
  padding: 24px!important;
}
.fee-banner {
  display: flex!important;
  align-items: center!important;
  justify-content: space-between!important;
  padding: 12px 16px!important;
  border-left: 4px solid #ffa426!important;
  background: #fff4e6!important;
  border-radius: 8px!important;
}
.fee-banner span { color: #8a94ad!important; font-size: 12px!important; }
.fee-banner strong { color: #d97706!important; font-size: 15px!important; font-weight: 700!important; }

.form-field {
  display: grid!important;
  gap: 6px!important;
}
.form-field > label {
  color: #34395e!important;
  font-size: 12px!important;
  font-weight: 600!important;
}
.form-field > label b {
  color: #fc544b!important;
}
.form-field > p {
  margin: 0!important;
  color: #98a6ad!important;
  font-size: 11px!important;
}
.choice-list {
  display: grid!important;
  gap: 8px!important;
}
.choice-list label {
  display: flex!important;
  align-items: center!important;
  gap: 8px!important;
  color: #495057!important;
  font-size: 13px!important;
  cursor: pointer!important;
}
.choice-list input {
  width: auto!important;
}
.file-input {
  display: flex!important;
  align-items: center!important;
  gap: 10px!important;
  flex-wrap: wrap!important;
}
.file-input a {
  color: #6777ef!important;
  font-size: 13px!important;
  font-weight: 600!important;
}
.payment-section {
  display: grid!important;
  gap: 12px!important;
  padding: 16px!important;
  border: 1px solid #ffe2ad!important;
  border-radius: 8px!important;
  background: #fffaf0!important;
}
.payment-section h3 {
  margin: 0!important;
  color: #b45309!important;
  font-size: 14px!important;
  font-weight: 700!important;
}
.payment-section p {
  margin: 0!important;
  color: #b45309!important;
  font-size: 11px!important;
  opacity: 0.85!important;
}
.payment-section label {
  display: grid!important;
  gap: 6px!important;
}
.application-modal form > footer {
  display: flex!important;
  justify-content: flex-end!important;
  gap: 10px!important;
  padding-top: 16px!important;
  border-top: 1px solid #f4f6f9!important;
}

/* Scrim & Backdrop drawer */
.sidebar-scrim {
  position: fixed!important;
  inset: 0!important;
  z-index: 1040!important;
  background: rgba(15, 23, 42, 0.48)!important;
  backdrop-filter: blur(3px)!important;
}

.cms-pagination {
  display: flex!important;
  align-items: center!important;
  justify-content: flex-end!important;
  gap: 10px!important;
  flex-wrap: wrap!important;
  margin-top: 20px!important;
}
.cms-pagination button {
  border: 0!important;
  border-radius: 30px!important;
  background: #6777ef!important;
  color: #fff!important;
  padding: 8px 16px!important;
  font-size: 12px!important;
  font-weight: 700!important;
  box-shadow: 0 2px 6px #acb5f6!important;
  cursor: pointer!important;
}
.cms-pagination button:disabled {
  cursor: not-allowed!important;
  opacity: .5!important;
}
.cms-pagination span {
  color: #6c757d!important;
  font-size: 12px!important;
  font-weight: 700!important;
}
:global(.dark .cms-pagination span) {
  color: #cbd5e1!important;
}

/* Dark Mode Overrides */
:global(.dark .user-portal) {
  background: #0f172a!important;
  color: #cbd5e1!important;
}
:global(.dark .user-portal .main-content) {
  background: #0f172a!important;
}
:global(.dark .user-portal .navbar-bg),
:global(.dark .user-portal .main-navbar) {
  background: #111827!important;
  color: #f8fafc!important;
  box-shadow: 0 4px 24px rgba(0,0,0,.35)!important;
  border-color: #1f2937!important;
}
:global(.dark .user-portal .main-navbar .nav-link),
:global(.dark .user-portal .main-navbar i) {
  color: #f8fafc!important;
}
:global(.dark .user-portal .main-sidebar) {
  background: #111827!important;
  border-color: #1f2937!important;
}
:global(.dark .user-portal .sidebar-brand) {
  border-color: #1f2937!important;
}
:global(.dark .user-portal .sidebar-user) {
  background: #1f2937!important;
}
:global(.dark .user-portal .sidebar-user strong),
:global(.dark .user-portal .sidebar-user span) {
  color: #f8fafc!important;
}
:global(.dark .user-portal .sidebar-menu .nav-link) {
  color: #94a3b8!important;
}
:global(.dark .user-portal .sidebar-menu li.active > button .nav-link) {
  background: #1e293b!important;
  color: #93c5fd!important;
  border-color: #93c5fd!important;
}
:global(.dark .user-portal .sidebar-menu li.active > button .nav-link i) {
  color: #93c5fd!important;
}
:global(.dark .user-portal .cms-main-footer) {
  background: #111827!important;
  border-color: #1f2937!important;
  color: #cbd5e1!important;
}
:global(.dark .user-portal .cms-main-footer .footer-left) {
  color: #f8fafc!important;
}
:global(.dark .user-portal .user-kpis article),
:global(.dark .user-portal .data-table),
:global(.dark .user-portal .timeline-list article),
:global(.dark .user-portal .feed-list article),
:global(.dark .user-portal .transaction-summary),
:global(.dark .user-portal .profile-summary),
:global(.dark .user-portal .profile-form),
:global(.dark .user-portal .application-modal) {
  background: #1f2937!important;
  border-color: #334155!important;
  color: #cbd5e1!important;
}
:global(.dark .user-portal .page-heading h1),
:global(.dark .user-portal .user-kpis strong),
:global(.dark .user-portal .data-table td strong),
:global(.dark .user-portal .data-table th),
:global(.dark .user-portal .timeline-list strong),
:global(.dark .user-portal .feed-list strong),
:global(.dark .user-portal .profile-summary h2),
:global(.dark .user-portal .profile-form h2),
:global(.dark .user-portal .profile-form label),
:global(.dark .user-portal .application-modal h2) {
  color: #f8fafc!important;
}
:global(.dark .user-portal input),
:global(.dark .user-portal textarea),
:global(.dark .user-portal select),
:global(.dark .user-portal .list-toolbar label) {
  background: #111827!important;
  color: #f8fafc!important;
  border-color: #475569!important;
}
:global(.dark .user-portal .data-table th) {
  background: #111827!important;
}
:global(.dark .user-portal .dropdown-menu) {
  background: #1f2937!important;
  border-color: #334155!important;
  color: #cbd5e1!important;
}
:global(.dark .user-portal .dropdown-item) {
  color: #cbd5e1!important;
}
:global(.dark .user-portal .dropdown-item:hover) {
  background: #111827!important;
}
:global(.dark .user-portal .nav-link-user) {
  background: #1f2937!important;
  color: #f8fafc!important;
  border-color: #475569!important;
}
:global(.dark .user-portal .cms-top-icon) {
  background: #1f2937!important;
  border-color: #475569!important;
  color: #f8fafc!important;
}
:global(.dark .user-portal .cms-top-icon:hover) {
  background: #111827!important;
  border-color: #64748b!important;
}
:global(.dark .user-portal .cms-top-icon i) {
  color: #f8fafc!important;
}

/* Responsive Overrides */
@media (max-width: 991px) {
  .user-portal .main-wrapper {
    --portal-sidebar-width: 0px;
    --portal-page-pad: 16px;
  }
  .user-portal .main-sidebar {
    position: fixed!important;
    left: 0!important;
    top: 0!important;
    width: 260px!important;
    height: 100vh!important;
    z-index: 1045!important;
    transform: translateX(-100%)!important;
    transition: transform 0.3s ease!important;
    box-shadow: 0 8px 24px rgba(15, 23, 42, 0.15)!important;
    border-right: 1px solid #f4f6f9!important;
  }
  .user-portal .main-sidebar.mobile-sidebar-open {
    transform: translateX(0)!important;
  }
  .user-portal .navbar-bg,
  .user-portal .main-navbar {
    left: 0!important;
    width: 100%!important;
    max-width: 100%!important;
  }
  .user-portal .main-content,
  .user-portal .cms-main-footer {
    margin-left: 0!important;
    width: 100%!important;
    max-width: 100%!important;
  }
  .user-portal .user-kpis {
    grid-template-columns: repeat(2, minmax(0, 1fr))!important;
  }
  .profile-layout {
    grid-template-columns: 1fr!important;
  }
}

@media (max-width: 768px) {
  .portal-navbar-title {
    display: none!important;
  }
  .page-heading {
    flex-direction: column!important;
    align-items: flex-start!important;
    gap: 12px!important;
  }
  .page-heading .primary-command {
    width: 100%!important;
    justify-content: center!important;
  }
  .list-toolbar {
    grid-template-columns: 1fr!important;
  }
  .profile-form > div {
    grid-template-columns: 1fr!important;
  }
  .user-portal .cms-main-footer {
    flex-direction: column!important;
    text-align: center!important;
    gap: 8px!important;
  }
  .user-portal .cms-main-footer .footer-right {
    text-align: center!important;
    margin-left: 0!important;
  }
  
  /* Responsive Tables as Cards */
  .data-table {
    background: transparent!important;
    box-shadow: none!important;
    border: 0!important;
  }
  .data-table table,
  .data-table thead,
  .data-table tbody,
  .data-table tr,
  .data-table td {
    display: block!important;
    width: 100%!important;
  }
  .data-table thead {
    display: none!important;
  }
  .data-table tr {
    margin-bottom: 12px!important;
    padding: 16px!important;
    border-radius: 12px!important;
    background: #fff!important;
    box-shadow: 0 4px 15px rgba(0,0,0,.03)!important;
    border: 1px solid #edf0f5!important;
  }
  :global(.dark) .data-table tr {
    background: #1f2937!important;
    border-color: #334155!important;
  }
  .data-table td {
    display: grid!important;
    grid-template-columns: minmax(100px, 0.4fr) minmax(0, 1fr)!important;
    gap: 8px!important;
    align-items: start!important;
    padding: 8px 0!important;
    border-bottom: 1px solid #f4f6f9!important;
  }
  :global(.dark) .data-table td {
    border-bottom-color: #334155!important;
  }
  .data-table td:last-child {
    border-bottom: 0!important;
  }
  .data-table td::before {
    content: attr(data-label);
    color: #98a6ad;
    font-size: 10px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  .data-table td.empty-cell {
    display: block!important;
    padding: 24px 12px!important;
    text-align: center!important;
  }
  .data-table td.empty-cell::before {
    content: ""!important;
  }
}

@media (max-width: 576px) {
  .user-portal .main-wrapper {
    --portal-navbar-height: 64px;
    --portal-page-pad: 12px;
  }
  .page-heading h1 {
    font-size: 20px!important;
  }
  .user-kpis {
    grid-template-columns: 1fr!important;
  }
  .user-kpis article {
    min-height: auto!important;
  }
  .page-heading .primary-command,
  .otika-page-actions .btn {
    width: 100%!important;
    text-align: center!important;
  }
  .timeline-list article, .feed-list article {
    grid-template-columns: 38px 1fr!important;
    padding: 14px!important;
  }
  .feed-list article > button {
    grid-column: 2!important;
    justify-self: start!important;
    margin-top: 4px!important;
  }
  .transaction-summary {
    flex-direction: column!important;
    align-items: flex-start!important;
    gap: 6px!important;
  }
  .application-modal {
    max-height: 100vh!important;
    border-radius: 0!important;
  }
  .form-backdrop {
    padding: 0!important;
  }
  .application-modal > header {
    padding: 16px 20px!important;
  }
  .application-modal > form {
    padding: 20px!important;
  }
  .application-modal form > footer {
    flex-direction: column-reverse!important;
    gap: 8px!important;
  }
  .application-modal form > footer button {
    width: 100%!important;
  }
  .navbar .dropdown-menu.show {
    left: 8px!important;
    right: 8px!important;
  }
}

/* Athlete Dashboard Grid & Cards */
.athlete-registry-dashboard {
  margin-top: 30px!important;
}
.bg-gradient-primary-to-secondary {
  background: linear-gradient(135deg, #6777ef 0%, #3abaf4 100%)!important;
}
.athlete-badge {
  display: grid!important;
  place-items: center!important;
  width: 56px!important;
  height: 56px!important;
  background: rgba(255, 255, 255, 0.2)!important;
  border-radius: 12px!important;
  font-size: 28px!important;
}
.watermark-icon {
  position: absolute!important;
  right: -20px!important;
  bottom: -30px!important;
  font-size: 150px!important;
  opacity: 0.1!important;
  pointer-events: none!important;
}
.athlete-stats-grid {
  display: grid!important;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr))!important;
  gap: 20px!important;
  margin-top: 20px!important;
}
.athlete-detail-card {
  background: #fff!important;
  border: 1px solid #f0f2f8!important;
  border-radius: 12px!important;
  padding: 20px!important;
  box-shadow: 0 4px 15px rgba(0,0,0,.03)!important;
}
:global(.dark .athlete-detail-card) {
  background: #1e293b!important;
  border-color: #334155!important;
}
.athlete-detail-card header {
  display: flex!important;
  align-items: center!important;
  gap: 10px!important;
  margin-bottom: 16px!important;
  border-bottom: 1px solid #f0f2f8!important;
  padding-bottom: 10px!important;
}
:global(.dark .athlete-detail-card header) {
  border-color: #334155!important;
}
.athlete-detail-card header i {
  font-size: 20px!important;
}
.athlete-detail-card header h3 {
  margin: 0!important;
  font-size: 15px!important;
  font-weight: 700!important;
  color: #34395e!important;
}
:global(.dark .athlete-detail-card header h3) {
  color: #f8fafc!important;
}
.athlete-detail-card ul {
  list-style: none!important;
  padding: 0!important;
  margin: 0!important;
  display: grid!important;
  gap: 12px!important;
}
.athlete-detail-card li {
  display: flex!important;
  justify-content: space-between!important;
  align-items: center!important;
  font-size: 13px!important;
}
.athlete-detail-card li span {
  color: #8a94ad!important;
}
.athlete-detail-card li strong {
  color: #34395e!important;
  font-weight: 600!important;
}
:global(.dark .athlete-detail-card li strong) {
  color: #cbd5e1!important;
}
</style>
