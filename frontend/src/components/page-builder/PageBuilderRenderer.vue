<template>
  <component :is="componentTag" v-bind="rootBindings" :class="nodeClasses" :style="nodeStyle" @click="handleClick">
    <template v-if="node.type === 'layout.columns'">
      <PageBuilderRenderer
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :form-state="formState"
        :action-context="actionContext"
      />
    </template>

    <template v-else-if="node.type === 'layout.card'">
      <header v-if="node.attributes.header" class="border-b border-gray-100 px-5 py-3 font-semibold text-gray-900">{{ node.attributes.header }}</header>
      <div class="space-y-4" :class="node.attributes.header || node.attributes.footer ? 'p-5' : ''">
        <PageBuilderRenderer v-for="child in node.children" :key="child.id" :node="child" :form-state="formState" :action-context="actionContext" />
      </div>
      <footer v-if="node.attributes.footer" class="border-t border-gray-100 px-5 py-3 text-sm text-gray-500">{{ node.attributes.footer }}</footer>
    </template>

    <template v-else-if="node.type === 'ui.typography'">
      <span v-if="node.attributes.rich" v-html="node.attributes.text"></span>
      <template v-else>{{ node.attributes.text }}</template>
    </template>

    <template v-else-if="node.type === 'ui.image'">
      <img
        :src="mediaUrl(node.attributes.src)"
        :alt="node.attributes.alt || ''"
        loading="lazy"
        class="h-full w-full"
        :class="imageClass"
        @error="imageError = true"
      />
      <div v-if="imageError" class="absolute inset-0 flex items-center justify-center bg-gray-100 text-sm text-gray-500">{{ node.attributes.fallback || 'Image unavailable' }}</div>
      <div v-if="node.attributes.overlayText" class="absolute inset-x-0 bottom-0 bg-black/60 p-3 text-sm font-semibold text-white">{{ node.attributes.overlayText }}</div>
    </template>

    <template v-else-if="node.type === 'ui.accordion'">
      <div class="space-y-2">
        <section v-for="panel in node.attributes.panels || []" :key="panel.id" class="rounded-lg border border-gray-200 bg-white">
          <button type="button" class="flex w-full items-center justify-between px-4 py-3 text-left font-semibold text-gray-900" @click.stop="togglePanel(panel.id)">
            <span>{{ panel.title }}</span>
            <i :class="openPanels.has(panel.id) ? 'icofont-rounded-up' : 'icofont-rounded-down'"></i>
          </button>
          <div v-if="openPanels.has(panel.id)" class="space-y-4 border-t border-gray-100 p-4">
            <PageBuilderRenderer v-for="child in panel.children || []" :key="child.id" :node="child" :form-state="formState" :action-context="actionContext" />
          </div>
        </section>
      </div>
    </template>

    <template v-else-if="node.type === 'ui.tabs'">
      <div>
        <div class="flex flex-wrap gap-2 border-b border-gray-200">
          <button
            v-for="tab in node.attributes.tabs || []"
            :key="tab.id"
            type="button"
            class="px-3 py-2 text-sm font-semibold"
            :class="activeTab === tab.id ? 'border-b-2 border-primary-700 text-primary-700' : 'text-gray-500'"
            @click.stop="activeTab = tab.id"
          >{{ tab.title }}</button>
        </div>
        <div class="space-y-4 py-4">
          <PageBuilderRenderer v-for="child in activeTabChildren" :key="child.id" :node="child" :form-state="formState" :action-context="actionContext" />
        </div>
      </div>
    </template>

    <template v-else-if="node.type === 'form.dropdown'">
      <label class="block text-sm font-medium text-gray-700">
        <span>{{ node.attributes.label }}</span>
        <select v-model="formState[node.attributes.name]" class="mt-1 w-full rounded-lg border-gray-300 text-sm" :multiple="node.attributes.multiple">
          <option v-if="node.attributes.clearable" value="">Select...</option>
          <option v-for="opt in node.attributes.options || []" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
        </select>
      </label>
    </template>

    <template v-else-if="node.type === 'form.checkbox' || node.type === 'form.toggle'">
      <label class="flex items-center gap-2 text-sm font-medium text-gray-700">
        <input v-model="formState[node.attributes.name]" type="checkbox" class="rounded border-gray-300 text-primary-700" />
        <span>{{ node.attributes.label }}</span>
      </label>
    </template>

    <template v-else>
      <PageBuilderRenderer v-for="child in node.children" :key="child.id" :node="child" :form-state="formState" :action-context="actionContext" />
    </template>
  </component>
</template>

<script setup>
import { computed, ref, watchEffect } from 'vue'
import { mediaUrl } from '@/api/client.js'
import { runAction } from '@/page-builder/actions.js'
import { layoutStyle, settingsToClasses } from '@/page-builder/tailwind.js'

const props = defineProps({
  node: { type: Object, required: true },
  formState: { type: Object, default: () => ({}) },
  actionContext: { type: Object, default: () => ({}) },
})

const imageError = ref(false)
const openPanels = ref(new Set())
const activeTab = ref('')

watchEffect(() => {
  const panels = props.node.attributes?.panels || []
  if (props.node.type === 'ui.accordion' && panels.length && openPanels.value.size === 0) {
    openPanels.value = new Set([panels[0].id])
  }
  const tabs = props.node.attributes?.tabs || []
  if (props.node.type === 'ui.tabs' && tabs.length && !activeTab.value) {
    activeTab.value = tabs[0].id
  }
})

const componentTag = computed(() => {
  if (props.node.type === 'ui.typography') return props.node.settings?.as || 'p'
  if (props.node.type === 'layout.columns') return 'section'
  if (props.node.type === 'ui.image') return 'figure'
  if (props.node.type.startsWith('form.')) return 'div'
  return props.node.attributes?.href ? 'a' : 'section'
})

const nodeClasses = computed(() => {
  const base = [settingsToClasses(props.node.settings)]
  if (props.node.type === 'layout.row') base.push(props.node.settings?.display === 'grid' ? 'grid' : props.node.settings?.display === 'flex' ? 'flex' : '')
  if (props.node.type === 'layout.columns') base.push('grid')
  if (props.node.type === 'layout.column') base.push('min-w-0 space-y-4')
  if (props.node.type === 'layout.card') base.push('overflow-hidden bg-white')
  if (props.node.type === 'ui.image') base.push('relative overflow-hidden bg-gray-100')
  return base.filter(Boolean).join(' ')
})

const nodeStyle = computed(() => {
  return layoutStyle(props.node.settings)
})

const rootBindings = computed(() => {
  const attrs = {}
  if (props.node.attributes?.href) attrs.href = props.node.attributes.href
  if (props.node.attributes?.target) attrs.target = props.node.attributes.target
  return attrs
})

const imageClass = computed(() => {
  const fit = props.node.settings?.objectFit || 'cover'
  return fit === 'contain' ? 'object-contain' : 'object-cover'
})

const activeTabChildren = computed(() => {
  return (props.node.attributes?.tabs || []).find(t => t.id === activeTab.value)?.children || []
})

function togglePanel(id) {
  const next = new Set(openPanels.value)
  if (next.has(id)) next.delete(id)
  else {
    if (!props.node.attributes?.allowMultiple) next.clear()
    next.add(id)
  }
  openPanels.value = next
}

function handleClick() {
  if (props.node.attributes?.onClick) runAction(props.node.attributes.onClick, props.actionContext)
}
</script>
