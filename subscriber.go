package vindalu

import (
	log "github.com/euforia/simplelog"
	"github.com/nats-io/nats"

	"github.com/vindalu/vindalu/core"
)

type VindaluSubscriber struct {
	conn   *nats.Conn
	enConn *nats.EncodedConn
}

func NewVindaluSubscriber(servers []string, logger *log.Logger) (vs *VindaluSubscriber, err error) {
	vs = &VindaluSubscriber{}

	opts := nats.DefaultOptions
	opts.Servers = servers

	if vs.conn, err = opts.Connect(); err != nil {
		return
	}
	logger.Debug.Printf("nats client connected to: %v!\n", vs.conn.ConnectedUrl())

	vs.conn.Opts.ReconnectedCB = func(nc *nats.Conn) {
		logger.Debug.Printf("nats client reconnected to: %v!\n", nc.ConnectedUrl())
	}

	vs.conn.Opts.DisconnectedCB = func(_ *nats.Conn) {
		logger.Debug.Printf("nats client disconnected!\n")
	}

	vs.enConn, err = nats.NewEncodedConn(vs.conn, nats.JSON_ENCODER)

	return
}

func (vs *VindaluSubscriber) Subscribe(topic string) (ch chan *core.Event, err error) {
	// Goes no where as we do not want to allow writing (i.e publishing)
	//if err = vs.enConn.BindSendChan(topic, make(chan *core.Event)); err != nil {
	//	return
	//}

	ch = make(chan *core.Event)
	_, err = vs.enConn.BindRecvChan(topic, ch)
	return
}

func (vs *VindaluSubscriber) SubscribeQueueGroup(topic, qGroup string) (ch chan *core.Event, err error) {
	// Goes no where as we do not want to allow writing (i.e publishing)
	//if err = vs.enConn.BindSendChan(topic, make(chan *core.Event)); err != nil {
	//	return
	//}

	ch = make(chan *core.Event)
	_, err = vs.enConn.BindRecvQueueChan(topic, qGroup, ch)
	return
}

func (vs *VindaluSubscriber) Close() {
	vs.enConn.Close()
	vs.conn.Close()
}
