<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { checkTeapot, format, getToken, MojeMapa, pridatOznameni } from '../utils';
import SipkaZpet from '../components/SipkaZpet.vue';
import { onMounted, ref } from 'vue';
import axios from 'axios';
import Vysledek from '../components/Vysledek.vue';
import { useHead } from '@unhead/vue';
import Psani from '../components/Psani.vue';
import { getCas } from '../stores.ts'

const router = useRouter()
const route = useRoute()
const pismena: string = Array.isArray(route.params.pismena) ? route.params.pismena[0] : route.params.pismena
const cislo: string = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id

useHead({
    title: "Cvičení " + pismena
})

const text = ref([] as { id: number, znak: string, spatne: number, psat: boolean }[][]) // spatne: 0 ok, 1 spatne, 2 opraveno
const delkaTextu = ref(0)
const preklepy = ref(0)
const opravenePocet = ref(0)
const nejcastejsiChyby = ref()

const posledni = ref(false)
const klavesnice = ref("")
const typTextu = ref("")

const konec = ref(false)
const delkaNapsanehoTextu = ref(0)
const nacitamNovej = ref(false)

function get() {
    nacitamNovej.value = true
    axios.get("/cvic/" + encodeURIComponent(pismena) + "/" + cislo, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0, psat: true })
                delkaTextu.value++
            })
        })

        posledni.value = response.data.posledni
        klavesnice.value = response.data.klavesnice
        typTextu.value = response.data.typ
    }).catch(e => {
        if (!checkTeapot(e)) {
            pridatOznameni()
            router.back()
        }
    }).finally(() => {
        nacitamNovej.value = false
    })
}

onMounted(() => {
    get()
})

function restart() {
    text.value = [] as { id: number, znak: string, spatne: number, psat: boolean }[][]
    delkaTextu.value = 0

    get()
    konec.value = false
}

function konecTextu(o: number, p: number, n: MojeMapa, d: number) {
    opravenePocet.value = o
    preklepy.value = p
    nejcastejsiChyby.value = new MojeMapa(n)
    konec.value = true
    delkaNapsanehoTextu.value = d
}

async function prodlouzit() {
    nacitamNovej.value = true
    axios.get("/cvic/" + encodeURIComponent(pismena) + "/" + cislo, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        let pocetSlov = text.value.length
        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[pocetSlov + i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0, psat: true })
                delkaTextu.value++
            })
        })
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni()
        }
    }).finally(() => {
        nacitamNovej.value = false
    })

}
</script>

<template>
    <h1 class="nadpisSeSipkou" style="margin: 0; direction: ltr;">
        <SipkaZpet />
        Lekce: {{ format(pismena) }}
    </h1>
    <h2>Cvičení: {{ cislo }}</h2>

    <Psani v-if="!konec" @konec="konecTextu" @restart="restart" @prodlouzit="prodlouzit" :text :klavesnice :hide-klavesnice="false" :nacitamNovej
        :cas="getCas(typTextu)" :delkaTextu />

    <Vysledek v-else @restart="restart" :preklepy :opravenych="opravenePocet" :delkaTextu="delkaNapsanehoTextu" :cas="getCas(typTextu)" :cislo
        :posledni :pismena :nejcastejsiChyby />
</template>

<style scoped>
.zvukIcon {
    width: 45px;
    height: 35px;
    margin-top: 1px;
}

#zvukBtn {
    position: absolute;
    right: 30px;
    bottom: 20px;
    background-color: var(--tmave-fialova);
    border-radius: 100px;
    width: 55px;
    height: 55px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
}

#zvukBtn:hover {
    background-color: var(--fialova);
}

#flex {
    display: flex;
    flex-direction: column;
    align-items: center;
}

#nabidka {
    margin: 20px 0 6px 0;
    width: var(--sirka-textoveho-pole);
}

#cas {
    float: left;
    width: 150px;
    display: block;
    text-align: left;
}

#preklepy {
    float: right;
    display: block;
    width: 150px;
    text-align: right;
}

#capslock {
    display: inline-block;
    color: red;
    font-weight: bold;
}

#ramecek {
    padding: 10px;
    height: 200px;
    border-radius: 10px 10px 0 0;
    background-color: var(--tmave-fialova);
    width: var(--sirka-textoveho-pole);
    overflow: hidden;
}

#text {
    display: flex;
    flex-wrap: wrap;
    position: relative;
    transition: ease 0.2s;
    top: 0em;
}

#fade {
    mask-image: linear-gradient(180deg, var(--tmave-fialova) 75%, transparent 97%);
    -webkit-mask-image: linear-gradient(180deg, var(--tmave-fialova) 75%, transparent 97%);
    height: 190px;
}

.slovo {
    display: flex;
    flex-wrap: nowrap;
}

.pismeno {
    border-radius: 3px;
    display: inline-flex;
    font-family: 'Red Hat Mono', monospace;
    font-weight: 400;
    font-size: 1.56rem;
    line-height: 2.2rem;
    text-decoration: none;
    padding: 0 1px;
    margin-right: 1px;
    border-bottom: 3px solid rgba(255, 255, 255, 0);
    /* aby se nedojebala vyska lajny když jdu na dalsi radek*/
    color: var(--bila);
    transition: 60ms;
}

#progress {
    height: 20px;
    background-color: var(--fialova);
    width: 0;
    border-bottom-left-radius: 10px;
    transition: ease 0.22s;
    text-align: right;
}

#bar {
    background-color: var(--tmave-fialova);
    width: var(--sirka-textoveho-pole);
    border-radius: 0 0 10px 10px;
    overflow: hidden;
}

.spravnePismeno {
    color: #9c9c9c;
}

.podtrzeni {
    border-bottom: 3px solid var(--bila);
    border-radius: 0;
    transition: 60ms;
}

.spatnePismeno {
    color: #ff0000;
}

.opravenePismeno {
    color: #b1529c;
}

h2 {
    margin-top: 5px;
    color: rgb(194, 194, 194);
}
</style>