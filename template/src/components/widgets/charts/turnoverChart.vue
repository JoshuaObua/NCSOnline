<template>
  <div class="xl-50 col-xl-5 box-col-6">
    <div class="small-chart-widget chart-widgets-small">
      <div class="card">
        <div class="card-header pb-0">
          <h5>Turnover</h5>
        </div>
        <div class="card-body">
          <div class="chart-container">
            <div class="row">
              <div class="col-12">
                <div id="chart-widget7">
                  <apexchart type="area" height="320" ref="chart" :options="chartOptions" :series="series"></apexchart>
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
var primary = localStorage.getItem('--theme-deafult') || '#24695c';
var secondary = localStorage.getItem('--theme-secondary') || '#ba895d';
import chart from "../../chart/chart.json";
export default {
  data() {
    return {
      series: [
        {

          data: chart.monthDataSeries1.prices,
        },
      ],
      chartOptions: {
        chart: {
          height: 320,
          type: "area",
          zoom: {
            enabled: false,
          },
        },
        dataLabels: {
          enabled: false,
        },
        stroke: {
          curve: "straight",
        },
        fill: {
          colors: [primary],
          type: "gradient",
          gradient: {
            shade: "light",
            type: "vertical",
            shadeIntensity: 0.4,
            inverseColors: false,
            opacityFrom: 0.9,
            opacityTo: 0.8,
            stops: [0, 100],
          },
        },

        colors: [primary],
        labels: chart.monthDataSeries1.dates,
        xaxis: {
          type: "datetime",
        },
        yaxis: {
          opposite: false,
        },
        legend: {
          horizontalAlign: "left",
        },
      },
    };
  },
  mounted() {
    this.series[0].data = chart.monthDataSeries1.prices;
  },
};
</script>