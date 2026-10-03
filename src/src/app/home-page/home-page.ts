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
        {date: "2026-11-16", locked: false, releaseType: "RANDOM"},
        {date: "2026-11-17", locked: false, releaseType: "TIMED", releaseTime: ""},
        {date: "2026-11-18", locked: false, releaseType: "RANDOM"},
        {date: "2026-11-19", locked: false, releaseType: "RANDOM"},
        {date: "2026-11-20", locked: false, releaseType: "RANDOM"},
        {date: "2026-11-21", locked: false, releaseType: "RANDOM"},
        {date: "2026-11-22", locked: false, releaseType: "RANDOM"},
        {date: "2026-11-23", locked: false, releaseType: "RANDOM"},
        {date: "2026-11-24", locked: false, releaseType: "RANDOM"},
        {date: "2026-11-25", locked: true, releaseType: "RANDOM"},
        {date: "2026-11-26", locked: true, releaseType: "TIMED", releaseTime: "12:00"},
        {date: "2026-11-27", locked: true, releaseType: "RANDOM"},
        {date: "2026-11-28", locked: true, releaseType: "RANDOM"},
        {date: "2026-11-29", locked: true, releaseType: "TIMED", releaseTime: "04:30"},
        {date: "2026-11-30", locked: true, releaseType: "RANDOM"},
        {date: "2026-12-01", locked: true, releaseType: "RANDOM"},
        {date: "2026-12-02", locked: true, releaseType: "RANDOM"},
        {date: "2026-12-03", locked: true, releaseType: "TIMED", releaseTime: "16:12"},
        {date: "2026-12-04", locked: true, releaseType: "RANDOM"},
        {date: "2026-12-05", locked: true, releaseType: "RANDOM"},
        {date: "2026-12-06", locked: true, releaseType: "RANDOM"},
        {date: "2026-12-07", locked: true, releaseType: "TIMED", releaseTime: "09:40"},
        {date: "2026-12-08", locked: true, releaseType: "RANDOM"},
        {date: "2026-12-09", locked: true, releaseType: "RANDOM"},
    ];
}
