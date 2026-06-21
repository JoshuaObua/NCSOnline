<template>
  <div class="col-xl-5">
    <div class="card o-hidden">
      <div class="card-header pb-0">
        <h5>Total Earning</h5>
      </div>
      <div class="bar-chart-widget">
        <div class="top-content bg-primary"></div>
        <div class="bottom-content card-body">
          <div class="row">
            <div class="col-12">
              <div id="chart-widget5">
              <totalEarning1/>
              </div>
            </div>
          </div>
          <div class="row text-center">
            <div class="col-4 b-r-light">
              <div>
                <span class="font-primary"
                  >12%<i class="icon-angle-up f-12 ms-1"></i></span
                ><span class="text-muted block-bottom">Year</span>
                <h4 class="num m-0">
                  <span class="counter color-bottom"> <number
                class="bold counter color-bottom"
                ref="number1"
                :from="0"
                :to="73"
                :duration="5"
                :delay="0"
                easing="Power1.easeOut"
              /></span>M
                </h4>
              </div>
            </div>
            <div class="col-4 b-r-light">
              <div>
                <span class="font-primary"
                  >15%<i class="icon-angle-up f-12 ms-1"></i></span
                ><span class="text-muted block-bottom">Month</span>
                <h4 class="num m-0">
                  <span class="counter color-bottom"> <number
                class="bold counter color-bottom"
                ref="number1"
                :from="0"
                :to="13"
                :duration="5"
                :delay="0"
                easing="Power1.easeOut"
              /></span>M
                </h4>
              </div>
            </div>
            <div class="col-4">
              <div>
                <span class="font-primary"
                  >34%<i class="icon-angle-up f-12 ms-1"></i></span
                ><span class="text-muted block-bottom">Today</span>
                <h4 class="num m-0">
                  <span class="counter color-bottom"> <number
                class="bold counter color-bottom"
                ref="number1"
                :from="0"
                :to="187"
                :duration="5"
                :delay="0"
                easing="Power1.easeOut"
              /></span>M
                </h4>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import totalEarning1 from './totalEarning1';
export default {
    components:{
        totalEarning1
    },
  data() {
    return {
      
  }
  
},
}
</script>