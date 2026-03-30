let username = "nodejsdevclient"
let ip = "ws://127.0.0.1:8098/ws"
let wsServer = new WebSocket(ip)
wsServer.addEventListener("open", () => {
    wsServer.send("USER " + username)
    console.log("Connected")
})
wsServer.addEventListener("message", (event) => {
    console.log(event.data)
    message = event.data.split(" ")
    switch (message[0]) {
        case "MOTD":
            console.log("MOTD: ", message[1])
            process.exit()
            break;
    
        default:
            break;
    }

})