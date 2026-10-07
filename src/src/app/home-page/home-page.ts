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
    public cards = makeCards();
}

function makeCards(): Card[] {
    const a: Card[] = [];
    const now = Date.now();
    let d = new Date(now - 4 * 24 * 60 * 60 * 1000);

    for (let i = 0; i < 24; i++) {
        d = new Date(d.getTime() + 24 * 60 * 60 * 1000);

        if (d.getTime() < now) {
            a.push({date: d.toISOString().substring(0, 10), locked: false, image: "sample.jpeg"});
            continue;
        }

        if (Math.random() < 0.5) {
            a.push({date: d.toISOString().substring(0, 10), locked: true, releaseType: "RANDOM"});
            continue;
        }

        const h = String(Math.floor(Math.random() * 24)).padStart(2, "0");
        const m = String(Math.floor(Math.random() * 60)).padStart(2, "0");
        a.push({date: d.toISOString().substring(0, 10), locked: true, releaseType: "TIMED", releaseTime: `${h}:${m}`});
    }

    return a;
}
