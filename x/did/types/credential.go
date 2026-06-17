package types

const maxCredentialDataLength = 64 * 1024

func NewCredential(did, sid, hash, uri string, data []byte) Credential {

	return Credential{Did: did, Sid: sid, Hash: hash, Uri: uri, Data: data}
}
