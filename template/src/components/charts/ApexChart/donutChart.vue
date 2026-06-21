<template>
    <div class="col-sm-12 col-xl-6 box-col-6">
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>Donut Chart</h5>
                  </div>
                  <div class="card-body apex-chart">
                    <div id="donutchart">
                <apexchart
                  width="380"
                  type="donut"
                  :options="chartOptions"
                  :series="series"
                ></apexchart>
                    </div>
                  </div>
                </div>
              </div>
              
</template>
<script>
var primary = localStorage.getItem('--theme-deafult') || '#24695c';
var secondary = localStorage.getItem('--theme-secondary') || '#ba895d';
export default {
    data(){
        return{
          series: [44, 55, 41, 17, 15],
          chartOptions: {
            chart: {
              type: 'donut',
            },
            colors:[primary, secondary, '#222222', '#717171', '#e2c636']
          },

    
    }

        }
    }

</script>