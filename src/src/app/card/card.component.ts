import {Component, Input} from "@angular/core";
import {Card} from "./card.type";
import {NgClass} from "@angular/common";

@Component({
    imports: [NgClass],
    selector: "app-card",
    styleUrl: "./card.component.css",
    templateUrl: "./card.component.html",
})
export class CardComponent {
    @Input() card!: Card;

    public getDay(c: Card): string {
        const day = c.date.split("-")[2];
        if (day.startsWith("0")) {
            return day.substring(1);
        }
        return day;
    }

    public getMonth(c: Card): string {
        const month = c.date.split("-")[1];
        switch (month) {
            case "11":
                return "NOV";
            case "12":
                return "DEC";
            default:
                return "";
        }
    }
}
