<script setup lang="ts">
import axios from "axios";
import { onMounted, onUnmounted, ref } from "vue";
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

const sources: EventSource[] = []

useHead({
    title: "Seznam tříd"
})

onMounted(() => {
    if (!prihlasen.value && getToken() == null) {
        router.push("/")
    }
    get()
})

onUnmounted(() => {
    sources.forEach(s => {
        s.close()
    })
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

        if (sources.length == 0) {
            rocniky.value.forEach(r => {
                r.forEach(t => {
                    let s: EventSource
                    if (window.location.hostname == "jakopavouk.cz") {
                        s = new EventSource("/api/skola/zaci-stream/" + t.id)
                    } else {
                        s = new EventSource("http://127.0.0.1:1323/api/skola/zaci-stream/" + t.id)
                    }
                    s.onmessage = function () {
                        get()
                    }
                    sources.push(s)
                })
            })
        }
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
        <form>
            <h2>Vytvořit třídu</h2>
            <div>
                <h3>Ročník:</h3>
                <select v-model="rocnik">
                    <option v-for="v in moznostiRocnik" :value="v" :key="v">{{ v }}</option>
                </select>
            </div>

            <div>
                <h3>Písmeno:</h3>
                <select v-model="trida">
                    <option v-for="v in moznostiTrida" :value="v" :key="v">{{ v }}</option>
                </select>
            </div>

            <div>
                <h3>Skupina:</h3>
                <select v-model="skupina">
                    <option v-for="v in moznostiSkupina" :value="v" :key="v">{{ v }}</option>
                </select>
            </div>

            <button class="tlacitko" @click="vytvorit">Vytvořit</button>

            <span>Podle jména se třídy řadí do ročníků v seznamu tříd.</span>
            <span>
                Skupina je určena pro rozdělení třídy. To se může hodit,
                pokud vyučujete třídy po menších skupinkách.
                <br>
                Např.: <b>3.B&nbsp;￨&nbsp;1</b> a <b>3.B&nbsp;￨&nbsp;2</b>.
            </span>
        </form>
    </div>
    <div v-else-if="rocniky.size === 0 && nacitam">
        Načítám...
    </div>
    <div v-else style="background-color: var(--tmave-fialova); padding: 20px; border-radius: 10px; max-width: 450px; margin: 0 5vw;">
        <h2 style="font-size: 21px;">Vítejte v rozhraní pro učitele!</h2>
        <br>
        Zde se vám budou třídy řadit do ročníků.
        <br>
        Žáci se do nich mohou připojit pomocí 4místného kódu, který bude pro každou třídu automaticky vytvořen.
    </div>

    <div v-if="rocniky.size == 0 && !pridavani && !nacitam" id="text-prace">
        <span>Zatím tu nejsou žádné třídy. <br>První vytvoříte pomocí tohoto tlačítka.</span>
        <img src="../../assets/icony/sipkaOhnuta.svg" alt="Šipka na tlačítko" width="100">
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

form {
    display: flex;
    flex-direction: column;
    max-width: 450px;
    align-items: center;
    gap: 8px;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 15px 25px;
    margin: 0 5vw;
}

form>h2 {
    margin-bottom: 15px;
}

form>div {
    display: flex;
    width: 100%;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
}

form>span {
    opacity: 0.5;
    font-size: 14px;
    margin-bottom: -4px;
}

form span:first-of-type {
    margin-top: 8px;
}

form span:last-of-type {
    margin-bottom: 0;
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
    font-size: 22px;
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

@media screen and (max-width: 1100px) {
    #rocniky {
        width: 420px;
    }
}

@media screen and (max-width: 600px) {
    #rocniky {
        width: 200px;
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