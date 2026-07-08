<template>
  <div class="facilities-page min-h-screen bg-gray-50">
    <section class="bg-[#1a365d] py-16 md:py-20">
      <div class="max-w-7xl mx-auto px-4">
        <nav class="flex items-center gap-2 text-white/60 text-sm mb-4" aria-label="Breadcrumb">
          <router-link to="/" class="hover:text-white">Home</router-link>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true">
            <path d="m9 18 6-6-6-6" />
          </svg>
          <span class="text-white">Facilities</span>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true">
            <path d="m9 18 6-6-6-6" />
          </svg>
          <span class="text-white capitalize">{{ selectedRegionLabelLower }}</span>
        </nav>
        <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">Sports Facilities</h1>
        <p class="text-white/80 text-lg max-w-2xl">Explore world-class sports infrastructure across Uganda's regions</p>
      </div>
    </section>

    <section class="bg-white border-b border-gray-200 sticky top-[104px] md:top-[116px] z-30">
      <div class="max-w-7xl mx-auto px-4 py-4">
        <div class="flex flex-wrap gap-2">
          <button
            v-for="region in regions"
            :key="region.id"
            type="button"
            :class="[
              'inline-flex h-8 items-center justify-center gap-2 whitespace-nowrap rounded-md px-3 text-xs font-medium shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623]',
              selectedRegion === region.id
                ? 'bg-[#f5a623] hover:bg-[#e09612] text-white'
                : 'border border-gray-200 text-gray-600 hover:border-[#f5a623] hover:text-[#f5a623]'
            ]"
            @click="setRegion(region.id)"
          >
            {{ region.label }}
          </button>
        </div>
      </div>
    </section>

    <section class="bg-[#f5a623]/10 border-b border-[#f5a623]/20">
      <div class="max-w-7xl mx-auto px-4 py-4">
        <p class="text-[#1a365d] flex items-center gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-[#f5a623] flex-shrink-0" aria-hidden="true">
            <path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0" />
            <circle cx="12" cy="10" r="3" />
          </svg>
          <span>Showing facilities in the <strong class="capitalize">{{ selectedRegionLabelLower }}</strong> of Uganda</span>
        </p>
      </div>
    </section>

    <section class="py-12 md:py-16">
      <div class="max-w-7xl mx-auto px-4">
        <div v-if="loading" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div v-for="i in 6" :key="i" class="bg-white rounded-xl overflow-hidden shadow-sm">
            <div class="h-48 bg-gray-200 animate-pulse"></div>
            <div class="p-6 space-y-3">
              <div class="h-5 bg-gray-200 rounded w-2/3 animate-pulse"></div>
              <div class="h-4 bg-gray-200 rounded w-full animate-pulse"></div>
              <div class="h-4 bg-gray-200 rounded w-1/2 animate-pulse"></div>
            </div>
          </div>
        </div>

        <div v-else-if="filteredFacilities.length" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="facility in filteredFacilities"
            :key="facility.id"
            class="group bg-white rounded-xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300"
          >
            <div class="relative h-48 overflow-hidden">
              <img
                v-if="facility.image"
                :src="facility.image"
                :alt="facility.name"
                class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
              >
              <div v-else class="w-full h-full bg-gradient-to-br from-[#1a365d] to-[#244b80] flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" class="w-12 h-12 text-white/40" aria-hidden="true">
                  <path d="M3 21h18" />
                  <path d="M5 21V7l8-4v18" />
                  <path d="M19 21V11l-6-4" />
                  <path d="M9 9h1" />
                  <path d="M9 13h1" />
                  <path d="M9 17h1" />
                </svg>
              </div>
              <div class="inline-flex items-center rounded-md border border-transparent bg-[#1a365d] px-2.5 py-0.5 text-xs font-semibold text-white shadow absolute top-4 left-4">{{ facility.type }}</div>
              <div class="inline-flex items-center rounded-md border border-transparent bg-green-500 px-2.5 py-0.5 text-xs font-semibold text-white shadow absolute top-4 right-4">{{ facility.status }}</div>
            </div>

            <div class="p-6">
              <h3 class="font-bold text-xl text-[#1a365d] mb-2 group-hover:text-[#f5a623] transition-colors">{{ facility.name }}</h3>
              <p class="text-gray-600 text-sm mb-4 line-clamp-2">{{ facility.description }}</p>

              <div class="flex items-center gap-2 text-gray-500 text-sm mb-2">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-[#f5a623] flex-shrink-0" aria-hidden="true">
                  <path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0" />
                  <circle cx="12" cy="10" r="3" />
                </svg>
                <span>{{ facility.location }}</span>
              </div>

              <div v-if="facility.amenities.length" class="flex flex-wrap gap-1 mt-4">
                <div
                  v-for="amenity in visibleAmenities(facility)"
                  :key="amenity"
                  class="inline-flex items-center rounded-md border border-transparent bg-gray-100 px-2.5 py-0.5 text-xs font-semibold text-gray-700"
                >
                  {{ amenity }}
                </div>
                <div
                  v-if="remainingAmenityCount(facility) > 0"
                  class="inline-flex items-center rounded-md border border-transparent bg-gray-100 px-2.5 py-0.5 text-xs font-semibold text-gray-700"
                >
                  +{{ remainingAmenityCount(facility) }} more
                </div>
              </div>

              <div class="mt-4 pt-4 border-t border-gray-200 flex items-center justify-between">
                <a
                  v-if="facility.phone"
                  :href="`tel:${facility.phone}`"
                  class="flex items-center gap-1 text-sm text-gray-500 hover:text-[#f5a623]"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true">
                    <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z" />
                  </svg>
                  Contact
                </a>
                <span v-else></span>
                <a
                  v-if="facility.email"
                  :href="`mailto:${facility.email}`"
                  class="flex items-center gap-1 text-sm text-gray-500 hover:text-[#f5a623]"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true">
                    <path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7" />
                    <rect x="2" y="4" width="20" height="16" rx="2" />
                  </svg>
                  Email
                </a>
                <span v-else></span>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="text-center py-20 text-gray-400">
          <p>No facilities listed for this region yet.</p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { listFacilities } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const route = useRoute()
const router = useRouter()

const regions = [
  { id: 'central', label: 'Central Region' },
  { id: 'northern', label: 'Northern Region' },
  { id: 'eastern', label: 'Eastern Region' },
  { id: 'western', label: 'Western Region' },
  { id: 'southern', label: 'Southern Region' },
]

const facilities = ref([])
const loading = ref(false)

const fallbackFacilities = [
  {
    id: 'cricket-oval',
    name: 'Cricket Oval',
    slug: 'cricket-oval',
    type: 'Pitch',
    status: 'Available',
    region: 'central',
    description: 'International standard cricket facility with modern amenities',
    location: 'Lugogo, Kampala',
    amenities: ['Practice nets', 'Pavilion', 'Changing rooms', 'Scoreboard', 'Parking'],
    image: 'https://images.unsplash.com/photo-1531415074968-036ba1b575da?w=600&h=400&fit=crop',
    phone: '+256414254477',
    email: 'cricket@ncs.go.ug',
  },
  {
    id: 'gymnasium',
    name: 'Gymnasium',
    slug: 'gymnasium',
    type: 'Gymnasium',
    status: 'Available',
    region: 'central',
    description: 'Fully equipped gymnasium for athlete training and fitness',
    location: 'Lugogo Sports Complex, Kampala',
    amenities: ['Cardio equipment', 'Weight training', 'Personal training', 'Locker rooms', 'Showers'],
    image: 'https://images.unsplash.com/photo-1534438327276-14e5300c3a48?w=600&h=400&fit=crop',
    phone: '+256414254477',
    email: 'gym@ncs.go.ug',
  },
  {
    id: 'hockey-pitch',
    name: 'Hockey Pitch',
    slug: 'hockey-pitch',
    type: 'Pitch',
    status: 'Available',
    region: 'central',
    description: 'Synthetic turf hockey pitch meeting international standards',
    location: 'Lugogo, Kampala',
    amenities: ['Synthetic turf', 'Floodlights', 'Spectator stands', 'Changing rooms', 'Parking'],
    image: 'https://images.unsplash.com/photo-1546519638-68e109498ffc?w=600&h=400&fit=crop',
    phone: '+256414254477',
    email: 'hockey@ncs.go.ug',
  },
  {
    id: 'indoor-stadium',
    name: 'Indoor Stadium',
    slug: 'indoor-stadium',
    type: 'Stadium',
    status: 'Available',
    region: 'central',
    description: 'Multi-purpose indoor facility for various sports events',
    location: 'Lugogo Sports Complex, Kampala',
    amenities: ['Air conditioning', 'Modern lighting', '2000 spectator seats', 'Changing rooms', 'Media room'],
    image: 'https://images.unsplash.com/photo-1504450758481-7338eba7524a?w=600&h=400&fit=crop',
    phone: '+256414254477',
    email: 'indoor@ncs.go.ug',
  },
  {
    id: 'volleyball-courts',
    name: 'Volleyball Courts',
    slug: 'volleyball-courts',
    type: 'Court',
    status: 'Available',
    region: 'central',
    description: 'Professional volleyball courts for training and competitions',
    location: 'Lugogo, Kampala',
    amenities: ['Indoor courts', 'Outdoor courts', 'Floodlights', 'Scoreboards', 'Seating'],
    image: 'https://images.unsplash.com/photo-1612872087720-bb876e2e67d1?w=600&h=400&fit=crop',
    phone: '+256414254477',
    email: 'volleyball@ncs.go.ug',
  },
  {
    id: 'sports-shop',
    name: 'Sports Shop',
    slug: 'sports-shop',
    type: 'Shop',
    status: 'Available',
    region: 'central',
    description: 'Quality sports equipment and merchandise available',
    location: 'Lugogo Sports Complex, Kampala',
    amenities: ['Sports equipment', 'Athletic apparel', 'Team uniforms', 'Merchandise', 'Accessories'],
    image: 'https://images.unsplash.com/photo-1526948128573-703ee1aeb6fa?w=600&h=400&fit=crop',
    phone: '+256414254477',
    email: 'shop@ncs.go.ug',
  },
  {
    id: 'test-iteration7-facility',
    name: 'TEST_Iteration7_Facility',
    slug: 'test-iteration7-facility',
    type: 'Stadium',
    status: 'Available',
    region: 'central',
    description: 'Test facility description',
    location: 'Kampala',
    amenities: [],
    image: 'https://images.unsplash.com/photo-1540747913346-19e32dc3e97e?w=600&h=400&fit=crop',
    phone: '',
    email: '',
  },
]

const fallbackByName = new Map(fallbackFacilities.map(facility => [facility.name.toLowerCase(), facility]))

const selectedRegion = computed(() => normalizeRegionParam(route.query.region))
const selectedRegionLabel = computed(() => regions.find(region => region.id === selectedRegion.value)?.label || 'Central Region')
const selectedRegionLabelLower = computed(() => selectedRegionLabel.value.charAt(0).toLowerCase() + selectedRegionLabel.value.slice(1))

const filteredFacilities = computed(() => displayFacilities.value.filter(facility => facility.region === selectedRegion.value))

const displayFacilities = computed(() => {
  if (!facilities.value.length) return fallbackFacilities
  return facilities.value.map(normalizeFacility)
})

function normalizeFacility(item) {
  const fallback = fallbackByName.get(String(item.name || '').toLowerCase()) || {}
  return {
    id: item.id || fallback.id || item.slug || item.name,
    name: item.name || fallback.name || 'Sports Facility',
    slug: item.slug || fallback.slug || '',
    type: item.category || fallback.type || 'Facility',
    status: item.is_active === false ? 'Unavailable' : 'Available',
    region: inferRegion(item.category || fallback.region),
    description: item.description || fallback.description || 'Sports facility managed by the National Council of Sports Uganda',
    location: item.location || fallback.location || 'Kampala',
    amenities: Array.isArray(item.amenities) ? item.amenities : (fallback.amenities || []),
    image: item.image_url ? mediaUrl(item.image_url) : (fallback.image || ''),
    phone: item.phone || fallback.phone || '+256414254477',
    email: item.email || fallback.email || '',
  }
}

function inferRegion(value = '') {
  const text = String(value).toLowerCase()
  return regions.find(region => text.includes(region.id))?.id || 'central'
}

function normalizeRegionParam(value) {
  const raw = Array.isArray(value) ? value[0] : value
  const normalized = String(raw || 'central').trim().toLowerCase().replace(/\s+/g, '-').replace(/-region$/, '')
  return regions.some(region => region.id === normalized) ? normalized : 'central'
}

function setRegion(regionId) {
  const region = normalizeRegionParam(regionId)
  if (route.query.region && normalizeRegionParam(route.query.region) === region) return
  router.push({ path: '/facilities', query: { ...route.query, region } })
}

function visibleAmenities(facility) {
  return facility.amenities.slice(0, 3)
}

function remainingAmenityCount(facility) {
  return Math.max(facility.amenities.length - 3, 0)
}

onMounted(() => {
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), 1200)
  listFacilities({}, { signal: controller.signal }).then(response => {
    const items = response?.data?.data
    if (Array.isArray(items) && items.length) facilities.value = items
  }).catch(() => {
    facilities.value = []
  }).finally(() => {
    clearTimeout(timeout)
  })
})

watch(
  () => route.query.region,
  value => {
    const region = normalizeRegionParam(value)
    if (value && value !== region) {
      router.replace({ path: '/facilities', query: { ...route.query, region } })
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.facilities-page :is(h1, h2, h3, h4) {
  letter-spacing: 0;
}

.facilities-page :is(a, button) {
  text-decoration: none;
}
</style>
