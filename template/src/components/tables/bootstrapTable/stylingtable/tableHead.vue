<template>
               <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Table head options</h5><span>Use class<code>.table-primary</code> inside thead tr element.</span>
                  </div>
                  <div class="card-block row">
                    <div class="col-sm-12 col-lg-12 col-xl-12">
                      <div class="table-responsive">
                        <table class="table">
                          <thead class="table-primary">
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
          { no:3, firstname:'Larry', lastname:'the Bird', username:'@twitter' },
            ]
        }
    }
}
</script>