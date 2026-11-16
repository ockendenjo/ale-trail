import {Component, Input} from "@angular/core";
import {FullUnlockedStash} from "../stash-service/types";
import {StashMap} from "../stash-map/stash-map";

@Component({
    imports: [StashMap],
    selector: "app-live-stash",
    styleUrl: "./live-stash.css",
    templateUrl: "./live-stash.html",
})
export class LiveStash {
    @Input({required: true}) stash!: FullUnlockedStash;

    public getW3WAppUrl(stash: FullUnlockedStash): string {
        return `w3w://show?threewords=${stash.w3w}`;
    }

    public getW3WWebUrl(stash: FullUnlockedStash): string {
        return `https://w3w.co/${stash.w3w}`;
    }

    public getGoogleMapsUrl(stash: FullUnlockedStash): string {
        return `https://www.google.com/maps/search/?api=1&query=${stash.lat},${stash.lon}`;
    }

    protected getAppleMapsUrl(stash: FullUnlockedStash) {
        return `https://maps.apple.com/?ll=${stash.lat},${stash.lon}`;
    }

    public getGeoURI(stash: FullUnlockedStash): string {
        return `geo:${stash.lat},${stash.lon}`;
    }
}
