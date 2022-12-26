
<script>
import {mutations, store} from "@/services/store";

export default {
  props: ['followees'],
  data: function () {
    return {
      data: {},
      username: store.username,
      baseUrl: store.baseUrl,
      newUsername : ''
    }
  },
  methods: {
    async stopFollow(followId) {

      try {
        let username = store.username;
        let authToken = store.authToken
        let response = await this.$axios.delete(store.baseUrl + "/follows/" + followId,
            {
              headers: {[authToken]: username}
            });
        this.data = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.$emit('refreshProfile', 'VOID')

    }
  },
  mounted() {

  }
}
</script>
<template>
  <!-- The Modal -->
  <div class="modal fade" id="followeeModal">
    <div class="modal-dialog">
      <div class="modal-content">

        <!-- Modal Header -->
        <div class="modal-header">
          <h4 class="modal-title">Followed Users:</h4>
          <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
        </div>

        <!-- Modal body -->
        <div class="modal-body">
           <ul class="list-group list-group-flush">
            <li  v-for="followee in followees" class="list-group-item">{{followee.followee.username}}
              <svg class="feather" fill="white" stroke="black" style="float:right" @click="stopFollow(followee.id)"><use href="/feather-sprite-v4.29.0.svg#user-minus"/></svg>
            </li>
          </ul>
        </div>

        <!-- Modal footer -->
        <div class="modal-footer">
          <button type="button" class="btn btn-danger" data-bs-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>
