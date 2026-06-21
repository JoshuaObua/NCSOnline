<template>
    <div class="col-xl-7">
                <div class="card o-hidden">
                  <div class="card-header pb-0">
                    <h5>Marketing Expenses</h5>
                  </div>
                  <div class="bar-chart-widget">
                    <div class="bottom-content card-body">
                      <div class="row">
                        <div class="col-12">
                          <div id="chart-widget4">
                               <apexchart
              type="bar"
              height="375"
              ref="chart"
              :options="chartOptions"
              :series="series"
            ></apexchart>
                          </div>
                        </div>
                      </div>
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
             series: [{
        name: 'Net Profit',
        data: [44, 55, 57, 56, 61, 58, 63, 60, 66]
      }, {
        name: 'Revenue',
        data: [76, 85, 101, 98, 87, 105, 91, 114, 94]
      } ],
      chartOptions: {
          chart: {
        type: 'bar',
        height: 360,
        
      },
      plotOptions: {
        bar: {
          horizontal: false,
          columnWidth: '55%',
            e̶n̶d̶i̶n̶g̶S̶h̶a̶p̶e̶: 'rounded',
          borderRadius: 10,
        },
      },
      dataLabels: {
        enabled: false
      },
      stroke: {
        show: true,
        width: 2,
        colors: ['transparent']
      },
      xaxis: {
        categories: ['Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct'],
      },
      yaxis: {
        title: {
          text: '$ (thousands)'
        }
      },
 
      fill: {
        opacity: 1,
        colors: [primary, secondary],
        type: 'gradient',
        gradient: {
            shade: 'light',
            type: 'vertical',
            shadeIntensity: 0.4,
            inverseColors: false,
            opacityFrom: 0.9,
            opacityTo: 0.8,
            stops: [0, 100]
        }
    },  
    colors: [primary, secondary],
      tooltip: {
        y: {
          formatter: function (val) {
            return "$ " + val + " thousands"
          }
        }
      }
      }
        }
    }
}
</script>