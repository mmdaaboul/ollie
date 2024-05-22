func TestCheckAndCreateRecord(t *testing.T) {
	// Open a temporary BoltDB for testing
	db, err := bolt.Open("", 0666, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Call the CheckAndCreateRecord function
	err = CheckAndCreateRecord()
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the record was created correctly
	expectedDir, expectedFolder, _ := GetCurrentDirectoryAndFolderName()
	expectedRecords := map[string]string{expectedFolder: expectedDir}
	records, err := GetAllRecords()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(records, expectedRecords) {
		t.Errorf("GetAllRecords() returned %v, expected %v", records, expectedRecords)
	}
}