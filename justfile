ent-install:
	go get entgo.io/ent/cmd/ent

ent-gen:
	go generate ./ent

ent-new ENTITY:
	go run entgo.io/ent/cmd/ent new {{ENTITY}}

run:
	clear
	go run cmd/server/main.go

tidy:
  go mod tidy