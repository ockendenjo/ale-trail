import {Component, OnInit, signal} from "@angular/core";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {NgClass} from "@angular/common";
import {PageState} from "../page-state";
import {LiveStash} from "../live-stash/live-stash";
import {PreviousLocations} from "../stash-service/types";
import {StashService} from "../stash-service/stash.service";
import {MultiStashMap} from "../multi-stash-map/multi-stash-map";

@Component({
    imports: [FormsModule, ReactiveFormsModule, NgClass, LiveStash, MultiStashMap],
    selector: "app-location-page",
    styleUrl: "./previous-locations-page.component.css",
    templateUrl: "./previous-locations-page.component.html",
})
export class PreviousLocationsPage implements OnInit {
    public pg = signal<PageState<PreviousLocations>>({state: "LOADING"});

    constructor(private readonly stashSvc: StashService) {}

    ngOnInit() {
        this.stashSvc
            .getPreviousLocations()
            .then((pv) => {
                this.pg.set({state: "LOADED", data: pv});
            })
            .catch((e) => {
                this.pg.set({state: "ERROR", error: e});
            });
    }
}
