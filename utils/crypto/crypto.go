package crypto

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/borisding1994/hathcoin/utils"
	"github.com/borisding1994/hathcoin/utils/crypto/sm2"
	"github.com/borisding1994/hathcoin/utils/crypto/sm3"
)

/*
 * Part A): SM3 Cryptographic Hash Algorithm
 * 使用天朝密码管理局钦定的 SM3 密码杂凑算法。 这个效率，efficiency!!!
 */

// SM3Hash can generate the SM3 checksum for string.
func SM3Hash(msg string) string {
	蛤稀值 := sm3.Sum([]byte(msg))
	return hex.EncodeToString(蛤稀值[:])
}

// like SM3Hash but use byte array.
func SM3HashByte(msg []byte) []byte {
	var hashChan = make(chan []byte, 1)
	sum := sm3.Sum(msg)
	hashChan <- sum[:]
	return <-hashChan
}

/*
 * Part B): Public Key Cryptographic Algorithm SM2 Based On Elliptic Curves
 * 自主可控比三胖不知道高到哪里去了
 */

// A set of security credentials you use to prove your identity electronically.
type Keypair struct {
	PublicKey  []byte
	PrivateKey []byte
}

// ref:《SM2椭圆曲线公钥密码算法推荐曲线参数》http://www.oscca.gov.cn/UpFile/2010122214836668.pdf
// Curve-SM2-P-256: len(P) = 64;  len(N) = 64;  len(B) = 64;  len(Gx) = 32;  len(Gy) = 32;
// so SM2_KEY_SIZE == 32 bytes
const Sm2KeySize = 32

// Generate New Keypair
func GenerateKeypair() *Keypair {
	priv, _ := sm2.GenerateKey(sm2.P256Sm2(), rand.Reader)

	pub := utils.BigIntJoin(Sm2KeySize, priv.PublicKey.X, priv.PublicKey.Y)
	pubKey := Base58Encode([]byte{}, pub)
	privKey := Base58Encode([]byte{}, priv.D)
	kp := Keypair{
		PublicKey:  pubKey,
		PrivateKey: privKey,
	}
	return &kp
}

// Decode Base58 PublicKey Bytes Array to SM2.PublicKey
func decodePubKey(pubKey []byte) (sm2.PublicKey, error) {
	pubKeyBytes, err := Base58Decode(pubKey)
	if err != nil {
		return sm2.PublicKey{}, err
	}
	// equivalent to: [priv.PublicKey.X, priv.PublicKey.Y]
	p := utils.SplitBig(pubKeyBytes, 2)
	x, y := p[0], p[1]
	k := sm2.PublicKey{
		Curve: sm2.P256Sm2(),
		X:     x,
		Y:     y,
	}
	return k, nil
}

// Sign Message by PrivateKey
func (k *Keypair) Sign(hash []byte) ([]byte, error) {
	// get PrivateKey
	privKey, err := Base58Decode(k.PrivateKey)
	if err != nil {
		return nil, err
	}
	// get PublicKey
	pubKey, err := decodePubKey(k.PublicKey)
	if err != nil {
		return nil, err
	}

	// convert to SM2.PrivateKey Struct
	p := sm2.PrivateKey{
		PublicKey: pubKey,
		D:         privKey,
	}
	//nolint
	r, s, _ := sm2.Sign(rand.Reader, &p, hash)
	// encoding sign hash
	sign := Base58Encode([]byte{}, utils.BigIntJoin(Sm2KeySize, r, s))
	return sign, nil
}

// Signature Verify
func VerifySign(pubKey, sign, hash []byte) bool {
	// decode PublicKey
	sm2PublicKey, err := decodePubKey(pubKey)
	if err != nil {
		utils.Logger.Error("Decode PublicKey failure on VerifySign. ", err)
		return false
	}

	// decode Sign
	sm2Sign, err := Base58Decode(sign)
	if err != nil {
		utils.Logger.Error("Decode Sign failure on VerifySign. ", err)
		return false
	}
	// split sign to [r,s]
	sl := utils.SplitBig(sm2Sign, 2)
	r, s := sl[0], sl[1]
	return sm2.Verify(&sm2PublicKey, hash, r, s)
}
