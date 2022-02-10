// Code generated by smithy-go-codegen DO NOT EDIT.

package ssm

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Amazon Web Services Systems Manager calls this API operation when displaying all
// Application Manager OpsMetadata objects or blobs.
func (c *Client) ListOpsMetadata(ctx context.Context, params *ListOpsMetadataInput, optFns ...func(*Options)) (*ListOpsMetadataOutput, error) {
	if params == nil {
		params = &ListOpsMetadataInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ListOpsMetadata", params, optFns, c.addOperationListOpsMetadataMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListOpsMetadataOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListOpsMetadataInput struct {

	// One or more filters to limit the number of OpsMetadata objects returned by the
	// call.
	Filters []types.OpsMetadataFilter

	// The maximum number of items to return for this call. The call also returns a
	// token that you can specify in a subsequent call to get the next set of results.
	MaxResults int32

	// A token to start the list. Use this token to get the next set of results.
	NextToken *string

	noSmithyDocumentSerde
}

type ListOpsMetadataOutput struct {

	// The token for the next set of items to return. Use this token to get the next
	// set of results.
	NextToken *string

	// Returns a list of OpsMetadata objects.
	OpsMetadataList []types.OpsMetadata

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListOpsMetadataMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpListOpsMetadata{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpListOpsMetadata{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpListOpsMetadataValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListOpsMetadata(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

// ListOpsMetadataAPIClient is a client that implements the ListOpsMetadata
// operation.
type ListOpsMetadataAPIClient interface {
	ListOpsMetadata(context.Context, *ListOpsMetadataInput, ...func(*Options)) (*ListOpsMetadataOutput, error)
}

var _ ListOpsMetadataAPIClient = (*Client)(nil)

// ListOpsMetadataPaginatorOptions is the paginator options for ListOpsMetadata
type ListOpsMetadataPaginatorOptions struct {
	// The maximum number of items to return for this call. The call also returns a
	// token that you can specify in a subsequent call to get the next set of results.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListOpsMetadataPaginator is a paginator for ListOpsMetadata
type ListOpsMetadataPaginator struct {
	options   ListOpsMetadataPaginatorOptions
	client    ListOpsMetadataAPIClient
	params    *ListOpsMetadataInput
	nextToken *string
	firstPage bool
}

// NewListOpsMetadataPaginator returns a new ListOpsMetadataPaginator
func NewListOpsMetadataPaginator(client ListOpsMetadataAPIClient, params *ListOpsMetadataInput, optFns ...func(*ListOpsMetadataPaginatorOptions)) *ListOpsMetadataPaginator {
	if params == nil {
		params = &ListOpsMetadataInput{}
	}

	options := ListOpsMetadataPaginatorOptions{}
	if params.MaxResults != 0 {
		options.Limit = params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListOpsMetadataPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListOpsMetadataPaginator) HasMorePages() bool {
	return p.firstPage || p.nextToken != nil
}

// NextPage retrieves the next ListOpsMetadata page.
func (p *ListOpsMetadataPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListOpsMetadataOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	params.MaxResults = p.options.Limit

	result, err := p.client.ListOpsMetadata(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken && prevToken != nil && p.nextToken != nil && *prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opListOpsMetadata(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ssm",
		OperationName: "ListOpsMetadata",
	}
}
