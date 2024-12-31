export type Player = {
  Name: string;
  ID: string;
};

export type StanceType = string;

export type GameState = "waiting" | "inprogress" | "finished";

export interface PlayerSettings {
  Stance: StanceType;
}

export interface GameSettings {}

export interface TurnActionPoints {
  Guess: number;
  Skill: number;
}

export interface Rally {
  Turn: number;
  TurnActionPoints: { [playerId: string]: TurnActionPoints };
  Guesses: { [playerId: string]: string[] };
  Word: string;
}

export interface Game {
  ID: string;
  State: GameState;
  Score: { [playerId: string]: number };
  CurrentServer: number;
  Rally: Rally | null;
  Settings: GameSettings | null;
  PlayerCooldowns: { [playerId: string]: { [action: string]: number } };
}

export interface Player {
  ID: string;
  Name: string;
}

export interface Lobby {
  ID: string;
  Players: Player[];
  Game: Game | null;
  Host: string;
  MaxPlayers: number;
  PlayerSettings: { [playerId: string]: PlayerSettings };
}

export type WebSocketMessage = {
  //Why is this lowercase?
  event: string;
  data: string;
};
