<template>
                    <div class="col-sm-12">
                      <div class="card">
                        <div class="profile-post">
                          <div class="post-header">
                            <div class="media"><img class="img-thumbnail rounded-circle me-3" src="../../../assets/images/user/7.jpg" alt="Generic placeholder image">
                              <div class="media-body align-self-center" id="first-post-tour"><router-link to="/pages/social-app">
                                  <h5 class="user-name">Emay Walter</h5></router-link>
                                <h6>22 Hours ago</h6>
                              </div>
                            </div>
                            <div class="post-setting"><i class="fa fa-ellipsis-h"></i></div>
                          </div>
                          <div class="post-body">
                            <div class="img-container">
                              <div class="my-gallery" id="aniimated-thumbnials" itemscope="">
                                <div v-for="(src, index) in imgs" :key="index" class="" @click="() => showImg(index)" >
                                <img :src="src" class="img-fluid rounded">
                                </div>
                                <vue-easy-lightbox :index="index" :visible="visible" :imgs="imgs" @hide="handleHide" ></vue-easy-lightbox>
                              </div>
                            </div>
                            <div class="post-react mt-3">
                              <ul>
                                <li><img class="rounded-circle" src="../../../assets/images/user/3.jpg" alt=""></li>
                                <li><img class="rounded-circle" src="../../../assets/images/user/5.jpg" alt=""></li>
                                <li><img class="rounded-circle" src="../../../assets/images/user/1.jpg" alt=""></li>
                              </ul>
                              <h6>+5 people react this post</h6>
                            </div>
                            <p>Comfort is very important to me. I think people live better in big houses and in big clothes. Design and style should work toward making you look good and feel good without a lot of effort so you can get on with the things that matter. My shows are about the complete woman who swallows it all. Its a question of survival. Those fashion designers are just crazy; but arent we all? You can only go forward by making mistakes.</p>
                            <ul class="post-comment">
                              <li>
                                <label><a href="#"><i data-feather="heart"></i>&nbsp;&nbsp;Like<number class="bold counter" ref="number1" :from="0" :to="50" :duration="5" :delay="0" easing="Power1.easeOut" /></a></label>
                              </li>
                              <li>
                                <label><a href="#"><i data-feather="message-square"></i>&nbsp;&nbsp;Comment<number
                                class="bold counter"
                                ref="number1"
                               :from="0"
                               :to="70"
                               :duration="5"
                               :delay="0"
                                easing="Power1.easeOut" /></a></label>
                              </li>
                              <li>
                                <label><a href="#"><i data-feather="share"></i>&nbsp;&nbsp;share<number
                                class="bold counter"
                                ref="number1"
                                :from="0"
                                :to="20"
                                :duration="5"
                                :delay="0"
                                easing="Power1.easeOut" /></a></label>
                              </li>
                            </ul>
                          </div>
                        </div>
                      </div>
                    </div>   
</template>
<script>
export default {
   data(){
       return{
            visible: false,
      imgs: [
         require('../../../assets/images/user-profile/post1.jpg'),
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
      }
    }
}
</script>