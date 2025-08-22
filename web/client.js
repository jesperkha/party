const log = document.getElementById("log")

const wsUrl = (location.protocol === "https:" ? "wss://" : "ws://") + location.host + "/connect"
const ws = new WebSocket(wsUrl)

ws.onopen = () => {
	log.textContent += `Connected to ${wsUrl}\n`
}

ws.onmessage = event => {
	try {
		const data = JSON.parse(event.data)
		if (data.type === "welcome") {
			log.textContent += `Server says: ${data.content}\n`
		} else {
			log.textContent += `Unknown message: ${event.data}\n`
		}
	} catch (e) {
		log.textContent += `Failed to parse message: ${event.data}\n`
	}
}

ws.onerror = err => {
	log.textContent += `Error: ${err}\n`
}

ws.onclose = () => {
	log.textContent += `Connection closed\n`
}

function broadcastMsg() {
	const msg = JSON.stringify({ type: "broadcast", content: "Hello, everyone!" })
	ws.send(msg)
	log.textContent += `Sent broadcast message: ${msg}\n`
}
