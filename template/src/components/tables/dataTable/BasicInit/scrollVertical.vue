<template>
     <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Scroll - vertical</h5><span>This example shows the DataTables table body scrolling in the vertical direction. This can generally be seen as an alternative method to pagination for displaying a large table in a fairly small vertical area, and as such pagination has been disabled here (note that this is not mandatory, it will work just fine with pagination enabled as well!).</span>
                  </div>
                  <div class="card-body">
                    <div class="table-responsive">
                       <div id="basic-1_wrapper" class="dataTables_wrapper no-footer">
                      <table class="table display dataTables" id="basic-1">
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
                      </table>
                      </div>
                  </div>
                </div>
                </div>
     </div>
</template>
<script>
import { mapState } from "vuex";
export default {
  data() {
    return {
       elementsPerPage: 10,
         currentPage: 1,
          ascending: false,
          sortColumn: '',
          pageOptions: [5, 10, 15],
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
     get_rows() {
      var start = (this.currentPage-1) * this.elementsPerPage;
      var end = start + this.elementsPerPage;
      return this.tableItems.slice(start, end);
    },
  }
};
</script>