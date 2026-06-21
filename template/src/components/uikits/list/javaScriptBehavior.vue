<template>
     <div class="col-sm-12 col-xl-12 ">
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>JavaScript behavior</h5>
                  </div>
                  <div class="card-body">
                    <div class="row">
                      <div class="col-sm-4">
                          <div class="list-group" id="list-tab" role="tablist">
                              <a class="list-group-item list-group-item-action " id="list-home-list" v-on:click="show_panel('home')" :class="[ tab_panel == 'home' ? 'active' : '']">Home</a>
                              <a class="list-group-item list-group-item-action " id="list-home-list" v-on:click="show_panel('profile')" :class="[ tab_panel == 'profile' ? 'show active' : '']">profile</a>
                              <a class="list-group-item list-group-item-action " id="list-home-list" v-on:click="show_panel('messages')" :class="[ tab_panel == 'messages' ? 'show active' : '']">messages</a>
                              <a class="list-group-item list-group-item-action " id="list-home-list" v-on:click="show_panel('settings')" :class="[ tab_panel == 'settings' ? 'show active' : '']">settings</a>
                          </div>
                          </div>
                          <div class="col-sm-8">
                        <div class="tab-content" id="nav-tabContent">
                            <div class="tab-pane fade" id="list-home" role="tabpanel" :class="[ tab_panel == 'home' ? 'show active' : '']">
                                      1. Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
                                  </div>
                                  <div class="tab-pane fade" id="list-home" role="tabpanel" :class="[ tab_panel == 'profile' ? 'show active' : '']">
                                      2. Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled.
                                  </div>
                                  <div class="tab-pane fade" id="list-home" role="tabpanel" :class="[ tab_panel == 'messages' ? 'show active' : '']">
                                      3. Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
                                  </div>
                                  <div class="tab-pane fade" id="list-home" role="tabpanel" :class="[ tab_panel == 'settings' ? 'show active' : '']">
                                      4. Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled.
                                  </div>
                        
                          </div>
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
        tab_panel:'home'
      };
    },
    methods:{
      show_panel:function(tab_pane){

        this.tab_panel = tab_pane;
      }
    }
}
</script>