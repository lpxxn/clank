package internal

import "testing"

const testPort int = 54312

func TestServer1(t *testing.T) {
	ser, err := ParseServerMethodsFromProto([]string{"./testdata/"}, []string{"protos/api/student_api.proto"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ser.Start(54312); err != nil {
		t.Fatal(err)
	}
}
