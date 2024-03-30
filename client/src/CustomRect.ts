import {Scene, GameObjects} from "phaser";

/**
 * creates a rounded rectangle with graphics api.  
 */
export default class CustomRect extends GameObjects.Container {

    private readonly gfx;
    private readonly tmpWidth;
    private readonly tmpHeight;
    private readonly tmpRadius;
    private fillStyle: number = 0xffffff;


    constructor(scene: Scene, width: number=100, height: number=100, radius: number=25) {

        super(scene);

        this.tmpWidth = width;
        this.tmpHeight = height;
        this.tmpRadius = radius;

        this.gfx = scene.add.graphics();
        this.redrawGfx();

        // const container = scene.add.container();
        this.setSize(width, height);
        this.setInteractive({draggable: true})
        this.add(this.gfx);
        
        scene.add.existing(this);

    }

    private redrawGfx() {
        this.gfx.fillStyle(this.fillStyle, 1)
        this.gfx.lineStyle(2, 0x000000)
        this.gfx.fillRoundedRect(-this.tmpWidth/2,-this.tmpHeight/2, this.tmpWidth, this.tmpHeight, this.tmpRadius);
        this.gfx.strokeRoundedRect(-this.tmpWidth/2, -this.tmpHeight/2, this.tmpWidth, this.tmpHeight, this.tmpRadius);
    }

    setBackgroundColor(fillColor: number) {
        this.fillStyle = fillColor;
        this.redrawGfx();
    }

}

