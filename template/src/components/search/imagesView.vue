<template>
    <div class="info-block m-t-30">
    <p class="pb-4 col-sm-12 digits ">About 12,120 results (0.50 seconds)</p>
         <div class="my-gallery card-body row gallery-with-description p-0" itemscope="">
          <figure
            v-for="(src, index) in imgs"
            :key="index"
            class="col-xl-3 col-sm-6"
            @click="() => showImg(index)"
          >
          <a>
            <img :src="require('../../assets/images/big-lightgallry/'+src.image)" alt="Image description" class="img-fluid" />
            <div class="caption">
           <h4> {{src.title}}</h4>
            <p>{{src.descr}}</p>
            </div>
          </a>
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
</template>
<script>
export default {
    data(){
        return{
            lightBoxImages:[],
            visible: false,
          imgs: [
         {
            title:'Portfolio Title',
            descr:'is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry\'s standard dummy.',
            image:'01.jpg'
          },
          {
            title:'Portfolio Title',
            descr:'is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry\'s standard dummy.',
            image:'02.jpg'
          },
          {
            title:'Portfolio Title',
            descr:'is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry\'s standard dummy.',
            image:'03.jpg'
          },
          {
            title:'Portfolio Title',
            descr:'is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry\'s standard dummy.',
            image:'04.jpg'
          },
          {
            title:'Portfolio Title',
            descr:'is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry\'s standard dummy.',
            image:'05.jpg'
          },
          {
            title:'Portfolio Title',
            descr:'is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry\'s standard dummy.',
            image:'011.jpg'
          },
          {
            title:'Portfolio Title',
            descr:'is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry\'s standard dummy.',
            image:'010.jpg'
          },
          {
            title:'Portfolio Title',
            descr:'is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry\'s standard dummy.',
            image:'08.jpg'
          }
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
};
</script>
