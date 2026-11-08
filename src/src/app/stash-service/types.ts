interface BaseStash {
    date: string;
    hasKey: boolean;
    locked: boolean;
}

interface RandomTimedStash extends BaseStash {
    locked: true;
    releaseType: "RANDOM";
    hostNeeded?: true;
}

interface TimedStash extends BaseStash {
    locked: true;
    releaseType: "TIMED";
    releaseTime?: string;
    hostNeeded?: true;
}

interface UnlockedStash extends BaseStash {
    locked: false;
    image: string;
}

export type ShortStash = UnlockedStash | RandomTimedStash | TimedStash;

export type UIShortStash = ShortStash & {
    hasKey: boolean;
    displayImage: string;
    displayClass: string;
};

export interface StashList {
    stashes: ShortStash[];
}

export interface FullUnlockedStash extends BaseStash {
    locked: false;
    image: string;
    lat: number;
    lon: number;
    w3w: string;
    description: string;
    contents: string;
}

export type FullStash = FullUnlockedStash | RandomTimedStash | TimedStash;

export interface Stasher {
    name: string;
    email: string;
}

export interface StasherResponse {
    date: string;
    key: string;
}
