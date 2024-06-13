# oauth2-sample

Simplified OAuth2 HTTP Server sample.

## Configuration
Please, see below configuration environment variables:

| **Environment variable name**   | **Description**                            | **Required** |     **Default**      |
|---------------------------------|--------------------------------------------|:------------:|:--------------------:|
| HTTP_SERVER_PORT                | Server will be listening on provided port  |     [ ]      |         8080         |
| JWT_ISSUER                      | Issuer used in `iss` JWT claim             |     [ ]      | oauth2-server-issuer |
| JWT_TTL                         | JWT time-to-live                           |     [ ]      |          5m          |
| JWT_PRIVATE_KEY                 | used to sign JWT                           |     [x]      |                      |
| OAUTH_CREDENTIALS_CLIENT_ID     | client id used in basic auth               |     [x]      |                      |
| OAUTH_CREDENTIALS_CLIENT_SECRET | client secret used in basic auth           |     [x]      |                      |

## Start via Makefile commands
In case [Make](https://www.gnu.org/software/make/manual/) utility is installed you can start server easily via single command.
### 1. Run locally
The simplest way to start server locally is to run the command below (replace `${OAUTH_CREDENTIALS_CLIENT_ID}` and `${OAUTH_CREDENTIALS_CLIENT_SECRET}`):
```console
make local OAUTH_CREDENTIALS_CLIENT_ID=${OAUTH_CREDENTIALS_CLIENT_ID} \
           OAUTH_CREDENTIALS_CLIENT_SECRET=${OAUTH_CREDENTIALS_CLIENT_SECRET}
```
Providing `OAUTH_CREDENTIALS_CLIENT_ID` and `OAUTH_CREDENTIALS_CLIENT_SECRET` mandatory options are enough, JWT private key will be generated automatically via `openssl genrsa 2048` command. To customize server config provide additional options:
```console
make local JWT_PRIVATE_KEY=${JWT_PRIVATE_KEY} \
	JWT_ISSUER=${JWT_ISSUER} \
	JWT_TTL=${JWT_TTL} \
	HTTP_SERVER_PORT=${HTTP_SERVER_PORT} \
	OAUTH_CREDENTIALS_CLIENT_ID=${OAUTH_CREDENTIALS_CLIENT_ID} \
	OAUTH_CREDENTIALS_CLIENT_SECRET=${OAUTH_CREDENTIALS_CLIENT_SECRET}
```
Make HTTP request to server supplying credentials set before:
```console
curl -X POST 'localhost:8080/token' -u "${OAUTH_CREDENTIALS_CLIENT_ID}:${OAUTH_CREDENTIALS_CLIENT_SECRET}"
```

### 2. Deploy to minikube
Make sure minikube is installed, up and running:
```console
minikube status
```
Server configuration is located in [k8s/config.yaml](./k8s/config.yaml) config map. `JWT_PRIVATE_KEY` is generated via `openssl genrsa 2048`:
```console
make minikube OAUTH_CREDENTIALS_CLIENT_ID=${OAUTH_CREDENTIALS_CLIENT_ID} \
              OAUTH_CREDENTIALS_CLIENT_SECRET=${OAUTH_CREDENTIALS_CLIENT_SECRET}
```
In case private key need to be provided pass corresponding option additionally:
```console
make minikube JWT_PRIVATE_KEY=${JWT_PRIVATE_KEY} \
              OAUTH_CREDENTIALS_CLIENT_ID=${OAUTH_CREDENTIALS_CLIENT_ID} \
              OAUTH_CREDENTIALS_CLIENT_SECRET=${OAUTH_CREDENTIALS_CLIENT_SECRET}
```
Next step is to receive ingress IP:
```console
foo@bar:~$ minikube kubectl -- get ingress -n oauth2-sample
NAME                    CLASS   HOSTS   ADDRESS        PORTS   AGE
oauth2-sample-ingress   nginx   *       192.168.49.2   80      13m
```
Take IP from `ADDRESS` column and execute HTTP request via curl (e.g. 192.168.49.2 in output above instead of `${IP}` in command below):
```console
curl -X POST 'http://${IP}/token' -u "${OAUTH_CREDENTIALS_CLIENT_ID}:${OAUTH_CREDENTIALS_CLIENT_SECRET}"
```
## Start manually
In case full control is required server can be started via running all commands manually.
### 1. Run locally
Run `main.go` passing corresponding environment variables (`JWT_PRIVATE_KEY` is generated via `openssl genrsa 2048`):
```console
OAUTH_CREDENTIALS_CLIENT_ID=${OAUTH_CREDENTIALS_CLIENT_ID} \
OAUTH_CREDENTIALS_CLIENT_SECRET=${OAUTH_CREDENTIALS_CLIENT_SECRET} go run main.go
```
With full list of options:
```console
JWT_PRIVATE_KEY=${JWT_PRIVATE_KEY} \
JWT_ISSUER=${JWT_ISSUER} \
JWT_TTL=${JWT_TTL} \
HTTP_SERVER_PORT=${HTTP_SERVER_PORT} \
OAUTH_CREDENTIALS_CLIENT_ID=${OAUTH_CREDENTIALS_CLIENT_ID} \
OAUTH_CREDENTIALS_CLIENT_SECRET=${OAUTH_CREDENTIALS_CLIENT_SECRET} go run main.go
```
### 1. Deploy to minikube
Run commands in a sequence:
```console
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
```
You can define private key from file when defining secret:
```console
minikube kubectl -- create secret generic -n oauth2-sample oauth2-sample-secrets \
            --from-literal=OAUTH_CREDENTIALS_CLIENT_ID=${OAUTH_CREDENTIALS_CLIENT_ID} \
            --from-literal=OAUTH_CREDENTIALS_CLIENT_SECRET=${OAUTH_CREDENTIALS_CLIENT_SECRET} \
            --from-file=JWT_PRIVATE_KEY=./private.pem"
```
Receive ingress IP:
```console
foo@bar:~$ minikube kubectl -- get ingress -n oauth2-sample
NAME                    CLASS   HOSTS   ADDRESS        PORTS   AGE
oauth2-sample-ingress   nginx   *       192.168.49.2   80      13m
```
Take IP from `ADDRESS` column and execute HTTP request via curl (e.g. 192.168.49.2 in output above instead of `${IP}` in command below):
```console
curl -X POST 'http://${IP}/token' -u "${OAUTH_CREDENTIALS_CLIENT_ID}:${OAUTH_CREDENTIALS_CLIENT_SECRET}"
```