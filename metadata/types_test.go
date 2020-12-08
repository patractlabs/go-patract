package metadata_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/patractlabs/go-patract/metadata"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	var raw metadata.Raw

	bz, _ := ioutil.ReadFile("../test/contracts/ink/erc721.json")

	err := json.Unmarshal(bz, &raw)
	assert.Nil(t, err)

	for _, rt := range raw.Types {
		t.Logf("type %v", rt.Path)
		for k, v := range rt.Def {
			t.Logf("def %s", k)
			t.Logf("def value %v", string(v))
		}
	}
}
