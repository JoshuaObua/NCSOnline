<template>
    <div class="form-group row">
        <label class="col-sm-3 col-form-label text-end">Default</label>
        <div class="col-xl-5 col-sm-9">
            <datepicker class="datepicker-here  digits" v-model="date" :format="format" />
        </div>
    </div>
</template>
<script>
import { ref } from 'vue';

export default {
    setup() {
        const date = ref();

        const format = (date) => {
            const day = date.getDate();
            const month = date.getMonth() + 1;
            const year = date.getFullYear();

            return ` ${day}/${month}/${year}`;
        }

        return {
            date,
            format,
        }
    }
}
</script>