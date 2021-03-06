package grpc

import (
	"context"
	"github.com/bxcodec/faker/v3"
	"github.com/lobsterk/otus-abf/internal/models"
	"github.com/lobsterk/otus-abf/internal/repositories/mock"
	"github.com/lobsterk/otus-abf/internal/services/bucket"
	"github.com/lobsterk/otus-abf/internal/services/ip_guard"
	"github.com/lobsterk/otus-abf/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	masksRepository := &mock.MasksRepository{
		Rows: []models.Mask{
			{Id: 1, Mask: "123.23.44.55/8", ListId: 1},
			{Id: 2, Mask: "122.27.44.55/8", ListId: 1},
		},
	}
	whiteList := ip_guard.NewMemoryIpGuard(1, masksRepository)
	blackList := ip_guard.NewMemoryIpGuard(1, masksRepository)

	bucketIpRep := &mock.BucketsRepository{Data: map[string]uint{"123.23.44.55": 2, "123.21.44.55": 1}}
	bucketLoginRep := &mock.BucketsRepository{Data: map[string]uint{"test_login_1": 1, "test_login_2": 1}}
	bucketPasswordRep := &mock.BucketsRepository{Data: map[string]uint{"test_password_1": 1, "test_password_2": 1}}

	bucketIp := bucket.NewBucket("ip", bucketIpRep, 3)
	bucketLogin := bucket.NewBucket("login", bucketLoginRep, 3)
	bucketPassword := bucket.NewBucket("password", bucketPasswordRep, 3)

	server := initServer(NewServer(whiteList, blackList, bucketIp, bucketLogin, bucketPassword, func(err string) {}))

	go server.Serve(getListener())
	defer server.Stop()

	client, cc := initClientTest()
	defer cc.Close()

	t.Run("Check", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()

		r, err := client.Check(ctx, &api.CheckRequest{Login: "test_login_3", Password: "test_password_3", Ip: "123.23.44.55"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		r, err = client.Check(ctx, &api.CheckRequest{Login: "test_login_3", Password: "test_password_3", Ip: "123.23.44.55"})
		if err != nil {
			handlerError(err, t)
		}
		if r.Success {
			t.Error("not success")
		}
	})

	t.Run("AddWhiteMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = whiteList.DropMask("123.23.40.55/4")
		r, err := client.AddWhiteMask(ctx, &api.AddWhiteMaskRequest{Mask: "123.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := whiteList.Contains("123.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("123.23.40.55 not found")
		}
	})

	t.Run("AddBlackMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = blackList.DropMask("123.23.40.55/4")
		r, err := client.AddBlackMask(ctx, &api.AddBlackMaskRequest{Mask: "123.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := blackList.Contains("123.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Error("123.23.40.55 not found")
		}
	})

	t.Run("DropBlackMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = blackList.AddMask("123.23.40.55/4")
		r, err := client.DropBlackMask(ctx, &api.DropBlackMaskRequest{Mask: "123.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := blackList.Contains("123.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if ok {
			t.Error("123.23.40.55 is found")
		}
	})

	t.Run("DropWhiteMask", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()
		_, _ = whiteList.AddMask("123.23.40.55/4")
		r, err := client.DropWhiteMask(ctx, &api.DropWhiteMaskRequest{Mask: "123.23.40.55/4"})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		ok, err := whiteList.Contains("123.23.40.55")
		if err != nil {
			t.Error(err)
		}
		if ok {
			t.Error("123.23.40.55 is found")
		}
	})

	t.Run("ClearBucket", func(t *testing.T) {
		ip := faker.IPv4()
		login := faker.Name()
		bucketIpRep.Data = map[string]uint{ip: 1}
		bucketLoginRep.Data = map[string]uint{login: 2}

		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
		defer cancel()

		r, err := client.ClearBucket(ctx, &api.ClearBucketRequest{Ip: ip, Login: login})
		if err != nil {
			handlerError(err, t)
		}
		if !r.Success {
			t.Error("not success")
		}
		if _, ok := bucketIpRep.Data[ip]; ok {
			t.Error("_, ok := bucketIpRep.Data[ip]; ok")
		}
		if _, ok := bucketLoginRep.Data[login]; ok {
			t.Error("_, ok := bucketLoginRep.Data[login]; ok")
		}
	})
}

func initServer(calendarServer api.AntiBruteForceServer) *grpc.Server {
	server := grpc.NewServer()
	reflection.Register(server)

	api.RegisterAntiBruteForceServer(server, calendarServer)
	return server
}

func getListener() net.Listener {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	return lis
}

func initClientTest() (api.AntiBruteForceClient, *grpc.ClientConn) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	return api.NewAntiBruteForceClient(cc), cc
}

func handlerError(err error, t *testing.T) {
	statusErr, ok := status.FromError(err)
	if ok {
		if statusErr.Code() == codes.DeadlineExceeded {
			t.Errorf("Deadline exceeded!")
		} else {
			t.Errorf("undexpected error %s", statusErr.Message())
		}
	} else {
		t.Errorf("Error while calling RPC CheckHomework: %v", err)
	}
}
