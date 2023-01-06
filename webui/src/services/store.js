import { reactive } from 'vue'

export const store = reactive({
    baseUrl: __API_URL__,
    isNavOpen: false,
    isUserLogged : false,
    token: '',
    username: '',
    authToken: 'X-User-Session-Identifier',
    toggleCommentBoxes : []

});

export const mutations = {
    toggleNav() {
        store.isNavOpen = !store.isNavOpen
    },

    setUserData(data) {
        store.token = data
        store.username = data
        store.isUserLogged = true
    },
    toggleCommentBox(id) {
        if(store.toggleCommentBoxes.includes(id)){
            const index = store.toggleCommentBoxes.indexOf(id);
            if (index > -1) {
                store.toggleCommentBoxes.splice(index, 1);
            }
        }else{
            store.toggleCommentBoxes.push(id)
        }

    }
};