package constant

type CommandType string

var (
	Alive     CommandType = "alive"
	Help      CommandType = "help"
	Depth     CommandType = "depth"
	Arbitrage CommandType = "arbitrage"
)

type Exchange string

var (
	MAX   Exchange = "MAX"
	Rybit Exchange = "Rybit"
)
