package public_test

import (
	"github.com/xihui-forever/mutualRead/public"
	"testing"
)

func TestPublic(t *testing.T) {
	buf, err := public.Public.ReadFile("build/index.html")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(len(buf))
}
