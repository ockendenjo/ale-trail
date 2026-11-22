import {AfterViewInit, Component, Input} from "@angular/core";
import {Feature, Map as OLMap, View} from "ol";
import TileLayer from "ol/layer/Tile";
import {defaults as defaultControls, ScaleLine} from "ol/control";
import {OSM} from "ol/source";
import VectorSource from "ol/source/Vector";
import VectorLayer from "ol/layer/Vector";
import {fromLonLat} from "ol/proj";
import {Point} from "ol/geom";
import {Circle as CircleStyle, Fill, Stroke, Style} from "ol/style";
import {Extent} from "ol/extent";

@Component({
    imports: [],
    selector: "multi-stash-map",
    styleUrl: "./multi-stash-map.css",
    templateUrl: "./multi-stash-map.html",
})
export class MultiStashMap implements AfterViewInit {
    @Input({required: true}) locations: [number, number][] = [];

    private features: Feature[] = [];
    private mapView?: View;

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
        this.mapView = mapView;

        new OLMap({
            controls: defaultControls().extend([new ScaleLine()]),
            target: "stash-map",
            layers: [osmLayer, layer],
            keyboardEventTarget: document,
            view: mapView,
        });

        this.features = this.addFeatures(source);
        const ext = source.getExtent();
        mapView.fit(ext as Extent, {padding: [20, 20, 20, 20]});
    }

    private addFeatures(source: VectorSource): Feature[] {
        return this.locations.map((lonLat) => {
            const feature = new Feature({
                geometry: new Point(fromLonLat(lonLat)),
            });
            feature.setStyle(this.getStyle());
            source.addFeature(feature);
            return feature;
        });
    }

    private style?: Style;

    private getStyle(): Style {
        if (!this.style) {
            this.style = new Style({
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
        return this.style;
    }
}
