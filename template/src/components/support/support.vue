<template>
     <Breadcrumbs title="Support Ticket" />
         <div class="container-fluid support-ticket">
            <div class="row">
              <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header pb-0">
                    <h5>Ticket List</h5><span>List of ticket opend by customers</span>
                  </div>
                  <div class="card-body">
                    <div class="row">
                     <ticket/>
                     <ticket2/>
                    </div>
                    <supporttable/>
                  </div>
                </div>
              </div>
            </div>
         </div>
</template>
<script>
import ticket from "./ticket.vue"
import ticket2 from "./ticket2.vue"
import supporttable from "./supporttable.vue"
export default {
    components:{
        ticket,
        ticket2,
        supporttable
    },
    data(){
        return{
     
        }
    },
}
</script>