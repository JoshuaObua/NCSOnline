import apiClient from '@/api/client.js'

export async function runAction(action, context = {}) {
  if (!action?.type) return
  switch (action.type) {
    case 'link_to':
      if (action.target) window.location.assign(resolveTemplate(action.target, context))
      break
    case 'toggle_visibility':
      context.visibility?.toggle?.(action.target)
      break
    case 'submit_form':
      context.forms?.submit?.(action.form || 'default')
      break
    case 'trigger_api_call':
      return apiClient({
        method: action.method || 'GET',
        url: resolveTemplate(action.url, context),
        data: action.body || {},
      })
    default:
      console.warn(`Unknown page-builder action: ${action.type}`)
  }
}

export function resolveTemplate(value = '', context = {}) {
  return String(value).replace(/\{\{\s*([\w.]+)\s*\}\}/g, (_, key) => {
    const parts = key.split('.')
    let current = context
    for (const part of parts) current = current?.[part]
    return current ?? ''
  })
}
