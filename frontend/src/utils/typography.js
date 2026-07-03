// Registry of public-site elements the Appearance settings module can style,
// grouped to match the CMS UI. `selector` is the actual CSS target used when
// generating the injected stylesheet — kept here (not duplicated in the
// admin panel) so the mapping between a friendly label and real markup stays
// in one place.
export const TYPOGRAPHY_GROUPS = [
  {
    label: 'Global Layout',
    elements: [
      { key: 'body', label: 'Document Body', selector: 'body' },
      { key: 'code', label: 'Code Blocks', selector: 'code, pre' },
    ],
  },
  {
    label: 'Headings',
    elements: [
      { key: 'heading', label: 'Global Titles (H1–H3)', selector: 'h1, h2, h3' },
      { key: 'subtitle', label: 'Subtitles (H4–H6)', selector: 'h4, h5, h6' },
    ],
  },
  {
    label: 'Public Components',
    elements: [
      { key: 'hero_title', label: 'Hero / Slider Titles', selector: '.carousel-content h1' },
      { key: 'hero_description', label: 'Slider Descriptions', selector: '.carousel-content p' },
      { key: 'paragraph', label: 'Standard Paragraphs', selector: 'p' },
      { key: 'button', label: 'Button Text', selector: 'button, .carousel-btn' },
      {
        key: 'nav_menu',
        label: 'Navigation Menu Links',
        selector: 'nav[aria-label="Primary navigation"] a, nav[aria-label="Primary navigation"] button, '
          + 'nav[aria-label="Mobile navigation"] a, nav[aria-label="Mobile navigation"] button',
      },
    ],
  },
]

export const FONT_WEIGHTS = [
  { value: '300', label: '300 — Light' },
  { value: '400', label: '400 — Regular' },
  { value: '500', label: '500 — Medium' },
  { value: '600', label: '600 — Semibold' },
  { value: '700', label: '700 — Bold' },
]

export const FONT_SIZE_UNITS = ['px', 'rem']

export const WEB_SAFE_FONTS = [
  'Inter', 'Roboto', 'Arial', 'Helvetica', 'Georgia', 'Times New Roman', 'Verdana', 'Tahoma',
  'serif', 'sans-serif', 'monospace',
]

export function emptyTypographyElement() {
  return { family: '', size: '', unit: 'rem', weight: '', line_height: '', letter_spacing: '' }
}

export function normalizeTypography(value) {
  const out = {}
  for (const group of TYPOGRAPHY_GROUPS) {
    for (const el of group.elements) {
      out[el.key] = { ...emptyTypographyElement(), ...(value?.[el.key] || {}) }
    }
  }
  return out
}

const FONT_FORMAT_TO_CSS = { ttf: 'truetype', otf: 'opentype', woff: 'woff', woff2: 'woff2' }

function quoteFontFamily(name) {
  const trimmed = String(name || '').trim()
  if (!trimmed) return ''
  return /\s/.test(trimmed) && !trimmed.startsWith("'") ? `'${trimmed}'` : trimmed
}

export function buildFontFaceCSS(customFonts = []) {
  return customFonts.map(f => `@font-face {
  font-family: ${quoteFontFamily(f.font_name)};
  src: url('${f.file_url}') format('${FONT_FORMAT_TO_CSS[f.font_format] || f.font_format}');
  font-display: swap;
}`).join('\n')
}

// Builds the full injectable stylesheet: @font-face declarations for every
// custom font, followed by one rule per configured element. Elements with no
// configured properties are skipped entirely so the site's existing styling
// (Tailwind defaults) shows through untouched until an admin sets something.
export function buildTypographyCSS(config, customFonts = []) {
  const typography = normalizeTypography(config)
  const rules = TYPOGRAPHY_GROUPS.flatMap(g => g.elements)
    .map(el => {
      const c = typography[el.key]
      const decls = []
      if (c.family) decls.push(`font-family: ${quoteFontFamily(c.family)}, sans-serif`)
      if (c.size) decls.push(`font-size: ${c.size}${c.unit || 'px'}`)
      if (c.weight) decls.push(`font-weight: ${c.weight}`)
      if (c.line_height) decls.push(`line-height: ${c.line_height}`)
      if (c.letter_spacing) decls.push(`letter-spacing: ${c.letter_spacing}`)
      return decls.length ? `${el.selector} { ${decls.join('; ')}; }` : ''
    })
    .filter(Boolean)
    .join('\n')

  return [buildFontFaceCSS(customFonts), rules].filter(Boolean).join('\n\n')
}

// Injects (or updates) a single <style> tag so re-saving settings doesn't
// accumulate duplicate tags across the page's lifetime.
export function injectTypographyCSS(config, customFonts = []) {
  if (typeof document === 'undefined') return
  const css = buildTypographyCSS(config, customFonts)
  let tag = document.getElementById('cms-typography')
  if (!tag) {
    tag = document.createElement('style')
    tag.id = 'cms-typography'
    document.head.appendChild(tag)
  }
  tag.textContent = css
}
