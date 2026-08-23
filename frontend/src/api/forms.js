// Dynamic application-form API client.
import apiClient from './client'

// ── Departments ────────────────────────────────────────────────
export const listDepartments = () => apiClient.get('/api/v1/departments').then(r => r.data?.data ?? r.data)

// ── Admin: form templates ──────────────────────────────────────
export const adminListForms = (status = '') => {
  const q = status ? `?status=${encodeURIComponent(status)}` : ''
  return apiClient.get(`/api/v1/admin/forms${q}`).then(r => r.data?.data ?? r.data)
}
export const adminGetForm    = (id) => apiClient.get(`/api/v1/admin/forms/${id}`).then(r => r.data?.data ?? r.data)
export const adminCreateForm = (payload) => apiClient.post('/api/v1/admin/forms', payload).then(r => r.data?.data ?? r.data)
export const adminUpdateForm = (id, payload) => apiClient.put(`/api/v1/admin/forms/${id}`, payload).then(r => r.data?.data ?? r.data)
export const adminDeleteForm = (id) => apiClient.delete(`/api/v1/admin/forms/${id}`).then(r => r.data)

// ── Admin: submissions ─────────────────────────────────────────
export const adminListSubmissions = (params = {}) => {
  const qs = new URLSearchParams(params).toString()
  return apiClient.get(`/api/v1/admin/forms/submissions${qs ? '?' + qs : ''}`).then(r => r.data)
}
export const adminGetSubmission    = (id) => apiClient.get(`/api/v1/admin/forms/submissions/${id}`).then(r => r.data?.data ?? r.data)
export const adminReviewSubmission = (id, status, notes) =>
  apiClient.post(`/api/v1/admin/forms/submissions/${id}/review`, { status, notes }).then(r => r.data)
export const adminVerifySubmissionPayment = (id) =>
  apiClient.post(`/api/v1/admin/forms/submissions/${id}/verify-payment`).then(r => r.data)
export const adminUpdateSubmissionPaymentStatus = (id, status) =>
  apiClient.post(`/api/v1/admin/forms/submissions/${id}/payment-status`, { status }).then(r => r.data)

// ── Public portal ──────────────────────────────────────────────
export const portalListOpenForms = () =>
  apiClient.get('/api/v1/portal/forms/open').then(r => r.data?.data ?? r.data)
export const portalGetForm = (slug) =>
  apiClient.get(`/api/v1/portal/forms/${slug}`).then(r => r.data?.data ?? r.data)
export const portalSaveDraft = (templateId, answers) =>
  apiClient.post(`/api/v1/portal/forms/${templateId}/draft`, { answers }).then(r => r.data?.data ?? r.data)
export const portalGetSubmission = (id) =>
  apiClient.get(`/api/v1/portal/submissions/${id}`).then(r => r.data?.data ?? r.data)
export const portalListSubmissions = (params = {}) =>
  apiClient.get('/api/v1/portal/submissions', { params }).then(r => r.data)
export const portalSubmit = (id) =>
  apiClient.post(`/api/v1/portal/submissions/${id}/submit`).then(r => r.data?.data ?? r.data)
export const portalUploadPaymentProof = (id, payload = {}) =>
  apiClient.post(`/api/v1/portal/submissions/${id}/payment-proof`, payload).then(r => r.data)
export const portalInitiateMoMoPayment = (id, phoneNumber) =>
  apiClient.post(`/api/v1/portal/submissions/${id}/pay/momo`, { phone_number: phoneNumber }).then(r => r.data?.data ?? r.data)
export const portalGetPaymentStatus = (id) =>
  apiClient.get(`/api/v1/portal/submissions/${id}/pay/status`).then(r => r.data?.data ?? r.data)

// ── Transactions ───────────────────────────────────────────────
export const listMyTransactions = (params = {}) =>
  apiClient.get('/api/v1/transactions', { params }).then(r => r.data)
export const adminListTransactions = (params = {}) => {
  const qs = new URLSearchParams(params).toString()
  return apiClient.get(`/api/v1/admin/transactions${qs ? '?' + qs : ''}`).then(r => r.data)
}
export const adminGetTransactionKPIs = () =>
  apiClient.get('/api/v1/admin/transactions/kpis').then(r => r.data?.data ?? r.data)

export const FIELD_TYPES = [
  { value: 'short_text', label: 'Short Text',     icon: 'text' },
  { value: 'long_text',  label: 'Long Text',      icon: 'text' },
  { value: 'phone',      label: 'Phone Number',   icon: 'phone' },
  { value: 'email',      label: 'Email Address',  icon: 'email' },
  { value: 'number',     label: 'Number',         icon: 'number' },
  { value: 'date',       label: 'Date',           icon: 'date' },
  { value: 'file',       label: 'File Upload',    icon: 'file' },
  { value: 'image',      label: 'Image Upload',   icon: 'image' },
  { value: 'dropdown',   label: 'Dropdown',       icon: 'select' },
  { value: 'radio',      label: 'Radio Buttons',  icon: 'radio' },
  { value: 'checkbox',   label: 'Checkboxes',     icon: 'check' },
]

export const FIELD_TYPES_WITH_OPTIONS = ['dropdown', 'radio', 'checkbox']
