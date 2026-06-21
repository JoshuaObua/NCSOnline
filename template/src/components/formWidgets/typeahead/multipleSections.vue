<template>
     <div class="col-sm-12 col-md-6 ">
                  <div class="card">
                    <div class="card-header pb-0">
                      <h5>Multiple Sections with Headers</h5><span>Two datasets that are prefetched, stored, and searched on the client. Highlighting is enabled.</span>
                    </div>
                    <div class="card-body">
                      <div id="multiple-datasets">
                        <form class="theme-form">
                          <div class="form-group">
                            <vue3-simple-typeahead :items="list" class="form-control typeahead form-control" :placeholder="options.placeholder"  @onInput="onInput" @onBlur="onBlur" :minInputLength="options.minInputLength" />
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
				placeholder: 'NBA and NHL teams',
				minInputLength: 1,
			},
			list: [
      'In the Heat of the Night',
          'A Man for All Seasons',
          'Annie Hall',
          'Amadeus',
          'Out of Africa',
          'The Godfather Part II',
          'Shakespeare in Love',
          'It Happened One Night',
          'You Can not Take It with You',
          'Oliver!',
          'One Flew over the Cuckoo Nest'
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