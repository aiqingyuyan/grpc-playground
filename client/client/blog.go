package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
	"yanyu/go-grpc/generated/blog"
)

func initialize() (*grpc.ClientConn, blog.BlogServiceClient) {
	var options []grpc.DialOption

	options = append(options, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:50051", options...)

	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	blogClient := blog.NewBlogServiceClient(conn)

	return conn, blogClient
}

func BlogUnaryExample() {
	conn, client := initialize()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	req := &blog.GetBlogRequest{
		Id: 1,
	}

	resp, err := client.GetBlog(ctx, req)

	if err != nil {
		log.Fatalf("%v.GetBlog(_) = _, %v", client, err)
	}

	log.Println(resp)
}

func BlogClientStreamExample() {
	conn, client := initialize()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	log.Println("Sending blogs...")

	stream, err := client.SaveBlogs(ctx)
	if err != nil {
		log.Fatalf("%v.SaveBlogs(_) = _, %v", client, err)
	}

	for i := 0; i < 10; i++ {
		log.Printf("Sending blog %d", i)

		err := stream.Send(&blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog %d", i),
				Author: "Leeroy",
				Text: fmt.Sprintf("Leeroy Jekins %d", i),
			},
		})
		if err != nil {
			log.Printf("%v.Send(_) = _, %v", stream, err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v", stream, err)
	}

	log.Printf("Total number of blogs in server %d", resp.NumberOfBlogs)
}

func BlogServerStreamExample() {
	conn, client := initialize()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	log.Println("Listing blogs...")

	stream, err := client.ListBlogs(ctx, &blog.ListBlogsRequest{})
	if err != nil {
		log.Fatalf("%v.ListBlogs(_) = _, %v", client, err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("%v.Recv(_) = _, %v", stream, err)
		}

		log.Printf("blog: %v", resp.Blog)
	}
}

func BlogBiStreamExample() {
	conn, client := initialize()
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	requests := []blog.SaveBlogRequest{
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog leeroy"),
				Author: "Leeroy",
				Text: fmt.Sprintf("Leeroy Jekins !!!!!!"),
			},
		},
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog jane"),
				Author: "Jane",
				Text: fmt.Sprintf("Jane Doe test"),
			},
		},
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog john"),
				Author: "John",
				Text: fmt.Sprintf("John Doe test"),
			},
		},
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog leeroy 1"),
				Author: "Leeroy",
				Text: fmt.Sprintf("Leeeeeeeeeeroy Jekins !!!!!!"),
			},
		},
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog leeroy 2"),
				Author: "Leeroy",
				Text: fmt.Sprintf("Leeeeeeeroy Jekins Test !!!!!!"),
			},
		},
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog alice"),
				Author: "Alice",
				Text: fmt.Sprintf("Alice blog test !!!!!!"),
			},
		},
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog leeroy 3"),
				Author: "Leeroy",
				Text: fmt.Sprintf("what's up, leeroy!!!"),
			},
		},
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog foo"),
				Author: "Foo",
				Text: fmt.Sprintf("foo blog test !!!!!!"),
			},
		},
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog bar"),
				Author: "bar",
				Text: fmt.Sprintf("blog test for bar !!!!!!"),
			},
		},
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog leeroy 4"),
				Author: "Leeroy",
				Text: fmt.Sprintf("Leeeeeeeeeeroy Jekinsssssss !!!!!!"),
			},
		},
		blog.SaveBlogRequest{
			Blog: &blog.Blog{
				Title: fmt.Sprintf("test blog author 1"),
				Author: "Author 1",
				Text: fmt.Sprintf("test blog for author 1 !"),
			},
		},
	}

	log.Println("Saving blogs...")

	stream, err := client.GetAuthorWithMostBlogsOnSave(ctx)
	if err != nil {
		log.Fatalf("%v.GetAuthorWithMostBlogsOnSave(_) = _, %v", client, err)
	}

	done := make(chan struct{})
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Printf("%v.Recv(_) = _, %v", stream, err)
			}
			log.Printf("Author with most blogs: %v", resp.Author)
		}
	}()

	for _, req := range requests {
		log.Printf("Saving blog %v", req.Blog.Title)
		if err := stream.Send(&req); err != nil {
			log.Printf("%v.Send(_) = _, %v", stream, err)
		}
		time.Sleep(300 * time.Millisecond)
	}

	stream.CloseSend()

	<-done
}
