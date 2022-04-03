package datastorage

import "testing"

func TestMemStorage_ImplementsDataStorage(t *testing.T) {
	var _ DataStorage = &MemStorage{}
}

func TestMemStorage_Initialize(t *testing.T) {
	mem := MemStorage{}.Initialize()
	if mem == nil {
		t.Error("returned a nil ptr")
	}

	if mem.data == nil || len(mem.data) != 0 {
		t.Error("failed to initialize data map")
	}

	if mem.rwMu == nil {
		t.Error("failed to initialize rwMutex")
	}
}

func TestMemStorage_RetrieveData(t *testing.T) {
	mem := MemStorage{}.Initialize()
	testData := []byte("test data")
	mem.data["test"] = testData

	t.Run("existing data", func(t *testing.T) {
		data, err := mem.RetrieveData("test")
		if err != nil {
			t.Fatalf("received unexpected error: %v", err)
		}

		if string(testData) != string(data) {
			t.Errorf(
				"expected data: %s but got: %s",
				testData, data,
			)
		}
	})

	t.Run("nonexistent key", func(t *testing.T) {
		_, err := mem.RetrieveData("foo")
		if err == nil {
			t.Error("expected error but got none")
		}
	})
}

func TestMemStorage_StoreData(t *testing.T) {
	mem := MemStorage{}.Initialize()
	dataName := "test"
	testData := []byte("test data")

	t.Run("first write", func(t *testing.T) {
		err := mem.StoreData(dataName, testData)
		if err != nil {
			t.Fatalf("received unexpected error: %v", err)
		}

		if len(mem.data) != 1 {
			t.Errorf(
				"expected MemStorage data capacity of: %d, but got: %d",
				1, len(mem.data),
			)
		}

		data := mem.data[dataName]
		if string(testData) != string(data) {
			t.Errorf(
				"expected data: %s but got: %s",
				testData, data,
			)
		}
	})

	t.Run("overwrite", func(t *testing.T) {
		testData = []byte("foobar")
		err := mem.StoreData(dataName, testData)
		if err != nil {
			t.Fatalf("received unexpected error: %v", err)
		}

		if len(mem.data) != 1 {
			t.Errorf(
				"expected MemStorage data capacity of: %d, but got: %d",
				1, len(mem.data),
			)
		}

		data := mem.data[dataName]
		if string(testData) != string(data) {
			t.Errorf(
				"expected data: %s but got: %s",
				testData, data,
			)
		}
	})

	t.Run("new data", func(t *testing.T) {
		dataName2 := "test2"
		testData2 := []byte("qwerty")
		err := mem.StoreData(dataName2, testData2)
		if err != nil {
			t.Fatalf("received unexpected error: %v", err)
		}

		if len(mem.data) != 2 {
			t.Errorf(
				"expected MemStorage data capacity of: %d, but got: %d",
				2, len(mem.data),
			)
		}

		data := mem.data[dataName2]
		if string(testData2) != string(data) {
			t.Errorf(
				"expected data: %s but got: %s",
				testData2, data,
			)
		}
	})
}

func TestMemStorage_DeleteData(t *testing.T) {
	mem := MemStorage{}.Initialize()
	dataName := "test"
	testData := []byte("test data")

	err := mem.StoreData(dataName, testData)
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}

	t.Run("delete existing data", func(t *testing.T) {
		err = mem.DeleteData(dataName)
		if err != nil {
			t.Fatalf("received unexpected error: %v", err)
		}

		if len(mem.data) != 0 {
			t.Error("MemStorage data should be empty")
		}
	})

	t.Run("delete nonexistent data", func(t *testing.T) {
		err = mem.DeleteData(dataName)
		if err == nil {
			t.Error("expected error but got none")
		}

		if len(mem.data) != 0 {
			t.Error("MemStorage data should be empty")
		}
	})
}
