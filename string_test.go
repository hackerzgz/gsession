package randstr

import "testing"

func TestNewString(t *testing.T) {

	newCase := []struct {
		Type charType
		Len  int
	}{
		{
			CharDigit | CharLowerCase | CharUpperCase,
			10,
		},
		{
			CharLowerCase | CharUpperCase,
			10,
		},
		{
			CharDigit,
			10,
		},
	}

	for _, n := range newCase {
		str, err := NewString(n.Type, n.Len)
		if err != nil {
			t.Fatal("failed to new random string: ", err)
		}

		t.Log(str)
	}
}
