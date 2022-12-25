<script>
import {mutations, store} from "@/services/store";

export default {
  props: ['username', 'id', 'uploadDate', 'link', 'comments', 'likes'],
  data: function () {
    return {
      toggleCommentBoxes: [],
      commentText: ''

    }
  },
  methods: {
    getUserNameString(username) {
      if (username === store.username){
        return "You"
      }
      return username.charAt(0).toUpperCase() + username.slice(1);
    },
    userAlreadyLikedThisPhoto() {
      return this.likes && this.likes.map(x => x.user.username).includes(store.username)
    },
    commentX() {
      return this.comments && this.comments.length > 9 ? 8 : 13
    },
    findUserLikeId() {
      return this.likes.find(x => {
        return x.user.username === store.username
      }).id
    },
    firstUserLike() {
      return this.likes ? this.likes[0].user.username : ''
    },
    likeColorFill() {
      return this.userAlreadyLikedThisPhoto() ? 'red' : 'white'
    },
    likeColorText() {
      return this.userAlreadyLikedThisPhoto() ? 'white' : 'black'
    },
    likeX() {
      return this.likes && this.likes.length > 9 ? 8 : 13
    },
    toggleLike() {
      const config = {
        headers: {
          [store.authToken]: store.token
        }
      };
      if (this.userAlreadyLikedThisPhoto()) {
        try {
          this.$axios.delete(store.baseUrl + "/likes/" + this.findUserLikeId(), config)
        } catch (e) {
          this.errormsg = e.toString();
        }
      } else {
        try {
          this.$axios.post(store.baseUrl + "/likes",
              {
                user: {username: store.username},
                photoId: this.id
              },
              config)
        } catch (e) {
          this.errormsg = e.toString();
        }
      }
      this.$emit('refreshDataRequest', 'VOID')
    },
    openCommentBox(id) {
      mutations.toggleCommentBox(id)
      this.toggleCommentBoxes = store.toggleCommentBoxes
    },
    isCommentBoxOpen(id) {
      return store.toggleCommentBoxes.includes(id)
    },
    addComment() {
      try {
        const config = {
          headers: {
            [store.authToken]: store.token
          }
        };
        this.$axios.post(store.baseUrl + "/comments",
            {
              user: {username: store.username},
              text: this.commentText,
              photoId: this.id
            },
            config)
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.$emit('refreshDataRequest', 'VOID')

    }
  },
  mounted() {

  }
}
</script>
<style scoped>
h5, h6 {
  display: inline;
}

.feather {
  width: 34px;
  height: 34px;
}


</style>
<template>
  <div class="card">
    <div class="card-header">
      <i>{{ username }}</i>
    </div>
    <img class="card-img-top" :src="link" alt="Card image cap">

    <div class="card-body">
      <svg class="feather" stroke="black" :fill="likeColorFill()" @click="toggleLike()">
        <use href="/feather-sprite-v4.29.0.svg#heart"/>
        <text :x="likeX()" y="21" :stroke="likeColorText()" :fill="likeColorText()">{{ likes ? likes.length : 0 }}
        </text>
      </svg>
      <svg class="feather" stroke="black" fill="white" @click="openCommentBox(id)">
        <use href="/feather-sprite-v4.29.0.svg#message-circle"/>
        <text :x="commentX()" y="21" stroke="black" fill="black">{{ comments ? comments.length : 0 }}</text>
      </svg>
      <p v-if="likes" class="card-text">
        {{
          (likes && likes.length > 1) ? firstUserLike() + ' and others like this photo' : firstUserLike() + ' likes this photo'
        }} </p>
      <p v-if="!likes"></p>

      <div v-if="isCommentBoxOpen(id)">
        <div class="input-group input-group-sm mb-3">
          <input v-model="commentText" type="text" class="form-control" placeholder="Add a comment..."
                 aria-label="Add a comment..."
                 aria-describedby="basic-addon2">
          <div class="input-group-append input-group-append-sm">
            <button class="btn btn-outline-secondary " type="button" @click="addComment()">
              Post
            </button>
          </div>
        </div>

        <div class="card-body" v-for="comment in comments">
          <h5 class="card-title">{{ getUserNameString(comment.user.username) }} </h5>
          <button type="button" class="btn-close" aria-label="Close"></button>
          <p class="card-text">{{ comment.text }}</p>
        </div>
      </div>
    </div>

  </div>
</template>
