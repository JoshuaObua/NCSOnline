const GA_ID_PATTERN = /^G-[A-Z0-9]+$/
const SCRIPT_ID = 'ncs-google-analytics-script'
const INLINE_ID = 'ncs-google-analytics-inline'

let activeMeasurementId = ''

function sanitizeMeasurementId(value) {
  const id = String(value || '').trim().toUpperCase()
  return GA_ID_PATTERN.test(id) ? id : ''
}

export function uninstallGoogleAnalytics() {
  document.getElementById(SCRIPT_ID)?.remove()
  document.getElementById(INLINE_ID)?.remove()
  activeMeasurementId = ''
  if (window.gtag) delete window.gtag
  if (window.dataLayer) delete window.dataLayer
}

export function installGoogleAnalytics(settings = {}) {
  const enabled = !!settings.google_analytics_enabled
  const measurementId = sanitizeMeasurementId(settings.google_analytics_id)
  if (!enabled || !measurementId) {
    uninstallGoogleAnalytics()
    return ''
  }
  if (activeMeasurementId === measurementId && document.getElementById(SCRIPT_ID)) {
    return measurementId
  }
  uninstallGoogleAnalytics()
  activeMeasurementId = measurementId

  const script = document.createElement('script')
  script.id = SCRIPT_ID
  script.async = true
  script.src = `https://www.googletagmanager.com/gtag/js?id=${encodeURIComponent(measurementId)}`
  document.head.appendChild(script)

  const inline = document.createElement('script')
  inline.id = INLINE_ID
  inline.text = [
    'window.dataLayer = window.dataLayer || [];',
    'function gtag(){dataLayer.push(arguments);}',
    'gtag("js", new Date());',
    `gtag("config", ${JSON.stringify(measurementId)});`,
  ].join('\n')
  document.head.appendChild(inline)
  return measurementId
}

export function trackGoogleAnalyticsPageView(path = window.location.pathname) {
  if (!activeMeasurementId || typeof window.gtag !== 'function') return
  window.gtag('config', activeMeasurementId, { page_path: path })
}
