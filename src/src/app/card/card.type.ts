interface BaseCard {
    date: string;
    locked: boolean;
}

interface RandomCard extends BaseCard {
    locked: true;
    releaseType: "RANDOM";
    hostNeeded?: true;
}

interface TimedCard extends BaseCard {
    locked: true;
    releaseType: "TIMED";
    releaseTime?: string;
    hostNeeded?: true;
}

interface UnlockedCard extends BaseCard {
    locked: false;
    image: string;
}

export type Card = UnlockedCard | RandomCard | TimedCard;
