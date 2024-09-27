<script setup lang="ts">
import { onMounted, PropType, ref, computed } from "vue";
import { checkTeapot, getToken, pridatOznameni } from "../../utils";
import { moznostiRocnik, moznostiTrida, moznostiSkupina } from "../../stores";
import axios from "axios";
import { useRouter } from "vue-router";
import KodTridy from "../../components/KodTridy.vue";

const emit = defineEmits(["prejmenovatTridu", "refresh"])

type Trida = {
    id: number,
    jmeno: string,
    ucitelID: number,
    kod: string,
    zamknuta: boolean,
    klavesnice: string
}

const props = defineProps({
    trida: {
        type: Object as PropType<Trida>,
        required: true
    },
    pocetStudentu: {
        type: Number,
        required: true
    }
})

const router = useRouter()

const tridaRocnikUprava = ref()
const tridaJmenoUprava = ref()
const tridaSkupinaUprava = ref()
const klavesniceUprava = ref()

function smazatTridu(e: Event) {
    e.preventDefault()

    if (props.pocetStudentu != 0) {
        let zprava: string = ""
        if (props.pocetStudentu == 1) zprava = "1 student se už připojil do této třídy."
        else if (props.pocetStudentu <= 4) zprava = `${props.pocetStudentu} studenti se už připojili do této třídy.`
        else zprava = `${props.pocetStudentu} studentů se už připojilo do této třídy.`
        if (!confirm(zprava + " Opravdu ji chcete smazat?")) return
    }

    axios.post("/skola/zmena-tridy", { trida_id: props.trida.id, zmena: "smazat" }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(() => {
        router.push("/skola")
        pridatOznameni(`Třída "${props.trida.jmeno}", byla smazána.`)
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
    })
}

function postKlavesnice(e: Event) {
    e.preventDefault()

    axios.post("/skola/zmena-tridy", { trida_id: props.trida.id, zmena: "klavesnice", hodnota: klavesniceUprava.value ? "qwerty" : "qwertz" }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(() => {
        emit("refresh")
    }).catch(e => {
        console.log(e)
        pridatOznameni("Chyba serveru")
    })
}

onMounted(() => {
    let a = props.trida.jmeno.split(/[. ]/)
    tridaJmenoUprava.value = a[1]
    tridaRocnikUprava.value = a[0] + (isNaN(+a[0]) ? " " : ".")
    tridaSkupinaUprava.value = a[3] == undefined ? "-" : a[3]

    klavesniceUprava.value = props.trida.klavesnice == "qwerty"
})

const tridaJmeno = computed(() => {
    return `${tridaRocnikUprava.value}${tridaJmenoUprava.value}${tridaSkupinaUprava.value != '-' ? ' ￨ ' + tridaSkupinaUprava.value : ''}`
})

</script>
<template>
    <div id="uprava-tridy">
        <KodTridy :kod="trida.kod" :id="trida.id" :zamknuta="trida.zamknuta" />

        <form>
            <div>
                <h3>Ročník:</h3>
                <select v-model="tridaRocnikUprava">
                    <option v-for="v in moznostiRocnik" :value="v" :key="v">{{ v }}</option>
                </select>
            </div>

            <div>
                <h3>Písmeno:</h3>
                <select v-model="tridaJmenoUprava">
                    <option v-for="v in moznostiTrida" :value="v" :key="v">{{ v }}</option>
                </select>
            </div>

            <div>
                <h3>Skupina:</h3>
                <select v-model="tridaSkupinaUprava">
                    <option v-for="v in moznostiSkupina" :value="v" :key="v">{{ v }}</option>
                </select>
            </div>

            <button class="tlacitko" @click="emit('prejmenovatTridu', $event, tridaJmeno)" :disabled="tridaJmeno == trida.jmeno">Potvrdit</button>

            <span>Podle jména se třídy řadí do ročníků v seznamu tříd.</span>
            <span>
                Skupina je určena pro rozdělení třídy na skupiny žáků. Např.: <b>3.B&nbsp;￨&nbsp;1</b> a <b>3.B&nbsp;￨&nbsp;2</b> To se může hodit,
                pokud vyučujete třídy po skupinkách.
            </span>
        </form>

        <form>
            <div>
                <h3>Preference klávesnice:</h3>

                <input type="checkbox" id="toggle1" class="toggle-checkbox" v-model="klavesniceUprava" @change="postKlavesnice" />
                <label for="toggle1" class="toggle-contejner">
                    <div>Qwertz</div>
                    <div>Qwerty</div>
                </label>
            </div>

            <span>
                Klávesnice ovlivňuje pouze to, jak se generují texty typu: <b>Naučená písmena &amp; Slova</b>, v lekcích <b>t,z - y,x,m</b>. Tam se
                totiž
                množina již naučených písmen liší klávesami Z a Y.
            </span>
            <span>Žákům, kteří mají jinou klácesnici než jaká je nastavena v jejich třídě, se u třídy zobrazí varování.</span>
        </form>

        <form>
            <div>
                <h3>Smazání třídy:</h3>
                <button type="button" class="cervene-tlacitko" @click="smazatTridu">Smazat třídu</button>
            </div>
        </form>
    </div>
</template>
<style scoped>
.toggle-contejner div {
    padding: 7px;
}

.toggle-contejner {
    font-size: 16px;
    width: 160px;
    height: 34px;
}

#uprava-tridy {
    width: 60%;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
}

#uprava-tridy>form {
    display: flex;
    flex-direction: column;
    width: 100%;
    align-items: center;
    gap: 8px;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 15px 25px;
}

#uprava-tridy>form>div {
    display: flex;
    width: 100%;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
}

#uprava-tridy .tlacitko,
#uprava-tridy .cervene-tlacitko {
    margin: 0;
    height: 34px;
    padding: 0 25px;
    width: auto;
}

form>span {
    opacity: 0.5;
    font-size: 0.85rem;
    margin-bottom: -4px;
}

form span:first-of-type {
    margin-top: 8px;
}

form span:last-of-type {
    margin-bottom: 0;
}


/* firefox nenenene */
@supports(-webkit-tap-highlight-color: black) {
    select:hover {
        background-color: var(--svetle-fialova) !important;
    }

    select {
        padding-left: 5px;
    }
}

select:hover {
    background-color: var(--fialova);
}

select option {
    background-color: var(--fialova) !important;
    font-weight: 400;
}
</style>