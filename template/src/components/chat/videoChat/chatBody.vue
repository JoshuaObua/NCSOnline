<template>
    <div class="col call-chat-body">
          <div class="card">
            <div class="card-body p-0">
              <div class="row chat-box">
                  <div class="col chat-right-aside">
                  <div class="chat">
                    <div class="media chat-header clearfix">
                      <img class="rounded-circle" v-if="currentchat.thumb" :src="getImgUrl(currentchat.thumb)" alt=""/>
                      <div class="media-body">
                        <div class="about">
                          <div class="name">
                            {{ currentchat.name}}
                          </div>
                          <div class="status digits">
                            {{ currentchat.lastSeenDate }}
                          </div>
                      </div>
                      </div>
                      <ul class="list-inline float-start float-sm-right chat-menu-icons">
                        <li class="list-inline-item">
                          <a href="#"><i class="icon-search"></i></a>
                        </li>
                        <li class="list-inline-item">
                          <a href="#"><i class="icon-clip"></i></a>
                        </li>
                        <li class="list-inline-item">
                          <a href="#"><i class="icon-headphone-alt"></i></a>
                        </li>
                        <li class="list-inline-item">
                          <a href="#"><i class="icon-video-camera"></i></a>
                        </li>
                        <li class="list-inline-item toogle-bar" data-bs-toggle="collapse"  data-bs-target="#collapseExample" :class="[isActive ? 'active' : '', 'collapsible']" v-on:click="toggle">
                          <a href="javascript:void(0)"><i class="icon-menu"></i></a>
                        </li>
                      </ul>
                    </div>
                   <div class="chat-history ">
                            <div class="row" v-for="(chat, index) in currentChat.chat.video" :key="index" v-bind:class="{ clearfix: chat.sender == 0 }">
                              <div class="col text-center pe-0 call-content "  >
                                <div>
                                  <div class="total-time">
                                    <h2 class="digits">{{chat.time}}</h2>
                                  </div>
                                  <div class="call-icons">
                                    <ul class="list-inline">
                                      <li class="list-inline-item"><a href="javascript:void(0)"><i class="icon-video-camera"></i></a></li>
                                      <li class="list-inline-item"><a href="javascript:void(0)"><i class="icon-volume"></i></a></li>
                                      <li class="list-inline-item"><a href="javascript:void(0)"><i class="icon-user"></i></a></li>
                                    </ul>
                                  </div>
                                  <button class="btn btn-danger btn-block btn-lg">END CALL</button>
                                  <div class="receiver-img"><img src="../../../assets/images/other-images/receiver-img.jpg" alt=""></div>
                                </div>
                              </div>
                              <div class="col-sm-6 caller-img-sec">
                                <div class="caller-img"></div>
                              </div>
                            </div>
                          </div>
                       </div>
                    </div>
                <div class="col pl-0 chat-menu custom-scrollbar" :class="[isActive ? 'block' : 'none', 'content']" v-show="isActive" id="collapseExample">
                <chatMenu/>
                </div>
              </div>
            </div>
          </div>
        </div>
</template>
<script>
 import { mapState } from 'vuex';
 import chatMenu from "../chatList/chatMenu.vue"
export default {
    components:{chatMenu},
    data(){
        return{
         currentchat: [],
        chatmenutoogle: false,
        isActive: false,
        }
    },
    computed: {
      ...mapState({
        currentChat() { 
          return (this.currentchat = this.$store.getters['chat/currentChat']);
        },
      }),
    },
    methods: {
      getImgUrl(path) {
        return require('../../../assets/images/' + path);
      },
      toggle() {
      this.isActive = !this.isActive
    },
      }
}
</script>