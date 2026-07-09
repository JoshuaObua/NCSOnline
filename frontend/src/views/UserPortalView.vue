<template>
  <main class="otika-cms user-portal">
    <div class="otika-app">
      <div class="main-wrapper main-wrapper-1" :class="{ 'sidebar-mini': sidebarCollapsed }">
        <div class="navbar-bg"></div>
        <nav class="navbar navbar-expand-lg main-navbar sticky">
          <div class="form-inline mr-auto">
            <ul class="navbar-nav mr-3">
              <li>
                <button type="button" class="nav-link nav-link-lg cms-top-icon collapse-btn" :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'" @click="toggleSidebar">
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
                <span class="d-sm-none d-lg-inline-block">{{ firstName }}</span>
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

        <div class="main-sidebar sidebar-style-2">
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
                        <td><strong>{{ item.title }}</strong><small>{{ item.source === 'custom' ? 'Custom form' : 'Standard form' }}</small></td>
                        <td>{{ item.reference || 'Pending' }}</td>
                        <td><span class="status" :class="statusClass(item.status)">{{ titleize(item.status) }}</span></td>
                        <td>{{ titleize(item.payment_status || 'not required') }}</td>
                        <td>{{ formatDate(item.updated_at) }}</td>
                        <td><button v-if="item.source === 'custom' && item.status === 'DRAFT'" type="button" title="Continue application" @click="continueApplication(item)"><i class="icofont-rounded-right"></i></button></td>
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
                      <tr v-for="item in transactions" :key="`${item.source}-${item.id}`"><td><strong>{{ item.title }}</strong></td><td>{{ item.payment_reference || '-' }}</td><td>UGX {{ formatMoney(item.amount) }}</td><td><span class="status" :class="paymentClass(item.payment_status)">{{ titleize(item.payment_status || 'unpaid') }}</span></td><td>{{ titleize(item.payment_method || 'proof upload') }}</td></tr>
                    </tbody>
                  </table>
                </div>
              </template>

              <template v-else-if="section === 'profile'">
                <header class="page-heading"><div><p>My Profile</p><h1>Personal information</h1><span>Keep your identity details current for applications and communication.</span></div></header>
                <div class="profile-layout">
                  <article class="profile-summary">
                    <div class="profile-avatar"><img v-if="profile.avatar_url" :src="mediaUrl(profile.avatar_url)" alt="" /><span v-else>{{ initials }}</span></div>
                    <h2>{{ fullName }}</h2><p>{{ profile.email }}</p><span class="ordinary-badge">Ordinary user</span>
                  </article>
                  <form class="profile-form" @submit.prevent="saveProfile">
                    <h2>Profile details</h2>
                    <div><label>First name<input v-model.trim="profile.first_name" required /></label><label>Last name<input v-model.trim="profile.last_name" required /></label></div>
                    <label>Email address<input :value="profile.email" type="email" disabled /></label>
                    <label>Profile photo URL<input v-model.trim="profile.avatar_url" type="url" placeholder="https://..." /></label>
                    <button type="submit" class="primary-command" :disabled="savingProfile"><i class="icofont-save"></i> {{ savingProfile ? 'Saving...' : 'Save profile' }}</button>
                  </form>
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

    <div v-if="activeForm" class="form-backdrop" @click.self="closeApplication">
      <section class="application-modal" role="dialog" aria-modal="true" :aria-label="activeForm.title">
        <header>
          <div><small>{{ activeForm.department_name || 'NCS Application' }}</small><h2>{{ activeForm.title }}</h2><p>{{ activeForm.description }}</p></div>
          <button type="button" title="Close application" @click="closeApplication"><i class="icofont-close"></i></button>
        </header>
        <div v-if="formLoading" class="form-loading">Loading application form...</div>
        <form v-else @submit.prevent="submitApplication">
          <div class="fee-banner"><span>Application fee</span><strong>{{ activeForm.price_ugx ? `UGX ${formatMoney(activeForm.price_ugx)}` : 'Free' }}</strong></div>
          <div v-for="field in activeForm.fields || []" :key="field.id || field.field_key" class="form-field">
            <label :for="`field-${field.field_key}`">{{ field.label }} <b v-if="field.is_required">*</b></label>
            <p v-if="field.help_text">{{ field.help_text }}</p>
            <textarea v-if="field.field_type === 'long_text'" :id="`field-${field.field_key}`" v-model="formAnswers[field.field_key]" :required="field.is_required" :placeholder="field.placeholder" rows="4"></textarea>
            <select v-else-if="field.field_type === 'dropdown'" :id="`field-${field.field_key}`" v-model="formAnswers[field.field_key]" :required="field.is_required"><option value="">Select an option</option><option v-for="option in fieldOptions(field)" :key="option" :value="option">{{ option }}</option></select>
            <div v-else-if="field.field_type === 'radio'" class="choice-list"><label v-for="option in fieldOptions(field)" :key="option"><input v-model="formAnswers[field.field_key]" type="radio" :name="field.field_key" :value="option" :required="field.is_required" /> {{ option }}</label></div>
            <div v-else-if="field.field_type === 'checkbox'" class="choice-list"><label v-for="option in fieldOptions(field)" :key="option"><input v-model="formAnswers[field.field_key]" type="checkbox" :value="option" /> {{ option }}</label></div>
            <div v-else-if="['file','image'].includes(field.field_type)" class="file-input">
              <input :id="`field-${field.field_key}`" type="file" :accept="fieldAccept(field)" :required="field.is_required && !formAnswers[field.field_key]" @change="uploadFieldFile(field, $event)" />
              <span v-if="uploadingField === field.field_key">Uploading...</span><a v-else-if="formAnswers[field.field_key]" :href="mediaUrl(formAnswers[field.field_key])" target="_blank" rel="noopener">Uploaded file</a>
            </div>
            <input v-else :id="`field-${field.field_key}`" v-model="formAnswers[field.field_key]" :type="inputType(field.field_type)" :required="field.is_required" :placeholder="field.placeholder" />
          </div>
          <section v-if="activeForm.price_ugx > 0 && paymentRequired" class="payment-section">
            <h3>Payment proof</h3><p>Enter the reference from your payment receipt before submitting.</p>
            <label>Payment reference<input v-model.trim="paymentReference" required placeholder="Receipt or transaction reference" /></label>
            <label>Amount paid (UGX)<input v-model.number="paymentAmount" type="number" min="1" required /></label>
          </section>
          <p v-if="formError" class="portal-error">{{ formError }}</p>
          <footer>
            <button type="button" class="secondary-command" :disabled="formSaving" @click="saveDraft"><i class="icofont-save"></i> Save draft</button>
            <button type="submit" class="primary-command" :disabled="formSaving || !!uploadingField"><i class="icofont-paper-plane"></i> {{ formSaving ? 'Saving...' : 'Submit application' }}</button>
          </footer>
        </form>
      </section>
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getCurrentUser, updateMyProfile } from '@/api/auth.js'
import { listMyAuditLogs } from '@/api/account.js'
import { listMyLegacyApplications, listMyTransactions } from '@/api/applications.js'
import * as cms from '@/api/cms.js'
import {
  portalGetForm, portalListOpenForms, portalListSubmissions, portalSaveDraft,
  portalSubmit, portalUploadPaymentProof,
} from '@/api/forms.js'
import { mediaUrl } from '@/api/client.js'
import OpenFormsPanel from '@/components/portal/OpenFormsPanel.vue'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'
import { ensureOtikaStyles } from '@/utils/otikaAssets.js'

const router = useRouter()
const portalSectionIds = ['dashboard', 'apply', 'applications', 'activities', 'notifications', 'messages', 'transactions', 'profile']
const section = ref('dashboard')
const sidebarCollapsed = ref(localStorage.getItem('ncsms_sidebar_collapsed') === 'true')
const profileOpen = ref(false)
const loading = ref(true)
const error = ref('')
const success = ref('')
const openForms = ref([])
const dynamicSubmissions = ref([])
const legacyApplications = ref([])
const legacyTransactions = ref([])
const activities = ref([])
const notifications = ref([])
const profile = reactive({})
const applicationSearch = ref('')
const applicationStatus = ref('')
const savingProfile = ref(false)
const activeForm = ref(null)
const formLoading = ref(false)
const formSaving = ref(false)
const formError = ref('')
const formAnswers = reactive({})
const activeSubmission = ref(null)
const uploadingField = ref('')
const paymentReference = ref('')
const paymentAmount = ref(0)

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
const paymentRequired = computed(() => activeForm.value?.price_ugx > 0 && !['PAID', 'PROOF_UPLOADED'].includes(activeSubmission.value?.payment_status))

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
    listMyAuditLogs({ page: 1, per_page: 50 }), cms.listNotifications({ page: 1, per_page: 100 }),
  ])
  if (results[0].status === 'fulfilled') {
    Object.assign(profile, unwrap(results[0].value))
    localStorage.setItem('ncsms_user', JSON.stringify(profile))
  }
  openForms.value = results[1].status === 'fulfilled' ? asList(results[1].value) : []
  dynamicSubmissions.value = results[2].status === 'fulfilled' ? asList(results[2].value) : []
  legacyApplications.value = results[3].status === 'fulfilled' ? asList(results[3].value) : []
  legacyTransactions.value = results[4].status === 'fulfilled' ? asList(results[4].value) : []
  activities.value = results[5].status === 'fulfilled' ? asList(results[5].value) : []
  notifications.value = results[6].status === 'fulfilled' ? asList(results[6].value) : []
  if (results.some(item => item.status === 'rejected')) error.value = 'Some dashboard information could not be loaded. Refresh to try again.'
  loading.value = false
}

function select(id) {
  section.value = id
  profileOpen.value = false
  success.value = ''
  error.value = ''
}
function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('ncsms_sidebar_collapsed', String(sidebarCollapsed.value))
}
async function startApplication(form) {
  activeForm.value = { ...form, fields: [] }
  formLoading.value = true
  formError.value = ''
  clearAnswers()
  try {
    activeForm.value = await portalGetForm(form.slug)
    const draft = dynamicSubmissions.value.find(item => item.template_id === activeForm.value.id && item.status === 'DRAFT')
    activeSubmission.value = draft || null
    Object.assign(formAnswers, parseObject(draft?.answers))
    for (const field of activeForm.value.fields || []) {
      if (field.field_type === 'checkbox' && !Array.isArray(formAnswers[field.field_key])) formAnswers[field.field_key] = []
      else if (formAnswers[field.field_key] == null) formAnswers[field.field_key] = ''
    }
    paymentReference.value = draft?.payment_reference || ''
    paymentAmount.value = draft?.payment_amount_ugx || activeForm.value.price_ugx || 0
  } catch (err) {
    formError.value = apiError(err, 'Could not load this form.')
  } finally {
    formLoading.value = false
  }
}
async function continueApplication(item) {
  const form = openForms.value.find(entry => entry.id === item.template_id)
  if (form) await startApplication(form)
  else error.value = 'This application form is no longer open.'
}
async function saveDraft(showMessage = true) {
  formSaving.value = true
  formError.value = ''
  try {
    activeSubmission.value = await portalSaveDraft(activeForm.value.id, { ...formAnswers })
    if (showMessage) success.value = 'Your draft was saved.'
    await refreshSubmissions()
    return activeSubmission.value
  } catch (err) {
    formError.value = apiError(err, 'Could not save your draft.')
    return null
  } finally {
    formSaving.value = false
  }
}
async function submitApplication() {
  const missing = (activeForm.value?.fields || []).find(field => {
    if (!field.is_required) return false
    const value = formAnswers[field.field_key]
    return Array.isArray(value) ? value.length === 0 : !String(value ?? '').trim()
  })
  if (missing) {
    formError.value = `${missing.label} is required.`
    return
  }
  formSaving.value = true
  formError.value = ''
  try {
    let submission = await portalSaveDraft(activeForm.value.id, { ...formAnswers })
    if (activeForm.value.price_ugx > 0 && !['PAID', 'PROOF_UPLOADED'].includes(submission.payment_status)) {
      if (!paymentReference.value || Number(paymentAmount.value) <= 0) throw new Error('Enter a valid payment reference and amount.')
      await portalUploadPaymentProof(submission.id, paymentReference.value, Number(paymentAmount.value))
    }
    await portalSubmit(submission.id)
    success.value = 'Your application was submitted successfully.'
    closeApplication()
    await refreshSubmissions()
    section.value = 'applications'
  } catch (err) {
    formError.value = apiError(err, err.message || 'Could not submit your application.')
  } finally {
    formSaving.value = false
  }
}
async function uploadFieldFile(field, event) {
  const file = event.target.files?.[0]
  if (!file) return
  uploadingField.value = field.field_key
  formError.value = ''
  try {
    const res = await cms.uploadMedia(file, 'application')
    const uploaded = unwrap(res)
    formAnswers[field.field_key] = uploaded.url || uploaded.file_url || uploaded.path
  } catch (err) {
    formError.value = apiError(err, 'Could not upload this file.')
  } finally {
    uploadingField.value = ''
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
  } catch (err) { error.value = apiError(err, 'Could not update your profile.') }
  finally { savingProfile.value = false }
}
function closeApplication() { activeForm.value = null; activeSubmission.value = null; clearAnswers(); formError.value = '' }
function clearAnswers() { Object.keys(formAnswers).forEach(key => delete formAnswers[key]) }
function logout() { localStorage.removeItem('ncsms_access_token'); localStorage.removeItem('ncsms_user'); router.push('/login') }
function unwrap(value) { return value?.data?.data ?? value?.data ?? value ?? {} }
function asList(value) { const data = unwrap(value); return Array.isArray(data) ? data : [] }
function parseObject(value) { if (!value) return {}; if (typeof value === 'object') return value; try { return JSON.parse(value) } catch { return {} } }
function apiError(err, fallback) { return err.response?.data?.error?.message || fallback }
function titleize(value) { return String(value || '').toLowerCase().replaceAll('_', ' ').replaceAll('-', ' ').replace(/\b\w/g, char => char.toUpperCase()) }
function formatDate(value) { return value ? new Intl.DateTimeFormat('en-UG', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) : '-' }
function formatMoney(value) { return new Intl.NumberFormat('en-UG', { maximumFractionDigits: 0 }).format(Number(value || 0)) }
function statusClass(status) { if (status === 'APPROVED') return 'green'; if (status === 'REJECTED') return 'red'; if (['NEEDS_INFORMATION', 'PENDING_PAYMENT'].includes(status)) return 'amber'; return 'blue' }
function paymentClass(status) { return status === 'PAID' || status === 'VERIFIED' ? 'green' : status === 'REJECTED' ? 'red' : 'amber' }
function notificationIcon(item) { return item.icon_key ? `icofont-${item.icon_key}` : 'icofont-notification' }
function activityTitle(item) { return titleize(item.action || item.event_type || 'Portal activity') }
function activityDescription(item) { return item.endpoint ? `${item.method || 'Action'} ${item.endpoint}` : `Activity recorded from ${item.ip_address || 'your account'}.` }
function inputType(type) { return ({ phone: 'tel', email: 'email', number: 'number', date: 'date' })[type] || 'text' }
function fieldConfig(field) { return parseObject(field.config) }
function fieldOptions(field) { return Array.isArray(fieldConfig(field).options) ? fieldConfig(field).options : [] }
function fieldAccept(field) { const accept = fieldConfig(field).accept; return Array.isArray(accept) ? accept.join(',') : (field.field_type === 'image' ? 'image/*' : '.pdf,.jpg,.jpeg,.png,.webp') }

</script>

<style scoped>
.user-portal{min-height:100vh;background:#f5f7fb;color:#303345}.portal-sidebar{position:fixed;inset:0 auto 0 0;z-index:50;display:flex;flex-direction:column;width:255px;padding:18px 14px;background:#10233f;color:#fff}.portal-brand{display:flex;align-items:center;gap:11px;padding:0 8px 20px;color:#fff}.portal-brand img{width:42px;height:42px;object-fit:contain;background:#fff;border-radius:5px}.portal-brand span{font-size:16px;font-weight:800}.portal-brand small{display:block;color:#aebbd0;font-size:10px;font-weight:600}.portal-sidebar nav{display:grid;gap:4px;overflow:auto}.portal-sidebar nav button,.logout-button{display:grid;grid-template-columns:24px 1fr auto;align-items:center;gap:8px;width:100%;min-height:44px;padding:9px 11px;border:0;border-radius:5px;background:transparent;color:#c6d0df;text-align:left}.portal-sidebar nav button:hover,.portal-sidebar nav button.active{background:#f5a623;color:#182943}.portal-sidebar nav i,.logout-button i{font-size:18px}.portal-sidebar nav b{display:grid;place-items:center;min-width:20px;height:20px;padding:0 5px;border-radius:10px;background:#fff;color:#10233f;font-size:10px}.logout-button{margin-top:auto;border-top:1px solid rgba(255,255,255,.12);border-radius:0}.portal-workspace{min-height:100vh;margin-left:255px}.portal-navbar{position:sticky;top:0;z-index:30;display:flex;align-items:center;gap:14px;height:70px;padding:0 28px;border-bottom:1px solid #e8eaf0;background:#fff}.menu-button{display:none!important}.navbar-title small,.navbar-title strong{display:block}.navbar-title small{color:#7f879a;font-size:10px;font-weight:800;text-transform:uppercase}.navbar-title strong{color:#303345}.navbar-actions{display:flex;align-items:center;gap:7px;margin-left:auto}.navbar-actions>button,.menu-button{position:relative;place-items:center;width:38px;height:38px;border:1px solid #e1e4eb;border-radius:5px;background:#fff;color:#50566a}.navbar-actions>button{display:grid}.navbar-actions>button>span{position:absolute;top:-5px;right:-5px;display:grid;place-items:center;min-width:17px;height:17px;padding:0 4px;border-radius:9px;background:#e84f4f;color:#fff;font-size:9px}.navbar-actions .profile-button{display:flex;width:auto;min-width:112px;padding:4px 9px;gap:7px}.profile-button img,.profile-button b{display:grid;place-items:center;width:27px;height:27px;object-fit:cover;border-radius:50%;background:#10233f;color:#fff}.profile-button span{position:static!important;background:transparent!important;color:#303345!important;font-size:12px!important}.portal-content{max-width:1480px;margin:0 auto;padding:26px}.page-heading{display:flex;align-items:center;justify-content:space-between;gap:20px;margin-bottom:20px}.page-heading p{margin:0;color:#d17b00;font-size:11px;font-weight:800;text-transform:uppercase}.page-heading h1{margin:3px 0;color:#25283a;font-size:25px}.page-heading span{color:#767d90;font-size:13px}.primary-command,.secondary-command{display:inline-flex;align-items:center;justify-content:center;gap:7px;min-height:40px;padding:0 14px;border:0;border-radius:5px;background:#10233f;color:#fff;font-size:12px;font-weight:800}.secondary-command{border:1px solid #dfe2ea;background:#fff;color:#4e5569}.user-kpis{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:14px;margin-bottom:22px}.user-kpis article{display:flex;align-items:center;gap:12px;min-height:112px;padding:16px;border:1px solid #e4e7ee;border-radius:6px;background:#fff}.user-kpis article>span{display:grid;place-items:center;flex:0 0 46px;width:46px;height:46px;border-radius:5px;font-size:22px}.user-kpis .blue{background:#e8edff;color:#5266d8}.user-kpis .amber{background:#fff1df;color:#bd6c12}.user-kpis .green{background:#e5f8ee;color:#218b55}.user-kpis .cyan{background:#e2f6f8;color:#197f8d}.user-kpis small,.user-kpis strong,.user-kpis p{display:block;margin:0}.user-kpis small{color:#7d8497;font-size:11px;font-weight:800;text-transform:uppercase}.user-kpis strong{color:#292c3e;font-size:25px}.user-kpis p{color:#9297a7;font-size:11px}.portal-success,.portal-error{margin:0 0 15px;padding:11px 13px;border-radius:5px;font-size:12px}.portal-success{border:1px solid #a7f3d0;background:#ecfdf5;color:#047857}.portal-error{border:1px solid #fecaca;background:#fef2f2;color:#b91c1c}.open-forms-section{margin-top:8px}.forms-empty,.empty-state{padding:48px;border:1px dashed #ccd2de;border-radius:6px;background:#fff;text-align:center;color:#7c8497}.list-toolbar{display:grid;grid-template-columns:1fr 210px;gap:10px;margin-bottom:12px}.list-toolbar label{display:flex;align-items:center;gap:8px}.list-toolbar label,.list-toolbar select{height:41px;padding:0 11px;border:1px solid #dfe2ea;border-radius:5px;background:#fff}.list-toolbar input{width:100%;border:0;outline:0}.data-table{overflow:auto;border:1px solid #e2e5ec;border-radius:6px;background:#fff}.data-table table{width:100%;min-width:850px;border-collapse:collapse}.data-table th,.data-table td{padding:13px 15px;border-bottom:1px solid #edf0f5;text-align:left;font-size:12px}.data-table th{background:#f8f9fc;color:#646b7e;font-weight:800;text-transform:uppercase}.data-table td{color:#61687b}.data-table td strong,.data-table td small{display:block}.data-table td strong{color:#303345}.data-table td small{color:#8e94a5}.data-table td button{display:grid;place-items:center;width:32px;height:32px;border:1px solid #dfe2ea;border-radius:4px;background:#fff}.empty-cell{padding:45px!important;text-align:center!important}.status{display:inline-flex;padding:4px 7px;border-radius:4px;font-size:10px;font-weight:800;text-transform:uppercase}.status.blue{background:#e8edff;color:#5266d8}.status.green{background:#e5f8ee;color:#218b55}.status.amber{background:#fff1df;color:#bd6c12}.status.red{background:#fee9e8;color:#cf433d}.timeline-list,.feed-list{display:grid;gap:9px}.timeline-list article,.feed-list article{display:flex;align-items:flex-start;gap:12px;padding:15px;border:1px solid #e3e6ed;border-radius:6px;background:#fff}.timeline-list article>span,.feed-list article>span{display:grid;place-items:center;flex:0 0 38px;width:38px;height:38px;border-radius:5px;background:#e8edff;color:#5266d8}.timeline-list strong,.feed-list strong{color:#303345}.timeline-list p,.feed-list p{margin:3px 0;color:#70778a;font-size:12px}.timeline-list small,.feed-list small{color:#969bab;font-size:10px}.feed-list article.unread{border-left:3px solid #6777ef}.feed-list article>button{margin-left:auto;border:0;background:transparent;color:#5266d8}.empty-state{display:block!important}.empty-state>i{font-size:32px}.empty-state h3{margin:9px 0 0;color:#303345}.transaction-summary{display:flex;align-items:center;justify-content:space-between;margin-bottom:12px;padding:15px 18px;border-left:3px solid #47c363;background:#fff}.transaction-summary span{color:#73798c}.transaction-summary strong{color:#303345;font-size:19px}.profile-layout{display:grid;grid-template-columns:270px minmax(0,1fr);gap:16px}.profile-summary,.profile-form{padding:22px;border:1px solid #e3e6ed;border-radius:6px;background:#fff}.profile-summary{text-align:center}.profile-avatar{display:grid;place-items:center;width:88px;height:88px;margin:0 auto 12px;overflow:hidden;border-radius:50%;background:#10233f;color:#fff;font-size:25px;font-weight:800}.profile-avatar img{width:100%;height:100%;object-fit:cover}.profile-summary h2,.profile-form h2{margin:0;color:#303345;font-size:18px}.profile-summary p{color:#737a8e}.ordinary-badge{display:inline-flex;padding:5px 8px;border-radius:4px;background:#e8edff;color:#5266d8;font-size:10px;font-weight:800;text-transform:uppercase}.profile-form{display:grid;gap:14px}.profile-form>div{display:grid;grid-template-columns:1fr 1fr;gap:12px}.profile-form label{display:grid;gap:5px;color:#555c70;font-size:11px;font-weight:800}.profile-form input,.form-field input,.form-field textarea,.form-field select,.payment-section input{width:100%;padding:10px;border:1px solid #d9dde7;border-radius:5px;background:#fff;color:#303345}.profile-form input:disabled{background:#f2f4f8}.profile-form button{justify-self:start}.form-backdrop{position:fixed;inset:0;z-index:100;display:grid;place-items:center;padding:18px;background:rgba(16,35,63,.6)}.application-modal{width:min(760px,100%);max-height:94vh;overflow:auto;border-radius:6px;background:#fff;box-shadow:0 20px 60px rgba(0,0,0,.25)}.application-modal>header{display:flex;justify-content:space-between;gap:18px;padding:22px;border-bottom:1px solid #e7e9ef}.application-modal>header small{color:#d17b00;font-size:10px;font-weight:800;text-transform:uppercase}.application-modal>header h2{margin:3px 0;color:#303345;font-size:21px}.application-modal>header p{margin:0;color:#747b8e;font-size:12px}.application-modal>header button{display:grid;place-items:center;flex:0 0 36px;width:36px;height:36px;border:1px solid #dfe2ea;border-radius:5px;background:#fff}.application-modal>form{display:grid;gap:17px;padding:22px}.fee-banner{display:flex;align-items:center;justify-content:space-between;padding:12px 14px;border-left:3px solid #f5a623;background:#fff8e9}.fee-banner span{color:#6e7588}.fee-banner strong{color:#303345}.form-field{display:grid;gap:6px}.form-field>label,.payment-section label{color:#4b5266;font-size:12px;font-weight:800}.form-field>label b{color:#cf433d}.form-field>p{margin:0;color:#858b9d;font-size:11px}.choice-list{display:grid;gap:8px}.choice-list label{display:flex;align-items:center;gap:7px;color:#5d6477;font-size:12px}.choice-list input{width:auto}.file-input{display:flex;align-items:center;gap:10px}.file-input a{color:#5266d8;font-size:12px}.payment-section{display:grid;gap:11px;padding:15px;border:1px solid #f1d39e;border-radius:5px;background:#fffaf0}.payment-section h3{margin:0;color:#303345;font-size:15px}.payment-section p{margin:0;color:#73798b;font-size:11px}.payment-section label{display:grid;gap:5px}.application-modal form>footer{display:flex;justify-content:flex-end;gap:9px;padding-top:15px;border-top:1px solid #e7e9ef}.form-loading{padding:70px;text-align:center;color:#72798d}.nav-scrim{display:none}.user-portal :deep(button:disabled){cursor:not-allowed;opacity:.55}@media(max-width:1120px){.user-kpis{grid-template-columns:repeat(2,minmax(0,1fr))}}@media(max-width:820px){.portal-sidebar{transform:translateX(-100%);transition:transform .2s ease}.portal-sidebar.open{transform:translateX(0)}.portal-workspace{margin-left:0}.nav-scrim{position:fixed;inset:0;z-index:40;display:block;border:0;background:rgba(15,23,42,.48)}.menu-button{display:grid!important}.portal-navbar{padding:0 16px}.portal-content{padding:20px 16px}.profile-layout{grid-template-columns:1fr}}@media(max-width:600px){.navbar-title{display:none}.profile-button span{display:none!important}.navbar-actions .profile-button{min-width:38px;width:38px;padding:4px}.page-heading{align-items:flex-start;flex-direction:column}.user-kpis,.list-toolbar,.profile-form>div{grid-template-columns:1fr}.application-modal{max-height:100vh}.form-backdrop{padding:0}.application-modal>form,.application-modal>header{padding:17px}.application-modal form>footer{flex-direction:column-reverse}.application-modal form>footer button{width:100%}}
.user-portal{min-height:100vh!important;background:#f4f6f9!important;color:#34395e!important}.user-portal button{font-family:inherit}.user-portal .main-wrapper{min-height:100vh!important;display:flex!important;flex-direction:column!important}.user-portal .navbar-bg{position:fixed!important;top:0!important;right:0!important;left:260px!important;z-index:1030!important;height:70px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.user-portal .main-navbar{position:fixed!important;top:0!important;right:0!important;left:260px!important;z-index:1040!important;min-height:70px!important;background:#fff!important;color:#111827!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.user-portal .main-navbar .nav-link,.user-portal .main-navbar i{color:#111827!important}.user-portal .main-sidebar{position:fixed!important;left:0!important;top:0!important;width:260px!important;height:100vh!important;overflow:hidden!important}.user-portal .sidebar-mini .main-sidebar{width:65px!important}.user-portal #sidebar-wrapper{display:flex!important;flex-direction:column!important;height:100vh!important;min-height:0!important}.user-portal .sidebar-brand{display:flex!important;align-items:center!important;justify-content:center!important;flex:0 0 70px!important;height:70px!important}.user-portal .sidebar-brand a{display:flex!important;align-items:center!important;justify-content:center!important;width:100%!important}.user-portal .header-logo{max-width:84px!important;max-height:48px!important;object-fit:contain!important}.user-portal .sidebar-user{display:flex!important;gap:10px!important;align-items:center!important;flex:0 0 auto!important;min-width:0!important;margin:8px 14px 12px!important;padding:12px!important;border-radius:8px!important;background:#f8f9fa!important}.user-portal .sidebar-user img,.sidebar-avatar-fallback{width:38px!important;height:38px!important;flex:0 0 38px!important;border-radius:50%!important;object-fit:cover!important}.sidebar-avatar-fallback,.nav-avatar-fallback{display:inline-flex!important;align-items:center!important;justify-content:center!important;background:#6777ef!important;color:#fff!important;font-weight:700!important}.user-portal .sidebar-user div{min-width:0!important}.user-portal .sidebar-user strong{display:block!important;max-width:150px!important;overflow:hidden!important;text-overflow:ellipsis!important;white-space:nowrap!important;color:#34395e!important;font-size:13px!important}.user-portal .sidebar-user span:last-child{display:block!important;color:#98a6ad!important;font-size:11px!important}.user-portal .sidebar-menu{flex:1 1 auto!important;min-height:0!important;overflow-y:auto!important;overflow-x:hidden!important;padding-bottom:18px!important;scrollbar-width:thin!important}.user-portal .sidebar-menu::-webkit-scrollbar{width:6px}.user-portal .sidebar-menu::-webkit-scrollbar-thumb{background:#d7dbea;border-radius:999px}.user-portal .sidebar-menu button{border:0!important;background:transparent!important;width:100%!important;min-width:0!important;text-align:left!important}.user-portal .sidebar-menu .nav-link{position:relative!important;display:flex!important;align-items:center!important;gap:10px!important;width:100%!important;height:auto!important;min-height:38px!important;padding:9px 18px!important;line-height:1.2!important}.user-portal .sidebar-menu .nav-link i{flex:0 0 18px!important;width:18px!important;text-align:center!important}.user-portal .sidebar-menu .nav-link span{min-width:0!important;overflow:hidden!important;text-overflow:ellipsis!important;white-space:nowrap!important;font-size:12px!important}.user-portal .sidebar-menu li.active>button{color:#6777ef!important;font-weight:600!important}.user-portal .sidebar-menu li.active>button i{color:#6777ef!important}.portal-nav-badge{display:inline-flex!important;align-items:center!important;justify-content:center!important;min-width:20px!important;height:20px!important;margin-left:auto!important;padding:0 6px!important;border-radius:999px!important;background:#ffa426!important;color:#fff!important;font-size:10px!important;line-height:1!important}.user-portal .main-wrapper.sidebar-mini .navbar-bg,.user-portal .main-wrapper.sidebar-mini .main-navbar{left:65px!important}.user-portal .main-content{flex:1 0 auto!important;margin-left:260px!important;padding:92px 30px 30px!important;background:#f4f6f9!important}.user-portal .main-wrapper.sidebar-mini .main-content{margin-left:65px!important}.user-portal .cms-main-footer{flex:0 0 auto!important;margin-top:auto!important;margin-left:260px!important;border-top:1px solid #e4e6fc!important;background:#fff!important;color:#6c757d!important}.user-portal .main-wrapper.sidebar-mini .cms-main-footer{margin-left:65px!important}.user-portal .cms-main-footer .footer-left{font-weight:700!important;color:#34395e!important}.main-navbar button.nav-link{border:0!important;background:transparent!important}.cms-top-icon{display:inline-flex!important;align-items:center!important;justify-content:center!important;color:#111827!important}.cms-top-icon i{color:#111827!important;font-size:18px!important}.portal-navbar-title{display:flex!important;flex-direction:column!important;justify-content:center!important;min-height:48px!important;margin-left:4px!important;line-height:1.15!important}.portal-navbar-title small{color:#98a6ad!important;font-size:10px!important;font-weight:800!important;text-transform:uppercase!important}.portal-navbar-title strong{color:#34395e!important;font-size:15px!important}.portal-top-action{position:relative!important}.headerBadge1,.headerBadge2{position:absolute!important;top:8px!important;right:4px!important;background:#ffa426!important;color:#fff!important}.nav-avatar-fallback{width:30px!important;height:30px!important;border-radius:50%!important;font-size:11px!important}.navbar .dropdown-menu.show{display:block!important;position:absolute!important}.navbar .dropdown-item,.navbar .dropdown-item.has-icon{width:100%!important;border:0!important;background:transparent!important;text-align:left!important}.portal-theme-toggle{border:0!important;background:transparent!important}.otika-page-actions{display:flex!important;justify-content:flex-end!important;gap:10px!important;margin:-12px 0 18px!important}.otika-page-actions .btn{border:0!important;border-radius:30px!important;color:#fff!important;padding:8px 18px!important;font-size:12px!important;font-weight:600!important}.otika-page-actions .btn-primary{background:#6777ef!important;box-shadow:0 2px 6px #acb5f6!important}.otika-page-actions .btn-info{background:#3abaf4!important;box-shadow:0 2px 6px rgba(58,186,244,.35)!important}.cms-message,.cms-error{margin:0 0 15px!important;padding:12px 14px!important;border:0!important;border-radius:3px!important;font-size:12px!important;box-shadow:0 4px 25px rgba(0,0,0,.05)!important}.cms-message{background:#e8f7f0!important;color:#47c363!important}.cms-error{background:#fdeaea!important;color:#fc544b!important}.page-heading{display:flex!important;align-items:center!important;justify-content:space-between!important;gap:20px!important;margin-bottom:20px!important}.page-heading p{margin:0!important;color:#6777ef!important;font-size:11px!important;font-weight:800!important;text-transform:uppercase!important}.page-heading h1{margin:3px 0!important;color:#34395e!important;font-size:24px!important;font-weight:700!important}.page-heading span{color:#6c757d!important;font-size:13px!important}.primary-command,.secondary-command{display:inline-flex!important;align-items:center!important;justify-content:center!important;gap:7px!important;min-height:40px!important;padding:0 16px!important;border:0!important;border-radius:30px!important;background:#6777ef!important;color:#fff!important;font-size:12px!important;font-weight:600!important;box-shadow:0 2px 6px #acb5f6!important}.secondary-command{border:0!important;background:#f4f6f9!important;color:#34395e!important;box-shadow:none!important}.user-kpis{display:grid!important;grid-template-columns:repeat(4,minmax(0,1fr))!important;gap:18px!important;margin-bottom:22px!important}.user-kpis article{display:flex!important;align-items:center!important;gap:14px!important;min-height:124px!important;padding:18px!important;border:0!important;border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.1)!important}.user-kpis article>span{display:grid!important;place-items:center!important;flex:0 0 54px!important;width:54px!important;height:54px!important;border-radius:50%!important;font-size:25px!important}.user-kpis .blue{background:#e8edff!important;color:#6777ef!important}.user-kpis .amber{background:#fff4e6!important;color:#ffa426!important}.user-kpis .green{background:#e8f7f0!important;color:#47c363!important}.user-kpis .cyan{background:#eaf7ff!important;color:#3abaf4!important}.user-kpis small,.user-kpis strong,.user-kpis p{display:block!important;margin:0!important}.user-kpis small{color:#98a6ad!important;font-size:12px!important;font-weight:700!important;text-transform:uppercase!important}.user-kpis strong{color:#34395e!important;font-size:28px!important;font-weight:700!important}.user-kpis p{color:#6c757d!important;font-size:12px!important}.forms-empty,.empty-state{padding:48px!important;border:1px dashed #e4e6fc!important;border-radius:3px!important;background:#fdfdff!important;text-align:center!important;color:#98a6ad!important}.list-toolbar{display:grid!important;grid-template-columns:minmax(0,1fr) 210px!important;gap:10px!important;margin-bottom:12px!important}.list-toolbar label{display:flex!important;align-items:center!important;gap:8px!important}.list-toolbar label,.list-toolbar select{height:42px!important;padding:0 12px!important;border:1px solid #e4e6fc!important;border-radius:3px!important;background:#fdfdff!important;color:#495057!important}.list-toolbar input{width:100%!important;border:0!important;outline:0!important;background:transparent!important}.data-table{overflow:auto!important;border:0!important;border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.data-table table{width:100%!important;min-width:850px!important;border-collapse:collapse!important}.data-table th,.data-table td{padding:14px 18px!important;border-bottom:1px solid #f4f6f9!important;text-align:left!important;font-size:12px!important}.data-table th{background:#fff!important;color:#34395e!important;font-weight:700!important;text-transform:uppercase!important}.data-table td{color:#6c757d!important}.data-table td strong,.data-table td small{display:block!important}.data-table td strong{color:#34395e!important}.data-table td small{color:#98a6ad!important}.data-table td button{display:grid!important;place-items:center!important;width:32px!important;height:32px!important;border:0!important;border-radius:50%!important;background:#f4f6f9!important;color:#6777ef!important}.empty-cell{padding:45px!important;text-align:center!important}.status{display:inline-flex!important;padding:4px 8px!important;border-radius:30px!important;font-size:10px!important;font-weight:800!important;text-transform:uppercase!important}.status.blue{background:#e8edff!important;color:#6777ef!important}.status.green{background:#e8f7f0!important;color:#47c363!important}.status.amber{background:#fff4e6!important;color:#ffa426!important}.status.red{background:#fdeaea!important;color:#fc544b!important}.timeline-list,.feed-list{display:grid!important;gap:12px!important}.timeline-list article,.feed-list article{display:flex!important;align-items:flex-start!important;gap:12px!important;padding:18px!important;border:0!important;border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.timeline-list article>span,.feed-list article>span{display:grid!important;place-items:center!important;flex:0 0 38px!important;width:38px!important;height:38px!important;border-radius:50%!important;background:#e8edff!important;color:#6777ef!important}.timeline-list strong,.feed-list strong{color:#34395e!important}.timeline-list p,.feed-list p{margin:3px 0!important;color:#6c757d!important;font-size:12px!important}.timeline-list small,.feed-list small{color:#98a6ad!important;font-size:10px!important}.feed-list article.unread{border-left:3px solid #6777ef!important}.feed-list article>button{margin-left:auto!important;border:0!important;background:transparent!important;color:#6777ef!important}.empty-state{display:block!important}.empty-state>i{font-size:32px!important}.empty-state h3{margin:9px 0 0!important;color:#34395e!important}.transaction-summary{display:flex!important;align-items:center!important;justify-content:space-between!important;margin-bottom:12px!important;padding:18px 20px!important;border-left:3px solid #47c363!important;border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.08)!important}.transaction-summary span{color:#6c757d!important}.transaction-summary strong{color:#34395e!important;font-size:19px!important}.profile-layout{display:grid!important;grid-template-columns:270px minmax(0,1fr)!important;gap:18px!important}.profile-summary,.profile-form{padding:25px!important;border:0!important;border-radius:3px!important;background:#fff!important;box-shadow:0 4px 25px rgba(0,0,0,.1)!important}.profile-summary{text-align:center!important}.profile-avatar{display:grid!important;place-items:center!important;width:88px!important;height:88px!important;margin:0 auto 12px!important;overflow:hidden!important;border-radius:50%!important;background:#6777ef!important;color:#fff!important;font-size:25px!important;font-weight:800!important}.profile-avatar img{width:100%!important;height:100%!important;object-fit:cover!important}.profile-summary h2,.profile-form h2{margin:0!important;color:#34395e!important;font-size:18px!important}.profile-summary p{color:#6c757d!important}.ordinary-badge{display:inline-flex!important;padding:5px 10px!important;border-radius:30px!important;background:#e8edff!important;color:#6777ef!important;font-size:10px!important;font-weight:800!important;text-transform:uppercase!important}.profile-form{display:grid!important;gap:14px!important}.profile-form>div{display:grid!important;grid-template-columns:1fr 1fr!important;gap:12px!important}.profile-form label{display:grid!important;gap:7px!important;color:#34395e!important;font-size:12px!important;font-weight:600!important}.profile-form input,.form-field input,.form-field textarea,.form-field select,.payment-section input{width:100%!important;padding:10px 15px!important;border:1px solid #e4e6fc!important;border-radius:3px!important;background:#fdfdff!important;color:#495057!important;box-shadow:none!important}.profile-form input:disabled{background:#f4f6f9!important}.profile-form button{justify-self:start!important}.form-backdrop{position:fixed!important;inset:0!important;z-index:1080!important;display:grid!important;place-items:center!important;padding:18px!important;background:rgba(17,24,39,.52)!important}.application-modal{width:min(760px,100%)!important;max-height:94vh!important;overflow:auto!important;border-radius:3px!important;background:#fff!important;box-shadow:0 20px 60px rgba(0,0,0,.25)!important}.application-modal>header{display:flex!important;justify-content:space-between!important;gap:18px!important;padding:24px 25px!important;border-bottom:1px solid #f9f9f9!important}.application-modal>header small{color:#6777ef!important;font-size:10px!important;font-weight:800!important;text-transform:uppercase!important}.application-modal>header h2{margin:3px 0!important;color:#34395e!important;font-size:21px!important}.application-modal>header p{margin:0!important;color:#6c757d!important;font-size:12px!important}.application-modal>header button{display:grid!important;place-items:center!important;flex:0 0 36px!important;width:36px!important;height:36px!important;border:0!important;border-radius:50%!important;background:#f4f6f9!important;color:#34395e!important}.application-modal>form{display:grid!important;gap:17px!important;padding:25px!important}.fee-banner{display:flex!important;align-items:center!important;justify-content:space-between!important;padding:12px 14px!important;border-left:3px solid #ffa426!important;background:#fff4e6!important}.fee-banner span{color:#6c757d!important}.fee-banner strong{color:#34395e!important}.form-field{display:grid!important;gap:6px!important}.form-field>label,.payment-section label{color:#34395e!important;font-size:12px!important;font-weight:600!important}.form-field>label b{color:#fc544b!important}.form-field>p{margin:0!important;color:#98a6ad!important;font-size:11px!important}.choice-list{display:grid!important;gap:8px!important}.choice-list label{display:flex!important;align-items:center!important;gap:7px!important;color:#6c757d!important;font-size:12px!important}.choice-list input{width:auto!important}.file-input{display:flex!important;align-items:center!important;gap:10px!important}.file-input a{color:#6777ef!important;font-size:12px!important}.payment-section{display:grid!important;gap:11px!important;padding:15px!important;border:1px solid #ffe2ad!important;border-radius:3px!important;background:#fffaf0!important}.payment-section h3{margin:0!important;color:#34395e!important;font-size:15px!important}.payment-section p{margin:0!important;color:#6c757d!important;font-size:11px!important}.payment-section label{display:grid!important;gap:5px!important}.application-modal form>footer{display:flex!important;justify-content:flex-end!important;gap:9px!important;padding-top:15px!important;border-top:1px solid #f9f9f9!important}.form-loading{padding:70px!important;text-align:center!important;color:#6c757d!important}.user-portal :deep(button:disabled){cursor:not-allowed!important;opacity:.55!important}.user-portal button:focus,.user-portal a:focus,.user-portal input:focus,.user-portal textarea:focus,.user-portal select:focus{outline:none!important}.user-portal button:focus-visible,.user-portal a:focus-visible,.user-portal input:focus-visible,.user-portal textarea:focus-visible,.user-portal select:focus-visible{outline:2px solid rgba(103,119,239,.45)!important;outline-offset:2px!important}.user-portal .sidebar-menu button:focus-visible,.user-portal .main-navbar button:focus-visible{outline:0!important;box-shadow:0 0 0 3px rgba(103,119,239,.22)!important}.user-portal .sidebar-mini .main-sidebar,.user-portal .sidebar-mini #sidebar-wrapper{overflow:visible!important}.user-portal .sidebar-mini .sidebar-menu{overflow-y:auto!important;overflow-x:visible!important}.user-portal .sidebar-mini .sidebar-user strong,.user-portal .sidebar-mini .sidebar-user span:not(.sidebar-avatar-fallback),.user-portal .sidebar-mini .sidebar-menu .menu-header,.user-portal .sidebar-mini .portal-nav-badge{display:none!important}.user-portal .sidebar-mini .sidebar-user{justify-content:center!important;margin:8px 8px 12px!important;padding:10px 6px!important}.user-portal .sidebar-mini .sidebar-menu .nav-link{justify-content:center!important;padding:10px!important}.user-portal .sidebar-mini .sidebar-menu .nav-link i{margin:0!important}.user-portal .sidebar-mini .sidebar-menu .nav-link span{display:block!important;position:absolute!important;left:62px!important;top:50%!important;z-index:1200!important;max-width:220px!important;padding:8px 12px!important;border-radius:4px!important;background:#111827!important;color:#fff!important;box-shadow:0 8px 24px rgba(15,23,42,.22)!important;font-size:12px!important;line-height:1!important;opacity:0!important;pointer-events:none!important;transform:translateY(-50%) translateX(-6px)!important;transition:opacity 140ms ease,transform 140ms ease!important;visibility:hidden!important}.user-portal .sidebar-mini .sidebar-menu .nav-link:hover span,.user-portal .sidebar-mini .sidebar-menu .nav-link:focus-visible span{opacity:1!important;transform:translateY(-50%) translateX(0)!important;visibility:visible!important}.user-portal .sidebar-mini .sidebar-menu .nav-link span::before{content:"";position:absolute;left:-5px;top:50%;width:10px;height:10px;background:#111827;transform:translateY(-50%) rotate(45deg)}:global(.dark) .user-portal{background:#0f172a!important;color:#e5e7eb!important}:global(.dark) .user-portal .main-content,:global(.dark) .user-portal .section-body{background:#0f172a!important;color:#e5e7eb!important}:global(.dark) .user-portal .navbar-bg,:global(.dark) .user-portal .main-navbar{background:#111827!important;color:#f8fafc!important;box-shadow:0 4px 24px rgba(0,0,0,.35)!important}:global(.dark) .user-portal .main-navbar .nav-link,:global(.dark) .user-portal .main-navbar i,:global(.dark) .user-portal .cms-top-icon,:global(.dark) .user-portal .cms-top-icon i{color:#f8fafc!important}:global(.dark) .user-portal .main-sidebar,:global(.dark) .user-portal #sidebar-wrapper{background:#111827!important}:global(.dark) .user-portal .sidebar-brand,:global(.dark) .user-portal .sidebar-user,:global(.dark) .user-portal .cms-main-footer,:global(.dark) .user-portal .user-kpis article,:global(.dark) .user-portal .data-table,:global(.dark) .user-portal .timeline-list article,:global(.dark) .user-portal .feed-list article,:global(.dark) .user-portal .transaction-summary,:global(.dark) .user-portal .profile-summary,:global(.dark) .user-portal .profile-form,:global(.dark) .user-portal .application-modal{background:#1f2937!important;border-color:#334155!important;color:#cbd5e1!important}:global(.dark) .user-portal .sidebar-user strong,:global(.dark) .user-portal .section-header h1,:global(.dark) .user-portal .page-heading h1,:global(.dark) .user-portal .user-kpis strong,:global(.dark) .user-portal .data-table th,:global(.dark) .user-portal .data-table td strong,:global(.dark) .user-portal .timeline-list strong,:global(.dark) .user-portal .feed-list strong,:global(.dark) .user-portal .empty-state h3,:global(.dark) .user-portal .transaction-summary strong,:global(.dark) .user-portal .profile-summary h2,:global(.dark) .user-portal .profile-form h2,:global(.dark) .user-portal .profile-form label,:global(.dark) .user-portal .application-modal h2,:global(.dark) .user-portal .payment-section h3,:global(.dark) .user-portal .cms-main-footer .footer-left{color:#f8fafc!important}:global(.dark) .user-portal .sidebar-user span,:global(.dark) .user-portal .section-header-breadcrumb,:global(.dark) .user-portal .breadcrumb-item,:global(.dark) .user-portal .page-heading span,:global(.dark) .user-portal .user-kpis p,:global(.dark) .user-portal .data-table td,:global(.dark) .user-portal .timeline-list p,:global(.dark) .user-portal .feed-list p,:global(.dark) .user-portal .profile-summary p{color:#cbd5e1!important}:global(.dark) .user-portal input,:global(.dark) .user-portal textarea,:global(.dark) .user-portal select,:global(.dark) .user-portal .form-control,:global(.dark) .user-portal .list-toolbar label{background:#111827!important;color:#f8fafc!important;border-color:#475569!important}:global(.dark) .user-portal .data-table th,:global(.dark) .user-portal .empty-state,:global(.dark) .user-portal .forms-empty,:global(.dark) .user-portal .profile-form input:disabled{background:#111827!important;border-color:#334155!important}:global(.dark) .user-portal .dropdown-menu{background:#1f2937!important;border-color:#334155!important;color:#e5e7eb!important}:global(.dark) .user-portal .dropdown-title,:global(.dark) .user-portal .dropdown-item,:global(.dark) .user-portal .dropdown-item span{color:#e5e7eb!important}:global(.dark) .user-portal .dropdown-item:hover{background:#111827!important}@media(max-width:1120px){.user-kpis{grid-template-columns:repeat(2,minmax(0,1fr))!important}}@media(max-width:991px){.user-portal .navbar-bg,.user-portal .main-navbar,.user-portal .main-wrapper.sidebar-mini .navbar-bg,.user-portal .main-wrapper.sidebar-mini .main-navbar{left:0!important}.user-portal .navbar-bg{height:116px!important}.user-portal .main-navbar{right:0!important;min-height:70px!important;padding:8px 14px!important;flex-wrap:wrap!important}.user-portal .main-navbar .form-inline.mr-auto{flex:1 1 auto!important;min-width:0!important}.user-portal .main-navbar .navbar-nav{flex-direction:row!important;align-items:center!important;gap:4px!important;margin-right:0!important}.user-portal .navbar-right{margin-left:auto!important;flex-direction:row!important}.user-portal .main-sidebar{position:relative!important;width:100%!important;height:auto!important;left:0!important;top:0!important;z-index:1!important;box-shadow:none!important}.user-portal #sidebar-wrapper{height:auto!important;max-height:none!important}.user-portal .main-content,.user-portal .main-wrapper.sidebar-mini .main-content{margin-left:0!important;padding:24px 16px!important;padding-top:24px!important}.user-portal .cms-main-footer,.user-portal .main-wrapper.sidebar-mini .cms-main-footer{margin-left:0!important}.user-portal .sidebar-menu{display:grid!important;grid-template-columns:repeat(2,minmax(0,1fr))!important;gap:4px 8px!important;max-height:60vh!important;padding-bottom:18px!important}.user-portal .sidebar-menu .menu-header{grid-column:1/-1!important}.user-portal .sidebar-user{margin:8px 16px 14px!important}.user-portal .section-header{align-items:flex-start!important;gap:10px!important;flex-direction:column!important}.user-portal .section-header-breadcrumb{margin-left:0!important}.otika-page-actions{justify-content:flex-start!important;flex-wrap:wrap!important}.profile-layout{grid-template-columns:1fr!important}}@media(max-width:768px){.portal-navbar-title{display:none!important}.navbar .dropdown-menu.show{position:fixed!important;top:62px!important;left:12px!important;right:12px!important;width:auto!important;max-width:none!important}.list-toolbar{grid-template-columns:1fr!important}.profile-form>div{grid-template-columns:1fr!important}}@media(max-width:600px){.user-portal .sidebar-menu{grid-template-columns:1fr!important}.user-portal .main-content{padding:18px 12px!important}.user-portal .section-header h1{font-size:20px!important;line-height:1.25!important}.page-heading{align-items:flex-start!important;flex-direction:column!important}.user-kpis{grid-template-columns:1fr!important}.application-modal{max-height:100vh!important}.form-backdrop{padding:0!important}.application-modal>form,.application-modal>header{padding:17px!important}.application-modal form>footer{flex-direction:column-reverse!important}.application-modal form>footer button{width:100%!important}}
</style>
