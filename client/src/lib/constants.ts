// export const WS_Messages: Map<string, WebSocketMessage> = new Map([
//   ['joinlobby', { type: 'joinlobby', payload: 'Hello World' }],
// ]);

import type { Stance } from "./types";

export const STANCES: Stance[] = [
  {
    id: "tennis",
    name: "Tennis",
  },
  {
    id: "volleyball",
    name: "Volleyball",
  },
  {
    id: "football",
    name: "Football",
  },
];
