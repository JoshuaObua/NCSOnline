import apiClient from './client.js'

// ── Public CMS (no auth required) ────────────────────────────────

export function listPosts(params = {}) {
  return apiClient.get('/api/v1/cms/posts', { params })
}

export function getPost(slug) {
  return apiClient.get(`/api/v1/cms/posts/${slug}`)
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

export function adminListSlides() {
  return apiClient.get('/api/v1/admin/cms/slides', { params: { active: 'false' } })
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
  return apiClient.post('/api/v1/admin/media/upload', fd, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// ── Fun Facts ─────────────────────────────────────────────────────

export function listFunFacts() { return apiClient.get('/api/v1/cms/fun-facts') }
export function adminListFunFacts() { return apiClient.get('/api/v1/admin/cms/fun-facts') }
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
