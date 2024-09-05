package types

func NewCredential(did, sid, hash, uri string) Credential {

	return Credential{Did: did, Sid: sid, Hash: hash, Uri: uri}
}
