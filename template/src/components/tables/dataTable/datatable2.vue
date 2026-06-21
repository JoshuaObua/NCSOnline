<template>
          <div class="card-body">
            <div class="row">
              <div class="delete-datatable col-xl-6">
                <div role="group" class="input-group datatable-btn" style="width: 50%;">
                  <input type="text" placeholder="Type to Search"  class="form-control" v-model="filterQuery" />
                  <div class="input-group-append">
                    <button type="button" class="btn btn-secondary disabled" :disabled="!filter"  @click="filter = ''" >
                      Clear
                    </button>
                  </div>
                </div>
              </div>
              <div class="col-xl-6">
                <div class="dataTables_length" id="basic-1_length" style="text-align: right" >
                  <label >Show
                    <select name="basic-1_length" aria-controls="basic-1" class="" v-model="perPage">
                      <option v-for="option in pageOptions" :key="option" :options="pageOptions"> {{ option }}
                      </option>
                    </select> entries</label>
                </div>
              </div>
            </div>
            <div class="table-responsive">
              <div id="basic-1_wrapper" class="dataTables_wrapper no-footer">
                <table class="display dataTable" id="basic-1">
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
                    <tr v-for="row in get_rows" :key="row">
                      <td>{{ row.name }}</td>
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
                  <li  class="page-item" v-for="i in num_pages()" :key="i" v-bind:class="[i == currentPage ? 'active' : '']" v-on:click="change_page(i)">
                    <a class="page-link">{{ i }}</a>
                  </li>
                   <li class="page-item"><a class="page-link" href="#">Next</a></li>
                </ul>
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
      sortColumn: "",
      perPage: 10,
      pageOptions: [10, 25, 50, 100],
      filterQuery: "",
      allData: [],
    };
  },
  watch: {
    filterQuery(search) {
      var filteredData = this.tableItems.filter((item) => {
        return (
          item.name.toLowerCase().includes(search.toLowerCase()) ||
          item.position.toLowerCase().includes(search.toLowerCase()) ||
          item.office.toLowerCase().includes(search.toLowerCase()) ||
          String(item.age).toLowerCase().includes(search.toLowerCase()) ||
          item.startdate.toLowerCase().includes(search.toLowerCase()) ||
          item.salary.toLowerCase().includes(search.toLowerCase())||
          item.option.toLowerCase().includes(search.toLowerCase())
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