
# go-neo4j-query-builder

## Overview

This project, **go-neo4j-query-builder**, was developed as part of **Vibe Coding**, where I explored various ways to simplify the process of generating **Neo4j Cypher queries**. The idea is to abstract the complexity of Cypher query syntax into a more user-friendly format, allowing developers to easily build graph queries without needing to directly write Cypher.

The core functionality is to enable users to define graph queries in a simplified JSON format, which are then automatically translated into corresponding Cypher queries that can be executed against a Neo4j database.

This project was developed with the help of AI tools, particularly ChatGPT, which assisted in designing, refining, and optimizing the approach to generating Cypher queries. AI played a significant role in streamlining the development process and ensuring the code's efficiency.

## Project Features

- **Simplified Query Language**: Allows users to define graph queries in a simplified JSON format.
- **Dynamic Cypher Query Generation**: Converts the JSON input into a valid Cypher query.
- **Customizable Relationships**: Supports defining various relationships between nodes.
- **Command-Line Interface**: Can be run via the command line to accept JSON input and generate Cypher queries.

## Getting Started

### Prerequisites

Ensure that you have the following before running the project locally:

- Go 1.16+ installed on your system.
- A running Neo4j instance (locally or remotely).
- A code editor or IDE to make any modifications.

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/TharunNo1/go-neo4j-query-builder.git
   cd go-neo4j-query-builder
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the application:

   ```bash
   go build
   ```

### Usage

To generate a Cypher query, you can pass either a JSON file or JSON string via the command line.

**Example 1: Using a JSON file**
```bash
go run main.go -json=query.json
```

**Example 2: Using a JSON string**
```bash
go run main.go -data='{"from": "Person", "where": {"name": "John"}, "relationships": [{"type": "KNOWS", "to": "Person"}], "return": ["p.name", "f.name"]}'
```

### Example JSON Input

```json
{
  "from": "Person",
  "where": {
    "name": "John"
  },
  "relationships": [
    {
      "type": "KNOWS",
      "to": "Person"
    }
  ],
  "return": ["p.name", "f.name"]
}
```

### JSON Fields:

- **from**: The label of the starting node (e.g., `Person`).
- **where**: Conditions to filter the nodes (e.g., `name = 'John'`).
- **relationships**: Defines the relationships between nodes (e.g., `KNOWS` from `Person` to `Person`).
- **return**: Specifies which fields to return in the result (e.g., `p.name`, `f.name`).

## Project Structure

```plaintext
go-neo4j-query-builder/
├── main.go             # Main entry point for the application
├── querybuilder/       # Package responsible for generating Cypher queries
│   ├── builder.go      # Core logic for building Cypher queries
│   ├── model.go        # Definitions for query models (MiniQuery, Relationship)
│   └── validate.go     # Query validation logic
├── README.md           # Project documentation
└── go.mod              # Go module definition
```

## How It Works

1. **Parsing JSON Input**: The user provides a JSON file or string that defines the query.
2. **Query Validation**: The JSON input is validated to ensure that all required fields are present and correct.
3. **Cypher Query Generation**: The `BuildCypher` function takes the validated query and converts it into a Cypher query string.
4. **Output**: The generated Cypher query is printed and can be used to interact with a Neo4j database.

## Contributing

Feel free to fork this repository, submit pull requests, or open issues if you find any bugs or have suggestions for new features.

### How to Contribute

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes and commit them (`git commit -am 'Add new feature'`).
4. Push to your branch (`git push origin feature-branch`).
5. Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

This project was developed as part of my journey in **Vibe Coding**, where I aimed to create a tool that simplifies working with Neo4j by converting a simplified JSON format into Cypher queries. The development process was significantly aided by AI, particularly ChatGPT, which helped me design the query structure, generate code, and optimize the project.
