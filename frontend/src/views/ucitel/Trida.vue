<script setup lang="ts">
import { onMounted } from 'vue';
import { useRoute } from 'vue-router';

const id = useRoute().params.id

onMounted(() => {
    
})

</script>
<template>
    <h1>{{ id }}</h1>
</template>
<style scoped>

</style>