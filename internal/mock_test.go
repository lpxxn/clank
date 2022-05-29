package internal

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/Knetic/govaluate"
)

const testPort int = 54312

func TestServerDesc(t *testing.T) {

	grpcSchema := &grpcSchema{
		ImportPath: []string{"./testdata/"},
		ProtoPath:  []string{"protos/api/student_api.proto"},
		Servers: GrpcServerDescriptionList{
			&GrpcServerDescription{
				Name: "api.StudentSrv",
				Methods: []*GrpcMethodDescription{
					&GrpcMethodDescription{
						Name:            "NewStudent",
						DefaultResponse: `{"code":"OK","desc":"OK","data":"eyJzdHVkZW50TGlzdCI6W3sibmFtZSI6InRlc3QiLCJhZ2UiOjF9LHsibmFtZSI6InRlc3QyIiwiYWdlIjoyfV19"}`,
						DefaultMetaData: map[string]string{
							"userID":   "123",
							"userName": "test",
						},
						Conditions: []*ResponseConditionDescription{
							&ResponseConditionDescription{
								Condition: `"$request.name" == "test"`,
								Response:  `{"code":"OK","desc":"OKHAHA","data":"eyJzdHVkZW50TGlzdCI6W3sibmFtZSI6ImhlaWhlaSIsImFnZSI6MX0seyJuYW1lIjoiaGFoYWhhIiwiYWdlIjo5fV19"}`,
							},

							&ResponseConditionDescription{
								Condition: `"$request.name" == "abc" && $request.id == 111`,
								Response:  `{"code":"OK","desc":"OKabc","data":"eyJzdHVkZW50TGlzdCI6W3sibmFtZSI6ImhlaWhlaSIsImFnZSI6MX0seyJuYW1lIjoiaGFoYWhhIiwiYWdlIjo5fV19"}`,
							},
						},
					},
					&GrpcMethodDescription{ /// {"id":1}
						Name:            "StudentByID",
						DefaultResponse: `{"studentList":[{"name":"test","age":1},{"name":"test2","age":2}]}`,
						Conditions: []*ResponseConditionDescription{
							&ResponseConditionDescription{
								Condition: "$request.id == 111",
								Response:  `{"studentList":[{"name":"test1111","age":111}]}`, // `{"a": adfe}`,
							},

							&ResponseConditionDescription{
								Condition: `$request.id == 456`,
								Response: `{"studentList":[{"name":"{{RandFixLenString 3}}","id": {{RandInt64}},"age":{{ RandInt32 }}}, 
															{"name":"{{RandString 3 10}}","id": {{RandInt64}},"age":{{ RandInt32 }}}, 
															{"name":"{{RandString 3 10}}","id": {{RandInt64}},"age":{{ RandInt32 }}}]}`,
							},
						},
					},
				},
			},
		},
	}
	schema := &SchemaDescription{
		SchemaDescriptionBase: SchemaDescriptionBase{
			Kind: GRPC,
			Port: testPort,
		},
		Server: grpcSchema,
	}
	t.Log(schema)
	t.Log(schema.Validate())
	ser, err := ParseServerMethodsFromProto(grpcSchema.ImportPath, grpcSchema.ProtoPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateGrpcServiceInputAndOutput(grpcSchema.Servers, ser); err != nil {
		t.Fatal(err)
	}
	if err := SetOutputFunc(grpcSchema.Servers, ser); err != nil {
		t.Fatal(err)
	}
	if err := ser.StartWithPort(schema.Port); err != nil {
		t.Fatal(err)
	}
}

func TestServer2(t *testing.T) {
	ser, err := ParseServerMethodsFromProto([]string{"./testdata/"}, []string{"protos/api/student_api.proto"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ser.StartWithPort(testPort); err != nil {
		t.Fatal(err)
	}
}

func TestTemplate(t *testing.T) {
	str1 := `{"studentList":[{"name":"{{RandFixLenString 3}}","ids": {{RandInt64Slice 5}},"age":{{ RandInt32 }}}, 
															{"name":"{{RandString 3 10}}","id": {{RandInt64}},"age":{{ RandInt32 }}}, 
															{"name":"{{RandString 3 10}}","id": {{RandInt64}},"age":{{ RandInt32 }}}]}`

	b, err := GenerateDefaultTemplate(str1)
	t.Logf("body: %s, err: %+v", string(b), err)

	str2 := `{"ids":{{RandInt64Slice 10}}}`
	b, err = GenerateDefaultTemplate(str2)

	t.Logf("body: %s, err: %+v", string(b), err)

	temp, err := template.New("").Parse(`{{range .DataFields}}{{println "," .}} {{end}}`)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := temp.Execute(&buf, map[string]interface{}{"DataFields": []string{"A", "B", "C"}}); err != nil {
		t.Fatal(err)
	}
	t.Log(buf.String())
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
	str := `$request.id > 0 || $request.id == -1 || "$request.name" == "test" || $request.t == false`
	var re = regexp.MustCompile(`\$request.(?P<parameter>[.\w]+)`)
	match := re.FindAllStringSubmatch(str, -1)
	idx := re.SubexpIndex("parameter")
	for _, matchItem := range match {
		t.Log(matchItem[idx])
	}

	str = strings.ReplaceAll(str, "$request.id", "123")
	str = strings.ReplaceAll(str, "$request.name", `test`)
	str = strings.ReplaceAll(str, "$request.t", `true`)
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

func TestRe(t *testing.T) {
	grpcParamRegex = regexp.MustCompile(`.*{{.*}}.*`)
	s := grpcParamRegex.FindAllString("asdfP{{asdf", -1)
	t.Log(s, len(s))

	s = grpcParamRegex.FindAllString("{{asdf}}", -1)
	t.Log(s, len(s))

	s = grpcParamRegex.FindAllString("}}asdfasdf", -1)
	t.Log(s, len(s))

	s = grpcParamRegex.FindAllString("}}asdfasdf", -1)
	t.Log(s, len(s))

	s = grpcParamRegex.FindAllString("asdfas{{asdf}}sdafas{{asdfs}}asdf}}asdf", -1)
	t.Log(s, len(s))

	m := grpcParamRegex.MatchString("asdfP{{asdf")
	t.Log(m)

	m = grpcParamRegex.MatchString("{{asdf}}")
	t.Log(m)

	m = grpcParamRegex.MatchString("}}asdfasdf")
	t.Log(m)

	m = grpcParamRegex.MatchString("asdfP{{asdf")
	t.Log(m)

	m = grpcParamRegex.MatchString("asdfas{{asdf}}sdafas{{asdfs}}asdf}}asdf")
	t.Log(m)

}

func TestServer3(t *testing.T) {
	ser, err := ParseServerMethodsFromProtoset("./testdata/protos/test.protoset")
	if err != nil {
		t.Fatal(err)
	}
	if err := ser.StartWithPort(testPort); err != nil {
		t.Fatal(err)
	}
}

func TestJson(t *testing.T) {
	val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"], "info":{"name":"test","age":1}}`)
	t.Log(jsonIterator.Get(val, "Colors", 0).ToString())
	t.Log(jsonIterator.Get(val, "info", "name").ToString())
	g1 := jsonIterator.Get(val, "info", "name")
	t.Log(g1.GetInterface(), g1.ValueType(), g1.ToString())
	g1 = jsonIterator.Get(val, "info", "age")
	t.Log(g1.GetInterface(), g1.ValueType(), g1.ToString(), fmt.Sprintf("%v", g1.GetInterface()))
	t.Log(strings.ReplaceAll(string(val), "ID", fmt.Sprintf("%v", g1.GetInterface())))
	t.Log(strings.ReplaceAll("$request.id == 111", "$request.id", fmt.Sprintf("%v", g1.GetInterface())))
	g1 = jsonIterator.Get(val, "info", "age1")
	t.Log(g1)

	type mk struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	m := map[mk][]int{}
	mk1 := mk{ID: 1, Name: "Reds"}
	mk2 := mk{ID: 2, Name: "Blue"}
	if _, ok := m[mk1]; !ok {
		m[mk1] = []int{1}
	}
	if _, ok := m[mk1]; ok {
		m[mk1] = append(m[mk1], 2)
	}
	if _, ok := m[mk2]; !ok {
		m[mk2] = []int{1}
	}
	if _, ok := m[mk2]; ok {
		m[mk2] = append(m[mk2], 2)
	}

	t.Log(m)
}
