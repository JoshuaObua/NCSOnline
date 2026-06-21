import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

import 'bootstrap/dist/css/bootstrap.css'
import "bootstrap"
import "bootstrap/dist/js/bootstrap.min.js";
import 'bootstrap-icons/font/bootstrap-icons.css'
import './assets/scss/app.scss';
import 'v-calendar/dist/style.css';
import "vue-multiselect/dist/vue-multiselect.css"
import 'vue3-simple-typeahead/dist/vue3-simple-typeahead.css'; 
import 'form-wizard-vue3/dist/form-wizard-vue3.css';
import 'vue3-perfect-scrollbar/dist/vue3-perfect-scrollbar.css'
import 'sweetalert2/dist/sweetalert2.min.css';
import 'vue3-tour/dist/vue3-tour.css'
import 'vue-rate/dist/vue-rate.css'
import '@vuepic/vue-datepicker/dist/main.css'

import Breadcrumbs from './layout/breadCrumbs.vue';
import VueApexCharts from "vue3-apexcharts";
import VCalendar from 'v-calendar';
import VueKanban from 'vue-kanban'
import VueNumber from "vue-number-animation";
import Multiselect from 'vue-multiselect'
import VueFeather from "vue-feather";
import Lightbox from 'vue-easy-lightbox'
import Notifications from '@kyvg/vue3-notification'
import SimpleTypeahead from 'vue3-simple-typeahead';
import Wizard from 'form-wizard-vue3';
import PerfectScrollbar from 'vue3-perfect-scrollbar'
import VueSweetalert2 from 'vue-sweetalert2';
import Toaster from "@meforma/vue-toaster";
import {VueMasonryPlugin} from 'vue-masonry';
import Vue3Tour from 'vue3-tour'
import { quillEditor } from "vue3-quill";
import rate from 'vue-rate'
import Datepicker from '@vuepic/vue-datepicker';
import vueChartist from "vue-chartist"
import VueCountdown from '@chenfengyuan/vue-countdown';

import { createI18n } from 'vue-i18n'
import English from "./locales/en.json"
import Español from "./locales/es.json"
import Deutsch from "./locales/de.json"
import Français from "./locales/fr.json"
import Português from "./locales/pt.json"
import 简体中文 from "./locales/cn.json"
import لعربية from "./locales/ae.json"
import German from "./locales/ge.json"
import Russian from "./locales/ru.json"
import Arabic from "./locales/ar.json"
const i18n = createI18n({ legacy: false, // you must specify 'legacy: false' option
  locale: '',
  messages: {
   English: English,
    Español: Español,
    Deutsch: Deutsch,
    Français: Français,   
   Português: Português,
    简体中文: 简体中文,
    لعربية: لعربية,
    German: German,
    Russian: Russian,
    Arabic: Arabic
    }
  })

createApp(App)
.use(store)
.use(router)
.use(VueNumber)
.use(vueChartist)
.use(require("vue-chartist"))
.use(VueApexCharts)
.use(VueKanban)
.use(VCalendar, {})
.use(Lightbox)
.use(PerfectScrollbar)
.use(VueSweetalert2)
.use(Wizard)
.use(rate)
.use(i18n)
.use(Notifications)
.use(SimpleTypeahead)
.use(Toaster)
.use(quillEditor)
.use(require("vue-chartist"))
.use(Vue3Tour)
.use(VueMasonryPlugin)
.component(VueFeather.name, VueFeather)
.component(VueCountdown.name, VueCountdown)
.component('Breadcrumbs', Breadcrumbs)
.component('Datepicker', Datepicker)
.component('multiselect', Multiselect)
.mount('#app')