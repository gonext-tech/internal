package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gonext-tech/internal/models"
	"github.com/gonext-tech/internal/views/components"
	"github.com/gonext-tech/internal/views/shop_views"
	"github.com/labstack/echo/v4"
)

type ShopService interface {
	GetALL(limit, page int, orderBy, sortBy, project, status, searchTerm string) ([]models.Shop, models.Meta, error)
	Fetch(string) ([]models.Shop, error)
	GetID(id, name string) (models.Shop, error)
	Create(models.Shop) (models.Shop, error)
	Update(models.Shop) (models.Shop, error)
	Delete(models.Shop) (models.Shop, error)
}

type ShopHandler struct {
	ShopServices    ShopService
	ProjectServices ProjectService
	UploadServices  UploadService
}

func NewShopHandler(ss ShopService, ps ProjectService, us UploadService) *ShopHandler {
	return &ShopHandler{
		ShopServices:    ss,
		ProjectServices: ps,
		UploadServices:  us,
	}
}

func (sh *ShopHandler) Fetch(c echo.Context) error {
	//projectName := c.QueryParam("project_name")
	//memberships, _ := sh.ShopServices.Fetch(projectName)
	//return renderView(c, components.MembershipResult(memberships))
	return nil
}

func (sh *ShopHandler) ListPage(c echo.Context) error {
	isError = false

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20
	}
	orderBy := c.QueryParam("orderBy")
	if orderBy == "" {
		orderBy = "desc"
	}

	sortBy := c.QueryParam("sortBy")
	if sortBy == "" {
		sortBy = "id"
	}
	status := c.QueryParam("status")
	project := c.QueryParam("project_name")
	if project == "" {
		project = ""
	}

	searchTerm := c.QueryParam("searchTerm")
	response, meta, err := sh.ShopServices.GetALL(limit, page, orderBy, sortBy, project, status, searchTerm)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var params models.ParamResponse
	if searchTerm != "" {
		params.Search = searchTerm
	}
	if status != "" {
		params.Status = status
	}
	params.Page = page
	params.Limit = limit
	params.SortBy = sortBy
	params.OrderBy = orderBy
	titlePage := fmt.Sprintf(
		"Shops (%d)", meta.TotalCount)
	return renderView(c, shop_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		shop_views.List(titlePage, response, meta, params),
	))
}

func (sh *ShopHandler) SearchUser(c echo.Context) error {
	isError = false
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 5
	}
	orderBy := c.QueryParam("orderBy")
	if orderBy == "" {
		orderBy = "desc"
	}
	sortBy := c.QueryParam("sortBy")
	if sortBy == "" {
		sortBy = "id"
	}
	status := c.QueryParam("status")
	project := c.QueryParam("project_name")
	searchTerm := c.QueryParam("searchTerm")
	shops, _, err := sh.ShopServices.GetALL(limit, page, orderBy, sortBy, project, status, searchTerm)
	if err != nil {
		isError = false
		errorMsg = "can't fetch shops"
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
	}

	return renderView(c, components.ShopResult(shops))
}

func (ss *ShopHandler) ViewPage(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	shop, err := ss.ShopServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("can't fetch shop with id: %s", id)
	}
	if isError {
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/shop")
	}
	titlePage := fmt.Sprintf(
		"Shop | %s", shop.Name)
	return renderView(c, shop_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		shop_views.View(shop),
	))
}

func (sh *ShopHandler) CreatePage(c echo.Context) error {
	isError = false
	titlePage := "Membership | Create"
	projects, _, _ := sh.ProjectServices.GetALL(50, 1, "desc", "id", "", "")
	return renderView(c, shop_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		shop_views.Create(projects),
	))
}

func (mh *ShopHandler) CreateHandler(c echo.Context) error {
	var shop models.Shop
	if err := c.Bind(&shop); err != nil {
		setFlashmessages(c, "error", errorMsg)
		return mh.CreatePage(c)
	}
	_, err := mh.ShopServices.Create(shop)
	if err != nil {
		setFlashmessages(c, "error", "Can't create shop")
		return mh.CreatePage(c)
	}
	imageURLs := UploadImage(c, mh.UploadServices, shop.ProjectName, fmt.Sprintf("shop/%d", shop.ID))

	if len(imageURLs) > 0 {
		shop.Image = imageURLs[0]
		_, err = mh.ShopServices.Update(shop)
		if err != nil {
			setFlashmessages(c, "error", "Can't create customer")
			return mh.CreatePage(c)
		}
	}

	setFlashmessages(c, "success", "shop created successfully!!")

	return c.Redirect(http.StatusSeeOther, "/shop")
}

func (sh *ShopHandler) UpdatePage(c echo.Context) error {
	isError = false
	titlePage := "Shop | Update"
	id := c.Param("id")
	projectName := c.Param("name")
	shop, err := sh.ShopServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("shop with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
	}

	projects, _, _ := sh.ProjectServices.GetALL(50, 1, "desc", "id", "", "")
	return renderView(c, shop_views.Index(
		titlePage,
		c.Get(email_key).(string),
		fromProtected,
		isError,
		getFlashmessages(c, "error"),
		getFlashmessages(c, "success"),
		shop_views.Update(shop, projects),
	))
}

func (sh *ShopHandler) UpdateHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	shop, err := sh.ShopServices.GetID(id, projectName)
	log.Println("errrrr", err)
	if err != nil {
		errorMsg = fmt.Sprintf("shop with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}

	if err := c.Bind(&shop); err != nil {
		errorMsg = "cannot parse the shop body"
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}
	shop.Owner = models.User{}
	imageURLs := UploadImage(c, sh.UploadServices, "internal", fmt.Sprintf("shop/%d", shop.ID))
	if len(imageURLs) > 0 {
		shop.Image = imageURLs[0]
	}
	_, err = sh.ShopServices.Update(shop)
	if err != nil {
		errorMsg = fmt.Sprintf("shop with id %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return sh.UpdatePage(c)
	}
	setFlashmessages(c, "success", "shop updated successfully!!")

	return c.Redirect(http.StatusSeeOther, "/shop")
}

func (sh *ShopHandler) DeleteHandler(c echo.Context) error {
	isError = false
	id := c.Param("id")
	projectName := c.Param("name")
	shop, err := sh.ShopServices.GetID(id, projectName)
	if err != nil {
		errorMsg = fmt.Sprintf("shop with %s not found", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/shop")
	}
	_, err = sh.ShopServices.Delete(shop)
	if err != nil {
		errorMsg = fmt.Sprintf("couldnt delete shop with id %s", id)
		setFlashmessages(c, "error", errorMsg)
		return c.Redirect(http.StatusSeeOther, "/shop")
	}
	setFlashmessages(c, "success", "shop successfully deleted!!")
	return c.Redirect(http.StatusSeeOther, "/shop")
}
