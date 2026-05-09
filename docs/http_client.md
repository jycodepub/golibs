# HttpClient Usages

## Package
- `github.com/jycodepub/golibs/net`

## Type
### 1. `HttpClient`
Methods:
- `func (c *HttpClient) Get(url string, ctx *RequestContext) (HttpResponse, error)`
- `func (c *HttpClient) Post(url string, ctx *RequestContext) (HttpResponse, error)`
- `func (c *HttpClient) Put(url string, ctx *RequestContext) (HttpResponse, error)`
- `func (c *HttpClient) Delete(url string, ctx *RequestContext) (HttpResponse, error)`
- `func (c *HttpClient) SubmitRequest(method string, url string, ctx *RequestContext) (HttpResponse, error)`

### 2. `RequestContext`
Fields:
- `Payload`
- `AuthToken`
- `Headers map[string]string`
Methods:
- `func (ctx *RequestContext) AddHeader(key string, value string) *RequestContext`
- `func (ctx *RequestContext) AddFormField(key string, value string) *RequestContext`
- `func (ctx *RequestContext) AddPayLoad(payload string) *RequestContext`
- `func (ctx *RequestContext) AddPayLoadFile(filepath string) *RequestContext`
- `func (ctx *RequestContext) AddToken(token string) *RequestContext`
- `func (ctx *RequestContext) AddTokenFile(filepath string) *RequestContext`

### 3. `HttpResponse`
Fields:
- `Code int`
- `Body string`

### 4. `PayLoad`
Fields:
- `Body string`
- `BodyFile string`
- `Form url.Values`

### 5. `AuthToken`
Fields:
- `Token string`
- `TokenFile string`

## Examples
### 1. Create `HttpClient`
```go
client := net.NewHttpClient()
```

### 2. Create `RequestContext`
```go
ctx := net.RequestContext{}
ctx.AddHeader("User-Agent", "Mozilla/5.0")
ctx.AddHeader("Accept", "text/html,application/json")
ctx.AddToken("token")
ctx.AddPayLoad("payload")
```

### 3. `GET`
```go
resp, err := client.Get("http://www.google.com", nil)
resp, err := client.Get("http://www.google.com", ctx)
```

### 4. `POST`
```go
resp, err := client.Post("http://www.google.com", ctx)
```

### 5. `PUT`
```go
resp, err := client.Put("http://www.google.com", ctx)
```

### 6. `DELETE`
```go
resp, err := client.Delete("http://www.google.com", nil)
resp, err := client.Delete("http://www.google.com", ctx)
```

