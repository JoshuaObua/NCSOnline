<template>
    <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Hoverable rows</h5><span>Use a class <code>table-hover</code> to enable a hover state on table rows within a <code>tbody</code>.</span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-hover">
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
                    <h5>Contextual classes</h5><span>Use contextual classes to color table rows or individual cells. you may use Classes <code>table-primary</code>,<code>table-secondary</code>,<code>table-success</code>,<code>table-info</code>,<code>table-warning</code>,<code>table-danger</code>,<code>table-light</code>,<code>table-active</code> in<code>tr</code></span>
                  </div>
                  <div class="table-responsive">
                    <table class="table">
                      <thead>
                        <tr>
                          <th scope="col">Class</th>
                          <th scope="col">Heading</th>
                          <th scope="col">Heading</th>
                        </tr>
                      </thead>
                      <tbody v-for="item2 in items2" :key="item2">
                        <tr :class="item2.class">
                          <th scope="row">{{item2.rowVariant}}</th>
                          <td>{{item2.heading}}</td>
                          <td>{{item2.heading}}</td>
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
         items2: [
          { class:'table-primary', heading:'Cell', rowVariant: 'primary' },
          { class:'table-secondary', heading:'Cell', rowVariant: 'secondary' },
          { class:'table-success', heading:'Cell', rowVariant: 'success' },
          { class:'table-info', heading:'Cell', rowVariant: 'info' },
          { class:'table-warning', heading:'Cell', rowVariant: 'warning' },
          { class:'table-danger', heading:'Cell', rowVariant: 'danger' },
          { class:'table-light', heading:'Cell', rowVariant: 'light' },
          { class:'table-active', heading:'Cell', rowVariant: 'active' },
        ],
        }
    }
}
</script>