<template>
  <div class="col-xl-6 xl-50 box-col-12">
    <div class="widget-joins card widget-arrow">
      <div class="row">
        <div class="col-sm-6 pe-0">
          <div class="media border-after-xs">
            <div class="align-self-center me-3 ">
              <span class="widget-t mb-1">Sale</span>
              <h5 class="mb-0">Today</h5>
            </div>
            <div class="media-body align-self-center">
              <i class="font-primary" data-feather="arrow-down"> <vue-feather type="arrow-down"></vue-feather></i>
            </div>
            <div class="media-body counter1">
              <h5 class="mb-0">
                $<span class="counter"
                  ><number
                    class="bold transition"
                    ref="number1"
                    :from="0"
                    :to="25698"
                    :duration="5"
                    :delay="0"
                    easing="Power1.easeOut"
                /></span>
              </h5>
              <span class="mb-1">-$2658(36%)</span>
            </div>
          </div>
        </div>
        <div class="col-sm-6 ps-0">
          <div class="media">
            <div class="align-self-center me-3 ">
              <span class="widget-t mb-1">Sale</span>
              <h5 class="mb-0">Month</h5>
            </div>
            <div class="media-body align-self-center">
              <i class="font-primary" data-feather="arrow-up"><vue-feather type="arrow-up"></vue-feather></i>
            </div>
            <div class="media-body ps-2 counter1">
              <h5 class="mb-0">
                $<span class="counter"
                  ><number
                    class="bold transition"
                    ref="number1"
                    :from="0"
                    :to="6954"
                    :duration="5"
                    :delay="0"
                    easing="Power1.easeOut"
                /></span>
              </h5>
              <span class="mb-1">+$369(15%)</span>
            </div>
          </div>
        </div>
        <weekGeneral/>
      </div>
    </div>
  </div>
</template>
<script>

import weekGeneral from "./weekGeneral.vue"
export default {
    components:{
    
        weekGeneral
    }
};
</script>