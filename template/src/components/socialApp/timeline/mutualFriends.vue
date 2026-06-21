<template>
    <div class="col-md-12 col-sm-6">
                            <div class="card">
                              <div class="card-header">
                                <h5 class="p-0">
                                  <button class="btn btn-link ps-0" data-bs-toggle="collapse" data-bs-target="#collapseicon1" :aria-expanded="isActive?'true':'false'"  aria-controls="collapseicon1"  :class="[isActive ? 'active' : '', 'collapsible']" v-on:click="toggle">Mutual Friends</button>
                                </h5>
                              </div>
                              <div  :class="[isActive ? 'block' : ' show']"  v-show="isActive">
                                <div class="card-body social-status filter-cards-view">
                                  <form>
                                    <div class="form-group m-0">
                                      <input class="form-control-social" type="search" placeholder="Search Contacts..">
                                    </div>
                                  </form>
                                  <div class="media" v-for="friend in friends" :key="friend"><img class="img-50 rounded-circle m-r-15" :src="getImgUrl(friend.image)" alt="">
                                    <div class=" social-status" :class="friend.StatusClass"></div>
                                    <div class="media-body"><a href=""><router-link to="/users/profile"><span class="f-w-600 d-block">{{friend.name}}</span>
                                    </router-link></a><span class="d-block">{{friend.email}}</span></div>
                                  </div>
                                  
                                </div>
                              </div>
                            </div>
                          </div>
</template>
<script>
export default {
    data(){
        return{
             isActive: true,
            friends:[
              {image:"2.png",StatusClass:"social-online", name:"Bucky Barnes", email:"winter@gmail.com"},
              {image:"10.jpg",StatusClass:"social-busy", name:"Sarah Loren", email:"barnes@gmail.com"},
              {image:"6.jpg",StatusClass:"social-offline", name:"Jason Borne", email:"jasonb@gmail.com"},
              {image:"8.jpg",StatusClass:"social-offline", name:"Comeren Diaz", email:"comere@gmail.com"},
              {image:"14.png",StatusClass:"social-online", name:"Andew Jon", email:"andrewj@gmail.com"},
              {image:"4.jpg",StatusClass:"social-busy", name:"Johny Waston", email:"johny@gmail.com"},
              {image:"3.png",StatusClass:"social-offline", name:"Johny William", email:"johnyw@gmail.com"},
              {image:"2.png",StatusClass:"social-online", name:"Bucky Barnes", email:"winter@gmail.com"},
              {image:"10.jpg",StatusClass:"social-busy", name:"Sarah Loren", email:"barnes@gmail.com"},
              {image:"6.jpg",StatusClass:"social-offline", name:"Jason Borne", email:"jasonb@gmail.com"},
              {image:"8.jpg",StatusClass:"social-offline", name:"Comeren Diaz", email:"comere@gmail.com"},
              {image:"14.png",StatusClass:"social-online", name:"Andew Jon", email:"andrewj@gmail.com"},
              {image:"4.jpg",StatusClass:"social-busy", name:"Johny Waston", email:"johny@gmail.com"},
              {image:"3.png",StatusClass:"social-offline", name:"Johny William", email:"johnyw@gmail.com"},       
            ]
        }
    },
    methods:{
       getImgUrl(path) {
      return require("../../../assets/images/user/" + path);
    },
      toggle() {
      this.isActive = !this.isActive
    },
    }
}
</script>