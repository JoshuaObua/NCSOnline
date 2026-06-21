<template>
  <div class="row">

    <div class="col-sm-6">
      <div class="dataTables_length ticket-length" id="basic-1_length">
        <label>Show
          <select name="basic-1_length" aria-controls="basic-1" class="" v-model="perPage" @click="paged($event)">
            <option v-for="option in pageOptions" :key="option" :options="pageOptions"> {{ option }}
            </option>
          </select> entries</label>
      </div>
    </div>
  </div>
  <table class=" display dataTable" id="basic-1">
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
    <tbody v-if="!get_rows.length">
      <tr class="odd">
                  <td valign="top" colspan="6" class="dataTables_empty">
                    No results match your search query
                  </td>
                </tr>
    </tbody>
    <tbody v-if="get_rows.length">
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
  <ul class="pagination mt-3" v-if="get_rows.length">
    <li class="page-item"><a class="page-link" v-on:click="prev()">Previous</a></li>
    <li class="page-item" v-for="i in num_pages()" :key="i" v-bind:class="[i == currentPage ? 'active' : '']"
      v-on:click="change_page(i)"><a class="page-link">{{ i }}</a></li>
    <li class="page-item"><a class="page-link" v-on:click="change()">Next</a></li>
  </ul>
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
      perPage: 10,
      pageOptions: [10, 25, 50, 100],
      filter: null,
      filterQuery: "",
      allData: [],
      start: 0,
      end: 10,
      active: false,
      pageOptions: [5, 10, 15, 20],
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
    },
    get_rows() {
      if (this.active) {
        this.start = 0
      }
      else {
        this.start = (this.currentPage - 1) * this.elementsPerPage;
      }
      this.end = this.start + this.elementsPerPage;
      return this.allData.slice(this.start, this.end);
    },
  },
  methods: {
    num_pages() {
      return Math.ceil(this.tableItems.length / this.elementsPerPage);
    },
    change_page(page) {
      this.currentPage = page;
    },
    change() {
      if (
        this.currentPage < Math.ceil(this.allData.length / this.elementsPerPage)
      ) {
        this.currentPage++;
      }
    },
    prev() {
      if (this.currentPage > 1) {
        this.currentPage--;  
      }
    },
    paged(e) {
      this.active = !this.active
      this.elementsPerPage = e.target.value
    },
    onFiltered(filteredItems) {

      this.tableItems = filteredItems.length;
      this.currentPage = 1;
    }
  },
  created() {
    this.allData = this.tableItems;
  },
};
</script>