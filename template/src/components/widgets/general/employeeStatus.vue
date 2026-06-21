<template>
    <div class="col-xl-6 xl-100 box-col-12">
                <div class="card employee-status">
                  <div class="card-header pb-0 d-flex justify-content-between align-items-center">
                    <h5>EMPLOYEE STATUS</h5>
                  </div>
                  <div class="card-body">
                    <div class="user-status table-responsive">
                      <table class="table table-bordernone">
                        <thead>
                          <tr>
                            <th scope="col">Name</th>
                            <th scope="col">Designation</th>
                            <th scope="col">Skill Level</th>
                            <th scope="col">Experience</th>
                          </tr>
                        </thead>
<employeeStatus1/>
                      </table>
                    </div>
                  </div>
                </div>
              </div>
</template>
<script>
import employeeStatus1 from "./employeeStatus1.vue"
export default {
    components:{
        employeeStatus1
    }
}
</script>