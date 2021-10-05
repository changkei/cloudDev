// Package main imlements a client for movieinfo service
package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloudnativecourse/lab7-docker/lab5-grpc/movieapi"
	"google.golang.org/grpc"
)

const (
	address      = "localhost:50051"
	defaultTitle = "Pulp fiction"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Contact the server and print out its response.
	title := defaultTitle
	if len(os.Args) > 1 {
		title = os.Args[1]
	}
	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s %d %s %v", title, r.GetYear(), r.GetDirector(), r.GetCast())

	newTitle := "Cloud Native Architecture Episode 5: gPRC Strikes Back"
	cast := make([]string, 0)
	cast = append(cast, "Keith Chang", "Landon Gibson", "James Tallett")
	director := "Arun Ravindran"
	var year int32 = 2021

	newr, newerr := c.SetMovieInfo(ctx, &movieapi.MovieData{Title: newTitle, Year: year, Director: director, Cast: cast})
	if newerr != nil {
		log.Fatalf("could not set movie info: %v", newerr)
	}
	log.Printf("Status: %s", newr.GetDebug())
	out, outerr := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: newTitle})
	if outerr != nil {
		log.Fatalf("could not get movie info: %v", outerr)
	}
	log.Printf("Movie Info for %s %d %s %v", newTitle, out.GetYear(), out.GetDirector(), out.GetCast())

}
