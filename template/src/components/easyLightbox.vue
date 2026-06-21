<template>
  <div class="my-gallery card-body row gallery-with-description" itemscope="">
    <figure 
      v-for="(src, index) in masonryItems"
      :key="index"
      class="col-xl-3 col-sm-6"
      @click="() => showImg(index)" >
      <a>
      <img :src="require('../assets/images/big-lightgallry/' + src.image)"  alt="Image description"  class="img-fluid" />
      <div class="caption">
        <h4>{{ src.title }}</h4>
        <p>{{ src.descr }}</p>
      </div>
      </a>
    </figure >
  </div>
  <vue-easy-lightbox
    :visible="visible"
    :index="index"
    :imgs="lightBoxImages"
    @hide="handleHide"
  >
  </vue-easy-lightbox>
</template>
<script>
import { mapState } from "vuex";
export default {
  data() {
    return {
      lightBoxImages: [],
      visible: false,
      
    };
  },
  computed: {
    ...mapState({
      masonryItems:(state)=>state.masonry.imgs,
         
    }),
  },
  methods: {
    showImg(index) {
      this.index = index;
      this.visible = true;
    },
    handleHide() {
      this.visible = false;
    },
  },
  mounted() {
    this.masonryItems.forEach((item) => {
      this.lightBoxImages.push(
        require("@/assets/images/big-lightgallry/" + item.image)
      );
    });
   
  },
};
</script>