package akismet

import (
	"testing"
)

func TestVerifyKey(t *testing.T) {
	ak := &Akismet{"key", "http://example.com", TypePadDomain}
	v, err := ak.CommentCheck(Comment{UserIp: "1.2.3.4", UserAgent: "Chrome", CommentContent: "This is a test"})
	if err != nil {
		t.Fatal(err)
	}
	println(v)
}
