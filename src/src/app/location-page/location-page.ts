import {Component, OnInit, signal} from "@angular/core";
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {ActivatedRoute} from "@angular/router";
import {take} from "rxjs";
import {NgClass} from "@angular/common";
import {PageState} from "../page-state";
import {KeyService} from "../key-service/key.service";
import {LiveStash} from "../live-stash/live-stash";
import {FullUnlockedStash} from "../stash-service/types";

@Component({
    imports: [FormsModule, ReactiveFormsModule, NgClass, LiveStash],
    selector: "app-location-page",
    styleUrl: "./location-page.css",
    templateUrl: "./location-page.html",
})
export class LocationPage implements OnInit {
    public date: string = "";
    private key: string = "";
    public pg = signal<LocationPageState>({state: "LOADING"});

    public form: FormGroup;
    public coords = "";

    constructor(
        ar: ActivatedRoute,
        private readonly formBuilder: FormBuilder,
        keySvc: KeyService,
    ) {
        ar.params.pipe(take(1)).subscribe((x) => {
            this.date = x["date"];
            this.key = keySvc.getKey(this.date);
        });

        this.form = this.formBuilder.group({
            lat: [undefined, Validators.required],
            lon: [undefined, [Validators.required]],
            w3w: ["", [Validators.required]],
            description: ["", [Validators.required, Validators.maxLength(200)]],
            contents: ["", [ Validators.maxLength(200)]],
        });
    }

    ngOnInit() {
        fetch(`/api/stashes/${this.date}/location?key=${this.key}`, {method: "GET"})
            .then((r) => r.json())
            .then((j: StashDetails) => {
                if (j.lat === 0) {
                    j.lat = null as unknown as number;
                }
                if (j.lon === 0) {
                    j.lon = null as unknown as number;
                }
                this.form.setValue(j);
                this.pg.set({state: "LOADED", data: undefined});
            })
            .catch((e) => {
                this.pg.set({state: "ERROR", error: e});
            });
    }

    public submit(): void {
        this.pg.set({state: "SAVING"});
        fetch(`/api/stashes/${this.date}/location?key=${this.key}`, {
            method: "PUT",
            body: JSON.stringify(this.form.value),
        })
            .then((r) => {
                if (r.status === 204) {
                    this.pg.set({state: "SAVED"});
                } else {
                    this.pg.set({state: "ERROR", error: new Error("error saving")});
                }
            })
            .catch((e) => {
                this.pg.set({state: "ERROR", error: e});
            });
    }

    public editAgain(): void {
        this.pg.set({state: "LOADED", data: undefined});
    }

    public getStash(): FullUnlockedStash | undefined {
        return {
            ...this.form.value,
            date: this.date,
            hasKey: false,
            image: "",
            locked: false,
        };
    }

    public updateLatLon(): void {
        const s = this.coords.trim();
        const parts = s.split(" ");
        if (parts.length != 4) {
            return;
        }

        const p1 = parts[0] + " " + parts[1];
        const p2 = parts[2] + " " + parts[3];

        try {
            const c1 = parseDm(p1);
            const c2 = parseDm(p2);

            this.form.get("lat")?.setValue(c1);
            this.form.get("lon")?.setValue(c2);
            this.coords = "";
        } catch {
            return;
        }
    }
}

type LocationPageState = PageState<void> | StateSaving | StateSaved;

interface StateSaving {
    state: "SAVING";
}

interface StateSaved {
    state: "SAVED";
}

interface StashDetails {
    lat: number;
    lon: number;
    w3w: string;
    description: string;
    contents: string;
}

function parseDm(input: string): number {
    const match = input.trim().match(/^([NSWE])\s*(\d{1,3})[°\s]+(\d+(?:\.\d+)?)['′\s]*$/);

    if (!match) throw new Error(`Invalid coordinate: ${input}`);

    const [, hemisphere, degrees, minutes] = match;
    const sign = "SW".includes(hemisphere) ? -1 : 1;

    const d = sign * (Number(degrees) + Number(minutes) / 60);
    const n = Math.pow(10, 6);
    return Math.floor(d * n) / n;
}
