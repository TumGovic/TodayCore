// This file adapts key generation from Xray-core (main/commands/all/vlessenc.go
// and mldsa65.go), licensed under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package main

import (
	"crypto/ecdh"
	"crypto/mlkem"
	"crypto/rand"
	"encoding/base64"
	"os"
	"strings"

	"github.com/sagernet/sing-box/log"

	"github.com/cloudflare/circl/sign/mldsa/mldsa65"
	"github.com/spf13/cobra"
)

func init() {
	commandGenerate.AddCommand(commandGenerateVLESSEncryption)
	commandGenerate.AddCommand(commandGenerateMLDSA65KeyPair)
}

// Mirrors `xray vlessenc`.
var commandGenerateVLESSEncryption = &cobra.Command{
	Use:   "vlessenc",
	Short: "Generate VLESS Encryption decryption/encryption pair",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := generateVLESSEncryption()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func generateVLESSEncryption() error {
	x25519Key, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	var seed [64]byte
	_, err = rand.Read(seed[:])
	if err != nil {
		return err
	}
	mlkemKey, err := mlkem.NewDecapsulationKey768(seed[:])
	if err != nil {
		return err
	}
	encode := base64.RawURLEncoding.EncodeToString
	pair := func(serverKey, clientKey []byte) string {
		decryption := strings.Join([]string{"mlkem768x25519plus", "native", "600s", encode(serverKey)}, ".")
		encryption := strings.Join([]string{"mlkem768x25519plus", "native", "0rtt", encode(clientKey)}, ".")
		return "\"decryption\": \"" + decryption + "\"\n\"encryption\": \"" + encryption + "\"\n"
	}
	os.Stdout.WriteString("Choose one Authentication to use, do not mix them. Ephemeral key exchange is Post-Quantum safe anyway.\n\n")
	os.Stdout.WriteString("Authentication: X25519, not Post-Quantum\n")
	os.Stdout.WriteString(pair(x25519Key.Bytes(), x25519Key.PublicKey().Bytes()))
	os.Stdout.WriteString("\nAuthentication: ML-KEM-768, Post-Quantum\n")
	os.Stdout.WriteString(pair(seed[:], mlkemKey.EncapsulationKey().Bytes()))
	return nil
}

// Mirrors `xray mldsa65`: the seed is for an Xray REALITY server
// ("mldsa65Seed"), the verify key is for REALITY clients ("mldsa65_verify").
var commandGenerateMLDSA65KeyPair = &cobra.Command{
	Use:   "mldsa65-keypair",
	Short: "Generate REALITY ML-DSA-65 key pair",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := generateMLDSA65Key()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func generateMLDSA65Key() error {
	var seed [32]byte
	_, err := rand.Read(seed[:])
	if err != nil {
		return err
	}
	publicKey, _ := mldsa65.NewKeyFromSeed(&seed)
	os.Stdout.WriteString("Seed: " + base64.RawURLEncoding.EncodeToString(seed[:]) + "\n")
	os.Stdout.WriteString("Verify: " + base64.RawURLEncoding.EncodeToString(publicKey.Bytes()) + "\n")
	return nil
}
