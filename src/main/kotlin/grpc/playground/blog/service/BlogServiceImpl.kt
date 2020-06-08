package grpc.playground.blog.service

import com.example.protobuf.blog.Blog
import com.example.protobuf.blog.BlogServiceGrpcKt
import com.example.protobuf.blog.GetBlogRequest
import com.example.protobuf.blog.GetBlogResponse

class BlogServiceImpl : BlogServiceGrpcKt.BlogServiceCoroutineImplBase() {

	private var blogCollection: List<Blog>

	init {
		blogCollection = generateBlogsData()
	}

	private fun generateBlogsData(): List<Blog> {
		val blogList = mutableListOf<Blog>()

		for (i in 1..5) {
			blogList.add(
				Blog
					.newBuilder()
					.setId(i)
					.setTitle("Blog $i")
					.setAuthor("Author $i")
					.setText("This is a piece of test text $i")
					.build()
			)
		}

		return blogList
	}

	override suspend fun getBlog(request: GetBlogRequest): GetBlogResponse =
		request.id
			.let { blogId ->
				blogCollection.find { blog ->
					blog.id == blogId
				}
			}
			.let { blog ->
				GetBlogResponse.newBuilder()
					.setBlog(blog)
					.build()
			}
}
