package dota2

import (
	"github.com/Philipp15b/go-steam"
	"github.com/Philipp15b/go-steam/protocol"
	"github.com/Philipp15b/go-steam/protocol/gamecoordinator"
	protobuf2 "github.com/Philipp15b/go-steam/protocol/protobuf"
	"github.com/Philipp15b/go-steam/protocol/steamlang"
	"github.com/Philipp15b/go-steam/tf2/protocol/protobuf"
	dota2 "github.com/bittermandel/go-dota2/objects"
	"io/ioutil"
	"log"
)

var AppId uint32 = 570

type GCReadyEvent struct{}

type Dota2 struct {
	client *steam.Client
	logOnDetails *steam.LogOnDetails
	loggedOn bool
	richPresences	[]*dota2.RichPresence
}

func New(client *steam.Client, logOnDetails *steam.LogOnDetails) (dota *Dota2) {
	dota = &Dota2{client: client, logOnDetails: logOnDetails, loggedOn: false}
	client.GC.RegisterPacketHandler(dota)
	client.RegisterPacketHandler(dota)
	go func() {
		for event := range client.Events() {
			dota.handleEvent(event)
		}
	}()

	log.Println("Connecting to Steam")
	steam.CMServers[1] = []string{"162.254.198.132:27018", "162.254.198.131:27017","162.254.198.131:27019","162.254.198.132:27017","162.254.198.130:27019","162.254.198.133:27017","162.254.198.130:27018","162.254.198.132:27019","162.254.198.131:27018","162.254.198.133:27019","162.254.198.133:27018","162.254.198.130:27017","162.254.196.68:27018","162.254.196.84:27019","162.254.197.42:27018","162.254.197.42:27019","162.254.197.180:27018","162.254.197.40:27018","162.254.196.67:27017","155.133.248.52:27017","162.254.196.84:27017","162.254.197.180:27017","162.254.197.180:27019","155.133.248.53:27018","155.133.248.52:27019","155.133.248.53:27017","162.254.196.67:27019","162.254.197.181:27017","162.254.197.40:27019","155.133.248.50:27019","162.254.196.67:27018","162.254.197.42:27017","162.254.197.181:27019","162.254.197.40:27017","155.133.248.51:27018","155.133.248.53:27019","155.133.248.51:27017","155.133.248.51:27019","155.133.248.52:27018","155.133.248.50:27018","162.254.196.68:27019","162.254.196.83:27017","162.254.197.181:27018","155.133.248.50:27017","162.254.196.83:27018","162.254.196.68:27017","162.254.196.83:27019","162.254.196.84:27018","185.25.182.77:27019","185.25.182.76:27017","185.25.182.76:27018","185.25.182.77:27017","185.25.182.77:27018","185.25.182.76:27019","162.254.192.109:27017","162.254.192.109:27019","162.254.192.108:27018","162.254.192.108:27017","162.254.192.101:27018","162.254.192.100:27019","162.254.192.101:27017","162.254.192.109:27018","162.254.192.101:27019","162.254.192.100:27018","162.254.192.108:27019","162.254.192.100:27017","155.133.246.69:27017","155.133.246.69:27019","155.133.246.69:27018","155.133.246.68:27019","155.133.246.68:27018","155.133.246.68:27017","146.66.155.101:27018","146.66.155.101:27019","146.66.155.100:27018","146.66.155.101:27017","146.66.155.100:27017","146.66.155.100:27019","162.254.193.7:27017","162.254.193.47:27017","162.254.193.6:27018","162.254.193.47:27019","162.254.193.6:27019","162.254.193.46:27019","162.254.193.46:27018","162.254.193.47:27018","162.254.193.7:27018","162.254.193.6:27017","162.254.193.7:27019","162.254.193.46:27017","155.133.230.34:27018","155.133.230.34:27019","155.133.230.34:27017","155.133.230.50:27018","155.133.230.50:27017","155.133.230.50:27019","162.254.195.83:27017","162.254.195.66:27017","162.254.195.82:27017","162.254.195.82:27019"}
	client.ConnectEurope()
	log.Println("Connected to Steam")

	return dota
}

func (dota *Dota2) HandlePacket(packet *protocol.Packet){
	//log.Printf("%s\n", packet.EMsg)
	switch packet.EMsg {
	case steamlang.EMsg_ClientRichPresenceInfo:
		body := new(protobuf2.CMsgClientRichPresenceInfo)
		packet.ReadProtoMsg(body)
		rps := []*dota2.RichPresence{}
		for _, rp := range body.RichPresence {
			rps = append(rps, &dota2.RichPresence{SteamID: *rp.SteamidUser, RichPresenceKV: string(rp.RichPresenceKv)})
		}
		dota.richPresences = rps
		log.Printf("Rich Presences: %+v", dota.richPresences)
	}
}

func (dota *Dota2) HandleGCPacket(packet *gamecoordinator.GCPacket) {
	if packet.AppId != AppId {
		return
	}
	switch protobuf.EGCBaseClientMsg(packet.MsgType) {
	case protobuf.EGCBaseClientMsg_k_EMsgGCClientWelcome:
		dota.handleWelcome(packet)
	default:
		log.Printf("Packet %+v\n", packet)
	}
}

func (dota *Dota2) handleEvent(event interface{}) {
	//log.Printf("%+v\n", reflect.ValueOf(event).Type())
	switch e := event.(type) {
	case *steam.ConnectedEvent:
		dota.client.Auth.LogOn(dota.logOnDetails)
	case *steam.LogOnFailedEvent:
		log.Printf("LogOn failed. Reason: %v", e.Result)
	case *steam.LoggedOnEvent:
		log.Printf("Logged on (%v) with SteamId %v and account flags %v", e.Result, e.ClientSteamId, e.AccountFlags)
		dota.loggedOn = true
		// Ask for RichPresence as soon as we're logged in.
		packet := protocol.ClientMsgProtobuf{Header: &steamlang.MsgHdrProtoBuf{
			Msg:   steamlang.EMsg_ClientRichPresenceRequest,
			Proto: &protobuf2.CMsgProtoBufHeader{RoutingAppid: &AppId},
		}, Body:&protobuf2.CMsgClientRichPresenceRequest{SteamidRequest:[]uint64{76561197993261133, 76561198074261550}}}
		dota.client.Write(&packet)
	case *steam.DisconnectedEvent:
		log.Printf("Disconnected.")
	case *steam.MachineAuthUpdateEvent:
		// Save MachineAuth hash for future use (skip AuthCode refresh)
		err := ioutil.WriteFile("./.sentry", e.Hash, 0666)
		if err != nil {
			panic(err)
		}
	}
}

func (dota *Dota2) handleWelcome(packet *gamecoordinator.GCPacket) {
	dota.client.Emit(&GCReadyEvent{})
}