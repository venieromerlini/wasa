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
  methods: {
    goToHomePage() {
        router.push("/home")
    },

    async submit() {
      this.loading = true;
      try {
        let response = await this.$axios.post(store.baseUrl +"/session",
            {
              name: this.username
            });

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
      <button @click="submit" type="button" class="btn btn-primary">Sign In</button>
    </form>

  </div>

</template>
<style>
</style>
