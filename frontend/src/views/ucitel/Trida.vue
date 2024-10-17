<script setup lang="ts">
import axios from "axios";
import { onMounted, ref, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { checkTeapot, getToken, pridatOznameni, naJednoDesetiny } from "../../utils";
import SipkaZpet from "../../components/SipkaZpet.vue";
import ZadaniPrace from "./ZadaniPrace.vue";
import { useHead } from "@unhead/vue";
import Tooltip from "../../components/Tooltip.vue";
import NastaveniTridy from "./NastaveniTridy.vue";
import KodTridy from "../../components/KodTridy.vue";
import { mobil } from "../../stores";

const id = useRoute().params.id

type Prace = { id: number, text: string, cas: number, datum: Date, prumerneCPM: number, prumernaPresnost: number, StudentuDokoncilo: number }

const trida = ref({} as { id: number, jmeno: string, ucitelID: number, kod: string, zamknuta: boolean, klavesnice: string })
const prace = ref([] as Prace[])
const studenti = ref([] as { id: number, jmeno: string, email: string, cpm: number }[])
const vsechnyTridy = ref([] as { id: number, jmeno: string, kod: string, zamknuta: boolean, pocet_studentu: number, klavesnice: string, pocet_praci: number }[])

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
const copyPraciIndex = ref(-1)
const copyTrida = ref(0)

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

        vsechnyTridy.value = []
        for (let i = 0; i < response.data.ostatniTridy.length; i++) {
            const trida = response.data.ostatniTridy[i]
            vsechnyTridy.value.push({ id: trida.id, jmeno: trida.jmeno, kod: trida.kod, zamknuta: trida.zamknuta, pocet_studentu: trida.pocet_studentu, klavesnice: trida.klavesnice, pocet_praci: trida.pocet_praci })
        }
        vsechnyTridy.value.sort((a: { jmeno: string; }, b: { jmeno: string; }) => a.jmeno.localeCompare(b.jmeno))

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
    if (mobil.value) return
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

function zadatDoJineTridy() {
    if (copyTrida.value == trida.value.id) {
        if (!confirm("Chcete práci zadat znovu do této třídy? Budou tak dvě práce se stejným textem.")) return
    }

    axios.post("/skola/pridat-praci", {
        "cas": prace.value[copyPraciIndex.value].cas,
        "trida_id": copyTrida.value,
        "text": prace.value[copyPraciIndex.value].text
    }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(() => {
        if (copyTrida.value == trida.value.id) get()
        copyPraciIndex.value = -1
    }).catch(e => {
        console.log(e)
        pridatOznameni("Chyba serveru")
    })
}

const posledniRychlostPrace = computed(() => {
    if (prace.value.length == 0) return -1

    for (let i = 0; i < prace.value.length; i++) {
        if (prace.value[i].prumerneCPM != -1) return prace.value[i].prumerneCPM
    }
    return -1
})

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
                <KodTridy :id="trida.id" :kod="trida.kod" :zamknuta="trida.zamknuta" />
                <h3>Jak se mohou žáci připojit?</h3>
                <ol>
                    <li>Sdělte žákům tento 4-místný kód</li>
                    <li>Žáci se zapíšou:
                        <ul>
                            <li>V záložce "Domů"</li>
                            <li>Uplně dole</li>
                            <li>Tlačítkem "Zapsat se"</li>
                        </ul>
                    </li>
                </ol>
            </div>
        </div>
        <div v-if="selectnutej != -1 && !mobil" class="detail">
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
        <div v-else-if="!mobil" class="detail" id="pred-kliknutim">
            <img src="../../assets/pavoucekBezPozadi.svg" alt="Pavouk">
            <h2 v-if="studenti.length == 0" style="font-size: 1.1rem;">Tady budou statistiky studentů!</h2>
            <h2 v-else>Vyberte studenta!</h2>
        </div>
    </div>
    <div v-else-if="tab == 'prace' && copyPraciIndex == -1" id="pulic-praci">
        <div id="prace-uprava-kontejner">
            <div v-for="v, i in prace" :key="v.id" class="uprava-pill" :style="{ opacity: (smazatPraciID == v.id || smazatPraciID == 0) ? 1 : 0.4 }">
                <div class="copy-btn" @click="copyPraciIndex = i">
                    <img src="../../assets/icony/copy.svg" alt="Kopírovat práci">
                </div>
                <div v-if="smazatPraciID != v.id" class="smazat-btn" @click="smazatPraciID = v.id">
                    <img src="../../assets/icony/trash.svg" alt="Smazat práci">
                </div>
                <div v-else class="smazat-btn" @click="smazatPraci(v)" @mouseleave="smazatPraciID = 0">
                    <img src="../../assets/icony/right.svg" alt="Smazat práci">
                </div>
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
            <div v-if="prace.length == 0" id="text-prace">
                <span>Zatím nejsou žádné zadané práce. <br>První vytvoříte pomocí tohoto tlačítka.</span>
                <img src="../../assets/icony/sipkaOhnuta.svg" alt="Šipka na tlačítko" width="100">
            </div>
        </div>
    </div>
    <div v-else-if="tab == 'prace' && copyPraciIndex != -1" id="copy-menu">
        <select v-model="copyTrida">
            <option :value="0">Vyberte třídu</option>
            <option v-for="t in vsechnyTridy" :value="t.id" :key="t.id">{{ t.jmeno }}</option>
        </select>
        <div>
            <button class="tlacitko" @click="copyPraciIndex = -1">Zrušit</button>
            <button class="tlacitko" :disabled="copyTrida == 0" @click="zadatDoJineTridy">Zadat práci</button>
        </div>
    </div>

    <ZadaniPrace v-else-if="tab == 'zadani'" :tridaID="trida.id" @zadano="zadano" :posledniRychlost="posledniRychlostPrace" />
    <NastaveniTridy v-else-if="tab == 'nastaveni'" ref="nastaveni" :trida="trida"
        :pocetStudentu="vsechnyTridy.find(t => t.id === trida.id)!.pocet_studentu" @prejmenovatTridu="prejmenovatTridu" @refresh="get" />

    <div v-if="tab == 'prace' || tab == 'zadani'" id="pridat" @click="tab = (tab == 'prace' ? 'zadani' : 'prace')"
        :style="{ transform: tab == 'zadani' ? 'rotate(-45deg)' : 'rotate(0deg)' }">
        <img src="../../assets/icony/plus.svg" alt="Přidat">
    </div>

</template>
<style scoped>
.uprava-pill {
    height: 60px;
    width: 28px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    transition: 0.2s;
}

.smazat-btn,
.copy-btn {
    width: 100%;
    padding: 5px;
    height: 28px;
    cursor: pointer;
    border-radius: 6px;
    background-color: var(--tmave-fialova);
    transition: 0.15s;
}

.smazat-btn:hover,
.copy-btn:hover {
    background-color: var(--fialova);
}

.copy-btn>img,
.smazat-btn>img {
    position: relative;
}

#copy-menu {
    background-color: var(--tmave-fialova);
    padding: 15px;
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 15px;
}

#copy-menu>select {
    min-width: 230px;
}

#copy-menu>div {
    display: flex;
    gap: 10px;
    justify-content: center;
}

#copy-menu .tlacitko {
    margin: 0;
    height: 34px;
}

#prace-uprava-kontejner {
    display: flex;
    flex-direction: column;
    gap: 15px;
}

#pulic-praci {
    display: flex;
    justify-content: space-between;
    width: 550px;
    gap: 8px;
}

.udaj2 {
    font-size: 26px !important;
    min-width: 95px !important;
}

.statistika span b {
    font-family: "Red Hat Mono";
    font-size: 29px;
}

.statistika span {
    font-size: 19px;
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
    font-size: 19px;
}

.udaje>div>span {
    font-size: 13px;
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
    font-size: 14px;
}

#posledni-prace>div>div:not(.udaje) {
    margin-top: 20px;
    font-size: 24px;
}

#posledni-prace>div>h4 {
    font-weight: 500;
    font-size: 19px;
}

.udaj {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 80%;
    align-self: center;
}

.udaj>h4 {
    font-size: 18px;
}

.udaj>div {
    font-size: 24px;
    font-weight: 500;
}

.udaj>div>span {
    font-size: 18px;
}

#text-zaci {
    height: 380px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    align-items: center;
    justify-content: center;
}

#text-zaci>ol {
    text-align: left;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

#text-zaci ul {
    margin-left: 24px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-top: 3px;
}

#text-zaci>h3 {
    margin-top: 25px;
    margin-bottom: 4px;
    font-size: 21px;
}

#text-prace {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
    position: fixed;
    bottom: 30px;
    right: 120px;
}

#text-prace>img {
    margin-left: 230px;
}

#prace-kontejner {
    display: flex;
    flex-direction: column;
    gap: 15px;
    width: 100%;
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
    height: 60px;
}

.nadpis-prace {
    display: flex;
    flex-direction: column;
    align-items: start;
}

.nadpis-prace h2 {
    font-size: 22px;
    font-weight: 500;
}

.nadpis-prace h3 {
    font-size: 18px;
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
    font-weight: 400;
}

#prepinac-tabu input:checked+label {
    font-weight: 500;
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
    font-size: 16px;
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
    font-size: 16px;
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
    font-size: 14px;
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
    font-size: 16px;
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
    font-size: 14px;
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
    font-size: 21px;
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

@media screen and (max-width: 900px) {
    #pulic {
        width: 400px;
    }

    #kontejner {
        height: auto;
        overflow-y: visible;
        padding-right: 0;
    }

    .nadpis-prace h2 {
        font-size: 18px;
    }

    .nadpis-prace h3 {
        font-size: 14px;
    }

    .prace {
        height: 48px;
        width: 400px;
        padding: 8px;
    }

    .statistika *:not(#tooltip) {
        font-size: 21px !important;
        height: 22px !important;
    }

    .statistika span {
        min-width: 90px;
    }

    .udaj2 {
        min-width: 60px !important;
    }

    .uprava-pill {
        display: none;
    }

    #prace-kontejner {
        width: auto;
    }

    #pulic-praci {
        width: auto;
        gap: 0;
    }

    .blok {
        padding: 4px 8px;
        height: 48px;
    }

    .blok h3 {
        height: 18px;
        position: relative;
        bottom: 2px;
    }

    #text-prace {
        bottom: 40px;
        right: 20px;
    }

    #text-prace>img {
        margin-left: 20px;
        transform: rotate(8deg);
    }
}
</style>