import {Component, OnInit, signal} from "@angular/core";
import {ActivatedRoute, RouterLink} from "@angular/router";
import {take} from "rxjs";
import {DatePipe} from "@angular/common";
import {StashService} from "../stash-service/stash.service";
import {PageState} from "../page-state";
import {FullStash, FullUnlockedStash} from "../stash-service/types";
import {StashMap} from "../stash-map/stash-map";

@Component({
    imports: [DatePipe, RouterLink, StashMap],
    selector: "app-stash-page",
    styleUrl: "./stash-page.css",
    templateUrl: "./stash-page.html",
})
export class StashPage implements OnInit {
    public date: string = "";
    public st = signal<PageState<FullStash>>({state: "LOADING"});

    constructor(
        ar: ActivatedRoute,
        private readonly stashSvc: StashService,
    ) {
        ar.params.pipe(take(1)).subscribe((x) => {
            this.date = x["date"];
        });
    }

    ngOnInit() {
        this.stashSvc
            .getStash(this.date)
            .then((r) => {
                this.st.set({state: "LOADED", data: r});
            })
            .catch((err) => {
                this.st.set({state: "ERROR", error: err});
            });
    }

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
