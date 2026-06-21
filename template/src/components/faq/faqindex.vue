<template>
    <Breadcrumbs main="" title="FAQ" />
        <div class="container-fluid">
            <div class="faq-wrap">
              <div class="row">
                  <faq/>
                <faq1/>
                <featured/>
                <articles/>
              </div>
            </div>
        </div>
</template>
<script>
import faq from "./faq.vue"
import faq1 from "./faq1.vue"
import featured from "./featured.vue"
import articles from "./articles.vue"
export default {
    components:{
          faq,
          faq1,
          featured,
          articles
    }
}
</script>