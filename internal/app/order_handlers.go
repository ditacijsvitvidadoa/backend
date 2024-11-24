package app

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/cookie"
	"github.com/ditacijsvitvidadoa/backend/internal/email_sender"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/ditacijsvitvidadoa/backend/internal/storage/requests"
	"github.com/ditacijsvitvidadoa/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (a *App) AddOrder(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	patronymic := r.FormValue("patronymic")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	postalType := r.FormValue("postal_type")
	city := r.FormValue("city")
	receivingType := r.FormValue("receiving_type")

	if firstName == "" || lastName == "" || patronymic == "" || phone == "" || email == "" || postalType == "" || city == "" {
		sendError(w, http.StatusBadRequest, "one or more values are empty")
		return
	}

	var order = entities.Order{
		OrderId:       utils.GenerateRandomNumber(),
		Status:        1,
		FirstName:     firstName,
		LastName:      lastName,
		Patronymic:    patronymic,
		Phone:         phone,
		Email:         email,
		PostalType:    postalType,
		City:          city,
		ReceivingType: receivingType,
		Date:          time.Now(),
	}

	if receivingType == "Branches" {
		postalInfo := r.FormValue("postal_info")
		if postalInfo == "" {
			sendError(w, http.StatusBadRequest, "postal_info is empty")
			return
		}
		order.PostalInfo = postalInfo
	} else if receivingType == "Courier" {
		street := r.FormValue("street")
		house := r.FormValue("house")
		apartment := r.FormValue("apartment")
		floor := r.FormValue("floor")

		order.PostalInfo = entities.CourierPostalInfo{
			Street:    street,
			House:     house,
			Apartment: apartment,
			Floor:     floor,
		}
	}

	var products []entities.Product
	for i := 0; ; i++ {
		title := r.FormValue(fmt.Sprintf("products[%d][Title]", i))
		if title == "" {
			break
		}

		priceStr := r.FormValue(fmt.Sprintf("products[%d][Price]", i))
		discountStr := r.FormValue(fmt.Sprintf("products[%d][Discount]", i))
		countStr := r.FormValue(fmt.Sprintf("products[%d][Count]", i))

		log.Println("countStr", countStr)

		count, err := strconv.Atoi(countStr)
		if err != nil {
			sendError(w, http.StatusBadRequest, "invalid count")
			return
		}

		price, err := strconv.Atoi(priceStr)
		if err != nil {
			sendError(w, http.StatusBadRequest, "invalid price")
			return
		}

		discount, err := strconv.Atoi(discountStr)
		if err != nil {
			sendError(w, http.StatusBadRequest, "invalid discount")
			return
		}

		imageUrls := r.Form[fmt.Sprintf("products[%d][Image_urls][]", i)]

		products = append(products, entities.Product{
			Title:     title,
			Price:     price,
			Discount:  discount,
			ImageUrls: imageUrls,
			Count:     count,
		})
	}
	order.Products = products

	cookieValue, err := cookie.GetSessionValue(r, "session")
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userId, err := cookie.GetUserIDFromCookie(cookieValue)
	if err == nil {
		if userId != "" {
			order.UserId = userId
		}
	} else {
		fmt.Println("err", err)
	}

	log.Println("order", order)

	orderId, err := requests.CreateNewOrder(a.client, order)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println(orderId)

	go func() {
		email_sender.SendOrderConfirmation(order.Email, order.FirstName, order.OrderId)
	}()

	sendOk(w)
}

func (a *App) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := requests.GetAll(a.client, "Orders")
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}

	sendResponse(w, orders)
}

func (a *App) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")

	objectID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Failed to create objectId")
	}

	deleteCount, err := requests.DeleteByObjectID(a.client, "Orders", objectID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if deleteCount == 0 {
		sendError(w, http.StatusBadRequest, "nothing to delete")
		return
	}

	sendOk(w)
}

func (a *App) ChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	statusStr := r.URL.Query().Get("status")
	productIDStr := r.PathValue("id")

	objectID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Failed to create objectId")
	}

	status, err := strconv.Atoi(statusStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid status")
		return
	}

	updated, err := requests.UpdateOrderStatus(a.client, objectID, status)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if updated == 0 {
		sendError(w, http.StatusBadRequest, "nothing to update")
		return
	}

	sendOk(w)
}

func (a *App) ArchiveOrder(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")

	objectID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Failed to create objectId")
	}

	updated, err := requests.ArchiveOrder(a.client, objectID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if updated == 0 {
		sendError(w, http.StatusBadRequest, "nothing to update")
	}

	sendOk(w)
}

func (a *App) GetArchiveOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := requests.GetAll(a.client, "ArchiveOrders")
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}

	sendResponse(w, orders)
}

func (a *App) RefreshArchiveOrder(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")

	objectID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Failed to create objectId")
	}

	updated, err := requests.RefreshOrder(a.client, objectID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if updated == 0 {
		sendError(w, http.StatusBadRequest, "nothing to update")
		return
	}

	sendOk(w)
}
