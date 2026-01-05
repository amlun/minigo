package id

import (
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// Init initializes the global Snowflake node.
// It reads node id from env `SNOWFLAKE_NODE_ID`, default is 1.
func Init() {
	if node != nil {
		return
	}
	var nodeID int64 = 1
	if v := os.Getenv("SNOWFLAKE_NODE_ID"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			nodeID = parsed
		}
	}
	var err error
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("failed to init snowflake node: %v", err)
	}
}

// NextID returns next unique int64 id.
func NextID() int64 {
	if node == nil {
		Init()
	}
	return node.Generate().Int64()
}
