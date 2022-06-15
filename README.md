<p align="center">
  <img src="/clank.png" height="100">
  <h2 align="center">
    Clank is your assistant for mock http & grpc services
  </h2>
</p>

## Features

- [x] Mock Http server 
- [x] Mock Grpc server
- [x] Returns the specified response data according to custom conditions  
- [x] Template function
- [x] Get value from Http&Grpc request
- [x] Send Http request after mock server response  


## build & install
```
go install github.com/lpxxn/clank/clank@latest
```

### build from source
```
git clone git@github.com:lpxxn/clank.git
cd clank
make build
```

## run clank
in the example directory run `clank` 
```
CLANK_LOG_LEVEL=info clank --yaml grpc_serv.yaml
```
![run clank](/asset/clank_run.png)

## run in docker
```
docker run --rm -v $(pwd):$(pwd) -v $GOPATH:$GOPATH -w $(pwd) -p 54312:54312 lpxxn/clank:latest --yaml ./grpc_serv.yaml 
```
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


### conditions
default response is required, use `condition` to define conditions, you must specify the response too, when the condition is matched, the response will be returned, if not, the default response will be returned.    
eg:
```
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
```
if `id` field in the request is `111`, the response will be `{"studentList":[{"name":"test1111","age":111}]}`    
supported values
* $request.xxx   get value from request data, eg: `$request.id` `$request.name`
* $header.xxx    get value from header metadata eg: `$header.clientID`
* $response.xxx  get value from response data eg: `$response.code` `$response.data.id`

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
supported values
* $param.xxx   get value from url data, eg: url: `/user/:userID/order/:orderNo` you can use `$param.userID` to get the value
* $body.xxx    get value from post body eg: body is `{"name": "hello"}` use `$body.name` to get value
* $query.xxx   get value from url query data eg: `/user?id=123` use `$query.id` to get value
* $form.xxx    get value from form
* $header.xxx  get value from request header
* $response.xxx get value from response body

### Send Http request after mock server response
you can use  `httpCallback` to send http request after mock server response
```
httpCallback:
  - method: GET
    url: https://github.com/lpxxn/clank?userName=$request.name
    body: |-
      {"desc": $response.desc, "data": "$response.data"}
```
you can use server method values to get value and send customer data.

### Log
use `CLANK_LOG_LEVEL=info` environment variable to set log level, default is `error`    
eg: `CLANK_LOG_LEVEL=info clank --yaml http_serv.yaml`
