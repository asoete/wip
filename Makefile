# -- General
.DEFAULT_GOAL := compile

# -- MAIN Config

WEB.ADDRESS := 127.0.0.1:8080

# -- CONTROL Config

CTL.ADDRESS := 127.0.0.1:8100
CTL.TOKEN := let-the-dev-times-roll
CTL.TOKEN_HEADER := X-WiP-Ctl-Token: $(CTL.TOKEN)


include make.d/*.mk
