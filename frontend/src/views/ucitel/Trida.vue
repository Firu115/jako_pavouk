<script setup lang="ts">
import axios from "axios";
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { checkTeapot, getToken, pridatOznameni, naJednoDesetiny } from "../../utils";
import SipkaZpet from "../../components/SipkaZpet.vue";
import { moznostiRocnik, moznostiTrida, moznostiSkupina } from "../../stores";
import ZadaniPrace from "./ZadaniPrace.vue";
import { useHead } from "@unhead/vue";
import Tooltip from "../../components/Tooltip.vue";

const id = useRoute().params.id

type Prace = { id: number, text: string, cas: number, datum: Date, prumerneCPM: number, prumernaPresnost: number, StudentuDokoncilo: number }

const trida = ref({} as { id: number, jmeno: string, ucitelID: number, kod: string, zamknuta: boolean })
const prace = ref([] as Prace[])
const studenti = ref([] as { id: number, jmeno: string, email: string, cpm: number }[])
const vsechnyTridy = ref([] as { id: string, jmeno: string, kod: string, zamknuta: boolean, pocet_studentu: number, pocet_praci: number }[])

const tab = ref("zaci") // zaci, prace, zadani

const selectnutej = ref(-1)
const studentOznacenej = ref({ jmeno: "...", email: "...@...", dokonceno: 0, daystreak: 0, rychlost: -1, uspesnost: -1, klavesnice: "QWERTZ", nejcastejsiChyby: new Map })
const cpmVPracich = ref(new Map<number, number>())
const presnostVPracich = ref(new Map<number, number>())

const upravaStudenta = ref(false)
const upravaTridy = ref("nic")
const jmenoUprava = ref()
const tridaJmenoUprava = ref()
const tridaRocnikUprava = ref()
const tridaSkupinaUprava = ref()

const studentTridaZmena = ref()

const nacitamStudenta = ref(false)

const router = useRouter()

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
            let p: Prace = { id: prace1.id, text: prace1.text, cas: prace1.cas, datum: new Date(prace1.datum), prumerneCPM: prace1.prumerne_cpm, prumernaPresnost: prace1.prumerna_presnost, StudentuDokoncilo: prace1.studentu_dokoncilo }
            prace.value.push(p)
        }
        prace.value.sort((a: any, b: any) => b.datum.getTime() - a.datum.getTime())

        let a = trida.value.jmeno.split(/[\. ]/)
        tridaJmenoUprava.value = a[1]
        tridaRocnikUprava.value = a[0] + (isNaN(+a[0]) ? " " : ".")
        tridaSkupinaUprava.value = a[2] == undefined ? "-" : a[2].slice(1)

        vsechnyTridy.value = response.data.ostatniTridy.sort((a: any, b: any) => a.jmeno.localeCompare(b.jmeno))
        studentTridaZmena.value = trida.value.id

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
    }).then(_ => {
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

function prejmenovatTridu(e: Event) {
    e.preventDefault()

    let staryJmeno = trida.value.jmeno
    trida.value.jmeno = `${tridaRocnikUprava.value}${tridaJmenoUprava.value}${tridaSkupinaUprava.value != '-' ? ' S' + tridaSkupinaUprava.value : ''}`
    axios.post("/skola/zmena-tridy", { trida_id: trida.value.id, zmena: "jmeno", hodnota: `${tridaRocnikUprava.value}${tridaJmenoUprava.value}${tridaSkupinaUprava.value != '-' ? ' S' + tridaSkupinaUprava.value : ''}` }, {
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

function smazatTridu(e: Event) {
    e.preventDefault()

    axios.post("/skola/zmena-tridy", { trida_id: trida.value.id, zmena: "smazat" }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(() => {
        router.push("/skola")
        pridatOznameni(`Třída "${trida.value.jmeno}", byla smazána.`)
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
    })
}

function copy() {
    navigator.clipboard.writeText(trida.value.kod)
    pridatOznameni("Zkopírováno!", undefined, "copy")
}

function zadano() {
    tab.value = "prace"
    get()
}

</script>
<template>
    <h1 class="nadpis-se-sipkou" style="margin: 0; direction: ltr;">
        <SipkaZpet />
        Třída: {{ trida.jmeno == undefined ? "-.-" : trida.jmeno }}
    </h1>

    <div id="dashboard">
        <div id="prepinac-tabu">
            <h3>Přepnout pohled:</h3>
            <button v-if="tab == 'zaci'" class="tlacitko" @click="tab = 'prace'">Práce</button>
            <button v-else class="tlacitko" @click="tab = 'zaci'">Žáci</button>
        </div>
        <div id="kod">
            <div>
                <span @click="copy" :class="{ 'zamknuty-kod': trida.zamknuta }">{{ trida.kod == undefined ? "------" : trida.kod }}</span>
                <Tooltip zprava="Uzamknout/Odemknout přístup novým žákům do této třídy" :sirka="192" :vzdalenost="12">
                    <img v-if="!trida.zamknuta" src="../../assets/icony/zamekOpen.svg" alt="Odemčená třída" @click="zamek()">
                    <img v-else src="../../assets/icony/zamekClosed.svg" alt="Zamčená třída" @click="zamek()">
                </Tooltip>
            </div>
            <span :class="{ 'zamknuty-kod': trida.zamknuta }">jakopavouk.cz/zapis/{{ trida.kod }}</span>
        </div>
        <div id="uprava-tridy">
            <form v-if="upravaTridy == 'nic'">
                <button class="tlacitko" @click="upravaTridy = 'jmeno'">Přejmenovat</button>
                <button class="tlacitko" @click="upravaTridy = 'smazat'">
                    <img src="../../assets/icony/trash.svg" alt="Smazat" width="25" height="25" style="filter: brightness(0.9);">
                </button>
            </form>
            <form v-else-if="upravaTridy == 'jmeno'" id="prejmenovani-tridy">
                <div>
                    <select v-model="tridaRocnikUprava">
                        <option v-for="v in moznostiRocnik" :value="v">{{ v }}</option>
                    </select>
                    <select v-model="tridaJmenoUprava">
                        <option v-for="v in moznostiTrida" :value="v">{{ v }}</option>
                    </select>
                    <select v-model="tridaSkupinaUprava">
                        <option v-for="v in moznostiSkupina" :value="v">{{ v }}</option>
                    </select>
                </div>
                <div>
                    <button class="tlacitko" @click="prejmenovatTridu"
                        :disabled="`${tridaRocnikUprava}${tridaJmenoUprava}${tridaSkupinaUprava != '-' ? ' S' + tridaSkupinaUprava : ''}` == trida.jmeno">Potvrdit</button>
                    <button class="tlacitko" @click="upravaTridy = 'nic'">Zpět</button>
                </div>
            </form>
            <form v-else-if="upravaTridy == 'smazat'">
                <button class="cervene-tlacitko" @click="smazatTridu">Opravdu smazat?</button>
                <button class="tlacitko" @click="upravaTridy = 'nic'">
                    <img src="../../assets/icony/plus.svg" alt="Smazat" width="28" height="28"
                        style="filter: brightness(0.9); transform: rotate(45deg);">
                </button>
            </form>
        </div>
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
            <div v-if="studenti.length == 0" id="text-zaci">Tady uvidíte všechny žáky... <br>Sdělte jim kód nahoře, aby se mohli připojit.</div>
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
                            <option v-for="v in vsechnyTridy" :value="v.id">{{ v.jmeno }}</option>
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
                <div v-if="prace.length == 0" style="margin-top: 40px; width: 100%;">
                    <span>Zatím žádné práce nebyly zadány</span>
                </div>
                <div v-else v-for="p, i in prace.slice(0, 3)">
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
            <h2 v-if="studenti.length != 0">Vyberte studenta!</h2>
        </div>
    </div>
    <div v-else-if="tab == 'prace'" id="prace-kontejner">
        <div v-for="v, i in prace" class="prace">
            <Tooltip :zprava="v.text.slice(0, 100) + '...'" :sirka="300" :vzdalenost="3">
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

    <ZadaniPrace v-else-if="tab == 'zadani'" :tridaID="trida.id" @zadano="zadano" />

    <div v-if="tab != 'zaci'" id="pridat" @click="tab = (tab == 'prace' ? 'zadani' : 'prace')"
        :style="{ transform: tab == 'zadani' ? 'rotate(-45deg)' : 'rotate(0deg)' }">
        <img src="../../assets/icony/plus.svg" alt="Přidat">
    </div>
</template>
<style scoped>
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

#posledni-prace>h3 {
    font-size: 1.15rem;
    font-weight: 500;
    grid-area: nadpis;
    align-self: center;
    justify-self: center;
    margin-bottom: 3px;
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
    align-items: center;
    justify-content: center;
}

#prace-kontejner {
    width: 75%;
    display: flex;
    flex-direction: column;
    gap: 15px;
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
}

#pridat:hover {
    background-color: var(--fialova);
}

.zamknuty-kod {
    color: rgba(255, 255, 255, 0.3);
    text-decoration: line-through;
}

#uprava-tridy {
    width: 250px;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 10px 15px;
}

#uprava-tridy form {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    height: 100%;
}

#uprava-tridy form .tlacitko,
#uprava-tridy form .cervene-tlacitko {
    margin: 0;
    padding: 0 9px;
}

#uprava-tridy form .tlacitko:has(img) {
    display: flex;
    align-items: center;
    justify-content: center;
}

#uprava-tridy form .tlacitko:nth-child(2) {
    padding: 1px;
    aspect-ratio: 1/1;
}

#prepinac-tabu .tlacitko {
    width: 100px;
    align-self: center;
}

#prejmenovani-tridy {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

#prejmenovani-tridy div {
    display: flex;
    gap: 10px;
}

#prejmenovani-tridy div .tlacitko {
    height: 28px;
    min-width: 50px;
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

#prepinac-tabu {
    background-color: var(--tmave-fialova);
    padding: 10px 15px;
    border-radius: 10px;
    width: 250px;
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
    border-radius: 8px;
}

#kod>span {
    transition: 0.2s;
}

#kod div span:hover {
    background-color: var(--fialova);
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
    width: 115%;
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
    gap: 40px;
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