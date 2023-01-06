<script setup></script>

<script>
import {store} from "@/services/store";
import router from "@/router";

export default {
  data: function () {
    return {
      loading: false,
      users: [],
      bans: [],
      followees: [],
      username: store.username,
      baseUrl: store.baseUrl
    }
  },
  methods: {
    back(){
      router.push("/home")
    },
    async findAllUsers(){
      try {
        let username = store.username;
        let authToken = store.authToken
        let response = await this.$axios.get(store.baseUrl + "/users",
            {
              headers: {[authToken]: username}
            });
        this.users = response.data.filter((e) => e.username != store.username);
      } catch (e) {
        this.errormsg = e.toString();
      }
    },
    async findAllBanned(){
      try {
        let username = store.username;
        let authToken = store.authToken
        let response = await this.$axios.get(store.baseUrl + "/bans?username=" + username,
            {
              headers: {[authToken]: username}
            });
        this.bans = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.$emit('refreshProfile', 'VOID')
    },
    async findAllFollowee() {
      try {
        let username = store.username;
        let authToken = store.authToken
        let response = await this.$axios.get(store.baseUrl + "/home/" + username,
            {
              headers: {[authToken]: username}
            });
        this.followees = response.data.followees;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    isFollowee(username){
        return this.followees.find(e => e.followee.username === username)
    },
    async follow(usernameToFollow) {
      try {
        let username = store.username;
        let authToken = store.authToken
        await this.$axios.post(store.baseUrl + "/follows",
            {
              user: {
                username: store.username
              },
              followee: {
                username: usernameToFollow
              }
            },
            {
              headers: {[authToken]: username}
            });
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.refresh()
      //this.$emit('refreshProfile', 'VOID')
    },
    refresh(){
      this.findAllUsers()
      this.findAllBanned()
      this.findAllFollowee()
    },



  },
  mounted() {
    this.refresh()
  }
}
</script>
<style scoped>
hr.hr-text::before {
  content: attr(data-content);
  display: inline-block;
  background: #fff;
  font-weight: bold;
  font-size: 0.85rem;
  color: #999;
  border-radius: 30rem;
  padding: 0.2rem 2rem;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
hr.solid {
  border-top: 2px solid #999;
}

</style>
<template>
  <div class="container">
    <div class="row">
      <div class="col-12">
        <hr class="solid">
      </div>
    </div>
    <div class="row">
      <div class="col-12">
        <h1 class="h2">{{ username }}, add a friend!</h1>
      </div>
    </div>
    <div class="row">
      <div class="col-12">
        <ul class="list-group">
          <li v-for="user in users" class="list-group-item" :key="user.username">{{user.username}}
            <svg v-if="!isFollowee(user.username)" @click="follow(user.username)" class="feather float-right" fill="white" stroke="black" ><use href="/feather-sprite-v4.29.0.svg#user-plus"/></svg>
            <svg v-if="isFollowee(user.username)" class="feather float-right" fill="white" stroke="green" ><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
          </li>
        </ul>
        &nbsp;</div>
    </div>
  </div>



</template>

<style scoped>
.float-right{
  float: right;
}
</style>
