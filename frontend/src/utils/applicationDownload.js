function titleize(value) {
  return String(value || '').toLowerCase().replaceAll('_', ' ').replaceAll('-', ' ').replace(/\b\w/g, char => char.toUpperCase())
}

function parseObject(value) {
  if (!value) return {}
  if (typeof value === 'object') return value
  try { return JSON.parse(value) } catch { return {} }
}

function escapeHtml(value) {
  return String(value ?? '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;')
}

function fileLabel(value) {
  const text = String(value || '')
  const clean = text.split('?')[0].split('#')[0]
  return decodeURIComponent(clean.split('/').pop() || 'Uploaded file')
}

function displayValue(value) {
  if (Array.isArray(value)) return value.map(displayValue).join(', ') || '-'
  if (value && typeof value === 'object') return JSON.stringify(value)
  const text = String(value ?? '').trim()
  if (!text) return '-'
  if (/(\.pdf|\.png|\.jpe?g|\.webp|\/uploads\/|\/media\/)/i.test(text)) return fileLabel(text)
  return text
}

function formatDate(value) {
  if (!value) return '-'
  try {
    return new Intl.DateTimeFormat('en-UG', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
  } catch {
    return String(value)
  }
}

export function downloadApplicationForm(record, options = {}) {
  const source = options.source || record?.source || 'custom'
  const title = record?.template_title || record?.title || titleize(record?.application_type || record?.form_type || 'Application')
  const reference = record?.submission_reference || record?.application_reference || record?.reference || 'Pending'
  const answers = parseObject(record?.answers ?? record?.form_data)
  const rows = Object.entries(answers).map(([key, value]) => `
    <tr><th>${escapeHtml(titleize(key))}</th><td>${escapeHtml(displayValue(value))}</td></tr>
  `).join('')
  const html = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>${escapeHtml(title)} - ${escapeHtml(reference)}</title>
  <style>
    body{font-family:Arial,sans-serif;margin:32px;color:#1f2937}
    h1{margin:0 0 4px;font-size:24px}
    .meta{margin:0 0 22px;color:#6b7280}
    table{width:100%;border-collapse:collapse}
    th,td{padding:10px 12px;border:1px solid #e5e7eb;text-align:left;vertical-align:top}
    th{width:32%;background:#f9fafb}
  </style>
</head>
<body>
  <h1>${escapeHtml(title)}</h1>
  <p class="meta">Reference: ${escapeHtml(reference)} | Source: ${escapeHtml(titleize(source))} | Downloaded: ${escapeHtml(formatDate(new Date()))}</p>
  <table>
    <tbody>
      <tr><th>Applicant</th><td>${escapeHtml(record?.applicant_name || record?.applicant_email || 'Portal user')}</td></tr>
      <tr><th>Status</th><td>${escapeHtml(titleize(record?.status))}</td></tr>
      <tr><th>Payment Status</th><td>${escapeHtml(titleize(record?.payment_status || 'Not required'))}</td></tr>
      ${rows || '<tr><td colspan="2">No response data is available.</td></tr>'}
    </tbody>
  </table>
</body>
</html>`
  const blob = new Blob([html], { type: 'text/html;charset=utf-8' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = `${String(reference || title).replace(/[^a-z0-9-]+/gi, '_')}_application.html`
  document.body.appendChild(link)
  link.click()
  URL.revokeObjectURL(link.href)
  link.remove()
}
