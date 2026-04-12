<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue';

  const serverURI = useTemplateRef("serverURI")
  const connectStatus = useTemplateRef("connectStatus")
  let wsServer: WebSocket = new WebSocket("")
  let messages = ref({
    userMessages: []
  })
  onMounted(() => {
    const messagesLocalStorage = localStorage.getItem("messages")
    if (messagesLocalStorage != null) {
      messages.value = JSON.parse(messagesLocalStorage)
    }
  })
  function addEventListeners(websocket: WebSocket) {
    websocket.addEventListener("open", () => {
      console.log("open")
      if (connectStatus.value == null) return
      connectStatus.value.innerText = "Connected"
    })
    websocket.addEventListener("close", () => {
      if (connectStatus.value == null) return
      connectStatus.value.innerText = "Disconnected"
    })
  }
  function connect() {
    if (serverURI.value == undefined || connectStatus.value == null) {console.log("uhh"); return}
    if (serverURI.value.value == "") {
      connectStatus.value.innerText = "No empty feilds!" 
      setTimeout(() => {if (connectStatus.value == null) {console.log("uhh"); return}; connectStatus.value.innerText = "Disconnected"}, 2000)
      return
    }
    wsServer = new WebSocket(serverURI.value?.value)
    connectStatus.value.innerText = "Connecting"
    addEventListeners(wsServer) // this is needed otherwise it does not work
  }
</script>

<template>
  <div ref="connectStatus">Disconnected</div>
  <input ref="serverURI" placeholder="Server URL"/>
  <button v-on:click="connect()">Connect</button>
</template>

<style scoped></style>
