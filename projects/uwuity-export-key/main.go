package main

import (
	"bytes"
	"crypto/ed25519"
	"flag"
	"fmt"

	"github.com/linxGnu/grocksdb"
	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-media-repo/cmd/utilities/_common"
	"github.com/t2bot/matrix-media-repo/homeserver_interop"
	"github.com/t2bot/matrix-media-repo/util"
	"github.com/youmark/pkcs8"
)

func main() {
	outFormat := flag.String("format", "mmr", "Output format: 'mmr', 'synapse', or 'dendrite'")
	outFile := flag.String("output", "./signing.key", "Output file path")
	dbPath := flag.String("input", "./data.db", "Path to DB directory")
	flag.Parse()

	key, err := extractKey(*dbPath)
	if err != nil {
		logrus.Fatalf("Error: %v", err)
	}

	pubKey := key.PrivateKey.Public().(ed25519.PublicKey)
	pubKeyB64 := util.EncodeUnpaddedBase64ToString(pubKey)

	logrus.Infof("Extracted Key-ID: ed25519:%s", key.KeyVersion)
	logrus.Infof("Derived Public Key: %s", pubKeyB64)

	_common.EncodeSigningKeys([]*homeserver_interop.SigningKey{key}, *outFormat, *outFile)
}

func extractKey(dbPath string) (*homeserver_interop.SigningKey, error) {
	opts := grocksdb.NewDefaultOptions()
	defer opts.Destroy()

	cfNames, err := grocksdb.ListColumnFamilies(opts, dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list column families: %w", err)
	}

	cfOpts := make([]*grocksdb.Options, len(cfNames))
	for i := range cfOpts {
		cfOpts[i] = opts
	}

	db, handles, err := grocksdb.OpenDbColumnFamilies(opts, dbPath, cfNames, cfOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	defer db.Close()

	var globalCF *grocksdb.ColumnFamilyHandle
	for i, name := range cfNames {
		if name == "global" {
			globalCF = handles[i]
			defer globalCF.Destroy()
		} else {
			handles[i].Destroy()
		}
	}

	if globalCF == nil {
		return nil, fmt.Errorf("column family 'global' not found")
	}

	readOpts := grocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()

	val, err := db.GetCF(readOpts, globalCF, []byte("keypair"))
	if err != nil || val.Size() == 0 {
		return nil, fmt.Errorf("key 'keypair' not found")
	}
	defer val.Free()

	keyVersionBytes, derBytes, found := bytes.Cut(val.Data(), []byte{0xFF})
	if !found {
		return nil, fmt.Errorf("invalid keypair format")
	}

	parsedKey, err := pkcs8.ParsePKCS8PrivateKey(derBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER: %w", err)
	}

	privKey, ok := parsedKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parsed key is not an ed25519.PrivateKey")
	}

	return &homeserver_interop.SigningKey{
		KeyVersion: string(keyVersionBytes),
		PrivateKey: privKey,
	}, nil
}
