import {Component, Signal, signal} from "@angular/core";
import {ActivatedRoute} from "@angular/router";
import {take} from "rxjs";
import {getDayStr, getMonthStr} from "../dateUtils";

@Component({
    imports: [],
    selector: "app-stasher-page",
    styleUrl: "./stasher-page.css",
    templateUrl: "./stasher-page.html",
})
export class StasherPage {
    private date: string = "";
    public st: Signal<SaveState> = signal("READY");

    constructor(private readonly ar: ActivatedRoute) {
        ar.params.pipe(take(1)).subscribe((x) => {
            this.date = x["date"];
        });
    }

    public get day(): string {
        const day = this.date.split("-")[2];
        return getDayStr(day);
    }

    public get month(): string {
        const month = this.date.split("-")[1];
        return getMonthStr(month);
    }

    public get before(): string {
        const beforeStr = new Date(new Date(this.date).getTime() - 24 * 60 * 60 * 1000).toISOString().substring(0, 10);
        const parts = beforeStr.split("-");
        const day = parts[2];
        const month = parts[1];
        return getDayStr(day) + " " + getMonthStr(month);
    }
}

type SaveState = "READY" | "SAVING" | "SAVED" | "CONFLICT" | "ERROR";
