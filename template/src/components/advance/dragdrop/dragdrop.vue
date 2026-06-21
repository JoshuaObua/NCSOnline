<template>
    <Breadcrumbs main="Bonus UI" title="drag and  drop"/>
   <div class="container-fluid">
              <div class="card">
                  <div class="card-header pb-0">
                      <h5>Basic Drag & Drop</h5>
                  </div>
                  <div class="card-body">
                      <div class="drag-drop-vue row">
                             <draggable :class="item.class+' '+'col-sm-6'"  v-for="(item,index) in basic" :key="index" group="shared" @start="drag = true" @end="drag = false">
                             <div class="drag-drop-item" v-for="(one,index) in item.children" :key="index">{{one.title}}</div>
                         </draggable>
                      </div>
                  </div>
              </div>
               <div class="card">
                  <div class="card-header pb-0">
                      <h5>Copy Drag & Drop</h5>
                  </div>
                  <div class="card-body">
                      <div class="drag-drop-vue row">
                           <draggable :class="text.class+' '+'col-sm-6'" v-for="(text,index) in dropOne" :key="index" group="shared" @start="drags = true" @end="drags = false">
                             <div class="drag-drop-item" v-for="(items,index) in text.children" :key="index">{{items.title}}</div>
                         </draggable>
                      </div>
                  </div>
               </div>
   </div>
</template>
<script>
import { VueDraggableNext } from 'vue-draggable-next'
export default {
    components:{
     draggable: VueDraggableNext,
    },
    data(){
        return{
            drag:false,
            drags:false,
            basic:[
                {class:"pull-left",children:[{title:"Lorem ipsum dolor sit amet, consectetur"}]},
                {class:"pull-left",children:[{title:"sed do eiusmod tempor incididunt ut labore et dolore magna."}]},
                {class:"pull-left",children:[{title:"Ut enim ad minim veniam, quis nostrud exercitation"}]},
                {class:"pull-left",children:[{title:"ullamco laboris nisi ut aliquip ex ea commodo consequat."}]},
                {class:"pull-right",children:[{title:"Duis aute irure dolor in reprehenderit in voluptate velit esse"}]},
                {class:"pull-right",children:[{title:"cillum dolore eu fugiat nulla pariatur."}]},
                {class:"pull-right",children:[{title:"Ut enim ad minim veniam, quis nostrud exercitation"}]},
                {class:"pull-right",children:[{title:"ullamco laboris nisi ut aliquip ex ea commodo consequat."}]}
            ],
           dropOne:[
            {class:"pull-left",children:[{title:"Lorem ipsum dolor sit amet, consectetuer adipiscing elit."}]},
            {class:"pull-left",children:[{title:"Aenean commodo ligula eget dolor. Aenean massa."}]},
            {class:"pull-left",children:[{title:"Lorem ipsum dolor sit amet, consectetur"}]},
            {class:"pull-left",children:[{title:"sed do eiusmod tempor incididunt ut labore et dolore magna."}]},
            {class:"pull-left",children:[{title:"Ut enim ad minim veniam, quis nostrud exercitation"}]},
            {class:"pull-left",children:[{title:"ullamco laboris nisi ut aliquip ex ea commodo consequat."}]},
            {class:"pull-right",children:[{title:"Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus."}]},
            {class:"pull-right",children:[{title:"Donec quam felis, ultricies nec, pellentesque eu, pretium quis, sem."}]},
            {class:"pull-right",children:[{title:"Duis aute irure dolor in reprehenderit in voluptate velit esse"}]},
            {class:"pull-right",children:[{title:"cillum dolore eu fugiat nulla pariatur."}]},
            {class:"pull-right",children:[{title:"Excepteur sint occaecat cupidatat non proident,"}]},
            {class:"pull-right",children:[{title:"sunt in culpa qui officia deserunt mollit anim id est laborum."}]}
           ],
        }
    }
}
</script>