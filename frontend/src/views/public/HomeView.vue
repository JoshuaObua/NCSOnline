<template>
  <div class="home-redesign">
    <template v-for="section in visibleSections" :key="section.id">
      <section v-if="section.id === 'hero'" id="section-hero">
        <PublicSlideshow :slideshow="slideshow" :slides="displaySlides" />
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
                  <img
                    v-if="leadershipLead?.image_url"
                    :src="mediaUrl(leadershipLead.image_url)"
                    :alt="leadershipLead.name"
                    class="w-20 h-20 rounded-full object-cover border-2 border-white shadow-md"
                    data-testid="homepage-chairperson-image"
                  />
                  <div
                    v-else
                    class="w-20 h-20 rounded-full bg-[#1a365d]/10 flex items-center justify-center text-[#1a365d] border-2 border-white shadow-md"
                    data-testid="homepage-chairperson-image"
                    aria-hidden="true"
                  >
                    <i class="icofont-user-alt-3 text-2xl"></i>
                  </div>
                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div v-for="(member, index) in displayedLeadershipMembers" :key="member.name || index">
                      <p class="text-sm text-gray-500">{{ member.title }}</p>
                      <p class="font-medium text-[#1a365d]">{{ member.name }}</p>
                    </div>
                  </div>
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
              <div v-for="item in home.core_functions" :key="item" class="flex items-start gap-3 p-4 bg-gray-50 rounded-lg hover:bg-[#f5a623]/5 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-[#f5a623] flex-shrink-0 mt-0.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21.801 10A10 10 0 1 1 17 3.335"/><path d="m9 11 3 3L22 4"/></svg>
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
            <div v-for="(fact, index) in displayedFacts" :key="fact.id || fact.label" class="text-center p-6 rounded-xl bg-white/5 backdrop-blur border border-white/10 hover:border-[#f5a623]/40 hover:bg-white/10 transition-all">
              <div class="w-14 h-14 rounded-xl bg-[#f5a623]/20 text-[#f5a623] flex items-center justify-center mx-auto mb-4">
                <i :class="[fact.icon || 'icofont-chart-growth', 'text-2xl']" aria-hidden="true"></i>
              </div>
              <div class="text-3xl md:text-4xl font-bold text-white mb-1">{{ counterValue(fact, index) }}</div>
              <div class="text-sm text-white/70">{{ fact.label }}</div>
              <p v-if="fact.description" class="mt-3 text-xs leading-relaxed text-white/55">{{ fact.description }}</p>
            </div>
          </div>
        </div>
      </section>

      <section v-else-if="section.id === 'news'" id="section-news" class="py-16 md:py-24 bg-gray-50">
        <div class="max-w-7xl mx-auto px-4">
          <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-8">
            <div>
              <span class="inline-block px-4 py-1.5 bg-[#f5a623]/10 text-[#f5a623] text-sm font-semibold rounded-full mb-3">Stay Updated</span>
              <h2 class="text-3xl md:text-4xl font-bold text-[#1a365d]">Latest News</h2>
            </div>
            <router-link
              to="/news"
              class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-[#1a365d] border shadow-sm h-9 px-4 py-2 border-[#1a365d] text-[#1a365d] hover:bg-[#1a365d] hover:text-white group"
              data-testid="homepage-view-all-news-button"
            >
              View All News
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" aria-hidden="true"><path d="m9 18 6-6-6-6"/></svg>
            </router-link>
          </div>

          <div class="flex flex-wrap gap-2 mb-8" data-testid="homepage-news-filters">
            <button
              v-for="category in homepageNewsCategories"
              :key="category.value"
              type="button"
              :data-testid="`homepage-news-filter-${category.value}`"
              :aria-pressed="activeNewsCategory === category.value"
              :class="activeNewsCategory === category.value ? 'bg-[#1a365d] text-white' : 'bg-white text-gray-600 hover:bg-[#f5a623]/10 hover:text-[#f5a623]'"
              class="rounded-full px-4 py-2 text-sm transition-all"
              @click="activeNewsCategory = category.value"
            >
              {{ category.label }}
            </button>
          </div>

          <div v-if="featuredNewsPost" :class="secondaryNewsPosts.length ? 'grid lg:grid-cols-2 gap-8' : 'grid gap-8 max-w-2xl'">
            <router-link class="group cursor-pointer" data-testid="homepage-featured-news-card" :to="`/news/${featuredNewsPost.slug}`">
              <div class="relative overflow-hidden rounded-2xl shadow-lg h-full bg-[#1a365d]">
                <div class="aspect-[4/3] overflow-hidden">
                  <img
                    v-if="featuredNewsPost.cover_image_url"
                    :src="mediaUrl(featuredNewsPost.cover_image_url)"
                    :alt="featuredNewsPost.title"
                    class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
                  />
                  <div v-else class="w-full h-full flex items-center justify-center bg-gradient-to-br from-[#1a365d] to-[#2d4a6f] text-white/35 text-5xl font-bold">NCS</div>
                </div>
                <div class="absolute inset-0 bg-gradient-to-t from-[#1a365d] via-[#1a365d]/50 to-transparent"></div>
                <div class="absolute bottom-0 left-0 right-0 p-6 md:p-8">
                  <div class="inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold transition-colors border-transparent shadow bg-[#f5a623] text-white mb-3">
                    {{ newsCategoryLabel(featuredNewsPost.category) }}
                  </div>
                  <h3 class="text-xl md:text-2xl font-bold text-white mb-3 group-hover:text-[#f5a623] transition-colors">{{ featuredNewsPost.title }}</h3>
                  <p class="text-white/80 mb-4 line-clamp-2">{{ featuredNewsPost.excerpt }}</p>
                  <div class="flex items-center justify-between gap-4">
                    <div class="flex items-center gap-2 text-white/70">
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M8 2v4"/><path d="M16 2v4"/><rect width="18" height="18" x="3" y="4" rx="2"/><path d="M3 10h18"/></svg>
                      <span class="text-sm">{{ formatDate(featuredNewsPost.published_at) }}</span>
                    </div>
                    <span class="hidden sm:flex items-center gap-2 text-[#f5a623] font-medium transition-all duration-300 -translate-x-2 opacity-0 group-hover:translate-x-0 group-hover:opacity-100">
                      Read More
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                    </span>
                  </div>
                </div>
              </div>
            </router-link>

            <div class="space-y-6">
              <router-link
                v-for="post in secondaryNewsPosts"
                :key="post.id"
                class="group flex gap-4 bg-white rounded-xl p-4 shadow-sm hover:shadow-lg transition-all duration-300 cursor-pointer border border-gray-100 hover:border-[#f5a623]/30"
                :data-testid="`homepage-news-card-${post.id}`"
                :to="`/news/${post.slug}`"
              >
                <div class="w-32 h-24 md:w-40 md:h-28 flex-shrink-0 overflow-hidden rounded-lg bg-[#1a365d]">
                  <img
                    v-if="post.cover_image_url"
                    :src="mediaUrl(post.cover_image_url)"
                    :alt="post.title"
                    class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                  />
                  <div v-else class="w-full h-full flex items-center justify-center text-white/35 text-xl font-bold">NCS</div>
                </div>
                <div class="flex-1 flex flex-col justify-between min-w-0">
                  <div>
                    <div class="inline-flex items-center rounded-md border px-2.5 py-0.5 font-semibold transition-colors border-transparent bg-gray-100 text-gray-600 mb-2 text-xs">
                      {{ newsCategoryLabel(post.category) }}
                    </div>
                    <h4 class="font-semibold text-[#1a365d] group-hover:text-[#f5a623] transition-colors line-clamp-2">{{ post.title }}</h4>
                  </div>
                  <div class="flex items-center justify-between gap-3 mt-2">
                    <div class="flex items-center gap-2 text-gray-500">
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3 h-3" aria-hidden="true"><path d="M8 2v4"/><path d="M16 2v4"/><rect width="18" height="18" x="3" y="4" rx="2"/><path d="M3 10h18"/></svg>
                      <span class="text-xs">{{ formatDate(post.published_at) }}</span>
                    </div>
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-[#f5a623] transition-all duration-300 -translate-x-2 opacity-0 group-hover:translate-x-0 group-hover:opacity-100" aria-hidden="true"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                  </div>
                </div>
              </router-link>
            </div>
          </div>

          <p v-else class="rounded-xl bg-white p-8 text-center text-gray-500 shadow-sm">News will appear here when published.</p>
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
          <div><span class="section-kicker">{{ home.involved_title }}</span><h2>{{ home.involved_subtitle }}</h2><p class="home-lead">{{ home.involved_text }}</p><div class="home-actions"><router-link :to="home.contact_url" class="home-btn home-btn-gold">{{ home.contact_label }}</router-link></div></div>
          <aside class="quick-contact"><h3>Quick Contact</h3><div><i class="icofont-location-pin"></i><p><b>Visit Us</b><span>{{ contact.address }}</span><span>{{ contact.postal_address }}</span></p></div><div><i class="icofont-phone"></i><p><b>Call Us</b><span>{{ contact.phone }}</span><span v-if="contact.fax">Fax: {{ contact.fax }}</span></p></div><div><i class="icofont-email"></i><p><b>Email Us</b><a :href="`mailto:${contact.email}`">{{ contact.email }}</a></p></div></aside>
        </div>
      </section>

      <section v-else-if="section.id === 'events'" id="section-events" class="home-section bg-white">
        <div class="home-shell">
          <div class="section-heading row-heading">
            <div><span class="section-kicker">{{ home.events.eyebrow }}</span><h2>{{ home.events.title }}</h2><p>{{ home.events.intro }}</p></div>
            <router-link to="/events">All events -></router-link>
          </div>
          <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            <router-link v-for="(event, index) in displayEvents.slice(0,6)" :key="event.id" :to="`/events/${event.slug}`" class="group relative bg-gradient-to-br from-gray-50 to-white rounded-2xl p-6 border border-gray-100 hover:border-[#f5a623]/30 hover:shadow-xl transition-all duration-500">
              <div class="absolute top-0 right-0 w-20 h-20 opacity-5 rounded-bl-full" :class="eventAccent(index).bg"></div>
              <div class="flex items-start gap-4">
                <div class="w-12 h-12 rounded-xl flex items-center justify-center group-hover:scale-110 transition-transform duration-300" :class="[eventAccent(index).softBg, eventAccent(index).text]">
                  <i :class="[eventAccent(index).icon, 'text-xl']" aria-hidden="true"></i>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="inline-flex items-center rounded-md px-2.5 py-0.5 font-semibold bg-gray-100 text-gray-600 mb-2 text-xs">{{ event.category || 'NCS Event' }}</div>
                  <h3 class="font-semibold text-[#1a365d] group-hover:text-[#f5a623] transition-colors mb-3 line-clamp-2">{{ event.title }}</h3>
                  <div class="space-y-2">
                    <div class="flex items-center gap-2 text-gray-500 text-sm"><i class="icofont-calendar text-[#f5a623]"></i><span>{{ formatDate(event.event_date) }}</span></div>
                    <div class="flex items-center gap-2 text-gray-500 text-sm"><i class="icofont-location-pin text-[#f5a623]"></i><span>{{ event.location || 'Venue TBA' }}</span></div>
                  </div>
                </div>
              </div>
              <div class="absolute bottom-6 right-6 opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-x-2 group-hover:translate-x-0"><i class="icofont-rounded-right text-[#f5a623]"></i></div>
            </router-link>
          </div>
        </div>
      </section>

      <section v-else-if="section.id === 'cta'" id="section-cta" class="home-cta"><div class="home-shell"><div><span>National Council of Sports Uganda</span><h2>In Case You Need Instant Help</h2></div><div class="home-cta-actions"><button type="button" class="home-btn home-btn-outline-light" @click="liveChatOpen = true"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg> Chat With Us Live</button><router-link to="/contact-us" class="home-btn home-btn-gold">Contact NCS</router-link></div></div></section>

      <LiveChatWidget v-if="liveChatOpen" @close="liveChatOpen = false" />

      <section v-else-if="section.id === 'faq_facts'" id="section-faq" class="home-section bg-soft">
        <div class="home-shell faq-facts-grid"><div><span class="section-kicker">{{ home.faq_eyebrow }}</span><h2>{{ home.faq_title }}</h2><div class="space-y-4 mt-6">
            <div
              v-for="faq in displayFAQs.slice(0,10)"
              :key="faq.id"
              class="border border-gray-100 rounded-xl px-6 shadow-sm hover:shadow-md transition-shadow"
              :class="isFaqOpen(faq.id) ? 'border-[#f5a623]/30 shadow-lg' : ''"
            >
              <h3 class="flex">
                <button
                  type="button"
                  class="flex flex-1 items-center justify-between gap-4 text-sm text-left text-[#1a365d] font-semibold hover:text-[#f5a623] py-5 transition-all"
                  :class="isFaqOpen(faq.id) ? 'text-[#f5a623]' : ''"
                  :aria-expanded="isFaqOpen(faq.id)"
                  @click="toggleFaq(faq.id)"
                >
                  <span>{{ faq.question }}</span>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 shrink-0 text-gray-400 transition-transform duration-200" :class="{ 'rotate-180': isFaqOpen(faq.id) }" aria-hidden="true"><path d="m6 9 6 6 6-6"/></svg>
                </button>
              </h3>
              <div class="grid overflow-hidden text-sm transition-[grid-template-rows] duration-300 ease-in-out" :class="isFaqOpen(faq.id) ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]'">
                <div class="overflow-hidden">
                  <div class="pt-0 text-gray-600 pb-5" v-html="faq.answer"></div>
                </div>
              </div>
            </div>
            <p v-if="!displayFAQs.length" class="empty-state">FAQs will appear here when published.</p>
          </div><router-link to="/faqs" class="text-link mt-6 inline-block">View All FAQs →</router-link></div><div><span class="section-kicker">{{ home.facts_eyebrow }}</span><h2>{{ home.facts_title }}</h2><div class="relative bg-gradient-to-br from-[#1a365d] to-[#2d4a6f] rounded-2xl p-8 md:p-10 text-white overflow-hidden mt-2" @mouseenter="pauseFacts" @mouseleave="resumeFacts"><div class="absolute inset-0 opacity-10"><div class="absolute top-0 right-0 w-64 h-64 rounded-full bg-[#f5a623] blur-3xl"></div><div class="absolute bottom-0 left-0 w-48 h-48 rounded-full bg-[#f5a623] blur-3xl"></div></div><div class="relative z-10"><div class="flex items-center gap-3 mb-6"><div class="w-12 h-12 rounded-full bg-[#f5a623]/20 flex items-center justify-center"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6 text-[#f5a623]" aria-hidden="true"><path d="M15 14c.2-1 .7-1.7 1.5-2.5 1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5.7.7 1.3 1.5 1.5 2.5"></path><path d="M9 18h6"></path><path d="M10 22h4"></path></svg></div><span class="text-[#f5a623] font-semibold">Sports History</span></div><div class="min-h-[120px] relative"><p v-for="(fact, index) in cmsFacts" :key="fact.id || index" class="text-lg md:text-xl leading-relaxed transition-all duration-500" :class="index === factIndex ? 'opacity-100 translate-y-0 relative' : 'opacity-0 translate-y-4 absolute top-0 left-0 right-0'">&quot;{{ fact.value }}&quot;</p><p v-if="!cmsFacts.length" class="text-lg md:text-xl leading-relaxed opacity-100">Fun facts will appear here when published.</p></div><div class="flex gap-2 mt-8"><button v-for="(fact, index) in cmsFacts" :key="'dot-' + (fact.id || index)" type="button" class="rounded-full transition-all duration-300" :class="index === factIndex ? 'w-8 h-2 bg-[#f5a623]' : 'w-2 h-2 bg-white/30 hover:bg-white/50'" :aria-label="`Go to fact ${index + 1}`" @click="goToFact(index)"></button></div></div></div></div></div>
      </section>
    </template>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { getSettings, getSlideshow, listAssociations, listEvents, listFAQs, listFunFacts, listPosts } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'
import PublicSlideshow from '@/components/public/PublicSlideshow.vue'
import LiveChatWidget from '@/components/public/LiveChatWidget.vue'

const defaultSections = ['hero','about','stats','news','find_sport','get_involved','events','cta','faq_facts'].map(id => ({ id, visible:true }))
const home = reactive({
  sections: defaultSections,
  about:{eyebrow:'About NCS',title:'Developing Sports Excellence Since 1964',intro:'The National Council of Sports (NCS) is a statutory body established to develop, promote, and control sports in Uganda under the Ministry of Education and Sports.',body:'Established under the <strong>National Council of Sports Act (Chapter 48)</strong>, assented on 22 June 1964 and commenced on 25 June 1964, NCS serves as the apex regulator for sports development in Uganda, now updated by the <strong>National Sports Act, 2023</strong>.',leadership_label:'View Current Membership',leadership_url:'/team',core_title:'Core Functions of NCS',core_intro:'As mandated by the National Sports Act, NCS performs the following key functions:',mandate_label:'Read Full Mandate',mandate_url:'/pages/the-mandate'},
  leadership:{members:[{name:'Mr. Ambrose Tashobya',title:'Chairperson',image_url:''},{name:'Dr. Bernard Patrick Ogwel',title:'General Secretary',image_url:''}]},
  milestones:[{value:'60+',label:'Years of Excellence',description:'Established in 1964 and still driving national sport.',icon:'icofont-award'},{value:'54+',label:'Sports Associations',description:'Recognised bodies supported across Uganda.',icon:'icofont-trophy'},{value:'32+',label:'Sports Facilities',description:'Facilities and venues supporting athletes and federations.',icon:'icofont-building-alt'},{value:'100K+',label:'Athletes Reached',description:'Athletes, coaches, and administrators served through NCS programs.',icon:'icofont-users-alt-5'}],
  values:[{title:'Our Mission',text:'Maximizing opportunities for all Ugandans to participate and excel in Sports.',icon:'icofont-dart',featured:true},{title:'Our Vision',text:'A centre of excellence for promotion and development of Sports.',icon:'icofont-eye'},{title:'Integrity',text:'Upholding the highest standards of ethics and fair play in all sporting activities.',icon:'icofont-shield'},{title:'Inclusivity',text:'Ensuring sports opportunities are accessible to all Ugandans regardless of background.',icon:'icofont-people'},{title:'Excellence',text:'Striving for the highest standards in athlete development and sports administration.',icon:'icofont-award'},{title:'Global Recognition',text:'Positioning Uganda as a leading sports nation on the African and world stage.',icon:'icofont-globe',featured:true}],
  core_functions:['Developing, promoting, and controlling sports on a national basis, including training and staffing','Recognizing sports disciplines and registering national sports organizations','Regulating associations and federations, awarding medals, certificates, trophies, and incentives','Approving international and national competitions and festivals','Facilitating Ugandan athletes participation in international competitions','Encouraging cooperation among associations and stimulating interest at all levels','Sponsoring scholarships for coaches and organizers','Advising on external sports relations and promoting sportsmanship','Arranging facilities with local authorities'],
  events:{eyebrow:'Upcoming Events',title:'NCS Calendar',intro:'Follow national competitions, federation events, athlete development programs and major sports gatherings.'},
  stats_title:'Sports Excellence in Numbers',stats_intro:'Driving the development of sports across Uganda through dedicated programs and world-class facilities',finder_eyebrow:'Discover your federation',finder_title:'Find Your Sport',finder_intro:'Search across all 50+ National Sports Associations and Federations recognised by NCS. Tap any card to see the president, secretary, address, phone and website.',involved_title:'Get Involved',involved_subtitle:"Be Part of Uganda's Sports Excellence",involved_text:"Whether you're an athlete, coach, sports association, or enthusiast, the National Council of Sports welcomes you to join us in developing and promoting sports across Uganda.",contact_label:'Contact Us',contact_url:'/contact-us',faq_eyebrow:'Got Questions?',faq_title:'Frequently Asked Questions',facts_eyebrow:'Did You Know?',facts_title:'Fun Facts'
})
const contact = reactive({ phone:'+256 414254477 / 343688', email:'info@ncs.go.ug', address:'Plot 2-10, Coronation Avenue', postal_address:'P.O. Box 20077, Lugogo, Kampala - UGANDA', fax:'+256 414 258350' })
const slideshow = reactive({ name:'Homepage Hero', slug:'homepage-hero', transition_effect:'fade', transition_duration:700, autoplay_speed:6500, pause_on_hover:true })
const posts=ref([]), events=ref([]), slides=ref([]), facts=ref([]), faqs=ref([]), associations=ref([])
const sportQuery=ref(''), activeCategory=ref('All Sports')
const activeNewsCategory=ref('all')
const liveChatOpen=ref(false)
const defaultNewsFilters = ['General', 'Infrastructure', 'Events', 'International', 'Football']
const fallbackEvents = [
  { id:'event-1', title:'NCS Hosts CAA Heroes Luncheon', slug:'ncs-hosts-caa-heroes-luncheon', category:'Athletics', location:'Kampala', event_date:'2026-02-28T09:00:00Z' },
  { id:'event-2', title:'Uganda Volleyball Federation U20 Africa Championships', slug:'uganda-volleyball-federation-u20-africa-championships', category:'Volleyball', location:'Cameroon', event_date:'2026-02-05T09:00:00Z' },
  { id:'event-3', title:'Ugandan Cyclists to Take on the World On Rwandan Soil', slug:'ugandan-cyclists-to-take-on-the-world-on-rwandan-soil', category:'Cycling', location:'Rwanda', event_date:'2026-01-15T09:00:00Z' },
  { id:'event-4', title:'Cricket Cranes Target High-Performance Gains in South Africa', slug:'cricket-cranes-target-high-performance-gains-in-south-africa', category:'Cricket', location:'South Africa', event_date:'2026-01-20T09:00:00Z' },
]
const eventAccents = [
  { bg:'bg-blue-500', softBg:'bg-blue-500/10', text:'text-blue-500', icon:'icofont-trophy' },
  { bg:'bg-purple-500', softBg:'bg-purple-500/10', text:'text-purple-500', icon:'icofont-megaphone' },
  { bg:'bg-green-500', softBg:'bg-green-500/10', text:'text-green-500', icon:'icofont-bicycle' },
  { bg:'bg-red-500', softBg:'bg-red-500/10', text:'text-red-500', icon:'icofont-ui-press' },
]
const visibleSections=computed(() => (home.sections?.length ? home.sections : defaultSections).filter(item=>item.visible!==false))
const displaySlides=computed(() => slides.value.length ? slides.value : [{id:'default',title:'Welcome to National Council of Sports',subtitle:'National Council of Sports',description:'A Centre of Excellence for Promotion and Development of Sports',button_text:'Discover NCS',button_url:'/pages/the-mandate',image_url:''}])
const animatedCounterValues=ref([])
const factIndex=ref(0)
let factTimer=null
const openFaqIds=ref(new Set())
function isFaqOpen(id){return openFaqIds.value.has(id)}
function toggleFaq(id){const next=new Set(openFaqIds.value);next.has(id)?next.delete(id):next.add(id);openFaqIds.value=next}
const displayedFacts=computed(() => Array.isArray(home.milestones) && home.milestones.length ? home.milestones : facts.value)
const displayedLeadershipMembers=computed(() => Array.isArray(home.leadership?.members) ? home.leadership.members.slice(0, 2) : [])
const leadershipLead=computed(() => displayedLeadershipMembers.value[0] || null)
const cmsFacts=computed(() => facts.value)
const displayFAQs=computed(() => faqs.value)
const displayEvents=computed(() => events.value.length ? events.value : fallbackEvents)
const homepageNewsCategories=computed(() => {
  const labels = new Map([['all', 'All']])
  defaultNewsFilters.forEach(label => labels.set(newsCategoryKey(label), label))
  posts.value.forEach(post => {
    const key = newsCategoryKey(post.category)
    if (key && !labels.has(key)) labels.set(key, newsCategoryLabel(post.category))
  })
  return Array.from(labels, ([value, label]) => ({ value, label }))
})
const filteredHomePosts=computed(() => activeNewsCategory.value === 'all'
  ? posts.value
  : posts.value.filter(post => newsCategoryKey(post.category) === activeNewsCategory.value))
const featuredNewsPost=computed(() => filteredHomePosts.value[0] || null)
const secondaryNewsPosts=computed(() => filteredHomePosts.value.slice(1, 4))
const sportCategories=computed(() => ['All Sports',...new Set(associations.value.map(item=>item.category).filter(Boolean))])
const filteredAssociations=computed(() => { const q=sportQuery.value.trim().toLowerCase(); return associations.value.filter(item => (activeCategory.value==='All Sports'||item.category===activeCategory.value) && (!q||[item.name,item.category,item.description].some(value=>String(value||'').toLowerCase().includes(q)))) })
function formatDate(value){return value?new Date(value).toLocaleDateString('en-UG',{day:'numeric',month:'short',year:'numeric'}):'Date TBA'}
function newsCategoryKey(value){return String(value || 'general').trim().toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '') || 'general'}
function newsCategoryLabel(value){return String(value || 'General').replace(/[-_]/g, ' ').replace(/\s+/g, ' ').trim().replace(/(^|\s)\S/g, part => part.toUpperCase())}
function parseCounter(value){const raw=String(value||'0');const match=raw.match(/^([\d,.]+)\s*(.*)$/);return {number:Number((match?.[1]||'0').replace(/,/g,''))||0,suffix:match?.[2]||''}}
function animateCounters(){const targets=displayedFacts.value.map(item=>parseCounter(item.value));const started=performance.now();const duration=1100;function tick(now){const progress=Math.min((now-started)/duration,1);const eased=1-Math.pow(1-progress,3);animatedCounterValues.value=targets.map(item=>Math.round(item.number*eased));if(progress<1)requestAnimationFrame(tick)};requestAnimationFrame(tick)}
function counterValue(fact,index){const parsed=parseCounter(fact.value);const value=animatedCounterValues.value[index];return `${value ?? parsed.number}${parsed.suffix}`}
function eventAccent(index){return eventAccents[index % eventAccents.length]}
function stopFactRotation(){if(factTimer){clearInterval(factTimer);factTimer=null}}
function startFactRotation(){stopFactRotation();if(cmsFacts.value.length<2)return;factTimer=setInterval(()=>{factIndex.value=(factIndex.value+1)%cmsFacts.value.length},5000)}
function pauseFacts(){stopFactRotation()}
function resumeFacts(){startFactRotation()}
function goToFact(index){factIndex.value=index;startFactRotation()}
function mergeHome(value){if(!value||typeof value!=='object')return;for(const [key,val] of Object.entries(value)){if(val&&typeof val==='object'&&!Array.isArray(val)&&home[key]&&typeof home[key]==='object')Object.assign(home[key],val);else home[key]=val}}
const NON_NEWS_CATEGORIES = ['page', 'case_study', 'project']
function dataOf(result){return result.status==='fulfilled'?result.value.data?.data:null}
function itemsOf(result){const value=dataOf(result);return Array.isArray(value)?value:(value?.items||[])}
watch(displayedFacts, animateCounters, { immediate:true, deep:true })
onMounted(async()=>{const results=await Promise.allSettled([listPosts({status:'published',per_page:6}),listEvents({status:'published',per_page:6}),getSlideshow('homepage-hero'),listFunFacts(),listFAQs(),listAssociations(),getSettings('homepage'),getSettings('contact')]);posts.value=(dataOf(results[0])?.items||[]).filter(p=>!NON_NEWS_CATEGORIES.includes(p.category));events.value=dataOf(results[1])?.items||[];const show=dataOf(results[2])||{};Object.assign(slideshow,show);slides.value=show.slides||[];facts.value=itemsOf(results[3]);faqs.value=itemsOf(results[4]);associations.value=itemsOf(results[5]);mergeHome(dataOf(results[6])?.value);Object.assign(contact,dataOf(results[7])?.value||{});startFactRotation();if(faqs.value.length)openFaqIds.value=new Set([faqs.value[0].id])})
onBeforeUnmount(stopFactRotation)
</script>

<style scoped>
.home-redesign{color:#334155}.home-shell{width:min(80rem,100%);margin:auto;padding-left:1.25rem;padding-right:1.25rem}.home-section{padding:5rem 0}.bg-soft{background:#f7f9fc}.home-section h2,.home-cta h2{color:#1a365d;font-size:clamp(2rem,4vw,3rem);font-weight:800;line-height:1.12}.home-section p{line-height:1.75}.section-kicker,.home-pill{display:inline-block;margin-bottom:1rem;border-radius:999px;background:rgb(245 166 35/.13);padding:.4rem .85rem;color:#d88700;font-size:.8rem;font-weight:800}.home-hero{position:relative;height:clamp(31rem,65vw,43rem);overflow:hidden;background:#1a365d}.home-hero-slide{position:absolute;inset:0;opacity:0;transition:opacity .7s}.home-hero-slide.active{opacity:1;z-index:1}.home-hero-slide>img,.home-hero-fallback{width:100%;height:100%;object-fit:cover}.home-hero-fallback{background:linear-gradient(120deg,#1a365d,#274d7e)}.home-hero-overlay{position:absolute;inset:0;background:linear-gradient(90deg,rgb(15 31 61/.94),rgb(26 54 93/.68),transparent)}.home-hero-content{position:absolute;inset:0;display:flex;align-items:center;color:white}.home-hero-content>div{max-width:43rem}.home-hero h1{font-size:clamp(2.5rem,6vw,4.6rem);font-weight:850;line-height:1.05}.home-hero p{max-width:39rem;margin:1.25rem 0 2rem;font-size:clamp(1.05rem,2vw,1.3rem)}.home-actions{display:flex;flex-wrap:wrap;gap:.8rem;margin-top:1.6rem}.home-btn{display:inline-flex;gap:.5rem;min-height:44px;align-items:center;justify-content:center;border-radius:.55rem;padding:.7rem 1.2rem;font-size:.86rem;font-weight:800;transition:.2s}.home-btn-gold{background:#f5a623;color:#172b4d}.home-btn-gold:hover{background:#ffc154}.home-btn-navy{margin-top:1rem;background:#1a365d;color:white}.home-btn-outline-light{border:1px solid rgb(255 255 255/.75);color:white}.home-btn-outline-navy{border:1px solid #1a365d;color:#1a365d}.hero-arrow{position:absolute;z-index:3;top:50%;width:2.8rem;height:2.8rem;border-radius:50%;background:rgb(255 255 255/.16);color:white;font-size:2rem}.hero-prev{left:1rem}.hero-next{right:1rem}.hero-dots{position:absolute;z-index:3;bottom:1.5rem;left:50%;display:flex;gap:.5rem;transform:translateX(-50%)}.hero-dots button{width:.65rem;height:.65rem;border-radius:50%;background:rgb(255 255 255/.5)}.hero-dots button.active{width:2rem;border-radius:1rem;background:#f5a623}.about-grid{display:grid;grid-template-columns:1fr 1fr;gap:4rem}.home-lead{margin:1.2rem 0;font-size:1.08rem}.milestone-grid{display:flex;flex-wrap:wrap;gap:1.5rem;margin-top:2rem}.milestone{display:grid;grid-template-columns:auto auto;align-items:center;gap:.25rem .5rem}.milestone i{color:#f5a623;font-size:1.4rem}.milestone strong{color:#1a365d;font-size:2rem}.milestone span{grid-column:1/-1;font-size:.76rem}.value-grid{display:grid;grid-template-columns:1fr 1fr;gap:1rem}.value-card{border:1px solid #e8edf3;border-radius:1rem;background:#f8fafc;padding:1.25rem}.value-card.featured{background:#1a365d;color:white}.value-card h3{margin:.6rem 0;color:#1a365d;font-weight:800}.value-card.featured h3{color:white}.value-card p{font-size:.82rem}.value-icon{display:flex;width:2.8rem;height:2.8rem;align-items:center;justify-content:center;border-radius:.7rem;background:rgb(26 54 93/.1);color:#1a365d;font-size:1.3rem}.featured .value-icon{background:rgb(245 166 35/.2);color:#f5a623}.core-functions{margin-top:4rem;border-top:1px solid #e8edf3;padding-top:4rem;text-align:center}.section-heading{max-width:45rem;margin:0 auto 2.5rem;text-align:center}.section-heading p{margin-top:.7rem}.row-heading{display:flex;max-width:none;align-items:end;justify-content:space-between;text-align:left}.row-heading a,.text-link{color:#d88700;font-weight:800}.core-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:.8rem;margin-bottom:2rem;text-align:left}.core-grid div{display:flex;align-items:center;gap:.7rem;border:1px solid #e5eaf0;border-radius:.7rem;padding:1rem;font-weight:650}.core-grid span{display:flex;width:2rem;height:2rem;flex:none;align-items:center;justify-content:center;border-radius:50%;background:#f5a623;color:#1a365d}.stats-section,.home-cta{background:#1a365d;color:white}.section-heading.light h2{color:white}.section-heading.light p{color:rgb(255 255 255/.7)}.stats-grid{display:grid;grid-template-columns:repeat(4,1fr);gap:1rem}.stats-grid div{text-align:center}.stats-grid i{display:block;color:#f5a623;font-size:1.8rem}.stats-grid strong{display:block;margin:.45rem 0;color:white;font-size:2.5rem}.stats-grid span{color:rgb(255 255 255/.72)}.card-grid{display:grid;gap:1.25rem}.card-grid-3{grid-template-columns:repeat(3,1fr)}.content-card,.image-card{overflow:hidden;border:1px solid #e5eaf0;border-radius:1rem;background:white;box-shadow:0 12px 32px rgb(15 31 61/.07)}.content-card>img,.content-card>.card-placeholder,.image-card>img,.image-card>.card-placeholder{width:100%;height:13rem;object-fit:cover}.card-placeholder{display:flex;align-items:center;justify-content:center;background:#e8eef5;color:#9aa9ba;font-size:3rem}.content-card-body,.image-card>div:last-child{padding:1.25rem}.content-card-body>span{color:#d88700;font-size:.72rem;font-weight:800;text-transform:uppercase}.content-card h3,.image-card h3{margin:.5rem 0;color:#1a365d;font-size:1.15rem;font-weight:800}.content-card p,.image-card p{display:-webkit-box;overflow:hidden;color:#64748b;font-size:.85rem;-webkit-box-orient:vertical;-webkit-line-clamp:2}.finder-section{position:relative;overflow:hidden;background:#0f1f3d;color:white}.finder-search{position:relative;max-width:42rem;margin:0 auto 1rem}.finder-search i{position:absolute;top:50%;left:1.1rem;transform:translateY(-50%);color:#94a3b8}.finder-search input{width:100%;border-radius:1rem;background:white;padding:1rem 1.2rem 1rem 3rem;color:#1a365d}.filter-chips{display:flex;flex-wrap:wrap;justify-content:center;gap:.5rem;margin:1rem 0 2rem}.filter-chips button{border-radius:999px;background:rgb(255 255 255/.1);padding:.45rem .9rem;font-size:.78rem}.filter-chips button.active{background:#f5a623;color:#1a365d}.finder-meta{display:flex;justify-content:space-between;margin-bottom:1rem;color:rgb(255 255 255/.72);font-size:.85rem}.finder-meta a{color:#f5a623}.association-grid{display:grid;grid-template-columns:repeat(4,1fr);gap:1rem}.association-card{display:flex;gap:1rem;border:1px solid rgb(255 255 255/.12);border-radius:1rem;background:rgb(255 255 255/.07);padding:1rem}.association-card img,.association-logo{width:3.2rem;height:3.2rem;flex:none;object-fit:contain;border-radius:.7rem;background:white;padding:.25rem}.association-logo{display:flex;align-items:center;justify-content:center;color:#1a365d}.association-card>div:last-child{min-width:0}.association-card span{color:#f5a623;font-size:.68rem}.association-card h3{margin:.15rem 0;color:white;font-size:.9rem;font-weight:800}.association-card p{display:-webkit-box;overflow:hidden;color:rgb(255 255 255/.62);font-size:.72rem;-webkit-box-orient:vertical;-webkit-line-clamp:2}.association-card details{margin-top:.6rem;font-size:.72rem}.association-card summary{cursor:pointer;color:#f5a623}.association-card ul{margin:.5rem 0}.association-card a{color:#f5a623}.finder-empty{grid-column:1/-1;padding:3rem;text-align:center;background:rgb(255 255 255/.05);border-radius:1rem}.finder-empty i{font-size:2.5rem;color:rgb(255 255 255/.3)}.involved-grid{display:grid;grid-template-columns:1fr 1fr;align-items:center;gap:4rem}.quick-contact{border-radius:1rem;background:#f4f7fb;padding:2rem}.quick-contact h3{margin-bottom:1.4rem;color:#1a365d;font-size:1.4rem;font-weight:800}.quick-contact>div{display:flex;gap:1rem;margin-top:1rem}.quick-contact i{color:#f5a623;font-size:1.4rem}.quick-contact p{display:flex;flex-direction:column}.quick-contact span,.quick-contact a{color:#64748b;font-size:.85rem}.event-card{display:flex;overflow:hidden;border-radius:1rem;background:white;box-shadow:0 8px 24px rgb(15 31 61/.07)}.event-card time{display:flex;width:7rem;flex:none;align-items:center;justify-content:center;background:#1a365d;padding:1rem;color:white;text-align:center;font-weight:800}.event-card>div{padding:1.2rem}.event-card h3{color:#1a365d;font-weight:800}.event-card p{color:#64748b;font-size:.82rem}.event-card a{color:#d88700;font-size:.8rem;font-weight:800}.associations-strip{background:#f5f7fa}.logo-strip{display:grid;grid-template-columns:repeat(5,1fr);gap:1rem}.logo-strip>div{display:flex;min-height:8rem;flex-direction:column;align-items:center;justify-content:center;gap:.6rem;border-radius:.8rem;background:white;padding:1rem;text-align:center;box-shadow:0 6px 20px rgb(15 31 61/.05)}.logo-strip img{width:3.5rem;height:3.5rem;object-fit:contain}.logo-strip i{font-size:2rem;color:#1a365d}.logo-strip span{font-size:.72rem;font-weight:700}.help-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:1rem}.help-grid a{border:1px solid #e5eaf0;border-radius:1rem;padding:1.5rem;transition:.2s}.help-grid a:hover{transform:translateY(-4px);box-shadow:0 12px 28px rgb(15 31 61/.08)}.help-grid i{color:#f5a623;font-size:2rem}.help-grid h3{margin:.8rem 0;color:#1a365d;font-weight:800}.help-grid p{color:#64748b;font-size:.85rem}.home-cta{padding:3rem 0}.home-cta .home-shell{display:flex;align-items:center;justify-content:space-between;gap:2rem}.home-cta h2{margin-top:.3rem;color:white;font-size:2rem}.home-cta span{color:#f5a623;font-size:.8rem;font-weight:800}.home-cta-actions{display:flex;flex-wrap:wrap;gap:.75rem}.faq-facts-grid{display:grid;grid-template-columns:1fr 1fr;gap:4rem}.empty-state{grid-column:1/-1;padding:2rem;text-align:center;color:#64748b}
@media(max-width:1024px){.about-grid,.involved-grid,.faq-facts-grid{grid-template-columns:1fr}.association-grid{grid-template-columns:repeat(2,1fr)}.logo-strip{grid-template-columns:repeat(3,1fr)}}
@media(max-width:720px){.home-section{padding:3.5rem 0}.home-hero-overlay{background:rgb(15 31 61/.82)}.hero-arrow{display:none}.about-grid{gap:2.5rem}.value-grid,.card-grid-3,.stats-grid,.help-grid,.core-grid{grid-template-columns:1fr}.stats-grid{grid-template-columns:1fr 1fr}.association-grid{grid-template-columns:1fr}.logo-strip{grid-template-columns:1fr 1fr}.row-heading,.home-cta .home-shell{align-items:flex-start;flex-direction:column}.milestone-grid{display:grid;grid-template-columns:1fr 1fr}.finder-meta{gap:1rem}.event-card{flex-direction:column}.event-card time{width:100%}}
</style>
