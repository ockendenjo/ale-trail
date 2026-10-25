import {Component, OnInit, signal} from "@angular/core";
import {CardComponent} from "../card/card.component";
import {Card} from "../card/card.type";
import {PageState} from "../page-state";

@Component({
    imports: [CardComponent],
    selector: "app-home-page",
    styleUrl: "./home-page.css",
    templateUrl: "./home-page.html",
})
export class HomePage implements OnInit {
    public st = signal<PageState<Card[]>>({state: "LOADING"});

    ngOnInit() {
        fetch("/api/cards")
            .then((res) => res.json())
            .then((res: CardsResponse) => {
                this.st.set({state: "LOADED", data: res.cards});
            })
            .catch((error) => {
                this.st.set({state: "ERROR", error: error});
            });
    }
}

interface CardsResponse {
    cards: Card[];
}
