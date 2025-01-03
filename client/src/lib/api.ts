import { lobby, player } from "./store";
import type { Player, PlayerSettings, WebSocketIncomingMessage, WebSocketOutgoingMessage } from "./types";

const BASE_URL = "http://localhost:8080";

const socket = new WebSocket("ws://localhost:8080/ws");
console.log("SOCKET ME");
socket.addEventListener("open", function (event) {
  console.log("WebSocket open.");
});

socket.addEventListener("close", function (event) {
  console.log("WebSocket close.");
});

socket.addEventListener("message", function (event) {
  let eventData: WebSocketIncomingMessage = JSON.parse(event.data);
  console.log("EVENT DATA", eventData, eventData.event);
  switch (eventData.event) {
    case "lobby":
      console.log("WHATCHAMACALIT", JSON.parse(eventData.data));
      lobby.set(JSON.parse(eventData.data));
      break;
    // Should add an error case: Or default = error
  }
});

const sendMessage = (message: WebSocketOutgoingMessage) => {
  if (socket.readyState <= 1) {
    socket.send(JSON.stringify(message));
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

export async function getPlayerLobby(p: Player) {
  let response = await fetch(`${BASE_URL}/getlobby?playerID=${p.ID}`);
  let lobbyResponse = await response.json();
  lobby.set(lobbyResponse);
}
