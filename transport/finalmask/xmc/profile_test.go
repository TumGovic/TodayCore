// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package xmc

import (
	"bytes"
	"testing"
)

func TestProfilesFromConfigRejectsEmpty(t *testing.T) {
	if _, err := profilesFromConfig(nil); err == nil {
		t.Fatal("expected empty profiles error")
	}
}

func TestProfilesFromConfig(t *testing.T) {
	uuid := bytes.Repeat([]byte{0x2a}, 16)
	profiles, err := profilesFromConfig([]*Profile{
		{
			Username:          "SignedUser",
			Uuid:              uuid,
			TexturesValue:     "textures-value",
			TexturesSignature: "textures-signature",
		},
	})
	if err != nil {
		t.Fatalf("build explicit profile: %v", err)
	}
	if len(profiles) != 1 || profiles[0].Username != "SignedUser" {
		t.Fatalf("unexpected profile: %+v", profiles)
	}
	if profiles[0].TexturesValue != "textures-value" || profiles[0].TexturesSignature != "textures-signature" {
		t.Fatalf("textures were not preserved: %+v", profiles[0])
	}
}
