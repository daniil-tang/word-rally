<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { player } from "$lib/store";
  import type { Player } from "$lib/types";
  import { createPlayer, registerConnection } from "$lib/api";

  let isNewPlayer = true;
  let _player: Player = { Name: "", ID: "" };

  onMount(async () => {
    let storedPlayer = JSON.parse(localStorage.getItem("word-rally-player-test") ?? "{}");
    if (storedPlayer?.ID) {
      _player = storedPlayer;
      await savePlayer();
      isNewPlayer = false;
    }
  });

  async function savePlayer() {
    await createPlayer(_player);
    localStorage.setItem("word-rally-player-test", JSON.stringify($player));
    registerConnection($player);
  }

  async function handleSubmit() {
    await savePlayer();
    isNewPlayer = false;

    startGame();
  }

  function startGame() {
    goto("/lobby");
  }
</script>

<main>
  <div class="container">
    {#if isNewPlayer}
      <form on:submit|preventDefault={handleSubmit} class="nes-container with-title">
        <h1>Welcome to Word Rally!</h1>
        <div class="nes-field">
          <label for="name_field">Your Name</label>
          <input type="text" id="name_field" class="nes-input" bind:value={_player.Name} required />
        </div>
        <button type="submit" class="nes-btn is-primary">Start Game</button>
      </form>
    {:else}
      <h2 class="title">Welcome back, {$player.Name}!</h2>
      <button on:click|preventDefault={startGame} class="nes-btn is-primary">Start Game</button>
    {/if}
  </div>
</main>

<style>
  main {
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
    font-size: 2rem;
  }

  button {
    margin-top: 1rem;
  }
</style>
