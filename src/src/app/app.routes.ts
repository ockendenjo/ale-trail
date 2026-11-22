import {CanActivateFn, Routes} from "@angular/router";
import {HomePage} from "./home-page/home-page";
import {StasherPage} from "./stasher-page/stasher-page";
import {StashPage} from "./stash-page/stash-page";
import {inject} from "@angular/core";
import {KeyService} from "./key-service/key.service";
import {LocationPage} from "./location-page/location-page";
import {PreviousLocationsPage} from "./previous-locations-page/previous-locations-page.component";

const hasKey: CanActivateFn = (route) => {
    const date = route.params["date"] || "";
    if (!date) {
        return false;
    }

    const keySvc = inject(KeyService);
    const qKey = route.queryParams["key"];
    if (qKey) {
        keySvc.setKey(date, qKey);
        return true;
    }

    return Boolean(keySvc.getKey(date));
};

const interceptKey: CanActivateFn = (route) => {
    const date = route.params["date"] || "";
    if (!date) {
        return false;
    }

    const qKey = route.queryParams["key"];
    if (qKey) {
        const keySvc = inject(KeyService);
        keySvc.setKey(date, qKey);
    }
    return true;
};

export const routes: Routes = [
    {path: "", component: HomePage, pathMatch: "full"},
    {path: "previous-locations", component: PreviousLocationsPage},
    {path: ":date", component: StashPage, pathMatch: "full"},
    {path: ":date/stasher", component: StasherPage, canActivate: [interceptKey]},
    {path: ":date/location", component: LocationPage, canActivate: [hasKey]},
    {path: "**", redirectTo: ""},
];
