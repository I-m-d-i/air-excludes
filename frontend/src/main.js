import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import Axios from 'axios';
import router from './router'
export const bus = new Vue()

//Axios.defaults.baseURL = "http://127.0.0.1:8080/";
Vue.config.productionTip = true;
Axios.defaults.withCredentials = true;

new Vue({
  bus,
  vuetify,
  router,
  render: h => h(App)
}).$mount('#app')
