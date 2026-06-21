<template>
     <div class="col-xl-3 col-lg-12 xl-35">
                  <div class="default-according style-1 faq-accordion job-accordion">
                    <div class="row">
                      <div class="col-xl-12">
                        <div class="card">
                          <div class="card-header">
                            <h5 class="p-0">
                              <button class="btn btn-link ps-0" data-bs-toggle="collapse"   data-bs-target="#collapseExample" :aria-expanded="isActive?'true':'false'"  :class="[isActive ? 'active' : '']" v-on:click="toggle">About Me</button>
                            </h5>
                          </div>
                          <div   :class="[isActive ? 'block' : ' show']" v-show="isActive">
                            <div class="card-body post-about">
                                 <ul>
                                <li>
                                  <div class="icon"><i data-feather="briefcase"><vue-feather type="briefcase" class="facebook"></vue-feather></i></div>
                                  <div>
                                    <h5>UX desginer at Pixelstrap</h5>
                                    <p>banglore - 2022</p>
                                  </div>
                                </li>
                                <li>
                                  <div class="icon"><i data-feather="book"><vue-feather type="book" class="facebook"></vue-feather></i></div>
                                  <div>
                                    <h5>studied computer science</h5>
                                    <p>at london univercity - 2015</p>
                                  </div>
                                </li>
                                <li>
                                  <div class="icon"><i data-feather="heart"><vue-feather type="heart" class="facebook"></vue-feather></i></div>
                                  <div>
                                    <h5>relationship status</h5>
                                    <p>single</p>
                                  </div>
                                </li>
                                <li>
                                  <div class="icon"><i data-feather="map-pin"><vue-feather type="map-pin" class="facebook"></vue-feather></i></div>
                                  <div>
                                    <h5>lived in london</h5>
                                    <p>last 5 year</p>
                                  </div>
                                </li>
                                <li>
                                  <div class="icon"><i data-feather="droplet"><vue-feather type="droplet" class="facebook"></vue-feather></i></div>
                                  <div>
                                    <h5>blood group</h5>
                                    <p>O+ positive</p>
                                  </div>
                                </li>
                              </ul>
                              <div class="social-network theme-form"><span class="f-w-600">Social Networks</span>
                                <button class="btn social-btn btn-fb mb-2 text-center"><i class="fa fa-facebook m-r-5"></i>Facebook</button>
                                <button class="btn social-btn btn-twitter mb-2 text-center"><i class="fa fa-twitter m-r-5"></i>Twitter</button>
                                <button class="btn social-btn btn-google text-center"><i class="fa fa-dribbble m-r-5"></i>Dribbble</button>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                      <followerView/>
                      <latesPhotos/>
                    </div>
                  </div>
     </div>
</template>
<script>

import followerView from "./followerView.vue"
import latesPhotos from "./latesPhotos.vue"
export default {
    components:{
        followerView,
        latesPhotos
    },
    data(){
        return{
            isActive: true,
        }
    },
     methods: {
     toggle() {
      this.isActive = !this.isActive
    },
}
}
</script>