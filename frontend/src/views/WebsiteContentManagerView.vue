<template>
  <main class="min-h-screen bg-[#f6f8fb] text-[#1f2937]">
    <div class="cms-shell" :class="{ collapsed: sidebarCollapsed }">
      <aside class="cms-sidebar">
        <div class="cms-sidebar-head">
          <router-link to="/" class="cms-brand">
            <img src="/main-logo.png" alt="NCS" />
            <span>Website CMS</span>
          </router-link>
          <button
            type="button"
            class="cms-collapse-toggle"
            @click="toggleSidebar"
            :aria-expanded="!sidebarCollapsed ? 'true' : 'false'"
            :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'"
          >
            <i :class="sidebarCollapsed ? 'icofont-rounded-right' : 'icofont-rounded-left'" aria-hidden="true"></i>
          </button>
        </div>
        <nav class="cms-nav" aria-label="Content manager sections" @click.capture="expandSidebarOnNav">
          <button v-for="item in topSections" :key="item.id" type="button" :class="{ active: active === item.id }" :title="item.label" @click="active = item.id">
            <i :class="item.icon" aria-hidden="true"></i>
            <span>{{ item.label }}</span>
          </button>
          <div
            class="cms-nav-group"
            :class="{ expanded: homepageGroupOpen }"
            @mouseenter="homepageGroupOpen = true"
            @mouseleave="homepageGroupOpen = false"
            @focusin="homepageGroupOpen = true"
            @focusout="handleHomepageGroupFocusOut"
          >
            <button
              type="button"
              class="cms-nav-parent"
              :aria-expanded="homepageGroupOpen ? 'true' : 'false'"
              aria-controls="homepage-management-submenu"
              title="Homepage Management"
              @click="homepageGroupOpen = !homepageGroupOpen"
            >
              <i class="icofont-home" aria-hidden="true"></i>
              <span>Homepage Management</span>
              <i class="icofont-rounded-down" aria-hidden="true"></i>
            </button>
            <div id="homepage-management-submenu" class="cms-nav-submenu">
              <button v-for="item in homepageSections" :key="item.id" type="button" :class="{ active: active === item.id }" @click="active = item.id">
                <i :class="item.icon" aria-hidden="true"></i>
                <span>{{ item.label }}</span>
              </button>
            </div>
          </div>
          <div
            class="cms-nav-group"
            :class="{ expanded: slideshowGroupOpen }"
            @mouseenter="slideshowGroupOpen = true"
            @mouseleave="slideshowGroupOpen = false"
            @focusin="slideshowGroupOpen = true"
            @focusout="handleSlideshowGroupFocusOut"
          >
            <button
              type="button"
              class="cms-nav-parent"
              :aria-expanded="slideshowGroupOpen ? 'true' : 'false'"
              aria-controls="slideshow-management-submenu"
              title="Slideshow Manager"
              @click="slideshowGroupOpen = !slideshowGroupOpen"
            >
              <i class="icofont-image" aria-hidden="true"></i>
              <span>Slideshow Manager</span>
              <i class="icofont-rounded-down" aria-hidden="true"></i>
            </button>
            <div id="slideshow-management-submenu" class="cms-nav-submenu">
              <button v-for="item in slideshowSections" :key="item.id" type="button" :class="{ active: active === item.id }" @click="active = item.id">
                <i :class="item.icon" aria-hidden="true"></i>
                <span>{{ item.label }}</span>
              </button>
            </div>
          </div>
          <div
            class="cms-nav-group"
            :class="{ expanded: blogsGroupOpen }"
            @mouseenter="blogsGroupOpen = true"
            @mouseleave="blogsGroupOpen = false"
            @focusin="blogsGroupOpen = true"
            @focusout="handleBlogsGroupFocusOut"
          >
            <button
              type="button"
              class="cms-nav-parent"
              :aria-expanded="blogsGroupOpen ? 'true' : 'false'"
              aria-controls="blogs-management-submenu"
              title="Blogs Management"
              @click="blogsGroupOpen = !blogsGroupOpen"
            >
              <i class="icofont-newspaper" aria-hidden="true"></i>
              <span>Blogs Management</span>
              <i class="icofont-rounded-down" aria-hidden="true"></i>
            </button>
            <div id="blogs-management-submenu" class="cms-nav-submenu">
              <button v-for="item in blogSections" :key="item.id" type="button" :class="{ active: active === item.id }" @click="active = item.id">
                <i :class="item.icon" aria-hidden="true"></i>
                <span>{{ item.label }}</span>
              </button>
            </div>
          </div>
          <button v-for="item in contentSections" :key="item.id" type="button" :class="{ active: active === item.id }" :title="item.label" @click="active = item.id">
            <i :class="item.icon" aria-hidden="true"></i>
            <span>{{ item.label }}</span>
          </button>
        </nav>
        <button type="button" class="cms-logout" @click="logout" title="Sign out">
          <i class="icofont-logout" aria-hidden="true"></i>
          <span>Sign out</span>
        </button>
      </aside>

      <section class="cms-main">
        <header class="cms-header">
          <div>
            <p class="cms-kicker">Public website content</p>
            <h1>{{ currentSection.label }}</h1>
          </div>
          <div class="cms-actions">
            <ThemeToggle />
            <router-link to="/" target="_blank">Preview website</router-link>
            <button type="button" @click="loadAll">Refresh</button>
          </div>
        </header>

        <p v-if="message" class="cms-message">{{ message }}</p>
        <p v-if="error" class="cms-error">{{ error }}</p>

        <section v-if="active === 'overview'" class="cms-grid">
          <article v-for="card in overviewCards" :key="card.label" class="cms-card metric">
            <span>{{ card.label }}</span>
            <strong>{{ card.value }}</strong>
          </article>
        </section>

        <section v-else-if="active === 'homepage'" class="cms-panel">
          <div class="cms-panel-head"><h2>Homepage Sections</h2><button @click="saveHomepage">Save homepage</button></div>
          <div class="cms-two">
            <label>Stats title<input v-model="homepage.stats_title" /></label>
            <label>Stats intro<input v-model="homepage.stats_intro" /></label>
            <label>Facilities eyebrow<input v-model="homepage.facilities.eyebrow" /></label>
            <label>Facilities title<input v-model="homepage.facilities.title" /></label>
            <label class="wide">Facilities intro<textarea v-model="homepage.facilities.intro"></textarea></label>
            <label>Events eyebrow<input v-model="homepage.events.eyebrow" /></label>
            <label>Events title<input v-model="homepage.events.title" /></label>
            <label class="wide">Events intro<textarea v-model="homepage.events.intro"></textarea></label>
          </div>
        </section>

        <section v-else-if="active === 'core'" class="cms-panel">
          <div class="cms-panel-head"><h2>Core Functions</h2><button @click="saveHomepage">Save functions</button></div>
          <div class="cms-list-editor">
            <div v-for="(_, index) in homepage.core_functions" :key="index" class="cms-row">
              <textarea v-model="homepage.core_functions[index]"></textarea>
              <button type="button" @click="homepage.core_functions.splice(index, 1)">Remove</button>
            </div>
            <button type="button" @click="homepage.core_functions.push('')">Add function</button>
          </div>
        </section>

        <section v-else-if="active === 'slideshow-manager'" class="cms-panel">
          <SlideshowManager :slideshow="slideshow" :slides="slides" @refresh="loadAll" @message="setMsg" @error="setErr" />
        </section>

        <section v-else-if="active === 'create-post'" class="cms-panel">
          <BlogPostEditor :model="postForm" :categories="blogCategories" @save="savePost" />
        </section>

        <section v-else-if="active === 'manage-posts'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Posts</h2><button type="button" @click="resetPostForm(); active = 'create-post'">New post</button></div>
          <ContentTable :items="posts" title-key="title" subtitle-key="status" @edit="editPost" @delete="removePost" />
        </section>

        <section v-else-if="active === 'create-blog-categories'" class="cms-panel">
          <EditorForm title="Create Blog Categories" :model="blogCategoryForm" :fields="blogCategoryFields" @save="saveBlogCategory" />
        </section>

        <section v-else-if="active === 'manage-blog-categories'" class="cms-panel">
          <div class="cms-panel-head"><h2>Manage Categories</h2><button type="button" @click="resetBlogCategoryForm(); active = 'create-blog-categories'">New category</button></div>
          <ContentTable :items="blogCategories" title-key="name" subtitle-key="slug" @edit="editBlogCategory" @delete="removeBlogCategory" />
        </section>

        <section v-else-if="active === 'comments'" class="cms-panel">
          <div class="cms-panel-head">
            <h2>Comment Moderation Queue</h2>
            <select v-model="commentStatus" @change="loadComments">
              <option value="">All comments</option>
              <option value="pending">Pending</option>
              <option value="approved">Approved</option>
              <option value="flagged">Flagged</option>
            </select>
          </div>
          <div class="cms-table">
            <article v-for="comment in comments" :key="comment.id" class="cms-table-row comment-row">
              <div>
                <strong>{{ comment.user_name }} on {{ comment.post_title }}</strong>
                <span>{{ comment.status }} · {{ comment.body }}</span>
              </div>
              <div>
                <button type="button" @click="moderateComment(comment, 'approve')">Approve</button>
                <button type="button" @click="moderateComment(comment, 'flag')">Flag</button>
                <button type="button" @click="deleteComment(comment)">Delete</button>
              </div>
            </article>
            <p v-if="!comments.length" class="cms-empty">No comments in this queue.</p>
          </div>
        </section>

        <section v-else-if="active === 'events'" class="cms-panel">
          <EditorForm title="Events" :model="eventForm" :fields="eventFields" @save="saveEvent" />
          <ContentTable :items="events" title-key="title" subtitle-key="category" @edit="editEvent" @delete="removeEvent" />
        </section>

        <section v-else-if="active === 'facilities'" class="cms-panel">
          <EditorForm title="Facilities" :model="facilityForm" :fields="facilityFields" @save="saveFacility" />
          <ContentTable :items="facilities" title-key="name" subtitle-key="category" @edit="editFacility" @delete="removeFacility" />
        </section>

        <section v-else-if="active === 'associations'" class="cms-panel">
          <EditorForm title="Sports Associations" :model="associationForm" :fields="associationFields" @save="saveAssociation" />
          <ContentTable :items="associations" title-key="name" subtitle-key="category" @edit="editAssociation" @delete="removeAssociation" />
        </section>

        <section v-else-if="active === 'facts'" class="cms-panel">
          <EditorForm title="Sports Excellence Numbers" :model="factForm" :fields="factFields" @save="saveFact" />
          <ContentTable :items="facts" title-key="label" subtitle-key="value" @edit="editFact" @delete="removeFact" />
        </section>

        <section v-else-if="active === 'faqs'" class="cms-panel">
          <EditorForm title="FAQs" :model="faqForm" :fields="faqFields" @save="saveFAQ" />
          <ContentTable :items="faqs" title-key="question" subtitle-key="category" @edit="editFAQ" @delete="removeFAQ" />
        </section>

        <section v-else-if="active === 'resources'" class="cms-panel">
          <EditorForm title="Resource Centre" :model="resourceForm" :fields="resourceFields" @save="saveResource" />
          <ContentTable :items="resources" title-key="title" subtitle-key="category" @edit="editResource" @delete="removeResource" />
        </section>

        <section v-else-if="active === 'menus'" class="cms-panel">
          <div class="cms-panel-head"><h2>Dynamic Menu Builder</h2><button @click="saveMenus">Save menus</button></div>
          <div class="cms-menu-builder-grid">
            <article>
              <h3>Main menu</h3>
              <MenuBuilder v-model="mainMenuTree" />
            </article>
            <article>
              <h3>Footer menu</h3>
              <MenuBuilder v-model="footerMenuTree" />
            </article>
          </div>
        </section>

        <section v-else-if="active === 'settings'" class="cms-panel">
          <div class="cms-panel-head"><h2>Contact and Footer Settings</h2><button @click="saveSettings">Save settings</button></div>
          <div class="cms-two">
            <label>Phone<input v-model="contact.phone" /></label>
            <label>Email<input v-model="contact.email" /></label>
            <label class="wide">Address<textarea v-model="contact.address"></textarea></label>
            <label class="wide">Footer about<textarea v-model="footer.about"></textarea></label>
          </div>
        </section>
      </section>
    </div>
  </main>
</template>

<script setup>
import { computed, defineComponent, h, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import * as cms from '@/api/cms.js'
import BlogPostEditor from '@/components/cms/BlogPostEditor.vue'
import MenuBuilder from '@/components/cms/MenuBuilder.vue'
import SlideshowManager from '@/components/cms/SlideshowManager.vue'
import ThemeToggle from '@/components/theme/ThemeToggle.vue'
import { normalizeMenuTree, toCmsMenuItems } from '@/utils/menuTree.js'

const router = useRouter()
const active = ref('overview')
const message = ref('')
const error = ref('')
const sidebarCollapsed = ref(localStorage.getItem('ncsms_sidebar_collapsed') === 'true')
const homepageGroupOpen = ref(false)
const blogsGroupOpen = ref(false)
const slideshowGroupOpen = ref(false)

const topSections = [
  { id:'overview', label:'Overview', icon:'icofont-dashboard-web' },
]
const homepageSections = [
  { id:'homepage', label:'Homepage', icon:'icofont-home' },
  { id:'core', label:'Core Functions', icon:'icofont-check-circled' },
  { id:'events', label:'Upcoming Events', icon:'icofont-calendar' },
  { id:'facilities', label:'Sports Facilities', icon:'icofont-stadium' },
  { id:'associations', label:'Sports Associations', icon:'icofont-trophy' },
  { id:'facts', label:'Excellence Numbers', icon:'icofont-chart-bar-graph' },
]
const slideshowSections = [
  { id:'slideshow-manager', label:'Homepage Hero', icon:'icofont-slidshare' },
]
const blogSections = [
  { id:'create-post', label:'Create New Post', icon:'icofont-edit' },
  { id:'manage-posts', label:'Manage Posts', icon:'icofont-list' },
  { id:'create-blog-categories', label:'Create Blog Categories', icon:'icofont-folder-open' },
  { id:'manage-blog-categories', label:'Manage Categories', icon:'icofont-tags' },
]
const contentSections = [
  { id:'comments', label:'Comment Moderation', icon:'icofont-speech-comments' },
  { id:'faqs', label:'FAQs', icon:'icofont-question-circle' },
  { id:'resources', label:'Resources', icon:'icofont-download' },
  { id:'menus', label:'Menus', icon:'icofont-navigation-menu' },
  { id:'settings', label:'Contact/Footer', icon:'icofont-settings' },
]
const sections = [...topSections, ...homepageSections, ...slideshowSections, ...blogSections, ...contentSections]
const currentSection = computed(() => sections.find(s => s.id === active.value) || sections[0])

const homepage = reactive({
  stats_title:'Sports Excellence in Numbers',
  stats_intro:'',
  facilities:{eyebrow:'World-Class Infrastructure', title:'Our Sports Facilities', intro:'', button_label:'Explore All Facilities'},
  events:{eyebrow:'Upcoming Events', title:'NCS Calendar', intro:''},
  core_functions:[],
})
const contact = reactive({ phone:'', email:'', address:'' })
const footer = reactive({ about:'' })

const posts = ref([]), events = ref([]), slides = ref([]), facilities = ref([]), associations = ref([]), facts = ref([]), faqs = ref([]), resources = ref([])
const blogCategories = ref([])
const comments = ref([])
const commentStatus = ref('pending')
const slideshow = reactive({ id:'homepage-hero', name:'Homepage Hero', slug:'homepage-hero', transition_effect:'fade', transition_duration:700, autoplay_speed:6500, pause_on_hover:true, is_active:true })
const mainMenuTree = ref([]), footerMenuTree = ref([])

const postForm = reactive({ id:'', title:'', slug:'', category:'news', excerpt:'', content:'', cover_image_url:'', status:'draft', meta_title:'', meta_description:'', focus_keywords:'' })
const blogCategoryForm = reactive({ id:'', name:'', slug:'', description:'', sort_order:0, is_active:true })
const eventForm = reactive({ id:'', title:'', slug:'', category:'', location:'', event_date:'', description:'', status:'published' })
const facilityForm = reactive({ id:'', name:'', slug:'', category:'', description:'', image_url:'', sort_order:0, is_active:true })
const associationForm = reactive({ id:'', name:'', slug:'', category:'', president:'', secretary:'', phone:'', website_url:'', description:'', logo_url:'', sort_order:0, is_active:true })
const factForm = reactive({ id:'', label:'', value:'', icon:'icofont-chart-growth', sort_order:0, is_active:true })
const faqForm = reactive({ id:'', question:'', answer:'', category:'General', sort_order:0, is_active:true })
const resourceForm = reactive({ id:'', title:'', category:'Guidelines', file_url:'', description:'', sort_order:0, is_active:true })

const blogCategoryFields = fields(['name','slug','sort_order','is_active'], ['description'])
const eventFields = fields(['title','slug','category','location','event_date','status'], ['description'])
const facilityFields = fields(['name','slug','category','image_url','sort_order','is_active'], ['description'])
const associationFields = fields(['name','slug','category','president','secretary','phone','website_url','logo_url','sort_order','is_active'], ['description'])
const factFields = fields(['label','value','icon','sort_order','is_active'])
const faqFields = fields(['question','category','sort_order','is_active'], ['answer'])
const resourceFields = fields(['title','category','file_url','sort_order','is_active'], ['description'])

const overviewCards = computed(() => [
  { label:'Posts', value:posts.value.length },
  { label:'Events', value:events.value.length },
  { label:'Facilities', value:facilities.value.length },
  { label:'Associations', value:associations.value.length },
  { label:'Facts', value:facts.value.length },
  { label:'FAQs', value:faqs.value.length },
])

onMounted(() => {
  if (!localStorage.getItem('ncsms_access_token')) router.replace('/login')
  loadAll()
})

function fields(short = [], long = []) {
  return [...short.map(name => ({ name, type: name === 'is_active' ? 'checkbox' : 'input' })), ...long.map(name => ({ name, type:'textarea' }))]
}
function data(res) { return res?.data?.data ?? {} }
function setMsg(text) { message.value = text; error.value = ''; setTimeout(() => { message.value = '' }, 2500) }
function setErr(err) { error.value = err.response?.data?.error?.message || err.message || 'Action failed' }
function copyInto(target, source) { Object.keys(target).forEach(k => { target[k] = source?.[k] ?? (typeof target[k] === 'boolean' ? false : '') }) }
function clean(payload) { return Object.fromEntries(Object.entries(payload).filter(([k]) => k !== 'id')) }

async function loadAll() {
  try {
    const results = await Promise.allSettled([
      cms.getSettings('homepage'), cms.getSettings('contact'), cms.getSettings('footer'), cms.getMenu('main'), cms.getMenu('footer'), cms.adminGetSlideshow('homepage-hero'),
      cms.adminListPosts({ per_page:200 }), cms.adminListEvents({ per_page:200 }), cms.adminListSlides(), cms.adminListFacilities(), cms.adminListAssociations(),
      cms.adminListFunFacts(), cms.adminListFAQs(), cms.adminListResources({ per_page:200 }), cms.adminListBlogCategories(), cms.adminListComments({ status: commentStatus.value, per_page:50 }),
    ])
    Object.assign(homepage, data(results[0].value)?.value || {})
    Object.assign(contact, data(results[1].value)?.value || {})
    Object.assign(footer, data(results[2].value)?.value || {})
    mainMenuTree.value = normalizeMenuTree(data(results[3].value)?.items || [])
    footerMenuTree.value = normalizeMenuTree(data(results[4].value)?.items || [])
    const show = data(results[5].value) || {}
    Object.assign(slideshow, show)
    posts.value = data(results[6].value)?.items || []
    events.value = data(results[7].value)?.items || []
    slides.value = show.slides || data(results[8].value) || []
    facilities.value = data(results[9].value) || []
    associations.value = data(results[10].value) || []
    facts.value = data(results[11].value) || []
    faqs.value = data(results[12].value) || []
    resources.value = data(results[13].value) || []
    blogCategories.value = data(results[14].value) || []
    comments.value = data(results[15].value)?.items || []
  } catch (err) { setErr(err) }
}

async function saveHomepage() { try { await cms.adminUpdateSettings('homepage', JSON.parse(JSON.stringify(homepage))); setMsg('Homepage saved') } catch (err) { setErr(err) } }
async function saveMenus() { try { await cms.adminUpdateMenu('main', toCmsMenuItems(mainMenuTree.value)); await cms.adminUpdateMenu('footer', toCmsMenuItems(footerMenuTree.value)); setMsg('Menus saved') } catch (err) { setErr(err) } }
async function saveSettings() { try { await cms.adminUpdateSettings('contact', { ...contact }); await cms.adminUpdateSettings('footer', { ...footer }); setMsg('Settings saved') } catch (err) { setErr(err) } }
async function createSlide() { try { await cms.adminCreateSlide({ title:'New slide', subtitle:'National Council of Sports', description:'', image_url:'', button_text:'Learn More', button_url:'/', sort_order:slides.value.length + 1, is_active:true }); await loadAll(); setMsg('Slide added') } catch (err) { setErr(err) } }
function editSlide(item) { active.value = 'homepage'; Object.assign(homepage, { hero_quick_edit: item.title }) }
async function removeSlide(item) { if(confirm('Delete this slide?')) { await cms.adminDeleteSlide(item.id); await loadAll() } }

async function savePost() { await saveEntity(postForm, cms.adminCreatePost, cms.adminUpdatePost, 'Post saved') }
async function saveBlogCategory() { await saveEntity(blogCategoryForm, cms.adminCreateBlogCategory, cms.adminUpdateBlogCategory, 'Category saved') }
async function saveEvent() { await saveEntity(eventForm, cms.adminCreateEvent, cms.adminUpdateEvent, 'Event saved') }
async function saveFacility() { await saveEntity(facilityForm, cms.adminCreateFacility, cms.adminUpdateFacility, 'Facility saved') }
async function saveAssociation() { await saveEntity(associationForm, cms.adminCreateAssociation, cms.adminUpdateAssociation, 'Association saved') }
async function saveFact() { await saveEntity(factForm, cms.adminCreateFunFact, cms.adminUpdateFunFact, 'Fact saved') }
async function saveFAQ() { await saveEntity(faqForm, cms.adminCreateFAQ, cms.adminUpdateFAQ, 'FAQ saved') }
async function saveResource() { await saveEntity(resourceForm, cms.adminCreateResource, cms.adminUpdateResource, 'Resource saved') }
async function saveEntity(form, createFn, updateFn, ok) { try { form.id ? await updateFn(form.id, clean(form)) : await createFn(clean(form)); await loadAll(); setMsg(ok) } catch (err) { setErr(err) } }

function editPost(item) { copyInto(postForm, item); active.value = 'create-post' }
function editBlogCategory(item) { copyInto(blogCategoryForm, item); active.value = 'create-blog-categories' }
function editEvent(item) { copyInto(eventForm, item); if (item.event_date) eventForm.event_date = item.event_date.slice(0, 16) }
function editFacility(item) { copyInto(facilityForm, item) }
function editAssociation(item) { copyInto(associationForm, item) }
function editFact(item) { copyInto(factForm, item) }
function editFAQ(item) { copyInto(faqForm, item) }
function editResource(item) { copyInto(resourceForm, item) }

async function removePost(item) { await removeEntity(item, cms.adminDeletePost) }
async function removeBlogCategory(item) { await removeEntity(item, cms.adminDeleteBlogCategory) }
async function removeEvent(item) { await removeEntity(item, cms.adminDeleteEvent) }
async function removeFacility(item) { await removeEntity(item, cms.adminDeleteFacility) }
async function removeAssociation(item) { await removeEntity(item, cms.adminDeleteAssociation) }
async function removeFact(item) { await removeEntity(item, cms.adminDeleteFunFact) }
async function removeFAQ(item) { await removeEntity(item, cms.adminDeleteFAQ) }
async function removeResource(item) { await removeEntity(item, cms.adminDeleteResource) }
async function removeEntity(item, fn) { if (confirm('Delete this item?')) { await fn(item.id); await loadAll(); setMsg('Deleted') } }

function logout() {
  localStorage.removeItem('ncsms_access_token')
  localStorage.removeItem('ncsms_user')
  router.push('/login')
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('ncsms_sidebar_collapsed', String(sidebarCollapsed.value))
}

function expandSidebarOnNav() {
  if (sidebarCollapsed.value) {
    sidebarCollapsed.value = false
    localStorage.setItem('ncsms_sidebar_collapsed', 'false')
  }
}

function handleHomepageGroupFocusOut(event) {
  if (!event.currentTarget.contains(event.relatedTarget)) homepageGroupOpen.value = false
}

function handleBlogsGroupFocusOut(event) {
  if (!event.currentTarget.contains(event.relatedTarget)) blogsGroupOpen.value = false
}

function handleSlideshowGroupFocusOut(event) {
  if (!event.currentTarget.contains(event.relatedTarget)) slideshowGroupOpen.value = false
}

function resetPostForm() {
  Object.assign(postForm, { id:'', title:'', slug:'', category:'news', excerpt:'', content:'', cover_image_url:'', status:'draft', meta_title:'', meta_description:'', focus_keywords:'' })
}

function resetBlogCategoryForm() {
  Object.assign(blogCategoryForm, { id:'', name:'', slug:'', description:'', sort_order:0, is_active:true })
}

async function loadComments() {
  try {
    const res = await cms.adminListComments({ status: commentStatus.value, per_page:50 })
    comments.value = data(res)?.items || []
  } catch (err) { setErr(err) }
}

async function moderateComment(comment, action) {
  try {
    action === 'approve' ? await cms.adminApproveComment(comment.id) : await cms.adminFlagComment(comment.id)
    await loadComments()
    setMsg(`Comment ${action === 'approve' ? 'approved' : 'flagged'}`)
  } catch (err) { setErr(err) }
}

async function deleteComment(comment) {
  if (!confirm('Delete this comment permanently?')) return
  try {
    await cms.adminDeleteComment(comment.id)
    await loadComments()
    setMsg('Comment deleted')
  } catch (err) { setErr(err) }
}

const ContentTable = defineComponent({
  props: { items:Array, titleKey:String, subtitleKey:String },
  emits: ['edit', 'delete'],
  setup(props, { emit }) {
    return () => h('div', { class:'cms-table' }, (props.items || []).map(item => h('article', { class:'cms-table-row' }, [
      h('div', [h('strong', item[props.titleKey] || 'Untitled'), h('span', item[props.subtitleKey] || '')]),
      h('div', [h('button', { onClick:() => emit('edit', item) }, 'Edit'), h('button', { onClick:() => emit('delete', item) }, 'Delete')]),
    ])))
  },
})

const EditorForm = defineComponent({
  props: { title:String, model:Object, fields:Array },
  emits: ['save'],
  setup(props, { emit }) {
    return () => h('form', { class:'cms-editor', onSubmit:e => { e.preventDefault(); emit('save') } }, [
      h('div', { class:'cms-panel-head' }, [h('h2', props.title), h('button', { type:'submit' }, props.model.id ? 'Update' : 'Create')]),
      h('div', { class:'cms-two' }, props.fields.map(field => h('label', { class: field.type === 'textarea' ? 'wide' : '' }, [
        field.name.replaceAll('_', ' '),
        field.type === 'textarea'
          ? h('textarea', { value: props.model[field.name], onInput:e => props.model[field.name] = e.target.value })
          : field.type === 'checkbox'
            ? h('input', { type:'checkbox', checked: !!props.model[field.name], onChange:e => props.model[field.name] = e.target.checked })
            : h('input', { value: props.model[field.name], onInput:e => props.model[field.name] = e.target.value }),
      ]))),
    ])
  },
})
</script>

<style scoped>
.cms-shell{display:grid;grid-template-columns:18rem 1fr;min-height:100vh;transition:grid-template-columns 200ms ease}.cms-sidebar{position:sticky;top:0;height:100vh;background:#10233f;color:white;padding:1rem;display:flex;flex-direction:column;min-width:0;overflow:hidden}.cms-sidebar-head{display:flex;align-items:center;justify-content:space-between;gap:.5rem;margin-bottom:1rem}.cms-brand{display:flex;align-items:center;gap:.7rem;color:white;font-weight:800;min-width:0;overflow:hidden}.cms-brand img{flex:0 0 auto;width:2.6rem;height:2.6rem;object-fit:contain;background:white;border-radius:.35rem}.cms-collapse-toggle{flex:0 0 auto;display:flex;align-items:center;justify-content:center;width:2.25rem;height:2.25rem;border-radius:.4rem;background:rgb(255 255 255/.08);color:white;border:1px solid rgb(255 255 255/.14)}.cms-collapse-toggle:hover,.cms-collapse-toggle:focus-visible{background:#f5a623;color:#10233f}.cms-nav{display:grid;gap:.25rem;overflow:auto}.cms-nav button,.cms-logout{display:flex;align-items:center;gap:.65rem;border-radius:.45rem;padding:.65rem .75rem;color:rgb(255 255 255/.78);text-align:left}.cms-nav button.active,.cms-nav button:hover,.cms-nav button:focus-visible{background:#f5a623;color:#10233f}.cms-nav-group{display:grid}.cms-nav-parent{width:100%;justify-content:space-between}.cms-nav-submenu{display:grid;gap:.25rem;max-height:0;opacity:0;overflow:hidden;transform:translateZ(0);transition:max-height 220ms ease,opacity 180ms ease,padding 220ms ease;padding-left:.75rem}.cms-nav-group.expanded .cms-nav-submenu{max-height:24rem;opacity:1;padding-top:.25rem;padding-bottom:.25rem}.cms-nav-submenu button{font-size:.86rem;padding-left:1rem}.cms-logout{margin-top:auto;background:rgb(255 255 255/.08)}.cms-main{padding:2rem;min-width:0}.cms-header{display:flex;justify-content:space-between;gap:1rem;align-items:center;margin-bottom:1.5rem}.cms-kicker{text-transform:uppercase;letter-spacing:.18em;color:#d88700;font-size:.72rem;font-weight:800}.cms-header h1{font-size:2rem;font-weight:850;color:#1a365d}.cms-actions{display:flex;gap:.5rem;align-items:center}.cms-actions a,.cms-actions button,.cms-panel-head button,.cms-list-editor>button,.cms-editor button{border-radius:.45rem;background:#1a365d;color:white;padding:.6rem .9rem;font-size:.85rem;font-weight:700}.cms-grid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:1rem}.cms-card,.cms-panel{background:white;border:1px solid #e5e7eb;border-radius:.5rem;padding:1rem}.metric span{color:#64748b;font-size:.8rem}.metric strong{display:block;color:#1a365d;font-size:2rem}.cms-panel{display:grid;gap:1rem}.cms-panel-head{display:flex;justify-content:space-between;align-items:center;gap:1rem}.cms-panel h2{font-size:1.15rem;font-weight:800;color:#1a365d}.cms-menu-builder-grid{display:grid;gap:1rem}.cms-menu-builder-grid article{display:grid;gap:.75rem;border:1px solid #e5e7eb;border-radius:.5rem;padding:1rem}.cms-menu-builder-grid h3{font-weight:800;color:#1a365d}.cms-two{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:1rem}.cms-two label{display:grid;gap:.35rem;font-size:.76rem;font-weight:800;text-transform:capitalize;color:#475569}.cms-two .wide{grid-column:1/-1}.cms-two input,.cms-two textarea,.cms-list-editor textarea{width:100%;border:1px solid #cbd5e1;border-radius:.45rem;padding:.65rem;text-transform:none;font-weight:500;color:#111827}.cms-two textarea,.cms-list-editor textarea{min-height:6rem}.cms-list-editor{display:grid;gap:.75rem}.cms-row{display:grid;grid-template-columns:1fr auto;gap:.75rem}.cms-row button,.cms-table-row button{border:1px solid #d1d5db;border-radius:.4rem;padding:.45rem .7rem}.cms-table{display:grid;gap:.5rem}.cms-table-row{display:flex;justify-content:space-between;gap:1rem;align-items:center;border:1px solid #e5e7eb;border-radius:.45rem;padding:.75rem}.cms-table-row strong{display:block;color:#1a365d}.cms-table-row span{font-size:.8rem;color:#64748b}.cms-table-row div:last-child{display:flex;gap:.4rem}.cms-message{background:#ecfdf5;color:#047857;border:1px solid #a7f3d0;border-radius:.45rem;padding:.75rem}.cms-error{background:#fef2f2;color:#b91c1c;border:1px solid #fecaca;border-radius:.45rem;padding:.75rem}:global(.dark) .cms-main{background:#0f172a;color:#e5e7eb}:global(.dark) .cms-card,:global(.dark) .cms-panel,:global(.dark) .cms-menu-builder-grid article{background:#111827;border-color:#334155}:global(.dark) .cms-header h1,:global(.dark) .cms-panel h2,:global(.dark) .cms-menu-builder-grid h3,:global(.dark) .cms-table-row strong,:global(.dark) .metric strong{color:#f8fafc}:global(.dark) .cms-two input,:global(.dark) .cms-two textarea,:global(.dark) .cms-list-editor textarea{background:#0f172a;color:#f8fafc;border-color:#475569}@media(max-width:900px){.cms-shell{grid-template-columns:1fr}.cms-sidebar{position:relative;height:auto}.cms-grid,.cms-two{grid-template-columns:1fr}.cms-header{align-items:flex-start;flex-direction:column}}
@media(min-width:901px){.cms-shell.collapsed{grid-template-columns:4.5rem 1fr}.cms-shell.collapsed .cms-sidebar{padding:1rem .6rem}.cms-shell.collapsed .cms-brand span,.cms-shell.collapsed .cms-nav button span,.cms-shell.collapsed .cms-logout span,.cms-shell.collapsed .cms-nav-parent .icofont-rounded-down,.cms-shell.collapsed .cms-nav-submenu{display:none}.cms-shell.collapsed .cms-nav button,.cms-shell.collapsed .cms-nav-parent,.cms-shell.collapsed .cms-logout{justify-content:center}.cms-shell.collapsed .cms-sidebar-head{justify-content:center;flex-direction:column;gap:.6rem}}
.cms-panel-head select,.cms-two select{border:1px solid #cbd5e1;border-radius:.45rem;padding:.65rem;color:#111827;background:white}.cms-empty{color:#64748b;border:1px dashed #cbd5e1;border-radius:.45rem;padding:1rem;text-align:center}.comment-row span{max-width:54rem;display:block;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}:deep(.blog-editor){display:grid;gap:1rem}:deep(.inline-field){display:flex;gap:.5rem;align-items:center;text-transform:none}:deep(.inline-field input){width:auto}:deep(.counter){float:right;color:#64748b;font-weight:700}:deep(.counter.warn){color:#b45309}:deep(.editor-toolbar button){border-radius:.4rem;border:1px solid #cbd5e1;padding:.55rem .8rem;font-weight:800}:deep(.editor-toolbar){display:flex;flex-wrap:wrap;gap:.4rem}:deep(.editor-toolbar .active){background:#1a365d;color:white}:deep(.rich-editor){border:1px solid #cbd5e1;border-radius:.5rem;background:white;padding:1rem;min-height:18rem}:deep(.rich-editor .ProseMirror){min-height:16rem;outline:none}:deep(.rich-editor h2){font-size:1.5rem}:deep(.rich-editor ul){list-style:disc;padding-left:1.25rem}:deep(.rich-editor img){max-width:100%;border-radius:.5rem}:deep(.draft-state){font-size:.8rem;color:#64748b}
</style>
