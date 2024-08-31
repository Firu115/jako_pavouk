<script setup lang="ts">
import axios from 'axios';
import { checkTeapot, getToken, pridatOznameni, naJednoDesetiny } from '../../utils';
import { onMounted, ref } from 'vue';
import { useHead } from "unhead"
import router from '../../router';

useHead({
    title: "Třída"
})

const trida = ref({} as { id: number, jmeno: string, ucitel_id: number, kod: string, zamknuta: boolean, pocet_studentu: number })
const praceNove = ref([] as { id: number, cislo: number, datum: string, cpm: number, presnost: number }[])
const praceDoko = ref([] as { id: number, cislo: number, datum: string, cpm: number, presnost: number }[])

const nacitam = ref(false)

onMounted(() => {
    get()
})

function get() {
    nacitam.value = true
    axios.get("/skola/trida/", {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        trida.value = response.data.trida

        praceNove.value = []
        response.data.prace.sort((a: any, b: any) => b.datum.localeCompare(a.datum))
        for (let i = 0; i < response.data.prace.length; i++) {
            const prace1 = response.data.prace[i]
            let p = { id: prace1.id, cislo: response.data.prace.length - i, datum: new Date(prace1.datum).toLocaleDateString('cs-CZ'), cpm: prace1.cpm, presnost: prace1.presnost }

            if (prace1.cpm != -1) praceDoko.value.push(p)
            else praceNove.value.push(p)
        }

    }).catch(e => {
        if (checkTeapot(e)) return
        if (e.response.data.error == "sql: no rows in result set") {
            pridatOznameni("Taková třída neexistuje")
            router.push("/skola")
            return
        }
        pridatOznameni("Chyba serveru")
    }).finally(() => {
        nacitam.value = false
    })
}

</script>
<template>
    <h1>Třída {{ trida.jmeno == undefined ? "-.-" : trida.jmeno }}</h1>

    <div v-if="!nacitam" id="kontejner">
        <h2>Čeká na dokončení</h2>
        <div v-if="praceNove.length != 0" class="praceKontejner">
            <RouterLink :to="`/prace/${v.id}`" v-for="v in praceNove" class="prace">
                <div class="nadpisPrace">
                    <h3>Práce {{ v.cislo }}</h3>
                    <h4>{{ v.datum }}</h4>
                </div>
                <img class="play" src="../../assets/icony/start.svg" alt="Dokonceno!">
            </RouterLink>
        </div>
        <span v-else>Žádné</span>

        <h2>Dokončené</h2>
        <div v-if="praceDoko.length != 0" class="praceKontejner">
            <div v-for="v in praceDoko" class="prace hotova">
                <div class="nadpisPrace">
                    <h3>Práce {{ v.cislo }}</h3>
                    <h4>{{ v.datum }}</h4>
                </div>
                <div class="statistika">
                    <span><b>{{ naJednoDesetiny(v.cpm) }}</b> CPM</span>
                    <span><b>{{ naJednoDesetiny(v.presnost) }}</b> %</span>
                </div>
            </div>
        </div>
        <span v-else>Žádné</span>
    </div>

    <span v-if="praceNove.length == 0 && praceDoko.length == 0 && !nacitam" id="textZaci">Zatím tu nejsou žádné zadané práce.</span>
</template>
<style scoped>
.statistika span b {
    font-family: 'Red Hat Mono';
    font-size: 1.8rem;
}

.statistika span {
    font-size: 1.2rem;
    display: flex;
    align-items: baseline;
    gap: 5px;
    justify-content: end;
    height: 34px;
    align-self: center;
    position: relative;
    top: -2px;

    min-width: 115px;
}

.statistika {
    display: flex;
}

.play {
    height: 32px;
    margin-right: 12px;
    margin-bottom: 1px;

}

h2 {
    margin-top: 10px;
    margin-left: 5px;
}

.hotova {
    color: var(--seda);
    opacity: 80%;
    cursor: default !important;
}

#kontejner {
    display: flex;
    flex-direction: column;
    gap: 20px;
    text-align: left;
}

#kontejner>span {
    align-self: center;
    margin-bottom: 20px;
}

.praceKontejner {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 15px;
}

.prace {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 6px 12px;
    cursor: pointer;
    transition: 0.1s;
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 500px;
    height: 64px;
    text-decoration: none;
    color: var(--bila);
}

.prace:not(.hotova):hover {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.nadpisPrace {
    display: flex;
    flex-direction: column;
    align-items: start;
}

.nadpisPrace h3 {
    font-size: 1.4rem;
    font-weight: 500;
}

.nadpisPrace h4 {
    font-size: 1.1rem;
    font-weight: 200;
    margin: 0;
}
</style>