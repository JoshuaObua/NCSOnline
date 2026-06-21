<template>
    <div class="col-md-12">
					<div class="card">
						<div class="card-header">
							<h5 class="card-title">Basic</h5>
						</div>
						<div class="card-body">
                             <div class="table-responsive">
                             <div id="basic-1_wrapper" class="dataTables_wrapper no-footer">
                                 <datatable/>
                             </div>
                        </div>
                     </div>
                </div>
    </div>
</template>
<script>
import datatable from "../dataTable/datatable.vue"
export default {
  components:{
    datatable
  }
}
</script>