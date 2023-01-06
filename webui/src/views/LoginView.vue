<script setup>
import LoadingSpinner from "@/components/LoadingSpinner.vue";
</script>

<script>
import {mutations, store} from '@/services/store.js'
import router from "@/router";
export default {
  data: function () {
    return {
      errormsg: null,
      loading: false,
      data: null,
      username: ''
    }
  },
  computed: {

  },
  methods: {
    goToHomePage() {
        router.push("/home")
    },
    isDisabled() {
      return this.username.length < 3;
    },
    async submit() {
      this.loading = true;
      try {
        //let username = store.username;
        let authToken = store.authToken
        let response = await this.$axios.post(store.baseUrl +"/session",
            {
              name: this.username
            },
            {
              headers: {
                [authToken]: 'not-yet',
                'Content-Type': 'application/json'
              }
            }
            );

        mutations.setUserData(response.data.identifier);

      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false
      this.goToHomePage()
    }
  },
  mounted() {

  }
}
</script>

<template>
  <LoadingSpinner  :loading="loading"/>
  <div class="vue-tempalte" v-if="!data">

    <h3>WASA</h3>
    <form>

      <div class="form-group">
        <label>Username</label>
        <input type="text" class="form-control form-control-lg" v-model="username"/>
      </div>
      <button @click="submit"
              :disabled="isDisabled(username)"
              type="button" class="btn btn-primary">Sign In</button>
    </form>

  </div>

</template>
<style>
</style>
