package internal

import "testing"

const testPort int = 54312

func TestServer1(t *testing.T) {
	ser, err := ParseServerMethodsFromProto([]string{"./testdata/"}, []string{"protos/api/student_api.proto"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ser.Start(testPort); err != nil {
		t.Fatal(err)
	}
}

func TestServer2(t *testing.T) {
	ser, err := ParseServerMethodsFromProtoset("./testdata/protos/test.protoset")
	if err != nil {
		t.Fatal(err)
	}
	if err := ser.Start(testPort); err != nil {
		t.Fatal(err)
	}
}
