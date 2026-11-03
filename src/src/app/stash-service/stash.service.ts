import {Service} from "@angular/core";
import {FullStash, ShortStash, StashList} from "./types";

@Service()
export class StashService {
    public async listStashes(): Promise<ShortStash[]> {
        const r = await fetch("/api/stashes");
        const j = (await r.json()) as StashList;
        return j.stashes;
    }

    public async getStash(date: string): Promise<FullStash> {
        const r = await fetch(`/api/stashes/${date}`);
        return await r.json();
    }
}
