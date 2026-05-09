# Mock Server
## Package
`github.com/jycodepub/golibs/net`

## Type
### 1. `Server`
Fields:
- `addr string`
- `handler http.HandlerFunc`

Methods:
- `func NewMockServer(addr string, configPath string, respdir string) *Server`
- `func (s *Server) Start()`

## Configuration
### Config File Format
```json
[
  {
    "method": "GET",
    "url": "/test",
    "contentType": "application/json",
    "responses": ["response1", "response2"],
    "responseFiles": ["file1", "file2"]
  }
]
```

## Examples
```go
server := net.NewMockServer("localhost:8080", "config.json", "responses" )
server.Start()
```
