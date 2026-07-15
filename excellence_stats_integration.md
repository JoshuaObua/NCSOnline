# Sports Excellence Stats Integration Guide

This guide provides instructions on integrating the dynamic stats endpoint with the public landing page to replace static placeholders with live database counters.

---

## 1. API Endpoint Details

- **Endpoint**: `/api/v1/cms/stats`
- **Method**: `GET`
- **Authentication**: None (Public Endpoint)
- **Rate Limit**: Governed by the Nginx `portal_cms_limit` zone (120 requests/minute, burst = 20)

---

## 2. API Response Structure

The endpoint returns a structured JSON payload:

```json
{
  "success": true,
  "data": {
    "years_of_excellence": 62,
    "sports_associations": 52,
    "sports_facilities": 10,
    "athletes_reached": 100000
  }
}
```

### Data Fields:
| Field Name | Type | Description |
| :--- | :--- | :--- |
| `years_of_excellence` | `integer` | Dynamic difference: `Current Year - 1964` |
| `sports_associations` | `integer` | Count of registered federations (falls back to active CMS associations if empty) |
| `sports_facilities` | `integer` | Count of active facilities registered in the CMS database |
| `athletes_reached` | `integer` | Count of registered athletes in the database (exact count) |

---

## 3. Frontend Integration Example (Vue 3 / Fetch API)

Here is a ready-to-use integration script to dynamically inject the stats into the homepage template.

### Vue Component Script:
```javascript
import { ref, onMounted } from 'vue'

const stats = ref({
  years_of_excellence: '60+',
  sports_associations: '54+',
  sports_facilities: '32+',
  athletes_reached: '0+'
})

async function fetchStats() {
  try {
    const response = await fetch('/api/v1/cms/stats')
    const result = await response.json()
    if (result.success && result.data) {
      const data = result.data
      
      stats.value = {
        years_of_excellence: `${data.years_of_excellence}+`,
        sports_associations: `${data.sports_associations}+`,
        sports_facilities: `${data.sports_facilities}+`,
        // Format large numbers elegantly or display exact count
        athletes_reached: data.athletes_reached >= 1000 ? `${Math.round(data.athletes_reached / 1000)}K+` : `${data.athletes_reached}+`
      }
    }
  } catch (error) {
    console.error('Failed to load excellence stats:', error)
  }
}

onMounted(() => {
  fetchStats()
})
```

### HTML Template:
```html
<section id="section-stats" class="py-16 md:py-20 bg-[#1a365d] text-white">
  <div class="max-w-7xl mx-auto px-4">
    
    <!-- Header -->
    <div class="text-center max-w-2xl mx-auto mb-12">
      <h2 class="text-3xl md:text-4xl font-bold text-white mb-3">Sports Excellence in Numbers</h2>
      <p class="text-white/70 leading-relaxed">
        Driving the development of sports across Uganda through dedicated programs and world-class facilities
      </p>
    </div>

    <!-- Cards Grid -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
      
      <!-- Card 1: Years of Excellence -->
      <div class="text-center p-6 rounded-xl bg-white/5 backdrop-blur border border-white/10 hover:border-[#f5a623]/40 hover:bg-white/10 transition-all">
        <div class="w-14 h-14 rounded-xl bg-[#f5a623]/20 text-[#f5a623] flex items-center justify-center mx-auto mb-4">
          <i class="icofont-award text-2xl" aria-hidden="true"></i>
        </div>
        <div class="text-3xl md:text-4xl font-bold text-white mb-1">{{ stats.years_of_excellence }}</div>
        <div class="text-sm text-white/70">Years of Excellence</div>
      </div>

      <!-- Card 2: Sports Associations -->
      <div class="text-center p-6 rounded-xl bg-white/5 backdrop-blur border border-white/10 hover:border-[#f5a623]/40 hover:bg-white/10 transition-all">
        <div class="w-14 h-14 rounded-xl bg-[#f5a623]/20 text-[#f5a623] flex items-center justify-center mx-auto mb-4">
          <i class="icofont-trophy text-2xl" aria-hidden="true"></i>
        </div>
        <div class="text-3xl md:text-4xl font-bold text-white mb-1">{{ stats.sports_associations }}</div>
        <div class="text-sm text-white/70">Sports Associations</div>
      </div>

      <!-- Card 3: Sports Facilities -->
      <div class="text-center p-6 rounded-xl bg-white/5 backdrop-blur border border-white/10 hover:border-[#f5a623]/40 hover:bg-white/10 transition-all">
        <div class="w-14 h-14 rounded-xl bg-[#f5a623]/20 text-[#f5a623] flex items-center justify-center mx-auto mb-4">
          <i class="icofont-stadium text-2xl" aria-hidden="true"></i>
        </div>
        <div class="text-3xl md:text-4xl font-bold text-white mb-1">{{ stats.sports_facilities }}</div>
        <div class="text-sm text-white/70">Sports Facilities</div>
      </div>

      <!-- Card 4: Athletes Reached -->
      <div class="text-center p-6 rounded-xl bg-white/5 backdrop-blur border border-white/10 hover:border-[#f5a623]/40 hover:bg-white/10 transition-all">
        <div class="w-14 h-14 rounded-xl bg-[#f5a623]/20 text-[#f5a623] flex items-center justify-center mx-auto mb-4">
          <i class="icofont-users-alt-5 text-2xl" aria-hidden="true"></i>
        </div>
        <div class="text-3xl md:text-4xl font-bold text-white mb-1">{{ stats.athletes_reached }}</div>
        <div class="text-sm text-white/70">Athletes Reached</div>
        <p class="mt-3 text-xs leading-relaxed text-white/55">
          Athletes, coaches, and administrators served through NCS programs.
        </p>
      </div>

    </div>
  </div>
</section>
