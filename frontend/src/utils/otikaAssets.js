const CSS_FILES = ['app.min.css', 'style.css', 'components.css', 'custom.css']
const FONT_FILES = ['fa-solid-900.woff2', 'fa-regular-400.woff2']

export function ensureOtikaStyles() {
  if (typeof document === 'undefined') return

  for (const font of FONT_FILES) {
    const id = `otika-font-${font}`
    if (document.getElementById(id)) continue
    const link = document.createElement('link')
    link.id = id
    link.rel = 'preload'
    link.as = 'font'
    link.type = 'font/woff2'
    link.crossOrigin = 'anonymous'
    link.href = `/otika-assets/fonts/webfonts/${font}`
    document.head.appendChild(link)
  }

  for (const file of CSS_FILES) {
    const id = file === 'app.min.css' ? 'otika-admin-css' : `otika-admin-${file}`
    if (document.getElementById(id)) continue
    const link = document.createElement('link')
    link.id = id
    link.rel = 'stylesheet'
    link.href = `/otika-assets/css/${file}`
    document.head.appendChild(link)
  }
}
