<template>
 <div class="mode">
    <i
      class="fa fa-moon-o"
      v-show="mixLayout == 'light-only'"
      @click="customizeMixLayout('dark-only')"
    ></i>
    <i
      class="fa fa-lightbulb-o"
      v-show="mixLayout == 'dark-only'"
      @click="customizeMixLayout('light-only')"
    ></i>
  </div>
</template>

<script>

  export default {
    components:{
        
    },
    data() {
      return {
        darkMode: false,
         mixLayout: 'light-only',
         actives: false,
      };
    },
    methods: {
        customizeMixLayout(val) {
        this.mixLayout = val;
        this.$store.dispatch('layout/setLayout', val);
      },
    },
    mounted() {
    const savedMode = localStorage.getItem('mode');
    if (savedMode) {
      this.mixLayout = savedMode;
      this.actives = savedMode === 'dark-only';
      this.$store.dispatch('layout/setLayout', savedMode);
    }
  },
  watch: {
    mixLayout(newValue) {
      localStorage.setItem('mode', newValue);
    },
  },
    
  };
</script>
