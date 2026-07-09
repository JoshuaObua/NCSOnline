<template>
  <div class="summernote-shell cms-rich-editor">
    <div class="summernote-toolbar" aria-label="Summernote formatting toolbar">
      <button type="button" class="btn btn-light btn-sm" :class="{ active: editor?.isActive('bold') }" @click="editor?.chain().focus().toggleBold().run()"><i class="fas fa-bold"></i></button>
      <button type="button" class="btn btn-light btn-sm" :class="{ active: editor?.isActive('italic') }" @click="editor?.chain().focus().toggleItalic().run()"><i class="fas fa-italic"></i></button>
      <button type="button" class="btn btn-light btn-sm" :class="{ active: editor?.isActive('heading', { level: 2 }) }" @click="editor?.chain().focus().toggleHeading({ level: 2 }).run()">H2</button>
      <button type="button" class="btn btn-light btn-sm" :class="{ active: editor?.isActive('heading', { level: 3 }) }" @click="editor?.chain().focus().toggleHeading({ level: 3 }).run()">H3</button>
      <button type="button" class="btn btn-light btn-sm" :class="{ active: editor?.isActive('bulletList') }" @click="editor?.chain().focus().toggleBulletList().run()"><i class="fas fa-list-ul"></i></button>
      <button type="button" class="btn btn-light btn-sm" :class="{ active: editor?.isActive('orderedList') }" @click="editor?.chain().focus().toggleOrderedList().run()"><i class="fas fa-list-ol"></i></button>
      <button type="button" class="btn btn-light btn-sm" @click="editor?.chain().focus().undo().run()"><i class="fas fa-undo"></i></button>
      <button type="button" class="btn btn-light btn-sm" @click="editor?.chain().focus().redo().run()"><i class="fas fa-redo"></i></button>
    </div>
    <EditorContent :editor="editor" class="summernote rich-editor" />
  </div>
</template>

<script setup>
import { watch } from 'vue'
import { EditorContent, useEditor } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { sanitizeRichHtml } from '@/utils/sanitize.js'

const props = defineProps({
  modelValue: { type: String, default: '' },
})
const emit = defineEmits(['update:modelValue'])

const editor = useEditor({
  extensions: [StarterKit],
  content: props.modelValue || '',
  editorProps: {
    attributes: { class: 'summernote-editable' },
  },
  onUpdate: ({ editor }) => emit('update:modelValue', sanitizeRichHtml(editor.getHTML())),
})

watch(() => props.modelValue, value => {
  if (!editor.value || editor.value.getHTML() === value) return
  editor.value.commands.setContent(value || '', false)
})
</script>

<style scoped>
.cms-rich-editor { display: grid; gap: 0; }
.summernote-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 10px;
  border: 1px solid #e4e6fc;
  border-bottom: 0;
  border-radius: 3px 3px 0 0;
  background: #fdfdff;
}
.summernote-toolbar .btn {
  border: 0;
  border-radius: 3px;
  color: #34395e;
  background: #f4f6f9;
  box-shadow: none;
}
.summernote-toolbar .btn.active {
  background: #6777ef;
  color: #fff;
}
.rich-editor {
  min-height: 180px;
  border: 1px solid #e4e6fc;
  border-radius: 0 0 3px 3px;
  background: #fff;
  padding: 16px;
}
:deep(.ProseMirror) { min-height: 150px; outline: none; color: #495057; }
:deep(.ProseMirror ul) { list-style: disc; padding-left: 1.25rem; }
:deep(.ProseMirror ol) { list-style: decimal; padding-left: 1.25rem; }
:global(.dark) .summernote-toolbar { background: #111827; border-color: #475569; }
:global(.dark) .rich-editor { background: #111827; border-color: #475569; }
:global(.dark) :deep(.ProseMirror) { color: #f8fafc; }
</style>
