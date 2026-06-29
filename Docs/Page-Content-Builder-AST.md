# Page Content Builder AST

This builder uses one recursive, JSON-serializable node shape for every block:

```ts
type UUID = string

type BuilderNode = {
  id: UUID
  type:
    | 'layout.row'
    | 'layout.columns'
    | 'layout.column'
    | 'layout.card'
    | 'ui.typography'
    | 'ui.image'
    | 'ui.accordion'
    | 'ui.tabs'
    | 'form.dropdown'
    | 'form.checkbox'
    | 'form.toggle'
  settings: Record<string, unknown>
  attributes: Record<string, unknown>
  children: BuilderNode[]
}
```

Structural blocks may contain any node. Leaf nodes such as typography, images, and form fields must not contain children. Accordions and tabs store panel children in `attributes.panels[].children` and `attributes.tabs[].children`, so each disclosure body is still a full recursive node list.

## JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://ncs.local/schemas/page-builder-node.json",
  "title": "PageBuilderNode",
  "type": "object",
  "required": ["id", "type", "settings", "attributes", "children"],
  "additionalProperties": false,
  "properties": {
    "id": { "type": "string", "format": "uuid" },
    "type": {
      "type": "string",
      "enum": [
        "layout.row",
        "layout.columns",
        "layout.column",
        "layout.card",
        "ui.typography",
        "ui.image",
        "ui.accordion",
        "ui.tabs",
        "form.dropdown",
        "form.checkbox",
        "form.toggle"
      ]
    },
    "settings": {
      "type": "object",
      "additionalProperties": true,
      "properties": {
        "padding": { "$ref": "#/$defs/spacing" },
        "margin": { "$ref": "#/$defs/spacing" },
        "responsive": {
          "type": "object",
          "additionalProperties": { "type": "object" }
        }
      }
    },
    "attributes": {
      "type": "object",
      "additionalProperties": true,
      "properties": {
        "onClick": { "$ref": "#/$defs/action" },
        "panels": {
          "type": "array",
          "items": { "$ref": "#/$defs/disclosurePanel" }
        },
        "tabs": {
          "type": "array",
          "items": { "$ref": "#/$defs/disclosurePanel" }
        }
      }
    },
    "children": {
      "type": "array",
      "items": { "$ref": "#" }
    }
  },
  "$defs": {
    "spacing": {
      "type": "object",
      "properties": {
        "all": { "type": ["integer", "string"] },
        "x": { "type": ["integer", "string"] },
        "y": { "type": ["integer", "string"] },
        "top": { "type": ["integer", "string"] },
        "right": { "type": ["integer", "string"] },
        "bottom": { "type": ["integer", "string"] },
        "left": { "type": ["integer", "string"] }
      }
    },
    "action": {
      "type": "object",
      "required": ["type"],
      "properties": {
        "type": {
          "enum": ["link_to", "toggle_visibility", "submit_form", "trigger_api_call"]
        },
        "target": { "type": "string" },
        "url": { "type": "string" },
        "method": { "type": "string" },
        "body": { "type": "object" }
      }
    },
    "disclosurePanel": {
      "type": "object",
      "required": ["id", "title", "children"],
      "properties": {
        "id": { "type": "string", "format": "uuid" },
        "title": { "type": "string" },
        "children": { "type": "array", "items": { "$ref": "#" } }
      }
    }
  }
}
```

## Rendering Architecture

The renderer is recursive:

```ts
function renderNode(node, context) {
  const classes = settingsToClasses(node.settings)

  switch (node.type) {
    case 'layout.row':
    case 'layout.column':
      return container(classes, node.children.map(child => renderNode(child, context)))

    case 'layout.columns':
      return grid(node.settings.columns, node.children.map(child => renderNode(child, context)))

    case 'layout.card':
      return card(node.attributes.header, node.children.map(child => renderNode(child, context)))

    case 'ui.accordion':
      return accordion(node.attributes.panels.map(panel => ({
        title: panel.title,
        body: panel.children.map(child => renderNode(child, context))
      })))

    case 'ui.tabs':
      return tabs(node.attributes.tabs.map(tab => ({
        title: tab.title,
        body: tab.children.map(child => renderNode(child, context))
      })))

    case 'form.dropdown':
    case 'form.checkbox':
    case 'form.toggle':
      return formControl(node, context.formState)

    default:
      return atomicElement(node)
  }
}
```

The implementation lives in:

- `frontend/src/page-builder/schema.js`
- `frontend/src/page-builder/tailwind.js`
- `frontend/src/page-builder/actions.js`
- `frontend/src/components/page-builder/PageBuilderRenderer.vue`
- `frontend/src/views/PageBuilderView.vue`

## Example Payload

The working example is exported as `exampleAst` in `frontend/src/page-builder/schema.js`. It renders:

`Row -> Columns -> Column(Image) + Column(Card -> Accordion -> Row -> Typography + Dropdown + Checkbox + Toggle)`
