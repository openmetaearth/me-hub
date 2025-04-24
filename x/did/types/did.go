package types

const DidLength = 13

func NewDidInfo(did, address, pubkey string, status DidStatus) DidInfo {
	return DidInfo{
		Did:     did,
		Address: address,
		Pubkey:  pubkey,
		Status:  status,
	}
}
