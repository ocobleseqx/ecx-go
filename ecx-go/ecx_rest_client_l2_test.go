package ecx

import (
	"context"
	"ecx-go/v3/internal/api"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	baseURL = "http://localhost:8888"
)

func TestGetL2Connection(t *testing.T) {
	//Given
	respBody := api.L2ConnectionResponse{}
	if err := readJSONData("./test-fixtures/ecx_l2connection_get_resp.json", &respBody); err != nil {
		assert.Failf(t, "Cannont read test response due to %s", err.Error())
	}
	connID := "connId"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ecx/v3/l2/connections/%s", baseURL, connID),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)

	//When
	ecxClient := NewClient(context.Background(), baseURL, testHc)
	conn, err := ecxClient.GetL2Connection(connID)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, conn, "Client should return a response")
	verifyL2Connection(t, *conn, respBody)
}

func TestCreateL2Connection(t *testing.T) {
	//Given
	respBody := api.CreateL2ConnectionResponse{}
	if err := readJSONData("./test-fixtures/ecx_l2connection_post_resp.json", &respBody); err != nil {
		assert.Failf(t, "Cannont read test response due to %s", err.Error())
	}
	reqBody := api.L2ConnectionRequest{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ecx/v3/l2/connections", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()
	newConnection := L2Connection{
		Name:                "name",
		ProfileUUID:         "profileUUID",
		Speed:               666,
		SpeedUnit:           "MB",
		Notifications:       []string{"janek@equinix.com", "marek@equinix.com"},
		PurchaseOrderNumber: "orderNumber",
		PortUUID:            "primaryPortUUID",
		VlanSTag:            100,
		VlanCTag:            101,
		ZSidePortUUID:       "primaryZSidePortUUID",
		ZSideVlanSTag:       200,
		ZSideVlanCTag:       201,
		SellerRegion:        "EMEA",
		SellerMetroCode:     "AM",
		AuthorizationKey:    "authorizationKey"}

	//When
	ecxClient := NewClient(context.Background(), baseURL, testHc)
	conn, err := ecxClient.CreateL2Connection(newConnection)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, conn, "Client should return a response")
	verifyL2ConnectionRequest(t, *conn, reqBody)
	assert.Equal(t, conn.UUID, respBody.PrimaryConnectionID, "UUID matches")
}

func TestCreateRedundantL2Connection(t *testing.T) {
	//Given
	respBody := api.CreateL2ConnectionResponse{}
	if err := readJSONData("./test-fixtures/ecx_l2connection_post_resp.json", &respBody); err != nil {
		assert.Failf(t, "Cannont read test response due to %s", err.Error())
	}
	reqBody := api.L2ConnectionRequest{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ecx/v3/l2/connections", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()
	newPriConn := L2Connection{
		Name:                "name",
		ProfileUUID:         "profileUUID",
		Speed:               666,
		SpeedUnit:           "MB",
		Notifications:       []string{"janek@equinix.com", "marek@equinix.com"},
		PurchaseOrderNumber: "orderNumber",
		PortUUID:            "primaryPortUUID",
		VlanSTag:            100,
		VlanCTag:            101,
		ZSidePortUUID:       "primaryZSidePortUUID",
		ZSideVlanSTag:       200,
		ZSideVlanCTag:       201,
		SellerRegion:        "EMEA",
		SellerMetroCode:     "AM",
		AuthorizationKey:    "authorizationKey"}
	newSecConn := L2Connection{
		Name:          "secName",
		PortUUID:      "secondaryPortUUID",
		VlanSTag:      690,
		VlanCTag:      691,
		ZSidePortUUID: "secondaryZSidePortUUID",
		ZSideVlanSTag: 717,
		ZSideVlanCTag: 718}

	//When
	ecxClient := NewClient(context.Background(), baseURL, testHc)
	conn, err := ecxClient.CreateL2RedundantConnection(newPriConn, newSecConn)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, conn, "Client should return a response")
	verifyRedundantL2ConnectionRequest(t, newPriConn, newSecConn, reqBody)
	assert.Equal(t, conn.UUID, respBody.PrimaryConnectionID, "UUID matches")
	assert.Equal(t, conn.RedundantUUID, respBody.SecondaryConnectionID, "RedundantUUID matches")
}

func TestDeleteL2Connection(t *testing.T) {
	//Given
	respBody := api.DeleteL2ConnectionResponse{}
	if err := readJSONData("./test-fixtures/ecx_l2connection_delete_resp.json", &respBody); err != nil {
		assert.Failf(t, "Cannont read test response due to %s", err.Error())
	}
	connID := "connId"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ecx/v3/l2/connections/%s", baseURL, connID),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		})
	defer httpmock.DeactivateAndReset()

	//When
	ecxClient := NewClient(context.Background(), baseURL, testHc)
	err := ecxClient.DeleteL2Connection(connID)

	//Then
	assert.Nil(t, err, "Client should not return an error")
}

func verifyL2Connection(t *testing.T, conn L2Connection, resp api.L2ConnectionResponse) {
	assert.Equal(t, resp.UUID, conn.UUID, "UUID matches")
	assert.Equal(t, resp.Name, conn.Name, "Name matches")
	assert.Equal(t, resp.SellerServiceUUID, conn.ProfileUUID, "Name matches")
	assert.Equal(t, resp.Speed, conn.Speed, "Speed matches")
	assert.Equal(t, resp.SpeedUnit, conn.SpeedUnit, "SpeedUnit matches")
	assert.Equal(t, resp.Status, conn.Status, "Status matches")
	assert.ElementsMatch(t, resp.Notifications, conn.Notifications, "Notifications match")
	assert.Equal(t, resp.PurchaseOrderNumber, conn.PurchaseOrderNumber, "PurchaseOrderNumber match")
	assert.Equal(t, resp.PortUUID, conn.PortUUID, "PrimaryPortUUID matches")
	assert.Equal(t, resp.VlanSTag, conn.VlanSTag, "PrimaryVlanSTag matches")
	assert.Equal(t, resp.VlanCTag, conn.VlanCTag, "PrimaryVlanCTag matches")
	assert.Equal(t, resp.ZSidePortUUID, conn.ZSidePortUUID, "PrimaryZSidePortUUID matches")
	assert.Equal(t, resp.ZSideVlanSTag, conn.ZSideVlanSTag, "PrimaryZSideVlanSTag matches")
	assert.Equal(t, resp.ZSideVlanCTag, conn.ZSideVlanCTag, "PrimaryZSideVlanCTag matches")
	assert.Equal(t, resp.SellerMetroCode, conn.SellerMetroCode, "SellerMetroCode matches")
	assert.Equal(t, resp.AuthorizationKey, conn.AuthorizationKey, "AuthorizationKey matches")
	assert.Equal(t, resp.RedundantUUID, conn.RedundantUUID, "RedundantUUID key matches")
}

func verifyL2ConnectionRequest(t *testing.T, conn L2Connection, req api.L2ConnectionRequest) {
	assert.Equal(t, conn.Name, req.PrimaryName, "Name matches")
	assert.Equal(t, conn.ProfileUUID, req.ProfileUUID, "ProfileUUID matches")
	assert.Equal(t, conn.Speed, req.Speed, "Speed matches")
	assert.Equal(t, conn.SpeedUnit, req.SpeedUnit, "SpeedUnit matches")
	assert.ElementsMatch(t, conn.Notifications, req.Notifications, "Notifications match")
	assert.Equal(t, conn.PurchaseOrderNumber, req.PurchaseOrderNumber, "PurchaseOrderNumber matches")
	assert.Equal(t, conn.PortUUID, req.PrimaryPortUUID, "PrimaryPortUUID matches")
	assert.Equal(t, conn.VlanSTag, req.PrimaryVlanSTag, "PrimaryVlanSTag matches")
	assert.Equal(t, conn.VlanCTag, req.PrimaryVlanCTag, "PrimaryVlanCTag matches")
	assert.Equal(t, conn.ZSidePortUUID, req.PrimaryZSidePortUUID, "PrimaryZSidePortUUID matches")
	assert.Equal(t, conn.ZSideVlanSTag, req.PrimaryZSideVlanSTag, "PrimaryZSideVlanSTag matches")
	assert.Equal(t, conn.ZSideVlanCTag, req.PrimaryZSideVlanCTag, "PrimaryZSideVlanCTag matches")
	assert.Equal(t, conn.SellerRegion, req.SellerRegion, "SellerRegion matches")
	assert.Equal(t, conn.SellerMetroCode, req.SellerMetroCode, "SellerMetroCode matches")
	assert.Equal(t, conn.AuthorizationKey, req.AuthorizationKey, "Authorization key matches")
}

func verifyRedundantL2ConnectionRequest(t *testing.T, primary L2Connection, secondary L2Connection, req api.L2ConnectionRequest) {
	verifyL2ConnectionRequest(t, primary, req)
	assert.Equal(t, secondary.Name, req.SecondaryName, "SecondaryName matches")
	assert.Equal(t, secondary.PortUUID, req.SecondaryPortUUID, "SecondaryPortUUID matches")
	assert.Equal(t, secondary.VlanSTag, req.SecondaryVlanSTag, "SecondaryVlanSTag matches")
	assert.Equal(t, secondary.VlanCTag, req.SecondaryVlanCTag, "SecondaryVlanCTag matches")
	assert.Equal(t, secondary.ZSidePortUUID, req.SecondaryZSidePortUUID, "SecondaryZSidePortUUID matches")
	assert.Equal(t, secondary.ZSideVlanSTag, req.SecondaryZSideVlanSTag, "SecondaryZSideVlanSTag matches")
	assert.Equal(t, secondary.ZSideVlanCTag, req.SecondaryZSideVlanCTag, "SecondaryZSideVlanCTag matches")
}
