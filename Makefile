.PHONY: dev
dev: go-dev svelte-dev

.PHONY: go-dev
go-dev:
	@bash -c 'cd server && go run main.go & pid1=$!; \
		trap "kill $$pid1" SIGINT; \
		wait $$pid1'

.PHONY: svelte-dev
svelte-dev:
	@bash -c 'cd client && pnpm dev & pid2=$!; \
		trap "kill $$pid2" SIGINT; \
		wait $$pid2'

# .PHONY: wait
# wait:
# 	wait