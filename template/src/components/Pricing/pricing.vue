<template>
      <Breadcrumbs title="Pricing"/>
          <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header pb-0"> 
                    <h5>Simple Pricing Card</h5>
                  </div>
                  <div class="card-body row pricing-content">
                    <div class="col-xl-3 col-sm-6 xl-50 box-col-6">
                      <div class="card text-center pricing-simple">
                        <div class="card-body">
                          <h3>Standard</h3>
                          <h1>$15</h1>
                          <h6 class="mb-0">Standard Price</h6>
                        </div><a class="btn btn-lg btn-primary btn-block" href="javascript:void(0)">
                          <h5 class="mb-0">Purchase</h5></a>
                      </div>
                    </div>
                    <div class="col-xl-3 col-sm-6 xl-50 box-col-6">
                      <div class="card text-center pricing-simple">
                        <div class="card-body">
                          <h3>Business</h3>
                          <h1>$25</h1>
                          <h6 class="mb-0">Business Price</h6>
                        </div><a class="btn btn-lg btn-primary btn-block" href="javascript:void(0)">
                          <h5 class="mb-0">Purchase</h5></a>
                      </div>
                    </div>
                    <div class="col-xl-3 col-sm-6 xl-50 box-col-6">
                      <div class="card text-center pricing-simple">
                        <div class="card-body">
                          <h3>Premium</h3>
                          <h1>$35</h1>
                          <h6 class="mb-0">Premium Price</h6>
                        </div><a class="btn btn-lg btn-primary btn-block" href="javascript:void(0)">
                          <h5 class="mb-0">Purchase</h5></a>
                      </div>
                    </div>
                    <div class="col-xl-3 col-sm-6 xl-50 box-col-6">
                      <div class="card text-center pricing-simple">
                        <div class="card-body">
                          <h3>Extra</h3>
                          <h1>$45</h1>
                          <h6 class="mb-0">Extra Price</h6>
                        </div><a class="btn btn-lg btn-primary btn-block" href="javascript:void(0)">
                          <h5 class="mb-0">Purchase</h5></a>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <pricing2/>
            </div>
          </div>

</template>
<script>
import pricing2 from "./pricing2.vue"
export default {
    components:{
        pricing2
    }
}
</script>