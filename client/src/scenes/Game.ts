import { Scene } from "phaser";
import { SceneNames } from "../SceneData";

export class Game extends Scene {

    constructor() {
        super(SceneNames.GAME);
    }

    preload() {}

    create() {
        
        this.scene.launch(SceneNames.PLAY_AS);
       
    }
}
