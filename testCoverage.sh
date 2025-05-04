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
