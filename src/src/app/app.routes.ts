import {Routes} from "@angular/router";
import {HomePage} from "./home-page/home-page";
import {StasherPage} from "./stasher-page/stasher-page";
import {LocationPage} from "./location-page/location-page";

export const routes: Routes = [
    {path: "", component: HomePage},
    {path: ":date/stasher", component: StasherPage},
    {path: ":date/location", component: LocationPage},
    {path: "**", redirectTo: ""},
];
