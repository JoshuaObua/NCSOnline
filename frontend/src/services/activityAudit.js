import apiClient from '@/api/client.js'

const ACTIVITY_ENDPOINT = '/api/v1/account/activity-events'
const PREVIEW_TOKENS = new Set(['local-portal-preview-token', 'local-cms-preview-token'])
const MAX_LABEL_LENGTH = 120

let lastNavigationKey = ''
let lastNavigationAt = 0
let lastClickKey = ''
let lastClickAt = 0
let clickAuditorInstalled = false

export function recordActivityEvent(payload = {}) {
  if (typeof window === 'undefined') return Promise.resolve()
  const token = window.localStorage?.getItem('ncsms_access_token')
  if (!token || PREVIEW_TOKENS.has(token)) return Promise.resolve()
  return apiClient.post(ACTIVITY_ENDPOINT, sanitizeActivityPayload(payload)).catch(() => {})
}

export function recordNavigation(to, from) {
  if (!to || !from || to.fullPath === from.fullPath) return
  const now = Date.now()
  const key = `${from.fullPath}->${to.fullPath}`
  if (key === lastNavigationKey && now - lastNavigationAt < 1500) return
  lastNavigationKey = key
  lastNavigationAt = now
  const section = cleanLabel(to.query?.section || to.params?.resource || '')
  recordActivityEvent({
    type: 'navigation',
    action: 'navigate',
    from_path: from.fullPath || from.path || '',
    to_path: to.fullPath || to.path || '',
    from_page: routePageName(from),
    page_name: routePageName(to),
    route_name: String(to.name || ''),
    section,
    resource: 'Portal Navigation',
    metadata: {
      query_section: section,
      params_resource: cleanLabel(to.params?.resource || ''),
      browser_title: typeof document !== 'undefined' ? document.title : '',
    },
  })
}

export function installActivityAuditor(router) {
  router.afterEach((to, from) => recordNavigation(to, from))
  installClickAuditor()
}

function installClickAuditor() {
  if (clickAuditorInstalled || typeof document === 'undefined') return
  clickAuditorInstalled = true
  document.addEventListener('click', event => {
    const target = event.target?.closest?.('a,button,[role="button"],[data-audit-label]')
    if (!target) return
    if (target.closest?.('input,textarea,select,[contenteditable="true"]')) return
    const label = cleanLabel(target.getAttribute('data-audit-label') || target.getAttribute('aria-label') || target.textContent || target.title || '')
    if (!label) return
    const href = target.getAttribute('href') || target.dataset?.href || ''
    const path = typeof window !== 'undefined' ? `${window.location.pathname}${window.location.search}` : ''
    const key = `${path}:${label}:${href}`
    const now = Date.now()
    if (key === lastClickKey && now - lastClickAt < 1000) return
    lastClickKey = key
    lastClickAt = now
    recordActivityEvent({
      type: 'click',
      action: 'click',
      path,
      page_name: cleanLabel(document.title.replace(/ - NCS Uganda$/, '')) || 'Portal',
      label,
      resource: 'Frontend Interaction',
      metadata: {
        target_tag: target.tagName?.toLowerCase?.() || '',
        target_path: href ? cleanPath(href) : '',
      },
    })
  }, true)
}

function sanitizeActivityPayload(payload) {
  return {
    type: cleanLabel(payload.type || payload.action || 'interaction', 40),
    action: cleanLabel(payload.action || payload.type || 'interaction', 40),
    from_path: cleanPath(payload.from_path || ''),
    to_path: cleanPath(payload.to_path || payload.path || ''),
    path: cleanPath(payload.path || payload.to_path || ''),
    page_name: cleanLabel(payload.page_name || '', MAX_LABEL_LENGTH),
    from_page: cleanLabel(payload.from_page || '', MAX_LABEL_LENGTH),
    label: cleanLabel(payload.label || '', MAX_LABEL_LENGTH),
    route_name: cleanLabel(payload.route_name || '', 80),
    section: cleanLabel(payload.section || '', 80),
    resource: cleanLabel(payload.resource || 'Frontend Activity', 80),
    occurred_at: new Date().toISOString(),
    metadata: sanitizeMetadata(payload.metadata || {}),
  }
}

function sanitizeMetadata(metadata) {
  const safe = {}
  for (const [key, value] of Object.entries(metadata)) {
    const normalizedKey = cleanLabel(key, 50)
    if (!normalizedKey || /password|token|secret|credential/i.test(normalizedKey)) continue
    if (typeof value === 'string') safe[normalizedKey] = cleanLabel(value, 160)
    else if (typeof value === 'number' || typeof value === 'boolean') safe[normalizedKey] = value
  }
  return safe
}

function routePageName(route) {
  const section = cleanLabel(route?.query?.section || route?.params?.resource || '')
  if (section) return titleize(section)
  const name = String(route?.name || '')
  if (name) return titleize(name.replace(/([a-z])([A-Z])/g, '$1 $2'))
  return titleize(String(route?.path || 'Portal').split('/').filter(Boolean).pop() || 'Portal')
}

function cleanLabel(value, maxLength = MAX_LABEL_LENGTH) {
  const normalized = String(value || '').replace(/\s+/g, ' ').trim()
  return normalized.length > maxLength ? normalized.slice(0, maxLength) : normalized
}

function cleanPath(value) {
  const raw = String(value || '').trim()
  if (!raw) return ''
  let path = raw
  try {
    if (/^https?:\/\//i.test(raw)) path = new URL(raw).pathname + new URL(raw).search
  } catch {}
  path = path.split('#')[0]
  return path.length > 240 ? path.slice(0, 240) : path
}

function titleize(value) {
  return cleanLabel(value)
    .replace(/[-_]/g, ' ')
    .replace(/\b\w/g, character => character.toUpperCase())
}
