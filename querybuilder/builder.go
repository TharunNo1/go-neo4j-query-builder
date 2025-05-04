package querybuilder

import (
	"fmt"
	"strings"
)

// BuildCypher builds a Cypher query from a MiniQuery
func BuildCypher(q MiniQuery) (string, error) {
	// Step 1: Validate the query
	if err := ValidateQuery(q); err != nil {
		return "", err
	}

	// Step 2: Start constructing the Cypher query
	var queryParts []string

	// Start MATCH clause with the root node
	matchClause := fmt.Sprintf("MATCH (p:%s)", q.From)

	// Build relationship path with unique variable names
	for i, r := range q.Relationships {
		nextVar := fmt.Sprintf("f%d", i+1)
		matchClause += fmt.Sprintf("-[:%s]->(%s:%s)", r.Type, nextVar, r.To)
	}

	// Add WHERE clause if present
	if len(q.Where) > 0 {
		var whereParts []string
		for k, v := range q.Where {
			switch val := v.(type) {
			case string:
				whereParts = append(whereParts, fmt.Sprintf("p.%s = '%s'", k, val))
			default:
				whereParts = append(whereParts, fmt.Sprintf("p.%s = %v", k, v))
			}
		}
		matchClause += fmt.Sprintf(" WHERE %s", strings.Join(whereParts, " AND "))
	}

	// Add the MATCH clause
	queryParts = append(queryParts, matchClause)

	// Handle RETURN clause
	if len(q.Return) > 0 {
		queryParts = append(queryParts, fmt.Sprintf("RETURN %s", strings.Join(q.Return, ", ")))
	} else {
		// Auto-generate default return fields
		returnFields := []string{"p"}
		for i := range q.Relationships {
			returnFields = append(returnFields, fmt.Sprintf("f%d", i+1))
		}
		queryParts = append(queryParts, fmt.Sprintf("RETURN %s", strings.Join(returnFields, ", ")))
	}

	return strings.Join(queryParts, " "), nil
}
