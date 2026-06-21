<template>
  <div class="col-xl-3 col-sm-6 box-col-3 chart_data_right second">
    <div class="card income-card card-primary">
      <div class="card-body">
        <div class="round-progress knob-block text-center">
          <div class="progress-circle">
            <div style="display: inline; width: 65px; height: 65px">
              <vueknob :width="knobSize" :height="knobSize" color-bg="rgb(36, 105, 92)" :value-min="0" :value-max="100"
                :angle-start="Math.PI" :angle-end="-Math.PI" color-fg="#eeeeee" value="60"
                @value-changed="colorChanged('b', $event)">
                <canvas width="65" height="65"></canvas>
              </vueknob>
            </div>
          </div>
          <h5>$4,55,462</h5>
          <p>Our Annual Income</p>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import vueknob from './vueknob.vue'
export default {
  components: {
    vueknob
  },
  data: function () {
    return {
      knobSize: 250,
      b: 127,
      direction: "N",
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