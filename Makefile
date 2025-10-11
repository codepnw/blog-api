include dev.env

FEATURE_DIRS = domains repositories usecases handlers
MIGRATE_DIR = internal/database/migrations
MIGRATE_DB = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# ----------- DOCKER -----------
docker-up:
	@docker compose --env-file=dev.env up

# ----------- MIGRATIONS -----------
migrate-create:
	@migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name)

migrate-up:
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_DIR) up

migrate-down:	
	@migrate -database $(MIGRATE_DB) -path $(MIGRATE_DIR) down 1

# ----------- FEATURE MANAGE -----------
ft-new:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Please specify feature name, e.g. 'make ft-new name=post'"; \
		exit 1; \
	fi; \
	for dir in $(FEATURE_DIRS); do \
	 	mkdir -p internal/$$dir/$(name) && \
		case $$dir in \
			domains) pkg="$(name)domain";; \
			repositories) pkg="$(name)repo";; \
			usecases) pkg="$(name)usecase";; \
			handlers) pkg="$(name)handler";; \
		esac; \
		FILE=internal/$$dir/$(name)/$(name).go; \
		if [ ! -f $$FILE ]; then \
			echo "package $$pkg" > $$FILE; \
			echo "Created: $$FILE"; \
		else \
			echo "Skipped (already exists): $$FILE"; \
		fi; \
	done

ft-del:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Please specify feature name, e.g. 'make ft-del name=post'"; \
		exit 1; \
	fi; \
	for dir in $(FEATURE_DIRS); do \
		if [ -d internal/$$dir/$(name) ]; then \
			rm -rf internal/$$dir/$(name); \
			echo "🗑️  Removed: internal/$$dir/$(name)"; \
		else \
			echo "⚠️  Not found: internal/$$dir/$(name)"; \
		fi; \
	done