const spacingKeys = {
  all: 'p',
  x: 'px',
  y: 'py',
  top: 'pt',
  right: 'pr',
  bottom: 'pb',
  left: 'pl',
}

const marginKeys = {
  all: 'm',
  x: 'mx',
  y: 'my',
  top: 'mt',
  right: 'mr',
  bottom: 'mb',
  left: 'ml',
}

const allowedBreakpoints = new Set(['sm', 'md', 'lg', 'xl', '2xl'])

export function settingsToClasses(settings = {}) {
  const classes = []
  pushSpacing(classes, settings.padding, spacingKeys)
  pushSpacing(classes, settings.margin, marginKeys)
  if (settings.background) classes.push(`bg-${settings.background}`)
  if (settings.color) classes.push(`text-${settings.color}`)
  if (settings.weight) classes.push(`font-${settings.weight}`)
  if (settings.leading) classes.push(`leading-${settings.leading}`)
  if (settings.tracking) classes.push(`tracking-${settings.tracking}`)
  if (settings.radius) classes.push(`rounded-${settings.radius}`)
  if (settings.shadow) classes.push(`shadow-${settings.shadow}`)
  if (settings.border) classes.push('border border-gray-200')
  if (settings.gap != null) classes.push(`gap-${settings.gap}`)
  if (settings.justify) classes.push(`justify-${settings.justify}`)
  if (settings.align) classes.push(`items-${settings.align}`)
  if (settings.wrap) classes.push('flex-wrap')
  if (settings.hoverShadow) classes.push(`hover:shadow-${settings.hoverShadow}`)
  if (settings.classes) classes.push(settings.classes)
  pushResponsive(classes, settings.responsive)
  return classes.filter(Boolean).join(' ')
}

export function layoutStyle(settings = {}) {
  const style = {}
  if (settings.display === 'grid') style.display = 'grid'
  if (settings.display === 'flex') style.display = 'flex'
  if (settings.columns?.length) style.gridTemplateColumns = settings.columns.join(' ')
  if (settings.aspect) style.aspectRatio = settings.aspect
  if (settings.backgroundImage) style.backgroundImage = `url("${settings.backgroundImage}")`
  if (settings.backgroundVideo) style.position = 'relative'
  if (settings.width) style.width = settings.width
  if (settings.maxWidth) style.maxWidth = settings.maxWidth
  return style
}

function pushSpacing(classes, value, keys) {
  if (!value || typeof value !== 'object') return
  for (const [key, prefix] of Object.entries(keys)) {
    if (value[key] != null && safeToken(value[key])) classes.push(`${prefix}-${value[key]}`)
  }
}

function pushResponsive(classes, responsive = {}) {
  if (!responsive || typeof responsive !== 'object') return
  for (const [breakpoint, config] of Object.entries(responsive)) {
    if (!allowedBreakpoints.has(breakpoint)) continue
    const nested = settingsToClasses(config)
    nested.split(' ').filter(Boolean).forEach(token => classes.push(`${breakpoint}:${token}`))
  }
}

function safeToken(value) {
  return /^[a-zA-Z0-9_./:%\-[\]]+$/.test(String(value))
}
