import { Input, GameObjects, Scene } from "phaser";

export default class CustomButton extends GameObjects.Container {

    private readonly _width;
    private readonly _height;
    private readonly radius;
    private readonly gfx;
    private readonly text;
    private readonly oriFillColor = 0xffffff;
    private curFillColor;
    private readonly hoverTintFillColor = 0xff0000;
    private readonly oriTextColor = "black";
    private readonly hoverTintTextColor = "white";

    constructor(scene: Scene, width=100, height=100, radius=25,title='untitled') {
        super(scene);

        this._width = width;
        this._height = height;
        this.radius = radius;

        this.curFillColor = this.oriFillColor;

        // add gfx
        this.gfx = scene.add.graphics();
        this.drawGraphics();
        this.add(this.gfx);

        // add text
        this.text = scene.add.text(0, 0, title, {
            align: "center",
            color: this.oriTextColor
        });
        this.text.setOrigin(0.5);
        this.add(this.text);

        // set size
        this.setSize(width, height);

        this.setInteractive()
        // add to scene
        scene.add.existing(this);
    }

    private drawGraphics() {
        this.gfx.fillStyle(this.curFillColor)
        this.gfx.lineStyle(2, 0x000000)
        this.gfx.fillRoundedRect(-this._width/2,-this._height/2, this._width, this._height, this.radius);
        this.gfx.strokeRoundedRect(-this._width/2, -this._height/2, this._width, this._height, this.radius);
    }

    setDraggable(draggable: boolean) {

        if (draggable) {

            this.scene.input.setDraggable(this);

            this.scene.input.on('drag', (ptr: any, go: any, x: number, y: number) => {
                go.x = x;
                go.y = y;
            })
        } else {
            this.scene.input.off('drag');
        }
    }

    setOnClick(fn: () => void) {
        
        this.on(Input.Events.POINTER_OVER, (ptr: any) => {
            this.curFillColor = this.hoverTintFillColor;
            this.text.setColor(this.hoverTintTextColor)

            this.drawGraphics();
        })
        
        this.on(Input.Events.POINTER_OUT, (ptr: any) => {
            this.curFillColor = this.oriFillColor;
            this.text.setColor(this.oriTextColor)
            this.drawGraphics();
        })

        this.on(Input.Events.POINTER_DOWN, (ptr: any) => {
            this.setScale(0.95)
        }) 

        this.on(Input.Events.POINTER_UP, (ptr: any) => {
            this.setScale(1)
            fn();
        }) 

    }

}
