// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue';
import Vuetify from 'vuetify';
import Vuex from 'vuex';
import '../node_modules/vuetify/dist/vuetify.min.css';
import App from './App';
import router from './router';
import modules from './store';

Vue.use(Vuetify);
Vue.use(Vuex);
Vue.config.productionTip = false;

const store = new Vuex.Store({
  modules,
});


/* eslint-disable no-new */
new Vue({
  el: '#app',
  store,
  router,
  template: '<App/>',
  components: { App },
});
