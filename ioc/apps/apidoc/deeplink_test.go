package apidoc

import "testing"

func TestScalarDeepLink(t *testing.T) {
	got := ScalarDeepLink("/api/x/v1/apidoc/ui.html", "蓝图编排", "QueryBlueprint")
	want := "/api/x/v1/apidoc/ui.html#tag/蓝图编排/QueryBlueprint"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if ScalarOperationHash("", "QueryBlueprint") != "" {
		t.Fatal("empty tag should not invent hash")
	}
}
