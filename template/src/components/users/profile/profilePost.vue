<template>
     <div class="col-xl-9 col-lg-12 xl-65">
                  <div class="row">
                    <div class="col-sm-12">
                      <div class="card">
                        <div class="profile-post">
                          <div class="post-header">
                            <div class="media"><img class="img-thumbnail rounded-circle me-3" src="../../../assets/images/user/7.jpg" alt="Generic placeholder image">
                              <div class="media-body align-self-center"><router-link to="/pages/social-app">
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
                            <p>Dressing is a way of life. My customers are successful working women. I want people to be afraid of the women I dress. Age is something only in your head or a stereotype. Age means nothing when you are passionate about something. There has to be a balance between your mental satisfaction and the financial needs of your company.</p>
                            <ul class="post-comment">
                              <li>
                                <label><a href="#"><i data-feather="heart"><vue-feather type="heart"></vue-feather></i>&nbsp;&nbsp;Like<number
                               class="bold counter"
                               ref="number1"
                               :from="0"
                               :to="50"
                               :duration="5"
                               :delay="0"
                               easing="Power1.easeOut" /></a></label>
                              </li>
                              <li>
                                <label><a href="#"><i data-feather="message-square"><vue-feather type="message-square"></vue-feather></i>&nbsp;&nbsp;Comment<number
                                class="bold counter"
                                ref="number1"
                               :from="0"
                               :to="70"
                               :duration="5"
                               :delay="0"
                                easing="Power1.easeOut" /></a></label>
                              </li>
                              <li>
                                <label><a href="#"><i data-feather="share"><vue-feather type="share"></vue-feather></i>&nbsp;&nbsp;share<number
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
                <profilePost1/>
             </div>
     </div>
</template>
<script>
import profilePost1 from "./profilePost1.vue"
export default {
    components:{
     profilePost1
    },
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