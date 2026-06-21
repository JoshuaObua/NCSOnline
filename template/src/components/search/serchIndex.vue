<template>
  <Breadcrumbs main="Search Results" title="Search Website" />
 
    <div class="container-fluid search-page ">
      <div class="row">
        <div class="col-sm-12">
          <div class="card">
              <div class="card-header pb-0">
              <form class="search-form">
                <div class="form-group mb-0"> 
                    <div class="input-group">
                    
                 <span class="input-group-text"><i class="icon-search"></i></span> 
                 <input
                    class="form-control-plaintext"
                    type="search"
                    placeholder="Search.."
                  />
                  </div>
                </div>
              </form>
                  </div>
                   <div class="card-body">
                       <ul class="nav nav-tabs search-list" id="top-tab" role="tablist">
							<li class="nav-item" v-on:click="changetab('all')">
                                <a class="nav-link" :class="{'show active': showtab == 'all' }" id="top-all" data-bs-toggle="tab" role="tab" aria-controls="all" aria-selected="false"><i class="icon-target"></i>All</a>
							</li>
							<li class="nav-item" v-on:click="changetab('images')">
								<a class="nav-link" :class="{'show active': showtab == 'images' }" id="top-images" data-toggle="tab" role="tab" aria-controls="images" aria-selected="false" ><i class="icon-image"></i>Images </a>
							</li>
							<li class="nav-item" v-on:click="changetab('videos')">
								<a class="nav-link" :class="{'show active': showtab == 'videos' }" id="top-videos" data-toggle="tab" role="tab" aria-controls="videos" aria-selected="false" ><i class="icon-video-clapper"></i>Videos</a>
							</li>
                            <li class="nav-item" v-on:click="changetab('maps')">
								<a class="nav-link" :class="{'show active': showtab == 'maps' }" id="top-maps" data-toggle="tab" role="tab" aria-controls="maps" aria-selected="false" ><i class="icon-map-alt"> </i>Maps</a>
							</li>
							<li class="nav-item" v-on:click="changetab('settings')">
								<a class="nav-link" :class="{'show active': showtab == 'settings' }" id="top-settings" data-toggle="tab" role="tab" aria-controls="settings" aria-selected="false" ><i class="icon-settings"></i>Settings</a>
							</li>
						</ul>
                         <div class="tab-content" id="top-tabContent">
                              <div class="tab-pane fade" :class="{'show active': showtab == 'all' }" id="timeline" role="tabpanel" aria-labelledby="all">
                               <allView/>
                              </div>
                              <div class="tab-pane fade" :class="{'show active': showtab == 'images' }" id="about" role="tabpanel" aria-labelledby="images">
                                  <imagesView/>
                              </div>
                              <div class="tab-pane fade" :class="{'show active': showtab == 'videos' }" id="friends" role="tabpanel" aria-labelledby="videos">
                               <videoView/>
                              </div>
                              <div class="tab-pane fade" :class="{'show active': showtab == 'maps' }" id="friends" role="tabpanel" aria-labelledby="maps">
                               <videoView/>
                              </div>
                              <div class="tab-pane fade" :class="{'show active': showtab == 'settings' }" id="photos" role="tabpanel" aria-labelledby="settings">
                              <videoView/> 
                              </div>
                            </div>      
                   </div>
          </div>
        </div>
      </div>
    </div>
 
</template>
<script>
import allView from "./allView.vue"
import imagesView from "./imagesView.vue"
import videoView from "./videoView.vue"
export default {
    components:{
      allView,
      imagesView,
      videoView
    },
     data() {
      return{
        showtab:'all',
      }
      },
     methods:{
      changetab(tab) {
        this.showtab = tab;
      }
    }
}
</script>