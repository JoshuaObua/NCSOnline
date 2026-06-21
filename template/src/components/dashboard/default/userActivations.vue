<template>
  <div class="col-xl-4 box-col-12 des-xl-100">
    <div class="row">
      <div class="col-xl-12 box-col-6 des-xl-50">
        <div class="card">
          <div class="card-header">
            <div class="header-top d-sm-flex align-items-center">
              <h5>User Activations</h5>
              <div class="center-content">
                <p>Yearly User 24.65k</p>
              </div>
            </div>
          </div>
          <div class="card-body p-0">
            <apexchart type="bar" height="265" :series="series" :options="chartOptions"></apexchart>
          </div>
        </div>
      </div>
      <transactionView />
    </div>
  </div>
</template>
<script>
var primary = localStorage.getItem('--theme-deafult') || '#24695c';
var secondary = localStorage.getItem('--theme-secondary') || '#ba895d';
import transactionView from '../default/transactionView.vue';
export default {
  components: {
    transactionView
  },
  data() {
    return {
      series: [
        {
          name: "Yearly Profit",
          data: [
            {
              x: "Jan",
              y: 1500,
              fillColor: primary,
            },
            {
              x: "Feb",
              y: 3000,
              fillColor: primary,
            },
            {
              x: "Mar",
              y: 1800,
              fillColor: primary,
            },
            {
              x: "Apr",
              y: 3000,
              fillColor: secondary,
            },
            {
              x: "May",
              y: 1800,
              fillColor: primary,
            },
            {
              x: "Jun",
              y: 1500,
              fillColor: primary,
            },
            {
              x: "Jul",
              y: 2500,
              fillColor: primary,
            },
            {
              x: "Sep",
              y: 1500,
              fillColor: secondary,
            },
            {
              x: "Oct",
              y: 2000,
              fillColor: primary,
            },
          ],
        },
      ],
      chartOptions: {
        chart: {
          height: 250,
          width: 60,
          type: "bar",
          toolbar: {
            show: false,
          },
        },
        plotOptions: {
          bar: {
            horizontal: false,
            columnWidth: "30%",
            e̶n̶d̶i̶n̶g̶S̶h̶a̶p̶e̶: 'rounded',
            startingShape: "rounded",
            borderRadius: 4,
            colors: {
              backgroundBarColors: ["#e5edef"],
              backgroundBarOpacity: 1,
              backgroundBarRadius: 9,
            },
          },
        },
        stroke: {
          show: false,
        },
        dataLabels: {
          enabled: false,
        },
        fill: {
          opacity: 1,
        },
        xaxis: {

          axisBorder: {
            show: false,
          },
          labels: {
            show: true,
          },
          axisTicks: {
            show: false,
          },
        },
        yaxis: {
          labels: {
            show: false,
          },
        },
      },
    };
  },
};
</script>