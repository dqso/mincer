package nc_adapter

func (p Producer) SendPayloadToClient(clientId uint64, payloadData []byte) error {
	// TODO log
	return p.server.SendPayloadToClient(clientId, payloadData)
}
