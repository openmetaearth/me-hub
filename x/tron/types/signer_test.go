package types

import (
	"crypto/sha256"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
)

func TestTronAddressFromSignature(t *testing.T) {
	dataHash := sha256.Sum256([]byte("tron - test -data"))
	key, _ := crypto.GenerateKey()
	signature, _ := NewTronSignature(dataHash[:], key)
	type args struct {
		hash      []byte
		signature []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test normal",
			args: args{
				hash:      dataHash[:],
				signature: signature,
			},
			want:    address.PubkeyToAddress(key.PublicKey).String(),
			wantErr: false,
		},
		{
			name: "test signature is error",
			args: args{
				hash:      dataHash[:],
				signature: []byte("error signature"),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test signature is nil",
			args: args{
				hash:      dataHash[:],
				signature: nil,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TronAddressFromSignature(tt.args.hash, tt.args.signature)
			if (err != nil) != tt.wantErr {
				t.Errorf("TronAddressFromSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TronAddressFromSignature() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateTronSignature(t *testing.T) {
	dataHash := sha256.Sum256([]byte("tron - test -data"))
	key, _ := crypto.GenerateKey()
	signature, _ := NewTronSignature(dataHash[:], key)
	type args struct {
		hash       []byte
		signature  []byte
		ethAddress string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test normal",
			args: args{
				hash:       dataHash[:],
				signature:  signature,
				ethAddress: address.PubkeyToAddress(key.PublicKey).String(),
			},
			wantErr: false,
		},
		{
			name: "test hash error",
			args: args{
				hash:       []byte("test hash error"),
				signature:  signature,
				ethAddress: address.PubkeyToAddress(key.PublicKey).String(),
			},
			wantErr: true,
		},
		{
			name: "test hash is nil",
			args: args{
				hash:       nil,
				signature:  signature,
				ethAddress: address.PubkeyToAddress(key.PublicKey).String(),
			},
			wantErr: true,
		},
		{
			name: "test signature error",
			args: args{
				hash:       dataHash[:],
				signature:  []byte("test hash signature"),
				ethAddress: address.PubkeyToAddress(key.PublicKey).String(),
			},
			wantErr: true,
		},
		{
			name: "test signature is nil",
			args: args{
				hash:       dataHash[:],
				signature:  nil,
				ethAddress: address.PubkeyToAddress(key.PublicKey).String(),
			},
			wantErr: true,
		},
		{
			name: "test address error",
			args: args{
				hash:       dataHash[:],
				signature:  signature,
				ethAddress: "error address",
			},
			wantErr: true,
		},
		{
			name: "test address is nil",
			args: args{
				hash:       dataHash[:],
				signature:  signature,
				ethAddress: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateTronSignature(tt.args.hash, tt.args.signature, tt.args.ethAddress); (err != nil) != tt.wantErr {
				t.Errorf("ValidateTronSignature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
