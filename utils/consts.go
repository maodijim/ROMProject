package utils

const (
	zoneId      = "ZoneId"
	version     = "Version"
	accId       = "accId"
	serverId    = "serverId"
	ip          = "ip"
	domain      = "domain"
	lineGrp     = "LineGrp"
	device      = "Device"
	deviceId    = "deviceId"
	clientVer   = "ClientVer"
	langZone    = "LangZone"
	ipAndPort   = "IpPort"
	phone       = "phone"
	safeDevice  = "SafeDevice"
	sha1Str     = "sha1Str"
	accessToken = "AccessToken"
	resourceVer = "ResVer"
	platV       = "PlatformVer"
	appPreVer   = "AppPreVer"
	lang        = "Lang"
	model       = "Model"
	phVer       = "PhoneVer"
	authoriz    = "Authoriz"
	authParams  = "AuthParams"
	region      = "region"
	char        = "char"
	gameServer  = "gameServer"
	authServer  = "authServer"

	zone      = "40"
	aapPreVer = "54"
	did       = "77241cb1d3dd91876c3e2e2d8504dca8"
	Received  = "received"
	Sent      = "sent"
)

var (
	CipherKey = []byte{95, 27, 5, 20, 131, 4, 8, 88}
	TcpFlag   = map[uint][]byte{
		0: []byte{0},
		1: []byte{1},
		2: []byte{2},
		3: []byte{3},
	}
	FlagByName = map[string][]byte{
		"noAction": []byte{0},
		"compress": []byte{1},
		"encrypt":  []byte{2},
		"both":     []byte{3},
	}
)
