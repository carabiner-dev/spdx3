package base

type PreNode struct {
	SPDXID string `json:"spdxID"`
	ID     string `json:"@id"`
	Type   string `json:"type"`
}

func (pn *PreNode) GetSPDXID() string {
	return pn.SPDXID
}

func (pn *PreNode) GetID() string {
	return pn.ID
}

func (pn *PreNode) GetType() string {
	return pn.Type
}
