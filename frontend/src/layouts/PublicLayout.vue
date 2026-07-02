<template>
  <div class="public-site min-h-screen flex flex-col">

    <a href="#main-content" class="skip-link">Skip to main content</a>
    <a href="#public-footer" class="skip-link skip-link-secondary">Skip to footer</a>

    <!-- ── Header (NCS Header spec) ─────────────────────────── -->
    <header :class="['fixed top-0 left-0 right-0 z-50 transition-all duration-500', scrolled ? 'bg-white shadow-lg' : 'bg-white']">

      <!-- Topbar (navy strip) -->
      <div class="bg-[#1a365d] text-white py-2 overflow-hidden">
        <div class="max-w-7xl mx-auto px-4 flex justify-between items-center">
          <div class="flex-1 overflow-hidden mr-4" aria-label="NCS highlights">
            <div class="public-marquee-container">
              <div class="public-marquee-content">
                <span class="inline-flex items-center gap-8 text-sm">
                  <template v-for="copy in 2" :key="copy">
                    <template v-for="(msg, idx) in headerSettings.marquee" :key="`${copy}-${idx}`">
                      <span>{{ msg }}</span>
                      <span class="text-[#f5a623]" aria-hidden="true">•</span>
                    </template>
                  </template>
                </span>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-3 flex-shrink-0">
            <a v-for="net in topbarSocialLinks" :key="net.key" :href="net.url" target="_blank" rel="noopener noreferrer" :aria-label="net.label" class="hover:text-[#f5a623] transition-colors duration-300">
              <svg v-if="net.key==='facebook'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"/></svg>
              <svg v-else-if="net.key==='twitter'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z"/></svg>
              <svg v-else-if="net.key==='linkedin'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"/><rect width="4" height="12" x="2" y="9"/><circle cx="4" cy="4" r="2"/></svg>
              <svg v-else-if="net.key==='youtube'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M2.5 17a24.12 24.12 0 0 1 0-10 2 2 0 0 1 1.4-1.4 49.56 49.56 0 0 1 16.2 0A2 2 0 0 1 21.5 7a24.12 24.12 0 0 1 0 10 2 2 0 0 1-1.4 1.4 49.55 49.55 0 0 1-16.2 0A2 2 0 0 1 2.5 17"/><path d="m10 15 5-3-5-3z"/></svg>
              <svg v-else-if="net.key==='webmail'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"/><rect x="2" y="4" width="20" height="16" rx="2"/></svg>
            </a>
          </div>
        </div>
      </div>

      <!-- Main bar -->
      <div class="max-w-7xl mx-auto px-4 py-3">
        <div class="flex items-center justify-between">

          <!-- Logo (bare image, no container) -->
          <router-link to="/" class="flex items-center gap-3 group flex-shrink-0" :aria-label="`${siteIdentity.name || 'NCS'} home`">
            <img :src="resolveAsset(siteIdentity.logoUrl) || '/main-logo.png'" :alt="`${siteIdentity.name || 'NCS'} - National Council of Sports`" class="h-14 md:h-16 w-auto object-contain" />
          </router-link>

          <!-- Desktop nav -->
          <nav class="hidden lg:flex items-center gap-1" aria-label="Primary navigation">
            <template v-for="item in visibleMenuItems" :key="item.id || item.label">
              <div v-if="(item.children?.length) || item.mega" class="relative group">
                <button type="button" :class="['flex items-center gap-1 px-3 py-2 font-medium transition-colors duration-300 text-sm', isActiveTopLevel(item) ? 'text-[#f5a623]' : 'text-[#1a365d] hover:text-[#f5a623]']" aria-haspopup="menu" aria-expanded="false">
                  {{ item.label }}
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 group-hover:rotate-180 transition-transform duration-200" aria-hidden="true"><path d="m6 9 6 6 6-6"/></svg>
                </button>
                <div class="absolute left-0 top-full pt-1 opacity-0 invisible group-hover:opacity-100 group-hover:visible group-focus-within:opacity-100 group-focus-within:visible transition-all duration-200 min-w-[220px] z-50">
                  <div class="bg-white rounded-xl shadow-xl border border-gray-100 py-2 overflow-hidden">
                    <component
                      :is="isExternalNavLink(sub) ? 'a' : 'router-link'"
                      v-for="sub in (item.children || item.megaItems || [])"
                      :key="sub.id || sub.label"
                      v-bind="linkAttrs(sub)"
                      class="flex items-center gap-2 px-4 py-2.5 text-sm text-[#1a365d] hover:bg-[#f5a623]/10 hover:text-[#f5a623] transition-colors"
                    >
                      <span v-if="sub.icon" class="text-base">{{ sub.icon }}</span>
                      {{ sub.label }}
                    </component>
                  </div>
                </div>
              </div>
              <component
                :is="isExternalNavLink(item) ? 'a' : 'router-link'"
                v-else
                v-bind="linkAttrs(item)"
                :class="['px-3 py-2 font-medium transition-colors duration-300 text-sm', isActiveLink(item) ? 'text-[#f5a623]' : 'text-[#1a365d] hover:text-[#f5a623]']"
              >{{ item.label }}</component>
            </template>
          </nav>

          <!-- Right controls -->
          <div class="flex items-center gap-2">
            <form class="relative transition-all duration-300" :class="searchOpen ? 'w-44 md:w-52' : 'w-0'" @submit.prevent="submitSearch">
              <label for="public-site-search" class="sr-only">Search NCS website</label>
              <input id="public-site-search" v-model="siteSearch" type="search" placeholder="Search..."
                :class="['flex h-9 rounded-md border bg-transparent px-3 py-1 text-sm shadow-sm transition-all duration-300 border-[#1a365d]/20 focus:border-[#f5a623] focus:outline-none focus-visible:ring-1 focus-visible:ring-[#f5a623]', searchOpen ? 'w-full opacity-100' : 'w-0 opacity-0 pointer-events-none']" />
            </form>
            <button type="button" class="inline-flex items-center justify-center h-9 w-9 rounded-md text-[#1a365d] hover:text-[#f5a623] hover:bg-[#f5a623]/10 transition-colors" :aria-expanded="searchOpen" aria-controls="public-site-search" aria-label="Search website" @click="searchOpen = !searchOpen">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><path d="m21 21-4.34-4.34"/><circle cx="11" cy="11" r="8"/></svg>
            </button>
            <button type="button" aria-pressed="false" aria-label="Accessibility tools" title="Accessibility tools" class="inline-flex items-center justify-center w-10 h-10 rounded-md text-[#1a365d] hover:text-[#f5a623] hover:bg-[#f5a623]/10 transition-colors" @click="openAccessibility">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5" aria-hidden="true"><circle cx="16" cy="4" r="1"/><path d="m18 19 1-7-6 1"/><path d="m5 8 3-3 5.5 3-2.36 3.5"/><path d="M4.24 14.5a5 5 0 0 0 6.88 6"/><path d="M13.76 17.5a5 5 0 0 0-6.88-6"/></svg>
            </button>
            <ThemeToggle />

            <a :href="portalUrl('/login')" class="hidden md:inline-flex h-8 items-center gap-2 rounded-md px-3 text-xs font-medium text-[#1a365d] hover:text-[#f5a623] hover:bg-[#f5a623]/10 transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" x2="3" y1="12" y2="12"/></svg>
              Login
            </a>
            <!-- Mobile menu toggle -->
            <button type="button" class="inline-flex items-center justify-center h-9 w-9 lg:hidden rounded-md text-[#1a365d] hover:text-[#f5a623] hover:bg-[#f5a623]/10 transition-colors" @click="mobileOpen = !mobileOpen" :aria-label="mobileOpen ? 'Close navigation menu' : 'Open navigation menu'" :aria-expanded="mobileOpen" aria-controls="mobile-navigation">
              <svg v-if="!mobileOpen" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6" aria-hidden="true"><path d="M4 12h16"/><path d="M4 18h16"/><path d="M4 6h16"/></svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-6 h-6" aria-hidden="true"><path d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Mobile menu (native details/summary per spec) -->
      <div :class="['lg:hidden absolute top-full left-0 right-0 bg-white shadow-xl transition-all duration-500 overflow-hidden', mobileOpen ? 'max-h-[80vh] opacity-100' : 'max-h-0 opacity-0 pointer-events-none']">
        <nav id="mobile-navigation" class="p-4 space-y-2 max-h-[70vh] overflow-y-auto" aria-label="Mobile navigation">
          <template v-for="item in visibleMenuItems" :key="item.id || item.label">
            <div v-if="(item.children?.length) || item.mega">
              <details class="group">
                <summary class="flex items-center justify-between px-4 py-3 text-[#1a365d] font-medium cursor-pointer hover:bg-[#f5a623]/10 rounded-lg transition-colors">
                  {{ item.label }}
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 group-open:rotate-180 transition-transform" aria-hidden="true"><path d="m6 9 6 6 6-6"/></svg>
                </summary>
                <div class="pl-4 mt-1 space-y-1">
                  <component
                    :is="isExternalNavLink(sub) ? 'a' : 'router-link'"
                    v-for="sub in (item.children || item.megaItems || [])"
                    :key="sub.id || sub.label"
                    v-bind="linkAttrs(sub)"
                    class="block px-4 py-2 text-[#1a365d]/80 hover:text-[#f5a623] transition-colors text-sm"
                    @click="mobileOpen = false"
                  >{{ sub.label }}</component>
                </div>
              </details>
            </div>
            <div v-else>
              <component
                :is="isExternalNavLink(item) ? 'a' : 'router-link'"
                v-bind="linkAttrs(item)"
                class="block px-4 py-3 text-[#1a365d] font-medium hover:bg-[#f5a623]/10 rounded-lg transition-colors"
                @click="mobileOpen = false"
              >{{ item.label }}</component>
            </div>
          </template>

          <div class="flex gap-3 mt-4 pt-4 border-t border-gray-100">
            <a :href="portalUrl('/login')" class="flex-1 py-2.5 text-sm font-semibold text-center border border-gray-300 rounded-lg text-[#1a365d] hover:border-[#f5a623] hover:text-[#f5a623] flex items-center justify-center gap-1.5" @click="mobileOpen = false">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" x2="3" y1="12" y2="12"/></svg>
              Login
            </a>
            <a :href="portalUrl('/register')" class="flex-1 py-2.5 text-sm font-semibold text-center border border-[#f5a623] text-[#f5a623] hover:bg-[#f5a623]/10 rounded-lg flex items-center justify-center gap-1.5" @click="mobileOpen = false">
              <i class="icofont-ui-user-group"></i> Register
            </a>
          </div>
        </nav>
      </div>
    </header>

    <!-- Page content -->
    <main id="main-content" tabindex="-1" class="flex-1">
      <router-view />
    </main>

    <!-- ── Footer ──────────────────────────────────────────── -->
    <footer id="public-footer" tabindex="-1" style="background-color: #252641;" aria-label="Website footer">

      <!-- Newsletter bar (NCS Newsletter Subscription spec) -->
      <div class="border-b border-white/10">
        <div class="max-w-7xl mx-auto px-4 py-8 md:py-12">
          <div class="flex flex-col md:flex-row items-center justify-between gap-6">
            <div class="text-center md:text-left">
              <h3 class="text-2xl md:text-3xl font-bold text-white mb-2">Stay Updated with NCS</h3>
              <p class="text-white/90">Subscribe to our newsletter for the latest sports news and events</p>
            </div>
            <form class="flex w-full md:w-auto gap-2" @submit.prevent="submitNewsletter">
              <label for="newsletter-email" class="sr-only">Email address for NCS updates</label>
              <input
                id="newsletter-email"
                v-model="newsletterEmail"
                :disabled="newsletterStatus === 'sending'"
                type="email"
                placeholder="Enter your email"
                autocomplete="email"
                required
                class="flex h-9 w-full rounded-md border px-3 py-1 text-base shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm flex-1 md:w-72 bg-white/20 border-white/30 text-white placeholder:text-white/70 focus:bg-white focus:text-[#1a365d] focus:placeholder:text-gray-500"
              />
              <button
                v-click-once
                type="submit"
                :disabled="newsletterStatus === 'sending'"
                :aria-label="newsletterStatus === 'sending' ? 'Subscribing…' : 'Subscribe'"
                class="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all duration-150 active:scale-95 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-60 shadow h-9 py-2 px-6 bg-[#1a365d] hover:bg-[#1a365d]/90 text-white"
              >
                <svg v-if="newsletterStatus !== 'sending'" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="M14.536 21.686a.5.5 0 0 0 .937-.024l6.5-19a.496.496 0 0 0-.635-.635l-19 6.5a.5.5 0 0 0-.024.937l7.93 3.18a2 2 0 0 1 1.112 1.11z"/><path d="m21.854 2.147-10.94 10.939"/></svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 animate-spin" aria-hidden="true"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
              </button>
            </form>
          </div>
          <p v-if="newsletterMessage" :class="['text-center md:text-right text-sm mt-3', newsletterStatus === 'error' ? 'text-red-300' : 'text-emerald-300']" role="status">{{ newsletterMessage }}</p>
        </div>
      </div>

      <!-- Links grid (NCS Footer Section Design spec) -->
      <div class="max-w-7xl mx-auto px-4 py-12 md:py-16">
        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-10 md:gap-8">

          <!-- Column 1 — Brand + about + social -->
          <div>
            <div class="flex items-center gap-3 mb-6">
              <img :src="footerLogo" :alt="`${siteIdentity.name || 'NCS'} - National Council of Sports`" class="h-14 w-auto object-contain" :class="!siteIdentity.footerWhiteLogoUrl && !siteIdentity.whiteLogoUrl ? 'brightness-0 invert' : ''" />
            </div>
            <p class="text-white/80 text-sm leading-relaxed mb-6">{{ footerSettings.about }}</p>
            <div class="flex gap-3">
              <a v-for="net in footerSocialLinks" :key="net.key" :href="net.url" target="_blank" rel="noopener noreferrer" :aria-label="net.label" class="w-10 h-10 rounded-full bg-white/10 flex items-center justify-center text-white hover:bg-[#f5a623] hover:text-white transition-all duration-300 hover:-translate-y-1">
                <i :class="[net.icon, 'text-base']" aria-hidden="true"></i>
              </a>
            </div>
          </div>

          <!-- Column 2 — Contact Us -->
          <div>
            <h4 class="text-lg font-bold mb-6 flex items-center gap-2 text-white">
              <span class="w-8 h-0.5 bg-[#f5a623]"></span>Contact Us
            </h4>
            <div class="space-y-4">
              <div v-if="contactAddressLines.length" class="flex items-start gap-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623] flex-shrink-0 mt-0.5" aria-hidden="true"><path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0"/><circle cx="12" cy="10" r="3"/></svg>
                <div class="text-sm text-white/80">
                  <p v-for="(line, i) in contactAddressLines" :key="i">{{ line }}</p>
                </div>
              </div>
              <div v-if="contact.phone" class="flex items-center gap-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623] flex-shrink-0" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
                <a :href="`tel:${primaryPhone}`" class="text-sm text-white/80 hover:text-[#f5a623] transition-colors">{{ contact.phone }}</a>
              </div>
              <div v-if="contact.email" class="flex items-center gap-3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5 text-[#f5a623] flex-shrink-0" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"/><rect x="2" y="4" width="20" height="16" rx="2"/></svg>
                <a :href="`mailto:${contact.email}`" class="text-sm text-white/80 hover:text-[#f5a623] transition-colors break-all">{{ contact.email }}</a>
              </div>
            </div>
            <a v-if="headerSettings.webmail_url" :href="headerSettings.webmail_url" target="_blank" rel="noopener noreferrer" class="inline-flex items-center justify-center gap-2 mt-6 h-9 px-4 py-2 rounded-md text-sm font-medium bg-[#f5a623] hover:bg-[#e09612] text-white shadow transition-colors">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4" aria-hidden="true"><path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7"/><rect x="2" y="4" width="20" height="16" rx="2"/></svg>
              Webmail
            </a>
          </div>

          <!-- Column 3 — Quick Links (dynamic from footerSettings.columns[0]) -->
          <div>
            <h4 class="text-lg font-bold mb-6 flex items-center gap-2 text-white">
              <span class="w-8 h-0.5 bg-[#f5a623]"></span>{{ quickLinksColumn.title || 'Quick Links' }}
            </h4>
            <ul class="space-y-3">
              <li v-for="(link, idx) in quickLinksColumn.links" :key="idx">
                <component
                  :is="isExternalLink(link.url) ? 'a' : 'router-link'"
                  v-bind="isExternalLink(link.url) ? { href: link.url, target: '_blank', rel: 'noopener' } : { to: link.url || '/' }"
                  class="flex items-center gap-2 text-sm text-white/80 hover:text-[#f5a623] transition-colors group"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 group-hover:translate-x-1 transition-transform" aria-hidden="true"><path d="m9 18 6-6-6-6"/></svg>
                  {{ link.label }}
                </component>
              </li>
              <li v-if="!quickLinksColumn.links?.length" class="text-xs text-white/40 italic">No quick links configured yet.</li>
            </ul>
          </div>

          <!-- Column 4 — Documents (dynamic from footerSettings.columns[1]) -->
          <div>
            <h4 class="text-lg font-bold mb-6 flex items-center gap-2 text-white">
              <span class="w-8 h-0.5 bg-[#f5a623]"></span>{{ documentsColumn.title || 'Documents' }}
            </h4>
            <ul class="space-y-3">
              <li v-for="(doc, idx) in documentsColumn.links" :key="idx">
                <a :href="doc.url" target="_blank" rel="noopener noreferrer" class="flex items-center gap-2 text-sm text-white/80 hover:text-[#f5a623] transition-colors group">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-red-400 flex-shrink-0" aria-hidden="true"><path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/><path d="M10 9H8"/><path d="M16 13H8"/><path d="M16 17H8"/></svg>
                  <span class="flex-1 line-clamp-1">{{ doc.label }}</span>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3 h-3 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0" aria-hidden="true"><path d="M15 3h6v6"/><path d="M10 14 21 3"/><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/></svg>
                </a>
              </li>
              <li v-if="!documentsColumn.links?.length" class="text-xs text-white/40 italic">No documents added yet.</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Bottom bar -->
      <div class="border-t border-white/10">
        <div class="max-w-screen-xl mx-auto px-6 py-5 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-gray-500">
          <span>&copy; {{ currentYear }} {{ footerSettings.copyright }}</span>
          <div class="flex gap-5">
            <router-link to="/faqs" class="hover:text-gray-300 transition-colors">FAQs</router-link>
            <router-link to="/contact-us" class="hover:text-gray-300 transition-colors">Contact</router-link>
          </div>
        </div>
      </div>
    </footer>
    <PublicAccessibilityMenu />
    <AppPreloader />
  </div>
</template>

<script setup>
import { ref, computed, nextTick, onMounted, onUnmounted, reactive, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useRouter } from 'vue-router'
import { getMenu, getSettings } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'
import PublicAccessibilityMenu from '@/components/public/PublicAccessibilityMenu.vue'
import AppPreloader from '@/components/public/AppPreloader.vue'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'

const mobileOpen = ref(false)
const searchOpen = ref(false)
const siteSearch = ref('')
const scrolled = ref(false)
const route = useRoute()
const router = useRouter()
const intranetUrl = (import.meta.env?.VITE_INTRANET_URL || 'http://localhost:9081').replace(/\/$/, '')
const currentYear = computed(() => new Date().getFullYear())

function onScroll() { scrolled.value = window.scrollY > 20 }
function openAccessibility() { window.dispatchEvent(new CustomEvent('open-accessibility-menu')) }
function portalUrl(path = '') { return `${intranetUrl}${path}` }
function submitSearch() {
  const query = siteSearch.value.trim()
  if (query) router.push({ path: '/news', query: { search: query } })
}
onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onScroll))

watch(() => route.fullPath, async () => {
  mobileOpen.value = false
  await nextTick()
  document.getElementById('main-content')?.focus({ preventScroll: true })
})

const defaultMenu = [
  { label: 'Home',            url: '/' },
  { label: 'News',            url: '/news' },
  { label: 'Events',          url: '/events' },
  { label: 'Careers',         url: '/careers' },
  { label: 'Facilities',      url: '/facilities' },
  { label: 'Associations',    url: '/associations' },
  { label: 'Resource Centre', url: '/resource-centre' },
  { label: 'Invest with Us',  url: '/invest' },
]

const defaultFooter = {
  about: 'The National Council of Sports is the statutory body mandated to develop, promote, and control sports in Uganda.',
  copyright: 'National Council of Sports, Uganda. All rights reserved.',
  columns: [
    // Column 0 — rendered as Quick Links (chevron icons)
    { title: 'Quick Links', links: [
      { label: 'Home',           url: '/' },
      { label: 'Associations',   url: '/associations' },
      { label: 'Invest With Us', url: '/invest' },
      { label: 'Careers',        url: '/careers' },
      { label: 'Contact',        url: '/contact-us' },
      { label: 'About NCS',      url: '/about/mandate' },
      { label: 'Latest News',    url: '/news' },
    ]},
    // Column 1 — rendered as Documents (PDF/external file icons)
    { title: 'Documents', links: [
      { label: 'NCS Annual Report 2023/2024',   url: 'https://www.ncs.go.ug/files/NCS%20ANNUAL%20REPORT%202023-2024_0.pdf' },
      { label: 'NCS Annual Report 2022/2023',   url: 'https://www.ncs.go.ug/files/NCS%20ANNUAL%20REPORT%202022-2023.pdf' },
      { label: 'NCS Strategic Plan 2020-2025',  url: 'https://www.ncs.go.ug/files/NCS%20STRATEGIC%20PLAN%202020%202025.pdf' },
      { label: 'Policy Guidelines',             url: 'https://www.ncs.go.ug/files/POLICY%20GUIDELINES.pdf' },
      { label: 'National Development Plan IV',  url: 'https://ncs.go.ug/files/NATIONAL%20DEVELOPMENT%20PLAN%20IV.pdf' },
    ]},
  ]
}

const menuItems = ref(defaultMenu)
const footerSettings = reactive({ ...defaultFooter, columns: [...defaultFooter.columns] })

// Quick Links uses the first dynamic column (CMS), Documents uses the second
const quickLinksColumn = computed(() => footerSettings.columns?.[0] || defaultFooter.columns[0])
const documentsColumn = computed(() => footerSettings.columns?.[1] || defaultFooter.columns[1])

// Address rendered as multi-line block: split contact.address on newlines (or commas),
// then append postal_address as a separate line if present.
const contactAddressLines = computed(() => {
  const lines = []
  if (contact.address) {
    const parts = String(contact.address).split(/\r?\n/).map(s => s.trim()).filter(Boolean)
    lines.push(...(parts.length ? parts : [String(contact.address).trim()]))
  }
  if (contact.postal_address) {
    String(contact.postal_address).split(/\r?\n/).forEach(s => { const t = s.trim(); if (t) lines.push(t) })
  }
  return lines
})

// First number from "+256 414254477 / 343688" so tel: works on mobile.
const primaryPhone = computed(() => {
  if (!contact.phone) return ''
  return String(contact.phone).split(/[\/,]/)[0].replace(/\s+/g, '')
})

function isExternalLink(url) {
  return /^(https?:|mailto:|tel:)/i.test(url || '')
}

// Newsletter subscription
const newsletterEmail = ref('')
const newsletterStatus = ref('idle')   // 'idle' | 'sending' | 'success' | 'error'
const newsletterMessage = ref('')

async function submitNewsletter() {
  const email = newsletterEmail.value.trim()
  if (!email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
    newsletterStatus.value = 'error'
    newsletterMessage.value = 'Please enter a valid email address.'
    return
  }
  newsletterStatus.value = 'sending'
  newsletterMessage.value = ''
  try {
    // Forward-compatible: if the backend endpoint exists it will accept this; otherwise we
    // gracefully queue locally so the admin can still recover the address later.
    const r = await fetch(`${import.meta.env.VITE_API_BASE_URL || ''}/api/v1/cms/newsletter/subscribe`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, source: 'public_footer' }),
    })
    if (r.ok) {
      newsletterStatus.value = 'success'
      newsletterMessage.value = 'Thanks — you are subscribed. Watch your inbox for updates.'
      newsletterEmail.value = ''
      return
    }
    // 404 or 5xx → queue locally so the email isn't lost
    throw new Error(`HTTP ${r.status}`)
  } catch {
    try {
      const queueKey = 'ncs_newsletter_queue'
      const queue = JSON.parse(localStorage.getItem(queueKey) || '[]')
      if (!queue.some(item => item.email === email)) {
        queue.push({ email, queued_at: new Date().toISOString(), source: 'public_footer' })
        localStorage.setItem(queueKey, JSON.stringify(queue))
      }
      newsletterStatus.value = 'success'
      newsletterMessage.value = 'Thanks — we have recorded your interest. We will follow up shortly.'
      newsletterEmail.value = ''
    } catch {
      newsletterStatus.value = 'error'
      newsletterMessage.value = 'Could not subscribe right now. Please try again later.'
    }
  }
}

const contact = reactive({
  phone: '+256 414254477 / 343688',
  whatsapp: '',
  email: 'info@ncs.go.ug',
  location: '',
  address: 'Plot 2-10, Coronation Avenue',
  postal_address: 'P.O. Box 20077, Lugogo\nKampala - UGANDA',
  fax: '',
  hours: '',
  mapUrl: '',
  social: { facebook: '', twitter: '', linkedin: '', instagram: '', youtube: '' }
})
const headerSettings = reactive({
  marquee: [
    'Welcome to National Council of Sports Uganda',
    'A centre of excellence for promotion and development of Sports.',
    'Maximizing opportunities for all Ugandans to participate and excel in Sports.',
    'Established 1964',
  ],
  webmail_url: 'https://mail.umcs.go.ug/',
  social: { facebook: '', twitter: '', linkedin: '', instagram: '', youtube: '' }
})
// Site identity (logos + favicon + name) — backed by cms_settings.site
const siteIdentity = reactive({
  name: 'NCS Uganda',
  logoUrl: '',
  whiteLogoUrl: '',
  footerLogoUrl: '',
  footerWhiteLogoUrl: '',
  faviconUrl: '',
})
function resolveAsset(path) {
  if (!path) return ''
  return mediaUrl(path)
}
// Footer logo: prefer white footer logo → footer logo → white logo → main logo → bundled fallback
const footerLogo = computed(() => {
  const chain = [siteIdentity.footerWhiteLogoUrl, siteIdentity.footerLogoUrl, siteIdentity.whiteLogoUrl, siteIdentity.logoUrl]
  for (const p of chain) { if (p) return resolveAsset(p) }
  return '/main-logo.png'
})
function applyFavicon(url) {
  if (!url) return
  const href = resolveAsset(url) || url
  let link = document.querySelector('link[rel="icon"]')
  if (!link) { link = document.createElement('link'); link.rel = 'icon'; document.head.appendChild(link) }
  link.href = href
}
const socialDefaults = {
  facebook:  'https://facebook.com/NCSUganda',
  twitter:   'https://twitter.com/NCSUganda1',
  linkedin:  'https://linkedin.com/company/ncsuganda',
  youtube:   'http://www.youtube.com/@NCSUgTV',
  instagram: '',
}
const socialNetworks = computed(() => [
  { key:'facebook',  label:'Facebook',     icon:'icofont-facebook',  url:contact.social.facebook  || socialDefaults.facebook },
  { key:'twitter',   label:'X / Twitter',  icon:'icofont-twitter',   url:contact.social.twitter   || socialDefaults.twitter },
  { key:'linkedin',  label:'LinkedIn',     icon:'icofont-linkedin',  url:contact.social.linkedin  || socialDefaults.linkedin },
  { key:'instagram', label:'Instagram',    icon:'icofont-instagram', url:contact.social.instagram || socialDefaults.instagram },
  { key:'youtube',   label:'YouTube',      icon:'icofont-youtube',   url:contact.social.youtube   || socialDefaults.youtube },
])
// Topbar social row (Header spec): plain icons w/ hover gold — FB, X, LinkedIn, YouTube, Webmail
const topbarSocialLinks = computed(() => {
  const rows = [
    { key:'facebook',  label:'Facebook',  url: headerSettings.social.facebook  || contact.social.facebook  || socialDefaults.facebook },
    { key:'twitter',   label:'Twitter',   url: headerSettings.social.twitter   || contact.social.twitter   || socialDefaults.twitter },
    { key:'linkedin',  label:'LinkedIn',  url: headerSettings.social.linkedin  || contact.social.linkedin  || socialDefaults.linkedin },
    { key:'youtube',   label:'YouTube',   url: headerSettings.social.youtube   || contact.social.youtube   || socialDefaults.youtube },
  ]
  if (headerSettings.webmail_url) {
    rows.push({ key:'webmail', label:'Webmail', url: headerSettings.webmail_url })
  }
  return rows.filter(r => r.url)
})

// Footer-spec social row: FB, X, LinkedIn, YouTube, Webmail (Instagram dropped to match the design)
const footerSocialLinks = computed(() => {
  const rows = [
    { key:'facebook',  label:'Facebook on NCS Uganda',  icon:'icofont-facebook',  url: headerSettings.social.facebook  || contact.social.facebook  || socialDefaults.facebook },
    { key:'twitter',   label:'X / Twitter',             icon:'icofont-twitter',   url: headerSettings.social.twitter   || contact.social.twitter   || socialDefaults.twitter },
    { key:'linkedin',  label:'LinkedIn',                icon:'icofont-linkedin',  url: headerSettings.social.linkedin  || contact.social.linkedin  || socialDefaults.linkedin },
    { key:'youtube',   label:'YouTube',                 icon:'icofont-youtube',   url: headerSettings.social.youtube   || contact.social.youtube   || socialDefaults.youtube },
  ]
  if (headerSettings.webmail_url) {
    rows.push({ key:'webmail', label:'Webmail', icon:'icofont-email', url: headerSettings.webmail_url })
  }
  return rows.filter(r => r.url)
})

function isActiveLink(item) {
  if (!item?.url) return false
  return route.path === item.url
}
function isActiveTopLevel(item) {
  const children = item.children || item.megaItems || []
  return children.some(c => c.url && route.path.startsWith(c.url))
}

function filterVisible(items) {
  return (items || [])
    .filter(item => !item.hidden)
    .map(item => ({
      ...item,
      children: item.children ? filterVisible(item.children) : item.children,
      megaItems: item.megaItems ? filterVisible(item.megaItems) : item.megaItems,
    }))
}
const visibleMenuItems = computed(() => filterVisible(menuItems.value))

function isExternalNavLink(item) {
  const url = item?.url || ''
  return item?.target === '_blank' || isExternalLink(url) || url.startsWith('#')
}

// Only include the keys relevant to each mode — passing an explicit
// href="undefined" alongside router-link's own :to would clobber the
// href it auto-generates, so external and internal links get distinct
// attribute sets rather than the same props toggled on/off.
function linkAttrs(item) {
  if (isExternalNavLink(item)) {
    return { href: item?.url || '/', target: item?.target, rel: item?.target === '_blank' ? 'noopener noreferrer' : undefined }
  }
  return { to: item?.url || '/' }
}

async function loadMenu() {
  try {
    const r = await getMenu('main')
    const items = r.data?.data?.items || r.data?.data?.Items || []
    if (items.length) menuItems.value = items
  } catch { /* keep defaults */ }
}

async function loadFooterSettings() {
  try {
    const r = await getSettings('footer')
    const v = r.data?.data?.value
    if (v && typeof v === 'object' && Object.keys(v).length) {
      if (typeof v.about     === 'string') footerSettings.about = v.about
      if (typeof v.copyright === 'string') footerSettings.copyright = v.copyright
      if (Array.isArray(v.columns) && v.columns.length) footerSettings.columns = v.columns
    }
  } catch { /* show defaults */ }
}

async function loadContact() {
  try {
    const r = await getSettings('contact')
    const v = r.data?.data?.value
    if (v && typeof v === 'object' && Object.keys(v).length) {
      if (typeof v.phone   === 'string') contact.phone = v.phone
      if (typeof v.whatsapp === 'string') contact.whatsapp = v.whatsapp
      if (typeof v.email   === 'string') contact.email = v.email
      if (typeof v.fax     === 'string') contact.fax = v.fax
      if (typeof v.location === 'string') contact.location = v.location
      if (typeof v.postal_address === 'string') contact.postal_address = v.postal_address
      if (typeof v.address === 'string') contact.address = v.address
      if (typeof v.hours   === 'string') contact.hours = v.hours
      if (typeof v.mapUrl  === 'string') contact.mapUrl = v.mapUrl
      if (v.social) Object.assign(contact.social, v.social)
    }
  } catch { /* keep defaults */ }
}

async function loadHeaderSettings() {
  const applyHeaderValue = (v) => {
    if (Array.isArray(v?.marquee) && v.marquee.length) headerSettings.marquee = v.marquee.filter(Boolean)
    if (typeof v?.webmail_url === 'string') headerSettings.webmail_url = v.webmail_url
    if (v?.social && typeof v.social === 'object') Object.assign(headerSettings.social, v.social)
  }
  try {
    const r = await getSettings('header')
    applyHeaderValue(r.data?.data?.value)
  } catch { /* keep defaults */ }
  try {
    const r = await getSettings('homepage')
    applyHeaderValue(r.data?.data?.value?.topbar)
  } catch { /* homepage topbar is optional */ }
}

async function loadSiteIdentity() {
  try {
    const r = await getSettings('site')
    const v = r.data?.data?.value
    if (v && typeof v === 'object') {
      if (typeof v.name === 'string')               siteIdentity.name = v.name
      if (typeof v.logoUrl === 'string')            siteIdentity.logoUrl = v.logoUrl
      if (typeof v.whiteLogoUrl === 'string')       siteIdentity.whiteLogoUrl = v.whiteLogoUrl
      if (typeof v.footerLogoUrl === 'string')      siteIdentity.footerLogoUrl = v.footerLogoUrl
      if (typeof v.footerWhiteLogoUrl === 'string') siteIdentity.footerWhiteLogoUrl = v.footerWhiteLogoUrl
      if (typeof v.faviconUrl === 'string')         siteIdentity.faviconUrl = v.faviconUrl
      applyFavicon(siteIdentity.faviconUrl)
    }
  } catch { /* first-time, no record yet */ }
}

onMounted(() => {
  loadMenu()
  loadFooterSettings()
  loadContact()
  loadHeaderSettings()
  loadSiteIdentity()
})
</script>

<style scoped>
.public-site > main { padding-top: 7rem; }

/* Marquee — single continuous track, replicates the Header spec's CSS animation */
.public-marquee-container { display: block; min-width: 0; }
.public-marquee-content { display: inline-block; white-space: nowrap; animation: public-marquee-scroll 34s linear infinite; }
@keyframes public-marquee-scroll { from { transform: translateX(0); } to { transform: translateX(-50%); } }
@media (prefers-reduced-motion: reduce) { .public-marquee-content { animation: none; } }
</style>
