import {Component, signal} from "@angular/core";
import {ActivatedRoute, Router, RouterLink} from "@angular/router";
import {take} from "rxjs";
import {getDayStr, getMonthStr} from "../dateUtils";
import {FormBuilder, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";
import {NgClass} from "@angular/common";
import {StashService} from "../stash-service/stash.service";
import {KeyService} from "../key-service/key.service";
import {Contents} from "../contents/contents";
import {StashLocation} from "../stash-location/stash-location";

@Component({
    imports: [ReactiveFormsModule, NgClass, Contents, StashLocation, RouterLink],
    selector: "app-stasher-page",
    styleUrl: "./stasher-page.css",
    templateUrl: "./stasher-page.html",
})
export class StasherPage {
    private date: string = "";
    private key?: string;
    public st = signal<SaveState>("READY");

    public form: FormGroup;

    constructor(
        ar: ActivatedRoute,
        private readonly formBuilder: FormBuilder,
        private readonly stashSvc: StashService,
        private readonly keySvc: KeyService,
        private readonly router: Router,
    ) {
        ar.params.pipe(take(1)).subscribe((x) => {
            this.date = x["date"];
            const key = keySvc.getKey(this.date);
            if (key) {
                this.key = key;
                //TODO: Check key
                this.st.set("LOADED");
            }
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
        this.stashSvc
            .setStasher(this.date, value)
            .then((r) => {
                switch (r.status) {
                    case 200:
                    case 204:
                        this.st.set("LOADED");
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

    public deleteStasher(): void {
        this.st.set("RELEASING");
        const key = this.keySvc.getKey(this.date);

        this.stashSvc
            .deleteStasher(this.date, key)
            .then(() => {
                this.router.navigateByUrl("");
            })
            .catch(() => {
                this.st.set("RELERROR");
            });
    }
}

type SaveState = "READY" | "SAVING" | "CONFLICT" | "ERROR" | "LOADING" | "LOADED" | "RELEASING" | "RELERROR";
