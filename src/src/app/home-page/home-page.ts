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
        {date: "2026-11-23"},
        {date: "2026-11-24"},
        {date: "2026-11-25"},
        {date: "2026-11-26"},
        {date: "2026-11-27"},
        {date: "2026-11-28"},
        {date: "2026-11-29"},
        {date: "2026-11-30"},
        {date: "2026-12-01"},
        {date: "2026-12-02"},
        {date: "2026-12-03"},
        {date: "2026-12-04"},
        {date: "2026-12-05"},
        {date: "2026-12-06"},
        {date: "2026-12-07"},
        {date: "2026-12-08"},
        {date: "2026-12-09"},
        {date: "2026-12-10"},
        {date: "2026-12-11"},
        {date: "2026-12-12"},
        {date: "2026-12-13"},
        {date: "2026-12-14"},
        {date: "2026-12-15"},
        {date: "2026-12-16"},
    ];
}
