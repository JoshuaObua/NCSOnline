<template>
  <div id="app">

    <router-view />
    <div class="loader-wrapper" v-if='showLoader'>
      <div class="theme-loader">
        <div class="loader-p"></div>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      showLoader: false
    }
  },
  watch: {
    $route() {
      this.showLoader = true;
      setTimeout(() => {
        this.showLoader = false
      }, 2000);
    }
  }
}
</script>
<style lang="scss">

</style>
