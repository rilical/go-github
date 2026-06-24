package output

import (
	"bytes"
	_ "embed"
	"fmt"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// FeedSchema is the embedded JSON Schema for the change feed, exported so
// external validators and CLI consumers can read it without importing the file.
//
//go:embed feed.schema.json
var FeedSchema []byte

const schemaURL = "mem://github.com/github/openapi-change-feed/feed.schema.json"

var (
	schemaOnce   sync.Once
	batchSchema  *jsonschema.Schema
	recordSchema *jsonschema.Schema
	schemaErr    error
)

func compileSchemas() {
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(FeedSchema))
	if err != nil {
		schemaErr = fmt.Errorf("parse embedded feed schema: %w", err)
		return
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource(schemaURL, doc); err != nil {
		schemaErr = fmt.Errorf("register feed schema: %w", err)
		return
	}
	if batchSchema, err = c.Compile(schemaURL); err != nil {
		schemaErr = fmt.Errorf("compile feed schema: %w", err)
		return
	}
	if recordSchema, err = c.Compile(schemaURL + "#/$defs/record"); err != nil {
		schemaErr = fmt.Errorf("compile record schema: %w", err)
		return
	}
}

func schemas() (batch, record *jsonschema.Schema, err error) {
	schemaOnce.Do(compileSchemas)
	return batchSchema, recordSchema, schemaErr
}

// ValidateFeed checks that a serialized ChangeBatch (changes.json) conforms to
// the embedded JSON Schema, including the severity enum and the required summary
// counts, not merely that the top-level keys are present.
func ValidateFeed(b []byte) error {
	batch, _, err := schemas()
	if err != nil {
		return err
	}
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("feed not valid JSON: %w", err)
	}
	if err := batch.Validate(inst); err != nil {
		return fmt.Errorf("feed failed schema validation: %w", err)
	}
	return nil
}

// ValidateFeedLine checks a single feed.jsonl record against the record schema.
func ValidateFeedLine(b []byte) error {
	_, record, err := schemas()
	if err != nil {
		return err
	}
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("feed line not valid JSON: %w", err)
	}
	if err := record.Validate(inst); err != nil {
		return fmt.Errorf("feed line failed schema validation: %w", err)
	}
	return nil
}
