<template>
    <div class="col-md-12">
					<div class="card">
						<div class="card-header">
							<h5 class="card-title"> Filter</h5>
						</div>
						<div class="card-body">
                <div class="table-responsive">
                <div id="basic-1_wrapper" class="dataTables_wrapper no-footer">    
                <div class="row filter-smart">
								<div class="col-sm-2"><input class="form-control" v-model="filtername" placeholder="Name"/>	</div>
								<div class="col-sm-2">	<input class="form-control" v-model="filterposition" placeholder="Position" /></div>
								<div class="col-sm-2"><input class="form-control" v-model="filteoffice" placeholder="Office" /></div>
								<div class="col-sm-2"><input class="form-control" v-model="filterage" placeholder="Age" /></div>
								<div class="col-sm-2"><input class="form-control" v-model="filterstartdate" placeholder="Start Date" /></div>
								<div class="col-sm-2"><input class="form-control" v-model="filtersalary" placeholder="Salary" /></div>
							</div>
            <table class=" display dataTable" id="basic-1"  >
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
         <tr v-for="row in get_rows" :key="row" >
                <td >{{ row.name }}</td>
                <td>{{ row.position }}</td>
                <td>{{ row.office }}</td>
                <td>{{ row.age }}</td>
                <td>{{ row.startdate }}</td>
                <td>{{ row.salary }}</td>
              </tr>
            </tbody>
          </table>
                   <ul class="pagination">
                   <li class="page-item"><a class="page-link" href="#">Previous</a></li>
                <li class="page-item" v-for="i in num_pages()" :key="i" v-bind:class="[i == currentPage ? 'active' : '']" v-on:click="change_page(i)"><a class="page-link">{{i}}</a></li>
                 <li class="page-item"><a class="page-link" href="#">Next</a></li>
                   </ul>
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
      perPage: 10,
      pageOptions: [10, 25, 50, 100],
      filtername: "",
      filterposition:"",
      filteroffice:"",
      filterage:"",
      filterstartdate:"",
      filtersalary:"",
      allData: [],
    };
  },
  watch: {
    filteroffice(search) {
      var filteredData = this.tableItems.filter((item) => {
        return (
          item.office.toLowerCase().includes(search.toLowerCase()) 
        );
      });
      search=="" ? this.allData = this.tableItems :this.allData = filteredData
    },
     filterage(search) {
      var filteredData = this.tableItems.filter((item) => {
        return (
          String(item.age).toLowerCase().includes(search.toLowerCase()) 
        );
      });
      search=="" ? this.allData = this.tableItems :this.allData = filteredData
    },
     filterstartdate(search) {
      var filteredData = this.tableItems.filter((item) => {
        return (
          item.startdate.toLowerCase().includes(search.toLowerCase()) 
        );
      });
      search=="" ? this.allData = this.tableItems :this.allData = filteredData
    },
     filtersalary(search) {
      var filteredData = this.tableItems.filter((item) => {
        return (
          item.salary.toLowerCase().includes(search.toLowerCase())
        );
      });
      search=="" ? this.allData = this.tableItems :this.allData = filteredData
    },
     filterposition(search) {
      var filteredData = this.tableItems.filter((item) => {
        return (
          item.position.toLowerCase().includes(search.toLowerCase()) 
        );
      });
      search=="" ? this.allData = this.tableItems :this.allData = filteredData
    },
    filtername(search) {
      var filteredData = this.tableItems.filter((item) => {
        return (
          item.name.toLowerCase().includes(search.toLowerCase())
        );
      });
      search=="" ? this.allData = this.tableItems :this.allData = filteredData
    },
  },
  computed: {
    ...mapState({
      tableItems: (state) => state.table.items,
    }),

    columns() {
      if (this.allData.length == 0) {
        return [];
      }
      return Object.keys(this.allData[0]);
    },
    get_rows() {
      var start = (this.currentPage - 1) * this.elementsPerPage;
      var end = start + this.elementsPerPage;
      return this.allData.slice(start, end);
    },
  },
  methods: {
    num_pages() {
      return Math.ceil(this.allData.length / this.elementsPerPage);
    },
    change_page(page) {
      this.currentPage = page;
    },
  },
  created() {
    this.allData = this.tableItems;
  },
};
</script>