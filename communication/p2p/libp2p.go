package p2p

import (
	"context"
	"encoding/json"
	comm "github.com/ChainSafe/chainbridge-core/communication"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type Libp2pCommunication struct {
	h                    host.Host
	protocolID           protocol.ID
	streamManager        *StreamManager
	logger               zerolog.Logger
	subscriptionManagers map[comm.ChainBridgeMessageType]*SessionSubscriptionManager
	subscriberLocker     *sync.Mutex
}

func NewCommunication(h host.Host, protocolID protocol.ID) Libp2pCommunication {
	logger := log.With().Str("Module", "communication").Str("Peer", h.ID().Pretty()).Logger()
	c := Libp2pCommunication{
		h:                    h,
		protocolID:           protocolID,
		streamManager:        NewStreamManager(),
		logger:               logger,
		subscriptionManagers: make(map[comm.ChainBridgeMessageType]*SessionSubscriptionManager),
		subscriberLocker:     &sync.Mutex{},
	}
	c.startProcessingStream()
	return c
}

/** Communication interface methods **/

// Broadcast sends
func (c *Libp2pCommunication) Broadcast(
	peers peer.IDSlice,
	msg []byte,
	msgType comm.ChainBridgeMessageType,
	sessionID string,
) {
	hostID := c.h.ID()
	wMsg := comm.WrappedMessage{
		MessageType: msgType,
		SessionID:   sessionID,
		Payload:     msg,
		From:        hostID,
	}
	marshaledMsg, err := json.Marshal(wMsg)
	if err != nil {
		c.logger.Error().Err(err).Str("SessionID", sessionID).Msg("unable to marshal message")
		return
	}
	c.logger.Debug().Str("MsgType", msgType.String()).Str("SessionID", sessionID).Msg(
		"broadcasting message",
	)
	for _, p := range peers {
		if hostID != p {
			stream, err := c.h.NewStream(context.TODO(), p, c.protocolID)
			if err != nil {
				c.logger.Error().Err(err).Str("MsgType", msgType.String()).Str("SessionID", sessionID).Msgf(
					"unable to open stream toward %s", p.Pretty(),
				)
				return
			}

			err = WriteStreamWithBuffer(marshaledMsg, stream)
			if err != nil {
				c.logger.Error().Str("To", string(p)).Err(err).Msg("unable to send message")
				return
			}
			c.logger.Trace().Str(
				"From", string(wMsg.From)).Str(
				"To", p.Pretty()).Str(
				"MsgType", msgType.String()).Str(
				"SessionID", sessionID).Msg(
				"message sent",
			)
			c.streamManager.AddStream(sessionID, stream)
		}
	}
}

// Subscribe
func (c *Libp2pCommunication) Subscribe(
	msgType comm.ChainBridgeMessageType,
	sessionID string,
	channel chan *comm.WrappedMessage,
) string {
	c.subscriberLocker.Lock()
	defer c.subscriberLocker.Unlock()

	subManager, ok := c.subscriptionManagers[msgType]
	if !ok {
		subManager = NewSessionSubscriptionManager()
		c.subscriptionManagers[msgType] = subManager
	}

	sID := subManager.Subscribe(sessionID, channel)
	c.logger.Info().Str("SessionID", sessionID).Msgf("subscribed to message type %s", msgType)
	return sID
}

// UnSubscribe
func (c *Libp2pCommunication) UnSubscribe(
	msgType comm.ChainBridgeMessageType,
	sessionID string,
	subID string,
) {
	c.subscriberLocker.Lock()
	defer c.subscriberLocker.Unlock()

	subManager, ok := c.subscriptionManagers[msgType]
	if !ok {
		c.logger.Debug().Msgf("cannot find the given channels %s", msgType.String())
		return
	}
	if subManager == nil {
		return
	}

	subManager.UnSubscribe(sessionID, subID)
}

// ReleaseStream
func (c *Libp2pCommunication) ReleaseStream(sessionID string) {
	c.streamManager.ReleaseStream(sessionID)
	c.logger.Info().Str("SessionID", sessionID).Msg("released stream")
}

/** Helper methods **/

func (c Libp2pCommunication) startProcessingStream() {
	c.h.SetStreamHandler(c.protocolID, func(s network.Stream) {
		msg, err := c.processMessageFromStream(s)
		if err != nil {
			c.logger.Error().Err(err).Str("StreamID", s.ID()).Msg("unable to process message")
			return
		}

		subscribers := c.getSubscribers(msg.MessageType, msg.SessionID)
		for _, sub := range subscribers {
			sub <- msg
		}
	})
}

func (c *Libp2pCommunication) getSubscribers(
	msgType comm.ChainBridgeMessageType, sessionID string,
) []chan *comm.WrappedMessage {
	c.subscriberLocker.Lock()
	defer c.subscriberLocker.Unlock()

	messageIDSubscriber, ok := c.subscriptionManagers[msgType]
	if !ok {
		c.logger.Debug().Msgf("fail to find subscription manager for message type %s", msgType)
		return nil
	}

	return messageIDSubscriber.GetSubscribers(sessionID)
}

func (c *Libp2pCommunication) processMessageFromStream(s network.Stream) (*comm.WrappedMessage, error) {
	msgBytes, err := ReadStreamWithBuffer(s)
	if err != nil {
		c.streamManager.AddStream("UNKNOWN", s)
		return nil, err
	}

	var wrappedMsg comm.WrappedMessage
	if err := json.Unmarshal(msgBytes, &wrappedMsg); nil != err {
		c.streamManager.AddStream("UNKNOWN", s)
		return nil, err
	}

	c.streamManager.AddStream(wrappedMsg.SessionID, s)

	c.logger.Trace().Str(
		"From", string(wrappedMsg.From)).Str(
		"MsgType", wrappedMsg.MessageType.String()).Str(
		"SessionID", wrappedMsg.SessionID).Msg(
		"processed message",
	)

	return &wrappedMsg, nil
}
