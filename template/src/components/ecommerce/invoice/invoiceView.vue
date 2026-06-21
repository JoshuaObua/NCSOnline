<template>
     <Breadcrumbs main="Ecommerce" title="Invoice" />
     <div class="container-fluid">
       <div class="container invoice">
              <div class="row">
                <div class="col-sm-12">
                  <div class="card">
                    <div class="card-body">                            
                      <div>
                        <div>
                          <div class="row invo-header">
                            <div class="col-sm-6">
                              <div class="media">
                                <div class="media-left"><img class="media-object img-60" src="../../../assets/images/logo/logo-1.png" alt=""></div>
                                <div class="media-body m-l-20">
                                  <h4 class="media-heading f-w-600">Viho</h4>
                                  <p>hello@viho.in<br><span class="digits">289-335-6503</span></p>
                                </div>
                              </div>
                            </div>
                            <div class="col-sm-6">
                              <div class="text-md-end text-xs-center">
                                <h3>Invoice #<span class="digits counter"><number
                                  class="bold counter"
                                  ref="number1"
                                  :from="0"
                                  :to="1069"
                                  :duration="5"
                                  :delay="0"
                                  easing="Power1.easeOut"
                                 /></span></h3>
                                <p>Issued: May<span class="digits"> 27, 2023</span><br>Payment Due: June <span class="digits">27, 2023</span></p>
                              </div>
                            </div>
                          </div>
                        </div>
                       <div class="row invo-profile">
                          <div class="col-xl-4">
                            <div class="media">
                              <div class="media-left"><img class="media-object rounded-circle img-60" src="../../../assets/images/user/1.jpg" alt=""></div>
                              <div class="media-body m-l-20">
                                <h4 class="media-heading f-w-600">Johan Deo</h4>
                                <p>JohanDeo@gmail.com<br><span class="digits">555-555-5555</span></p>
                              </div>
                            </div>
                          </div>
                          <div class="col-xl-8">
                            <div class="text-xl-end" id="project">
                              <h6>Project Description</h6>
                              <p>You're only as good as your last collection, which is an enormous pressure. Jeans represent democracy in fashion. Fashion is about dressing according to what's fashionable.</p>
                            </div>
                          </div>
                        </div>
                        <div>
                         
                              <invoiceTable/>
                           
                          <div class="row mt-3">
                            <div class="col-xl-8">
                              <div>
                                <p class="legal"><strong>Thank you for your business!</strong>  Payment is expected within 31 days; please process this invoice within that time. There will be a 5% interest charge per month on late invoices.</p>
                              </div>
                            </div>
                          </div>
                        </div>
                        <div class="col-sm-12 text-center mt-3">
                        <button class="btn btn btn-primary me-2" type="button" @click="printWindow()">Print</button>
                        <button class="btn btn-secondary" type="button">Cancel</button>
                      </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
       </div>
     </div>
</template>
<script>
import invoiceTable from "./invoiceTable.vue"
export default {
  components:{
    invoiceTable
  },
   methods: {
      printWindow() {
        window.print();
      }
    }
}
</script>