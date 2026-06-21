<template>
     <div class="col-xl-4 box-col-4">
                <div class="card testimonial text-center">
                  <div class="card-body">
                    <div>
                      <div class="item4"><i class="fa fa-quote-left"></i>
                        <p>Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock.</p>
                        <img class="img-80" src="../../../assets/images/dashboard/boy-2.png" alt="">
                        <h5 class="font-primary">Poio klot</h5><span>Developer</span>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="card xl-none">
                  <div class="ecommerce-widget card-body">
                    <div class="row">
                      <div class="col-6 "><span>New Order</span>
                        <h3 class="total-num counter"><number
                       class="bold counter"
                       ref="number1"
                       :from="0"
                       :to="25639"
                       :duration="5"
                       :delay="0"
                       easing="Power1.easeOut" /></h3>
                      </div>
                      <div class="col-6">
                        <div class="text-end">
                          <ul>
                            <li>Profit<span class="product-stts font-primary m-l-10">8989<i class="icon-angle-up f-12 m-l-10"></i></span></li>
                            <li>Loss<span class="product-stts font-primary m-l-10">2560<i class="icon-angle-down f-12 m-l-10"></i></span></li>
                          </ul>
                        </div>
                      </div>
                    </div>
                    <div class="progress-showcase">
                      <div class="progress sm-progress-bar">
                        <div class="progress-bar bg-primary" role="progressbar" style="width: 70%" aria-valuenow="25" aria-valuemin="0" aria-valuemax="100">                                               </div>
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
};
</script>