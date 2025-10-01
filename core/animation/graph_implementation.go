//go::build js && wasm

package animation

import (
	"fmt"
	"log"
	"time"
)

// NewGraph creates a new animation graph
func NewGraph() *Graph {
	return &Graph{
		Nodes:      make(map[string]*Node),
		Edges:      make(map[string][]*Edge),
		Parameters: make(map[string]interface{}),
		IsPlaying:  false,
	}
}

// AddNode add an animation node to the graph
func (g *Graph) AddNode(node *Node) {
	g.Nodes[node.Name] = node
	if g.Edges[node.Name] == nil {
		g.Edges[node.Name] = make([]*Edge, 0)
	}
}

// AddEdge add an animation edge to the graph
func (g *Graph) AddEdge(edge *Edge) error {
	// validate node exists
	if _, exists := g.Nodes[edge.From]; !exists {
		return fmt.Errorf("node %s does not exist", edge.From)
	}
	if _, exists := g.Nodes[edge.To]; !exists {
		log.Printf("Node %s does not exist", edge.To)
		return fmt.Errorf("node %s does not exist", edge.To)
	}

	g.Edges[edge.From] = append(
		g.Edges[edge.From],
		edge,
	)

	return nil
}

// SetNode immediately transition to the node
func (g *Graph) SetNode(nodeName string) error {
	newNode, exists := g.Nodes[nodeName]
	if !exists {
		return fmt.Errorf("node %s does not exist", nodeName)
	}

	// Exit the current node
	if g.CurrentNode != nil && g.CurrentNode.OnExit != nil {
		g.CurrentNode.OnExit()
	}

	// Enter the new node
	oldNodeName := ""
	if g.CurrentNode != nil {
		oldNodeName = g.CurrentNode.Name
	}

	g.CurrentNode = newNode
	g.CurrentNode.LocalTime = 0
	g.CurrentNode.IsActive = true

	if g.CurrentNode.OnEnter != nil {
		g.CurrentNode.OnEnter()
	}

	// Trigger callbacks
	if g.OnNodeChanged != nil {
		g.OnNodeChanged(oldNodeName, nodeName)
	}

	return nil
}

// Update advances the graph by delta time
func (g *Graph) Update(deltaTime time.Duration) {
	if !g.IsPlaying || g.CurrentNode == nil {
		return
	}

	g.GlobalTime += deltaTime
	g.CurrentNode.LocalTime += deltaTime

	// Update the current node
	if g.CurrentNode.OnUpdate != nil {
		g.CurrentNode.OnUpdate(deltaTime)
	}

	// Check for transitions
	g.checkEdges()

	// Update animation clip
	if g.CurrentNode.Clip != nil {
		g.UpdateAnimationClip(g.CurrentNode, deltaTime)
	}
}

// UpdateAnimationClip updates the animation clip
func (g *Graph) UpdateAnimationClip(node *Node, deltaTime time.Duration) {
	clip := node.Clip
	if clip == nil {
		return
	}

	// Update the current time
	node.LocalTime += deltaTime

	// Check if the animation is complete
	if node.LocalTime >= clip.Duration {
		node.LocalTime = 0
		if clip.IsLooping {
			node.LocalTime = 0
		} else {
			node.IsActive = false
		}
	}

}

// checkEdges checks for transitions
func (g *Graph) checkEdges() {
	if g.CurrentNode == nil {
		return
	}

	edges := g.Edges[g.CurrentNode.Name]
	for _, edge := range edges {
		if edge.Condition.Evaluate(g) {
			g.StartEdgeSwitch(edge)
			break
		}
	}
}

// StartEdgeSwitch starts a transition
func (g *Graph) StartEdgeSwitch(edge *Edge) {
	edge.IsActive = true
	edge.Progress = 0

	if g.OnEdge != nil {
		g.OnEdge(edge)
	}

	// For immediate transitions, just switch to the target node
	if edge.Duration == 0 {
		g.SetNode(edge.To)
		return
	}

	// TODO: Handle blended transitions
	g.handleBlendedEdges(edge)
}

// handleBlendedEdges handles blended transitions
// TODO: Implement blended transitions
func (g *Graph) handleBlendedEdges(edge *Edge) {
}
