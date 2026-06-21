<template>
  <a
    class="nav-link"
    :class="{ 'active show': customizer == 'colors' }"
    id="c-pills-home-tab"
    @click="openCustomizerSetting('colors')"
    href="javascript:void(0)"
  >
    <div class="settings">
     <i class="icofont icofont-color-bucket"></i>
    </div>
    <span>Color Picker</span>
  </a>
</template>

<script>
  import { mapState } from 'vuex';

  export default {
    name: 'CustomSetting',
    computed: {
      ...mapState({
        customizer: (state) => state.menu.customizer,
      }),
    },
    methods: {
      openCustomizerSetting(val) {
        this.$store.state.menu.customizer = val;
      },
    },
  };
</script>
