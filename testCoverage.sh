#!/bin/bash

echo "🚀 Running CLI Test Scenarios..."
echo

total=0
passed=0

# Define tests as: description | expected | input JSON string (escaped)

tests=(
  "Simple filter|MATCH (p:Person) WHERE p.name = 'Alice' RETURN p.name|{\"from\":\"Person\",\"where\":{\"name\":\"Alice\"},\"return\":[\"p.name\"]}"
  "One hop relationship|MATCH (p:Person)-[:KNOWS]->(f1:Person) WHERE p.name = 'Alice' RETURN p.name, f1.name|{\"from\":\"Person\",\"relationships\":[{\"type\":\"KNOWS\",\"to\":\"Person\"}],\"where\":{\"name\":\"Alice\"},\"return\":[\"p.name\",\"f1.name\"]}"
  "Two hop relationship|MATCH (p:Person)-[:KNOWS]->(f1:Person)-[:KNOWS]->(f2:Person) WHERE p.name = 'Alice' RETURN p.name, f1.name, f2.name|{\"from\":\"Person\",\"relationships\":[{\"type\":\"KNOWS\",\"to\":\"Person\"},{\"type\":\"KNOWS\",\"to\":\"Person\"}],\"where\":{\"name\":\"Alice\"},\"return\":[\"p.name\",\"f1.name\",\"f2.name\"]}"
  "Only start node|MATCH (p:City) RETURN p|{\"from\":\"City\"}"
  "Numeric filter|MATCH (p:Employee) WHERE p.age = 30 RETURN p.name, p.age|{\"from\":\"Employee\",\"where\":{\"age\":30},\"return\":[\"p.name\",\"p.age\"]}"
  "Auto return with one relationship|MATCH (p:Person)-[:FRIEND]->(f1:Person) RETURN p, f1|{\"from\":\"Person\",\"relationships\":[{\"type\":\"FRIEND\",\"to\":\"Person\"}]} "
  "Empty query|missing 'From' node type|{\"from\":\"\",\"where\":{\"name\":\"Alice\"},\"return\":[\"p.name\"]}"
  "Complex WHERE clause|MATCH (p:Person) WHERE p.age = 30 AND p.name = 'Alice' RETURN p.name, p.age|{\"from\":\"Person\",\"where\":{\"name\":\"Alice\",\"age\":30},\"return\":[\"p.name\",\"p.age\"]}"
  "Multiple conditions in WHERE|MATCH (p:Person) WHERE p.age = 30 AND p.city = 'NYC' AND p.name = 'Alice' RETURN p.name, p.age, p.city|{\"from\":\"Person\",\"where\":{\"name\":\"Alice\",\"age\":30,\"city\":\"NYC\"},\"return\":[\"p.name\",\"p.age\",\"p.city\"]}"
  "Multiple hops (3 hops)|MATCH (p:Person)-[:KNOWS]->(f1:Person)-[:KNOWS]->(f2:Person)-[:KNOWS]->(f3:Person) WHERE p.name = 'Alice' RETURN p.name, f1.name, f2.name, f3.name|{\"from\":\"Person\",\"relationships\":[{\"type\":\"KNOWS\",\"to\":\"Person\"},{\"type\":\"KNOWS\",\"to\":\"Person\"},{\"type\":\"KNOWS\",\"to\":\"Person\"}],\"where\":{\"name\":\"Alice\"},\"return\":[\"p.name\",\"f1.name\",\"f2.name\",\"f3.name\"]}"
  "Complex relationship and WHERE with multiple filters|MATCH (p:Person)-[:KNOWS]->(f1:Person)-[:LIVES_IN]->(f2:City) WHERE p.name = 'Alice' AND f2.name = 'New York' RETURN p.name, f1.name, f2.name|{\"from\":\"Person\",\"relationships\":[{\"type\":\"KNOWS\",\"to\":\"Person\"},{\"type\":\"LIVES_IN\",\"to\":\"City\"}],\"where\":{\"name\":\"Alice\",\"f2.name\":\"New York\"},\"return\":[\"p.name\",\"f1.name\",\"f2.name\"]}"
  "Query with OR condition|MATCH (p:Person) WHERE (p.name = 'Alice' OR p.name = 'Bob') RETURN p.name|{\"from\":\"Person\",\"where\":{\"name\":[\"Alice\", \"Bob\"]},\"return\":[\"p.name\"]}"
  "Invalid From field (missing node type)|missing 'From' node type|{\"from\":\"\",\"where\":{\"name\":\"Alice\"},\"return\":[\"p.name\"]}"
  "Relationship with multiple types|MATCH (p:Person)-[:KNOWS|:FRIEND]->(f:Person) WHERE p.name = 'Alice' RETURN p.name, f.name|{"from":"Person","relationships":[{"type":"KNOWS","to":"Person"},{"type":"FRIEND","to":"Person"}],"where":{"name":"Alice"},"return":["p.name","f1.name"]}"
  "One hop with filter on related node|MATCH (p:Person)-[:KNOWS]->(f1:Person) WHERE p.name = 'Alice' AND f1.age = 30 RETURN p.name, f1.name|{\"from\":\"Person\",\"relationships\":[{\"type\":\"KNOWS\",\"to\":\"Person\"}],\"where\":{\"name\":\"Alice\",\"f1.age\":30},\"return\":[\"p.name\",\"f1.name\"]}"
  "Empty result|MATCH (p:Person) WHERE p.name = 'NonExistent' RETURN p.name|{\"from\":\"Person\",\"where\":{\"name\":\"NonExistent\"},\"return\":[\"p.name\"]}"
  "Invalid relationship type|missing relationship type or target node|{\"from\":\"Person\",\"relationships\":[{\"type\":\"\",\"to\":\"Person\"}],\"where\":{\"name\":\"Alice\"},\"return\":[\"p.name\"]}"
  "Empty WHERE clause|MATCH (p:Person) RETURN p|{\"from\":\"Person\",\"where\":{},\"return\":[\"p\"]}"
  "Filter with nested map|MATCH (p:Person) WHERE p.address.city = 'New York' RETURN p.name, p.address|{\"from\":\"Person\",\"where\":{\"address.city\":\"New York\"},\"return\":[\"p.name\",\"p.address\"]}"
  "Complex RETURN with dynamic nodes|MATCH (p:Person)-[:KNOWS]->(f1:Person)-[:LIVES_IN]->(f2:City) WHERE p.name = 'Alice' RETURN p.name, f1.name, f2.name|{\"from\":\"Person\",\"relationships\":[{\"type\":\"KNOWS\",\"to\":\"Person\"},{\"type\":\"LIVES_IN\",\"to\":\"City\"}],\"where\":{\"name\":\"Alice\"},\"return\":[\"p.name\",\"f1.name\",\"f2.name\"]}"
  "Join two nodes with multiple relationships|MATCH (p:Person)-[:KNOWS]->(f1:Person)-[:KNOWS]->(f2:Person) WHERE p.name = 'Alice' RETURN p.name, f1.name, f2.name|{\"from\":\"Person\",\"relationships\":[{\"type\":\"KNOWS\",\"to\":\"Person\"},{\"type\":\"KNOWS\",\"to\":\"Person\"}],\"where\":{\"name\":\"Alice\"},\"return\":[\"p.name\",\"f1.name\",\"f2.name\"]}"
)


# Path to your compiled Go CLI binary
BIN=./neo4jQueryBuilder

for test in "${tests[@]}"; do
    IFS='|' read -r desc expected input <<< "$test"
    ((total++))
    echo "[$total] $desc"

    result=$($BIN -data="$input")

    if [[ "$result" == "$expected" ]]; then
        echo "  ✅ Passed"
        ((passed++))
    else
        echo "  ❌ Failed"
        echo " Command: $BIN -data="$input""
        echo "    Expected: $expected"
        echo "    Got     : $result"
    fi
    echo
done

# Summary
coverage=$((100 * passed / total))
echo "🧪 Total Passed: $passed / $total"
echo "📊 Custom Coverage: $coverage%"

if [[ $coverage -eq 100 ]]; then
    echo "🎉 All tests passed with full coverage!"
else
    echo "⚠️ Some tests failed. Coverage below 100%."
fi
