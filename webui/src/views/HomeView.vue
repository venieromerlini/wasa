<script setup>
import Photo from "@/components/Photo.vue";</script>
<script>
import {store} from "@/services/store";
import router from "@/router";

export default {
  data: function () {
    return {
      loading: false,
      data: {},
      dataListOptions: [],
      username: store.username,
      baseUrl: store.baseUrl
    }
  },
  methods: {
    goToUsersPage(){
      router.push("/users"
      );
    },
    async refresh() {
      this.loading = true;
      this.errormsg = null;
      try {
        let username = store.username;
        let authToken = store.authToken
        let response = await this.$axios.get(store.baseUrl + "/home/" + username,
            {
              headers: {[authToken]: username}
            });
        this.data = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;

    },
  },
  mounted() {
    this.refresh()
    var myCollapse = document.getElementById('sidebarMenu')
    var bsCollapse = new bootstrap.Collapse(myCollapse, {
      toggle: false
    })
    bsCollapse.hide()
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
.fullWidthButton {
  width: 100%;
}
.feather{
  width:40px;
  height:40px;
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
      <div class="col-12" >
        <button type="button" class="btn btn-outline-secondary btn-lg btn-block fullWidthButton" @click="goToUsersPage">
          <svg class="feather" stroke="black" fill="white"><use href="/feather-sprite-v4.29.0.svg#zoom-in"/></svg>
        </button>
      </div>
    </div>
    <div class="row">
      <div class="col-12">
        <hr class="solid">
      </div>
    </div>
    <div class="row">
      <div class="col-12">
        <h1 class="h2">Welcome, {{ username }}</h1>
      </div>
    </div>
    <div class="row">
      <div class="col-12">
        <hr class="solid">
      </div>
    </div>
    <div class="row">
      <div class="col-12">
        <Photo v-for="photo in data['photos']" :key="photo.id"
               :show-details="true"
               :username="photo.user.username"
               :id="photo.id"
               :uploadDate="photo.uploadDate"
               :link="baseUrl+ photo.link"
               :comments="photo.comments"
               :likes="photo.likes"
               @refreshDataRequest="refresh"
        />
      </div>
    </div>
  </div>

<!--    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>-->

</template>

