<template>
     <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Basic Table</h5><span>Use a class<code>table</code> to any table.</span>
                  </div>
                  <div class="table-responsive">
                    <table class="table" >
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
              <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Inverse Table</h5><span>Use a class <code>table-inverse</code> inside table element.</span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-inverse">
                      <thead>
                        <tr>
                          <th scope="col">#</th>
                          <th scope="col">First Name</th>
                          <th scope="col">Last Name</th>
                          <th scope="col">Date</th>
                          <th scope="col">Invoice</th>
                          <th scope="col">Action</th>
                        </tr>
                      </thead>
                      <tbody v-for="detail in details" :key="detail">
                        <tr>
                          <th scope="row">{{detail.no}}</th>
                          <td>{{detail.firstname}}</td>
                          <td>{{detail.lastname}}</td>
                          <td>{{detail.date}}</td>
                          <td>{{detail.invoice}}</td>
                        <td><vue-feather type="more-vertical"></vue-feather></td>
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
        details:[
          {no:1,firstname:'Bobby',lastname:'Davis',date:'Nov 14, 2021', invoice:"$2,410",role:'Admin'},
          {no:2,firstname:'Christopher',lastname:'Neal',date:'Nov 21, 2021', invoice:"$2,450", role:'User'},
          {no:3,firstname:'Monkey',lastname:'Karry',date:'October 5, 2021', invoice:"$3,500", role:'admin'},
          {no:4,firstname:'James',lastname:'White',date:'Nov 24, 2021', invoice:"$2,250",role:'user'},
          {no:5,firstname:'Aaron',lastname:'James',date:'Nov 25, 2021', invoice:"$7,452",role:'admin'},
        ],
        }
    }
}
</script>