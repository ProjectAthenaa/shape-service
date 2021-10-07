compileProto:
	protoc --go_out=./protos --go_opt=paths=source_relative --go-grpc_out=./protos --go-grpc_opt=paths=source_relative ./Shape.proto

rollout:
	doctl kubernetes cluster kubeconfig save --expiry-seconds 600 athena
	kubectl rollout restart deployments shape -n antibots
	kubectl rollout status deployments shape -n antibots

build:
	docker build --build-arg GH_TOKEN=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  -t athena/shape_local:1.0 .

run:
	export DEBUG=1
	export PATH=$PATH:/usr/local/go/bin
	export REDIS_URL=rediss://default:ncwzkvkuy09khcry@shape-do-user-9104051-0-8888.b.db.ondigitalocean.com:25061
	go
	go run .

runDocker:
	docker build --build-arg GH_TOKEN=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  -t athena/shape_local:1.0 .
	docker run -e DEBUG=1 -p 3000:3000 --name shape_local athena/shape_local:1.0

tidy:
	go mod tidy -compat=1.17