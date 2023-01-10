
<script>
import {mutations, store} from "@/services/store";

export default {
  data: function () {
    return {
      data: {},
      username: store.username,
      baseUrl: store.baseUrl,
      newUsername : '',
      change : 0,
    }
  },
  methods: {
    async changeUsername() {

      try {

        let username = store.username;
        let authToken = store.authToken
        let response = await this.$axios.put(store.baseUrl + "/users/" + username,
            {
              username: this.newUsername
            },
            {
              headers: {[authToken]: username}
            });
        this.data = response.data;
        mutations.setUserData(this.data.username)
        this.change = 1
      } catch (e) {
        this.change = -1
        console.log(e)
        this.errormsg = e.toString();
      }
      this.$emit('refreshProfile', 'VOID')

    },
    resetAlert(){
      this.change = 0
    }
  },
  mounted() {
  }
}
</script>
<template>
  <!-- The Modal -->
  <div class="modal fade" id="changeUsernameModal">
    <div class="modal-dialog">
      <div class="modal-content">

        <!-- Modal Header -->
<!--        <div class="modal-header">-->
<!--          <h4 class="modal-title">Modal Heading</h4>-->
<!--          <button type="button" class="btn-close" data-bs-dismiss="modal"></button>-->
<!--        </div>-->

        <!-- Modal body -->
        <div class="modal-body">
          <form>
            <div class="mb-3">
              <label for="usernameId" class="form-label">New Username</label>
              <input type="text" v-model="newUsername" class="form-control" id="usernameId" >
<!--              <div id="usernameHelp" class="form-text">Insert the new username</div>-->
            </div>
            <div v-if="change>0" class="alert alert-success" role="alert">
              Username changed successfully
            </div>
            <div v-if="change<0" class="alert alert-danger" role="alert">
              Couldn't change username!
            </div>
          </form>
        </div>

        <!-- Modal footer -->
        <div class="modal-footer">
          <button type="submit" class="btn btn-primary" @click="changeUsername">Change</button>
          <button type="button" class="btn btn-danger" data-bs-dismiss="modal" @click="resetAlert">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>
