# Find Your Sport Search & Filter Integration Guide

This guide provides documentation and integration examples for the public Sports Federation Search & Filter API, enabling real-time lookup of recognized National Sports Associations on the homepage.

---

## 1. API Endpoint Details

- **Endpoint**: `/api/v1/cms/associations`
- **Method**: `GET`
- **Authentication**: None (Public Endpoint)
- **Rate Limit**: Governed by the Nginx `portal_cms_limit` zone (120 requests/minute, burst = 20)

### Query Parameters:
| Parameter | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `active` | `string` | `"true"` | Filters active federations. Set to `"false"` to retrieve inactive records. |
| `search` | `string` | `""` | Filters by search text query (matches name, abbreviation, or description case-insensitively). |
| `category` | `string` | `""` | Filters by federation category type (e.g. `"Ball Sports"`, `"Combat Sports"`, `"Other"`). |

---

## 2. API Response Structure

The endpoint returns an array of recognized national sports associations:

```json
[
  {
    "id": "assoc_fufa",
    "name": "Federation of Uganda Football Associations",
    "slug": "fufa",
    "abbreviation": "FUFA",
    "description": "FUFA is the governing body of football in Uganda, affiliated with FIFA and CAF.",
    "logo_url": "/uploads/logos/fufa.png",
    "website_url": "https://www.fufa.co.ug",
    "category": "Ball Sports",
    "president": "Moses Magogo",
    "secretary": "Edgar Watson",
    "address": "Mengo, Kampala",
    "phone": "+256 414 272700",
    "sort_order": 1,
    "is_active": true,
    "created_at": "2026-07-15T12:00:00Z",
    "updated_at": "2026-07-15T12:00:00Z"
  }
]
```

---

## 3. Frontend Integration Example (Vue 3 / Fetch API)

Below is a complete implementation script and template structure for the "Find Your Sport" interactive search block.

### Vue Component Script:
```javascript
import { ref, onMounted, watch } from 'vue'

const federations = ref([])
const totalCount = ref(0)
const searchQuery = ref('')
const selectedCategory = ref('All Sports')
const categories = ref(['All Sports', 'Ball Sports', 'Combat Sports', 'Other'])
const loading = ref(false)
let searchTimeout = null

async function loadFederations() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    params.append('active', 'true')
    if (searchQuery.value.trim()) {
      params.append('search', searchQuery.value.trim())
    }
    if (selectedCategory.value !== 'All Sports') {
      params.append('category', selectedCategory.value)
    }

    const response = await fetch(`/api/v1/cms/associations?${params.toString()}`)
    const data = await response.json()
    
    federations.value = Array.isArray(data) ? data : []
    totalCount.value = federations.value.length
  } catch (error) {
    console.error('Failed to load sports federations:', error)
  } finally {
    loading.value = false
  }
}

// Watch inputs to reload list dynamically
watch(selectedCategory, () => {
  loadFederations()
})

function onSearchInput() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    loadFederations()
  }, 350) // Debounce search inputs
}

onMounted(() => {
  loadFederations()
})
```

### HTML Template:
```html
<section id="section-find_sport" class="home-section finder-section">
  <div class="home-shell">
    
    <!-- Heading -->
    <div class="section-heading light">
      <span class="section-kicker gold">Discover your federation</span>
      <h2>Find Your Sport</h2>
      <p>
        Search across all National Sports Associations and Federations recognised by NCS. 
        Tap any card to see the president, secretary, address, phone and website.
      </p>
    </div>

    <!-- Searchbox -->
    <div class="finder-search">
      <i class="icofont-search-1" aria-hidden="true"></i>
      <label for="sport-search" class="sr-only">Search for a sport or federation</label>
      <input 
        id="sport-search" 
        type="search" 
        v-model="searchQuery" 
        @input="onSearchInput" 
        placeholder="Try 'rugby', 'football', 'tennis'" 
      />
    </div>

    <!-- Filter Chips -->
    <div class="filter-chips" aria-label="Sport categories">
      <button 
        v-for="cat in categories" 
        :key="cat" 
        type="button" 
        :class="{ active: selectedCategory === cat }"
        @click="selectedCategory = cat"
      >
        {{ cat }}
      </button>
    </div>

    <!-- Meta Details / Count -->
    <div class="finder-meta">
      <span>
        Showing <strong>{{ federations.length }}</strong> of {{ totalCount }} federations
      </span>
      <a href="/associations" class="view-all-link">View directory →</a>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-8 text-white/50">
      Searching recognized sports federations...
    </div>

    <!-- Association Cards Grid -->
    <div v-else class="association-grid">
      <article v-for="fed in federations" :key="fed.id" class="association-card">
        <div class="association-logo">
          <i class="icofont-trophy"></i>
        </div>
        <div>
          <span>{{ fed.category || 'Other' }}</span>
          <h3>{{ fed.name }}</h3>
          <p>{{ fed.description }}</p>
          
          <details>
            <summary>Contact details</summary>
            <ul>
              <li v-if="fed.president"><b>President:</b> {{ fed.president }}</li>
              <li v-if="fed.secretary"><b>Secretary:</b> {{ fed.secretary }}</li>
              <li v-if="fed.address"><b>Address:</b> {{ fed.address }}</li>
              <li v-if="fed.phone"><b>Phone:</b> {{ fed.phone }}</li>
              <li v-if="fed.website_url">
                <b>Website:</b> 
                <a :href="fed.website_url" target="_blank" class="text-[#f5a623] hover:underline">
                  {{ fed.website_url }}
                </a>
              </li>
            </ul>
          </details>
        </div>
      </article>
      
      <!-- Empty state -->
      <div v-if="!federations.length" class="text-center py-8 text-white/50 w-full col-span-full">
        No registered federations match your query.
      </div>
    </div>

  </div>
</section>
```
