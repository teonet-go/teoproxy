
import TeoProxyClient from "./teoproxy.js";
import { Command } from "./teoproxy.js";

const teoFortune = "8agv3IrXQk7INHy5rVlbCxMWVmOOCoQgZBF";

let teo = new TeoProxyClient();

// console.debug("connect to:", teoFortune);
teo.connect("fortune-gui.teonet.dev", teoFortune, nextMessage);

teo.onmessage = (pac) => {
    console.debug("onmessage", pac);
    if (pac.cmd == Command.SendTo) {
        document.getElementById("fortune").innerHTML = pac.data;
    }
}

export function nextMessage() {
    teo.cmd.sendTo(teoFortune, "fortb");
}
