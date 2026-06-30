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
            <form v-if="isAuthenticated" class="comment-form" @submit.prevent="submitComment(null)">
              <textarea v-model="commentBody" maxlength="1500" placeholder="Share a thoughtful comment"></textarea>
              <button type="submit" :disabled="submitting">Submit for moderation</button>
            </form>
            <div v-else class="login-cta">
              <p>Sign in to join the conversation on this article.</p>
              <router-link to="/login">Log in to comment</router-link>
            </div>
            <div class="comments-list">
              <article v-for="comment in comments" :key="comment.id" class="comment-card">
                <strong>{{ comment.user_name || 'Reader' }}</strong>
                <p>{{ comment.body }}</p>
                <button v-if="isAuthenticated" type="button" @click="replyTo = replyTo === comment.id ? '' : comment.id">Reply</button>
                <form v-if="replyTo === comment.id" class="comment-form compact" @submit.prevent="submitComment(comment.id)">
                  <textarea v-model="replyBody" maxlength="1500" placeholder="Write a reply"></textarea>
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
import { computed, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getPost, listPostComments, submitPostComment } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'
import { sanitizeRichHtml } from '@/utils/sanitize.js'

const route = useRoute()
const post = ref(null)
const comments = ref([])
const commentBody = ref('')
const replyBody = ref('')
const replyTo = ref('')
const submitting = ref(false)
const loading = ref(true)
const isAuthenticated = computed(() => !!localStorage.getItem('ncsms_access_token'))
const safeContent = computed(() => sanitizeRichHtml(post.value?.content || ''))
const flatCommentCount = computed(() => comments.value.reduce((total, c) => total + 1 + (c.replies?.length || 0), 0))

function formatDate(d) {
  return new Date(d).toLocaleDateString('en-UG', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(async () => {
  try {
    const res = await getPost(route.params.slug)
    post.value = res.data.data
    await loadComments()
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
  const body = parentId ? replyBody.value : commentBody.value
  if (!body.trim()) return
  submitting.value = true
  try {
    await submitPostComment(route.params.slug, { body, parent_id: parentId })
    if (parentId) {
      replyBody.value = ''
      replyTo.value = ''
    } else {
      commentBody.value = ''
    }
    alert('Thanks. Your comment is waiting for moderation.')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.comment-form{display:grid;gap:.75rem;margin-bottom:1.5rem}.comment-form textarea{min-height:7rem;border:1px solid #d1d5db;border-radius:.5rem;padding:.85rem;color:#111827}.comment-form button,.login-cta a{justify-self:start;background:#112b4e;color:white;border-radius:.45rem;padding:.65rem 1rem;font-weight:700}.comment-form.compact textarea{min-height:5rem}.login-cta{display:flex;align-items:center;justify-content:space-between;gap:1rem;background:#f8fafc;border:1px solid #e5e7eb;border-radius:.5rem;padding:1rem;margin-bottom:1.5rem}.comments-list{display:grid;gap:1rem}.comment-card{border:1px solid #e5e7eb;border-radius:.5rem;padding:1rem;background:white}.comment-card strong{color:#112b4e}.comment-card p{margin:.5rem 0;color:#4b5563}.comment-card>button{font-size:.85rem;font-weight:700;color:#f48c06}.reply-list{display:grid;gap:.75rem;margin-top:.9rem;padding-left:1rem;border-left:3px solid #facc15}.reply{background:#f8fafc}
</style>
