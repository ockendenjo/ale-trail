import {Component} from "@angular/core";
import {CardComponent} from "../card/card.component";
import {Card} from "../card/card.type";

@Component({
    imports: [CardComponent],
    selector: "app-home-page",
    styleUrl: "./home-page.css",
    templateUrl: "./home-page.html",
})
export class HomePage {
    public cards: Card[] = [
        {date: "2026-11-23", locked: false},
        {date: "2026-11-24", locked: false},
        {date: "2026-11-25", locked: true},
        {date: "2026-11-26", locked: true},
        {date: "2026-11-27", locked: true},
        {date: "2026-11-28", locked: true},
        {date: "2026-11-29", locked: true},
        {date: "2026-11-30", locked: true},
        {date: "2026-12-01", locked: true},
        {date: "2026-12-02", locked: true},
        {date: "2026-12-03", locked: true},
        {date: "2026-12-04", locked: true},
        {date: "2026-12-05", locked: true},
        {date: "2026-12-06", locked: true},
        {date: "2026-12-07", locked: true},
        {date: "2026-12-08", locked: true},
        {date: "2026-12-09", locked: true},
        {date: "2026-12-10", locked: true},
        {date: "2026-12-11", locked: true},
        {date: "2026-12-12", locked: true},
        {date: "2026-12-13", locked: true},
        {date: "2026-12-14", locked: true},
        {date: "2026-12-15", locked: true},
        {date: "2026-12-16", locked: true},
    ];
}
