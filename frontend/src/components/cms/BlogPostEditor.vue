<template>
  <form class="blog-editor" @submit.prevent="submit">
    <div class="cms-panel-head">
      <h2>{{ model.id ? 'Edit Blog Post' : 'Create New Post' }}</h2>
      <button type="submit">{{ model.id ? 'Update post' : 'Create post' }}</button>
    </div>

    <div class="cms-two">
      <label>Title<input v-model="model.title" required @input="syncSlug" /></label>
      <label>
        Category
        <select v-model="model.category">
          <option v-for="cat in categories" :key="cat.slug || cat.value" :value="cat.slug || cat.value">{{ cat.name || cat.label }}</option>
        </select>
      </label>
      <label>
        <span class="inline-field"><input v-model="manualSlug" type="checkbox" /> Manual slug override</span>
        <input v-model="model.slug" :readonly="!manualSlug" pattern="^[a-z0-9-_]+$" required />
      </label>
      <label>
        Workflow status
        <select v-model="model.status">
          <option value="draft">Draft</option>
          <option value="approved">Approved</option>
          <option value="published">Published</option>
        </select>
      </label>
      <label>SEO title<input v-model="model.meta_title" maxlength="180" /></label>
      <label>Focus keywords<input v-model="model.focus_keywords" placeholder="sports, federations, youth" /></label>
      <label class="wide">
        Meta description <span class="counter" :class="{ warn: model.meta_description.length >= 150 }">{{ model.meta_description.length }}/160</span>
        <textarea v-model="model.meta_description" maxlength="160"></textarea>
      </label>
      <label class="wide">Excerpt<textarea v-model="model.excerpt" maxlength="500"></textarea></label>
    </div>

    <DropzoneUpload
      v-model="model.cover_image_url"
      label="featured image"
      hint="PNG, JPG, or WEBP. Max 5MB."
      accept="image/png,image/jpeg,image/webp"
      :max-image-mb="5"
    />

    <div class="editor-toolbar" aria-label="Rich text formatting">
      <button type="button" @click="editor?.chain().focus().toggleBold().run()" :class="{ active: editor?.isActive('bold') }">B</button>
      <button type="button" @click="editor?.chain().focus().toggleItalic().run()" :class="{ active: editor?.isActive('italic') }">I</button>
      <button type="button" @click="editor?.chain().focus().toggleHeading({ level: 2 }).run()">H2</button>
      <button type="button" @click="editor?.chain().focus().toggleBulletList().run()">List</button>
      <button type="button" @click="addImage">Image</button>
    </div>
    <EditorContent :editor="editor" class="rich-editor" />
    <p class="draft-state" aria-live="polite">{{ draftState }}</p>
  </form>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { EditorContent, useEditor } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { debounce } from '@/utils/menuTree.js'
import { sanitizeRichHtml } from '@/utils/sanitize.js'
import DropzoneUpload from '@/components/cms/DropzoneUpload.vue'

const props = defineProps({ model: { type: Object, required: true }, categories: { type: Array, default: () => [] } })
const emit = defineEmits(['save'])
const manualSlug = ref(false)
const draftState = ref('')
const draftKey = 'ncs_blog_editor_draft'

const editor = useEditor({
  extensions: [StarterKit],
  content: props.model.content || '',
  onUpdate: ({ editor }) => {
    props.model.content = sanitizeRichHtml(editor.getHTML())
    saveDraft()
  },
})

const saveDraft = debounce(() => {
  localStorage.setItem(draftKey, JSON.stringify(props.model))
  draftState.value = 'Draft saved locally'
}, 450)

watch(() => props.model.id, async () => {
  await nextTick()
  editor.value?.commands.setContent(props.model.content || '', false)
  manualSlug.value = !!props.model.id
}, { flush: 'post' })

watch(props.model, saveDraft, { deep: true })

onMounted(() => {
  if (!props.model.id) {
    try {
      const draft = JSON.parse(localStorage.getItem(draftKey) || '{}')
      if (draft.title && !props.model.title) Object.assign(props.model, draft)
    } catch {}
  }
})

onBeforeUnmount(() => editor.value?.destroy())

function slugify(value) {
  return String(value || '').toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '')
}

function syncSlug() {
  if (!manualSlug.value) props.model.slug = slugify(props.model.title)
}

function addImage() {
  const url = window.prompt('Image URL')
  if (url) editor.value?.chain().focus().insertContent(`<img src="${url}" alt="">`).run()
}

function submit() {
  props.model.content = sanitizeRichHtml(editor.value?.getHTML() || props.model.content)
  emit('save')
}
</script>
