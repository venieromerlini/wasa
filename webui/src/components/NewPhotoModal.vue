
<script>
import {mutations, store} from "@/services/store";

export default {
  props: ['followers', 'bans'],
  data: function () {
    return {
      username: store.username,
      baseUrl: store.baseUrl,
      file : '',
      uploadDoneOk: false
    }
  },
  methods: {
    handleFileUpload(){
      this.file =  this.$refs.file.files[0];
    },
    async uploadPhoto() {
      try {
        let username = store.username;
        let authToken = store.authToken

        const formData = new FormData();
        formData.append('file', this.file);
        formData.append("username",  store.username);
        let response = await this.$axios.post("/photos", formData,
            {
              headers: {
                [authToken]: username,
                'Content-Type': 'multipart/form-data'
              }
            });
        this.uploadDoneOk = true
      } catch (e) {
        console.error(e)
        this.errormsg = e.toString();
      }
      this.$emit('refreshProfile', 'VOID')
    },
  },
  mounted() {

  }
}
</script>
<template>
  <!-- The Modal -->
  <div class="modal fade" id="newPhotoModal">
    <div class="modal-dialog">
      <div class="modal-content">

        <!-- Modal Header -->
        <div class="modal-header">
          <h4 class="modal-title">Upload a Photo</h4>
          <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
        </div>

        <!-- Modal body -->
        <div class="modal-body">
          <div class="mb-3">
            <label for="file" class="form-label"></label>
            <input class="form-control" type="file" id="file"  ref="file" @change="handleFileUpload()">
          </div>
        </div>
        <!-- Modal footer -->
        <div class="modal-footer">
          <button v-if="!uploadDoneOk" type="submit" class="btn btn-primary" @click="uploadPhoto">Upload</button>
          <button v-if="uploadDoneOk" type="submit" class="btn success" disabled>Done!</button>

          <button type="button" class="btn btn-danger" data-bs-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>
</template>
