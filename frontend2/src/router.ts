import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			name: 'Domu',
			component: () => import('./views/Domu.vue')
		},
		 {
			path: '/o-nas',
			name: 'O nas',
			component: () => import('./views/ONas.vue')
		},
		{
			path: '/lekce',
			name: 'Vsechny lekce',
			component: () => import('./views/SeznamLekci.vue')
		},
        {
			path: '/prihlaseni',
			name: 'Prihlaseni',
			component: () => import('./views/Prihlaseni.vue')
		},
        {
			path: '/registrace',
			name: 'Registrace',
			component: () => import('./views/Registrace.vue')
		},
        {
			path: '/ucet',
			name: 'Ucet',
			component: () => import('./views/Ucet.vue')
		},
		/*
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
			path: '/:pathMatch(.*)*',
			name: '404',
			component: () => import('./views/404.vue')
		} */
	]
})

export default router