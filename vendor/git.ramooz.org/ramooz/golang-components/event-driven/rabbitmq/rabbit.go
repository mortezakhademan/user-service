package rabbitmq

import (
	"fmt"
	error2 "git.ramooz.org/ramooz/golang-components/error-handler"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

const (
	DIRECT  Kind = iota // DIRECT a event goes to the consumerQueues whose binding key exactly matches the routing key of the event.
	FANOUT              // FANOUT exchanges can be useful when the same event needs to be sent to one or more consumerQueues with consumers who may process the same event in different ways.
	TOPIC               // TOPIC exchange is similar to direct exchange, but the routing is done according to the routing pattern. Instead of using fixed routing key, it uses wildcards.
	HEADERS             // HEADERS exchange routes events based on arguments containing headers and optional values. It uses the event header attributes for routing.
)

const (
	delayReconnectTime = 5 * time.Second
)

// NewConnection create a rabbitmq connection object
func NewConnection(serviceName string, options *Options, done chan os.Signal) (*Connection, error) {
	opts, err := validateOptions(serviceName, options)
	if err != nil {
		return nil, err
	}
	connObj := &Connection{
		ServiceCallerName: serviceName,
		ConnOpt:           opts,
		done:              done,
		alive:             true,
		consumerQueues:    make(map[string]EventHandler),
	}
	go connObj.handleReconnect(opts.UriAddress)
	for {
		if connObj.isConnected {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return connObj, nil
}

// NewEncodedConn will wrap an existing Connection and utilize the appropriate registered encoder
func NewEncodedConn(c *Connection, encType string) (*EncodedConn, error) {
	if c == nil {
		return nil, error2.New(ERROR_NIL_CONNECTION, nil)
	}
	ec := &EncodedConn{Conn: c, Enc: EncoderForType(encType)}
	if ec.Enc == nil {
		return nil, fmt.Errorf("no encoder registered for '%s'", encType)
	}
	return ec, nil
}

// connect dial to rabbitMQ server and declare exchange
func (c *Connection) connect() bool {
	conn, err := amqp.Dial(c.ConnOpt.UriAddress)
	if err != nil {
		log.Printf("rabbitMQ on dial got error %v", err)
		return false
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("rabbitMQ on connect to channel got error %v", err)
		return false
	}
	c.updateConnection(conn, ch)
	c.isConnected = true
	return true
}

// updateConnection update connection and channel in memory
func (c *Connection) updateConnection(connection *amqp.Connection, channel *amqp.Channel) {
	c.conn = connection
	c.channel = channel
	c.notifyClose = make(chan *amqp.Error)
	c.channel.NotifyClose(c.notifyClose)
}

// handleReconnect if closing rabbitMQ try to connect rabbitMQ continuously
func (c *Connection) handleReconnect(addr string) {
	for c.alive {
		c.isConnected = false
		now := time.Now()
		log.Printf("attempting to connect to rabbitMQ %v", addr)
		retryCount := 0
		for !c.connect() {
			if !c.alive {
				return
			}
			select {
			case <-c.done:
				return
			case <-time.After(delayReconnectTime + time.Duration(retryCount)*time.Second):
				log.Printf("cannot connect to rabbitMQ try connecting to rabbitMQ (next try after %v)...", delayReconnectTime+time.Duration(retryCount)*time.Second)
				if retryCount != 10 {
					retryCount++
				}
			}
		}
		log.Printf("connected to rabbitMQ after %v second", time.Since(now).Seconds())
		select {
		case <-c.done:
			return
		case <-c.notifyClose:
		}
	}
}

// ExchangeDeclare declare new exchange with specific kind (direct, topic, fanout, headers)
func (c *Connection) ExchangeDeclare(exchange string, kind Kind) error {
	if checkElementInSlice(c.exchanges, exchange) {
		return error2.New(ERROR_EXCHANGE_ALREADY_EXISTS, nil)
	}
	c.exchanges = append(c.exchanges, exchange)
	if err := c.channel.ExchangeDeclare(
		exchange,
		kind.String(),
		c.ConnOpt.DurableExchange,
		c.ConnOpt.AutoDelete,
		false,
		c.ConnOpt.NoWait,
		nil); err != nil {
		return err
	}
	return nil
}

// DeclarePublisherQueue declare new queue and bind queue and bind exchange with routing key
func (c *Connection) DeclarePublisherQueue(queue, exchange string, routingKey ...string) error {
	return c.queueDeclare(queue, exchange, routingKey...)
}

// DeclareConsumerQueue declare new queue and bind queue and bind exchange with routing key
func (c *Connection) DeclareConsumerQueue(eventHandler EventHandler, queue, exchange string, routingKey ...string) error {
	if _, ok := c.consumerQueues[queue]; ok {
		return error2.New(ERROR_QUEUE_ALREADY_EXISTS, nil)
	} else {
		c.consumerQueues[queue] = eventHandler
	}
	return c.queueDeclare(queue, exchange, routingKey...)
}

func (c *Connection) queueDeclare(queue, exchange string, routingKey ...string) error {
	if _, err := c.channel.QueueDeclare(
		queue,
		c.ConnOpt.DurableExchange,
		c.ConnOpt.AutoDelete,
		c.ConnOpt.ExclusiveQueue,
		c.ConnOpt.NoWait,
		nil,
	); err != nil {
		return err
	}

	for _, key := range routingKey {
		if err := c.channel.QueueBind(
			queue,
			key,
			exchange,
			c.ConnOpt.NoWait,
			nil,
		); err != nil {
			return err
		}
	}
	return nil
}

// IsConnected check rabbitMQ client is connected
func (c *Connection) IsConnected() bool {
	return c.isConnected
}

// GetEventConnection return connection and channel object of amqp
func (c *Connection) GetEventConnection() (*amqp.Connection, *amqp.Channel) {
	return c.conn, c.channel
}

// GetExchangeList return list of exchanges
func (c *Connection) GetExchangeList() []string {
	return c.exchanges
}

// GetQueueList return list of consumerQueues with handlers
func (c *Connection) GetQueueList() map[string]EventHandler {
	return c.consumerQueues
}

// Ack delegates an acknowledgement through the Acknowledger interface that the client or server has finished work on a delivery.
func (d Delivery) Ack(multiple bool) error {
	return d.Acknowledger.Ack(d.DeliveryTag, multiple)
}

// Nack negatively acknowledge the delivery of message(s) identified by the delivery tag from either the client or server.
func (d Delivery) Nack(multiple bool, requeue bool) error {
	return d.Acknowledger.Nack(d.DeliveryTag, multiple, requeue)
}

// Reject delegates a negatively acknowledgement through the Acknowledger interface.
func (d Delivery) Reject(requeue bool) error {
	return d.Acknowledger.Reject(d.DeliveryTag, requeue)
}

// Close stop rabbitMQ client
func (c *Connection) Close() error {
	if !c.isConnected {
		return nil
	}
	c.alive = false
	if err := c.channel.Close(); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	c.isConnected = false
	log.Printf("gracefully stopped rabbitMQ connection")
	return nil
}

// String exchange type as string
func (k Kind) String() string {
	switch k {
	case DIRECT:
		return "direct"
	case FANOUT:
		return "fanout"
	case TOPIC:
		return "topic"
	case HEADERS:
		return "headers"
	default:
		return "topic"
	}
}

func validateOptions(serviceName string, newOpt *Options) (*Options, error) {
	opt := getDefaultOptions()
	if len(serviceName) == 0 {
		return nil, error2.New(ERROR_SERVICE_NAME, nil)
	}
	if len(newOpt.UriAddress) == 0 {
		return nil, error2.New(ERROR_URI_ADDRESS, nil)
	} else {
		opt.UriAddress = newOpt.UriAddress
	}
	if newOpt.DurableExchange != opt.DurableExchange {
		opt.DurableExchange = newOpt.DurableExchange
	}
	if newOpt.AutoAck != opt.AutoAck {
		opt.AutoAck = newOpt.AutoAck
	}
	if newOpt.AutoDelete != opt.AutoDelete {
		opt.AutoDelete = newOpt.AutoDelete
	}
	if newOpt.NoWait != opt.NoWait {
		opt.NoWait = newOpt.NoWait
	}
	if newOpt.ExclusiveQueue != opt.ExclusiveQueue {
		opt.ExclusiveQueue = newOpt.ExclusiveQueue
	}
	return opt, nil
}

func checkElementInSlice(slice []string, newElement string) bool {
	if slice != nil {
		for _, s := range slice {
			if s == newElement {
				return true
			}
		}
	}
	return false
}
