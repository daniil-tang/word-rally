<script lang="ts">
  import { goto } from "$app/navigation";
  import { STANCES } from "$lib/constants";
  import { player, lobby } from "$lib/store";
  import { onMount } from "svelte";

  // Refactor to fetch lobby stuff on every load
  onMount(() => {
    if (!$player?.ID) goto("/");
  });
  let opponent = $lobby.Players.find((p) => p.ID != $player.ID);
</script>

<div>
  <h1>Game Lobby</h1>
  <p>Code: {$lobby.ID}</p>
  <div class="inline-flex">
    <div class="nes-container player-container">
      <h3>{$player.Name}</h3>
      Select your stance:
      {#each STANCES as stance, i}
        <label>
          <input
            type="radio"
            class="nes-radio"
            name="answer"
            checked={stance.id == $lobby.PlayerSettings[$player.ID].Stance}
          />
          <span>{stance.name}</span>
        </label>
      {/each}
      <button class="nes-btn ready-button">Ready</button>
    </div>
    <div class="nes-container player-container"><h3>{opponent?.Name}</h3></div>
  </div>
  {#if $player.ID == $lobby.Host}
    <button type="button" class={"nes-btn is-primary start-button"}>Start Game</button>
  {/if}
</div>

<style>
  .inline-flex {
    display: inline-flex;
  }
  .player-container {
    margin: 10px;
    width: 400px;
  }
  .ready-button {
    display: block;
    margin: auto;
    margin-top: 10px;
  }
  .start-button {
    display: block;
    margin: auto;
    margin-top: 10px;
  }
</style>
