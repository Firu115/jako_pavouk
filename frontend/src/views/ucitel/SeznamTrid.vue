<script setup lang="ts">
import axios from 'axios';
import { onMounted, ref } from 'vue';
import { checkTeapot, getToken, pridatOznameni } from '../../utils';

const rocniky = ref({ string: [] as { id: number, jmeno: string, ucitel_id: number, kod: string, zamknuta: boolean, pocet_studentu: number }[] })
const pridavani = ref(false)

const rocnik = ref("1.")
const trida = ref("A")

onMounted(() => {
    get()
})

function get() {
    axios.get("/skola/tridy", {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        rocniky.value = response.data.tridy
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
    <div id="rocniky" v-if="!pridavani && Object.keys(rocniky).length != 0">
        <div v-for="tridy, i in rocniky" class="rocnik">
            <h2 v-if="i != 'string'">{{ i }}. ročník</h2>
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
                    <option value="1.">1.</option>
                    <option value="2.">2.</option>
                    <option value="3.">3.</option>
                    <option value="4.">4.</option>
                    <option value="5.">5.</option>
                    <option value="6.">6.</option>
                    <option value="7.">7.</option>
                    <option value="8.">8.</option>
                    <option value="9.">9.</option>
                    <option value="Prima ">Prima</option>
                    <option value="Sekunda ">Sekunda</option>
                    <option value="Tercie ">Tercie</option>
                    <option value="Kvarta ">Kvarta</option>
                    <option value="Kvinta ">Kvinta</option>
                    <option value="Sexta ">Sexta</option>
                    <option value="Septima ">Septima</option>
                    <option value="Oktáva ">Oktáva</option>
                </select>
                <h3>Třída:</h3>
                <select v-model="trida">
                    <option value="A">A</option>
                    <option value="B">B</option>
                    <option value="C">C</option>
                    <option value="D">D</option>
                    <option value="E">E</option>
                    <option value="F">F</option>
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
    position: absolute;
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