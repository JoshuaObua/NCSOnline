<template>
    <div class="col-xl-6 xl-50 box-col-12">
                <div class="widget-joins card">
                  <div class="row">
                    <div class="col-sm-6 pe-0">
                      <div class="media border-after-xs">
                        <div class="align-self-center me-3">68%<i class="fa fa-angle-up ms-2"> </i></div>
                        <div class="media-body details ps-3"><span class="mb-1">New</span>
                          <h5 class="mb-0 counter"><number class="bold transition" ref="number1" :from="0" :to="6982" :duration="5" :delay="0" easing="Power1.easeOut"/></h5>
                        </div>
                        <div class="media-body align-self-center"><i class="font-primary float-end ms-2" data-feather="shopping-bag"><vue-feather class="shopping" type="shopping-bag"></vue-feather></i></div>
                      </div>
                    </div>
                    <div class="col-sm-6 ps-0">
                      <div class="media">
                          <div class="align-self-center me-3">12%<i class="fa fa-angle-down ms-2"></i></div>
                        <div class="media-body details ps-3"><span class="mb-1">Pending</span>
                          <h5 class="mb-0 counter"><number class="bold transition" ref="number1" :from="0" :to="783" :duration="5" :delay="0" easing="Power1.easeOut"/></h5>
                        </div>
                        <div class="media-body align-self-center"><i class="font-primary float-end ms-3" data-feather="layers"><vue-feather class="shopping" type="layers"></vue-feather></i></div>
                      </div>
                    </div>
                    <div class="col-sm-6 pe-0">
                      <div class="media border-after-xs">
                        <div class="align-self-center me-3">68%<i class="fa fa-angle-up ms-2"></i></div>
                        <div class="media-body details ps-3 pt-0"><span class="mb-1">Done</span>
                          <h5 class="mb-0 counter"><number class="bold transition" ref="number1" :from="0" :to="3674" :duration="5" :delay="0" easing="Power1.easeOut"/></h5>
                        </div>
                        <div class="media-body align-self-center"><i class="font-primary float-end ms-2" data-feather="shopping-cart"><vue-feather class="shopping" type="shopping-cart"></vue-feather></i></div>
                      </div>
                    </div>
                    <div class="col-sm-6 ps-0">
                      <div class="media">
                        <div class="align-self-center me-3">68%<i class="fa fa-angle-up ms-2"></i></div>
                        <div class="media-body details ps-3 pt-0"><span class="mb-1">Cancel</span>
                          <h5 class="mb-0 counter"><number class="bold transition" ref="number1" :from="0" :to="69" :duration="5" :delay="0" easing="Power1.easeOut"/></h5>
                        </div>
                        <div class="media-body align-self-center"><i class="font-primary float-end ms-2" data-feather="dollar-sign"><vue-feather class="shopping" type="dollar-sign"></vue-feather></i></div>
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