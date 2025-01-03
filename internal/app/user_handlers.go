package app

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/cash"
	"github.com/ditacijsvitvidadoa/backend/internal/cookie"
	"github.com/ditacijsvitvidadoa/backend/internal/email_sender"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/ditacijsvitvidadoa/backend/internal/storage/requests"
	"github.com/ditacijsvitvidadoa/backend/internal/utils"
	"github.com/google/uuid"
	password2 "github.com/vzglad-smerti/password_hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
)

func (a *App) createUserAccount(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	patronymic := r.FormValue("patronymic")
	phoneNumber := r.FormValue("phoneNumber")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashPassword, err := password2.Hash(password)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	isExists, err := requests.IsEmailExists(a.client, email)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if isExists {
		sendWarning(w, http.StatusConflict)
		return
	}

	userID, err := utils.GetNextUserID(a.client)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to generate user ID.")
		return
	}

	newUser := entities.User{
		UserID:           userID,
		Password:         hashPassword,
		FullName:         entities.FullName{FirstName: firstName, LastName: lastName, Patronymic: patronymic},
		Phone:            phoneNumber,
		Email:            email,
		MarketingConsent: false,
		Cart:             []entities.CartItem{},
		Favourites:       []string{},
	}

	_, err = requests.CreateNewUser(a.client, newUser)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create user account.")
		return
	}

	sendOk(w)
}

func (a *App) getProfileInfo(w http.ResponseWriter, r *http.Request) {
	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	UserInfo, err := requests.GetUserByID(a.client, userObjectId)
	if err != nil {
		sendError(w, http.StatusNoContent, "Failed to retrieve user info from storage.")
		return
	}

	// Отладочные сообщения
	fmt.Println("User info retrieved:", UserInfo)

	// Проверка на наличие PostalServiceInfo и PostalInfo
	if UserInfo.PostalServiceInfo != nil && UserInfo.PostalServiceInfo.PostalInfo != nil {
		var postalInfoArray []map[string]string
		switch postalInfo := UserInfo.PostalServiceInfo.PostalInfo.(type) {
		case map[string]interface{}:
			for key, value := range postalInfo {
				postalInfoArray = append(postalInfoArray, map[string]string{
					"Key":   key,
					"Value": value.(string),
				})
			}
			UserInfo.PostalServiceInfo.PostalInfo = postalInfoArray
		case string:
			fmt.Println("PostalInfo is a string:", postalInfo)
		default:
			fmt.Println("Unknown type of PostalInfo")
		}
	} else {
		fmt.Println("PostalServiceInfo or PostalInfo is nil")
	}

	// Отправка ответа
	sendResponse(w, UserInfo)
}

func (a *App) logIn(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	log.Printf("Email: %s, Password: %s\n", email, password)

	if email == "" || password == "" {
		sendError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	userId, err := requests.LogInAccount(a.client, email, password)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateJWT(userId.Hex())
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error generating token: %s", err))
		return
	}

	sessionID := uuid.New().String()
	cookieValue := fmt.Sprintf("userId=%s|token=%s", userId.Hex(), token)
	sessionKey := fmt.Sprintf("session:%s:%s", userId.Hex(), sessionID)

	err = cash.SaveSessionToRedis(a.cash.Conn, sessionKey, token)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Could not save session to redis: %s", err))
		return
	}

	log.Printf("Session Key: %s, Cookie Value: %s\n", sessionKey, cookieValue)
	cookie.SetCookie(w, "session", userId.Hex(), token)

	sendOk(w)
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	cookieValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value.")
		return
	}

	parts := strings.Split(cookieValue, "|")
	if len(parts) != 2 {
		sendError(w, http.StatusUnauthorized, "Invalid cookie format")
		return
	}

	var userId, token string
	for _, part := range parts {
		keyValue := strings.Split(part, "=")
		if len(keyValue) != 2 {
			sendError(w, http.StatusUnauthorized, "Invalid key-value format in cookie")
			return
		}

		switch keyValue[0] {
		case "userId":
			userId = keyValue[1]
		case "token":
			token = keyValue[1]
		default:
			sendError(w, http.StatusUnauthorized, "Unexpected key in cookie")
			return
		}
	}

	if token == "" {
		sendError(w, http.StatusUnauthorized, "Token not found in cookie")
		return
	}

	if userId == "" {
		sendError(w, http.StatusUnauthorized, "User ID not found in cookie")
		return
	}

	err = cash.DeleteSessionByToken(a.cash.Conn, token)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting session: %s", err))
		return
	}

	cookie.ClearSessionCookie(w)

	sendOk(w)
}

func (a *App) checkAuthentication(w http.ResponseWriter, r *http.Request) {
	valid, _, err := cookie.CheckCookie(a.cash.Conn, r)

	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error checking cookie: %s", err))
		return
	}
	if !valid {
		sendError(w, http.StatusUnauthorized, "Invalid cookie")
		return
	}

	sendOk(w)
}

func (a *App) PurchasesHistory(w http.ResponseWriter, r *http.Request) {
	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	orders, err := requests.GetPurchaseHistory(a.client, userId)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to get purchase history")
		return
	}

	var historyResponse []struct {
		OrderDetails struct {
			OrderId int    `json:"order_id"`
			Date    string `json:"date"`
			Status  string `json:"status"`
		} `json:"order_details"`
		Products []entities.Product `json:"products"`
	}

	for _, order := range orders {
		status := getStatusText(order.Status)

		responseItem := struct {
			OrderDetails struct {
				OrderId int    `json:"order_id"`
				Date    string `json:"date"`
				Status  string `json:"status"`
			} `json:"order_details"`
			Products []entities.Product `json:"products"`
		}{
			OrderDetails: struct {
				OrderId int    `json:"order_id"`
				Date    string `json:"date"`
				Status  string `json:"status"`
			}{
				OrderId: order.OrderId,
				Date:    order.Date.Format("2006-01-02"),
				Status:  status,
			},
			Products: order.Products,
		}

		historyResponse = append(historyResponse, responseItem)
	}

	sendResponse(w, historyResponse)
}

func getStatusText(status int) string {
	switch status {
	case 0:
		return "Не обработано"
	case 1:
		return "В процессе"
	case 2:
		return "Завершено"
	default:
		return "Неизвестно"
	}
}

func (a *App) updateFirstName(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("firstname")

	if firstName == "" {
		sendError(w, http.StatusBadRequest, "firstName is required")
		return
	}

	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	err = requests.UpdateUserProfileField(a.client, userObjectId, "FullName.FirstName", firstName)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user profile: %s", err))
		return
	}

	sendOk(w)
}

func (a *App) updateLastName(w http.ResponseWriter, r *http.Request) {
	lastName := r.FormValue("lastname")

	if lastName == "" {
		sendError(w, http.StatusBadRequest, "lastName is required")
		return
	}

	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	err = requests.UpdateUserProfileField(a.client, userObjectId, "FullName.LastName", lastName)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user profile: %s", err))
		return
	}

	sendOk(w)
}

func (a *App) updatePatronymic(w http.ResponseWriter, r *http.Request) {
	patronymic := r.FormValue("patronymic")

	if patronymic == "" {
		sendError(w, http.StatusBadRequest, "patronymic is required")
		return
	}

	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	err = requests.UpdateUserProfileField(a.client, userObjectId, "FullName.Patronymic", patronymic)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user profile: %s", err))
		return
	}

	sendOk(w)
}

func (a *App) updatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	phoneNumber := r.FormValue("phone")

	if phoneNumber == "" {
		sendError(w, http.StatusBadRequest, "phone-number is required")
		return
	}

	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	err = requests.UpdateUserProfileField(a.client, userObjectId, "Phone", phoneNumber)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user profile: %s", err))
		return
	}

	sendOk(w)
}

func (a *App) updateEmail(w http.ResponseWriter, r *http.Request) {
	Email := r.FormValue("email")

	if Email == "" {
		sendError(w, http.StatusBadRequest, "email is required")
		return
	}

	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	err = requests.UpdateUserProfileField(a.client, userObjectId, "Email", Email)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user profile: %s", err))
		return
	}

	sendOk(w)
}

func (a *App) updatePassword(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")

	if password == "" {
		sendError(w, http.StatusBadRequest, "password is required")
		return
	}

	hashPassword, err := password2.Hash(password)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error hashing password: %s", err))
		return
	}

	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	err = requests.UpdateUserProfileField(a.client, userObjectId, "Password", hashPassword)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user profile: %s", err))
		return
	}

	sendOk(w)
}

func (a *App) updateMarketingConsent(w http.ResponseWriter, r *http.Request) {
	marketingConsent := r.FormValue("marketing-consent")
	var marketingConsentBool bool

	switch marketingConsent {
	case "yes":
		marketingConsentBool = true
	case "no":
		marketingConsentBool = false
	default:
		marketingConsentBool = false
	}

	if marketingConsent == "" {
		sendError(w, http.StatusBadRequest, "marketingConsent is required")
		return
	}

	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	err = requests.UpdateUserProfileField(a.client, userObjectId, "MarketingConsent", marketingConsentBool)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user profile: %s", err))
		return
	}

	sendOk(w)
}

func (a *App) addOrUpdatePostalServiceInfo(w http.ResponseWriter, r *http.Request) {
	postalType := r.FormValue("postal_type")
	city := r.FormValue("city")
	cityRef := r.FormValue("city_ref")
	receivingType := r.FormValue("receiving_type")

	if postalType == "" || city == "" || cityRef == "" || receivingType == "" {
		sendError(w, http.StatusBadRequest, "All fields are required.")
		return
	}

	var postalInfo interface{}

	if receivingType == "Branches" {
		postalInfoValue := r.FormValue("postal_info")
		log.Println(postalInfoValue)
		if postalInfoValue == "" {
			sendError(w, http.StatusBadRequest, "postal_info is required for Branches.")
			return
		}
		postalInfo = postalInfoValue
	} else if receivingType == "Courier" {
		street := r.FormValue("street")
		house := r.FormValue("house")
		apartment := r.FormValue("apartment")
		floor := r.FormValue("floor")
		log.Println(street, house, apartment, floor)

		if street == "" || house == "" {
			sendError(w, http.StatusBadRequest, "street, house are required for Courier.")
			return
		}

		postalInfo = map[string]string{
			"street":    street,
			"house":     house,
			"apartment": apartment,
			"floor":     floor,
		}
	} else {
		sendError(w, http.StatusBadRequest, "Invalid receiving_type. Must be 'Branches' or 'Courier'.")
		return
	}

	sessionValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Unable to retrieve session value. Please ensure you are logged in.")
		return
	}

	userId, err := cookie.GetUserIDFromCookie(sessionValue)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "Failed to retrieve user ID from session cookie.")
		return
	}

	postalServiceInfo := entities.PostalServiceInfo{
		PostalType:    postalType,
		City:          city,
		CityRef:       cityRef,
		ReceivingType: receivingType,
		PostalInfo:    postalInfo,
	}

	err = requests.AddOrUpdatePostalServiceInfo(a.client, userObjectId, postalServiceInfo)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding or updating postal service info: %s", err))
		return
	}

	sendOk(w)
}

func (a *App) sendToSupport(w http.ResponseWriter, r *http.Request) {
	Name := r.FormValue("Name")
	Email := r.FormValue("Email")
	Phone := r.FormValue("Phone")
	Title := r.FormValue("Title")
	Description := r.FormValue("Description")

	go func() {
		err := email_sender.SendSupportQuestion(Name, Email, Phone, Title, Description)
		if err != nil {
			sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error sending support question: %s", err))
			return
		}
	}()

	sendOk(w)
}
