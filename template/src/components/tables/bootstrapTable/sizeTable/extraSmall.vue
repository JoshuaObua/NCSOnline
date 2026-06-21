<template>             
            
              <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Extra Small table</h5><span>Example of small table, Add <code>.table-xs</code> class to the <code>.table</code> to create a table with extra small spacing. Small table all rows have <code>1.5rem</code> height</span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-xs">
                      <thead>
                        <tr>
                          <th scope="col">#</th>
                          <th scope="col">First</th>
                          <th scope="col">Last</th>
                          <th scope="col">Handle</th>
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