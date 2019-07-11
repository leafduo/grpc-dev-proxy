package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func invoke(w http.ResponseWriter, target string, service string, method string, headers []string, message string) error {
	ctx := context.Background()
	network := "tcp"
	var creds credentials.TransportCredentials
	var opts []grpc.DialOption
	cc, err := grpcurl.BlockingDial(ctx, network, target, creds, opts...)
	if err != nil {
		err := xerrors.Errorf("Failed to dial target host %q: %w", target, err)
		reportError(w, err, 500)
		return err
	}
	refClient := grpcreflect.NewClient(ctx, reflectpb.NewServerReflectionClient(cc))
	descSource := grpcurl.DescriptorSourceFromServer(ctx, refClient)

	format := "json"
	in := strings.NewReader(message)
	output := &bytes.Buffer{}
	rf, formatter, err := grpcurl.RequestParserAndFormatterFor(grpcurl.Format(format), descSource, true, true, in)
	if err != nil {
		err := xerrors.Errorf("failed to construct request parser and formatter for %s: %w", format, err)
		reportError(w, err, 500)
		return err
	}
	h := grpcurl.NewDefaultEventHandler(output, descSource, formatter, true)

	symbol := fmt.Sprintf("%s/%s", service, method)
	err = grpcurl.InvokeRPC(ctx, descSource, cc, symbol, headers, h, rf.Next)
	if err != nil {
		err := xerrors.Errorf("Error invoking method %q: %w", symbol, err)
		reportError(w, err, 500)
		return err
	}

	if h.Status.Code() != codes.OK {
		errorMessage := fmt.Sprintf("\nERROR:\n  Code: %s\n  Message: %s\n", h.Status.Code().String(), h.Status.Message())
		_, _ = io.WriteString(w, errorMessage)
		if len(h.Status.Details()) > 0 {
			_, _ = io.WriteString(w, fmt.Sprintf("  Details: %s\n", h.Status.Details()))
		}
	}

	_, _ = w.Write(output.Bytes())

	return nil
}
