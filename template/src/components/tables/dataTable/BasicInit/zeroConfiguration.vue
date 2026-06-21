<template>
  <div class="col-sm-12 ">
    <div class="card">
      <div class="card-header">
        <h5>Zero Configuration</h5>
        <span>DataTables has most features enabled by default, so all you need to
          do to use it with your own tables is to call the construction
          function:<code>$().DataTable();</code>.</span><span>Searching, ordering and paging goodness will be
          immediately added to
          the table, as shown in this example.</span>
      </div>
      <div class="card-body">
        <div class="basic-table-overflow">
          <div id="basic-1_wrapper" class="dataTables_wrapper no-footer">

            <datatable />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>

import datatable from "../datatable.vue"
import { mapState } from "vuex";
export default {
  components: {
    datatable
  },

};
</script>