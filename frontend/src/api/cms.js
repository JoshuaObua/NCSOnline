import apiClient from './client.js'

// ── Public CMS (no auth required) ────────────────────────────────

export function listPosts(params = {}) {
  return apiClient.get('/api/v1/cms/posts', { params })
}

export function getPost(slug) {
  return apiClient.get(`/api/v1/cms/posts/${slug}`)
}

export function listBlogCategories(active = true) {
  return apiClient.get('/api/v1/cms/blog/categories', { params: { active: active ? 'true' : 'false' } })
}

export function listPostComments(slug) {
  return apiClient.get(`/api/v1/cms/posts/${slug}/comments`)
}

export function submitPostComment(slug, data) {
  return apiClient.post(`/api/v1/cms/posts/${slug}/comments`, data)
}

export function listEvents(params = {}) {
  return apiClient.get('/api/v1/cms/events', { params })
}

export function getEvent(slug) {
  return apiClient.get(`/api/v1/cms/events/${slug}`)
}

export function listCareers(params = {}) {
  return apiClient.get('/api/v1/cms/careers', { params })
}

export function getCareer(id) {
  return apiClient.get(`/api/v1/cms/careers/${id}`)
}

// ── Admin CMS (auth required) ─────────────────────────────────────

export function adminListPosts(params = {}) {
  return apiClient.get('/api/v1/admin/cms/posts', { params })
}

export function adminCreatePost(data) {
  return apiClient.post('/api/v1/admin/cms/posts', data)
}

export function adminUpdatePost(id, data) {
  return apiClient.put(`/api/v1/admin/cms/posts/${id}`, data)
}

export function adminDeletePost(id) {
  return apiClient.delete(`/api/v1/admin/cms/posts/${id}`)
}

export function adminListBlogCategories() {
  return apiClient.get('/api/v1/admin/cms/blog/categories', { params: { active: 'false' } })
}

export function adminCreateBlogCategory(data) {
  return apiClient.post('/api/v1/admin/cms/blog/categories', data)
}

export function adminUpdateBlogCategory(id, data) {
  return apiClient.put(`/api/v1/admin/cms/blog/categories/${id}`, data)
}

export function adminDeleteBlogCategory(id) {
  return apiClient.delete(`/api/v1/admin/cms/blog/categories/${id}`)
}

export function adminListProjectCategories() {
  return apiClient.get('/api/v1/admin/cms/blog/categories', { params: { active: 'false', content_type: 'project' } })
}

export function adminCreateProjectCategory(data) {
  return apiClient.post('/api/v1/admin/cms/blog/categories', { ...data, content_type: 'project' })
}

export function adminUpdateProjectCategory(id, data) {
  return apiClient.put(`/api/v1/admin/cms/blog/categories/${id}`, { ...data, content_type: 'project' })
}

export function adminDeleteProjectCategory(id) {
  return apiClient.delete(`/api/v1/admin/cms/blog/categories/${id}`)
}

export function adminListCaseStudyCategories() {
  return apiClient.get('/api/v1/admin/cms/blog/categories', { params: { active: 'false', content_type: 'case_study' } })
}

export function adminCreateCaseStudyCategory(data) {
  return apiClient.post('/api/v1/admin/cms/blog/categories', { ...data, content_type: 'case_study' })
}

export function adminUpdateCaseStudyCategory(id, data) {
  return apiClient.put(`/api/v1/admin/cms/blog/categories/${id}`, { ...data, content_type: 'case_study' })
}

export function adminDeleteCaseStudyCategory(id) {
  return apiClient.delete(`/api/v1/admin/cms/blog/categories/${id}`)
}

function adminListScopedCategories(contentType) {
  return apiClient.get('/api/v1/admin/cms/blog/categories', { params: { active: 'false', content_type: contentType } })
}

function adminCreateScopedCategory(contentType, data) {
  return apiClient.post('/api/v1/admin/cms/blog/categories', { ...data, content_type: contentType })
}

function adminUpdateScopedCategory(contentType, id, data) {
  return apiClient.put(`/api/v1/admin/cms/blog/categories/${id}`, { ...data, content_type: contentType })
}

export function adminListFAQCategories() { return adminListScopedCategories('faq') }
export function adminCreateFAQCategory(data) { return adminCreateScopedCategory('faq', data) }
export function adminUpdateFAQCategory(id, data) { return adminUpdateScopedCategory('faq', id, data) }
export function adminDeleteFAQCategory(id) { return apiClient.delete(`/api/v1/admin/cms/blog/categories/${id}`) }
export function adminListResourceCategories() { return adminListScopedCategories('resource') }
export function adminCreateResourceCategory(data) { return adminCreateScopedCategory('resource', data) }
export function adminUpdateResourceCategory(id, data) { return adminUpdateScopedCategory('resource', id, data) }
export function adminDeleteResourceCategory(id) { return apiClient.delete(`/api/v1/admin/cms/blog/categories/${id}`) }
export function adminListCareerCategories() { return adminListScopedCategories('career') }
export function adminCreateCareerCategory(data) { return adminCreateScopedCategory('career', data) }
export function adminUpdateCareerCategory(id, data) { return adminUpdateScopedCategory('career', id, data) }
export function adminDeleteCareerCategory(id) { return apiClient.delete(`/api/v1/admin/cms/blog/categories/${id}`) }
export function adminListTeamDepartments() { return apiClient.get('/api/v1/admin/cms/departments') }
export function adminCreateTeamDepartment(data) { return apiClient.post('/api/v1/admin/cms/departments', data) }
export function adminUpdateTeamDepartment(id, data) { return apiClient.put(`/api/v1/admin/cms/departments/${id}`, data) }
export function adminDeleteTeamDepartment(id) { return apiClient.delete(`/api/v1/admin/cms/departments/${id}`) }
export function adminListFacilityCategories() { return adminListScopedCategories('facility') }
export function adminCreateFacilityCategory(data) { return adminCreateScopedCategory('facility', data) }
export function adminUpdateFacilityCategory(id, data) { return adminUpdateScopedCategory('facility', id, data) }
export function adminDeleteFacilityCategory(id) { return apiClient.delete(`/api/v1/admin/cms/blog/categories/${id}`) }
export function adminListEventCategories() { return adminListScopedCategories('event') }
export function adminCreateEventCategory(data) { return adminCreateScopedCategory('event', data) }
export function adminUpdateEventCategory(id, data) { return adminUpdateScopedCategory('event', id, data) }
export function adminDeleteEventCategory(id) { return apiClient.delete(`/api/v1/admin/cms/blog/categories/${id}`) }
export function adminListInvestCategories() { return adminListScopedCategories('investment') }
export function adminCreateInvestCategory(data) { return adminCreateScopedCategory('investment', data) }
export function adminUpdateInvestCategory(id, data) { return adminUpdateScopedCategory('investment', id, data) }
export function adminDeleteInvestCategory(id) { return apiClient.delete(`/api/v1/admin/cms/blog/categories/${id}`) }
export function adminListFederationCategories() { return adminListScopedCategories('federation') }
export function adminCreateFederationCategory(data) { return adminCreateScopedCategory('federation', data) }
export function adminUpdateFederationCategory(id, data) { return adminUpdateScopedCategory('federation', id, data) }
export function adminDeleteFederationCategory(id) { return apiClient.delete(`/api/v1/admin/cms/blog/categories/${id}`) }

export function adminListNewsletterSubscribers(params = {}) {
  return apiClient.get('/api/v1/admin/cms/newsletter/subscribers', { params })
}

export function adminExportNewsletterSubscribers(format = 'csv') {
  return apiClient.get('/api/v1/admin/cms/newsletter/subscribers/export', { params: { format }, responseType: 'blob' })
}

export function adminListComments(params = {}) {
  return apiClient.get('/api/v1/admin/cms/comments', { params })
}

export function adminApproveComment(id) {
  return apiClient.post(`/api/v1/admin/cms/comments/${id}/approve`)
}

export function adminFlagComment(id) {
  return apiClient.post(`/api/v1/admin/cms/comments/${id}/flag`)
}

export function adminDeleteComment(id) {
  return apiClient.delete(`/api/v1/admin/cms/comments/${id}`)
}

export function adminListEvents(params = {}) {
  return apiClient.get('/api/v1/admin/cms/events', { params })
}

export function adminCreateEvent(data) {
  return apiClient.post('/api/v1/admin/cms/events', data)
}

export function adminUpdateEvent(id, data) {
  return apiClient.put(`/api/v1/admin/cms/events/${id}`, data)
}

export function adminDeleteEvent(id) {
  return apiClient.delete(`/api/v1/admin/cms/events/${id}`)
}

export function adminListCareers(params = {}) {
  return apiClient.get('/api/v1/admin/cms/careers', { params })
}

export function adminCreateCareer(data) {
  return apiClient.post('/api/v1/admin/cms/careers', data)
}

export function adminUpdateCareer(id, data) {
  return apiClient.put(`/api/v1/admin/cms/careers/${id}`, data)
}

export function adminDeleteCareer(id) {
  return apiClient.delete(`/api/v1/admin/cms/careers/${id}`)
}

// ── Slides ────────────────────────────────────────────────────────

export function listSlides(activeOnly = true) {
  return apiClient.get('/api/v1/cms/slides', { params: { active: activeOnly ? 'true' : 'false' } })
}

export function getSlideshow(slug = 'homepage-hero', activeOnly = true) {
  return apiClient.get(`/api/v1/cms/slideshows/${slug}`, { params: { active: activeOnly ? 'true' : 'false' } })
}

export function adminListSlides() {
  return apiClient.get('/api/v1/admin/cms/slides', { params: { active: 'false' } })
}

export function adminGetSlideshow(slug = 'homepage-hero') {
  return apiClient.get(`/api/v1/admin/cms/slideshows/${slug}`, { params: { active: 'false' } })
}

export function adminUpdateSlideshow(slug, data) {
  return apiClient.put(`/api/v1/admin/cms/slideshows/${slug}`, data)
}

export function adminCreateSlide(data) {
  return apiClient.post('/api/v1/admin/cms/slides', data)
}

export function adminUpdateSlide(id, data) {
  return apiClient.put(`/api/v1/admin/cms/slides/${id}`, data)
}

export function adminDeleteSlide(id) {
  return apiClient.delete(`/api/v1/admin/cms/slides/${id}`)
}

// ── Menus ─────────────────────────────────────────────────────────

export function getMenu(name) {
  return apiClient.get(`/api/v1/cms/menus/${name}`)
}

export function adminUpdateMenu(name, items) {
  return apiClient.put(`/api/v1/admin/cms/menus/${name}`, items)
}

// ── Settings (generic key/value JSON store) ──────────────────────

export function getSettings(key) {
  return apiClient.get(`/api/v1/cms/settings/${key}`)
}

export function adminUpdateSettings(key, value) {
  return apiClient.put(`/api/v1/admin/cms/settings/${key}`, value)
}

// ── Media Upload ──────────────────────────────────────────────────

export function uploadMedia(file) {
  const fd = new FormData()
  fd.append('file', file)
  return apiClient.post('/api/v1/media/upload', fd, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// ── Fun Facts ─────────────────────────────────────────────────────

export function listFunFacts() { return apiClient.get('/api/v1/cms/fun-facts') }
export function adminListFunFacts() { return apiClient.get('/api/v1/admin/cms/fun-facts', { params: { active: 'false' } }) }
export function adminCreateFunFact(data) { return apiClient.post('/api/v1/admin/cms/fun-facts', data) }
export function adminUpdateFunFact(id, data) { return apiClient.put(`/api/v1/admin/cms/fun-facts/${id}`, data) }
export function adminDeleteFunFact(id) { return apiClient.delete(`/api/v1/admin/cms/fun-facts/${id}`) }

// ── FAQs ──────────────────────────────────────────────────────────

export function listFAQs(params = {}) { return apiClient.get('/api/v1/cms/faqs', { params }) }
export function adminListFAQs(params = {}) { return apiClient.get('/api/v1/admin/cms/faqs', { params }) }
export function adminCreateFAQ(data) { return apiClient.post('/api/v1/admin/cms/faqs', data) }
export function adminUpdateFAQ(id, data) { return apiClient.put(`/api/v1/admin/cms/faqs/${id}`, data) }
export function adminDeleteFAQ(id) { return apiClient.delete(`/api/v1/admin/cms/faqs/${id}`) }

// ── Resources ─────────────────────────────────────────────────────

export function listResources(params = {}) { return apiClient.get('/api/v1/cms/resources', { params }) }
export function adminListResources(params = {}) { return apiClient.get('/api/v1/admin/cms/resources', { params }) }
export function adminCreateResource(data) { return apiClient.post('/api/v1/admin/cms/resources', data) }
export function adminUpdateResource(id, data) { return apiClient.put(`/api/v1/admin/cms/resources/${id}`, data) }
export function adminDeleteResource(id) { return apiClient.delete(`/api/v1/admin/cms/resources/${id}`) }

// ── Facilities ────────────────────────────────────────────────────

export function listFacilities(params = {}) { return apiClient.get('/api/v1/cms/facilities', { params }) }
export function adminListFacilities() { return apiClient.get('/api/v1/admin/cms/facilities', { params: { active: 'false' } }) }
export function adminCreateFacility(data) { return apiClient.post('/api/v1/admin/cms/facilities', data) }
export function adminUpdateFacility(id, data) { return apiClient.put(`/api/v1/admin/cms/facilities/${id}`, data) }
export function adminDeleteFacility(id) { return apiClient.delete(`/api/v1/admin/cms/facilities/${id}`) }

// ── Associations ──────────────────────────────────────────────────

export function listAssociations(params = {}) { return apiClient.get('/api/v1/cms/associations', { params }) }
export function adminListAssociations() { return apiClient.get('/api/v1/admin/cms/associations', { params: { active: 'false' } }) }
export function adminCreateAssociation(data) { return apiClient.post('/api/v1/admin/cms/associations', data) }
export function adminUpdateAssociation(id, data) { return apiClient.put(`/api/v1/admin/cms/associations/${id}`, data) }
export function adminDeleteAssociation(id) { return apiClient.delete(`/api/v1/admin/cms/associations/${id}`) }

// ── Invest with Us ────────────────────────────────────────────────

export function listInvest() { return apiClient.get('/api/v1/cms/invest') }
export function adminListInvest() { return apiClient.get('/api/v1/admin/cms/invest') }
export function adminCreateInvest(data) { return apiClient.post('/api/v1/admin/cms/invest', data) }
export function adminUpdateInvest(id, data) { return apiClient.put(`/api/v1/admin/cms/invest/${id}`, data) }
export function adminDeleteInvest(id) { return apiClient.delete(`/api/v1/admin/cms/invest/${id}`) }

// ── Team Members ──────────────────────────────────────────────────

export function listTeam() { return apiClient.get('/api/v1/cms/team') }
export function adminListTeam(params = {}) { return apiClient.get('/api/v1/admin/cms/team', { params: { active: 'false', ...params } }) }
export function adminCreateTeam(data) { return apiClient.post('/api/v1/admin/cms/team', data) }
export function adminUpdateTeam(id, data) { return apiClient.put(`/api/v1/admin/cms/team/${id}`, data) }
export function adminDeleteTeam(id) { return apiClient.delete(`/api/v1/admin/cms/team/${id}`) }
export function listInstitutionalDepartments() { return apiClient.get('/api/v1/cms/departments') }

// ── Admin RBAC / Roles ────────────────────────────────────────────────

export function adminListRoles() { return apiClient.get('/api/v1/admin/roles') }
export function adminGetRole(id) { return apiClient.get(`/api/v1/admin/roles/${id}`) }
export function adminCreateRole(data) { return apiClient.post('/api/v1/admin/roles', data) }
export function adminUpdateRole(id, data) { return apiClient.put(`/api/v1/admin/roles/${id}`, data) }
export function adminDeleteRole(id) { return apiClient.delete(`/api/v1/admin/roles/${id}`) }
export function adminListPermissions() { return apiClient.get('/api/v1/admin/permissions') }
export function adminAssignRolePermission(roleId, permissionId) {
  return apiClient.post(`/api/v1/admin/roles/${roleId}/permissions`, { permission_id: permissionId })
}
export function adminRemoveRolePermission(roleId, permissionId) {
  return apiClient.delete(`/api/v1/admin/roles/${roleId}/permissions/${permissionId}`)
}

export function adminListUsers(params = {}) { return apiClient.get('/api/v1/admin/users', { params }) }
export function adminGetUser(id) { return apiClient.get(`/api/v1/admin/users/${id}`) }
export function adminCreateUser(data) { return apiClient.post('/api/v1/admin/users', data) }
export function adminUpdateUser(id, data) { return apiClient.put(`/api/v1/admin/users/${id}`, data) }
export function adminDeleteUser(id) { return apiClient.delete(`/api/v1/admin/users/${id}`) }
export function adminActivateUser(id) { return apiClient.post(`/api/v1/admin/users/${id}/activate`) }
export function adminDeactivateUser(id) { return apiClient.post(`/api/v1/admin/users/${id}/deactivate`) }
export function adminResetUserPassword(id, newPassword) { return apiClient.post(`/api/v1/admin/users/${id}/reset-password`, { new_password: newPassword }) }
export function adminAssignUserRole(userId, roleId) { return apiClient.post(`/api/v1/admin/users/${userId}/roles`, { role_id: roleId }) }
export function adminRemoveUserRole(userId, roleId) { return apiClient.delete(`/api/v1/admin/users/${userId}/roles/${roleId}`) }

// ── Admin Audit Logs ──────────────────────────────────────────────────

export function adminListAuditLogs(params = {}) {
  return apiClient.get('/api/v1/admin/audit-logs', { params })
}

export function adminGetAuditLog(id) {
  return apiClient.get(`/api/v1/admin/audit-logs/${id}`)
}

// Compatibility placeholders for CMS sections whose backend modules are not
// available yet. Keeping named exports prevents production bundlers from
// replacing optional namespace lookups with undefined imports. Callers already
// handle rejected requests and preserve the existing empty/error states.
const unavailableCMSModule = (name) => Promise.reject(new Error(`${name} CMS module is not available`))
export function adminListServices() { return unavailableCMSModule('Services') }
export function adminCreateService() { return unavailableCMSModule('Services') }
export function adminUpdateService() { return unavailableCMSModule('Services') }
export function adminDeleteService() { return unavailableCMSModule('Services') }
export function adminListTickets() { return unavailableCMSModule('Support tickets') }
export function adminCreateTicket() { return unavailableCMSModule('Support tickets') }
export function adminUpdateTicket() { return unavailableCMSModule('Support tickets') }
export function adminDeleteTicket() { return unavailableCMSModule('Support tickets') }
export function adminListArticles() { return unavailableCMSModule('Knowledgebase') }
export function adminCreateArticle() { return unavailableCMSModule('Knowledgebase') }
export function adminUpdateArticle() { return unavailableCMSModule('Knowledgebase') }
export function adminDeleteArticle() { return unavailableCMSModule('Knowledgebase') }
