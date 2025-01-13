// export const WS_Messages: Map<string, WebSocketMessage> = new Map([
//   ['joinlobby', { type: 'joinlobby', payload: 'Hello World' }],
// ]);

import type { Stance, StanceData } from "./types";

export const STANCES: Stance[] = [
  {
    id: "tennis",
    name: "Tennis",
  },
  // {
  //   id: "volleyball",
  //   name: "Volleyball",
  // },
  {
    id: "football",
    name: "Football",
  },
];

export const GAME_STATE = {
  WAITING: "waiting",
  IN_PROGRESS: "inprogress",
  FINISHED: "finished",
};

export const STANCE_DATA: StanceData = {
  tennis: {
    skills: [
      {
        id: "ace",
        name: "Ace",
        description: "Instantly grants a correct guess. (Requires 1 guess point).",
      },
      {
        id: "fault",
        name: "Fault",
        description: "Makes opponent miss next turn.",
      },
    ],
  },
  // volleyball: {
  //   skills: [
  //     {
  //       id: "libero",
  //       name: "Libero",
  //     },
  //   ],
  // },
  football: {
    skills: [
      {
        id: "tackle",
        name: "Tackle",
        description: "Steals a correct guess from the opponent.",
      },
      {
        id: "goalkeeper",
        name: "Goalkeeper",
        description: "Block the opponent's next correct guess.",
      },
    ],
  },
};
