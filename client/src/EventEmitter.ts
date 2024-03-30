let emitter = new Phaser.Events.EventEmitter();

export function getEmitter() {
    return emitter;
}

export function destroyEmitter() {
    emitter.destroy()
}

export enum GameEvents {
    ON_CONNECTED="onConnected",
    ON_DISCONNECTED="onDisconnecteed",
    PLAY_AS_PLAYER="playAsPlayer",
}