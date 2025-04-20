MIGRATE_CMD = go run cmd/migrate/runner.go
RUN_CMD = go run cmd/oma/main.go

migrate:
	${MIGRATE_CMD} migrate

reset:
	${MIGRATE_CMD} reset

init:
	${RUN_CMD} init

diff:
	${RUN_CMD} diff

commit:
	${RUN_CMD} commit

plain:
	${RUN_CMD}

