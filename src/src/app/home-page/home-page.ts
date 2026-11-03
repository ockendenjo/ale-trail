import {Component, OnInit, signal} from "@angular/core";
import {CardComponent} from "../card/card.component";
import {PageState} from "../page-state";
import {StashService} from "../stash-service/stash.service";
import {ShortStash} from "../stash-service/types";

@Component({
    imports: [CardComponent],
    selector: "app-home-page",
    styleUrl: "./home-page.css",
    templateUrl: "./home-page.html",
})
export class HomePage implements OnInit {
    public st = signal<PageState<ShortStash[]>>({state: "LOADING"});

    constructor(private readonly stashSvc: StashService) {}

    ngOnInit() {
        this.stashSvc
            .listStashes()
            .then((res) => {
                this.st.set({state: "LOADED", data: res});
            })
            .catch((err) => {
                this.st.set({state: "ERROR", error: err});
            });
    }
}
