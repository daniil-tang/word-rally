<script lang="ts">
  import { goto } from "$app/navigation";
  import { updatePlayerSettings } from "$lib/api";
  import { STANCES } from "$lib/constants";
  import { player, lobby } from "$lib/store";
  import { onMount } from "svelte";
  import { get } from "svelte/store";

  // Refactor to fetch lobby stuff on every load
  $: selectedStance = $lobby.PlayerSettings[$player.ID].Stance;
  let playerReady = $lobby.PlayerSettings[$player.ID].Ready;

  $: {
    let _lobby = get(lobby);
    let _player = get(player);
    if (selectedStance != _lobby.PlayerSettings[_player.ID].Stance) {
      console.log("SEND UPDATES", {
        ..._lobby.PlayerSettings[_player.ID],
        Stance: selectedStance,
      });
      updatePlayerSettings(_lobby.ID, _player, {
        ..._lobby.PlayerSettings[_player.ID],
        Stance: selectedStance,
      });
    }
  }

  $: {
    if (!$player?.ID) goto("/");
    if (!$lobby?.ID) goto(`/`);
  }
  $: opponent = $lobby.Players.find((p) => p.ID != $player.ID);

  onMount(() => {
    // If there's no game create one
    // If there's a game go to game (Handle in reactive?)
    // Handle all game states here
  });

  function handleStartGame() {}
</script>

<div>
  <h1>Game Lobby</h1>
  <p>Code: {$lobby.ID}</p>
  <div class="inline-flex">
    <div class="nes-container player-container">
      <h3>{$player.Name}</h3>
      Select your stance:
      <div class="nes-container stance-container">
        <!--               checked={stance.id == $lobby.PlayerSettings[$player.ID].Stance}
 -->
        {#each STANCES as stance, i}
          <label>
            <input type="radio" class="nes-radio" name="player-stance" bind:group={selectedStance} value={stance.id} />
            <span>{stance.name}</span>
          </label>
        {/each}
      </div>
      <button class="nes-btn ready-button">Ready</button>
    </div>

    <div class="nes-container player-container opponent-container">
      {#if !opponent}
        <h3>Waiting for opponent...</h3>
      {:else}
        <h3>{opponent?.Name}</h3>
        <div class="nes-container stance-container">
          {#each STANCES as stance, i}
            <label>
              <input
                type="radio"
                class="nes-radio"
                name="opponent-stance"
                checked={stance.id == $lobby.PlayerSettings[opponent.ID].Stance}
              />
              <span>{stance.name}</span>
            </label>
          {/each}
        </div>
      {/if}
    </div>
  </div>
  {#if $player.ID == $lobby.Host}
    <button type="button" class={"nes-btn is-primary start-button"} on:click={handleStartGame}>Start Game</button>
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
  .stance-container {
    padding-left: 0;
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
  .opponent-container {
    pointer-events: none;
  }
</style>
