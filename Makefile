MIGRATE_CMD = go run cmd/db.go
RUN_CMD = go run main.go

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

