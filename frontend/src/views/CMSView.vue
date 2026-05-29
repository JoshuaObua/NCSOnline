<template>
  <LayoutDefault title="Content Management">
    <div class="flex -mx-6 -mt-6 -mb-6" style="min-height:calc(100vh - 64px)">

      <!-- ═══ CMS NAV SIDEBAR ═══════════════════════════════════════ -->
      <nav class="w-52 flex-shrink-0 bg-white border-r border-gray-200 overflow-y-auto py-2" style="min-height:calc(100vh - 64px)">
        <template v-for="g in navGroups" :key="g.id">
          <button @click="toggleGroup(g.id)"
            class="w-full flex items-center justify-between px-4 pt-4 pb-1 text-[10px] font-bold uppercase tracking-widest text-gray-400 hover:text-gray-600 transition-colors">
            <span>{{ g.label }}</span>
            <svg :class="collapsed.has(g.id) ? '-rotate-90' : ''" class="w-3 h-3 transition-transform" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5"/>
            </svg>
          </button>
          <template v-if="!collapsed.has(g.id)">
            <button v-for="item in g.items" :key="item.id"
              @click="navigate(item.id)"
              :class="activeSection===item.id ? 'bg-primary-50 text-primary-700 border-r-2 border-primary-600 font-medium' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'"
              class="w-full flex items-center px-4 py-2 text-sm transition-all text-left">
              <span class="truncate">{{ item.label }}</span>
            </button>
          </template>
        </template>
      </nav>

      <!-- ═══ CONTENT AREA ══════════════════════════════════════════ -->
      <div class="flex-1 min-w-0 overflow-y-auto p-6">

        <!-- Section header -->
        <div class="flex items-center justify-between mb-5">
          <div>
            <h2 class="text-xl font-bold text-gray-900">{{ currentNavItem?.label }}</h2>
            <p v-if="currentNavItem?.description" class="text-sm text-gray-500 mt-0.5">{{ currentNavItem.description }}</p>
          </div>
          <button v-if="creatableSections.includes(activeSection)" @click="openCreate"
            class="flex items-center gap-1.5 bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-4 py-2 rounded-lg transition-colors">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15"/></svg>
            New {{ currentNavItem?.singular || 'Item' }}
          </button>
        </div>

        <!-- Toast -->
        <Transition name="toast">
          <div v-if="toast" class="fixed top-5 right-5 z-[100] flex items-center gap-2 px-4 py-3 bg-gray-900 text-white text-sm rounded-xl shadow-xl">
            <svg class="w-4 h-4 text-green-400" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5"/></svg>
            {{ toast }}
          </div>
        </Transition>

        <!-- ═══ POSTS / BLOGS / CASE STUDIES / PAGES ═══════════════ -->
        <div v-if="['posts','blogs','case-studies','pages'].includes(activeSection)">
          <div class="flex gap-3 mb-4 flex-wrap">
            <select v-if="activeSection==='posts'" v-model="postCategory" @change="loadPosts()" class="text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white text-gray-700">
              <option value="">All Categories</option>
              <option value="news">News</option>
              <option value="announcement">Announcement</option>
              <option value="project">Project</option>
              <option value="blog">Blog</option>
              <option value="case_study">Case Study</option>
            </select>
            <select v-model="postStatus" @change="loadPosts()" class="text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white text-gray-700">
              <option value="">All Statuses</option>
              <option value="draft">Draft</option>
              <option value="published">Published</option>
            </select>
          </div>
          <ContentTable :items="posts" :loading="postsLoading"
            :cols="[{key:'title',label:'Title'},{key:'category',label:'Category'},{key:'status',label:'Status'},{key:'published_at',label:'Date'}]"
            @edit="editPost" @delete="id=>confirmDelete('post',id)" />
        </div>

        <!-- ═══ EVENTS ══════════════════════════════════════════════ -->
        <div v-else-if="activeSection==='events'">
          <ContentTable :items="events" :loading="eventsLoading"
            :cols="[{key:'title',label:'Title'},{key:'location',label:'Location'},{key:'status',label:'Status'},{key:'event_date',label:'Date'}]"
            @edit="editEvent" @delete="id=>confirmDelete('event',id)" />
        </div>

        <!-- ═══ CAREERS ════════════════════════════════════════════ -->
        <div v-else-if="activeSection==='careers'">
          <ContentTable :items="careers" :loading="careersLoading"
            :cols="[{key:'title',label:'Title'},{key:'department',label:'Department'},{key:'job_type',label:'Type'},{key:'status',label:'Status'}]"
            @edit="editCareer" @delete="id=>confirmDelete('career',id)" />
        </div>

        <!-- ═══ HERO SLIDES ════════════════════════════════════════ -->
        <div v-else-if="activeSection==='slides'">
          <div v-if="slidesLoading" class="space-y-2">
            <div v-for="i in 4" :key="i" class="h-20 bg-gray-100 rounded-xl animate-pulse"/>
          </div>
          <div v-else class="space-y-3">
            <div v-for="slide in slides" :key="slide.id"
              draggable="true"
              @dragstart="dragStart($event,slide)"
              @dragover.prevent
              @drop="dropSlide($event,slide)"
              class="bg-white rounded-xl border border-gray-200 shadow-sm p-4 flex items-center gap-4 cursor-grab active:opacity-60">
              <svg class="w-5 h-5 text-gray-300 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9h16.5m-16.5 6.75h16.5"/></svg>
              <img v-if="slide.image_url" :src="mediaUrl(slide.image_url)" class="w-16 h-10 object-cover rounded flex-shrink-0"/>
              <div v-else class="w-16 h-10 bg-gray-100 rounded flex-shrink-0 flex items-center justify-center text-gray-300 text-xs">No img</div>
              <div class="flex-1 min-w-0">
                <div class="font-medium text-gray-900 truncate">{{ slide.title }}</div>
                <div class="text-xs text-gray-400 truncate">{{ slide.subtitle }}</div>
              </div>
              <label class="flex items-center gap-1.5 cursor-pointer flex-shrink-0">
                <input type="checkbox" :checked="slide.is_active" @change="toggleSlideActive(slide)" class="rounded text-primary-600"/>
                <span class="text-xs text-gray-500">Active</span>
              </label>
              <div class="flex gap-1 flex-shrink-0">
                <button @click="editSlide(slide)" class="p-1.5 rounded hover:bg-gray-100 text-gray-400 hover:text-gray-700 transition-colors">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931z"/></svg>
                </button>
                <button @click="confirmDelete('slide',slide.id)" class="p-1.5 rounded hover:bg-red-50 text-gray-400 hover:text-red-600 transition-colors">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916"/></svg>
                </button>
              </div>
            </div>
            <div v-if="slides.length===0" class="text-center py-12 text-gray-400 bg-white rounded-xl border border-gray-200">No slides yet</div>
          </div>
        </div>

        <!-- ═══ MENU BUILDER ═══════════════════════════════════════ -->
        <div v-else-if="activeSection==='menus'">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div v-for="menuName in ['main','footer']" :key="menuName" class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
              <div class="px-5 py-4 border-b border-gray-100 flex items-center justify-between">
                <h3 class="font-semibold text-gray-900 capitalize">{{ menuName === 'main' ? 'Main Navigation' : 'Footer Links' }}</h3>
                <button @click="saveMenu(menuName)" :disabled="savingMenu===menuName"
                  class="text-sm bg-primary-600 hover:bg-primary-700 text-white px-3 py-1.5 rounded-lg transition-colors disabled:opacity-50">
                  {{ savingMenu===menuName ? 'Saving…' : 'Save' }}
                </button>
              </div>
              <div class="p-5 space-y-2">
                <div v-for="item in menus[menuName]" :key="item.id"
                  draggable="true"
                  @dragstart="menuDragStart($event,menuName,item)"
                  @dragover.prevent
                  @drop="menuDrop($event,menuName,item)"
                  class="border border-gray-200 rounded-lg overflow-hidden">
                  <div class="flex items-center gap-2 p-2 bg-gray-50 cursor-grab">
                    <svg class="w-4 h-4 text-gray-300 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9h16.5m-16.5 6.75h16.5"/></svg>
                    <input v-model="item.label" placeholder="Label" class="flex-1 min-w-0 text-sm border border-gray-200 rounded px-2 py-1 focus:outline-none focus:ring-1 focus:ring-primary-400"/>
                    <input v-model="item.url" placeholder="URL" class="flex-1 min-w-0 text-sm border border-gray-200 rounded px-2 py-1 focus:outline-none focus:ring-1 focus:ring-primary-400"/>
                    <label class="flex items-center gap-1 text-xs text-gray-500 cursor-pointer flex-shrink-0">
                      <input type="checkbox" v-model="item.mega" class="rounded text-primary-600 w-3 h-3"/> Mega
                    </label>
                    <button @click="addSubItem(menuName,item.id)" class="text-xs text-primary-600 hover:text-primary-800 flex-shrink-0">+Sub</button>
                    <button @click="removeMenuItem(menuName,item.id)" class="text-gray-400 hover:text-red-500 flex-shrink-0">
                      <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                  </div>
                  <!-- Sub items -->
                  <div v-if="item.children?.length" class="pl-6 pr-2 py-2 space-y-1.5 bg-white border-t border-gray-100">
                    <div v-for="child in item.children" :key="child.id" class="flex items-center gap-2">
                      <span class="text-gray-300 text-xs">↳</span>
                      <input v-model="child.label" placeholder="Label" class="flex-1 text-xs border border-gray-200 rounded px-2 py-1 focus:outline-none focus:ring-1 focus:ring-primary-400"/>
                      <input v-model="child.url" placeholder="URL" class="flex-1 text-xs border border-gray-200 rounded px-2 py-1 focus:outline-none focus:ring-1 focus:ring-primary-400"/>
                      <button @click="removeSubItem(menuName,item.id,child.id)" class="text-gray-400 hover:text-red-500 flex-shrink-0">
                        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                      </button>
                    </div>
                  </div>
                  <!-- Mega menu items -->
                  <div v-if="item.mega" class="pl-6 pr-2 py-2 bg-blue-50 border-t border-blue-100">
                    <p class="text-xs text-blue-600 mb-1.5 font-medium">Mega Menu Items</p>
                    <div v-for="(mg,mi) in (item.megaItems||[])" :key="mi" class="flex items-center gap-2 mb-1.5">
                      <input v-model="mg.label" placeholder="Label" class="flex-1 text-xs border border-blue-200 rounded px-2 py-1 focus:outline-none"/>
                      <input v-model="mg.url" placeholder="URL" class="flex-1 text-xs border border-blue-200 rounded px-2 py-1 focus:outline-none"/>
                      <input v-model="mg.icon" placeholder="Icon emoji" class="w-16 text-xs border border-blue-200 rounded px-2 py-1 focus:outline-none"/>
                      <button @click="item.megaItems.splice(mi,1)" class="text-blue-400 hover:text-red-500">
                        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                      </button>
                    </div>
                    <button @click="if(!item.megaItems)item.megaItems=[];item.megaItems.push({label:'',url:'',icon:''})" class="text-xs text-blue-600 hover:text-blue-800">+ Add Mega Item</button>
                  </div>
                </div>
                <button @click="addMenuItem(menuName)" class="w-full py-2 border-2 border-dashed border-gray-200 hover:border-primary-300 text-sm text-gray-400 hover:text-primary-600 rounded-lg transition-colors">
                  + Add Item
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- ═══ WIDGET BUILDER ════════════════════════════════════ -->
        <div v-else-if="activeSection==='widgets'">
          <div class="flex gap-5">
            <!-- Available widgets -->
            <div class="w-48 flex-shrink-0">
              <h3 class="text-sm font-semibold text-gray-700 mb-3">Available Widgets</h3>
              <div class="space-y-2">
                <div v-for="w in availableWidgets" :key="w.type"
                  draggable="true"
                  @dragstart="widgetDragStart($event,w)"
                  class="flex items-center gap-2 bg-white border border-gray-200 rounded-lg p-2.5 cursor-grab hover:border-primary-300 hover:bg-primary-50 transition-colors">
                  <span class="text-lg leading-none">{{ w.icon }}</span>
                  <span class="text-sm text-gray-700">{{ w.label }}</span>
                </div>
              </div>
            </div>
            <!-- Footer columns -->
            <div class="flex-1 grid grid-cols-2 xl:grid-cols-4 gap-4">
              <div v-for="col in 4" :key="col" class="min-h-48"
                @dragover.prevent
                @drop="dropWidget($event,col)">
                <h3 class="text-sm font-semibold text-gray-700 mb-2">Column {{ col }}</h3>
                <div class="min-h-32 bg-gray-50 border-2 border-dashed border-gray-200 rounded-xl p-2 space-y-2 transition-colors"
                  :class="widgetDropTarget===col ? 'border-primary-400 bg-primary-50' : ''"
                  @dragenter="widgetDropTarget=col" @dragleave="widgetDropTarget=null">
                  <div v-for="(w,wi) in widgetColumns[col]||[]" :key="wi"
                    class="bg-white border border-gray-200 rounded-lg p-2.5 flex items-start gap-2 group">
                    <span class="text-base leading-none mt-0.5">{{ w.icon }}</span>
                    <div class="flex-1 min-w-0">
                      <div class="text-xs font-medium text-gray-700">{{ w.label }}</div>
                      <input v-if="w.type==='html'" v-model="w.content" placeholder="HTML content…" class="mt-1 w-full text-xs border border-gray-200 rounded px-1.5 py-1 focus:outline-none"/>
                      <input v-else-if="w.type==='image'" v-model="w.imageUrl" placeholder="Image URL…" class="mt-1 w-full text-xs border border-gray-200 rounded px-1.5 py-1 focus:outline-none"/>
                      <input v-else-if="w.type==='about'" v-model="w.text" placeholder="About text…" class="mt-1 w-full text-xs border border-gray-200 rounded px-1.5 py-1 focus:outline-none"/>
                    </div>
                    <button @click="removeWidget(col,wi)" class="opacity-0 group-hover:opacity-100 text-gray-400 hover:text-red-500 transition-all flex-shrink-0">
                      <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg>
                    </button>
                  </div>
                  <div v-if="!(widgetColumns[col]||[]).length" class="flex items-center justify-center h-20 text-xs text-gray-400">
                    Drop widget here
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="mt-4 flex justify-end">
            <button @click="saveWidgets" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">
              Save Widget Layout
            </button>
          </div>
        </div>

        <!-- ═══ HOMEPAGE BUILDER ══════════════════════════════════ -->
        <div v-else-if="activeSection==='homepage-builder'">
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-5">
            <div class="lg:col-span-2 bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
              <div class="px-5 py-4 border-b border-gray-100">
                <h3 class="font-semibold text-gray-900">Section Order & Visibility</h3>
                <p class="text-xs text-gray-400 mt-0.5">Drag to reorder. Toggle to show/hide.</p>
              </div>
              <div class="divide-y divide-gray-50">
                <div v-for="sec in homeSections" :key="sec.id"
                  draggable="true"
                  @dragstart="homeDragStart($event,sec)"
                  @dragover.prevent
                  @drop="homeDrop($event,sec)"
                  class="flex items-center gap-3 px-5 py-3.5 hover:bg-gray-50 cursor-grab transition-colors">
                  <svg class="w-4 h-4 text-gray-300 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3.75 9h16.5m-16.5 6.75h16.5"/></svg>
                  <span class="text-xl leading-none">{{ sec.icon }}</span>
                  <div class="flex-1 min-w-0">
                    <div class="font-medium text-sm text-gray-900">{{ sec.label }}</div>
                    <div class="text-xs text-gray-400">{{ sec.description }}</div>
                  </div>
                  <label class="relative inline-flex items-center cursor-pointer flex-shrink-0">
                    <input type="checkbox" v-model="sec.visible" class="sr-only peer"/>
                    <div class="w-9 h-5 bg-gray-200 peer-focus:ring-2 peer-focus:ring-primary-300 rounded-full peer peer-checked:bg-primary-600 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-4"></div>
                  </label>
                </div>
              </div>
            </div>
            <div class="space-y-4">
              <div class="bg-white rounded-xl border border-gray-200 shadow-sm p-5">
                <h3 class="font-semibold text-gray-900 mb-3">Header Settings</h3>
                <div class="space-y-3">
                  <div>
                    <label class="text-xs font-medium text-gray-600 block mb-1">Hero Background</label>
                    <select v-model="homeSettings.heroStyle" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                      <option value="gradient">Gradient</option><option value="image">Image</option><option value="video">Video</option><option value="solid">Solid Color</option>
                    </select>
                  </div>
                  <div>
                    <label class="text-xs font-medium text-gray-600 block mb-1">Overlay Opacity</label>
                    <input type="range" min="0" max="100" v-model="homeSettings.overlayOpacity" class="w-full accent-primary-600"/>
                    <span class="text-xs text-gray-400">{{ homeSettings.overlayOpacity }}%</span>
                  </div>
                  <div>
                    <label class="text-xs font-medium text-gray-600 block mb-1">Counter Speed (ms)</label>
                    <input type="number" v-model="homeSettings.counterSpeed" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2"/>
                  </div>
                </div>
              </div>
              <button @click="saveHomeSettings" class="w-full bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium py-2.5 rounded-xl transition-colors">
                Save Page Settings
              </button>
            </div>
          </div>
        </div>

        <!-- ═══ TEAM MEMBERS ══════════════════════════════════════ -->
        <div v-else-if="activeSection==='team'">
          <div v-if="teamLoading" class="space-y-2"><div v-for="i in 4" :key="i" class="h-16 bg-gray-100 rounded-xl animate-pulse"/></div>
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            <div v-for="member in team" :key="member.id"
              class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
              <div class="h-40 bg-gray-100 overflow-hidden">
                <img v-if="member.image_url" :src="mediaUrl(member.image_url)" class="w-full h-full object-cover"/>
                <div v-else class="w-full h-full flex items-center justify-center text-4xl text-gray-300">👤</div>
              </div>
              <div class="p-3">
                <div class="font-semibold text-gray-900 text-sm">{{ member.full_name }}</div>
                <div class="text-xs text-gray-500 mt-0.5">{{ member.designation }}</div>
                <div class="flex gap-1 mt-3">
                  <button @click="editTeam(member)" class="flex-1 text-xs text-center py-1 border border-gray-200 rounded-lg hover:bg-gray-50 text-gray-600 transition-colors">Edit</button>
                  <button @click="confirmDelete('team',member.id)" class="flex-1 text-xs text-center py-1 border border-red-200 rounded-lg hover:bg-red-50 text-red-500 transition-colors">Delete</button>
                </div>
              </div>
            </div>
            <div v-if="team.length===0" class="col-span-full text-center py-16 text-gray-400 bg-white rounded-xl border border-gray-200">No team members yet</div>
          </div>
        </div>

        <!-- ═══ SERVICES ══════════════════════════════════════════ -->
        <div v-else-if="activeSection==='services'">
          <ContentTable :items="services" :loading="servicesLoading"
            :cols="[{key:'title',label:'Title'},{key:'category',label:'Category'},{key:'status',label:'Status'}]"
            @edit="editService" @delete="id=>confirmDelete('service',id)" />
        </div>

        <!-- ═══ CONTACT INFO ══════════════════════════════════════ -->
        <div v-else-if="activeSection==='contact'">
          <div class="max-w-2xl space-y-5">
            <div class="bg-white rounded-xl border border-gray-200 shadow-sm p-6 space-y-4">
              <h3 class="font-semibold text-gray-900">Contact Details</h3>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Phone</label><input v-model="contactSettings.phone" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="+256 ..."/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Email</label><input v-model="contactSettings.email" type="email" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="info@..."/></div>
                <div class="sm:col-span-2"><label class="text-xs font-medium text-gray-600 block mb-1">Address</label><textarea v-model="contactSettings.address" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="2" placeholder="Physical address"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Working Hours</label><input v-model="contactSettings.hours" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="Mon–Fri, 8am–5pm"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Google Maps Embed URL</label><input v-model="contactSettings.mapUrl" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="https://maps.google.com/..."/></div>
              </div>
            </div>
            <div class="bg-white rounded-xl border border-gray-200 shadow-sm p-6 space-y-4">
              <h3 class="font-semibold text-gray-900">Social Media</h3>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div v-for="s in ['facebook','twitter','linkedin','instagram','youtube']" :key="s">
                  <label class="text-xs font-medium text-gray-600 block mb-1 capitalize">{{ s }}</label>
                  <input v-model="contactSettings.social[s]" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" :placeholder="'https://'+s+'.com/...'"/>
                </div>
              </div>
            </div>
            <button @click="saveContactSettings" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-6 py-2.5 rounded-xl transition-colors">Save Contact Info</button>
          </div>
        </div>

        <!-- ═══ SUPPORT TICKETS ═══════════════════════════════════ -->
        <div v-else-if="activeSection==='tickets'">
          <ContentTable :items="tickets" :loading="ticketsLoading"
            :cols="[{key:'subject',label:'Subject'},{key:'department',label:'Department'},{key:'status',label:'Status'},{key:'created_at',label:'Date'}]"
            @edit="editTicket" @delete="id=>confirmDelete('ticket',id)" />
        </div>

        <!-- ═══ KNOWLEDGEBASE ═════════════════════════════════════ -->
        <div v-else-if="activeSection==='knowledgebase'">
          <ContentTable :items="kbArticles" :loading="kbLoading"
            :cols="[{key:'title',label:'Title'},{key:'topic',label:'Topic'},{key:'status',label:'Status'},{key:'created_at',label:'Date'}]"
            @edit="editKB" @delete="id=>confirmDelete('kb',id)" />
        </div>

        <!-- ═══ FUN FACTS ════════════════════════════════════════ -->
        <div v-else-if="activeSection==='fun-facts'">
          <ContentTable :items="funFacts" :loading="funFactsLoading"
            :cols="[{key:'label',label:'Label'},{key:'value',label:'Value'},{key:'icon',label:'Icon'},{key:'sort_order',label:'Order'}]"
            @edit="editFunFact" @delete="id=>confirmDelete('fun-fact',id)" />
        </div>

        <!-- ═══ FAQs ══════════════════════════════════════════════ -->
        <div v-else-if="activeSection==='faqs'">
          <div class="flex gap-3 mb-4">
            <select v-model="faqCategory" @change="loadFAQs()" class="text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white text-gray-700">
              <option value="">All Categories</option>
              <option value="general">General</option>
              <option value="licensing">Licensing</option>
              <option value="membership">Membership</option>
              <option value="payments">Payments</option>
            </select>
          </div>
          <ContentTable :items="faqs" :loading="faqsLoading"
            :cols="[{key:'question',label:'Question'},{key:'category',label:'Category'},{key:'sort_order',label:'Order'}]"
            @edit="editFAQ" @delete="id=>confirmDelete('faq',id)" />
        </div>

        <!-- ═══ RESOURCES ═════════════════════════════════════════ -->
        <div v-else-if="activeSection==='resources'">
          <div class="flex gap-3 mb-4">
            <select v-model="resourceCategory" @change="loadResources()" class="text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white text-gray-700">
              <option value="">All Categories</option>
              <option value="guidelines">Guidelines & Standards</option>
              <option value="reports">Reports</option>
              <option value="speeches">Speeches</option>
              <option value="press-release">Press Release</option>
              <option value="sports-rules">Sports Rules & Laws</option>
              <option value="downloads">Important Downloads</option>
            </select>
          </div>
          <ContentTable :items="resources" :loading="resourcesLoading"
            :cols="[{key:'title',label:'Title'},{key:'category',label:'Category'},{key:'file_url',label:'File'}]"
            @edit="editResource" @delete="id=>confirmDelete('resource',id)" />
        </div>

        <!-- ═══ FACILITIES ════════════════════════════════════════ -->
        <div v-else-if="activeSection==='facilities'">
          <ContentTable :items="facilities" :loading="facilitiesLoading"
            :cols="[{key:'name',label:'Name'},{key:'slug',label:'Slug'},{key:'sort_order',label:'Order'}]"
            @edit="editFacility" @delete="id=>confirmDelete('facility',id)" />
        </div>

        <!-- ═══ ASSOCIATIONS ══════════════════════════════════════ -->
        <div v-else-if="activeSection==='associations'">
          <ContentTable :items="associations" :loading="associationsLoading"
            :cols="[{key:'name',label:'Name'},{key:'website_url',label:'Website'},{key:'sort_order',label:'Order'}]"
            @edit="editAssociation" @delete="id=>confirmDelete('association',id)" />
        </div>

        <!-- ═══ INVEST WITH US ════════════════════════════════════ -->
        <div v-else-if="activeSection==='invest'">
          <ContentTable :items="investItems" :loading="investLoading"
            :cols="[{key:'title',label:'Title'},{key:'subtitle',label:'Subtitle'},{key:'sort_order',label:'Order'}]"
            @edit="editInvest" @delete="id=>confirmDelete('invest',id)" />
        </div>

        <!-- ═══ APPEARANCE ════════════════════════════════════════ -->
        <div v-else-if="activeSection==='appearance'">
          <div class="flex gap-1 border-b border-gray-200 mb-5">
            <button v-for="t in appearanceTabs" :key="t.id" @click="appearanceTab=t.id"
              :class="appearanceTab===t.id ? 'border-b-2 border-primary-600 text-primary-700 font-medium' : 'text-gray-500 hover:text-gray-700 border-b-2 border-transparent'"
              class="px-4 py-2.5 text-sm transition-colors">{{ t.label }}</button>
          </div>
          <!-- Topbar -->
          <div v-if="appearanceTab==='topbar'" class="max-w-xl space-y-4 bg-white rounded-xl border border-gray-200 p-6">
            <div class="grid grid-cols-2 gap-4">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Background Color</label><input type="color" v-model="appearance.topbar.bg" class="w-full h-9 border border-gray-200 rounded-lg cursor-pointer"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Text Color</label><input type="color" v-model="appearance.topbar.text" class="w-full h-9 border border-gray-200 rounded-lg cursor-pointer"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Phone</label><input type="text" v-model="appearance.topbar.phone" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Email</label><input type="email" v-model="appearance.topbar.email" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
            </div>
            <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="appearance.topbar.show" class="rounded text-primary-600"/> Show topbar</label>
            <button @click="saveAppearance" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">Save</button>
          </div>
          <!-- Navbar -->
          <div v-else-if="appearanceTab==='navbar'" class="max-w-xl space-y-4 bg-white rounded-xl border border-gray-200 p-6">
            <div class="grid grid-cols-2 gap-4">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Background Color</label><input type="color" v-model="appearance.navbar.bg" class="w-full h-9 border border-gray-200 rounded-lg cursor-pointer"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Text Color</label><input type="color" v-model="appearance.navbar.text" class="w-full h-9 border border-gray-200 rounded-lg cursor-pointer"/></div>
              <div class="col-span-2"><label class="text-xs font-medium text-gray-600 block mb-1">Logo Alignment</label>
                <select v-model="appearance.navbar.logoAlign" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                  <option value="left">Left</option><option value="center">Center</option>
                </select>
              </div>
            </div>
            <div class="space-y-2">
              <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="appearance.navbar.sticky" class="rounded text-primary-600"/> Sticky navbar</label>
              <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="appearance.navbar.transparentOnHero" class="rounded text-primary-600"/> Transparent on hero section</label>
            </div>
            <button @click="saveAppearance" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">Save</button>
          </div>
          <!-- Footer -->
          <div v-else-if="appearanceTab==='footer'" class="max-w-xl space-y-4 bg-white rounded-xl border border-gray-200 p-6">
            <div class="grid grid-cols-2 gap-4">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Background Color</label><input type="color" v-model="appearance.footer.bg" class="w-full h-9 border border-gray-200 rounded-lg cursor-pointer"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Text Color</label><input type="color" v-model="appearance.footer.text" class="w-full h-9 border border-gray-200 rounded-lg cursor-pointer"/></div>
              <div class="col-span-2"><label class="text-xs font-medium text-gray-600 block mb-1">Copyright Text</label><input type="text" v-model="appearance.footer.copyright" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Footer Columns</label>
                <select v-model="appearance.footer.columns" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                  <option :value="1">1</option><option :value="2">2</option><option :value="3">3</option><option :value="4">4</option>
                </select>
              </div>
            </div>
            <button @click="saveAppearance" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">Save</button>
          </div>
          <!-- Colors -->
          <div v-else-if="appearanceTab==='colors'" class="max-w-xl space-y-4 bg-white rounded-xl border border-gray-200 p-6">
            <p class="text-sm text-gray-500">These override the default theme colors for the public site.</p>
            <div class="grid grid-cols-2 gap-4">
              <div v-for="c in [{k:'primary',l:'Primary'},{k:'secondary',l:'Secondary'},{k:'accent',l:'Accent'},{k:'dark',l:'Dark'}]" :key="c.k">
                <label class="text-xs font-medium text-gray-600 block mb-1">{{ c.l }}</label>
                <div class="flex items-center gap-2"><input type="color" v-model="appearance.colors[c.k]" class="w-10 h-9 border border-gray-200 rounded-lg cursor-pointer flex-shrink-0"/><input type="text" v-model="appearance.colors[c.k]" class="flex-1 text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 font-mono"/></div>
              </div>
            </div>
            <button @click="saveAppearance" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">Save Colors</button>
          </div>
          <!-- Typography -->
          <div v-else-if="appearanceTab==='typography'" class="max-w-xl space-y-4 bg-white rounded-xl border border-gray-200 p-6">
            <div class="grid grid-cols-1 gap-4">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Heading Font</label>
                <select v-model="appearance.typography.heading" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                  <option v-for="f in googleFonts" :key="f" :value="f">{{ f }}</option>
                </select>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Body Font</label>
                <select v-model="appearance.typography.body" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                  <option v-for="f in googleFonts" :key="f" :value="f">{{ f }}</option>
                </select>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Base Font Size</label>
                <select v-model="appearance.typography.size" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                  <option value="14px">Small (14px)</option><option value="16px">Normal (16px)</option><option value="18px">Large (18px)</option>
                </select>
              </div>
            </div>
            <button @click="saveAppearance" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">Save Typography</button>
          </div>
        </div>

        <!-- ═══ SITEMAP ══════════════════════════════════════════ -->
        <div v-else-if="activeSection==='sitemap'">
          <div class="mb-4 flex items-center justify-between">
            <p class="text-sm text-gray-500">All public pages and dynamic content slugs. Click any URL to open it.</p>
            <button @click="loadSitemap" class="text-xs bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 py-1.5 rounded-lg transition-colors">Refresh</button>
          </div>
          <div v-if="sitemapLoading" class="space-y-2">
            <div v-for="i in 6" :key="i" class="h-12 bg-gray-100 rounded-lg animate-pulse"/>
          </div>
          <div v-else class="space-y-5">
            <div v-for="group in sitemap" :key="group.id" class="bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
              <div class="px-5 py-3 border-b border-gray-100 bg-gray-50">
                <h3 class="font-semibold text-gray-800 text-sm">{{ group.label }}</h3>
                <p v-if="group.description" class="text-xs text-gray-500 mt-0.5">{{ group.description }}</p>
              </div>
              <table class="w-full text-sm">
                <thead class="bg-gray-50/40 border-b border-gray-100">
                  <tr>
                    <th class="text-left px-5 py-2.5 text-xs font-medium text-gray-500">Title</th>
                    <th class="text-left px-5 py-2.5 text-xs font-medium text-gray-500">Slug</th>
                    <th class="text-left px-5 py-2.5 text-xs font-medium text-gray-500">URL</th>
                    <th class="px-5 py-2.5 w-16"></th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-50">
                  <tr v-for="row in group.entries" :key="row.url" class="hover:bg-gray-50">
                    <td class="px-5 py-2.5 text-gray-800 font-medium">{{ row.title || '—' }}</td>
                    <td class="px-5 py-2.5 text-gray-500 font-mono text-xs">{{ row.slug || '—' }}</td>
                    <td class="px-5 py-2.5 text-gray-500 font-mono text-xs truncate max-w-md">{{ row.url }}</td>
                    <td class="px-5 py-2.5 text-right">
                      <a :href="row.url" target="_blank" class="text-primary-600 hover:text-primary-800 text-xs font-medium">Open ↗</a>
                    </td>
                  </tr>
                  <tr v-if="!group.entries.length">
                    <td colspan="4" class="px-5 py-6 text-center text-gray-400 text-sm">No entries</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- ═══ GENERAL SETTINGS ══════════════════════════════════ -->
        <div v-else-if="activeSection==='settings'">
          <div class="flex gap-1 border-b border-gray-200 mb-5 flex-wrap">
            <button v-for="t in settingsTabs" :key="t.id" @click="settingsTab=t.id"
              :class="settingsTab===t.id ? 'border-b-2 border-primary-600 text-primary-700 font-medium' : 'text-gray-500 hover:text-gray-700 border-b-2 border-transparent'"
              class="px-4 py-2.5 text-sm transition-colors">{{ t.label }}</button>
          </div>
          <!-- Site Identity -->
          <div v-if="settingsTab==='identity'" class="max-w-xl bg-white rounded-xl border border-gray-200 p-6 space-y-4">
            <div><label class="text-xs font-medium text-gray-600 block mb-1">Site Name</label><input type="text" v-model="siteSettings.name" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
            <div><label class="text-xs font-medium text-gray-600 block mb-1">Tagline</label><input type="text" v-model="siteSettings.tagline" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
            <div><label class="text-xs font-medium text-gray-600 block mb-1">Logo URL</label><input type="text" v-model="siteSettings.logoUrl" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="/main-logo.png"/></div>
            <div><label class="text-xs font-medium text-gray-600 block mb-1">Favicon URL</label><input type="text" v-model="siteSettings.faviconUrl" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="/favicon.png"/></div>
            <button @click="saveSettings" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">Save</button>
          </div>
          <!-- SEO -->
          <div v-else-if="settingsTab==='seo'" class="max-w-xl bg-white rounded-xl border border-gray-200 p-6 space-y-4">
            <div><label class="text-xs font-medium text-gray-600 block mb-1">Title Template</label><input type="text" v-model="siteSettings.titleTemplate" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="%s | NCS Uganda"/></div>
            <div><label class="text-xs font-medium text-gray-600 block mb-1">Default Meta Description</label><textarea v-model="siteSettings.metaDescription" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="3"/></div>
            <div><label class="text-xs font-medium text-gray-600 block mb-1">Robots.txt</label>
              <select v-model="siteSettings.robots" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                <option value="index,follow">Index, Follow</option><option value="noindex,nofollow">NoIndex, NoFollow</option>
              </select>
            </div>
            <button @click="saveSettings" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">Save</button>
          </div>
          <!-- SMTP -->
          <div v-else-if="settingsTab==='smtp'" class="max-w-xl bg-white rounded-xl border border-gray-200 p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div class="col-span-2"><label class="text-xs font-medium text-gray-600 block mb-1">SMTP Host</label><input type="text" v-model="siteSettings.smtpHost" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Port</label><input type="number" v-model="siteSettings.smtpPort" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div class="flex items-end pb-1"><label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="siteSettings.smtpSSL" class="rounded text-primary-600"/> Use SSL/TLS</label></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Username</label><input type="text" v-model="siteSettings.smtpUser" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Password</label><input type="password" v-model="siteSettings.smtpPass" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">From Email</label><input type="email" v-model="siteSettings.fromEmail" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">From Name</label><input type="text" v-model="siteSettings.fromName" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
            </div>
            <button @click="saveSettings" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">Save SMTP</button>
          </div>
          <!-- Custom CSS -->
          <div v-else-if="settingsTab==='css'" class="space-y-4">
            <div class="bg-white rounded-xl border border-gray-200 p-4">
              <label class="text-xs font-medium text-gray-600 block mb-2">Custom CSS</label>
              <textarea v-model="siteSettings.customCss" class="w-full font-mono text-xs bg-gray-900 text-green-300 border-0 rounded-lg p-4 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="14" placeholder="/* Your custom CSS here */"/>
            </div>
            <div class="bg-white rounded-xl border border-gray-200 p-4">
              <label class="text-xs font-medium text-gray-600 block mb-2">Custom JavaScript</label>
              <textarea v-model="siteSettings.customJs" class="w-full font-mono text-xs bg-gray-900 text-yellow-300 border-0 rounded-lg p-4 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="14" placeholder="// Your custom JS here"/>
            </div>
            <button @click="saveSettings" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">Save Scripts</button>
          </div>
          <!-- Modules -->
          <div v-else-if="settingsTab==='modules'" class="max-w-xl bg-white rounded-xl border border-gray-200 p-6">
            <p class="text-sm text-gray-500 mb-4">Enable or disable site modules globally.</p>
            <div class="space-y-3">
              <div v-for="mod in siteModules" :key="mod.key" class="flex items-center justify-between py-2 border-b border-gray-50 last:border-0">
                <div>
                  <div class="text-sm font-medium text-gray-800">{{ mod.label }}</div>
                  <div class="text-xs text-gray-400">{{ mod.description }}</div>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" v-model="siteSettings.modules[mod.key]" class="sr-only peer"/>
                  <div class="w-9 h-5 bg-gray-200 peer-focus:ring-2 peer-focus:ring-primary-300 rounded-full peer peer-checked:bg-primary-600 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:after:translate-x-4"></div>
                </label>
              </div>
            </div>
            <button @click="saveSettings" class="mt-5 bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg transition-colors">Save Modules</button>
          </div>
        </div>

      </div><!-- /content area -->
    </div><!-- /flex wrapper -->

    <!-- ═══ SHARED MODAL ══════════════════════════════════════════ -->
    <Teleport to="body">
      <div v-if="showModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4 overflow-y-auto">
        <div class="bg-white rounded-2xl shadow-2xl w-full max-w-2xl my-4" @click.stop>
          <div class="flex items-center justify-between px-6 py-5 border-b border-gray-100">
            <h2 class="text-lg font-bold text-gray-900">{{ editingId ? 'Edit' : 'New' }} {{ currentNavItem?.singular||'Item' }}</h2>
            <button @click="showModal=false" class="text-gray-400 hover:text-gray-600"><svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/></svg></button>
          </div>
          <div class="px-6 py-5 space-y-4 max-h-[70vh] overflow-y-auto">

            <!-- Posts / Blogs / Case Studies / Pages form -->
            <template v-if="['posts','blogs','case-studies','pages'].includes(activeSection)">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Title *</label><input v-model="form.title" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div class="grid grid-cols-2 gap-4">
                <div v-if="activeSection==='posts'"><label class="text-xs font-medium text-gray-600 block mb-1">Category</label>
                  <select v-model="form.category" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                    <option value="news">News</option><option value="announcement">Announcement</option><option value="project">Project</option><option value="blog">Blog</option><option value="case_study">Case Study</option>
                  </select>
                </div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Status</label>
                  <select v-model="form.status" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                    <option value="draft">Draft</option><option value="published">Published</option>
                  </select>
                </div>
                <div v-if="activeSection!=='pages'"><label class="text-xs font-medium text-gray-600 block mb-1">Visibility</label>
                  <select v-model="form.visibility" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                    <option value="public">Public</option><option value="private">Private</option>
                  </select>
                </div>
              </div>
              <div v-if="['posts','blogs','case-studies'].includes(activeSection)">
                <label class="text-xs font-medium text-gray-600 block mb-1">Tags (comma-separated)</label>
                <input v-model="form.tags" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="tag1, tag2"/>
              </div>
              <div v-if="activeSection==='blogs'">
                <label class="text-xs font-medium text-gray-600 block mb-1">Video URL (optional)</label>
                <input v-model="form.video_url" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="https://youtube.com/..."/>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Cover Image URL <span class="text-gray-400 font-normal">(1920×1280 recommended)</span></label><input v-model="form.cover_image_url" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="https://..."/>
                <img v-if="form.cover_image_url" :src="mediaUrl(form.cover_image_url)" class="mt-2 h-24 object-cover rounded-lg border border-gray-200"/>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Excerpt</label><textarea v-model="form.excerpt" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="2"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Content</label><textarea v-model="form.content" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none font-mono" rows="8"/></div>
              <div class="border-t border-gray-100 pt-4">
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3">SEO</p>
                <div class="space-y-3">
                  <div><label class="text-xs font-medium text-gray-600 block mb-1">Meta Title</label><input v-model="form.meta_title" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                  <div><label class="text-xs font-medium text-gray-600 block mb-1">Meta Description</label><textarea v-model="form.meta_description" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="2"/></div>
                </div>
              </div>
            </template>

            <!-- Events form -->
            <template v-else-if="activeSection==='events'">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Title *</label><input v-model="form.title" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Event Date</label><input v-model="form.event_date_str" type="datetime-local" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">End Date</label><input v-model="form.end_date_str" type="datetime-local" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Location</label><input v-model="form.location" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Status</label><select v-model="form.status" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2"><option value="draft">Draft</option><option value="published">Published</option></select></div>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Cover Image URL</label><input v-model="form.cover_image_url" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Description</label><textarea v-model="form.description" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="6"/></div>
            </template>

            <!-- Careers form -->
            <template v-else-if="activeSection==='careers'">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Job Title *</label><input v-model="form.title" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Department</label><input v-model="form.department" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Location</label><input v-model="form.location" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Category</label><select v-model="form.category" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2"><option value="jobs">Jobs</option><option value="tenders">Tenders</option><option value="internships">Internship Opportunity</option></select></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Job Type</label><select v-model="form.job_type" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2"><option value="full_time">Full Time</option><option value="part_time">Part Time</option><option value="contract">Contract</option><option value="internship">Internship</option></select></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Salary Range</label><input v-model="form.salary_range" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="e.g. UGX 2M – 3M"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Deadline</label><input v-model="form.deadline_str" type="date" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Status</label><select v-model="form.status" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2"><option value="draft">Draft</option><option value="published">Published</option></select></div>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Description *</label><textarea v-model="form.description" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="5"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Requirements</label><textarea v-model="form.requirements" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="4"/></div>
            </template>

            <!-- Slides form -->
            <template v-else-if="activeSection==='slides'">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Title *</label><input v-model="form.title" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Subtitle</label><input v-model="form.subtitle" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Description</label><textarea v-model="form.description" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="3"/></div>
              <div>
                <label class="text-xs font-medium text-gray-600 block mb-1">Background Image</label>
                <DropzoneUpload v-model="form.image_url" label="Drop slide image or click to upload" hint="1920×1080 recommended" preview-class="h-36"/>
              </div>
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Button Text</label><input v-model="form.button_text" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Button URL</label><input v-model="form.button_url" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Sort Order</label><input v-model.number="form.sort_order" type="number" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div class="flex items-end pb-1"><label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="form.is_active" class="rounded text-primary-600"/> Active</label></div>
              </div>
            </template>

            <!-- Team form -->
            <template v-else-if="activeSection==='team'">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Full Name *</label><input v-model="form.full_name" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Designation</label><input v-model="form.designation" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="e.g. General Secretary"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Photo URL <span class="text-gray-400 font-normal">(1:1 ratio recommended)</span></label><input v-model="form.image_url" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/>
                <img v-if="form.image_url" :src="mediaUrl(form.image_url)" class="mt-2 w-20 h-20 object-cover rounded-full border-2 border-gray-200"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Bio</label><textarea v-model="form.bio" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="3"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Sort Order</label><input v-model.number="form.sort_order" type="number" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
            </template>

            <!-- Services form -->
            <template v-else-if="activeSection==='services'">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Service Title *</label><input v-model="form.title" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Category</label><input v-model="form.category" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Icon (emoji)</label><input v-model="form.icon" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400" placeholder="🏆"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Status</label><select v-model="form.status" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2"><option value="draft">Draft</option><option value="published">Published</option></select></div>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Image URL</label><input v-model="form.image_url" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Description *</label><textarea v-model="form.description" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="4"/></div>
            </template>

            <!-- Tickets form -->
            <template v-else-if="activeSection==='tickets'">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Subject *</label><input v-model="form.subject" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Department</label>
                  <select v-model="form.department" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                    <option value="general">General</option><option value="technical">Technical</option><option value="billing">Billing</option><option value="compliance">Compliance</option>
                  </select>
                </div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Priority</label>
                  <select v-model="form.priority" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                    <option value="low">Low</option><option value="normal">Normal</option><option value="high">High</option><option value="urgent">Urgent</option>
                  </select>
                </div>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Message *</label><textarea v-model="form.message" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none" rows="5"/></div>
            </template>

            <!-- Knowledgebase form -->
            <template v-else-if="activeSection==='knowledgebase'">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Article Title *</label><input v-model="form.title" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Topic</label><input v-model="form.topic" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Status</label><select v-model="form.status" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2"><option value="draft">Draft</option><option value="published">Published</option></select></div>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Content *</label><textarea v-model="form.content" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none font-mono" rows="10"/></div>
            </template>

            <!-- Fun Fact form -->
            <template v-else-if="activeSection==='fun-facts'">
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Label *</label><input v-model="form.label" type="text" placeholder="e.g. Athletes Registered" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Value *</label><input v-model="form.value" type="text" placeholder="e.g. 5,000+" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Icon (icofont class)</label><input v-model="form.icon" type="text" placeholder="icofont-users" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Sort Order</label><input v-model.number="form.sort_order" type="number" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              </div>
              <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="form.is_active" class="rounded text-primary-600"/> Active</label>
            </template>

            <!-- FAQ form -->
            <template v-else-if="activeSection==='faqs'">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Question *</label><input v-model="form.question" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Answer *</label><textarea v-model="form.answer" rows="4" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none"/></div>
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Category</label>
                  <select v-model="form.category" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                    <option value="general">General</option><option value="licensing">Licensing</option>
                    <option value="membership">Membership</option><option value="payments">Payments</option>
                  </select>
                </div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Sort Order</label><input v-model.number="form.sort_order" type="number" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              </div>
              <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="form.is_active" class="rounded text-primary-600"/> Active</label>
            </template>

            <!-- Resource form -->
            <template v-else-if="activeSection==='resources'">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Title *</label><input v-model="form.title" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Category</label>
                  <select v-model="form.category" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2">
                    <option value="guidelines">Guidelines & Standards</option><option value="reports">Reports</option>
                    <option value="speeches">Speeches</option><option value="press-release">Press Release</option>
                    <option value="sports-rules">Sports Rules & Laws</option><option value="downloads">Important Downloads</option>
                  </select>
                </div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Sort Order</label><input v-model.number="form.sort_order" type="number" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              </div>
              <div>
                <label class="text-xs font-medium text-gray-600 block mb-1">PDF File</label>
                <DropzoneUpload v-model="form.file_url" accept=".pdf,application/pdf" label="Drop PDF or click to upload" hint="Max 32MB · PDF files only"/>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Description</label><textarea v-model="form.description" rows="2" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none"/></div>
              <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="form.is_active" class="rounded text-primary-600"/> Active</label>
            </template>

            <!-- Facility form -->
            <template v-else-if="activeSection==='facilities'">
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Name *</label><input v-model="form.name" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Sort Order</label><input v-model.number="form.sort_order" type="number" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              </div>
              <div>
                <label class="text-xs font-medium text-gray-600 block mb-1">Image</label>
                <DropzoneUpload v-model="form.image_url" label="Drop facility image or click to upload" hint="Landscape orientation recommended" preview-class="h-32"/>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Description</label><textarea v-model="form.description" rows="3" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none"/></div>
              <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="form.is_active" class="rounded text-primary-600"/> Active</label>
            </template>

            <!-- Association form -->
            <template v-else-if="activeSection==='associations'">
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Name *</label><input v-model="form.name" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Website URL</label><input v-model="form.website_url" type="url" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Sort Order</label><input v-model.number="form.sort_order" type="number" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              </div>
              <div>
                <label class="text-xs font-medium text-gray-600 block mb-1">Logo</label>
                <DropzoneUpload v-model="form.logo_url" label="Drop association logo or click to upload" hint="Square format recommended" preview-class="h-28 object-contain bg-gray-50"/>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Description</label><textarea v-model="form.description" rows="3" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none"/></div>
              <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="form.is_active" class="rounded text-primary-600"/> Active</label>
            </template>

            <!-- Invest form -->
            <template v-else-if="activeSection==='invest'">
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Title *</label><input v-model="form.title" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Subtitle</label><input v-model="form.subtitle" type="text" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
              <div>
                <label class="text-xs font-medium text-gray-600 block mb-1">Image</label>
                <DropzoneUpload v-model="form.image_url" label="Drop image or click to upload" preview-class="h-32"/>
              </div>
              <div><label class="text-xs font-medium text-gray-600 block mb-1">Content</label><textarea v-model="form.content" rows="5" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400 resize-none"/></div>
              <div class="grid grid-cols-2 gap-4">
                <div><label class="text-xs font-medium text-gray-600 block mb-1">Sort Order</label><input v-model.number="form.sort_order" type="number" class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary-400"/></div>
                <div class="flex items-end pb-1"><label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer"><input type="checkbox" v-model="form.is_active" class="rounded text-primary-600"/> Active</label></div>
              </div>
            </template>

            <p v-if="formError" class="text-sm text-red-500 bg-red-50 border border-red-200 rounded-lg px-3 py-2">{{ formError }}</p>
          </div>
          <div class="px-6 py-4 border-t border-gray-100 flex justify-end gap-3">
            <button @click="showModal=false" class="px-4 py-2 text-sm text-gray-600 hover:text-gray-800 border border-gray-200 rounded-lg">Cancel</button>
            <button @click="saveItem" :disabled="saving" class="bg-primary-600 hover:bg-primary-700 text-white text-sm font-medium px-5 py-2 rounded-lg disabled:opacity-50 transition-colors">
              {{ saving ? 'Saving…' : (editingId ? 'Update' : 'Create') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete confirm -->
    <Teleport to="body">
      <div v-if="deleteTarget" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
        <div class="bg-white rounded-2xl p-6 max-w-sm w-full shadow-2xl">
          <h3 class="font-semibold text-gray-900 mb-2">Confirm Delete</h3>
          <p class="text-sm text-gray-500 mb-6">This action cannot be undone.</p>
          <div class="flex justify-end gap-3">
            <button @click="deleteTarget=null" class="px-4 py-2 text-sm border border-gray-200 rounded-lg text-gray-600 hover:bg-gray-50">Cancel</button>
            <button @click="doDelete" :disabled="deleting" class="bg-red-600 hover:bg-red-700 text-white text-sm font-medium px-4 py-2 rounded-lg disabled:opacity-50 transition-colors">{{ deleting ? 'Deleting…' : 'Delete' }}</button>
          </div>
        </div>
      </div>
    </Teleport>

  </LayoutDefault>
</template>

<script setup>
import { ref, reactive, computed, onMounted, defineComponent, h } from 'vue'
import LayoutDefault from '@/components/layout/LayoutDefault.vue'
import DropzoneUpload from '@/components/ui/DropzoneUpload.vue'
import * as cmsApi from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

// ── ContentTable render component ─────────────────────────────────
const ContentTable = defineComponent({
  name: 'ContentTable',
  props: { items: Array, loading: Boolean, cols: Array },
  emits: ['edit','delete'],
  setup(props, { emit }) {
    return () => {
      if (props.loading) {
        return h('div', { class:'space-y-2' },
          [1,2,3,4,5].map(i => h('div', { key:i, class:'h-14 bg-gray-100 rounded-lg animate-pulse' })))
      }
      return h('div', { class:'bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden' }, [
        h('table', { class:'w-full text-sm' }, [
          h('thead', { class:'bg-gray-50 border-b border-gray-100' },
            h('tr', {}, [
              ...(props.cols||[]).map((col,ci) =>
                h('th', { key:ci, class:'text-left px-4 py-3 font-medium text-gray-500 '+(ci>0?'hidden md:table-cell':'') }, col.label)
              ),
              h('th', { class:'px-4 py-3 w-20' })
            ])
          ),
          h('tbody', { class:'divide-y divide-gray-50' }, [
            ...(props.items||[]).map(item =>
              h('tr', { key:item.id, class:'hover:bg-gray-50' }, [
                ...(props.cols||[]).map((col,ci) =>
                  h('td', { key:ci, class:'px-4 py-3 '+(ci>0?'text-gray-500 hidden md:table-cell':'font-medium text-gray-900') }, [
                    col.key==='status' ? h('span', { class:'inline-flex px-2 py-0.5 rounded-full text-xs font-medium capitalize '+(item.status==='published'?'bg-green-100 text-green-700':'bg-yellow-100 text-yellow-700') }, item.status||'—') :
                    (col.key==='created_at'||col.key==='published_at'||col.key==='event_date') ? (item[col.key] ? new Date(item[col.key]).toLocaleDateString('en-UG',{day:'numeric',month:'short',year:'numeric'}) : '—') :
                    col.key==='category' ? h('span', { class:'capitalize' }, (item.category||'').replace(/_/g,' ')||'—') :
                    String(item[col.key]||'—')
                  ])
                ),
                h('td', { class:'px-4 py-3' },
                  h('div', { class:'flex items-center gap-1 justify-end' }, [
                    h('button', { onClick:()=>emit('edit',item), class:'p-1.5 rounded hover:bg-gray-100 text-gray-400 hover:text-gray-700 transition-colors' },
                      h('svg', { class:'w-4 h-4', fill:'none', viewBox:'0 0 24 24', 'stroke-width':'1.5', stroke:'currentColor' },
                        h('path', { 'stroke-linecap':'round', 'stroke-linejoin':'round', d:'M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931z' })
                      )
                    ),
                    h('button', { onClick:()=>emit('delete',item.id), class:'p-1.5 rounded hover:bg-red-50 text-gray-400 hover:text-red-600 transition-colors' },
                      h('svg', { class:'w-4 h-4', fill:'none', viewBox:'0 0 24 24', 'stroke-width':'1.5', stroke:'currentColor' },
                        h('path', { 'stroke-linecap':'round', 'stroke-linejoin':'round', d:'M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916' })
                      )
                    ),
                  ])
                ),
              ])
            ),
            (props.items||[]).length===0 ? h('tr', {}, h('td', { colspan:((props.cols||[]).length+1), class:'px-4 py-12 text-center text-gray-400 text-sm' }, 'No items found')) : null,
          ])
        ])
      ])
    }
  }
})

// ── Nav definition ─────────────────────────────────────────────────
const navGroups = [
  { id:'content', label:'Content', items:[
    { id:'posts',       label:'Posts & Articles', singular:'Post',         icon:'📰', description:'News, announcements and articles' },
    { id:'blogs',       label:'Blogs',            singular:'Blog Post',    icon:'✍️', description:'Long-form blog content' },
    { id:'events',      label:'Events',           singular:'Event',        icon:'📅', description:'Upcoming and past events' },
    { id:'case-studies',label:'Case Studies',     singular:'Case Study',   icon:'📁', description:'Project case studies' },
    { id:'pages',       label:'Static Pages',     singular:'Page',         icon:'📄', description:'About, Privacy, Terms pages' },
    { id:'careers',       label:'Careers / Jobs',  singular:'Job',          icon:'💼', description:'Jobs, Tenders & Internships' },
    { id:'fun-facts',     label:'Fun Facts',        singular:'Fun Fact',     icon:'🔢', description:'Statistics displayed on homepage' },
    { id:'faqs',          label:'FAQs',             singular:'FAQ',          icon:'❓', description:'Frequently asked questions' },
    { id:'resources',     label:'Resource Centre',  singular:'Resource',     icon:'📥', description:'Guidelines, Reports, PDFs for download' },
    { id:'facilities',    label:'Facilities',       singular:'Facility',     icon:'🏟', description:'Sports facilities and venues' },
    { id:'associations',  label:'Associations',     singular:'Association',  icon:'🤝', description:'Sports associations and federations' },
    { id:'invest',        label:'Invest with Us',   singular:'Invest Item',  icon:'💰', description:'Investment opportunities content' },
  ]},
  { id:'homepage', label:'Homepage', items:[
    { id:'slides',            label:'Hero Slides',      singular:'Slide',    icon:'🖼', description:'Manage homepage carousel' },
    { id:'homepage-builder',  label:'Page Builder',     singular:null,       icon:'🏗', description:'Configure homepage sections' },
  ]},
  { id:'navigation', label:'Navigation', items:[
    { id:'menus',   label:'Menu Builder',   singular:null, icon:'☰', description:'Drag & drop menu management' },
    { id:'widgets', label:'Widget Builder', singular:null, icon:'🧩', description:'Footer widget areas' },
  ]},
  { id:'people', label:'People & Services', items:[
    { id:'team',     label:'Team Members', singular:'Member',  icon:'👥', description:'Staff and leadership' },
    { id:'services', label:'Services',     singular:'Service', icon:'🔧', description:'Services offered' },
  ]},
  { id:'support', label:'Contact & Support', items:[
    { id:'contact',       label:'Contact Info',    singular:null,     icon:'📞', description:'Phone, address, map' },
    { id:'tickets',       label:'Support Tickets', singular:'Ticket', icon:'🎫', description:'Customer support tickets' },
    { id:'knowledgebase', label:'Knowledgebase',   singular:'Article',icon:'📖', description:'Help articles' },
  ]},
  { id:'appearance', label:'Appearance', items:[
    { id:'appearance', label:'Appearance', singular:null, icon:'🎨', description:'Topbar, navbar, footer, colors' },
  ]},
  { id:'sitemap', label:'Sitemap', items:[
    { id:'sitemap', label:'Page Sitemap', singular:null, icon:'', description:'All public pages, slugs, and URLs' },
  ]},
  { id:'settings', label:'Settings', items:[
    { id:'settings', label:'Site Settings', singular:null, icon:'⚙️', description:'Identity, SEO, SMTP, modules' },
  ]},
]

const creatableSections = ['posts','blogs','events','case-studies','pages','careers','slides','team','services','tickets','knowledgebase','fun-facts','faqs','resources','facilities','associations','invest']

const activeSection = ref('posts')
const collapsed = ref(new Set())
const toast = ref('')

const currentNavItem = computed(() => {
  for (const g of navGroups) {
    const item = g.items.find(i => i.id === activeSection.value)
    if (item) return item
  }
  return null
})

function toggleGroup(id) {
  const s = new Set(collapsed.value)
  s.has(id) ? s.delete(id) : s.add(id)
  collapsed.value = s
}

function navigate(id) {
  activeSection.value = id
  if (id==='posts')      { postCategory.value=''; postStatus.value=''; loadPosts() }
  else if (id==='blogs') { postCategory.value='blog'; postStatus.value=''; loadPosts() }
  else if (id==='case-studies') { postCategory.value='case_study'; postStatus.value=''; loadPosts() }
  else if (id==='pages') { postCategory.value='page'; postStatus.value=''; loadPosts() }
  else if (id==='events')   loadEvents()
  else if (id==='careers')  loadCareers()
  else if (id==='slides')   loadSlides()
  else if (id==='menus')    loadMenus()
  else if (id==='team')     loadTeam()
  else if (id==='services') loadServices()
  else if (id==='tickets')  loadTickets()
  else if (id==='knowledgebase') loadKB()
  else if (id==='fun-facts')    loadFunFacts()
  else if (id==='faqs')         loadFAQs()
  else if (id==='resources')    loadResources()
  else if (id==='facilities')   loadFacilities()
  else if (id==='associations') loadAssociations()
  else if (id==='invest')       loadInvest()
  else if (id==='sitemap')      loadSitemap()
}

function showToast(msg) {
  toast.value = msg
  setTimeout(() => { toast.value = '' }, 3000)
}

function fmt(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-UG', { day:'numeric', month:'short', year:'numeric' })
}

// ── Posts / Blogs / Case Studies / Pages ──────────────────────────
const posts = ref([])
const postsLoading = ref(false)
const postCategory = ref('')
const postStatus = ref('')

async function loadPosts() {
  postsLoading.value = true
  try {
    const r = await cmsApi.adminListPosts({ category: postCategory.value, status: postStatus.value, per_page: 100 })
    posts.value = r.data.data?.items || []
  } finally { postsLoading.value = false }
}

function editPost(p) { editingId.value=p.id; form.value={...p, visibility:p.visibility||'public'}; formError.value=''; showModal.value=true }

// ── Events ────────────────────────────────────────────────────────
const events = ref([])
const eventsLoading = ref(false)
async function loadEvents() {
  eventsLoading.value = true
  try { const r = await cmsApi.adminListEvents({ per_page:100 }); events.value = r.data.data?.items||[] }
  finally { eventsLoading.value = false }
}
function editEvent(e) { editingId.value=e.id; form.value={...e, event_date_str:e.event_date?.slice(0,16)||'', end_date_str:e.end_date?.slice(0,16)||''}; formError.value=''; showModal.value=true }

// ── Careers ───────────────────────────────────────────────────────
const careers = ref([])
const careersLoading = ref(false)
async function loadCareers() {
  careersLoading.value = true
  try { const r = await cmsApi.adminListCareers({ per_page:100 }); careers.value = r.data.data?.items||[] }
  finally { careersLoading.value = false }
}
function editCareer(c) { editingId.value=c.id; form.value={...c, deadline_str:c.deadline_at?.slice(0,10)||''}; formError.value=''; showModal.value=true }

// ── Slides ────────────────────────────────────────────────────────
const slides = ref([])
const slidesLoading = ref(false)
let dragSlide = null
async function loadSlides() {
  slidesLoading.value = true
  try { const r = await cmsApi.adminListSlides(); slides.value = r.data.data||[] }
  finally { slidesLoading.value = false }
}
function editSlide(s) { editingId.value=s.id; form.value={...s}; formError.value=''; showModal.value=true }
async function toggleSlideActive(slide) {
  await cmsApi.adminUpdateSlide(slide.id, {...slide, is_active:!slide.is_active})
  await loadSlides()
}
function dragStart(e,item) { dragSlide=item }
function dropSlide(e,target) {
  if (!dragSlide||dragSlide.id===target.id) return
  const arr=[...slides.value]
  const from=arr.findIndex(s=>s.id===dragSlide.id), to=arr.findIndex(s=>s.id===target.id)
  arr.splice(to,0,arr.splice(from,1)[0])
  arr.forEach((s,i)=>s.sort_order=i)
  slides.value=arr
  Promise.all(arr.map(s=>cmsApi.adminUpdateSlide(s.id,s)))
  dragSlide=null
}

// ── Menus ─────────────────────────────────────────────────────────
const menus = reactive({ main:[], footer:[] })
const savingMenu = ref('')
let menuDragSrc = null
async function loadMenus() {
  try {
    const [main,footer] = await Promise.all([cmsApi.getMenu('main'),cmsApi.getMenu('footer')])
    menus.main   = main.data.data?.items   || main.data.data?.Items   || []
    menus.footer = footer.data.data?.items || footer.data.data?.Items || []
  } catch { menus.main=[]; menus.footer=[] }
}
function makeId() { return Math.random().toString(36).slice(2,10) }
function addMenuItem(n) { menus[n].push({id:makeId(),label:'',url:'',children:[],mega:false,megaItems:[]}) }
function removeMenuItem(n,id) { menus[n]=menus[n].filter(i=>i.id!==id) }
function addSubItem(n,pid) { const item=menus[n].find(i=>i.id===pid); if(item){ item.children=[...(item.children||[]),{id:makeId(),label:'',url:''}] } }
function removeSubItem(n,pid,cid) { const item=menus[n].find(i=>i.id===pid); if(item){ item.children=item.children.filter(c=>c.id!==cid) } }
function menuDragStart(e,n,item) { menuDragSrc={n,item} }
function menuDrop(e,n,target) {
  if(!menuDragSrc||menuDragSrc.item.id===target.id||menuDragSrc.n!==n) return
  const arr=menus[n], from=arr.findIndex(i=>i.id===menuDragSrc.item.id), to=arr.findIndex(i=>i.id===target.id)
  arr.splice(to,0,arr.splice(from,1)[0])
  menus[n]=[...arr]; menuDragSrc=null
}
async function saveMenu(n) {
  savingMenu.value=n
  try { await cmsApi.adminUpdateMenu(n,menus[n]); showToast('Menu saved!') }
  catch {} finally { savingMenu.value='' }
}

// ── Widget Builder ────────────────────────────────────────────────
const availableWidgets = [
  {type:'html',label:'HTML Block',icon:'📝'},{type:'image',label:'Image',icon:'🖼'},
  {type:'newsletter',label:'Newsletter',icon:'📧'},{type:'recent-posts',label:'Recent Posts',icon:'📰'},
  {type:'social',label:'Social Links',icon:'🔗'},{type:'contact',label:'Contact Info',icon:'📞'},
  {type:'hours',label:'Opening Hours',icon:'🕐'},{type:'about',label:'About Text',icon:'📖'},
  {type:'tags',label:'Tag Cloud',icon:'🏷'},
]
const widgetColumns = reactive({1:[],2:[],3:[],4:[]})
const widgetDropTarget = ref(null)
let dragWidget = null
function widgetDragStart(e,w) { dragWidget={...w,content:'',imageUrl:'',text:''} }
function dropWidget(e,col) {
  widgetDropTarget.value=null
  if(!dragWidget) return
  if(!widgetColumns[col]) widgetColumns[col]=[]
  widgetColumns[col]=[...widgetColumns[col],{...dragWidget}]
  dragWidget=null
}
function removeWidget(col,idx) { widgetColumns[col]=widgetColumns[col].filter((_,i)=>i!==idx) }
function saveWidgets() {
  localStorage.setItem('cms_widgets',JSON.stringify(widgetColumns))
  showToast('Widget layout saved!')
}

// ── Homepage Builder ──────────────────────────────────────────────
const homeSections = reactive([
  {id:'hero',label:'Hero Section',icon:'🎯',description:'Slideshow hero banner',visible:true},
  {id:'services',label:'Services',icon:'🔧',description:'Services offered grid',visible:true},
  {id:'counters',label:'Counter Up',icon:'📊',description:'Statistics and numbers',visible:true},
  {id:'projects',label:'Projects',icon:'📂',description:'Featured projects',visible:true},
  {id:'blog',label:'Blog Area',icon:'📰',description:'Latest blog posts',visible:false},
  {id:'testimonials',label:'Testimonials',icon:'💬',description:'Client testimonials',visible:true},
  {id:'team',label:'Team',icon:'👥',description:'Team members section',visible:false},
  {id:'contact',label:'Contact',icon:'📞',description:'Contact form and info',visible:true},
  {id:'newsletter',label:'Newsletter',icon:'📧',description:'Email subscription area',visible:false},
])
const homeSettings = reactive({ heroStyle:'gradient', overlayOpacity:50, counterSpeed:2000 })
let homeDragSrc = null
function homeDragStart(e,sec) { homeDragSrc=sec }
function homeDrop(e,target) {
  if(!homeDragSrc||homeDragSrc.id===target.id) return
  const from=homeSections.findIndex(s=>s.id===homeDragSrc.id), to=homeSections.findIndex(s=>s.id===target.id)
  homeSections.splice(to,0,homeSections.splice(from,1)[0])
  homeDragSrc=null
}
function saveHomeSettings() {
  localStorage.setItem('cms_home',JSON.stringify({sections:homeSections,settings:homeSettings}))
  showToast('Home page settings saved!')
}

// ── Team ──────────────────────────────────────────────────────────
const team = ref([])
const teamLoading = ref(false)
async function loadTeam() {
  teamLoading.value=true
  try { const r=await cmsApi.adminListTeam?.({per_page:50}); team.value=r.data.data?.items||[] }
  catch { team.value=[] } finally { teamLoading.value=false }
}
function editTeam(m) { editingId.value=m.id; form.value={...m}; formError.value=''; showModal.value=true }

// ── Services ──────────────────────────────────────────────────────
const services = ref([])
const servicesLoading = ref(false)
async function loadServices() {
  servicesLoading.value=true
  try { const r=await cmsApi.adminListServices?.({per_page:50}); services.value=r.data.data?.items||[] }
  catch { services.value=[] } finally { servicesLoading.value=false }
}
function editService(s) { editingId.value=s.id; form.value={...s}; formError.value=''; showModal.value=true }

// ── Tickets ───────────────────────────────────────────────────────
const tickets = ref([])
const ticketsLoading = ref(false)
async function loadTickets() {
  ticketsLoading.value=true
  try { const r=await cmsApi.adminListTickets?.({per_page:50}); tickets.value=r.data.data?.items||[] }
  catch { tickets.value=[] } finally { ticketsLoading.value=false }
}
function editTicket(t) { editingId.value=t.id; form.value={...t}; formError.value=''; showModal.value=true }

// ── Knowledgebase ─────────────────────────────────────────────────
const kbArticles = ref([])
const kbLoading = ref(false)
async function loadKB() {
  kbLoading.value=true
  try { const r=await cmsApi.adminListArticles?.({per_page:50}); kbArticles.value=r.data.data?.items||[] }
  catch { kbArticles.value=[] } finally { kbLoading.value=false }
}
function editKB(a) { editingId.value=a.id; form.value={...a}; formError.value=''; showModal.value=true }

// ── Fun Facts ─────────────────────────────────────────────────────
const funFacts = ref([])
const funFactsLoading = ref(false)
async function loadFunFacts() {
  funFactsLoading.value = true
  try { const r = await cmsApi.adminListFunFacts(); funFacts.value = r.data.data || [] }
  catch { funFacts.value = [] } finally { funFactsLoading.value = false }
}
function editFunFact(f) { editingId.value=f.id; form.value={...f}; formError.value=''; showModal.value=true }

// ── FAQs ──────────────────────────────────────────────────────────
const faqs = ref([])
const faqsLoading = ref(false)
const faqCategory = ref('')
async function loadFAQs() {
  faqsLoading.value = true
  try { const r = await cmsApi.adminListFAQs({ category: faqCategory.value }); faqs.value = r.data.data || [] }
  catch { faqs.value = [] } finally { faqsLoading.value = false }
}
function editFAQ(f) { editingId.value=f.id; form.value={...f}; formError.value=''; showModal.value=true }

// ── Resources ─────────────────────────────────────────────────────
const resources = ref([])
const resourcesLoading = ref(false)
const resourceCategory = ref('')
async function loadResources() {
  resourcesLoading.value = true
  try { const r = await cmsApi.adminListResources({ category: resourceCategory.value }); resources.value = r.data.data || [] }
  catch { resources.value = [] } finally { resourcesLoading.value = false }
}
function editResource(res) { editingId.value=res.id; form.value={...res}; formError.value=''; showModal.value=true }

// ── Facilities ────────────────────────────────────────────────────
const facilities = ref([])
const facilitiesLoading = ref(false)
async function loadFacilities() {
  facilitiesLoading.value = true
  try { const r = await cmsApi.adminListFacilities(); facilities.value = r.data.data || [] }
  catch { facilities.value = [] } finally { facilitiesLoading.value = false }
}
function editFacility(f) { editingId.value=f.id; form.value={...f}; formError.value=''; showModal.value=true }

// ── Associations ──────────────────────────────────────────────────
const associations = ref([])
const associationsLoading = ref(false)
async function loadAssociations() {
  associationsLoading.value = true
  try { const r = await cmsApi.adminListAssociations(); associations.value = r.data.data || [] }
  catch { associations.value = [] } finally { associationsLoading.value = false }
}
function editAssociation(a) { editingId.value=a.id; form.value={...a}; formError.value=''; showModal.value=true }

// ── Invest with Us ────────────────────────────────────────────────
const investItems = ref([])
const investLoading = ref(false)
async function loadInvest() {
  investLoading.value = true
  try { const r = await cmsApi.adminListInvest(); investItems.value = r.data.data || [] }
  catch { investItems.value = [] } finally { investLoading.value = false }
}
function editInvest(i) { editingId.value=i.id; form.value={...i}; formError.value=''; showModal.value=true }

// ── Sitemap ───────────────────────────────────────────────────────
const sitemap = ref([])
const sitemapLoading = ref(false)
async function loadSitemap() {
  sitemapLoading.value = true
  const origin = window.location.origin
  try {
    const [postsR, eventsR, careersR, slidesR, faqsR, resourcesR, facilitiesR, associationsR, investR] = await Promise.allSettled([
      cmsApi.adminListPosts({ per_page: 200 }),
      cmsApi.adminListEvents({ per_page: 200 }),
      cmsApi.adminListCareers({ per_page: 200 }),
      cmsApi.adminListSlides(),
      cmsApi.adminListFAQs(),
      cmsApi.adminListResources(),
      cmsApi.adminListFacilities(),
      cmsApi.adminListAssociations(),
      cmsApi.adminListInvest(),
    ])
    const safe = (res, key) => {
      if (res.status !== 'fulfilled') return []
      const d = res.value.data
      if (key) return d.data?.[key] || d.data || []
      return d.data || []
    }
    const posts = safe(postsR, 'items')
    const events = safe(eventsR, 'items')
    const careers = safe(careersR, 'items')
    const faqs = safe(faqsR)
    const resources = safe(resourcesR)
    const facilities = safe(facilitiesR)
    const associations = safe(associationsR)
    const investRows = safe(investR)

    sitemap.value = [
      {
        id: 'static',
        label: 'Static Pages',
        description: 'Always available — built-in routes',
        entries: [
          { title: 'Home',          slug: '/',                url: origin + '/' },
          { title: 'Apply (License Portal)', slug: '/apply',  url: origin + '/apply' },
          { title: 'News',          slug: '/news',            url: origin + '/news' },
          { title: 'Events',        slug: '/events',          url: origin + '/events' },
          { title: 'Careers',       slug: '/careers',         url: origin + '/careers' },
          { title: 'Projects',      slug: '/projects',        url: origin + '/projects' },
          { title: 'Case Studies',  slug: '/case-studies',    url: origin + '/case-studies' },
          { title: 'Facilities',    slug: '/facilities',      url: origin + '/facilities' },
          { title: 'Associations',  slug: '/associations',    url: origin + '/associations' },
          { title: 'Resource Centre',slug:'/resource-centre', url: origin + '/resource-centre' },
          { title: 'FAQs',          slug: '/faqs',            url: origin + '/faqs' },
          { title: 'Invest with Us',slug: '/invest',          url: origin + '/invest' },
          { title: 'Sign In',       slug: '/login',           url: origin + '/login' },
        ]
      },
      {
        id: 'posts',
        label: `Posts & Articles (${posts.length})`,
        description: 'Dynamic content from the Posts module',
        entries: posts.map(p => ({ title: p.title, slug: p.slug, url: `${origin}/news/${p.slug}` }))
      },
      {
        id: 'events',
        label: `Events (${events.length})`,
        entries: events.map(e => ({ title: e.title, slug: e.slug, url: `${origin}/events/${e.slug}` }))
      },
      {
        id: 'careers',
        label: `Careers (${careers.length})`,
        entries: careers.map(c => ({ title: c.title, slug: c.id, url: `${origin}/careers/${c.id}` }))
      },
      {
        id: 'facilities',
        label: `Facilities (${facilities.length})`,
        entries: facilities.map(f => ({ title: f.name, slug: f.slug, url: `${origin}/facilities/${f.slug}` }))
      },
      {
        id: 'associations',
        label: `Associations (${associations.length})`,
        entries: associations.map(a => ({ title: a.name, slug: a.slug, url: `${origin}/associations/${a.slug}` }))
      },
      {
        id: 'resources',
        label: `Resource Centre (${resources.length})`,
        entries: resources.map(r => ({ title: r.title, slug: r.category, url: `${origin}/resource-centre#${r.category}` }))
      },
      {
        id: 'faqs',
        label: `FAQs (${faqs.length})`,
        entries: faqs.map(f => ({ title: f.question, slug: f.category, url: `${origin}/faqs#${(f.category||'general')}` }))
      },
      {
        id: 'invest',
        label: `Invest with Us (${investRows.length})`,
        entries: investRows.map(i => ({ title: i.title, slug: '/invest', url: `${origin}/invest#${i.id}` }))
      },
    ]
  } finally {
    sitemapLoading.value = false
  }
}

// ── Appearance settings ───────────────────────────────────────────
const appearanceTabs = [
  {id:'topbar',label:'Topbar'},{id:'navbar',label:'Navbar'},{id:'footer',label:'Footer'},
  {id:'colors',label:'Colors'},{id:'typography',label:'Typography'},
]
const appearanceTab = ref('topbar')
const googleFonts = ['Inter','Roboto','Open Sans','Lato','Poppins','Montserrat','Raleway','Nunito','Source Sans Pro','Playfair Display','Merriweather']
const appearance = reactive({
  topbar:  { bg:'#112b4e', text:'#ffffff', phone:'', email:'', show:true },
  navbar:  { bg:'#ffffff', text:'#112b4e', logoAlign:'left', sticky:true, transparentOnHero:false },
  footer:  { bg:'#0f1c2e', text:'#9ca3af', copyright:'© 2025 NCS Uganda. All rights reserved.', columns:4 },
  colors:  { primary:'#112b4e', secondary:'#1e5799', accent:'#e8a020', dark:'#0a1628' },
  typography: { heading:'Poppins', body:'Inter', size:'16px' },
})
function saveAppearance() {
  localStorage.setItem('cms_appearance',JSON.stringify(appearance))
  showToast('Appearance saved!')
}

// ── Contact settings ──────────────────────────────────────────────
const contactSettings = reactive({ phone:'', email:'', address:'', hours:'', mapUrl:'', social:{facebook:'',twitter:'',linkedin:'',instagram:'',youtube:''} })
function saveContactSettings() {
  localStorage.setItem('cms_contact',JSON.stringify(contactSettings))
  showToast('Contact info saved!')
}

// ── Site settings ─────────────────────────────────────────────────
const settingsTabs = [
  {id:'identity',label:'Site Identity'},{id:'seo',label:'SEO'},{id:'smtp',label:'SMTP'},
  {id:'css',label:'Custom CSS/JS'},{id:'modules',label:'Modules'},
]
const settingsTab = ref('identity')
const siteModules = [
  {key:'blog',label:'Blog',description:'Enable blog section on public site'},
  {key:'events',label:'Events',description:'Show events section'},
  {key:'careers',label:'Careers',description:'Display job listings'},
  {key:'testimonials',label:'Testimonials',description:'Show testimonials section'},
  {key:'newsletter',label:'Newsletter',description:'Newsletter subscription form'},
  {key:'team',label:'Team Section',description:'Display team members'},
  {key:'services',label:'Services',description:'Show services section'},
  {key:'caseStudies',label:'Case Studies',description:'Show case studies section'},
  {key:'knowledgebase',label:'Knowledgebase',description:'Public help articles'},
  {key:'support',label:'Support Tickets',description:'Customer support portal'},
]
const siteSettings = reactive({
  name:'NCS Uganda', tagline:'National Council of Sports', logoUrl:'/main-logo.png', faviconUrl:'/favicon.png',
  titleTemplate:'%s | NCS Uganda', metaDescription:'', robots:'index,follow',
  smtpHost:'', smtpPort:587, smtpSSL:false, smtpUser:'', smtpPass:'', fromEmail:'', fromName:'',
  customCss:'', customJs:'',
  modules:{ blog:true, events:true, careers:true, testimonials:true, newsletter:false, team:false, services:false, caseStudies:true, knowledgebase:false, support:false },
})
function saveSettings() {
  const s={...siteSettings}; delete s.smtpPass
  localStorage.setItem('cms_settings',JSON.stringify(s))
  showToast('Settings saved!')
}

// ── Modal ─────────────────────────────────────────────────────────
const showModal  = ref(false)
const editingId  = ref(null)
const form       = ref({})
const formError  = ref('')
const saving     = ref(false)
const deleteTarget = ref(null)
const deleting   = ref(false)

function openCreate() {
  editingId.value=null; formError.value=''
  const s=activeSection.value
  if (s==='posts')        form.value={category:'news',status:'draft',visibility:'public',title:'',content:'',excerpt:'',cover_image_url:'',tags:'',meta_title:'',meta_description:''}
  else if (s==='blogs')   form.value={category:'blog',status:'draft',visibility:'public',title:'',content:'',excerpt:'',cover_image_url:'',tags:'',video_url:'',meta_title:'',meta_description:''}
  else if (s==='case-studies') form.value={category:'case_study',status:'draft',visibility:'public',title:'',content:'',excerpt:'',cover_image_url:'',tags:'',meta_title:'',meta_description:''}
  else if (s==='pages')   form.value={category:'page',status:'draft',visibility:'public',title:'',content:'',meta_title:'',meta_description:''}
  else if (s==='events')  form.value={title:'',description:'',location:'',status:'draft',cover_image_url:'',event_date_str:'',end_date_str:''}
  else if (s==='careers') form.value={title:'',description:'',department:'',location:'',category:'jobs',job_type:'full_time',salary_range:'',requirements:'',status:'draft',deadline_str:''}
  else if (s==='slides')  form.value={title:'',subtitle:'',description:'',image_url:'',button_text:'',button_url:'',sort_order:0,is_active:true}
  else if (s==='team')    form.value={full_name:'',designation:'',image_url:'',bio:'',sort_order:0}
  else if (s==='services') form.value={title:'',description:'',category:'',icon:'',image_url:'',status:'published'}
  else if (s==='tickets') form.value={subject:'',department:'general',priority:'normal',message:''}
  else if (s==='knowledgebase') form.value={title:'',topic:'',content:'',status:'draft'}
  else if (s==='fun-facts')    form.value={label:'',value:'',icon:'',sort_order:0,is_active:true}
  else if (s==='faqs')         form.value={question:'',answer:'',category:'general',sort_order:0,is_active:true}
  else if (s==='resources')    form.value={title:'',category:'downloads',file_url:'',description:'',sort_order:0,is_active:true}
  else if (s==='facilities')   form.value={name:'',description:'',image_url:'',sort_order:0,is_active:true}
  else if (s==='associations') form.value={name:'',description:'',logo_url:'',website_url:'',sort_order:0,is_active:true}
  else if (s==='invest')       form.value={title:'',subtitle:'',content:'',image_url:'',sort_order:0,is_active:true}
  showModal.value=true
}

async function saveItem() {
  formError.value=''; saving.value=true
  try {
    const s=activeSection.value
    if (['posts','blogs','case-studies','pages'].includes(s)) {
      if (!form.value.title) { formError.value='Title is required'; return }
      editingId.value ? await cmsApi.adminUpdatePost(editingId.value,form.value) : await cmsApi.adminCreatePost(form.value)
      await loadPosts()
    } else if (s==='events') {
      if (!form.value.title) { formError.value='Title is required'; return }
      const d={...form.value}
      if (d.event_date_str) { d.event_date=new Date(d.event_date_str).toISOString(); delete d.event_date_str }
      if (d.end_date_str)   { d.end_date=new Date(d.end_date_str).toISOString(); delete d.end_date_str }
      editingId.value ? await cmsApi.adminUpdateEvent(editingId.value,d) : await cmsApi.adminCreateEvent(d)
      await loadEvents()
    } else if (s==='careers') {
      if (!form.value.title||!form.value.description) { formError.value='Title and description are required'; return }
      const d={...form.value}
      if (d.deadline_str) { d.deadline_at=new Date(d.deadline_str).toISOString(); delete d.deadline_str }
      editingId.value ? await cmsApi.adminUpdateCareer(editingId.value,d) : await cmsApi.adminCreateCareer(d)
      await loadCareers()
    } else if (s==='slides') {
      if (!form.value.title) { formError.value='Title is required'; return }
      editingId.value ? await cmsApi.adminUpdateSlide(editingId.value,form.value) : await cmsApi.adminCreateSlide(form.value)
      await loadSlides()
    } else if (s==='team') {
      if (!form.value.full_name) { formError.value='Name is required'; return }
      try { editingId.value ? await cmsApi.adminUpdateTeam?.(editingId.value,form.value) : await cmsApi.adminCreateTeam?.(form.value) } catch {}
      await loadTeam()
    } else if (s==='services') {
      if (!form.value.title) { formError.value='Title is required'; return }
      try { editingId.value ? await cmsApi.adminUpdateService?.(editingId.value,form.value) : await cmsApi.adminCreateService?.(form.value) } catch {}
      await loadServices()
    } else if (s==='tickets') {
      if (!form.value.subject) { formError.value='Subject is required'; return }
      try { editingId.value ? await cmsApi.adminUpdateTicket?.(editingId.value,form.value) : await cmsApi.adminCreateTicket?.(form.value) } catch {}
      await loadTickets()
    } else if (s==='knowledgebase') {
      if (!form.value.title) { formError.value='Title is required'; return }
      try { editingId.value ? await cmsApi.adminUpdateArticle?.(editingId.value,form.value) : await cmsApi.adminCreateArticle?.(form.value) } catch {}
      await loadKB()
    } else if (s==='fun-facts') {
      if (!form.value.label||!form.value.value) { formError.value='Label and value are required'; return }
      editingId.value ? await cmsApi.adminUpdateFunFact(editingId.value,form.value) : await cmsApi.adminCreateFunFact(form.value)
      await loadFunFacts()
    } else if (s==='faqs') {
      if (!form.value.question||!form.value.answer) { formError.value='Question and answer are required'; return }
      editingId.value ? await cmsApi.adminUpdateFAQ(editingId.value,form.value) : await cmsApi.adminCreateFAQ(form.value)
      await loadFAQs()
    } else if (s==='resources') {
      if (!form.value.title) { formError.value='Title is required'; return }
      editingId.value ? await cmsApi.adminUpdateResource(editingId.value,form.value) : await cmsApi.adminCreateResource(form.value)
      await loadResources()
    } else if (s==='facilities') {
      if (!form.value.name) { formError.value='Name is required'; return }
      editingId.value ? await cmsApi.adminUpdateFacility(editingId.value,form.value) : await cmsApi.adminCreateFacility(form.value)
      await loadFacilities()
    } else if (s==='associations') {
      if (!form.value.name) { formError.value='Name is required'; return }
      editingId.value ? await cmsApi.adminUpdateAssociation(editingId.value,form.value) : await cmsApi.adminCreateAssociation(form.value)
      await loadAssociations()
    } else if (s==='invest') {
      if (!form.value.title) { formError.value='Title is required'; return }
      editingId.value ? await cmsApi.adminUpdateInvest(editingId.value,form.value) : await cmsApi.adminCreateInvest(form.value)
      await loadInvest()
    }
    showModal.value=false
    showToast('Saved successfully!')
  } catch (err) {
    formError.value=err.response?.data?.error?.message||'Failed to save.'
  } finally { saving.value=false }
}

function confirmDelete(type,id) { deleteTarget.value={type,id} }
async function doDelete() {
  if (!deleteTarget.value) return
  deleting.value=true
  try {
    const {type,id}=deleteTarget.value
    if (type==='post')   { await cmsApi.adminDeletePost(id);   await loadPosts() }
    if (type==='event')  { await cmsApi.adminDeleteEvent(id);  await loadEvents() }
    if (type==='career') { await cmsApi.adminDeleteCareer(id); await loadCareers() }
    if (type==='slide')  { await cmsApi.adminDeleteSlide(id);  await loadSlides() }
    try {
      if (type==='team')    { await cmsApi.adminDeleteTeam?.(id);    await loadTeam() }
      if (type==='service') { await cmsApi.adminDeleteService?.(id); await loadServices() }
      if (type==='ticket')  { await cmsApi.adminDeleteTicket?.(id);  await loadTickets() }
      if (type==='kb')      { await cmsApi.adminDeleteArticle?.(id); await loadKB() }
    } catch {}
    if (type==='fun-fact')   { await cmsApi.adminDeleteFunFact(id);   await loadFunFacts() }
    if (type==='faq')        { await cmsApi.adminDeleteFAQ(id);        await loadFAQs() }
    if (type==='resource')   { await cmsApi.adminDeleteResource(id);   await loadResources() }
    if (type==='facility')   { await cmsApi.adminDeleteFacility(id);   await loadFacilities() }
    if (type==='association'){ await cmsApi.adminDeleteAssociation(id);await loadAssociations() }
    if (type==='invest')     { await cmsApi.adminDeleteInvest(id);     await loadInvest() }
    deleteTarget.value=null
    showToast('Deleted.')
  } catch {} finally { deleting.value=false }
}

// ── Load saved settings from localStorage ────────────────────────
function loadLocalSettings() {
  try {
    const a=localStorage.getItem('cms_appearance'); if(a) Object.assign(appearance,JSON.parse(a))
    const c=localStorage.getItem('cms_contact');    if(c) Object.assign(contactSettings,JSON.parse(c))
    const s=localStorage.getItem('cms_settings');   if(s) Object.assign(siteSettings,JSON.parse(s))
    const w=localStorage.getItem('cms_widgets');    if(w) { const d=JSON.parse(w); Object.keys(d).forEach(k=>widgetColumns[k]=d[k]) }
    const h=localStorage.getItem('cms_home');
    if (h) {
      const d=JSON.parse(h)
      if (d.sections) { homeSections.splice(0,homeSections.length,...d.sections) }
      if (d.settings) Object.assign(homeSettings,d.settings)
    }
  } catch {}
}

onMounted(() => {
  loadLocalSettings()
  loadPosts()
})
</script>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.3s ease; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(-12px); }
</style>
