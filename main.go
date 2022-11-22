package main

import "fmt"

type Token struct {
	typeName string
	value    string
}

// abstract syntax tree node type declaration
type ASTNode struct {
	typeName     string
	value        string
	params       []ASTNode //child of node in AST
	currentValue int       //used to store next index value used in recursive function
}

// abstract syntax tree type declaration
type AST struct {
	typeName string
	body     []ASTNode
}

// list of all language types
var typeNames = []string{"NumberLiteral",
	"WordOperator",
	"CallBackExpression",
	"Root",
	"Terminator",
}

// used to allow funtions to be passed into traverse function to determine
//how each type of node should be transformed when added to the new AST
type ASTFunction func(ASTNode, ASTNode, AST)
type ASTMethods map[string]ASTFunction

func contains(slice []string, targetValue string) bool {
	// check if slice contains value in it

	for _, value := range slice {
		if value == targetValue {
			return true
		}
	}

	return false
}

func tokenize(inputExpression string) []Token {
	// generate list of tokens from input string
	var current int = 0
	var tokens []Token

	var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	var digits = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	var whitespace = " "

	for current < len(inputExpression) {
		var char string = string(inputExpression[current])

		if char == whitespace {
			current++
			continue // skip whitespace
		} else if char == "(" {
			// append token for left parenthesis
			newToken := Token{
				typeName: "left_paren",
				value:    "(",
			}

			tokens = append(tokens, newToken)
			current++
			continue
		} else if char == ")" {
			// append token for right parenthesis
			newToken := Token{
				typeName: "right_paren",
				value:    ")",
			}

			tokens = append(tokens, newToken)
			current++
			continue
		} else if contains(digits, char) {
			// if a digit is found, iterate through the next characters to store
			//consecutive digits in the same token
			var digitToken string

			for contains(digits, char) {
				digitToken += char
				current++
				char = string(inputExpression[current])
			}

			newToken := Token{
				typeName: "number",
				value:    digitToken,
			}

			tokens = append(tokens, newToken)
			continue
		} else if contains(alphabet, char) {
			// if a latter is found, iterate through the next characters to store
			//consecutive letters in the same token
			var word string

			for contains(alphabet, char) {
				word += char
				current++
				char = string(inputExpression[current])
			}

			newToken := Token{
				typeName: "word",
				value:    word,
			}

			tokens = append(tokens, newToken)
			continue
		} else {
			panic("Invalid character detected " + char)
		}
	}

	return tokens
}

func parse(tokens []Token, current int) ASTNode {
	// recursive base case
	if current >= len(tokens)-1 {
		terminator := ASTNode{
			typeName:     "Terminator",
			value:        "",
			params:       []ASTNode{},
			currentValue: -1,
		}

		return terminator
	}

	token := tokens[current]

	if token.typeName == "number" {
		newNode := ASTNode{
			typeName:     "NumberLiteral",
			value:        token.value,
			params:       []ASTNode{},
			currentValue: current + 1,
		}

		return newNode
	} else if token.typeName == "word" {
		newNode := ASTNode{
			typeName:     "WordOperator",
			value:        token.value,
			params:       []ASTNode{},
			currentValue: current + 1,
		}

		return newNode
	} else if token.typeName == "left_paren" {
		// if a left parenthesis is detected we open a new expression

		current++ // skip past parenthesis
		token = tokens[current]

		// create new node
		newNode := ASTNode{
			typeName:     "CallExpression",
			value:        token.value,
			params:       []ASTNode{},
			currentValue: 0,
		}

		// go to next token and start to build an abstract tree using recursion
		current++
		token = tokens[current]

		for token.typeName != "right_paren" {
			newNode.params = append(newNode.params, parse(tokens, current))
			current++
			token = tokens[current]
		}

		newNode.currentValue = current + 1
		return newNode
	}

	panic("Unexpected token found: " + token.typeName)
}

func generateAST(tokens []Token) AST {
	var ast AST
	ast.typeName = "Program"

	// add root node to start of AST
	rootNode := ASTNode{
		typeName:     "Root",
		value:        "",
		params:       []ASTNode{},
		currentValue: 0,
	}
	ast.body = append(ast.body, rootNode)

	// get current node from the parse function
	currentNode := parse(tokens, 0)

	for currentNode.typeName != "Terminator" {
		// append current node to the body of the AST
		ast.body = append(ast.body, currentNode)

		// use the currentValue field inside the node to determine the index
		//that must be used to call the parse function again
		currentNode = parse(tokens, currentNode.currentValue)
	}

	return ast
}

func traverseAST(tree AST, methods ASTMethods) {
	// extract all of the nodes from the tree body
	nodes := tree.body

	// for each node in the tree, run the method assigned to its type
	//and do the same for its children in a recursive manner
	traverseArray(nodes, nodes[0], methods) //nodes[0] is the root node hence it is the parent
}

func traverseNode(node ASTNode, parent ASTNode, methods ASTMethods) {
	// add current node to the new tree by calling assigned method
	method := methods[node.typeName]

	if method != nil {
		method(node, parent)
	}

	// if node is call expression, call traverse array function recursively to
	// add each parameter to the new AST
	if node.typeName == "CallExpression" {
		traverseArray(node.params, node, methods)
	} else if !contains(typeNames, node.typeName) {
		panic("Type Error (Invalid type): " + node.typeName)
	}
}

// iterate through all of the children of a node
func traverseArray(nodes []ASTNode, parent ASTNode, methods ASTMethods) {
	for _, node := range nodes {
		traverseNode(node, parent, methods)
	}
}


func NumberLiteralTraverse(node ASTNode, parent ASTNode, newTree AST){
	newNode := ASTNode{
		typeName: "NumberLiteral", 
		value: node.value, 
	}

	newTree.body = append(newTree.body, newNode)
}

func WordOperatorTraverse(node ASTNode, parent ASTNode, newTree AST){
	newNode := ASTNode{
		typeName: "WordOperator", 
		value: node.value, 
	}

	newTree.body = append(newTree.body, newNode)
}


func CallExpressionTraverse(node ASTNode, parent ASTNode, tree AST){
	
}

// transform tree into new AST
func transform(tree AST) {
	methods := map[string]ASTFunction{
		"NumberLiteral":  NumberLiteralTraverse,
		"WordOperator": WordOperatorTraverse,
		"CallExpression": CallExpressionTraverse,
	}
}


// generate code in new language from transformed AST
func generateCode() {

}


func main() {
	expression := "(123 234 456) add (123)"

	tokens := tokenize(expression)
	ast := generateAST(tokens)

	methods := map[string]ASTFunction{
		"NumberLiteral":  func(node1 ASTNode, node2 ASTNode) { fmt.Print(node1) },
		"WordOperator":   func(node1 ASTNode, node2 ASTNode) { fmt.Print(node1) },
		"CallExpression": func(node1 ASTNode, node2 ASTNode) { fmt.Print(node1) },
	}

	traverseAST(ast, methods)
}
