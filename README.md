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


## mock rpc server
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
        defaultResponse: '{"studentList":[{"name":"test","age":1},{"name":"{{RandString 3 10}}","age":{{ RandInt32 }}}]}'
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

### proto and protoset file
you can use `importPath` and `protoPath` to import proto files
```
importPath:
  - ../internal/testdata/
protoPath:
  - protos/api/student_api.proto
```
or you can use `protosetPath` to import protoset files

```
protosetPath: ../internal/testdata/protos/test.protoset
```

### rpc server & method
you can use `servers` to define rpc server and method
```
servers:
  - name: api.StudentSrv
    methods:
      - name: StudentByID
        defaultResponse: '{"studentList":[{"name":"test","age":1},{"name":"{{RandString 3 10}}","age":{{ RandInt32 }}}]}'
```
`{{ RandInt32 }}` is a template func, you can use it to generate random int32 value    
`{{ RandInt64 }}` is a template func, you can use it to generate random int64 value   
`{{ RandString 3 10 }}` is a template func, you can use it to generate random string value   
`{{ RandFixLenString 3 }}` is a template func, you can use it to generate random fixed length string value    

use `condition` to define conditions, in `condition` you can use `$request.xxx` to get request data

```
        conditions:
          - condition: '$request.id == 111'
            response: '{"studentList":[{"name":"test1111","age":111}]}'
          - condition: '"$header.x-header" == "test"'
            response: '{"studentList":[{"name":"header","age":222}]}'
```
if `condition` is true, the response will be used.


## mock http server

in the example directory, you can run the following command to mock http server

```
clank --yaml http_serv.yaml
```

look at the yaml 
```
kind: http
port: 9527
server:
  methods:
    - name: testApi
      path: /test
      method: GET
      defaultResponse: |-
        {"code": 0,"message": "success",
           "data": {"name": "Jerry","age": 18}
        }
    - name: testApi2
      path: /user/:userID/order/:orderNo
      method: POST
      defaultResponse: |-
        {
          "code": 0, "message": "success",
          "data": {
            "orderNo": "$param.orderNo",
            "userID": $param.userID,
            "desc": "{{RandString 5 20}}"
                  }
        }
    - name: testApi3
      path: /user/:userID/createOrder
      method: POST
      defaultResponse: |-
        {
          "code": 0, "message": "success",
          "data": {
            "orderNo": "OrderNo{{RandString 5 10}}",
            "userID": $param.userID,
            "desc": "{{RandString 5 20}}"
          }
        }
      conditions:
        - condition: '$query.userID == 1'
          response: |-
            {
              "code": 0, "message": "success",
              "data": {
                "orderNo": "OrderNo{{RandString 5 10}}",
                "userID": $param.userID,
                "desc": "query.userID == 1"
              }
            }
        - condition: '$param.userID == 1'
          response: |-
            {
              "code": 0, "message": "success",
              "data": {
                "orderNo": "OrderNo{{RandString 5 10}}",
                "userID": $param.userID,
                "desc": "param.userID == 1"
              }
            }
        - condition: '$body.userID == 1 && $query.userID == 2'
          response: |-
            {
              "code": 0, "message": "success",
              "data": {
                "orderNo": "OrderNo{{RandString 5 10}}",
                "userID": $body.userID,
                "queryUserID": $query.userID,
                "desc": "body.userID == 1 && query.userID == 2"
              }
            }
        - condition: '$body.userID == 1'
          response: |-
            {
              "code": 0, "message": "success",
              "data": {
                "orderNo": "OrderNo{{RandString 5 10}}",
                "userID": $param.userID,
                "desc": "body.userID == 1"
              }
            }
```