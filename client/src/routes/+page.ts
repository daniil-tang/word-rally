import { initWS } from "$lib/api";

export const ssr = false;

export const load = async () => {
  initWS();
};
