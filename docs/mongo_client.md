# MongoDB Client
## Package
- `github.com/jycodepub/golibs/mongo`

## Type
### 1. `Client`
Methods:
- `func NewClient(connectString string) *Client`
- `func (c *Client) Close()`
- `func (c *Client) Insert(ctx context.Context, database string, collection string, document interface{}) error`
- `func (c *Client) Query(ctx context.Context, database string, collection string, filter interface{}) (*mongo.Cursor, error)`
- `func (c *Client) GetCollection(database string, collection string) *mongo.Collection`
- `func (c *Client) QueryForStruct(ctx context.Context, database string, collection string, filter interface{}, result interface{}) error`
- `func (e *DataNotFound) Error() string`