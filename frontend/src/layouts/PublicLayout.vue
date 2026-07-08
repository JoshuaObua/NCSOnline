<template>
  <div v-if="maintenanceActive" class="maintenance-screen">
    <div class="maintenance-shell" role="status" aria-live="polite">
      <section class="maintenance-copy" aria-labelledby="public-maintenance-title">
        <div class="maintenance-brand">
          <img :src="maintenanceLogo" :alt="`${siteIdentity.name || 'NCS'} logo`" class="maintenance-logo" />
          <span>{{ siteIdentity.name || 'NCS Uganda' }}</span>
        </div>
        <p class="maintenance-kicker">Public Website Maintenance</p>
        <h1 id="public-maintenance-title">{{ maintenanceTitle }}</h1>
        <p class="maintenance-message">{{ maintenanceMessage }}</p>
        <div class="maintenance-actions" aria-label="Maintenance contact options">
          <a v-if="maintenanceEmailHref" :href="maintenanceEmailHref">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="m22 7-8.991 5.727a2 2 0 0 1-2.009 0L2 7" />
              <rect x="2" y="4" width="20" height="16" rx="2" />
            </svg>
            Email NCS
          </a>
          <a v-if="maintenancePhoneHref" :href="maintenancePhoneHref" class="secondary">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z" />
            </svg>
            Call NCS
          </a>
        </div>
      </section>

      <aside class="maintenance-panel" aria-label="Maintenance status">
        <div class="maintenance-icon" aria-hidden="true">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.1-3.1a6 6 0 0 1-7.8 7.8l-5.6 5.6a2.1 2.1 0 0 1-3-3l5.6-5.6a6 6 0 0 1 7.8-7.8l-3.1 3.1Z" />
          </svg>
        </div>
        <p class="maintenance-status-label">Current Status</p>
        <h2>Scheduled upgrades are in progress</h2>
        <div class="maintenance-status-grid">
          <div>
            <span>Expected return</span>
            <strong>{{ maintenanceExpectedEnd || 'Shortly' }}</strong>
          </div>
          <div>
            <span>Started</span>
            <strong>{{ maintenanceStartedAt || 'In progress' }}</strong>
          </div>
        </div>
        <ol class="maintenance-steps">
          <li class="complete"><span></span>Updates started</li>
          <li class="active"><span></span>Quality checks</li>
          <li><span></span>Website restored</li>
        </ol>
      </aside>
    </div>
    <p class="maintenance-footnote">&copy; 1964 - {{ currentYear }} National Council of Sports. All Rights Reserved.</p>
  </div>
  <div v-else class="public-site min-h-screen flex flex-col">

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
            <img :src="headerLogo" :alt="`${siteIdentity.name || 'NCS'} - National Council of Sports`" class="h-14 md:h-16 w-auto object-contain" />
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
                  <div class="bg-white rounded-xl shadow-xl border border-gray-100 py-2">
                    <template v-for="sub in (item.children || item.megaItems || [])" :key="sub.id || sub.label">
                      <div v-if="sub.children?.length" class="relative group/sub">
                        <button type="button" class="w-full flex items-center justify-between gap-2 px-4 py-2.5 text-sm text-[#1a365d] hover:bg-[#f5a623]/10 hover:text-[#f5a623] transition-colors" aria-haspopup="menu" aria-expanded="false">
                          <span class="flex items-center gap-2"><span v-if="sub.icon" class="text-base">{{ sub.icon }}</span>{{ sub.label }}</span>
                          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 flex-shrink-0" aria-hidden="true"><path d="m9 18 6-6-6-6"/></svg>
                        </button>
                        <div class="absolute left-full top-0 pl-1 opacity-0 invisible group-hover/sub:opacity-100 group-hover/sub:visible group-focus-within/sub:opacity-100 group-focus-within/sub:visible transition-all duration-200 min-w-[200px] z-50">
                          <div class="bg-white rounded-xl shadow-xl border border-gray-100 py-2">
                            <component
                              :is="isExternalNavLink(grandsub) ? 'a' : 'router-link'"
                              v-for="grandsub in sub.children"
                              :key="grandsub.id || grandsub.label"
                              v-bind="linkAttrs(grandsub)"
                              class="flex items-center gap-2 px-4 py-2.5 text-sm text-[#1a365d] hover:bg-[#f5a623]/10 hover:text-[#f5a623] transition-colors"
                            >{{ grandsub.label }}</component>
                          </div>
                        </div>
                      </div>
                      <component
                        v-else
                        :is="isExternalNavLink(sub) ? 'a' : 'router-link'"
                        v-bind="linkAttrs(sub)"
                        class="flex items-center gap-2 px-4 py-2.5 text-sm text-[#1a365d] hover:bg-[#f5a623]/10 hover:text-[#f5a623] transition-colors"
                      >
                        <span v-if="sub.icon" class="text-base">{{ sub.icon }}</span>
                        {{ sub.label }}
                      </component>
                    </template>
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

            <div v-if="currentUser" class="account-menu">
              <button type="button" class="account-trigger" aria-haspopup="menu" :aria-expanded="accountOpen" @click="accountOpen = !accountOpen">
                <img v-if="currentUser.avatar_url" :src="currentUser.avatar_url" alt="" referrerpolicy="no-referrer" />
                <span v-else>{{ accountInitials }}</span>
              </button>
              <div v-show="accountOpen" class="account-dropdown" role="menu">
                <router-link to="/account/profile" role="menuitem">Profile Overview</router-link>
                <router-link to="/account/settings" role="menuitem">Settings & Security</router-link>
                <router-link to="/account/activities" role="menuitem">My Audit Activities</router-link>
                <button type="button" role="menuitem" @click="logoutAccount">Sign out</button>
              </div>
            </div>

            <a v-else :href="portalUrl('/login')" class="hidden md:inline-flex h-8 items-center gap-2 rounded-md px-3 text-xs font-medium text-[#1a365d] hover:text-[#f5a623] hover:bg-[#f5a623]/10 transition-colors">
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
                  <template v-for="sub in (item.children || item.megaItems || [])" :key="sub.id || sub.label">
                    <details v-if="sub.children?.length" class="group/sub">
                      <summary class="flex items-center justify-between px-4 py-2 text-[#1a365d]/90 font-medium text-sm cursor-pointer hover:text-[#f5a623]">
                        {{ sub.label }}
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-3.5 h-3.5 group-open/sub:rotate-180 transition-transform" aria-hidden="true"><path d="m6 9 6 6 6-6"/></svg>
                      </summary>
                      <div class="pl-4 mt-1 space-y-1">
                        <component
                          :is="isExternalNavLink(grandsub) ? 'a' : 'router-link'"
                          v-for="grandsub in sub.children"
                          :key="grandsub.id || grandsub.label"
                          v-bind="linkAttrs(grandsub)"
                          class="block px-4 py-2 text-[#1a365d]/70 hover:text-[#f5a623] transition-colors text-sm"
                          @click="mobileOpen = false"
                        >{{ grandsub.label }}</component>
                      </div>
                    </details>
                    <component
                      v-else
                      :is="isExternalNavLink(sub) ? 'a' : 'router-link'"
                      v-bind="linkAttrs(sub)"
                      class="block px-4 py-2 text-[#1a365d]/80 hover:text-[#f5a623] transition-colors text-sm"
                      @click="mobileOpen = false"
                    >{{ sub.label }}</component>
                  </template>
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
      <div class="bg-[#f5a623]">
        <div class="max-w-7xl mx-auto px-4 py-8 md:py-12">
          <div class="flex flex-col md:flex-row items-center justify-between gap-6">
            <div class="text-center md:text-left">
              <h3 class="text-2xl md:text-3xl font-bold text-white mb-2">Stay Updated with NCS</h3>
              <p class="text-white/90">Subscribe to our newsletter for the latest sports news and events</p>
            </div>
            <form class="flex w-full md:w-auto gap-2" @submit.prevent="submitNewsletter">
              <label for="newsletter-email" class="sr-only">Email address for NCS updates</label>
              <input v-model="newsletterWebsite" type="text" tabindex="-1" autocomplete="off" class="newsletter-hp" aria-hidden="true" />
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
          <p v-if="newsletterMessage" :class="['text-center md:text-right text-sm mt-3 font-medium', newsletterStatus === 'error' ? 'text-red-800' : 'text-[#1a365d]']" role="status">{{ newsletterMessage }}</p>
        </div>
      </div>

      <!-- Links grid (NCS Footer Section Design spec) -->
      <div class="max-w-7xl mx-auto px-4 py-12 md:py-16">
        <div class="grid md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-10 md:gap-8">

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

          <!-- Columns 3+ — fully dynamic footer link columns (title + links both
               editable in the CMS under Settings → Footer Columns). Any number
               of columns can exist; each renders identically. -->
          <div v-for="(column, colIdx) in footerSettings.columns" :key="colIdx">
            <h4 class="text-lg font-bold mb-6 flex items-center gap-2 text-white">
              <span class="w-8 h-0.5 bg-[#f5a623]"></span>{{ column.title || 'Links' }}
            </h4>
            <ul class="space-y-3">
              <li v-for="(link, idx) in column.links" :key="idx">
                <component
                  :is="isExternalLink(link.url) ? 'a' : 'router-link'"
                  v-bind="isExternalLink(link.url) ? { href: link.url, target: '_blank', rel: 'noopener' } : { to: link.url || '/' }"
                  class="flex items-center gap-2 text-sm text-white/80 hover:text-[#f5a623] transition-colors group"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 group-hover:translate-x-1 transition-transform" aria-hidden="true"><path d="m9 18 6-6-6-6"/></svg>
                  {{ link.label }}
                </component>
              </li>
              <li v-if="!column.links?.length" class="text-xs text-white/40 italic">No links configured yet.</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Bottom bar -->
      <div class="border-t border-white/10">
        <div class="max-w-7xl mx-auto px-4 py-6">
          <div class="flex flex-col md:flex-row items-center justify-between gap-4">
            <p class="text-sm text-white/60 text-center md:text-left">&copy; 1964 - {{ currentYear }} National Council of Sports. All Rights Reserved.</p>
            <div class="flex items-center gap-2 text-sm text-white/60">
              <span>Designed &amp; Developed by</span>
              <span class="text-[#f5a623] font-medium">NCS Uganda</span>
            </div>
          </div>
        </div>
      </div>
    </footer>
    <PublicAccessibilityMenu />
    <ChatBotWidget />
    <AppPreloader />
  </div>
</template>

<script setup>
import { ref, computed, nextTick, onMounted, onUnmounted, reactive, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useRouter } from 'vue-router'
import { getMaintenanceStatus, getMenu, getSettings, getTypography, listFonts } from '@/api/cms.js'
import { injectTypographyCSS } from '@/utils/typography.js'
import { mediaUrl } from '@/api/client.js'
import { installAnalyticsTracker, trackPageView, uninstallAnalyticsTracker } from '@/utils/analyticsTracker.js'
import { installGoogleAnalytics, trackGoogleAnalyticsPageView, uninstallGoogleAnalytics } from '@/utils/googleAnalytics.js'
import PublicAccessibilityMenu from '@/components/public/PublicAccessibilityMenu.vue'
import AppPreloader from '@/components/public/AppPreloader.vue'
import ChatBotWidget from '@/components/public/ChatBotWidget.vue'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'
import { useTheme } from '@/composables/useTheme.js'
import { animatePublicPage, cleanupPublicMotion } from '@/utils/publicMotion.js'

const mobileOpen = ref(false)
const searchOpen = ref(false)
const accountOpen = ref(false)
const siteSearch = ref('')
const scrolled = ref(false)
const accountUser = ref(readAccountUser())
const route = useRoute()
const router = useRouter()
const { isDark } = useTheme()
const maintenanceActive = ref(false)
const maintenanceInfo = ref({})
// Priority: explicit intranet URL, then the browser host plus VITE_INTRANET_PORT.
// Mirrors resolveApiBase() in api/client.js so neither URL is pinned to a
// specific host at build time.
function resolveIntranetUrl() {
  const envBase = import.meta.env?.VITE_INTRANET_URL
  if (envBase) return envBase.replace(/\/$/, '')
  if (typeof window !== 'undefined' && window.location?.hostname) {
    const port = import.meta.env?.VITE_INTRANET_PORT
    return port
      ? `${window.location.protocol}//${window.location.hostname}:${port}`
      : window.location.origin
  }
  return ''
}
const intranetUrl = resolveIntranetUrl()
const currentYear = computed(() => new Date().getFullYear())
const maintenanceTitle = computed(() => maintenanceInfo.value?.display_meta?.custom_title || "We'll be right back")
const maintenanceMessage = computed(() => maintenanceInfo.value?.display_meta?.custom_message
  || maintenanceInfo.value?.reason
  || 'The National Council of Sports website is temporarily unavailable while scheduled maintenance is in progress.')
const maintenanceExpectedEnd = computed(() => formatMaintenanceTime(maintenanceInfo.value?.expected_end))
const maintenanceStartedAt = computed(() => formatMaintenanceTime(maintenanceInfo.value?.scheduled_start || maintenanceInfo.value?.changed_at))
const maintenanceEmailHref = computed(() => contact.email ? `mailto:${contact.email}` : '')
const maintenancePhoneHref = computed(() => primaryPhone.value ? `tel:${primaryPhone.value}` : '')
const currentUser = computed(() => accountUser.value)
const accountName = computed(() => {
  const user = currentUser.value || {}
  return `${user.first_name || ''} ${user.last_name || ''}`.trim() || user.email || 'Account'
})
const accountInitials = computed(() => accountName.value.split(/\s+/).slice(0, 2).map(part => part[0] || '').join('').toUpperCase() || 'U')

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
  accountUser.value = readAccountUser()
  mobileOpen.value = false
  accountOpen.value = false
  trackPageView(route.path)
  trackGoogleAnalyticsPageView(route.path)
  await nextTick()
  const main = document.getElementById('main-content')
  main?.focus({ preventScroll: true })
  animatePublicPage(main)
})

function logoutAccount() {
  localStorage.removeItem('ncsms_access_token')
  localStorage.removeItem('ncsms_user')
  accountUser.value = null
  accountOpen.value = false
  router.push('/')
}

function readAccountUser() {
  try { return JSON.parse(localStorage.getItem('ncsms_user') || 'null') } catch { return null }
}

const defaultMenu = [
  { label: 'Home',            url: '/' },
  { label: 'News',            url: '/news' },
  { label: 'Events',          url: '/events' },
  { label: 'Careers',         url: '/careers' },
  { label: 'Facilities',      url: '/facilities' },
  { label: 'Associations',    url: '/associations' },
  { label: 'Resource Centre', url: '/resource-centre' },
  { label: 'Sports Rules',    url: '/sports-rules' },
  { label: 'Press Releases',  url: '/press-releases' },
  { label: 'NCS Reports',     url: '/reports' },
  { label: 'NCS Speeches',    url: '/speeches' },
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
const footerMenuItems = ref([])
const footerSettings = reactive({ ...defaultFooter, columns: [...defaultFooter.columns] })

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
const newsletterWebsite = ref('')
const newsletterStatus = ref('idle')   // 'idle' | 'sending' | 'success' | 'error'
const newsletterMessage = ref('')

async function submitNewsletter() {
  const email = newsletterEmail.value.trim()
  if (newsletterWebsite.value.trim()) {
    newsletterStatus.value = 'success'
    newsletterMessage.value = 'Thanks - you are subscribed. Watch your inbox for updates.'
    newsletterEmail.value = ''
    return
  }
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
      body: JSON.stringify({ email, source: 'public_footer', website: newsletterWebsite.value }),
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
// Header logo swaps with the light/dark toggle: white logo on dark backgrounds,
// main logo on light — falling back to the other if only one has been uploaded.
const headerLogo = computed(() => {
  const chain = isDark.value
    ? [siteIdentity.whiteLogoUrl, siteIdentity.logoUrl]
    : [siteIdentity.logoUrl, siteIdentity.whiteLogoUrl]
  for (const p of chain) { if (p) return resolveAsset(p) }
  return '/main-logo.png'
})
const maintenanceLogo = computed(() => {
  const chain = [siteIdentity.whiteLogoUrl, siteIdentity.logoUrl, siteIdentity.footerWhiteLogoUrl, siteIdentity.footerLogoUrl]
  for (const p of chain) { if (p) return resolveAsset(p) }
  return '/main-logo.png'
})

function formatMaintenanceTime(value) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleString('en-UG', {
    weekday: 'short',
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function applyFavicon(url) {
  if (!url) return
  const href = resolveAsset(url) || url
  let link = document.querySelector('link[rel="icon"]')
  if (!link) { link = document.createElement('link'); link.rel = 'icon'; document.head.appendChild(link) }
  link.href = href
}
function setMetaTag(selector, attrs) {
  let tag = document.querySelector(selector)
  if (!tag) { tag = document.createElement('meta'); document.head.appendChild(tag) }
  Object.entries(attrs).forEach(([k, v]) => tag.setAttribute(k, v))
}
async function loadSeoDefaults() {
  try {
    const r = await getSettings('seo')
    const v = r.data?.data?.value
    if (!v || typeof v !== 'object') return
    if (v.meta_description) setMetaTag('meta[name="description"]', { name: 'description', content: v.meta_description })
    if (v.meta_keywords) setMetaTag('meta[name="keywords"]', { name: 'keywords', content: v.meta_keywords })
    if (v.meta_image_url) setMetaTag('meta[property="og:image"]', { property: 'og:image', content: resolveAsset(v.meta_image_url) })
  } catch { /* no SEO defaults saved yet */ }
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
function urlMatchesRoute(url) {
  // A dropdown's own highlighted state should never be driven by a child
  // pointing at "/" — the Home link already handles that via isActiveLink,
  // and a "/" child inside a dropdown is virtually always a stray/unfinished
  // menu entry, not an intentional "this section includes Home" case.
  if (!url || url === '/') return false
  return route.path === url || route.path.startsWith(`${url}/`)
}
function isActiveTopLevel(item) {
  const children = item.children || item.megaItems || []
  return children.some(c => urlMatchesRoute(c.url) || isActiveTopLevel(c))
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

async function loadFooterMenu() {
  try {
    const r = await getMenu('footer')
    const items = r.data?.data?.items || r.data?.data?.Items || []
    footerMenuItems.value = items.filter(item => !item.hidden)
  } catch { /* keep empty, bottom bar falls back to defaults */ }
}

async function loadTypography() {
  try {
    const [settingsRes, fontsRes] = await Promise.all([getTypography(), listFonts()])
    const value = settingsRes.data?.data?.value || {}
    const fonts = fontsRes.data?.data || []
    injectTypographyCSS(value, fonts)
  } catch { /* leave default site typography untouched */ }
}

async function loadThirdPartyIntegrations() {
  try {
    const r = await getSettings('third_party')
    installGoogleAnalytics(r.data?.data?.value || {})
  } catch {
    uninstallGoogleAnalytics()
  }
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

async function loadMaintenanceStatus() {
  try {
    const r = await getMaintenanceStatus()
    const scoped = r.data?.data?.public_cms
    // Staff already signed into the CMS can keep browsing the public site
    // normally while maintenance is on, so they can verify it live.
    if (scoped?.is_active && !localStorage.getItem('ncsms_access_token')) {
      maintenanceInfo.value = scoped
      maintenanceActive.value = true
    }
  } catch { /* if the status check fails, fail open and show the site */ }
}

onMounted(() => {
  loadMaintenanceStatus()
  installAnalyticsTracker()
  trackPageView(route.path)
  loadMenu()
  loadFooterMenu()
  loadFooterSettings()
  loadContact()
  loadHeaderSettings()
  loadSiteIdentity()
  loadSeoDefaults()
  loadTypography()
  loadThirdPartyIntegrations()
  nextTick(() => animatePublicPage(document.getElementById('main-content')))
})
onUnmounted(() => {
  uninstallAnalyticsTracker()
  uninstallGoogleAnalytics()
  cleanupPublicMotion()
})
</script>

<style scoped>
.maintenance-screen {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: center;
  background: linear-gradient(135deg, #102b4d 0%, #1a365d 52%, #0d223d 100%);
  color: #fff;
  padding: clamp(1.25rem, 4vw, 3.5rem);
}
.maintenance-screen::before {
  position: absolute;
  inset: 0;
  content: "";
  pointer-events: none;
  background-image:
    linear-gradient(rgb(255 255 255 / 0.045) 1px, transparent 1px),
    linear-gradient(90deg, rgb(255 255 255 / 0.045) 1px, transparent 1px);
  background-size: 52px 52px;
  mask-image: linear-gradient(to bottom right, #000, transparent 78%);
}
.maintenance-shell {
  position: relative;
  z-index: 1;
  width: min(100%, 1160px);
  margin: auto;
  display: grid;
  grid-template-columns: minmax(0, 1.08fr) minmax(320px, 0.72fr);
  align-items: center;
  gap: clamp(2rem, 7vw, 6.5rem);
}
.maintenance-copy { min-width: 0; }
.maintenance-brand {
  display: flex;
  align-items: center;
  gap: 0.9rem;
  margin-bottom: clamp(2.2rem, 7vh, 5rem);
  color: rgb(255 255 255 / 0.86);
  font-size: 0.82rem;
  font-weight: 700;
  text-transform: uppercase;
}
.maintenance-logo {
  display: block;
  width: auto;
  max-width: 168px;
  height: 60px;
  object-fit: contain;
  border-radius: 8px;
  background: #fff;
  padding: 0.35rem 0.7rem;
}
.maintenance-kicker {
  display: inline-flex;
  align-items: center;
  gap: 0.55rem;
  margin: 0 0 1.1rem;
  color: #f5a623;
  font-size: 0.78rem;
  font-weight: 800;
  text-transform: uppercase;
}
.maintenance-kicker::before {
  width: 28px;
  height: 2px;
  content: "";
  background: #f5a623;
}
.maintenance-copy h1 {
  max-width: 720px;
  margin: 0;
  color: #fff;
  font-size: clamp(2.65rem, 6vw, 5.35rem);
  font-weight: 800;
  line-height: 1.02;
  letter-spacing: 0;
  overflow-wrap: anywhere;
}
.maintenance-message {
  max-width: 650px;
  margin: 1.5rem 0 0;
  color: rgb(255 255 255 / 0.76);
  font-size: clamp(1rem, 1.8vw, 1.18rem);
  line-height: 1.75;
}
.maintenance-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-top: 2rem;
}
.maintenance-actions a {
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.55rem;
  border: 1px solid #f5a623;
  border-radius: 6px;
  background: #f5a623;
  color: #102b4d;
  padding: 0.75rem 1.1rem;
  font-size: 0.88rem;
  font-weight: 800;
  transition: background-color 160ms ease, color 160ms ease, border-color 160ms ease;
}
.maintenance-actions a:hover { background: #ffc154; border-color: #ffc154; }
.maintenance-actions a.secondary {
  border-color: rgb(255 255 255 / 0.34);
  background: transparent;
  color: #fff;
}
.maintenance-actions a.secondary:hover { border-color: #fff; background: rgb(255 255 255 / 0.08); }
.maintenance-actions svg { width: 18px; height: 18px; flex: 0 0 auto; }
.maintenance-panel {
  border: 1px solid rgb(255 255 255 / 0.22);
  border-radius: 8px;
  background: #fff;
  color: #1a365d;
  padding: clamp(1.5rem, 4vw, 2.4rem);
  box-shadow: 0 28px 80px rgb(4 18 36 / 0.34);
}
.maintenance-icon {
  width: 56px;
  height: 56px;
  display: grid;
  place-items: center;
  border-radius: 8px;
  background: #fff5e4;
  color: #e2920f;
}
.maintenance-icon svg { width: 29px; height: 29px; }
.maintenance-status-label {
  margin: 1.5rem 0 0.5rem;
  color: #e2920f;
  font-size: 0.75rem;
  font-weight: 800;
  text-transform: uppercase;
}
.maintenance-panel h2 {
  margin: 0;
  color: #1a365d;
  font-size: clamp(1.35rem, 2.6vw, 1.75rem);
  font-weight: 800;
  line-height: 1.25;
  letter-spacing: 0;
}
.maintenance-status-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1px;
  margin-top: 1.5rem;
  overflow: hidden;
  border: 1px solid #e5e9ef;
  border-radius: 6px;
  background: #e5e9ef;
}
.maintenance-status-grid div {
  min-width: 0;
  background: #f8fafc;
  padding: 1rem;
}
.maintenance-status-grid span,
.maintenance-status-grid strong { display: block; }
.maintenance-status-grid span {
  margin-bottom: 0.35rem;
  color: #6b7280;
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
}
.maintenance-status-grid strong {
  color: #1a365d;
  font-size: 0.88rem;
  line-height: 1.45;
  overflow-wrap: anywhere;
}
.maintenance-steps {
  display: grid;
  gap: 0.9rem;
  margin: 1.5rem 0 0;
  padding: 1.35rem 0 0;
  border-top: 1px solid #e5e9ef;
  list-style: none;
}
.maintenance-steps li {
  display: grid;
  grid-template-columns: 12px minmax(0, 1fr);
  align-items: center;
  gap: 0.7rem;
  color: #8a94a3;
  font-size: 0.86rem;
  font-weight: 700;
}
.maintenance-steps li span {
  width: 10px;
  height: 10px;
  border: 2px solid #cbd5e1;
  border-radius: 50%;
}
.maintenance-steps li.complete { color: #237a4b; }
.maintenance-steps li.complete span { border-color: #2fa96b; background: #2fa96b; }
.maintenance-steps li.active { color: #1a365d; }
.maintenance-steps li.active span {
  border-color: #f5a623;
  background: #f5a623;
  box-shadow: 0 0 0 4px rgb(245 166 35 / 0.18);
}
.maintenance-footnote {
  position: relative;
  z-index: 1;
  width: min(100%, 1160px);
  margin: clamp(2rem, 7vh, 4.5rem) auto 0;
  color: rgb(255 255 255 / 0.48);
  font-size: 0.76rem;
}
@media (max-width: 820px) {
  .maintenance-screen { justify-content: flex-start; }
  .maintenance-shell { grid-template-columns: 1fr; gap: 2.2rem; }
  .maintenance-brand { margin-bottom: 2.5rem; }
  .maintenance-copy h1 { font-size: clamp(2.45rem, 12vw, 4.2rem); }
  .maintenance-panel { max-width: 620px; }
}
@media (max-width: 480px) {
  .maintenance-screen { padding: 1.1rem; }
  .maintenance-brand { align-items: flex-start; flex-direction: column; gap: 0.65rem; }
  .maintenance-logo { height: 52px; max-width: 150px; }
  .maintenance-actions { display: grid; grid-template-columns: 1fr; }
  .maintenance-actions a { width: 100%; }
  .maintenance-status-grid { grid-template-columns: 1fr; }
  .maintenance-footnote { line-height: 1.6; }
}
.public-site > main { padding-top: 7rem; }
.newsletter-hp { position: absolute; left: -9999px; width: 1px; height: 1px; opacity: 0; }
.account-menu{position:relative}.account-trigger{width:38px;height:38px;border:1px solid #e5e7eb;border-radius:999px;background:white;color:#112b4e;display:grid;place-items:center;overflow:hidden;font-weight:900}.account-trigger img{width:100%;height:100%;object-fit:cover}.account-dropdown{position:absolute;right:0;top:calc(100% + 10px);width:230px;background:white;border:1px solid #e5e7eb;border-radius:8px;box-shadow:0 18px 40px rgba(15,23,42,.16);padding:.45rem;z-index:70}.account-dropdown a,.account-dropdown button{display:block;width:100%;border:0;background:transparent;border-radius:6px;padding:.7rem .75rem;text-align:left;color:#112b4e;font-size:.88rem;font-weight:800}.account-dropdown a:hover,.account-dropdown button:hover{background:#f8fafc;color:#f48c06}

/* Marquee — single continuous track, replicates the Header spec's CSS animation */
.public-marquee-container { display: block; min-width: 0; }
.public-marquee-content {
  display: inline-block;
  white-space: nowrap;
  animation: public-marquee-scroll 34s linear infinite;
  /* Keep the announcements ticker scrolling continuously even when the OS
     requests reduced motion. The global reduced-motion reset in style.css
     forces `animation-duration: 0.01ms !important; animation-iteration-count:
     1 !important` on `*`, which would otherwise freeze the marquee. Re-assert
     with !important — the class selector's specificity beats the universal
     (*) reset, so the ticker always runs. */
  animation-duration: 34s !important;
  animation-iteration-count: infinite !important;
  animation-timing-function: linear !important;
}
@keyframes public-marquee-scroll { from { transform: translateX(0); } to { transform: translateX(-50%); } }
</style>
