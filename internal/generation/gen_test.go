package generation

import "testing"

type testNodeKind uint32

func (t testNodeKind) String() string {
	return "STD__NodeFoo"
}

type testTokenKind uint32

func (t testTokenKind) String() string {
	return "TokMisleadingString"
}

func TestManifestNodeKey_prefersStringerOverRawUint(t *testing.T) {
	if got := manifestNodeKey(testNodeKind(32)); got != "STD__NodeFoo" {
		t.Fatalf("manifestNodeKey(testNodeKind(32)) = %q, want STD__NodeFoo (not decimal)", got)
	}
}

func TestManifestTokenKey_prefersRawUintEvenWithStringer(t *testing.T) {
	if got := manifestTokenKey(testTokenKind(32)); got != "32" {
		t.Fatalf("manifestTokenKey(testTokenKind(32)) = %q, want 32", got)
	}
}

func TestManifestTokenKey_plainUint32_noStringer(t *testing.T) {
	var u uint32 = 32
	if got := manifestTokenKey(u); got != "32" {
		t.Fatalf("manifestTokenKey(uint32(32)) = %q, want 32", got)
	}
}
