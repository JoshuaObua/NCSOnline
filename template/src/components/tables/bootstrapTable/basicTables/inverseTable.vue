<template>
    <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Inverse Table with Primary background</h5><span>Use a class <code>.bg-info, .bg-success, .bg-warning and .bg-danger classes.</code> with light text on dark backgrounds inside table element.<br>To set the light background color use .bg-[color] class where [color] is the value of your selected color from stack color palette. So for teal color background class will be .bg-teal</span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-striped bg-primary">
                      <thead class="tbl-strip-thad-bdr">
                        <tr>
                          <th scope="col">#</th>
                          <th scope="col">First Name</th>
                          <th scope="col">Email</th>
                          <th scope="col">Gender</th>
                          <th scope="col">Country</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <th scope="row">1</th>
                          <td>Marybelle</td>
                          <td>mdellascala0@opensource.org</td>
                          <td>Female</td>
                          <td>USA</td>
                        </tr>
                        <tr>
                          <th scope="row">2</th>
                          <td>Octavia </td>
                          <td>otohill1@google.co.jp</td>
                          <td>Female</td>
                          <td>USA</td>
                        </tr>
                        <tr>
                          <th scope="row">3</th>
                          <td>Jennie </td>
                          <td>jgeroldi2@slideshare.net</td>
                          <td>Female</td>
                          <td>UK</td>
                        </tr>
                        <tr>
                          <th scope="row">4</th>
                          <td>vwharby3@mozilla.org</td>
                          <td>@mdothe</td>
                          <td>user</td>
                          <td>AUS</td>
                        </tr>
                        <tr>
                          <th scope="row">5</th>
                          <td>Olenka</td>
                          <td>obrawleyc@adobe.com</td>
                          <td>Female</td>
                          <td>IND</td>
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