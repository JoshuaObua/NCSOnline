<template>
  <Breadcrumbs main="Bonus UI" title="Range Slider" />
  <div class="container-fluid">
    <div class="row">
      <div class="col-sm-12">
        <div class="card">
          <div class="card-header pb-0">
            <h5>Ion Range Slider</h5>
          </div>
          <div class="card-body range-space">
            <form class="theme-form form-label-align-right range-slider">
              <div class="form-group mt-3 me-1 row">
                <label class="col-md-2 col-form-label sm-left-text" for="u-range-01">Default</label>
                <div class="col-md-10 range-item">
                  <VueSlider v-model="one.value" :data="one.data" :marks="true" :tooltip="'always'"
                    :tooltip-placement="'top'"></VueSlider>
                </div>
              </div>
              <div class="form-group mt-5 me-1 row">
                <label class="col-md-2 control-label sm-left-text" for="u-range-02">Min-Max Value</label>
                <div class="col-md-10 range-item">
                  <VueSlider v-model="two.value" :data="two.data" :marks="true" :tooltip="'always'"
                    :tooltip-placement="'top'"></VueSlider>
                </div>
              </div>
              <div class="form-group mt-5 me-1 row">
                <label class="col-md-2 control-label sm-left-text" for="u-range-03">Prefix</label>
                <div class="col-md-10 range-item">
                  <VueSlider v-model="three.value" :data="three.data1" :marks="true" :tooltip="'always'"
                    :tooltip-placement="'top'"></VueSlider>
                </div>
              </div>
              <div class="form-group mt-5 me-1 row">
                <label class="col-md-2 control-label sm-left-text" for="u-range-04">Nagative value</label>
                <div class="col-md-10 range-item">
                  <VueSlider v-model="four.value" :data="four.data2" :marks="true" :tooltip="'always'"
                    :tooltip-placement="'top'"></VueSlider>
                </div>
              </div>
              <div class="form-group mt-5 me-1 row">
                <label class="col-md-2 control-label sm-left-text" for="u-range-05">Steps</label>
                <div class="col-md-10 range-item">
                  <VueSlider v-model="foour.value" :data="foour.data3" :marks="true" :tooltip="'always'"
                    :tooltip-placement="'top'"></VueSlider>
                </div>
              </div>
              <div class="form-group mt-5 me-1 row">
                <label class="col-md-2  control-label sm-left-text" for="u-range-06">Custom Values</label>
                <div class="col-md-10 range-item">
                  <VueSlider v-model="five.value" :data="five.data" :marks="true" :tooltip="'always'"
                    :tooltip-placement="'top'"></VueSlider>
                </div>
              </div>
              <div class="form-group mt-5 me-1 row">
                <label class="col-md-2 control-label sm-left-text" for="u-range-07">Prettify Numbers</label>
                <div class="col-md-10 range-item">
                  <VueSlider v-model="six.value" :data="six.data" :marks="true" :tooltip="'always'"
                    :tooltip-placement="'top'"></VueSlider>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import VueSlider from 'vue-slider-component'
import 'vue-slider-component/theme/antd.css'
export default {
  components: {
    VueSlider
  },
  data() {
    return {
      one: {
        value: 10,
        data: [10, 100, 200, 300, 400]
      },
      two: {
        value: 550,
        data: [100, 200, 300, 400, 500, 550, 700, 800, 1000]
      },
      three: {
        value: [200, 800],
        data1: [10, 100, 200, 250, 500, 750, 800, 1000],

      },
      four: {
        value: [-500, 500],
        data2: [-1000, -720, -500, -250, 0, 250, 500, 720, 1000],
      },
      foour: {
        value: [-250, 250],
        data3: [-1000, -720, -500, -250, 0, 250, 500, 720, 1000],
      },
      five: {
        value: 'April',
        data: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "July", "Aug", "Sept", "Oct", "Nov", "Dec"]
      },
      six: {
        value: 2000,
        data: [1000, 2000, 3250, 5500, 7750, 10000]
      }
    }
  }
}
</script>