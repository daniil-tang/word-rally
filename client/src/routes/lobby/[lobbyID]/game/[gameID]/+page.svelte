<script lang="ts">
  import { goto } from "$app/navigation";
  import { createGame, endTurn, guess, useSkill } from "$lib/api";
  import { GAME_STATE, STANCE_DATA } from "$lib/constants";
  import { actionLog, lobby, player } from "$lib/store";
  import { onMount } from "svelte";

  $: {
    if (!$player?.ID) goto("/");
    if (!$lobby) goto(`/lobby`);
    // Goto game screen
    if (!$lobby?.Game || $lobby?.Game?.State == GAME_STATE.WAITING) goto(`/lobby/${$lobby.ID}`);
  }

  $: opponent = $lobby.Players.find((p) => p.ID != $player.ID) ?? { ID: "", Name: "" }; //Hacky

  let guessValue = "";
  let playerGuesses: Boolean[] = [];
  let opponentGuesses: Boolean[] = [];
  let endGameDialog: HTMLDialogElement | null = null;
  $: {
    const _lobby = $lobby;
    const _player = $player;
    const word = _lobby.Game?.Rally?.Word ?? "";
    const _playerGuessData = _lobby.Game?.Rally?.Guesses[_player.ID] ?? [];
    const _opponentGuessData = _lobby.Game?.Rally?.Guesses[opponent.ID] ?? [];
    if (word) {
      // Create new arrays to trigger reactivity
      playerGuesses = [...word].map((char, index) => _playerGuessData[index] === char.charCodeAt(0));
      opponentGuesses = [...word].map((char, index) => _opponentGuessData[index] === char.charCodeAt(0));
    }
    if (_lobby.Game?.State == GAME_STATE.FINISHED) {
      if (endGameDialog) {
        endGameDialog.showModal();
      }
    }
  }

  onMount(() => {
    actionLog.set([]);
  });

  async function handleGuess() {
    await guess($lobby.ID, $player, guessValue);
    guessValue = "";
  }

  function handleEndTurn() {
    endTurn($lobby.ID, $player);
  }

  async function handleReturnToLobby() {
    // await createGame($lobby.ID, $player);
    goto(`/lobby/${$lobby.ID}`);
  }

  async function handleSkillClick(skillId: string) {
    await useSkill($lobby.ID, $player, skillId);
  }
</script>

<div>
  <dialog class="nes-dialog end-game-dialog" id="end-game-dialog" bind:this={endGameDialog}>
    <form method="dialog">
      <menu class="dialog-menu end-game-dialog-menu">
        <h1>
          {#if $lobby?.Game?.Score[$player.ID] == 3}
            You Win! :D
          {:else}
            {"You Lose :("}
          {/if}
        </h1>
        <!-- {#if $player.ID == $lobby?.Host} -->
        <button class="nes-btn is-primary back-to-lobby-button" on:click={handleReturnToLobby}>Return to Lobby</button>
        <!-- {:else}
          <p>Waiting for host to return you to lobby</p>
        {/if} -->
      </menu>
    </form>
  </dialog>
  <div class="score-container">
    <div>{$player.Name}: {$lobby.Game?.Score[$player.ID]}</div>
    <div>First to 3</div>
    <div>{opponent.Name}: {$lobby.Game?.Score[opponent.ID]}</div>
  </div>
  <div class="grid-container">
    <div class="player-container">
      <div class="nes-container with-title player-board">
        <p class="title">{$player.Name}'s Board</p>
        <div class="guess-container">
          {#each playerGuesses as playerGuess, i}
            <div class={`guess-box ${playerGuess ? "correct-guess" : ""}`}>
              {playerGuess ? `${$lobby.Game?.Rally?.Word[i]}` : ""}
            </div>
          {/each}
        </div>
      </div>
      <div class="nes-container with-title player-action">
        <p class="title">{$player.Name}'s Actions</p>
        <div class="action-points-container">
          <span>Guess Points: {$lobby.Game?.Rally?.TurnActionPoints[$player?.ID]?.Guess}</span>
          <span>Skill Points: {$lobby.Game?.Rally?.TurnActionPoints[$player?.ID]?.Skill}</span>
        </div>
        <div class="skills-container">
          {#each STANCE_DATA[$lobby.PlayerSettings[$player.ID].Stance].skills as skillData}
            <!--Disable button if on cooldown  -->
            {#if ($lobby.Game?.PlayerCooldowns?.[$player.ID]?.[skillData.id] ?? 0) > 0}
              <button disabled class={`nes-btn is-disabled`}
                >Cooldown: {$lobby.Game?.PlayerCooldowns?.[$player.ID]?.[skillData.id]}</button
              >
            {:else}
              <div class="tooltip-container">
                <button
                  on:click={() => handleSkillClick(skillData.id)}
                  class={`nes-btn ${($lobby.Game?.PlayerCooldowns?.[$player.ID]?.[skillData.id] ?? 0) > 0 ? "is-disabled" : ""}`}
                  >{skillData.name}</button
                >
                <div class="tooltip-text">{skillData.description}</div>
              </div>
            {/if}
          {/each}
        </div>
      </div>
    </div>

    <div class="nes-container turn-container">
      {#if $lobby.Players[$lobby.Game?.Rally?.Turn].ID == $player.ID}
        <h4>Your turn</h4>
        <div class="guess-form">
          <p>Guess your next letter:</p>
          <input
            type="text"
            id="player_guess"
            bind:value={guessValue}
            maxlength="1"
            class="nes-input player-guess-input"
          />
          <button
            class={`nes-btn guess-button ${guessValue ? "" : "is-disabled"}`}
            on:click={handleGuess}
            disabled={!guessValue}>Guess</button
          >
        </div>
        <button class="nes-btn end-turn-button" on:click={handleEndTurn}>End Turn</button>
      {:else}
        <h4>{opponent.Name}'s turn</h4>
      {/if}
    </div>

    <div class="player-container">
      <div class="nes-container with-title player-board">
        <p class="title">{opponent?.Name}'s Board</p>
        <div class="guess-container">
          {#each opponentGuesses as opponentGuess}
            <div class={`guess-box ${opponentGuess ? "correct-guess" : ""}`}></div>
          {/each}
        </div>
      </div>

      <div class="nes-container with-title player-action">
        <p class="title">{opponent?.Name}'s Actions</p>
        <div class="action-points-container">
          <span>Guess Points: {$lobby.Game?.Rally?.TurnActionPoints[opponent?.ID]?.Guess}</span>
          <span>Skill Points: {$lobby.Game?.Rally?.TurnActionPoints[opponent?.ID]?.Skill}</span>
        </div>
      </div>
    </div>
  </div>
  <br />
  <br />
  <div class="nes-container with-title">
    <h3 class="title">Action Logs</h3>

    <div class="log-container">
      {#each $actionLog as log, index}
        <div class="log-item">
          {log}
        </div>
      {/each}
    </div>
  </div>
  <br />
  <br />
  <div class="nes-text instructions">
    <h2>How to play</h2>
    <p>Players take turns guessing the word.</p>
    <p>First to guess the word wins the rally.</p>
    <p>First to 3 wins the game!</p>
    <br />
    <h2>Stances and Skills</h2>

    <div>
      <h3>Tennis</h3>
      <ul>
        <li>
          <h4>Ace</h4>
          <p><strong>Cooldown:</strong> 4 turns</p>
          <p>Instantly grants a correct guess. (Requires 1 guess point).</p>
        </li>
        <li>
          <h4>Fault</h4>
          <p><strong>Cooldown:</strong> 5 turns</p>
          <p>Makes opponent miss next turn.</p>
        </li>
      </ul>
    </div>

    <div>
      <h3>Football</h3>
      <ul>
        <li>
          <h4>Tackle</h4>
          <p><strong>Cooldown:</strong> 3 turns</p>
          <p>Steals a correct guess from the opponent.</p>
        </li>
        <li>
          <h4>Goalkeeper</h4>
          <p><strong>Cooldown:</strong> 4 turns</p>
          <p>Block the opponent's next correct guess.</p>
        </li>
      </ul>
    </div>
  </div>
</div>

<style>
  .score-container {
    display: flex;
    justify-content: space-between;
  }

  .grid-container {
    display: inline-grid;
    grid-template-columns: 1fr max-content 1fr;
    gap: 10px;
  }

  .player-container {
    min-width: 400px;
  }

  .player-board {
  }

  .player-action {
    margin-top: 10px;
  }

  .turn-container {
    min-width: 400px;
    /* display: flex;
    flex-direction: column;
    justify-content: flex-end; */
  }

  .guess-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, 50px);
    gap: 1px;
    width: 100%;
    text-align: center;
  }

  .guess-box {
    width: 40px;
    height: 40px;
    text-align: center;
    border: 1px solid black;
    text-transform: uppercase;
    display: flex;
    justify-content: center;
    align-items: center;
    color: white;
  }

  .correct-guess {
    background-color: green;
  }

  .guess-form {
    display: inline;
  }
  .player-guess-input {
    width: auto;
    text-transform: uppercase;
    /* display: block; */
    /* margin: auto; */
    max-width: 100px;
  }

  .guess-button {
    /* margin-top: 20px;
    display: block; */
  }

  .end-turn-button {
    margin-top: 100px;
    display: block;
    width: 100%;
  }

  .action-points-container {
    display: flex;
    width: 100%;
    justify-content: space-between;
    font-size: 0.6rem;
  }

  .end-game-dialog {
    width: 600px;
    height: 400px;
    text-align: center;
  }

  .end-game-dialog-menu {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
  }

  /* .back-to-lobby-button {
    align-self: flex-end;
  } */
  .log-container {
    width: 100%;
    height: 300px;
    border: 2px solid #2d2f31;
    padding: 10px;
    /* background-color: #2d2f31; */
    border-radius: 8px;
    overflow-y: auto;
  }

  .log-item {
    margin-bottom: 10px;
    /* color: #f8f8f2; */
    font-size: 14px;
    /* background-color: #444; */
    padding: 5px;
    border-radius: 4px;
  }

  .tooltip-container {
    position: relative;
    display: inline-block;
  }

  /* Tooltip text style */
  .tooltip-text {
    visibility: hidden;
    background-color: #000;
    color: #fff;
    text-align: center;
    border-radius: 5px;
    padding: 5px;
    position: absolute;
    z-index: 1;
    bottom: 125%; /* Position the tooltip above the button */
    left: 50%;
    margin-left: -60px;
    font-family: "Press Start 2P", cursive; /* NES-style pixel font */
    font-size: 12px;
    width: 120px;
    text-transform: uppercase;
    letter-spacing: 1px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  }

  /* Show the tooltip on hover */
  .tooltip-container:hover .tooltip-text {
    visibility: visible;
  }
</style>
