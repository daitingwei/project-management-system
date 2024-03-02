package kk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	brokers        []string
	topic         string
	groupID       string
	logstashHost  string
	logstashPort  int
	maxRetries    int
	retryInterval time.Duration
	conn          net.Conn
	mu            sync.Mutex
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewKafkaConsumer(brokers []string, topic, groupID, logstashHost string, logstashPort int) *KafkaConsumer {
	return &KafkaConsumer{
		brokers:        brokers,
		topic:         topic,
		groupID:       groupID,
		logstashHost:  logstashHost,
		logstashPort:  logstashPort,
		maxRetries:    3,
		retryInterval: 1 * time.Second,
	}
}

func (c *KafkaConsumer) connectToLogstash() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	addr := fmt.Sprintf("%s:%d", c.logstashHost, c.logstashPort)
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to logstash at %s: %w", addr, err)
	}
	c.conn = conn
	return nil
}

func (c *KafkaConsumer) sendToLogstash(data []byte) error {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()

	if conn == nil {
		return fmt.Errorf("not connected to logstash")
	}

	data = append(data, '\n')
	_, err := conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send data to logstash: %w", err)
	}
	return nil
}

func (c *KafkaConsumer) sendToLogstashWithRetry(data []byte) error {
	var lastErr error
	for i := 0; i < c.maxRetries; i++ {
		if err := c.sendToLogstash(data); err != nil {
			lastErr = err
			log.Printf("Attempt %d: failed to send to logstash: %v", i+1, err)
			time.Sleep(c.retryInterval)

			if reconnErr := c.connectToLogstash(); reconnErr != nil {
				log.Printf("Failed to reconnect to logstash: %v", reconnErr)
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to send to logstash after %d retries: %w", c.maxRetries, lastErr)
}

func (c *KafkaConsumer) closeLogstashConnection() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	if err := c.connectToLogstash(); err != nil {
		log.Printf("Initial connection to logstash failed: %v, will retry later", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.brokers,
		Topic:          c.topic,
		GroupID:        c.groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		MaxWait:        1 * time.Second,
		StartOffset:    kafka.FirstOffset,
		CommitInterval: time.Second,
	})

	ctx, c.cancel = context.WithCancel(ctx)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer reader.Close()
		defer c.closeLogstashConnection()

		log.Printf("Kafka consumer started, reading from topic: %s, group: %s", c.topic, c.groupID)

		for {
			select {
			case <-ctx.Done():
				log.Printf("Kafka consumer stopping due to context cancellation")
				return
			default:
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						log.Printf("Kafka consumer context cancelled, exiting")
						return
					}
					log.Printf("Error reading message: %v", err)
					time.Sleep(3 * time.Second)
					continue
				}

				log.Printf("Received message: topic=%s, partition=%d, offset=%d, size=%d",
					msg.Topic, msg.Partition, msg.Offset, len(msg.Value))

				if err := c.sendToLogstashWithRetry(msg.Value); err != nil {
					log.Printf("Failed to send message to logstash after retries: %v", err)
					continue
				}

				log.Printf("Successfully sent message to logstash, offset: %d", msg.Offset)
			}
		}
	}()

	return nil
}

func (c *KafkaConsumer) Stop() error {
	log.Printf("Stopping Kafka consumer...")

	if c.cancel != nil {
		c.cancel()
	}

	c.wg.Wait()

	c.closeLogstashConnection()

	log.Printf("Kafka consumer stopped")
	return nil
}

type KafkaLogMessage struct {
	Type     string         `json:"type"`
	Action   string         `json:"action"`
	Time     string         `json:"time"`
	Msg      string         `json:"msg"`
	Field    map[string]any `json:"field"`
	FuncName string         `json:"funcName"`
}

func (c *KafkaConsumer) SendTestLog(funcName, msg string, fields map[string]any) error {
	logMsg := KafkaLogMessage{
		Type:     "info",
		Action:   "test",
		Time:     time.Now().Format("2006-01-02 15:04:05"),
		Msg:      msg,
		Field:    fields,
		FuncName: funcName,
	}

	data, err := json.Marshal(logMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal log message: %w", err)
	}

	return c.sendToLogstashWithRetry(data)
}
