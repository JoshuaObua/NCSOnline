<template>
     <div class="form-group row mb-0">
                          <label class="col-sm-3 col-form-label text-end">Permanently visible Datepicker</label>
                          <div class="col-sm-3">
                              <v-date-picker v-model="date"  />
                          </div>
                        </div>
</template>
<script>
import Datepicker from "vue3-datepicker";
export default {
  components: {
    Datepicker,
   
  },
    data(){
        return{
           date: new Date(),
            date1: new Date(),
             date3: new Date(),
            format: 'MM/dd/yyyy',
        format2: 'MMMM yyyy',orientation: "top right",
        disabledDates:{
          days: [6, 0],
        },
    }
},
}
</script>