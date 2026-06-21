<template>
     <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Extra Large Tables</h5><span>Example of Extra large table, Add <code>.table-xl</code> class to the <code>.table</code> to create a table with extra large spacing. Extra larger table all rows have <code>1.25rem</code> height.</span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-xl">
                      <thead>
                        <tr>
                          <th>#</th>
                          <th>First Name</th>
                          <th>Last Name</th>
                          <th>Username</th>
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