<template>
  <div class="col-xl-4">
    <div class="card o-hidden">
      <div class="chart-widget-top">
        <div class="row card-body">
          <div class="col-5">
            <h6 class="f-w-600 font-primary sales">SALE</h6>
            <span class="num sales"
              ><span class="counter"><number
                class="bold "
                ref="number1"
                :from="0"
                :to="90"
                :duration="5"
                :delay="0"
                easing="Power1.easeOut"
              /></span>%<i class="icon-angle-up f-12"></i
            ></span>
          </div>
          <div class="col-7 text-end">
            <h4 class="num total-value">
              $ <span class="counter">
                   <number
                class="bold "
                ref="number1"
                :from="0"
                :to="3654"
                :duration="5"
                :delay="0"
                easing="Power1.easeOut"
              /></span>.00
            </h4>
          </div>
        </div>
        <div>
          <div id="chart-widget1">
           <saleChart1/>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import saleChart1 from './saleChart1.vue'

export default {
    components:{
       saleChart1,
     
    },
    data(){
        return{
        }
    }
};
</script>