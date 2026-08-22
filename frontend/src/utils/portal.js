// Keep public-site account and application journeys on the portal. The host is
// configured at deploy time, with smart fallback for production and local development.
export function portalUrl(path = '') {
  const base = (
    import.meta.env?.VITE_PORTAL_URL ||
    import.meta.env?.VITE_INTRANET_URL
  )?.replace(/\/$/, '')
  if (base) return `${base}${path}`

  if (typeof window !== 'undefined' && window.location?.hostname) {
    const isLocal = ['localhost', '127.0.0.1', '0.0.0.0'].includes(window.location.hostname)
    if (isLocal) {
      const port = import.meta.env?.VITE_INTRANET_PORT || 9081
      return `${window.location.protocol}//${window.location.hostname}:${port}${path}`
    }
    // Remote / production domain:
    if (window.location.hostname.includes('ncsweb.')) {
      return `${window.location.protocol}//${window.location.hostname.replace('ncsweb.', 'ncsportal.')}${path}`
    }
    return `https://ncsportal.atenimedia.com${path}`
  }

  return path
}

export function portalApiUrl(path = '') {
  const base = (
    import.meta.env?.VITE_PORTAL_API_URL ||
    import.meta.env?.VITE_PORTAL_URL ||
    import.meta.env?.VITE_INTRANET_URL
  )?.replace(/\/$/, '')
  if (base) return `${base}${path}`

  // Production uses the website gateway so visitors never need direct access
  // to the portal's host port. Vite proxies this prefix during development.
  return `/portal-api${path}`
}
