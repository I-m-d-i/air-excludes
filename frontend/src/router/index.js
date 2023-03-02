import Vue from 'vue'
import VueRouter from 'vue-router'
import Excludes from '../components/Excludes.vue'
import Login from '../components/Login.vue'
import RegistrationPage from '../components/Registration.vue'
import axios from "axios";

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'home',
        component: Excludes,
        meta: {
            requiresAuth: true
        },
    },
    {
        path: '/login',
        name: 'login',
        component: () => import('../components/Login.vue'),
        meta: {
            requiresAuth: false
        },
    },
    {
        path: "/registration",
        name: "registration",
        component:()=> import('../components/Registration.vue'),
        meta: {
            requiresAuth: false
        },
    },
]

const router = new VueRouter({
    routes: routes
})

router.beforeEach((to, from, next) => {
    if (to.name === "login" || to.name === "registration") {
        axios
            .post("/api/auth/authRequire", {})
            .then(response => {
                router.replace({path: '/',
                    name: 'home'})
                return
            })
            .catch(e => {
            });
    }
    if (to.matched.some(record => record.meta.requiresAuth)) {
        axios
            .post("/api/auth/authRequire", {})
            .then(response => {
                next()
            })
            .catch(e => {
                if (e.response.status === 401) {
                    next({
                        path: '/login',
                        name: 'login',
                        params: {nextUrl: from}
                    })
                }
            });
    } else {
        next()
    }
})

export default router
