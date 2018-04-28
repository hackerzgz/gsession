package gsession

import (
	"testing"

	"github.com/hackez/gsession/ticker"
)

func TestSession_LoginTime(t *testing.T) {
	s := Session("202CB962AC59075B964B07152D234B7072883001")
	tim, err := s.LoginTime()
	if err != nil {
		t.Fatal("failed to format login time:", err)
	}

	t.Log(tim)

	if tim != "Apr 26 15:47:10" {
		t.Fatal("not ")
	}
}

func TestSessionPool_Generate(t *testing.T) {

	pool := New(ticker.TenSecTicker, 20)

	session, err := pool.Generate("1046")
	if err != nil {
		t.Fatal("failed to generate session pool:", err)
	}

	t.Log("session:", session)

	pool.Free()
}

func BenchmarkSessionPool_Generate(b *testing.B) {
	pool := New(ticker.TenSecTicker, 20)
	defer pool.Free()
	// 1442 ns/op => 1s 生成 693 个ID

	for i := 0; i < b.N; i++ {
		_, err := pool.Generate("1046")
		if err != nil {
			b.Fatal("failed to generate session pool:", err)
		}
	}

}
