package oss

import "net/url"

type Objects struct {
	List      []Object
	NextToken string
}

type Object struct {
	path string
}

func NewObject(path string) Object {
	return Object{path}
}

func (obj Object) ToUrl(bucket *Bucket) url.URL {
	url := bucket.ToUrl()
	url.Path = obj.path
	return url
}
