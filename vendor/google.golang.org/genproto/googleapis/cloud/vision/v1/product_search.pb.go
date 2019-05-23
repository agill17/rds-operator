// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/vision/v1/product_search.proto

package vision

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Parameters for a product search request.
type ProductSearchParams struct {
	// The bounding polygon around the area of interest in the image.
	// Optional. If it is not specified, system discretion will be applied.
	BoundingPoly *BoundingPoly `protobuf:"bytes,9,opt,name=bounding_poly,json=boundingPoly,proto3" json:"bounding_poly,omitempty"`
	// The resource name of a [ProductSet][google.cloud.vision.v1.ProductSet] to
	// be searched for similar images.
	//
	// Format is:
	// `projects/PROJECT_ID/locations/LOC_ID/productSets/PRODUCT_SET_ID`.
	ProductSet string `protobuf:"bytes,6,opt,name=product_set,json=productSet,proto3" json:"product_set,omitempty"`
	// The list of product categories to search in. Currently, we only consider
	// the first category, and either "homegoods-v2", "apparel-v2", or "toys-v2"
	// should be specified. The legacy categories "homegoods", "apparel", and
	// "toys" are still supported, but these should not be used for new products.
	ProductCategories []string `protobuf:"bytes,7,rep,name=product_categories,json=productCategories,proto3" json:"product_categories,omitempty"`
	// The filtering expression. This can be used to restrict search results based
	// on Product labels. We currently support an AND of OR of key-value
	// expressions, where each expression within an OR must have the same key. An
	// '=' should be used to connect the key and value.
	//
	// For example, "(color = red OR color = blue) AND brand = Google" is
	// acceptable, but "(color = red OR brand = Google)" is not acceptable.
	// "color: red" is not acceptable because it uses a ':' instead of an '='.
	Filter               string   `protobuf:"bytes,8,opt,name=filter,proto3" json:"filter,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProductSearchParams) Reset()         { *m = ProductSearchParams{} }
func (m *ProductSearchParams) String() string { return proto.CompactTextString(m) }
func (*ProductSearchParams) ProtoMessage()    {}
func (*ProductSearchParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_4fdf2c80d2106c63, []int{0}
}

func (m *ProductSearchParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProductSearchParams.Unmarshal(m, b)
}
func (m *ProductSearchParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProductSearchParams.Marshal(b, m, deterministic)
}
func (m *ProductSearchParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProductSearchParams.Merge(m, src)
}
func (m *ProductSearchParams) XXX_Size() int {
	return xxx_messageInfo_ProductSearchParams.Size(m)
}
func (m *ProductSearchParams) XXX_DiscardUnknown() {
	xxx_messageInfo_ProductSearchParams.DiscardUnknown(m)
}

var xxx_messageInfo_ProductSearchParams proto.InternalMessageInfo

func (m *ProductSearchParams) GetBoundingPoly() *BoundingPoly {
	if m != nil {
		return m.BoundingPoly
	}
	return nil
}

func (m *ProductSearchParams) GetProductSet() string {
	if m != nil {
		return m.ProductSet
	}
	return ""
}

func (m *ProductSearchParams) GetProductCategories() []string {
	if m != nil {
		return m.ProductCategories
	}
	return nil
}

func (m *ProductSearchParams) GetFilter() string {
	if m != nil {
		return m.Filter
	}
	return ""
}

// Results for a product search request.
type ProductSearchResults struct {
	// Timestamp of the index which provided these results. Products added to the
	// product set and products removed from the product set after this time are
	// not reflected in the current results.
	IndexTime *timestamp.Timestamp `protobuf:"bytes,2,opt,name=index_time,json=indexTime,proto3" json:"index_time,omitempty"`
	// List of results, one for each product match.
	Results []*ProductSearchResults_Result `protobuf:"bytes,5,rep,name=results,proto3" json:"results,omitempty"`
	// List of results grouped by products detected in the query image. Each entry
	// corresponds to one bounding polygon in the query image, and contains the
	// matching products specific to that region. There may be duplicate product
	// matches in the union of all the per-product results.
	ProductGroupedResults []*ProductSearchResults_GroupedResult `protobuf:"bytes,6,rep,name=product_grouped_results,json=productGroupedResults,proto3" json:"product_grouped_results,omitempty"`
	XXX_NoUnkeyedLiteral  struct{}                              `json:"-"`
	XXX_unrecognized      []byte                                `json:"-"`
	XXX_sizecache         int32                                 `json:"-"`
}

func (m *ProductSearchResults) Reset()         { *m = ProductSearchResults{} }
func (m *ProductSearchResults) String() string { return proto.CompactTextString(m) }
func (*ProductSearchResults) ProtoMessage()    {}
func (*ProductSearchResults) Descriptor() ([]byte, []int) {
	return fileDescriptor_4fdf2c80d2106c63, []int{1}
}

func (m *ProductSearchResults) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProductSearchResults.Unmarshal(m, b)
}
func (m *ProductSearchResults) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProductSearchResults.Marshal(b, m, deterministic)
}
func (m *ProductSearchResults) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProductSearchResults.Merge(m, src)
}
func (m *ProductSearchResults) XXX_Size() int {
	return xxx_messageInfo_ProductSearchResults.Size(m)
}
func (m *ProductSearchResults) XXX_DiscardUnknown() {
	xxx_messageInfo_ProductSearchResults.DiscardUnknown(m)
}

var xxx_messageInfo_ProductSearchResults proto.InternalMessageInfo

func (m *ProductSearchResults) GetIndexTime() *timestamp.Timestamp {
	if m != nil {
		return m.IndexTime
	}
	return nil
}

func (m *ProductSearchResults) GetResults() []*ProductSearchResults_Result {
	if m != nil {
		return m.Results
	}
	return nil
}

func (m *ProductSearchResults) GetProductGroupedResults() []*ProductSearchResults_GroupedResult {
	if m != nil {
		return m.ProductGroupedResults
	}
	return nil
}

// Information about a product.
type ProductSearchResults_Result struct {
	// The Product.
	Product *Product `protobuf:"bytes,1,opt,name=product,proto3" json:"product,omitempty"`
	// A confidence level on the match, ranging from 0 (no confidence) to
	// 1 (full confidence).
	Score float32 `protobuf:"fixed32,2,opt,name=score,proto3" json:"score,omitempty"`
	// The resource name of the image from the product that is the closest match
	// to the query.
	Image                string   `protobuf:"bytes,3,opt,name=image,proto3" json:"image,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProductSearchResults_Result) Reset()         { *m = ProductSearchResults_Result{} }
func (m *ProductSearchResults_Result) String() string { return proto.CompactTextString(m) }
func (*ProductSearchResults_Result) ProtoMessage()    {}
func (*ProductSearchResults_Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_4fdf2c80d2106c63, []int{1, 0}
}

func (m *ProductSearchResults_Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProductSearchResults_Result.Unmarshal(m, b)
}
func (m *ProductSearchResults_Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProductSearchResults_Result.Marshal(b, m, deterministic)
}
func (m *ProductSearchResults_Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProductSearchResults_Result.Merge(m, src)
}
func (m *ProductSearchResults_Result) XXX_Size() int {
	return xxx_messageInfo_ProductSearchResults_Result.Size(m)
}
func (m *ProductSearchResults_Result) XXX_DiscardUnknown() {
	xxx_messageInfo_ProductSearchResults_Result.DiscardUnknown(m)
}

var xxx_messageInfo_ProductSearchResults_Result proto.InternalMessageInfo

func (m *ProductSearchResults_Result) GetProduct() *Product {
	if m != nil {
		return m.Product
	}
	return nil
}

func (m *ProductSearchResults_Result) GetScore() float32 {
	if m != nil {
		return m.Score
	}
	return 0
}

func (m *ProductSearchResults_Result) GetImage() string {
	if m != nil {
		return m.Image
	}
	return ""
}

// Information about the products similar to a single product in a query
// image.
type ProductSearchResults_GroupedResult struct {
	// The bounding polygon around the product detected in the query image.
	BoundingPoly *BoundingPoly `protobuf:"bytes,1,opt,name=bounding_poly,json=boundingPoly,proto3" json:"bounding_poly,omitempty"`
	// List of results, one for each product match.
	Results              []*ProductSearchResults_Result `protobuf:"bytes,2,rep,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *ProductSearchResults_GroupedResult) Reset()         { *m = ProductSearchResults_GroupedResult{} }
func (m *ProductSearchResults_GroupedResult) String() string { return proto.CompactTextString(m) }
func (*ProductSearchResults_GroupedResult) ProtoMessage()    {}
func (*ProductSearchResults_GroupedResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_4fdf2c80d2106c63, []int{1, 1}
}

func (m *ProductSearchResults_GroupedResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProductSearchResults_GroupedResult.Unmarshal(m, b)
}
func (m *ProductSearchResults_GroupedResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProductSearchResults_GroupedResult.Marshal(b, m, deterministic)
}
func (m *ProductSearchResults_GroupedResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProductSearchResults_GroupedResult.Merge(m, src)
}
func (m *ProductSearchResults_GroupedResult) XXX_Size() int {
	return xxx_messageInfo_ProductSearchResults_GroupedResult.Size(m)
}
func (m *ProductSearchResults_GroupedResult) XXX_DiscardUnknown() {
	xxx_messageInfo_ProductSearchResults_GroupedResult.DiscardUnknown(m)
}

var xxx_messageInfo_ProductSearchResults_GroupedResult proto.InternalMessageInfo

func (m *ProductSearchResults_GroupedResult) GetBoundingPoly() *BoundingPoly {
	if m != nil {
		return m.BoundingPoly
	}
	return nil
}

func (m *ProductSearchResults_GroupedResult) GetResults() []*ProductSearchResults_Result {
	if m != nil {
		return m.Results
	}
	return nil
}

func init() {
	proto.RegisterType((*ProductSearchParams)(nil), "google.cloud.vision.v1.ProductSearchParams")
	proto.RegisterType((*ProductSearchResults)(nil), "google.cloud.vision.v1.ProductSearchResults")
	proto.RegisterType((*ProductSearchResults_Result)(nil), "google.cloud.vision.v1.ProductSearchResults.Result")
	proto.RegisterType((*ProductSearchResults_GroupedResult)(nil), "google.cloud.vision.v1.ProductSearchResults.GroupedResult")
}

func init() {
	proto.RegisterFile("google/cloud/vision/v1/product_search.proto", fileDescriptor_4fdf2c80d2106c63)
}

var fileDescriptor_4fdf2c80d2106c63 = []byte{
	// 485 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xcf, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0x95, 0x76, 0x4b, 0xa9, 0xcb, 0x0e, 0x98, 0x31, 0xa2, 0x08, 0xa9, 0xd5, 0x04, 0x52,
	0x25, 0x84, 0xa3, 0xad, 0xa7, 0x01, 0xa7, 0xee, 0x30, 0x71, 0x00, 0x55, 0x01, 0x71, 0xe0, 0x12,
	0xb9, 0x89, 0x67, 0x2c, 0x25, 0x7e, 0x91, 0xed, 0x54, 0x94, 0x3f, 0x87, 0x1b, 0x47, 0xfe, 0x0b,
	0xfe, 0x24, 0x8e, 0xa8, 0xfe, 0x01, 0x2b, 0xac, 0xe2, 0xc7, 0x4e, 0xc9, 0xb3, 0xbf, 0xef, 0xf3,
	0xfc, 0x7d, 0x7e, 0x46, 0x8f, 0x39, 0x00, 0xaf, 0x59, 0x56, 0xd6, 0xd0, 0x55, 0xd9, 0x4a, 0x68,
	0x01, 0x32, 0x5b, 0x9d, 0x64, 0xad, 0x82, 0xaa, 0x2b, 0x4d, 0xa1, 0x19, 0x55, 0xe5, 0x7b, 0xd2,
	0x2a, 0x30, 0x80, 0x8f, 0x9c, 0x98, 0x58, 0x31, 0x71, 0x62, 0xb2, 0x3a, 0x49, 0x1f, 0x78, 0x08,
	0x6d, 0x45, 0x46, 0xa5, 0x04, 0x43, 0x8d, 0x00, 0xa9, 0x5d, 0x56, 0xfa, 0x68, 0x47, 0x09, 0xce,
	0xa0, 0x61, 0x46, 0xad, 0xbd, 0x6c, 0xf6, 0x57, 0x27, 0x29, 0x34, 0x53, 0x2b, 0x51, 0x32, 0x9f,
	0x34, 0xf6, 0x49, 0x36, 0x5a, 0x76, 0x97, 0x99, 0x11, 0x0d, 0xd3, 0x86, 0x36, 0xad, 0x13, 0x1c,
	0x7f, 0x8d, 0xd0, 0xdd, 0x85, 0x23, 0xbc, 0xb6, 0x80, 0x05, 0x55, 0xb4, 0xd1, 0xf8, 0x05, 0x3a,
	0x58, 0x42, 0x27, 0x2b, 0x21, 0x79, 0xd1, 0x42, 0xbd, 0x4e, 0x86, 0x93, 0x68, 0x3a, 0x3a, 0x7d,
	0x48, 0xae, 0xb7, 0x48, 0xe6, 0x5e, 0xbc, 0x80, 0x7a, 0x9d, 0xdf, 0x5e, 0x5e, 0x89, 0xf0, 0x18,
	0x8d, 0x7e, 0x9e, 0xd1, 0x24, 0xf1, 0x24, 0x9a, 0x0e, 0x73, 0xd4, 0x86, 0xa2, 0x06, 0x3f, 0x41,
	0x38, 0x08, 0x4a, 0x6a, 0x18, 0x07, 0x25, 0x98, 0x4e, 0x06, 0x93, 0xfe, 0x74, 0x98, 0xdf, 0xf1,
	0x3b, 0xe7, 0x3f, 0x36, 0xf0, 0x11, 0x8a, 0x2f, 0x45, 0x6d, 0x98, 0x4a, 0x6e, 0x59, 0x94, 0x8f,
	0x8e, 0xbf, 0xec, 0xa1, 0xc3, 0x2d, 0x2b, 0x39, 0xd3, 0x5d, 0x6d, 0x34, 0x3e, 0x43, 0x48, 0xc8,
	0x8a, 0x7d, 0x28, 0x36, 0xe6, 0x93, 0x9e, 0x35, 0x92, 0x06, 0x23, 0xa1, 0x33, 0xe4, 0x4d, 0xe8,
	0x4c, 0x3e, 0xb4, 0xea, 0x4d, 0x8c, 0x5f, 0xa2, 0x81, 0x72, 0x94, 0x64, 0x7f, 0xd2, 0x9f, 0x8e,
	0x4e, 0x67, 0xbb, 0x1a, 0x70, 0x5d, 0x65, 0xe2, 0xbe, 0x79, 0x60, 0x60, 0x85, 0xee, 0x07, 0xa7,
	0x5c, 0x41, 0xd7, 0xb2, 0xaa, 0x08, 0xf8, 0xd8, 0xe2, 0x9f, 0xfe, 0x13, 0xfe, 0xc2, 0x31, 0x7c,
	0x95, 0x7b, 0x1e, 0xbd, 0xb5, 0xaa, 0x53, 0x40, 0xb1, 0xfb, 0xc5, 0x67, 0x68, 0xe0, 0x25, 0x49,
	0x64, 0x9b, 0x30, 0xfe, 0x43, 0xb5, 0x3c, 0xe8, 0xf1, 0x21, 0xda, 0xd7, 0x25, 0x28, 0xd7, 0xbd,
	0x5e, 0xee, 0x82, 0xcd, 0xaa, 0x68, 0x28, 0x67, 0x49, 0xdf, 0x5e, 0x84, 0x0b, 0xd2, 0xcf, 0x11,
	0x3a, 0xd8, 0x3a, 0xc3, 0xef, 0xc3, 0x14, 0xfd, 0xf7, 0x30, 0x5d, 0xb9, 0x90, 0xde, 0xcd, 0x2f,
	0x64, 0xfe, 0x11, 0xa5, 0x25, 0x34, 0x3b, 0x10, 0x73, 0xbc, 0xfd, 0x32, 0x36, 0x93, 0xb2, 0x88,
	0xde, 0x3d, 0xf7, 0x6a, 0x0e, 0x35, 0x95, 0x9c, 0x80, 0xe2, 0x19, 0x67, 0xd2, 0xce, 0x51, 0xe6,
	0xb6, 0x68, 0x2b, 0xf4, 0xaf, 0xef, 0xf4, 0x99, 0xfb, 0xfb, 0x16, 0x45, 0x9f, 0x7a, 0x7b, 0x17,
	0xe7, 0x6f, 0x5f, 0x2d, 0x63, 0x9b, 0x32, 0xfb, 0x1e, 0x00, 0x00, 0xff, 0xff, 0x93, 0xfe, 0xaa,
	0xbb, 0x63, 0x04, 0x00, 0x00,
}
