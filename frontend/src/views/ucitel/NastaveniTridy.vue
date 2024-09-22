<script setup lang="ts">
import { onMounted, PropType, ref, computed } from "vue";
import { checkTeapot, getToken, pridatOznameni } from "../../utils";
import { moznostiRocnik, moznostiTrida, moznostiSkupina } from "../../stores";
import axios from "axios";
import { useRouter } from "vue-router";
import KodTridy from "../../components/KodTridy.vue";

const emit = defineEmits(["prejmenovatTridu"])

type Trida = {
    id: number,
    jmeno: string,
    ucitelID: number,
    kod: string,
    zamknuta: boolean
}

const props = defineProps({
    trida: {
        type: Object as PropType<Trida>,
        required: true
    },
    pocetStudentu: Number
})

const router = useRouter()

const tridaRocnikUprava = ref()
const tridaJmenoUprava = ref()
const tridaSkupinaUprava = ref()

function smazatTridu(e: Event) {
    e.preventDefault()

    if (props.pocetStudentu != 0) {
        if (!confirm(`${props.pocetStudentu} studentů se už připojilo do této třídy. Opravdu ji chcete smazat?`)) return
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

onMounted(() => {
    let a = props.trida.jmeno.split(/[. ]/)
    tridaJmenoUprava.value = a[1]
    tridaRocnikUprava.value = a[0] + (isNaN(+a[0]) ? " " : ".")
    tridaSkupinaUprava.value = a[3] == undefined ? "-" : a[3]
})

const tridaJmeno = computed(() => {
    return `${tridaRocnikUprava.value}${tridaJmenoUprava.value}${tridaSkupinaUprava.value != '-' ? ' ￨ ' + tridaSkupinaUprava.value : ''}`
})

</script>
<template>
    <div id="nastaveni">
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
                    <button class="tlacitko" @click="emit('prejmenovatTridu', $event, tridaJmeno)"
                        :disabled="tridaJmeno == trida.jmeno">Potvrdit</button>
                </div>
            </form>

            <h3>Smazat třídu:</h3>
            <form>
                <button type="button" class="cervene-tlacitko" @click="smazatTridu">Smazat třídu</button>
            </form>
        </div>
    </div>
</template>
<style scoped>
#nastaveni {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 15px;
    width: 100%;
}

#uprava-tridy {
    width: 70%;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 10px 15px;
    display: flex;
    flex-direction: column;
    gap: 10px;
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
    padding: 0 15px;
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

</style>