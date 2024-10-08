<script setup lang="ts">
import axios from "axios";
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { checkTeapot, getToken, pridatOznameni, naJednoDesetiny } from "../../utils";
import SipkaZpet from "../../components/SipkaZpet.vue";
import ZadaniPrace from "./ZadaniPrace.vue";
import { useHead } from "@unhead/vue";
import Tooltip from "../../components/Tooltip.vue";
import NastaveniTridy from "./NastaveniTridy.vue";
import KodTridy from "../../components/KodTridy.vue";

const id = useRoute().params.id

type Prace = { id: number, text: string, cas: number, datum: Date, prumerneCPM: number, prumernaPresnost: number, StudentuDokoncilo: number }

const trida = ref({} as { id: number, jmeno: string, ucitelID: number, kod: string, zamknuta: boolean, klavesnice: string })
const prace = ref([] as Prace[])
const studenti = ref([] as { id: number, jmeno: string, email: string, cpm: number }[])
const vsechnyTridy = ref([] as { id: string, jmeno: string, kod: string, zamknuta: boolean, pocet_studentu: number, pocet_praci: number }[])

const tab = ref("zaci") // zaci, prace, zadani, nastaveni

const selectnutej = ref(-1)
const studentOznacenej = ref({ jmeno: "...", email: "...@...", dokonceno: 0, daystreak: 0, rychlost: -1, uspesnost: -1, klavesnice: "QWERTZ", nejcastejsiChyby: new Map })
const cpmVPracich = ref(new Map<number, number>())
const presnostVPracich = ref(new Map<number, number>())

const upravaStudenta = ref(false)
const jmenoUprava = ref()

const studentTridaZmena = ref()

const nacitamStudenta = ref(false)

const router = useRouter()

const smazatPraciID = ref(0)

const nastaveni = ref()

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
        studenti.value.sort((a: { jmeno: string; }, b: { jmeno: string; }) => a.jmeno.localeCompare(b.jmeno))

        prace.value = []
        for (let i = 0; i < response.data.prace.length; i++) {
            const prace1 = response.data.prace[i]
            let p: Prace = { id: prace1.id, text: prace1.text, cas: prace1.cas, datum: new Date(prace1.datum), prumerneCPM: prace1.prumerne_cpm, prumernaPresnost: prace1.prumerna_presnost, StudentuDokoncilo: prace1.studentu_dokoncilo }
            prace.value.push(p)
        }
        prace.value.sort((a: { datum: Date; }, b: { datum: Date; }) => b.datum.getTime() - a.datum.getTime())

        vsechnyTridy.value = response.data.ostatniTridy.sort((a: { jmeno: string; }, b: { jmeno: string; }) => a.jmeno.localeCompare(b.jmeno))
        studentTridaZmena.value = trida.value.id

        useHead({
            title: trida.value.jmeno
        })
    }).catch(e => {
        console.log(e)
        if (e.response.data.error == "sql: no rows in result set") {
            pridatOznameni("Taková třída neexistuje")
            router.push("/skola")
            return
        }
        pridatOznameni("Chyba serveru")
    })
}

function select(id: number) {
    if (selectnutej.value == id) { //unselect
        selectnutej.value = -1
        return
    }
    upravaStudenta.value = false
    selectnutej.value = id
    studentOznacenej.value = { jmeno: "...", email: "...@...", dokonceno: 0, daystreak: 0, rychlost: -1, uspesnost: -1, klavesnice: "QWERTZ", nejcastejsiChyby: new Map }
    cpmVPracich.value = new Map<number, number>()
    presnostVPracich.value = new Map<number, number>()

    nacitamStudenta.value = true

    axios.get("/skola/student/" + id, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        studentOznacenej.value = response.data.student

        for (const key in response.data.cpmVPracich) {
            cpmVPracich.value.set(+key, response.data.cpmVPracich[key])
        }
        for (const key in response.data.presnostVPracich) {
            presnostVPracich.value.set(+key, response.data.presnostVPracich[key])
        }
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
        selectnutej.value = -1
    }).finally(() => {
        nacitamStudenta.value = false
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
        return
    }
    axios.post("/skola/student/", { jmeno: jmenoUprava.value, id: selectnutej.value }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(() => {
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

function zmenaStudentTridy(e: Event) {
    e.preventDefault()

    if (studentTridaZmena.value == trida.value.id) {
        upravaStudenta.value = false
        return
    }

    axios.post("/skola/student/", { trida_id: parseInt(studentTridaZmena.value), id: selectnutej.value }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(() => {
        upravaStudenta.value = false
        get()
        selectnutej.value = -1
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
        selectnutej.value = -1
    })
}

function prejmenovatTridu(e: Event, novyJmeno: string) {
    e.preventDefault()

    let staryJmeno = trida.value.jmeno
    trida.value.jmeno = novyJmeno
    axios.post("/skola/zmena-tridy", { trida_id: trida.value.id, zmena: "jmeno", hodnota: novyJmeno }, {
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

function zadano() {
    tab.value = "prace"
    get()
}

function smazatPraci(prace: Prace) {
    if (prace.StudentuDokoncilo != 0) {
        if (!confirm(`Tuto práci už dokončilo ${prace.StudentuDokoncilo} studentů! Opravdu ji chcete smazat?`)) return
    }

    axios.delete("/skola/smazat-praci/" + prace.id, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(() => {
        get()
        pridatOznameni("Práce byla smazána.")
        smazatPraciID.value = 0
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
    })
}

</script>
<template>
    <h1 class="nadpis-se-sipkou" style="margin: 0; direction: ltr;">
        <SipkaZpet />
        Třída: {{ trida.jmeno == undefined ? "-.-" : trida.jmeno }}
    </h1>

    <div v-if="tab != 'zadani'" id="prepinac-tabu">
        <input id="r1" type="radio" value="zaci" v-model="tab">
        <label for="r1">Studenti</label>

        <input id="r2" type="radio" value="prace" v-model="tab">
        <label for="r2">Práce</label>

        <input id="r3" type="radio" value="nastaveni" v-model="tab">
        <label for="r3">Nastavení</label>

        <span></span>
    </div>

    <div v-if="tab == 'zaci'" id="pulic">
        <div id="kontejner">
            <div v-for="st in studenti" class="blok" @click="select(st.id)" :class="{ oznaceny: selectnutej == st.id }" :key="st.id">
                <div>
                    <h3>{{ st.jmeno }}</h3>
                    <h4>{{ st.email }}</h4>
                </div>
                <span><b>{{ naJednoDesetiny(st.cpm) }}</b> <span style="font-size: 0.95rem;">CPM</span></span>
            </div>
            <div v-if="studenti.length == 0 && trida.jmeno" id="text-zaci">
                Tady uvidíte všechny žáky... <br>
                <KodTridy :id="trida.id" :kod="trida.kod" :zamknuta="trida.zamknuta" />
                Sdělte jim tento kód, aby se mohli připojit.
            </div>
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
                    <input v-model="jmenoUprava" type="text" placeholder="Doporučeno: Příjmení Jméno">
                    <div>
                        <select v-model="studentTridaZmena">
                            <option :value="-1">Odebrat</option>
                            <option v-for="v in vsechnyTridy" :value="v.id" :key="v.id">{{ v.jmeno }}</option>
                        </select>
                        <button type="submit" @click="zmenaJmena($event); zmenaStudentTridy($event)" class="tlacitko" id="ulozit">Potvrdit</button>
                    </div>
                </form>
            </div>

            <hr>

            <div class="udaj">
                <h4>Průměrná rychlost:</h4>
                <div>{{ naJednoDesetiny(studentOznacenej.rychlost) == -1 ? "-" : naJednoDesetiny(studentOznacenej.rychlost) }} <span>CPM</span></div>
            </div>

            <div class="udaj">
                <h4>Průměrná přesnost:</h4>
                <div>{{ naJednoDesetiny(studentOznacenej.uspesnost) == -1 ? "-" : naJednoDesetiny(studentOznacenej.uspesnost) }} <span>%</span></div>
            </div>

            <div class="udaj">
                <h4>Dní v řadě procvičuje:</h4>
                <div>{{ naJednoDesetiny(studentOznacenej.daystreak) }}</div>
            </div>

            <div class="udaj">
                <h4>Postup v kurzu:</h4>
                <div>{{ naJednoDesetiny(studentOznacenej.dokonceno) }} <span>%</span></div>
            </div>

            <hr>

            <div id="posledni-prace">
                <div v-if="prace.length == 0" style="margin-top: 40px; width: 100%; border: none;">
                    <span>Zatím žádné práce nebyly zadány</span>
                </div>
                <div v-else v-for="p, i in prace.slice(0, 3)" :key="i">
                    <h4>Práce {{ prace.length - i }}</h4>
                    <span>{{ p.datum.toLocaleDateString("cs-CZ") }}</span>

                    <div v-if="cpmVPracich.get(p.id) != undefined || nacitamStudenta" class="udaje">
                        <div>{{ cpmVPracich.get(p.id) == undefined ? "..." : naJednoDesetiny(cpmVPracich.get(p.id)!) }} <span>CPM</span></div>
                        <div>{{ presnostVPracich.get(p.id) == undefined ? "..." : naJednoDesetiny(presnostVPracich.get(p.id)!) }} <span>%</span></div>
                    </div>
                    <div v-else>
                        <Tooltip zprava="Nedokončil/a" :sirka="120" :vzdalenost="-3"><span>---</span></Tooltip>
                    </div>
                </div>
            </div>
        </div>
        <div v-else class="detail" id="pred-kliknutim">
            <img src="../../assets/pavoucekBezPozadi.svg" alt="Pavouk">
            <h2 v-if="studenti.length == 0" style="font-size: 1.1rem;">Tady budou statistiky studentů!</h2>
            <h2 v-else>Vyberte studenta!</h2>
        </div>
    </div>
    <div v-else-if="tab == 'prace'" id="pulic-praci">
        <div id="prace-smazat-kontejner">
            <div v-for="v in prace" :key="v.id">
                <img v-if="smazatPraciID == 0" src="../../assets/icony/trash.svg" alt="Smazat práci" @click="smazatPraciID = v.id">
                <img v-else :style="{ visibility: smazatPraciID == v.id ? 'visible' : 'hidden' }" style="padding: 6px; width: 40px;"
                    src="../../assets/icony/right.svg" alt="Smazat práci" @click="smazatPraci(v)" @mouseleave="smazatPraciID = 0">
            </div>
        </div>
        <div id="prace-kontejner">
            <div v-for="v, i in prace" :key="i" class="prace">
                <Tooltip :zprava="`<b>${v.cas / 60} min</b> | ${v.text.slice(0, 100)}...`" :sirka="300" :vzdalenost="3">
                    <div class="nadpis-prace">
                        <h2>Práce {{ prace.length - i }}</h2>
                        <h3>{{ v.datum.toLocaleDateString("cs-CZ") }}</h3>
                    </div>
                </Tooltip>

                <div class="statistika">
                    <Tooltip v-if="v.prumerneCPM != -1" zprava="Průměrná rychlost studentů" :sirka="160" :vzdalenost="5">
                        <span><b>{{ naJednoDesetiny(v.prumerneCPM) }}</b> CPM</span>
                    </Tooltip>
                    <Tooltip v-if="v.prumernaPresnost != -1" zprava="Průměrná přesnost studentů" :sirka="160" :vzdalenost="5">
                        <span><b>{{ naJednoDesetiny(v.prumernaPresnost) }}</b> %</span>
                    </Tooltip>
                    <Tooltip zprava="Studentů kteří mají hotovo" :sirka="160" :vzdalenost="5">
                        <span class="udaj2"><b>{{ v.StudentuDokoncilo }}</b>/<b>{{ studenti.length }}</b></span>
                    </Tooltip>
                </div>
            </div>
            <span v-if="prace.length == 0" id="text-zaci">Zatím tu nejsou žádné zadané práce. <br>První vytvoříte pomocí tlačítka dole.</span>
        </div>
    </div>


    <ZadaniPrace v-else-if="tab == 'zadani'" :tridaID="trida.id" @zadano="zadano" :posledniRychlost="prace[0].prumerneCPM" />
    <NastaveniTridy v-else-if="tab == 'nastaveni'" ref="nastaveni" :trida="trida"
        :pocetStudentu="vsechnyTridy.find(t => t.id === String(trida.id))!.pocet_studentu" @prejmenovatTridu="prejmenovatTridu" @refresh="get" />

    <div v-if="tab == 'prace' || tab == 'zadani'" id="pridat" @click="tab = (tab == 'prace' ? 'zadani' : 'prace')"
        :style="{ transform: tab == 'zadani' ? 'rotate(-45deg)' : 'rotate(0deg)' }">
        <img src="../../assets/icony/plus.svg" alt="Přidat">
    </div>

</template>
<style scoped>
#prace-smazat-kontejner {
    display: flex;
    flex-direction: column;
    gap: 15px;
}

#prace-smazat-kontejner>div {
    padding: 10px 0px 10px 10px;
    display: flex;
}

#prace-smazat-kontejner>div>img {
    height: 40px;
    cursor: pointer;
    transition: 0.15s;
    padding: 10px;
}

#prace-smazat-kontejner>div>img:hover {
    height: 40px;
    padding: 7px;
}

#pulic-praci {
    display: flex;
    justify-content: space-between;
    width: 80%;
    margin-right: 50px;
}

.udaj2 {
    font-size: 1.6rem !important;
    min-width: 95px !important;
}

.statistika span b {
    font-family: "Red Hat Mono";
    font-size: 1.8rem;
}

.statistika span {
    font-size: 1.2rem;
    display: flex;
    align-items: baseline;
    gap: 5px;
    justify-content: end;
    height: 34px;
    position: relative;
    top: -2px;

    min-width: 115px;
}

.statistika {
    display: flex;
}

.udaje {
    display: flex;
    flex-direction: column;
    margin-top: 12px;
}

.udaje>div {
    font-weight: 500;
    font-size: 1.2rem;
}

.udaje>div>span {
    font-size: 0.8rem;
}

hr {
    width: 90%;
    align-self: center;
    margin: 5px 0;
}

form input::placeholder {
    color: rgba(255, 255, 255, 0.42);
}

#posledni-prace {
    display: flex;
    height: 100%;
    margin-top: 5px;
}

#posledni-prace>div:nth-child(1),
#posledni-prace>div:nth-child(2) {
    border-right: 2px solid #afafaf;
}

#posledni-prace>* {
    width: 33%;
    padding-top: 5px;
    height: 100%;
}

#posledni-prace>div>span {
    font-size: 0.9rem;
}

#posledni-prace>div>div:not(.udaje) {
    margin-top: 20px;
    font-size: 1.5rem;
}

#posledni-prace>div>h4 {
    font-weight: 500;
    font-size: 1.2rem;
}

.udaj {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 80%;
    align-self: center;
}

.udaj>h4 {
    font-size: 1.1rem;
}

.udaj>div {
    font-size: 1.5rem;
    font-weight: 500;
}

.udaj>div>span {
    font-size: 1.1rem;
}

#text-zaci {
    height: 380px;
    display: flex;
    flex-direction: column;
    gap: 20px;
    align-items: center;
    justify-content: center;
}

#prace-kontejner {
    display: flex;
    flex-direction: column;
    gap: 15px;
    width: 90%;
}

.prace {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 6px 12px;
    transition: 0.1s;
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 100%;
}

.nadpis-prace {
    display: flex;
    flex-direction: column;
    align-items: start;
}

.nadpis-prace h2 {
    font-size: 1.4rem;
    font-weight: 500;
}

.nadpis-prace h3 {
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
    user-select: none;
}

#pridat:hover {
    background-color: var(--fialova);
}

#prepinac-tabu {
    display: flex;
    padding: 8px;
    border-radius: 100px;
    background-color: var(--tmave-fialova);
    margin: 20px 0;
    height: 46px;
}

#prepinac-tabu input {
    display: none;
}

#prepinac-tabu label {
    padding: 5px;
    width: 100px;
    z-index: 1;
    cursor: pointer;
    color: #c5c5c5;
    transition: 0.15s;
}

#prepinac-tabu input:checked+label {
    font-weight: 400;
    color: var(--bila);
}

#prepinac-tabu span {
    position: absolute;
    background-color: var(--fialova);
    width: 100px;
    height: 30px;
    border-radius: 100px;
    transition: 0.15s ease-out;
}

input[id="r1"]:checked~span {
    transform: translateX(0px);
}

input[id="r2"]:checked~span {
    transform: translateX(100px);
}

input[id="r3"]:checked~span {
    transform: translateX(200px);
}

#pred-kliknutim {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 40px;
    gap: 10px;
}

#pred-kliknutim img {
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
    max-height: calc(100vh - 90px - 60px - 40px - 25px - 30px - 5px);
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
    gap: 3px;
}

.detail h3 {
    font-size: 1rem;
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
    gap: 5px;
    width: 100%;
}

#vrsek form>div {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
}

#vrsek select {
    height: 30px;
    font-size: 1rem;
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
    font-size: 0.9rem;
    margin-top: 0;
    width: 80px;
    height: 30px;
    padding: 0;
}

#vrsek input {
    width: 100%;
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
    border-radius: 3px;
    padding: 1px;
}

/* Customizes the thumb of the scrollbar */
::-webkit-scrollbar-thumb {
    background: var(--fialova);
    border-radius: 3px;
}

/* Changes the thumb color on hover */
::-webkit-scrollbar-thumb:hover {
    background: var(--svetle-fialova);
    cursor: ns-resize;
}
</style>