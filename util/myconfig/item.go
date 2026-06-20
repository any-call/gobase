package myconfig

const (
	TypeString  = "string"
	TypeBool    = "bool"
	TypeInt     = "int"
	TypeInt8    = "int8"
	TypeInt16   = "int16"
	TypeInt32   = "int32"
	TypeInt64   = "int64"
	TypeUint    = "uint"
	TypeUint8   = "uint8"
	TypeUint16  = "uint16"
	TypeUint32  = "uint32"
	TypeUint64  = "uint64"
	TypeFloat32 = "float32"
	TypeFloat64 = "float64"
	TypeJSON    = "json"
	TypeNil     = "nil"
)

const rootKey = "$"

// Item 是配置模型转换后可用于表存储的一条数据。
type Item struct {
	Key   string
	Value string
	Type  string
}
