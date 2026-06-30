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
