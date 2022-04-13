package datastorage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_ImplementsDataStorage(t *testing.T) {
	var _ DataStorage = &MemStorage{}
}

func TestMemStorage_Initialize(t *testing.T) {
	mem := MemStorage{}.Initialize()
	assert.NotNil(t, mem)
	assert.NotNil(t, mem.data)
	assert.NotZero(t, mem.data)
	assert.NotNil(t, mem.rwMu)
}

func TestMemStorage_RetrieveData(t *testing.T) {
	// Embed test data
	mem := MemStorage{}.Initialize()
	testData := []byte("test data")
	mem.data["test"] = testData

	t.Run("existing data", func(t *testing.T) {
		// Check correct data returned
		data, err := mem.RetrieveData("test")
		require.NoError(t, err)
		assert.Equal(t, testData, data)
	})

	t.Run("nonexistent key", func(t *testing.T) {
		// Key does not exist, should error
		_, err := mem.RetrieveData("foo")
		assert.Error(t, err)
	})
}

func TestMemStorage_StoreData(t *testing.T) {
	mem := MemStorage{}.Initialize()
	dataName := "test"
	testData := []byte("test data")

	t.Run("first write", func(t *testing.T) {
		// Store data + check directly for correctness
		err := mem.StoreData(dataName, testData)
		require.NoError(t, err)
		assert.Len(t, mem.data, 1)

		data := mem.data[dataName]
		assert.Equal(t, testData, data)
	})

	t.Run("overwrite", func(t *testing.T) {
		// Ensure StoreData properly overwrites with a duplicate key, does not
		// create a second key value pair
		testData = []byte("foobar")
		err := mem.StoreData(dataName, testData)
		require.NoError(t, err)
		assert.Len(t, mem.data, 1)

		data := mem.data[dataName]
		assert.Equal(t, testData, data)
	})

	t.Run("new data", func(t *testing.T) {
		// Check StoreData can store more than 1 unique data instance
		dataName2 := "test2"
		testData2 := []byte("qwerty")
		err := mem.StoreData(dataName2, testData2)
		require.NoError(t, err)

		assert.Len(t, mem.data, 2)
		data := mem.data[dataName2]
		assert.Equal(t, testData2, data)
	})
}

func TestMemStorage_DeleteData(t *testing.T) {
	mem := MemStorage{}.Initialize()
	dataName := "test"
	testData := []byte("test data")

	// Embed test data
	err := mem.StoreData(dataName, testData)
	require.NoError(t, err)

	t.Run("delete existing data", func(t *testing.T) {
		// Check data actually deleted
		err = mem.DeleteData(dataName)
		require.NoError(t, err)

		assert.Len(t, mem.data, 0)
	})

	t.Run("delete nonexistent data", func(t *testing.T) {
		// Check DeleteData errors on bad key
		err = mem.DeleteData(dataName)
		assert.Error(t, err)
		assert.Len(t, mem.data, 0)
	})
}

func TestMemStorage_ThreadSafe(t *testing.T) {
	// -race test flag will detect race conditions
	var wg sync.WaitGroup
	start := make(chan struct{})

	mem := MemStorage{}.Initialize()
	dataName := "test"
	testData := []byte("test data")

	for i := 0; i < 69; i++ {
		wg.Add(1)

		go func(itr int) {
			defer wg.Done()
			<-start
			if itr%3 == 0 {
				_ = mem.StoreData(dataName, testData)
			} else if itr%3 == 1 {
				_, _ = mem.RetrieveData(dataName)
			} else {
				_ = mem.DeleteData(dataName)
			}
		}(i)
	}

	close(start) // Start after all workers created
	wg.Wait()    // Wait for workers to complete
}
