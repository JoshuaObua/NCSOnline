export function slugifyFieldKey(value) {
  return String(value || '')
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 80)
}

export function parseFieldConfig(value) {
  if (!value) return {}
  if (typeof value === 'object' && !Array.isArray(value)) return { ...value }
  try {
    const parsed = JSON.parse(value)
    return parsed && typeof parsed === 'object' && !Array.isArray(parsed) ? parsed : {}
  } catch {
    return {}
  }
}

export function normalizeFormField(field = {}, index = 0) {
  const config = parseFieldConfig(field.config)
  const options = Array.isArray(config.options)
    ? config.options.map(option => typeof option === 'object' ? option.label ?? option.value ?? '' : option)
    : []

  return {
    id: field.id || `field-${Date.now()}-${index}`,
    field_key: field.field_key || '',
    field_type: field.field_type || 'short_text',
    label: field.label || '',
    placeholder: field.placeholder || '',
    help_text: field.help_text || '',
    is_required: !!field.is_required,
    options_text: options.filter(Boolean).join('\n'),
    accepted_types: Array.isArray(config.accept)
      ? config.accept.join(',')
      : String(config.accept || ''),
    config,
  }
}

export function buildTemplatePayload(form = {}) {
  const usedKeys = new Set()
  const fields = (form.fields || []).map((field, index) => {
    const baseKey = slugifyFieldKey(field.field_key || field.label) || `field-${index + 1}`
    let fieldKey = baseKey
    let suffix = 2
    while (usedKeys.has(fieldKey)) {
      fieldKey = `${baseKey}-${suffix}`
      suffix += 1
    }
    usedKeys.add(fieldKey)

    const config = { ...parseFieldConfig(field.config) }
    if (['dropdown', 'radio', 'checkbox'].includes(field.field_type)) {
      config.options = String(field.options_text || '')
        .split(/\r?\n/)
        .map(option => option.trim())
        .filter(Boolean)
    } else {
      delete config.options
    }
    if (['file', 'image'].includes(field.field_type) && field.accepted_types) {
      config.accept = String(field.accepted_types)
        .split(',')
        .map(type => type.trim())
        .filter(Boolean)
    } else {
      delete config.accept
    }

    return {
      field_key: fieldKey,
      field_type: field.field_type,
      label: String(field.label || '').trim(),
      placeholder: String(field.placeholder || '').trim(),
      help_text: String(field.help_text || '').trim(),
      is_required: !!field.is_required,
      config,
    }
  })

  return {
    department_id: form.department_id || '',
    title: String(form.title || '').trim(),
    description: String(form.description || '').trim(),
    banner_image_url: form.banner_image_url || '',
    price_ugx: Math.max(0, Number(form.price_ugx) || 0),
    status: form.status || 'DRAFT',
    fields,
  }
}

export function parseSubmissionAnswers(value) {
  if (!value) return {}
  if (typeof value === 'object' && !Array.isArray(value)) return value
  try {
    const parsed = JSON.parse(value)
    return parsed && typeof parsed === 'object' && !Array.isArray(parsed) ? parsed : {}
  } catch {
    return {}
  }
}

export function normalizePagedResponse(value) {
  const body = value?.data && !Array.isArray(value.data) && value.data.success !== undefined
    ? value.data
    : value
  return {
    items: Array.isArray(body?.data) ? body.data : Array.isArray(body) ? body : [],
    meta: body?.meta || {},
  }
}
