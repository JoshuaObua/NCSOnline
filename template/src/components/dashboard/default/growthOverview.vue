<template>
  <div class="col-xl-8 box-col-12 des-xl-100">
    <div class="row">
      <div class="col-xl-6 col-50 box-col-6 des-xl-50">
        <div class="card">
          <div class="card-header">
            <div class="header-top d-sm-flex align-items-center">
              <h5>Growth Overview</h5>
              <div class="center-content">
                <p class="d-flex align-items-center">
                  <i class="toprightarrow-primary fa fa-arrow-up me-2 d-grid"></i>80% Growth </p>
              </div>
            </div>
          </div>
          <div class="card-body p-0 me-3 ms-4 growth-chart ">
            <div id="chart-dashbord">
           <apexchart type="radialBar" height="380" :series="series" :options="chartOptions" ></apexchart>
           </div>
          </div>
        </div>
      </div>
      <latestActivity />
      <recentOrders />
    </div>
  </div>
</template>
<script>
import latestActivity from "../default/latestActivity.vue";
import recentOrders from "../default/recentOrders.vue";
var primary = localStorage.getItem('--theme-deafult') || '#24695c';
var secondary = localStorage.getItem('--theme-secondary') || '#ba895d';
export default {
  components: {
    latestActivity,
    recentOrders,
  },
  data() {
    return {
      series: [76, 67, 61, 90],
          chartOptions: {
            chart: {
              height: 390,
              type: 'radialBar',
            },
            plotOptions: {
              radialBar: {
                offsetY: 0,
                startAngle: 0,
                endAngle: 270,
                 borderradius:20,
                hollow: {
                  margin: 5,
                  size: '30%',
                  background: 'transparent',
                  image: undefined,
                },
                dataLabels: {
                  name: {
                    show: false,
                  },
                  value: {
                    show: false,
                  }
                }
              }
            },
            colors: [primary, secondary, primary, secondary],
            labels: ['Total Order', 'Total product ', 'Quantity', 'Page views'],
            legend: {
              show: true,
              floating: true,
              fontSize: '14px',
              position: 'right',
              offsetX: 200,
              fontweight: 400,
              offsetY: 15,
              fontfamily: 'Roboto',
              labels: {
                useSeriesColors: true,
               
              },
              markers: {
                size: 0
              },
              formatter: function(seriesName, opts) {
                return seriesName + ":  " + opts.w.globals.series[opts.seriesIndex]
              },
              itemMargin: {
                vertical: 3,
                 horizontal: 2
              }
            },
          },     
    };
  },
};
</script>