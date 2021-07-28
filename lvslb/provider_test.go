package lvslb_test

import (
	"testing"

	"github.com/jeremmfr/terraform-provider-lvslb/lvslb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := lvslb.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = lvslb.Provider()
}
