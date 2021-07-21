package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func RPCInvokeHandler(services map[string]*Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Request must be JSON", http.StatusUnsupportedMediaType)
			return
		}

		input, err := getInput(r.Body)
		if err != nil {
			http.Error(w, "Failed to parse input data", http.StatusBadRequest)
			return
		}

		svc, ok := services[input.Service]
		if ok {
			for _, md := range svc.methods {
				method := input.Service + "." + input.Method
				if md.GetFullyQualifiedName() == method {
					descSource, err := grpcurl.DescriptorSourceFromFileDescriptors(md.GetFile())
					if err != nil {
						http.Error(w, "Failed to create descriptor source: "+err.Error(), http.StatusInternalServerError)
						return
					}
					results, err := invokeRPC(r.Context(), method, svc.cc, descSource, input)
					if err != nil {
						http.Error(w, "Unexpected error: "+err.Error(), http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					enc := json.NewEncoder(w)
					enc.SetIndent("", "  ")
					enc.Encode(results)
					return
				}
			}
		}

		http.NotFound(w, r)
	})
}

func getInput(body io.Reader) (*rpcInput, error) {
	var input rpcInput

	js, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(js, &input); err != nil {
		return nil, err
	}

	return &input, nil
}

func invokeRPC(ctx context.Context, methodName string, cc grpc.ClientConnInterface, descSource grpcurl.DescriptorSource, input *rpcInput) (*rpcResult, error) {
	reqStats := rpcRequestStats{
		Total: len(input.Data),
	}
	requestFunc := func(m proto.Message) error {
		if len(input.Data) == 0 {
			return io.EOF
		}
		reqStats.Sent++
		req := input.Data[0]
		input.Data = input.Data[1:]
		if err := jsonpb.Unmarshal(bytes.NewReader([]byte(req)), m); err != nil {
			return status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil
	}

	hdrs := make([]string, len(input.Metadata))
	for i, hdr := range input.Metadata {
		hdrs[i] = fmt.Sprintf("%s: %s", hdr.Name, hdr.Value)
	}

	if input.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		timeout := time.Duration(input.TimeoutSeconds * float32(time.Second))
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	result := rpcResult{
		descSource: descSource,
		Requests:   &reqStats,
	}
	if err := grpcurl.InvokeRPC(ctx, descSource, cc, methodName, hdrs, &result, requestFunc); err != nil {
		return nil, err
	}

	return &result, nil
}

type rpcInput struct {
	TimeoutSeconds float32           `json:"timeout_seconds"`
	Metadata       []rpcMetadata     `json:"metadata"`
	Data           []json.RawMessage `json:"data"`
	Service        string            `json:"service"`
	Method         string            `json:"method"`
}

type rpcMetadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type rpcResponseElement struct {
	Data    json.RawMessage `json:"message"`
	IsError bool            `json:"isError"`
}

type rpcRequestStats struct {
	Total int `json:"total"`
	Sent  int `json:"sent"`
}

type rpcError struct {
	Code    uint32               `json:"code"`
	Name    string               `json:"name"`
	Message string               `json:"message"`
	Details []rpcResponseElement `json:"details"`
}

type rpcResult struct {
	descSource grpcurl.DescriptorSource
	Headers    []rpcMetadata        `json:"headers"`
	Error      *rpcError            `json:"error"`
	Responses  []rpcResponseElement `json:"responses"`
	Requests   *rpcRequestStats     `json:"requests"`
	Trailers   []rpcMetadata        `json:"trailers"`
}

func (*rpcResult) OnResolveMethod(*desc.MethodDescriptor) {}

func (*rpcResult) OnSendHeaders(metadata.MD) {}

func (r *rpcResult) OnReceiveHeaders(md metadata.MD) {
	r.Headers = responseMetadata(md)
}

func (r *rpcResult) OnReceiveResponse(m proto.Message) {
	r.Responses = append(r.Responses, responseToJSON(r.descSource, m))
}

func (r *rpcResult) OnReceiveTrailers(stat *status.Status, md metadata.MD) {
	r.Trailers = responseMetadata(md)
	r.Error = toRpcError(r.descSource, stat)
}

func responseMetadata(md metadata.MD) []rpcMetadata {
	keys := make([]string, 0, len(md))
	for k := range md {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ret := make([]rpcMetadata, 0, len(md))
	for _, k := range keys {
		vals := md[k]
		for _, v := range vals {
			if strings.HasSuffix(k, "-bin") {
				v = base64.StdEncoding.EncodeToString([]byte(v))
			}
			ret = append(ret, rpcMetadata{Name: k, Value: v})
		}
	}
	return ret
}

func toRpcError(descSource grpcurl.DescriptorSource, stat *status.Status) *rpcError {
	if stat.Code() == codes.OK {
		return nil
	}

	details := stat.Proto().Details
	msgs := make([]rpcResponseElement, len(details))
	for i, d := range details {
		msgs[i] = responseToJSON(descSource, d)
	}

	return &rpcError{
		Code:    uint32(stat.Code()),
		Name:    stat.Code().String(),
		Message: stat.Message(),
		Details: msgs,
	}
}

func responseToJSON(descSource grpcurl.DescriptorSource, msg proto.Message) rpcResponseElement {
	anyResolver := grpcurl.AnyResolverFromDescriptorSourceWithFallback(descSource)
	jsm := jsonpb.Marshaler{EmitDefaults: true, OrigName: true, Indent: "  ", AnyResolver: anyResolver}
	var b bytes.Buffer
	if err := jsm.Marshal(&b, msg); err == nil {
		return rpcResponseElement{Data: json.RawMessage(b.Bytes())}
	} else {
		b, _ := json.Marshal(err.Error())
		return rpcResponseElement{Data: json.RawMessage(b), IsError: true}
	}
}
