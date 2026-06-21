<template>
  <div class="row">
    <div class="delete-datatable col-sm-6">
      <div role="group" class="input-group datatable-btn ticket-table">
        <input placeholder="Search..." class="form-control" v-model="filterQuery" />

        <button type="button" class="btn btn-secondary" @click="resetInput">
          Clear</button>

      </div>
    </div>
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
  <div class="table-responsive mt-3">
    <table class="table display dataTable" id="basic-6">
      <thead>
        <tr>
          <th>Name</th>
          <th>Position</th>
          <th>Salary</th>
          <th>Office</th>
          <th>Skill</th>
          <th>Extn.</th>
          <th>E-mail</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in get_rows" :key="row">
          <td>
            <div class="media"><img class="rounded-circle img-30 me-3"
                :src="require('../../assets/images/' + row.profile)" alt="Generic placeholder image">
              <div class="media-body align-self-center">
                <div>{{ row.name }}</div>
              </div>
            </div>
          </td>
          <td>{{ row.position }}</td>
          <td>{{ row.salary }}</td>
          <td>{{ row.office }}</td>
          <td>
            <div class="progress-showcase">
              <div class="progress sm-progress-bar">
                <div class="progress-bar" :class="row.StatusClass" role="progressbar" :style="{ 'width': row.skill }"
                  aria-valuenow="50" aria-valuemin="0" aria-valuemax="100"></div>
              </div>
            </div>
          </td>
          <td>{{ row.extn }}</td>
          <td>{{ row.email }}</td>
        </tr>
      </tbody>
    </table>
    <ul class="pagination mt-3">
      <li class="page-item"><a class="page-link" v-on:click="prev()">Previous</a></li>
      <li class="page-item" v-for="i in num_pages()" :key="i" v-bind:class="[i == currentPage ? 'active' : '']"
        v-on:click="change_page(i)">
        <a class="page-link">{{ i }}</a>
      </li>
      <li class="page-item"><a class="page-link" v-on:click="change()">Next</a></li>
    </ul>
  </div>
</template>
<script>
import { mapState } from "vuex";
import ticket from "./ticket.vue"
import ticket2 from "./ticket2.vue"
export default {
  components: {
    ticket,
    ticket2
  },
  data() {
    return {
      elementsPerPage: 10,
      currentPage: 1,
      ascending: false,
      sortColumn: "",
      perPage: 10,
      pageOptions: [10, 25],
      filterQuery: "",
      allData: [],
      start: 0,
      end: 10,
      active: false
    }
  },
  watch: {
    filterQuery(search) {
      var filteredData = this.supportItems.filter((item) => {
        return (
          item.name.toLowerCase().includes(search.toLowerCase()) ||
          item.position.toLowerCase().includes(search.toLowerCase()) ||
          item.office.toLowerCase().includes(search.toLowerCase()) ||
          item.extn.toLowerCase().includes(search.toLowerCase()) ||
          item.salary.toLowerCase().includes(search.toLowerCase()) ||
          item.email.toLowerCase().includes(search.toLowerCase())
        );
      });
      search == "" ? this.allData = this.supportItems : this.allData = filteredData

    },
  },
  computed: {
    ...mapState({
      supportItems: (state) => state.support.items,
    }),
    columns() {
      if (this.allData.length == 0) {
        return [];
      }
      return Object.keys(this.allData[0]);
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
    paged(e) {
      this.active = !this.active
      this.elementsPerPage = e.target.value

    },
    getImageUrl(path) {
      return require("../../assets/images/user/" + path);
    },
    num_pages() {
      return Math.ceil(this.allData.length / this.elementsPerPage);
    },
    change_page(page) {
      this.currentPage = page;
    },
    change() {
      this.currentPage++;
    },
    prev() {
      this.currentPage--;
    },
    resetInput() {

      this.filterQuery = "";
    },
  },

  created() {
    this.allData = this.supportItems;
  },
}
</script>