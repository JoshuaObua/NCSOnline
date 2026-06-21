<template>
  <div class="col-sm-12 ">
    <div class="card">
      <div class="card-header">
        <h5>Alternative pagination</h5><span>The default page control presented by DataTables (forward and backward
          buttons with up to 7 page numbers in-between) is fine for most situations, but there are cases where you may
          wish to customise the options presented to the end user. This is done through DataTables' extensible
          pagination mechanism, the pagingType option.</span>
      </div>
      <div class="card-body">
        <div class="table-responsive">
          <div id="basic-1_wrapper" class="dataTables_wrapper no-footer">
            <table class="table display  dataTable " id="basic-1">
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
                  <td>{{ row.name }}</td>
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
              <li class="page-item"><a class="page-link" href="#">Previous</a></li>
              <li class="page-item" v-for="i in num_pages()" :key="i" v-bind:class="[i == currentPage ? 'active' : '']"
                v-on:click="change_page(i)"><a class="page-link">{{ i}}</a></li>
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
      ascending: false,
      sortColumn: '',
      pageOptions: [5, 10, 15],
    };
  },
  computed: {
    ...mapState({
      tableItems: (state) => state.table.items,
    }),
    columns() {
      if (this.tableItems.length == 0) {
        return [];
      }
      return Object.keys(this.tableItems[0])
    }
  },
  methods: {
    num_pages() {
      return Math.ceil(this.tableItems.length / this.elementsPerPage);
    },
    change_page(page) {
      this.currentPage = page;
    },
    get_rows() {
      var start = (this.currentPage - 1) * this.elementsPerPage;
      var end = start + this.elementsPerPage;
      return this.tableItems.slice(start, end);
    },
  }
};
</script>