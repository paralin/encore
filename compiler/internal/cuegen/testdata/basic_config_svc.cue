// Code generated by encore. DO NOT EDIT.
//
// The contents of this file are generated from the structs used in
// conjunction with Encore's `config.Load[T]()` function. This file
// automatically be regenerated if the data types within the struct
// are changed.
//
// For more information about this file, see:
// https://encore.dev/docs/develop/config
package svc

Name:     string // The users name
Port:     uint16
ReadOnly: bool // true if we're in read only mode

// MagicNumber is complicated and requires
// a multi-line comment to explain it.
MagicNumber: int
ID:          string // An ID
PublicKey:   bytes