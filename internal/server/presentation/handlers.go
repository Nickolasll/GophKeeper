package presentation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/Nickolasll/goph-keeper/internal/server/domain"
)

// @Summary Регистрация нового пользователя по логину и паролю
// @ID auth-register
// @Tags Auth
// @Accept json
// @Param payload body registrationPayload true "Логин и Пароль"
// @Success 200
// @Failure 400 "Некорректный формат данных"
// @Failure 409 "Логин уже занят"
// @Header 200 {string} Authorization eyJhbGciOiJI...qIScZUU8P0Zhck "JWT"
// @Router /auth/register [post]
func registrationHandler(w http.ResponseWriter, r *http.Request) { //nolint: dupl
	var payload registrationPayload
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	token, err := app.Registration.Do(payload.Login, payload.Password)
	if err != nil {
		if errors.Is(err, domain.ErrLoginAlreadyInUse) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}

	w.Header().Set("Authorization", string(token))
	w.WriteHeader(http.StatusOK)
}

// @Summary Авторизация пользователя по логину и паролю
// @ID auth-login
// @Tags Auth
// @Accept json
// @Param payload body registrationPayload true "Логин и Пароль"
// @Success 200
// @Failure 400 "Некорректный формат данных"
// @Failure 401 "Неправильный логин или пароль"
// @Header 200 {string} Authorization eyJhbGciOiJI...qIScZUU8P0Zhck "JWT"
// @Router /auth/login [post]
func loginHandler(w http.ResponseWriter, r *http.Request) { //nolint: dupl
	var payload registrationPayload
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	token, err := app.Login.Do(payload.Login, payload.Password)
	if err != nil {
		if errors.Is(err, domain.ErrLoginOrPasswordIsInvalid) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}

	w.Header().Set("Authorization", string(token))
	w.WriteHeader(http.StatusOK)
}

// @Summary Создать и зашифровать текстовые данные
// @ID text-create
// @Tags Text
// @Accept plain
// @Param data body string true "Текст для сохранения"
// @Success 201
// @Failure 400 "Некорректный формат данных"
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Header 201 {string} Location 020cb30c-c495-4a18-ac09-fd68c6f7c941 "UUID ресурса"
// @Router /text/create [post]
// @Security ApiKeyAuth
func createTextHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	body, err := parseBody(textType, r)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	textID, err := app.CreateText.Do(userID, string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Location", textID.String())
	w.WriteHeader(http.StatusCreated)
}

// @Summary Обновить и зашифровать существующие текстовые данные
// @ID text-update
// @Tags Text
// @Accept plain
// @Param text_id path string true "Text ID"
// @Param data body string true "Текст для сохранения"
// @Success 200
// @Failure 400 "Некорректный формат данных или идентификатора"
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Failure 404 "Не найдено"
// @Router /text/{text_id} [post]
// @Security ApiKeyAuth
func updateTextHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	id, err := getRouteID(r, "textID")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}

	body, err := parseBody(textType, r)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	err = app.UpdateText.Do(userID, id, string(body))
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Получить все расшифрованные текстовые данные
// @ID text-all
// @Tags Text
// @Success 200 {object} GetAllTextsResponse
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Router /text/all [get]
// @Security ApiKeyAuth
func getAllTextsHandler(w http.ResponseWriter, _ *http.Request, userID uuid.UUID) {
	textsResponse := []textResponse{}
	w.Header().Set(contentTypeHeader, jsonType)
	texts, err := app.GetAllTexts.Do(userID)
	if err != nil {
		log.Error(err)
		err = responseError(w, err.Error())
		if err != nil {
			log.Error(err)
		}

		return
	}

	for _, v := range texts {
		respItem := textResponse{
			ID:      v.ID.String(),
			Content: string(v.Content),
		}
		textsResponse = append(textsResponse, respItem)
	}

	response := GetAllTextsResponse{
		Status: true,
	}
	response.Data.Texts = textsResponse

	err = makeResponse(w, http.StatusOK, response)
	if err != nil {
		log.Error(err)
		err = responseError(w, err.Error())

		if err != nil {
			log.Error(err)
		}

		return
	}
}

// @Summary Получение публичного ключа для валидации JWT на клиенте
// @ID auth-certs
// @Tags Auth
// @Success 200 {object} object{kty=string,alg=string,x=string,y=string,use=string,kid=string}
// @Router /auth/certs [get]
func getCertsHandler(w http.ResponseWriter, _ *http.Request) {
	certs, err := joseService.GetCerts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	resp, err := json.Marshal(certs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
	}
}

// @Summary Запрос состояния сервиса
// @Tags Status
// @ID health
// @Success 200
// @Failure 500
// @Router /health [get]
func getHealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// @Summary Создать и зашифровать бинарные данные
// @ID binary-create
// @Tags Binary
// @Accept mpfd
// @Param data body []byte true "Содержимое файла"
// @Success 201
// @Failure 400 "Некорректный формат данных"
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Header 201 {string} Location 020cb30c-c495-4a18-ac09-fd68c6f7c941 "UUID ресурса"
// @Router /binary/create [post]
// @Security ApiKeyAuth
func createBinaryHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	body, err := parseBody(binaryType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	binID, err := app.CreateBinary.Do(userID, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Location", binID.String())
	w.WriteHeader(http.StatusCreated)
}

// @Summary Обновить и зашифровать существующие бинарные данные
// @ID binary-update
// @Tags Binary
// @Accept mpfd
// @Param binary_id path string true "Binary ID"
// @Param data body []byte true "Содержимое файла"
// @Success 200
// @Failure 400 "Некорректный формат данных или идентификатора"
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Failure 404 "Не найдено"
// @Router /binary/{binary_id} [post]
// @Security ApiKeyAuth
func updateBinaryHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	id, err := getRouteID(r, "binaryID")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	body, err := parseBody(binaryType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	err = app.UpdateBinary.Do(userID, id, body)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Получить все расшифрованные бинарные данные
// @ID binary-all
// @Tags Binary
// @Success 200 {object} GetAllBinariesResponse
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Router /binary/all [get]
// @Security ApiKeyAuth
func getAllBinariesHandler(w http.ResponseWriter, _ *http.Request, userID uuid.UUID) {
	binariesResponse := []binaryResponse{}
	w.Header().Set(contentTypeHeader, jsonType)
	binaries, err := app.GetAllBinaries.Do(userID)
	if err != nil {
		log.Error(err)
		err = responseError(w, err.Error())
		if err != nil {
			log.Error(err)
		}

		return
	}

	for _, v := range binaries {
		respItem := binaryResponse{
			ID:      v.ID.String(),
			Content: v.Content,
		}
		binariesResponse = append(binariesResponse, respItem)
	}

	response := GetAllBinariesResponse{
		Status: true,
	}
	response.Data.Binaries = binariesResponse

	err = makeResponse(w, http.StatusOK, response)
	if err != nil {
		log.Error(err)
		err = responseError(w, err.Error())

		if err != nil {
			log.Error(err)
		}

		return
	}
}

// @Summary Создать и зашифровать логин и пароль
// @ID credentials-create
// @Tags Credentials
// @Accept json
// @Param data body credentialsPayload true "Наименование, логин и пароль"
// @Success 201
// @Failure 400 "Некорректный формат данных"
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Header 201 {string} Location 020cb30c-c495-4a18-ac09-fd68c6f7c941 "UUID ресурса"
// @Router /credentials/create [post]
// @Security ApiKeyAuth
func createCredentialsHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var payload credentialsPayload
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	credID, err := app.CreateCredentials.Do(
		userID,
		payload.Name,
		payload.Login,
		payload.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Location", credID.String())
	w.WriteHeader(http.StatusCreated)
}

// @Summary Обновить и зашифровать существующий логин и пароль
// @ID credentials-update
// @Tags Credentials
// @Accept json
// @Param credentials_id path string true "Credentials ID"
// @Param data body credentialsPayload true "Наименование, логин и пароль"
// @Success 200
// @Failure 400 "Некорректный формат данных или идентификатора"
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Failure 404 "Не найдено"
// @Router /credentials/{credentials_id} [post]
// @Security ApiKeyAuth
func updateCredentialsHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var payload credentialsPayload
	id, err := getRouteID(r, "credID")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	err = app.UpdateCredentials.Do(
		userID,
		id,
		payload.Name,
		payload.Login,
		payload.Password,
	)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Получить все расшифрованные логины и пароли
// @ID credentials-all
// @Tags Credentials
// @Success 200 {object} GetAllCredentialsResponse
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Router /credentials/all [get]
// @Security ApiKeyAuth
func getAllCredentialsHandler(w http.ResponseWriter, _ *http.Request, userID uuid.UUID) {
	credResponse := []credentialsResponse{}
	w.Header().Set(contentTypeHeader, jsonType)
	credentials, err := app.GetAllCredentials.Do(userID)
	if err != nil {
		log.Error(err)
		err = responseError(w, err.Error())
		if err != nil {
			log.Error(err)
		}

		return
	}

	for _, v := range credentials {
		respItem := credentialsResponse{
			ID:       v.ID.String(),
			Name:     string(v.Name),
			Login:    string(v.Login),
			Password: string(v.Password),
		}
		credResponse = append(credResponse, respItem)
	}

	response := GetAllCredentialsResponse{
		Status: true,
	}
	response.Data.Credentials = credResponse

	err = makeResponse(w, http.StatusOK, response)
	if err != nil {
		log.Error(err)
		err = responseError(w, err.Error())

		if err != nil {
			log.Error(err)
		}

		return
	}
}

// @Summary Создать и зашифровать банковскую карту
// @ID bank-card-create
// @Tags BankCard
// @Accept json
// @Param data body bankCardPayload true "Номер, срок действия, cvv код, ФИО держателя карты"
// @Success 201
// @Failure 400 "Некорректный формат данных"
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Header 201 {string} Location 020cb30c-c495-4a18-ac09-fd68c6f7c941 "UUID ресурса"
// @Router /bank_card/create [post]
// @Security ApiKeyAuth
func createBankCardHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var payload bankCardPayload
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	cardID, err := app.CreateBankCard.Do(
		userID,
		payload.Number,
		payload.ValidThru,
		payload.CVV,
		payload.CardHolder,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)

		return
	}
	w.Header().Add("Location", cardID.String())
	w.WriteHeader(http.StatusCreated)
}

// @Summary Обновить и зашифровать существующую банковскую карту
// @ID bank-card-update
// @Tags BankCard
// @Accept json
// @Param bank_card_id path string true "Bank Card ID"
// @Param data body bankCardPayload true "Номер, срок действия, cvv код, ФИО держателя карты"
// @Success 200
// @Failure 400 "Некорректный формат данных или идентификатора"
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Failure 404 "Не найдено"
// @Router /bank_card/{bank_card_id} [post]
// @Security ApiKeyAuth
func updateBankCardHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var payload bankCardPayload
	id, err := getRouteID(r, "cardID")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	body, err := parseBody(jsonType, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	payload, err = payload.Load(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)

		return
	}
	err = app.UpdateBankCard.Do(
		userID,
		id,
		payload.Number,
		payload.ValidThru,
		payload.CVV,
		payload.CardHolder,
	)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(err)
		}

		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Получить все расшифрованные банковские карты
// @ID bank-card-all
// @Tags BankCard
// @Success 200 {object} GetAllBankCardsResponse
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Router /bank_card/all [get]
// @Security ApiKeyAuth
func getAllBankCardsHandler(w http.ResponseWriter, _ *http.Request, userID uuid.UUID) {
	bankCardsResponse := []bankCardResponse{}
	w.Header().Set(contentTypeHeader, jsonType)
	bankCards, err := app.GetAllBankCards.Do(userID)
	if err != nil {
		log.Error(err)
		err = responseError(w, err.Error())
		if err != nil {
			log.Error(err)
		}

		return
	}

	for _, v := range bankCards {
		respItem := bankCardResponse{
			ID:         v.ID.String(),
			Number:     string(v.Number),
			ValidThru:  string(v.ValidThru),
			CVV:        string(v.CVV),
			CardHolder: string(v.CardHolder),
		}
		bankCardsResponse = append(bankCardsResponse, respItem)
	}

	response := GetAllBankCardsResponse{
		Status: true,
	}
	response.Data.BankCards = bankCardsResponse

	err = makeResponse(w, http.StatusOK, response)
	if err != nil {
		log.Error(err)
		err = responseError(w, err.Error())

		if err != nil {
			log.Error(err)
		}

		return
	}
}

// @Summary Получить все расшифрованные данные пользователя
// @ID all
// @Tags All
// @Success 200 {object} GetAllResponse
// @Failure 401 "Нет токена авторизации, либо токен невалиден"
// @Router /all [get]
// @Security ApiKeyAuth
func getAllHandler(w http.ResponseWriter, _ *http.Request, userID uuid.UUID) {
	credResponse := []credentialsResponse{}
	bankCardsResponse := []bankCardResponse{}
	textsResponse := []textResponse{}
	binariesResponse := []binaryResponse{}
	w.Header().Set(contentTypeHeader, jsonType)
	texts, bankCards, binaries, credentials, err := app.GetAll.Do(userID)
	if err != nil {
		log.Error(err)
		err = responseError(w, err.Error())
		if err != nil {
			log.Error(err)
		}

		return
	}

	for _, v := range bankCards {
		respItem := bankCardResponse{
			ID:         v.ID.String(),
			Number:     string(v.Number),
			ValidThru:  string(v.ValidThru),
			CVV:        string(v.CVV),
			CardHolder: string(v.CardHolder),
		}
		bankCardsResponse = append(bankCardsResponse, respItem)
	}
	for _, v := range credentials {
		respItem := credentialsResponse{
			ID:       v.ID.String(),
			Name:     string(v.Name),
			Login:    string(v.Login),
			Password: string(v.Password),
		}
		credResponse = append(credResponse, respItem)
	}
	for _, v := range texts {
		respItem := textResponse{
			ID:      v.ID.String(),
			Content: string(v.Content),
		}
		textsResponse = append(textsResponse, respItem)
	}
	for _, v := range binaries {
		respItem := binaryResponse{
			ID:      v.ID.String(),
			Content: v.Content,
		}
		binariesResponse = append(binariesResponse, respItem)
	}

	response := GetAllResponse{
		Status: true,
	}
	response.Data.BankCards = bankCardsResponse
	response.Data.Credentials = credResponse
	response.Data.Texts = textsResponse
	response.Data.Binaries = binariesResponse

	err = makeResponse(w, http.StatusOK, response)
	if err != nil {
		log.Error(err)
		err = responseError(w, err.Error())

		if err != nil {
			log.Error(err)
		}

		return
	}
}
