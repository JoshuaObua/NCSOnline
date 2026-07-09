<template>
  <form class="blog-editor otika-blog-editor" @submit.prevent="submit">
    <div class="card blog-editor-card">
      <div class="card-header">
        <h4>{{ editorTitle }}</h4>
        <div class="card-header-action">
          <button type="button" class="btn btn-icon icon-left btn-light" @click="restoreDraft"><i class="fas fa-history"></i> Draft</button>
          <button type="submit" class="btn btn-icon icon-left btn-primary"><i class="fas fa-save"></i> {{ model.id ? 'Update' : 'Create' }}</button>
        </div>
      </div>

      <div class="card-body">
        <div class="form-section">
          <div class="section-title mt-0">Post Details</div>
          <div class="row">
            <div class="form-group col-lg-8">
              <label>Title</label>
              <input v-model="model.title" class="form-control" required placeholder="Enter post title" @input="syncSlug" />
            </div>
            <div class="form-group col-lg-4">
              <label>Category</label>
              <select v-model="model[categoryField]" class="form-control selectric">
                <option value="">None</option>
                <option v-for="cat in categories" :key="cat.slug || cat.value" :value="cat.slug || cat.value">{{ cat.name || cat.label }}</option>
              </select>
            </div>
            <div class="form-group col-lg-8">
              <label>Slug</label>
              <input v-model="model.slug" class="form-control" :readonly="!manualSlug" pattern="^[a-z0-9_\-]+$" required placeholder="post-url-slug" />
            </div>
            <div class="form-group col-lg-4 blog-check-field">
              <label class="d-block">Slug Control</label>
              <div class="custom-control custom-checkbox">
                <input id="manual-slug-toggle" v-model="manualSlug" type="checkbox" class="custom-control-input" />
                <label class="custom-control-label" for="manual-slug-toggle">Manual slug override</label>
              </div>
            </div>
            <div class="form-group col-lg-4">
              <label>Workflow status</label>
              <select v-model="model.status" class="form-control selectric">
                <option value="draft">Draft</option>
                <option value="approved">Approved</option>
                <option value="published">Published</option>
              </select>
            </div>
            <div class="form-group col-lg-8">
              <label>Excerpt <span class="counter" :class="{ warn: model.excerpt.length >= 450 }">{{ model.excerpt.length }}/500</span></label>
              <textarea v-model="model.excerpt" class="form-control compact-textarea" maxlength="500" placeholder="Short summary shown on listing cards"></textarea>
            </div>
          </div>
        </div>

        <div class="form-section">
          <div class="section-title">Featured Media</div>
          <DropzoneUpload
            v-model="model.cover_image_url"
            label="featured image"
            hint="PNG, JPG, or WEBP. Max 5MB."
            accept="image/png,image/jpeg,image/webp"
            :max-image-mb="5"
          />
        </div>

        <div class="form-section">
          <div class="section-title">Content</div>
          <div class="summernote-shell">
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
          <p class="draft-state" aria-live="polite">{{ draftState }}</p>
        </div>

        <div class="form-section seo-card">
          <div class="section-title">SEO Meta</div>
          <div class="row">
            <div class="form-group col-lg-7">
              <label>SEO title</label>
              <input v-model="model.meta_title" class="form-control" maxlength="180" placeholder="Search result title" />
            </div>
            <div class="form-group col-lg-5">
              <label>Focus keywords</label>
              <div class="tag-input" @click="tagInput?.focus()">
                <span v-for="tag in keywordTags" :key="tag" class="badge badge-primary">
                  {{ tag }}
                  <button type="button" aria-label="Remove tag" @click.stop="removeTag(tag)">&times;</button>
                </span>
                <input
                  ref="tagInput"
                  v-model="tagDraft"
                  type="text"
                  placeholder="Type tag then Enter or Space"
                  @blur="commitTag"
                  @keydown.enter.prevent="commitTag"
                  @keydown.space.prevent="commitTag"
                  @keydown.tab="commitTag"
                />
              </div>
            </div>
            <div class="form-group col-12">
              <label>Meta description <span class="counter" :class="{ warn: model.meta_description.length >= 150 }">{{ model.meta_description.length }}/160</span></label>
              <textarea v-model="model.meta_description" class="form-control compact-textarea" maxlength="160" placeholder="Short SEO description for search engines"></textarea>
            </div>
          </div>
        </div>
      </div>

      <div class="card-footer text-right">
        <span class="draft-state footer-state">{{ draftState }}</span>
        <button type="submit" class="btn btn-primary mr-1">{{ model.id ? 'Update Post' : 'Create Post' }}</button>
      </div>
    </div>
  </form>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { EditorContent, useEditor } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { debounce } from '@/utils/menuTree.js'
import { sanitizeRichHtml } from '@/utils/sanitize.js'
import DropzoneUpload from '@/components/cms/DropzoneUpload.vue'

const props = defineProps({ model: { type: Object, required: true }, categories: { type: Array, default: () => [] }, categoryField: { type: String, default: 'category' } })
const emit = defineEmits(['save'])
const manualSlug = ref(false)
const draftState = ref('')
const tagDraft = ref('')
const tagInput = ref(null)
const categoryLabels = {
  page: 'Static Page',
  project: 'Project Post',
  case_study: 'Case Study',
  news: 'News Post',
  blog: 'Blog Post',
}
const editorTitle = computed(() => `${props.model.id ? 'Edit' : 'Create'} ${categoryLabels[props.model.category] || 'Post'}`)
const draftKey = computed(() => `ncs_blog_editor_draft_${props.model.category || 'post'}`)
const keywordTags = computed(() => splitTags(props.model.focus_keywords))

const editor = useEditor({
  extensions: [StarterKit],
  content: props.model.content || '',
  editorProps: {
    attributes: {
      class: 'summernote-editable',
    },
  },
  onUpdate: ({ editor }) => {
    props.model.content = sanitizeRichHtml(editor.getHTML())
    saveDraft()
  },
})

const saveDraft = debounce(() => {
  localStorage.setItem(draftKey.value, JSON.stringify(props.model))
  draftState.value = 'Draft saved locally'
}, 450)

watch(() => props.model.id, async () => {
  await nextTick()
  editor.value?.commands.setContent(props.model.content || '', false)
  manualSlug.value = !!props.model.id
  tagDraft.value = ''
}, { flush: 'post' })

watch(props.model, saveDraft, { deep: true })

onMounted(() => {
  ensureSummernoteAssets()
  if (!props.model.id) restoreDraft(false)
})

onBeforeUnmount(() => editor.value?.destroy())

function slugify(value) {
  return String(value || '').toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-+|-+$/g, '')
}

function syncSlug() {
  if (!manualSlug.value) props.model.slug = slugify(props.model.title)
}

function splitTags(value) {
  return String(value || '').split(',').map(tag => tag.trim()).filter(Boolean)
}

function commitTag() {
  const value = tagDraft.value.trim().replace(/,+$/g, '')
  if (!value) {
    tagDraft.value = ''
    return
  }
  const tags = new Set(keywordTags.value.map(tag => tag.toLowerCase()))
  if (!tags.has(value.toLowerCase())) {
    props.model.focus_keywords = [...keywordTags.value, value].join(', ')
  }
  tagDraft.value = ''
}

function removeTag(tag) {
  props.model.focus_keywords = keywordTags.value.filter(item => item !== tag).join(', ')
}

function restoreDraft(showMessage = true) {
  try {
    const draft = JSON.parse(localStorage.getItem(draftKey.value) || '{}')
    if (draft.title && (!props.model.title || showMessage)) {
      Object.assign(props.model, draft)
      editor.value?.commands.setContent(props.model.content || '', false)
      if (showMessage) draftState.value = 'Draft restored'
    }
  } catch {}
}

// Only the CSS is loaded — class names like .summernote-toolbar/.summernote-editable
// are reused purely for visual consistency with the original theme, but the actual
// editor below is TipTap (EditorContent), not Summernote. The real summernote-bs4.js
// bundle expects a global jQuery to attach itself to (`$.fn.summernote = ...`), which
// this app never loads, so including it here only threw an uncaught TypeError on
// every mount without doing anything.
function ensureSummernoteAssets() {
  addAsset('link', '/otika-assets/bundles/summernote/summernote-bs4.css')
  addAsset('link', '/otika-assets/bundles/jquery-selectric/selectric.css')
}

function addAsset(tag, href) {
  const attr = tag === 'link' ? 'href' : 'src'
  if (document.querySelector(`${tag}[${attr}="${href}"]`)) return
  const el = document.createElement(tag)
  el[attr] = href
  if (tag === 'link') el.rel = 'stylesheet'
  if (tag === 'script') el.async = true
  document.head.appendChild(el)
}

function submit() {
  commitTag()
  props.model.content = sanitizeRichHtml(editor.value?.getHTML() || props.model.content)
  emit('save')
}
</script>

<style scoped>
.otika-blog-editor {
  display: block;
}

.blog-editor-card {
  border: 0;
  border-radius: 3px;
  box-shadow: 0 4px 25px rgba(0, 0, 0, .1);
}

.blog-editor-card .card-header h4 {
  color: #34395e;
  font-size: 16px;
  font-weight: 700;
}

.form-section {
  padding-bottom: 18px;
}

.form-section + .form-section {
  border-top: 1px solid #f9f9f9;
  padding-top: 18px;
}

.section-title {
  color: #34395e;
  font-size: 14px;
  font-weight: 700;
}

.form-group label {
  color: #34395e;
  font-size: 12px;
  font-weight: 600;
}

.form-control {
  height: auto;
  border: 1px solid #e4e6fc;
  border-radius: 3px;
  background: #fdfdff;
  color: #495057;
  padding: 10px 15px;
  box-shadow: none;
}

.form-control:focus {
  border-color: #6777ef;
  box-shadow: 0 2px 6px #acb5f6;
}

.compact-textarea {
  min-height: 92px;
  resize: vertical;
}

.blog-check-field {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding-bottom: 8px;
}

.summernote-shell {
  overflow: hidden;
  border: 1px solid #e4e6fc;
  border-radius: 3px;
  background: #fff;
}

.summernote-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  border-bottom: 1px solid #e4e6fc;
  background: #f4f6f9;
  padding: 10px;
}

.summernote-toolbar .btn {
  min-width: 34px;
  border: 1px solid #e4e6fc;
  border-radius: 3px;
  background: #fff;
  color: #34395e;
  box-shadow: none;
}

.summernote-toolbar .btn.active {
  background: #6777ef;
  color: #fff;
}

.summernote {
  min-height: 320px;
  padding: 18px;
  background: #fff;
}

:deep(.summernote-editable) {
  min-height: 280px;
  outline: none;
  color: #34395e;
  line-height: 1.75;
}

:deep(.summernote-editable h2) {
  font-size: 1.55rem;
  font-weight: 800;
}

:deep(.summernote-editable h3) {
  font-size: 1.25rem;
  font-weight: 750;
}

:deep(.summernote-editable ul),
:deep(.summernote-editable ol) {
  padding-left: 1.3rem;
}

.seo-card {
  padding-bottom: 0;
}

.counter {
  float: right;
  color: #98a6ad;
  font-weight: 700;
}

.counter.warn {
  color: #ffa426;
}

.tag-input {
  display: flex;
  min-height: 44px;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  border: 1px solid #e4e6fc;
  border-radius: 3px;
  background: #fdfdff;
  padding: 6px 8px;
  cursor: text;
}

.tag-input:focus-within {
  border-color: #6777ef;
  box-shadow: 0 2px 6px #acb5f6;
}

.tag-input .badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border-radius: 30px;
  padding: 6px 10px;
  font-size: 11px;
}

.tag-input .badge button {
  border: 0;
  background: transparent;
  color: inherit;
  font-size: 14px;
  line-height: 1;
  padding: 0;
}

.tag-input input {
  min-width: 12rem;
  flex: 1 1 12rem;
  border: 0;
  background: transparent;
  color: #495057;
  outline: 0;
  padding: 4px;
}

.draft-state {
  margin: 10px 0 0;
  color: #98a6ad;
  font-size: 12px;
  font-weight: 600;
}

.footer-state {
  float: left;
  margin-top: 9px;
}

@media (max-width: 768px) {
  .blog-editor-card .card-header {
    align-items: flex-start;
    flex-direction: column;
    gap: 12px;
  }

  .card-header-action {
    width: 100%;
  }

  .card-header-action .btn {
    margin-bottom: 6px;
  }
}

:global(.dark) .blog-editor-card,
:global(.dark) .summernote-shell,
:global(.dark) .summernote {
  background: #1f2937 !important;
  color: #e5e7eb !important;
}

:global(.dark) .blog-editor-card .card-header h4,
:global(.dark) .section-title,
:global(.dark) .form-group label,
:global(.dark) :deep(.summernote-editable) {
  color: #f8fafc !important;
}

:global(.dark) .form-control,
:global(.dark) .tag-input,
:global(.dark) .tag-input input {
  background: #111827 !important;
  color: #f8fafc !important;
  border-color: #475569 !important;
}

:global(.dark) .summernote-toolbar {
  background: #111827 !important;
  border-color: #334155 !important;
}

:global(.dark) .summernote-toolbar .btn {
  background: #1f2937 !important;
  color: #e5e7eb !important;
  border-color: #475569 !important;
}

:global(.dark) .summernote-toolbar .btn.active {
  background: #6777ef !important;
  color: #fff !important;
}
</style>
