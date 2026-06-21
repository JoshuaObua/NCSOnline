<template>
  <div>
    <apexchart type="area" height="425" ref="chart" :options="chartOptions" :series="series"></apexchart>
  </div>
</template>
<script>
var primary = localStorage.getItem('--theme-deafult') || '#24695c';

import chart from "../../chart/chart.json"

export default {
  data() {
    return {
      series: [
        {
          data: chart.data
        }
      ],
      chartOptions: {
        chart: {
          id: "area-datetime",
          type: "area",
          height: 425,
          zoom: {
            autoScaleYaxis: true,
          },
          toolbar: {
            show: false,
          },
        },
        annotations: {
          yaxis: [
            {
              y: 50,
              borderColor: primary,
              label: {
                show: true,
                text: "Support",
                style: {
                  color: "#fff",
                  background: primary,
                },
              },
            },
          ],
          xaxis: [
            {
              x: new Date("15 Nov 2012").getTime(),
              borderColor: primary,
              yAxisIndex: 50,
              label: {
                show: false,
                text: "$859,432",
                style: {
                  color: "#fff",
                  background: primary,
                },
              },
            },
          ],
        },
        dataLabels: {
          enabled: false,
        },
        markers: {
          size: 0,
          style: "hollow",
        },
        xaxis: {
          type: "datetime",
          min: new Date("01 Apr 2012").getTime(),
          tickAmount: 6,
          axisTicks: {
            show: false,
          },
          axisBorder: {
            show: false,
          },
        },
        tooltip: {
          x: {
            format: "dd MMM yyyy",
          },
        },
        fill: {
          type: "gradient",
          gradient: {
            shadeIntensity: 1,
            opacityFrom: 0.7,
            opacityTo: 0.9,
            stops: [0, 100],
          },
        },
        colors: [primary],
      },


    };
  },

};
</script>