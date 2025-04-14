package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Drug struct {
	DrugID          string   `json:"drugId"`
	Name            string   `json:"name"`
	Manufacturer    string   `json:"manufacturer"`
	BatchNumber     string   `json:"batchNumber"`
	MfgDate         string   `json:"mfgDate"`
	ExpiryDate      string   `json:"expiryDate"`
	Composition     string   `json:"composition"`
	CurrentOwner    string   `json:"currentOwner"` // Cipla, Medlife, Apollo
	Status          string   `json:"status"`       // InProduction, InTransit, Delivered, Recalled
	History         []string `json:"history"`      // Format: "timestamp|event|from|to|details"
	IsRecalled      bool     `json:"isRecalled"`
	InspectionNotes []string `json:"inspectionNotes"`
}

type SmartContract struct {
	contractapi.Contract
}

func timestamp() string {
	return time.Now().GoString()
}

// ============== MANUFACTURER FUNCTIONS ==============
func (s *SmartContract) RegisterDrug(ctx contractapi.TransactionContextInterface,
	drugID string, name string, batchNumber string, mfgDate string, expiryDate string, composition string) error {

	// TODO: Verify caller is CiplaMSP
	clientMsp, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMsp == "CDSCO" {
		return fmt.Errorf("Drug regulator cannot register drug")
	}
	// TODO: Check if drug exists
	drugBytes, err := ctx.GetStub().GetState(drugID)
	if err != nil {
		return err
	}
	if drugBytes != nil {
		return fmt.Errorf("drug with drug id already exists")
	}
	// TODO: Initialize drug object with all fields
	drug := &Drug{
		DrugID:          drugID,
		Name:            name,
		Manufacturer:    clientMsp,
		BatchNumber:     batchNumber,
		MfgDate:         mfgDate,
		ExpiryDate:      expiryDate,
		Composition:     composition,
		CurrentOwner:    clientMsp,
		Status:          "InProduction",
		History:         []string{fmt.Sprintf("%s|%s|%s|%s|%s", timestamp(), "Manufactured", clientMsp, clientMsp, "manufactured the drug")},
		IsRecalled:      false,
		InspectionNotes: []string{},
	}

	drugBytes, err = json.Marshal(drug)
	if err != nil {
		return err
	}

	// TODO: Save to ledger
	err = ctx.GetStub().PutState(drugID, drugBytes)
	if err != nil {
		return err
	}

	// TODO: Add creation event to history
	err = ctx.GetStub().SetEvent("CreateDrug", drugBytes)
	if err != nil {
		return err
	}

	return nil
}

// ============== DISTRIBUTION FUNCTIONS ==============
func (s *SmartContract) ShipDrug(ctx contractapi.TransactionContextInterface, drugID string, to string) error {
	// TODO: Verify current owner is caller
	clientMsp, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}

	drugBytes, err := ctx.GetStub().GetState(drugID)
	if err != nil {
		return err
	}
	var drug = &Drug{}
	err = json.Unmarshal(drugBytes, drug)
	if err != nil {
		return err
	}

	if drug.CurrentOwner != clientMsp {
		return fmt.Errorf("owner can only ship drugs")
	}
	// TODO: Update CurrentOwner and Status
	drug.CurrentOwner = to
	drug.Status = "shipped"
	// TODO: Add shipment record to history
	drug.History = append(drug.History, fmt.Sprintf("%s|%s|%s|%s|%s", timestamp(), "Shipped", clientMsp, to, "shipping the drug"))

	drugBytes, err = json.Marshal(drug)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(drugID, drugBytes)
	if err != nil {
		return err
	}
	// TODO: Emit shipment event
	err = ctx.GetStub().SetEvent("shipDrug", drugBytes)
	if err != nil {
		return err
	}
	return nil
}

// ============== REGULATOR FUNCTIONS ==============
func (s *SmartContract) RecallDrug(ctx contractapi.TransactionContextInterface, drugID string, reason string) error {
	// TODO: Verify caller is CDSCOMSP
	clientMsp, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientMsp != "CDSCO" {
		return fmt.Errorf("Drug regulator can only recall drugs")
	}
	// TODO: Set IsRecalled=true
	drugBytes, err := ctx.GetStub().GetState(drugID)
	if err != nil {
		return err
	}

	var drug = &Drug{}
	err = json.Unmarshal(drugBytes, drug)
	if err != nil {
		return err
	}

	drug.IsRecalled = true
	drug.History = append(drug.History, fmt.Sprintf("%s|%s|%s|%s|%s", timestamp(), "Recalled", clientMsp, drug.Manufacturer, "recalling the drug"))
	drug.InspectionNotes = append(drug.InspectionNotes, reason)
	// TODO: Add recall note to InspectionNotes
	drugBytes, err = json.Marshal(drug)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(drugID, drugBytes)
	if err != nil {
		return err
	}

	return nil
}

// ============== COMMON FUNCTIONS ==============
func (s *SmartContract) TrackDrug(ctx contractapi.TransactionContextInterface, drugID string) (string, error) {
	// TODO: Return full drug history as JSON
	drugBytes, err := ctx.GetStub().GetState(drugID)
	if err != nil {
		return "", err
	}

	var drug = &Drug{}
	err = json.Unmarshal(drugBytes, drug)
	if err != nil {
		return "", err
	}

	var buf *bytes.Buffer
	err = json.Indent(buf, drugBytes, "", " ")
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
