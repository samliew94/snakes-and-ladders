interface ZoneMap {
    gfx: Phaser.GameObjects.Graphics;
    zone: Phaser.GameObjects.Zone;       
}

export default class ZoneService {

    private readonly scene;
    private numCardsInZone = 0;
    private gfxes: Phaser.GameObjects.Graphics[] = [];
    private zones: ZoneMap[] = [];
    

    constructor(scene: Phaser.Scene) {
        this.scene = scene;

        this.createDropZone();
    }

    private createDropZone() {

        const x = 600;
        const y = 200;
        const n = 300;

        for (let i = 0 ; i < 2; i++ ){  

            const gfx = this.scene.add.graphics();
            const zone = this.scene.add.zone(x, i === 0 ? y : 600, n, n).setRectangleDropZone(n, n);

            this.zones.push({
                gfx,
                zone
            })
            
            this.drawZoneGraphics(this.zones[i])
    
        }   

        this.scene.input.on('dragenter', (ptr: any, go: any, dropZone: Phaser.GameObjects.Zone) => {

            const zone = this.zones.find(x=>x.zone == dropZone);

            if (zone) {
                this.drawZoneGraphics(zone, "blue");
            }

        })
        
        this.scene.input.on('dragleave', (ptr: any, go: any, dropZone: any) => {

            const zone = this.zones.find(x=>x.zone == dropZone);

            if (zone) {
                this.drawZoneGraphics(zone, "yellow");
            }

        })
        
        this.scene.input.on('drop', (ptr: any, go: any, dropZone: Phaser.GameObjects.Zone) => {
            
            go.x = dropZone.x + (this.numCardsInZone*5);
            go.y = dropZone.y + (this.numCardsInZone*5);

            this.numCardsInZone += 1;

        })

    }

    private drawZoneGraphics(zoneMap: ZoneMap, color: "yellow" | "blue" = "yellow") {

        const {gfx, zone} = zoneMap;

        gfx.clear();
        gfx.lineStyle(2, color === "yellow" ? 0xffff00 : 0x00ffff)
        gfx.strokeRect(zone.x - zone.input?.hitArea.width/2, zone.y-zone.input?.hitArea.height/2, zone.input?.hitArea.width, zone.input?.hitArea.height);

    }   

}