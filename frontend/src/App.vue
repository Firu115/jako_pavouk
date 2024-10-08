<script setup lang="ts">
import { onMounted, ref } from "vue";
import MenuLink from "./components/MenuLink.vue";
import { mobil, prihlasen, role, tokenJmeno } from "./stores";
import { jeToRobot, getToken, oznameni, pridatOznameni } from "./utils";
import { useHead } from "unhead";
import axios from "axios";
import { useRouter } from "vue-router";

useHead({
    titleTemplate: (title?: string) => title == "" || title == undefined ? "Psaní všemi deseti zdarma | Jako Pavouk" : `${title} | Jako Pavouk`
})


const router = useRouter()
const mobilMenu = ref(false)

onMounted(() => {
    console.log("%cCo sem koukáš koloušku?", "color: white; font-size: x-large"); // troulin

    if (getToken()) {
        axios.get("/token-expirace", {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        }).then(response => {
            role.value = response.data.role

            if (response.data.jePotrebaVymenit) {
                localStorage.removeItem(tokenJmeno)
                prihlasen.value = false
                router.push("/prihlaseni")
                pridatOznameni("Z bezpečnostních důvodů jsme tě odhlásili ze sítě 🕸️", 8000)
            } else {
                prihlasen.value = true
            }
        }).catch(e => {
            if (!(e.response && e.response.status == 418)) {
                console.log(e)
                pridatOznameni("Chyba serveru")
            }
        })
    } else if (!jeToRobot(navigator.userAgent) && window.location.host !== "localhost:5173") { //test jestli to neni bot + počítají se jen na produkci
        axios.post("/navsteva")
    }

    window.addEventListener("resize", function () {
        mobil.value = document.body.clientWidth <= 900
    })
})

</script>

<template>
    <header>
        <div id="menu-mobilni-btn" @click="mobilMenu = !mobilMenu"><img id="menuIcon" src="./assets/icony/menu.svg" alt="Menu" width="40" height="40">
        </div>
        <nav :class="{ 'mobil-hidden': !mobilMenu }" @click="mobilMenu = !mobilMenu">
            <MenuLink jmeno="Domů" cesta="/" />
            <MenuLink jmeno="Jak psát" cesta="/jak-psat" />
            <MenuLink jmeno="Lekce" cesta="/lekce" />
            <MenuLink jmeno="Procvičování" cesta="/procvic" />
            <MenuLink jmeno="Test psaní" cesta="/test-psani" />
            <MenuLink v-if="role == 'student'" jmeno="Škola" cesta="/trida" />
            <MenuLink v-else-if="role == 'ucitel'" jmeno="Škola" cesta="/skola" />
            <MenuLink jmeno="O nás" cesta="/o-nas" />
            <MenuLink v-if="!prihlasen" jmeno="Přihlásit se" cesta="/prihlaseni" />
            <MenuLink v-else jmeno="Můj účet" cesta="/ucet" />
        </nav>
    </header>
    <div id="view">
        <RouterView :key="$route.fullPath" />
    </div>

    <div id="alerty">
        <TransitionGroup name="list">
            <div v-for="(o, i) in oznameni" class="alert" :key="i">
                <img v-if="o.typ == 'vykricnik'" src="./assets/icony/alert.svg" alt="Vykřičník">
                <img v-else-if="o.typ == 'copy'" src="./assets/icony/copy.svg" alt="Zkopírováno">
                <span v-html="o.text"></span>
            </div>
        </TransitionGroup>
    </div>
</template>

<style scoped>
/* na tu animaci oznameni */
.list-move {
    transition: all 0.2s ease;
}

.list-enter-active,
.list-leave-active {
    transition: all 0.1s ease;
}

.list-enter-from,
.list-leave-to {
    opacity: 0;
    transform: translateX(50px);
}

.list-leave-active {
    position: absolute;
}

#alerty {
    position: fixed;
    right: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    justify-content: end;
    gap: 10px;
    padding: 20px;
    min-height: 100px;
    pointer-events: none;
    width: 100vw;
    z-index: 1000;
}

.alert {
    min-height: 60px;
    background-color: var(--tmave-fialova);
    min-width: 100px;
    max-width: min(85%, 330px);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 5px;
    padding: 10px 20px 10px 20px;
    gap: 15px;
    box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
}

.alert img {
    width: 24px;
}

nav {
    position: fixed;
    left: 10px;
    top: 10px;
    width: var(--sirka-menu);
    height: calc(100vh - 20px);
    flex-shrink: 0;
    border-radius: 10px;
    background-color: var(--tmave-fialova);
    transition: transform ease-in-out 0.3s;
    overflow: hidden;
}

#menu-mobilni-btn {
    display: none;
}

@media screen and (max-width: 1100px) {
    .mobil-hidden {
        transform: translateX(-250px);
        transition: transform ease-in-out 0.3s;
    }

    #menu-mobilni-btn {
        background-color: var(--tmave-fialova);
        border-radius: 100px;
        padding: 10px;
        display: block;
        position: fixed;
        right: 10px;
        top: 10px;
        width: 60px;
        height: 60px;
        box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
        z-index: 1000;
    }

    nav {
        border-radius: 10px;
        background-color: var(--tmave-fialova);
        display: flex;
        flex-direction: column;
        z-index: 10;
        box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
    }

    #view {
        padding-top: 30px;
        margin-left: 0;
        margin-bottom: 50px;
        text-align: center;
        width: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }
}
</style>
