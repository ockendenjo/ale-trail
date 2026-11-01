import {CanActivateFn, Routes} from "@angular/router";
import {HomePage} from "./home-page/home-page";
import {StasherPage} from "./stasher-page/stasher-page";
import {LocationPage} from "./location-page/location-page";

const hasKey: CanActivateFn = (route) => {
    return Boolean(route.queryParams["key"]);
};

export const routes: Routes = [
    {path: "", component: HomePage, pathMatch: "full"},
    {path: ":date/stasher", component: StasherPage},
    {path: ":date/location", component: LocationPage, canActivate: [hasKey]},
    {path: "**", redirectTo: ""},
];
