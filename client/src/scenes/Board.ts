import { Input, Scene, GameObjects } from 'phaser';
import { getSocket } from '../SocketHandler';

export class Board extends Scene
{
    private readonly xBeginsAt = 100;
    private readonly offsetY = 100
    private loadingTextObj: GameObjects.Text;
    private rollContainer: GameObjects.Container;
    private playerNameTextObj: GameObjects.Text;
    private turnTextObj: GameObjects.Text;
    private p1PosTextObj:GameObjects.Text;
    private p2PosTextObj:GameObjects.Text;
    private lastMessagesTextObjs : GameObjects.Text[] = [];

    constructor () {
        super('Board');
    }

    preload() {
        
    }

    create () {
        
        // draw the 5x5 board
        const n = this.xBeginsAt;
        let tileId = 1
        for (let y = n; y <= n*5; y += n) {
            
            for (let x = n; x <= n*5; x += n) {
                
                this.drawTile(x, y+this.offsetY, tileId)
                
                if (tileId === 3) {
                    this.addTextTileTeleport(x, y+this.offsetY, true, 9);
                } else if (tileId === 14) {
                    this.addTextTileTeleport(x, y+this.offsetY, true, 21);
                } else if (tileId === 12) {
                    this.addTextTileTeleport(x, y+this.offsetY, false, 6);
                } else if (tileId === 23) {
                    this.addTextTileTeleport(x, y+this.offsetY, false, 2);
                }
                
                this.add.text(x, y+this.offsetY, String(tileId++), {
                    align: 'center',
                    font: '24px Courier',
                    color: tileId-1 === 25 ? 'white' : 'black'
                }).setOrigin(0.5)
                
            }
            
        }
        
        this.initTextPlayerName()
        this.initTextPlayerTurn()
        this.initPlayerPositions()
        this.initTextLastMessages();
        
        this.createLoadingText();
        this.createBtnRoll()
        
        this.events.on("board", this.onWebsocketMessage, this);
        
        // call back when switch is called.
        this.events.on(Phaser.Scenes.Events.WAKE, () => {
            this.onWake();
        }) 

        this.onWake();

    }

    private onWake() {

        this.initTextPlayerName()
        this.initTextPlayerTurn()
        this.initPlayerPositions()
        this.initTextLastMessages();

        this.loadingTextObj.setVisible(true)
        this.rollContainer.setVisible(false)

        this.sendConnected();

    }

    private createLoadingText() {

        this.loadingTextObj = this.add.text(512, 80, "Loading ... Please wait", {
            align: "center",
            font: "52px New Courier",
            color: "white"
        }).setOrigin(0.5);
    }

    private drawTile(x: number, y: number, id: number) {

        const width = 100;
        const height = 100;

        const gfx = this.add.graphics();
        gfx.fillStyle(id === 1 ? 0x00ff00 : id === 25 ? 0xff0000 : 0xffffff)
        gfx.lineStyle(2, 0x000000)
        gfx.fillRoundedRect(-width/2, -height/2, width, height, 0)
        gfx.strokeRoundedRect(-width/2, -height/2, width, height, 0)

        const container = this.add.container(x, y, [gfx])
        container.setSize(width, height);

    }

    // event emitted from socket
    private onWebsocketMessage(jsonParsed: any) {   
        
        const {player, turn, positions, lastMessages} = jsonParsed;

        if (player !== undefined) {
            this.setTextPlayerName(player)
            this.loadingTextObj.setVisible(false)
            this.rollContainer.setVisible(true)
        }

        if (turn !== undefined) {
            this.setTextPlayerTurn(turn)
        }

        if (positions?.length) {
            this.setPlayerPositions(positions)
        }

        if (lastMessages) {
            this.setTextLastMessages(lastMessages);
        }

        

    }

    /**
     * load latest data from DDB
     */
    private sendConnected(){
        
        getSocket().send(JSON.stringify({"action":"connected"}))

    }

    private createBtnRoll() {

        const width = 200
        const height = 100;
        const x = 775;
        const y = 550;

        // draw roll button
        const gfx = this.add.graphics();
        gfx.fillStyle(0xffff00);
        gfx.lineStyle(1, 0x000000);
        gfx.fillRoundedRect(-width/2, -height/2, width, height, 50)
        gfx.strokeRoundedRect(-width/2, -height/2, width, height, 50)

        const textObj = this.add.text(0, 0, "Roll Dice", {
            color: "black",
            font: "40px New Courier"
        }).setOrigin(0.5)

        const container = this.add.container(x, y, [gfx, textObj])
        container.setSize(width, height)

        container.setInteractive()

        container.on(Input.Events.POINTER_OUT, (ptr: any) => {
            container.setScale(1)
        }) 

        container.on(Input.Events.POINTER_DOWN , (ptr: any) => {
            container.setScale(0.95)
        })

        container.on(Input.Events.POINTER_UP , (ptr: any) => {
            container.setScale(1)
            getSocket().send(JSON.stringify({"action":"roll"}))
        })

        container.setVisible(false)

        this.rollContainer = container;

    }
    
    private initTextPlayerName() {

        if (this.playerNameTextObj) {
            this.playerNameTextObj.setVisible(false)
            return;
        }

        this.playerNameTextObj = this.add.text(50, 50, ``, {
            font: '24px Courier',
            color: "white"
        }).setVisible(false)

    }

    private setTextPlayerName(player: number) {

        this.playerNameTextObj.text = `You are Player ${player+1}`
        this.playerNameTextObj.setVisible(true)

    }

    private initTextPlayerTurn() {

        if (this.turnTextObj ){
            this.turnTextObj.setVisible(false)
            return;
        }

        this.turnTextObj = this.add.text(50, 90, ``, {
            font: '16px Courier',
            color: "white"
        }).setVisible(false)
        
    }

    private setTextPlayerTurn(turn: boolean = false) {

        this.turnTextObj.text = `It's Player ${turn ? '1' : '2'}'s turn`
        this.turnTextObj.setVisible(true)
        
    }

    private initPlayerPositions() {

        if (this.p1PosTextObj && this.p2PosTextObj)  {
            this.p1PosTextObj.setVisible(false)
            this.p2PosTextObj.setVisible(false)
            return;
        }

        this.p1PosTextObj = this.add.text(0, 0, `P1`, {
            align: 'center',
            font: '20px Courier',
            color: 'black'
        })
        .setOrigin(0.5)
        .setVisible(false)

        this.p2PosTextObj = this.add.text(0, 0, `P2`, {
            align: 'center',
            font: '20px Courier',
            color: 'black'
        })
        .setOrigin(0.5)
        .setVisible(false)

    }

    private setPlayerPositions(positions: number[] = [0,0]) {
        
        const [p1, p2] = positions;

        const p1Row = Math.floor(p1 / 5)
        const p1Col = p1 % 5
        const p1x = this.xBeginsAt * (p1Col+1) - 16
        const p1y = 100 * (p1Row+1) + this.offsetY - 32

        this.p1PosTextObj.setPosition(p1x, p1y)
        this.p1PosTextObj.setVisible(true)

        const p2Row = Math.floor(p2 / 5)
        const p2Col = p2 % 5
        const p2x = this.xBeginsAt * (p2Col+1) + 16
        const p2y = 100 * (p2Row+1) + this.offsetY - 32

        this.p2PosTextObj.setPosition(p2x, p2y)
        this.p2PosTextObj.setVisible(true)

    }

    private initTextLastMessages() {

        if (this.lastMessagesTextObjs.length) {
            this.lastMessagesTextObjs.forEach(x=>x.setVisible(false))
            return;
        }

        let y = 150
        // initialize 5 textObjs
        for (let i = 0; i < 5; i++, y+=65) {

            const textObj = this.add.text(560, y, ``, {
                align: 'left',
                font: '16px Courier',
                color: 'white',
                wordWrap: {
                    width: 450
                },
                lineSpacing: 4
            }).setVisible(false)
            
            this.lastMessagesTextObjs.push(textObj)                    
        }            

    }

    private setTextLastMessages(texts: string[] = []) {

        let j = 0

        for (let i = texts.length-1; i > -1; i--, j++) {
            
            this.lastMessagesTextObjs[j].text = texts[i]
            this.lastMessagesTextObjs[j].setVisible(true)

        }
        
    }

    private addTextTileTeleport(x: number, y: number, forward: boolean, targetTileId: number, ) {

        const text = forward ? `Go To ${targetTileId}` : `Back To ${targetTileId}`;

        this.add.text(x, y+36, text, {
            align: 'center',
            font: '12px Courier',
            color: forward ? "green" : "red",
        }).setOrigin(0.5)

    }

}
