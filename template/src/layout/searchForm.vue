<template>
    <input type="text" placeholder="Search.." v-on:keyup="searchterm" v-model="terms" title="" autofocus>
    <div :class="searchResult ? 'Typeahead-menu is-open' : 'Typeahead-menu '" v-if="menuItems.length">
        <div class="ProfileCard u-cf" v-for="(menuItem, index) in menuItems.slice(0, 8)" :key="index">
            <div class="ProfileCard-avatar header-search"><vue-feather :type="menuItem.icon"></vue-feather></div>
            <div class="ProfileCard-details">
                <div class="ProfileCard-realName"><span @click="removeFix()"><router-link :to="{ path: menuItem.path }"
                            class="realname">{{ menuItem.title }}</router-link></span>
                </div>
            </div>
        </div>
    </div>
    <div :class="searchResultEmpty ? 'Typeahead-menu is-open' : 'Typeahead-menu' " >
          <div class="tt-dataset tt-dataset-0">
            <div class="EmptyMessage">
                Opps!! There are no result found.
            </div>
          </div>
        </div>
</template>
<script>
var body = document.getElementsByTagName('body')[0];
 import { mapState } from 'vuex';
export default {
    data(){
    return{
       filtered: false,
        terms: '',
        searchResult: false,
        searchResultEmpty: false,
    }
  },
    computed: {
      ...mapState({
        menuItems: (state) => state.menu.searchData,
        searchOpen: (state) => state.menu.searchOpen
      }),
    },
    watch: {
      menuItems: function() {
        this.terms ? this.addFix() : this.removeFix();
        if (!this.menuItems.length) this.searchResultEmpty = true;
        else this.searchResultEmpty = false;
      },
    },
  methods:{
     collapseFilter() {
        this.filtered = !this.filtered;
      },
       searchterm: function() {
        this.$store.dispatch('menu/searchTerm', this.terms);
      },
      addFix() {
        body.classList.add('offcanvas');
        this.searchResult = true;
      },
      removeFix() {
        body.classList.remove('offcanvas');
        this.searchResult = false;
        this.terms = '';
      },
  },
}
</script>