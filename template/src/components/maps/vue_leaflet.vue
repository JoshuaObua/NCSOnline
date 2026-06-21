<template>
    <Breadcrumbs main="Maps" title="Vue Leaflet"/>
        <div class="container-fluid">
            <div class="row">
                <div class="col-xl-6 ">
                    <div class="card">
                        <div class="card-header pb-0">
                          <h5>Simple Map</h5>
                      </div>
                      <div class="card-body">
                      <simplemap/>
                      </div>
                    </div>
                </div>
                <div class="col-xl-6 ">
                    <div class="card">
                        <div class="card-header pb-0">
                            <h5>Multi-language support</h5>
                        </div>
                        <div class="card-body">
                            <multilanguage/>
                        </div>
                    </div>
                </div>
                  <div class="col-xl-12 ">
                    <div class="card">
                        <div class="card-header pb-0">
                            <h5>Map With Marker</h5>
                        </div>
                        <div class="card-body">
                           <multilanguage2/>
                        </div>
                    </div>
                </div>
            </div>
        </div>
</template>
<script>
import simplemap from "./maps/simplemap.vue"
import multilanguage from "./maps/multilanguage.vue"
import multilanguage2 from "./maps/multilanguage2.vue"
export default {
    components:{
        simplemap,
        multilanguage,
        multilanguage2
    }
}
</script>