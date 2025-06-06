import { createRouter, createWebHistory } from "vue-router"
import { prihlasen } from "./stores";

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: "/",
            component: () => import("./views/Domu.vue"),
        },
        {
            path: "/o-nas",
            component: () => import("./views/ONas.vue")
        },
        {
            path: "/kurz",
            component: () => import("./views/SeznamLekci.vue")
        },
        {
            path: "/procvic",
            component: () => import("./views/SeznamProcvicovani.vue")
        },
        {
            path: "/jak-psat",
            component: () => import("./views/Teorie.vue"),
        },
        {
            path: "/skolni-system",
            component: () => import("./views/InfoSkola.vue"),
        },
        {
            path: "/skola",
            component: () => import("./views/ucitel/SeznamTrid.vue"),
        },
        {
            path: "/skola/:id",
            component: () => import("./views/ucitel/Trida.vue"),
        },
        {
            path: "/trida",
            component: () => import("./views/student/Trida.vue"),
        },
        {
            path: "/prace/:id",
            component: () => import("./views/student/Prace.vue"),
        },
        {
            path: "/zapis/:kod",
            component: () => import("./views/student/Zapis.vue"),
        },
        {
            path: "/zapis/",
            component: () => import("./views/student/Zapis.vue"),
        },
        {
            path: "/prihlaseni",
            component: () => import("./views/Prihlaseni.vue")
        },
        {
            path: "/registrace",
            component: () => import("./views/Registrace.vue")
        },
        {
            path: "/zapomenute-heslo",
            component: () => import("./views/ZapomenuteHeslo.vue")
        },
        {
            path: "/nastaveni",
            component: () => import("./views/Nastaveni.vue"),
        },
        {
            path: "/statistiky",
            component: () => import("./views/Statistiky.vue"),
        },
        {
            path: "/klavesnice",
            component: () => import("./views/VyberKlavesnice.vue"),
        },
        {
            path: "/lekce/:pismena",
            component: () => import("./views/Lekce.vue"),
        },
        {
            path: "/lekce/:pismena/:id",
            component: () => import("./views/Cviceni.vue"),
        },
        {
            path: "/procvic/:id",
            component: () => import("./views/Procvic.vue"),
        },
        {
            path: "/test-psani",
            component: () => import("./views/TestPsani.vue"),
        },
        {
            path: "/prvni-psani",
            component: () => import("./views/PrvniPsani.vue"),
        },
        {
            path: "/schvalovani-textu",
            component: () => import("./views/SchvalovaniTextu.vue"),
        },
        {
            path: "/:pathMatch(.*)*",
            component: () => import("./views/404.vue")
        }
    ],
    scrollBehavior(_, __, savedPos) {
        if (savedPos) return savedPos
        return { top: 0, behavior: "smooth" } // aby scroll nezustaval dole na strankach kde se nescrolluje
    },
})

router.beforeEach((to, _, next) => { // kdyz potrebuje auth tak => prihlaseni
    if (to.meta.requireAuth) {
        if (!prihlasen) {
            next("/prihlaseni")
        } else {
            to.fullPath = to.fullPath.toLocaleLowerCase()
            to.path = to.path.toLocaleLowerCase()
            next()
        }
    } else {
        to.fullPath = to.fullPath.toLocaleLowerCase()
        to.path = to.path.toLocaleLowerCase()
        next()
    }
})

export default router