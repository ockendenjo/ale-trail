import {Injectable} from "@angular/core";
import {FullStash, ShortStash, Stasher, StasherResponse, StashList} from "./types";
import {KeyService} from "../key-service/key.service";

@Injectable({providedIn: "root"})
export class StashService {
    constructor(private readonly keySvc: KeyService) {}

    public async listStashes(): Promise<ShortStash[]> {
        const r = await fetch("/api/stashes");
        const j = (await r.json()) as StashList;
        const datesWithKey = this.keySvc.getDatesWithKey();

        j.stashes.forEach((s) => {
            s.hasKey = datesWithKey.has(s.date);
        });
        return j.stashes;
    }

    public async getStash(date: string): Promise<FullStash> {
        const r = await fetch(`/api/stashes/${date}`);
        return await r.json();
    }

    public setStasher(date: string, stasher: Stasher) {
        return fetch(`/api/stashes/${date}/stasher`, {method: "POST", body: JSON.stringify(stasher)}).then(
            async (r) => {
                if (r.status === 200) {
                    const res = (await r.json()) as StasherResponse;
                    this.keySvc.setKey(res.date, res.key);
                }
                return r;
            },
        );
    }
}
