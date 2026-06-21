<template>
    <div class="email-sidebar"><a class="btn btn-primary email-aside-toggle"  href="javascript:void(0)" @click="collapseFilter()">To Do filter</a>
   <div class="email-left-aside" :class="filtered?'open':''">
       <div class="card">
           <div class="card-body">
               <div class="email-app-sidebar left-bookmark custom-scrollbar">
                   <div class="media align-items-center">
                       <div class="media-size-email">
                           <img class="me-3 rounded-circle" src="../../assets/images/user/user.png" alt="">
                       </div>
                       <div class="media-body">
                           <h6 class="f-w-600">MARK JENCO</h6>
                           <p>Markjecno@gmail.com</p>
                       </div>
                   </div>
                   <ul class="nav main-menu" role="tablist">
                       <li class="nav-item">
                         <newTodo/>
                       </li>
                       <li class="nav-item">
                           <a href="Javascript:void(0)">
                               <span class="iconbg badge-light-primary"> <vue-feather type="file-plus"></vue-feather></span>
                              <span class="title  ms-2">
                                  All Task
                              </span>
                           </a>
                       </li>
                       <li class="nav-item">
                           <a href="Javascript:void(0)">
                                <span class="iconbg badge-light-success">
                                  <vue-feather type="check-circle"></vue-feather>
                               </span>
                               <span class="title ms-2">
                                   Important 
                                </span>
                               <span class="badge badge-success"> 2 </span>
                           </a>
                       </li>
                       <li class="nav-item">
                           <a href="Javascript:void(0)">
                                <span class="iconbg badge-light-danger">
                                   <vue-feather type="cast"></vue-feather>
                               </span>
                               <span class="title ms-2"> 
                                 Pending 
                                </span>
                               <span class="badge badge-danger"> 3 </span>
                           </a>
                       </li>
                       <li class="nav-item">
                           <a href="Javascript:void(0)">
                                <span class="iconbg badge-light-info">
                                  <vue-feather type="activity"></vue-feather>
                               </span>
                               <span class="title ms-2">
                                  In Progress
                                </span>
                           </a>
                       </li>
                       <li class="nav-item">
                           <a href="Javascript:void(0)">
                                <span class="iconbg badge-light-danger">
                                  <vue-feather type="trash"></vue-feather>
                               </span>
                              <span class="title ms-2">
                            Trash
                                </span>
                           </a>
                       </li>
                   </ul>
               </div>
           </div>
       </div>
   </div>
    </div>
</template>
<script>
import newTodo from "./newTodo.vue"
export default {
    components:{
newTodo
    },
      data(){
      return {
        filtered: false,
      };
    },
    methods:{
         collapseFilter() {
        this.filtered = !this.filtered;
      },
    }
}
</script>