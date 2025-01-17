# Setup
- Have the following installed
  - NodeJS
  - pnpm
  - Go
- Run `make dev -j`

# Gameplay
- Players take turns guessing a letter each turn until one player completes the word, thereby winning the rally.
- Each turn, players have 1 guess point and 1 skill point by default
  - Guessing a letter consumes 1 guess point
  - Activating a skill consumes 1 skill point
- First player to win 3 rallies wins the game.