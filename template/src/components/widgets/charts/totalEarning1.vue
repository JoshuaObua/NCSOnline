<template>
<apexchart
                  type="radialBar"
                  height="315"
                  ref="chart"
                  :options="chartOptions"
                  :series="series"
                ></apexchart>
</template>
<script>

export default {
  data() {
    return {
      series: [70],
      chartOptions: {
        chart: {
            height: 350,
            type: 'radialBar',
        },
        plotOptions: {
            radialBar: {
                hollow: {
                    margin: 15,
                    size: '70%',
                    image: require('@/assets/images/email-template/success.png'),   
                    imageWidth: 64,
                    imageHeight: 64,
                    imageClipped: false
                },
                dataLabels: {
                    name: {
                        show: false,
                        color: '#fff'
                    },
                    value: {
                        show: true,
                        color: '#333',
                        offsetY: 70,
                        fontSize: '22px'
                    }
                }
            }
        },
        colors:["#999999"],
        stroke: {
            lineCap: 'round'
        },
        labels: ['Volatility'],
    }
  }
  
},
}
</script>