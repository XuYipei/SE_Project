<template>
  <div id="home">
    <!-- <img src="./assets/logo.png"> -->
    <div :style ="bg">
    </div>

    <p>user: {{this.disp['user']}}</p>
    <p>userContent: {{this.disp['userContent']}}</p>
    <p>comment: {{this.disp['comment']}}</p>
    <p>commentContent: {{this.disp['commentContent']}}</p>
    <p>problem: {{this.disp['problem']}}</p>
    <p>problemContent: {{this.disp['problemContent']}}</p>
    <button v-on:click="registerBtn">register</button>
    <button v-on:click="loginBtn">login</button>
    <button v-on:click="updateBtn">update</button>
    <router-view/>
  </div>
</template>

<script>
import VueAxios from 'vue-axios';
import axios from 'axios';

export default {
  name: 'Home',
  // created: function () {
  //   this.getMonitor()  
  // },
  data: function () {
    return {
      disp: {'user': 'Loading......'},
      resp: {"data": [{"data": ""}]},
      bg: {
        backgroundImage: "url(" + require("./assets/background.png") + ")",
        backgroundRepeat: "no-repeat",
        backgroundPosition: "center",
        backgroundSize: "100% 100%",
      }
    }
  },
  mounted: function () {
    this.getMonitor()
    window.setInterval(() => {
      setTimeout(this.getMonitor(), 0);
    }, 6000000);
  },
  methods: {
    getProblems () {
      this.resp = {"data": [{"data": ""}]}
      this.disp = 'Loading......'
      axios
        .get('http://localhost:5000/api/problems')
        .then(response => (this.resp = response))
        .catch(function (error) { // 请求失败处理1
          console.log(error);
        });
    },
    getMonitor: function() {
      this.disp = 'Loading......'
      // this.resp = {"data": ""}
      axios
        .get('http://localhost:5000/api/monitor')
        .then(response => (this.disp = response.data["data"]))
        .catch(function (error) { // 请求失败处理1
          console.log(error);
        });
      // this.disp = this.resp
    },
    registerBtn: function() {
      this.resp = {"data": [{"data": ""}]}
      axios
      .post('http://localhost:5000/api/user/register', {
        username: "tester2",
        password: "123"
      })
      .then(response => (this.resp = response))
      .catch(function (error) { // 请求失败处理
        console.log(error);
      });
    },
    loginBtn: function() {
      this.resp = {"data": [{"data": ""}]}
      axios
      .post('http://localhost:5000/api/user/login', {
        username: "tester2",
        password: "123"
      })
      .then(response => (this.resp = response))
      .catch(function (error) { // 请求失败处理
        console.log(error);
      });
    },
    updateBtn: function() {
      this.resp = {"data": [{"data": ""}]}
      axios
      .post('http://localhost:5000/api/user/update', {
        id: "c81e728d9d4c2f636f067f89cc14862c",
        gender: true,
        description: 'a system tester',
        username: 'tester2'
      },{
        headers: {
          'Authorization': 'c81e728d9d4c2f636f067f89cc14862c'
        }
      })
      .then(response => (this.resp = response))
      .catch(function (error) { // 请求失败处理
        console.log(error);
      });
    },
  }
}
</script>

<style>
#home {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
