export const pageBuilderVersion = 2

export const blockRegistry = [
  { type: 'section', label: 'Main Section / Row', category: 'Layout', icon: 'icofont-layout', container: true, accepts: ['*'], props: { background: '#ffffff', padding: '48px 0', margin: '0', media: '' } },
  { type: 'grid', label: 'Grid / Columns', category: 'Layout', icon: 'icofont-columns', container: true, accepts: ['column'], props: { template: '1fr 1fr', gap: '24px' } },
  { type: 'column', label: 'Column', category: 'Layout', icon: 'icofont-split-h', container: true, accepts: ['*'], props: { width: '1fr' }, hidden: true },
  { type: 'sidebar_layout', label: 'Sidebar Layout', category: 'Layout', icon: 'icofont-sidebar', container: true, accepts: ['column'], props: { side: 'right', template: '2fr 1fr', gap: '24px' } },
  { type: 'card', label: 'Card Container', category: 'Layout', icon: 'icofont-card', container: true, accepts: ['*'], props: { title: 'Card title', footer: '', padding: '20px' } },
  { type: 'accordion', label: 'Accordion / Collapse', category: 'Layout', icon: 'icofont-list', container: true, accepts: ['accordion_panel'], props: { allowMultiple: false } },
  { type: 'accordion_panel', label: 'Accordion Pane', category: 'Layout', icon: 'icofont-rounded-down', container: true, accepts: ['*'], props: { title: 'Pane title', open: true }, hidden: true },
  { type: 'tabs', label: 'Tabs Container', category: 'Layout', icon: 'icofont-tabs', container: true, accepts: ['tab_panel'], props: { active: 0 } },
  { type: 'tab_panel', label: 'Tab Panel', category: 'Layout', icon: 'icofont-ui-folder', container: true, accepts: ['*'], props: { title: 'Tab' }, hidden: true },
  { type: 'h1', label: 'Main Title (H1)', category: 'Typography', icon: 'icofont-heading', props: { text: 'Main title', align: 'left' } },
  { type: 'h2', label: 'Section Heading (H2)', category: 'Typography', icon: 'icofont-heading', props: { text: 'Section heading', align: 'left' } },
  { type: 'h3', label: 'Sub-heading (H3)', category: 'Typography', icon: 'icofont-heading', props: { text: 'Sub-heading', align: 'left' } },
  { type: 'h5', label: 'Minor Heading (H5)', category: 'Typography', icon: 'icofont-heading', props: { text: 'Minor heading', align: 'left' } },
  { type: 'paragraph', label: 'Paragraph / Lead Text', category: 'Typography', icon: 'icofont-paragraph', props: { text: 'Write body copy here.', lead: false } },
  { type: 'blockquote', label: 'Blockquote', category: 'Typography', icon: 'icofont-quote-left', props: { text: 'Pull quote text.', cite: 'Source' } },
  { type: 'summernote_text', label: 'Rich Text Block', category: 'Typography', icon: 'icofont-edit', props: { html: '<p>Rich text content</p>', editor: 'summernote-inline' } },
  { type: 'button_primary', label: 'Primary Button', category: 'Buttons', icon: 'icofont-rounded-right', props: { label: 'Primary action', href: '#', variant: 'primary', icon: '', iconPosition: 'left' } },
  { type: 'button_secondary', label: 'Secondary Button', category: 'Buttons', icon: 'icofont-rounded-right', props: { label: 'Secondary action', href: '#', variant: 'secondary', icon: '', iconPosition: 'left' } },
  { type: 'button_ghost', label: 'Outline / Ghost Button', category: 'Buttons', icon: 'icofont-rounded-right', props: { label: 'Ghost action', href: '#', variant: 'ghost', icon: '', iconPosition: 'left' } },
  { type: 'button_link', label: 'Text / Link Button', category: 'Buttons', icon: 'icofont-link', props: { label: 'Text link', href: '#', variant: 'link' } },
  { type: 'button_icon', label: 'Icon Button', category: 'Buttons', icon: 'icofont-search', props: { label: 'Icon action', icon: 'icofont-search', shape: 'circle' } },
  { type: 'button_icon_enhanced', label: 'Icon-Enhanced Button', category: 'Buttons', icon: 'icofont-star', props: { label: 'Icon action', href: '#', icon: 'icofont-star', iconPosition: 'left', variant: 'primary' } },
  { type: 'fab', label: 'Floating Action Button', category: 'Buttons', icon: 'icofont-plus', props: { icon: 'icofont-plus', position: 'bottom-right' } },
  { type: 'button_group', label: 'Button Group Container', category: 'Buttons', icon: 'icofont-ui-clip-board', container: true, accepts: ['button_primary', 'button_secondary', 'button_ghost', 'button_link', 'button_icon', 'button_icon_enhanced'], props: { direction: 'horizontal', unifiedRadius: true } },
  { type: 'image', label: 'Single Image', category: 'Media', icon: 'icofont-image', props: { src: '', alt: '', lazy: true, lightbox: false, fit: 'cover' } },
  { type: 'carousel', label: 'Image Carousel / Slider', category: 'Media', icon: 'icofont-slidshare', container: true, accepts: ['image'], props: { transition: 'slide', autoplay: false } },
  { type: 'video', label: 'Video Player', category: 'Media', icon: 'icofont-video', props: { source: '', provider: 'html5', controls: true } },
  { type: 'icon_block', label: 'Icon Block', category: 'Media', icon: 'icofont-star', props: { icon: 'icofont-star', framed: true, size: '48px' } },
  { type: 'audio', label: 'Audio Player', category: 'Media', icon: 'icofont-audio', props: { source: '', title: 'Audio clip' } },
  { type: 'divider', label: 'Divider / Separator', category: 'Feedback', icon: 'icofont-minus', props: { style: 'solid', weight: '1px', spacing: '24px', centerIcon: '' } },
  { type: 'alert', label: 'Alert / Notification Banner', category: 'Feedback', icon: 'icofont-warning-alt', props: { tone: 'info', text: 'Notification text', dismissible: true } },
  { type: 'progress', label: 'Progress Bar', category: 'Feedback', icon: 'icofont-chart-flow', props: { value: 65, label: 'Progress' } },
]

export const visibleBlocks = blockRegistry.filter(block => !block.hidden)
export const blockMap = Object.fromEntries(blockRegistry.map(block => [block.type, block]))

export function createNode(type) {
  const spec = blockMap[type] || blockMap.paragraph
  const node = {
    id: `pb_${Date.now().toString(36)}_${Math.random().toString(36).slice(2, 7)}`,
    type: spec.type,
    props: structuredCloneSafe(spec.props || {}),
    children: [],
  }
  if (type === 'grid') node.children = [createNode('column'), createNode('column')]
  if (type === 'sidebar_layout') node.children = [createNode('column'), createNode('column')]
  if (type === 'accordion') node.children = [createNode('accordion_panel')]
  if (type === 'tabs') node.children = [createNode('tab_panel')]
  if (type === 'button_group') node.children = [createNode('button_primary'), createNode('button_ghost')]
  if (type === 'card') node.children = [createNode('paragraph')]
  if (type === 'section') node.children = [createNode('h2'), createNode('paragraph')]
  return node
}

export function canAccept(parentType, childType) {
  const parent = blockMap[parentType]
  if (!parent?.container) return false
  return parent.accepts?.includes('*') || parent.accepts?.includes(childType)
}

export function removeNode(tree, nodeId) {
  for (let i = 0; i < tree.length; i++) {
    if (tree[i].id === nodeId) return tree.splice(i, 1)[0]
    const found = removeNode(tree[i].children || [], nodeId)
    if (found) return found
  }
  return null
}

export function findNode(tree, nodeId) {
  for (const node of tree) {
    if (node.id === nodeId) return node
    const found = findNode(node.children || [], nodeId)
    if (found) return found
  }
  return null
}

export function insertNode(tree, parentId, node, index = -1) {
  const list = parentId ? findNode(tree, parentId)?.children : tree
  if (!list) return false
  if (index < 0 || index > list.length) list.push(node)
  else list.splice(index, 0, node)
  return true
}

export function duplicateNode(node) {
  const copy = structuredCloneSafe(node)
  const rekey = item => {
    item.id = `pb_${Date.now().toString(36)}_${Math.random().toString(36).slice(2, 7)}`
    ;(item.children || []).forEach(rekey)
  }
  rekey(copy)
  return copy
}

function structuredCloneSafe(value) {
  return typeof structuredClone === 'function' ? structuredClone(value) : JSON.parse(JSON.stringify(value))
}
