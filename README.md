<p align="center">
  <img src="/clank.png" height="100">
  <h2 align="center">
    Clank is your assistant for mock http & grpc services
  </h2>
</p>

## Features

- [x] Mock Http server 
- [x] Mock Grpc server
- [x] Template function
- [x] Get value from Http&Grpc request
- [x] Response data based on custom conditions
- [ ] Custom variable parameters
- [ ] Send Http or Grpc request after mock server response  


## build & install


## mock rpc
in the example directory, you can run the following command to mock rpc server

```
clank --yaml grpc_serv.yaml
```

look at the yaml 
```
kind: grpc
port: 54312
importPath:
  - ../internal/testdata/
protoPath:
  - protos/api/student_api.proto
servers:
  - name: api.StudentSrv
    methods:
      - name: StudentByID
        defaultResponse: '{"studentList":[{"name":"test","age":1},{"name":"test2","age":2}]}'
        conditions:
          - condition: '$request.id == 111'
            response: '{"studentList":[{"name":"test1111","age":111}]}'
          - condition: '"$header.x-header" == "test"'
            response: '{"studentList":[{"name":"header","age":222}]}'
          - condition: $request.id == 456
            response: |-
              {"studentList":[{"name":"{{RandFixLenString 3}}","id": {{RandInt64}},"age":{{ RandInt32 }}}, 
              	{"name":"{{RandString 3 10}}","id": {{RandInt64}},"age":{{ RandInt32 }}}, 
              	{"name":"{{RandString 3 10}}","id": {{RandInt64}},"age":{{ RandInt32 }}}]}
      - name: NewStudent
        defaultResponse: '{"code":"OK","desc":"OK","data":"eyJzdHVkZW50TGlzdCI6W3sibmFtZSI6InRlc3QiLCJhZ2UiOjF9LHsibmFtZSI6InRlc3QyIiwiYWdlIjoyfV19"}'
        conditions:
          - condition: '"$request.name" == "test"'
            response: '{"code":"OK","desc":"OKHAHA","data":"eyJzdHVkZW50TGlzdCI6W3sibmFtZSI6ImhlaWhlaSIsImFnZSI6MX0seyJuYW1lIjoiaGFoYWhhIiwiYWdlIjo5fV19"}'
          - condition: '"$request.name" == "abc" && $request.id == 111'
            response: '{"code":"OK","desc":"OKabc","data":"eyJzdHVkZW50TGlzdCI6W3sibmFtZSI6ImhlaWhlaSIsImFnZSI6MX0seyJuYW1lIjoiaGFoYWhhIiwiYWdlIjo5fV19"}'
          - condition: '"$header.x-header" == "test"'
            response: '{"code":"OK","desc":"OKHeader","data":"eyJzdHVkZW50TGlzdCI6W3sibmFtZSI6ImhlaWhlaSIsImFnZSI6MX0seyJuYW1lIjoiaGFoYWhhIiwiYWdlIjo5fV19"}'
      - name: AllStudent
        defaultResponse: '{"studentList": [{"id":111,"name":"abc","age":1298498081},{"id":222,"name":"def","age":2019727887}]}'

```
