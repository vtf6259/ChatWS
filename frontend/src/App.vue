<script setup lang="ts">
import { useTemplateRef } from 'vue';

  const serverURI = useTemplateRef("serverURI")
  let wsServer: WebSocket
  function connect() {
    if (serverURI.value == undefined) return
    wsServer = new WebSocket(serverURI.value?.value)
  }
</script>

<template>
  <input ref="serverURI" placeholder="Server URL"/>
  <button v-on:click="connect()">Connect</button>
</template>

<style scoped></style>
