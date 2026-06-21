<template>
    <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Vertical Borders</h5><span>Example of <code>Vertical </code> table borders. This is a default table border style attached to <code> .table</code> class.All borders have the same grey color and style, table itself doesn't have a border, but you can add this border using<code> .table-border-vertical</code>class added to the table with <code>.table</code>class.</span>
                  </div>
                  <div class="table-responsive mb-0">
                    <table class="table table-border-vertical">
                      <thead>
                        <tr>
                          <th scope="col">#</th>
                          <th scope="col">First Name</th>
                          <th scope="col">Last Name</th>
                          <th scope="col">Username</th>
                          <th scope="col">Role</th>
                          <th scope="col">Country</th>
                        </tr>
                      </thead>
                      <tbody v-for="item in items" :key="item">
                        <tr>
                          <th scope="row">{{item.no}}</th>
                          <td>{{item.firstname}}</td>
                          <td>{{item.lastname}}</td>
                          <td>{{item.username}}</td>
                          <td>{{item.role}}</td>
                          <td>{{item.country}}</td>
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
          { no:1, firstname:'Alexander', lastname:'Orton', username:'@mdorton', role:'Admin', country:'USA' },
          { no:2, firstname:'John Deo', lastname:'Deo', username:'@johndeo', role:'User', country:'USA' },
          { no:3, firstname:'Randy Orton', lastname:'the Bird', username:'@twitter', role:'admin', country:'UK' },
          { no:4, firstname:'Randy Mark', lastname:'Ottandy', username:'@mdothe', role:'user', country:'AUS' },
          { no:5, firstname:'Ram Jacob', lastname:'Thornton', username:'@twitter', role:'admin', country:'IND' }
        ],
        }
    }
}
</script>