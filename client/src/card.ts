import Phaser from "phaser";

export default class Card {
    private readonly scene;
    private readonly id;
    private readonly x;
    private readonly y;
    private cardContainer: Phaser.GameObjects.Container;

    constructor(scene: Phaser.Scene, id: number, x: number, y: number) {
        this.scene = scene;
        this.id = id;
        this.x = x;
        this.y = y;

        this.drawCard();
    }

    private drawCard() {
        const width = 100;
        const height = 100;
        const radius = 10;
        
        const container = this.scene.add.container(this.x, this.y);

        // const rect = this.scene.add.rectangle(0 ,0, width, height, 0xffffff);
        const rect = this.scene.add.graphics();
        rect.fillStyle(0xffffff, 1)
        rect.lineStyle(1, 0x000000)
        rect.fillRoundedRect(-width/2, -height/2, width, height, radius);
        rect.strokeRoundedRect(-width/2, -height/2, width, height, radius);

        const text = this.scene.add.text(0, 0, this.getCardFullName(this.id), {
            align: "center",
            color: this.getCardColor(this.id)
        })
        text.setOrigin(0.5)
        
        container.add([rect, text]);

        container.setSize(width, height);
        container.setInteractive({draggable: true});
        
        this.scene.input.on('dragstart', (ptr: any, go: any) => {
            this.scene.children.bringToTop(go);
        })

        this.scene.input.on('drag', (ptr: any, go: any, dragX: any, dragY:any) => {
            go.x = dragX;
            go.y = dragY;
        })

    }

    private getCardSuiteId(id: number) {
        const suiteId = Math.floor(id / 13);
        return suiteId;
    }

    private getCardSuiteName(suiteId: number) {

        suiteId = this.getCardSuiteId(suiteId);

        switch (suiteId) {
            case 0:
                return "Hearts";
            case 1:
                return "Diamonds";
            case 2:
                return "Clubs";
            default:
                return "Spades";
        }
    }

    private getCardValueName(id: number) {
        const idx = id % 13;

        if (idx <= 8) {
            return String(idx + 2);
        } else if (idx === 9) {
            return "Jack";
        } else if (idx === 10) {
            return "Queen";
        } else if (idx === 11) {
            return "King";
        } else if (idx === 12) {
            return "Ace";
        }
    }

    private getCardFullName(id: number) {
        return this.getCardValueName(id) + "\nof\n" + this.getCardSuiteName(id);
    }

    private getCardColor(id: number) {
        const suite = this.getCardSuiteId(id);
        console.log(suite);

        return suite <= 1 ? "red" : "black";
    }
}
