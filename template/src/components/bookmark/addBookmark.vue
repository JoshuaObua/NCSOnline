<template>
    <div class="col-xl-3 box-col-3 xl-30">
                  <div class="email-sidebar"><a class="btn btn-primary email-aside-toggle" href="javascript:void(0)" @click="collapseFilter()">bookmark filter</a>
                    <div class="email-left-aside" :class="filtered?'open':''">
                      <div class="card">
                        <div class="card-body">
                          <div class="email-app-sidebar left-bookmark">
                            <div class="media">
                              <div class="media-size-email"><img class="me-3 rounded-circle" src="../../assets/images/user/user.png" alt=""></div>
                              <div class="media-body">
                                <h6 class="f-w-600">MARK JENCO</h6>
                                <p>Markjecno@gmail.com</p>
                              </div>
                            </div>
                           <ul class="nav main-menu" role="tablist">
                            <newBookmark/>
                            <li class="nav-item"><span class="main-title"> Views</span></li>
                            <li><a id="pills-created-tab" data-toggle="pill" href="javascript:void(0)" role="tab" aria-controls="pills-created" aria-selected="true"><span class="title" v-on:click="say('pills_created')"> Created by me</span></a></li>
                            <li><a class="show" id="pills-favourites-tab" data-toggle="pill" href="javascript:void(0)" role="tab" aria-controls="pills-favourites" aria-selected="false"><span class="title" v-on:click="say('pills_favourites')"> Favourites ({{ favourite.length }})</span></a></li>
                            <li><a class="show" id="pills-shared-tab" data-toggle="pill" href="javascript:void(0)" role="tab" aria-controls="pills-shared" aria-selected="false"><span class="title" v-on:click="say('pills_shared')"> Shared with me</span></a></li>
                            <li><a class="show" id="pills-bookmark-tab" data-toggle="pill" href="javascript:void(0)" role="tab" aria-controls="pills-bookmark" aria-selected="false"><span class="title" v-on:click="say('pills_bookmark')"> My bookmark</span></a></li>
                            <li>
                              <hr/>
                            </li>
                            <li><span class="main-title"> Tags<span class="pull-right"><tagsmodal/></span></span></li>
                            <li><a class="show" id="pills-notification-tab" data-toggle="pill" href="javascript:void(0)" role="tab" aria-controls="pills-notification" aria-selected="false"><span class="title" @click="say('pills_notification')"> notification</span></a></li>
                            <li><a class="show" id="pills-newsletter-tab" data-toggle="pill" href="javascript:void(0)" role="tab" aria-controls="pills-newsletter" aria-selected="false"><span class="title" @click="say('pills_newsletter')"> Newsletter</span></a></li>
                          </ul>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div></div> 
                <div class="col-xl-9 col-md-12 box-col-9 xl-70">
                  <div class="email-right-aside bookmark-tabcontent">
                    <div class="card email-body radius-left">
                      <div class="ps-0">
                        <div class="tab-content">
                          <div class="tab-pane fade" v-bind:class="(activeclass==='pills_created') ? 'active show': ''" id="pills-created" role="tabpanel" aria-labelledby="pills-created-tab">
                            <div class="card mb-0">
                              <div class="card-header d-flex">
                                <h6 class="mb-0">Created by me</h6>
                                <ul>
                                  <li><a class="grid-bookmark-view" href="javascript:void(0);"><vue-feather type="grid" @click="liststyle = !liststyle"></vue-feather></a></li> 
                                  <li><a class="list-layout-view" href="javascript:void(0);"><vue-feather type="list" @click="liststyle = !liststyle"></vue-feather></a></li>  
                                </ul>
                              </div>
                              <div class="card-body pb-0">
                                <div class="details-bookmark text-center" v-bind:class="{ 'list-bookmark' : liststyle}">
                                  <div class="row" id="bookmarkData">
                                    <div class="col-xl-3 col-sm-6 xl-50 box-col-6" v-for="(item,index) in bookmark" :key="index">
                                      <div class="card card-with-border bookmark-card o-hidden">
                                        <div class="details-website">
                                          <img class="img-fluid" :src='getImgUrl(item.image)' alt="">
                                          <div class="favourite-icon favourite_0" :class="(getActive(item.id))? 'favourite' : ''" @click="setFavourite(item)"><a href="javascript:void(0)"><i class="fa fa-star"></i></a></div>
                                          <div class="desciption-data">
                                            <div class="title-bookmark">
                                              <h6 class="title_0">{{ item.title }}</h6>
                                              <p class="weburl_0">{{ item.website_url }}</p>
                                              <div class="hover-block">
                                                <ul>
                                                  <editModel/>
                                                    
                                                 
                                                  <li><a href="#"><vue-feather type="link"></vue-feather></a></li>
                                                  <li><a href="#"><vue-feather type="share-2"></vue-feather></a></li>
                                                  <li><a @click="removeBookmark(index)"><vue-feather type="trash-2"></vue-feather></a></li>
                                                  <li class="pull-right text-end"><a href="#"><vue-feather type="tag"></vue-feather></a></li>
                                                </ul>
                                              </div>
                                              <div class="content-general">
                                                <p class="desc_0"> {{ item.desc }}</p><span class="collection_0">General</span>
                                              </div>
                                            </div>
                                          </div>
                                        </div>
                                      </div>
                                    </div>
                                  </div>
                                </div>
                              </div>
                            </div>
                          </div>
                        <div class="fade tab-pane" v-bind:class="(activeclass==='pills_favourites') ? 'active show': ''" id="pills-favourites" role="tabpanel" aria-labelledby="pills-favourites-tab">
                            <div class="card mb-0">
                              <div class="card-header d-flex">
                                <h6 class="mb-0">Favourites</h6>
                                <ul>
                                  <li><a class="grid-bookmark-view" href="#"><vue-feather type="grid" @click="liststyle = !liststyle"></vue-feather></a></li>
                                  <li><a class="list-layout-view" href="#"><vue-feather type="list" @click="liststyle = !liststyle"></vue-feather></a></li>
                                </ul>
                              </div>
                              <div class="card-body pb-0" v-if="favourite.length">
                                <div class="details-bookmark text-center" v-bind:class="{ 'list-bookmark' : liststyle}">
                                  <div class="row" id="bookmarkData">
                                    <div class="col-xl-3 col-md-4 xl-50" v-for="(item,index) in favourite" :key="index">
                                      <div class="card card-with-border bookmark-card o-hidden">
                                        <div class="details-website">
                                          <img class="img-fluid" :src='getImgUrl(item.image)' alt="">
                                          <div class="favourite-icon favourite_0 favourite" @click="setFavourite(item)"><a href="#"><i class="fa fa-star"></i></a></div>
                                          <div class="desciption-data">
                                            <div class="title-bookmark">
                                              <h6 class="title_0">{{ item.title }}</h6>
                                              <p class="weburl_0">{{ item.website_url }}</p>
                                              <div class="hover-block">
                                                 <ul>
                                                  <li><a href="#" data-toggle="modal" data-target="#edit-bookmark" v-on:click="toogle = !toogle"><vue-feather type="edit-2" ></vue-feather></a></li>
                                                  <li><a href="#"><vue-feather type="link"></vue-feather></a></li>
                                                  <li><a href="#"><vue-feather type="share-2"></vue-feather></a></li>
                                                  <li><a @click="removeFavourite(index)"><vue-feather type="trash-2"></vue-feather></a></li>
                                                  <li class="pull-right text-end"><vue-feather type="tag"></vue-feather></li>
                                                </ul>
                                              </div>
                                              <div class="content-general">
                                                <p class="desc_0"> {{ item.desc }}</p><span class="collection_0">General</span>
                                              </div>
                                            </div>
                                          </div>
                                        </div>
                                      </div>
                                    </div>
                                  </div>
                                </div>
                              </div>
                              <div class="card-body" v-else>
                                <div class="details-bookmark text-center">
                                  <div class="row" id="favouriteData"></div>
                                  <div class="no-favourite"><span>No Bookmarks Found.</span></div>
                                </div>
                              </div>
                            </div>
                        </div>
                        <div class="fade tab-pane" v-bind:class="(activeclass==='pills_shared') ? 'active show': ''" id="pills-shared" role="tabpanel" aria-labelledby="pills-shared-tab">
                            <shared/>
                        </div>
                        <div class="fade tab-pane" v-bind:class="(activeclass==='pills_bookmark') ? 'active show': ''" id="pills-bookmark" role="tabpanel" aria-labelledby="pills-bookmark-tab">
                                <div class="card mb-0">
                              <div class="card-header d-flex">
                                <h6 class="mb-0">My Bookmark</h6>
                                <ul>
                                  <li><a class="grid-bookmark-view" href="javascript:void(0);"><vue-feather type="grid" @click="liststyle = !liststyle"></vue-feather></a></li> 
                                  <li><a class="list-layout-view" href="javascript:void(0);"><vue-feather type="list" @click="liststyle = !liststyle"></vue-feather></a></li>  
                                </ul>
                              </div>
                              <div class="card-body pb-0">
                                <div class="details-bookmark text-center" v-bind:class="{ 'list-bookmark' : liststyle}">
                                  <div class="row" id="bookmarkData">
                                    <div class="col-xl-3 col-md-4 xl-50" v-for="(item,index) in bookmark" :key="index">
                                      <div class="card card-with-border bookmark-card o-hidden">
                                        <div class="details-website">
                                          <img class="img-fluid" :src='getImgUrl(item.image)' alt="">
                                          <div class="favourite-icon favourite_0"><a href="javascript:void(0)"><i class="fa fa-star"></i></a></div>
                                          <div class="desciption-data">
                                            <div class="title-bookmark">
                                              <h6 class="title_0">{{ item.title }}</h6>
                                              <p class="weburl_0">{{ item.website_url }}</p>
                                              <div class="hover-block">
                                                <ul>
                                                  <li><a href="" onclick="editBookmark(0)" data-toggle="modal" data-target="#edit-bookmark"><vue-feather type="edit-2"></vue-feather></a></li>
                                                  <li><a href="#"><vue-feather type="link"></vue-feather></a></li>
                                                  <li><a href="#"><vue-feather type="share-2"></vue-feather></a></li>
                                                  <li><a @click="removeBookmark(index)"><vue-feather type="trash-2"></vue-feather></a></li>
                                                  <li class="pull-right text-end"><a href="#"><vue-feather type="tag"></vue-feather></a></li>
                                                </ul>
                                              </div>
                                              <div class="content-general">
                                                <p class="desc_0"> {{ item.desc }}</p><span class="collection_0">General</span>
                                              </div>
                                            </div>
                                          </div>
                                        </div>
                                      </div>
                                    </div>
                                  </div>
                                </div>
                              </div>
                            </div>
                        </div>
                        <div class="fade tab-pane" v-bind:class="(activeclass==='pills_notification') ? 'active show': ''" id="pills-notification" role="tabpanel" aria-labelledby="pills-notification-tab">
                          <notification/>
                        </div>
                         <div class="fade tab-pane" v-bind:class="(activeclass==='pills_newsletter') ? 'active show': ''" id="pills-newsletter" role="tabpanel" aria-labelledby="pills-newsletter-tab">
                           <newsletter/>
                         </div>
                         
                        </div>
                      </div>
                    </div>
                  </div>
                </div>       
</template>
<script>
import notification from "./notification.vue"
import newsletter from "./newsletter.vue"
import shared from "./shared.vue"
import newBookmark from "./newBookmark.vue"
import tagsmodal from "./tagsmodal.vue"
import editModel from "./editModel.vue"
import { mapState } from 'vuex';
export default {
    components:{
      notification,
      newsletter,
      shared,
      newBookmark,
      tagsmodal,
      editModel
    },
     data(){
      return {
        liststyle: false,
        data : [],
        activeclass : 'pills_created',
        favourite : [],
        filtered: false,
        filter: false,
      };
    },
    computed: {
      ...mapState({
        bookmark: state => state.common.bookmark,
      })
    },
    methods:{
      getImgUrl(path) {
        return require('@/assets/images/'+path);
      },
      say: function (message) {
        this.activeclass = message;
      },
       collapseFilter() {
        this.filtered = !this.filtered;
      },
      collapse(){
           this.filter = !this.filter;
        },
        removeBookmark(index){
          this.bookmark.splice(index,1)
        },
        removeFavourite(index){
          this.favourite.splice(index,1)
        },
      setFavourite(item) {
        let objIndex = this.favourite.findIndex((obj => obj.id === item.id));
        if (objIndex > -1) {
          this.favourite.splice(objIndex, 1);
        } else {
          this.favourite.push(item);
        }
      },
      getActive(itemId) {
        let objIndex = this.favourite.findIndex((obj => obj.id === itemId));
        if (objIndex > -1) {
          return true;
        } else {
          return false;
        }
      }
    }
}
</script>