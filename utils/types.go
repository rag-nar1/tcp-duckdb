package utils

var typeMapping = map[string]string{
	// Numeric Types
	"SMALLINT":        	"SMALLINT",
	"INTEGER":         	"INTEGER",
	"INT":            	"INTEGER",
	"BIGINT":          	"BIGINT",
	"DECIMAL":         	"DECIMAL",
	"NUMERIC":         	"DECIMAL",
	"REAL":            	"REAL",
	"FLOAT4":          	"REAL",
	"DOUBLE PRECISION": "DOUBLE",
	"FLOAT8":          	"DOUBLE",
	"SERIAL":          	"INTEGER AUTOINCREMENT",
	"BIGSERIAL":       	"BIGINT AUTOINCREMENT",

	// Character Types
	"VARCHAR":  "VARCHAR",
	"TEXT":     "TEXT",
	"CHAR":     "VARCHAR",

	// Date & Time Types
	"DATE":                  	"DATE",
	"TIME":                  	"TIME",
	"TIMESTAMP":             	"TIMESTAMP",
	"TIMESTAMPTZ":           	"TIMESTAMP WITH TIME ZONE",
	"TIMESTAMP WITH TIME ZONE": "TIMESTAMP WITH TIME ZONE",

	// Boolean Type
	"BOOLEAN": "BOOLEAN",

	// JSON & Array Types
	"JSON":  "JSON",
	"JSONB": "JSON",
	"ARRAY": "UNSUPPORTED", // Needs conversion

	// Other Types
	"UUID":  "VARCHAR(36)",
	"BYTEA": "BLOB",
	"INET":  "VARCHAR",
	"CIDR":  "VARCHAR",
	"ENUM":  "VARCHAR",
}

func DbTypeMap(postgres string) string {
	return typeMapping[postgres]
}