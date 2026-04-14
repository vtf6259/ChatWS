<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue';

  const serverURI = useTemplateRef("serverURI")
  const connectStatus = useTemplateRef("connectStatus")
  const username = useTemplateRef("username")
  const pubKeyRef = ref<HTMLInputElement | null>(null)
  const privKeyRef = ref<HTMLInputElement | null>(null)
  let pgpPrivKey: string
  let pgpPubKey: string
  let wsServer: WebSocket = new WebSocket("ws://0.0.0.0/aasjdioajsdoiajsdiajsdo")
  let messages = ref(<string[]>[])
  function genRecoverPhrase(pubkey: string, privkey: string) {
    const encode = btoa(pubkey + ":" + privkey)
    
  }
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
      wsServer.send("USER " + username.value?.value)
      connectStatus.value.innerText = "Connected"
    })
    websocket.addEventListener("close", () => {
      if (connectStatus.value == null) return
      connectStatus.value.innerText = "Disconnected"
    })
    websocket.addEventListener("message",(e) => {
      messages.value.push(e.data)
      let message = e.data.split(" ") 
      console.log(message)
    })
  }
  function connect() {
    if (serverURI.value == null || connectStatus.value == null || username.value == null || pubKeyRef.value == null || privKeyRef.value == null) {console.log("uhh"); return}
    if (serverURI.value.value == "" || username.value.value == "" || privKeyRef.value.value == "" || pubKeyRef.value.value == "") {
      connectStatus.value.innerText = "No empty feilds!" 
      setTimeout(() => {if (connectStatus.value == null) {console.log("uhh"); return}; connectStatus.value.innerText = "Disconnected"}, 2000)
      return
    }
    pgpPrivKey = atob(privKeyRef.value.value)
    pgpPubKey = atob(pubKeyRef.value.value)
    console.log(pgpPrivKey, "\n", pgpPubKey)
    wsServer = new WebSocket(serverURI.value?.value)
    connectStatus.value.innerText = "Connecting"
    addEventListeners(wsServer) // this is needed otherwise it does not work
  }
</script>

<template>
  <div ref="connectStatus">Disconnected</div>
  <div class="flexrow gap10px">
    <input class="flexrowMember" ref="serverURI" placeholder="Server URL">
    <input class="flexrowMember" ref="username" placeholder="Username">
    <input class="flexrowMember" ref="pubKeyRef" placeholder="Base64 PGP Public Key">
    <input class="flexrowMember" ref="privKeyRef" placeholder="Base64 PGP Private Key">
    <button class="flexrowMember" v-on:click="connect()">Connect</button>
  </div>
  <ul>
    <li v-for="message in messages">{{ message }}</li>
  </ul>
</template>

<style scoped></style>
