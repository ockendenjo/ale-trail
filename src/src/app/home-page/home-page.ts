import {Component, OnInit, signal} from "@angular/core";
import {CardComponent} from "../card/card.component";
import {PageState} from "../page-state";
import {StashService} from "../stash-service/stash.service";
import {UIShortStash} from "../stash-service/types";
import {gitCount} from "../build";

@Component({
    imports: [CardComponent],
    selector: "app-home-page",
    styleUrl: "./home-page.css",
    templateUrl: "./home-page.html",
})
export class HomePage implements OnInit {
    public st = signal<PageState<PageData>>({state: "LOADING"});

    constructor(private readonly stashSvc: StashService) {}

    ngOnInit() {
        this.stashSvc
            .listStashes()
            .then((res) => {
                const data: PageData = {
                    stashes: res,
                    hostsNeeded: res.some((i) => i.locked && i.hostNeeded),
                };
                this.st.set({state: "LOADED", data});
            })
            .catch((err) => {
                this.st.set({state: "ERROR", error: err});
            });
    }

    protected readonly gitCount = gitCount;
}

interface PageData {
    stashes: UIShortStash[];
    hostsNeeded: boolean;
}
