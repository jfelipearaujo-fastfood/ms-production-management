package dummy

type Dummy struct {
	Name string
}

func NewDummy(name string) *Dummy {
	return &Dummy{
		Name: name,
	}
}

func (d *Dummy) GetName() string {
	return d.Name
}
