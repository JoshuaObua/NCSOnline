<template>
     <div class="col-sm-12 ">
                 <div class="card card-overflow">
                    <div class="card-header pb-0">
                      <h5>RTL Support</h5><span>Dataset shows in Right Side.</span>
                    </div>
                    <div class="card-body">
                      <div id="rtl-support">
                        <form class="theme-form">
                          <div class="form-group">
                             <vue3-simple-typeahead :items="list" class="form-control typeahead form-control" dir="rtl" :placeholder="options.placeholder"  @onInput="onInput" @onBlur="onBlur" :minInputLength="options.minInputLength" />
                          </div>
                        </form>
                      </div>
                    </div>
                  </div>
                </div>
</template>
<script>
export default {
	created() {
		this.listFiltered = this.list;
	},
	data() {
		return {
			options: {
				placeholder: 'Countries',
				minInputLength: 1,
			},
			list: [
				  "India",
                  "USA",
                  "Australia",
                  "UEA",
                  "China"
			],
			listFiltered: [],
		};
	},
	methods: {
		onInput(event) {
			this.listFiltered = event.items;
		},
		onBlur(event) {
			this.listFiltered = event.items;
		},
	},
};
</script>