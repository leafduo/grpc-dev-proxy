package client

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func Invoke(target string, service string, method string, headers []string, message string) (string, error) {
	ctx := context.Background()
	network := "tcp"
	var creds credentials.TransportCredentials
	var opts []grpc.DialOption
	cc, err := grpcurl.BlockingDial(ctx, network, target, creds, opts...)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to dial target host %q", target)
	}
	refClient := grpcreflect.NewClient(ctx, reflectpb.NewServerReflectionClient(cc))
	descSource := grpcurl.DescriptorSourceFromServer(ctx, refClient)

	format := "json"
	in := strings.NewReader(message)
	includeSeparators := true
	emitDefaults := true
	output := &bytes.Buffer{}
	rf, formatter, err := grpcurl.RequestParserAndFormatterFor(grpcurl.Format(format), descSource, emitDefaults, includeSeparators, in)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to construct request parser and formatter for %s", format)
	}
	h := grpcurl.NewDefaultEventHandler(output, descSource, formatter, true)

	symbol := fmt.Sprintf("%s/%s", service, method)
	err = grpcurl.InvokeRPC(ctx, descSource, cc, symbol, headers, h, rf.Next)
	if err != nil {
		logrus.WithError(err).Errorf("Error invoking method %q", symbol)
	}

	if h.Status.Code() != codes.OK {
		errorMessage := fmt.Sprintf("\nERROR:\n  Code: %s\n  Message: %s\n", h.Status.Code().String(), h.Status.Message())
		output.WriteString(errorMessage)
		if len(h.Status.Details()) > 0 {
			output.WriteString(fmt.Sprintf("  Details: %s\n", h.Status.Details()))
		}
	}

	return output.String(), nil
}
