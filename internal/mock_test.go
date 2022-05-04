package internal

import (
	"regexp"
	"strings"
	"testing"

	"github.com/Knetic/govaluate"
	jsonIter "github.com/json-iterator/go"
)

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
								Condition: "$request.obj.name == 123 || $request.id == 456",
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
	t.Log(schema.Validate())
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

func TestValuate1(t *testing.T) {
	expression, err := govaluate.NewEvaluableExpression("10 > 0")
	if err != nil {
		t.Fatal(err)
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

	expression, err = govaluate.NewEvaluableExpression(`10 > "afd"`)
	if err != nil {
		t.Fatal(err)
	}
	result, err = expression.Evaluate(nil)
	if err != nil {
		t.Log(err)
	}
	t.Log(result)

	parameters := make(map[string]interface{}, 8)
	parameters["$request.id"] = -1
	str := `$request.id > 0 || $request.id == -1 || $request.name == "test"`
	var re = regexp.MustCompile(`\$request.(?P<parameter>[.\w]+)`)
	match := re.FindAllStringSubmatch(str, -1)
	idx := re.SubexpIndex("parameter")
	for _, matchItem := range match {
		t.Log(matchItem[idx])
	}
	str = strings.ReplaceAll(str, "$request.id", "123")
	str = strings.ReplaceAll(str, "$request.name", `"test"`)
	t.Log(str)
	expression, err = govaluate.NewEvaluableExpression(str)
	if err != nil {
		t.Fatal(err)
	}
	result, err = expression.Evaluate(parameters)
	if err != nil {
		t.Log(err)
	}
	t.Log(result)
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

func TestJson(t *testing.T) {
	val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"], "info":{"name":"test","age":1}}`)
	t.Log(jsonIter.Get(val, "Colors", 0).ToString())
	t.Log(jsonIter.Get(val, "info", "name").ToString())
	t.Log(jsonIter.Get(val, "info", "age").ToInt())
}
