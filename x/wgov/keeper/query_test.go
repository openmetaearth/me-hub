package keeper

import (
	"os"
	"strings"
	"testing"
)

func TestMeTallyResultUsesCacheContextForActiveProposalTally(t *testing.T) {
	source, err := os.ReadFile("query.go")
	if err != nil {
		t.Fatalf("read query.go: %v", err)
	}

	text := string(source)
	if !strings.Contains(text, "cacheCtx, _ := ctx.CacheContext()") {
		t.Fatalf("MeTallyResult must create a cache context before tallying active proposals")
	}
	if !strings.Contains(text, "q.Tally(cacheCtx, proposal)") {
		t.Fatalf("MeTallyResult must tally active proposals against cacheCtx")
	}
	if strings.Contains(text, "q.Tally(ctx, proposal)") {
		t.Fatalf("MeTallyResult must not tally active proposals against the write-through ctx")
	}
}
