package main

import "strings"

// Key formats keys (for use in a Key/Value store).
type Key struct {
	Base      string
	Model     string
	Attribute string
	Value     string
}

// NewKey makes a new KeyMaker.
func NewKey(model, attribute, value string) *Key {
	k := new(Key)

	k.Base = "gogocal.hcpss.org"
	k.Model = model
	k.Attribute = attribute
	k.Value = value

	return k
}

// NewKeyFromString creates a new Key given a serialized key string.
func NewKeyFromString(sKey string) *Key {
	k := new(Key)

	parts := strings.Split(sKey, ":")
	k.Base = parts[0]
	k.Model = parts[1]
	k.Attribute = parts[2]
	k.Value = parts[3]

	return k
}

// KeysEqual returns true if the keys are equal and false otherwise.
func KeysEqual(key1, key2 *Key) bool {
	return key1.String() == key2.String()
}

// String is the string representation of a KeyMaker which should be in the
// format "base:model:attribute:value".
func (k Key) String() string {
	parts := []string{k.Base, k.Model, k.Attribute, k.Value}

	return strings.Join(parts, ":")
}

// KeyList represents a list (slice) of keys.
type KeyList []*Key

// KeyListFromSlice creates a KeyList from  slice of string keys.
func KeyListFromSlice(strings []string) KeyList {
	var kl KeyList

	for _, s := range strings {
		key := NewKeyFromString(s)
		kl = append(kl, key)
	}

	return kl
}

// KeyListDiff returns the Keys in keys1 that are not in keys2 as a KeyList.
func KeyListDiff(keys1, keys2 KeyList) KeyList {
	var diff KeyList

	for i := range keys1 {
		found := false

		for j := range keys2 {
			if KeysEqual(keys1[i], keys2[j]) {
				found = true
			}
		}

		if !found {
			diff = append(diff, keys1[i])
		}
	}

	return diff
}

// ToSlice converts a KeyList to a slice of strings.
func (kl KeyList) ToSlice() []string {
	var s []string

	for _, key := range kl {
		s = append(s, key.String())
	}

	return s
}

// Merge adds values in keys to the KeyList without duplicating any.
func (kl *KeyList) Merge(keys KeyList) {
	for _, k := range keys {
		found := false

		for _, l := range *kl {
			if KeysEqual(k, l) {
				found = true
			}
		}

		if !found {
			*kl = append(*kl, k)
		}
	}
}
