<template>
    <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Border Bottom Color </h5><span>EExample of a table head border with<code>custom </code>  color. According to the custom color system options, you can set any of predefined colors by adding  <code>.border-bottom-*</code> class to the table head row. This technique works with all border styles demonstrated above.</span>
                  </div>
                  <div class="table-responsive">
                    <table class="table">
                      <thead>
                        <tr class="border-bottom-primary">
                          <th scope="col">#</th>
                          <th scope="col">First Name</th>
                          <th scope="col">Last Name</th>
                          <th scope="col">Username</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr class="border-bottom-secondary">
                          <th scope="row">3</th>
                          <td>Jacob</td>
                          <td>Thornton</td>
                          <td>@fat</td>
                        </tr>
                        <tr class="border-bottom-success">
                          <th scope="row">3</th>
                          <td>Jacob</td>
                          <td>Thornton</td>
                          <td>@fat</td>
                        </tr>
                        <tr class="border-bottom-info">
                          <th scope="row">3</th>
                          <td>Jacob</td>
                          <td>Thornton</td>
                          <td>@fat</td>
                        </tr>
                        <tr class="border-bottom-warning">
                          <th scope="row">3</th>
                          <td>Jacob</td>
                          <td>Thornton</td>
                          <td>@fat</td>
                        </tr>
                        <tr class="border-bottom-danger">
                          <th scope="row">3</th>
                          <td>Jacob</td>
                          <td>Thornton</td>
                          <td>@fat</td>
                        </tr>
                        <tr>
                          <th scope="row">3</th>
                          <td>Jacob</td>
                          <td>Thornton</td>
                          <td>@fat</td>
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
    },
}
</script>