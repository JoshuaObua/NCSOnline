<template>
  <div class="card">
    <div class="card-header pb-0">
      <h5>Default buttons</h5>
      <span>Add <code>.btn-pill</code> class for edge button</span>
    </div>
    <div class="card-body btn-showcase">
      <button class="btn btn-pill btn-primary" type="button">Primary Button</button>
      <button class="btn btn-pill btn-secondary" type="button">Secondary Button</button>
      <button class="btn btn-pill btn-success" type="button">Success Button</button>
      <button class="btn btn-pill btn-info" type="button">Info Button</button>
      <button class="btn btn-pill btn-warning" type="button">Warning Button</button>
      <button class="btn btn-pill btn-danger" type="button">Danger Button</button>
      <button class="btn btn-pill btn-light" type="button">Light Button</button>

    </div>
  </div>
  <div class="card">
    <div class="card-header pb-0">
      <h5>Large buttons</h5>
      <span>Add <code>.btn-pill</code> and <code>.btn-lg</code> class for large button</span>
    </div>
    <div class="card-body btn-showcase">
      <button class="btn btn-pill btn-primary btn-lg" type="button">Primary Button</button>
      <button class="btn btn-pill btn-secondary btn-lg" type="button">Secondary Button</button>
      <button class="btn btn-pill btn-success btn-lg" type="button">Success Button</button>
      <button class="btn btn-pill btn-info btn-lg" type="button">Info Button</button>
      <button class="btn btn-pill btn-warning btn-lg" type="button">Warning Button</button>
      <button class="btn btn-pill btn-danger btn-lg" type="button">Danger Button</button>
      <button class="btn btn-pill btn-light btn-lg" type="button">Light Button</button>

    </div>
  </div>
  <div class="card">
    <div class="card-header pb-0">
      <h5>Small buttons</h5>
      <span>Add <code>.btn-pill</code> and <code>.btn-sm</code> class for small button</span>
    </div>
    <div class="card-body btn-showcase">
      <button class="btn btn-pill btn-primary btn-sm" type="button">Primary Button</button>
      <button class="btn btn-pill btn-secondary btn-sm" type="button">Secondary Button</button>
      <button class="btn btn-pill btn-success btn-sm" type="button">Success Button</button>
      <button class="btn btn-pill btn-info btn-sm" type="button">Info Button</button>
      <button class="btn btn-pill btn-warning btn-sm" type="button">Warning Button</button>
      <button class="btn btn-pill btn-danger btn-sm" type="button">Danger Button</button>
      <button class="btn btn-pill btn-light btn-sm" type="button">Light Button</button>

    </div>
  </div>
  <div class="card">
    <div class="card-header pb-0">
      <h5>Active Buttons</h5>
      <span>Add <code>.active</code> for active state</span>
    </div>
    <div class="card-body btn-showcase">
      <button class="btn btn-pill btn-primary active" type="button">Active</button>
      <button class="btn btn-pill btn-secondary active" type="button">Active</button>
      <button class="btn btn-pill btn-success active" type="button">Active</button>
      <button class="btn btn-pill btn-info active" type="button">Active</button>
      <button class="btn btn-pill btn-warning active" type="button">Active</button>
      <button class="btn btn-pill btn-danger active" type="button">Active</button>
      <button class="btn btn-pill btn-light active txt-dark" type="button">Active</button>

    </div>
  </div>
  <div class="card">
    <div class="card-header pb-0">
      <h5>Disabled buttons</h5>
      <span>Add <code>.disabled</code> class or <code>disabled="disabled"</code> attribute for disabled button</span>
    </div>
    <div class="card-body btn-showcase">
      <button class="btn btn-pill btn-primary disabled" type="button">Disabled</button>
      <button class="btn btn-pill btn-secondary disabled" type="button">Disabled</button>
      <button class="btn btn-pill btn-success disabled" type="button">Disabled</button>
      <button class="btn btn-pill btn-info disabled" type="button">Disabled</button>
      <button class="btn btn-pill btn-warning disabled" type="button">Disabled</button>
      <button class="btn btn-pill btn-danger disabled" type="button">Disabled</button>
      <button class="btn btn-pill btn-light disabled" type="button">Disabled</button>

    </div>
  </div>
  <div class="card">
    <div class="card-header pb-0">
      <h5>Custom state buttons</h5>
      <span>The <code>.btn</code> class used with <code>&lt;button&gt;</code>, <code>&lt;a&gt;</code> and
        <code>&lt;input&gt;</code> elements.</span>
    </div>
    <div class="card-body btn-showcase"><a class="btn btn-pill btn-primary" href="javascript:void(0)"
        data-bs-toggle="tooltip" title="btn btn-primary" role="button">Link</a>
      <input class="btn btn-pill btn-secondary" type="button" value="Input" data-bs-toggle="tooltip"
        title="btn btn-secondary">
      <input class="btn btn-pill btn-success" type="submit" value="Submit" data-bs-toggle="tooltip"
        title="btn btn-success">
      <button class="btn btn-pill btn-info" type="submit">Button</button>

    </div>
  </div>
</template>
<script>
import { Tooltip } from 'bootstrap/dist/js/bootstrap.esm.min.js'
export default {
  mounted() {

    new Tooltip(document.body, {
      selector: "[data-bs-toggle='tooltip']",
    })
  }
}
</script>