<template>
    <div class="col-sm-12 box-col-12">
                <div class="donut-chart-widget">
                  <div class="card">
                    <div class="card-header pb-0">
                      <h5>Browser Uses</h5>
                    </div>
                    <div class="card-body">
                      <div id="chart-widget13">
                          <browserUses1/>
                      </div>
                      
                    </div>
                  </div>
                </div>
    </div>
</template>
<script>
import browserUses1 from './browserUses1.vue'
export default {
    components:{
        browserUses1
    }
}
</script>