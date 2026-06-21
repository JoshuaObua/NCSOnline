<template>
     <div class="col-sm-12 col-md-6 ">
                  <div class="card card-overflow">
                    <div class="card-header pb-0">
                      <h5>Custom Templates</h5><span>Custom templates give you full control over how suggestions get rendered</span>
                    </div>
                    <div class="card-body">
                      <div id="custom-templates">
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
				placeholder: 'Oscar winners',
				minInputLength: 1,
			},
			list: [
				  'Lawrence of Arabia',
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