<template>
    <Breadcrumbs main="Ecommerce" title="Payment Details"/>
    <div class="container-fluid credit-card payment-details ">
            <div class="row">
                <creditCard/>
                <debitCard/>
                <codCard/>
                <emiCard/>
                <netBanking/>
            </div>
    </div>
</template>
<script>
import creditCard from "./creditCard.vue"
import debitCard from "./debitCard.vue"
import codCard from "./codCard.vue"
import emiCard from "./emiCard.vue"
import netBanking from "./netBanking.vue"
export default {
    components:{
        creditCard,
        debitCard,
        codCard,
        emiCard,
        netBanking
    }
}
</script>