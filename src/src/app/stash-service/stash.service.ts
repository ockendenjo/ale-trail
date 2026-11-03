import {Service} from "@angular/core";
import {ShortStash, StashList} from "./types";

@Service()
export class StashService {
    public listStashes(): Promise<ShortStash[]> {
        return fetch("/api/stashes")
            .then((r) => r.json())
            .then((j: StashList) => {
                return j.stashes;
            });
    }
}
