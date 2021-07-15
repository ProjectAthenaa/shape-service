compileProto:
	protoc --go_out=./protos --go_opt=paths=source_relative --go-grpc_out=./protos --go-grpc_opt=paths=source_relative ./Shape.proto

rollout:
	doctl kubernetes cluster kubeconfig save --expiry-seconds 600 athena
	kubectl rollout restart deployments shape -n antibots
	kubectl rollout status deployments shape -n antibots