<template>
  <div class="col-xl-6 box-col-6">
    <div class="card">
      <div class="card-header">
        <div class="
            header-top
            d-sm-flex
            justify-content-between
            align-items-center
          ">
          <h5>Total Selling</h5>
          <div class="center-content">
            <ul class="week-date">
              <li class="font-primary">Today</li>
              <li>Month</li>
            </ul>
          </div>
        </div>
      </div>
      <div class="card-body chart-block p-0">
        <div id="chart-dash-2-line">
          <apexchart type="line" height="470" :series="series" :options="chartOptions"></apexchart>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import chart from "../../chart/chart.json"
var primary = localStorage.getItem('--theme-deafult') || '#24695c';
var secondary = localStorage.getItem('--theme-secondary') || '#ba895d';
export default {
  data() {
    return {
      series: chart.total,
      chartOptions: {
        chart: {
          height: 470,
          type: "line",
          toolbar: {
            show: false,
          },
        },
        stroke: {
          curve: "smooth",
          width: [5, 2],
          dashArray: [0, 8],
        },
        fill: {
          type: "solid",
          opacity: [0.35, 1],
        },
        labels: ["2010", "2011", "2012", "2013"],
        markers: {
          size: 5,
        },
        responsive: [
          {
            breakpoint: 991,
            options: {
              chart: {
                height: 300
              }
            }
          },
          {
            breakpoint: 1500,
            options: {
              chart: {
                height: 325
              }
            }
          }
        ],
        yaxis: [
          {
            labels: {
              formatter: function (value) {
                return value + "k";
              },
            },
          },
        ],
        tooltip: {
          shared: true,
          intersect: false,
          y: {
            formatter: function (y) {
              if (typeof y !== "undefined") {
                return y.toFixed(0) + " points";
              }
              return y;
            },
          },
        },
        colors: [primary, secondary],
        legend: {
          show: false,
        },
      },
    };
  },
};
</script>