package config

type HttpConfig struct {
	Address        string            `yaml:"address"`
	RemoteTrust    RemoteTrustConfig `yaml:"remotetrust"`
	StopWaitSecond int               `yaml:"stopwaitsecond"`
	Cors           CorsConfig        `yaml:"cors"`
}

func (h *HttpConfig) setDefault(global *GlobalConfig) {
	if h.Address == "" {
		h.Address = "localhost:2689"
	}

	if h.StopWaitSecond <= 0 {
		h.StopWaitSecond = 10
	}

	h.RemoteTrust.setDefault(global)
	h.Cors.setDefault()
}

func (h *HttpConfig) check(co *CorsOrigin) ConfigError {
	err := h.RemoteTrust.check()
	if err != nil && err.IsError() {
		return err
	}

	err = h.Cors.check(co)
	if err != nil && err.IsError() {
		return err
	}

	return nil
}
