<script setup lang="ts">
import axios from "axios";
import { onMounted, ref } from "vue";
import { checkTeapot, getToken, pridatOznameni } from "../../utils";
import { moznostiRocnik, moznostiTrida, moznostiSkupina, prihlasen } from "../../stores";
import { useHead } from "@unhead/vue";
import { useRouter } from "vue-router";

const router = useRouter()

type Trida = { id: number, jmeno: string, ucitel_id: number, kod: string, zamknuta: boolean, pocet_studentu: number, pocet_praci: number }
const rocniky = ref(new Map<string, Trida[]>())
const pridavani = ref(false)

const rocnik = ref("1.")
const trida = ref("A")
const skupina = ref("-")

const nacitam = ref(true)

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
            rocniky.value.set(key, response.data.tridy[key].sort((a: { jmeno: string }, b: { jmeno: string }) => a.jmeno.localeCompare(b.jmeno)))
        })
    }).catch(e => {
        if (checkTeapot(e)) return
        if (e.response.status == 401) {
            pridatOznameni("Sem nemáš přístup")
            router.push("/")
            return
        }
        console.log(e)
        pridatOznameni("Chyba serveru")
    }).finally(() => {
        nacitam.value = false
    })
}

function vytvorit(e: Event) {
    e.preventDefault()

    axios.post("/skola/create-trida", { jmeno: `${rocnik.value}${trida.value}${skupina.value != '-' ? ' ￨ ' + skupina.value : ''}` }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(() => {
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
        <div v-for="[rocnik, tridy] in rocniky" :key="rocnik" class="rocnik">
            <h2>{{ rocnik }}{{ isNaN(+rocnik) ? "" : ". ročník" }}</h2>
            <div id="kontejner">
                <div class="blok" v-for="t in tridy" :key="t.id" @click="$router.push('/skola/' + t.id)">
                    <h3>{{ t.jmeno }}</h3>

                    <hr style="margin: 0 8px 8px 8px; border: #c0c0c0 1px solid;">

                    <div style="display: flex; justify-content: space-around;">
                        <div class="statistiky">
                            <span v-if="t.pocet_studentu == 1"><b>{{ t.pocet_studentu }}</b> student</span>
                            <span v-else-if="t.pocet_studentu >= 2 && t.pocet_studentu <= 4"><b>{{ t.pocet_studentu }}</b> studenti</span>
                            <span v-else><b>{{ t.pocet_studentu }}</b> studentů</span>

                            <span v-if="t.pocet_praci == 0 || t.pocet_praci > 4"><b>{{ t.pocet_praci }}</b> prací</span>
                            <span v-else><b>{{ t.pocet_praci }}</b> práce</span>
                        </div>

                        <img v-if="!t.zamknuta" src="../../assets/icony/zamekOpen.svg" alt="Odemčená třída">
                        <img v-else src="../../assets/icony/zamekClosed.svg" alt="Zamčená třída">
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div v-else-if="pridavani">
        <form id="pridat-formular">
            <h2 style="margin-bottom: 15px;">Vytvořit třídu</h2>

            <div>
                <h3>Ročník:</h3>
                <select v-model="rocnik" style="margin-right: 10px;">
                    <option v-for="v in moznostiRocnik" :key="v" :value="v">{{ v }}</option>
                </select>
                <h3>Třída:</h3>
                <select v-model="trida" style="margin-right: 10px;">
                    <option v-for="v in moznostiTrida" :key="v" :value="v">{{ v }}</option>
                </select>
                <h3>Skupina:</h3>
                <select v-model="skupina">
                    <option v-for="v in moznostiSkupina" :key="v" :value="v">{{ v }}</option>
                </select>
            </div>
            <button class="tlacitko" @click="vytvorit">Vytvořit</button>
        </form>
    </div>
    <div v-else-if="rocniky.size === 0 && nacitam">
        Načítám...
    </div>
    <div v-else style="background-color: var(--tmave-fialova); padding: 20px; border-radius: 10px; max-width: 450px;">
        <b style="font-size: 1.3rem;">Vítejte v rozhraní pro učitele!</b>
        <br><br>
        Zatím tu nejsou žádné třídy, a tak svou první můžete vytvořit pomocí tlačítka dole.
        <br>
        Žáci se do ní mohou připojit pomocí 6ti-místného kódu, který bude automaticky vytvořen.
    </div>

    <div id="pridat" @click="pridavani = !pridavani" :style="{ transform: pridavani ? 'rotate(-45deg)' : 'rotate(0deg)' }">
        <img src="../../assets/icony/plus.svg" alt="Přidat">
    </div>
</template>
<style scoped>
.tlacitko {
    width: 100px;
    align-self: center;
}

#pridat-formular {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 20px 30px;
    display: flex;
    flex-direction: column;
    gap: 5px;
}

#pridat-formular div {
    display: flex;
    align-items: center;
    gap: 5px;
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
    box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
}

#pridat:hover {
    background-color: var(--fialova);
}

#rocniky {
    display: flex;
    flex-direction: column;
    gap: 25px;
    width: 640px;
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
    width: max(30%, 200px);
    position: relative;
}

.blok img {
    width: 30px;
    height: 30px;
    position: relative;
    top: 5px;
    right: 3px;
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
    flex-wrap: wrap;
    gap: 20px;
}

@media screen and (max-width: 1100px) {
    #rocniky {
        width: 420px;
    }
}

@media screen and (max-width: 600px) {
    #rocniky {
        width: 200px;
    }
}
</style>