import { get } from "svelte/store";
import { lobby, player, websocket } from "./store";
import type { Player, PlayerSettings, WebSocketIncomingMessage, WebSocketOutgoingMessage } from "./types";

const BASE_URL = "http://localhost:8080";

export async function initWS() {
  const socket = new WebSocket("ws://localhost:8080/ws");
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
        lobby.set(JSON.parse(eventData.data));
        break;
      // Should add an error case: Or default = error
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

export async function getPlayerLobby(p: Player) {
  let response = await fetch(`${BASE_URL}/getlobby?playerID=${p.ID}`);
  let lobbyResponse = await response.json();
  lobby.set(lobbyResponse);
}
