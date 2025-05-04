package querybuilder

import (
	"fmt"
	"sort"
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
		keys := make([]string, 0, len(q.Where))
		for k := range q.Where {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := q.Where[k]
			// Check if the key starts with "f" to signify a relationship node like f1, f2
			if strings.HasPrefix(k, "f") {
				// Handle filtering for relationship nodes (f1, f2, ...)
				for i := range q.Relationships {
					nextVar := fmt.Sprintf("f%d", i+1)
					// If the key matches a relationship node (f1, f2,...), build the condition
					if nextVar == k {
						switch val := v.(type) {
						case string:
							whereParts = append(whereParts, fmt.Sprintf("%s.%s = '%s'", nextVar, strings.Split(k, ".")[1], val))
						case []interface{}:
							var orConditions []string
							for _, item := range val {
								orConditions = append(orConditions, fmt.Sprintf("%s.%s = '%v'", nextVar, strings.Split(k, ".")[1], item))
							}
							whereParts = append(whereParts, "("+strings.Join(orConditions, " OR ")+")")
						default:
							whereParts = append(whereParts, fmt.Sprintf("%s.%s = %v", nextVar, strings.Split(k, ".")[1], v))
						}
					}
				}
			} else {
				// Handle filtering for root node (p)
				switch val := v.(type) {
				case string:
					whereParts = append(whereParts, fmt.Sprintf("p.%s = '%s'", k, val))
				case []interface{}:
					var orConditions []string
					for _, item := range val {
						orConditions = append(orConditions, fmt.Sprintf("p.%s = '%v'", k, item))
					}
					whereParts = append(whereParts, "("+strings.Join(orConditions, " OR ")+")")
				default:
					whereParts = append(whereParts, fmt.Sprintf("p.%s = %v", k, v))
				}
			}
		}
		// Join the conditions with AND
		matchClause += fmt.Sprintf(" WHERE %s", strings.Join(whereParts, " AND "))
	}

	// if len(q.Where) > 0 {
	// 	var whereParts []string
	// 	for k, v := range q.Where {
	// 		// Check if the key starts with "f" to indicate a relationship node
	// 		if strings.HasPrefix(k, "f") {
	// 			for i := range q.Relationships {
	// 				nextVar := fmt.Sprintf("f%d", i+1)
	// 				if nextVar == k {
	// 					switch val := v.(type) {
	// 					case string:
	// 						whereParts = append(whereParts, fmt.Sprintf("%s.%s = '%s'", nextVar, k, val))
	// 					case []interface{}:
	// 						var orConditions []string
	// 						for _, item := range val {
	// 							orConditions = append(orConditions, fmt.Sprintf("%s.%s = '%v'", nextVar, k, item))
	// 						}
	// 						whereParts = append(whereParts, "("+strings.Join(orConditions, " OR ")+")")
	// 					default:
	// 						whereParts = append(whereParts, fmt.Sprintf("%s.%s = %v", nextVar, k, v))
	// 					}
	// 				}
	// 			}
	// 		} else {
	// 			// For the root node "p", handle conditions as usual
	// 			switch val := v.(type) {
	// 			case string:
	// 				whereParts = append(whereParts, fmt.Sprintf("p.%s = '%s'", k, val))
	// 			case []interface{}:
	// 				var orConditions []string
	// 				for _, item := range val {
	// 					orConditions = append(orConditions, fmt.Sprintf("p.%s = '%v'", k, item))
	// 				}
	// 				whereParts = append(whereParts, "("+strings.Join(orConditions, " OR ")+")")
	// 			default:
	// 				whereParts = append(whereParts, fmt.Sprintf("p.%s = %v", k, v))
	// 			}
	// 		}
	// 	}
	// 	matchClause += fmt.Sprintf(" WHERE %s", strings.Join(whereParts, " AND "))
	// }

	// if len(q.Where) > 0 {
	// 	var whereParts []string
	// 	for k, v := range q.Where {
	// 		switch val := v.(type) {
	// 		case string:
	// 			// If it's a single string, create the condition as usual
	// 			whereParts = append(whereParts, fmt.Sprintf("p.%s = '%s'", k, val))
	// 		case []interface{}:
	// 			// If it's an array (multiple values), create OR conditions
	// 			var orConditions []string
	// 			for _, item := range val {
	// 				orConditions = append(orConditions, fmt.Sprintf("p.%s = '%v'", k, item))
	// 			}
	// 			whereParts = append(whereParts, "("+strings.Join(orConditions, " OR ")+")")
	// 		default:
	// 			// For other types, just create a condition
	// 			whereParts = append(whereParts, fmt.Sprintf("p.%s = %v", k, v))
	// 		}
	// 	}
	// 	matchClause += fmt.Sprintf(" WHERE %s", strings.Join(whereParts, " AND "))
	// }

	// if len(q.Where) > 0 {
	// 	var whereParts []string
	// 	for k, v := range q.Where {
	// 		switch val := v.(type) {
	// 		case string:
	// 			whereParts = append(whereParts, fmt.Sprintf("p.%s = '%s'", k, val))
	// 		default:
	// 			whereParts = append(whereParts, fmt.Sprintf("p.%s = %v", k, v))
	// 		}
	// 	}
	// 	matchClause += fmt.Sprintf(" WHERE %s", strings.Join(whereParts, " AND "))
	// }

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
