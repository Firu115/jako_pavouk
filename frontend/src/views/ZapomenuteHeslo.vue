<script setup lang="ts">
import axios from "axios";
import { onMounted, ref } from "vue";
import { pridatOznameni } from "../utils";
import { useHead } from "unhead";
import { useRoute, useRouter } from "vue-router";

useHead({
    title: "Zapomenuté heslo"
})

const router = useRouter()

const predvyplenyEmail = useRoute().query.email as string

const heslo = ref("")
const kod = ref("")
const email = ref("")
const spatnyEmail = ref(false)
const spatnyKod = ref(false)
const spatnyHeslo = ref(false)

const posilame = ref(false)

const state = ref("email")

onMounted(() => {
    if (predvyplenyEmail != undefined) {
        email.value = predvyplenyEmail
    }
})

function zmena() { // pokud zacnu znova psat tak zrusim znaceni spatnyho inputu
    spatnyEmail.value = false
    spatnyKod.value = false
    spatnyHeslo.value = false
}

function poslatEmail(e: Event) {
    e.preventDefault(); //aby se nerefreshla stranka
    if (spatnyEmail.value) return

    posilame.value = true
    axios.post("/zmena-hesla", {
        "email": email.value,
    }).then(() => {
        state.value = "kod"
        posilame.value = false
    }).catch(e => {
        if (e.response.data.error.toLowerCase().search("email") != -1) {
            spatnyEmail.value = true
            pridatOznameni("Tenhle email ještě neznáme")
            return
        }
        pridatOznameni()
        console.log(e)
    })
}

function overitZmenu(e: Event) {
    e.preventDefault(); //aby se nerefreshla stranka
    if (spatnyHeslo.value) {
        pridatOznameni("Heslo musí být alespoň 5 znaků. Toť vše")
        return
    }
    if (spatnyHeslo.value || spatnyKod.value) return

    axios.post("/overeni-zmeny-hesla", {
        "email": email.value,
        "heslo": heslo.value,
        "kod": kod.value
    }).then(() => {
        state.value = "konec"
    }).catch(e => {
        console.log(e)
        if (e.response.data.error.toLowerCase().search("kod") != -1) spatnyKod.value = true
        else if (e.response.data.error.toLowerCase().search("heslo") != -1) spatnyHeslo.value = true
        pridatOznameni()
    })
}

function chekujUdaje(jaky: string) {
    if (jaky === "email" && email.value) spatnyEmail.value = !/^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/g.test(email.value) //test jestli email
    else if (jaky === "heslo" && heslo.value !== undefined) spatnyHeslo.value = !/^(?=.*[a-zA-Z]).{5,128}$/.test(heslo.value)
    else if (jaky === "kod" && kod.value !== undefined) spatnyKod.value = !/^\d{5}$/.test(kod.value) //kod 5 dlouhy
}

function presmerovat(e: Event) {
    e.preventDefault()
    router.push("/prihlaseni")
}

</script>

<template>
    <h2>Zapomenuté heslo</h2>
    <form v-if="state === 'email'">
        <h3 style="margin-bottom: 20px;">Na e-mail ti bude zaslán <br> ověřovací kód pro obnovení hesla.</h3>
        <h3 class="nadpis">Zadej e-mail:</h3>
        <input :class="{ 'spatnej-input': spatnyEmail }" :oninput="zmena" type="text" v-model="email" placeholder="Např: pepa@zdepa.cz"
            inputmode="email">
        <button type="submit" class="tlacitko" @click="poslatEmail" :disabled="posilame">{{ posilame ? ". . ." : "Poslat e-mail" }}</button>
    </form>
    <form v-else-if="state === 'kod'">
        <h3 style="margin-bottom: 20px;">Zkontroluj prosím svou<br> e-mailovou schránku.</h3>
        <h3 class="nadpis">Kód z e-mailu:</h3>
        <input style="margin-bottom: 20px;" :class="{ 'spatnej-input': spatnyKod }" @:input="chekujUdaje('kod')" type="text" inputmode="numeric"
            v-model.trim="kod" placeholder="Např: 12345">
        <h3 class="nadpis">Nové heslo:</h3>
        <input :class="{ 'spatnej-input': spatnyHeslo }" @:input="chekujUdaje('heslo')" type="password" v-model="heslo"
            placeholder="Rozhodně ne 'Pepa123'">
        <button type="submit" class="tlacitko" @click="overitZmenu">Potvrdit</button>
    </form>
    <form v-else-if="state === 'konec'">
        <img src="../assets/pavoucekBezPozadi.svg" alt="Pavouk">
        <h3>Heslo úspěšně změněno!</h3>
        <h3><br>Tentokrát si heslo<br> prosím pamatuj. Díky!</h3>
        <button v-if="predvyplenyEmail === undefined" style="margin-top: 25px;" class="tlacitko" @click="presmerovat">Přihlásit</button>
    </form>

    <p v-if="state !== 'konec'">
        Nemáš ještě účet?
        <router-link to="/registrace">Registrace</router-link>
    </p>
</template>

<style scoped>
@import "../loginRegisterForma.css";

img {
    height: 150px !important;
    margin-bottom: 20px;
    align-self: center !important;
    margin-right: 0 !important;
}
</style>
