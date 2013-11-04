package phpserialize

import (
	"fmt"
)

const TYPE_VALUE_SEPARATOR = ':'
const VALUES_SEPARATOR = ';'

type PhpObject struct {
	members   map[interface{}]interface{}
	className string
}

func NewPhpObject() *PhpObject {
	d := &PhpObject{
		members: make(map[interface{}]interface{}),
	}
	return d
}

func (obj *PhpObject) GetClassName() string {
	return obj.className
}

func (obj *PhpObject) SetClassName(cName string) {
	obj.className = cName
}

func (obj *PhpObject) GetMembers() map[interface{}]interface{} {
	return obj.members
}

func (obj *PhpObject) GetPrivateMemberValue(memberName string) (interface{}, bool) {
	key := fmt.Sprintf("\x00%s\x00%s", obj.className, memberName)
	v, ok := obj.members[key]
	return v, ok
}

func (obj *PhpObject) SetPrivateMemberValue(memberName string, value interface{}) {
	key := fmt.Sprintf("\x00%s\x00%s", obj.className, memberName)
	obj.members[key] = value
}

func (obj *PhpObject) GetProtectedMemberValue(memberName string) (interface{}, bool) {
	key := "\x00*\x00" + memberName
	v, ok := obj.members[key]
	return v, ok
}

func (obj *PhpObject) SetProtectedMemberValue(memberName string, value interface{}) {
	key := "\x00*\x00" + memberName
	obj.members[key] = value
}

func (obj *PhpObject) GetPublicMemberValue(memberName string) (interface{}, bool) {
	v, ok := obj.members[memberName]
	return v, ok
}

func (obj *PhpObject) SetPublicMemberValue(memberName string, value interface{}) {
	obj.members[memberName] = value
}
