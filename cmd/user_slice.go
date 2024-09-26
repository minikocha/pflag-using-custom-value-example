package cmd

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type userSliceValue struct {
	value   *[]User
	changed bool
}

func newUserSliceValue(users []User) *userSliceValue {
	usv := &userSliceValue{}
	usv.value = &users
	*usv.value = users

	return usv
}

func (v *userSliceValue) String() string {
	buf := ""
	for i, u := range *v.value {
		if i > 0 {
			buf = buf + ", "
		}

		buf = buf + fmt.Sprintf(`{Id: %d, Name: "%s"}`, u.Id, u.Name)
	}

	return "[" + buf + "]"
}

func (v *userSliceValue) Set(val string) error {
	in := []User{}

	if strings.HasPrefix(val, "[") {
		if err := yaml.Unmarshal([]byte(val), &in); err != nil {
			return err
		}
	} else {
		u := User{}
		if err := yaml.Unmarshal([]byte(val), &u); err != nil {
			return err
		}

		in = append(in, u)
	}

	if !v.changed {
		*v.value = in
		v.changed = true

		return nil
	}

	// If the same Id exists in both current and new, overwrite with new one.
	*v.value = Merge(in, *v.value)
	v.changed = true

	return nil
}

func (v *userSliceValue) Type() string {
	return "UserSlice"
}

func Merge(src []User, dest []User) []User {
	for i, d := range dest {
		for j, s := range src {
			if d.Id == s.Id {
				dest[i] = s
				src = append(src[:j], src[j+1:]...)
				continue
			}
		}
	}

	if len(src) > 0 {
		dest = append(dest, src...)
	}

	return dest
}
