<template>
  <div class="col-sm-12 ">
    <div class="card">
      <div class="card-header">
        <h5>Feature Unable /Disable</h5>
        <span>Disabling features that you don't
          wish to use for a particular table
          is easily done by setting a variable
          in the initialisation object</span>
        <span>In the following example only the search feature is left enabled (which it is by default).</span>
      </div>
      <div class="card-body">
        <div class="table-responsive">
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