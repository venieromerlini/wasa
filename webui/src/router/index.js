import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from "@/views/LoginView.vue";
import ProfileView from "@/views/ProfileView.vue";
import UsersView from "@/views/UsersView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/login', component: LoginView},
		{path: '/home', component: HomeView},
		{path: '/profile', component: ProfileView},
		{path: '/users', component: UsersView},

		/*{path: '/link2', component: HomeView},
		{path: '/some/:id/link', component: HomeView},
		*/
	],


})

export default router
