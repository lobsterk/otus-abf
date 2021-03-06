package bucket

import (
	"github.com/lobsterk/otus-abf/internal/repositories/mock"
	log2 "github.com/lobsterk/otus-abf/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestBucket_Hold(t *testing.T) {
	log2.Logger, _ = getLogger("logs/logs.log", "info")

	asRep := mock.BucketsRepository{
		Data: map[string]uint{"test": 1},
	}
	bucket := NewBucket("test", &asRep, 3)
	ok, err := bucket.Hold("test")
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Error("ok")
	}
	if asRep.Data["test"] != 2 {
		t.Error("count not increment")
	}

	ok, err = bucket.Hold("test")
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Error("ok")
	}
	if asRep.Data["test"] != 3 {
		t.Error("count not increment")
	}

	ok, err = bucket.Hold("test")
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("ok")
	}
	if asRep.Data["test"] != 3 {
		t.Error("count not increment")
	}
}

func TestBucket_Drop(t *testing.T) {
	log2.Logger, _ = getLogger("logs/logs.log", "info")

	asRep := mock.BucketsRepository{
		Data: map[string]uint{"test": 1},
	}
	bucket := NewBucket("test", &asRep, 3)

	if err := bucket.Clear("test"); err != nil {
		t.Error(err)
	}
	if _, ok := asRep.Data["test"]; ok {
		t.Error("bucket not clear")
	}

}

func getLogger(logFile, logLevel string) (logger *zap.Logger, err error) {
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	return zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(level),
		OutputPaths: []string{"stdout", logFile},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message", // <--
		},
	}.Build()
}
