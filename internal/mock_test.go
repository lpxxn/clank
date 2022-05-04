package internal

import "testing"

const testPort int = 54312

func TestServerDesc(t *testing.T) {
	schema := &SchemaDescription{
		Kind: GRPC,
		Port: testPort,
		Servers: []*ServerDescription{
			&ServerDescription{
				Name: "api.StudentSrv",
				Methods: []*MethodDescription{
					&MethodDescription{
						Name:            "NewStudent",
						DefaultResponse: `{"code":"OK","desc":"OK","data":"eyJzdHVkZW50TGlzdCI6W3sibmFtZSI6InRlc3QiLCJhZ2UiOjF9LHsibmFtZSI6InRlc3QyIiwiYWdlIjoyfV19"}`,
					},
					&MethodDescription{ /// {"id":1}
						Name:            "StudentByID",
						DefaultResponse: `{"studentList":[{"name":"test","age":1},{"name":"test2","age":2}]}`,
						Conditions: []*ResponseConditionDescription{
							&ResponseConditionDescription{
								Condition: "$request.id == 111",
								Response:  `{"studentList":[{"name":"test1111","age":111}]}`,
							},

							&ResponseConditionDescription{
								Condition: "$request.id == 123 || $request.id == 456",
								Response:  `{"studentList":[{"name":"123||456","age":123456}]}`,
							},
						},
					},
				},
			},
		},
		ImportPath: []string{"./testdata/"},
		ProtoPath:  []string{"protos/api/student_api.proto"},
	}
	t.Log(schema)
}

func TestServer2(t *testing.T) {
	ser, err := ParseServerMethodsFromProto([]string{"./testdata/"}, []string{"protos/api/student_api.proto"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ser.Start(testPort); err != nil {
		t.Fatal(err)
	}
}

func TestServer3(t *testing.T) {
	ser, err := ParseServerMethodsFromProtoset("./testdata/protos/test.protoset")
	if err != nil {
		t.Fatal(err)
	}
	if err := ser.Start(testPort); err != nil {
		t.Fatal(err)
	}
}
