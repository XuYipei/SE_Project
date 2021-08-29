package state

import (
	"context"

	"github.com/tuplz/tuplz-be/internal/pkg/tools"
	"go.mongodb.org/mongo-driver/mongo"
)

func extractCur(cur *mongo.Cursor) []StateStruct {
	result := make([]StateStruct, 0)
	for cur.Next(context.Background()) {
		var state StateStruct
		cur.Decode(&state)
		result = append(result, state)
	}
	return result
}

func generateKey() string {
	stateCount += 1
	return tools.GenerateMd5(stateCount)
}
