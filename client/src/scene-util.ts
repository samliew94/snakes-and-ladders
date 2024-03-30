import { Scene } from "phaser";

export function stopScenesExcept(scene: Scene, keys: string[]) {

    const manager = scene.scene.manager;

    manager.scenes.forEach(x=>{

        const key = x.scene.key;

        if (keys.includes(key)) {
            return
        }

        manager.remove(key)

    })

}