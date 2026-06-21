<template>
     <div class="col-sm-12 ">
          <div class="card">
            <div class="card-header pb-0">
              <h5>MASONRY GALLERY</h5>
            </div>
            <div class="card-body photoswipe-pb-responsive" >
              <div v-masonry  id="aniimated-thumbnials" itemscope="" class="masonry-container my-gallery row grid gallery">
                <figure
                   v-masonry-tile 
                   class="col-md-3 col-sm-6 grid-item"  itemprop="associatedMedia" itemscope=""
                  :key="index"
                  v-for="(item, index) in imagearray"
                  @click="() => showImg(index)"
                  >
                  <a><img :src="require('../../../assets/images/masonry/'+item.image)" class="img-fluid" /></a>
                </figure>
              </div>
               <vue-easy-lightbox
               :index="index"
                  :visible="visible"
                  :imgs="lightBoxImages" 
                  @hide="handleHide"
                    >
                  </vue-easy-lightbox>
            </div>
          </div>
        </div>
</template>
<script>
export default {
    data() {
    return {   
      lightBoxImages:[],
      visible: false,
      imagearray: [
          {  
           image: '1.jpg'
          },
          {       
            image:'2.jpg'
          },
          {        
           image: '3.jpg'
          },
          {    
            image:'4.jpg'
          },
          {  
           image: '5.jpg'
          },
          {       
           image: '6.jpg'
          },
          {       
            image:'8.jpg'
          },
          {        
            image:'9.jpg'
          },
          {      
            image:'10.jpg'
          },
          {      
            image:'11.jpg'
          },
          {         
           image: '12.jpg'
          },
          {          
           image: '14.jpg'
          },
          {          
           image: '15.jpg'
          },
          {
           image: '13.jpg'
          }
        ]
    };
  },
  mounted() {
    if (typeof this.$redrawVueMasonry === 'function') {
        this.$redrawVueMasonry();
      }
       this.imagearray.forEach(item=>{
            this.lightBoxImages.push(require('../../../assets/images/masonry/'+item.image))
        })
  },
    methods: {
      showImg (index) {
        this.index = index
        this.visible = true
      },
      handleHide () {
        this.visible = false
      },
    },
}
</script>