package mrender_gltf

type Animation struct {
	Channels []AnimationKey `json:"channels"`
	Name     string         `json:"name"`
}

type AnimationKey struct {
	Sampler int             `json:"sampler"`
	Target  AnimationTarget `json:"target"`
}

type AnimationTarget struct {
	Node int    `json:"node"`
	Path string `json:"path"`
}
