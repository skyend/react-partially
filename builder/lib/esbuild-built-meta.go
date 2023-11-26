package lib

type EsbuildBuiltMeta struct {
	Outputs map[string]struct {
		Imports []struct {
			Path string `json:"path"`
			Kind string `json:"kind"`
		} `json:"imports"`

		Exports []interface{} `json:"exports"`

		EntryPoint string `json:"entryPoint"`

		Inputs map[string]struct {
			BytesInOutput int `json:"bytesInOutput"`
		}

		Bytes int `json:"bytes"`
	} `json:"outputs"`
}
