<template>
  <div>
    <div v-if="loading" class="bg-cream pt-24 pb-16">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 animate-pulse">
        <div class="h-12 bg-gray-300 rounded mb-3 w-3/4"></div>
        <div class="h-4 bg-gray-300 rounded w-1/3"></div>
      </div>
    </div>

    <div v-else-if="page">
      <div class="relative min-h-[300px] flex items-end">
        <div v-if="page.cover_image_url" class="absolute inset-0 bg-cover bg-center" :style="`background-image: url('${mediaUrl(page.cover_image_url)}')`"></div>
        <div v-else class="absolute inset-0 bg-gradient-to-br from-[#112b4e] to-[#1e4080]"></div>
        <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent"></div>
        <div class="relative z-10 max-w-5xl mx-auto px-4 sm:px-6 pb-10 pt-28 w-full">
          <nav class="flex items-center gap-2 text-sm text-white/70 mb-4">
            <router-link to="/" class="hover:text-white transition-colors">Home</router-link>
            <span>/</span>
            <span class="text-white truncate max-w-xs">{{ page.title }}</span>
          </nav>
          <h1 class="text-3xl md:text-5xl font-bold text-white leading-tight">{{ page.title }}</h1>
          <p v-if="page.excerpt" class="mt-4 max-w-3xl text-white/80 text-lg">{{ page.excerpt }}</p>
        </div>
      </div>

      <div class="bg-white">
        <div class="max-w-5xl mx-auto px-4 sm:px-6 py-14">
          <div v-if="builderBlocks.length" class="page-builder-render">
            <PageBuilderBlock v-for="block in builderBlocks" :key="block.id || block.title" :block="block" />
          </div>
          <div v-else class="prose prose-gray max-w-none prose-headings:text-[#112b4e] prose-a:text-[#F48C06]" v-html="page.content"></div>
        </div>
      </div>
    </div>

    <div v-else class="bg-cream pt-32 pb-20 text-center text-gray-400">
      <p class="text-lg font-medium mb-4">Page not found.</p>
      <router-link to="/" class="text-sm font-semibold text-accent hover:text-[#d47b05]">Go to Home</router-link>
    </div>
  </div>
</template>

<script setup>
import { computed, h, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getPost } from '@/api/cms.js'
import { mediaUrl } from '@/api/client.js'

const route = useRoute()
const page = ref(null)
const loading = ref(true)
const builderBlocks = computed(() => parseBuilder(page.value?.content))
const PageBuilderBlock = { props: { block: { type:Object, required:true } }, setup(props) { return () => renderBlock(props.block) } }

function parseBuilder(content) {
  try {
    const parsed = JSON.parse(content || '{}')
    return parsed?.type === 'ncs-page-builder' && Array.isArray(parsed.blocks) ? parsed.blocks : []
  } catch {
    return []
  }
}

function textNode(tag, text, attrs = {}) {
  return text ? h(tag, attrs, text) : null
}

function renderBlock(block) {
  if (block.type === 'section') {
    return h('section', { class:'builder-section' }, [
      textNode('span', block.kicker, { class:'builder-kicker' }),
      textNode('h2', block.title),
      textNode('p', block.text, { class:'builder-copy' }),
    ])
  }
  if (block.type === 'image') {
    return h('figure', { class:'builder-image' }, [
      block.src ? h('img', { src:mediaUrl(block.src), alt:block.alt || block.caption || '' }) : null,
      textNode('figcaption', block.caption),
    ])
  }
  if (block.type === 'row') {
    return h('section', { class:'builder-row' }, (block.columns || []).map(column => h('article', { class:'builder-card' }, [
      column.icon ? h('i', { class:[column.icon, 'builder-card-icon'] }) : null,
      textNode('h3', column.title),
      textNode('p', column.text),
    ])))
  }
  if (block.type === 'accordion') {
    return h('section', { class:'builder-accordion' }, (block.items || []).map(item => h('details', [
      h('summary', item.title || 'Accordion item'),
      textNode('p', item.text),
    ])))
  }
  if (block.type === 'dropdown') {
    return h('details', { class:'builder-dropdown' }, [
      h('summary', block.label || 'Dropdown'),
      textNode('p', block.text),
    ])
  }
  if (block.type === 'icon-card') {
    return h('article', { class:'builder-icon-card' }, [
      h('i', { class:[block.icon || 'icofont-star', 'builder-card-icon'] }),
      textNode('h3', block.title),
      textNode('p', block.text),
    ])
  }
  if (block.type === 'html') {
    return h('section', { class:'builder-html', innerHTML:block.html || '' })
  }
  return h('section', { class:'builder-section' }, textNode('p', block.text || ''))
}

async function load(slug) {
  loading.value = true
  try {
    const res = await getPost(slug)
    const post = res.data.data
    if (post?.category !== 'page') {
      page.value = null
    } else {
      page.value = post
      document.title = post.meta_title || `${post.title} | NCS Uganda`
    }
  } catch {
    page.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => load(route.params.slug))
watch(() => route.params.slug, (slug) => slug && load(slug))
</script>

<style scoped>
.page-builder-render{display:grid;gap:2rem}.builder-section{display:grid;gap:.7rem}.builder-section h2{color:#112b4e;font-size:clamp(1.8rem,3vw,2.55rem);font-weight:800;line-height:1.12}.builder-kicker{width:max-content;border-radius:999px;background:rgb(245 166 35/.12);padding:.35rem .8rem;color:#d88700;font-size:.78rem;font-weight:800;text-transform:uppercase}.builder-copy,.builder-card p,.builder-icon-card p,.builder-dropdown p,.builder-accordion p{color:#64748b;line-height:1.8}.builder-image{overflow:hidden;border-radius:1rem;background:#f8fafc;box-shadow:0 12px 32px rgb(15 31 61/.08)}.builder-image img{width:100%;max-height:32rem;object-fit:cover}.builder-image figcaption{padding:.9rem 1rem;color:#64748b;font-size:.88rem}.builder-row{display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:1rem}.builder-card,.builder-icon-card{border:1px solid #e8edf3;border-radius:1rem;background:#fff;padding:1.35rem;box-shadow:0 8px 24px rgb(15 31 61/.06)}.builder-card h3,.builder-icon-card h3{margin:.7rem 0;color:#112b4e;font-weight:800}.builder-card-icon{display:inline-flex;color:#f5a623;font-size:2rem}.builder-accordion,.builder-dropdown{display:grid;gap:.75rem}.builder-accordion details,.builder-dropdown{border:1px solid #e8edf3;border-radius:.75rem;background:#f8fafc;padding:1rem}.builder-accordion summary,.builder-dropdown summary{cursor:pointer;color:#112b4e;font-weight:800}.builder-accordion p,.builder-dropdown p{margin-top:.8rem}.builder-html{overflow:auto}.builder-html :deep(img){max-width:100%;border-radius:.75rem}@media(max-width:720px){.builder-row{grid-template-columns:1fr}}
</style>
