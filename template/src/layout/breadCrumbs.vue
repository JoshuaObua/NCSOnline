<template>
<div class="container-fluid">
    <div class="page-header">
      <div class="row">
        <div class="col-sm-6">
          <h3 class="bread">{{ title }}</h3>
          <ol class="breadcrumb">
            <li class="breadcrumb-item">
              <router-link :to="{ path: '/' }" >Home</router-link>
            </li>
             <li class="breadcrumb-item" v-if="sub">{{ sub }}</li>
            <li class="breadcrumb-item" v-if="main">{{ main }}</li>
            <li class="breadcrumb-item active" v-if="title">{{ title }}</li>
          </ol>
        </div>
         <div class="col-sm-6">
                  <div class="bookmark">
                    <ul>
                      <li><router-link to="/app/bookmark" data-container="body"  data-bs-placement="top" title="Bookmark" > <vue-feather type="heart"></vue-feather></router-link></li>
                      <li><router-link to="/app/chat" data-container="body"  data-placement="top" title="Chat"><vue-feather type="message-square"></vue-feather></router-link> </li>
                      <li><router-link to="/icons/flag" data-container="body"  data-placement="top" title="Icons" ><vue-feather type="command"></vue-feather></router-link> </li>
                      <li><router-link to="/learning/list" data-container="body"  data-placement="top" title="Learning" > <vue-feather type="layers"></vue-feather></router-link></li>
                      <li><a href="javascript:void(0)"><vue-feather class="bookmark-search"  type="star" @click="collapseFilter()"></vue-feather></a>
                        <form class="form-inline search-form">
                          <div class="form-group form-control-search " :class="filtered?'open':''">
                           <searchForm/>
                          </div>
                        </form>
                      </li>
                    </ul>
                  </div>
         </div>
      </div>
    </div>
</div>
</template>
<script>
 import searchForm from "./searchForm.vue"
export default {
  props: {
    title: {
      default: "home",
    },
    main: {
      default: "",
    },
    sub:{
      default:""
    }
  },
  components:{
    searchForm
  },
  data(){
    return{
       filtered: false,
       
    }
  },
   
  methods:{
     collapseFilter() {
        this.filtered = !this.filtered;
      },
  },
};
</script>
