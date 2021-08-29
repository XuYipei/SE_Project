// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from '@/App'
import router from './router'
import { Button, Table } from 'ant-design-vue'

import VueAxios from 'vue-axios';
import axios from 'axios';

Vue.use(Button)
Vue.use(Table)
Vue.config.productionTip = false

new Vue({
  el: '#app',
  router,
  data () {
    return {
      msg: '2333'
    }
  },
  components: {
    App
  },
  template: '<App/>'
})

/* eslint-disable no-new */
// new Vue({
//   el: '#app',
//   router,
//   components: { App },
//   template: '<App/>'
// })
