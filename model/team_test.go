// Copyright (c) 2015 Spinpunch, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"strings"
	"testing"
)

func TestTeamJson(t *testing.T) {
	o := Team{Id: NewId(), DisplayName: NewId()}
	json := o.ToJson()
	ro := TeamFromJson(strings.NewReader(json))

	if o.Id != ro.Id {
		t.Fatal("Ids do not match")
	}
}

func TestTeamIsValid(t *testing.T) {
	o := Team{}

	if err := o.IsValid(); err == nil {
		t.Fatal("should be invalid")
	}

	o.Id = NewId()
	if err := o.IsValid(); err == nil {
		t.Fatal("should be invalid")
	}

	o.CreateAt = GetMillis()
	if err := o.IsValid(); err == nil {
		t.Fatal("should be invalid")
	}

	o.UpdateAt = GetMillis()
	if err := o.IsValid(); err == nil {
		t.Fatal("should be invalid")
	}

	o.Email = strings.Repeat("01234567890", 20)
	if err := o.IsValid(); err == nil {
		t.Fatal("should be invalid")
	}

	o.Email = "corey@hulen.com"
	o.DisplayName = strings.Repeat("01234567890", 20)
	if err := o.IsValid(); err == nil {
		t.Fatal("should be invalid")
	}

	o.DisplayName = "1234"
	o.Name = "ZZZZZZZ"
	if err := o.IsValid(); err == nil {
		t.Fatal("should be invalid")
	}

	o.Name = "zzzzz"
	o.Type = TEAM_OPEN
	if err := o.IsValid(); err != nil {
		t.Fatal(err)
	}
}

func TestTeamPreSave(t *testing.T) {
	o := Team{DisplayName: "test"}
	o.PreSave()
	o.Etag()
}

func TestTeamPreUpdate(t *testing.T) {
	o := Team{DisplayName: "test"}
	o.PreUpdate()
}
