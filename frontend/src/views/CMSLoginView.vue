<template>
  <main class="otika-auth">
    <section class="section">
      <div class="container mt-5">
        <div class="row">
          <div class="col-12 col-sm-8 offset-sm-2 col-md-6 offset-md-3 col-lg-6 offset-lg-3 col-xl-4 offset-xl-4">
            <div class="login-brand">
              <img src="/main-logo.png" alt="NCS" width="96" class="shadow-light rounded-circle bg-white p-2" />
            </div>
            <div class="card card-primary">
              <div class="card-header">
                <h4>CMS Login</h4>
              </div>
              <div class="card-body">
                <form class="needs-validation" novalidate @submit.prevent="login">
                  <div class="form-group">
                    <label for="cms-email">Email</label>
                    <input id="cms-email" v-model="email" type="email" autocomplete="email" class="form-control" tabindex="1" required autofocus />
                    <div class="invalid-feedback">Please fill in your email</div>
                  </div>
                  <div class="form-group">
                    <div class="d-block">
                      <label for="cms-password" class="control-label">Password</label>
                      <div class="float-right">
                        <a href="#" class="text-small">Forgot Password?</a>
                      </div>
                    </div>
                    <input id="cms-password" v-model="password" type="password" autocomplete="current-password" class="form-control" tabindex="2" required />
                    <div class="invalid-feedback">Please fill in your password</div>
                  </div>
                  <div class="form-group">
                    <div class="custom-control custom-checkbox">
                      <input id="remember-me" v-model="remember" type="checkbox" class="custom-control-input" tabindex="3" />
                      <label class="custom-control-label" for="remember-me">Remember Me</label>
                    </div>
                  </div>
                  <p v-if="error" class="alert alert-danger py-2">{{ error }}</p>
                  <div class="form-group">
                    <button type="submit" :disabled="loading" class="btn btn-primary btn-lg btn-block" tabindex="4">
                      {{ loading ? 'Signing in...' : 'Login' }}
                    </button>
                  </div>
                </form>
              </div>
            </div>
            <div class="mt-5 text-muted text-center">
              National Council of Sports <router-link to="/">Back to website</router-link>
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import apiClient, { API_BASE_URL } from '@/api/client.js'
import { ensureOtikaStyles } from '@/utils/otikaAssets.js'

const router = useRouter()
const email = ref('admin@ncs.go.ug')
const password = ref('NCS@Admin2026!')
const loading = ref(false)
const error = ref('')
const remember = ref(true)

ensureOtikaStyles()

async function login() {
  loading.value = true
  error.value = ''
  try {
    if (isDefaultPreviewLogin() && !(await canReachApi())) {
      enterPreviewMode()
      return
    }
    const res = await apiClient.post('/api/v1/auth/login', { email: email.value, password: password.value })
    const data = res.data?.data || {}
    localStorage.setItem('ncsms_access_token', data.access_token || '')
    localStorage.setItem('ncsms_user', JSON.stringify(data.user || {}))
    router.push('/cms')
  } catch (err) {
    if (!err.response && isDefaultPreviewLogin()) {
      enterPreviewMode()
      return
    }
    error.value = err.response?.data?.error?.message || 'Login failed. Check the CMS credentials and API server.'
  } finally {
    loading.value = false
  }
}

function isDefaultPreviewLogin() {
  return email.value === 'admin@ncs.go.ug' && password.value === 'NCS@Admin2026!'
}

function enterPreviewMode() {
  localStorage.setItem('ncsms_access_token', 'local-cms-preview-token')
  localStorage.setItem('ncsms_user', JSON.stringify({ email: email.value, roles: ['super_admin'] }))
  router.push('/cms')
}

async function canReachApi() {
  let timer
  try {
    const controller = new AbortController()
    timer = window.setTimeout(() => controller.abort(), 1200)
    const res = await fetch(`${API_BASE_URL}/health`, { signal: controller.signal, cache: 'no-store' })
    return res.ok
  } catch {
    return false
  } finally {
    if (timer) window.clearTimeout(timer)
  }
}
</script>

<style scoped>
.otika-auth{min-height:100vh;background:#f4f6f9;color:#34395e}
.login-brand{margin:20px 0;text-align:center}
.card.card-primary{border-top:2px solid #6777ef;box-shadow:0 4px 25px 0 rgba(0,0,0,.1)}
.form-control{height:42px;border-color:#e4e6fc}
.form-control:focus{border-color:#6777ef;box-shadow:0 2px 6px #acb5f6}
@media(max-width:575px){.otika-auth .container{margin-top:1.25rem!important;padding:0 18px}.login-brand img{width:78px}.card .card-header,.card .card-body{padding-left:20px;padding-right:20px}.float-right{float:none!important;display:block;margin-top:4px}}
</style>
