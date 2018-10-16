// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package inspect

// [START dlp_inspect_file]
import (
	"context"
	"io/ioutil"
	"fmt"
	"errors"
	"log"

	"cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func inspectFile(projectID, filepath, fileType string) error {
	// projectID := "my-project-id";
	// filePath := "path/to/image.png";
	// fileType := "IMAGE"

	// Initialize client
	client, err := dlp.NewClient(context.Background())
	if err != nil {
		return err
	}
	defer client.Close() // Closing the client safely stops any background threads or connections.

	// Prepare the request
	// Set the item for the request.
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	var itemType dlppb.ByteContentItem_BytesType
	switch fileType {
	case "IMAGE":
		itemType = dlppb.ByteContentItem_IMAGE
	case "TEXT_UTF8":
		itemType = dlppb.ByteContentItem_TEXT_UTF8
	default:
		return errors.New(fmt.Sprintf("invalid ByteType for ByteContentItem: '%s'", fileType))
	}
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_ByteItem{
			ByteItem: &dlppb.ByteContentItem{
				Type: itemType,
				Data: data,
			},
		},
	}
	// Set the required InfoTypes for the inspection config.
	var infoTypes []*dlppb.InfoType
	for _, it := range []string{"PHONE_NUMBER", "EMAIL_ADDRESS", "CREDIT_CARD_NUMBER"} {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
	}
	// Set the inspection configuration for the request.
	config := &dlppb.InspectConfig{
		InfoTypes:    infoTypes,
		IncludeQuote: true,
	}

	// Create and send the request.
	req := &dlppb.InspectContentRequest{
		Parent:        "projects/" + projectID,
		Item:          item,
		InspectConfig: config,
	}
	resp, err := client.InspectContent(context.Background(), req)
	if err != nil {
		return err
	}

	// Process the results.
	result := resp.Result
	log.Printf("Findings: %d\n", len(result.Findings))
	for _, f := range result.Findings {
		log.Printf("\tQoute: %s\n", f.Quote)
		log.Printf("\tInfo type: %s\n", f.InfoType.Name)
		log.Printf("\tLikelihood: %s\n", f.Likelihood)
	}
	return nil
}

// [END dlp_inspect_file]
