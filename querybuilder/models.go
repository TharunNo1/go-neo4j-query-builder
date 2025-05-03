package querybuilder

// MiniQuery represents the structure of the simplified query DSL
type MiniQuery struct {
	From         string
	Where        map[string]interface{}
	Relationships []Relationship
	Aggregate    string
	Return       []string
	OrderBy      []string
	Limit        int
	With         []string
	Create       string
	Merge        string
	Delete       string
	Union        []MiniQuery
	Optional     bool
}

// Relationship represents a relationship between nodes in Neo4j
type Relationship struct {
	Type       string
	To         string
	Properties map[string]interface{}
	BiDir      bool // For bidirectional relationships
}
