package snowflake

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// Init should be called on application startup
func Init(machineID int64) {
	var err error
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		log.Fatalf("Snowflake Init Failed: %v", err)
	}
}

// GenerateID returns a globally unique 64-bit identifier
func GenerateID() int64 {
	if node == nil {
		Init(1) // Fallback for auto-init
	}
	return node.Generate().Int64()
}
