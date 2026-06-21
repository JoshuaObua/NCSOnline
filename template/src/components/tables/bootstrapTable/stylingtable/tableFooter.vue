<template>
              <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Table Footer Styling</h5><span>Use class <code>table-info</code>,<code>table-success</code>,<code>table-success</code>,<code>table-info</code>,<code>table-danger</code>,<code>table-primary</code>,<code>table-secondary</code>,<code>table-light</code>,<code>table-active</code>and also use <code>bg-*</code> inside tfoot  element.</span>
                  </div>
                  <div class="card-block row">
                    <div class="col-sm-12 col-lg-12 col-xl-12">
                      <div class="table-responsive table-border-radius">
                        <table class="table">
                          <thead>
                            <tr>
                              <th scope="col">#</th>
                              <th scope="col">First Name</th>
                              <th scope="col">Last Name</th>
                              <th scope="col">Username</th>
                            </tr>
                          </thead>
                          <tbody v-for="item3 in items3" :key="item3">
                            <tr>
                              <th scope="row">{{item3.no}}</th>
                              <td>{{item3.firstname}}</td>
                              <td>{{item3.lastname}}</td>
                              <td>{{item3.username}}</td>
                            </tr>
                            
                          </tbody>
                          <tfoot class="table-success">
                            <tr>
                              <th scope="row">3</th>
                              <td>Larry</td>
                              <td>the Bird</td>
                              <td>@twitter</td>
                            </tr>
                          </tfoot>
                        </table>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
</template>
<script>
export default {
    data(){
        return{
        items3:[
          { no:1, firstname:'Mark', lastname:'Otto', username:'@mdo', },
          { no:2, firstname:'Jacob', lastname:'Thornton', username:'@fat' },
            ]
        }
    }
}
</script>