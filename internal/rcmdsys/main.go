package rcmdsys

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
)

func Upd(userId string, probId string, score float64) error {
	conn, err := grpc.Dial("localhost:5010", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return err
	}
	defer conn.Close()

	client := NewRcmdsysClient(conn)
	rep, err := client.Upd(context.Background(), &Record{
		UserId: userId,
		ProbId: probId,
		Score:  score,
	})
	if rep.Status != "success" || err != nil {
		err = errors.New("fail")
	}
	return err
}

func Query(userId string, maxLength int) ([]string, error) {
	conn, err := grpc.Dial("localhost:5010", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	if maxLength == 0 {
		maxLength = 10000
	}

	client := NewRcmdsysClient(conn)
	rep, err := client.Query(context.Background(), &UId{
		Id:        userId,
		MaxLength: int32(maxLength),
	})
	return rep.Id, err
}
