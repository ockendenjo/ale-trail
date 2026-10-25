import {Routes} from "@angular/router";
import {HomePage} from "./home-page/home-page";
import {StasherPage} from "./stasher-page/stasher-page";

export const routes: Routes = [
    {path: "", component: HomePage},
    {path: ":date/stasher", component: StasherPage},
];
