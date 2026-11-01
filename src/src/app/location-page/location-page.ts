import {Component} from "@angular/core";
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from "@angular/forms";
import {ActivatedRoute} from "@angular/router";
import {take} from "rxjs";
import {NgClass} from "@angular/common";

@Component({
    imports: [FormsModule, ReactiveFormsModule, NgClass],
    selector: "app-location-page",
    styleUrl: "./location-page.css",
    templateUrl: "./location-page.html",
})
export class LocationPage {
    public date: string = "";
    private key: string = "";

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

    public submit(): void {}
}
