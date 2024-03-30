import { Scene, GameObjects, Input } from "phaser";
import CustomButton from "../CustomButton";
import { authenticate } from "../HttpHandler";
import { connectSocket, getSocket } from "../SocketHandler";
import { getEmitter } from "../EventEmitter";
import { SceneNames } from "../SceneData";
import { stopScenesExcept } from "../scene-util";
import { createPublicKey } from "crypto";

export class PlayAs extends Scene {

    private connectingTextObj: GameObjects.Text;
    private playerContainers: GameObjects.Container[] = [];

    constructor() {
        super("PlayAs");
    }

    preload() {
        // load assets
    }

    create() {

        this.add.text(512, 150, "Snakes and Ladder Online Multiplayer", {
            align: "center",
            font: "40px Courier",
            color: "#ffffff",
        }).setOrigin(0.5);

        this.add .text(512, 250, "Click to Play As:", {
            align: "center",
            font: "32px Courier",
            color: "#ffffff",
        }).setOrigin(0.5);

        this.add.text(512, 550, "Built with:\n\nPhaser 3.8 (Typescript) \nAWS Gateway Websocket API\n AWS Lambda (AWS SDK for GoLang v2)\nAWS DynamoDB",
            {
                lineSpacing: 10,
                align: "center",
                font: "italic 24px Courier",
                color: "#ffffff",
            }
        ).setOrigin(0.5);

        this.createLinkedInButton();

        this.createPlayerButton(0, 512, 300);
        this.createPlayerButton(1, 512, 360);
        this.createConnectingText();
            
        this.events.on("connect" , this.onConnect, this)
        this.events.on("disconnect" , this.onDisconnect, this)

    }

    createLinkedInButton() {
        
        const textObj = this.add.text(0, 0, `Developed by Sam Liew\nLinkedIn Profile:\nwww.linkedin.com/in/samliew94`, {
            align: 'center',
            font: '18px New Courier',
            color: "black",
            // lineSpacing: 8,
        }).setOrigin(0.5)

        const w = 350
        const h = 75

        const gfx = this.add.graphics()

        function drawGfx(color: "white" | "blue"){
            gfx.fillStyle(color === "white" ? 0xffffff : 0x0077b5);
            gfx.lineStyle(1, 0x000000);
            gfx.fillRoundedRect(-w/2, -h/2, w, h, 10)
            gfx.strokeRoundedRect(-w/2, -h/2, w, h, 10)
        }

        drawGfx("white")

        const container = this.add.container(512, 700, [gfx, textObj])
        container.setSize(w, h)
        container.setInteractive()


        container.on(Input.Events.POINTER_OVER , (_: any) => {
            drawGfx("blue") // red
            textObj.setColor("white")
        })
        
        container.on(Input.Events.POINTER_OUT , (_: any) => {
            drawGfx("white") // red
            textObj.setColor("black")

        })

        container.on(Input.Events.POINTER_DOWN , (_: any) => {
            container.setScale(0.95)
        })

        container.on(Input.Events.POINTER_UP , (_: any) => {
            container.setScale(1)

            window.open(`https://www.linkedin.com/in/samliew94`, `_blank`)
        })  


    }

    createPlayerButton(player: 0 | 1, x: number, y: number){

        const textObj = this.add.text(0, 0, `Player ${player+1}`, {
            align: 'center',
            font: '18px New Courier',
            color: "black",
        }).setOrigin(0.5)

        const w = 100
        const h = 50

        const gfx = this.add.graphics()

        function drawGfx(color: "white" | "green" | "yellow"){
            gfx.fillStyle(color === "white" ? 0xffffff : color === "green" ? 0x0141cf : 0xfd6600);
            gfx.lineStyle(1, 0x000000);
            gfx.fillRoundedRect(-w/2, -h/2, w, h, 10)
            gfx.strokeRoundedRect(-w/2, -h/2, w, h, 10)
        }

        drawGfx("white")

        const container = this.add.container(x, y, [gfx, textObj])
        container.setSize(w, h)
        container.setInteractive()


        container.on(Input.Events.POINTER_OVER , (_: any) => {
            drawGfx(player === 0 ? "green" : "yellow") // red
            textObj.setColor("white")
        })
        
        container.on(Input.Events.POINTER_OUT , (_: any) => {
            drawGfx("white") // red
            textObj.setColor("black")

        })

        container.on(Input.Events.POINTER_DOWN , (_: any) => {
            container.setScale(0.95)
        })

        container.on(Input.Events.POINTER_UP , (_: any) => {

            container.setScale(1)

            this.playerContainers.forEach(x=>x.setVisible(false))
            this.connectingTextObj.setVisible(true);

            authenticate(player).then(t => connectSocket(this, t))
                
        })  

        this.playerContainers.push(container);

        this.events.on(Phaser.Scenes.Events.WAKE, () => {
            this.onWake();
        }) 

    }

    onWake() {

        this.playerContainers.forEach(x=>x.setVisible(true))
        this.connectingTextObj.setVisible(false);
        
    }

    createConnectingText() {

        this.connectingTextObj = this.add.text(512, 350, "Connecting...",{
            align: "center",
            font: "52px New Courier",
            color: "white"
        }).setOrigin(0.5).setVisible(false)

    }

    onConnect () {

        this.scene.switch(SceneNames.BOARD)

    }
    
    onDisconnect () {

    }
   

    
}
