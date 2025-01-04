<script lang="ts">
  import { guess } from "$lib/api";
  import { lobby, player } from "$lib/store";

  $: opponent = $lobby.Players.find((p) => p.ID != $player.ID) ?? { ID: "", Name: "" }; //Hacky

  let guessValue = "";
  let playerGuesses: Boolean[] = [];
  let opponentGuesses: Boolean[] = [];

  $: {
    const _lobby = $lobby;
    const _player = $player;
    const word = _lobby.Game?.Rally?.Word ?? "";
    const _playerGuessData = _lobby.Game?.Rally?.Guesses[_player.ID] ?? [];
    const _opponentGuessData = _lobby.Game?.Rally?.Guesses[opponent.ID] ?? [];
    if (word) {
      for (const [index, char] of [...word].entries()) {
        playerGuesses[index] = _playerGuessData[index] == char;
        opponentGuesses[index] = _opponentGuessData[index] == char;
      }
    }
  }

  async function handleGuess() {
    await guess($lobby.ID, $player, guessValue);
  }

  function endTurn() {}
</script>

<div>
  <div></div>
  <div class="grid-container">
    <div class="player-container">
      <div class="nes-container with-title player-board">
        <p class="title">{$player.Name}'s Board</p>
        <div class="guess-container">
          {#each playerGuesses as playerGuess}
            <div class={`guess-box ${playerGuess ?? "correct-guess"}`}></div>
          {/each}
        </div>
      </div>
      <div class="nes-container with-title player-action">
        <p class="title">{$player.Name}'s Actions</p>
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
          <button class={`nes-btn guess-button ${guessValue ? "" : "is-disabled"}`} on:click={handleGuess}>Guess</button
          >
        </div>
        <button class="nes-btn end-turn-button">End Turn</button>
      {:else}
        <h4>{opponent.Name}'s turn</h4>
      {/if}
    </div>

    <div class="player-container">
      <div class="nes-container with-title player-board">
        <p class="title">{opponent?.Name}'s Board</p>
        <div class="guess-container">
          {#each opponentGuesses as opponentGuess}
            <div class={`guess-box ${opponentGuess ?? "correct-guess"}`}></div>
          {/each}
        </div>
      </div>

      <div class="nes-container with-title player-action">
        <p class="title">{opponent?.Name}'s Actions</p>
      </div>
    </div>
  </div>
</div>

<style>
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
</style>
