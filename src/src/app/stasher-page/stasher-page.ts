import {Component, signal} from "@angular/core";
import {ActivatedRoute} from "@angular/router";
import {take} from "rxjs";
import {getDayStr, getMonthStr} from "../dateUtils";
import {FormBuilder, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";
import {NgClass} from "@angular/common";

@Component({
    imports: [ReactiveFormsModule, NgClass],
    selector: "app-stasher-page",
    styleUrl: "./stasher-page.css",
    templateUrl: "./stasher-page.html",
})
export class StasherPage {
    private date: string = "";
    public st = signal<SaveState>("READY");

    public form: FormGroup;

    constructor(
        private readonly ar: ActivatedRoute,
        private readonly formBuilder: FormBuilder,
    ) {
        ar.params.pipe(take(1)).subscribe((x) => {
            this.date = x["date"];
        });

        this.form = this.formBuilder.group({
            name: ["", Validators.required],
            email: ["", [Validators.required, Validators.email]],
        });
    }

    public get day(): string {
        const day = this.date.split("-")[2];
        return getDayStr(day);
    }

    public get month(): string {
        const month = this.date.split("-")[1];
        return getMonthStr(month);
    }

    public get before(): string {
        const beforeStr = new Date(new Date(this.date).getTime() - 24 * 60 * 60 * 1000).toISOString().substring(0, 10);
        const parts = beforeStr.split("-");
        const day = parts[2];
        const month = parts[1];
        return getDayStr(day) + " " + getMonthStr(month);
    }

    public submit(): void {
        const value = this.form.value;

        this.st.set("SAVING");
        fetch(`/api/stashes/${this.date}/stasher`, {method: "POST", body: JSON.stringify(value)})
            .then((r) => {
                switch (r.status) {
                    case 204:
                        this.st.set("SAVED");
                        break;
                    case 409:
                        this.st.set("CONFLICT");
                        break;
                    default:
                        this.st.set("ERROR");
                }
            })
            .catch(() => {
                this.st.set("ERROR");
            });
    }
}

type SaveState = "READY" | "SAVING" | "SAVED" | "CONFLICT" | "ERROR";
