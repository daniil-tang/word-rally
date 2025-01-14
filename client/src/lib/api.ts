import { get } from "svelte/store";
import { actionLog, lobby, player, websocket } from "./store";
import type { Player, PlayerSettings, WebSocketIncomingMessage, WebSocketOutgoingMessage } from "./types";
const BASE_URL = import.meta.env.VITE_SERVER_URL || "http://localhost:8080";
const WS_URL = import.meta.env.VITE_WEBSOCKET_URL || "ws://localhost:8080/ws";
// const BASE_URL = "http://localhost:8080";
console.log("ENV VAR", BASE_URL);
export async function initWS() {
  // const socket = new WebSocket("ws://localhost:8080/ws");
  const socket = new WebSocket(WS_URL);
  socket.addEventListener("open", function (event) {
    console.log("WebSocket open.");
    websocket.set(socket);
  });

  socket.addEventListener("close", function (event) {
    console.log("WebSocket close.");
    reconnect();
  });

  socket.addEventListener("message", function (event) {
    let eventData: WebSocketIncomingMessage = JSON.parse(event.data);
    console.log("EVENT DATA", eventData, eventData.event);
    switch (eventData.event) {
      case "lobby":
        let lob = JSON.parse(eventData.data)?.ID ? JSON.parse(eventData.data) : null;
        lobby.set(lob);
        console.log("RALLY", JSON.parse(eventData.data)?.Game?.Rally?.StatusEffects);
        break;
      // Should add an error case: Or default = error
      case "actionlog":
        console.log("ACTION LOG!", event);
        // let data = JSON.parse(eventData.data);
        actionLog.update((log) => {
          return [...log, eventData.data.message];
        });
        break;
    }
  });
}

function reconnect() {
  setTimeout(() => {
    console.log("Reconnecting to WebSocket...");
    initWS(); // Try to reinitialize the WebSocket connection
  }, 1000); // Reconnect after 1 second
}

const sendMessage = (message: WebSocketOutgoingMessage) => {
  const ws = get(websocket);
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(message));
  } else {
    console.warn("WebSocket is not open, cannot send message.");
  }
};

export async function createPlayer(p: Player) {
  let response = await fetch(`${BASE_URL}/createplayer`, { method: "POST", body: JSON.stringify(p) });
  let playerResponse = await response.json();
  player.set(playerResponse);
}

export async function createLobby(hostPlayer: Player) {
  let response = await fetch(`${BASE_URL}/createlobby`, {
    method: "POST",
    body: JSON.stringify(hostPlayer),
  });
  let lobbyResponse = await response.json();
  console.log("RESPONSI", lobbyResponse);
  lobby.set(lobbyResponse);
}

export async function registerConnection(p: Player) {
  sendMessage({
    Event: "registerconnection",
    Data: JSON.stringify({
      player: p,
    }),
  });
}

export async function joinLobby(lobbyID: string, p: Player) {
  sendMessage({
    Event: "joinlobby",
    Data: JSON.stringify({
      lobbyID,
      player: p,
    }),
  });
}

export async function leaveLobby(lobbyID: string, p: Player) {
  sendMessage({
    Event: "leavelobby",
    Data: JSON.stringify({
      lobbyID,
      player: p,
    }),
  });
}

export async function updatePlayerSettings(lobbyID: string, p: Player, playerSettings: PlayerSettings) {
  sendMessage({
    Event: "updateplayersettings",
    Data: JSON.stringify({
      lobbyID,
      player: p,
      playerSettings,
    }),
  });
}

// Create this on lobby mount..?
export async function createGame(lobbyID: string, p: Player) {
  sendMessage({
    Event: "creategame",
    Data: JSON.stringify({
      lobbyID,
      player: p,
    }),
  });
}

export async function startGame(lobbyID: string, p: Player) {
  sendMessage({
    Event: "startgame",
    Data: JSON.stringify({
      lobbyID,
      player: p,
    }),
  });
}

export async function guess(lobbyID: string, p: Player, guess: string) {
  console.log("GUESS", guess.toUpperCase().charCodeAt(0));
  sendMessage({
    Event: "playeraction",
    Data: JSON.stringify({
      lobbyID,
      player: p,
      action: "guess",
      actionDetails: {
        guessedLetters: [guess.toUpperCase().charCodeAt(0)],
      },
    }),
  });
}

export async function endTurn(lobbyID: string, p: Player) {
  sendMessage({
    Event: "playeraction",
    Data: JSON.stringify({
      lobbyID,
      player: p,
      action: "endturn",
      actionDetails: {},
    }),
  });
}

export async function useSkill(lobbyID: string, p: Player, skillId: string) {
  sendMessage({
    Event: "playeraction",
    Data: JSON.stringify({
      lobbyID,
      player: p,
      action: "useskill",
      actionDetails: {
        skillUsed: skillId,
      },
    }),
  });
}

export async function getPlayerLobby(p: Player) {
  let response = await fetch(`${BASE_URL}/getlobby?playerID=${p.ID}`);
  let lobbyResponse = await response.json();
  lobby.set(lobbyResponse);
}
