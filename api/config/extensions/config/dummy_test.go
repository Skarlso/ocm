package config_test

import (
	"ocm.software/ocm/api/config"
	"ocm.software/ocm/api/config/cpi"
	"ocm.software/ocm/api/utils/runtime"
)

const (
	DummyType   = "Dummy"
	DummyTypeV1 = DummyType + "/v1"
)

func RegisterAt(reg cpi.ConfigTypeScheme) {
	reg.Register(cpi.NewConfigType[*Config](DummyType))
	reg.Register(cpi.NewConfigType[*Config](DummyTypeV1))
}

// Config describes a a dummy config
type Config struct {
	runtime.ObjectVersionedType `json:",inline"`
	Alice                       string `json:"alice,omitempty"`
	Bob                         string `json:"bob,omitempty"`
}

// NewConfig creates a new memory ConfigSpec
func NewConfig(a, b string) *Config {
	return &Config{
		ObjectVersionedType: runtime.NewVersionedTypedObject(DummyType),
		Alice:               a,
		Bob:                 b,
	}
}

func (a *Config) GetType() string {
	return DummyType
}

func (a *Config) Info() string {
	return "dummy config"
}

func (a *Config) ApplyTo(ctx config.Context, target interface{}) error {
	d, ok := target.(*dummyContext)
	if ok {
		d.applied = append(d.applied, a)
		return nil
	}
	c, ok := target.(*dummyTarget)
	if ok {
		c.used = ctx
		return nil
	}
	return cpi.ErrNoContext(DummyType)
}

////////////////////////////////////////////////////////////////////////////////

type dummyTarget struct {
	used config.Context
}

////////////////////////////////////////////////////////////////////////////////

func newDummy(ctx config.Context) *dummyContext {
	d := &dummyContext{
		config: ctx,
	}
	d.update()
	return d
}

type dummyContext struct {
	config         config.Context
	lastGeneration int64
	applied        []*Config
}

func (d *dummyContext) getApplied() []*Config {
	d.update()
	return d.applied
}

func (d *dummyContext) update() error {
	gen, err := d.config.ApplyTo(d.lastGeneration, d)
	d.lastGeneration = gen
	return err
}
