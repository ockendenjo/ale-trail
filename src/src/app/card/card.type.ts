interface BaseCard {
    date: string;
    locked: boolean;
}

interface RandomCard extends BaseCard {
    locked: true;
    releaseType: "RANDOM";
}

interface TimedCard extends BaseCard {
    locked: true;
    releaseType: "TIMED";
    releaseTime: string;
}

interface UnlockedCard extends BaseCard {
    locked: false;
    image: string;
}

export type Card = UnlockedCard | RandomCard | TimedCard;
