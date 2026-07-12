import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import test from 'node:test'
import { isOrdinaryUser, portalDestination, roleNames } from '../src/utils/portalAuth.js'

const read = relativePath => readFileSync(fileURLToPath(new URL(relativePath, import.meta.url)), 'utf8')

test('role-aware routing keeps roleless and legacy user accounts in the applicant portal', () => {
  assert.deepEqual(roleNames({ roles: [{ name: 'admin' }, 'content_manager'] }), ['admin', 'content_manager'])
  assert.equal(isOrdinaryUser({ roles: [] }), true)
  assert.equal(isOrdinaryUser({ roles: [{ name: 'user' }] }), true)
  assert.equal(isOrdinaryUser({ roles: [{ name: 'admin' }] }), false)
  assert.equal(portalDestination({ roles: [] }), '/dashboard')
  assert.equal(portalDestination({ roles: ['super_admin'] }), '/portal')
})

test('ordinary-user dashboard exposes requested navigation and live application flow', () => {
  const view = read('../src/views/UserPortalView.vue')
  const wizard = read('../src/views/ApplicationWizardView.vue')
  const router = read('../src/router/index.js')
  const formsApi = read('../src/api/forms.js')
  const server = read('../../backend/cmd/server/main.go')
  const handler = read('../../backend/internal/handlers/forms.go')

  for (const label of ['Dashboard', 'Apply Now', 'My Applications', 'My Activities', 'Notifications', 'Messages', 'My Transactions']) {
    assert.match(view, new RegExp(label))
  }
  for (const shellClass of ['otika-cms user-portal', 'main-sidebar', 'main-navbar', 'main-footer']) {
    assert.match(view, new RegExp(shellClass))
  }
  assert.match(view, /ensureOtikaStyles/)
  assert.match(view, /ThemeToggle/)
  assert.match(view, /<OpenForms/)
  assert.match(view, /portalListOpenForms/)
  assert.match(view, /ApplicationWizard/)
  assert.match(router, /\/dashboard\/apply\/:slug/)
  assert.match(wizard, /application-wizard-page/)
  assert.match(wizard, /buildSectionSteps/)
  assert.match(wizard, /Step \{\{ currentStepIndex \+ 1 \}\} of \{\{ steps.length \}\}/)
  assert.match(wizard, /portalSaveDraft/)
  assert.match(wizard, /portalUploadPaymentProof/)
  assert.match(wizard, /portalSubmit/)
  assert.match(wizard, /uploadFieldFile/)
  assert.match(formsApi, /portalListSubmissions/)
  assert.match(server, /r\.Get\("\/", h\.Forms\.PortalListSubmissions\)/)
  assert.match(handler, /ListUserSubmissions/)
})

test('admin dashboard and application queue use operational APIs and KPI contracts', () => {
  const portal = read('../src/views/WebsiteContentManagerView.vue')
  const dashboard = read('../src/components/portal/AdminDashboardPanel.vue')
  const applications = read('../src/components/portal/AdminApplicationsPanel.vue')
  const applicationDetail = read('../src/views/ApplicationDetailView.vue')
  const backend = read('../../backend/internal/handlers/dashboard.go')

  assert.match(portal, /<AdminDashboardPanel/)
  assert.match(portal, /<AdminApplicationsPanel/)
  assert.match(dashboard, /getAdminDashboard/)
  for (const key of ['total_applications', 'total_athletes', 'total_profiles', 'total_users', 'open_forms']) {
    assert.match(backend, new RegExp(`"${key}"`))
  }
  assert.match(applications, /adminListApplications/)
  assert.match(applications, /adminListSubmissions/)
  assert.match(applications, /AdminApplicationDetail/)
  assert.match(applicationDetail, /adminReviewApplication/)
  assert.match(applicationDetail, /adminReviewSubmission/)
  assert.match(applicationDetail, /adminUpdateSubmissionPaymentStatus/)
})

test('self-registration creates an intentionally roleless ordinary user', () => {
  const authService = read('../../backend/internal/services/auth_service.go')
  const registerSection = authService.slice(authService.indexOf('func (s *AuthService) Register'), authService.indexOf('func (s *AuthService) Login'))

  assert.doesNotMatch(registerSection, /AssignRole/)
  assert.match(registerSection, /u\.Roles = roles/)
})
