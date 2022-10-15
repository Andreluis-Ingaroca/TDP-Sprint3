/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a medicine
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a medicine
type Medicine struct {
	Name   	       string `json:"name"`
	Concentration  string `json:"concentration"`
	Form           string `json:"form"`
	Expiration     string `json:"expiration"`
	Quantity       string `json:"quantity"`
	Code           string `json:"code"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Medicine
}

// InitLedger adds a base set of medicines to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	medicines := []Medicine{
		Medicine{Name: "Amoxicilina", Concentration: "250mg/5ml", Form: "Jarabe", Expiration: "31/12/2023", Quantity: "100", Code: "1234567"},
		Medicine{Name: "Ibuprofeno", Concentration: "400mg", Form: "Tableta", Expiration: "31/12/2024", Quantity: "100", Code: "1234567"},
	}

	for i, medicine := range medicines {
		medicineAsBytes, _ := json.Marshal(medicine)
		err := ctx.GetStub().PutState("MEDICINE"+strconv.Itoa(i), medicineAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateMedicine adds a new medicine to the world state with given details
func (s *SmartContract) CreateMedicine(ctx contractapi.TransactionContextInterface, medicineNumber string, name string, concentration string, form string, expiration string, quantity string, code string) error {
	medicine := Medicine{
		Name: name,
		Concentration: concentration,
		Form: form,
		Expiration: expiration,
		Quantity: quantity,
		Code: code,
	}

	medicineAsBytes, _ := json.Marshal(medicine)

	return ctx.GetStub().PutState(medicineNumber, medicineAsBytes)
}

// QueryMedicine returns the medicine stored in the world state with given id
func (s *SmartContract) QueryMedicine(ctx contractapi.TransactionContextInterface, medicineNumber string) (*Medicine, error) {
	medicineAsBytes, err := ctx.GetStub().GetState(medicineNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if medicineAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", medicineNumber)
	}

	medicine := new(Medicine)
	_ = json.Unmarshal(medicineAsBytes, medicine)

	return medicine, nil
}

// QueryAllMedicines returns all medicines found in world state
func (s *SmartContract) QueryAllMedicines(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		medicine := new(Medicine)
		_ = json.Unmarshal(queryResponse.Value, medicine)

		queryResult := QueryResult{Key: queryResponse.Key, Record: medicine}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeCarOwner updates the owner field of medicine with given id in world state
func (s *SmartContract) ChangeMedicineQuantity(ctx contractapi.TransactionContextInterface, medicineNumber string, newQuantity string) error {
	medicine, err := s.QueryMedicine(ctx, medicineNumber)

	if err != nil {
		return err
	}

	medicine.Quantity = newQuantity
	medicineAsBytes, _ := json.Marshal(medicine)

	return ctx.GetStub().PutState(medicineNumber, medicineAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
