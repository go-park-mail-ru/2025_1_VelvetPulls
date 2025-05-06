package http_test

// func TestGetContacts_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockContactUC := mocks.NewMockIContactUsecase(ctrl)
// 	mockSessionUC := mocks.NewMockISessionUsecase(ctrl)

// 	userID := uuid.New()
// 	expectedContacts := []model.Contact{
// 		{
// 			ID:       uuid.MustParse("140768b8-1f0d-49a6-b7bd-f1f594dda332"),
// 			Username: "contact1",
// 		},
// 		{
// 			ID:       uuid.MustParse("4ee74e92-d4e6-4486-80db-f16d84a91100"),
// 			Username: "contact2",
// 		},
// 	}

// 	mockSessionUC.EXPECT().
// 		CheckLogin(gomock.Any(), "valid-token").
// 		Return(userID.String(), nil)

// 	mockContactUC.EXPECT().
// 		GetUserContacts(gomock.Any(), userID).
// 		Return(expectedContacts, nil)

// 	router := mux.NewRouter()
// 	delivery.NewContactController(router, mockContactUC, mockSessionUC)

// 	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
// 	req.AddCookie(&http.Cookie{
// 		Name:  "token",
// 		Value: "valid-token",
// 	})
// 	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	var resp utils.JSONResponse
// 	err := json.Unmarshal(rr.Body.Bytes(), &resp)
// 	assert.NoError(t, err)
// 	assert.True(t, resp.Status)

// 	var actualContacts []model.Contact
// 	jsonData, err := json.Marshal(resp.Data)
// 	assert.NoError(t, err)
// 	err = json.Unmarshal(jsonData, &actualContacts)
// 	assert.NoError(t, err)

// 	assert.Equal(t, expectedContacts, actualContacts)
// }

// func TestAddContact_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockContactUC := mocks.NewMockIContactUsecase(ctrl)
// 	mockSessionUC := mocks.NewMockISessionUsecase(ctrl)

// 	userID := uuid.New()
// 	contactData := model.RequestContact{
// 		Username: "new-contact",
// 	}

// 	mockSessionUC.EXPECT().
// 		CheckLogin(gomock.Any(), "valid-token").
// 		Return(userID.String(), nil)

// 	mockContactUC.EXPECT().
// 		AddUserContact(gomock.Any(), userID, contactData.Username).
// 		Return(nil)

// 	router := mux.NewRouter()
// 	delivery.NewContactController(router, mockContactUC, mockSessionUC)

// 	body, err := json.Marshal(contactData)
// 	assert.NoError(t, err)

// 	req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.AddCookie(&http.Cookie{
// 		Name:  "token",
// 		Value: "valid-token",
// 	})
// 	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	var resp utils.JSONResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &resp)
// 	assert.NoError(t, err)
// 	assert.True(t, resp.Status)
// 	assert.Equal(t, "Contact added successfully", resp.Data)
// }

// func TestDeleteContact_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockContactUC := mocks.NewMockIContactUsecase(ctrl)
// 	mockSessionUC := mocks.NewMockISessionUsecase(ctrl)

// 	userID := uuid.New()
// 	contactData := model.RequestContact{
// 		Username: "contact-to-delete",
// 	}

// 	mockSessionUC.EXPECT().
// 		CheckLogin(gomock.Any(), "valid-token").
// 		Return(userID.String(), nil)

// 	mockContactUC.EXPECT().
// 		RemoveUserContact(gomock.Any(), userID, contactData.Username).
// 		Return(nil)

// 	router := mux.NewRouter()
// 	delivery.NewContactController(router, mockContactUC, mockSessionUC)

// 	body, err := json.Marshal(contactData)
// 	assert.NoError(t, err)

// 	req := httptest.NewRequest(http.MethodDelete, "/contacts", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.AddCookie(&http.Cookie{
// 		Name:  "token",
// 		Value: "valid-token",
// 	})
// 	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	var resp utils.JSONResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &resp)
// 	assert.NoError(t, err)
// 	assert.True(t, resp.Status)
// 	assert.Equal(t, "Contact deleted successfully", resp.Data)
// }

// func TestGetContacts_Unauthorized(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockContactUC := mocks.NewMockIContactUsecase(ctrl)
// 	mockSessionUC := mocks.NewMockISessionUsecase(ctrl)

// 	router := mux.NewRouter()
// 	delivery.NewContactController(router, mockContactUC, mockSessionUC)

// 	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
// 	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusUnauthorized, rr.Code)
// }

// func TestAddContact_InvalidBody(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockContactUC := mocks.NewMockIContactUsecase(ctrl)
// 	mockSessionUC := mocks.NewMockISessionUsecase(ctrl)

// 	// The middleware should check the token first
// 	mockSessionUC.EXPECT().
// 		CheckLogin(gomock.Any(), "valid-token").
// 		Return(uuid.New().String(), nil)

// 	// We expect no calls to contactUsecase since the request should fail validation
// 	mockContactUC.EXPECT().
// 		AddUserContact(gomock.Any(), gomock.Any(), gomock.Any()).
// 		Times(0)

// 	router := mux.NewRouter()
// 	delivery.NewContactController(router, mockContactUC, mockSessionUC)

// 	req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBufferString("invalid json"))
// 	req.Header.Set("Content-Type", "application/json")
// 	req.AddCookie(&http.Cookie{
// 		Name:  "token",
// 		Value: "valid-token",
// 	})
// 	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusBadRequest, rr.Code)
// }
