package cmd

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type userValue User

func newUserValue(id int, name string) *userValue {
	u := &User{Id: id, Name: name}

	return (*userValue)(u)
}

func (v *userValue) String() string {
	return fmt.Sprintf(`{Id: %d, Name: "%s"}`, v.Id, v.Name)
}

func (v *userValue) Set(val string) error {
	u := &User{}
	if err := yaml.Unmarshal([]byte(val), u); err != nil {
		return err
	}

	*v = (userValue)(*u)

	return nil
}
func (v *userValue) Type() string {
	return "User"
}
