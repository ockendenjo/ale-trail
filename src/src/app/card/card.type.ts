interface BaseCard {
    date: string;
    locked: boolean;
}

interface RandomCard extends BaseCard {
    releaseType: "RANDOM";
}

interface TimedCard extends BaseCard {
    releaseType: "TIMED";
    releaseTime: string;
}

export type Card = RandomCard | TimedCard;
