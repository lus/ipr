package shared

// Machine represents a machine which reports its IP address
type Machine struct {
	Name    string `json:"name"`
	Token   string `json:"token"`
	Address string `json:"address"`
	Updated int64  `json:"updated"`
}

// MachineRepository represents a service interface to provide machine storage methods
type MachineRepository interface {
	All() ([]*Machine, error)
	Lookup(name string) (*Machine, error)
	Upsert(machine *Machine) error
	Delete(name string) error
}
