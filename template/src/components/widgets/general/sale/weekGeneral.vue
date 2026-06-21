<template>
    <div class="col-sm-6 pe-0">
          <div class="media border-after-xs">
            <div class="align-self-center me-3 ">
              <span class="widget-t mb-1">Sale</span>
              <h5 class="mb-0">Week</h5>
            </div>
            <div class="media-body align-self-center">
              <i class="font-primary" data-feather="arrow-up"><vue-feather type="arrow-up"></vue-feather></i>
            </div>
            <div class="media-body counter1">
              <h5 class="mb-0">
                $<span class="counter"
                  ><number
                    class="bold transition"
                    ref="number1"
                    :from="0"
                    :to="63147"
                    :duration="5"
                    :delay="0"
                    easing="Power1.easeOut"
                /></span>
              </h5>
              <span class="mb-1">+$69(65%)</span>
            </div>
          </div>
        </div>
        <div class="col-sm-6 ps-0">
          <div class="media">
            <div class="align-self-center me-3 ">
              <span class="widget-t mb-1">Sale</span>
              <h5 class="mb-0">Year</h5>
            </div>
            <div class="media-body align-self-center ps-3">
              <i class="font-primary" data-feather="arrow-up"><vue-feather type="arrow-up"></vue-feather></i>
            </div>
            <div class="media-body ps-2 counter1">
              <h5 class="mb-0">
                $<span class="counter">
                  <number
                    class="bold transition"
                    ref="number1"
                    :from="0"
                    :to="963198"
                    :duration="5"
                    :delay="0"
                    easing="Power1.easeOut"
                /></span>
              </h5>
              <span class="mb-1">+$3654(90%) </span>
            </div>
          </div>
        </div>
</template>
<script>

export default {
    components:{
       
    }
}
</script>