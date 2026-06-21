<template>
  <div class="col-sm-12 ">
    <div class="card">
      <div class="card-header">
        <h5>DOM Positioning</h5><span>
          When customising DataTables for your own usage, you might find that the default position of the feature
          elements (filter input etc) is not quite to your
          liking. To address this issue DataTables takes inspiration from the CSS 3 Advanced Layout Module and provides
          the <code class="option" title="DataTables initialisation option">dom:option</code> initialisation parameter
          which can be set
          to indicate where you wish particular features to appear in the DOM. You can also specify <code>div</code>
          wrapping containers (with an id and / or class) to
          provide complete layout flexibility.</span><span>Each HTML control element presented by DataTables is denoted
          by a single character in the <code class="option" title="DataTables initialisation option">dom:option</code>
          option. For example the <code>l</code> option is used for the <code>L</code>ength changing
          input option.</span>
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
              <li class="page-item">
                <a class="page-link" href="#" aria-label="Previous">
                  <span aria-hidden="true">&laquo;</span>
                </a>
              </li>
              <li class="page-item" v-for="i in num_pages()" :key="i" v-bind:class="[i == currentPage ? 'active' : '']"
                v-on:click="change_page(i)"><a class="page-link">{{ i}}</a></li>
              <li class="page-item">
                <a class="page-link" href="#" aria-label="Next">
                  <span aria-hidden="true">&raquo;</span>
                </a>
              </li>
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