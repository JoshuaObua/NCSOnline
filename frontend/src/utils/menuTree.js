export const MAX_MENU_DEPTH = 3

export function createMenuItem(overrides = {}) {
  return normalizeMenuItem({
    id: cryptoSafeId(),
    title: 'New menu item',
    url: '/',
    target: '_self',
    hidden: false,
    children: [],
    ...overrides,
  })
}

export function normalizeMenuTree(items = []) {
  return items.map(item => normalizeMenuItem(item)).filter(Boolean)
}

export function normalizeMenuItem(item) {
  if (!item || typeof item !== 'object') return null
  const title = sanitizeText(item.title ?? item.label ?? 'Untitled')
  const url = sanitizeUrl(item.url ?? '/')
  const target = item.target === '_blank' ? '_blank' : '_self'
  return {
    id: sanitizeText(item.id || cryptoSafeId()),
    title,
    url,
    target,
    hidden: item.hidden === true,
    children: normalizeMenuTree(item.children || item.megaItems || []),
  }
}

export function toCmsMenuItems(items = []) {
  return normalizeMenuTree(items).map(item => ({
    id: item.id,
    label: item.title,
    title: item.title,
    url: item.url,
    target: item.target,
    hidden: item.hidden,
    children: toCmsMenuItems(item.children),
  }))
}

// Derives a UI-only "link type" from the item's url — not persisted, so it
// can never drift out of sync with the url it describes.
export function inferLinkType(url, pageOptions = []) {
  const value = String(url ?? '')
  if (value.startsWith('#')) return 'anchor'
  if (pageOptions.some(page => page.url === value)) return 'internal'
  return 'custom'
}

export function cloneTree(items = []) {
  return normalizeMenuTree(JSON.parse(JSON.stringify(items)))
}

export function maxDepth(items = [], depth = 1) {
  return items.reduce((max, item) => Math.max(max, depth, maxDepth(item.children || [], depth + 1)), depth)
}

export function sanitizeText(value) {
  return String(value ?? '')
    .replace(/[<>]/g, '')
    .replace(/[\u0000-\u001f\u007f]/g, '')
    .trim()
    .slice(0, 160)
}

export function sanitizeUrl(value) {
  const raw = String(value ?? '/').trim()
  if (!raw) return '/'
  if (/^(javascript|data|vbscript):/i.test(raw)) return '/'
  if (/^(https?:|mailto:|tel:|\/|#)/i.test(raw)) return raw.slice(0, 2048)
  return `/${raw.replace(/^\/+/, '')}`.slice(0, 2048)
}

export function debounce(fn, wait = 350) {
  let timer
  return (...args) => {
    clearTimeout(timer)
    timer = setTimeout(() => fn(...args), wait)
  }
}

function cryptoSafeId() {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) return crypto.randomUUID()
  return `menu_${Date.now()}_${Math.random().toString(16).slice(2)}`
}
