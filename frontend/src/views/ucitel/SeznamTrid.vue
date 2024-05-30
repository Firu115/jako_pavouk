<script setup lang="ts">
import axios from 'axios';
import { onMounted, ref } from 'vue';
import { checkTeapot, getToken, pridatOznameni } from '../../utils';

const rocniky = ref({ string: [] as { id: number, jmeno: string, ucitel_id: number, kod: string, zamknuta: boolean, pocet_studentu: number }[] })

onMounted(() => {
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
})

</script>
<template>
    <h1>Třídy</h1>
    <div id="rocniky">
        <div v-for="tridy, i in rocniky" class="rocnik">
            <h2>{{ i }}. ročník</h2>
            <div id="kontejner">
                <div class="blok" v-for="t in tridy" @click="$router.push('/skola/' + t.id)">
                    <h3>{{ t.jmeno }}
                        <img v-if="!t.zamknuta" src="../../assets/icony/zamekOpen.svg" alt="Odemčená třída">
                        <img v-else src="../../assets/icony/zamekClosed.svg" alt="Zamčená třída">
                    </h3>

                    <div class="statistiky">
                        <span v-if="t.pocet_studentu == 1"><b>{{ t.pocet_studentu }}</b> student</span>
                        <span v-if="t.pocet_studentu >= 2 && t.pocet_studentu <= 4"><b>{{ t.pocet_studentu }}</b>
                            studenti</span>
                        <span v-else><b>{{ t.pocet_studentu }}</b> studentů</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<style scoped>
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