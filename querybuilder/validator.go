package querybuilder

import "errors"

// ValidateQuery validates the MiniQuery before generating the Cypher query
func ValidateQuery(q MiniQuery) error {
	if q.From == "" {
		return errors.New("missing 'From' node type")
	}

	// Validate 'Where' conditions
	for key, value := range q.Where {
		if key == "" || value == nil {
			return errors.New("invalid 'Where' clause")
		}
	}

	// Validate relationships
	for _, r := range q.Relationships {
		if r.Type == "" || r.To == "" {
			return errors.New("missing relationship type or target node")
		}
	}

	// Validate aggregate, return, order, and limit...
	// Other validations as required

	return nil
}
