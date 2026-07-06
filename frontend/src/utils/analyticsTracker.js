import { API_BASE_URL } from '@/api/client.js'

const SESSION_KEY = 'ncs_analytics_session_id'
let enteredAt = Date.now()
let interactions = 0
let lastPath = ''
let clickHandler
let exitHandler

function sessionID() {
  let value = sessionStorage.getItem(SESSION_KEY)
  if (!value) {
    value = crypto?.randomUUID?.() || `${Date.now()}-${Math.random().toString(16).slice(2)}`
    sessionStorage.setItem(SESSION_KEY, value)
  }
  return value
}

function payload(eventName, path = window.location.pathname) {
  return {
    event_name: eventName,
    session_id: sessionID(),
    path,
    referrer: document.referrer || '',
    screen_resolution: `${window.screen?.width || 0}x${window.screen?.height || 0}`,
    language: navigator.language || '',
    session_duration_seconds: Math.max(0, Math.round((Date.now() - enteredAt) / 1000)),
    interaction_count: interactions,
  }
}

function postAnalytics(data, keepalive = true) {
  const body = JSON.stringify(data)
  if (keepalive && navigator.sendBeacon) {
    const blob = new Blob([body], { type: 'application/json' })
    if (navigator.sendBeacon(`${API_BASE_URL}/api/v1/analytics/collect`, blob)) return
  }
  fetch(`${API_BASE_URL}/api/v1/analytics/collect`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body,
    keepalive,
    credentials: 'include',
  }).catch(() => {})
}

export function trackPageView(path = window.location.pathname) {
  if (path === lastPath) return
  if (lastPath) postAnalytics(payload('page_exit', lastPath))
  enteredAt = Date.now()
  interactions = 0
  lastPath = path
  postAnalytics(payload('page_view', path), false)
}

export function installAnalyticsTracker() {
  if (clickHandler) return
  clickHandler = (event) => {
    const link = event.target?.closest?.('a,button')
    if (!link) return
    interactions += 1
    const href = link.getAttribute?.('href') || ''
    postAnalytics({ ...payload('click'), path: href || window.location.pathname })
  }
  exitHandler = () => {
    if (lastPath) postAnalytics(payload('page_exit', lastPath))
  }
  document.addEventListener('click', clickHandler, { passive: true })
  window.addEventListener('pagehide', exitHandler)
}

export function uninstallAnalyticsTracker() {
  if (clickHandler) document.removeEventListener('click', clickHandler)
  if (exitHandler) window.removeEventListener('pagehide', exitHandler)
  clickHandler = null
  exitHandler = null
}
