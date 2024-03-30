import { Game } from './scenes/Game';
import Phaser from 'phaser';
import { PlayAs } from './scenes/PlayAs';
import { Board } from './scenes/Board';

//  Find out more information about the Game Config at:
//  https://newdocs.phaser.io/docs/3.70.0/Phaser.Types.Core.GameConfig
const config: Phaser.Types.Core.GameConfig = {
    type: Phaser.AUTO,
    width: 1024,
    height: 768,
    parent: 'game-container',
    backgroundColor: '#000000',
    scale: {
        mode: Phaser.Scale.FIT,
        autoCenter: Phaser.Scale.CENTER_BOTH
    },
    scene: [
        Game,
        PlayAs,
        Board,
    ]
};

export default new Phaser.Game(config);
