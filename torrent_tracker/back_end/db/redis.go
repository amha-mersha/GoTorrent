package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/amha-mersha/GoTorrent/domains"
	"github.com/redis/go-redis/v9"
)

type RedisInst struct {
	client *redis.Client
}

func InitRedis(address, password string, DB uint8) (*RedisInst, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       int(DB),
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Redis: %v", err)
	}
	return &RedisInst{client}, nil
}

func (redis *RedisInst) AddPeer(infoHash string, peer domains.Peer) error {
	existingPeers, err := redis.GetPeers(infoHash)
	if err != nil {
		return fmt.Errorf("failed to retrieve existing peers: %v", err)
	}

	for _, existingPeer := range existingPeers {
		if existingPeer.PeerID == peer.PeerID {
			return nil
		}
	}

	peerJSON, err := json.Marshal(peer)
	if err != nil {
		return fmt.Errorf("failed to marshal peer: %v", err)
	}

	err = redis.client.LPush(context.Background(), infoHash, peerJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to add peer to Redis: %v", err)
	}

	return nil
}

func (redis *RedisInst) GetPeers(infoHash string) ([]domains.Peer, error) {
	peerJSONs, err := redis.client.LRange(context.Background(), infoHash, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get peers from redis: %v", err)
	}

	var peers []domains.Peer
	for _, peerJSON := range peerJSONs {
		var peer domains.Peer
		err := json.Unmarshal([]byte(peerJSON), &peer)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal peer: %v", err)
		}
		peers = append(peers, peer)
	}
	return peers, nil
}

func (redis *RedisInst) RemovePeer(infoHash string, peerID string) error {
	peerJSONs, err := redis.client.LRange(context.Background(), infoHash, 0, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get peers from Redis: %v", err)
	}

	var updatedPeers []string
	for _, peerJSON := range peerJSONs {
		var peer domains.Peer
		err = json.Unmarshal([]byte(peerJSON), &peer)
		if err != nil {
			return fmt.Errorf("failed to deserialize peer: %v", err)
		}
		if peer.PeerID != peerID {
			updatedPeers = append(updatedPeers, peerJSON)
		}
	}

	err = redis.client.Del(context.Background(), infoHash).Err()
	if err != nil {
		return fmt.Errorf("failed to delete old peer list: %v", err)
	}
	err = redis.client.RPush(context.Background(), infoHash, updatedPeers).Err()
	if err != nil {
		return fmt.Errorf("failed to update peer list in Redis: %v", err)
	}
	return nil
}

func (redis *RedisInst) ClearPeers(infoHash string) error {
	err := redis.client.Del(context.Background(), infoHash).Err()
	if err != nil {
		return fmt.Errorf("failed to clear peers for infoHash %s: %v", infoHash, err)
	}
	return nil
}
