import {Injectable} from "@angular/core";
import {environment} from "../../environments/environment.dev";

@Injectable({providedIn: "root"})
export class KeyService {
    public setKey(date: string, stashKey: string) {
        if (!date || !stashKey) {
            return;
        }
        const record = this.readValue();
        record[date] = stashKey;
        this.writeValue(record);
    }

    public removeKey(date: string) {
        const record = this.readValue();
        delete record[date];
        this.writeValue(record);
    }

    public getKey(date: string): string {
        const record = this.readValue();
        return record[date];
    }

    private readValue(): Record<string, string> {
        const rawValue = localStorage.getItem(environment.storageKey) || "{}";
        return JSON.parse(rawValue) || {};
    }

    private writeValue(v: Record<string, string>) {
        const json = JSON.stringify(v);
        localStorage.setItem(environment.storageKey, json);
    }
}
