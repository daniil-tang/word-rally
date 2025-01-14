<script lang="ts">
  import { goto } from "$app/navigation";
  import { createGame, leaveLobby, startGame, updatePlayerSettings } from "$lib/api";
  import { GAME_STATE, STANCES } from "$lib/constants";
  import { player, lobby } from "$lib/store";
  import { onMount } from "svelte";

  $: {
    if (!$player?.ID) goto("/");
    if (!$lobby) goto(`/lobby`);
    // Goto game screen
    if ($lobby?.Game?.State == GAME_STATE.IN_PROGRESS) goto(`/lobby/${$lobby.ID}/game/${$lobby.Game.ID}`);
    // if ($lobby?.Game?.State == GAME_STATE.FINISHED) goto(`/`); // Prolly don't need this..?
  }

  $: selectedStance = $lobby?.PlayerSettings?.[$player.ID]?.Stance;
  // $: playerReady = $lobby?.PlayerSettings[$player.ID].Ready;

  $: {
    if (selectedStance) {
      const _lobby = $lobby;
      const _player = $player;

      const currentStance = _lobby.PlayerSettings[_player.ID].Stance;
      if (selectedStance !== currentStance) {
        updatePlayerSettings(_lobby.ID, _player, {
          ..._lobby.PlayerSettings[_player.ID],
          Stance: selectedStance,
        });
      }
    }
  }

  $: opponent = $lobby.Players.find((p) => p.ID != $player.ID);

  function handlePlayerReady() {
    updatePlayerSettings($lobby.ID, $player, {
      ...$lobby.PlayerSettings[$player.ID],
      Ready: !$lobby.PlayerSettings[$player.ID].Ready,
    });
  }

  async function handleStartGame() {
    await createGame($lobby.ID, $player);
    await startGame($lobby.ID, $player);
  }

  async function handleLeaveLobby() {
    await leaveLobby($lobby.ID, $player);
  }
</script>

<div>
  <div class="topbar">
    <h1>Game Lobby</h1>
    <button type="button" class={`nes-btn`} on:click={handleLeaveLobby}>Leave</button>
  </div>
  <p>Host: {$lobby.Players.find((p) => p.ID == $lobby.Host)?.Name}</p>
  <p>Code: {$lobby.ID}</p>
  <div class="inline-flex">
    <div class="nes-container player-container">
      <h3>{$player.Name}</h3>
      Select your stance:
      <div class="nes-container stance-container">
        {#each STANCES as stance, i}
          <label>
            <input type="radio" class="nes-radio" name="player-stance" bind:group={selectedStance} value={stance.id} />
            <span>{stance.name}</span>
          </label>
        {/each}
      </div>
      <button
        class={`nes-btn ready-button ${$lobby.PlayerSettings[$player.ID].Ready ? "is-success" : ""}`}
        on:click={handlePlayerReady}>Ready</button
      >
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
        {#if !opponent || !$lobby.PlayerSettings[opponent.ID].Ready}
          <div class={`opponent-ready-text nes-text is-error`}>Not Ready</div>
        {:else}
          <div class={`opponent-ready-text  nes-text is-success`}>Ready</div>
        {/if}
      {/if}
    </div>
  </div>
  {#if $player.ID == $lobby.Host}
    <button
      type="button"
      class={`nes-btn is-primary start-button ${!opponent || !($lobby.PlayerSettings[$player.ID].Ready && $lobby.PlayerSettings[opponent.ID].Ready) ? "is-disabled" : ""}`}
      on:click={handleStartGame}
      disabled={!opponent || !($lobby.PlayerSettings[$player.ID].Ready && $lobby.PlayerSettings[opponent.ID].Ready)}
      >Start Game</button
    >
  {/if}
</div>

<style>
  .topbar {
    display: flex;
    justify-content: space-between;
  }
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
  .opponent-ready-text {
    display: block;
    margin: auto;
    margin-top: 10px;
    width: max-content;
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
