<template>
  <div class="col-sm-12 ">
    <div class="card">
      <div class="card-header">
       <h5>Multiple tables</h5><span>Often you might wish to initialise multiple tables with a single statement. This is trivially done by using a jQuery selector which will pick up multiple tables.</span><span>The tables are independent for user control (i.e. user controlled paging on one table does not effect the others), but they do share the initialisation parameters given (for example if you specific the Spanish language file, all tables will be shown in Spanish). Additionally, the API can be used to manipulate both together, or independently.</span><span>The example below shows two tables initialised with a single line of code, through the use of the table.display selector (i.e. select all elements which have the class of table.display (which is suitable in this example, you might wish to use a different selector).</span>
      </div>
      <div class="card-body">
        <div class="table-responsive">
           <div id="basic-1_wrapper" class="dataTables_wrapper no-footer">
          <table class="table display dataTable" id="basic-1">
            <thead>
              <tr>
              <th>Name</th>
              <th>Position</th>
              <th>Office</th>
              <th>Age</th>
              <th>Start date</th>
              <th>Salary</th>
              </tr>
            </thead>
            <tbody>        
         <tr v-for="row in get_rows()" :key="row">
                <td >{{ row.name }}</td>
                <td>{{ row.position }}</td>
                <td>{{ row.office }}</td>
                <td>{{ row.age }}</td>
                <td>{{ row.startdate }}</td>
                <td>{{ row.salary }}</td>
              </tr>
            </tbody>
            <thead>
              <tr>
              <th>Name</th>
              <th>Position</th>
              <th>Office</th>
              <th>Age</th>
              <th>Start date</th>
              <th>Salary</th>
              </tr>
            </thead>
          </table>
  <ul class="pagination">
    <li class="page-item">
      <a class="page-link" href="#" aria-label="Previous">
        <span aria-hidden="true">&laquo;</span>
      </a>
    </li>
      <li class="page-item" v-for="i in num_pages()" :key="i"
         v-bind:class="[i == currentPage ? 'active' : '']"
         v-on:click="change_page(i)"><a class="page-link">{{i}}</a></li>
      <li class="page-item">
        <a class="page-link" href="#" aria-label="Next">
        <span aria-hidden="true">&raquo;</span>
        </a>
         </li>
         </ul>
       </div>
      </div>
    </div>
    <multipleTables2/>
  </div>
  </div>
</template>
<script>
import multipleTables2 from "./multipleTables2.vue"
import { mapState } from "vuex";
export default {
    components:{multipleTables2},
  data() {
    return {
       elementsPerPage: 10,
         currentPage: 1,
          ascending: false,
          sortColumn: '',
    };
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
  }
};
</script>