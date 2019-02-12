package resourcestore_test

// func testList(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	r := NewMockResourceStore(ctrl)

// 	tests := []struct {
// 		name        string
// 		expectError error
// 	}{
// 		{
// 			name:        "test list keys",
// 			expectError: nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		r.EXPECT().List(gomock.Any()).Return(t.expectError)

// 		NewKVResourceStore()

// 		assert.Equal(t, test.expectError, err)
// 	}
// }

// func TestMyThing(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	r := mockRes.NewMockResourceStore(mockCtrl)

// 	storeConfig := &libkvStore.Config{}

// 	mockS := new(libkvStoreMock.Mock)

// 	mockS.On("List", "").Return()

// 	kvStore := resourcestore.NewKVResourceStore(mockStore)

// 	mockObj := r.EXPECT().List("cluster/gojek/")
// }
