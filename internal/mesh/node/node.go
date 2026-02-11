package node

const nodeEndpointPort = 51830 // default '51820'
const nodeInterface = "wgi"
const nodePrivKeyFname = "/etc/i12e/wg-priv-key"

func GetNodeInterface() string {
	return nodeInterface
}
