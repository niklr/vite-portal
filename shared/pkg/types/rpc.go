package types

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RpcAppInfoResponse struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
}

type RpcViteNodeInfoResponse struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	NetID                 int     `json:"netId"`
	Version               int     `json:"version"`
	Address               string  `json:"address"`
	PeerCount             int     `json:"peerCount"`
	Height                uint64  `json:"height"`
	Nodes                 int     `json:"nodes"`
	Latency               []int64 `json:"latency"` // [0,1,12,24]
	BroadCheckFailedRatio float32 `json:"broadCheckFailedRatio"`
}

type RpcViteProcessInfoResponse struct {
	BuildVersion  string `json:"build_version"`
	CommitVersion string `json:"commit_version"`
	NodeName      string `json:"nodeName"`
	RewardAddress string `json:"rewardAddress"`
	Pid           int    `json:"pid"`
}
