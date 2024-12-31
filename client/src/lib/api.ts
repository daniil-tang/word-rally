import { lobby, player } from "./store";
import type { Player, WebSocketMessage } from "./types";

const BASE_URL = "http://localhost:8080";

const socket = new WebSocket("ws://localhost:8080/ws");

socket.addEventListener("open", function (event) {
  console.log("WebSocket open.");
});

socket.addEventListener("message", function (event) {
  let eventData: WebSocketMessage = JSON.parse(event.data);
  console.log("EVENT DATA", eventData, eventData.event);
  switch (eventData.event) {
    case "lobby":
      lobby.set(JSON.parse(eventData.data));
      break;
    // Should add an error case: Or default = error
  }
});

const sendMessage = (message: WebSocketMessage) => {
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

export async function joinLobby(lobbyID: string, p: Player) {
  sendMessage({
    Event: "joinlobby",
    Data: JSON.stringify({
      lobbyID,
      player: p,
    }),
  });
}
