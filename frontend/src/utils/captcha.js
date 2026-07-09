import { getSettings } from '@/api/cms.js'

let cachedSettings = null
let recaptchaLoader = null
let turnstileLoader = null

function readValue(response) {
  return response?.data?.data?.value || response?.data?.value || {}
}

async function loadSettings() {
  if (cachedSettings) return cachedSettings
  try {
    const response = await getSettings('captcha')
    cachedSettings = readValue(response) || {}
  } catch {
    cachedSettings = { captcha_provider: 'none' }
  }
  return cachedSettings
}

function loadScript(src, attr = 'src') {
  const existing = document.querySelector(`script[${attr}="${src}"]`)
  if (existing) return Promise.resolve()
  return new Promise((resolve, reject) => {
    const script = document.createElement('script')
    script.src = src
    script.async = true
    script.defer = true
    script.onload = resolve
    script.onerror = reject
    document.head.appendChild(script)
  })
}

async function executeRecaptcha(settings, action) {
  const siteKey = settings.recaptcha_site_key
  if (!siteKey) throw new Error('reCAPTCHA is not configured.')
  const src = `https://www.google.com/recaptcha/api.js?render=${encodeURIComponent(siteKey)}`
  recaptchaLoader ||= loadScript(src)
  await recaptchaLoader
  return new Promise((resolve, reject) => {
    window.grecaptcha?.ready(() => {
      window.grecaptcha.execute(siteKey, { action }).then(resolve).catch(reject)
    })
  })
}

async function executeTurnstile(settings) {
  const siteKey = settings.cloudflare_site_key
  if (!siteKey) throw new Error('Cloudflare Turnstile is not configured.')
  turnstileLoader ||= loadScript('https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit')
  await turnstileLoader
  return new Promise((resolve, reject) => {
    const mount = document.createElement('div')
    mount.style.position = 'fixed'
    mount.style.left = '-9999px'
    mount.style.width = '1px'
    mount.style.height = '1px'
    document.body.appendChild(mount)
    const cleanup = () => {
      try { window.turnstile?.remove(widgetId) } catch {}
      mount.remove()
    }
    const widgetId = window.turnstile.render(mount, {
      sitekey: siteKey,
      size: 'invisible',
      callback: token => {
        cleanup()
        resolve(token)
      },
      'error-callback': () => {
        cleanup()
        reject(new Error('Bot verification failed. Please try again.'))
      },
      'expired-callback': () => {
        cleanup()
        reject(new Error('Bot verification expired. Please try again.'))
      },
    })
    window.turnstile.execute(widgetId)
  })
}

export async function captchaPayload(action) {
  const settings = await loadSettings()
  const provider = settings.captcha_provider || 'none'
  if (provider === 'none') return {}
  const token = provider === 'google_recaptcha'
    ? await executeRecaptcha(settings, action)
    : await executeTurnstile(settings)
  return {
    captcha_provider: provider,
    captcha_token: token,
    captcha_action: action,
  }
}

export function clearCaptchaSettingsCache() {
  cachedSettings = null
}
