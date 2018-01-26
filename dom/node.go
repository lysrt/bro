package dom

const (
	ErrorNode   NodeType = "error"
	TextNode             = "text"
	ElementNode          = "element"
	CommentNode          = "comment"
)

type NodeType string

type Node struct {
	Parent      *Node
	FirstChild  *Node
	LastChild   *Node
	PrevSibling *Node
	NextSibling *Node

	Type NodeType

	Tag        string
	Attributes map[string]string

	TextContent string
}
