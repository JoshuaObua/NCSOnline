<template>
    <div class="col-sm-12 ">
                <div class="card">
                  <div class="card-header">
                    <h5>Text or background utilities</h5><span>Regular table background variants are not available with the inverse table, however, you may use Classes <code>bg-dark</code>,<code>bg-warning</code>,<code>bg-success</code>,<code>bg-info</code>,<code>bg-danger</code>,<code>bg-primary</code>,<code>bg-secondary</code>,<code>bg-light</code> in<code>td</code></span>
                  </div>
                  <div class="table-responsive">
                    <table class="table table-borderedfor">
                      <thead>
                        <tr>
                          <th scope="col">#</th>
                          <th scope="col">Heading</th>
                          <th scope="col">Heading</th>
                        </tr>
                      </thead>
                      <tbody v-for="item2 in items2" :key="item2">
                        <tr :class="item2.class">
                          <td >{{item2.rowVariant}}</td>
                          <td >{{item2.heading}}</td>
                          <td >{{item2.heading}}</td>
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
              items2: [
          { class:'bg-primary', heading:'primary', rowVariant: '1' },
          { class:'bg-secondary', heading:'secondary', rowVariant: '2' },
          { class:'bg-success', heading:'success', rowVariant: '3' },
          { class:'bg-info', heading:'info', rowVariant: '4' },
          { class:'bg-warning', heading:'warning', rowVariant: '5' },
          { class:'bg-danger', heading:'danger', rowVariant: '6' },
          { class:'bg-light', heading:'light', rowVariant: '7' },
        
        ],
        }
    }
}
</script>