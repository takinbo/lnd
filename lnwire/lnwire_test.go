package lnwire

import (
	"encoding/hex"

	"net"

	"github.com/roasbeef/btcd/btcec"
	"github.com/roasbeef/btcd/chaincfg/chainhash"
	"github.com/roasbeef/btcd/txscript"
	"github.com/roasbeef/btcd/wire"
)

// Common variables and functions for the message tests

var (
	revHash = [32]byte{
		0xb7, 0x94, 0x38, 0x5f, 0x2d, 0x1e, 0xf7, 0xab,
		0x4d, 0x92, 0x73, 0xd1, 0x90, 0x63, 0x81, 0xb4,
		0x4f, 0x2f, 0x6f, 0x25, 0x88, 0xa3, 0xef, 0xb9,
		0x6a, 0x49, 0x18, 0x83, 0x31, 0x98, 0x47, 0x53,
	}

	maxUint32 uint32 = (1 << 32) - 1
	maxUint24 uint32 = (1 << 24) - 1
	maxUint16 uint16 = (1 << 16) - 1

	// For debugging, writes to /dev/shm/
	// Maybe in the future do it if you do "go test -v"
	WRITE_FILE = false
	filename   = "/dev/shm/serialized.raw"

	// preimage: 9a2cbd088763db88dd8ba79e5726daa6aba4aa7e
	// echo -n | openssl sha256 | openssl ripemd160 | openssl sha256 | openssl ripemd160
	revocationHashBytes, _ = hex.DecodeString("4132b6b48371f7b022a16eacb9b2b0ebee134d41")
	revocationHash         [20]byte

	// preimage: "hello world"
	redemptionHashBytes, _ = hex.DecodeString("5b315ebabb0d8c0d94281caa2dfee69a1a00436e")
	redemptionHash         [20]byte

	// preimage: "next hop"
	nextHopBytes, _ = hex.DecodeString("94a9ded5a30fc5944cb1e2cbcd980f30616a1440")
	nextHop         [20]byte

	privKeyBytes, _ = hex.DecodeString("9fa1d55217f57019a3c37f49465896b15836f54cb8ef6963870a52926420a2dd")
	privKey, pubKey = btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	address         = pubKey

	//  Delivery PkScript
	// Privkey: f2c00ead9cbcfec63098dc0a5f152c0165aff40a2ab92feb4e24869a284c32a7
	// PKhash: n2fkWVphUzw3zSigzPsv9GuDyg9mohzKpz
	deliveryPkScript, _ = hex.DecodeString("76a914e8048c0fb75bdecc91ebfb99c174f4ece29ffbd488ac")

	//  Change PkScript
	// Privkey: 5b18f5049efd9d3aff1fb9a06506c0b809fb71562b6ecd02f6c5b3ab298f3b0f
	// PKhash: miky84cHvLuk6jcT6GsSbgHR8d7eZCu9Qc
	changePkScript, _ = hex.DecodeString("76a914238ee44bb5c8c1314dd03974a17ec6c406fdcb8388ac")

	// echo -n | openssl sha256
	// This stuff gets reversed!!!
	shaHash1Bytes, _ = hex.DecodeString("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	shaHash1, _      = chainhash.NewHash(shaHash1Bytes)
	outpoint1        = wire.NewOutPoint(shaHash1, 0)
	// echo | openssl sha256
	// This stuff gets reversed!!!
	shaHash2Bytes, _ = hex.DecodeString("01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b")
	shaHash2, _      = chainhash.NewHash(shaHash2Bytes)
	outpoint2        = wire.NewOutPoint(shaHash2, 1)
	// create inputs from outpoint1 and outpoint2
	inputs = []*wire.TxIn{wire.NewTxIn(outpoint1, nil, nil), wire.NewTxIn(outpoint2, nil, nil)}

	// Commitment Signature
	tx           = wire.NewMsgTx(1)
	emptybytes   = new([]byte)
	sigStr, _    = txscript.RawTxInSignature(tx, 0, *emptybytes, txscript.SigHashAll, privKey)
	commitSig, _ = btcec.ParseSignature(sigStr, btcec.S256())

	// Funding TX Sig 1
	sig1privKeyBytes, _ = hex.DecodeString("927f5827d75dd2addeb532c0fa5ac9277565f981dd6d0d037b422be5f60bdbef")
	sig1privKey, _      = btcec.PrivKeyFromBytes(btcec.S256(), sig1privKeyBytes)
	sigStr1, _          = txscript.RawTxInSignature(tx, 0, *emptybytes, txscript.SigHashAll, sig1privKey)
	commitSig1, _       = btcec.ParseSignature(sigStr1, btcec.S256())
	// Funding TX Sig 2
	sig2privKeyBytes, _ = hex.DecodeString("8a4ad188f6f4000495b765cfb6ffa591133a73019c45428ddd28f53bab551847")
	sig2privKey, _      = btcec.PrivKeyFromBytes(btcec.S256(), sig2privKeyBytes)
	sigStr2, _          = txscript.RawTxInSignature(tx, 0, *emptybytes, txscript.SigHashAll, sig2privKey)
	commitSig2, _       = btcec.ParseSignature(sigStr2, btcec.S256())
	// Slice of Funding TX Sigs
	ptrFundingTXSigs = append(*new([]*btcec.Signature), commitSig1, commitSig2)

	// TxID
	txid = new(chainhash.Hash)
	// Reversed when displayed
	txidBytes, _ = hex.DecodeString("fd95c6e5c9d5bcf9cfc7231b6a438e46c518c724d0b04b75cc8fddf84a254e3a")
	_            = copy(txid[:], txidBytes)

	someAlias, _ = NewAlias("012345678901234567890")
	someSig, _   = btcec.ParseSignature(sigStr, btcec.S256())
	someSigBytes = someSig.Serialize()

	someAddress = &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8333}

	someChannelID = ChannelID{
		BlockHeight: maxUint24,
		TxIndex:     maxUint24,
		TxPosition:  maxUint16,
	}

	someRGB = RGB{
		red:   255,
		green: 255,
		blue:  255,
	}
)
