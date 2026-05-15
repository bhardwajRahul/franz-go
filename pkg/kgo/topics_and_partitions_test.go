package kgo

import (
	"testing"
)

func TestDelTopicsWithCoexistingPartitionPauses(t *testing.T) {
	// delTopics must clear all=true even when partition-level pauses exist for
	// the same topic. Previously the modified pausedPartitions struct was never
	// written back, leaving all=true in the map when len(pps.m) > 0.
	m := make(pausedTopics)

	// Add a topic-level pause and a partition-level pause for the same topic.
	m.addTopics("topic-a")
	m.addPartitions(map[string][]int32{"topic-a": {0}})

	if !m["topic-a"].all {
		t.Fatal("expected topic-a to be topic-paused after addTopics")
	}

	// Removing the topic-level pause must clear all=true even though partition
	// pauses still exist.
	m.delTopics("topic-a")

	if pps, ok := m["topic-a"]; ok && pps.all {
		t.Fatal("delTopics did not clear all=true when partition-level pauses coexist")
	}

	// The partition-level pause must still be present.
	if _, ok := m["topic-a"].m[0]; !ok {
		t.Fatal("delTopics unexpectedly removed the partition-level pause")
	}
}
