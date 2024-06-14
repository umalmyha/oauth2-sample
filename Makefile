JWT_PRIVATE_KEY ?= $$(openssl genrsa 2048)

.PHONY: local
local:
	@if [ -z "$(OAUTH_CREDENTIALS_CLIENT_ID)" ] || [ -z "$(OAUTH_CREDENTIALS_CLIENT_SECRET)" ]; then \
      		echo "both OAUTH_CREDENTIALS_CLIENT_ID and OAUTH_CREDENTIALS_CLIENT_SECRET arguments must be provided"; \
      		exit 1; \
	fi
	@JWT_PRIVATE_KEY=$(JWT_PRIVATE_KEY) \
	JWT_ISSUER=${JWT_ISSUER} \
	JWT_TTL=${JWT_TTL} \
	HTTP_SERVER_PORT=${HTTP_SERVER_PORT} \
	OAUTH_CREDENTIALS_CLIENT_ID=${OAUTH_CREDENTIALS_CLIENT_ID} \
	OAUTH_CREDENTIALS_CLIENT_SECRET=${OAUTH_CREDENTIALS_CLIENT_SECRET} go run main.go

.PHONY: minikube
minikube:
	@if [ -z "$(OAUTH_CREDENTIALS_CLIENT_ID)" ] || [ -z "$(OAUTH_CREDENTIALS_CLIENT_SECRET)" ]; then \
  		echo "both OAUTH_CREDENTIALS_CLIENT_ID and OAUTH_CREDENTIALS_CLIENT_SECRET arguments must be provided"; \
  		exit 1; \
  	 fi
	minikube addons enable ingress
	minikube kubectl -- apply -f k8s/namespace.yaml
	minikube kubectl -- apply -f k8s/config.yaml
	minikube kubectl -- create secret generic -n oauth2-sample oauth2-sample-secrets \
	            --from-literal=OAUTH_CREDENTIALS_CLIENT_ID=${OAUTH_CREDENTIALS_CLIENT_ID} \
	            --from-literal=OAUTH_CREDENTIALS_CLIENT_SECRET=${OAUTH_CREDENTIALS_CLIENT_SECRET} \
	            --from-literal=JWT_PRIVATE_KEY="${JWT_PRIVATE_KEY}"
	minikube kubectl -- apply -f k8s/deployment.yaml
	minikube kubectl -- apply -f k8s/service.yaml
	minikube kubectl -- apply -f k8s/ingress.yaml