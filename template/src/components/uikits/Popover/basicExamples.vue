<template>
  <div class="col-sm-12 ">
    <div class="card">
      <div class="card-header pb-0">
        <h5>Basic Examples</h5>
        <span>lorem ipsum dolor sit amet, consectetur adipisicing elit</span>
      </div>
      <div class="card-body btn-showcase">
        <button class="btn btn-primary example-popover" type="button" ref="tolek" data-bs-toggle="popover"
          title="Popover title" data-bs-content="And here's some amazing content. It's very engaging. Right?">
          Click to toggle popover
        </button>
        <a class="example-popover btn btn-primary" tabindex="0" role="button" ref="tole" data-bs-toggle="popover"
          data-bs-trigger="focus" title="Popover title"
          data-bs-content="And here's some amazing content. It's very engaging. Right?">Dismissible popover</a>

        <button class="example-popover btn btn-primary" type="button" ref="to" data-bs-trigger="hover"
          data-container="body" data-bs-toggle="popover" data-bs-placement="bottom" title="Popover title"
          data-offset="-20px -20px" data-bs-content="Vivamus sagittis lacus vel augue laoreet rutrum faucibus.">
          On Hover Tooltip
        </button>
      </div>
    </div>
  </div>
  <div class="col-sm-12 ">
    <div class="card">
      <div class="card-header pb-0">
        <h5>Direction</h5>
        <span>lorem ipsum dolor sit amet, consectetur adipisicing elit</span>
      </div>
      <div class="card-body btn-showcase">
        <button class="example-popover btn btn-primary" type="button" ref="top" data-container="body"
          data-bs-toggle="popover" data-bs-placement="top" title="Popover title"
          data-bs-content="Vivamus sagittis lacus vel augue laoreet rutrum faucibus.">
          Popover on top
        </button>
        <button class="example-popover btn btn-primary" type="button" ref="right" data-container="body"
          data-bs-toggle="popover" data-bs-placement="right"
          data-bs-content="Vivamus sagittis lacus vel augue laoreet rutrum faucibus.">
          Popover on right
        </button>
        <button class="example-popover btn btn-primary" type="button" ref="bottom" data-container="body"
          data-bs-toggle="popover" data-bs-placement="bottom"
          data-bs-content="Vivamus sagittis lacus vel augue laoreet rutrum faucibus.">
          Popover on bottom
        </button>
        <button class="example-popover btn btn-primary" type="button" ref="left" data-container="body"
          data-bs-toggle="popover" data-bs-placement="left"
          data-bs-content="Vivamus sagittis lacus vel augue laoreet rutrum faucibus.">
          Popover on left
        </button>
      </div>
    </div>
  </div>
</template>
<script>
import { Popover } from 'bootstrap/dist/js/bootstrap.esm.min.js'
export default {
  mounted() {

    new Popover(this.$refs.tole)
    new Popover(this.$refs.to)
    new Popover(document.body, {
      selector: "[data-bs-toggle='popover']",
    })
  }
}
</script>