<template>
  <div class="rich-editor border border-gray-200 rounded-lg overflow-hidden focus-within:ring-2 focus-within:ring-primary-400 focus-within:border-primary-400">
    <!-- Toolbar -->
    <div class="flex flex-wrap items-center gap-0.5 bg-gray-50 border-b border-gray-200 px-2 py-1.5">
      <button type="button" title="Bold" @click="editor?.chain().focus().toggleBold().run()"
        :class="editor?.isActive('bold') ? 'bg-gray-200 text-gray-900' : 'text-gray-500 hover:bg-gray-200 hover:text-gray-800'"
        class="px-2 py-1 rounded text-sm font-bold transition-colors">B</button>
      <button type="button" title="Italic" @click="editor?.chain().focus().toggleItalic().run()"
        :class="editor?.isActive('italic') ? 'bg-gray-200 text-gray-900' : 'text-gray-500 hover:bg-gray-200 hover:text-gray-800'"
        class="px-2 py-1 rounded text-sm italic transition-colors">I</button>
      <button type="button" title="Strike" @click="editor?.chain().focus().toggleStrike().run()"
        :class="editor?.isActive('strike') ? 'bg-gray-200 text-gray-900' : 'text-gray-500 hover:bg-gray-200 hover:text-gray-800'"
        class="px-2 py-1 rounded text-sm line-through transition-colors">S</button>
      <div class="w-px h-5 bg-gray-200 mx-1"></div>
      <button type="button" title="Heading 2" @click="editor?.chain().focus().toggleHeading({level:2}).run()"
        :class="editor?.isActive('heading',{level:2}) ? 'bg-gray-200 text-gray-900' : 'text-gray-500 hover:bg-gray-200 hover:text-gray-800'"
        class="px-2 py-1 rounded text-xs font-bold transition-colors">H2</button>
      <button type="button" title="Heading 3" @click="editor?.chain().focus().toggleHeading({level:3}).run()"
        :class="editor?.isActive('heading',{level:3}) ? 'bg-gray-200 text-gray-900' : 'text-gray-500 hover:bg-gray-200 hover:text-gray-800'"
        class="px-2 py-1 rounded text-xs font-bold transition-colors">H3</button>
      <div class="w-px h-5 bg-gray-200 mx-1"></div>
      <button type="button" title="Bullet List" @click="editor?.chain().focus().toggleBulletList().run()"
        :class="editor?.isActive('bulletList') ? 'bg-gray-200 text-gray-900' : 'text-gray-500 hover:bg-gray-200 hover:text-gray-800'"
        class="px-2 py-1 rounded text-sm transition-colors">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 6.75h12M8.25 12h12m-12 5.25h12M3.75 6.75h.007v.008H3.75V6.75zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zM3.75 12h.007v.008H3.75V12zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm-.375 5.25h.007v.008H3.75v-.008zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z"/>
        </svg>
      </button>
      <button type="button" title="Ordered List" @click="editor?.chain().focus().toggleOrderedList().run()"
        :class="editor?.isActive('orderedList') ? 'bg-gray-200 text-gray-900' : 'text-gray-500 hover:bg-gray-200 hover:text-gray-800'"
        class="px-2 py-1 rounded text-sm transition-colors">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 6.75h12M8.25 12h12m-12 5.25h12M3.75 6.75h.007v.008H3.75V6.75zM3.75 12h.007v.008H3.75V12zM3.75 17.25h.007v.008H3.75v-.008z"/>
        </svg>
      </button>
      <button type="button" title="Blockquote" @click="editor?.chain().focus().toggleBlockquote().run()"
        :class="editor?.isActive('blockquote') ? 'bg-gray-200 text-gray-900' : 'text-gray-500 hover:bg-gray-200 hover:text-gray-800'"
        class="px-2 py-1 rounded text-sm transition-colors">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 01.865-.501 48.172 48.172 0 003.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0012 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018z"/>
        </svg>
      </button>
      <div class="w-px h-5 bg-gray-200 mx-1"></div>
      <button type="button" title="Horizontal Rule" @click="editor?.chain().focus().setHorizontalRule().run()"
        class="px-2 py-1 rounded text-xs text-gray-500 hover:bg-gray-200 hover:text-gray-800 transition-colors">—</button>
      <button type="button" title="Clear formatting" @click="editor?.chain().focus().clearNodes().unsetAllMarks().run()"
        class="px-2 py-1 rounded text-xs text-gray-500 hover:bg-gray-200 hover:text-gray-800 transition-colors ml-auto">Clear</button>
    </div>

    <!-- Editor area -->
    <editor-content
      :editor="editor"
      class="prose prose-sm max-w-none min-h-[220px] px-4 py-3 text-gray-800 focus:outline-none [&_.tiptap]:outline-none [&_.tiptap]:min-h-[200px]"
    />
  </div>
</template>

<script setup>
import { watch, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'

const props = defineProps({
  modelValue: { type: String, default: '' }
})
const emit = defineEmits(['update:modelValue'])

const editor = useEditor({
  content: props.modelValue,
  extensions: [StarterKit],
  editorProps: {
    attributes: { class: 'tiptap focus:outline-none' }
  },
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  }
})

watch(() => props.modelValue, (val) => {
  if (editor.value && editor.value.getHTML() !== val) {
    editor.value.commands.setContent(val || '', false)
  }
})

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>
