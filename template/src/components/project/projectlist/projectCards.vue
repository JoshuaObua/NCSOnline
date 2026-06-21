<template>
  <div class="col-md-12 project-list">
    <div class="card">
      <div class="row">
        <div class="col-md-6 p-0">
          <ul class="nav nav-tabs border-tab" id="top-tab" role="tablist">
             <li class="nav-item">
               <a class="nav-link active" id="top-home-tab" data-bs-toggle="tab" href="#top-home" role="tab" aria-controls="top-home" aria-selected="true"><vue-feather type="target"></vue-feather>All</a>
               </li>
             <li class="nav-item">
               <a class="nav-link" id="profile-top-tab" data-bs-toggle="tab" href="#top-profile" role="tab" aria-controls="top-profile" aria-selected="false"><vue-feather type="info"></vue-feather>Doing</a>
               </li>
            <li class="nav-item">
              <a class="nav-link" id="contact-top-tab" data-bs-toggle="tab" href="#top-contact" role="tab" aria-controls="top-contact" aria-selected="false"><vue-feather type="check-circle"></vue-feather>Done</a>
              </li>
          </ul>
       </div>
        <div class="col-md-6 p-0">
          <div class="form-group mb-0 me-0"></div>
          <router-link to="createproject"><a class="btn btn-primary" href="">
              <vue-feather class="me-1" type="plus-square"></vue-feather>Create New Project</a></router-link>
        </div>
      </div> 
    </div> 
     </div> 
      <div class="col-sm-12 ">
         <div class="card">
             <div class="card-body">
                <div class="tab-content" id="top-tabContent">
                   <div class="tab-pane fade show active" id="top-home" role="tabpanel" aria-labelledby="top-home-tab">
                    <projectList/>
                  </div>  
                  <div class="tab-pane fade" id="top-profile" role="tabpanel" aria-labelledby="profile-top-tab">
                    <projectDoing/>
                  </div>    
                   <div class="tab-pane fade" id="top-contact" role="tabpanel" aria-labelledby="contact-top-tab">
                    <projectDone/>
                  </div>           
               </div>
             </div>
         </div>
     </div>
</template>
<script>
import projectList from "./projectList.vue";
import projectDoing from './projectDoing.vue';
import projectDone from "./projectDone.vue"
export default {
  components: {
    projectList,
    projectDoing,
    projectDone
  },
  data() {
    return {
      showtab:'all',
    };
  },
  methods: {
   changetab(tab) {
        this.showtab = tab;
      }
  },
};
</script>