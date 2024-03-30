import { Scene } from "phaser";

export default function createText(scene: Scene, x: number, y: number, text: string) {

    const textObj = scene.add.text(0, 0, text, {
        font: "16px Arial",
        color: "#ffffff",
        align: "center"
    });

    scene.add.container(x, y, [textObj]);
    
}