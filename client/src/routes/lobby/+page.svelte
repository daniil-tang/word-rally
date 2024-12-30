<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { player } from "$lib/store";

  let lobbyId = "";

  onMount(() => {
    console.log("PLAYER", $player);
    if (!$player?.ID) goto("/");
  });

  async function joinLobby() {
    // TODO: Add API call to join lobby
    // For now, just navigate to the game page
    await goto(`/game/${lobbyId}`);
  }
</script>

<div class="container">
  <h1>Join a lobby</h1>
  <div class="nes-field">
    <input placeholder="Lobby ID" type="text" id="lobby_id" class="nes-input" bind:value={lobbyId} />
  </div>
  <button class="nes-btn is-primary" on:click={joinLobby}>Join</button>
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
</style>
