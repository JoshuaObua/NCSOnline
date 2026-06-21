import { createStore } from 'vuex'
import layout from './modules/layout';
import menu from './modules/menu'
import chartData from "./modules/chart"
import common from './modules/common';
import products from './modules/products'
import email from "./modules/email"
import chat from "./modules/chat"
import todo from "./modules/todo"
import table from './modules/table';
import icon from "./modules/icon"
import jobs from "./modules/jobs"
import courses from "./modules/courses"
import support from "./modules/support"
import order from "./modules/order"
import masonry from "./modules/masonry"
import kanban from "./modules/kanban"
import filemaneger from "./modules/filemaneger"
import language from "./modules/language"
export default createStore({
  state: {
  },
  getters: {
  },
  mutations: {
  },
  actions: {
  },
  modules: {
    layout,
    menu,
    chartData,
    common,
    products,
    email,
    chat,
    todo,
    language,
    table,
    icon,
    jobs,
    courses,
    support,
    order,
    masonry,
    kanban,
    filemaneger
  }
})
