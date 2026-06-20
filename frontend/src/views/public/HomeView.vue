<template>
  <div>

    <!-- ── HERO ────────────────────────────────────────────── -->
    <section class="bg-cream relative overflow-hidden">
      <!-- Background decorative circles -->
      <div class="absolute -top-24 -left-24 w-80 h-80 bg-sky-100 rounded-full opacity-50 pointer-events-none"></div>
      <div class="absolute top-10 right-0 w-56 h-56 bg-yellow-50 rounded-full opacity-60 pointer-events-none"></div>

      <div class="max-w-screen-xl px-8 mx-auto flex flex-col lg:flex-row items-center gap-10 lg:gap-6 py-14 lg:py-20">

        <!-- Left: text -->
        <div class="w-full lg:w-1/2 text-center lg:text-left">
          <div class="hero-badge inline-flex items-center gap-2 bg-primary-700/10 border border-primary-700/20 rounded-full px-4 py-1.5 text-sm text-primary-700 font-medium mb-6">
            <span class="w-2 h-2 rounded-full bg-primary-600 animate-pulse flex-shrink-0"></span>
            {{ displaySlides[currentSlide]?.subtitle || 'Official Sports Regulatory Authority — Uganda' }}
          </div>
          <h1 class="hero-title text-4xl sm:text-5xl font-bold text-darken mb-5 leading-tight">
            {{ currentTitle.main }}
            <span v-if="currentTitle.accent" class="text-accent"> {{ currentTitle.accent }}</span>
          </h1>
          <p class="hero-desc text-gray-600 text-lg mb-8 max-w-lg mx-auto lg:mx-0 leading-relaxed">
            {{ displaySlides[currentSlide]?.description || 'The National Council of Sports registers, licenses and regulates sports organisations, federations and clubs across Uganda.' }}
          </p>
          <div class="hero-btns flex flex-wrap gap-4 justify-center lg:justify-start">
            <router-link
              :to="displaySlides[currentSlide]?.button_url || '/apply'"
              class="bg-accent hover:bg-yellow-600 text-white font-bold px-8 py-3.5 rounded-full shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-300 text-sm"
            >
              {{ displaySlides[currentSlide]?.button_text || 'Apply for License' }}
            </router-link>
            <router-link to="/news" class="flex items-center gap-3 group">
              <div class="w-12 h-12 bg-white rounded-full shadow-md flex items-center justify-center group-hover:shadow-lg transition-shadow flex-shrink-0">
                <svg class="w-4 h-4 ml-0.5 text-accent" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M8 5v14l11-7z"/>
                </svg>
              </div>
              <span class="text-sm font-semibold text-gray-700 group-hover:text-darken transition-colors">Latest News</span>
            </router-link>
          </div>

          <!-- Slide dots -->
          <div v-if="displaySlides.length > 1" class="flex items-center gap-2 mt-8 justify-center lg:justify-start">
            <button
              v-for="(_, i) in displaySlides"
              :key="i"
              @click="goToSlide(i)"
              :aria-label="`Show slide ${i + 1} of ${displaySlides.length}`"
              :aria-current="currentSlide === i ? 'true' : undefined"
              class="transition-all rounded-full"
              :class="currentSlide === i ? 'w-6 h-2.5 bg-accent' : 'w-2.5 h-2.5 bg-gray-300 hover:bg-gray-400'"
            ></button>
          </div>
        </div>

        <!-- Right: image + floating elements -->
        <div class="w-full lg:w-1/2 relative mt-6 lg:mt-0">
          <!-- Decorative shapes -->
          <div class="absolute -top-5 -right-5 pub-floating-slow pointer-events-none z-0">
            <svg class="w-20 h-20 opacity-60" viewBox="0 0 79 79" fill="none">
              <path d="M35.26 2.24C37.6-.1 41.4-.1 43.74 2.24L76.76 35.26C79.1 37.6 79.1 41.4 76.76 43.74L43.74 76.76C41.4 79.1 37.6 79.1 35.26 76.76L2.24 43.74C-.1 41.4-.1 37.6 2.24 35.26L35.26 2.24Z" fill="#29B9E7"/>
            </svg>
          </div>
          <div class="absolute bottom-16 left-2 w-8 h-8 bg-accent/25 rounded-full animate-ping pointer-events-none z-0"></div>
          <div class="absolute top-1/2 -left-3 w-4 h-4 bg-primary-400/40 rounded-full animate-pulse pointer-events-none z-0"></div>

          <!-- Slide image -->
          <Transition name="slide-fade" mode="out-in">
            <div :key="currentSlide" class="relative z-10">
              <div
                v-if="displaySlides[currentSlide]?.image_url"
                class="w-full h-72 sm:h-80 lg:h-[420px] rounded-3xl overflow-hidden shadow-2xl"
              >
                <img
                  :src="mediaUrl(displaySlides[currentSlide].image_url)"
                  :alt="displaySlides[currentSlide].title"
                  class="w-full h-full object-cover"
                />
              </div>
              <div
                v-else
                class="w-full h-72 sm:h-80 lg:h-[420px] rounded-3xl overflow-hidden shadow-2xl bg-gradient-to-br from-primary-700 via-primary-800 to-primary-900 flex items-center justify-center relative"
              >
                <div class="absolute top-8 left-8 w-28 h-28 bg-white/5 rounded-full"></div>
                <div class="absolute bottom-8 right-8 w-36 h-36 bg-white/5 rounded-full"></div>
                <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-48 h-48 bg-white/5 rounded-full"></div>
                <img src="/main-logo.png" alt="NCS" class="w-32 h-32 object-contain opacity-75 z-10 relative" />
              </div>
            </div>
          </Transition>

          <!-- Floating badge: Federations -->
          <div class="absolute top-8 -left-6 pub-floating-slow z-20 hidden sm:block">
            <div class="bg-white rounded-2xl shadow-xl px-4 py-3 flex items-center gap-3">
              <div class="w-10 h-10 bg-primary-100 rounded-xl flex items-center justify-center flex-shrink-0">
                <i class="icofont-trophy text-primary-700 text-lg"></i>
              </div>
              <div>
                <div class="font-bold text-darken text-sm leading-none">50+</div>
                <div class="text-xs text-gray-500 mt-0.5">Federations</div>
              </div>
            </div>
          </div>

          <!-- Floating badge: Athletes -->
          <div class="absolute bottom-10 -right-5 pub-floating z-20 hidden sm:block">
            <div class="bg-white rounded-2xl shadow-xl px-4 py-3">
              <div class="font-bold text-darken text-sm">10,000+</div>
              <div class="text-xs text-gray-400 mt-0.5">Licensed Athletes</div>
              <div class="flex gap-1 mt-2">
                <div class="w-2 h-2 bg-green-400 rounded-full"></div>
                <div class="w-2 h-2 bg-green-300 rounded-full"></div>
                <div class="w-2 h-2 bg-green-200 rounded-full"></div>
              </div>
            </div>
          </div>

          <!-- Floating chart icon -->
          <div class="absolute top-8 right-8 pub-floating hidden md:flex pointer-events-none z-20">
            <div class="w-14 h-14 rounded-xl bg-red-400/20 backdrop-blur-sm flex items-center justify-center">
              <svg class="w-7 h-7 text-red-500/60" fill="currentColor" viewBox="0 0 24 24">
                <path d="M4 12h4v8H4v-8zm6-8h4v16h-4V4zm6 4h4v12h-4V8z"/>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Wave separator: cream → white -->
      <div class="text-white relative z-10 -mt-6 sm:-mt-12 pointer-events-none">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1200 120" preserveAspectRatio="none" class="w-full h-14 sm:h-20 lg:h-24 block">
          <path d="M600,112.77C268.63,112.77,0,65.52,0,7.23V120H1200V7.23C1200,65.52,931.37,112.77,600,112.77Z" fill="currentColor"/>
        </svg>
      </div>
    </section>

    <!-- ── FEATURE CARDS ────────────────────────────────────── -->
    <section class="bg-white pb-20 px-4">
      <div class="max-w-screen-xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center max-w-2xl mx-auto mb-16 mt-4">
          <span class="section-tag">What We Do</span>
          <h2 class="text-3xl font-bold text-darken mt-2 mb-3">How We Support <span class="text-accent">Sports in Uganda</span></h2>
          <p class="text-gray-500 text-sm leading-relaxed">A comprehensive sports management platform serving federations, clubs, coaches and athletes nationwide.</p>
        </div>
        <div class="grid md:grid-cols-3 gap-6 md:gap-5 mt-20">
          <div class="quick-link bg-white shadow-xl p-6 text-center rounded-2xl hover:shadow-2xl transition-all hover:-translate-y-1">
            <div class="bg-primary-700 rounded-full w-16 h-16 flex items-center justify-center mx-auto shadow-lg transform -translate-y-12 mb-2">
              <i class="icofont-certificate text-white text-2xl"></i>
            </div>
            <h3 class="font-semibold text-darken text-base mb-3">Register Sports Organisations</h3>
            <p class="text-gray-500 text-sm px-2 leading-relaxed">Register federations, associations, clubs and coaches to ensure legal compliance and recognition across Uganda.</p>
          </div>
          <div class="quick-link bg-white shadow-xl p-6 text-center rounded-2xl hover:shadow-2xl transition-all hover:-translate-y-1">
            <div class="rounded-full w-16 h-16 flex items-center justify-center mx-auto shadow-lg transform -translate-y-12 mb-2" style="background:#F48C06;">
              <i class="icofont-shield text-white text-2xl"></i>
            </div>
            <h3 class="font-semibold text-darken text-base mb-3">License &amp; Regulate</h3>
            <p class="text-gray-500 text-sm px-2 leading-relaxed">Issue licenses, enforce standards, and ensure all sports organisations operate within the law.</p>
          </div>
          <div class="quick-link bg-white shadow-xl p-6 text-center rounded-2xl hover:shadow-2xl transition-all hover:-translate-y-1">
            <div class="rounded-full w-16 h-16 flex items-center justify-center mx-auto shadow-lg transform -translate-y-12 mb-2" style="background:#29B9E7;">
              <i class="icofont-chart-growth text-white text-2xl"></i>
            </div>
            <h3 class="font-semibold text-darken text-base mb-3">Develop &amp; Promote</h3>
            <p class="text-gray-500 text-sm px-2 leading-relaxed">Invest in sports infrastructure, training programmes, and the promotion of excellence nationwide.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ── ABOUT NCS ─────────────────────────────────────────── -->
    <section class="bg-white py-4 px-4">
      <div class="max-w-screen-xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center max-w-3xl mx-auto mb-10">
          <span class="section-tag">About NCS</span>
          <h2 class="text-3xl font-bold text-darken mt-2 mb-4">What is <span class="text-accent">NCS?</span></h2>
          <p class="text-gray-500 about-text leading-relaxed">
            The National Council of Sports (NCS) is Uganda's national sports authority mandated by the National Council of Sports Act.
            We oversee the registration and regulation of all sports organisations, federations, associations, and clubs across the country —
            promoting sports development for the benefit of all Ugandans.
          </p>
        </div>
        <div class="flex flex-col md:flex-row justify-center gap-6">
          <!-- For Federations -->
          <div class="relative md:w-5/12 h-60 rounded-2xl overflow-hidden shadow-lg group cursor-pointer">
            <div class="absolute inset-0 bg-gradient-to-br from-primary-700 to-primary-900">
              <div class="absolute -top-8 -left-8 w-32 h-32 bg-white/5 rounded-full"></div>
              <div class="absolute -bottom-8 -right-8 w-40 h-40 bg-white/5 rounded-full"></div>
            </div>
            <div class="absolute inset-0 bg-primary-900/30 group-hover:bg-primary-900/10 transition-colors duration-300"></div>
            <div class="absolute inset-0 flex flex-col items-center justify-center z-10 text-center px-6">
              <h3 class="uppercase text-white font-bold text-lg mb-4 tracking-wider">FOR FEDERATIONS</h3>
              <router-link
                to="/apply"
                class="rounded-full text-white border border-white/60 text-sm px-7 py-2.5 font-semibold hover:bg-white hover:text-primary-700 transition-all duration-300 focus:outline-none"
              >Register Today</router-link>
            </div>
          </div>
          <!-- For Clubs & Athletes -->
          <div class="relative md:w-5/12 h-60 rounded-2xl overflow-hidden shadow-lg group cursor-pointer">
            <div class="absolute inset-0" style="background: linear-gradient(135deg, #F48C06 0%, #d47000 100%);">
              <div class="absolute -top-8 -left-8 w-32 h-32 bg-white/10 rounded-full"></div>
              <div class="absolute -bottom-8 -right-8 w-40 h-40 bg-white/10 rounded-full"></div>
            </div>
            <div class="absolute inset-0 bg-black/20 group-hover:bg-black/5 transition-colors duration-300"></div>
            <div class="absolute inset-0 flex flex-col items-center justify-center z-10 text-center px-6">
              <h3 class="uppercase text-white font-bold text-lg mb-4 tracking-wider">FOR CLUBS &amp; ATHLETES</h3>
              <router-link
                to="/apply"
                class="rounded-full text-white text-sm px-7 py-2.5 font-semibold hover:opacity-90 transition-all duration-300"
                style="background: rgba(17,43,78,0.85)"
              >Get Your License</router-link>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ── STATS ROW ─────────────────────────────────────────── -->
    <section class="bg-cream py-20 px-4">
      <div class="max-w-screen-xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-12">
          <span class="section-tag">By The Numbers</span>
          <h2 class="text-3xl font-bold text-darken mt-2">Sports in Uganda at a Glance</h2>
        </div>
        <div class="about-stats grid grid-cols-2 md:grid-cols-4 gap-8">
          <div v-for="(stat, i) in displayedStats" :key="stat.label || i" class="home-stat-item text-center">
            <div class="stat-num text-4xl md:text-5xl font-bold text-darken mb-2" :data-target="parseStatValue(stat.value)">{{ stat.value }}</div>
            <div class="w-8 h-1 bg-accent rounded-full mx-auto mb-2"></div>
            <div class="text-sm font-medium text-gray-500">{{ stat.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- ── ALTERNATING FEATURE ───────────────────────────────── -->
    <section class="bg-white py-20 px-4">
      <div class="max-w-screen-xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col md:flex-row items-center gap-12 md:gap-16">
          <!-- Text -->
          <div class="md:w-1/2 relative">
            <div class="w-14 h-14 bg-accent/20 rounded-full absolute -top-4 -left-4 animate-pulse pointer-events-none"></div>
            <h2 class="font-bold text-3xl text-darken relative z-10 mb-5 leading-tight">
              Everything you need for <span class="text-accent">sports governance</span> in Uganda
            </h2>
            <p class="text-gray-500 mb-5 leading-relaxed">NCS's sports management platform helps sports bodies manage registration, licensing, and compliance — all in one secure cloud-based system, accessible from anywhere.</p>
            <div class="space-y-3 mb-6">
              <div class="flex items-start gap-3">
                <div class="flex-shrink-0 bg-white shadow rounded-full p-2"><svg class="w-4 h-4 text-darken" fill="currentColor" viewBox="0 0 27 26"><rect width="11.8" height="11.8" rx="2" fill="#2F327D"/><rect y="14.2" width="11.8" height="11.8" rx="2" fill="#2F327D"/><rect x="14.8" width="11.8" height="11.8" rx="2" fill="#2F327D"/><rect x="14.8" y="14.2" width="11.8" height="11.8" rx="2" fill="#F48C06"/></svg></div>
                <p class="text-gray-500 text-sm">Administrators get a clear overview of all federations and clubs.</p>
              </div>
              <div class="flex items-start gap-3">
                <div class="flex-shrink-0 bg-white shadow rounded-full p-2"><svg class="w-4 h-4" fill="currentColor" viewBox="0 0 28 26"><rect x="8" y="6" width="20" height="20" rx="2" fill="#2F327D"/><rect width="21" height="21" rx="2" fill="#F48C06"/></svg></div>
                <p class="text-gray-500 text-sm">Applications and renewals handled digitally — no paperwork.</p>
              </div>
              <div class="flex items-start gap-3">
                <div class="flex-shrink-0 bg-white shadow rounded-full p-2"><i class="icofont-people text-darken text-base"></i></div>
                <p class="text-gray-500 text-sm">Coaches and athletes can track their own license status online.</p>
              </div>
            </div>
            <router-link to="/apply" class="inline-flex items-center gap-2 text-accent font-semibold text-sm hover:gap-3 transition-all">
              Learn More
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M17 8l4 4m0 0l-4 4m4-4H3"/></svg>
            </router-link>
          </div>
          <!-- Image -->
          <div class="md:w-1/2 relative">
            <div class="w-24 h-24 rounded-xl bg-sky-300/25 absolute -top-3 -left-3 pub-floating pointer-events-none"></div>
            <div class="relative z-10 rounded-2xl overflow-hidden shadow-2xl bg-gradient-to-br from-primary-700 to-primary-900 h-72 flex items-center justify-center">
              <div class="absolute top-6 left-6 w-16 h-16 bg-white/10 rounded-full"></div>
              <div class="absolute bottom-6 right-6 w-24 h-24 bg-white/10 rounded-full"></div>
              <div class="text-center text-white relative z-10 px-6">
                <div class="text-5xl font-bold mb-1">112</div>
                <div class="text-primary-200 text-sm mb-3">Districts Covered</div>
                <div class="w-10 h-0.5 bg-accent mx-auto mb-3"></div>
                <div class="text-primary-300 text-xs">Nationwide Sports Reach</div>
              </div>
            </div>
            <div class="w-28 h-28 bg-accent/20 rounded-xl pub-floating absolute -bottom-3 -right-3 pointer-events-none"></div>
            <button aria-label="Play video about NCS impact" class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-14 h-14 bg-white rounded-full shadow-xl flex items-center justify-center hover:scale-110 transition-transform z-20">
              <svg class="w-5 h-5 ml-0.5 text-accent" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path d="M8 5v14l11-7z"/></svg>
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- ── FUN FACTS / STATS STRIP ──────────────────────────── -->
    <section v-if="displayedFunFacts.length" ref="funFactsRef" class="py-16 px-4" style="background-color: #252641;">
      <div class="max-w-screen-xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-10">
          <span class="section-tag" style="color: #4a7fad;">Sports in Uganda</span>
          <h2 class="text-3xl font-bold text-white mt-2">Sports at a Glance</h2>
        </div>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-8">
          <div v-for="(fact, i) in displayedFunFacts" :key="fact.id || i" class="fun-fact-card text-center">
            <div v-if="fact.icon" class="text-accent text-4xl mb-3"><i :class="fact.icon"></i></div>
            <div v-else class="w-10 h-0.5 bg-accent mx-auto mb-4"></div>
            <div
              class="fun-fact-num text-4xl md:text-5xl font-bold text-white mb-2"
              :data-target="parseStatValue(fact.value)"
              :data-suffix="statSuffix(fact.value)"
            >{{ fact.value }}</div>
            <div class="text-sm font-medium text-gray-400 uppercase tracking-wide">{{ fact.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- ── UPCOMING EVENTS ──────────────────────────────────── -->
    <section ref="eventsRef" class="py-20 bg-white px-4">
      <div class="max-w-screen-xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-end mb-12">
          <div>
            <span class="section-tag">What's Happening</span>
            <h2 class="text-3xl font-bold text-darken mt-2">Upcoming Events</h2>
          </div>
          <router-link to="/events" class="text-sm font-semibold text-accent hover:text-yellow-600 flex items-center gap-1 transition-colors">
            All events
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M17 8l4 4m0 0l-4 4m4-4H3"/></svg>
          </router-link>
        </div>
        <div class="space-y-4">
          <router-link
            v-for="ev in events"
            :key="ev.id"
            :to="`/events/${ev.slug}`"
            class="event-card group flex items-center gap-5 p-5 rounded-2xl border border-gray-100 hover:border-primary-200 hover:shadow-md bg-white transition-all"
          >
            <div class="flex-shrink-0 w-16 h-16 bg-primary-700 text-white rounded-2xl flex flex-col items-center justify-center text-center shadow-md">
              <span class="text-[10px] font-semibold uppercase leading-none tracking-wide">{{ monthShort(ev.event_date) }}</span>
              <span class="text-2xl font-bold leading-tight">{{ new Date(ev.event_date).getDate() }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="font-semibold text-darken group-hover:text-primary-700 transition-colors truncate text-base">{{ ev.title }}</h3>
              <p v-if="ev.location" class="text-sm text-gray-400 mt-1 flex items-center gap-1">
                <i class="icofont-location-pin text-accent text-sm"></i>
                {{ ev.location }}
              </p>
            </div>
            <div class="flex-shrink-0 w-8 h-8 rounded-full bg-gray-50 group-hover:bg-primary-50 flex items-center justify-center transition-colors">
              <svg class="w-4 h-4 text-gray-400 group-hover:text-primary-600 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/></svg>
            </div>
          </router-link>
          <div v-if="!events.length" class="text-center py-12 text-gray-400">
            <i class="icofont-calendar text-5xl text-gray-200 block mb-3"></i>
            <p>No upcoming events.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ── LATEST NEWS ──────────────────────────────────────── -->
    <section ref="newsRef" class="py-20 bg-gray-50 px-4">
      <div class="max-w-screen-xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-12">
          <span class="section-tag">Stay Informed</span>
          <h2 class="text-3xl font-bold text-darken mt-2">Latest News &amp; Resources</h2>
          <p class="text-gray-500 mt-2 text-sm">Stay up to date with developments in Ugandan sports</p>
        </div>

        <div v-if="postsLoading" class="flex flex-col lg:flex-row gap-8">
          <div class="lg:w-6/12 h-80 bg-gray-200 rounded-2xl animate-pulse"></div>
          <div class="lg:w-6/12 space-y-5">
            <div v-for="i in 3" :key="i" class="flex gap-4">
              <div class="w-28 h-24 bg-gray-200 rounded-xl animate-pulse flex-shrink-0"></div>
              <div class="flex-1 space-y-2">
                <div class="h-4 bg-gray-200 rounded animate-pulse w-3/4"></div>
                <div class="h-3 bg-gray-200 rounded animate-pulse w-full"></div>
              </div>
            </div>
          </div>
        </div>

        <div v-else-if="posts.length" class="flex flex-col lg:flex-row lg:gap-16">
          <!-- Featured large article -->
          <div class="lg:w-6/12 mb-10 lg:mb-0">
            <router-link :to="`/news/${posts[0].slug}`" class="news-card group block">
              <div class="w-full h-60 sm:h-72 rounded-2xl overflow-hidden mb-5 shadow-sm">
                <img
                  v-if="posts[0].cover_image_url"
                  :src="mediaUrl(posts[0].cover_image_url)"
                  :alt="posts[0].title"
                  class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                />
                <div v-else class="w-full h-full bg-gradient-to-br from-primary-100 to-primary-200 flex items-center justify-center">
                  <i class="icofont-newspaper text-primary-300 text-5xl"></i>
                </div>
              </div>
              <span class="inline-block bg-yellow-300 text-darken font-semibold px-4 py-0.5 text-xs rounded-full mb-3">{{ posts[0].category || 'NEWS' }}</span>
              <h3 class="text-gray-800 font-bold text-xl mb-2 group-hover:text-primary-700 transition-colors leading-snug line-clamp-2">{{ posts[0].title }}</h3>
              <p v-if="posts[0].excerpt" class="text-gray-500 text-sm line-clamp-2 mb-3">{{ posts[0].excerpt }}</p>
              <span class="text-accent text-sm font-semibold">Read more →</span>
            </router-link>
          </div>

          <!-- Side articles -->
          <div class="lg:w-6/12 flex flex-col justify-between gap-5">
            <div v-for="post in posts.slice(1, 4)" :key="post.id" class="flex gap-4 news-card group">
              <div class="w-28 h-24 rounded-xl overflow-hidden flex-shrink-0 shadow-sm">
                <img
                  v-if="post.cover_image_url"
                  :src="mediaUrl(post.cover_image_url)"
                  :alt="post.title"
                  class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                />
                <div v-else class="w-full h-full bg-gradient-to-br from-primary-100 to-primary-200 flex items-center justify-center">
                  <i class="icofont-newspaper text-primary-300 text-2xl"></i>
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <span class="inline-block bg-yellow-300 text-darken font-semibold px-3 py-0.5 text-[10px] rounded-full mb-1.5">{{ post.category || 'NEWS' }}</span>
                <router-link :to="`/news/${post.slug}`">
                  <h4 class="text-gray-800 font-semibold text-sm sm:text-base group-hover:text-primary-700 transition-colors leading-snug line-clamp-2">{{ post.title }}</h4>
                </router-link>
                <p class="text-gray-400 mt-1.5 text-xs line-clamp-2">{{ post.excerpt }}</p>
              </div>
            </div>
            <router-link to="/news" class="self-end text-sm font-semibold text-accent hover:text-yellow-600 transition-colors flex items-center gap-1 mt-2">
              All News &amp; Resources →
            </router-link>
          </div>
        </div>

        <div v-else class="text-center py-16 text-gray-400">
          <i class="icofont-newspaper text-5xl text-gray-200 block mb-3"></i>
          <p>No news articles yet.</p>
        </div>
      </div>
    </section>

    <!-- ── CTA BANNER ───────────────────────────────────────── -->
    <section ref="ctaRef" class="relative overflow-hidden py-20 px-4" style="background: linear-gradient(135deg, #252641 0%, #2F327D 60%, #4B4F9E 100%);">
      <!-- Decorative -->
      <div class="absolute top-0 left-0 w-40 h-40 bg-white/5 rounded-full -translate-x-1/2 -translate-y-1/2 animate-pulse pointer-events-none"></div>
      <div class="absolute bottom-0 right-0 w-56 h-56 bg-white/5 rounded-full translate-x-1/3 translate-y-1/3 animate-pulse pointer-events-none"></div>
      <div class="absolute top-1/2 right-20 w-20 h-20 bg-accent/10 rounded-full animate-ping pointer-events-none"></div>

      <div class="max-w-3xl mx-auto text-center relative z-10">
        <span class="section-tag" style="color:#4a7fad;">Get Started</span>
        <h2 class="cta-title text-3xl md:text-4xl font-bold text-white mt-3 mb-4 leading-tight">
          Ready to Register Your <span class="text-accent">Sports Organisation?</span>
        </h2>
        <p class="cta-desc text-lg text-gray-300 mb-8 max-w-2xl mx-auto leading-relaxed">
          Apply online for your sports license, register your federation, club or association with the National Council of Sports Uganda.
        </p>
        <div class="flex flex-wrap gap-4 justify-center">
          <router-link
            to="/apply"
            class="bg-accent hover:bg-yellow-600 text-white font-bold px-10 py-3.5 rounded-full shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-300 text-sm"
          >Start Application</router-link>
          <router-link
            v-if="!isAuthenticated"
            to="/login"
            class="border border-white/40 hover:border-white text-white font-semibold px-8 py-3.5 rounded-full transition-all text-sm hover:bg-white/10"
          >Sign In</router-link>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { gsap } from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { listPosts, listEvents, listSlides, listFunFacts } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'
import { useAuthStore } from '@/stores/auth.js'

gsap.registerPlugin(ScrollTrigger)

const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)

const posts        = ref([])
const events       = ref([])
const slides       = ref([])
const funFacts     = ref([])
const postsLoading = ref(true)

// ── Default data ──────────────────────────────────────────────
const defaultStats = [
  { label: 'Sports Federations', value: '50+'    },
  { label: 'Registered Clubs',   value: '500+'   },
  { label: 'Licensed Athletes',  value: '10,000+'},
  { label: 'Districts Covered',  value: '112'    },
]

const displayedStats = computed(() =>
  funFacts.value.length > 0
    ? funFacts.value.slice(0, 4).map(f => ({ label: f.label, value: f.value }))
    : defaultStats
)

const displayedFunFacts = computed(() =>
  funFacts.value.length > 0
    ? funFacts.value
    : defaultStats.map((s, i) => ({ id: String(i), label: s.label, value: s.value, icon: '', is_active: true }))
)

function parseStatValue(v) {
  if (!v) return 0
  const n = parseInt(String(v).replace(/[^0-9]/g, ''), 10)
  return isNaN(n) ? 0 : n
}
function statSuffix(v) {
  const s = String(v || '')
  if (s.includes('K') || s.includes('k')) return 'K+'
  if (s.includes('+')) return '+'
  return ''
}

// ── Slideshow ─────────────────────────────────────────────────
const currentSlide = ref(0)
let slideTimer = null

const defaultSlide = {
  id: 'default',
  title: 'Building a World-Class Sports Nation',
  subtitle: 'Official Sports Regulatory Authority — Uganda',
  description: 'The National Council of Sports registers, licenses and regulates sports organisations, federations and clubs across Uganda. Apply online for your sports license today.',
  button_text: 'Apply for License',
  button_url: '/apply',
  image_url: '',
}
const displaySlides = computed(() => slides.value.length > 0 ? slides.value : [defaultSlide])

const currentTitle = computed(() => {
  const t = displaySlides.value[currentSlide.value]?.title || ''
  if (!t) return { main: 'Building a World-Class', accent: 'Sports Nation' }
  const words = t.split(' ')
  if (words.length <= 3) return { main: t, accent: '' }
  return { main: words.slice(0, -2).join(' '), accent: words.slice(-2).join(' ') }
})

function nextSlide() { currentSlide.value = (currentSlide.value + 1) % displaySlides.value.length; restartTimer() }
function prevSlide() { currentSlide.value = (currentSlide.value - 1 + displaySlides.value.length) % displaySlides.value.length; restartTimer() }
function goToSlide(i) { currentSlide.value = i; restartTimer() }
function restartTimer() {
  clearInterval(slideTimer)
  if (displaySlides.value.length > 1) slideTimer = setInterval(nextSlide, 6000)
}

// ── Date helpers ──────────────────────────────────────────────
function monthShort(d) {
  if (!d) return '?'
  return new Date(d).toLocaleDateString('en-UG', { month: 'short' })
}

// ── Section refs ──────────────────────────────────────────────
const funFactsRef = ref(null)
const newsRef     = ref(null)
const eventsRef   = ref(null)
const ctaRef      = ref(null)

// ── GSAP animations ───────────────────────────────────────────
function initAnimations() {
  const tl = gsap.timeline({ defaults: { ease: 'power3.out' } })
  tl.fromTo('.hero-badge', { opacity: 0, y: 30 }, { opacity: 1, y: 0, duration: 0.7 })
    .fromTo('.hero-title', { opacity: 0, y: 40 }, { opacity: 1, y: 0, duration: 0.8 }, '-=0.4')
    .fromTo('.hero-desc',  { opacity: 0, y: 30 }, { opacity: 1, y: 0, duration: 0.7 }, '-=0.5')
    .fromTo('.hero-btns',  { opacity: 0, y: 20 }, { opacity: 1, y: 0, duration: 0.6 }, '-=0.4')

  gsap.fromTo('.quick-link',
    { opacity: 0, y: 30 },
    { opacity: 1, y: 0, duration: 0.6, stagger: 0.15, ease: 'power2.out',
      scrollTrigger: { trigger: '.quick-link', start: 'top 85%' } }
  )

  gsap.fromTo('.about-text',
    { opacity: 0, y: 30 },
    { opacity: 1, y: 0, duration: 0.8, ease: 'power2.out',
      scrollTrigger: { trigger: '.about-text', start: 'top 80%' } }
  )

  gsap.fromTo('.home-stat-item',
    { opacity: 0, scale: 0.85 },
    { opacity: 1, scale: 1, duration: 0.6, stagger: 0.12, ease: 'back.out(1.5)',
      scrollTrigger: { trigger: '.about-stats', start: 'top 78%', onEnter: animateCounters } }
  )

  if (funFactsRef.value) {
    gsap.fromTo('.fun-fact-card',
      { opacity: 0, y: 30 },
      { opacity: 1, y: 0, duration: 0.6, stagger: 0.12, ease: 'power2.out',
        scrollTrigger: { trigger: funFactsRef.value, start: 'top 80%', onEnter: animateFunFactCounters } }
    )
  }

  gsap.fromTo('.news-card',
    { opacity: 0, y: 40 },
    { opacity: 1, y: 0, duration: 0.7, stagger: 0.12, ease: 'power2.out',
      scrollTrigger: { trigger: newsRef.value, start: 'top 78%' } }
  )

  gsap.fromTo('.event-card',
    { opacity: 0, x: -30 },
    { opacity: 1, x: 0, duration: 0.6, stagger: 0.1, ease: 'power2.out',
      scrollTrigger: { trigger: eventsRef.value, start: 'top 78%' } }
  )

  gsap.fromTo('.cta-title, .cta-desc',
    { opacity: 0, y: 30 },
    { opacity: 1, y: 0, duration: 0.8, stagger: 0.2, ease: 'power2.out',
      scrollTrigger: { trigger: ctaRef.value, start: 'top 80%' } }
  )
}

function animateCounters() {
  document.querySelectorAll('.stat-num').forEach(el => {
    const target = parseInt(el.dataset.target || '0', 10)
    if (!target) return
    const suffix = target >= 10000 ? 'K+' : target >= 100 ? '+' : ''
    const displayTarget = target >= 10000 ? target / 1000 : target
    gsap.fromTo({ val: 0 }, { val: displayTarget, duration: 1.8, ease: 'power2.out',
      onUpdate: function () { el.textContent = Math.floor(this.targets()[0].val) + suffix }
    })
  })
}

function animateFunFactCounters() {
  document.querySelectorAll('.fun-fact-num').forEach(el => {
    const target = parseInt(el.dataset.target || '0', 10)
    const suffix = el.dataset.suffix || ''
    if (!target) return
    const displayTarget = suffix.startsWith('K') ? target / 1000 : target
    gsap.fromTo({ val: 0 }, { val: displayTarget, duration: 2, ease: 'power2.out',
      onUpdate: function () { el.textContent = Math.floor(this.targets()[0].val) + suffix }
    })
  })
}

onMounted(async () => {
  try {
    const [pr, er, sr, fr] = await Promise.all([
      listPosts({ status: 'published', per_page: 4 }),
      listEvents({ status: 'published', per_page: 4 }),
      listSlides(),
      listFunFacts()
    ])
    posts.value      = pr.data.data?.items || []
    events.value     = er.data.data?.items || []
    slides.value     = sr.data.data || []
    funFacts.value   = fr.data.data || []
  } catch { /* graceful fallback */ } finally {
    postsLoading.value = false
  }
  restartTimer()
  setTimeout(initAnimations, 100)
})

onUnmounted(() => {
  clearInterval(slideTimer)
  ScrollTrigger.getAll().forEach(t => t.kill())
})
</script>

<style scoped>
.slide-fade-enter-active, .slide-fade-leave-active { transition: opacity 0.6s ease; }
.slide-fade-enter-from, .slide-fade-leave-to { opacity: 0; }
</style>
