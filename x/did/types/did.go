package types

func NewDidDocument(did, pubkey string, status DidStatus) DidDocument {
	info := DidInfo{
		Did:    did,
		Pubkey: pubkey,
		Status: status,
	}

	return DidDocument{Info: info, Vcs: []Credential{}}
}
