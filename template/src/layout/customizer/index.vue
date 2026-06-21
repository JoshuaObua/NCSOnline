<template>
    <div class="customizer-links " :class="{ open: customizer }">
        <div class="nav flex-column nac-pills " id="c-pills-tab" role="tablist" aria-orientation="vertical">
            <customizertool/>
            <customizersetting/>
            <customizercolor/>
        </div>
    <div class="customizer-contain" :class="{ open: customizer }">
      <div class="tab-content" id="c-pills-tabContent">
          <customizerconfiguration/>
          <div class="customizer-body custom-scrollbar ">
              <customlayout/>
              <customsetting/>
              <customcolor/>
          </div>
      </div>
    </div></div>
</template>
<script>
 import { mapState } from 'vuex';
 import customizertool from "./customizertool.vue"
 import customizersetting from "./customizersetting.vue"
 import customizercolor from "./customizercolor.vue"
 import customizerconfiguration from "./customizerconfiguration.vue"
 import customlayout from "./customlayout.vue"
 import customsetting from "./customsetting.vue"
 import customcolor from "./customcolor.vue"
export default {
    components:{
       customizertool,
       customizersetting,
       customizercolor,
       customizerconfiguration,
       customlayout,
       customsetting,
       customcolor
    },
     computed: {
      ...mapState({
        customizer: (state) => state.menu.customizer,
      }),
    },
}
</script>