import { writable } from "svelte/store";
import type { Player, Lobby } from "./types";

export const player = writable<Player>();

export const lobby = writable<Lobby>();
