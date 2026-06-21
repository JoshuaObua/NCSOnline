<template>
     <div class="card-block row">
                    <div class="col-sm-12 col-lg-12 col-xl-12">
                      <div class="table-responsive">
                        <table class="table table-responsive-sm" >
                          <thead >
                            <tr>
                              <th scope="col">#</th>
                              <th scope="col">Table heading</th>
                              <th scope="col">Table heading</th>
                              <th scope="col">Table heading</th>
                              <th scope="col">Table heading</th>
                              <th scope="col">Table heading</th>
                            </tr>
                          </thead>
                         <tbody v-for="item4 in items4" :key="item4">
                            <tr>
                              <th scope="row">{{item4.no}}</th>
                              <td>{{item4.Tableheading}}</td>
                              <td>{{item4.Tableheading}}</td>
                              <td>{{item4.Tableheading}}</td>
                              <td>{{item4.Tableheading}}</td>
                              <td>{{item4.Tableheading}}</td>
                            </tr>
                          </tbody>
                        </table>
                      </div>
                    </div>
                  </div>
                   <div class="card-block row">
                    <div class="col-sm-12 col-lg-12 col-xl-12">
                      <div class="table-responsive">
                        <table class="table table-responsive-sm" >
                          <thead >
                            <tr>
                              <th scope="col">#</th>
                              <th scope="col">Table heading</th>
                              <th scope="col">Table heading</th>
                              <th scope="col">Table heading</th>
                              <th scope="col">Table heading</th>
                              <th scope="col">Table heading</th>
                            </tr>
                          </thead>
                         <tbody v-for="item4 in items4" :key="item4">
                            <tr>
                              <th scope="row">{{item4.no}}</th>
                              <td>{{item4.Tableheading}}</td>
                              <td>{{item4.Tableheading}}</td>
                              <td>{{item4.Tableheading}}</td>
                              <td>{{item4.Tableheading}}</td>
                              <td>{{item4.Tableheading}}</td>
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
              items4: [
          { no:1,Tableheading :'Table cell'},
          { no:2,Tableheading :'Table cell'},
          {no:3,Tableheading :'Table cell'}
        ]
        }
    }
}
</script>