package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/protobuf/proto"
)

const (
	Etcd2NodeFilePrefix = "."
)

func transform(n *client.Node) *mvccpb.KeyValue {
	if n.Dir {
		// fmt.Fprintf(os.Stderr, "Dir: %s\n", n.Key)
		return nil
	}

	key := strings.Replace(n.Key, fmt.Sprintf("/%s", Etcd2NodeFilePrefix), "/", 1)

	value, err := base64.StdEncoding.DecodeString(n.Value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not base64 decode key %s: %s\n", n.Key, err)
		return nil
	}

	if n.TTL != 0 {
		fmt.Fprintf(os.Stderr, "TTL key: %s\n", n.Key)
	} else {
		fmt.Fprintf(os.Stderr, "Key: %s\n", n.Key)
	}
	kv := &mvccpb.KeyValue{
		Key:            []byte(key),
		Value:          []byte(value),
		CreateRevision: int64(n.CreatedIndex),
		ModRevision:    int64(n.ModifiedIndex),
		// We can't get version from etcd2 nodes. Assuming all KVs has version 1.
		Version: 1,
	}
	return kv
}

func ExitError(err error) {
	fmt.Fprintf(os.Stderr, "Err: %v\n", err)
	os.Exit(1)
}

func main() {
	reader := os.Stdin
	decoder := json.NewDecoder(reader)
	writer := os.Stdout

	fmt.Fprintf(os.Stderr, "Starting...\n")
	defer fmt.Fprintf(os.Stderr, "Exiting...\n")
	defer os.Stdout.Close()

	buf := make([]byte, 8)
	for {
		node := &client.Node{}
		if err := decoder.Decode(node); err != nil {
			if err == io.EOF {
				return
			}
			ExitError(err)
		}

		kv := transform(node)
		if kv == nil {
			continue
		}

		data, err := proto.Marshal(kv)
		if err != nil {
			ExitError(err)
		}
		binary.LittleEndian.PutUint64(buf, uint64(len(data)))
		if _, err := writer.Write(buf); err != nil {
			ExitError(err)
		}
		if _, err := writer.Write(data); err != nil {
			ExitError(err)
		}
	}
}
