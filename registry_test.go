package main

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/orbs-network/orbs-contract-sdk/go/testing/unit"
	"github.com/stretchr/testify/require"
)

var (
	key  = "abc"
	meta = `{
		"contentRegistryItem": {
		  "id": "9d5cdf280cdc400bacf0eeca7bbb6e51",
		  "type": "picture",
		  "title": "Britain CWC Cricket",
		  "author": {
			"name": "Aijaz Rahi",
			"title": "STAFF",
			"id": "STF"
		  },
		  "url": "https://mapi.associatedpress.com/v1/items/9d5cdf280cdc400bacf0eeca7bbb6e51/preview/preview.jpg",
		  "credit": "ASSOCIATED PRESS",
		  "copyright": "Copyright 2019 The Associated Press. All rights reserved",
		  "createdDateTime": "2019-07-11T16:00:43",
		  "publishedDateTime": "2019-07-11T16:00:43",
		  "description": "England’s captain Eoin Morgan bats during the Cricket World Cup semi-final match between England and Australia at Edgbaston in Birmingham, England, Thursday, July 11, 2019. (AP Photo/Aijaz Rahi)",
		  "rightsModel": {
			"id": "ap42",
			"name": "editorialOnly",
			"restrictions": "no online or web use"
		  }
		}
	  }`
)

func TestVerify(t *testing.T) {
	InServiceScope(nil, nil, func(m Mockery) {
		require.Panics(t, func() {
			verify(key)
		})
	})
}

func TestRegisterAndSearch(t *testing.T) {
	InServiceScope(nil, nil, func(m Mockery) {
		s := register(key, meta)
		if s == "" {
			t.Error("Failed to create hash array")
		} else {
			fmt.Println(s)
		}

		s = register("aaa", strings.Replace(meta, "9d5cdf280cdc400bacf0eeca7bbb6e51", "9d5cdf280cdc400bacf0eeca7bbb6xyz", 1))
		if s == "" {
			t.Error("Failed to create hash array")
		} else {
			fmt.Println(s)
		}

		s = search(key, 25)
		fmt.Println(s)
	})
}
