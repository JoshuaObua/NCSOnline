<template>
    <div class="col-xl-4 col-sm-6 col-sm-12">
                <div class="card browser-widget">
                  <div class="media card-body">
                    <div class="media-img"><img src="../../../assets/images/dashboard/chrome.png" alt=""></div>
                    <div class="media-body align-self-center">
                      <div>
                        <p>Daily </p>
                        <h4><span class="counter">
                          <number
                       class="bold counter"
                       ref="number1"
                       :from="0"
                       :to="36"
                       :duration="5"
                       :delay="0"
                       easing="Power1.easeOut" /></span>%</h4>
                      </div>
                      <div>
                        <p>Month </p>
                        <h4><span class="counter"> <number
                       class="bold counter"
                       ref="number1"
                       :from="0"
                       :to="96"
                       :duration="5"
                       :delay="0"
                       easing="Power1.easeOut" /></span>%</h4>
                      </div>
                      <div>
                        <p>Week </p>
                        <h4><span class="counter">
                           <number
                       class="bold counter"
                       ref="number1"
                       :from="0"
                       :to="46"
                       :duration="5"
                       :delay="0"
                       easing="Power1.easeOut" /></span>%</h4>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
</template>
<script>
export default {
    
}
</script>