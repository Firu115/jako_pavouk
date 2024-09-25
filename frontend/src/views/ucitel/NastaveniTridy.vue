<script setup lang="ts">
import { onMounted, PropType, ref, computed } from "vue";
import { checkTeapot, getToken, pridatOznameni } from "../../utils";
import { moznostiRocnik, moznostiTrida, moznostiSkupina } from "../../stores";
import axios from "axios";
import { useRouter } from "vue-router";
import KodTridy from "../../components/KodTridy.vue";
import Tooltip from "../../components/Tooltip.vue";

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
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
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
    <KodTridy :kod="trida.kod" :id="trida.id" :zamknuta="trida.zamknuta" />
    <div id="uprava-tridy">
        <h3>Přejmenovat třídu:</h3>
        <form id="prejmenovani-tridy">
            <div>
                <select v-model="tridaRocnikUprava">
                    <option v-for="v in moznostiRocnik" :value="v" :key="v">{{ v }}</option>
                </select>
                <select v-model="tridaJmenoUprava">
                    <option v-for="v in moznostiTrida" :value="v" :key="v">{{ v }}</option>
                </select>
                <select v-model="tridaSkupinaUprava">
                    <option v-for="v in moznostiSkupina" :value="v" :key="v">{{ v }}</option>
                </select>
            </div>
            <div>
                <button class="tlacitko" @click="emit('prejmenovatTridu', $event, tridaJmeno)" :disabled="tridaJmeno == trida.jmeno">Potvrdit</button>
            </div>
        </form>

        <hr>

        <Tooltip zprava="Klávesnice ovlivňuje pouze to, jak se generují texty typu: <b>Naučená písmena / Slova</b>, v lekcích t,z - y,x,m. Tam se totiž množina již naučených písmen liší klávesami Z a Y." :sirka="380" :vzdalenost="8">
            <h3>Preference klávesnice:</h3>
        </Tooltip>
        <form>
            <input type="checkbox" id="toggle1" class="toggle-checkbox" v-model="klavesniceUprava" />
            <label for="toggle1" class="toggle-contejner">
                <div>Qwertz</div>
                <div>Qwerty</div>
            </label>

            <button class="tlacitko" @click="postKlavesnice" :disabled="klavesniceUprava == props.trida.klavesnice">Potvrdit</button>
        </form>

        <hr>

        <h3>Smazat třídu:</h3>
        <form>
            <button type="button" class="cervene-tlacitko" @click="smazatTridu">Smazat třídu</button>
        </form>

        <span style="opacity: 0.5; font-size: 0.9rem; margin-top: 16px; border-top: 2px solid gray; padding-top: 8px;">Nastavení budu muset ještě zkrášlit... :)</span>
    </div>
</template>
<style scoped>
hr {
    width: 80%;
    margin: 10px;
    align-self: center;
}

.toggle-contejner div {
    padding: 7px;
}

.toggle-contejner {
    font-size: 16px;
    width: 190px;
    height: 34px;
}

#uprava-tridy {
    width: 70%;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 10px 15px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin-top: 15px;
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
    height: 34px;
    padding: 0 15px;
    width: auto;
}

#uprava-tridy form .tlacitko:has(img) {
    display: flex;
    align-items: center;
    justify-content: center;
}

#uprava-tridy form .tlacitko:nth-child(2) {
    padding: 1px;
    aspect-ratio: 1/1;
    width: 40px;
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