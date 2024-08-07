<script setup lang="ts">
import axios from "axios"
import BlokLekce from "../components/BlokLekce.vue";
import Rada from "../components/Rada.vue"
import { onMounted, onUnmounted, ref } from "vue"
import { Oznacene, checkTeapot, getToken, pridatOznameni, napovedaKNavigaci } from "../utils";
import { useHead } from "unhead"
import { useRouter } from "vue-router";
import { mobil, prihlasen } from "../stores"

useHead({
    title: "Lekce",
    meta: [
        {
            name: "description",
            content: "Jako Pavouk nabízí spoustu lekcí s řadou typů cviční. Od jednotlivých písmen, přes slova až k celým větám. Naučte se psát všemi deseti!",
        },
    ],
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz/lekce"
        }
    ]
})

const lekce = ref([[]] as { id: number, pismena: string, cislo: number }[][])
const dokoncene = ref([] as number[])
const o = new Oznacene()
const prvniNedokoncena = ref(1)

const router = useRouter()

onMounted(() => {
    const header = getToken() ? { headers: { Authorization: `Bearer ${getToken()}` } } : {}
    axios.get("/lekce", header)
        .then(response => {
            lekce.value = response.data.lekce
            dokoncene.value = response.data.dokoncene
            o.setMax(lekce.value.join(',').split(',').length) // pocet lekci

            let counter = 1
            let nebylaNedoko = true
            for (let i = 0; i < lekce.value.length; i++) {
                for (let j = 0; j < lekce.value[i].length; j++) {
                    lekce.value[i][j]['cislo'] = counter
                    counter += 1
                    if (dokoncene.value.includes(lekce.value[i][j]['id']) && nebylaNedoko) prvniNedokoncena.value += 1
                    else nebylaNedoko = false
                }
            }
        }).catch(e => {
            if (!checkTeapot(e)) {
                pridatOznameni()
                console.log(e)
            }
        })
    document.addEventListener("keydown", e1)
    document.addEventListener("keyup", e2)
    document.addEventListener("mousemove", zrusitVyber)
})

onUnmounted(() => {
    document.removeEventListener("keydown", e1)
    document.removeEventListener("keyup", e2)
    document.removeEventListener("mousemove", zrusitVyber)
})

let jede = false
let ms = 120

function e1(e: KeyboardEvent) {
    if (e.key == "ArrowUp" || e.key == "ArrowLeft") {
        e.preventDefault()
        if (jede) return

        if (o.index.value == 0) o.index.value = prvniNedokoncena.value + 1
        o.mensi()
        let lekceE: HTMLElement | null = document.querySelector(`[i="${o.index.value}"]`)

        jede = true

        let scroll = document.body.getBoundingClientRect().top
        window.scrollTo({ top: lekceE?.offsetTop! - 200, behavior: Math.abs(-scroll - lekceE?.offsetTop!) > 600 ? "instant" : "smooth" })
        setTimeout(() => { jede = false }, ms)
    } else if (e.key == "ArrowDown" || e.key == "ArrowRight") {
        e.preventDefault()
        if (jede) return

        if (o.index.value == 0) o.index.value = prvniNedokoncena.value - 1
        o.vetsi()
        let lekceE: HTMLElement | null = document.querySelector(`[i="${o.index.value}"]`)

        jede = true

        let scroll = document.body.getBoundingClientRect().top
        window.scrollTo({ top: lekceE?.offsetTop! - 200, behavior: Math.abs(-scroll - lekceE?.offsetTop!) > 600 ? "instant" : "smooth" })
        setTimeout(() => { jede = false }, ms)
    } if (e.key == "Enter") {
        e.preventDefault()
        let lekceE: HTMLElement | null = document.querySelector(`.oznacene`)
        if (lekceE == null || o.bezOznaceni) {
            o.bezOznaceni = true
            o.index.value = prvniNedokoncena.value
            lekceE = document.querySelector(`[i="${o.index.value}"]`)
            window.scrollTo({ top: lekceE?.offsetTop! - 200 })
        } else lekceE?.click()
    } else if (e.key == "Tab") {
        e.preventDefault()
        napovedaKNavigaci()
    }
}

function e2(e: KeyboardEvent) {
    if (e.key == "Enter") {
        e.preventDefault()
        let lekceE: HTMLElement | null = document.querySelector(`.oznacene`)
        router.push(lekceE?.getAttribute("href")!)  // chromium sus
    }
}

function zrusitVyber() {
    o.index.value = 0
}

</script>

<template>
    <h1>Lekce</h1>
    <div id="seznam">
        <Rada v-if="(!prihlasen || (dokoncene.length < 2 && lekce.length != 1)) && !mobil" />
        <h2>Střední řada</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 4" pismena="..." :jeDokoncena="false" />
        <!-- jen aby tam něco bylo než se to načte -->
        <BlokLekce v-else v-for="l in lekce[0]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" />
        <h2>Horní řada</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 5" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[1]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" />
        <h2>Dolní řada</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 3" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[2]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" />
        <h2>Diakritika</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 5" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[3]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" />
        <h2>Závěr kurzu</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 2" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[4]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" />
        <h2>Pro programátory</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 2" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[5]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" />
    </div>
</template>

<style scoped>
#seznam {
    display: flex;
    flex-direction: column;
    gap: 20px;
    text-align: left;
}

h2 {
    margin-top: 10px;
    margin-left: 5px;
}

.nohover {
    pointer-events: none;
    transition-duration: 0.1s;
}

@media screen and (max-width: 1100px) {
    #seznam {
        width: 70vw;
        align-items: center;
    }

    h2 {
        align-self: start;
    }
}
</style>