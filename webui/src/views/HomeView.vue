<script setup>
import Photo from "@/components/Photo.vue";</script>
<script>
import {store} from "@/services/store";

export default {
  data: function () {
    return {
      loading: false,
      data: {},
      username: store.username,
      baseUrl: store.baseUrl
    }
  },
  methods: {
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
        console.log(this.data)
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
      console.table(this.data)
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

<template>
  <div>

    <h1 class="h2">Welcome, {{ username }}</h1>
    <Photo v-for="(photo, index)  in data['photos']"
           :key="index"
           :showDetails="true"
           :username="photo.user.username"
           :id="photo.id"
           :uploadDate="photo.uploadDate"
           :link="baseUrl+ photo.link"
           :comments="photo.comments"
           :likes="photo.likes"
           @refreshDataRequest="refresh"
    />
    <span>{{ data }}</span>


    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<style>
</style>
