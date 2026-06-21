<template>
    <apexchart
              type="area"
              height="180"
              ref="chart"
              :options="chartOptions"
              :series="series"
            ></apexchart>
</template>
<script>
var primary = localStorage.getItem('--theme-deafult') || '#24695c';
var secondary = localStorage.getItem('--theme-secondary') || '#ba895d';
export default {
    data(){
        return{
             series: [{
            data: [24, 55, 21, 67, 22, 43, 21]
        }],
        chartOptions: {
             chart: {
            toolbar: {
                show: false
            },
            height: 170,
            type: 'area'
        },
        dataLabels: {
            enabled: false
        },
        stroke: {
            curve: 'smooth'
        },
         yaxis: {
          show: false,
        },
        xaxis: {
            show: false,
            type: 'datetime',
            categories: ["2018-09-19T00:00:00", "2018-09-19T01:30:00", "2018-09-19T02:30:00", "2018-09-19T03:30:00", "2018-09-19T04:30:00", "2018-09-19T05:30:00", "2018-09-19T06:30:00"],
            labels: {
                show: false,
            },
            axisBorder: {
                show: false,
            },
        },
        grid: {
            show: false,
            padding: {
                left: 0,
                right: 0,
                bottom: -40
            }
        },
        fill: {
            colors: [primary],
            type: 'gradient',
            gradient: {
                shade: 'light',
                type: 'vertical',
                shadeIntensity: 0.4,
                inverseColors: false,
                opacityFrom: 0.8,
                opacityTo: 0.2,
                stops: [0, 100]
            },
        },
        colors: [primary],
        
        tooltip: {
            x: {
                format: 'dd/MM/yy HH:mm',
                enabled: false,
            }
        }
        }
        }
    }
}
</script>