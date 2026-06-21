<template>
    <div class="col-sm-6 col-xl-3 xl-25 col-lg-6 box-col-6">
                <div class="card social-widget-card">
                  <div class="card-body">
                    <div class="redial-social-widget radial-bar-70" data-label="50%"><i class="fa fa-linkedin font-primary"></i></div>
                    <h5 class="b-b-light">linkedin</h5>
                    <div class="row">
                      <div class="col text-center b-r-light"><span>Post</span>
                        <h4 class="counter mb-0"> 
                          <number class="bold counter" ref="number1" :from="0" :to="1234" :duration="5" :delay="0" easing="Power1.easeOut" /></h4>
                      </div>
                      <div class="col text-center"><span>linkedin</span>
                        <h4 class="counter mb-0"> 
                          <number class="bold counter" ref="number1" :from="0" :to="4369" :duration="5" :delay="0" easing="Power1.easeOut" /></h4>
                      </div>
                    </div>
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