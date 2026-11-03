import {Component, OnInit, signal} from "@angular/core";
import {ActivatedRoute, RouterLink} from "@angular/router";
import {take} from "rxjs";
import {DatePipe} from "@angular/common";
import {StashService} from "../stash-service/stash.service";
import {PageState} from "../page-state";
import {FullStash} from "../stash-service/types";

@Component({
    imports: [DatePipe, RouterLink],
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
}
