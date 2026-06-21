<template>
     <Breadcrumbs sub="Forms" main="Form Widgets" title="Datepicker"/>
            <div class="container-fluid">
            <div class="card">
              <div class="card-header"><h5>Date Picker</h5></div>
              <div class="card-body">
                <div class="row">
                  <div class="col-md-12">
                    <div class="date-picker">
                     <form class="theme-form">
                       <defaultDate/>
                       <multipleDate/>
                        <monthDate/>
                       <minimumMaximum/>
                       <rangeDate/>
                       <disableDate/>
                        <orientationDate/>
                       <permanentlyVisible/>
                      </form>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            </div>
</template>
<script>
import Datepicker from "vue3-datepicker";
import defaultDate from "./defaultDate.vue"
import multipleDate from "./multipleDate.vue"
import monthDate from "./monthDate.vue"
import minimumMaximum from "./minimumMaximum.vue"
import rangeDate from "./rangeDate.vue"
import disableDate from "./disableDate.vue"
import orientationDate from "./orientationDate.vue"
import permanentlyVisible from "./permanentlyVisible.vue"
export default {
  components: {
    Datepicker,
    defaultDate,
    multipleDate,
    monthDate,
    minimumMaximum,
    rangeDate,
    disableDate,
    orientationDate,
    permanentlyVisible
  },
  
}
</script>