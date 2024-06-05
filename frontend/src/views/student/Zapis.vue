<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { checkTeapot, getToken, pridatOznameni } from '../../utils';
import axios from 'axios';

const router = useRouter()

const kodURL = useRoute().params.kod as string
const state = ref("kod") // kod, jmeno

const kod = ref("")
const jmeno = ref("")

onMounted(() => {
    if (kodURL != undefined) {
        kod.value = kodURL
    }
})

function potvrditKod(e: Event) {
    e.preventDefault()
    if (!/^[a-zA-Z0-9]{6}$/.test(kod.value)) {
        pridatOznameni("Kód je neplatný")
        return
    }
    axios.get("/skola/test-tridy/" + kod.value, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(_ => {
        state.value = 'jmeno'
    }).catch(e => {
        if (e.response.data.error == "Takova trida neexistuje") {
            pridatOznameni("Taková třída neexistuje")
            return
        }
        if (e.response.data.error == "Trida je zamcena") {
            pridatOznameni("Tato třída je zamčená")
            return
        }
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
    })
}

function zapsatSe(e: Event) {
    e.preventDefault()
    axios.post("/skola/zapis", { kod: kod.value, jmeno: jmeno.value }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(_ => {
        router.push("/trida")
    }).catch(e => {
        if (!checkTeapot(e)) {
            if (e.response.data.error == "Uz jsi ve tride") {
                pridatOznameni("Už jsi v jedné třídě")
                return
            }
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
    })
}

</script>
<template>
    <h1>Zápis do třídy</h1>
    <div id="kontejner">
        <img src="../../assets/pavoukSkola.svg" alt="Pavouk před školou">
        <form v-if="state == 'kod'">
            <h2>Zadej kód:</h2>
            <input type="text" v-model="kod" placeholder="ABC123">
            <button class="tlacitko" @click="potvrditKod">Potvrdit</button>
        </form>
        <form v-else>
            <h2>Zadej své jméno:</h2>
            <input id="jmeno" type="text" v-model="jmeno" placeholder="Novák Honza">
            <p>(Nejlépe příjmení první)</p>
            <button class="tlacitko" @click="zapsatSe">Zapsat se</button>
        </form>
    </div>
</template>
<style scoped>
.tlacitko {
    margin: 5px;
    width: 120px;
    align-self: center;
}

h2 {
    padding-left: 5px;
    font-size: 1.3em;
}

#kontejner {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    display: flex;
    margin-top: 20px;
    padding: 40px 25px 20px 25px;
    gap: 10px;
    width: 550px;
    justify-content: space-around;
}

#kontejner img {
    width: 250px;
}

form {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
}

form p {
    padding-left: 5px;
    font-size: 0.9rem;
    margin-bottom: 8px;
}

input {
    width: 190px;
    height: 50px;
    background-color: var(--fialova);
    border: 0;
    border-radius: 5px;
    color: var(--bila);
    padding: 10px;
    font-weight: normal;
    font-size: 1.9em;
    transition: all 0.15s cubic-bezier(0.5, 0, 0.5, 1) 0s;
    margin: 5px;
}

#jmeno {
    height: 40px;
    font-size: 1em;
}

input:focus {
    outline: none !important;
    width: 200px;
    margin: 5px 0;
}

form input::placeholder {
    color: rgba(255, 255, 255, 0.42);
}
</style>