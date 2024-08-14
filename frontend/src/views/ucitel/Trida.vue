<script setup lang="ts">
import axios from 'axios';
import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { checkTeapot, getToken, pridatOznameni, naJednoDesetiny } from '../../utils';
import SipkaZpet from '../../components/SipkaZpet.vue';
import { moznostiRocnik, moznostiTrida } from '../../stores';
import ZadaniPrace from './ZadaniPrace.vue';
import router from '../../router';
import { useHead } from '@unhead/vue';

const id = useRoute().params.id

type Prace = { id: number, text: string, cas: number, datum: Date, prumer_cpm: number }

const trida = ref({} as { id: number, jmeno: string, ucitel_id: number, kod: string, zamknuta: boolean, pocet_studentu: number })
const prace = ref([] as Prace[])
const studenti = ref([] as { id: number, jmeno: string, email: string, cpm: number }[])

const tab = ref("zaci") // zaci, prace, zadani

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

        prace.value = []
        for (let i = 0; i < response.data.prace.length; i++) {
            const prace1 = response.data.prace[i]
            let p: Prace = { id: prace1.id, text: prace1.text, cas: prace1.cas, datum: new Date(prace1.datum), prumer_cpm: 0 }
            prace.value.push(p)
        }
        prace.value.sort((a: any, b: any) => b.datum.getTime() - a.datum.getTime())

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
    if (jmenoUprava.value == studentOznacenej.value.jmeno) {
        upravaStudenta.value = false
        return
    }

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

function zadano() {
    tab.value = 'prace'
    get()
}

</script>
<template>
    <h1 class="nadpisSeSipkou" style="margin: 0; direction: ltr;">
        <SipkaZpet />
        Třída: {{ trida.jmeno == undefined ? "-.-" : trida.jmeno }}
    </h1>

    <div id="dashboard">
        <div v-if="tab == 'zaci'" id="prepinacTabu">
            <button class="tlacitko" @click="tab = 'prace'">Práce</button>
        </div>
        <div v-else id="prepinacTabu">
            <button class="tlacitko" @click="tab = 'zaci'">Žáci</button>
        </div>
        <div id="kod">
            <div>
                <span @click="copy" :class="{ zamknutyKod: trida.zamknuta }">{{ trida.kod == undefined ? "------" :
                    trida.kod }}</span>
                <img v-if="!trida.zamknuta" src="../../assets/icony/zamekOpen.svg" alt="Odemčená třída" @click="zamek()">
                <img v-else src="../../assets/icony/zamekClosed.svg" alt="Zamčená třída" @click="zamek()">
            </div>
            <span :class="{ zamknutyKod: trida.zamknuta }">jakopavouk.cz/zapis/{{ trida.kod }}</span>
        </div>
        <form id="upravaTridy">
            <div>
                <select v-model="tridaRocnikUprava" style="margin-right: 10px;">
                    <option v-for="v in moznostiRocnik" :value="v">{{ v }}</option>
                </select>
                <select v-model="tridaJmenoUprava">
                    <option v-for="v in moznostiTrida" :value="v">{{ v }}</option>
                </select>
            </div>
            <button class="tlacitko" @click="prejmenovatTridu" :disabled="`${tridaRocnikUprava}${tridaJmenoUprava}` == trida.jmeno">Potvrdit</button>
        </form>
    </div>

    <div v-if="tab == 'zaci'" id="pulic">
        <div id="kontejner">
            <div v-for="st in studenti" class="blok" @click="select(st.id)" :class="{ oznaceny: selectnutej == st.id }">
                <div>
                    <h3>{{ st.jmeno }}</h3>
                    <h4>{{ st.email }}</h4>
                </div>
                <span><b>{{ naJednoDesetiny(st.cpm) }}</b> <span style="font-size: 0.95rem;">CPM</span></span>
            </div>
            <div v-if="studenti.length == 0" id="textZaci">Tady uvidíte všechny žáky... <br>Sdělte jim kód nahoře.</div>
        </div>
        <div v-if="selectnutej != -1" class="detail">
            <div id="vrsek">
                <img src="../../assets/pavoucekBezPozadi.svg" alt="Pavouk" width="200" height="126">
                <div v-if="!upravaStudenta">
                    <h2>
                        <span>{{ studentOznacenej.jmeno }}</span>
                        <img id="upravit" @click="upravaStudenta = true; jmenoUprava = studentOznacenej.jmeno" src="../../assets/icony/upravit.svg"
                            alt="Upravit">
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
            <h2 v-if="studenti.length != 0">Vyberte studenta!</h2>
        </div>
    </div>

    <div v-else-if="tab == 'prace'" id="praceKontejner">
        <div v-for="v, i in prace" class="prace">
            <div class="nadpisPrace">
                <h2>Práce {{ prace.length - i }}</h2>
                <h3>{{ v.datum.toLocaleDateString() }}</h3>
            </div>
            <span>{{ v.prumer_cpm }} CPM</span>
        </div>
        <span v-if="prace.length == 0" id="textZaci">Zatím tu nejsou žádné zadané práce. <br>První vytvoříte pomocí tlačítka dole.</span>
    </div>

    <ZadaniPrace v-else-if="tab == 'zadani'" :tridaID="trida.id" @zadano="zadano" />

    <div v-if="tab != 'zaci'" id="pridat" @click="tab = (tab == 'prace' ? 'zadani' : 'prace')"
        :style="{ transform: tab == 'zadani' ? 'rotate(-45deg)' : 'rotate(0deg)' }">
        <img src="../../assets/icony/plus.svg" alt="Přidat">
    </div>
</template>
<style scoped>
#textZaci {
    height: 400px;
    display: flex;
    align-items: center;
    justify-content: center;
}

#praceKontejner {
    width: 75%;
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
    max-width: 100%;
}

.nadpisPrace {
    display: flex;
    flex-direction: column;
    align-items: start;
}

.nadpisPrace h2 {
    font-size: 1.4rem;
    font-weight: 500;
}

.nadpisPrace h3 {
    font-size: 1.1rem;
    font-weight: 200;
    margin: 0;
}

#pridat {
    background-color: var(--tmave-fialova);
    border-radius: 100%;
    width: 55px;
    height: 55px;
    position: fixed;
    right: 30px;
    bottom: 25px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 13px;
    cursor: pointer;
    transition: background-color 0.15s, transform 0.3s;
    z-index: 10;
    box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
}

#pridat:hover {
    background-color: var(--fialova);
}

.zamknutyKod {
    color: rgba(255, 255, 255, 0.3);
    text-decoration: line-through;
}

#upravaTridy {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 10px 15px;
    display: flex;
    flex-direction: column;
    gap: 5px;
}

#upravaTridy div {
    display: flex;
    align-items: center;
}

#upravaTridy .tlacitko,
#prepinacTabu .tlacitko {
    width: 100px;
    margin-top: 5px;
    align-self: center;
}

#upravaTridy select {
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

#upravaTridy select:hover {
    background-color: var(--svetle-fialova);
}

#upravaTridy select option {
    font-family: "Red Hat Mono";
    background-color: var(--fialova) !important;
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
    padding: 10px;
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
    margin-bottom: -50px;
    /* countruju view margin */
}

#kontejner {
    width: 430px;
    display: flex;
    gap: 10px;
    height: calc(100vh - 90px - 60px - 40px - 25px - 30px - 5px);
    /* celá obrazovka - všechno co je nad seznamem zaku */
    flex-direction: column;
    overflow-y: scroll;
    padding-right: 10px;

    scrollbar-gutter: stable;
    scrollbar-width: auto;
}

.detail {
    width: 410px;
    height: 380px;
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
    font-weight: 500;
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