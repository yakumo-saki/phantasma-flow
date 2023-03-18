package util

import "github.com/ulule/deepcopier"

func DeepCopy(from interface{}, to interface{}) {
	// err := copier.CopyWithOption(&to, &from, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	deepcopier.Copy(from).To(to)
	// return nil
}
