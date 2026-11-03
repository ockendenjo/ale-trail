import {Component, Input} from "@angular/core";
import {getDayStr, getMonthStr} from "../dateUtils";
import {RouterLink} from "@angular/router";
import {ShortStash} from "../stash-service/types";

@Component({
    imports: [RouterLink],
    selector: "app-card",
    styleUrl: "./card.component.css",
    templateUrl: "./card.component.html",
})
export class CardComponent {
    @Input() card!: ShortStash;

    public getDay(c: ShortStash): string {
        const day = c.date.split("-")[2];
        return getDayStr(day);
    }

    public getMonth(c: ShortStash): string {
        const month = c.date.split("-")[1];
        return getMonthStr(month);
    }

    public getLink(c: ShortStash): string | undefined {
        if (!c.locked) {
            return c.date;
        }
        if (c.locked && c.hostNeeded) {
            return c.date + "/stasher";
        }
        return undefined;
    }
}
