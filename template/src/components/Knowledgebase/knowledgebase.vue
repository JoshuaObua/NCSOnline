<template>
  <Breadcrumbs  title="Knowledgebase" />
       <div class="container-fluid faq-section">
            <div class="row">
                <div class="col-12 ">
                <div class="knowledgebase-bg">
               
                  </div>
                <div class="knowledgebase-search">
                  <div>
                    <h3>How Can I help you?</h3>
                    <form class="form-inline" action="#" method="get">
                      <div class="form-group w-100 mb-0">
                        <i class="fa fa-search"></i>
                        <input class="form-control-plaintext w-100" type="text" placeholder="Type question here" title="">
                      </div>
                    </form>
                  </div>
                </div>
              </div>
              <faq/>
              <featured/>
            </div>

       </div>
</template>
<script>
import faq from "../faq/faq.vue"
import featured from "./featured.vue"
export default {
    components:{
      faq,
      featured
    }
}
</script>