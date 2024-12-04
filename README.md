# web3-auth-svc

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Description

RariMe Auth service designed to authorize users with Iden3 AuthV2 ZK-proofs and issue JWT tokens based on it.
This JWT can be used on other internal or external service to authenticate user for executing endpoints.

Frontend firstly should request Base64-encoded challenge using `v1/authorize/{address}/challenge` request.
Then generate AuthV2 ZK proof with received challenge as decoded big-endian value. Using this proof
execute `v1/authorize` request and receive JWT (refresh and access) tokens in response and also in cookies.

## Usage

To integrate on other service use the [pkg/auth](./pkg/auth) package.
It contains client and client config to execute `v1/validate` requests and example of `grants` that can be used
with `Authenticates` method to check uses access.

Example:

Add middleware to endpoints that require auth:

```go
package middleware

import (
	"net/http"

	"github.com/rarimo/web3-auth-svc/pkg/auth"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func AuthMiddleware(client *auth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claim, err := client.ValidateJWT(r)
			if err != nil {
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			// Save claims somewhere (probably in request context)
			ctx := handlers.CtxClaim(claim)(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
```

How protected endpoints definition looks like:

```go
r.Route("/integrations/your-service", func (r chi.Router) {
    r.Route("/v1", func (r chi.Router) {
        r.Post("/unprotected", handlers.Unprotected)
        r.With(middleware.AuthMiddleware(s.client)).Get("/protected", handlers.Protected)
    })
})
```

Then, use parsed claims in handler to allow users execute business logic:

```go
if !auth.Authenticates([]resources.Claim{claim}, auth.UserGrant("address")) {
    ape.RenderErr(w, problems.Unauthorized())
    return
}
```

## Install

  ```
  git clone github.com/rarimo/web3-auth-svc
  cd web3-auth-svc
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  ./main run service
  ```

## Documentation

We do use openapi:json standard for API. We use swagger for documenting our API.

To open online documentation, go to [swagger editor](http://localhost:8080/swagger-editor/) here is how you can start it

```
  cd docs
  npm install
  npm start
```

To build documentation use `npm run build` command,
that will create open-api documentation in `web_deploy` folder.

To generate resources for Go models run `./generate.sh` script in root folder.
use `./generate.sh --help` to see all available options.

Note: if you are using Gitlab for building project `docs/spec/paths` folder must not be
empty, otherwise only `Build and Publish` job will be passed.

## Running from docker

Make sure that docker installed.

use `docker run ` with `-p 8080:80` to expose port 80 to 8080

  ```
  docker build -t github.com/rarimo/web3-auth-svc .
  docker run -e KV_VIPER_FILE=/config.yaml github.com/rarimo/web3-auth-svc
  ```

## Running from Source

* Set up environment value with config file path `KV_VIPER_FILE=./config.yaml`
* Provide valid config file
* Launch the service with `run service` command

## Contact

Responsible Oleg Fomenko
The primary contact for this project is t.me/of_dl
