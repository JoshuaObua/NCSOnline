<template>
  <div>
    <!-- Loading skeleton -->
    <div v-if="loading" class="bg-cream pt-24 pb-16">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 animate-pulse">
        <div class="h-6 bg-gray-300 rounded w-1/4 mb-4"></div>
        <div class="h-12 bg-gray-300 rounded mb-3 w-3/4"></div>
        <div class="h-4 bg-gray-300 rounded w-1/3"></div>
      </div>
    </div>

    <div v-else-if="post">
      <!-- Hero image banner with dark overlay + title -->
      <div class="relative min-h-[360px] flex items-end" :style="post.cover_image_url ? '' : ''">
        <div
          v-if="post.cover_image_url"
          class="absolute inset-0 bg-cover bg-center"
          :style="`background-image: url('${mediaUrl(post.cover_image_url)}')`"
        ></div>
        <div v-else class="absolute inset-0 bg-gradient-to-br from-[#112b4e] to-[#1e4080]"></div>
        <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent"></div>

        <div class="relative z-10 max-w-4xl mx-auto px-4 sm:px-6 pb-12 pt-32 w-full">
          <!-- Breadcrumb -->
          <nav class="flex items-center gap-2 text-sm text-white/60 mb-5">
            <router-link to="/" class="hover:text-white transition-colors">Home</router-link>
            <span>/</span>
            <router-link to="/news" class="hover:text-white transition-colors">News</router-link>
            <span>/</span>
            <span class="text-white/80 truncate max-w-xs">{{ post.title }}</span>
          </nav>

          <!-- Meta + title -->
          <div class="flex items-center gap-3 mb-4 flex-wrap">
            <span class="text-xs font-semibold bg-yellow-300 text-yellow-900 px-3 py-1 rounded-full capitalize">{{ post.category || 'News' }}</span>
            <span v-if="post.published_at" class="text-sm text-white/70">{{ formatDate(post.published_at) }}</span>
            <span v-if="post.author_name" class="text-sm text-white/70">By {{ post.author_name }}</span>
            <span v-if="post.view_count" class="text-sm text-white/60 flex items-center gap-1">
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.964-7.178z"/><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
              {{ post.view_count }}
            </span>
          </div>
          <h1 class="text-4xl md:text-5xl font-bold text-white leading-tight">{{ post.title }}</h1>
        </div>
      </div>

      <!-- Article content -->
      <div class="bg-white">
        <div class="max-w-4xl mx-auto px-4 sm:px-6 py-14">
          <!-- Excerpt block -->
          <p v-if="post.excerpt" class="text-xl text-gray-500 leading-relaxed mb-10 border-l-4 border-[#F48C06] pl-6 italic">{{ post.excerpt }}</p>

          <!-- Body -->
          <div
            class="prose prose-gray max-w-none prose-headings:text-[#112b4e] prose-a:text-[#F48C06] prose-a:no-underline hover:prose-a:underline prose-img:rounded-xl prose-p:text-gray-600 prose-p:leading-relaxed"
            v-html="safeContent"
          ></div>

          <section class="mt-14 pt-8 border-t border-gray-100">
            <div class="flex items-center justify-between gap-4 mb-6">
              <h2 class="text-2xl font-bold text-[#112b4e]">Comments</h2>
              <span class="text-sm text-gray-400">{{ flatCommentCount }} visible</span>
            </div>
            <div v-if="currentUser" class="commenter-chip">
              <img v-if="currentUser.avatar_url" :src="currentUser.avatar_url" alt="" referrerpolicy="no-referrer" />
              <span v-else>{{ commenterInitials }}</span>
              <strong>{{ commenterName }}</strong>
            </div>
            <form class="comment-form" @submit.prevent="submitComment(null)">
              <input v-model="website" type="text" tabindex="-1" autocomplete="off" class="hp-field" aria-hidden="true" />
              <textarea
                ref="commentTextarea"
                v-model="commentBody"
                maxlength="1500"
                placeholder="Share a thoughtful comment"
                :readonly="!isAuthenticated"
                @focus="ensureCommentAuth"
                @click="ensureCommentAuth"
              ></textarea>
              <button type="submit" :disabled="submitting || authLoading" @click="ensureCommentAuth">
                {{ isAuthenticated ? 'Submit for moderation' : (authLoading ? 'Signing in...' : 'Sign in with Google to comment') }}
              </button>
            </form>
            <div v-show="!isAuthenticated" ref="googleButton" class="google-button-wrap" aria-live="polite"></div>
            <div class="comments-list">
              <article v-for="comment in comments" :key="comment.id" class="comment-card">
                <strong>{{ comment.user_name || 'Reader' }}</strong>
                <p>{{ comment.body }}</p>
                <button type="button" @click="toggleReply(comment.id)">Reply</button>
                <form v-if="replyTo === comment.id" class="comment-form compact" @submit.prevent="submitComment(comment.id)">
                  <textarea v-model="replyBody" maxlength="1500" placeholder="Write a reply" :readonly="!isAuthenticated" @focus="ensureCommentAuth" @click="ensureCommentAuth"></textarea>
                  <button type="submit" :disabled="submitting">Submit reply</button>
                </form>
                <div v-if="comment.replies?.length" class="reply-list">
                  <article v-for="reply in comment.replies" :key="reply.id" class="comment-card reply">
                    <strong>{{ reply.user_name || 'Reader' }}</strong>
                    <p>{{ reply.body }}</p>
                  </article>
                </div>
              </article>
            </div>
          </section>

          <!-- Footer -->
          <div class="mt-14 pt-8 border-t border-gray-100 flex items-center justify-between flex-wrap gap-4">
            <router-link
              to="/news"
              class="text-sm font-semibold text-[#112b4e] hover:text-[#F48C06] flex items-center gap-2 transition-colors"
            >
              ← Back to News
            </router-link>
            <span v-if="post.published_at" class="text-xs text-gray-400">Published {{ formatDate(post.published_at) }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="bg-cream pt-32 pb-20 text-center text-gray-400">
      <p class="text-lg font-medium mb-4">Article not found.</p>
      <router-link to="/news" class="text-sm font-semibold text-accent hover:text-[#d47b05]">← Browse all articles</router-link>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import Swal from 'sweetalert2'
import { getPost, listPostComments, submitPostComment } from '@/api/cms.js'
import { loginWithGoogleCredential } from '@/api/auth.js'
import { mediaUrl } from '@/api/client.js'
import { sanitizePlainText, sanitizeRichHtml } from '@/utils/sanitize.js'

const route = useRoute()
const post = ref(null)
const comments = ref([])
const commentBody = ref('')
const replyBody = ref('')
const replyTo = ref('')
const submitting = ref(false)
const authLoading = ref(false)
const loading = ref(true)
const currentUser = ref(readStoredUser())
const website = ref('')
const googleButton = ref(null)
const commentTextarea = ref(null)
const isAuthenticated = computed(() => !!localStorage.getItem('ncsms_access_token') && !!currentUser.value)
const safeContent = computed(() => sanitizeRichHtml(post.value?.content || ''))
const flatCommentCount = computed(() => comments.value.reduce((total, c) => total + 1 + (c.replies?.length || 0), 0))
const commenterName = computed(() => {
  const u = currentUser.value || {}
  return `${u.first_name || ''} ${u.last_name || ''}`.trim() || u.email || 'Reader'
})
const commenterInitials = computed(() => commenterName.value.split(/\s+/).slice(0, 2).map(p => p[0] || '').join('').toUpperCase() || 'R')

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(async () => {
  try {
    const res = await getPost(route.params.slug)
    post.value = res.data.data
    await loadComments()
    await initGoogleSignIn()
  } catch {
    post.value = null
  } finally {
    loading.value = false
  }
})

async function loadComments() {
  const res = await listPostComments(route.params.slug)
  comments.value = res.data?.data || []
}

async function submitComment(parentId) {
  if (!isAuthenticated.value) {
    await ensureCommentAuth()
    if (!isAuthenticated.value) return
  }
  if (website.value.trim()) return
  const body = parentId ? replyBody.value : commentBody.value
  if (!body.trim()) return
  submitting.value = true
  try {
    const cleanBody = sanitizePlainText(body).slice(0, 1500)
    await submitPostComment(route.params.slug, { body: cleanBody, parent_id: parentId, website: website.value })
    if (parentId) {
      replyBody.value = ''
      replyTo.value = ''
    } else {
      commentBody.value = ''
    }
    await Swal.fire({
      icon: 'success',
      title: 'Comment submitted',
      text: 'Thanks. Your comment is waiting for moderation.',
      confirmButtonColor: '#6777ef',
    })
  } finally {
    submitting.value = false
  }
}

function readStoredUser() {
  try { return JSON.parse(localStorage.getItem('ncsms_user') || 'null') } catch { return null }
}

function storeAuth(result) {
  const data = result?.data?.data || result?.data || {}
  if (data.access_token) localStorage.setItem('ncsms_access_token', data.access_token)
  if (data.user) {
    localStorage.setItem('ncsms_user', JSON.stringify(data.user))
    currentUser.value = data.user
  }
}

async function initGoogleSignIn() {
  if (isAuthenticated.value || !googleClientId()) return
  await loadGoogleScript()
  window.google?.accounts?.id?.initialize({
    client_id: googleClientId(),
    callback: handleGoogleCredential,
    ux_mode: 'popup',
    auto_select: false,
    cancel_on_tap_outside: true,
  })
  renderGoogleButton()
  window.google?.accounts?.id?.prompt()
}

async function ensureCommentAuth() {
  if (isAuthenticated.value || authLoading.value) return
  if (!googleClientId()) {
    await Swal.fire({ icon: 'info', title: 'Sign-in unavailable', text: 'Google sign-in is not configured yet.', confirmButtonColor: '#6777ef' })
    return
  }
  await initGoogleSignIn()
  window.google?.accounts?.id?.prompt()
}

async function handleGoogleCredential(response) {
  if (!response?.credential) return
  authLoading.value = true
  try {
    const result = await loginWithGoogleCredential(response.credential)
    storeAuth(result)
    await nextTick()
    commentTextarea.value?.focus()
  } catch {
    await Swal.fire({ icon: 'error', title: 'Sign-in failed', text: 'Google sign-in could not be verified.', confirmButtonColor: '#6777ef' })
  } finally {
    authLoading.value = false
  }
}

function renderGoogleButton() {
  if (!googleButton.value || !window.google?.accounts?.id) return
  googleButton.value.innerHTML = ''
  window.google.accounts.id.renderButton(googleButton.value, {
    theme: 'outline',
    size: 'large',
    type: 'standard',
    text: 'signin_with',
    shape: 'rectangular',
    width: 260,
  })
}

function loadGoogleScript() {
  if (window.google?.accounts?.id) return Promise.resolve()
  return new Promise((resolve, reject) => {
    const existing = document.querySelector('script[src="https://accounts.google.com/gsi/client"]')
    if (existing) {
      existing.addEventListener('load', resolve, { once: true })
      existing.addEventListener('error', reject, { once: true })
      return
    }
    const script = document.createElement('script')
    script.src = 'https://accounts.google.com/gsi/client'
    script.async = true
    script.defer = true
    script.onload = resolve
    script.onerror = reject
    document.head.appendChild(script)
  })
}

function googleClientId() {
  return import.meta.env.VITE_GOOGLE_CLIENT_ID || ''
}

async function toggleReply(id) {
  if (!isAuthenticated.value) {
    await ensureCommentAuth()
    if (!isAuthenticated.value) return
  }
  replyTo.value = replyTo.value === id ? '' : id
}
</script>

<style scoped>
.comment-form{display:grid;gap:.75rem;margin-bottom:1.5rem}.comment-form textarea{min-height:7rem;border:1px solid #d1d5db;border-radius:.5rem;padding:.85rem;color:#111827;background:white}.comment-form textarea[readonly]{cursor:pointer;background:#f8fafc}.comment-form button{justify-self:start;background:#112b4e;color:white;border-radius:.45rem;padding:.65rem 1rem;font-weight:700}.comment-form button:disabled{opacity:.65}.comment-form.compact textarea{min-height:5rem}.commenter-chip{display:inline-flex;align-items:center;gap:.6rem;background:#f8fafc;border:1px solid #e5e7eb;border-radius:999px;padding:.35rem .8rem .35rem .4rem;margin-bottom:.8rem}.commenter-chip img,.commenter-chip span{width:2rem;height:2rem;border-radius:999px;display:grid;place-items:center;background:#112b4e;color:white;font-size:.75rem;font-weight:800}.google-button-wrap{min-height:44px;margin:-.75rem 0 1.5rem}.hp-field{position:absolute;left:-9999px;width:1px;height:1px;opacity:0}.comments-list{display:grid;gap:1rem}.comment-card{border:1px solid #e5e7eb;border-radius:.5rem;padding:1rem;background:white}.comment-card strong{color:#112b4e}.comment-card p{margin:.5rem 0;color:#4b5563}.comment-card>button{font-size:.85rem;font-weight:700;color:#f48c06}.reply-list{display:grid;gap:.75rem;margin-top:.9rem;padding-left:1rem;border-left:3px solid #facc15}.reply{background:#f8fafc}
</style>
