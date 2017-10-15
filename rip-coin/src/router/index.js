import Vue from 'vue';
import Router from 'vue-router';
import RipStream from '@/components/RipStream';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'RipStream',
      component: RipStream,
    },
  ],
});
