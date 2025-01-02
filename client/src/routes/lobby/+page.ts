import { getPlayerLobby } from "$lib/api";
import { player } from "$lib/store";
import { get } from "svelte/store";
export const ssr = false;

export const load = async () => {
  await getPlayerLobby(get(player));
};
