package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

const key = "counter"

type counterResp struct {
	CurrentValue int64 `json:"current_value"`
}

func writeCounterResp(w io.Writer, resp counterResp) {
	marshaled, err := json.Marshal(resp)
	if err != nil {
		log.Errorln("failed to marshal counter response", err)
	}
	w.Write(marshaled)
}

func writeInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}

func countGetHandler(w http.ResponseWriter, r *http.Request) {
	conn := RedisConnContext(r.Context())
	defer conn.Close()

	cmd := conn.Get(r.Context(), key)
	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Errorln("failed when get query", err)
		writeInternalServerError(w)
		return
	}

	var (
		resp counterResp
		err  error
	)

	if cmd.Err() == nil {
		resp.CurrentValue, err = cmd.Int64()
		if err != nil {
			log.Errorln("failed when obtain get result from redis", err)
			writeInternalServerError(w)
			return
		}
	} else {
		resp.CurrentValue = 0
	}

	writeCounterResp(w, resp)
}

func countAddHandler(w http.ResponseWriter, r *http.Request) {
	conn := RedisConnContext(r.Context())
	defer conn.Close()

	cmd := conn.Incr(r.Context(), key)
	if err := cmd.Err(); err != nil {
		log.Errorln("failed when incr query", err)
		writeInternalServerError(w)
		return
	}

	writeCounterResp(w, counterResp{
		CurrentValue: cmd.Val(),
	})
}

func countDecHandler(w http.ResponseWriter, r *http.Request) {
	conn := RedisConnContext(r.Context())
	defer conn.Close()

	cmd := conn.Decr(r.Context(), key)
	if err := cmd.Err(); err != nil {
		log.Errorln("failed when decr query", err)
		writeInternalServerError(w)
		return
	}

	writeCounterResp(w, counterResp{
		CurrentValue: cmd.Val(),
	})
}
