<script setup>
import {RouterLink, RouterView} from 'vue-router'

</script>

<script>
import router from './router'
import {mutations, store} from '@/services/store.js'

export default {

  data() {
    return {
      isUserLogged: store.isUserLogged
    }
  },
  methods:{
    refresh() {
      this.isUserLogged = store.isUserLogged
    }
  },
  beforeMount() {
    if (!store.isUserLogged){
      router.push("/login")
    }

  }

}
</script>


<template>

    <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
      <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASA APP</a>
      <button class="navbar-toggler position-absolute d-md-none collapsed" @click="refresh" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
    </header>




    <div class="container-fluid">
      <div class="row">
        <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
          <div class="position-sticky pt-3 sidebar-sticky">
            <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
              <span>settings</span>
            </h6>
            <ul class="nav flex-column">
              <li class="nav-item" v-if="!isUserLogged">
                <RouterLink to="/login" class="nav-link">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                  Login
                </RouterLink>
              </li>

              <li class="nav-item" v-if="isUserLogged">
                <RouterLink to="/home" class="nav-link">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                  Home
                </RouterLink>
              </li>
              <li class="nav-item" v-if="isUserLogged">
                <RouterLink to="/profile" class="nav-link">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
                  Profile
                </RouterLink>
              </li>

            </ul>
          </div>
        </nav>

        <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
          <RouterView />
        </main>
      </div>
    </div>

</template>

<style>
</style>
