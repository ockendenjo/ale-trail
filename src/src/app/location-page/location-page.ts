import {Component, OnInit, signal} from "@angular/core";
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {ActivatedRoute} from "@angular/router";
import {take} from "rxjs";
import {NgClass} from "@angular/common";
import {PageState} from "../page-state";

@Component({
    imports: [FormsModule, ReactiveFormsModule, NgClass],
    selector: "app-location-page",
    styleUrl: "./location-page.css",
    templateUrl: "./location-page.html",
})
export class LocationPage implements OnInit {
    public date: string = "";
    private key: string = "";
    public pg = signal<LocationPageState>({state: "LOADING"});

    public form: FormGroup;

    constructor(
        private readonly ar: ActivatedRoute,
        private readonly formBuilder: FormBuilder,
    ) {
        ar.params.pipe(take(1)).subscribe((x) => {
            this.date = x["date"];
        });

        ar.queryParams.pipe(take(1)).subscribe((x) => {
            this.key = x["key"];
        });

        this.form = this.formBuilder.group({
            lat: [undefined, Validators.required],
            lon: [undefined, [Validators.required]],
            w3w: ["", [Validators.required]],
            description: ["", [Validators.required, Validators.maxLength(200)]],
            contents: ["", [Validators.required, Validators.maxLength(200)]],
        });
    }

    ngOnInit() {
        fetch(`/api/stashes/${this.date}/location?key=${this.key}`, {method: "GET"})
            .then((r) => r.json())
            .then((j: StashDetails) => {
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
