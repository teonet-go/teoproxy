// Import TeoProxyClient class and Command enum
import TeoProxyClient from "./teoproxy.js";
import { Command } from "./teoproxy.js";

// The teoFortune Teonet peer name
const teoFortune = "8agv3IrXQk7INHy5rVlbCxMWVmOOCoQgZBF";

// Create TeoProxy client object
let teo = new TeoProxyClient();

// Connect to Teonet proxy websocket and Teonet peer api.
teo.connect("fortune-gui.teonet.dev", teoFortune, nextMessage);

// On message received from Teonet peer
teo.onmessage = (pac) => {
    console.debug("onmessage", pac);
    if (pac.cmd == Command.SendTo) {
        document.getElementById("fortune").innerHTML = pac.data;
    }
}

// Send next message request
export function nextMessage() {
    teo.cmd.sendTo(teoFortune, "fortb");
}
