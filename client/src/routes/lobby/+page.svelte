<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { lobby, player } from "$lib/store";
  import { createLobby, joinLobby } from "$lib/api";

  let lobbyId = "";

  // onMount(() => {
  //   if (!$player?.ID) goto("/");
  //   if ($lobby) goto(`/lobby/${$lobby.ID}`);
  // });

  $: {
    if (!$player?.ID) goto("/");
    if ($lobby) goto(`/lobby/${$lobby.ID}`);
  }

  async function handleJoinLobby() {
    // TODO: Add API call to join lobby
    // For now, just navigate to the game page
    await joinLobby(lobbyId, $player);
    // console.log("JOIN LOBBY?", $lobby);
    // if ($lobby) goto(`/lobby/${$lobby.ID}`);
  }

  async function handleCreateLobby() {
    await createLobby($player);
    // console.log("RES", $lobby);
    // if ($lobby) goto(`/lobby/${$lobby.ID}`);
  }
</script>

<div class="inline-flex">
  <div class="container">
    <h1>Join a lobby</h1>
    <div class="nes-field">
      <input placeholder="Lobby ID" type="text" id="lobby_id" class="nes-input" bind:value={lobbyId} />
    </div>
    <button class="nes-btn is-primary" on:click={handleJoinLobby}>Join</button>
  </div>

  <div class="vertical-line"></div>

  <div class="container flex-vertical-align">
    <button class="nes-btn is-primary" on:click={handleCreateLobby}>Create Lobby</button>
  </div>
</div>

<style>
  :global(body) {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .container {
    max-width: 600px;
    padding: 1rem;
    text-align: center;
  }

  h1 {
    margin-bottom: 2rem;
    font-size: 1.5rem;
  }

  .nes-field {
    margin-bottom: 1rem;
  }

  button {
    margin-top: 1rem;
  }

  .vertical-line {
    border-left: 4px solid black;
    /* height: 500px; */
  }

  .inline-flex {
    display: inline-flex;
  }

  .flex-vertical-align {
    display: flex;
    align-items: center;
  }
</style>
