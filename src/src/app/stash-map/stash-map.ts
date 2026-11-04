import {AfterViewInit, Component, Input} from "@angular/core";
import {FullUnlockedStash} from "../stash-service/types";
import {Feature, Map as OLMap, View} from "ol";
import TileLayer from "ol/layer/Tile";
import {defaults as defaultControls, ScaleLine} from "ol/control";
import {OSM} from "ol/source";
import VectorSource from "ol/source/Vector";
import VectorLayer from "ol/layer/Vector";
import {fromLonLat} from "ol/proj";
import {Point} from "ol/geom";
import {Circle as CircleStyle, Fill, Stroke, Style} from "ol/style";

@Component({
    imports: [],
    selector: "app-stash-map",
    styleUrl: "./stash-map.css",
    templateUrl: "./stash-map.html",
})
export class StashMap implements AfterViewInit {
    @Input({required: true}) stash!: FullUnlockedStash;

    private map?: OLMap;

    ngAfterViewInit() {
        setTimeout(() => {
            this.setupMap();
        }, 0);
    }

    private setupMap() {
        const osmLayer = new TileLayer({
            source: new OSM(),
            opacity: 0.8,
        });

        const source = new VectorSource({wrapX: false});
        const layer = new VectorLayer({
            source: source,
            visible: true,
        });

        const mapView = new View({maxZoom: 19});
        mapView.setZoom(13);

        this.map = new OLMap({
            controls: defaultControls().extend([new ScaleLine()]),
            target: "stash-map",
            layers: [osmLayer, layer],
            keyboardEventTarget: document,
            view: mapView,
        });

        this.addFeature(source);

        mapView.setCenter(fromLonLat([this.stash.lon, this.stash.lat]));
        mapView.setZoom(16);
    }

    private addFeature(source: VectorSource) {
        const feature = new Feature({
            geometry: new Point(fromLonLat([this.stash.lon, this.stash.lat])),
        });

        feature.setStyle(this.getStyle());
        source.addFeature(feature);
    }

    private getStyle(): Style {
        return new Style({
            image: new CircleStyle({
                radius: 12,
                fill: new Fill({
                    color: "#7b3306",
                }),
                stroke: new Stroke({
                    color: "#ffffff",
                    width: 2,
                }),
            }),
        });
    }
}
