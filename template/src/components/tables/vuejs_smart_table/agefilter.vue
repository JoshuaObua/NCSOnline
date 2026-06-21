<template>
    <div class="col-md-12">
					<div class="card">
						<div class="card-header">
							<h5 class="card-title">Age Filter</h5>
						</div>
						<div class="card-body">
                             <div class="table-responsive">
                             <div id="basic-1_wrapper" class="dataTables_wrapper no-footer">
                                 <div class="row ">
                                   <div class="col-sm-6">
							                     <div class="col-sm-2">
                                     <div class="input-group bootstrap-touchspin">
                                      <button type="button" class="btn btn-primary btn-square bootstrap-touchspin-down" @click="num8--" ><i class="fa fa-minus"></i></button>
                                      <input class="touchspin form-control" type="text" v-model="num8">
                                      <button type="button" class="btn btn-primary btn-square bootstrap-touchspin-down" @click="num8++" ><i class="fa fa-plus"></i></button>
                                   </div>
                                   </div>
                                   </div>
                                   <div class="col-sm-6">
                                   <div class="col-sm-2">
                                     <div class="input-group bootstrap-touchspin">
                                 <button type="button" class="btn btn-primary btn-square bootstrap-touchspin-down" @click="num--" ><i class="fa fa-minus"></i></button>
                                 <input class="touchspin form-control" type="text" v-model="num">
                                  <button type="button" class="btn btn-primary btn-square bootstrap-touchspin-down" @click="num++" ><i class="fa fa-plus"></i></button>
                                   </div>
                                   </div>
                                   </div>
                                 </div>
                                 <datatable/>
                             </div>
                        </div>
                     </div>
                </div>
    </div>
</template>
<script>
import datatable from "../dataTable/datatable.vue"
import { mapState } from "vuex";
export default {
   components:{
    datatable
  },
  data() {
    return {
       elementsPerPage: 10,
         currentPage: 1,
          ascending: false,
          sortColumn: '',
           perPage: 10,
          pageOptions: [ 10, 25,50,100],
           filter: null,
            num8:19,
            num:22
    };
  },
   watch:{
      num8:function(newValue){
        if(newValue >= 100) {
          this.num8 = 100;
        } else if(newValue <= 0) {
          this.num8 = 0;
        }
      },
       num:function(newValue){
        if(newValue >= 100) {
          this.num = 100;
        } else if(newValue <= 0) {
          this.num = 0;
        }
      },
     },
   computed: {
    ...mapState({
      tableItems:(state)=>state.table.items,    
    }),
   columns() {
      if (this.tableItems.length == 0) {
        return [];
      }
      return Object.keys(this.tableItems[0])
    }
   },
  methods:{
      num_pages() {
      return Math.ceil(this.tableItems.length / this.elementsPerPage);
    },
     change_page(page) {
      this.currentPage = page;
    },
     get_rows() {
      var start = (this.currentPage-1) * this.elementsPerPage;
      var end = start + this.elementsPerPage;
      return this.tableItems.slice(start, end);
    },
   
  },
  created() {
    this.allData = this.tableItems;
  },
};
</script>