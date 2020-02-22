package main

import (
	"basicAPI/proto"
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//this is our server code so it needs the server struct
//this is so server can implement the interface generated in the service.pb.go
type server struct{}

func main() {

	//open a tcp listener at a specific port. If theres an error the err will get assigned something and we can fire the panic function
	//otherwise carry on.
	//this listerner will listen for clients requests on port 4040
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	//create a server using a function provided in the grpc package from google
	OurServer := grpc.NewServer()

	//now assign this server in our proto package using a function generated in the service.pb.go
	//we pass it the server we created and the location of our server struct.
	proto.RegisterAddServiceServer(OurServer, &server{})

	//because we'll be serializing and deserializing data we import the reflection package from google and register our server
	reflection.Register(OurServer)

	//now we tell our server to run.
	//this is a fancy if statement that runs the server but if the serve call returns an error it will assign it
	//to err2 and if err2 isn't empty (caught an error) we fire the panic function.
	if err2 := OurServer.Serve(listener); err2 != nil {
		panic(err2)
	}
}

//to implement the interface that we generated in service.pb.go we need this server package to have the two functions the interface does: add, and multiply
//these functions need to have a reciever of server (unimplementedAddServiceServer), which we are implementing now. we use a pointer *server to the server.
//if we look at what add takes in in the service.pb.go:
//-------------------------func (*UnimplementedAddServiceServer) Add(ctx context.Context, req *Request) (*Response, error) {
//------------------------------return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
//-------------------------}
//we see it takes in a context, and a request, and returns a response and error
// the reason our function's pointers have proto. infront of the request and response is because that is where our service.pb.go is stored, in a package called proto.
//if we had named the package service we would use *service.Request
func (s *server) Add(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	//the request was generated get functions so adding is fairly simple

	a := request.GetA()
	b := request.GetB()

	result := a + b

	//now we return a pointer to the proto package's response structure were we assign our result to it's result and we return the error as 'nil' (no error)
	return &proto.Response{Result: result}, nil

}

//now we do the same with multiply.

func (s *server) Multiply(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	//user the getters from the proto package we imported again.
	a := request.GetA()
	b := request.GetB()

	result := a * b

	return &proto.Response{Result: result}, nil

}
