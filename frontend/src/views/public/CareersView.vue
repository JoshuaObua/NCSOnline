<template>
  <div class="careers-page bg-white">
    <section class="bg-[#1a365d] py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4">
        <nav class="flex items-center gap-2 text-white/60 text-sm mb-4" aria-label="Breadcrumb">
          <router-link to="/" class="hover:text-white">Home</router-link>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true">
            <path d="m9 18 6-6-6-6" />
          </svg>
          <span class="text-white">Careers</span>
        </nav>
        <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">Careers at NCS</h1>
        <p class="text-white/80 text-lg max-w-2xl">Explore job opportunities and business tenders with the National Council of Sports</p>
      </div>
    </section>

    <section class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="grid md:grid-cols-2 gap-4">
          <a href="#jobs" class="flex items-center justify-between p-4 bg-[#1a365d]/5 rounded-xl hover:bg-[#1a365d]/10 transition-colors group">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-lg bg-[#1a365d] flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true">
                  <path d="M16 20V4a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16" />
                  <rect width="20" height="14" x="2" y="6" rx="2" />
                </svg>
              </div>
              <div>
                <h3 class="font-semibold text-[#1a365d]">Job Opportunities</h3>
                <p class="text-sm text-gray-500">{{ jobs.length }} open positions</p>
              </div>
            </div>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#1a365d] group-hover:translate-x-1 transition-transform" aria-hidden="true">
              <path d="M5 12h14" />
              <path d="m12 5 7 7-7 7" />
            </svg>
          </a>

          <a href="#tenders" class="flex items-center justify-between p-4 bg-[#f5a623]/5 rounded-xl hover:bg-[#f5a623]/10 transition-colors group">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-lg bg-[#f5a623] flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-white" aria-hidden="true">
                  <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z" />
                  <path d="M14 2v4a2 2 0 0 0 2 2h4" />
                  <path d="M10 9H8" />
                  <path d="M16 13H8" />
                  <path d="M16 17H8" />
                </svg>
              </div>
              <div>
                <h3 class="font-semibold text-[#1a365d]">Tenders &amp; Procurement</h3>
                <p class="text-sm text-gray-500">{{ tenders.length }} active tenders</p>
              </div>
            </div>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623] group-hover:translate-x-1 transition-transform" aria-hidden="true">
              <path d="M5 12h14" />
              <path d="m12 5 7 7-7 7" />
            </svg>
          </a>
        </div>
      </div>
    </section>

    <section id="jobs" class="py-12 md:py-16">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex items-center justify-between mb-8">
          <div>
            <h2 class="text-2xl md:text-3xl font-bold text-[#1a365d] flex items-center gap-3">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-7 h-7 text-[#f5a623]" aria-hidden="true">
                <path d="M16 20V4a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16" />
                <rect width="20" height="14" x="2" y="6" rx="2" />
              </svg>
              Job Opportunities
            </h2>
            <p class="text-gray-600 mt-1">Current openings at the National Council of Sports</p>
          </div>
        </div>

        <div class="space-y-4">
          <div v-for="job in jobs" :key="job.title" class="bg-white rounded-xl p-6 shadow-sm hover:shadow-lg transition-all border border-gray-100 hover:border-[#f5a623]/30">
            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-2">
                  <div class="inline-flex items-center rounded-md border border-transparent bg-gray-100 px-2.5 py-0.5 text-xs font-semibold text-gray-700 transition-colors">{{ job.type }}</div>
                  <div class="inline-flex items-center rounded-md border border-transparent bg-[#1a365d]/10 px-2.5 py-0.5 text-xs font-semibold text-[#1a365d] shadow transition-colors">{{ job.department }}</div>
                </div>
                <h3 class="text-xl font-semibold text-[#1a365d] mb-2">{{ job.title }}</h3>
                <p class="text-gray-600 text-sm mb-3">{{ job.description }}</p>
                <div class="flex flex-wrap items-center gap-4 text-sm text-gray-500">
                  <span class="flex items-center gap-1">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true">
                      <path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0" />
                      <circle cx="12" cy="10" r="3" />
                    </svg>
                    {{ job.location }}
                  </span>
                  <span class="flex items-center gap-1">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true">
                      <circle cx="12" cy="12" r="10" />
                      <polyline points="12 6 12 12 16 14" />
                    </svg>
                    Deadline: {{ job.deadline }}
                  </span>
                </div>
              </div>

              <router-link to="/login">
                <button type="button" class="inline-flex h-9 items-center justify-center gap-2 whitespace-nowrap rounded-md bg-[#1a365d] px-4 py-2 text-sm font-medium text-white shadow transition-colors hover:bg-[#1a365d]/90 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623]">
                  Apply Now
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 ml-2" aria-hidden="true">
                    <path d="M5 12h14" />
                    <path d="m12 5 7 7-7 7" />
                  </svg>
                </button>
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="tenders" class="py-12 md:py-16 bg-gray-100">
      <div class="max-w-7xl mx-auto px-4">
        <div class="flex items-center justify-between mb-8">
          <div>
            <h2 class="text-2xl md:text-3xl font-bold text-[#1a365d] flex items-center gap-3">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-7 h-7 text-[#f5a623]" aria-hidden="true">
                <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z" />
                <path d="M14 2v4a2 2 0 0 0 2 2h4" />
                <path d="M10 9H8" />
                <path d="M16 13H8" />
                <path d="M16 17H8" />
              </svg>
              Tenders &amp; Procurement
            </h2>
            <p class="text-gray-600 mt-1">Current business opportunities with NCS</p>
          </div>
          <router-link to="/invest">
            <button type="button" class="inline-flex h-9 items-center justify-center gap-2 whitespace-nowrap rounded-md border border-[#f5a623] px-4 py-2 text-sm font-medium text-[#f5a623] shadow-sm transition-colors hover:bg-[#f5a623] hover:text-white focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623]">Submit Proposal</button>
          </router-link>
        </div>

        <div class="space-y-4">
          <div v-for="tender in tenders" :key="tender.ref" class="bg-white rounded-xl p-6 shadow-sm hover:shadow-lg transition-all border border-gray-100 hover:border-[#f5a623]/30">
            <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-2">
                  <div class="inline-flex items-center rounded-md border border-transparent bg-green-100 px-2.5 py-0.5 text-xs font-semibold text-green-700 shadow transition-colors">{{ tender.status }}</div>
                  <div class="inline-flex items-center rounded-md border border-transparent bg-gray-100 px-2.5 py-0.5 text-xs font-semibold text-gray-700 transition-colors">{{ tender.category }}</div>
                </div>
                <h3 class="text-xl font-semibold text-[#1a365d] mb-2">{{ tender.title }}</h3>
                <div class="flex flex-wrap items-center gap-4 text-sm text-gray-500">
                  <span class="flex items-center gap-1">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true">
                      <path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z" />
                      <path d="M14 2v4a2 2 0 0 0 2 2h4" />
                      <path d="M10 9H8" />
                      <path d="M16 13H8" />
                      <path d="M16 17H8" />
                    </svg>
                    Ref: {{ tender.ref }}
                  </span>
                  <span class="flex items-center gap-1">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true">
                      <path d="M8 2v4" />
                      <path d="M16 2v4" />
                      <rect width="18" height="18" x="3" y="4" rx="2" />
                      <path d="M3 10h18" />
                    </svg>
                    Closing: {{ tender.closing }}
                  </span>
                </div>
              </div>

              <router-link to="/invest">
                <button type="button" class="inline-flex h-9 items-center justify-center gap-2 whitespace-nowrap rounded-md bg-[#f5a623] px-4 py-2 text-sm font-medium text-white shadow transition-colors hover:bg-[#e09612] focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623]">
                  Submit Bid
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 ml-2" aria-hidden="true">
                    <path d="M5 12h14" />
                    <path d="m12 5 7 7-7 7" />
                  </svg>
                </button>
              </router-link>
            </div>
          </div>
        </div>

        <div class="mt-8 bg-[#1a365d] rounded-xl p-6 text-white">
          <div class="flex items-start gap-4">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-8 h-8 text-[#f5a623] flex-shrink-0" aria-hidden="true">
              <path d="M6 22V4a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v18Z" />
              <path d="M6 12H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h2" />
              <path d="M18 9h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-2" />
              <path d="M10 6h4" />
              <path d="M10 10h4" />
              <path d="M10 14h4" />
              <path d="M10 18h4" />
            </svg>
            <div>
              <h3 class="font-semibold text-lg mb-2">Procurement Information</h3>
              <p class="text-white/80 text-sm mb-4">All tender submissions must be accompanied by valid registration documents. For detailed tender documents and submission guidelines, please register or login to your account.</p>
              <div class="flex flex-wrap gap-3">
                <router-link to="/register">
                  <button type="button" class="inline-flex h-9 items-center justify-center gap-2 whitespace-nowrap rounded-md bg-[#f5a623] px-4 py-2 text-sm font-medium text-white shadow transition-colors hover:bg-[#e09612] focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-white">Register Now</button>
                </router-link>
                <router-link to="/login">
                  <button type="button" class="inline-flex h-9 items-center justify-center gap-2 whitespace-nowrap rounded-md border border-white px-4 py-2 text-sm font-medium text-white shadow-sm transition-colors hover:bg-white hover:text-[#1a365d] focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-white">Login</button>
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
const jobs = [
  {
    type: 'Full-time',
    department: 'Sports Officers',
    title: 'Sports Development Officer',
    description: 'Coordinate sports programs and athlete development initiatives across regions.',
    location: 'Kampala',
    deadline: 'January 30, 2026',
  },
  {
    type: 'Full-time',
    department: 'Finance',
    title: 'Finance Officer',
    description: 'Manage financial operations and budget planning for NCS programs.',
    location: 'Kampala',
    deadline: 'February 15, 2026',
  },
  {
    type: 'Full-time',
    department: 'ICT',
    title: 'ICT Support Specialist',
    description: 'Provide technical support and maintain NCS IT infrastructure.',
    location: 'Kampala',
    deadline: 'February 28, 2026',
  },
]

const tenders = [
  {
    status: 'Open',
    category: 'Goods',
    title: 'Supply of Sports Equipment for Regional Centers',
    ref: 'NCS/PROC/2026/001',
    closing: 'February 10, 2026',
  },
  {
    status: 'Open',
    category: 'Works',
    title: 'Construction of Volleyball Courts - Eastern Region',
    ref: 'NCS/PROC/2026/002',
    closing: 'February 25, 2026',
  },
  {
    status: 'Open',
    category: 'Services',
    title: 'Provision of Catering Services for National Events',
    ref: 'NCS/PROC/2026/003',
    closing: 'March 5, 2026',
  },
]
</script>

<style scoped>
.careers-page :is(h1, h2, h3, h4) {
  letter-spacing: 0;
}

.careers-page :is(a, button) {
  text-decoration: none;
}
</style>
