<template>
    <div class="xl-50 col-xl-7 box-col-6">
                <div class="small-chart-widget chart-widgets-small">
                  <div class="card">
                    <div class="card-header pb-0">
                      <h5>Uses</h5>
                    </div>
                    <div class="card-body">
                      <div class="chart-container">
                        <div class="row">
                          <div class="col-12">
                            <div id="chart-widget9">
                                 <usesChart1/>
                            </div>
                          </div>
                        </div>
                      </div>
                      
                    </div>
                  </div>
                </div>
              </div>
</template>
<script>
import usesChart1 from "./usesChart1.vue"
export default {
    components:{
     usesChart1
    },
    data(){
        return{
        }
    }
}
</script>