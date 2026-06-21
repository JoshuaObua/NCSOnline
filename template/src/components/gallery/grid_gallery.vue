<template>
 <Breadcrumbs main="Gallery" title="Grid Gallery"/>
         <div class="container-fluid">
            <div class="row">
              <div class="col-sm-12">
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>IMAGE GALLERY</h5>
                  </div>
                  <div class="gallery my-gallery card-body row" itemscope=" ">
                       <figure
                         v-for="(src, index) in imgs"
                        :key="index"
                         class="col-xl-3 col-md-4 col-sm-6 xl-33 grid-item"
                         @click="() => showImg(index)"
                         >
                        <img :src="require('../../assets/images/big-lightgallry/'+src.image)" alt="Image description" class="img-thumbnail" />
                        <div class="caption">
                        <h4> {{src.title}}</h4>
                        <p>{{src.descr}}</p>
                        </div>
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
         </div>
</template>
<script>
export default {
    data(){
        return{
             lightBoxImages:[],
            visible: false,
          imgs: [
         {
            image:'01.jpg'
          },
          {
            image:'02.jpg'
          },
          {
            image:'03.jpg'
          },
          {
            image:'04.jpg'
          },
          {
            image:'05.jpg'
          },
          {
            image:'011.jpg'
          },
          {
            image:'010.jpg'
          },
          {
            image:'08.jpg'
          },
          {
            image: '013.jpg'
          },
          {
            image: '014.jpg'
          },
          {
            image: '015.jpg'
          },
          {
            image: '016.jpg'
          },
      ]
        }
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
    mounted(){
        this.imgs.forEach(item=>{
            this.lightBoxImages.push(require('@/assets/images/big-lightgallry/'+item.image))
        })
    }
}
</script>