
<script>
import {mutations, store} from "@/services/store";

export default {
  props: ['followers', 'bans'],
  data: function () {
    return {
      username: store.username,
      baseUrl: store.baseUrl,
      newUsername : ''
    }
  },
  methods: {
    async unban(bannedUser) {
      try {
        let username = store.username;
        let authToken = store.authToken
        let response = await this.$axios.delete(store.baseUrl + "/bans/" + this.findBanByBannerUsername(bannedUser).id,
            {
              headers: {[authToken]: username}
            });
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.$emit('refreshProfile', 'VOID')
    },
    async ban(bannedUsername) {
      try {
        let username = store.username;
        let authToken = store.authToken
        let response = await this.$axios.post(store.baseUrl + "/bans",
            {
                  user: {
                    username: store.username
                  },
                  banned: {
                    username: bannedUsername
                  }
            },
            {
              headers: {[authToken]: username}
            });
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.$emit('refreshProfile', 'VOID')
    },
    alsoBanned(username){
      return this.bans.find(e => e.banned.username === username)
    },
    findBanByBannerUsername(username){
      return this.bans.find(e => e.banned.username === username)
    }
  },
  mounted() {

  }
}
</script>
<template>
  <!-- The Modal -->
  <div class="modal fade" id="followerModal">
    <div class="modal-dialog">
      <div class="modal-content">

        <!-- Modal Header -->
        <div class="modal-header">
          <h4 class="modal-title">Your Followers:</h4>
          <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
        </div>

        <!-- Modal body -->
        <div class="modal-body">
           <ul class="list-group list-group-flush">
            <li  v-for="follower in followers" class="list-group-item">{{follower.user.username}}
              <svg v-if="!alsoBanned(follower.user.username)" class="feather" fill="white" stroke="green" style="float:right" @click="ban(follower.user.username)">
                <use href="/feather-sprite-v4.29.0.svg#unlock"/>
              </svg>
              <svg v-if="alsoBanned(follower.user.username)" class="feather" fill="white" stroke="red" style="float:right" @click="unban(follower.user.username)">
                <use href="/feather-sprite-v4.29.0.svg#lock"/>
              </svg>
            </li>
<!--             <li  v-for="ban in this.bans" class="list-group-item">{{ban.banned.username}}-->
<!--               <svg class="feather" fill="white" stroke="red" style="float:right" @click="unban(ban.id)">-->
<!--                 <use href="/feather-sprite-v4.29.0.svg#lock"/>-->
<!--               </svg>-->
<!--             </li>-->
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
