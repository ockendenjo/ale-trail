import {Injectable} from "@angular/core";
import {FullStash, Stasher, StasherResponse, StashList, UIShortStash} from "./types";
import {KeyService} from "../key-service/key.service";

@Injectable({providedIn: "root"})
export class StashService {
    constructor(private readonly keySvc: KeyService) {}

    public async listStashes(): Promise<UIShortStash[]> {
        const r = await fetch("/api/stashes");
        const j = (await r.json()) as StashList;
        const datesWithKey = this.keySvc.getDatesWithKey();

        return j.stashes.map((i) => {
            const hasKey = datesWithKey.has(i.date);
            const displayImage = this.getRandomImage();
            const displayClass = this.getRandomClass();
            return {...i, hasKey, displayImage, displayClass};
        });
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

    private getRandomImage(): string {
        const r = Math.random();
        if (r < 0.25) {
            return "bottles.png";
        }
        if (r < 0.5) {
            return "cans.png";
        }
        if (r < 0.75) {
            return "pints.png";
        }
        return "kegs.png";
    }

    private getRandomClass(): string {
        const r = Math.random();
        const d = Math.floor(r * 8);
        return "fill" + String(d);
    }
}
