// Keep public-site account and application journeys on the portal. The host is
// configured at deploy time, with a same-host fallback for local development.
export function portalUrl(path = '') {
  const base = import.meta.env?.VITE_INTRANET_URL?.replace(/\/$/, '')
  if (base) return `${base}${path}`

  if (typeof window !== 'undefined' && window.location?.hostname) {
    const port = import.meta.env?.VITE_INTRANET_PORT
    const origin = port
      ? `${window.location.protocol}//${window.location.hostname}:${port}`
      : window.location.origin
    return `${origin}${path}`
  }

  return path
}

export function portalApiUrl(path = '') {
  const base = (
    import.meta.env?.VITE_PORTAL_API_URL || import.meta.env?.VITE_INTRANET_URL
  )?.replace(/\/$/, '')
  if (base) return `${base}${path}`

  // Production uses the website gateway so visitors never need direct access
  // to the portal's host port. Vite proxies this prefix during development.
  return `/portal-api${path}`
}
