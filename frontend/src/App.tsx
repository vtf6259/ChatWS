import { useState } from "react"

let connected: boolean = false


function connectWS(url: string, webSocket: WebSocket, setWebSocket: React.Dispatch<React.SetStateAction<WebSocket>>) {
  if(connected == true) {
    console.log("already connected")
    return
  }
  document.getElementById("status")!.innerText = "Connecting"
  let tmpWebSocket = new WebSocket(url)
  tmpWebSocket.addEventListener("message", (event) => {
    console.log(event.data)
  })
  tmpWebSocket.addEventListener("open", () => {
    document.getElementById("status")!.innerText = "Connected"
  })
  setWebSocket(tmpWebSocket)

  connected = true
}

export default function App() {
  const [webSocket, setWebSocket] = useState(new WebSocket(""))
  const [url, setURL] = useState('')
  const [messageToSend, setMessageToSend] = useState('')
  return (
    <>
    <h1>FolderChat</h1>
    <p id="status">Disconnected</p>
    <input type="text" id="serverURL" placeholder='Server url (ex "http://127.0.0.1:8098/ws")' value={url} onChange={(e) => {setURL(e.target.value)}}/>
    <button onClick={() => {connectWS(url, webSocket, setWebSocket)}}>Connect</button>
    <input placeholder="Text to send" value={messageToSend} onChange={(e) => {setMessageToSend(e.target.value)}} type="text" id="sendText"/>
    <button onClick={() => {webSocket.send(messageToSend)}}>Send</button>
    <button onClick={() => {webSocket.send("closing") /* to make go not panik */; webSocket.close(); connected = false; document.getElementById("status")!.innerText = "Disconnected"}}>Disconnect</button>
    <button onClick={() => {console.log(webSocket)}}>debug</button>
    </>
  )
}