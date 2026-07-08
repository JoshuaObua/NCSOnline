<template>
  <div class="min-h-screen bg-gray-50">
    <section class="bg-cream px-4 py-16 text-center">
      <span class="section-tag">Leadership</span>
      <h1 class="mt-3 text-4xl font-bold text-darken md:text-5xl">Governing Council</h1>
      <p class="mt-4 text-gray-500">The council providing strategic oversight for sports development in Uganda.</p>
    </section>

    <section class="mx-auto max-w-7xl px-4 py-16">
      <div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
        <div class="order-2 lg:order-1 lg:col-span-2">
          <div class="rounded-lg bg-white p-6 shadow-sm md:p-8">
            <div class="mb-7 flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
              <div>
                <h2 class="flex items-center gap-3 text-2xl font-bold text-[#1a365d]">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-7 w-7 flex-shrink-0 text-[#f5a623]" aria-hidden="true">
                    <path d="M20 13c0 5-3.5 7.5-8 9-4.5-1.5-8-4-8-9V5l8-3 8 3v8Z"></path>
                    <path d="m9 12 2 2 4-4"></path>
                  </svg>
                  Council Members
                </h2>
                <p class="mt-2 max-w-2xl text-gray-600">Meet the appointed members responsible for the Council's strategic direction and governance.</p>
              </div>
              <span v-if="!loading && council.length" class="inline-flex w-fit items-center rounded-md bg-[#1a365d]/5 px-3 py-2 text-sm font-semibold text-[#1a365d]">
                {{ memberCountLabel }}
              </span>
            </div>

            <div v-if="loading" class="grid gap-5 sm:grid-cols-2" aria-label="Loading council members">
              <div v-for="index in 4" :key="index" class="overflow-hidden rounded-lg border border-gray-100">
                <div class="aspect-[4/3] animate-pulse bg-gray-100"></div>
                <div class="space-y-3 p-5">
                  <div class="h-5 w-2/3 animate-pulse rounded bg-gray-100"></div>
                  <div class="h-4 w-1/2 animate-pulse rounded bg-gray-100"></div>
                </div>
              </div>
            </div>

            <div v-else-if="council.length" class="grid gap-5 sm:grid-cols-2">
              <article v-for="(member, index) in council" :key="member.id" class="overflow-hidden rounded-lg border border-gray-100 bg-white transition-shadow duration-300 hover:shadow-md">
                <div class="relative aspect-[4/3] overflow-hidden bg-[#1a365d]/5">
                  <img
                    v-if="member.image_url"
                    :src="mediaUrl(member.image_url)"
                    :alt="member.full_name"
                    class="h-full w-full object-cover"
                  />
                  <div v-else class="flex h-full w-full items-center justify-center" aria-hidden="true">
                    <span class="flex h-24 w-24 items-center justify-center rounded-full bg-[#1a365d] text-3xl font-bold text-white">{{ initials(member.full_name) }}</span>
                  </div>
                  <span class="absolute left-4 top-4 rounded-md bg-white px-2.5 py-1 text-xs font-bold text-[#1a365d] shadow-sm">
                    {{ memberNumber(index) }}
                  </span>
                </div>

                <div class="p-5">
                  <h3 class="text-lg font-bold text-[#1a365d]">{{ member.full_name }}</h3>
                  <p class="mt-1 text-sm font-semibold text-[#c47700]">{{ member.designation }}</p>

                  <div v-if="member.bio" class="mt-4 border-t border-gray-100 pt-4">
                    <div
                      v-show="expandedIds.has(member.id)"
                      :id="`council-bio-${member.id}`"
                      class="prose-bio mb-3 text-sm leading-6 text-gray-600"
                      v-html="sanitizeRichHtml(member.bio)"
                    ></div>
                    <button
                      type="button"
                      class="inline-flex items-center gap-2 text-sm font-semibold text-[#1a365d] transition-colors hover:text-[#c47700]"
                      :aria-expanded="expandedIds.has(member.id)"
                      :aria-controls="`council-bio-${member.id}`"
                      @click="toggleBio(member.id)"
                    >
                      {{ expandedIds.has(member.id) ? 'Hide profile' : 'View profile' }}
                      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :class="['h-4 w-4 transition-transform', expandedIds.has(member.id) ? 'rotate-180' : '']" aria-hidden="true">
                        <path d="m6 9 6 6 6-6"></path>
                      </svg>
                    </button>
                  </div>
                </div>
              </article>
            </div>

            <div v-else class="rounded-lg border border-dashed border-gray-200 bg-gray-50 px-6 py-14 text-center">
              <span class="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-[#1a365d]/10 text-[#1a365d]">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-7 w-7" aria-hidden="true">
                  <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
                  <circle cx="9" cy="7" r="4"></circle>
                  <path d="M19 8v6"></path>
                  <path d="M22 11h-6"></path>
                </svg>
              </span>
              <h3 class="mt-4 text-lg font-bold text-[#1a365d]">Council profiles are being prepared</h3>
              <p class="mx-auto mt-2 max-w-md text-sm leading-6 text-gray-500">Published Governing Council profiles will appear here as they become available.</p>
            </div>
          </div>
        </div>

        <aside class="order-1 space-y-6 lg:order-2">
          <div class="rounded-lg bg-white p-6 shadow-sm">
            <h2 class="mb-4 font-bold text-[#1a365d]">About NCS</h2>
            <nav class="space-y-2" aria-label="About NCS">
              <router-link to="/pages/the-mandate" class="flex items-center gap-2 rounded-md p-3 text-gray-600 transition-colors hover:bg-gray-50">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4" aria-hidden="true"><path d="m16 16 3-8 3 8c-.87.65-1.92 1-3 1s-2.13-.35-3-1Z"></path><path d="m2 16 3-8 3 8c-.87.65-1.92 1-3 1s-2.13-.35-3-1Z"></path><path d="M7 21h10"></path><path d="M12 3v18"></path><path d="M3 7h2c2 0 5-1 7-2 2 1 5 2 7 2h2"></path></svg>
                Mandate
              </router-link>
              <router-link to="/governing-council" class="flex items-center gap-2 rounded-md bg-[#f5a623]/10 p-3 font-medium text-[#c47700]" aria-current="page">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4" aria-hidden="true"><path d="M20 13c0 5-3.5 7.5-8 9-4.5-1.5-8-4-8-9V5l8-3 8 3v8Z"></path><path d="m9 12 2 2 4-4"></path></svg>
                The Council
              </router-link>
              <router-link to="/team" class="flex items-center gap-2 rounded-md p-3 text-gray-600 transition-colors hover:bg-gray-50">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4" aria-hidden="true"><path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"></path><path d="M14 2v4a2 2 0 0 0 2 2h4"></path><path d="M10 9H8"></path><path d="M16 13H8"></path><path d="M16 17H8"></path></svg>
                The Staff
              </router-link>
            </nav>
          </div>

          <div class="rounded-lg bg-[#1a365d] p-6 text-white">
            <h2 class="font-bold">Council's Role</h2>
            <p class="mt-3 text-sm leading-6 text-white/75">The Governing Council provides strategic oversight and supports accountable delivery of the National Council of Sports mandate.</p>
            <div class="mt-5 flex items-center gap-3 border-t border-white/10 pt-5 text-sm text-white/80">
              <span class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-md bg-[#f5a623] text-[#1a365d]">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4" aria-hidden="true"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10"></path></svg>
              </span>
              <span>Leadership, governance and accountability</span>
            </div>
          </div>

          <div class="rounded-lg bg-white p-6 shadow-sm">
            <h2 class="font-bold text-[#1a365d]">Contact the Secretariat</h2>
            <p class="mt-3 text-sm leading-6 text-gray-500">For official correspondence concerning the Council, contact the NCS Secretariat.</p>
            <div class="mt-4 space-y-3 text-sm">
              <a :href="`mailto:${contactEmail}`" class="flex items-center gap-2 text-gray-600 transition-colors hover:text-[#c47700]">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 flex-shrink-0 text-[#f5a623]" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"></path><rect x="2" y="4" width="20" height="16" rx="2"></rect></svg>
                <span class="min-w-0 break-all">{{ contactEmail }}</span>
              </a>
              <a :href="`tel:${contactPhone}`" class="flex items-center gap-2 text-gray-600 transition-colors hover:text-[#c47700]">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 flex-shrink-0 text-[#f5a623]" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 1 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                {{ contactPhone }}
              </a>
            </div>
          </div>
        </aside>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { getSettings, listCouncil } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'
import { sanitizeRichHtml } from '@/utils/sanitize.js'

const council = ref([])
const loading = ref(true)
const expandedIds = ref(new Set())
const contactEmail = ref('info@ncs.go.ug')
const contactPhone = ref('+256 414254477')

const memberCountLabel = computed(() => `${council.value.length} Council ${council.value.length === 1 ? 'Member' : 'Members'}`)

function toggleBio(id) {
  const next = new Set(expandedIds.value)
  next.has(id) ? next.delete(id) : next.add(id)
  expandedIds.value = next
}

function initials(name) {
  const parts = String(name || '')
    .replace(/\([^)]*\)/g, '')
    .split(/\s+/)
    .filter(part => part && !/^(dr|mr|mrs|ms|eng|hon)\.?$/i.test(part))
  return parts.slice(0, 2).map(part => part.charAt(0).toUpperCase()).join('') || 'NCS'
}

function memberNumber(index) {
  return `Member ${String(index + 1).padStart(2, '0')}`
}

onMounted(async () => {
  try {
    const [councilRes, contactRes] = await Promise.allSettled([
      listCouncil(),
      getSettings('contact'),
    ])
    council.value = councilRes.status === 'fulfilled' ? (councilRes.value.data?.data || []) : []
    if (contactRes.status === 'fulfilled') {
      const value = contactRes.value.data?.data?.value || {}
      if (value.email) contactEmail.value = value.email
      if (value.phone) contactPhone.value = value.phone
    }
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.prose-bio :deep(p) {
  margin: 0 0 0.5rem;
}

.prose-bio :deep(p:last-child) {
  margin-bottom: 0;
}
</style>
