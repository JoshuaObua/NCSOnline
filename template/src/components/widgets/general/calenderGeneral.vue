<template>
  <div class="col-xl-6 xl-100 box-col-12">
    <div class="card">
      <div class="cal-date-widget card-body">
        <div class="row">
          <div class="col-xl-6 col-xs-12 col-md-6 col-sm-6">
            <div class="cal-info text-center">
              <div>
                <h2>24</h2>
                <div class="d-inline-block">
                  <span class="b-r-dark pe-3">March</span
                  ><span class="ps-3">2018</span>
                </div>
                <p class="f-16">
                  There is no minimum donation, any sum is appreciated
                </p>
              </div>
            </div>
          </div>
          <div class="col-xl-6 col-xs-12 col-md-6 col-sm-6">
            <div class="cal-datepicker">
              <div class="datepicker-here float-sm-end" data-language="en">
                <v-date-picker v-model="date"  />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  components: {},
  data(){
    return{
      
    }
  }
};
</script>