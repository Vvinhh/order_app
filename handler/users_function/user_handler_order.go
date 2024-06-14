// Inside handler package (user_handler_order.go)

package handler

import (
	"github.com/MinhSang97/order_app/pkg/id"
	"github.com/MinhSang97/order_app/usecases"
	"github.com/MinhSang97/order_app/usecases/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"

	"github.com/MinhSang97/order_app/usecases/req"
	"github.com/MinhSang97/order_app/usecases/res"
)

// UsersOrder godoc
// @Summary Users can order
// @Description Users can order
// @Tags usersOrder
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param order_date body string true "OrderDate"
// @Param total_price body float64 true "TotalPrice"
// @Param status body string true "Status"
// @Param address body string true "Address"
// @Param payment_method body string true "PaymentMethod"
// @Param payment_date body string true "PaymentDate"
// @Param amount body float64 true "Amount"
// @Param discount_code_id body string true "DiscountCodeId"
// @Param item_id body []string true "ItemID"
// @Param quantity body []int true "Quantity"
// @Param price body []float64 true "Price"
// @Success 200 {object} res.Response
// @Failure 400 {object} res.Response
// @Failure 403 {object} res.Response
// @Failure 500 {object} res.Response
// @Router /v1/api/users/order/{user_id} [post]
func UsersOrder() func(*gin.Context) {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		if user_id == "" || user_id == ":user_id" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "không tìm thấy user_id",
			})
			return
		}

		var validate *validator.Validate
		validate = validator.New(validator.WithRequiredStructEnabled())
		reqBody := req.ReqOrder{}
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, res.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request body",
				Data:       nil,
			})
			return
		}

		if len(reqBody.ItemID) != len(reqBody.Quantity) || len(reqBody.ItemID) != len(reqBody.Price) {
			c.JSON(http.StatusBadRequest, res.Response{
				StatusCode: http.StatusBadRequest,
				Message:    "ItemID, Quantity, and Price must have the same number of elements",
				Data:       nil,
			})
			return
		}

		order := make(map[string]req.OrderItem)
		for i, item := range reqBody.ItemID {
			order[item] = req.OrderItem{
				Quantity: reqBody.Quantity[i],
				Price:    reqBody.Price[i],
			}
		}

		if err := validate.Struct(reqBody); err != nil {
			c.JSON(http.StatusForbidden, res.Response{
				StatusCode: http.StatusForbidden,
				Message:    err.Error(),
				Data:       nil,
			})
			return
		}
		// Create order_id
		order_id := id.OrderID()

		orderData := dto.OrderDto{
			OrderID:        order_id,
			OrderDate:      reqBody.OrderDate,
			TotalPrice:     reqBody.TotalPrice,
			Status:         reqBody.Status,
			Address:        reqBody.Address,
			PaymentMethod:  reqBody.PaymentMethod,
			PaymentDate:    reqBody.PaymentDate,
			Amount:         reqBody.Amount,
			DiscountCodeId: reqBody.DiscountCodeId,
			ItemID:         reqBody.ItemID,
			Quantity:       reqBody.Quantity,
			Price:          reqBody.Price,
		}

		// Save order to database
		data := orderData.ToPayload().ToModel()
		uc := usecases.NewUsersUseCase()

		orderDB, err := uc.AddOrderUsersOrder(c.Request.Context(), user_id, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, res.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
			})
			return
		}

		c.JSON(http.StatusOK, res.Response{
			StatusCode: http.StatusOK,
			Message:    "Order successfully",
			Data:       "OrderID: " + strconv.FormatInt(orderDB.OrderID, 10),
		})
	}
}