<script setup lang="ts">
import { ref, watch } from 'vue';

const text = ref("")

const textarea = ref<HTMLElement | null>(null)
const div = ref<HTMLElement>()

function scrollDiv() {
    div.value?.scrollTo(0, textarea.value!.scrollTop)
}

watch(text, () => {
    let t = text.value
        .replace(/\n$/g, '\n\n')
        .replace(/(\r\n|\r|\n)/g, "↵\n")
        .replace(/(^ )|( {2,})|(↵\n)|( $)|!(.(\r\n|\r|\n))|([^A-Za-z0-9ěščřžýáíéůúťďňóĚŠČŘŽÝÁÍÉŮÚŤĎŇÓ ,.!?;:_=+\-*/%()[\]{}<>])/g, `<m style='background-color: rgba(255, 0, 0, 0.4); font-family: "Red Hat Mono", monospace; border-radius: 3px'>$&</m>`)
    // krejzy https://codersblock.com/blog/highlight-text-inside-a-textarea/
    div.value!.innerHTML = t
})

defineExpose({ text })

</script>
<template>
    <div ref="div"></div>
    <textarea ref="textarea" placeholder="Text který budou žáci psát..." v-model="text" @scroll="scrollDiv"/>
</template>
<style scoped>
div {
    position: absolute;
    padding: 10px;
    width: 410px;
    height: calc(100vh - 90px - 60px - 40px - 25px - 30px - 5px - 40px - 15px); /* celá obrazovka - všechno co je nad tím */
    text-align: start;
    top: 275px;
    word-wrap: break-word;
    z-index: 1;
    background-color: var(--tmave-fialova);
    border-radius: 8px;
    font-size: 1rem !important;
    overflow-y: auto;
    text-decoration: none;

    color: rgb(156, 156, 156);
    font-family: "Red Hat Mono", monospace;

    white-space: pre-wrap;
    word-wrap: break-word;
    line-height: 1.2rem;
}

textarea {
    background-color: transparent;
    color: blue;
    border: 0;
    height: 100%;
    width: 100%;
    font-size: 1rem !important;
    text-decoration: none;
    color: var(--bila);
    padding: 10px;
    resize: none;
    z-index: 2;
    line-height: 1.2rem;
    font-family: "Red Hat Mono", monospace;
}

::placeholder {
    color: rgb(175, 175, 175);
    opacity: 1;
}

textarea::-webkit-scrollbar {
    width: 10px;
}

/* Customizes the track of the scrollbar */
textarea::-webkit-scrollbar-track {
    background: var(--tmave-fialova);
    /* Color of the track */
    border-radius: 3px;
    padding: 1px;
    cursor: default;
}

/* Customizes the thumb of the scrollbar */
textarea::-webkit-scrollbar-thumb {
    background: var(--fialova);
    /* Color of the thumb */
    border-radius: 3px;
    cursor: default;
}

/* Changes the thumb color on hover */
textarea::-webkit-scrollbar-thumb:hover {
    background: var(--svetle-fialova);
    /* Darker color on hover */
}

div::-webkit-scrollbar {
    width: 10px;
    background: transparent;
}
</style>