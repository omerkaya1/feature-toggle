package rest

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	createFeatureToggle struct {
		Inverted      bool      `json:"inverted"`
		Active        bool      `json:"active"`
		DisplayName   string    `json:"displayName,omitempty"`
		TechnicalName string    `json:"technicalName"`
		Description   string    `json:"description,omitempty"`
		ExpiresOn     time.Time `json:"expiresOn,omitempty"`
		CustomerIDs   []string  `json:"customerIDs"`
	}
	updateFeatureToggle struct {
		Active      string    `json:"active"`
		DisplayName string    `json:"displayName"`
		Description string    `json:"description"`
		ExpiresOn   time.Time `json:"expiresOn"`
	}
)

func (ft createFeatureToggle) Validate() error {
	if ft.TechnicalName == "" {
		return errors.New("empty technical name")
	}
	if len(ft.CustomerIDs) == 0 {
		return errors.New("empty customer id list")
	}
	return nil
}

func (ft updateFeatureToggle) Validate() error {
	if ft.DisplayName == "" {
		return errors.New("empty display name")
	}
	if ft.Description == "" {
		return errors.New("empty description")
	}
	if ft.Active == "" {
		return errors.New("empty status")
	}
	if ft.ExpiresOn.IsZero() {
		return errors.New("empty expiration date")
	}
	return nil
}

func (s *Server) getFeatures(c echo.Context) error {
	if c.QueryParam("customer") != "" && c.QueryParam("active") != "" {
		return s.getAllCustomerFeaturesByStatus(c)
	}
	if c.QueryParam("customer") != "" {
		return s.getAllCustomerFeatures(c)
	}
	return s.getAllFeatures(c)
}

func (s *Server) getAllFeatures(c echo.Context) error {
	result, err := s.store.GetFeatures(c.Request().Context())
	if err != nil {
		s.log.Errorf("failure to get user features by status: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get features")
	}
	return c.JSON(http.StatusOK, result)
}

func (s *Server) getAllCustomerFeatures(c echo.Context) error {
	result, err := s.store.GetUserFeatures(c.Request().Context(), c.QueryParam("customer"))
	if err != nil {
		s.log.Errorf("failure to get user features by status: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get features")
	}
	return c.JSON(http.StatusOK, result)
}

func (s *Server) getAllCustomerFeaturesByStatus(c echo.Context) error {
	val, err := strconv.ParseBool(c.QueryParam("active"))
	if err != nil {
		s.log.Errorf("failure to parse a query parameter: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameter")
	}
	result, err := s.store.GetUserFeaturesByStatus(c.Request().Context(), c.QueryParam("customer"), val)
	if err != nil {
		s.log.Errorf("failure to get user features by status: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get features")
	}
	return c.JSON(http.StatusOK, result)
}

func (s *Server) createFeature(c echo.Context) error {
	var ft createFeatureToggle
	if err := c.Bind(&ft); err != nil {
		s.log.Errorf("failure to unmarshal the payload: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "malformed payload")
	}
	if err := ft.Validate(); err != nil {
		s.log.Errorf("failure to validate the payload data: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload")
	}
	id, err := s.store.CreateFeature(c.Request().Context(), ft.Inverted, ft.Active, ft.DisplayName, ft.TechnicalName,
		ft.Description, ft.CustomerIDs, ft.ExpiresOn)
	if err != nil {
		s.log.Errorf("failure to create a feature toggle: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save a new feature toggle")
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"id": id})
}

func (s *Server) updateFeature(c echo.Context) error {
	var uft updateFeatureToggle
	if err := c.Bind(&uft); err != nil {
		s.log.Errorf("failure to unmarshal the payload: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "malformed payload")
	}
	if err := uft.Validate(); err != nil {
		s.log.Errorf("failure to validate the payload data: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload")
	}
	val, err := strconv.ParseBool(uft.Active)
	if err != nil {
		s.log.Errorf("failure to parse a query parameter: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameter")
	}
	if err = s.store.UpdateFeature(c.Request().Context(), c.Param("name"), uft.Description, uft.DisplayName,
		uft.ExpiresOn, val); err != nil {
		s.log.Errorf("failure to update a feature: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update a feature")
	}
	return c.JSON(http.StatusOK, "feature was updated")
}

func (s *Server) deleteFeature(c echo.Context) error {
	if err := s.store.DeleteFeature(c.Request().Context(), c.Param("name")); err != nil {
		s.log.Errorf("failure to delete a feature: %s", err)
		return echo.NewHTTPError(http.StatusBadRequest, "could not delete a feature")
	}
	return c.JSON(http.StatusOK, "feature was deleted")
}
