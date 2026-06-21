<template>
  <div class="col-xl-8 col-md- des-xl-100 box-col-12">
    <div class="row">
      <div class="col-xl-3 col-sm-6 box-col-3 chart_data_right">
        <div class="card income-card card-secondary">
          <div class="card-body align-items-center">
            <div class="round-progress knob-block text-center">
              <div class="progress-circle">
                <div style="display: inline; width: 65px; height: 65px">
                  <vueknob :width="knobSize" :height="knobSize" color-bg="rgb(186, 137, 93)" :value-min="0"
                    :value-max="100" :angle-start="Math.PI" :angle-end="-Math.PI" color-fg="#eeeeee" value="60"
                    @value-changed="colorChanged('b', $event)">
                  </vueknob>
                </div>
              </div>
              <h5>$9,84,235</h5>
              <p>Our Annual Income</p>
            </div>
          </div>
        </div>
      </div>
      <ecormmerceIncome2 />
      <sellingProduct />
      <totalSelling />
      <ourTarget />
    </div>
  </div>
</template>
<script>
import ourTarget from '../e-commerce/ourTarget.vue'
import ecormmerceIncome2 from '../e-commerce/ecormmerceIncome2.vue'
import sellingProduct from '../e-commerce/sellingProduct.vue'
import totalSelling from '../e-commerce/totalSelling.vue'
import vueknob from './vueknob.vue'

export default {
  components: {
    ecormmerceIncome2,
    sellingProduct,
    totalSelling,
    ourTarget,
    vueknob
  },
  data() {
    return {
      knobSize: 250,
      b: 127,
    };
  },
  methods: {
    colorChanged(color, e) {
      this[color] = e;
    },

    directionChanged(e) {
      this.direction = e;
    },

    to2hex(n) {
      let s = n.toString(16);
      if (s.length < 2) s = "0" + s;
      return s;
    },


    directionToValue(direction) {
      const idx = directions.findIndex((el) => el === direction.toUpperCase());
      return idx === -1 ? 0 : idx;
    },


    valueToDirection(value) {
      return directions[Math.trunc(value % 8)];
    },
  },
};
</script>