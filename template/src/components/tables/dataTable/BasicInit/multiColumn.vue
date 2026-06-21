<template>
  <div class="col-sm-12 ">
    <div class="card">
      <div class="card-header">
        <h5>Multi-column ordering</h5><span>DataTables allows ordering by multiple columns at the same time, which can
          be activated in a number of different ways:</span><span>User shift click on a column (added the clicked column
          as a secondary, tertiary etc ordering column).</span><span>On a per-column basis (i.e. order by a specific
          column and then a secondary column if the data in the first column is identical), through
          the<code>columns.orderData option.</code></span><span>Using the <code>columns.orderData option</code>option to
          specify a multiple column order by default (for example [ 0, 1 ]).</span><span>Through the
          <code>order() API </code>API method.</span><span>Note that, the ability for the user to shift click to order
          multiple columns can be disabled through the <code>orderMulti option</code>option.</span><span>The example
          below shows the first column having a secondary order applied to the second column in the table, vice-versa
          for the second column being tied directly to the first and the salary column to the first name column.</span>
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
  }
};
</script>