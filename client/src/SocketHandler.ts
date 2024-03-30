import { Scene, Scenes } from "phaser";
import { config } from "./config"
import { SceneNames } from "./SceneData";

let socket: WebSocket

let sceneManager: Scenes.SceneManager;

export function connectSocket(scene: Scene, token:string){

    sceneManager = scene.scene.manager;
    
    const url = config.urls.ws + `?token=${token}`
    // const url = config.urls.ws

    socket = new WebSocket(url);

    socket.onopen = () => {
        console.log("connected to AWS websocket")        
        scene.events.emit("connect")
    }

    socket.onmessage = (event) => {
        
        const data = JSON.parse(event.data)
        console.log(`recv message from socket:`);
        console.log(data);

        sceneManager.getScene(SceneNames.BOARD).events.emit("board", data)

    }

    socket.onclose = () => {
        console.log("disconnected from AWS websocket")
        const board = sceneManager.getScene(SceneNames.BOARD)
        board.scene.switch(SceneNames.PLAY_AS)
    }

}

export function getSocket() {
    return socket;
}
