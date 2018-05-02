package objutils

import (
	"testing"
)

func Test_UnderLineToCamel(t *testing.T) {
	src := "user_id_id"
	dest := "UserIdId"

	out, _ := UnderLineToCamel(src)
	if dest != out {
		t.Error(out + "," + dest)
	}

	src = "user5"
	dest = "User5"

	out, _ = UnderLineToCamel(src)
	if dest != out {
		t.Error(out + "," + dest)
	}
}

func Test_CamelToUnderLine(t *testing.T) {
	src := "UserId"
	dest := "user_id"

	out, _ := CamelToUnderLine(src)
	if dest != out {
		t.Error(out + "," + dest)
	}

	src = "User5"
	dest = "user5"

	out, _ = CamelToUnderLine(src)
	if dest != out {
		t.Error(out + "," + dest)
	}
}
