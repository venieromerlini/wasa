<script setup>
import Photo from "@/components/Photo.vue";
import ChangeUsernameModal from "@/components/ChangeUsernameModal.vue";
import FolloweesModal from "@/components/FolloweesModal.vue";
import FollowersModal from "@/components/FollowersModal.vue";
import NewPhotoModal from "@/components/NewPhotoModal.vue";</script>

<script>
import {store} from "@/services/store";

export default {
  data: function () {
    return {
      loading: false,
      data: {},
      photos: [],
      followees: [],
      followers: [],
      bans: [],
      username: store.username,
      baseUrl: store.baseUrl
    }
  },
  methods: {
    async deletePhoto(photoId) {
      try {
        let username = store.username;
        let authToken = store.authToken
        await this.$axios.delete(store.baseUrl + "/photos/" + photoId,
            {
              headers: {[authToken]: username}
            });
      } catch (e) {
        this.errormsg = e.toString();
      }
      await this.refresh()
    },
    async findAllBanned() {
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
    async refresh() {
      this.username = store.username
      this.loading = true;
      this.errormsg = null;
      let username = store.username;
      let authToken = store.authToken
      try {
        let response = await this.$axios.get(store.baseUrl + "/users/" + username,
            {
              headers: {[authToken]: username}
            });
        //this.data = response.data;
        this.photos = response.data ? response.data.photos : [];
        this.followers = response.data ? response.data.followers : [];
        this.followees = response.data ? response.data.followees: [];
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;

      try {
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

</style>
<template>
  <div class="container">
    <div class="row">
      <div class="col-12">
        <hr class="solid">
      </div>
    </div>
    <div class="row">
      <div class="col-8">&nbsp;</div>
      <div class="col-4"><h1 class="h2" id="profileId">{{ username }}</h1></div>
    </div>
    <div class="row">
      <div class="col-8">&nbsp;</div>
      <div class="col-4">
        <button type="button" class="btn btn-outline-secondary border-0 btn-sm" data-bs-toggle="modal"
                data-bs-target="#changeUsernameModal">
          edit username
        </button>
      </div>
    </div>
    <div class="row">
      <div class="col-12">
        <hr class="solid">
      </div>
    </div>
    <div class="row">
      <div class="col-4">
        <p class="text-center">Photos:
          <button type="button" class="btn btn-outline-primary border-0 btn-sm" data-bs-toggle="modal"
                  data-bs-target="#newPhotoModal">
            <b>{{ photos ? photos.length : 0 }}</b>
            <svg class="feather" fill="white" stroke="black">
              <use href="/feather-sprite-v4.29.0.svg#plus-circle"/>
            </svg>
          </button>
        </p>
      </div>
      <div class="col-4">
        <p class="text-center">Followee:
          <button type="button" class="btn btn-outline-primary border-0 btn-sm" data-bs-toggle="modal"
                  data-bs-target="#followeeModal">
            <b>{{ followees ? followees.length : 0 }}</b>
            <svg class="feather" fill="white" stroke="black">
              <use href="/feather-sprite-v4.29.0.svg#zoom-in"/>
            </svg>
          </button>
        </p>
      </div>
      <div class="col-4">
        <p class="text-center">Followers:
          <button type="button" class="btn btn-outline-primary border-0 btn-sm" data-bs-toggle="modal"
                  data-bs-target="#followerModal">
            <b>{{ followers ? followers.length : 0 }}</b>
            <svg class="feather" fill="white" stroke="black">
              <use href="/feather-sprite-v4.29.0.svg#zoom-in"/>
            </svg>
          </button>
        </p>
      </div>
    </div>

    <div class="row" v-for="photo in photos" :key="photo.id">
      <div class="row">
        <div class="col-12">
          <Photo
                 :show-details="false"
                 :username="photo.user.username"
                 :id="photo.id"
                 :uploadDate="photo.uploadDate"
                 :link="baseUrl+ photo.link"
          />
        </div>
      </div>
      <div class="row">
        <div class="col-8">
          &nbsp;
        </div>
        <div class="col-4">
          <button type="button" class="btn btn-primary" @click="deletePhoto(photo.id)">
            <svg class="feather"  fill="white" stroke="black">
              <use href="/feather-sprite-v4.29.0.svg#trash-2" />
            </svg>
            Delete
          </button>

        </div>
        <div class="row">
          <div class="col-12">
            <hr class="solid">
          </div>
        </div>
      </div>
    </div>


  </div>


  <!--  <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#myModal">-->
  <!--    Open modal-->
  <!--  </button>-->
  <!-- The Modal -->
  <ChangeUsernameModal @refreshProfile="refresh"></ChangeUsernameModal>
  <NewPhotoModal @refreshProfile="refresh"></NewPhotoModal>
  <FolloweesModal :followees="followees"
                  @refreshProfile="refresh"></FolloweesModal>
  <FollowersModal :followers="followers"
                  :bans="bans"
                  @refreshProfile="refresh"></FollowersModal>


</template>

<style>
</style>
