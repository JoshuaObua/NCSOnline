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
import { sanitizeRichHtml } from '@/utils/sanitize.js'

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
  const props = block.props || {}
  const children = () => (block.children || []).map(child => renderBlock(child))
  if (props.text || props.html || block.children) {
    if (block.type === 'section') {
      return h('section', { class:'builder-section', style:{ background:props.background, padding:props.padding, margin:props.margin } }, children().length ? children() : [
        textNode('h2', props.title || props.text),
        textNode('p', props.text, { class:'builder-copy' }),
      ])
    }
    if (block.type === 'grid' || block.type === 'sidebar_layout') {
      return h('section', { class:'builder-grid', style:{ gridTemplateColumns:props.template || '1fr 1fr', gap:props.gap || '1rem' } }, children())
    }
    if (block.type === 'column') return h('div', { class:'builder-column' }, children())
    if (block.type === 'card') return h('article', { class:'builder-card' }, [textNode('h3', props.title), ...children(), textNode('footer', props.footer)])
    if (block.type === 'accordion') return h('section', { class:'builder-accordion' }, children())
    if (block.type === 'accordion_panel') return h('details', { open:props.open !== false }, [h('summary', props.title || 'Pane'), ...children()])
    if (block.type === 'tabs') return h('section', { class:'builder-tabs' }, children())
    if (block.type === 'tab_panel') return h('article', { class:'builder-tab-panel' }, [textNode('h3', props.title), ...children()])
    if (['h1','h2','h3','h5'].includes(block.type)) return textNode(block.type, props.text, { class:`builder-${block.type}`, style:{ textAlign:props.align || 'left' } })
    if (block.type === 'paragraph') return textNode('p', props.text, { class:props.lead ? 'builder-lead' : 'builder-copy' })
    if (block.type === 'plain_text') return textNode('div', props.text, { class:'builder-plain-text' })
    if (block.type === 'list') {
      const tag = props.ordered === true || props.ordered === 'true' ? 'ol' : 'ul'
      const items = String(props.items || '').split(/\r?\n/).map(item => item.trim()).filter(Boolean)
      return h(tag, { class:['builder-list', tag === 'ol' ? 'is-ordered' : ''] }, items.map(item => h('li', item)))
    }
    if (block.type === 'blockquote') return h('blockquote', { class:'builder-quote' }, [textNode('p', props.text), textNode('cite', props.cite)])
    if (block.type === 'summernote_text' || block.type === 'raw_html') return h('section', { class:'builder-html', innerHTML:sanitizeRichHtml(props.html || '') })
    if (block.type?.startsWith('button_') || block.type === 'fab') return h('a', { class:['builder-button', `builder-button-${props.variant || 'primary'}`], href:props.href || '#' }, [props.icon ? h('i', { class:props.icon }) : null, props.label || props.icon || 'Action'])
    if (block.type === 'button_group') return h('div', { class:['builder-button-group', props.direction === 'vertical' ? 'vertical' : ''] }, children())
    if (block.type === 'image') return h('figure', { class:'builder-image' }, [props.src ? h('img', { src:mediaUrl(props.src), alt:props.alt || '', loading:props.lazy === false ? 'eager' : 'lazy' }) : null, textNode('figcaption', props.caption)])
    if (block.type === 'carousel') return h('section', { class:'builder-carousel' }, children())
    if (block.type === 'video') return props.provider === 'html5' ? h('video', { class:'builder-video', src:mediaUrl(props.source), controls:props.controls !== false }) : h('iframe', { class:'builder-video', src:props.source, loading:'lazy', allowfullscreen:true })
    if (block.type === 'audio') return h('audio', { class:'builder-audio', src:mediaUrl(props.source), controls:true })
    if (block.type === 'icon_block') return h('div', { class:'builder-icon-block' }, h('i', { class:props.icon || 'icofont-star', style:{ fontSize:props.size || '48px' } }))
    if (block.type === 'divider') return h('div', { class:'builder-divider', style:{ borderTopStyle:props.style || 'solid', borderTopWidth:props.weight || '1px', margin:`${props.spacing || '24px'} 0` } }, props.centerIcon ? h('i', { class:props.centerIcon }) : null)
    if (block.type === 'alert') return h('aside', { class:['builder-alert', `tone-${props.tone || 'info'}`] }, props.text || 'Notification')
    if (block.type === 'progress') return h('div', { class:'builder-progress' }, [textNode('span', props.label), h('div', h('i', { style:{ width:`${props.value || 0}%` } }))])
  }
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
    // Only static pages that have been published are shown on the public site;
    // drafts/approved-but-unpublished pages resolve to "Page not found".
    if (post?.category !== 'page' || post?.status !== 'published') {
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
.builder-plain-text{white-space:pre-wrap;color:#64748b;line-height:1.8}.builder-list{display:grid;gap:.8rem;margin:0;padding-left:1.35rem;color:#475569;line-height:1.65;list-style:disc}.builder-list.is-ordered{list-style:decimal}.builder-list li::marker{color:#f48c06;font-weight:900}
.page-builder-render{display:grid;gap:2rem}.builder-section{display:grid;gap:.7rem}.builder-section h1,.builder-h1{color:#112b4e;font-size:clamp(2.2rem,4vw,3.6rem);font-weight:900;line-height:1.05}.builder-section h2,.builder-h2{color:#112b4e;font-size:clamp(1.8rem,3vw,2.55rem);font-weight:800;line-height:1.12}.builder-h3{color:#112b4e;font-size:1.45rem;font-weight:800}.builder-h5{color:#475569;font-size:.95rem;font-weight:900;text-transform:uppercase}.builder-kicker{width:max-content;border-radius:999px;background:rgb(245 166 35/.12);padding:.35rem .8rem;color:#d88700;font-size:.78rem;font-weight:800;text-transform:uppercase}.builder-copy,.builder-lead,.builder-card p,.builder-icon-card p,.builder-dropdown p,.builder-accordion p{color:#64748b;line-height:1.8}.builder-lead{font-size:1.2rem;color:#334155}.builder-grid,.builder-row{display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:1rem}.builder-column{display:grid;gap:1rem;min-width:0}.builder-card,.builder-icon-card,.builder-tab-panel{border:1px solid #e8edf3;border-radius:1rem;background:#fff;padding:1.35rem;box-shadow:0 8px 24px rgb(15 31 61/.06)}.builder-card h3,.builder-icon-card h3,.builder-tab-panel h3{margin:.7rem 0;color:#112b4e;font-weight:800}.builder-card footer{margin-top:1rem;color:#94a3b8;font-size:.85rem}.builder-card-icon{display:inline-flex;color:#f5a623;font-size:2rem}.builder-image{overflow:hidden;border-radius:1rem;background:#f8fafc;box-shadow:0 12px 32px rgb(15 31 61/.08)}.builder-image img{width:100%;max-height:32rem;object-fit:cover}.builder-image figcaption{padding:.9rem 1rem;color:#64748b;font-size:.88rem}.builder-accordion,.builder-dropdown,.builder-tabs{display:grid;gap:.75rem}.builder-accordion details,.builder-dropdown{border:1px solid #e8edf3;border-radius:.75rem;background:#f8fafc;padding:1rem}.builder-accordion summary,.builder-dropdown summary{cursor:pointer;color:#112b4e;font-weight:800}.builder-accordion p,.builder-dropdown p{margin-top:.8rem}.builder-quote{border-left:4px solid #f5a623;padding:1rem 1.25rem;background:#fff7ed;color:#334155;border-radius:.75rem}.builder-quote cite{display:block;margin-top:.6rem;color:#92400e;font-weight:800}.builder-button-group{display:flex;flex-wrap:wrap;gap:.5rem}.builder-button-group.vertical{flex-direction:column;align-items:flex-start}.builder-button{display:inline-flex;align-items:center;gap:.45rem;border-radius:.5rem;padding:.65rem 1rem;font-weight:900}.builder-button-primary{background:#112b4e;color:#fff}.builder-button-secondary{background:#e5e7eb;color:#111827}.builder-button-ghost{border:1px solid #112b4e;color:#112b4e}.builder-button-link{color:#f48c06;padding-left:0}.builder-video{width:100%;aspect-ratio:16/9;border:0;border-radius:1rem;background:#111827}.builder-audio{width:100%}.builder-icon-block{display:inline-grid;place-items:center;width:80px;height:80px;border-radius:1rem;background:#eef2ff;color:#3730a3}.builder-divider{border-top-color:#cbd5e1;text-align:center;color:#f48c06}.builder-alert{border-radius:.75rem;padding:1rem;font-weight:800}.tone-info{background:#eff6ff;color:#1d4ed8}.tone-success{background:#ecfdf5;color:#047857}.tone-warning{background:#fffbeb;color:#b45309}.tone-danger,.tone-critical{background:#fef2f2;color:#b91c1c}.builder-progress span{display:block;margin-bottom:.4rem;font-weight:800;color:#112b4e}.builder-progress div{height:10px;background:#e5e7eb;border-radius:999px;overflow:hidden}.builder-progress i{display:block;height:100%;background:#f48c06}.builder-html{overflow:auto}.builder-html :deep(img){max-width:100%;border-radius:.75rem}@media(max-width:720px){.builder-grid,.builder-row{grid-template-columns:1fr!important}}
</style>
