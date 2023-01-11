package main

type Logs struct {
	LogName      string
	PassCount    string
	CookiesCount string
	UserId       string
	Tag          string
	DateTime     string
	Wallets      string
	Telegram     string
	Steam        string
}

type Dashboard struct {
	USA          int
	GER          int
	AUS          int
	UK           int
	RO           int
	BR           int
	CountryCount int
	Subtime      string
	WalletCount  int
	CookieCount  int
	MinMoth1     string
	MinMoth2     string
	MinMoth3     string
	MinMoth4     string
	MinMoth5     string
	WalletMonth1 int
	WalletMonth2 int
	WalletMonth3 int
	WalletMonth4 int
	WalletMonth5 int
	CookieMonth1 int
	CookieMonth2 int
	CookieMonth3 int
	CookieMonth4 int
	CookieMonth5 int
	DayName1     string
	DayName2     string
	DayName3     string
	DayName4     string
	DayName5     string
	DayName6     string
	DayName7     string
	DayName8     string
	DayName9     string
	DayName10    string
	DayName11    string
	DayName12    string
	DayLog1      int
	DayLog2      int
	DayLog3      int
	DayLog4      int
	DayLog5      int
	DayLog6      int
	DayLog7      int
	DayLog8      int
	DayLog9      int
	DayLog10     int
	DayLog11     int
	DayLog12     int
	MonthName1   string
	MonthName2   string
	MonthName3   string
	MonthName4   string
	MonthValue1  int
	MonthValue2  int
	MonthValue3  int
	MonthValue4  int
}

type LogsForPage struct {
	Subtime     string
	Alllogs     []ELogs
	LogsCount   int
	PassCount   int
	SteamCount  int
	WalletCount int
}

type ELogs struct {
	LogName      string
	PassCount    int
	CookiesCount string
	UserId       string
	Tag          string
	DateTime     string
	Wallets      string
	Telegram     string
	Steam        string
	LogSize      float64
}
