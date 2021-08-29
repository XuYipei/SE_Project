<template>
  <!-- <div id="app" class="note" :style="note"></div> -->
  <!-- <img src="./assets/logo.png"> -->
  <div id="app">
    <div id="deviceInfo">
      <a-table key="deviceTable" 
        :columns="deviceColums"
        :dataSource="deviceData"
        :rowKey="(record,index)=>{return index}">
        <span slot="name">设备名</span>
        <span slot="usage">使用率</span>
      </a-table>
    </div>
    <div id="databaseInfo">
      <a-table key="databaseTable" 
        :columns="databaseColums"
        :dataSource="databaseData">
        <span slot="name">表名</span>
        <span slot="count">数据条目</span>
        <span slot="size">已用空间</span>
      </a-table>
    </div>
    <div id="databaseOpt">
      <button v-on:click="synclockBtn">synclock</button>
      <button v-on:click="dumpBtn">dump</button>
      <button v-on:click="storeBtn">store</button>
      <button v-on:click="unlockBtn">unlock</button>
    </div>
    <router-view/>
  </div>
</template>

<style>
  /* #app{
    background:url("./assets/background.png");
    width:100%;	
    height:100%;
    position:fixed;
    background-size:100% 100%;
    margin: 0;
    padding: 0;
  } */
  #deviceInfo{
    /* background:url("./assets/logo.png"); */
    margin-top: 10%;
    margin-left: 5%;
    position: center center;
  }
  #databaseInfo{
    margin-top: 10%;
    margin-left: 5%;
    position: center center;
  }
</style>

<script>
  import VueAxios from 'vue-axios';
  import axios from 'axios'

  export default {
    name: 'App',
    // created: function () {
    //   this.getMonitor()  
    // },
    data: function () {
      return {
        disp: {'user': 'Loading......'},
        resp: {"data": [{"data": ""}]},
        
        databaseColums: [
          {　
            dataIndex: 'name',
            align: 'center',
            slots: { title: 'name'}
          },
          {　
            dataIndex: 'count',
            align: 'center',
            slots: { title: 'count'}
          },
          {　
            dataIndex: 'size',
            align: 'center',
            slots: { title: 'size'}
          },
        ], 
        databaseData: [{
          "name": "user",
          "count": 4,
          "size": 716.0,
          "key": "0",
        }],

        deviceColums: [
          {　
            dataIndex: 'name',
            align: 'center',
            slots: { title: 'name'}
          },
          {　
            dataIndex: 'usage',
            align: 'center',
            slots: { title: 'usage'}
          },
        ],
        deviceData: [
          {
            "name": "CPU",
            "usage": 0,
          },
          {
            "name": "memory",
            "usage": 0,
          },
        ]
      }
    },
    mounted: function () {
      this.getMonitor()
      window.setInterval(() => {
        setTimeout(this.getDatabaseMonitor(), 0);
      }, 6000000);
    },
    methods: {
      getDatabaseMonitor: function() {
        axios
          .get('http://localhost:5000/api/monitor/database')
          .then(response => (this.databaseData = response.data["data"]))
          .catch(function (error) {console.log(error);});
      },
      getDeviceMonitor: function() {
        axios
          .get('http://localhost:5000/api/monitor/device')
          .then(response => (this.deviceData = response.data["data"]))
          .catch(function (error) {console.log(error);});
      },
      getMonitor: function() {
        this.getDatabaseMonitor()
        this.getDeviceMonitor()
      },
      synclockBtn: function() {
        axios
        .post('http://localhost:5000/api/monitor/synclock', {})
        .then(response => (this.resp = response))
        .catch(function (error) {console.log(error);});
      },
      unlockBtn: function() {
        axios
        .post('http://localhost:5000/api/monitor/unlock', {})
        .then(response => (this.resp = response))
        .catch(function (error) {console.log(error);});
      },
      dumpBtn: function() {
        axios
        .post('http://localhost:5000/api/monitor/dump', {})
        .then(response => (this.resp = response))
        .catch(function (error) {console.log(error);});
      },
      storeBtn: function() {
        axios
        .post('http://localhost:5000/api/monitor/store', {})
        .then(response => (this.resp = response))
        .catch(function (error) {console.log(error);});
      },
    }
  }
</script>
