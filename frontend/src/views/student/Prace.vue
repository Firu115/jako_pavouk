<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";
import { checkTeapot, getToken, MojeMapa, pridatOznameni } from "../../utils";
import SipkaZpet from "../../components/SipkaZpet.vue";
import { onMounted, ref } from "vue";
import axios from "axios";
import Vysledek from "../../components/Vysledek.vue";
import { useHead } from "@unhead/vue";
import Psani from "../../components/Psani.vue";

const router = useRouter()
const route = useRoute()
const cislo: string = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id

useHead({
    title: "Práce ve třídě"
})

const text = ref([] as { id: number, znak: string, spatne: number, psat: boolean }[][]) // spatne: 0 ok, 1 spatne, 2 opraveno
const delkaTextu = ref(0)
const preklepy = ref(0)
const opravenePocet = ref(0)
const nejcastejsiChyby = ref()

const klavesnice = ref("")
const cas = ref(0)

const konec = ref(false)
const delkaNapsanehoTextu = ref(0)
const nacitamNovej = ref(false)

const rawText = ref([] as string[])

function get() {
    nacitamNovej.value = true
    axios.get("/skola/get-praci/" + cislo, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        rawText.value = response.data.text

        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0, psat: true })
                delkaTextu.value++
            })
        })

        klavesnice.value = response.data.klavesnice
        cas.value = response.data.cas

        while (delkaTextu.value <= 250) {
            prodlouzit()
        }
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

function konecTextu(o: number, p: number, n: MojeMapa, d: number) {
    opravenePocet.value = o
    preklepy.value = p
    nejcastejsiChyby.value = new MojeMapa(n)
    konec.value = true
    delkaNapsanehoTextu.value = d
}

function prodlouzit() {
    nacitamNovej.value = true

    let pocetSlov = text.value.length
    let lastSlovo = text.value[pocetSlov - 1]
        let lastPismeno = lastSlovo[lastSlovo.length - 1]
        if (lastPismeno.znak != " ") {
            delkaTextu.value++
            text.value[pocetSlov - 1].push({ id: delkaTextu.value, znak: " ", spatne: 0, psat: true })
        }

    rawText.value.forEach((slovo: string, i: number) => {
        text.value.push([])
        const slovoArr = [...slovo]
        slovoArr.forEach(pismeno => {
            text.value[pocetSlov + i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0, psat: true })
            delkaTextu.value++
        })
    })

    nacitamNovej.value = false
}
</script>

<template>
    <h1 class="nadpis-se-sipkou" style="margin: 0; direction: ltr;">
        <SipkaZpet />
        Práce ve třídě
    </h1>

    <Psani v-if="!konec" @konec="konecTextu" @prodlouzit="prodlouzit" :text :klavesnice :hide-klavesnice="false" :nacitamNovej :cas="cas" :delkaTextu
        :resetBtn="false" />

    <Vysledek v-else :preklepy :opravenych="opravenePocet" :delkaTextu="delkaNapsanehoTextu" :cas="cas" :cislo :posledni="true" :nejcastejsiChyby
        :pismena="'pracepraceprace'" />
</template>

<style scoped>
</style>