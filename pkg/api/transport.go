package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/anabiozz/lapkins-api/pkg/auth"
	"github.com/anabiozz/lapkins-api/pkg/cookies"
	"github.com/anabiozz/lapkins-api/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

// ***************** GET PRODUCT ********************

type getProductRequest struct {
	SKU string
}

type getProductResponse struct {
	Product *model.Product `json:"product"`
	Err     error          `json:"err"`
}

func decodeGetProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getProductRequest{}
	req.SKU = r.URL.Query().Get("sku")
	return req, nil
}

func encodeGetProductResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getProductResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Product)
}

func makeGetProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getProductRequest)
		product, err := s.GetProduct(ctx, req.SKU)
		return getProductResponse{Err: err, Product: product}, nil
	}
}

// ***************** UPDATE PRODUCT ********************

type updateProductRequest struct {
	Product *model.Product
}

type updateProductResponse struct {
	Err error `json:"err"`
}

func decodeUpdateProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := updateProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeUpdateProductResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(updateProductResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func makeUpdateProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateProductRequest)
		err = s.UpdateProduct(ctx, req.Product)
		return updateProductResponse{Err: err}, nil
	}
}

// **************** GET PRODUCTS *******************

type getProductsRequest struct {
}

type getProductsResponse struct {
	Products []*model.Product `json:"products"`
	Err      error            `json:"err"`
}

func makeGetProductsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		products, err := s.GetProducts(ctx)
		return getProductsResponse{Err: err, Products: products}, nil
	}
}

func decodeGetProductsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getProductsRequest{}
	return req, nil
}

func encodeGetProductsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getProductsResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Products)
}

// **************************** GET CATEGORIES *************************

type getCategoriesRequest struct{}

type getCategoriesResponse struct {
	Categories []*model.Category `json:"categories"`
	Err        error             `json:"err"`
}

func makeGetCategoriesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		categories, err := s.GetCategories(ctx)
		return getCategoriesResponse{Err: err, Categories: categories}, nil
	}
}

func decodeGetCategoriesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getCategoriesRequest{}
	return req, nil
}

func encodeGetCategoriesResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getCategoriesResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Categories)
}

// **************************** GET ATTRIBUTES *************************

type addAttributeRequest struct {
	SKU       string
	Attribute *model.NameValue
}

type addAttributeResponse struct {
	Err error `json:"err"`
}

func makeAddAttributeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addAttributeRequest)
		err = s.AddAttribute(ctx, req.SKU, req.Attribute)
		return addAttributeResponse{Err: err}, nil
	}
}

func decodeAddAttributeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := addAttributeRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeAddAttributeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(addAttributeResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

// **************************** REMOVE ATTRIBUTES *************************

type removeAttributeRequest struct {
	SKU       string
	Attribute string
}

type removeAttributeResponse struct {
	Err error `json:"err"`
}

func makeRemoveAttributeEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeAttributeRequest)
		err = s.RemoveAttribute(ctx, req.SKU, req.Attribute)
		return removeAttributeResponse{Err: err}, nil
	}
}

func decodeRemoveAttributeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := removeAttributeRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeRemoveAttributeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(removeAttributeResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

// **************************** ADD CATEGORY *************************

type addCategoryRequest struct {
	SKU      string
	Category *model.Category
}

type addCategoryResponse struct {
	Err error `json:"err"`
}

func makeAddCategoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addCategoryRequest)
		err = s.AddCategory(ctx, req.SKU, req.Category)
		return addCategoryResponse{Err: err}, nil
	}
}

func decodeAddCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := addCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeAddCategoryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(addCategoryResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

// **************************** REMOVE CATEGORY *************************

type removeCategoryRequest struct {
	SKU      string
	Category *model.Category
}

type removeCategoryResponse struct {
	Err error `json:"err"`
}

func makeRemoveCategoryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeCategoryRequest)
		err = s.RemoveCategory(ctx, req.SKU, req.Category)
		return removeCategoryResponse{Err: err}, nil
	}
}

func decodeRemoveCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := removeCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeRemoveCategoryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(removeCategoryResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

type addProductRequest struct {
	SKU            string `json:"sku"`
	UserID         string `json:"user_id"`
	IsLoggedIn     bool   `json:"is_logged_in"`
	IsTmpUserIDSet bool   `json:"is_tmp_user_id_set"`
}

type addProductResponse struct {
	Err                error  `json:"err"`
	UserID             string `json:"user_id"`
	SetTmpUserIDCookie bool   `json:"set_tmp_user_id_cookie"`
}

func makeAddProductEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addProductRequest)
		setTmpUserIDCookie, userID, err := s.AddProduct(ctx, req.SKU, req.UserID, req.IsLoggedIn, req.IsTmpUserIDSet)
		return addProductResponse{Err: err, SetTmpUserIDCookie: setTmpUserIDCookie, UserID: userID}, nil
	}
}

type getHeaderCartInfoRequest struct {
	UserID string `json:"user_id"`
	Err    error  `json:"err"`
}

type getHeaderCartInfoResponse struct {
	Info *model.HeaderCartInfo `json:"info"`
	Err  error                 `json:"err"`
}

func makeGetHeaderCartInfo(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getHeaderCartInfoRequest)
		info, err := s.GetHeaderCartInfo(ctx, req.UserID)
		return getHeaderCartInfoResponse{Err: err, Info: info}, nil
	}
}

type getCartRequest struct {
	UserID string `json:"user_id"`
	Err    error  `json:"err"`
}

type getCartResponse struct {
	Cart []*model.CartProduct `json:"cart"`
	Err  error                `json:"err"`
}

func makeGetCart(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCartRequest)
		cart, err := s.LoadCart(ctx, req.UserID)
		return getCartResponse{Err: err, Cart: cart}, nil
	}
}

type increaseProductQtyRequest struct {
	UserID string `json:"user_id"`
	SKU    string `json:"sku"`
}

type increaseProductQtyResponse struct {
	Err error `json:"err"`
}

func makeIncreaseProductQty(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(increaseProductQtyRequest)
		err = s.IncreaseProductQuantity(ctx, req.UserID, req.SKU)
		return increaseProductQtyResponse{Err: err}, nil
	}
}

type decreaseProductQtyRequest struct {
	UserID string `json:"user_id"`
	SKU    string `json:"sku"`
}

type decreaseProductQtyResponse struct {
	Err error `json:"err"`
}

func makeDecreaseProductQty(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(decreaseProductQtyRequest)
		err = s.DecreaseProductQuantity(ctx, req.UserID, req.SKU)
		return decreaseProductQtyResponse{Err: err}, nil
	}
}

type removeProductRequest struct {
	UserID string `json:"user_id"`
	SKU    string `json:"sku"`
}

type removeProductResponse struct {
	Err error `json:"err"`
}

func makeRemoveProduct(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeProductRequest)
		err = s.RemoveProduct(ctx, req.UserID, req.SKU)
		return removeProductResponse{Err: err}, nil
	}
}

func decodeAddProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := addProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	userID, isLoggedIn, err := auth.GetUserID(r)
	if err != nil {
		return nil, err
	}
	req.UserID = userID
	if userID != "" && isLoggedIn {
		req.IsLoggedIn = true
	}
	tmpUserID, err := cookies.GetCookieValue(r, "tmp-user-id")
	if err != nil {
		if err != http.ErrNoCookie {
			return nil, err
		}
	}
	if tmpUserID != "" {
		req.IsTmpUserIDSet = true
	}
	return req, nil
}

func encodeAddProductResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(addProductResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}

	if res.SetTmpUserIDCookie {
		http.SetCookie(w, &http.Cookie{
			Path:    "/",
			Name:    "tmp-user-id",
			Value:   res.UserID,
			Expires: time.Now().Add(168 * time.Hour),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

func decodeGetHeaderCartInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getHeaderCartInfoRequest{}
	userID, _, err := auth.GetUserID(r)
	if err != nil {
		return nil, errBadRequest("%s", err)
	}
	req.UserID = userID
	return req, nil
}

func encodeGetHeaderCartInfoResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getHeaderCartInfoResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Info)
}

func decodeGetCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := getCartRequest{}
	userID, _, err := auth.GetUserID(r)
	if err != nil {
		return nil, errBadRequest("%s", err)
	}
	req.UserID = userID
	return req, nil
}

func encodeGetCartResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getCartResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Cart)
}

func decodeIncreaseProductQtyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := increaseProductQtyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	userID, _, err := auth.GetUserID(r)
	if err != nil {
		return nil, errBadRequest("%s", err)
	}
	req.UserID = userID
	return req, nil
}

func encodeIncreaseProductQtyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(increaseProductQtyResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

func decodeDecreaseProductQtyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := decreaseProductQtyRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	userID, _, err := auth.GetUserID(r)
	if err != nil {
		return nil, errBadRequest("%s", err)
	}
	req.UserID = userID
	return req, nil
}

func encodeDecreaseProductQtyResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(decreaseProductQtyResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

func decodeRemoveProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := removeProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	userID, _, err := auth.GetUserID(r)
	if err != nil {
		return nil, errBadRequest("%s", err)
	}
	req.UserID = userID
	return req, nil
}

func encodeRemoveProductResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(removeProductResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}

type addOrderRequest struct {
	Order *model.Order
}

type addOrderResponse struct {
	Err error `json:"err"`
}

func makeAddOrderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(addOrderRequest)
		err = s.AddOrder(ctx, req.Order)
		return addOrderResponse{Err: err}, nil
	}
}

func decodeAddOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := addOrderRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeAddOrderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(addOrderResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(true)
}


type registerRequest struct {
	Input *model.UserInput
}

type registerResponse struct {
	User *model.UserOutput
	Err  error
}

func makeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		user, err := s.Register(ctx, req.Input)
		return registerResponse{User: user, Err: err}, nil
	}
}

type loginRequest struct {
	Input     *model.UserInput
	TmpUserID string `json:"tmp_user_id"`
}

type loginResponse struct {
	User                 *model.UserOutput
	UnsetTmpUserIDCookie bool `json:"unset_tmp_user_id_cookie"`
	Err                  error
}

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		user, unsetTmpUserIDCookie, err := s.Login(ctx, req.Input, req.TmpUserID)
		return loginResponse{User: user, UnsetTmpUserIDCookie: unsetTmpUserIDCookie, Err: err}, nil
	}
}

type refreshTokenRequest struct {
	Token string
}

type refreshTokenResponse struct {
	User *model.UserOutput
	Err  error
}

func makeRefreshTokenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(refreshTokenRequest)
		user, err := s.RefreshToken(ctx, req.Token)
		return refreshTokenResponse{User: user, Err: err}, nil
	}
}

type getUsersResponse struct {
	Users []*model.User
	Err   error
}

func makeGetUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		users, err := s.GetUsers(ctx)
		return getUsersResponse{Users: users, Err: err}, nil
	}
}



func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := registerRequest{Input: &model.UserInput{}}
	if err := json.NewDecoder(r.Body).Decode(req.Input); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	return req, nil
}

func encodeRegisterResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(registerResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	http.SetCookie(w, &http.Cookie{
		Path:    "/",
		Name:    "token",
		Value:   res.User.Token,
		Expires: res.User.ExpirationTime,
	})
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.User)
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := loginRequest{Input: &model.UserInput{}}
	if err := json.NewDecoder(r.Body).Decode(req.Input); err != nil {
		return nil, errBadRequest("failed to decode JSON request: %v", err)
	}
	tmpUserID, err := cookies.GetCookieValue(r, "tmp-user-id")
	if err != nil {
		if err != http.ErrNoCookie {
			return nil, err
		}
	}
	req.TmpUserID = tmpUserID
	return req, nil
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(loginResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Path:    "/",
		Name:    "token",
		Value:   res.User.Token,
		Expires: res.User.ExpirationTime,
	})

	if res.UnsetTmpUserIDCookie {
		http.SetCookie(w, &http.Cookie{
			Path:    "/",
			Name:    "tmp-user-id",
			Value:   "",
			Expires: time.Unix(0, 0),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.User)
}

func decodeRefreshTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := refreshTokenRequest{}
	token, err := auth.GetToken(r)
	if err != nil {
		return nil, err
	}
	req.Token = token
	return req, nil
}

func encodeRefreshTokenResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(refreshTokenResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.User)
}

func decodeGetUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	header := r.Header.Get("Authorization")
	_, err := auth.Check(header)
	if err != nil {
		return nil, errUnauthorized("%s", err)
	}
	return r, nil
}

func encodeGetUsersResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(getUsersResponse)
	if res.Err != nil {
		encodeError(ctx, res.Err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res.Users)
}


// ****************** Errors *********************

func decodeError(res *http.Response) error {
	defaultErr := &serviceError{
		code: res.StatusCode,
	}
	tmp := &serviceError{}
	err := json.NewDecoder(res.Body).Decode(tmp)
	if err != nil {
		return defaultErr
	}
	defaultErr.message = tmp.message
	return defaultErr
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	defaultErr := &serviceError{code: http.StatusInternalServerError, message: "internal error"}
	if err, ok := err.(*serviceError); ok {
		defaultErr.code = err.code
		defaultErr.message = err.message
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(defaultErr.code)
	_ = json.NewEncoder(w).Encode(defaultErr)
}
