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

test('audit log viewer uses friendly fields and hides raw API endpoints', () => {
  const portal = read('../src/views/WebsiteContentManagerView.vue')
  const backendModel = read('../../backend/internal/models/models.go')
  const backendMiddleware = read('../../backend/internal/middleware/middleware.go')
  const router = read('../src/router/index.js')
  const activityAudit = read('../src/services/activityAudit.js')
  const server = read('../../backend/cmd/server/main.go')
  const auditHandler = read('../../backend/internal/handlers/audit.go')

  assert.match(portal, /auditActivityName\(log\)/)
  assert.match(portal, /auditLocation\(log\)/)
  assert.match(portal, /auditClientSummary\(log\)/)
  assert.match(portal, /auditDetailRows\(selectedAuditLog\)/)
  assert.doesNotMatch(portal, /\{\{\s*log\.method\s*\|\|\s*log\.action\s*\}\}\s*\{\{\s*log\.endpoint/)
  assert.doesNotMatch(portal, /<pre>\{\{ selectedAuditLog \}\}<\/pre>/)
  assert.match(backendModel, /Timestamp\s+time\.Time\s+`json:"timestamp"`/)
  assert.match(backendModel, /AccessMedium\s+string\s+`json:"access_medium"`/)
  assert.match(backendModel, /Signature\s+string\s+`json:"signature,omitempty"`/)
  assert.match(backendMiddleware, /canonicalAuditAction\(eventType, method\)/)
  assert.match(router, /installActivityAuditor\(router\)/)
  assert.match(activityAudit, /recordNavigation\(to, from\)/)
  assert.match(activityAudit, /ACTIVITY_ENDPOINT = '\/api\/v1\/account\/activity-events'/)
  assert.match(activityAudit, /document\.addEventListener\('click'/)
  assert.match(activityAudit, /password\|token\|secret\|credential/i)
  assert.match(server, /h\.Audit\.SetWriter\(auditWriter\)/)
  assert.match(server, /r\.Post\("\/account\/activity-events", h\.Audit\.RecordFrontendActivity\)/)
  assert.match(auditHandler, /Action:\s+normalized\.Action/)
  assert.match(auditHandler, /EventType:\s+normalized\.EventType/)
  assert.match(auditHandler, /"ui:" \+ kind/)
})

test('command center uses live system resources and functional service actions', () => {
  const panel = read('../src/components/cms/SystemCommandCenterPanel.vue')
  const operatorApi = read('../src/api/operator.js')
  const routes = read('../../backend/cmd/server/main.go')
  const handler = read('../../backend/internal/handlers/operator.go')
  const compose = read('../../docker-compose.yml')

  assert.match(panel, /getSystemResources/)
  assert.match(panel, /getServiceLogs/)
  assert.match(panel, /runServiceAction/)
  assert.match(panel, /displayServiceStatus\(service = \{\}\)[\s\S]+service\.status \|\| 'unavailable'/)
  assert.match(panel, /canRunServiceAction\(svc, 'restart'\)/)
  assert.match(panel, /dockerUnavailable\.value/)
  assert.match(operatorApi, /\/api\/v1\/admin\/system\/resources/)
  assert.match(operatorApi, /\/api\/v1\/admin\/system\/service-logs/)
  assert.match(operatorApi, /\/api\/v1\/admin\/system\/services\/action/)
  assert.match(routes, /r\.Get\("\/admin\/system\/resources", h\.Operator\.Resources\)/)
  assert.match(routes, /r\.Post\("\/admin\/system\/services\/action", h\.Operator\.ServiceAction\)/)
  assert.match(handler, /opsServiceCatalog\(\)[\s\S]+"worker"/)
  assert.match(handler, /composeProject = "ncsportal"/)
  assert.match(compose, /\/var\/run\/docker\.sock:\/var\/run\/docker\.sock/)
})

test('NAMIS equipment uses a dedicated linked entry page with a complete form', () => {
  const manager = read('../src/components/portal/NamisManagerPanel.vue')
  const entry = read('../src/views/NamisRegistryEntryView.vue')
  const router = read('../src/router/index.js')

  assert.match(router, /\/portal\/namis\/:resource\/new/)
  assert.match(router, /NamisRegistryCreate[^\n]+PortalManagerView/)
  assert.match(manager, /params: \{ resource \}/)
  assert.match(manager, /const resource = props\.tab \|\| activeTab\.value/)
  assert.match(manager, /quantity_received/)
  assert.match(manager, /quantity_distributed/)
  assert.match(entry, /Back to \{\{ listLabel \}\}/)
  assert.match(entry, /createNsmisDomain\(resource\.value/)
  assert.match(entry, /Quantity distributed cannot exceed quantity received/)
  const dashboard = read('../src/views/WebsiteContentManagerView.vue')
  assert.match(dashboard, /active === 'namis-registry-new'/)
  assert.match(dashboard, /useRoute\(\)/)
  assert.match(dashboard, /\(\) => \[route\.name, route\.params\.resource, route\.query\.section\]/)
  assert.match(dashboard, /name === 'NamisRegistryCreate'[\s\S]+active\.value = 'namis-registry-new'/)
  assert.match(dashboard, /name === 'PortalDashboard'[\s\S]+requestedSection[\s\S]+active\.value = sections\.some/)
})

test('every Sports Registry resource has a complete standalone form schema', () => {
  const config = read('../src/utils/namisRegistryConfig.js')
  const entry = read('../src/views/NamisRegistryEntryView.vue')
  for (const resource of ['athletes','clubs','coaches','competitions','competition-results','medals','talent','national-team','technical-officials','medical-records','safeguarding-records','anti-doping','disbursements','accountabilities','equipment']) {
    assert.match(config, new RegExp(`['\"]?${resource.replaceAll('-', '\\-')}['\"]?\\s*:`))
  }
  assert.match(entry, /createNsmisDomain\(resource\.value/)
  assert.match(entry, /loadingReferences/)
  assert.match(entry, /definition\.sensitive/)
})

test('self-registration creates an intentionally roleless ordinary user', () => {
  const authService = read('../../backend/internal/services/auth_service.go')
  const registerSection = authService.slice(authService.indexOf('func (s *AuthService) Register'), authService.indexOf('func (s *AuthService) Login'))

  assert.doesNotMatch(registerSection, /AssignRole/)
  assert.match(registerSection, /u\.Roles = roles/)
})
