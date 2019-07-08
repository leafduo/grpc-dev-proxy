package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func listServices(w http.ResponseWriter, target string) error {
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
	services, err := descSource.ListServices()
	if err != nil {
		reportError(w, err, 500)
	}

	_, _ = io.WriteString(w, "Supported services on target host:\n\n")
	for _, serviceName := range services {
		_, _ = io.WriteString(w, serviceName)
		_, _ = io.WriteString(w, "\n")
	}

	_, _ = io.WriteString(w, "\n")
	_, _ = io.WriteString(w, "Put service name in HTTP header `Service` to specify which service to invoke")

	return nil
}

func listMethods(w http.ResponseWriter, target string, service string) error {
	return nil
}
