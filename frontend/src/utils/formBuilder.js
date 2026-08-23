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

export function parseJsonArray(value) {
  if (!value) return []
  if (Array.isArray(value)) return value
  try {
    const parsed = typeof value === 'string' ? JSON.parse(value) : value
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

export function normalizeFormSection(section = {}, index = 0) {
  return {
    id: String(section.id || section.section_id || `section-${index + 1}`).trim(),
    title: String(section.title || section.name || (index === 0 ? 'Application Details' : `Section ${index + 1}`)).trim(),
    subtitle: String(section.subtitle || '').trim(),
    description: String(section.description || '').trim(),
  }
}

export function normalizeFormSections(value = [], fields = []) {
  const rawSections = parseJsonArray(value?.sections ?? value)
  const seen = new Set()
  const sections = rawSections.map((section, index) => {
    const normalized = normalizeFormSection(section, index)
    let id = normalized.id || `section-${index + 1}`
    let suffix = 2
    while (seen.has(id)) {
      id = `${normalized.id || `section-${index + 1}`}-${suffix}`
      suffix += 1
    }
    seen.add(id)
    return { ...normalized, id }
  }).filter(section => section.id && section.title)

  if (!sections.length) {
    sections.push(normalizeFormSection({}, 0))
  }

  const sectionIds = new Set(sections.map(section => section.id))
  for (const field of fields || []) {
    const config = parseFieldConfig(field.config)
    const sectionId = field.section_id || config.section_id
    if (sectionId && !sectionIds.has(sectionId)) {
      sections.push(normalizeFormSection({ id: sectionId, title: 'Application Details' }, sections.length))
      sectionIds.add(sectionId)
    }
  }

  return sections
}

export function normalizeFormField(field = {}, index = 0, fallbackSectionId = '') {
  const config = parseFieldConfig(field.config)
  const options = Array.isArray(config.options)
    ? config.options.map(option => typeof option === 'object' ? option.label ?? option.value ?? '' : option)
    : []
  const sectionId = field.section_id || config.section_id || fallbackSectionId

  return {
    id: field.id || `field-${Date.now()}-${index}`,
    field_key: field.field_key || '',
    field_type: field.field_type || 'short_text',
    label: field.label || '',
    placeholder: field.placeholder || '',
    help_text: field.help_text || '',
    is_required: !!field.is_required,
    section_id: sectionId,
    options_text: options.filter(Boolean).join('\n'),
    accepted_types: Array.isArray(config.accept)
      ? config.accept.join(',')
      : String(config.accept || ''),
    config,
  }
}

export function buildTemplatePayload(form = {}) {
  const usedKeys = new Set()
  const sections = normalizeFormSections(form.sections || [], form.fields || [])
  const sectionIds = new Set(sections.map(section => section.id))
  const fallbackSectionId = sections[0]?.id || 'section-1'
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
    config.section_id = sectionIds.has(field.section_id) ? field.section_id : fallbackSectionId
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

  const allowedPaymentMethods = Array.isArray(form.allowed_payment_methods) && form.allowed_payment_methods.length
    ? form.allowed_payment_methods
    : ['OVER_THE_COUNTER', 'MOBILE_MONEY']

  return {
    department_id: form.department_id || '',
    slug: String(form.slug || '').trim(),
    title: String(form.title || '').trim(),
    description: String(form.description || '').trim(),
    banner_image_url: form.banner_image_url || '',
    price_ugx: Math.max(0, Number(form.price_ugx) || 0),
    allowed_payment_methods: allowedPaymentMethods,
    status: form.status || 'DRAFT',
    sections,
    fields,
  }
}

export function buildSectionSteps(template = {}) {
  const sections = normalizeFormSections(template.sections || [], template.fields || [])
  const fallbackSectionId = sections[0]?.id || 'section-1'
  const sectionIds = new Set(sections.map(section => section.id))
  const fields = (template.fields || []).map((field, index) => {
    const normalized = normalizeFormField(field, index, fallbackSectionId)
    if (!sectionIds.has(normalized.section_id)) normalized.section_id = fallbackSectionId
    return normalized
  })
  return sections.map(section => ({
    ...section,
    fields: fields.filter(field => field.section_id === section.id),
  }))
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
