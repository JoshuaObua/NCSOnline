<template>
      <div class="row"> 
              <div class=" col-xl-6">
                <div class="dataTables_length" id="basic-1_length"  >
                  <label >Show
                    <select name="basic-1_length" aria-controls="basic-1" class="" v-model="perPage" @click="paged($event)">
                      <option v-for="option in pageOptions" :key="option" :options="pageOptions"> {{ option }}
                      </option>
                    </select> entries</label >
                </div>
              </div>  
            </div>
                  <div class="order-history table-responsive">
                      <div id="basic-1_wrapper" class="dataTables_wrapper no-footer">  
                      <table class="table table-bordernone display dataTable no-footer" id="basic-1" role="grid" aria-describedby="basic-1_info">
                        <thead>
                          <tr role="row"  >
                            <th scope="col">Prdouct</th>
                            <th scope="col" >Prdouct name</th>
                            <th scope="col" >Size</th>
                            <th scope="col" >Color</th>
                            <th scope="col" >Article number</th>
                            <th scope="col">Units</th>
                            <th scope="col">Price</th>
                            <th scope="col" ><i class="fa fa-angle-down"></i></th>
                          </tr>
                        </thead>
                        <tbody> 
                           <tr role="row" v-for="row in get_rows" :key="row">
                            <td class="sorting_1"><router-link to="/ecommerce/indexProduct"><img class="img-fluid img-30" :src="getImgUrl(row.images)" alt="#"></router-link></td>
                            <td><div class="product-name"><a href="#">{{row.product1}}</a>
                            <div class="order-process"><span class="order-process-circle" :class="getStatusClass(row.orderStatus)"></span>{{row.product2}}</div>
                            </div></td>
                            <td>{{row.size}}</td>
                            <td>{{row.color}}</td>
                            <td>{{row.articlenumber}}</td>
                            <td>{{row.units}}</td>
                            <td>{{row.price}}</td>
                            <td><vue-feather type="more-vertical"></vue-feather></td>
                          </tr>  
                        </tbody>
                      </table>
                    </div>
                  </div>
</template>
<script>
import { mapState } from "vuex";
export default {
     data(){
       return{
         elementsPerPage: 10,
         currentPage: 1,
          ascending: false,
          sortColumn: '',
           perPage: 10,
           pageOptions: [5, 10, 15, 20],
          allData: [],
           filterQuery: "",
       }
     },
     computed: {
    ...mapState({
      orderItems:(state)=>state.order.product,    
    }),
    columns() {
      if (this.allData.length == 0) {
        return [];
      }
      return Object.keys(this.allData[0]);
    },
    get_rows() {
      if(this.active){
        this.start=0
      }
    else{
      this.start = (this.currentPage - 1) * this.elementsPerPage;
    }
      this.end = this.start + this.elementsPerPage;
  return this.allData.slice(this.start,this.end);
    },
     },
      watch: {
    filterQuery(search) {
      var filteredData = this.orderItems.filter((product) => {
        return (
          product.product1.toLowerCase().includes(search.toLowerCase()) ||
           product.product2.toLowerCase().includes(search.toLowerCase()) ||
            product.size.toLowerCase().includes(search.toLowerCase()) ||
             product.color.toLowerCase().includes(search.toLowerCase()) ||
              product.articlenumber.toLowerCase().includes(search.toLowerCase()) ||
            product.price.toLowerCase().includes(search.toLowerCase())
        );
      });
      search=="" ? this.allData = this.orderItems :this.allData = filteredData
    },
  },
   created() {
    this.allData = this.orderItems;
  },
     methods:{
       getStatusClass(orderStatus){
         if(orderStatus==='shipped'){
           return 'shipped-order'
         }else if(orderStatus=='processing')
         return ''
         else if(orderStatus ==='canclled')
         return 'cancel-order'
         else return ''
       },
       paged(e){
        this.active=!this.active
        this.elementsPerPage =e.target.value

      },
        num_pages() {
      return Math.ceil(this.allData.length / this.elementsPerPage);
    },
    change_page(page) {
      this.currentPage = page;
    },
    change(){
      this.currentPage ++ ;
    },
    prev(){
      this.currentPage -- ;
    },
     getImgUrl(path) {
      return require("../../../assets/images/product/" + path);
    },
     }
}
</script>