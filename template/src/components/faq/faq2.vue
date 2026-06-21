<template>
  <div class="faq-title">
    <h6>Selling And Buying</h6>
  </div>
  <div class="card">
    <div class="card-header">
      <h5 class="mb-0">
        <button class="btn btn-link collapsed" data-bs-toggle="collapse" data-bs-target="#collapseicon9"
          :aria-expanded="isActive ? 'true' : 'false'" v-on:click="open()"><vue-feather
            type="help-circle"></vue-feather> WordPress Site Maintenance ?</button>
      </h5>
    </div>
    <div :class="[isActive ? 'block' : 'none', 'content']" v-show="isActive">
      <div class="card-body">
        <p>We’ve just published a detailed on how to backup your WordPress website, however, if you’re in a hurry,
          here’s a quick solution.</p>
      </div>
    </div>
  </div>
  <div class="card">
    <div class="card-header">
      <h5 class="mb-0">
        <button class="btn btn-link collapsed" data-bs-toggle="collapse" data-bs-target="#collapseicon10"
          :aria-expanded="isActive1 ? 'true' : 'false'" aria-controls="collapseicon2" v-on:click="open1()"><vue-feather
            type="help-circle"></vue-feather>Meta Tags in WordPress ?</button>
      </h5>
    </div>
    <div :class="[isActive1 ? 'block' : 'none', 'content']" v-show="isActive1">
      <div class="card-body">
        <p>Manually adding meta tags in WordPress is relatively simple. For this demo, we’ll use the example from the
          WordPress Codex. Imagine you’re Harriet Smith, a veterinarian who blogs about their animal stories on
          WordPress.
        </p>
      </div>
    </div>
  </div>
  <div class="card">
    <div class="card-header">
      <h5 class="mb-0">
        <button class="btn btn-link collapsed" data-bs-toggle="collapse" data-bs-target="#collapseicon11"
          :aria-expanded="isActive2 ? 'true' : 'false'" aria-controls="collapseicon2" v-on:click="open2()"><vue-feather
            type="help-circle"></vue-feather>Validating a Website ?</button>
      </h5>
    </div>
    <div :class="[isActive2 ? 'block' : 'none', 'content']" v-show="isActive2">
      <div class="card-body">
        <p>Validating a website is the process of ensuring that the pages on the website conform to the norms or
          standards defined by various organizations.</p>
      </div>
    </div>
  </div>
  <div class="card">
    <div class="card-header">
      <h5 class="mb-0">
        <button class="btn btn-link collapsed" data-bs-toggle="collapse" data-bs-target="#collapseicon12"
          :aria-expanded="isActive3 ? 'true' : 'false'" aria-controls="collapseicon2" v-on:click="open3()"><vue-feather
            type="help-circle"></vue-feather>Know Your Sources ?</button>
      </h5>
    </div>
    <div :class="[isActive3 ? 'block' : 'none', 'content']" v-show="isActive3">
      <div class="card-body">
        <p>As of version 4.0, you can have WordPress automatically install the language of your choice during the
          installation process.</p>
      </div>
    </div>
  </div>
  <useraccounts />
</template>
<script>
import useraccounts from "./useraccounts.vue"
export default {
  components: {
    useraccounts
  },
  data() {
    return {
      isActive: false,
      isActive1: false,
      isActive2: false,
      isActive3: false,
      lorem: "Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus."
    }
  },
  methods: {
    open() {
      this.isActive = !this.isActive
    },
    open1() {
      this.isActive1 = !this.isActive1
    },
    open2() {
      this.isActive2 = !this.isActive2
    },
    open3() {
      this.isActive3 = !this.isActive3
    }
  }
}
</script>