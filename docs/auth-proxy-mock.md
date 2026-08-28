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
 | MAKE OPTION       | DEFAULT VALUE               | PASSED TO                           | DESCRIPTION                                          |
 | ---               | ---                         | ---                                 | ---                                                  |
 | `APM.LISTEN_ADDR` | `127.0.0.1:8888`            | `auth-proxy-mock --listen-on=...`   | Accept requests on this address                      |
 | `WEB.ADDRESS`     | `127.0.0.1:8080`            | `auth-proxy-mock --forward-to=...`  | Forward requests to this WiP service instance        |
 | `APM.REMOTE_USER` | `$USER`                     | `auth-proxy-mock --remote-user=...` | Inject this remote user in the HTTP request (header) |
 | `APM.SSO_SUB`     | `<md5sum($USER) | as-uuid>` | `auth-proxy-mock --sso-sub=...`     | Inject this sso sub in the HTTP request (header)     |
<!-- END make apm.start-server.usage.options -->

### `auth-proxy-mock` options

<!-- START bin/auth-proxy-mock --help -->
```bash
./bin/auth-proxy-mock --help
	# Usage of ./bin/auth-proxy-mock:
	#   -forward-to address
	#     	forward requests to this address (default "http://127.0.0.1:8080/")
	#   -listen-on address
	#     	listen on this address for requests to forward (default "127.0.0.1:8888")
	#   -remote-user username
	#     	set value for the REMOTE_USER http header to this username
	#   -sso-sub uuid
	#     	set value for the SSO_SUB http header to this uuid
```
<!-- END bin/auth-proxy-mock --help -->
