<template>
  <div class="col-sm-12 ">
    <div class="card">
      <div class="card-header">
        <h5>State saving</h5><span>
          DataTables has the option of being able to save the state of a table (its paging position, ordering state etc)
          so that is can be restored when the user
          reloads a page, or comes back to the page after visiting a sub-page. This state saving ability is enabled by
          the <code class="option" title="DataTables initialisation option">stateSave</code> option.</span><span>The
          built in state saving method uses the HTML5 <code>localStorage</code> and <code>sessionStorage</code> APIs for
          efficient storage of the data. Please
          note that this means that the built in state saving option <strong>will not work with IE6/7</strong> as these
          browsers do not support these APIs. Alternative
          options of using cookies or saving the state on the server through Ajax can be used through the <code
            class="option" title="DataTables initialisation option">stateSaveCallback</code> and <a
            href="//datatables.net/reference/option/stateLoadCallback"><code class="option"
              title="DataTables initialisation option">stateLoadCallback</code></a> options.</span><span>The duration
          for which the saved state is valid and can be used to restore the table state can be set using the <code
            class="option" title="DataTables initialisation option">stateDuration</code> initialisation
          parameter (2 hours by default). This parameter also controls if <code>localStorage</code> (0 or greater) or
          <code>sessionStorage</code> (-1) is used to store
          the data.</span><span>The example below simply shows state saving enabled in DataTables with the <code
            class="option" title="DataTables initialisation option">stateSave</code> option.</span>
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