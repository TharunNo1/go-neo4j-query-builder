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

	// Step 2: Generate Cypher query after validation
	var queryParts []string

	// MATCH Clause (start with node)
	matchClause := fmt.Sprintf("MATCH (p:%s)", q.From)

	// Relationships Clause (if present)
	if len(q.Relationships) > 0 {
		var relationshipClauses []string
		for _, r := range q.Relationships {
			relationshipClauses = append(relationshipClauses, fmt.Sprintf("-[:%s]->(f:%s)", r.Type, r.To))
		}
		matchClause = fmt.Sprintf("%s %s", matchClause, strings.Join(relationshipClauses, ""))
	}

	// WHERE Clause (if present) - AFTER the MATCH and RELATIONSHIPS clauses
	var whereParts []string
	for k, v := range q.Where {
		if str, ok := v.(string); ok {
			whereParts = append(whereParts, fmt.Sprintf("p.%s = '%s'", k, str))
		} else {
			whereParts = append(whereParts, fmt.Sprintf("p.%s = %v", k, v))
		}
	}

	if len(whereParts) > 0 {
		matchClause = fmt.Sprintf("%s WHERE %s", matchClause, strings.Join(whereParts, " AND "))
	}

	// Add the MATCH clause (with relationships and WHERE)
	queryParts = append(queryParts, matchClause)

	// RETURN Clause (if present) - Ensure that both p.name and f.name are included
	if len(q.Return) > 0 {
		// Join the fields in Return, ensure both p and f fields are correctly included
		returnFields := strings.Join(q.Return, ", ")
		queryParts = append(queryParts, fmt.Sprintf("RETURN %s", returnFields))
	} else {
		// Default to returning the nodes and their properties if no explicit fields in Return
		queryParts = append(queryParts, "RETURN p, f")
	}

	// Return the final Cypher query
	return strings.Join(queryParts, " "), nil
}
