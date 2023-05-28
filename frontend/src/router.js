import {createRouter, createWebHistory} from 'vue-router'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomeView
        },
        {
            path: '/about',
            name: 'about',
            component: () => import('./views/AboutView.vue')
        },
        {
            path: '/lekce',
            name: 'VsechnyLekce',
            component: () => import('./views/VsechnyLekce.vue')
        },
        {
            path: '/lekce/:pismena',
            name: 'lekce',
            component: () => import('./views/LekceView.vue')
        },
        {
            path: '/lekce/:pismena/:id',
            name: 'cviceni',
            component: () => import('./views/CviceniView.vue')
        },
        {
            path: '/login',
            name: 'login',
            component: () => import('./views/LoginView.vue')
        },
        {
            path: '/register',
            name: 'register',
            component: () => import('./views/RegisterView.vue')
        },
        {
            path: '/ucet',
            name: 'ucet',
            component: () => import('./views/UcetView.vue')
        },
        {
            path: '/:pathMatch(.*)*',
            name: '404',
            component: () => import('./views/404.vue')
        }
    ]
})

export default router
