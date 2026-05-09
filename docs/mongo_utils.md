# MongoDB Utils
## Package
- `github.com/jycodepub/golibs/mongo`

## Functions
- `func ListCollections(connectionUri string, database string) []string`
- `func CleanDB(connectionUri string, database string)`
- `func DumpDB(connectionUri string, database string, outputDir string)`
- `func RestoreDB(connectionUri string, database string, outputDir string)`
- `func ExportCollection(connectionUri string, database string, collectionName string, outputDir string)`
- `func ImportCollection(connectionUri string, database string, collectionName string, inputFile string)`
- `func CleanCollection(connectionUri string, database string, collection string)`
- `func DropCollections(connectionUri string, database string)`
- `func DropCollection(connectionUri string, database string, collection string)`
