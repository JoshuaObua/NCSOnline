<template>
  <div class="container-fluid dashboard-2-main">
    <div class="row">
      <div class="col-xl-6 box-col-12 des-xl-100 invoice-sec">
        
        <div class="card">
          <div class="card-header">
            <div
              class="
                header-top
                d-sm-flex
                justify-content-between
                align-items-center
              "
            >
              <h5>Invoice Overview</h5>
              <div class="center-content">
                <p class="d-sm-flex align-items-center">
                  <span class="m-r-10">$5,56548k</span
                  ><i class="toprightarrow-primary fa fa-arrow-up m-r-10 d-grid"></i
                  >94% More Than Last Year
                </p>
              </div>
            </div>
          </div>
          <div class="card-body p-0">
           <invoiceOver/>
          </div>
        </div>
      </div>
      <topDealer/>
      <ecormmerceIncome/>
      <activityTimeline/>
      <yearlyGrowth/>
     
    </div>
  </div>
</template>
<script>
import invoiceOver from '../e-commerce/invoiceOver.vue'
import ecormmerceIncome from '../e-commerce/ecormmerceIncome.vue'
import yearlyGrowth from '../e-commerce/yearlyGrowth.vue'
import topDealer from '../e-commerce/topDealer.vue'
import activityTimeline from '../e-commerce/activityTimeline.vue'

export default {
components:{
    invoiceOver,
    ecormmerceIncome,
     yearlyGrowth,
     topDealer,
     activityTimeline
}
};
</script>