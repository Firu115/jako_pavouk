<script setup lang="ts">
import axios from 'axios';
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { checkTeapot, getToken, pridatOznameni } from '../../utils';
import SipkaZpet from '../../components/SipkaZpet.vue';
import { moznostiRocnik, moznostiTrida } from '../../stores';
import ZadaniPrace from '../../components/ZadaniPrace.vue';
import router from '../../router';
import { useHead } from '@unhead/vue';

const id = useRoute().params.id

const trida = ref({} as { id: number, jmeno: string, ucitel_id: number, kod: string, zamknuta: boolean, pocet_studentu: number })
const studenti = ref([] as { id: number, jmeno: string, email: string, cpm: number }[])

const tab = ref("zaci")

const selectnutej = ref(-1)
const studentOznacenej = ref(ref({ jmeno: "...", email: "...@...", dokonceno: 0, daystreak: 0, medianRychlosti: -1, uspesnost: -1, klavesnice: "QWERTZ", nejcastejsiChyby: new Map }))

const upravaStudenta = ref(false)
const jmenoUprava = ref()
const tridaJmenoUprava = ref()
const tridaRocnikUprava = ref()

onMounted(() => {
    get()
})

function get() {
    axios.get("/skola/trida/" + id, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        trida.value = response.data.trida
        studenti.value = response.data.studenti
        studenti.value.sort((a: any, b: any) => a.jmeno.localeCompare(b.jmeno))

        let a = trida.value.jmeno.split(/[\. ]/)
        tridaJmenoUprava.value = a[1]
        tridaRocnikUprava.value = a[0] + (isNaN(+a[0]) ? " " : ".")

        useHead({
            title: trida.value.jmeno
        })
    }).catch(e => {
        if (checkTeapot(e)) return
        if (e.response.data.error == "sql: no rows in result set") {
            pridatOznameni("Taková třída neexistuje")
            router.push("/skola")
            return
        }
        console.log(e)
        pridatOznameni("Chyba serveru")
    })
}

function zamek() {
    trida.value.zamknuta = !trida.value.zamknuta

    axios.post("/skola/zmena-tridy", { trida_id: trida.value.id, zmena: "zamek" }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
        trida.value.zamknuta = !trida.value.zamknuta
    })
}

function select(id: number) {
    if (selectnutej.value == id) { //unselect
        selectnutej.value = -1
        return
    }
    upravaStudenta.value = false
    selectnutej.value = id
    axios.get("/skola/student/" + id, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        studentOznacenej.value = response.data
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
        selectnutej.value = -1
    })
}

function zmenaJmena(e: Event) {
    e.preventDefault()
    if (jmenoUprava.value == "" || jmenoUprava.value.length > 30) {
        pridatOznameni("Jméno musí být 1-30 znaků dlouhé")
        upravaStudenta.value = false
        return
    }
    axios.post("/skola/student/", { jmeno: jmenoUprava.value, id: selectnutej.value }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(_ => {
        upravaStudenta.value = false
        studentOznacenej.value.jmeno = jmenoUprava.value
        get()
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
        selectnutej.value = -1
    })
}

function prejmenovatTridu(e: Event) {
    e.preventDefault()

    let staryJmeno = trida.value.jmeno
    trida.value.jmeno = `${tridaRocnikUprava.value}${tridaJmenoUprava.value}`
    axios.post("/skola/zmena-tridy", { trida_id: trida.value.id, zmena: "jmeno", hodnota: `${tridaRocnikUprava.value}${tridaJmenoUprava.value}` }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
        trida.value.jmeno = staryJmeno
    })
}

function copy() {
    navigator.clipboard.writeText(trida.value.kod)
    pridatOznameni("Zkopírováno!", undefined, "copy")
}

</script>
<template>
    <h1 class="nadpisSeSipkou">
        <SipkaZpet />
        Třída: {{ trida.jmeno == undefined ? "-.-" : trida.jmeno }}
    </h1>
    <div id="dashboard">
        <div v-if="tab == 'zaci'" id="prepinacTabu">
            <h2>Zadat práci</h2>
            <button class="tlacitko" @click="tab = 'prace'">Práce</button>
        </div>
        <div v-else id="prepinacTabu">
            <h2>Zobrazit žáky</h2>
            <button class="tlacitko" @click="tab = 'zaci'">Zpět</button>
        </div>
        <div id="kod">
            <div>
                <span @click="copy" :class="{ zamknutyKod: trida.zamknuta }">{{ trida.kod == undefined ? "------" :
                    trida.kod }}</span>
                <img v-if="!trida.zamknuta" src="../../assets/icony/zamekOpen.svg" alt="Odemčená třída"
                    @click="zamek()">
                <img v-else src="../../assets/icony/zamekClosed.svg" alt="Zamčená třída" @click="zamek()">
            </div>
            <span :class="{ zamknutyKod: trida.zamknuta }">jakopavouk.cz/zapis/{{ trida.kod }}</span>
        </div>
        <form id="nastaveni">
            <div>
                <select v-model="tridaRocnikUprava" style="margin-right: 10px;">
                    <option v-for="v in moznostiRocnik" :value="v">{{ v }}</option>
                </select>
                <select v-model="tridaJmenoUprava">
                    <option v-for="v in moznostiTrida" :value="v">{{ v }}</option>
                </select>
            </div>
            <button class="tlacitko" @click="prejmenovatTridu"
                :disabled="`${tridaRocnikUprava}${tridaJmenoUprava}` == trida.jmeno">Potvrdit</button>
        </form>
    </div>
    <div v-if="tab == 'zaci'" id="pulic">
        <div id="kontejner">
            <div v-for="st in studenti" class="blok" @click="select(st.id)" :class="{ oznaceny: selectnutej == st.id }">
                <div>
                    <h3>{{ st.jmeno }}</h3>
                    <h4>{{ st.email }}</h4>
                </div>
                <span><b>{{ Math.round(st.cpm * 10) / 10 }}</b> <span style="font-size: 0.95rem;">CPM</span></span>
            </div>
            <div v-if="studenti.length == 0"
                style="height: 400px; display: flex; align-items: center; justify-content: center;">Zatím žádní
                studenti...</div>
        </div>
        <div v-if="selectnutej != -1" class="detail">
            <div id="vrsek">
                <img src="../../assets/pavoucekBezPozadi.svg" alt="Pavouk" width="200" height="126">
                <div v-if="!upravaStudenta">
                    <h2>
                        <span>{{ studentOznacenej.jmeno }}</span>
                        <img id="upravit" @click="upravaStudenta = true; jmenoUprava = studentOznacenej.jmeno"
                            src="../../assets/icony/upravit.svg" alt="Upravit">
                    </h2>
                    <h3>{{ studentOznacenej.email }}</h3>
                </div>
                <form v-else>
                    <input v-model="jmenoUprava" type="text">
                    <button type="submit" @click="zmenaJmena" class="tlacitko" id="ulozit">Uložit</button>
                </form>
            </div>
            <div>

            </div>
        </div>
        <div v-else class="detail" id="predKliknutim">
            <img src="../../assets/pavoucekBezPozadi.svg" alt="Pavouk">
            <h2>Vyber studenta!</h2>
        </div>
    </div>
    <ZadaniPrace v-else-if="tab == 'prace'" />
</template>
<style scoped>
.zamknutyKod {
    color: rgba(255, 255, 255, 0.3);
    text-decoration: line-through;
}

#nastaveni {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 10px 15px;
    display: flex;
    flex-direction: column;
    gap: 5px;
}

#nastaveni div {
    display: flex;
    align-items: center;
}

#nastaveni .tlacitko,
#prepinacTabu .tlacitko {
    width: 100px;
    margin-top: 5px;
    align-self: center;
}

#nastaveni select {
    border: none;
    border-radius: 5px;
    padding: 3px;
    font-size: 1.3rem;
    color: white;
    font-family: "Red Hat Mono";
    background-color: var(--fialova);
    cursor: pointer;
    transition: 0.2s;
}

#predKliknutim {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 40px;
    gap: 10px;
}

#predKliknutim img {
    width: 200px;
}

#upravit {
    width: 26px !important;
    height: 24px !important;
    cursor: pointer;
    margin: 2px 0 0 2px;
}

.tlacitko {
    width: auto;
    padding: 0 15px;
}

#prepinacTabu {
    background-color: var(--tmave-fialova);
    padding: 10px 15px;
    border-radius: 10px;
    width: 196px;
    display: flex;
    flex-direction: column;
}

#kod {
    background-color: var(--tmave-fialova);
    padding: 10px 15px;
    border-radius: 8px;
    transition: 0.1s;
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    width: 250px;
}

#kod div {
    font-size: 2rem;
    display: flex;
    gap: 5px;
    justify-content: center;
}

#kod div span {
    font-weight: 500;
    padding: 0 6px;
    transition: 0.2s;
}

#kod div span:hover {
    background-color: var(--fialova);
    border-radius: 8px;
    cursor: grab;
}

#kod img {
    width: 30px;
    height: 30px;
    position: relative;
    top: 4px;
    transition: 0.2s;
    cursor: pointer;
}

#kod img:hover {
    transform: scale(1.1);
}

#dashboard {
    width: 100%;
    height: 90px;
    margin: 20px 0 40px 0;
    display: flex;
    gap: 30px;
    justify-content: center;
}

#pulic {
    display: flex;
    justify-content: space-between;
    width: 860px;
}

#kontejner {
    width: 430px;
    display: flex;
    gap: 10px;
    flex-direction: column;
    height: 400px;
    overflow-y: scroll;
    padding-right: 10px;

    scrollbar-gutter: stable;
    scrollbar-width: auto;
}

.detail {
    width: 410px;
    height: 400px;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    padding: 10px 15px;
}

.detail h3 {
    font-size: 1.1rem;
    text-align: start;
}

#vrsek {
    margin-left: -8px;
    display: flex;
    align-items: center;
    gap: 5px;
    flex-wrap: wrap;
}

#vrsek img {
    width: 110px;
    height: 80px;
    object-fit: cover;
}

#vrsek h2 {
    display: flex;
    text-wrap: nowrap;
}

#vrsek h2 span {
    max-width: 265px;
    overflow: hidden;
    text-overflow: ellipsis;
    font-weight: 400;
}

#vrsek form {
    display: flex;
    flex-direction: column;
    align-items: center;
}

#vrsek div {
    max-width: 280px;
}

#vrsek h3 {
    text-overflow: ellipsis;
    overflow: hidden;
}

#vrsek h3:hover {
    overflow: visible;
}

#ulozit {
    transform: scale(0.8);
    margin-top: 5px;
    width: 100px;
    height: 35px;
}

#vrsek input {
    width: 280px;
    height: 30px;
    background-color: var(--fialova);
    border: 0;
    border-radius: 5px;
    color: var(--bila);
    padding: 10px;
    font-weight: normal;
    font-size: 1em;
}

#vrsek input:focus {
    outline: none !important;
}

.blok {
    border-radius: 10px;
    background-color: var(--tmave-fialova);
    padding: 6px 12px;
    cursor: pointer;
    transition: 0.1s;
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 100%;
}

.blok:hover {
    background-color: var(--fialova);
}

.blok h4 {
    font-size: 0.9rem;
    text-overflow: ellipsis;
    overflow: hidden;
    text-wrap: nowrap;
}

.blok h4:hover {
    overflow: visible;
}

.blok h3 {
    text-wrap: nowrap;
    font-weight: 400;
}

.blok span {
    font-size: 1.3rem;
}

.blok div {
    text-align: start;
    width: calc(100% - 120px);
}

.oznaceny {
    background-color: var(--fialova);
}

::-webkit-scrollbar {
    width: 10px;
    /* Width of the scrollbar */
}

/* Customizes the track of the scrollbar */
::-webkit-scrollbar-track {
    background: var(--tmave-fialova);
    /* Color of the track */
    border-radius: 3px;
    padding: 1px;
}

/* Customizes the thumb of the scrollbar */
::-webkit-scrollbar-thumb {
    background: var(--fialova);
    /* Color of the thumb */
    border-radius: 3px;
}

/* Changes the thumb color on hover */
::-webkit-scrollbar-thumb:hover {
    background: var(--svetle-fialova);
    /* Darker color on hover */
}
</style>