export const NODE_TYPES = [
  'layout.row',
  'layout.columns',
  'layout.column',
  'layout.card',
  'ui.typography',
  'ui.image',
  'ui.accordion',
  'ui.tabs',
  'form.dropdown',
  'form.checkbox',
  'form.toggle',
]

export const CONTAINER_TYPES = new Set([
  'layout.row',
  'layout.columns',
  'layout.column',
  'layout.card',
  'ui.accordion',
  'ui.tabs',
])

export const LEAF_TYPES = new Set([
  'ui.typography',
  'ui.image',
  'form.dropdown',
  'form.checkbox',
  'form.toggle',
])

export const STRUCTURAL_TYPES = new Set([
  'layout.row',
  'layout.columns',
  'layout.column',
  'layout.card',
])

export const ELEMENT_LABELS = {
  'layout.row': 'Row',
  'layout.columns': 'Columns',
  'layout.column': 'Column',
  'layout.card': 'Card',
  'ui.typography': 'Typography',
  'ui.image': 'Image',
  'ui.accordion': 'Accordion',
  'ui.tabs': 'Tabs',
  'form.dropdown': 'Dropdown',
  'form.checkbox': 'Checkbox',
  'form.toggle': 'Toggle',
}

export function uuid() {
  if (crypto?.randomUUID) return crypto.randomUUID()
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

export function createNode(type = 'ui.typography', overrides = {}) {
  const node = {
    id: uuid(),
    type,
    settings: {},
    attributes: {},
    children: [],
    ...overrides,
  }
  if (!node.settings) node.settings = {}
  if (!node.attributes) node.attributes = {}
  if (!Array.isArray(node.children)) node.children = []

  if (type === 'layout.row') {
    node.settings = { display: 'flex', gap: 6, padding: { top: 6, bottom: 6 }, ...node.settings }
  } else if (type === 'layout.columns') {
    node.settings = { columns: ['1fr', '1fr'], gap: 6, stackAt: 'md', ...node.settings }
    if (!node.children.length) {
      node.children = [createNode('layout.column'), createNode('layout.column')]
    }
  } else if (type === 'layout.card') {
    node.settings = { padding: { all: 5 }, border: true, shadow: 'sm', radius: 'lg', ...node.settings }
    node.attributes = { header: '', footer: '', ...node.attributes }
  } else if (type === 'ui.typography') {
    node.settings = { as: 'p', color: 'gray-700', weight: 'normal', leading: 'relaxed', ...node.settings }
    node.attributes = { text: 'Text content', rich: false, ...node.attributes }
  } else if (type === 'ui.image') {
    node.settings = { aspect: '16/9', objectFit: 'cover', radius: 'lg', ...node.settings }
    node.attributes = { src: '', alt: '', fallback: 'Image unavailable', ...node.attributes }
  } else if (type === 'ui.accordion') {
    node.attributes = { panels: [{ id: uuid(), title: 'Panel', children: [] }], allowMultiple: false, ...node.attributes }
  } else if (type === 'ui.tabs') {
    node.attributes = { tabs: [{ id: uuid(), title: 'Tab', children: [] }], ...node.attributes }
  } else if (type === 'form.dropdown') {
    node.attributes = { name: 'dropdown', label: 'Select option', options: [], multiple: false, clearable: true, ...node.attributes }
  } else if (type === 'form.checkbox') {
    node.attributes = { name: 'checkbox', label: 'Checkbox', defaultValue: false, ...node.attributes }
  } else if (type === 'form.toggle') {
    node.attributes = { name: 'toggle', label: 'Toggle', defaultValue: false, ...node.attributes }
  }
  return node
}

export function validateTree(root, path = 'root', seen = new Set()) {
  const errors = []
  if (!root || typeof root !== 'object') return [`${path} must be an object`]
  if (!root.id || typeof root.id !== 'string') errors.push(`${path}.id must be a UUID string`)
  if (seen.has(root.id)) errors.push(`${path}.id is duplicated`)
  seen.add(root.id)
  if (!NODE_TYPES.includes(root.type)) errors.push(`${path}.type is not supported`)
  if (!root.settings || typeof root.settings !== 'object' || Array.isArray(root.settings)) errors.push(`${path}.settings must be an object`)
  if (!root.attributes || typeof root.attributes !== 'object' || Array.isArray(root.attributes)) errors.push(`${path}.attributes must be an object`)
  if (!Array.isArray(root.children)) errors.push(`${path}.children must be an array`)
  if (LEAF_TYPES.has(root.type) && root.children?.length) errors.push(`${path} is a leaf node and cannot contain children`)
  ;(root.children || []).forEach((child, index) => errors.push(...validateTree(child, `${path}.children[${index}]`, seen)))
  ;(root.attributes?.panels || []).forEach((panel, index) => {
    ;(panel.children || []).forEach((child, childIndex) => errors.push(...validateTree(child, `${path}.attributes.panels[${index}].children[${childIndex}]`, seen)))
  })
  ;(root.attributes?.tabs || []).forEach((tab, index) => {
    ;(tab.children || []).forEach((child, childIndex) => errors.push(...validateTree(child, `${path}.attributes.tabs[${index}].children[${childIndex}]`, seen)))
  })
  return errors
}

export const exampleAst = createNode('layout.row', {
  settings: { display: 'block', padding: { top: 8, bottom: 8 }, background: 'gray-50' },
  children: [
    createNode('layout.columns', {
      settings: { columns: ['minmax(0, 1fr)', 'minmax(0, 1fr)'], gap: 8, stackAt: 'md' },
      children: [
        createNode('layout.column', {
          children: [
            createNode('ui.image', {
              attributes: {
                src: '/main-logo.png',
                alt: 'National Council of Sports',
                overlayText: 'National Sports Services',
              },
            }),
          ],
        }),
        createNode('layout.column', {
          children: [
            createNode('layout.card', {
              attributes: { header: 'Application preferences' },
              children: [
                createNode('ui.accordion', {
                  attributes: {
                    panels: [
                      {
                        id: uuid(),
                        title: 'Applicant details',
                        children: [
                          createNode('layout.row', {
                            children: [
                              createNode('ui.typography', {
                                settings: { as: 'h2', weight: 'bold', color: 'gray-900' },
                                attributes: { text: 'Tell us how to route your application' },
                              }),
                              createNode('form.dropdown', {
                                attributes: {
                                  name: 'sport_type',
                                  label: 'Sport category',
                                  options: [
                                    { label: 'Athletics', value: 'athletics' },
                                    { label: 'Football', value: 'football' },
                                    { label: 'Swimming', value: 'swimming' },
                                  ],
                                  validation: { required: true },
                                },
                              }),
                              createNode('form.checkbox', {
                                attributes: {
                                  name: 'confirm_accuracy',
                                  label: 'I confirm that this information is accurate',
                                  validation: { required: true },
                                },
                              }),
                              createNode('form.toggle', {
                                attributes: { name: 'receive_updates', label: 'Receive application updates' },
                              }),
                            ],
                          }),
                        ],
                      },
                    ],
                  },
                }),
              ],
            }),
          ],
        }),
      ],
    }),
  ],
})
