# (DEV) auth-proxy-mock

`auth-proxy-mock` is a development helper server which takes in a HTTP
requests, adds some additional HTTP headers and forwards the request to
the main WIP service.

This is mainly intended to mock OpenID connect OAuth provided by another
HTTP proxy (e.g. Apache).

## Build

makefile: `make.d/auth-proxy-mock.mk`

## USAGE

### Via `make` wrapper

`make apm.start-server`: start a proxy server and forward requests

<!-- START make apm.start-server.usage.options -->
 | MAKE OPTION       | DEFAULT VALUE               | PASSED TO                                                      | DESCRIPTION                                                         |
 | ---               | ---                         | ---                                                            | ---                                                                 |
 | `APM.LISTEN_ADDR` | `127.0.0.1:8888`            | `auth-proxy-mock --listen-on=...`                              | Accept requests on this address                                     |
 | `WIP.ADDRESS`     | `127.0.0.1:8080`            | `auth-proxy-mock --forward-to=...`                             | Forward requests to this WiP service instance                       |
 | `APM.REMOTE_USER` | `<md5sum($USER) | as-uuid>` | `auth-proxy-mock --header "REMOTE_USER:..."`                   | Inject this sso sub in the HTTP request (header)                    |
 | `APM.USERNAME`    | `$USER`                     | `auth-proxy-mock --header="OIDC_CLAIM_preferred_username:..."` | Inject this remote user in the HTTP request (header)                |
 | `APM.GROUPS`      | `group1,group2,...`         | `auth-proxy-mock --header="OIDC_CLAIM_memberOf:..."`           | Inject this comma seperated group list in the HTTP request (header) |
<!-- END make apm.start-server.usage.options -->

### `auth-proxy-mock` options

<!-- START bin/auth-proxy-mock --help -->
```bash
./bin/auth-proxy-mock --help
	# Usage of ./bin/auth-proxy-mock:
	#   -forward-to string
	#     	forward requests to this address (default "http://127.0.0.1:8080/")
	#   -header name:value
	#     	specify additional HTTP headers (name:value)
	#   -listen-on string
	#     	listen on this address for requests to forward (default "127.0.0.1:8888")
```
<!-- END bin/auth-proxy-mock --help -->
