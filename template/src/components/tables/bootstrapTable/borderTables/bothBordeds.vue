<template>
              <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Both Borders</h5><span>Example of <code>horizontal</code> table borders. This is a default table border style attached to <code> .table</code> class.All borders have the same grey color and style, table itself doesn't have a border, but you can add this border using<code> .table-bordered</code>class added to the table with <code>.table</code>class.</span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-bordered">
                      <thead>
                        <tr>
                          <th scope="col">#</th>
                          <th scope="col">First Name</th>
                          <th scope="col">Last Name</th>
                          <th scope="col">Username</th>
                        </tr>
                      </thead>
                     <tbody v-for="item in items" :key="item">
                        <tr>
                          <th scope="row">{{item.no}}</th>
                          <td>{{item.firstname}}</td>
                          <td>{{item.lastname}}</td>
                          <td>{{item.username}}</td>
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
               items: [
          { no:1, firstname:'Mark', lastname:'Otto', username:'@mdo', },
          { no:2, firstname:'Jacob', lastname:'Thornton', username:'@fat' },
          { no:3, firstname:'Larry', lastname:'the Bird', username:'@twitter' },
        ],
        }
    }
}
</script>