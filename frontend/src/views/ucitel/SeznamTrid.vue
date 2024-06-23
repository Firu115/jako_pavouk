<script setup lang="ts">
import axios from 'axios';
import { onMounted, ref } from 'vue';
import { checkTeapot, getToken, pridatOznameni } from '../../utils';
import { moznostiRocnik, moznostiTrida, prihlasen } from '../../stores';
import { useHead } from '@unhead/vue';
import { useRouter } from 'vue-router';

const router = useRouter()

type Trida = { id: number, jmeno: string, ucitel_id: number, kod: string, zamknuta: boolean, pocet_studentu: number }
const rocniky = ref(new Map<string, Trida[]>())
const pridavani = ref(false)

const rocnik = ref("1.")
const trida = ref("A")

useHead({
    title: "Seznam tříd"
})

onMounted(() => {
    if (!prihlasen.value && getToken() == null) {
        router.push("/")
    }
    get()
})

function get() {
    axios.get("/skola/tridy", {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        Object.keys(response.data.tridy).forEach(key => {
            rocniky.value.set(key, response.data.tridy[key].sort((a: any, b: any) => a.jmeno.localeCompare(b.jmeno)))
        })
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
    })
}

function vytvorit(e: Event) {
    e.preventDefault()

    axios.post("/skola/create-trida", { jmeno: `${rocnik.value}${trida.value}` }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(_ => {
        pridavani.value = false
        get()
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
    })
}

</script>
<template>
    <h1>Třídy</h1>
    <div id="rocniky" v-if="!pridavani && rocniky.size !== 0">
        <div v-for="[rocnik, tridy] in rocniky" class="rocnik">
            <h2>{{ rocnik }}{{ isNaN(+rocnik) ? "" : ". ročník" }}</h2>
            <div id="kontejner">
                <div class="blok" v-for="t in tridy" @click="$router.push('/skola/' + t.id)">
                    <h3>{{ t.jmeno }}
                        <img v-if="!t.zamknuta" src="../../assets/icony/zamekOpen.svg" alt="Odemčená třída">
                        <img v-else src="../../assets/icony/zamekClosed.svg" alt="Zamčená třída">
                    </h3>

                    <div class="statistiky">
                        <span v-if="t.pocet_studentu == 1"><b>{{ t.pocet_studentu }}</b> student</span>
                        <span v-else-if="t.pocet_studentu >= 2 && t.pocet_studentu <= 4"><b>{{ t.pocet_studentu }}</b>
                            studenti</span>
                        <span v-else><b>{{ t.pocet_studentu }}</b> studentů</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div v-else-if="pridavani">
        <form id="pridatFormular">
            <h2 style="margin-bottom: 10px;">Vytvořit třídu</h2>

            <div>
                <h3>Ročník:</h3>
                <select v-model="rocnik" style="margin-right: 10px;">
                    <option v-for="v in moznostiRocnik" :value="v">{{ v }}</option>
                </select>
                <h3>Třída:</h3>
                <select v-model="trida">
                    <option v-for="v in moznostiTrida" :value="v">{{ v }}</option>
                </select>
            </div>
            <button class="tlacitko" @click="vytvorit">Vytvořit</button>
        </form>
    </div>
    <div v-else>
        Zatím tu žádné nejsou!
    </div>

    <div id="pridat" @click="pridavani = !pridavani"
        :style="{ transform: pridavani ? 'rotate(-45deg)' : 'rotate(0deg)' }">
        <img src="../../assets/icony/plus.svg" alt="Přidat">
    </div>
</template>
<style scoped>
.tlacitko {
    width: 100px;
    align-self: center;
}

#pridatFormular {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 20px 30px;
    display: flex;
    flex-direction: column;
    gap: 5px;
}

#pridatFormular div {
    display: flex;
    align-items: center;
    gap: 5px;
}

#pridatFormular select {
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

#pridatFormular select:hover {
    background-color: var(--svetle-fialova) !important;
}

#pridatFormular select option {
    font-family: "Red Hat Mono";
    background-color: var(--fialova) !important;
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
}

#pridat:hover {
    background-color: var(--fialova);
}

#rocniky {
    display: flex;
    flex-direction: column;
    gap: 25px;
    width: 100%;
}

.rocnik {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: 12px;
}

.rocnik h2 {
    margin-left: 5px;
}

.blok {
    border-radius: 10px;
    background-color: var(--tmave-fialova);
    padding: 10px 15px;
    cursor: pointer;
    transition: 0.1s;
    width: 30%;
    position: relative;
}

.blok img {
    width: 24px;
    height: 24px;
    position: relative;
    top: 3px;
}

.blok h3 {
    margin-bottom: 6px;
    font-size: 1.4rem;
    font-weight: 500;
}

.blok .statistiky {
    display: flex;
    flex-direction: column;
}

.blok:hover {
    background-color: var(--fialova);
}

#kontejner {
    display: flex;
    width: 100%;
    gap: 36px;
    flex-wrap: wrap;
}
</style>