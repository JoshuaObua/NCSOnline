<template>
  <div class="home-redesign">
    <template v-for="section in visibleSections" :key="section.id">
      <section v-if="section.id === 'hero'" id="section-hero" class="relative h-[500px] md:h-[600px] lg:h-[700px] overflow-hidden -mt-28" aria-roledescription="carousel" aria-label="NCS highlights">
        <div :key="activeSlide.id || currentSlide" class="absolute inset-0 transition-all duration-700 ease-in-out opacity-100 scale-100 z-10">
          <!-- Dynamic background: image from CMS via bg-cover, or gradient fallback -->
          <div
            class="absolute inset-0 bg-cover bg-center transition-opacity duration-700"
            :style="activeSlide.image_url ? { backgroundImage: `url('${mediaUrl(activeSlide.image_url)}')` } : null"
            :class="!activeSlide.image_url ? 'bg-gradient-to-br from-[#1a365d] to-[#274d7e]' : ''"
            role="img"
            :aria-label="activeSlide.title || 'National Council of Sports Uganda'"
          ></div>
          <!-- Readability overlay -->
          <div class="absolute inset-0 bg-gradient-to-r from-[#1a365d]/90 via-[#1a365d]/70 to-transparent pointer-events-none"></div>
          <div class="relative z-20 h-full max-w-7xl mx-auto px-4 flex items-center">
            <div class="max-w-2xl transition-all duration-700 delay-200 translate-y-0 opacity-100">
              <span class="inline-block px-4 py-1.5 bg-[#f5a623] text-white text-sm font-semibold rounded-full mb-4">{{ activeSlide.subtitle || 'National Council of Sports' }}</span>
              <h1 class="text-3xl md:text-5xl lg:text-6xl font-bold text-white mb-4 leading-tight">{{ activeSlide.title || 'Welcome to National Council of Sports' }}</h1>
              <p class="text-lg md:text-xl text-white/90 mb-8">{{ activeSlide.description || 'A Centre of Excellence for Promotion and Development of Sports' }}</p>
              <div class="flex flex-wrap gap-4">
                <router-link :to="activeSlide.button_url || '/pages/the-mandate'" class="inline-flex items-center justify-center bg-[#f5a623] hover:bg-[#e09612] text-white px-8 py-3 text-lg font-semibold rounded-lg shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">{{ activeSlide.button_text || 'Learn More' }}</router-link>
                <router-link to="/my-portal" class="inline-flex items-center justify-center gap-2 border-2 border-white text-white hover:bg-white hover:text-[#1a365d] px-8 py-3 text-lg font-semibold rounded-lg transition-all duration-300 hover:-translate-y-1 bg-transparent">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" x2="3" y1="12" y2="12"/></svg>
                  My Portal
                </router-link>
              </div>
            </div>
          </div>
        </div>
        <button v-if="displaySlides.length > 1" type="button" class="absolute left-4 md:left-8 top-1/2 -translate-y-1/2 z-30 w-12 h-12 md:w-14 md:h-14 rounded-full bg-white/20 backdrop-blur-sm border border-white/30 flex items-center justify-center text-white hover:bg-[#f5a623] hover:border-[#f5a623] transition-all duration-300 group" aria-label="Previous slide" @click="moveSlide(-1)">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 group-hover:scale-110 transition-transform" aria-hidden="true"><path d="m15 18-6-6 6-6"></path></svg>
        </button>
        <button v-if="displaySlides.length > 1" type="button" class="absolute right-4 md:right-8 top-1/2 -translate-y-1/2 z-30 w-12 h-12 md:w-14 md:h-14 rounded-full bg-white/20 backdrop-blur-sm border border-white/30 flex items-center justify-center text-white hover:bg-[#f5a623] hover:border-[#f5a623] transition-all duration-300 group" aria-label="Next slide" @click="moveSlide(1)">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 group-hover:scale-110 transition-transform" aria-hidden="true"><path d="m9 18 6-6-6-6"></path></svg>
        </button>
        <div v-if="displaySlides.length > 1" class="absolute bottom-8 left-1/2 -translate-x-1/2 z-30 flex gap-3">
          <button v-for="(_, index) in displaySlides" :key="index" type="button" :class="['transition-all duration-300 rounded-full h-3', index===currentSlide ? 'w-10 bg-[#f5a623]' : 'w-3 bg-white/50 hover:bg-white/80']" :aria-label="`Show slide ${index + 1}`" @click="setSlide(index)"></button>
        </div>
        <div class="absolute bottom-0 left-0 right-0 h-24 bg-gradient-to-t from-white to-transparent z-20 pointer-events-none"></div>
      </section>

      <section v-else-if="section.id === 'about'" id="section-about" class="py-16 md:py-24 bg-white overflow-hidden">
        <div class="max-w-7xl mx-auto px-4">
          <div class="grid lg:grid-cols-2 gap-12 lg:gap-16 items-start">

            <!-- LEFT 50% -->
            <div>
              <span class="inline-block px-4 py-1.5 bg-[#f5a623]/10 text-[#f5a623] text-sm font-semibold rounded-full mb-4">{{ home.about.eyebrow }}</span>
              <h2 class="text-3xl md:text-4xl font-bold text-[#1a365d] mb-6">{{ home.about.title }}</h2>
              <p class="text-gray-600 text-lg leading-relaxed mb-6">{{ home.about.intro }}</p>
              <p class="text-gray-600 leading-relaxed mb-6" v-html="home.about.body"></p>

              <!-- Leadership card -->
              <div class="bg-gray-50 rounded-xl p-6 mb-8">
                <h3 class="font-semibold text-[#1a365d] mb-4 flex items-center gap-2">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623]" aria-hidden="true"><rect width="16" height="20" x="4" y="2" rx="2" ry="2"/><path d="M9 22v-4h6v4"/><path d="M8 6h.01"/><path d="M16 6h.01"/><path d="M12 6h.01"/><path d="M12 10h.01"/><path d="M12 14h.01"/><path d="M16 10h.01"/><path d="M16 14h.01"/><path d="M8 10h.01"/><path d="M8 14h.01"/></svg>
                  Current Leadership
                </h3>
                <div class="grid gap-4 md:grid-cols-[auto,1fr] md:items-center">
                  <img v-if="home.leadership.chairperson_image" :src="mediaUrl(home.leadership.chairperson_image)" :alt="home.leadership.chairperson_name" class="w-20 h-20 rounded-full object-cover border-2 border-white shadow-md" />
                  <div v-else class="w-20 h-20 rounded-full bg-[#1a365d]/10 flex items-center justify-center text-[#1a365d] border-2 border-white shadow-md"><i class="icofont-user-alt-3 text-2xl" aria-hidden="true"></i></div>
                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div><p class="text-sm text-gray-500">Chairperson</p><p class="font-medium text-[#1a365d]">{{ home.leadership.chairperson_name }}</p></div>
                    <div><p class="text-sm text-gray-500">General Secretary</p><p class="font-medium text-[#1a365d]">{{ home.leadership.secretary_name }}</p></div>
                  </div>
                </div>
                <div class="mt-4">
                  <router-link :to="home.about.leadership_url" class="inline-flex items-center gap-1 text-sm font-semibold text-[#f5a623] hover:text-[#e09612] transition-colors">
                    {{ home.about.leadership_label }}
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="m9 18 6-6-6-6"/></svg>
                  </router-link>
                </div>
              </div>

              <!-- Milestones -->
              <div class="flex flex-wrap gap-8">
                <div v-for="item in home.milestones" :key="item.label">
                  <div class="text-4xl font-bold text-[#1a365d]">{{ stripPlus(item.value) }}<span class="text-[#f5a623]">+</span></div>
                  <div class="text-gray-500 text-sm">{{ item.label }}</div>
                </div>
              </div>
            </div>

            <!-- RIGHT 50% — value cards 2x3 -->
            <div class="grid grid-cols-2 gap-4">
              <article v-for="item in home.values" :key="item.title" :class="['group p-5 rounded-xl border border-gray-100 hover:border-[#f5a623]/30 hover:shadow-lg transition-all duration-300 cursor-pointer', item.featured ? 'bg-[#1a365d] text-white' : 'bg-gray-50 hover:bg-white']">
                <div :class="['w-12 h-12 rounded-lg flex items-center justify-center mb-4 transition-colors', item.featured ? 'bg-[#f5a623]/20 text-[#f5a623]' : 'bg-[#1a365d]/10 text-[#1a365d] group-hover:bg-[#f5a623]/10 group-hover:text-[#f5a623]']">
                  <i :class="[item.icon, 'text-xl']" aria-hidden="true"></i>
                </div>
                <h3 :class="['font-semibold mb-2', item.featured ? 'text-white' : 'text-[#1a365d]']">{{ item.title }}</h3>
                <p :class="['text-sm', item.featured ? 'text-white/70' : 'text-gray-500']">{{ item.text }}</p>
              </article>
            </div>
          </div>

          <!-- Core Functions -->
          <div class="mt-16 pt-16 border-t border-gray-100">
            <div class="text-center mb-10">
              <h3 class="text-2xl md:text-3xl font-bold text-[#1a365d] mb-4">{{ home.about.core_title }}</h3>
              <p class="text-gray-600 max-w-2xl mx-auto">{{ home.about.core_intro }}</p>
            </div>
            <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div v-for="(item, index) in home.core_functions" :key="item" class="flex items-start gap-3 p-4 rounded-xl border border-gray-100 bg-white hover:border-[#f5a623]/30 hover:shadow-md transition-all">
                <span class="flex-none w-8 h-8 rounded-full bg-[#f5a623] text-[#1a365d] font-bold flex items-center justify-center text-sm">{{ index + 1 }}</span>
                <span class="text-sm text-gray-700 leading-snug">{{ item }}</span>
              </div>
            </div>
            <div class="text-center mt-8">
              <router-link :to="home.about.mandate_url" class="inline-flex items-center justify-center gap-2 border shadow-sm px-4 py-2 rounded-md text-sm font-medium border-[#1a365d] text-[#1a365d] hover:bg-[#1a365d] hover:text-white transition-colors">
                {{ home.about.mandate_label }}
              </router-link>
            </div>
          </div>
        </div>
      </section>

      <section v-else-if="section.id === 'stats'" id="section-stats" class="py-16 md:py-20 bg-[#1a365d] text-white">
        <div class="max-w-7xl mx-auto px-4">
          <div class="text-center max-w-2xl mx-auto mb-12">
            <h2 class="text-3xl md:text-4xl font-bold text-white mb-3">{{ home.stats_title }}</h2>
            <p class="text-white/70 leading-relaxed">{{ home.stats_intro }}</p>
          </div>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div v-for="fact in displayedFacts.slice(0,4)" :key="fact.id || fact.label" class="text-center p-6 rounded-xl bg-white/5 backdrop-blur border border-white/10 hover:border-[#f5a623]/40 hover:bg-white/10 transition-all">
              <div class="w-14 h-14 rounded-xl bg-[#f5a623]/20 text-[#f5a623] flex items-center justify-center mx-auto mb-4">
                <i :class="[fact.icon || 'icofont-chart-growth', 'text-2xl']" aria-hidden="true"></i>
              </div>
              <div class="text-3xl md:text-4xl font-bold text-white mb-1">{{ fact.value }}</div>
              <div class="text-sm text-white/70">{{ fact.label }}</div>
            </div>
          </div>
        </div>
      </section>

      <section v-else-if="section.id === 'news'" id="section-news" class="home-section bg-soft">
        <div class="home-shell"><div class="section-heading row-heading"><div><span class="section-kicker">Latest updates</span><h2>News from NCS</h2></div><router-link to="/news">View all news →</router-link></div>
          <div class="card-grid card-grid-3">
            <article v-for="post in posts.slice(0,3)" :key="post.id" class="content-card">
              <img v-if="post.cover_image_url" :src="mediaUrl(post.cover_image_url)" :alt="post.title" /><div v-else class="card-placeholder"><i class="icofont-newspaper"></i></div>
              <div class="content-card-body"><span>{{ post.category || 'News' }}</span><h3><router-link :to="`/news/${post.slug}`">{{ post.title }}</router-link></h3><p>{{ post.excerpt }}</p></div>
            </article>
            <p v-if="!posts.length" class="empty-state">News will appear here when published.</p>
          </div>
        </div>
      </section>

      <section v-else-if="section.id === 'find_sport'" id="section-find_sport" class="home-section finder-section">
        <div class="home-shell"><div class="section-heading light"><span class="section-kicker gold">{{ home.finder_eyebrow }}</span><h2>{{ home.finder_title }}</h2><p>{{ home.finder_intro }}</p></div>
          <div class="finder-search"><i class="icofont-search-1" aria-hidden="true"></i><label for="sport-search" class="sr-only">Search for a sport or federation</label><input id="sport-search" v-model="sportQuery" type="search" placeholder='Try "rugby", "football", "tennis"' /></div>
          <div class="filter-chips" aria-label="Sport categories"><button v-for="category in sportCategories" :key="category" type="button" :class="{ active:category===activeCategory }" @click="activeCategory=category">{{ category }}</button></div>
          <div class="finder-meta"><span>Showing <strong>{{ filteredAssociations.length }}</strong> of {{ associations.length }} federations</span><router-link to="/associations">View directory →</router-link></div>
          <div class="association-grid">
            <article v-for="item in filteredAssociations.slice(0,8)" :key="item.id" class="association-card">
              <img v-if="item.logo_url" :src="mediaUrl(item.logo_url)" :alt="`${item.name} logo`" /><div v-else class="association-logo"><i class="icofont-trophy"></i></div>
              <div><span>{{ item.category || 'Sport Federation' }}</span><h3>{{ item.name }}</h3><p>{{ item.description }}</p>
                <details><summary>Contact details</summary><ul><li v-if="item.president"><b>President:</b> {{ item.president }}</li><li v-if="item.secretary"><b>Secretary:</b> {{ item.secretary }}</li><li v-if="item.address"><b>Address:</b> {{ item.address }}</li><li v-if="item.phone"><b>Phone:</b> {{ item.phone }}</li></ul><a v-if="item.website_url" :href="item.website_url" target="_blank" rel="noopener">Visit website</a></details>
              </div>
            </article>
            <div v-if="!filteredAssociations.length" class="finder-empty"><i class="icofont-trophy"></i><p>No federations match your search. Try a different sport name.</p></div>
          </div>
        </div>
      </section>

      <section v-else-if="section.id === 'get_involved'" id="section-get-involved" class="home-section bg-white">
        <div class="home-shell involved-grid">
          <div><span class="section-kicker">{{ home.involved_title }}</span><h2>{{ home.involved_subtitle }}</h2><p class="home-lead">{{ home.involved_text }}</p><div class="home-actions"><router-link :to="home.register_url" class="home-btn home-btn-gold">{{ home.register_label }}</router-link><router-link :to="home.contact_url" class="home-btn home-btn-outline-navy">{{ home.contact_label }}</router-link></div></div>
          <aside class="quick-contact"><h3>Quick Contact</h3><div><i class="icofont-location-pin"></i><p><b>Visit Us</b><span>{{ contact.address }}</span><span>{{ contact.postal_address }}</span></p></div><div><i class="icofont-phone"></i><p><b>Call Us</b><span>{{ contact.phone }}</span><span v-if="contact.fax">Fax: {{ contact.fax }}</span></p></div><div><i class="icofont-email"></i><p><b>Email Us</b><a :href="`mailto:${contact.email}`">{{ contact.email }}</a></p></div></aside>
        </div>
      </section>

      <section v-else-if="section.id === 'events'" id="section-events" class="home-section bg-soft">
        <div class="home-shell"><div class="section-heading row-heading"><div><span class="section-kicker">Calendar</span><h2>Upcoming Events</h2></div><router-link to="/events">All events →</router-link></div><div class="card-grid card-grid-3"><article v-for="event in events.slice(0,3)" :key="event.id" class="event-card"><time :datetime="event.event_date">{{ formatDate(event.event_date) }}</time><div><h3>{{ event.title }}</h3><p><i class="icofont-location-pin"></i> {{ event.location }}</p><router-link :to="`/events/${event.slug}`">Event details →</router-link></div></article><p v-if="!events.length" class="empty-state">Upcoming events will appear here.</p></div></div>
      </section>

      <section v-else-if="section.id === 'facilities'" id="section-facilities" class="home-section bg-white">
        <div class="home-shell"><div class="section-heading row-heading"><div><span class="section-kicker">Venues</span><h2>Sports Facilities</h2></div><router-link to="/facilities">Explore facilities →</router-link></div><div class="card-grid card-grid-3"><router-link v-for="item in facilities.slice(0,3)" :key="item.id" :to="`/facilities/${item.slug}`" class="image-card"><img v-if="item.image_url" :src="mediaUrl(item.image_url)" :alt="item.name" /><div v-else class="card-placeholder"><i class="icofont-stadium"></i></div><div><h3>{{ item.name }}</h3><p>{{ item.description }}</p></div></router-link></div></div>
      </section>

      <section v-else-if="section.id === 'associations'" id="section-associations" class="home-section associations-strip">
        <div class="home-shell"><div class="section-heading"><span class="section-kicker">Recognised bodies</span><h2>National Sports Associations</h2></div><div class="logo-strip"><div v-for="item in associations.slice(0,10)" :key="item.id"><img v-if="item.logo_url" :src="mediaUrl(item.logo_url)" :alt="item.name" /><i v-else class="icofont-trophy"></i><span>{{ item.name }}</span></div></div></div>
      </section>

      <section v-else-if="section.id === 'help'" id="section-help" class="home-section bg-white">
        <div class="home-shell"><div class="section-heading"><span class="section-kicker">Help & Support</span><h2>How can we help?</h2><p>Find guidance, contact NCS, or access your services.</p></div><div class="help-grid"><router-link to="/my-portal"><i class="icofont-ui-user"></i><h3>My Portal</h3><p>Sign in to manage registrations, licensing and renewals.</p></router-link><router-link to="/resource-centre"><i class="icofont-download"></i><h3>Resource Centre</h3><p>Guidelines, forms, reports and rules.</p></router-link><router-link to="/contact-us"><i class="icofont-support"></i><h3>Contact Support</h3><p>Speak with the NCS team.</p></router-link></div></div>
      </section>

      <section v-else-if="section.id === 'cta'" id="section-cta" class="home-cta"><div class="home-shell"><div><span>National Council of Sports Uganda</span><h2>Manage your NCS services in one place.</h2></div><router-link to="/my-portal" class="home-btn home-btn-gold">Open My Portal</router-link></div></section>

      <section v-else-if="section.id === 'faq_facts'" id="section-faq" class="home-section bg-soft">
        <div class="home-shell faq-facts-grid"><div><span class="section-kicker">{{ home.faq_eyebrow }}</span><h2>{{ home.faq_title }}</h2><div class="faq-list"><details v-for="faq in faqs.slice(0,10)" :key="faq.id"><summary>{{ faq.question }}</summary><p>{{ faq.answer }}</p></details></div><router-link to="/faqs" class="text-link">View All FAQs →</router-link></div><div><span class="section-kicker">{{ home.facts_eyebrow }}</span><h2>{{ home.facts_title }}</h2><div class="fact-list"><div v-for="fact in displayedFacts.slice(0,10)" :key="fact.id || fact.label"><i :class="fact.icon || 'icofont-light-bulb'"></i><p><strong>{{ fact.value }}</strong><span>{{ fact.label }}</span></p></div></div></div></div>
      </section>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { getSettings, listAssociations, listEvents, listFacilities, listFAQs, listFunFacts, listPosts, listSlides } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const defaultSections = ['hero','about','stats','news','find_sport','get_involved','events','facilities','associations','help','cta','faq_facts'].map(id => ({ id, visible:true }))
const home = reactive({
  sections: defaultSections,
  about:{eyebrow:'About NCS',title:'Developing Sports Excellence Since 1964',intro:'The National Council of Sports (NCS) is a statutory body established to develop, promote, and control sports in Uganda under the Ministry of Education and Sports.',body:'Established under the <strong>National Council of Sports Act (Chapter 48)</strong>, assented on 22 June 1964 and commenced on 25 June 1964, NCS serves as the apex regulator for sports development in Uganda, now updated by the <strong>National Sports Act, 2023</strong>.',leadership_label:'View Current Membership',leadership_url:'/team',core_title:'Core Functions of NCS',core_intro:'As mandated by the National Sports Act, NCS performs the following key functions:',mandate_label:'Read Full Mandate',mandate_url:'/pages/the-mandate'},
  leadership:{chairperson_name:'Mr. Ambrose Tashobya',chairperson_image:'',secretary_name:'Dr. Bernard Patrick Ogwel'},
  milestones:[{value:'60+',label:'Years of Excellence',icon:'icofont-award'},{value:'54+',label:'Sports Associations',icon:'icofont-trophy'},{value:'32+',label:'Sports Facilities',icon:'icofont-stadium'}],
  values:[{title:'Our Mission',text:'Maximizing opportunities for all Ugandans to participate and excel in Sports.',icon:'icofont-dart',featured:true},{title:'Our Vision',text:'A centre of excellence for promotion and development of Sports.',icon:'icofont-eye'},{title:'Integrity',text:'Upholding the highest standards of ethics and fair play in all sporting activities.',icon:'icofont-shield'},{title:'Inclusivity',text:'Ensuring sports opportunities are accessible to all Ugandans regardless of background.',icon:'icofont-people'},{title:'Excellence',text:'Striving for the highest standards in athlete development and sports administration.',icon:'icofont-award'},{title:'Global Recognition',text:'Positioning Uganda as a leading sports nation on the African and world stage.',icon:'icofont-globe',featured:true}],
  core_functions:['Register and regulate national sports organisations','Develop and promote sports throughout Uganda','Advise government on sports policy and standards','Coordinate national and international sports participation','Manage and develop public sports facilities','Support athlete, coach and official development'],
  stats_title:'Sports Excellence in Numbers',stats_intro:'Driving the development of sports across Uganda through dedicated programs and world-class facilities',finder_eyebrow:'Discover your federation',finder_title:'Find Your Sport',finder_intro:'Search across all 50+ National Sports Associations and Federations recognised by NCS. Tap any card to see the president, secretary, address, phone and website.',involved_title:'Get Involved',involved_subtitle:"Be Part of Uganda's Sports Excellence",involved_text:"Whether you're an athlete, coach, sports association, or enthusiast, the National Council of Sports welcomes you to join us in developing and promoting sports across Uganda.",register_label:'My Portal',register_url:'/my-portal',contact_label:'Contact Us',contact_url:'/contact-us',faq_eyebrow:'Got Questions?',faq_title:'Frequently Asked Questions',facts_eyebrow:'Did You Know?',facts_title:'Fun Facts'
})
const contact = reactive({ phone:'+256 414254477 / 343688', email:'info@ncs.go.ug', address:'Plot 2-10, Coronation Avenue', postal_address:'P.O. Box 20077, Lugogo, Kampala - UGANDA', fax:'+256 414 258350' })
const posts=ref([]), events=ref([]), slides=ref([]), facts=ref([]), faqs=ref([]), facilities=ref([]), associations=ref([])
const currentSlide=ref(0), sportQuery=ref(''), activeCategory=ref('All Sports')
let slideTimer
const visibleSections=computed(() => (home.sections?.length ? home.sections : defaultSections).filter(item=>item.visible!==false))
const displaySlides=computed(() => slides.value.length ? slides.value : [{id:'default',title:'Welcome to National Council of Sports',subtitle:'National Council of Sports',description:'A Centre of Excellence for Promotion and Development of Sports',button_text:'Discover NCS',button_url:'/pages/the-mandate',image_url:''}])
const activeSlide=computed(()=>displaySlides.value[currentSlide.value]||displaySlides.value[0])
const displayedFacts=computed(() => facts.value.length ? facts.value : home.milestones)
const sportCategories=computed(() => ['All Sports',...new Set(associations.value.map(item=>item.category).filter(Boolean))])
const filteredAssociations=computed(() => { const q=sportQuery.value.trim().toLowerCase(); return associations.value.filter(item => (activeCategory.value==='All Sports'||item.category===activeCategory.value) && (!q||[item.name,item.category,item.description].some(value=>String(value||'').toLowerCase().includes(q)))) })
function moveSlide(step){currentSlide.value=(currentSlide.value+step+displaySlides.value.length)%displaySlides.value.length;restartTimer()}
function setSlide(index){currentSlide.value=index;restartTimer()}
function restartTimer(){clearInterval(slideTimer);if(displaySlides.value.length>1)slideTimer=setInterval(()=>moveSlide(1),6500)}
function formatDate(value){return value?new Date(value).toLocaleDateString('en-UG',{day:'numeric',month:'short',year:'numeric'}):'Date TBA'}
function stripPlus(value){return String(value||'').replace(/\+\s*$/,'')}
function mergeHome(value){if(!value||typeof value!=='object')return;for(const [key,val] of Object.entries(value)){if(val&&typeof val==='object'&&!Array.isArray(val)&&home[key]&&typeof home[key]==='object')Object.assign(home[key],val);else home[key]=val}}
function dataOf(result){return result.status==='fulfilled'?result.value.data?.data:null}
onMounted(async()=>{const results=await Promise.allSettled([listPosts({status:'published',per_page:6}),listEvents({status:'published',per_page:6}),listSlides(),listFunFacts(),listFAQs(),listFacilities(),listAssociations(),getSettings('homepage'),getSettings('contact')]);posts.value=dataOf(results[0])?.items||[];events.value=dataOf(results[1])?.items||[];slides.value=dataOf(results[2])||[];facts.value=dataOf(results[3])||[];faqs.value=dataOf(results[4])||[];facilities.value=dataOf(results[5])||[];associations.value=dataOf(results[6])||[];mergeHome(dataOf(results[7])?.value);Object.assign(contact,dataOf(results[8])?.value||{});restartTimer()})
onUnmounted(()=>clearInterval(slideTimer))
</script>

<style scoped>
.home-redesign{color:#334155}.home-shell{width:min(80rem,100%);margin:auto;padding-left:1.25rem;padding-right:1.25rem}.home-section{padding:5rem 0}.bg-soft{background:#f7f9fc}.home-section h2,.home-cta h2{color:#1a365d;font-size:clamp(2rem,4vw,3rem);font-weight:800;line-height:1.12}.home-section p{line-height:1.75}.section-kicker,.home-pill{display:inline-block;margin-bottom:1rem;border-radius:999px;background:rgb(245 166 35/.13);padding:.4rem .85rem;color:#d88700;font-size:.8rem;font-weight:800}.home-hero{position:relative;height:clamp(31rem,65vw,43rem);overflow:hidden;background:#1a365d}.home-hero-slide{position:absolute;inset:0;opacity:0;transition:opacity .7s}.home-hero-slide.active{opacity:1;z-index:1}.home-hero-slide>img,.home-hero-fallback{width:100%;height:100%;object-fit:cover}.home-hero-fallback{background:linear-gradient(120deg,#1a365d,#274d7e)}.home-hero-overlay{position:absolute;inset:0;background:linear-gradient(90deg,rgb(15 31 61/.94),rgb(26 54 93/.68),transparent)}.home-hero-content{position:absolute;inset:0;display:flex;align-items:center;color:white}.home-hero-content>div{max-width:43rem}.home-hero h1{font-size:clamp(2.5rem,6vw,4.6rem);font-weight:850;line-height:1.05}.home-hero p{max-width:39rem;margin:1.25rem 0 2rem;font-size:clamp(1.05rem,2vw,1.3rem)}.home-actions{display:flex;flex-wrap:wrap;gap:.8rem;margin-top:1.6rem}.home-btn{display:inline-flex;min-height:44px;align-items:center;justify-content:center;border-radius:.55rem;padding:.7rem 1.2rem;font-size:.86rem;font-weight:800;transition:.2s}.home-btn-gold{background:#f5a623;color:#172b4d}.home-btn-gold:hover{background:#ffc154}.home-btn-navy{margin-top:1rem;background:#1a365d;color:white}.home-btn-outline-light{border:1px solid rgb(255 255 255/.75);color:white}.home-btn-outline-navy{border:1px solid #1a365d;color:#1a365d}.hero-arrow{position:absolute;z-index:3;top:50%;width:2.8rem;height:2.8rem;border-radius:50%;background:rgb(255 255 255/.16);color:white;font-size:2rem}.hero-prev{left:1rem}.hero-next{right:1rem}.hero-dots{position:absolute;z-index:3;bottom:1.5rem;left:50%;display:flex;gap:.5rem;transform:translateX(-50%)}.hero-dots button{width:.65rem;height:.65rem;border-radius:50%;background:rgb(255 255 255/.5)}.hero-dots button.active{width:2rem;border-radius:1rem;background:#f5a623}.about-grid{display:grid;grid-template-columns:1fr 1fr;gap:4rem}.home-lead{margin:1.2rem 0;font-size:1.08rem}.milestone-grid{display:flex;flex-wrap:wrap;gap:1.5rem;margin-top:2rem}.milestone{display:grid;grid-template-columns:auto auto;align-items:center;gap:.25rem .5rem}.milestone i{color:#f5a623;font-size:1.4rem}.milestone strong{color:#1a365d;font-size:2rem}.milestone span{grid-column:1/-1;font-size:.76rem}.value-grid{display:grid;grid-template-columns:1fr 1fr;gap:1rem}.value-card{border:1px solid #e8edf3;border-radius:1rem;background:#f8fafc;padding:1.25rem}.value-card.featured{background:#1a365d;color:white}.value-card h3{margin:.6rem 0;color:#1a365d;font-weight:800}.value-card.featured h3{color:white}.value-card p{font-size:.82rem}.value-icon{display:flex;width:2.8rem;height:2.8rem;align-items:center;justify-content:center;border-radius:.7rem;background:rgb(26 54 93/.1);color:#1a365d;font-size:1.3rem}.featured .value-icon{background:rgb(245 166 35/.2);color:#f5a623}.core-functions{margin-top:4rem;border-top:1px solid #e8edf3;padding-top:4rem;text-align:center}.section-heading{max-width:45rem;margin:0 auto 2.5rem;text-align:center}.section-heading p{margin-top:.7rem}.row-heading{display:flex;max-width:none;align-items:end;justify-content:space-between;text-align:left}.row-heading a,.text-link{color:#d88700;font-weight:800}.core-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:.8rem;margin-bottom:2rem;text-align:left}.core-grid div{display:flex;align-items:center;gap:.7rem;border:1px solid #e5eaf0;border-radius:.7rem;padding:1rem;font-weight:650}.core-grid span{display:flex;width:2rem;height:2rem;flex:none;align-items:center;justify-content:center;border-radius:50%;background:#f5a623;color:#1a365d}.stats-section,.home-cta{background:#1a365d;color:white}.section-heading.light h2{color:white}.section-heading.light p{color:rgb(255 255 255/.7)}.stats-grid{display:grid;grid-template-columns:repeat(4,1fr);gap:1rem}.stats-grid div{text-align:center}.stats-grid i{display:block;color:#f5a623;font-size:1.8rem}.stats-grid strong{display:block;margin:.45rem 0;color:white;font-size:2.5rem}.stats-grid span{color:rgb(255 255 255/.72)}.card-grid{display:grid;gap:1.25rem}.card-grid-3{grid-template-columns:repeat(3,1fr)}.content-card,.image-card{overflow:hidden;border:1px solid #e5eaf0;border-radius:1rem;background:white;box-shadow:0 12px 32px rgb(15 31 61/.07)}.content-card>img,.content-card>.card-placeholder,.image-card>img,.image-card>.card-placeholder{width:100%;height:13rem;object-fit:cover}.card-placeholder{display:flex;align-items:center;justify-content:center;background:#e8eef5;color:#9aa9ba;font-size:3rem}.content-card-body,.image-card>div:last-child{padding:1.25rem}.content-card-body>span{color:#d88700;font-size:.72rem;font-weight:800;text-transform:uppercase}.content-card h3,.image-card h3{margin:.5rem 0;color:#1a365d;font-size:1.15rem;font-weight:800}.content-card p,.image-card p{display:-webkit-box;overflow:hidden;color:#64748b;font-size:.85rem;-webkit-box-orient:vertical;-webkit-line-clamp:2}.finder-section{position:relative;overflow:hidden;background:#0f1f3d;color:white}.finder-search{position:relative;max-width:42rem;margin:0 auto 1rem}.finder-search i{position:absolute;top:50%;left:1.1rem;transform:translateY(-50%);color:#94a3b8}.finder-search input{width:100%;border-radius:1rem;background:white;padding:1rem 1.2rem 1rem 3rem;color:#1a365d}.filter-chips{display:flex;flex-wrap:wrap;justify-content:center;gap:.5rem;margin:1rem 0 2rem}.filter-chips button{border-radius:999px;background:rgb(255 255 255/.1);padding:.45rem .9rem;font-size:.78rem}.filter-chips button.active{background:#f5a623;color:#1a365d}.finder-meta{display:flex;justify-content:space-between;margin-bottom:1rem;color:rgb(255 255 255/.72);font-size:.85rem}.finder-meta a{color:#f5a623}.association-grid{display:grid;grid-template-columns:repeat(4,1fr);gap:1rem}.association-card{display:flex;gap:1rem;border:1px solid rgb(255 255 255/.12);border-radius:1rem;background:rgb(255 255 255/.07);padding:1rem}.association-card img,.association-logo{width:3.2rem;height:3.2rem;flex:none;object-fit:contain;border-radius:.7rem;background:white;padding:.25rem}.association-logo{display:flex;align-items:center;justify-content:center;color:#1a365d}.association-card>div:last-child{min-width:0}.association-card span{color:#f5a623;font-size:.68rem}.association-card h3{margin:.15rem 0;color:white;font-size:.9rem;font-weight:800}.association-card p{display:-webkit-box;overflow:hidden;color:rgb(255 255 255/.62);font-size:.72rem;-webkit-box-orient:vertical;-webkit-line-clamp:2}.association-card details{margin-top:.6rem;font-size:.72rem}.association-card summary{cursor:pointer;color:#f5a623}.association-card ul{margin:.5rem 0}.association-card a{color:#f5a623}.finder-empty{grid-column:1/-1;padding:3rem;text-align:center;background:rgb(255 255 255/.05);border-radius:1rem}.finder-empty i{font-size:2.5rem;color:rgb(255 255 255/.3)}.involved-grid{display:grid;grid-template-columns:1fr 1fr;align-items:center;gap:4rem}.quick-contact{border-radius:1rem;background:#f4f7fb;padding:2rem}.quick-contact h3{margin-bottom:1.4rem;color:#1a365d;font-size:1.4rem;font-weight:800}.quick-contact>div{display:flex;gap:1rem;margin-top:1rem}.quick-contact i{color:#f5a623;font-size:1.4rem}.quick-contact p{display:flex;flex-direction:column}.quick-contact span,.quick-contact a{color:#64748b;font-size:.85rem}.event-card{display:flex;overflow:hidden;border-radius:1rem;background:white;box-shadow:0 8px 24px rgb(15 31 61/.07)}.event-card time{display:flex;width:7rem;flex:none;align-items:center;justify-content:center;background:#1a365d;padding:1rem;color:white;text-align:center;font-weight:800}.event-card>div{padding:1.2rem}.event-card h3{color:#1a365d;font-weight:800}.event-card p{color:#64748b;font-size:.82rem}.event-card a{color:#d88700;font-size:.8rem;font-weight:800}.associations-strip{background:#f5f7fa}.logo-strip{display:grid;grid-template-columns:repeat(5,1fr);gap:1rem}.logo-strip>div{display:flex;min-height:8rem;flex-direction:column;align-items:center;justify-content:center;gap:.6rem;border-radius:.8rem;background:white;padding:1rem;text-align:center;box-shadow:0 6px 20px rgb(15 31 61/.05)}.logo-strip img{width:3.5rem;height:3.5rem;object-fit:contain}.logo-strip i{font-size:2rem;color:#1a365d}.logo-strip span{font-size:.72rem;font-weight:700}.help-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:1rem}.help-grid a{border:1px solid #e5eaf0;border-radius:1rem;padding:1.5rem;transition:.2s}.help-grid a:hover{transform:translateY(-4px);box-shadow:0 12px 28px rgb(15 31 61/.08)}.help-grid i{color:#f5a623;font-size:2rem}.help-grid h3{margin:.8rem 0;color:#1a365d;font-weight:800}.help-grid p{color:#64748b;font-size:.85rem}.home-cta{padding:3rem 0}.home-cta .home-shell{display:flex;align-items:center;justify-content:space-between;gap:2rem}.home-cta h2{margin-top:.3rem;color:white;font-size:2rem}.home-cta span{color:#f5a623;font-size:.8rem;font-weight:800}.faq-facts-grid{display:grid;grid-template-columns:1fr 1fr;gap:4rem}.faq-list details{border-bottom:1px solid #dce3eb;padding:1rem 0}.faq-list summary{cursor:pointer;color:#1a365d;font-weight:750}.faq-list p{margin-top:.7rem;color:#64748b;font-size:.86rem}.fact-list{display:grid;grid-template-columns:1fr 1fr;gap:.8rem}.fact-list>div{display:flex;align-items:center;gap:.8rem;border-radius:.8rem;background:white;padding:1rem}.fact-list i{color:#f5a623;font-size:1.5rem}.fact-list p{display:flex;flex-direction:column}.fact-list strong{color:#1a365d;font-size:1.35rem}.fact-list span{font-size:.75rem}.empty-state{grid-column:1/-1;padding:2rem;text-align:center;color:#64748b}
@media(max-width:1024px){.about-grid,.involved-grid,.faq-facts-grid{grid-template-columns:1fr}.association-grid{grid-template-columns:repeat(2,1fr)}.logo-strip{grid-template-columns:repeat(3,1fr)}}
@media(max-width:720px){.home-section{padding:3.5rem 0}.home-hero-overlay{background:rgb(15 31 61/.82)}.hero-arrow{display:none}.about-grid{gap:2.5rem}.value-grid,.card-grid-3,.stats-grid,.help-grid,.core-grid{grid-template-columns:1fr}.stats-grid{grid-template-columns:1fr 1fr}.association-grid{grid-template-columns:1fr}.logo-strip{grid-template-columns:1fr 1fr}.row-heading,.home-cta .home-shell{align-items:flex-start;flex-direction:column}.milestone-grid{display:grid;grid-template-columns:1fr 1fr}.fact-list{grid-template-columns:1fr}.finder-meta{gap:1rem}.event-card{flex-direction:column}.event-card time{width:100%}}
</style>
