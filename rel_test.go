package main_test

import (
	"testing"

	"github.com/1414C/sqac/gen"
)

// TestGetAreFromAndToKeysOpt
//
// Check that TestGetAreFromAndToKeysOpt determines the correct result with common cases.
func TestGetAreFromAndToKeysOpt(t *testing.T) {

	var rel gen.Relation
	var fromRow gen.Info
	var toRow gen.Info
	var fromInfo []gen.Info
	var toInfo []gen.Info

	// (1,1)
	fromRow.Name = "ID"
	fromRow.Value = "uint64"
	fromRow.Required = true
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Foozle"
	fromRow.Value = "string"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Fuddle"
	fromRow.Value = "float32"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// to
	toRow.Name = "ToID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Volume"
	toRow.Value = "int64"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Weight"
	toRow.Value = "float64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult := rel.GetAreFromAndToKeysOpt("ID", "ToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetAreFromAndToKeysOpt expected false (1,1) - got %t", bResult)
	}

	// (1,0)
	fromInfo = []gen.Info{}
	fromRow.Name = "ID"
	fromRow.Value = "uint64"
	fromRow.Required = true
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Foozle"
	fromRow.Value = "string"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Fuddle"
	fromRow.Value = "float32"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// to
	toInfo = []gen.Info{}
	toRow.Name = "ToID"
	toRow.Value = "uint64"
	toRow.Required = false
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Volume"
	toRow.Value = "int64"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Weight"
	toRow.Value = "float64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetAreFromAndToKeysOpt("ID", "ToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetAreFromAndToKeysOpt expected false (1,0) - got %t", bResult)
	}

	// (1,0)
	toInfo = []gen.Info{}
	toRow.Name = "ID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetAreFromAndToKeysOpt("FromID", "ID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetAreFromAndToKeysOpt expected false (0,1) - got %t", bResult)
	}

	// (1,0)
	toInfo = []gen.Info{}
	toRow.Name = "AltToID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect true response
	bResult = rel.GetAreFromAndToKeysOpt("FromID", "AltToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetAreFromAndToKeysOpt expected false (0,1) - got %t", bResult)
	}

	// (0,0)
	toInfo = []gen.Info{}
	toRow.Name = "AltToID"
	toRow.Value = "uint64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect true response
	bResult = rel.GetAreFromAndToKeysOpt("FromID", "AltToID", fromInfo, toInfo)
	if bResult != true {
		t.Errorf("rel.GetAreFromAndToKeysOpt expected false (0,0) - got %t", bResult)
	}
}

// TestGetAreFromAndToKeysReq
//
// Check that TestGetAreFromAndToKeysReq determines the correct result with common cases.
func TestGetAreFromAndToKeysReq(t *testing.T) {

	var rel gen.Relation
	var fromRow gen.Info
	var toRow gen.Info
	var fromInfo []gen.Info
	var toInfo []gen.Info

	// (1,1)
	fromRow.Name = "ID"
	fromRow.Value = "uint64"
	fromRow.Required = true
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Foozle"
	fromRow.Value = "string"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Fuddle"
	fromRow.Value = "float32"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// to
	toRow.Name = "ToID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Volume"
	toRow.Value = "int64"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Weight"
	toRow.Value = "float64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// call with the key names/slices - expect true response
	bResult := rel.GetAreFromAndToKeysReq("ID", "ToID", fromInfo, toInfo)
	if bResult != true {
		t.Errorf("rel.GetAreFromAndToKeysReq expected true (1,1) - got %t", bResult)
	}

	// (1,0)
	fromInfo = []gen.Info{}
	fromRow.Name = "ID"
	fromRow.Value = "uint64"
	fromRow.Required = true
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Foozle"
	fromRow.Value = "string"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Fuddle"
	fromRow.Value = "float32"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// to
	toInfo = []gen.Info{}
	toRow.Name = "ToID"
	toRow.Value = "uint64"
	toRow.Required = false
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Volume"
	toRow.Value = "int64"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Weight"
	toRow.Value = "float64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetAreFromAndToKeysOpt("ID", "ToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetAreFromAndToKeysReq expected false (1,0) - got %t", bResult)
	}

	// (1,0)
	toInfo = []gen.Info{}
	toRow.Name = "ID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetAreFromAndToKeysReq("FromID", "ID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetAreFromAndToKeysReq expected false (0,1) - got %t", bResult)
	}

	// (1,0)
	toInfo = []gen.Info{}
	toRow.Name = "AltToID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetAreFromAndToKeysReq("FromID", "AltToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetAreFromAndToKeysReq expected false (0,1) - got %t", bResult)
	}

	// (0,0)
	toInfo = []gen.Info{}
	toRow.Name = "AltToID"
	toRow.Value = "uint64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetAreFromAndToKeysReq("FromID", "AltToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetAreFromAndToKeysReq expected false (0,0) - got %t", bResult)
	}
}

// TestGetIsFromKeyOptAndToKeyReq
//
// Check that GetIsFromKeyOptAndToKeyReq determines the correct result with common cases.
func TestGetIsFromKeyOptAndToKeyReq(t *testing.T) {

	var rel gen.Relation
	var fromRow gen.Info
	var toRow gen.Info
	var fromInfo []gen.Info
	var toInfo []gen.Info

	// (1,1)
	fromRow.Name = "ID"
	fromRow.Value = "uint64"
	fromRow.Required = true
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Foozle"
	fromRow.Value = "string"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Fuddle"
	fromRow.Value = "float32"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// to
	toRow.Name = "ToID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Volume"
	toRow.Value = "int64"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Weight"
	toRow.Value = "float64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult := rel.GetIsFromKeyOptAndToKeyReq("ID", "ToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetIsFromKeyOptAndToKeyReq expected false (1,1) - got %t", bResult)
	}

	// (1,0)
	fromInfo = []gen.Info{}
	fromRow.Name = "ID"
	fromRow.Value = "uint64"
	fromRow.Required = true
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Foozle"
	fromRow.Value = "string"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Fuddle"
	fromRow.Value = "float32"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// to
	toInfo = []gen.Info{}
	toRow.Name = "ToID"
	toRow.Value = "uint64"
	toRow.Required = false
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Volume"
	toRow.Value = "int64"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Weight"
	toRow.Value = "float64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetIsFromKeyOptAndToKeyReq("ID", "ToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetIsFromKeyOptAndToKeyReq expected false (1,0) - got %t", bResult)
	}

	// (0,1)
	toInfo = []gen.Info{}
	toRow.Name = "ID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect true response
	bResult = rel.GetIsFromKeyOptAndToKeyReq("FromID", "ID", fromInfo, toInfo)
	if bResult != true {
		t.Errorf("rel.GetIsFromKeyOptAndToKeyReq expected false (0,1) - got %t", bResult)
	}

	// (0,1)
	toInfo = []gen.Info{}
	toRow.Name = "AltToID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect true response
	bResult = rel.GetIsFromKeyOptAndToKeyReq("FromID", "AltToID", fromInfo, toInfo)
	if bResult != true {
		t.Errorf("rel.GetIsFromKeyOptAndToKeyReq expected false (0,1) - got %t", bResult)
	}

	// (0,0)
	toInfo = []gen.Info{}
	toRow.Name = "AltToID"
	toRow.Value = "uint64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetIsFromKeyOptAndToKeyReq("FromID", "AltToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetIsFromKeyOptAndToKeyReq expected false (0,0) - got %t", bResult)
	}
}

// TestGetIsFromKeyReqAndToKeyOpt
//
// Check that GetIsFromKeyReqAndToKeyOpt determines the correct result with common cases.
func TestGetIsFromKeyReqAndToKeyOpt(t *testing.T) {

	var rel gen.Relation
	var fromRow gen.Info
	var toRow gen.Info
	var fromInfo []gen.Info
	var toInfo []gen.Info

	// (1,1)
	fromRow.Name = "ID"
	fromRow.Value = "uint64"
	fromRow.Required = true
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Foozle"
	fromRow.Value = "string"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Fuddle"
	fromRow.Value = "float32"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// to
	toRow.Name = "ToID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Volume"
	toRow.Value = "int64"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Weight"
	toRow.Value = "float64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult := rel.GetIsFromKeyReqAndToKeyOpt("ID", "ToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetIsFromKeyReqAndToKeyOpt expected false (1,1) - got %t", bResult)
	}

	// (1,0)
	fromInfo = []gen.Info{}
	fromRow.Name = "ID"
	fromRow.Value = "uint64"
	fromRow.Required = true
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Foozle"
	fromRow.Value = "string"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Fuddle"
	fromRow.Value = "float32"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// to
	toInfo = []gen.Info{}
	toRow.Name = "ToID"
	toRow.Value = "uint64"
	toRow.Required = false
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Volume"
	toRow.Value = "int64"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Weight"
	toRow.Value = "float64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// call with the key names/slices - expect true response
	bResult = rel.GetIsFromKeyReqAndToKeyOpt("ID", "ToID", fromInfo, toInfo)
	if bResult != true {
		t.Errorf("rel.GetIsFromKeyReqAndToKeyOpt expected true (1,0) - got %t", bResult)
	}

	// (0,1)
	toInfo = []gen.Info{}
	toRow.Name = "ID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetIsFromKeyReqAndToKeyOpt("FromID", "ID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetIsFromKeyReqAndToKeyOpt expected false (0,1) - got %t", bResult)
	}

	// (0,1)
	toInfo = []gen.Info{}
	toRow.Name = "AltToID"
	toRow.Value = "uint64"
	toRow.Required = true
	toRow.IsKey = true
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetIsFromKeyReqAndToKeyOpt("FromID", "AltToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetIsFromKeyReqAndToKeyOpt expected false (0,1) - got %t", bResult)
	}

	// (0,0)
	toInfo = []gen.Info{}
	toRow.Name = "AltToID"
	toRow.Value = "uint64"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Foozle"
	toRow.Value = "string"
	toRow.Required = false
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	toRow.Name = "Fuddle"
	toRow.Value = "float32"
	toRow.Required = true
	toRow.IsKey = false
	toInfo = append(toInfo, toRow)
	toRow = gen.Info{}

	// from
	fromInfo = []gen.Info{}
	fromRow.Name = "FromID"
	fromRow.Value = "uint64"
	fromRow.Required = false
	fromRow.IsKey = true
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Volume"
	fromRow.Value = "int64"
	fromRow.Required = true
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	fromRow.Name = "Weight"
	fromRow.Value = "float64"
	fromRow.Required = false
	fromRow.IsKey = false
	fromInfo = append(fromInfo, fromRow)
	fromRow = gen.Info{}

	// call with the key names/slices - expect false response
	bResult = rel.GetIsFromKeyReqAndToKeyOpt("FromID", "AltToID", fromInfo, toInfo)
	if bResult != false {
		t.Errorf("rel.GetIsFromKeyReqAndToKeyOpt expected false (0,0) - got %t", bResult)
	}
}
