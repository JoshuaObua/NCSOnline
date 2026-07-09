const blockedTags = /<\s*(script|iframe|object|embed|style)[^>]*>.*?<\s*\/\s*\1\s*>/gis
const eventAttributes = /\s+on[a-z]+\s*=\s*(".*?"|'.*?'|[^\s>]+)/gi
const javascriptUrls = /(href|src)\s*=\s*(['"]?)\s*javascript:[^'"\s>]*\2/gi

export function sanitizeRichHtml(value = '') {
  return String(value)
    .replace(blockedTags, '')
    .replace(eventAttributes, '')
    .replace(javascriptUrls, '$1="#"')
}

export function sanitizePlainText(value = '') {
  const div = document.createElement('div')
  div.textContent = String(value)
  return div.innerHTML
}
