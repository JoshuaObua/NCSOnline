<template>
    <div class="col-sm-12">
                      <div class="card">
                        <div class="profile-post">
                          <div class="post-header">
                            <div class="media"><img class="img-thumbnail rounded-circle me-3" src="../../../assets/images/user/7.jpg" alt="Generic placeholder image">
                              <div class="media-body align-self-center"><router-link to="/pages/social-app">
                                  <h5 class="user-name">Emay Walter</h5></router-link>
                                <h6>2 Hours ago</h6>
                              </div>
                            </div>
                            <div class="post-setting"><i class="fa fa-ellipsis-h"></i></div>
                          </div>
                          <div class="post-body">
                            <div class="img-container">
                              <div class="my-gallery" id="aniimated-thumbnials" itemscope="">
                                <div v-for="(src, index) in imgs" :key="index"  @click="() => showImg(index)" >
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
                              <h6>+20 people react this post</h6>
                            </div>
                            <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.                                     </p>
                            <ul class="post-comment">
                              <li>
                                <label><a href="#"><i data-feather="heart"><vue-feather type="heart"></vue-feather></i>&nbsp;&nbsp;Like<number
                               class="bold counter"
                               ref="number1"
                               :from="0"
                               :to="407"
                               :duration="5"
                               :delay="0"
                               easing="Power1.easeOut" /></a></label>
                              </li>
                              <li>
                                <label><a href="#"><i data-feather="message-square"><vue-feather type="message-square"></vue-feather></i>&nbsp;&nbsp;Comment<number
                                class="bold counter"
                                ref="number1"
                               :from="0"
                               :to="302"
                               :duration="5"
                               :delay="0"
                                easing="Power1.easeOut" /></a></label>
                              </li>
                              <li>
                                <label><a href="#"><i data-feather="share"><vue-feather type="share"></vue-feather></i>&nbsp;&nbsp;share<number
                                class="bold counter"
                                ref="number1"
                                :from="0"
                                :to="180"
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
         require('../../../assets/images/user-profile/post4.jpg'),
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