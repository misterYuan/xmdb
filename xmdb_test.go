package xmdb

import (
	"log"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func cp1(src []bson.M) []bson.M {
	return append(src, bson.M{"1": 1})
}

func cp2(src []bson.M) []bson.M {
	return append(src, bson.M{"a": "b"})
}

type A struct {
	Name string
}

func (a *A) cp3(src []bson.M) []bson.M {
	return append(src, bson.M{a.Name: "haha"})
}

func TestOption(t *testing.T) {
	a := &A{"yaya"}
	res := FillPipe(cp1, cp2, a.cp3)
	log.Println(res)
}
