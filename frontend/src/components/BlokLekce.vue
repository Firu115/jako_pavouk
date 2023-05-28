<script>
export default {
    name: "blokLekce",
    props: {
        data: Object,
    },
    methods: {
        prihlaste_se() {
            alert('Nejprve se prosím přihlašte')
        }
    },
    computed: {
        formatovany_pismena() {
            let vratit = "";
            for (let i = 0; i < this.data.pismena.length; i++) {
                vratit += i < this.data.pismena.length - 1 ? this.data.pismena.at(i) + ", " : this.data.pismena.at(i);
            }
            return vratit;
        }
    }
}
</script>

<template>
    <router-link class="lekceBlok" :class="{hotovoBlok: true/* je_dokonceno */}" v-if="$ls.getItem('token')" :to="'/lekce/' + data.pismena">
        <h2>Lekce: {{ formatovany_pismena }} </h2>
        <img class="fajvka" v-if="je_dokonceno" src="@/assets/icony/right.svg" alt="Dokonceno!">
    </router-link>
    <a v-else class="lekceBlok" @click="prihlaste_se"><h2>Lekce: {{ formatovany_pismena }}</h2></a>
</template>

<style scoped>
.lekceBlok {
    color: var(--bila);
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding: 12px 20px 12px 25px;
    text-decoration: none;
    border-radius: 10px;
    width: 500px;
    background-color: var(--tmave-fialova);
    height: 64px;
    transition-duration: 0.2s;
}

.lekceBlok:hover {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.hotovoBlok {
    color: var(--seda);
    opacity: 80%;
}

.lekceBlok h2 {
    align-self: center;
    font-size: 24px;
}
</style>
