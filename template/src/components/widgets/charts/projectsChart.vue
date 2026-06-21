<template>
  <div class="col-xl-4">
    <div class="card o-hidden">
      <div class="chart-widget-top">
        <div class="row card-body">
          <div class="col-7">
            <h6 class="f-w-600 font-secondary sales">PROJECTS</h6>
            <span class="num sales"
              ><span class="counter"
                ><number
                  class="bold"
                  ref="number1"
                  :from="0"
                  :to="30"
                  :duration="5"
                  :delay="0"
                  easing="Power1.easeOut" /></span
              >%<i class="icon-angle-up f-12"
                ></i
            ></span>
          </div>
          <div class="col-5 text-end">
            <h4 class="num total-value counter">
              <number
                class="bold counter"
                ref="number1"
                :from="0"
                :to="12569"
                :duration="5"
                :delay="0"
                easing="Power1.easeOut"
              />
            </h4>
          </div>
        </div>
        <div id="chart-widget2">
          <div
            class="flot-chart-placeholder"
            id="chart-widget-top-second"
          ></div>
          <projectsChart1 />
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import projectsChart1 from "./projectsChart1.vue";

export default {
  components: {
    projectsChart1,

  },
};
</script>