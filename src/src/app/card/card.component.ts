import {Component, Input} from "@angular/core";
import {Card} from "./card.type";
import {getDayStr, getMonthStr} from "../dateUtils";
import {RouterLink} from "@angular/router";

@Component({
    imports: [RouterLink],
    selector: "app-card",
    styleUrl: "./card.component.css",
    templateUrl: "./card.component.html",
})
export class CardComponent {
    @Input() card!: Card;

    public getDay(c: Card): string {
        const day = c.date.split("-")[2];
        return getDayStr(day);
    }

    public getMonth(c: Card): string {
        const month = c.date.split("-")[1];
        return getMonthStr(month);
    }

    public getLink(c: Card): string | undefined {
        if (!c.locked) {
            return c.date;
        }
        if (c.locked && c.hostNeeded) {
            return c.date + "/stasher";
        }
        return undefined;
    }
}
