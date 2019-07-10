package main

import (
	"context"
	"fmt"
	"log"

	"github.com/maxbrain0/grpc-go-course/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create Blog
	fmt.Println("Creating a blog")
	blog := &blogpb.Blog{
		AuthorId: "Juancho",
		Title:    "My First Post",
		Content:  "This is my veeerrry special first post!",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})

	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	fmt.Printf("Blog has been created: %v\n", createBlogRes)

	// extact inserted id for next step
	blogID := createBlogRes.GetBlog().GetId()

	// read Blog
	fmt.Println("Reading the blog (post)")

	// bad read
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5d2400321bd494c65c023aa5"})

	if err2 != nil {
		fmt.Printf("Error occured while reading (finding) blog (post): %v \n", err2)
	}

	// good read
	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)

	if readBlogErr != nil {
		fmt.Printf("Error occured while reading (finding) blog (post): %v \n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v \n", readBlogRes.GetBlog())

	// update blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Eddy",
		Title:    "Been posting for days",
		Content:  "Yall's a buncha suckas!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})

	if updateErr != nil {
		fmt.Printf("Error occured while updating: %v\n", updateErr)
	}

	fmt.Printf("Blog was read: %v\n", updateRes)

	// delete blog
	delRes, delErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if delErr != nil {
		fmt.Printf("Error occured while deleting: %v\n", delErr)
	}

	fmt.Printf("Blog was deleted: %v\n", delRes)
}
