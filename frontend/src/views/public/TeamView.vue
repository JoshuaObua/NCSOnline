<template>
  <div class="bg-gray-50 min-h-screen">
    <section class="bg-cream py-16 px-4 text-center"><span class="section-tag">Leadership</span><h1 class="text-4xl md:text-5xl font-bold text-darken mt-3">The Staff</h1><p class="mt-4 text-gray-500">The Secretariat delivering NCS's mandate day to day.</p></section>

    <section class="max-w-7xl mx-auto px-4 py-16">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div class="lg:col-span-2 order-2 lg:order-1">
          <div class="bg-white rounded-xl p-6 md:p-8 shadow-sm">
            <div class="mb-6">
              <h2 class="text-2xl font-bold text-[#1a365d] mb-2 flex items-center gap-3">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-7 h-7 text-[#f5a623]" aria-hidden="true"><path d="M6 22V4a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v18Z"></path><path d="M6 12H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h2"></path><path d="M18 9h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-2"></path><path d="M10 6h4"></path><path d="M10 10h4"></path><path d="M10 14h4"></path><path d="M10 18h4"></path></svg>
                Staff Departments
              </h2>
              <p class="text-gray-600">Click on any department to view staff members. The Secretariat operates under the General Secretary.</p>
            </div>

            <div v-if="loading" class="text-center text-gray-400 py-8">Loading departments…</div>
            <div v-else class="grid md:grid-cols-2 gap-4">
              <div v-for="dept in orderedDepartments" :key="dept.id" class="bg-gray-50 rounded-xl overflow-hidden border border-transparent hover:border-[#f5a623]/30 transition-all duration-300">
                <button type="button" class="w-full p-5 flex items-start gap-4 text-left hover:bg-white transition-colors" @click="toggleDept(dept.id)">
                  <div :class="['w-12 h-12 rounded-lg flex items-center justify-center flex-shrink-0', deptMeta(dept.name).color]">
                    <svg v-if="deptMeta(dept.name).icon==='building2'" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true"><path d="M6 22V4a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v18Z"></path><path d="M6 12H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h2"></path><path d="M18 9h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-2"></path><path d="M10 6h4"></path><path d="M10 10h4"></path><path d="M10 14h4"></path><path d="M10 18h4"></path></svg>
                    <svg v-else-if="deptMeta(dept.name).icon==='users'" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path><path d="M16 3.128a4 4 0 0 1 0 7.744"></path><path d="M22 21v-2a4 4 0 0 0-3-3.87"></path><circle cx="9" cy="7" r="4"></circle></svg>
                    <svg v-else-if="deptMeta(dept.name).icon==='calculator'" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true"><rect width="16" height="20" x="4" y="2" rx="2"></rect><line x1="8" x2="16" y1="6" y2="6"></line><line x1="16" x2="16" y1="14" y2="18"></line><path d="M16 10h.01"></path><path d="M12 10h.01"></path><path d="M8 10h.01"></path><path d="M12 14h.01"></path><path d="M8 14h.01"></path><path d="M12 18h.01"></path><path d="M8 18h.01"></path></svg>
                    <svg v-else-if="deptMeta(dept.name).icon==='monitor'" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true"><rect width="20" height="14" x="2" y="3" rx="2"></rect><line x1="8" x2="16" y1="21" y2="21"></line><line x1="12" x2="12" y1="17" y2="21"></line></svg>
                    <svg v-else-if="deptMeta(dept.name).icon==='medal'" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true"><path d="M7.21 15 2.66 7.14a2 2 0 0 1 .13-2.2L4.4 2.8A2 2 0 0 1 6 2h12a2 2 0 0 1 1.6.8l1.6 2.14a2 2 0 0 1 .14 2.2L16.79 15"></path><path d="M11 12 5.12 2.2"></path><path d="m13 12 5.88-9.8"></path><path d="M8 7h8"></path><circle cx="12" cy="17" r="5"></circle><path d="M12 18v-2h-.5"></path></svg>
                    <svg v-else-if="deptMeta(dept.name).icon==='heart-pulse'" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true"><path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"></path><path d="M3.22 12H9.5l.5-1 2 4.5 2-7 1.5 3.5h5.27"></path></svg>
                    <svg v-else-if="deptMeta(dept.name).icon==='megaphone'" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true"><path d="m3 11 18-5v12L3 14v-3z"></path><path d="M11.6 16.8a3 3 0 1 1-5.8-1.6"></path></svg>
                    <svg v-else-if="deptMeta(dept.name).icon==='briefcase'" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true"><path d="M16 20V4a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"></path><rect width="20" height="14" x="2" y="6" rx="2"></rect></svg>
                    <svg v-else-if="deptMeta(dept.name).icon==='file-text'" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true"><path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"></path><path d="M14 2v4a2 2 0 0 0 2 2h4"></path><path d="M10 9H8"></path><path d="M16 13H8"></path><path d="M16 17H8"></path></svg>
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center justify-between">
                      <h3 class="font-bold mb-1 transition-colors text-[#1a365d]">{{ dept.name }}</h3>
                      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :class="['w-5 h-5 text-gray-400 transition-transform duration-300 flex-shrink-0', openIds.has(dept.id) ? 'rotate-180' : '']" aria-hidden="true"><path d="m6 9 6 6 6-6"></path></svg>
                    </div>
                    <p class="text-gray-600 text-sm mb-2 line-clamp-2">{{ dept.description }}</p>
                    <div class="inline-flex items-center rounded-md border px-2.5 py-0.5 font-semibold border-transparent bg-secondary text-secondary-foreground text-xs">{{ dept.staff_count }} Staff Members</div>
                  </div>
                </button>
                <div :class="['overflow-hidden transition-all duration-500', openIds.has(dept.id) ? 'max-h-[1000px] opacity-100' : 'max-h-0 opacity-0']">
                  <div class="border-t border-gray-200 p-5 bg-white">
                    <h4 class="text-base font-semibold text-[#1a365d] mb-4 flex items-center gap-2">
                      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623]" aria-hidden="true"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path><path d="M16 3.128a4 4 0 0 1 0 7.744"></path><path d="M22 21v-2a4 4 0 0 0-3-3.87"></path><circle cx="9" cy="7" r="4"></circle></svg>
                      Department Staff ({{ staffByDept(dept.id).length }})
                    </h4>
                    <p v-if="!staffByDept(dept.id).length" class="text-gray-500 text-sm py-4 text-center">Staff information coming soon.</p>
                    <div v-else class="space-y-3">
                      <div v-for="member in staffByDept(dept.id)" :key="member.id" class="flex items-center gap-3">
                        <img v-if="member.image_url" :src="mediaUrl(member.image_url)" :alt="member.full_name" class="w-10 h-10 rounded-full object-cover flex-shrink-0" />
                        <div v-else class="w-10 h-10 rounded-full bg-[#1a365d]/10 flex items-center justify-center text-[#1a365d] flex-shrink-0"><i class="icofont-user text-lg" aria-hidden="true"></i></div>
                        <div class="min-w-0"><p class="font-medium text-[#1a365d] truncate">{{ member.full_name }}</p><p class="text-gray-500 text-xs truncate">{{ member.designation }}</p></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="mt-6 p-4 bg-[#1a365d]/5 rounded-lg border-l-4 border-[#f5a623]">
              <p class="text-sm text-[#1a365d]"><strong>Note:</strong> All departments report to the General Secretary, {{ generalSecretaryName }}, who oversees the day-to-day operations of the National Council of Sports Secretariat.</p>
            </div>
          </div>
        </div>

        <div class="order-1 lg:order-2 space-y-6">
          <div class="bg-white rounded-xl p-6 shadow-sm">
            <h3 class="font-bold text-[#1a365d] mb-4">About NCS</h3>
            <nav class="space-y-2">
              <router-link to="/pages/the-mandate" class="flex items-center gap-2 p-3 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="m16 16 3-8 3 8c-.87.65-1.92 1-3 1s-2.13-.35-3-1Z"></path><path d="m2 16 3-8 3 8c-.87.65-1.92 1-3 1s-2.13-.35-3-1Z"></path><path d="M7 21h10"></path><path d="M12 3v18"></path><path d="M3 7h2c2 0 5-1 7-2 2 1 5 2 7 2h2"></path></svg>
                Mandate
              </router-link>
              <router-link to="/governing-council" class="flex items-center gap-2 p-3 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path><path d="M16 3.128a4 4 0 0 1 0 7.744"></path><path d="M22 21v-2a4 4 0 0 0-3-3.87"></path><circle cx="9" cy="7" r="4"></circle></svg>
                The Council
              </router-link>
              <router-link to="/team" class="flex items-center gap-2 p-3 rounded-lg bg-[#f5a623]/10 text-[#f5a623] font-medium">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"></path><path d="M14 2v4a2 2 0 0 0 2 2h4"></path><path d="M10 9H8"></path><path d="M16 13H8"></path><path d="M16 17H8"></path></svg>
                The Staff
              </router-link>
            </nav>
          </div>

          <div class="bg-[#1a365d] rounded-xl p-6 text-white">
            <h3 class="font-bold mb-4">Contact the Secretariat</h3>
            <div class="space-y-3 text-sm">
              <p class="text-white/70">For inquiries and official correspondence, please contact the NCS Secretariat.</p>
              <div class="flex items-center gap-2 text-white/80">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-[#f5a623]" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"></path><rect x="2" y="4" width="20" height="16" rx="2"></rect></svg>
                <a :href="`mailto:${contactEmail}`" class="hover:text-[#f5a623]">{{ contactEmail }}</a>
              </div>
              <div class="flex items-center gap-2 text-white/80">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-[#f5a623]" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>
                <span>{{ contactPhone }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
<script setup>
import { computed, onMounted, ref } from 'vue'
import { getSettings, listInstitutionalDepartments, listTeam } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const departments = ref([]), team = ref([]), loading = ref(true)
const openIds = ref(new Set())
const contactEmail = ref('info@ncs.go.ug')
const contactPhone = ref('+256 414254477')
const generalSecretaryName = ref('Dr. Bernard Patrick Ogwel')

const DEPT_META = {
  'Administration': { color: 'bg-blue-500', icon: 'building2' },
  'Human Resource': { color: 'bg-green-500', icon: 'users' },
  'Finance & Accounts': { color: 'bg-amber-500', icon: 'calculator' },
  'ICT': { color: 'bg-purple-500', icon: 'monitor' },
  'Sports Officers': { color: 'bg-red-500', icon: 'medal' },
  'Engineering & Facilities': { color: 'bg-pink-500', icon: 'heart-pulse' },
  'Public Relations & Communications': { color: 'bg-cyan-500', icon: 'megaphone' },
  'Legal & Compliance': { color: 'bg-slate-500', icon: 'briefcase' },
  'Procurement & Records': { color: 'bg-indigo-500', icon: 'file-text' },
  'Support Services': { color: 'bg-orange-500', icon: 'users' },
}
const DEPT_ORDER = Object.keys(DEPT_META)
function deptMeta(name) {
  return DEPT_META[name] || { color: 'bg-gray-500', icon: 'building2' }
}

const orderedDepartments = computed(() => {
  return [...departments.value].sort((a, b) => {
    const ai = DEPT_ORDER.indexOf(a.name), bi = DEPT_ORDER.indexOf(b.name)
    return (ai === -1 ? DEPT_ORDER.length : ai) - (bi === -1 ? DEPT_ORDER.length : bi)
  })
})

function staffByDept(deptId) {
  return team.value.filter(m => m.department_id === deptId)
}
function toggleDept(id) {
  const next = new Set(openIds.value)
  next.has(id) ? next.delete(id) : next.add(id)
  openIds.value = next
}

onMounted(async () => {
  try {
    const [deptRes, teamRes, contactRes, homepageRes] = await Promise.allSettled([
      listInstitutionalDepartments(),
      listTeam(),
      getSettings('contact'),
      getSettings('homepage'),
    ])
    departments.value = deptRes.status === 'fulfilled' ? (deptRes.value.data?.data || []) : []
    team.value = teamRes.status === 'fulfilled' ? (teamRes.value.data?.data || []) : []
    if (contactRes.status === 'fulfilled') {
      const value = contactRes.value.data?.data?.value || {}
      if (value.email) contactEmail.value = value.email
      if (value.phone) contactPhone.value = value.phone
    }
    if (homepageRes.status === 'fulfilled') {
      const members = homepageRes.value.data?.data?.value?.leadership?.members || []
      const secretary = members.find(m => /secretary/i.test(m.title || ''))
      if (secretary?.name) generalSecretaryName.value = secretary.name
    }
  } finally {
    loading.value = false
  }
})
</script>
