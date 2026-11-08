import {Component, OnDestroy, OnInit, signal} from "@angular/core";
import {RouterLink, RouterOutlet} from "@angular/router";
import {UpdateService} from "./update.service";
import {environment} from "../environments/environment.dev";
import {gitCount} from "./build";

@Component({
    imports: [RouterOutlet, RouterLink],
    selector: "app-root",
    styleUrl: "./app.css",
    templateUrl: "./app.html",
})
export class App implements OnInit, OnDestroy {
    protected readonly environment = environment;

    private readonly _visibilityListener: () => void;
    protected readonly title = signal("ale-trail");

    constructor(private readonly updateSvc: UpdateService) {
        this._visibilityListener = this.visibilityListener.bind(this);
    }

    private visibilityListener(): void {
        this.updateSvc
            .checkForUpdate()
            .then((b) => {
                if (b) {
                    console.log("Application update available");
                } else {
                    console.log("Application update not available");
                }
            })
            .catch(console.error);
    }

    ngOnInit() {
        document.addEventListener("visibilitychange", this._visibilityListener);
        window.addEventListener("focus", this._visibilityListener);
    }

    ngOnDestroy() {
        document.removeEventListener("visibilitychange", this._visibilityListener);
        window.removeEventListener("focus", this._visibilityListener);
    }

    protected readonly gitCount = gitCount;
}
