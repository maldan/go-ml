package test_test

/*func TestXXX3(t *testing.T) {
	userDb := mdb_goson.New[TestUser]("../../trash/db")
	// userDb.Insert(TestUser{Id: int(userDb.GenerateId()), Username: "lox", Password: "oglox"})
	sr := userDb.FindBy(mdb_goson.ArgsFind[TestUser]{
		FieldList: "Id",
		Where: func(u *TestUser) bool {
			return u.Id == 3
		},
	})

	ml_console.PrettyPrint(sr.Unpack())
}*/

/*func TestDelete(t *testing.T) {
	userDb := mdb_goson.New[TestUser]("../../trash/db")

	name := ml_crypto.UID(12)

	// Insert
	userDb.Insert(TestUser{Id: int(userDb.GenerateId()), Username: name, Password: "oglox"})

	// Find
	sr := userDb.FindBy(mdb_goson.ArgsFind[TestUser]{
		FieldList: "Username",
		Where: func(u *TestUser) bool {
			return u.Username == name
		},
	})
	if !sr.IsFound {
		t.Fatalf("fuck")
	}
	if sr.Unpack()[0].Username != name {
		t.Fatalf("fuck")
	}

	// Delete record
	userDb.DeleteBy(mdb_goson.ArgsFind[TestUser]{
		FieldList: "Username",
		Where: func(u *TestUser) bool {
			return u.Username == name
		},
	})

	// Find again
	sr = userDb.FindBy(mdb_goson.ArgsFind[TestUser]{
		FieldList: "Username",
		Where: func(u *TestUser) bool {
			return u.Username == name
		},
	})
	if sr.IsFound {
		t.Fatalf("fuck")
	}

	ml_console.PrettyPrint(sr.Unpack())
}*/
