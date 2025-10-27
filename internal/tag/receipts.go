package tag

type ReceiptType string

const (
	RSpawnChain ReceiptType = "spawn_chain"
	RInject     ReceiptType = "inject"
	RDiffuse    ReceiptType = "diffuse"
	RMetaBirth  ReceiptType = "meta_totelevation"
	RBackfeed   ReceiptType = "backfeed"
	RReconcile  ReceiptType = "reconcile"
	RQuench     ReceiptType = "quench"
)

type Receipt struct {
	Step    int         `json:"step"`
	Type    ReceiptType `json:"type"`
	Subject string      `json:"subject"`
	Note    string      `json:"note"`
	Value1  float64     `json:"value1,omitempty"`
	Value2  float64     `json:"value2,omitempty"`
}
